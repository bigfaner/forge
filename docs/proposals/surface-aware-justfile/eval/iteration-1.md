---
iteration: 1
title: "CTO Adversary Rubric Scoring — Iteration 1"
date: 2026-05-25
scorer: CTO adversary (blind annotated review)
---

# 评分报告：init-justfile Surface 感知 + 测试编排简化

## Phase 1 — 推理链审计

### 论证链路追踪

1. **Problem -> Solution**：两个问题（surface 不感知 + test.execution 冗余）-> Surface 感知配方生成 + 移除委托层。链路成立但存在**捆绑谬误**——两个问题是否必须同时解决？提案在"替代方案"中列出了"仅 surface 感知，保留 test.execution"但标记为"治标不治本"，该判定缺乏论证。保留 test.execution 不影响 surface 感知的核心价值，反而降低了迁移风险。

2. **Solution -> Evidence**：证据在 pre-revised 后有所改善（承认了非 just 示例的存在），但核心量化证据仍然缺失。"绝大多数模板变量最终都解析为 just 命令"——"绝大多数"是多少？70%？95%？没有数据。

3. **Evidence -> Success Criteria**：成功标准以 checklist 覆盖了主要交付物，但缺少性能基准（surface 规则加载耗时的具体数字）和端到端集成测试标准（不仅仅是 dry-run）。

4. **自相矛盾检测**：
   - **矛盾 A（严重）**：提案声称 `test.execution` 在 Go 层面"未被映射为独立字段"，暗示移除成本低。但依赖就绪度章节承认 run-tests SKILL.md 通过 `forge config get test.execution` 让 LLM agent 实际使用这些字段。这种"代码未实现但实际在用"的状态意味着移除影响面比声称的大——不是"无代码变更"，而是需要重写 run-tests 的整个编排逻辑。
   - **矛盾 B**：非功能需求声称 "just >= 1.0（支持 `[linux]`/`[windows]` recipe attribute）"，但 just 的 `[linux]`/`[windows]` recipe attribute 是在 just **1.4.0**（2023 年 3 月）引入的。just 1.0.0 不支持此功能。如果用户安装的是 just 1.0-1.3.x，生成的 justfile 将无法运行，且报错信息不清晰（just 会报告"unknown attribute"）。
   - **矛盾 C**：scope 迁移要求"阶段 1-4 必须在同一提交中完成"，同时又说"阶段 1 的 resolveScope() 保留一个版本的向后兼容逻辑"。如果所有阶段在同一提交中完成，为什么需要兼容层？兼容层的存在本身说明了原子迁移在实践中难以保证，两个声明互相矛盾。

---

## Phase 2 — Rubric Scoring

### 1. Problem Definition (75/110)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Problem stated clearly | 35/40 | 两个问题陈述明确。扣分：问题 1 和问题 2 的独立性和因果关系模糊——"两个相互关联的问题"暗示必须捆绑解决，但未论证为何不能只解决问题 1。 |
| Evidence provided | 15/40 | **严重不足**。"Web UI 的 e2e 测试必须先启动应用，但当前配方没有 surface 特定的启动逻辑"——这是一个逻辑推断，不是证据。没有提供：(1) 任何一个实际项目的 justfile 作为证据；(2) 用户反馈或 bug 报告；(3) 受影响项目的数量或比例。证据章节在 pre-revised 后确实诚实地承认了非 just 示例的存在，但这反而削弱了"委托层冗余"的论证力度——如果非 just 路径确实在使用，委托层就不是纯粹的"冗余转发"。 |
| Urgency justified | 25/30 | 与 v3.0.0 test profile 对齐的时机论证合理。但"cost of delay"分析不完整——如果不做这个提案，v3.0.0 能否正常发布？现有 test.execution 方案能撑多久？ |

