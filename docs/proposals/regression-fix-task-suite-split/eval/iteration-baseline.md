# Evaluation Report: iteration-baseline

**Date**: 2026-05-29
**Evaluator**: CTO Adversary (Baseline)
**Document**: `docs/proposals/regression-fix-task-suite-split/proposal.md`

---

## Phase 1: Reasoning Audit

### 1. Problem -> Solution Trace

**Stated Problem**: Agent 卡死，因为 quality-gate 创建单个 fix task 覆盖 20+ 跨文件失败。

**Proposed Solution**: 按**测试文件**（而非 test suite 目录）拆分 fix task。

**Gap**: Proposal 识别出"agent 卡死"的直接原因是"concise error 只展示输出尾部，agent 看不到完整失败列表"（引用 lesson 原文第 18-19 行）。这意味着 **agent 卡死的真正原因可能是信息不足，而非 scope 过大**。Proposal 自己承认了这一点（Phase 0），但仍然将 Phase 1（按文件拆分）作为主体方案，Phase 0 仅作为前置验证。逻辑上，如果 Phase 0 已能解决问题，Phase 1 就是过度工程。

**Verdict**: Solution 解决了一个更难的问题（scope 过大），而真正原因（信息不足）可能用更简单的方案就能解决。Phase 0/1 的优先级关系正确，但 Phase 1 在 Phase 0 验证之前就被完整设计了，这表明作者已经假设 Phase 0 不够。

### 2. Solution -> Evidence Trace

**Claim**: "按测试文件拆分是 fix task 的最小 scope 单元"。

**Evidence**: 无。Proposal 没有提供任何证据证明按文件拆分优于按目录拆分（现有行为）或按 suite 拆分（lesson 建议）。实际上，lesson 的建议是按 `FAIL forge-tests/<suite>` 行分组（suite 目录级别），而 proposal 改为按 `_test.go` 文件分组——粒度更细但没有论证为什么更细更好。

**Counter-evidence**: Source code 显示 `groupFilesByDir` 已经实现了按目录分组，且 proposal 自己说"无法识别测试文件时 fallback 到按目录分组"。如果按目录分组已存在且可用，按文件分组的增量价值未被量化。

### 3. Evidence -> Success Criteria Trace

SC[1] (SourceTaskID bug fix) 与 Problem 定义无直接关系——它是独立 bug，不解决 agent 卡死问题。
SC[2-3] 验证拆分行为但未验证 agent 是否因此不再卡死。
SC[8] 仅验证 Go 命名约定——MVP 的局限性被成功标准显式承认。

**Missing SC**: 没有任何 SC 验证"agent 不再卡死"这一核心目标。所有 SC 都是实现层面的验证，没有效果验证。

### 4. Self-Contradiction Check

**Contradiction 1**: Proposal 说"regression 路径使用独立软上限（10 个），不受 maxFixTasksPerStep 硬上限约束"。理由是"按测试文件拆分后 scope 已收窄，降低循环触发概率"。但 SourceTaskID bug 未修复前，cap 本身就不生效（proposal 自己说的）。先修 bug 再绕过 cap——这两步的逻辑关系需要更清晰的论证。

**Contradiction 2**: Risk 表说"同一根因 bug 导致多个测试文件失败时创建冲突修复任务"（Likelihood M），mitigation 是"在 task description 中列出相关任务 ID"。但 In Scope 中没有列出"在 description 中注入 RELATED_TASKS"这一 deliverable。这不是在 scope 中，mitigation 是空承诺。

**Contradiction 3**: Proposal 说 MVP 仅支持 Go，但 Constraints 中说"后续迭代按需扩展其他语言的提取模式"。Scope out 列表也说"非 Go 语言的测试文件识别作为后续迭代"。然而 `extractFileLineMap` 的签名设计（`func extractFileLineMap(output string) map[string][]string`）没有任何语言参数或扩展点——扩展时需要改签名或加参数，当前设计不具备扩展性。

### 5. SC Consistency Deep-Dive

**Cluster: Regression 路径拆分**
- SC[2]: 4 个文件 4 个 task
- SC[4]: regression 软上限 10，不受硬上限约束
- SC[5]: 最多 10 个 fix task，超出 fallback 按目录分组
- SC[6]: 无法识别时 fallback 按目录分组

