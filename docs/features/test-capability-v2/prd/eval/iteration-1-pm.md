# PRD Evaluation Report — Iteration 1 (PM)

**Feature**: test-capability-v2
**Date**: 2026-05-23
**Evaluator**: Senior PM (independent re-evaluation)
**Scoring Mode**: Mode B (No UI — prd-ui-functions.md absent)
**Total**: 913 / 1000

---

## Dimension 1: Background & Goals — 94 / 100

### Background has three elements (Reason/Target/Users) — 28 / 30

**Reason** (Why): 三大结构性缺陷清晰——双路径并存、测试深度不足、通用性有限。每个有具体的现状描述。v3.0.0 作为重构窗口的理由明确。

**Target** (What): 五条升级主线（管线统一、深度增强、通用扩展、评测补全、信息增强）结构清晰。定位声明 "管线只生成开发者手动编写成本高的复杂测试" 明确界定了功能边界。

**Users** (Who): Forge 用户（项目开发者）和 Forge 维护者，各有明确的使用场景描述。

**Deduction (-2)**: Target 部分五条升级方向中，"评测补全" 和 "信息增强" 在 Background 的三大缺陷中没有被定位为独立缺陷——它们更接近解决方案的一部分。Background 定义了三个问题但 Target 列了五个方向，问题与解之间的映射存在缝隙。

### Goals are quantified — 28 / 30

六个 Goals 中五个有量化指标：

| Goal | Metric | 可验证性 |
|------|--------|---------|
| 消除双路径困惑 | gen-test-cases 完全删除 | 清晰、二元可判定 |
| 提升测试深度 | 高风险测试数 ≥ 8，且 ≥ 低风险 × 1.5 | 量化、可测量 |
| 提升测试信息质量 | Fact Table 覆盖率提升 ≥ 20pp | 量化，公式在 Other Notes 定义 |
| 提升通用性 | ≥ 3 个新 Convention 文件 | 清晰、可计数 |
| 建立评测门禁 | eval 评分 ≥ 850/1000（gold standard 校准） | 量化 |
| 降低 Mobile 接入成本 | Maestro YAML 骨架 + deep link 测试 | 具体产出物 |

**Deduction (-2)**: "降低 Mobile 接入成本" 的 Metric 是产出物描述而非量化指标。没有衡量 "降低" 了多少成本，也没有与当前状态的对比基线。"尽力而为" 的定性定位进一步削弱了可验证性。

### Background and goals are logically consistent — 38 / 40

- 双路径问题 → 退休旧路径：直接对应
- 深度不足 → 风险驱动 + 边界衍生：因果关系清晰
- 通用性有限 → Convention 扩充 + test-guide：直接解决
- 评测补全 → eval-journey/eval-contract 评分指标：可追溯
- 信息增强 → Run-to-Learn → Fact Table 覆盖率：因果链完整

**Deduction (-2)**: Background 明确定位 "管线只生成开发者手动编写成本高的复杂测试"，但 Goals 中没有指标验证这个定位是否在实现中被保持。如果管线实际开始生成简单测试或仍然遗漏复杂场景，没有度量手段来发现偏离。

---

## Dimension 2: Flow Diagrams — 144 / 150

### Mermaid diagram exists — 50 / 50

一个大型 Mermaid flowchart 存在，包含 START/END 节点、决策菱形、处理矩形，使用中文标签，结构完整。

### Main path complete (start → end) — 48 / 50

Happy path 完整覆盖：START → SCENE_DETECT → DETECT → GEN_JOURNEY → EVAL_J → GEN_CONTRACT → EVAL_C → GEN_SCRIPTS → R2L_CHOICE → ENV_CHECK → CONFIDENCE → RUN_TESTS → REPORT → END。四个阶段全部体现。

**Deduction (-2)**: 阶段一步骤 1 "用户在项目中运行测试生成技能" 在图中合并到 START 节点，没有体现用户触发动作的入口。图中 SCENE_DETECT 同时处理了场景检测，导致两个独立步骤（用户触发 + 场景检测）被压缩为一个节点。

### Decision points + error branches covered — 46 / 50

决策点覆盖充分：SCENE_DETECT（未知/混合类型）、DETECT（Convention 存在）、EVAL_J/C（是否达标）、R2L_CHOICE（是否启用 Run-to-Learn）、FIX_DECIDE（自动修复）。

错误分支包括：SCENE_FAIL、REVISE_J/C（迭代修正）、PAUSE_J/C（3 轮耗尽）、ENV_FAIL（环境不就绪）。