### 2. Solution Clarity (90/120)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Approach is concrete | 38/40 | Surface 编排模式表格是亮点——5 种 surface 的编排序列一目了然。test 配方生成 fallback 链设计清晰。 |
| User-facing behavior described | 32/45 | init-justfile 的用户行为描述较好。但：(1) run-tests 的用户行为迁移描述不足——"不再配置 test.execution"对现有用户意味着什么迁移步骤？(2) surface-orchestration.yaml 的字段所有权规则复杂——"工具管理字段无条件覆盖"+"用户可编辑字段保留"的语义对用户是否可理解？(3) 重新生成合并规则有 5 个步骤，init-justfile 是 LLM 执行的 SKILL——LLM 能可靠执行如此复杂的合并逻辑吗？ |
| Technical direction clear | 20/35 | 后台进程管理（后台启动 + PID 追踪 + teardown）在"可行性评估"中有深入讨论，但位于文档后半部分而非方案核心——读者需要读到最后才能理解关键技术可行性。Probe 轮询逻辑的伪代码使用 `$(seq 1 $max_retries)` 是 bash 语法，但提案选择的方案是 just 原生平台 attribute（不依赖 bash）——伪代码与选定方案矛盾。 |

### 3. Industry Benchmarking (85/120)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Industry solutions referenced | 30/40 | 引用了 Docker Compose、Kubernetes、Cypress、Makefile、GitHub Actions 五个方案，且正确识别了 Forge 的关键差异（LLM agent 执行 vs 确定性代码执行）。扣分：对比表缺少**test framework 内建的编排支持**（如 Jest 的 `--runInBand`、Vitest 的 `pool`、Playwright 的 `webServer` 配置）——这些是 Forge 用户最可能遇到的前置方案。 |
| At least 3 meaningful alternatives | 20/30 | 4 个替代方案（不做/仅 surface/surface+去委托/Go 代码管理），覆盖了增量选择。但"仅 surface 感知，保留 test.execution"被标记为"治标不治本"——这是一个潜在的 straw-man 判定，因为保留 test.execution 的 surface 感知方案完全可以作为独立的第一步迭代。"Go 代码直接管理进程生命周期"方案虽然被列为替代，但提案的 LLM agent 确定性缓解措施最终采纳了其核心思想作为兜底——这说明该方案不是真正的"替代"，而是"延期"。 |
| Honest trade-off comparison | 18/25 | justfile 作为唯一抽象层的 trade-off 分析是亮点。但"已知局限"部分列出的两个局限都是可缓解的，没有列出**不可缓解的局限**——例如：just 配方体的跨平台复杂性导致用户难以理解和维护生成的配方。 |
| Chosen approach justified against benchmarks | 17/25 | Forge 方案与行业方案的核心差异定位清晰（LLM agent 执行 + justfile 文本协议），但没有回答：为什么不采用 Cypress `start-server-and-test` 的模式（在 SKILL 内部用 Node.js 脚本管理进程）？该方案同样适用于"启动 -> 等待 -> 测试 -> 清理"场景，且不依赖 just 的后台运行能力。 |

### 4. Requirements Completeness (85/110)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Scenario coverage | 33/40 | 7 个关键场景覆盖了 5 种 surface + 无 surface + 混合。遗漏：(1) `just dev` 启动失败（端口被占、依赖缺失）——虽然 probe 超时会捕获，但 probe 超时是 60 秒的等待，用户体验差；(2) 同类型多 surface 并发（如两个 web surface 使用不同端口）——混合项目场景只讨论了 web+api，未讨论 web+web；(3) init-justfile 多次运行（surface 配置变化后的增量更新/冲突处理）。 |
| Non-functional requirements | 28/40 | pre-revised 后新增了 NFR 表格（跨平台、向后兼容、可观测性、性能、可靠性、just 版本），覆盖面大幅改善。但：(1) "just >= 1.0" 的版本要求有误（实际需要 >= 1.4.0），这会导致兼容性声明不可靠；(2) "性能"只约束了 surface 规则文件加载（不超过 1 秒），未约束 run-tests 编排端到端耗时（如 probe 轮询的 60 秒超时是否可接受？）；(3) "可观测性"要求"结构化日志"但 run-tests 是 LLM agent——LLM 输出不是结构化的，此 NFR 的可实现性存疑。 |
| Constraints & dependencies | 24/30 | Surface 信息源优先级规则清晰。依赖就绪度分析诚实（承认 test.execution 在 agent 层面实际在用）。但 GetConfigValue 扩展作为关键依赖，评估过于简略——"需独立评审"是推诿，不是约束分析。 |

