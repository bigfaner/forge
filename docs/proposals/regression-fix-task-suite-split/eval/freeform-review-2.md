---
created: "2026-05-28"
reviewer: domain-expert/agent-cognitive-load-dx
document: proposal.md
---

# Freeform Review: Agent Cognitive Load Perspective

## Section 1: Background Assessment

这份提案要解决的核心问题不是技术性的——它是关于 LLM agent 的注意力管理。当一个 fix task 包含 4 个 test suite 的 20+ 失败时，agent 不是"处理不了这么多信息"，而是在高噪声环境中无法建立有效的因果推理链。当前的 single-task 模式把所有失败信号塞进一个 description，agent 的注意力被分散到多个不相关的文件和模块中，导致它无法确定"先修哪个、怎么修、修完怎么验证"。

提案的核心策略是按测试文件拆分 fix task。这个方向在认知负载层面是合理的——每个 task 聚焦到单一文件，减少了 agent 需要同时处理的注意力焦点。但"聚焦到单一文件"不等于"agent 能高效修复"。从 agent 执行者视角看，一个 fix task 的价值取决于三件事：(1) agent 能否快速定位 root cause（因果链完整性），(2) description 中的信息是否足以引导修复（信息充分性），(3) agent 在修复过程中是否会遇到不可见的障碍（并发安全性）。

提案在这些维度上的表现参差不齐。分组策略本身（按测试文件）是合理的注意力分割，但每个 task 传递给 agent 的信息质量和结构，以及多 task 并发时的协调机制，都存在值得深入审视的问题。

提案已历经 3 轮 adversarial 迭代，从 CI/CD 管线架构的视角被充分审查过。CI 专家关注的是系统正确性——`extractFileLineMap` 的解析能力、cap 机制的安全性、fallback 的完整性。这些是必要但不充分的条件。本审查从 agent 作为 fix task 消费者的角度，补充一个被忽视的维度：系统正确性之上的 agent 体验正确性。

## Section 2: Key Risk Identification

问题：上下文窗口的定义满足的是"搜索引擎"需求，不是"诊断"需求。

提案 Scope 部分规定："匹配行取前后各 2 行作为上下文窗口（共 5 行）"。这个设计假设关键信息集中在匹配行的紧邻位置。对于 Go 测试输出，这个假设部分成立——`file_test.go:42: expected status 200, got 500` 的前后行通常包含相关的断言信息。但对于需要 deep stack trace 才能定位 root cause 的场景，2 行上下文远远不够。

考虑这个真实场景：一个 HTTP handler 测试失败，错误输出是 `handler_test.go:42: expected status 200, got 500`。匹配行前后各 2 行可能给出的是测试文件中的断言上下文，但 root cause 在 `handler.go` 的某个 middleware 中——stack trace 可能跨越 8-10 行才到达实际出错的位置。Agent 拿到的 description 中只有 `handler_test.go` 的 5 行上下文，看不到 `handler.go:108` 的 middleware 逻辑。它必须自己做额外的文件读取和推理才能建立因果链。

这不是"信息量不够"的问题——增加上下文到 10 行也解决不了，因为 root cause 可能不在同一个 stack trace 的紧邻行中。这是"上下文窗口策略对 agent 诊断路径的匹配度"问题。提案的"前后各 2 行"是一个固定窗口，不考虑失败类型的差异。对于断言级别的失败（expected X, got Y），2 行窗口通常足够。但对于需要跨文件追踪的失败（如 panic、goroutine leak、deadlock），2 行窗口可能截断关键的因果链。

风险：per-file task 可能制造"信息孤岛"——agent 看到症状但看不到全貌。

提案按测试文件拆分 task，每个 task 只包含该文件相关的输出行。这个设计隐含一个假设：agent 处理 `handler_test.go` 的失败时，不需要看到 `middleware_test.go` 的失败信息。当两个测试文件的失败由同一个 root cause 引起时，这个假设不成立。

