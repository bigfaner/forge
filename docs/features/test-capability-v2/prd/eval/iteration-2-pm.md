# PRD Evaluation Report — Iteration 2 (PM)

**Feature**: test-capability-v2
**Date**: 2026-05-23
**Evaluator**: Senior PM (adversarial re-evaluation)
**Scoring Mode**: Mode B (No UI — prd-ui-functions.md absent)
**Total**: 955 / 1000

---

## Previous-Attack Resolution Status

| Iteration 1 Attack | Status | Evidence |
|---------------------|--------|----------|
| blindspot-1: Eval 门禁一致性（总分 vs 维度阈值） | **Resolved** | Story 4 AC 现在明确包含维度阈值："评分 ≥ 850/1000 且每维度不低于最低阈值（完整性 ≥ 120、语义纯度 ≥ 120、前置条件互斥性 ≥ 90、事实依据 ≥ 90、场景适配 ≥ 90、一致性 ≥ 90）" |
| blindspot-2: 断言严格度不可验证 | **Resolved** | Flow 步骤 14 现在标记为 **human-verified**："此约束由人工在 code review 中验证，无自动化检测机制" |
| blindspot-3: 交付门禁模糊 | **Unresolved** | Delivery Phasing 阶段 1 仍为 "2+ 个已有项目跑完整管线无报错"，未量化 |
| blindspot-4: 风险密度范围重叠 | **Resolved** | Per-Scenario Strategy 表范围不再重叠：High 13-20、Medium 8-12、Low 4-7 |
| D2: eval LLM 解析失败分支缺失 | **Resolved** | Mermaid 图新增 EVAL_J_SKIP / EVAL_C_SKIP 节点 |
| D3: test-guide 拒绝重试分支缺失 | **Resolved** | Mermaid 图新增 TEST_GUIDE 拒绝重试回路 |
| D3: 数据流缺失 | **Partially Resolved** | 新增 Data Flow Table（7 行），但 GEN_SCRIPTS → RUN_TESTS 和 CONFIDENCE 计算输入仍缺失 |
| D3: PRD 不存在未纳入 Exit Codes | **Unresolved** | Pipeline Exit Codes 表仍未包含此错误路径 |
| D6: eval-skipped 恢复不完整 | **Partially Resolved** | 扩展为 4 点策略，但审核标准和审核通过后的具体操作仍有模糊空间 |

**Resolution Summary**: 5/9 完全解决，2/9 部分解决，2/9 未解决。

---

## Dimension 1: Background & Goals — 98 / 100

### Background has three elements (Reason/Target/Users) — 30 / 30

**Reason** (Why): 三大结构性缺陷（双路径并存、测试深度不足、通用性有限）清晰具体，每个有现状描述。v3.0.0 重构窗口理由明确。

**Target** (What): 五条升级主线结构清晰。定位声明 "管线只生成开发者手动编写成本高的复杂测试" 明确界定功能边界。

**Users** (Who): Forge 用户（项目开发者）和 Forge 维护者各有明确场景。

Iteration 1 扣分点（问题与方案映射缝隙）已自然消解——五条升级方向与三个问题之间的映射关系更加明确。

### Goals are quantified — 28 / 30

六个 Goals 中五个有清晰量化指标：

| Goal | Metric | 可验证性 |
|------|--------|---------|
| 消除双路径困惑 | gen-test-cases 完全删除 | 二元可判定 |
| 提升测试深度 | 高风险 ≥ 13 且 ≥ 低风险 × 1.5 | 量化，绝对下限防低基数合规 |
| 提升测试信息质量 | Fact Table 覆盖率提升 ≥ 20pp | 量化，公式在 Other Notes 定义 |
| 提升通用性 | ≥ 3 个新 Convention 文件 | 可计数 |
| 建立评测门禁 | ≥ 850/1000，gold standard 校准 | 量化 |
| 降低 Mobile 接入成本 | Maestro YAML 骨架 + deep link 测试 | 产出物描述 |

**Deduction (-2)**: "降低 Mobile 接入成本" 的 Metric 仍然是产出物描述而非量化指标——没有衡量成本降低了多少，也没有与当前状态对比的基线。Scope 中 Mobile 定位为 "尽力而为" 部分解释了这一点，但作为 Goal 的 Metric 仍应包含量化成分。

