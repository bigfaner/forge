---
iteration: 4
title: "CTO Adversary Rubric Scoring — Iteration 4"
date: 2026-05-25
scorer: CTO adversary (blind annotated review)
baseline: iteration-3 (858/1000)
---

# 评分报告：init-justfile Surface 感知 + 测试编排简化

## Phase 1 — 推理链审计

### 论证链路追踪

1. **Problem → Solution**：两个问题的捆绑论证保持稳固。iteration-3 以来，论证链无结构性变化。三条捆绑理由（双重设计成本、参数分散、零迁移成本窗口）仍然成立。

2. **Solution → Evidence**：证据性质声明保持坦诚。8 个样本的量化数据仍然是唯一的具体数字，但这是客观约束（v3.0.0 未发布），而非论证缺陷。

3. **Evidence → Success Criteria**：14 条成功标准覆盖了范围定义中的所有主要交付物，包括 config schema 子方案（第 14 条）和 scope 迁移（第 13 条）。

4. **自相矛盾检测**：
   - **矛盾 A（iteration-3 矛盾 E，已修复）**：HARD-GATE "禁止重试"与最坏情况的矛盾已修复——最坏情况从"不必要的 teardown 后重试"修改为"执行不必要的 teardown 后中止"，并添加了明确注释"最坏情况的'不必要 teardown'仅指 LLM 可能先执行了完整的 teardown 再中止（而非跳过 teardown 直接重试）"。
   - **矛盾 B（iteration-3 blindspot-2，已修复）**：`wmic` 弃用问题已修复——Windows 命令行校验从 `wmic process where ProcessId=<pid> get CommandLine` 改为 `Get-CimInstance Win32_Process`（失败时回退 `tasklist /V`）。
   - **矛盾 C（iteration-3 blindspot-4，已修复）**：用户手动修改配方被覆盖的问题已修复——新增 `# user-customized` 标记保护机制，init-justfile 检测到此标记时跳过覆盖并输出警告。
   - **矛盾 D（iteration-3 blindspot-5，已修复）**：字母序消歧缺乏语义依据已修复——从"字母序排列第一个 key"改为"YAML 映射保留插入顺序"，并附带了理由"反映用户在配置中表达的逻辑优先级"。
   - **新矛盾 E**：`# user-customized` 标记保护机制的粒度不足——该机制是"配方级别"的（检查整个配方是否有标记），但用户可能只修改了配方的部分内容（如添加了一个环境变量），init-justfile 会因标记而完全跳过覆盖，导致用户错过 surface 规则的其他改进（如 probe 重试逻辑优化）。这是一个"全有或全无"的保护策略，缺乏细粒度控制。
   - **新矛盾 F**：probe 重试差异化机制（line 79）声明 run-tests 区分连接失败类型（ECONNREFUSED vs 超时 vs 端口冲突），但 `just probe` 是一个独立配方——run-tests 调用 `just probe` 后只获得退出码（0 或 1），无法区分失败类型。提案声称的"区分连接失败类型以加速定位"在当前架构下（run-tests 通过退出码与 just probe 交互）无法实现——run-tests 只能看到 probe 失败了（exit 1），无法知道是 ECONNREFUSED、超时还是 EADDRINUSE。这些差异化逻辑只能在 `just probe` 配方体内部实现（配方体内部可以区分并在 stderr 输出不同消息），但 run-tests 无法基于失败类型调整重试策略。

### Pre-score Anchors

- **Anchor 1**：iteration-4 相比 iteration-3 的提案文本无结构性变化——pre-revised 标记区域在 iteration-3 已评估。本次评估需确认 iteration-3 指出的 5 个 blindspot 的修复质量，并检查修复是否引入了新问题。
- **Anchor 2**：4 个 blindspot（wmic 弃用、字母序消歧、HARD-GATE 矛盾、用户编辑保护）的修复都是高质量的——直接回应了攻击点，没有引入明显的回退。
- **Anchor 3**：新矛盾 E（user-customized 全有或全无保护）和矛盾 F（probe 重试差异化与退出码架构不兼容）是本次评估发现的主要新问题。

---

## Phase 2 — Rubric Scoring