**Conflict**: SC[5] 说"超出 10 个按目录合并为综合 task"，但 SC[6] 说"无法识别时 fallback 按目录分组"。这两个 fallback 路径是否使用同一个 `addFixTask` 函数？如果是，那 `addFixTask` 也有 cap 限制（`maxFixTasksPerStep=3`），超过 3 个目录分组怎么办？这个边界条件未被覆盖。

**Cluster: 现有行为不变**
- SC[7]: compile/fmt/lint/unit-test 不受影响
- SC[1]: SourceTaskID 修复影响 `addSingleFixTask`

**Conflict**: SC[1] 修复 `addSingleFixTask` 中设置 `opts.SourceTaskID`。但 `addSingleFixTask` 被 `addFixTask` 调用，`addFixTask` 也被 compile/fmt/lint/unit-test 使用。修复 SourceTaskID 会让这些路径的 fix task 也带上 `"quality-gate"` 作为 SourceTaskID——这正确吗？proposal 只说"设置固定标识如 `quality-gate`"，但 compile 步骤的 fix task 也用 `"quality-gate"` 作为 source 是语义错误的。

---

## Phase 2: Rubric Scoring

### 1. Problem Definition (110 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Problem stated clearly | 32/40 | 核心问题"agent 卡死"清晰，但混淆了两个层面：scope 过大 vs 信息不足。Proposal 承认直接原因是信息不足（Phase 0），但主体方案解决的是 scope 过大（Phase 1）。问题定义自相矛盾。 |
| Evidence provided | 35/40 | 有具体案例（4 个 suite 20+ 失败），引用了 lesson 文档，但缺少重现步骤和 token/时间消耗的量化数据。"长时间无响应"和"被用户手动中断"是主观描述。 |
| Urgency justified | 22/30 | "每次多文件失败都会触发"是合理推理，但实际发生频率未知——是一个已知触发了 1 次的问题，还是反复出现的问题？缺少频率数据。 |

**Subtotal: 89/110**

### 2. Solution Clarity (120 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Approach is concrete | 35/40 | 函数名、签名、流程都有说明。但 `extractFileLineMap` 的 "叠加 `--- FAIL:` 块的缩进行解析"描述模糊——什么缩进？多少层？解析规则是什么？ |
| User-facing behavior described | 30/45 | 描述了 task 拆分后的结果，但**没有描述 agent（最终用户）实际体验的变化**。agent 拿到拆分后的 task 后是否不再卡死？卡死是否真的因为 scope 而不是信息不足？这个关键的用户体验问题没有回答。 |
| Technical direction clear | 30/35 | Go 代码路径清晰，但 `createFixTask` helper 的提取、`addRegressionFixTasks` 与 `addFixTask` 的调用关系图缺失。需要从文字描述中自行推断代码结构。 |

**Subtotal: 95/120**

### 3. Industry Benchmarking (120 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Industry solutions referenced | 30/40 | 引用了 GitHub Actions test grouping、JUnit XML、Sentry fingerprinting。但引用停留在名称层面，没有深入任何一种方案的具体实现或为什么在本项目中不适用。 |
| At least 3 meaningful alternatives | 22/30 | 7 个 alternatives（含 do nothing），但"Do nothing"被以"用户体验差"一句否决——这是稻草人论证。"改进 description 信息呈现"被标为 Phase 0，实际是本提案最有价值的替代方案，但被降级为先导验证而非正式替代。 |
| Honest trade-off comparison | 15/25 | Comparison table 的 Cons 列对选中方案（"同一根因 bug 创建多任务；加剧 claim priority 和 cross-feature pollution"）一笔带过，这两个问题在 Risk 表中 Likelihood 为 M 和 H——高可能性的严重副作用在对比表中轻描淡写。 |
| Chosen approach justified | 18/25 | 选中"按测试文件分组"但 lesson 建议的是按 suite 目录分组。从 suite 目录到测试文件是粒度细化，proposal 没有论证为什么需要更细的粒度。 |

**Subtotal: 85/120**