具体场景：`auth_test.go` 和 `user_test.go` 都因为 `auth.go` 的 token 验证 bug 而失败。`auth_test.go` 的失败信息是 "token expired immediately after creation"，`user_test.go` 的失败信息是 "user profile request returned 401"。这两个失败从不同角度指向同一个 root cause。如果同一个 agent 能同时看到两个失败信息，它可以更快地推断出 "token 过期导致所有需要认证的请求失败" 这个全貌。但拆分后，处理 `auth_test.go` 的 agent 只看到 token 过期，可能直接修改 token 过期逻辑；处理 `user_test.go` 的 agent 只看到 401，可能去改 user handler 的认证检查。两个 agent 各自修复了症状，但都没触及 root cause——或者更糟，两个修复互相矛盾。

提案的 Risk 表承认了这个问题（"同一根因 bug 导致多个测试文件失败时创建冲突修复任务"，Likelihood M, Impact M），缓解措施是 "在 task description 中列出其他可能编辑同一生产代码的相关任务 ID（RELATED_TASKS 字段），供 agent 避免并发冲突"。

但 `RELATED_TASKS` 字段的实际效果值得怀疑。它要求 agent 在执行修复前先"查看相关任务"——这增加了 agent 的认知步骤。更重要的是，agent 看到另一个 task 的 ID 后，它能否有效地利用这个信息？它需要：(1) 读取另一个 task 的 description，(2) 分析两个 failure 是否同源，(3) 决定是否协调修复策略。这是一个需要深度推理的 multi-step 过程，而 LLM agent 在处理"元任务信息"（关于其他 task 的信息）时的可靠性远低于处理"任务内信息"（直接相关的代码和错误）。

问题：description 格式示例展示的是一个精心挑选的理想场景，没有覆盖 agent 真正会遇到困难的场景。

提案给出的 description 示例：

```
handler_test.go
--- FAIL: TestGetUser (0.00s)
    handler_test.go:42: expected status 200, got 500
    handler_test.go:43: response body: {"error": "unauthorized"}
```

这是一个对 agent 极度友好的场景：失败信息明确（expected 200, got 500），行号精确（第 42 行），甚至给出了 response body 作为额外线索。Agent 看到这个 description 后的第一跳路径很清晰：读 `handler.go` 第 42 行附近，检查为什么返回 500。

但真实世界中的失败信息往往更晦涩。比如 Go 的 race condition 检测输出，一个 data race 的报告可能跨越 30+ 行、引用 4 个不同的 goroutine 和 2 个不同的文件。按文件分组后，`extractFileLineMap` 给 `handler_test.go` 的上下文窗口可能只包含 race report 的一个片段——agent 看到的是"data race at handler_test.go:42"，但看不到另一半在 `middleware.go:108` 的 write 操作。没有完整的 race pair，agent 根本无法定位 root cause。

再比如 panic 场景：`handler_test.go` 的测试触发了 panic，stack trace 经过 `router.go` → `middleware.go` → `handler.go` → `db.go`。按文件分组后，`handler_test.go` 的 task 只包含 stack trace 中引用 `handler_test.go` 的部分（加上前后 2 行上下文）。但 panic 的 root cause 在 `db.go:200`。Agent 拿到的 description 信息指向了一个与 root cause 无关的文件。

风险：`isTestFile` 和 `extractFileLineMap` 的双重失败模式会制造"幽灵任务"。

提案中存在两个独立的识别机制：`isTestFile`（命名约定匹配）识别测试文件，`extractFileLineMap`（输出解析）提取关联行。两个机制可以独立失败。当 `isTestFile` 成功但 `extractFileLineMap` 失败时（例如 Go 使用了 `go test -json` 输出格式），系统会为该测试文件创建一个 fix task，但 description 中只有文件名，没有失败信息。

从 agent 视角看，这是一个"幽灵任务"——task 告诉它"修复 handler_test.go 的失败"，但没有告诉它"哪里失败了、为什么失败"。Agent 的第一跳变成盲目的：它必须自己去运行测试、阅读测试代码、猜测失败原因。这比不拆分（至少有完整 output 可读）更糟糕——因为不拆分时 agent 至少能从完整 output 中自行定位；拆分后 agent 的信息被截断到单文件，它甚至不知道完整的 failure landscape 长什么样。