### 1. Problem Definition (90/110)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Problem stated clearly | 39/40 | 两个问题陈述明确，捆绑论证三条理由逻辑自洽。扣分：问题 1（surface 不感知）的"当前有多少种 surface 类型的项目受影响"仍然缺失——虽然 v3.0.0 未发布，但 Forge 内部的示例项目或测试 fixture 应该可以提供粗略数据。 |
| Evidence provided | 26/40 | 证据性质声明保持坦诚。8 个样本的量化数据局限性已充分声明。扣分：iteration-3 和 iteration-4 均未改善证据质量——没有新增代码审计发现、没有新的使用场景分析、没有来自内部 dogfooding 的数据。 |
| Urgency justified | 25/30 | 零迁移成本窗口仍然是最强的论证。扣分：未回答"v3.0.0 延期或 beta/preview 版本是否关闭零迁移窗口"——iteration-3 已指出，仍未修复。 |

### 2. Solution Clarity (106/120)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Approach is concrete | 40/40 | Surface 编排模式表格、fallback 链、scope 迁移 4 阶段、原子性约束、用户编辑保护——方案的具体程度已达到可直接拆分为任务的水平。 |
| User-facing behavior described | 43/45 | iteration-4 的改善：(1) `# user-customized` 标记保护使用户编辑不被意外覆盖（iteration-3 blindspot-4 修复）；(2) scope 消歧从字母序改为声明顺序（更符合直觉）。扣分：(1) `# user-customized` 的"全有或全无"粒度——用户添加了一个环境变量后，后续所有 surface 规则的改进都无法自动应用；(2) probe 重试差异化（line 79）声称 run-tests 区分连接失败类型，但 run-tests 通过退出码与 just probe 交互，只能看到 0/1，无法实现基于失败类型的差异化重试策略——这是用户期望的行为描述与实际可实现行为之间的差距。 |
| Technical direction clear | 23/35 | `Get-CimInstance` 替代 `wmic` 修复了 Windows 弃用问题（iteration-3 blindspot-2 修复）。扣分：(1) probe 重试差异化策略声称 run-tests 可以区分 ECONNREFUSED/超时/EADDRINUSE，但架构上 run-tests 只能获得退出码——需要在 `just probe` 配方体内部实现差异化重试（配方体根据 curl 的错误类型调整重试计数），而非 run-tests 层面实现。提案当前的描述暗示 run-tests 层面做差异化，这与架构不符。(2) 跨平台配方双变体（shell vs PowerShell）的维护成本仍未充分分析——init-justfile 的 LLM 需要为每个跨平台配方生成语法完全不同的两个版本。 |

### 3. Industry Benchmarking (107/120)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Industry solutions referenced | 38/40 | 9 个成熟方案对比保持完整（Docker Compose、K8s、Cypress、Makefile、GitHub Actions、Playwright、Vitest、Testcontainers）。扣分：仍缺少 Bazel/Please 等构建系统中"测试编排作为构建规则"的模式对比。 |
| At least 3 meaningful alternatives | 27/30 | 4 个替代方案 + "不做"。Go 代码方案作为"长期方向"保留是诚实的。扣分："Go 代码直接管理进程生命周期"方案标注为"采纳其核心思想作为兜底机制"——这不是完全独立的替代方案。 |
| Honest trade-off comparison | 21/25 | "为何不复用测试框架内建编排"和"为何不采用 Testcontainers 模式"的分析保持完整。扣分：just 配方跨平台双变体维护成本的量化分析仍然缺失。 |
| Chosen approach justified against benchmarks | 21/25 | Forge 的差异化定位（LLM agent + justfile 文本协议 + 框架无关）清晰。扣分：未讨论"为何不在 SKILL 内部用轻量脚本管理进程"的中间方案——这个选项介于纯 justfile 和 Go 子命令之间。 |

### 4. Requirements Completeness (95/110)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Scenario coverage | 38/40 | 7 个关键场景 + 多 surface 同类型 journey 过滤。扣分：(1) init-justfile 多次运行的增量更新场景仍不在"关键场景"中（仅在"仲裁规则"中通过 `# user-customized` 保护间接覆盖）；(2) `just probe` 配方体内部的连接失败类型差异化如何传达给 run-tests（通过 stderr 输出？通过特定退出码？）未在场景中说明。 |
| Non-functional requirements | 35/40 | just 版本 >= 1.4.0、可观测性格式、PID 存活检查保持完整。扣分：(1) "性能"仍只约束 init-justfile 加载时间，未约束 run-tests 端到端编排总耗时（dev 启动 + probe 等待 + test 执行 + teardown 清理的上限）；(2) "跨平台兼容"验证方式仍为"各平台手动验证"。 |
| Constraints & dependencies | 22/30 | Surface 信息源优先级规则、test.execution 引用审计清单、GetConfigValue 扩展保持完整。扣分：(1) just >= 1.4.0 版本检查机制仍未定义——init-justfile 在生成使用 `[linux]`/`[windows]` attribute 的配方前是否检查 just 版本？(2) `GetConfigValue` 扩展"不破坏现有键"的依据仍然缺乏——现有键是否有单元测试？ |