### 5. Solution Creativity (65/100)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Novelty over industry baseline | 25/40 | surface-orchestration.yaml 作为 init-justfile 和 run-tests 的共享契约是一个好的设计——在 SKILL（LLM 执行）的约束下，用文件代替 API 作为接口是务实的创新。但核心编排模式（dev -> probe -> test -> teardown）是标准的测试流水线，无创新。 |
| Cross-domain inspiration | 20/35 | 从 K8s readinessProbe 借鉴探针重试、从 Cypress 借鉴清理机制、从 Docker Compose 借鉴声明式编排——这些借鉴是明确的。但都是从同一领域（服务编排/测试）借鉴，没有跨域灵感。 |
| Simplicity of insight | 20/25 | "justfile 已经是抽象层，config 再包一层只是转发"——这个洞察简洁有力。但实际方案引入了 surface-orchestration.yaml（新的中间层）、字段所有权规则、重新生成合并规则——复杂度并未真正消除，只是从 config.yaml 转移到了另一个 YAML 文件和更复杂的生成逻辑。 |

### 6. Feasibility (70/100)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Technical feasibility | 30/40 | 后台进程管理方案（just 原生平台 attribute + PID 追踪 + 命令行匹配）在 pre-revised 后讨论深度显著提升。Windows 兼容性方案具体（PowerShell Start-Process -PassThru）。扣分：(1) just 版本要求错误（需要 >= 1.4.0 非 1.0.0）；(2) `just probe` 的伪代码使用 bash `$(seq ...)` 但选定方案不依赖 bash——实现路径与方案选择不一致；(3) `just dev` 后台启动后 exit 0 不可靠的问题被坦诚承认，但"依赖 probe 作为唯一检测机制"意味着 dev server 崩溃后用户需等待 60 秒才能收到反馈。 |
| Resource & timeline feasibility | 20/30 | "15-20 个编码任务"的估算范围过大（33% 不确定性）。scope 统一迁移涉及 7 个以上组件，要求"同一提交"完成——这在实践中意味着一个巨型 PR，代码审查困难。config schema 变更作为"独立子方案"在范围内但又需"独立评审"——如果子方案被拒，主提案是否仍然可行？ |
| Dependency readiness | 20/30 | Surface 检测已就位是事实。test.execution 引用审计在 pre-revised 后被补充，但审计清单只列了 4 个 skill——是否还有其他位置引用了 test.execution？（如文档、README、examples 目录？） |

### 7. Scope Definition (65/80)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| In-scope items are concrete | 25/30 | 每个范围内项都是可交付的。5 个 surface 规则文件 + SKILL.md 更新 + config schema 变更 + scope 迁移，粒度合理。 |
| Out-of-scope explicitly listed | 18/25 | 列了 6 项范围外（含 pre-revised 后新增的"回滚基础设施"和"Go 代码子命令"）。但缺少：(1) 现有项目的迁移指南——从 test.execution 过渡到 just 配方的用户操作步骤；(2) CI/CD 环境的适配——在无 just 的 CI runner 上如何工作？ |
| Scope is bounded | 22/25 | "向后兼容：无 surface 配置 -> 当前行为不变"是好的约束。"预计 15-20 个编码任务"提供了时间框架参考。但 config schema 子方案的边界模糊——它在范围内但需独立评审，如果评审不通过怎么办？ |

### 8. Risk Assessment (65/90)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Risks identified | 22/30 | 6 个风险覆盖了主要场景。遗漏：(1) just 版本兼容性风险（>= 1.4.0 要求可能排除了部分用户）；(2) surface-orchestration.yaml 的损坏/手动修改导致 run-tests 行为异常；(3) init-justfile 的 LLM 执行不完整（生成的配方遗漏 surface 特定逻辑）。 |
| Likelihood + impact rated | 21/30 | 大部分评估合理。但 "test [journey] 过滤与原生运行器不兼容" 标记为"中/高"——这是合理的诚实评估。"run-tests 无法感知 surface"标记为"低/高"——如果 run-tests 的 surface 感知依赖 LLM 正确读取 surface-orchestration.yaml，LLM 读取 YAML 的可靠性不是"低"可能性。 |
| Mitigations are actionable | 22/30 | 回滚计划在 pre-revised 后被补充（git revert），基本可操作。但 "LLM 组合语言模板 + surface 规则"不是缓解措施——它是设计本身。HARD-GATE 规则（probe 失败后禁止重试）是一个好的缓解，但它依赖 LLM 遵守 SKILL.md 中的指令——LLM 的指令遵从率不是 100%。 |