提案没有为 `extractFileLineMap` 返回空结果的情况定义 fallback 行为。Constraint 部分定义了 `isTestFile` 失败时的 fallback（按目录分组），但 `extractFileLineMap` 失败的 fallback 缺失。

问题：软上限的 overflow fallback 对 agent 来说是一种不可预测的体验退化。

提案规定："regression 路径最多创建 10 个 fix task，超出部分 fallback 到按目录分组"。当一个项目有 12 个测试文件失败时，前 10 个文件各自获得精确的 per-file task，后 2 个文件被合并成按目录分组的 task。

从 agent 视角看，这制造了两种截然不同的 task 体验：前 10 个 task 的 description 是精确的、聚焦的、信息密度高的；后 2 个 task 的 description 是模糊的、范围过大的、噪声多于信号的。同一次 regression run 中，agent 可能先拿到一个 per-file task（体验良好），然后拿到一个 directory-based task（体验退化）。这种不可预测的体验退化比一致的"不太好"更糟糕，因为它破坏了 agent 对 task 格式的预期稳定性。

更重要的是，overflow 后的 directory-based task 恰好是提案要解决的问题的复现——scope 过大的 fix task 导致 agent 卡死。这意味着提案只解决了"大部分情况"的问题，而把最坏情况留给了未改进的 fallback。这不是 fallback，这是部分回滚。

建议：在 task description 中增加"agent 执行提示"——不只是失败信息，还要包含推理引导。

当前的 description 格式是纯粹的信息展示：文件路径 + 失败输出行。这是一个被动的信息容器，agent 必须自己从中推导出修复策略。考虑到 LLM agent 的推理模式——它倾向于先理解任务结构，再执行——在 description 中增加轻量级的执行提示可以显著提高第一跳效率。

具体建议：在每个 fix task 的 description 开头增加一行结构化摘要：

```
TARGET: handler_test.go
FAILURE_COUNT: 3
FAILURES: TestGetUser, TestCreateUser, TestDeleteUser
SIGNAL: HTTP 500 errors in authentication-related tests
START_HERE: Read handler.go around line 108 (token validation)
```

`START_HERE` 不是强制指令，而是基于失败信息的推理提示。对于 agent 来说，这比让它自己从 output 行中推导出"应该先看哪个文件"要高效得多。`SIGNAL` 提供了跨失败的共性摘要，帮助 agent 快速判断是否为单一 root cause。

当然，自动生成 `START_HERE` 和 `SIGNAL` 需要 `extractFileLineMap` 具备比当前设计更强的分析能力。如果这超出了当前实现范围，至少应该确保 description 中包含完整的 stack trace（而不是截断到 2 行上下文），让 agent 有足够的信息自己做出推理。

风险：多 agent 并发时的"隐形竞态"——`RELATED_FILES` 不够，agent 需要知道的不仅是"哪些文件相关"，而是"谁正在编辑"。

提案的缓解措施 `RELATED_FILES`（列出其他 fix task 可能编辑的生产代码文件）试图解决并发冲突问题。但从 agent 的认知模型看，这个信息的可操作性很低。

Agent 的推理过程是："我需要修复 `handler_test.go` 的失败 → root cause 在 `auth.go` → 我要编辑 `auth.go`。"当 `RELATED_FILES` 告诉它"另一个 task 也可能编辑 `auth.go`"时，agent 面临一个决策困境：它应该等另一个 task 完成后再修？还是直接修？如果直接修，如何避免冲突？

LLM agent 没有锁机制、没有 merge 策略、没有 conflict resolution 的能力。它只能做两件事：(1) 忽略 `RELATED_FILES` 直接修（冲突概率高），(2) 先读另一个 task 的 description，尝试推理另一个 agent 的修复策略，然后尝试做"兼容"的修改。选项 (2) 的可靠性极低——它要求 agent 准确预测另一个 agent 的行为。