**Deduction (-4)**: Flow Description 文字中明确描述了 "eval 评分因 LLM 输出无法解析而失败" 的错误路径（"记录错误日志并重试评分一次；重试仍失败则跳过门禁，标记该 Journey/Contract 为 eval-skipped"），但 Mermaid 图中 EVAL_J 和 EVAL_C 节点只有 "≥ 阈值" 和 "< 阈值" 两个出口，完全没有体现 LLM 解析失败这个第三分支。这是文档内部不一致——文字描述的错误路径在图中缺失。

---

## Dimension 3: Flow Completeness — 177 / 200

### Flow steps describe complete business process — 65 / 70

四个阶段（管线准备、Journey-Contract 生成、测试生成与增强、执行与报告）覆盖完整业务流程。每阶段有步骤编号、输入输出、状态转换。前置条件明确（PRD 存在）。

场景类型检测规则表（10 种信号组合）提供了具体的检测逻辑。风险分级判定规则（High/Medium/Low）定义了分类标准。

**Deduction (-5)**: Run-to-Learn 的终止条件描述为 "≤ 3 轮或覆盖率达标"，但 "达标" 的具体阈值没有在 Flow 中就地定义。读者需要跨 section 查找 Goals（"提升 ≥ 20 个百分点"）和 Other Notes（覆盖率公式）才能确定停止条件。关键终止条件应该就地明确。

### Data flow documented — 60 / 70

Core Concepts 定义了关键数据结构（Journey、Step、Outcome、骨架测试、Fact Table）。Other Notes 定义了 Fact Table 的字段结构和覆盖率计算公式。Convention Schema 必需 Section 定义在 Functional Specs 中。

**Deduction (-10)**: 虽然数据结构定义清晰，但缺少技能间数据传递的显式说明。例如：gen-journeys 输出的 Journey 文档格式是什么？gen-contracts 期望接收什么格式？这些数据契约对于实现者至关重要。管线涉及 8 个模块之间的数据交接（Functional Specs 表格中的 8 个变更点），但没有数据流图或数据契约说明来描述模块间的接口。

### Exception handling and edge cases covered — 52 / 60

Run-to-Learn Failure Handling 表格覆盖了 4 种失败场景（编译失败、运行时崩溃、脏数据输出、API 写操作副作用），每种有检测方式和处理策略。eval 评分失败有处理（重试 → eval-skipped）。场景检测失败有处理。自动修复失败有处理（2 次上限）。兜底原则明确。

**Deduction (-4)**: PRD 不存在的错误在 Flow 中提到 "管线在步骤 1 报错并提示用户先完成 PRD 编写"，但没有被纳入 Mermaid 图、Pipeline Exit Codes 表格或专门的错误处理表中。与其他错误路径的结构化描述（如 Run-to-Learn Failure Handling 表格）相比，这个处理不够完整。

**Deduction (-4)**: test-guide 用户拒绝 Convention 草稿后的重试流程在 User Story 5 中明确描述（"基于用户反馈重新生成草稿，最多重试 2 次"），但在 Flow Description 中完全没有提及。Flow 中 TEST_GUIDE 只有 "用户审核确认" → GEN_JOURNEY 这一条路径，缺少拒绝重试的分支。

---

## Dimension 4: User Stories — 188 / 200

### Coverage: one story per target user — 48 / 50

Forge 用户：5 个 story（1/2/3/5/6），Forge 维护者：2 个 story（4/7）。两种角色都有覆盖。

**Deduction (-2)**: Background 中维护者需求为 "清晰的管线架构和可扩展的场景类型系统"。Story 7 覆盖了可扩展性，Story 4 覆盖了质量保障，但 "清晰的管线架构" 这个需求没有独立的 story。虽然管线统一（Story 1）间接服务于此，但 Story 1 是面向用户的 story，不是维护者视角的架构清晰性需求。

### Format correct (As a / I want / So that) — 50 / 50

所有 7 个 story 严格遵循 As a / I want / So that 格式。I want 部分描述具体行为（如 "管线根据功能的风险等级自动生成不同密度的测试矩阵"），而非模糊动词。

### AC per story (Given/When/Then) — 50 / 50

所有 story 都有 Given/When/Then 格式的 AC。部分 story 有多个 Given/When/Then 块覆盖不同子场景（Story 3: CLI + Mobile；Story 4: eval-journey + eval-contract；Story 5: 3 个场景；Story 6: 4 个条件）。格式规范且完整。