### 9. Success Criteria (60/80)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Criteria are measurable and testable | 38/55 | 11 条标准中约 7 条是明确可验证的（生成差异化配方、不生成 run、不再依赖 test.execution.run 等）。但：(1) "所有生成的配方通过 --dry-run 验证"——dry-run 只验证语法不验证运行时行为；(2) "语言模板与 surface 规则的配方职责边界清晰（语言级 vs 编排级），无同名冲突"——"清晰"和"无冲突"如何量化验证？(3) 混合项目端到端验证标准虽然详细，但没有指定测试框架或自动化方式——"端到端验证"是手动操作还是自动化测试？ |
| Coverage is complete | 22/25 | 覆盖了范围内的主要交付物。但缺少：(1) config schema 变更的成功标准（GetConfigValue 扩展是否通过单元测试？残留 test.execution 的处理是否被验证？）；(2) 性能基准的具体数字（"不超过 1 秒"在 NFR 中出现但成功标准未引用）。 |

### 10. Logical Consistency (70/90)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Solution addresses the stated problem | 30/35 | Surface 感知解决"编排流程不同"成立。去掉 test.execution 解决"冗余"也成立。但两者捆绑的必要性论证有 gap——"治标不治本"的判定缺乏深度分析。保留 test.execution + surface 感知完全可以作为 v3.0.0 方案，去掉委托层留到 v3.1.0。 |
| Scope <-> Solution <-> Success Criteria aligned | 22/30 | 范围、方案、标准基本对齐。但对齐缝隙：(1) 方案详细描述了 surface-orchestration.yaml 的字段所有权和重新生成语义，但成功标准中没有"重新生成后用户编辑字段正确保留"的验证标准；(2) 范围包含了 scope 统一迁移，但风险表中没有列出"迁移窗口中遗漏组件"的风险——虽然迁移顺序约束讨论了这个风险，但风险表中的遗漏导致风险评估维度不完整。 |
| Requirements <-> Solution coherent | 18/25 | 下游集成契约表格与方案一致。但"语言模板与 surface 规则的仲裁规则"在方案中定义了（surface 优先覆盖），需求分析中未列出"当语言模板已生成 test 配方时，surface 规则覆盖它的用户影响"作为一个场景——用户自定义的 test 配方会被覆盖吗？边界标记（boundary marker）能保护用户编辑吗？ |

---

## Scoring Summary

| Dimension | Score | Max |
|-----------|-------|-----|
| 1. Problem Definition | 75 | 110 |
| 2. Solution Clarity | 90 | 120 |
| 3. Industry Benchmarking | 85 | 120 |
| 4. Requirements Completeness | 85 | 110 |
| 5. Solution Creativity | 65 | 100 |
| 6. Feasibility | 70 | 100 |
| 7. Scope Definition | 65 | 80 |
| 8. Risk Assessment | 65 | 90 |
| 9. Success Criteria | 60 | 80 |
| 10. Logical Consistency | 70 | 90 |
| **Total** | **730** | **1000** |

---

## Phase 3 — Blindspot Hunt

### [blindspot-1] just 版本要求错误导致兼容性声明不可靠

提案声称 "just >= 1.0" 支持 `[linux]`/`[windows]` recipe attribute。实际上该功能在 just 1.4.0 才引入。如果用户安装了 just 1.0-1.3.x，生成的配方将报 "unknown attribute" 错误。这是一个**阻塞性事实错误**——非功能需求中的兼容性承诺无法兑现。

**引用**："just >= 1.0（支持 `[linux]`/`[windows]` recipe attribute 进行平台分支）"