### 4. Requirements Completeness (110 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Scenario coverage | 32/40 | Happy path、单文件、非标准命名、栈 trace 引用辅助文件、全部通过——覆盖合理。但缺少：(1) output 完全没有 `file:line` 模式的情况（`sourceFileRe` 无匹配）；(2) 同一文件在不同 `--- FAIL:` 块中出现多次的去重场景。 |
| Non-functional requirements | 28/40 | "10000 行 output < 50MB"的约束没有论证依据（为什么是 50MB？实际测试输出多大？）。"性能可忽略"是断言而非测量。缺少对 `extractFileLineMap` 在极端 case（如 100000 行 output）下的行为分析。 |
| Constraints & dependencies | 22/30 | 依赖命名约定已说明，但未说明 `sourceFileRe` 在非 Go 输出中的匹配率——`([\w][\w./-]*\.\w{1,10})(?::\d+){1,2}` 这个正则匹配 `file.go:42` 也能匹配 `file.py:42`，但 Python 测试输出的格式不同（`FAIL test_file.py::test_func`），实际匹配效果未知。 |

**Subtotal: 82/110**

### 5. Solution Creativity (100 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Novelty over industry baseline | 15/40 | Proposal 自己承认"无创新"。按文件分组是 CI 系统的常规做法。MVP 限 Go 的决策正确但不算创新。 |
| Cross-domain inspiration | 10/35 | Sentry fingerprinting 被引用但被拒绝，没有跨域借鉴。测试文件命名约定是最直接的机械匹配，没有从其他领域引入洞察。 |
| Simplicity of insight | 20/25 | Phase 0（改 description 信息呈现）是简洁的洞察，承认直接原因并先验证。但 Phase 1 的 `extractFileLineMap` + 上下文窗口 + 重叠去重的复杂度不低——为了解决"agent 看不到失败列表"问题，引入了 output 解析引擎。 |

**Subtotal: 45/100**

### 6. Feasibility (100 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Technical feasibility | 35/40 | 技术上完全可行。但 `extractFileLineMap` 的"叠加 `--- FAIL:` 块的缩进行解析"在 Go 测试输出中并非总是有缩进（`go test -v` 的输出格式取决于测试框架），解析规则的鲁棒性存疑。 |
| Resource & timeline | 22/30 | "Phase 0 半天 + Phase 1 1-2 天"的时间估计合理，但遗漏了测试编写时间。In Scope 列出"更新或新增单元测试覆盖新逻辑"——`extractFileLineMap` 的测试需要覆盖多种 Go 测试输出格式，这不是 trivial 的。 |
| Dependency readiness | 28/30 | 无外部依赖，`sourceFileRe` 已存在。依赖成熟度好。 |

**Subtotal: 85/100**

### 7. Scope Definition (80 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| In-scope items concrete | 22/30 | 大多数 item 是 deliverable，但 "关联算法" 列出了 5 步关联规则 + 示例，然后说"保留 maxFixTasksPerStep 用于非 regression 步骤"——这些实现细节混合在 scope 定义中，scope 和 solution 的边界模糊。 |
| Out-of-scope listed | 22/25 | 列出了 7 个 out-of-scope 项，含具体文档引用。但"RELATED_TASKS 字段注入"在 Risk mitigation 中提到却不在 In Scope 也不在 Out of Scope——悬空承诺。 |
| Scope bounded | 18/25 | Phase 0 + Phase 1 + SourceTaskID bug fix = 3 个独立工作包，在 2-3 天内可行。但 Phase 0 的验证结果可能改变 Phase 1 的必要性，scope 边界实际上依赖于 Phase 0 的结果。 |

**Subtotal: 62/80**

### 8. Risk Assessment (90 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Risks identified | 24/30 | 8 个 risks，覆盖面广。但缺少一个关键风险：**Phase 0 验证成功后 Phase 1 的投入浪费**。如果改进 description 就够了，Phase 1 的设计和实现成本完全浪费。 |
| Likelihood + impact rated | 22/30 | "claim priority" 评为 H/M 合理。"同一根因多任务"评为 M/M 偏乐观——20+ 失败跨 4 个 suite 很可能由 1-2 个根因引起，此时创建 4 个 task 中 2-3 个是浪费的。 |
| Mitigations actionable | 20/30 | "RELATED_TASKS 字段供 agent 避免并发冲突"——这个 mitigation 不在 scope 中，无法执行。"fix task ID 改为数字前缀"标记为长期方案，短期无行动。"沿袭现有 mark as skipped 工作流"是接受风险而非缓解。 |