### Background and goals are logically consistent — 40 / 40

五条升级方向与三个问题的因果关系清晰：
- 双路径 → 退休旧路径：直接对应
- 深度不足 → 风险驱动 + 边界衍生 + 场景差异化：完整解决方案
- 通用性有限 → Convention 扩充 + test-guide + Mobile Maestro：直接解决

评测补全和信息增强作为测试深度和通用性不足的衍生解决方案，逻辑链完整。定位声明不需要独立指标——它是一个设计约束，由整体架构保障。

---

## Dimension 2: Flow Diagrams — 150 / 150

### Mermaid diagram exists — 50 / 50

大型 Mermaid flowchart 存在，包含 START/END 节点、决策菱形、处理矩形、中文标签。

### Main path complete (start → end) — 50 / 50

Happy path 完整覆盖四个阶段：管线准备 → Journey-Contract 生成 → 测试生成与增强 → 执行与报告。START → SCENE_DETECT → DETECT → GEN_JOURNEY → EVAL_J → GEN_CONTRACT → EVAL_C → GEN_SCRIPTS → R2L_CHOICE → ENV_CHECK → CONFIDENCE → RUN_TESTS → REPORT → END。

### Decision points + error branches covered — 50 / 50

决策点：SCENE_DETECT（未知/混合类型）、DETECT（Convention 存在）、EVAL_J/C（评分是否达标）、R2L_CHOICE（是否启用 Run-to-Learn）、FIX_DECIDE（自动修复类型选择）。

错误分支：SCENE_FAIL、TEST_GUIDE 拒绝重试（≤ 2）、REVISE_J/C（迭代修正）、PAUSE_J/C（3 轮耗尽）、EVAL_J_SKIP/C_SKIP（LLM 解析失败 → eval-skipped）、R2L_DEGRADE（骨架测试失败降级）、ENV_FAIL（环境不就绪）。

Iteration 1 扣分点（eval LLM 解析失败分支缺失、test-guide 拒绝重试缺失）均已修正。

---

## Dimension 3: Flow Completeness — 187 / 200

### Flow steps describe complete business process — 67 / 70

四个阶段覆盖完整业务流程。前置条件明确（PRD 存在且含至少 1 个 User Story + 1 条 AC）。场景类型检测规则表（10 种信号组合）具体。风险分级判定规则（High/Medium/Low）定义清晰。

**Deduction (-3)**: PRD 质量前置检查（"PRD 必须包含至少 1 个 User Story（含 As a / I want / So that 结构），且每个 Story 至少包含 1 条 Acceptance Criteria。若 PRD 不存在或质量前置检查未通过，管线在步骤 1 报错并输出缺失项明细"）是一个重要的业务逻辑分支，但它作为前置条件段落存在，而非编号步骤。这个检查的失败路径（报错并输出缺失项）没有在步骤序列中明确体现，读者可能忽略这个门控点。

### Data flow documented — 65 / 70

新增 Data Flow Table 覆盖了 7 行主要数据传递路径：SCENE_DETECT → 多个消费步骤、TEST_GUIDE → gen-test-scripts、GEN_JOURNEY → gen-contracts + eval-journey、GEN_CONTRACT → gen-test-scripts + eval-contract + Run-to-Learn、R2L → gen-test-scripts、CONFIDENCE → run-tests + 测试报告。每行包含源步骤、产出数据、消费步骤、传递方式。显著改善。

**Deduction (-5)**: Data Flow Table 缺少两个关键路径：
1. GEN_SCRIPTS → RUN_TESTS：gen-test-scripts 生成的测试代码如何传递给 run-tests？通过文件系统路径约定？通过 session 元数据？
2. CONFIDENCE 计算输入：置信度评级的计算依赖 Fact Table 中 confirmed 事实的占比，但 Data Flow Table 中 CONFIDENCE 行只描述了输出传递，没有描述 Fact Table 作为计算输入的来源路径。

### Exception handling and edge cases covered — 55 / 60

Run-to-Learn Failure Handling 表格（4 种失败场景 + 兜底原则）结构化程度高。eval 评分失败处理完整（重试 → eval-skipped → 4 点降级策略）。场景检测失败处理有 Mermaid 分支。自动修复失败处理有 2 次上限 + 类型区分。test-guide 拒绝重试有 Mermaid 回路。