### [blindspot-2] probe 伪代码与选定平台方案的技术矛盾

Probe 轮询逻辑的伪代码使用 `$(seq 1 $max_retries)` 和 `curl -sf`——这是 bash 语法。但选定的后台进程管理方案是 just 原生平台 attribute（不依赖 bash）。这两个设计选择互相矛盾：如果采用平台 attribute，probe 配方也需要 `[linux]`/`[windows]` 变体，伪代码应展示两个变体而非仅 Unix 版本。

**引用**：选定方案为"使用 just 原生 `[linux]`/`[windows]` recipe attribute 实现平台分支"；但 probe 伪代码仅展示 bash 版本。

### [blindspot-3] surface-orchestration.yaml 的重新生成合并逻辑对 LLM 执行过于复杂

重新生成合并规则包含 5 个步骤（读取现有文件 -> 识别用户编辑字段 -> 生成新结构保留旧值 -> 新增条目写入默认值 -> 移除条目删除）。这些步骤需要 init-justfile（LLM 执行的 SKILL）精确执行——包括：(1) 区分"工具管理字段"和"用户可编辑字段"；(2) YAML 解析和合并；(3) 处理文件损坏的异常路径。LLM 执行如此复杂的文件合并操作的可靠性未被讨论。

**引用**："当 init-justfile 被重新执行时（如用户添加了新的 surface）：1. 读取现有 surface-orchestration.yaml，识别用户编辑过的字段..."

### [blindspot-4] scope 兼容层引入 type->key 映射的多义性

过渡期兼容层规定 `frontend -> 找到 type=web 的 key`，`backend -> 找到 type=api 的 key`。但如果项目有多个 type=web 的 surface（如 `admin-panel: web` + `marketing-site: web`），`frontend` 映射到哪个 key？提案未讨论多义性冲突。

**引用**："若 scope 值为旧枚举（`frontend`/`backend`），在 surfaces map 中查找对应的 surface key 并映射（`frontend` -> 找到 type=web 的 key，`backend` -> 找到 type=api 的 key）"

### [blindspot-5] `just dev` 后台进程的 stdout/stderr 丢失

`just dev` 使用 `nohup npm run dev > /dev/null 2>&1 &` 后台启动 dev server——stdout 和 stderr 被重定向到 /dev/null。如果 dev server 启动失败（如编译错误、依赖缺失），错误信息被静默丢弃。probe 超时后用户只看到"服务启动超时"，看不到失败原因。这对调试体验是灾难性的。

**引用**："nohup npm run dev > /dev/null 2>&1 & echo $! > .forge/dev-server.pid"

### [blindspot-6] `just probe` 在 Windows 上 curl 不可用的 fallback 未设计

提案在"Windows 兼容性"章节承认 `curl` 可能不可用，提出使用 PowerShell `Invoke-WebRequest` 作为 fallback。但 probe 轮询伪代码只展示了 `curl` 版本。如果采用平台 attribute 方案，`[windows]` 变体的 probe 应该使用 PowerShell 实现，但该实现未被展示或设计。跨平台 probe 实现的差异可能导致行为不一致（如 curl 的 `-sf` 静默模式 vs Invoke-WebRequest 的错误处理）。

**引用**："`curl` 可能不可用——`just probe` 需要考虑使用 PowerShell `Invoke-WebRequest` 作为 fallback"

### [blindspot-7] 成功标准缺少 surface-orchestration.yaml 重新生成的验证

方案详细设计了 surface-orchestration.yaml 的字段所有权和重新生成语义（工具管理字段无条件覆盖、用户可编辑字段保留），但成功标准中没有对应的验证项。例如："重新运行 init-justfile 后，用户修改的 probe_target 字段被正确保留"。

### [blindspot-8] api/web 合并为 service 的前瞻性声明缺乏收敛条件

提案声明"若后续验证两者确实无实质性差异，可合并为 service 规则"。但没有定义什么是"验证"——需要多少个迭代？什么数据点？没有收敛条件的声明是空头承诺。

**引用**："若后续验证两者确实无实质性差异，可合并为 `service` 规则并共享编排模板"