**Subtotal: 66/90**

### 9. Success Criteria (80 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Measurable and testable | 22/30 | SC[1-8] 大部分可通过单元测试验证。但 SC 没有验证核心目标——"agent 不再卡死"。SC 验证的是实现是否正确，而非方案是否有效。 |
| Coverage complete | 18/25 | 缺少：(1) Phase 0 的成功标准（description 改进后 agent 是否不再卡死？如何衡量？）；(2) `createFixTask` helper 提取的验证；(3) fallback 到 `addFixTask` 时的 cap 行为验证。 |
| SC internal consistency | 16/25 | SC[1] 修复 SourceTaskID 对 compile/fmt/lint/unit-test 路径有副作用（这些路径的 fix task 也会带上 `"quality-gate"` source），但 SC[7] 说这些路径不受影响——矛盾。SC[4-5] 的 regression 软上限与现有 cap 的交互在 fallback 场景下未定义。 |

**Subtotal: 56/80**

### 10. Logical Consistency (90 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Solution addresses stated problem | 25/35 | Proposal 承认直接原因是信息不足（Phase 0），但主体方案解决 scope 过大（Phase 1）。Solution 部分解决了 Problem，但可能过度解决——如果 Phase 0 足够，Phase 1 是不必要的复杂性。 |
| Scope <-> Solution <-> SC aligned | 20/30 | Scope 中的"关联算法 5 步"没有对应 SC。SC[1] 对应 Scope 中的 SourceTaskID 修复但影响超出 regression 路径。Risk mitigation 中承诺的 RELATED_TASKS 不在 Scope 也不在 SC 中。 |
| Requirements <-> Solution coherent | 20/25 | Requirements 中列出了多种场景，Solution 覆盖了大部分。但"栈 trace 引用辅助测试文件"场景在 Solution 中没有对应的特殊处理——它只是"两个文件各创建 fix task"，这不完全是一个 solution，更多是接受误报。 |

**Subtotal: 65/90**

---

## Phase 3: Blindspot Hunt

### [blindspot-1] SourceTaskID 修复的副作用未被分析

Proposal 要求修复 `addSingleFixTask` 中设置 `opts.SourceTaskID = "quality-gate"`。但 `addSingleFixTask` 被 compile/fmt/lint/unit-test 步骤也调用（通过 `addFixTask`）。修复后，这些步骤创建的 fix task 也会带 `SourceTaskID: "quality-gate"`。

在 `add.go` 中，`opts.SourceTaskID != ""` 会触发 dedup check（line 156-169）和 source resolution（line 179-191）。这意味着 compile 步骤创建的 fix-1 会阻止后续 compile fix task 的创建（因为 `hasActiveFixTasks` 会返回 fix-1）——这可能是预期行为，但 proposal 没有分析。

更严重的是：`SourceTaskID: "quality-gate"` 不是一个有效的 task ID。`FindTask(index, "quality-gate")` 会失败，但 `hasActiveFixTasks` 仍会基于这个值搜索。这可能导致所有步骤的 fix task 共享同一个 source group，使 dedup 逻辑意外触发跨步骤去重。

### [blindspot-2] extractFileLineMap 与 extractSourceFiles 的关系未定义

Proposal 说 `extractFileLineMap` "替代 `extractSourceFiles` 用于 regression 路径"，但没说两者是否共享代码。如果完全独立实现，同一个 `sourceFileRe` 正则会在两个函数中各维护一份。如果共享，`extractFileLineMap` 的上下文窗口逻辑会耦合到 `extractSourceFiles` 的调用者。

### [blindspot-3] Phase 0 和 Phase 1 的决策门控未定义

Phase 0 被定位为"先验证信息不足是否是卡死的真正原因"，但没有定义：
- Phase 0 成功的判断标准是什么？agent 不再卡死？agent 仍卡死但原因不同？
- 如果 Phase 0 验证后发现信息不足不是原因，Phase 1 的设计是否需要调整？
- 如果 Phase 0 验证发现信息不足是原因但 agent 仍偶尔卡住，Phase 1 是否仍然需要？