### AC verifiability & boundary coverage — 42 / 50

大部分 AC 量化可验证：
- Story 1: 文件路径列表 + 全局搜索无匹配 — 高度可验证
- Story 2: 3-5 Outcome, ≥ 1.5× — 量化可验证
- Story 3: Contract 占比 ≥ 80%, Maestro YAML — 可验证
- Story 4: 评分 ≥ 850/1000 — 可验证
- Story 7: 自动识别 + 回归验证 — 可验证

**Deduction (-4)**: Story 2 的 "3-5 个 Outcome（含必须衍生的边界 Outcome）" 范围过宽（3 到 5 差距 67%），且 "含必须衍生的边界 Outcome" 存在歧义——这 3-5 个是否包含必须衍生的？如果 3 个 Outcome 中只有 1 个是边界 Outcome，是否满足 "含必须衍生的边界 Outcome"？判定标准不明确。

**Deduction (-4)**: Story 5 中 "用户审核修改量 ≤ 20%" 标记为 **human-verified**，Story 6 中 "边界/异常 Outcome 占比 ≥ 30%" 也标记为 **human-verified**。这些关键质量指标无法自动化验证，作为 Acceptance Criteria 的可执行性受限。**human-verified** 标签虽然诚实地承认了限制，但没有提供任何半自动化验证方法来降低主观判断的方差。

---

## Dimension 5: Scenario Completeness — 129 / 150

### End-to-end scenario coverage — 55 / 60

Flow Description 四个阶段覆盖完整端到端流程。Per-Scenario Strategy 表覆盖 5 种场景类型。风险驱动密度覆盖 3 个等级。Run-to-Learn 覆盖完整迭代周期。

**Deduction (-5)**: 缺少 "全新项目首次接入" 的完整端到端场景。虽然 Story 5 涉及 Convention 自动生成，但 Flow 中 test-guide 只在 Convention 不存在时触发，没有描述首次接入的完整用户旅程——从空项目状态到首次测试执行的完整体验路径。对于 "降低接入成本" 这个目标，首次体验是关键场景。

### Implicit assumptions surfaced — 32 / 40

前置条件明确（PRD 存在）。依赖关系清晰（场景检测 → Convention → Journey → Contract → Scripts）。Run-to-Learn 有明确的执行环境依赖。

**Deduction (-4)**: 所有场景类型检测都依赖项目源码文件的信号（如 main.go、package.json）。对于非主流项目结构（如 Bazel monorepo、pnpm workspace、Nix flake），10 种信号组合可能无法匹配。文档只说 "无法匹配 → 暂停管线"，但没有讨论这种情况的预期频率，也没有提供自定义信号规则的机制（仅在 Story 7 的可扩展性中隐含覆盖）。

**Deduction (-4)**: eval 评分解析失败时，Pipeline Exit Codes 表将其归类为退出码 2（blocking — rubric 配置错误），但更常见的原因可能是 LLM 输出格式不稳定。这个假设将问题归因于静态配置错误而非 LLM 不确定性，可能导致用户误诊——修改 rubric 配置无法解决 LLM 输出格式不稳定的问题。

### Business-rules consistency — 42 / 50

- BIZ-error-reporting-001/002: Pipeline Exit Codes 表遵循 0/1/2 语义，有明确的终止点、退出码和语义说明。
- BIZ-task-lifecycle-003: 删除 test.graduate 和 test.gen-cases 任务类型——这些是系统保留类型，文档计划从注册中移除，符合规则。
- BIZ-quality-gate-001: Scope 中列有 "质量门禁更新以反映新管线"。

**Deduction (-4)**: BIZ-error-reporting-002 要求每个错误信息包含 "具体的失败原因" 和 "恢复提示"。Pipeline Exit Codes 表定义了退出码语义，但 Flow 中部分错误路径的描述不够具体——PAUSE_J/C 的 "暂停管线，输出评分+明细，用户决定" 没有保证输出内容包含恢复提示。与 ENV_FAIL 的 "输出缺失项+修复建议" 相比，标准不统一。

**Deduction (-4)**: "质量门禁更新以反映新管线" 在 In Scope 中列出但无细节。BIZ-quality-gate-001 定义了 compile → fmt → lint → unit/integration tests → e2e regression 的多阶段管线。新测试能力（Journey-Contract 管线）如何融入这个已有质量门禁？如果质量门禁只运行 unit/integration/e2e 而不运行 Journey-Contract 生成的测试，就存在覆盖缺口。