### 5. Solution Creativity (74/100)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Novelty over industry baseline | 31/40 | 规则文件物理独立但逻辑同构的设计（init-justfile 和 run-tests 各自持有独立规则文件副本但共享编排序列契约）是 Forge 特有的创新。`# user-customized` 标记保护虽然简单但有效。PID 存活检查在 probe 循环中加速崩溃检测保持优雅。 |
| Cross-domain inspiration | 23/35 | 借鉴来源略有扩展（Testcontainers 的 Ryuk sidecar → test-state.json 恢复机制）。scope 兼容层的声明顺序消歧来自 YAML 映射语义。 |
| Simplicity of insight | 20/25 | probe 循环中的 PID 存活检查、`# user-customized` 单行注释保护、scope 声明顺序消歧——都是简单但有效的方案。扣分：probe 重试差异化的描述过于复杂且与架构不符——如果差异化只能在配方体内部实现，不应声称 run-tests 做差异化。 |

### 6. Feasibility (84/100)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Technical feasibility | 38/40 | `Get-CimInstance` 替代 `wmic` 修复了 Windows 兼容性问题。`# user-customized` 保护机制简单可行。扣分：probe 重试差异化策略的实现路径不清晰——声称 run-tests 做差异化但 run-tests 只能获得退出码。 |
| Resource & timeline feasibility | 25/30 | config schema 子方案降级路径完整。scope 统一迁移的原子性约束明确。扣分：15-20 个编码任务的范围仍然偏大，scope 迁移的同一 PR 约束在实践中意味着巨型 PR。 |
| Dependency readiness | 21/30 | Surface 检测已就位。PowerShell 依赖已声明。`Get-CimInstance` 可用性合理（Windows PowerShell 5.x 和 PowerShell 7.x 均支持）。扣分：just >= 1.4.0 版本检查机制未定义。 |

### 7. Scope Definition (76/80)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| In-scope items are concrete | 29/30 | 每个范围内项都是可交付的。config schema 变更有降级路径、边界约束和影响面评估。scope 统一迁移有原子性约束和兼容层策略。 |
| Out-of-scope explicitly listed | 21/25 | 列了 6 项范围外。回滚通过 git revert。扣分：(1) 迁移指南（从 test.execution 到 just 配方的用户文档）是否在范围内仍未明确；(2) `# user-customized` 标记机制的"全有或全无"策略是否覆盖所有需要的细粒度保护场景——如果用户只想保留自定义的环境变量但接受 probe 逻辑更新，当前方案无法支持。 |
| Scope is bounded | 26/25 | 原子性约束 + 兼容层保留一个 minor version 的策略提供了好的时间约束。config schema 子方案边界清晰。 |

### 8. Risk Assessment (84/90)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Risks identified | 28/30 | 8 个风险覆盖了主要场景。HARD-GATE 被违反风险的缓解措施更精确了（去除了"重试"的矛盾描述）。扣分：(1) probe 重试差异化策略与退出码架构的不兼容性未作为风险列出——如果 run-tests 无法区分失败类型，line 79 描述的差异化行为无法实现，这是一个设计缺口；(2) `# user-customized` 保护可能导致用户错过 surface 规则改进的风险未列出。 |
| Likelihood + impact rated | 28/30 | 大部分评估合理。"run-tests 无法感知 surface"标为"低/高"——考虑到 surface 感知依赖 config.yaml 或 forge surfaces CLI（已就绪），评估可接受。 |
| Mitigations are actionable | 28/30 | HARD-GATE 分层兜底机制更精确（去除了矛盾描述）。回滚计划保持完整。扣分：回滚后用户自定义的 `# user-customized` 标记的配方处理未说明。 |