没有决策门控，Phase 0 和 Phase 1 实际上是串行执行的两个独立方案，而非验证驱动的条件分支。

### [blindspot-4] Surface inference 对拆分后 task 的影响

`addSingleFixTask` 调用 `inferSurface` 来推断 surface key/type。按文件拆分后，每个 fix task 只有单个测试文件路径。`inferSurface` 基于 `forgeconfig.MatchSurface` 匹配文件到 surface。测试文件（`*_test.go`）通常不直接匹配 surface 配置（surface 通常匹配 `src/` 或 `pkg/` 下的生产代码路径）。这意味着拆分后的 fix task 可能全部 surface inference 失败，导致所有 task 的 SurfaceKey/SurfaceType 为空——这可能影响 dispatcher 的 claim 逻辑。

### [blindspot-5] addRegressionFixTasks 的调用位置

Proposal 说"`runTestRegression` 失败时调用新函数替代 `addFixTask`"。但 source code 显示有两条 regression 路径：`runTestRegressionLegacy`（line 260）和 `runTestRegressionSurface`（line 289）。两条路径都调用 `addFixTask`。Proposal 只说替换一处，没有明确是否两条路径都替换。

---

## Score Summary

| Dimension | Score | Max |
|-----------|-------|-----|
| Problem Definition | 89 | 110 |
| Solution Clarity | 95 | 120 |
| Industry Benchmarking | 85 | 120 |
| Requirements Completeness | 82 | 110 |
| Solution Creativity | 45 | 100 |
| Feasibility | 85 | 100 |
| Scope Definition | 62 | 80 |
| Risk Assessment | 66 | 90 |
| Success Criteria | 56 | 80 |
| Logical Consistency | 65 | 90 |
| **Total** | **730** | **1000** |

---

## Top Attacks (prioritized by severity)

1. **Success Criteria**: SC 验证实现正确性但不验证方案有效性——没有 SC 衡量"agent 不再卡死"。核心目标在 SC 中完全缺失。**改进：增加端到端 SC，如"给定 20+ 跨文件失败的 regression output，拆分后的单个 fix task 可在 N 分钟内被 agent 完成处理"。**

2. **Logical Consistency**: SourceTaskID 修复用固定值 `"quality-gate"` 影响所有步骤（compile/fmt/lint/unit-test/regression），但 SC[7] 声称非 regression 步骤不受影响——矛盾。**改进：将 SourceTaskID 设置限定在 regression 路径的 `addRegressionFixTasks` 中，而非全局修改 `addSingleFixTask`。**

3. **Success Criteria**: Phase 0 没有对应的成功标准。In Scope 列出了 Phase 0 作为 deliverable，但 SC 全部针对 Phase 1。**改进：增加 SC "Phase 0 实施后，agent 能在 task description 中看到完整 `--- FAIL` 行列表而不必读 raw-output.txt"。**

4. **Logical Consistency**: Risk mitigation 承诺"在 task description 中列出相关任务 ID（RELATED_TASKS 字段）"，但该 deliverable 既不在 In Scope 也不在 Out of Scope。**改进：要么将 RELATED_TASKS 注入加入 In Scope，要么删除该 mitigation。**

5. **Scope Definition**: Phase 0 的结果可能使 Phase 1 不必要，但 Phase 1 已被完整设计。缺少决策门控定义。**改进：明确 Phase 0 的通过/失败标准，以及 Phase 1 的触发条件。**

6. **Industry Benchmarking**: Lesson 建议按 suite 目录分组（`FAIL forge-tests/<suite>`），Proposal 改为按测试文件分组。从目录到文件是粒度细化，但没有论证为什么更细更好。**改进：提供按目录 vs 按文件的场景对比，展示按文件分组的增量价值。**

7. **Logical Consistency**: `addRegressionFixTasks` 的软上限（10）超出后 fallback 到 `addFixTask`，但 `addFixTask` 有硬上限（3）。fallback 路径的 cap 交互未定义。**改进：明确 fallback 时是否绕过 `addFixTask` 的 cap 检查。**

8. **Requirements Completeness**: "10000 行 output < 50MB"约束无依据。Go `go test -v` 的实际输出大小未测量。**改进：用实际 regression 输出测量典型和峰值大小，用数据支撑约束。**