**Deduction (-5)**: PRD 不存在/质量不达标的错误路径在 Flow 文字中有描述（前置条件段落 + 步骤 1 报错），但：
1. 未纳入 Pipeline Exit Codes 表格——其他所有错误路径都有退出码定义
2. 未纳入 Mermaid 图——图中没有 PRD 检查的决策节点
3. 缺少结构化的错误描述表（与 Run-to-Learn Failure Handling 表格相比）

这个错误路径的文档化水平低于其他同类路径。

---

## Dimension 4: User Stories — 194 / 200

### Coverage: one story per target user — 50 / 50

Forge 用户：Story 1/2/3/5/6（5 个）。Forge 维护者：Story 4/7（2 个）。Story 1 间接服务维护者对管线架构清晰性的需求。覆盖充分。

### Format correct (As a / I want / So that) — 50 / 50

所有 7 个 story 严格遵循格式。I want 描述具体行为而非模糊动词。

### AC per story (Given/When/Then) — 50 / 50

所有 story 都有 Given/When/Then 格式的 AC。多场景 story 有多个 Given/When/Then 块。

### AC verifiability & boundary coverage — 44 / 50

**Deduction (-3)**: Story 6 有 2 个 AC 条件标记为 **human-verified**（修改量 ≤ 20%、Outcome 占比 ≥ 30%），这些关键质量指标无法自动化验证。虽然 Story 5 的 ≤ 20% 描述了具体方法（"diff --stat 统计用户修改行数占草稿总行数的比例"），Story 6 的 ≥ 30% 定义了基线（"初始静态侦察时的占比"），但它们仍然依赖人工执行和判断。在一个以自动化管线为核心的产品中，多个 **human-verified** AC 的积累降低了 AC 集合的整体可执行性。

**Deduction (-3)**: Story 2 的 "每个 Step 包含 3-5 个 Outcome" 范围较宽（3 到 5 差距 67%）。虽然增加了退化逻辑（无低风险 variant 时退化为绝对值验证 ≥ 13）和比较验证方法，但 3-5 的范围仍允许较低标准的合规——3 个 Outcome 的 Step 是否真的达到了 "深度测试" 的目标？

---

## Dimension 5: Scenario Completeness — 137 / 150

### End-to-end scenario coverage — 57 / 60

Flow Description 四个阶段覆盖完整端到端流程。Per-Scenario Strategy 表覆盖 5 种场景类型。风险驱动密度覆盖 3 个等级。Run-to-Learn 覆盖完整迭代周期。首次接入通过 test-guide 自动触发隐含覆盖。

**Deduction (-3)**: Delivery Phasing 分三阶段，但缺少跨阶段集成场景的讨论——阶段 1 退休旧路径后、阶段 2 新能力上线前是否存在空窗期？阶段 1 的门禁标准（"完整管线无报错"）中 eval-journey/eval-contract 是否已存在？如果不存在，阶段 1 的 "完整管线" 指的是什么？缺少阶段间衔接说明。

### Implicit assumptions surfaced — 35 / 40

前置条件明确（PRD 存在）。依赖关系清晰。Run-to-Learn 有执行环境依赖和超时保护。

**Deduction (-3)**: Pipeline Exit Codes 表中 "eval 评分解析失败（重试后仍失败）" 退出码 2 的语义描述为 "blocking — rubric 配置错误需人工修复"。但 Flow 步骤 5 的描述是 "eval 评分因 LLM 输出无法解析而失败"。两个位置对同一失败的归因不同——一个归因于 rubric 配置错误（静态问题），一个归因于 LLM 输出不稳定（动态问题）。这会导致用户错误地修改 rubric 而非调整 LLM 参数或重试。

**Deduction (-2)**: 场景检测依赖项目源码文件信号（10 种组合），非主流项目结构可能无法匹配。虽然降级策略合理（暂停管线，用户确认）且 Story 7 提供了长期扩展方案，但近期内的用户体验影响没有评估。

### Business-rules consistency — 45 / 50

BIZ-error-reporting-001/002: Pipeline Exit Codes 表遵循 0/1/2 语义。BIZ-task-lifecycle-003: 删除的系统保留类型与规则一致。