### 9. Success Criteria (77/80)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Criteria are measurable and testable | 53/55 | 14 条成功标准都是可验证的 checklist 或端到端测试。config schema 变更有了验证标准（第 14 条）。扣分：(1) "所有生成的配方通过 --dry-run 验证"——dry-run 只验证语法不验证运行时行为（iteration-1 到 iteration-4 均指出，但第 8 条"运行时端到端验证"部分补偿了这一缺口）；(2) "语言模板与 surface 规则的配方职责边界清晰（语言级 vs 编排级），无同名冲突"——"清晰"仍不是可量化的验证标准。 |
| Coverage is complete | 24/25 | 覆盖了范围内的主要交付物。扣分：`# user-customized` 保护机制的有效性验证未在成功标准中体现。 |

### 10. Logical Consistency (84/90)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Solution addresses the stated problem | 34/35 | Surface 感知解决编排差异成立，委托层简化消除冗余成立，捆绑论证三条理由逻辑自洽。HARD-GATE 矛盾已修复。 |
| Scope ↔ Solution ↔ Success Criteria aligned | 28/30 | config schema 变更、scope 迁移、用户编辑保护都有对应的成功标准。扣分：`# user-customized` 保护机制在方案和范围内有描述但成功标准无对应验证条目。 |
| Requirements ↔ Solution coherent | 22/25 | 下游集成契约表格与方案一致。scope 值域迁移细则完整。`# user-customized` 保护使仲裁规则更完整。扣分：(1) probe 重试差异化（line 79）声称 run-tests 区分连接失败类型，但 run-tests 通过退出码交互，无法实现基于失败类型的差异化策略——这是需求描述与方案架构之间的不一致；(2) 端口冲突预防的 best-effort 检查可能输出误导性错误信息的风险仍然存在。 |

---

## Scoring Summary

| Dimension | Score | Max |
|-----------|-------|-----|
| 1. Problem Definition | 90 | 110 |
| 2. Solution Clarity | 106 | 120 |
| 3. Industry Benchmarking | 107 | 120 |
| 4. Requirements Completeness | 95 | 110 |
| 5. Solution Creativity | 74 | 100 |
| 6. Feasibility | 84 | 100 |
| 7. Scope Definition | 76 | 80 |
| 8. Risk Assessment | 84 | 90 |
| 9. Success Criteria | 77 | 80 |
| 10. Logical Consistency | 84 | 90 |
| **Total** | **877** | **1000** |

---

## Phase 3 — Blindspot Hunt

### [blindspot-1] probe 重试差异化策略与退出码架构不兼容

提案 line 79 声明 run-tests 区分连接失败类型（ECONNREFUSED vs 超时 vs EADDRINUSE）以加速定位。但 run-tests 调用 `just probe` 后只能获得退出码（0 或 1），无法区分失败类型。这意味着提案描述的"run-tests 区分连接失败类型"在实际架构中无法实现。

**可能的修正路径**：(1) 差异化重试逻辑完全在 `just probe` 配方体内部实现——配方体内部区分 ECONNREFUSED（继续重试）和超时（增加计数权重），run-tests 无需感知失败类型；(2) 使用特定退出码传递失败类型（如 exit 2 = ECONNREFUSED，exit 3 = 超时，exit 4 = 端口冲突）——但这要求 run-tests 理解 just probe 的退出码语义，增加了耦合。

**引用**："**probe 重试差异化**：run-tests 区分连接失败类型以加速定位——(1) **连接被拒绝**（ECONNREFUSED）：端口无进程监听，继续重试；(2) **连接超时**：端口有进程但不响应 HTTP，增加计数权重（每次超时计为 2 次重试）加速退出；(3) **端口冲突**（日志 EADDRINUSE）：立即退出"

### [blindspot-2] `# user-customized` 全有或全无的粒度问题

用户在 justfile 的 `test` 配方中添加了 `# user-customized` 标记（因为修改了一个环境变量），后续 init-justfile 执行时会完全跳过该配方的覆盖。这意味着即使用户只做了一个小修改，也会错过 surface 规则的所有改进（如 probe 重试逻辑优化、新增的端口冲突检测等）。

这不是一个致命缺陷，但在长期使用中会导致用户的 justfile 逐渐与推荐的 surface 规则模板产生偏差。提案应至少说明：(1) 当 init-justfile 跳过覆盖时，应同时输出"当前配方版本与推荐版本的差异摘要"，帮助用户判断是否需要手动更新；(2) 或者提供一个 `--force-regenerate` 选项，允许用户在了解风险的情况下强制更新。

**引用**："**用户编辑保护**：init-justfile 在覆盖编排级配方前，检查 justfile 中是否包含用户手动编辑标记（`# user-customized` 注释行）。若目标配方已有此标记，init-justfile 跳过覆盖并输出警告"