更务实的方案可能不是在 agent 层面解决并发问题，而是在系统层面避免并发：当两个 fix task 共享生产代码文件时，系统应该顺序调度它们，而不是并发调度。这超出了当前提案的范围，但作为风险识别，它值得被记录——因为当前的缓解措施（`RELATED_FILES`）在 agent 体验层面提供了一个"看起来有缓解"但实际上难以操作的信号。

问题：`extractFileLineMap` 返回 `map[string][]string` 丢失了行号信息，agent 看到的是"无定位"的文本块。

函数签名 `extractFileLineMap(output string) map[string][]string` 返回每个文件关联的输出行内容，但没有保留行号。当这些行被写入 task description 时，agent 看到的是：

```
handler_test.go
--- FAIL: TestGetUser (0.00s)
    handler_test.go:42: expected status 200, got 500
```

注意这个例子中行号 `42` 是包含在原始输出行文本中的——`handler_test.go:42:` 是 test framework 输出的一部分。但 `extractFileLineMap` 的调用者（`addRegressionFixTasks`）在组装 description 时，不知道这一行在原始 output 中的位置（第几行），也不知道上下文窗口覆盖的是原始 output 的哪些行。

更严重的问题是：当 `extractFileLineMap` 提取的行经过上下文窗口扩展和重叠去重后，最终的行列表可能与原始 output 的行顺序不一致（取决于去重和合并的实现方式）。Agent 看到的 description 中的行可能是乱序的，这会干扰它的因果推理。

建议：将 description 格式设计为"结构化诊断报告"而非"原始 output 片段"。

当前的 description 是从 output 中截取的原始行片段——它保留了 test framework 的原始格式，但没有经过信息结构化。Agent 需要从这些原始行中自行提取：哪个测试失败了、断言是什么、预期值和实际值是什么、出错位置在哪里。

更 agent-friendly 的格式是将 description 设计为结构化的诊断报告：

```
## Failure Summary for handler_test.go

### Failed Tests
- TestGetUser (line 42): expected status 200, got 500
- TestCreateUser (line 78): expected non-nil response, got nil

### Common Pattern
Both failures involve HTTP 500 errors from the /api/* endpoints.

### Raw Output (for detailed stack traces)
--- FAIL: TestGetUser (0.00s)
    handler_test.go:42: expected status 200, got 500
    handler_test.go:43: response body: {"error": "unauthorized"}
--- FAIL: TestCreateUser (0.00s)
    handler_test.go:78: expected non-nil response, got nil
    handler_test.go:79: response body: {"error": "unauthorized"}
```

这种结构化格式让 agent 在第一秒就获得了 failure summary（不需要逐行解析 raw output），同时保留了 raw output 供深度分析。这不需要 LLM 分析——`extractFileLineMap` 已经在做行提取，结构化只是改变输出格式而不是增加分析能力。

当然，这增加了 `extractFileLineMap` 的实现复杂度（需要返回结构化数据而非扁平字符串），但从 agent 体验的角度，这是投入产出比最高的改进之一。

建议：同一文件内的多失败应按"疑似共同根因"做二级分组。

提案按测试文件拆分 task。但一个测试文件中可能有 5-10 个失败，其中 3 个是因为同一个函数的 bug，另外 2 个是因为另一个函数的 bug。如果 description 只是按出现顺序列出所有失败行，agent 仍然需要自己做"哪些失败可能是同源的"分组推理。

`extractFileLineMap` 在提取行时已经知道每行引用的行号（如 `handler_test.go:42`、`handler_test.go:78`）。一个简单的启发式是：行号相近的失败（如 42 和 43）可能同源，行号差距大的（如 42 和 78）可能不同源。在 description 中按行号聚类展示失败信息，可以帮助 agent 更高效地规划修复策略。

这不是要求做 root cause 分析——只是按行号距离做简单的聚类，这在 `extractFileLineMap` 内部用极低的成本就能实现，但对 agent 的推理效率提升显著。