**Deduction (-3)**: "质量门禁更新以反映新管线" 在 In Scope 中描述了 eval-journey → eval-contract → 置信度评级的多阶段门禁。但 BIZ-quality-gate-001 定义了完全不同的多阶段管线（compile → fmt → lint → unit/integration tests → e2e regression）。两个管线之间是什么关系？新测试管线的门禁结果是否汇入 BIZ-quality-gate-001 的流程？缺少集成说明导致两个管线可能各自运行而互不感知。

**Deduction (-2)**: PAUSE_J/C（eval 3 轮耗尽）的输出标准低于 ENV_FAIL。ENV_FAIL 保证 "输出缺失项+修复建议"，PAUSE_J/C 只保证 "输出评分+明细"。用户在 PAUSE_J/C 暂停后需要自行分析评分明细来决定下一步，缺乏直接的行动指导。

---

## Dimension 6: Edge Case Coverage — 93 / 100

### Error paths documented — 38 / 40

错误路径覆盖全面：场景检测失败、Convention 缺失、eval 不达标、eval 评分解析失败（eval-skipped）、环境不就绪、测试失败（区分脚本 vs Contract 问题）、Run-to-Learn 4 种失败场景。

**Deduction (-2)**: PRD 前置检查失败（不存在或质量不达标）未纳入 Pipeline Exit Codes 表格。作为管线启动的前提条件，应有正式的退出码定义（建议：退出码 1，retryable — 用户完成 PRD 后可重跑）。

### Boundary conditions covered — 32 / 35

明确的边界条件：eval 迭代 3 轮上限、自动修复 2 次上限、Convention 重试 2 次、Run-to-Learn 3 轮迭代、风险密度 3 级分类和 Outcome 数量范围（无重叠）、Per-Scenario Strategy 5 种场景类型差异化参数。

**Deduction (-3)**: 缺少并发生成边界条件——多个功能同时运行测试生成时，.forge/session.yaml 的并发写入、Fact Table 的并发读写、测试输出的文件锁冲突等问题没有讨论。对于支持多功能的 Forge 管线，这是一个实际会遇到的边界条件。

### Failure recovery described — 23 / 25

eval-skipped 降级策略已扩展为 4 点：（1）下游正常执行；（2）测试文件头部标记 eval-skipped: true + confidence: LOW；（3）报告中单独列出 eval-skipped 项；（4）用户审核后可手动清除 eval-skipped 标记，清除后置信度由 Fact Table 覆盖率重新计算。显著改善。

**Deduction (-2)**: eval-skipped 第 (4) 点中 "用户审核" 的内容仍不够具体——用户审核 Journey/Contract 的语义正确性？还是审核 LLM 评分输出？审核通过的标准是什么？虽然比 iteration 1 好很多，但审核通过的操作定义仍有模糊空间。

---

## Dimension 7: Scope Clarity — 96 / 100

### In-scope items are concrete deliverables — 34 / 35

每个 in-scope 条目指向具体模块、功能或文件。Eval Rubric 维度框架有子表格展开。Convention Schema 必需 Section 有字段级定义。

"质量门禁更新" 改为更具体的描述："将现有单一门禁替换为多阶段门禁（eval-journey → eval-contract → 置信度评级），每阶段独立 pass/fail 判定，门禁结果汇入统一质量报告"。

**Deduction (-1)**: 描述更具体但仍未指明具体修改的配置文件或代码入口。"统一质量报告" 的格式和位置也未定义。

### Out-of-scope explicitly lists deferred items — 29 / 30

11 个 out-of-scope 条目，每个有明确描述。两个条目附带注解解释边界。

**Deduction (-1)**: "已使用 gen-test-cases 项目的迁移工具" 被列为 Out of Scope 但没有提供替代方案。作为破坏性变更，至少应有一句话说明已有用户的升级路径（如 "用户需手动迁移至 Journey-Contract 路径" 或 "参考 Story 1 AC 中的删除清单进行清理"）。

### Scope consistent with functional specs and user stories — 33 / 35

In Scope 条目与 User Stories 基本一一对应。Functional Specs 8 个模块变更覆盖 Scope 条目。