### [blindspot-3] YAML 映射插入顺序假设的正确性

scope 兼容层的消歧规则从字母序改为"YAML 映射保留插入顺序"。这个假设依赖于 Go 的 `yaml.v3`（或 Forge 使用的 YAML 库）在反序列化 `map[string]string` 时保留插入顺序。但 Go 的 `map[string]string` 本身是无序的——即使 YAML 库按插入顺序解码，映射到 Go map 后顺序信息丢失。

除非 Forge 的 `surfaces` 字段使用了有序数据结构（如 `yaml.Node` 或第三方有序 map），否则"声明顺序靠前的 key"在 Go 运行时中不可用。这是一个隐含的技术假设，需要验证。

**引用**："若有多个匹配，优先匹配 surfaces map 中声明顺序靠前的 key（YAML 映射保留插入顺序，提供确定性选择，且反映用户在配置中表达的逻辑优先级）"

### [blindspot-4] just >= 1.4.0 版本检查的缺失链路

NFR 要求 just >= 1.4.0（因为 `[linux]`/`[windows]` recipe attribute 在此版本引入）。但提案未定义版本检查的触发时机：init-justfile 在生成使用这些 attribute 的配方前是否检查 just 版本？如果用户安装了 just 1.3.x，init-justfile 生成的包含 `[linux]`/`[windows]` attribute 的配方会在执行时报 "unknown attribute" 错误。

提案的 config schema 变更子方案涉及 `forge config get` 的扩展，可以在 `init-justfile` 的前置检查步骤中添加 `just --version` 检查。但这个检查是否在范围内、由谁实现、失败时如何处理（拒绝生成？降级为单平台配方？），均未说明。

**引用**："just >= 1.4.0（`[linux]`/`[windows]` recipe attribute 在 just 1.4.0 引入，1.0-1.3.x 不支持此功能，会报 'unknown attribute' 错误）"

### [blindspot-5] run-tests 编排总耗时无上限约束

NFR "性能"仅约束 init-justfile 的 surface 规则加载时间（不超过 1 秒），但未约束 run-tests 端到端编排的总耗时。对于一个 web surface 项目，完整编排序列为：dev 启动（后台）+ probe 重试（30 次 x 2 秒 = 最长 60 秒）+ test 执行 + teardown 清理。如果 dev server 启动缓慢（如大型 Next.js 项目首次编译），probe 可能持续接近 60 秒才成功。加上 test 执行时间，总编排时间可能超过 5 分钟。

虽然 config schema 子方案保留了 `test.timeout`（默认 300 秒），但这个超时是否覆盖了整个编排序列（dev + probe + test + teardown）还是仅覆盖 test 步骤？如果是前者，300 秒可能不够；如果是后者，probe 的 60 秒超时是否可配置？

**引用**："test: timeout: 300"

---

## Bias Detection Report

**Pre-revised annotated regions**: 10 annotated paragraphs/blocks (lines 78, 156, 174, 187, 203, 332, 335, 509, 524, 527)

Attacks found in annotated regions:
1. [Logical Consistency] probe 重试差异化与退出码架构不兼容 (blindspot-1) — line 79 pre-revised:medium 上下文区域
2. [Logical Consistency] YAML 插入顺序假设在 Go map 中不成立 (blindspot-3) — line 338 pre-revised:medium 上下文区域（消歧规则）
3. [Feasibility] just >= 1.4.0 版本检查触发时机未定义 (blindspot-4) — line 365 NFR 区域（无 pre-revised 标记但由 pre-revised:medium line 332/335 的兼容层修订间接暴露）

Annotated region attacks: 3 attack points / 10 annotated paragraphs = density 0.30

Unannotated regions: ~200 paragraphs