---

## Bias Detection Report

**Pre-revised annotated regions**: 9 annotated paragraphs/blocks

Attacks found in annotated regions:
1. [Feasibility] just 版本要求错误 (blindspot-1) — affects pre-revised NFR section (line 368)
2. [Feasibility] probe 伪代码与平台方案矛盾 (blindspot-2) — probe pseudocode (line 464-477) is not pre-revised but the platform attribute selection (line 420) is pre-revised context
3. [Solution Clarity] surface-orchestration.yaml 合并逻辑对 LLM 过于复杂 (blindspot-3) — line 188 pre-revised:high
4. [Logical Consistency] scope 兼容层多义性 (blindspot-4) — line 328 pre-revised:high
5. [Feasibility] stdout/stderr 丢失 (blindspot-5) — dev recipe example (line 428) is not pre-revised itself but part of pre-revised section context
6. [Feasibility] Windows probe fallback (blindspot-6) — line 435 pre-revised:high discusses Windows CMD note
7. [Solution Clarity] test.execution 移除的残留处理 (dimension 2 deduction) — line 605 pre-revised:medium
8. [Feasibility] test.execution 引用审计不完整 (dimension 6 deduction) — line 549 pre-revised:medium

Annotated region attacks: 6 attack points / 9 annotated paragraphs = density 0.67

Unannotated regions: ~183 paragraphs

Attacks in unannotated regions:
1. [Problem Definition] 缺乏量化证据
2. [Problem Definition] 两个问题捆绑的必要性
3. [Solution Clarity] run-tests 用户行为迁移不足
4. [Industry Benchmarking] test framework 内建编排未引用
5. [Industry Benchmarking] "仅 surface 感知"不是 straw-man 论证不足
6. [Requirements] 同类型多 surface 并发
7. [Requirements] 可观测性 NFR 对 LLM agent 不可实现
8. [Risk Assessment] LLM 指令遵从率不是 100%
9. [Success Criteria] dry-run 不验证运行时
10. [Logical Consistency] 用户自定义 test 配方被覆盖的风险
11. [blindspot-5] stdout/stderr 丢失
12. [blindspot-7] 重新生成验证缺失
13. [blindspot-8] api/web 合并收敛条件缺失

Unannotated region attacks: 13 attack points / ~183 paragraphs = density 0.07

**Ratio (annotated/unannotated)**: 9.6x

**Interpretation**: Annotated regions have significantly higher attack density. This is expected — pre-revised sections addressed the original iteration-0 weaknesses but introduced new surface area for attack. The high ratio suggests the pre-revision improved completeness but at the cost of introducing new verifiable claims that can be challenged. Two attacks are tagged `conflict-with-pre-revision`:

- `conflict-with-pre-revision-1`: [Feasibility] just version requirement error — pre-revision added the NFR table with "just >= 1.0", but this factual error was introduced by the revision itself.
- `conflict-with-pre-revision-2`: [Logical Consistency] scope compatibility layer ambiguity — pre-revision added the migration ordering section with compatibility layer, but introduced the multi-surface mapping ambiguity.

---

## Rating

**730/1000 — 及格线以上，存在可修复的结构性问题**

核心优势：
1. Surface 编排模式表格清晰、可操作
2. 后台进程管理的跨平台方案具体（just 原生平台 attribute + PowerShell）
3. 行业对标引用了 5 个成熟方案并正确定位了 Forge 的差异
4. LLM agent 执行确定性的分层防御设计务实

需要修复的问题（按优先级排序）：
1. **事实错误**：just 版本要求应为 >= 1.4.0 非 1.0.0
2. **技术矛盾**：probe 伪代码使用 bash 语法但选定方案不依赖 bash
3. **复杂度失控**：surface-orchestration.yaml 的重新生成合并逻辑对 LLM 执行过于复杂
4. **多义性**：scope 兼容层的 type->key 映射在多 surface 场景下有歧义
5. **调试体验**：`just dev` 后台启动时 stdout/stderr 被丢弃
6. **验证缺口**：成功标准缺少重新生成和 config schema 变更的验证项