---

## Dimension 6: Edge Case Coverage — 88 / 100

### Error paths documented — 36 / 40

主要错误路径覆盖充分：场景检测失败（未知/混合类型）、Convention 缺失、eval 不达标（自动迭代 + 3 轮耗尽）、eval 评分解析失败（重试 + eval-skipped）、环境不就绪、测试失败（区分脚本 vs Contract 问题，回退到不同阶段）、Run-to-Learn 4 种失败场景（编译失败、运行时崩溃、脏数据、API 写操作副作用）。

**Deduction (-4)**: PRD 缺失错误在 Flow 中提到（"管线在步骤 1 报错并提示"）但没有纳入 Pipeline Exit Codes 表格，也没有像 Run-to-Learn Failure Handling 那样有结构化的表格描述。PRD 不存在是管线启动的前提条件，应作为正式的错误路径文档化。

### Boundary conditions covered — 30 / 35

明确的边界条件：eval 迭代 3 轮上限、自动修复 2 次上限、Convention 重试 2 次、Run-to-Learn 3 轮迭代、风险密度 3 级分类和 Outcome 数量范围、Per-Scenario Strategy 中 5 种场景类型的差异化参数。

**Deduction (-5)**: 缺少并发生成边界条件——多个功能同时运行测试生成时管线状态如何管理？文件锁、竞态条件、共享 Convention 缓存的影响都没有讨论。也缺少大型项目性能边界——数百个 Journey 场景下的执行时间和资源消耗。

### Failure recovery described — 22 / 25

大部分失败场景有恢复路径：eval 失败 → 自动迭代修正；Run-to-Learn 失败 → 使用静态信息继续并降低置信度；环境失败 → 用户修复后重检测；测试失败 → 区分类型回退到不同阶段。Run-to-Learn 兜底原则明确。

**Deduction (-3)**: eval-skipped 标记的恢复路径不完整。LLM 评分解析失败后标记为 eval-skipped，文档说 "由用户手动审核"，但没有说明：(a) 用户审核的具体内容是什么？(b) 审核通过后的操作是什么？(c) 审核不通过是修改 rubric 还是修改 Journey/Contract？缺少审核后的操作流程削弱了门禁机制的保障价值。

---

## Dimension 7: Scope Clarity — 93 / 100

### In-scope items are concrete deliverables — 33 / 35

每个 in-scope 条目指向具体的模块、功能或文件。Eval Rubric 维度框架以子表格展开增加了具体性。Convention Schema 必需 Section 定义在 Functional Specs 中提供了明确的交付标准。

**Deduction (-2)**: "质量门禁更新以反映新管线" 不是一个具体可交付物。"更新" 是动词而非产出物——更新什么文件？增加什么步骤？达到什么效果？与其他条目（如 "新增 eval-journey 评测技能（含 rubric）"）的具体性相比，这一条过于模糊。

### Out-of-scope explicitly lists deferred items — 28 / 30

11 个 out-of-scope 条目，每个有明确描述。两个条目附带注解解释了边界：
- "合约 6 维度模型（schema）修改" 注：required_outcomes 是实例数据不属于 schema 变更
- "gen-test-scripts 的编译/lint 执行器核心逻辑变更" 注：Maestro YAML 属于场景差异化适配

**Deduction (-2)**: "已使用 gen-test-cases 项目的迁移工具" 被列为 Out of Scope 但没有提供不做迁移的理由。对于已使用 gen-test-cases 的用户，升级 v3.0.0 意味着旧路径完全删除且无迁移路径。作为一个破坏性变更，至少应有一句话说明替代方案或建议。

### Scope consistent with functional specs and user stories — 32 / 35

In Scope 条目与 User Stories 基本一一对应：
- 退休 gen-test-cases → Story 1
- eval-journey/eval-contract → Story 4
- 边界衍生 + 风险驱动 → Story 2
- 场景差异化 → Story 3
- Convention + test-guide → Story 5
- Run-to-Learn → Story 6
- 可扩展场景类型 → Story 7

Functional Specs 中 8 个模块变更覆盖了 Scope 条目。

**Deduction (-3)**: Functional Specs #8 "run-tasks: 清理 test.graduate 引用" 是具体代码变更但没有对应的 User Story。这是一个跨切变更，与 In Scope 中 "删除 test.graduate 任务类型和相关任务文件" 应更明确关联。此外，"置信度评级系统" 和 "场景特定执行环境就绪检测" 在 In Scope 中列出，它们的行为在 Flow Description 中描述但分散在多个 Story 的 AC 中（Story 6 的置信度、Story 3 的场景策略），没有独立的 Story 来统一验收。