**Deduction (-2)**: Functional Specs #8 "run-tasks: 清理 test.graduate 引用" 无对应 User Story。"置信度评级系统" 和 "场景特定执行环境就绪检测" 在 In Scope 中列出但分散在 Story 6 和 Story 3 的 AC 中，没有统一的验收入口。虽然可通过交叉覆盖间接验收，但缺少明确的 Story-to-Scope 追溯。

---

## Blindspot Attacks

### [blindspot-1] Delivery Phasing 阶段 1 门禁仍然模糊

**Quote**: Delivery Phasing 阶段 1: "门禁：2+ 个已有项目跑完整管线无报错"

**Issue**: Iteration 1 已指出此问题，当前版本未修正。"已有项目" 的特征未定义（什么场景类型？什么语言？多少 Journey？）、"完整管线" 的范围未定义（包含 Run-to-Learn 吗？包含自动修复吗？eval-journey/eval-contract 在阶段 1 就存在吗？）、"无报错" 的定义不明确（所有测试通过？管线不崩溃？退出码为 0？）。阶段 2 和 3 的门禁有量化指标，阶段 1 的门禁标准与其他两个阶段不统一。这不仅是模糊性——它可能导致阶段 1 在不同评审者之间有不同的通过判定。

**What must improve**: 使用量化标准——具体的项目特征（至少覆盖 CLI + API 两种场景类型）、管线范围（从 gen-journeys 到 run-tests 的完整流程）、通过标准（退出码 0 且无 eval-skipped 标记）。

### [blindspot-2] Pipeline Exit Codes 归因矛盾

**Quote**: Pipeline Exit Codes: "eval 评分解析失败（重试后仍失败）" — "blocking — rubric 配置错误需人工修复" vs Flow 步骤 5: "若 eval 评分因 LLM 输出无法解析而失败，记录错误日志并重试评分一次"

**Issue**: 退出码 2 将失败归因为 "rubric 配置错误"，Flow 将其归因为 "LLM 输出无法解析"。如果真实原因是 LLM 输出格式不稳定（如 token 限制截断、格式偏差），用户修改 rubric 无法解决问题。这个归因矛盾会直接误导用户到错误的修复路径。已在 Dimension 5 (Implicit assumptions) 中扣分，但作为 blindspot 记录以强调其跨 section 影响力。

**What must improve**: Pipeline Exit Codes 表中 "eval 评分解析失败" 的语义应改为更中性的归因（如 "blocking — eval 输出格式异常，检查 rubric 配置或 LLM 输出稳定性"），或拆分为两个退出码（一个配置错误，一个 LLM 输出异常）。

### [blindspot-3] Story 7 缺少新场景类型的评测门禁要求

**Quote**: Story 7 AC: "管线自动识别新场景类型，gen-journeys/gen-contracts/gen-test-scripts 按配置的策略和规则执行，无需修改任何管线技能代码"

**Issue**: Story 7 描述了新增场景类型的接入方式，但 AC 没有要求新场景类型通过 eval-journey/eval-contract 评测门禁。Flow 中 eval 的 rubric 评分维度包含 "场景适配（Scenario Fitness）" 维度——评估是否遵循场景类型的 required_outcomes 规则。但新增场景类型的配置文件（`detect_rules`、`strategy`、`env_check`、`required_outcomes`）中的规则是否需要被 rubric 理解？如果 eval 技能不知道新场景类型的 required_outcomes 列表，它无法在 "场景适配" 维度上正确评分。这意味着新场景类型可能通过了 Story 7 的 AC（能生成测试），但无法通过 eval 门禁（eval 不知道新场景的规则），导致所有新场景类型的 Journey/Contract 自动触发 eval 迭代修正。

**What must improve**: Story 7 应增加一条 AC：新场景类型的配置文件（包含 `required_outcomes`）应被 eval-journey/eval-contract 自动加载，确保 "场景适配" 维度能正确评分。或明确说明 eval rubric 如何动态适配新场景类型。

---

## Score Summary

| Dimension | Score | Max |
|-----------|-------|-----|
| 1. Background & Goals | 98 | 100 |
| 2. Flow Diagrams | 150 | 150 |
| 3. Flow Completeness | 187 | 200 |
| 4. User Stories | 194 | 200 |
| 5. Scenario Completeness | 137 | 150 |
| 6. Edge Case Coverage | 93 | 100 |
| 7. Scope Clarity | 96 | 100 |
| **Total** | **955** | **1000** |