Attacks in unannotated regions:
1. [Problem Definition] 证据质量未改善（iteration-3 到 iteration-4 无新数据）
2. [Solution Clarity] `# user-customized` 全有或全无粒度 (blindspot-2)
3. [Solution Clarity] 跨平台双变体维护成本未量化
4. [Industry Benchmarking] 缺少 Bazel/Please 构建规则编排模式
5. [Industry Benchmarking] 未讨论 SKILL 内轻量脚本中间方案
6. [Requirements] just >= 1.4.0 版本检查机制未定义 (blindspot-4)
7. [Requirements] GetConfigValue 现有键测试覆盖未说明
8. [Requirements] run-tests 端到端编排总耗时无上限 (blindspot-5)
9. [Feasibility] 15-20 个编码任务范围仍然偏大
10. [Success Criteria] dry-run 不验证运行时（4 个 iteration 均指出）
11. [Success Criteria] "清晰"仍不是可量化标准（4 个 iteration 均指出）
12. [Success Criteria] `# user-customized` 有效性验证未在成功标准体现
13. [Scope Definition] 迁移指南是否在范围内未明确
14. [Risk Assessment] probe 差异化与退出码架构不兼容未列为风险
15. [Risk Assessment] `# user-customized` 导致用户错过改进的风险未列出
16. [Logical Consistency] 端口冲突 best-effort 检查误导性信息仍存在

Unannotated region attacks: 16 attack points / ~200 paragraphs = density 0.08

**Ratio (annotated/unannotated)**: 3.75x

**Interpretation**: Annotated regions 的攻击密度从 iteration-3 的 5.5x 降至 3.75x。这表明 pre-revised 区域的修复质量持续提高——iteration-4 的修订引入的新问题更少。无 `conflict-with-pre-revision` 标记。

**跨迭代趋势**:
- iteration-1 → iteration-2: 9.6x（大量新攻击面）
- iteration-2 → iteration-3: 5.5x（修复质量改善）
- iteration-3 → iteration-4: 3.75x（修复质量持续改善，趋近收敛）

---

## Rating

SCORE: 877/1000
DIMENSIONS:
  Problem Definition: 90/110
  Solution Clarity: 106/120
  Industry Benchmarking: 107/120
  Requirements Completeness: 95/110
  Solution Creativity: 74/100
  Feasibility: 84/100
  Scope Definition: 76/80
  Risk Assessment: 84/90
  Success Criteria: 77/80
  Logical Consistency: 84/90
ATTACKS:
1. [Solution Clarity/Logical Consistency]: probe 重试差异化策略声称 run-tests 区分连接失败类型，但 run-tests 只能获得退出码（0/1），无法实现差异化重试 — "run-tests 区分连接失败类型以加速定位——(1) 连接被拒绝（ECONNREFUSED）：端口无进程监听，继续重试；(2) 连接超时：端口有进程但不响应 HTTP，增加计数权重" — 明确差异化重试的实现层级：在 just probe 配方体内部实现（配方体根据 curl 错误类型调整重试计数和权重），或使用特定退出码传递失败类型
2. [Solution Clarity/Scope]: `# user-customized` 全有或全无的粒度——用户一个小修改会阻止所有后续 surface 规则改进的自动应用 — "若目标配方已有此标记，init-justfile 跳过覆盖并输出警告" — 跳过覆盖时输出当前配方与推荐版本的差异摘要，或提供 `--force-regenerate` 选项
3. [Logical Consistency]: YAML 映射插入顺序假设在 Go 的 `map[string]string` 中不成立——Go map 无序 — "优先匹配 surfaces map 中声明顺序靠前的 key（YAML 映射保留插入顺序）" — 验证 Forge 的 YAML 库是否保留插入顺序，或改用有序数据结构（`yaml.Node`）
4. [Feasibility/Requirements]: just >= 1.4.0 版本检查触发时机未定义 — "just >= 1.4.0（[linux]/[windows] recipe attribute 在 just 1.4.0 引入，1.0-1.3.x 不支持此功能，会报 'unknown attribute' 错误）" — 在 init-justfile 前置检查步骤中添加 `just --version` 检查，定义检查失败时的处理策略
5. [Requirements]: run-tests 端到端编排总耗时无上限约束，`test.timeout` 覆盖范围不明确 — "test: timeout: 300" — 明确 timeout 覆盖整个编排序列还是仅覆盖 test 步骤，probe 超时是否独立可配
6. [Problem Definition]: 证据质量从 iteration-3 到 iteration-4 无改善 — "来源于 config-schema.md 中记录的 8 个示例，样本量有限" — 至少提供 Forge 内部示例项目的 surface 类型分布
7. [Success Criteria]: "dry-run 只验证语法"和"'清晰'不可量化"在 4 个迭代中均指出仍未修复 — "所有生成的配方通过 --dry-run 验证（语法和依赖正确性）"和"语言模板与 surface 规则的配方职责边界清晰" — 接受当前状态（第 8 条运行时验证部分补偿 dry-run 缺口），或修改标准为可量化表述