---

## Blindspot Attacks

### [blindspot-1] Eval 门禁一致性缺口 — Story 4 AC 与实际门禁逻辑不匹配

**Quote**: Story 4 AC: "评分 ≥ 850/1000 则通过" vs Other Notes Eval Gate Calibration: "每维度最低阈值（完整性 ≥ 120/200、语义纯度 ≥ 120/200、前置条件互斥性 ≥ 90/150...）"

**Issue**: 一个 Journey/Contract 可以总分 ≥ 850 但某个维度低于最低阈值（例如完整性 100/200 + 其余维度满分 150×5=750 = 总分 850）。Story 4 的 AC 只检查总分，不检查维度阈值。这意味着 AC 声称通过的场景在实际门禁逻辑中会被拒绝。这是文档内部的逻辑矛盾——两个 section 对同一事件定义了不同的通过条件。实现者如果只看 Story 4 的 AC 会遗漏维度检查。

**What must improve**: Story 4 的 AC 必须明确包含维度阈值检查条件，或说明 "总分 ≥ 850 且所有维度均达最低阈值" 才视为通过。

### [blindspot-2] "不得降低断言严格度" 不可验证

**Quote**: Flow Description 步骤 14: "修复后的测试不得降低断言严格度（如移除断言、放宽阈值）"

**Issue**: 这是一个明确的约束但没有描述检测机制。如何判断 "移除断言" 还是 "删除了不必要的断言"？如何量化 "放宽阈值"？如果无法自动检测，这个约束就是不可执行的。Reasoning audit 独立标记了这个问题——一个不可验证的约束等于没有约束。

**What must improve**: 需要描述具体的检测策略（如 diff 比较断言数量和匹配器类型），或明确这是一个由代码审查保障的非自动化约束并记录在 Risk Mitigation 中。

### [blindspot-3] 交付门禁定义模糊

**Quote**: Delivery Phasing 阶段 1: "门禁：2+ 个已有项目跑完整管线无报错"

**Issue**: "已有项目" 的特征未定义（什么场景类型？多少 Journey？什么语言？）、"完整管线" 的范围未定义（包含 Run-to-Learn 吗？包含自动修复吗？）、"无报错" 的定义不明确（所有测试通过？还是管线不崩溃？）。一个门禁标准如果定义模糊就无法客观判定是否通过。相比之下，阶段 2 和 3 的门禁有量化指标（"高风险 ≥ 低风险 × 1.5"、"≥ 3 个可执行测试"），阶段 1 的门禁标准明显薄弱。

**What must improve**: 交付门禁应使用与 Goals 相同的量化标准——具体的项目特征、管线范围和通过标准。

### [blindspot-4] 风险密度数字范围重叠导致等级区分可能失效

**Quote**: Per-Scenario Strategy: High "总测试数估算 10-20"、Medium "7-13"、Low "4-8"

**Issue**: High 的下限（10）与 Medium 的上限（13）有重合（10-13），Medium 的下限（7）与 Low 的上限（8）有重合（7-8）。一个有 10 个测试的 Journey 可能是 High 也可能是 Medium；一个有 7-8 个测试的 Journey 可能是 Medium 也可能是 Low。Goal 中的 "高风险 ≥ 低风险 × 1.5" 在极端情况下可能不被满足（High=10, Low=8 → 10/8=1.25 < 1.5），导致 Goal metric 与 Strategy range 不一致。Reasoning audit 独立标记了这个数字范围问题。

**What must improve**: 三个等级的范围应不重叠（如 High ≥ 12, Medium 8-11, Low ≤ 7），或明确说明范围是估算值而非硬约束，只有 Goal 中的 ≥ 1.5× 比率是硬性指标。

---

## Score Summary

| Dimension | Score | Max |
|-----------|-------|-----|
| 1. Background & Goals | 94 | 100 |
| 2. Flow Diagrams | 144 | 150 |
| 3. Flow Completeness | 177 | 200 |
| 4. User Stories | 188 | 200 |
| 5. Scenario Completeness | 129 | 150 |
| 6. Edge Case Coverage | 88 | 100 |
| 7. Scope Clarity | 93 | 100 |
| **Total** | **913** | **1000** |
