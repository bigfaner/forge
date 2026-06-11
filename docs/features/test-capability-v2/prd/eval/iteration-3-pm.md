# PRD Evaluation Report — Iteration 3 (PM)

**Feature**: test-capability-v2
**Date**: 2026-05-23
**Evaluator**: Senior PM (adversarial re-evaluation)
**Scoring Mode**: Mode B (No UI — prd-ui-functions.md absent)
**Total**: 940 / 1000

---

## Previous-Attack Resolution Status

### Iteration 2 PM Attacks

| # | Attack | Status | Evidence |
|---|--------|--------|----------|
| 1 | Delivery Phasing 阶段 1 门禁模糊 | **Resolved** | 阶段 1 门禁现包含具体项目特征（"需包含 1 个 CLI + 1 个 API 项目，各含 >= 3 个 PRD User Story"）、管线范围（gen-journeys → eval-journey → gen-contracts → eval-contract → gen-test-scripts → run-tests）、通过标准（"退出码 0 或仅含已知的 manual-only 标记项"） |
| 2 | Pipeline Exit Codes 归因矛盾 | **Resolved** | 退出码 2 语义改为 "blocking — LLM 输出无法解析为结构化评分（可能原因：rubric 配置错误或 LLM 输出格式异常），需人工检查 rubric 配置或 LLM prompt"，消除了单一归因问题 |
| 3 | Story 7 缺少新场景类型的评测门禁要求 | **Resolved** | Story 7 新增第三条 AC："新场景类型的 required_outcomes 配置自动反映到 eval rubric 的场景适配维度评分中...缺少任何一个必须 Outcome 扣 30 分/个；新增场景类型需提交至少 1 个人工标注的 gold standard 文档对用于校准" |

### Iteration 2 PM Blindspots

| # | Blindspot | Status | Evidence |
|---|-----------|--------|----------|
| 1 | Delivery Phasing 阶段 1 门禁模糊 | **Resolved** | 同 Attack 1 |
| 2 | Pipeline Exit Codes 归因矛盾 | **Resolved** | 同 Attack 2 |
| 3 | Story 7 缺少新场景类型的评测门禁要求 | **Resolved** | 同 Attack 3 |

### Iteration 2 Dimension Deductions — Resolution Check

| Dimension | Iteration 2 Deduction | Status | Evidence |
|-----------|----------------------|--------|----------|
| D1: Mobile Goal 非量化 | -2 pts | **Resolved** | Goal 改为量化指标："新 Mobile 项目从零到可运行 Maestro 测试 <= 30 分钟...deep link 测试覆盖 >= 2 个核心 Journey" |
| D3: PRD 前置检查不在步骤序列中 | -3 pts | **Resolved** | Mermaid 图新增 PRD_CHECK 决策节点（START 后第一个节点），Pipeline Exit Codes 表新增 PRD 不存在的退出码 1 |
| D3: Data Flow 缺 GEN_SCRIPTS → RUN_TESTS 和 CONFIDENCE 计算输入 | -5 pts | **Resolved** | Data Flow Table 新增 GEN_SCRIPTS → RUN_TESTS 行（"写入 tests/<journey>/ 目录，按 Convention discovery 规则可被发现"）和 EVAL_J/EVAL_C → revise 行；CONFIDENCE 的计算输入（Fact Table）在 Other Notes 覆盖率公式中明确 |
| D3: PRD 不存在未纳入 Exit Codes | -5 pts | **Resolved** | Pipeline Exit Codes 表新增第一行："PRD 不存在或质量前置检查未通过（缺少 User Story / 缺少 Acceptance Criteria）| 1 | retryable" |
| D4: Story 1 否定验证不可穷举 | -3 pts | **Resolved** | Story 1 AC 现包含具体删除清单（技能目录、评测命令、Rubric 文件及子 rubric、任务类型、任务文件）+ 3 条正向验证 AC |
| D4: Story 2 比较基数（3-5 Outcome 范围过宽） | -3 pts | **Resolved** | Story 2 AC 新增退化逻辑："若同一功能仅有高风险 Journey（无低风险 variant 可比较），退化为绝对值验证：高风险 Journey 总 Outcome 数 >= 13" |
| D5: 跨阶段集成场景缺失 | -3 pts | **Resolved** | Delivery Phasing 阶段 1 门禁明确管线范围包含 eval-journey/eval-contract |
| D5: Exit Codes 对同一失败归因矛盾 | -3 pts | **Resolved** | 退出码 2 语义已修正 |
| D5: BIZ-quality-gate-001 集成关系未说明 | -3 pts | **Resolved** | In Scope 新增"与 BIZ-quality-gate-001 的集成关系"段落（4 点说明） |
| D5: PAUSE_J/C 缺行动指导 | -2 pts | **Resolved** | 新增 PAUSE_J/PAUSE_C 恢复路径（3 选项：跳过门禁继续 / 放弃管线 / 修改后重跑） |
| D6: PRD 前置检查失败无退出码 | -2 pts | **Resolved** | Pipeline Exit Codes 表新增 PRD 前置检查失败退出码 |
| D6: eval-skipped 审核内容模糊 | -2 pts | **Resolved** | eval-skipped 降级策略第 (3) 点明确："提示用户人工审核 Journey/Contract 内容正确性" |
| D7: 质量门禁更新缺少具体文件 | -1 pt | **Partially Resolved** | 描述更具体但仍未指明具体配置文件路径 |
| D7: gen-test-cases 迁移工具缺替代方案 | -1 pt | **Unresolved** | 仍无替代方案说明 |

**Resolution Summary**: 14/16 完全解决，1/16 部分解决，1/16 未解决。

---

## Dimension 1: Background & Goals — 98 / 100

### Background has three elements (Reason/Target/Users) — 30 / 30

**Reason** (Why): 三大结构性缺陷（双路径并存、测试深度不足、通用性有限）清晰具体。v3.0.0 重构窗口理由明确。

**Target** (What): 五条升级主线结构清晰。定位声明"管线只生成开发者手动编写成本高的复杂测试"明确界定功能边界。

**Users** (Who): Forge 用户（项目开发者）和 Forge 维护者各有明确场景和需求描述。

### Goals are quantified — 28 / 30

六个 Goals 全部有清晰指标：

| Goal | Metric | 可验证性 |
|------|--------|---------|
| 消除双路径困惑 | gen-test-cases 完全删除 | 二元可判定 |
| 提升测试深度 | 高风险 >= 13 且 >= 低风险 x 1.5 | 量化，绝对下限防低基数合规 |
| 提升测试信息质量 | Fact Table 覆盖率提升 >= 20pp | 量化，公式在 Other Notes 定义 |
| 提升通用性 | >= 3 个新 Convention 文件 | 可计数 |
| 建立评测门禁 | >= 850/1000，gold standard 校准 | 量化 |
| 降低 Mobile 接入成本 | <= 30 分钟 + deep link 覆盖 >= 2 个核心 Journey | 量化 |

**Deduction (-2)**: "降低 Mobile 接入成本" 的 "<= 30 分钟" 依赖人工计时，缺少对"从零"起点的操作性定义。"30 分钟"包含哪些操作？安装 Maestro？编写 PRD？还是仅从运行命令到拿到测试结果？对照度量说明了"接入耗时（分钟）+ 自动生成 Journey 覆盖数"，但"接入耗时"的计时区间未定义，不同评审者可能计时不一致。

### Background and goals are logically consistent — 40 / 40

五条升级方向与三个问题的因果关系清晰且完整。定位声明作为设计约束由整体架构保障。BIZ-quality-gate-001 集成关系段落消除了两个管线之间的模糊地带。

---

## Dimension 2: Flow Diagrams — 150 / 150

### Mermaid diagram exists — 50 / 50

大型 Mermaid flowchart 存在，包含 START/END 节点、决策菱形、处理矩形。PRD_CHECK 作为新节点补充了前置检查。

### Main path complete (start → end) — 50 / 50

Happy path 完整覆盖四个阶段：PRD_CHECK → SCENE_DETECT → DETECT → GEN_JOURNEY → EVAL_J → GEN_CONTRACT → EVAL_C → GEN_SCRIPTS → R2L_CHOICE → ENV_CHECK → CONFIDENCE → RUN_TESTS → REPORT → END。

### Decision points + error branches covered — 50 / 50

决策点：PRD_CHECK（PRD 存在且质量通过）、SCENE_DETECT（未知/混合类型）、DETECT（Convention 存在）、EVAL_J/C（评分达标）、R2L_CHOICE（启用 Run-to-Learn）、FIX_DECIDE（自动修复类型选择，区分脚本问题 vs Contract 语义错误）。

错误分支：PRD_FAIL（报错退出码 1）、SCENE_FAIL、TEST_GUIDE 拒绝重试、REVISE_J/C（迭代修正）、PAUSE_J/C（3 轮耗尽 + 3 条恢复路径）、EVAL_J_SKIP/C_SKIP（LLM 解析失败）、R2L_DEGRADE（骨架测试降级）、ENV_FAIL、FIX_DECIDE 双回路线（GEN_SCRIPTS / GEN_CONTRACT）。

Iteration 2 的所有扣分点（PRD 检查缺失、eval LLM 解析失败缺失、FIX_DECIDE 回路粒度）均已修正。

---

## Dimension 3: Flow Completeness — 182 / 200

### Flow steps describe complete business process — 65 / 70

四个阶段覆盖完整业务流程。前置条件已升级为 Mermaid 图中的 PRD_CHECK 节点 + 文字描述。场景类型检测规则表格扩展至 16 种信号组合（含 Python/Java/Rust 生态）。风险分级判定规则清晰。PAUSE_J/PAUSE_C 恢复路径（3 选项）补充完整。

**Deduction (-5)**: Flow Description 步骤 6-9 存在结构性重复。步骤 6 "gen-contracts 从 Journey 生成 6 维度合约规范，自动衍生边界/异常 Outcome" 和步骤 8 "gen-contracts 从 Journey 生成 6 维度合约规范，自动衍生边界 Outcome。合约生成后执行 schema 验证..." 描述了同一个步骤但内容不同（步骤 6 无 schema 验证，步骤 8 有）。同样步骤 7 和步骤 9 都是 "eval-contract 评测"。PAUSE 恢复路径也重复出现了两次（步骤 7 后和步骤 9 后各一次）。读者无法确定哪一组步骤是权威描述——是步骤 6-7（简版）还是步骤 8-9（含 schema 验证的完整版）？这种重复不是"冗余补充"，而是造成了步骤序列的歧义。

### Data flow documented — 68 / 70

Data Flow Table 覆盖 9 行关键数据传递路径，新增了 EVAL_J/EVAL_C → revise（评分结果）、GEN_SCRIPTS → RUN_TESTS（测试代码文件路径）两行。每行包含源步骤、产出数据、消费步骤、传递方式。

**Deduction (-2)**: Data Flow Table 中 GEN_CONTRACT 产出 "静态 Fact Table（代码侦察结果）"，但 Fact Table 的产出方应该是独立的代码侦察步骤而非 gen-contracts 本身。如果代码侦察是 gen-contracts 的内部子步骤，应明确标注。同时，Data Flow 中缺少 CONFIDENCE 行的输入源——覆盖率公式引用 Fact Table（`confirmed/runtime 事实的 Outcome 数`），但 Data Flow 中没有 Fact Table → CONFIDENCE 的传递路径行。虽然 Other Notes 中的覆盖率公式隐含了这一关系，但 Data Flow Table 作为数据流文档应显式列出。

### Exception handling and edge cases covered — 49 / 60

Run-to-Learn Failure Handling 表格（4 种失败场景 + 兜底原则）完整。eval-skipped 降级策略扩展为 4 点且审核内容明确化（"提示用户人工审核 Journey/Contract 内容正确性"）。PRD 前置检查失败纳入 Mermaid 图和 Exit Codes 表。FIX_DECIDE 区分脚本问题和 Contract 语义错误两种回退路径。

**Deduction (-5)**: Flow 步骤 8 描述 gen-contracts 的 schema 验证失败处理（"验证失败则记录不符合项明细，自动重新生成一次...重试仍失败则暂停管线，输出不符合项供人工修正"），但这一失败路径未纳入 Mermaid 图和 Pipeline Exit Codes 表。Mermaid 图中 GEN_CONTRACT → EVAL_C 的直连没有 schema 验证失败的决策分支。这是与 Run-to-Learn Failure Handling 同级别的错误路径，文档化水平不一致。

**Deduction (-3)**: gen-test-scripts 验证失败（步骤 10："验证失败则自动重试生成一次，重试仍失败则标记该测试文件为 gen-failed 并跳过"）同样未纳入 Pipeline Exit Codes 表。gen-failed 标记是一种特殊的降级产物，但 Exit Codes 表未定义当所有测试文件都 gen-failed 时的管线行为——是退出码 0（报告含失败详情）还是退出码 1（可重试）？

**Deduction (-3)**: ENV_FAIL 的用户放弃路径未定义。Mermaid 图只展示了 "用户修复后重新检测 → ENV_CHECK" 循环，但如果用户无法修复环境（如缺失的服务依赖），没有"放弃管线"的退出路径。与 PAUSE_J/PAUSE_C 有 3 个恢复选项不同，ENV_FAIL 暗示用户只能无限重试。

---

## Dimension 4: User Stories — 194 / 200

### Coverage: one story per target user — 50 / 50

Forge 用户：Story 1/2/3/5/6（5 个）。Forge 维护者：Story 4/7（2 个）。所有 In-Scope 功能项均有对应 Story 覆盖。

### Format correct (As a / I want / So that) — 50 / 50

所有 7 个 story 严格遵循格式。I want 描述具体行为而非模糊动词。

### AC per story (Given/When/Then) — 50 / 50

所有 story 都有 Given/When/Then 格式的 AC。多场景 story 有多个 Given/When/Then 块。

### AC verifiability & boundary coverage — 44 / 50

**Deduction (-3)**: Story 5 有 2 个条件标记为 **human-verified**（"diff --stat 统计用户修改行数占草稿总行数的比例；目标 <= 20%"）。虽然描述了具体方法（diff --stat），仍依赖人工执行和判断。在自动化管线中，human-verified 指标可以作为质量参考但不应作为正式 AC 的通过判定标准——因为没有人能保证每次审核都使用完全相同的 diff 工具和统计口径。

**Deduction (-3)**: Story 6 "边界/异常 Outcome 占比 >= 30%（基线：初始静态侦察时的占比，标记为 **human-verified**）" 同样是 human-verified AC。且 "基线：初始静态侦察时的占比" 需要两次管线运行的结果对比——一次作为基线，一次作为验证——这增加了 AC 验证的操作成本。缺少自动化验证路径。

---

## Dimension 5: Scenario Completeness — 138 / 150

### End-to-end scenario coverage — 57 / 60

Flow Description 四个阶段覆盖完整端到端流程。Per-Scenario Strategy 表覆盖 5 种场景类型 x 3 维度。Delivery Phasing 三阶段门禁标准均已量化。BIZ-quality-gate-001 集成关系段落消除了两个管线之间的场景衔接模糊。

**Deduction (-3)**: Flow 步骤 6-9 的结构性重复（同一 gen-contracts + eval-contract 步骤出现两次，PAUSE 路径出现两次）对端到端场景理解造成干扰。读者需要自行判断哪个版本是最终描述。这种重复不是对场景的"深度补充"，而是对同一场景的不一致描述（步骤 6 无 schema 验证，步骤 8 有），构成了场景歧义。

### Implicit assumptions surfaced — 36 / 40

前置条件在 Mermaid 图和文字中均明确标注。依赖关系清晰。BIZ-quality-gate-001 集成关系明确。Fact Table 覆盖率公式和置信度评级计算规则清晰。

**Deduction (-2)**: Data Flow Table 中 CONFIDENCE 行描述输出传递（"嵌入生成的测试文件头部注释"），但置信度评级的计算依赖 Fact Table 中 confirmed 事实的占比。Fact Table 到 CONFIDENCE 的数据流依赖仅在 Other Notes 覆盖率公式中隐含，Data Flow Table 未显式列出这条输入路径。一个只看 Data Flow Table 的实现者会不知道 CONFIDENCE 节点需要读取 Fact Table。

**Deduction (-2)**: 场景检测信号表涵盖 16 种组合，但未定义匹配优先级。当一个项目同时匹配多个场景类型时（如 `Cargo.toml` + `clap` + `http.Handler` 既匹配 CLI 又匹配 API），表中每种信号组合是独立匹配还是互斥匹配？"且无前端入口"是部分信号的排除条件，但全局的匹配策略（第一个匹配 / 最佳匹配 / 全部匹配）未说明。

### Business-rules consistency — 45 / 50

BIZ-error-reporting-001/002: Pipeline Exit Codes 表遵循 0/1/2 语义。BIZ-task-lifecycle-003: 删除的系统保留类型（test.gen-cases、test.graduate）与规则一致。BIZ-quality-gate-001: 集成关系段落明确串行执行、独立判定。

**Deduction (-3)**: Pipeline Exit Codes 表中 FIX_DECIDE 修复耗尽的退出码为 0（"成功 — 报告含失败详情"），语义是"管线成功完成，只是测试有失败"。但 BIZ-quality-gate-001 的 e2e regression 步骤如果测试失败，通常会创建 fix task（P0）。新管线的 FIX_DECIDE 修复耗尽后退出码 0 会让上游的 BIZ-quality-gate-001 认为测试通过，可能导致质量门禁误判。两个系统对"测试失败但管线正常退出"的语义不一致。

**Deduction (-2)**: Out of Scope 包含 "已使用 gen-test-cases 项目的迁移工具"，但作为破坏性变更（删除 gen-test-cases），PRD 未提供已有用户的升级路径。In Scope 中删除清单（Story 1 AC）只描述了删除什么，没有描述已有项目如何过渡。

---

## Dimension 6: Edge Case Coverage — 93 / 100

### Error paths documented — 38 / 40

错误路径覆盖全面且结构化程度高：PRD 前置检查失败（Mermaid 图 + Exit Codes 表）、场景检测失败、Convention 缺失 + test-guide 拒绝重试、eval 不达标 + 3 轮迭代 + 3 条恢复路径、eval 评分解析失败（eval-skipped + 4 点降级策略）、gen-contracts schema 验证失败、gen-test-scripts 验证失败（gen-failed 标记）、环境不就绪、Run-to-Learn 4 种失败场景 + 兜底原则、测试失败（区分脚本 vs Contract 语义 + 2 次修复上限）。

**Deduction (-2)**: gen-contracts schema 验证失败路径（步骤 8："验证失败则记录不符合项明细，自动重新生成一次...重试仍失败则暂停管线"）未纳入 Pipeline Exit Codes 表和 Mermaid 图。这是一个与 eval 评分失败同级别的错误路径（涉及自动重试和管线暂停），但文档化水平不一致。

### Boundary conditions covered — 32 / 35

明确边界条件：eval 迭代 3 轮上限、gen-contracts schema 验证重试 1 次、gen-test-scripts 验证重试 1 次、FIX_DECIDE 自动修复 2 次、Convention 重试 2 次、Run-to-Learn 3 轮迭代、风险密度 3 级分类和 Outcome 数量范围（无重叠）、Per-Scenario Strategy 5 种场景类型差异化参数。

**Deduction (-3)**: 缺少并发生成边界条件——多个功能同时运行测试生成时，`.forge/session.yaml` 的并发写入、Fact Table 的并发读写、测试输出的文件锁冲突等问题没有讨论。Data Flow Table 中多个步骤（SCENE_DETECT、EVAL_J/C）都写入 `.forge/session.yaml`，并发场景下可能产生数据竞争。

### Failure recovery described — 23 / 25

eval-skipped 降级策略 4 点完整且审核目标明确化。PAUSE_J/PAUSE_C 恢复路径（3 选项）覆盖了关键恢复场景。FIX_DECIDE 双回路线（脚本问题 → GEN_SCRIPTS，Contract 语义错误 → GEN_CONTRACT）覆盖了测试失败的主要恢复路径。

**Deduction (-2)**: ENV_FAIL 的恢复路径只有"用户修复后重新检测 → ENV_CHECK"循环，没有放弃选项。与 PAUSE_J/PAUSE_C 有 3 条恢复路径不同，ENV_FAIL 暗示用户必须解决环境问题才能继续，缺乏"跳过环境检测，降级执行"或"放弃管线"的退出路径。在环境问题无法短期解决的情况下（如缺少外部服务依赖），用户没有明确的退出方式。

---

## Dimension 7: Scope Clarity — 85 / 100

### In-scope items are concrete deliverables — 33 / 35

每个 in-scope 条目指向具体模块、功能或文件。Eval Rubric 维度框架有子表格展开。Convention Schema 必需 Section 有字段级定义。BIZ-quality-gate-001 集成关系有 4 点说明。

**Deduction (-2)**: "质量门禁更新以反映新管线" 描述更具体了（"将现有单一门禁替换为多阶段门禁...门禁结果汇入统一质量报告"），但仍未指明具体修改的配置文件、代码入口或 hook 点。"统一质量报告" 的格式、位置和消费者也未定义。对于需要实现此条目的开发者，信息不足。

### Out-of-scope explicitly lists deferred items — 26 / 30

11 个 out-of-scope 条目，每个有明确描述。两个条目附带注解解释边界。

**Deduction (-2)**: "已使用 gen-test-cases 项目的迁移工具" 被列为 Out of Scope 但没有提供替代方案。作为破坏性变更（删除 gen-test-cases 及所有相关文件），已有用户需要知道如何过渡。Iteration 2 已指出此问题，Iteration 3 未修正。至少应有一句话说明已有用户的升级路径（如"用户需手动迁移至 Journey-Contract 路径"或"参考 Story 1 AC 中的删除清单进行清理"）。

**Deduction (-2)**: "执行环境自动准备与配置（仅做就绪检测）" 的边界不清晰。步骤 12 "场景特定环境就绪检测" 描述了检测行为，但 ENV_FAIL 后用户需要"修复环境"——修复过程中是否可以使用管线提供的辅助信息？如果 ENV_CHECK 输出的"修复建议"包含具体的安装命令，这是否已经超出了"仅做就绪检测"的边界？

### Scope consistent with functional specs and user stories — 26 / 35

In Scope 条目与 User Stories 基本对应。Functional Specs 8 个模块变更覆盖 Scope 条目。

**Deduction (-5)**: Flow Description 步骤 6-9 存在结构性重复，导致 Flow 描述与 Scope/Functional Specs 的对应关系出现歧义。Functional Specs #4 只有一个 "gen-contracts: 增加边界衍生能力"，但 Flow 中有两个 gen-contracts 步骤（步骤 6 和步骤 8，后者多出 schema 验证）。用户阅读 Flow 时会困惑：gen-contracts 的完整行为是哪个步骤？如果只读步骤 6 会遗漏 schema 验证，只读步骤 8 会遗漏与步骤 7 的上下文。

**Deduction (-4)**: Flow Description 步骤 6 和步骤 8 都描述 gen-contracts，但内容不同。步骤 6 是"从 Journey 生成 6 维度合约规范，自动衍生边界/异常 Outcome"，步骤 8 是同样的开头加上"合约生成后执行 schema 验证...验证失败则...暂停管线"。如果步骤 8 是步骤 6 的修正版（增加 schema 验证），则步骤 6 应被删除。如果步骤 6 和步骤 8 是两个不同的步骤（gen-contracts 被执行两次），则 Functional Specs 应解释为什么需要两次执行。当前文档未澄清这一关系。

---

## Blindspot Attacks

### [blindspot-1] Flow 步骤 6-9 结构性重复导致歧义

**Quote**: 步骤 6: "gen-contracts 从 Journey 生成 6 维度合约规范，自动衍生边界/异常 Outcome" vs 步骤 8: "gen-contracts 从 Journey 生成 6 维度合约规范，自动衍生边界/异常 Outcome。合约生成后执行 schema 验证（6 维度结构完整性 + Outcome Preconditions 互斥性检查）..."

**Issue**: gen-contracts 步骤出现两次（步骤 6 和步骤 8），eval-contract 步骤出现两次（步骤 7 和步骤 9），PAUSE 恢复路径也出现两次。步骤 8 比 6 多了 schema 验证描述，这暗示步骤 8 可能是步骤 6 的完整版/修正版，但步骤 6 未被删除。读者有两种理解：(1) gen-contracts 和 eval-contract 各只执行一次，步骤 8-9 是步骤 6-7 的完整版；(2) gen-contracts 和 eval-contract 各执行两次。两种理解导致完全不同的管线行为。这不是措辞问题，是结构歧义。

**What must improve**: 合并步骤 6-7 和步骤 8-9，保留一份完整的 gen-contracts + eval-contract 步骤描述（含 schema 验证）。删除重复的 PAUSE 路径。确保步骤编号连续且无重复。

### [blindspot-2] gen-contracts schema 验证失败未纳入 Pipeline Exit Codes 和 Mermaid 图

**Quote**: 步骤 8: "验证失败则记录不符合项明细，自动重新生成一次（将 schema 错误作为反馈注入 prompt），重试仍失败则暂停管线，输出不符合项供人工修正"

**Issue**: gen-contracts 的 schema 验证失败有完整的错误处理逻辑（重试一次 → 暂停管线），但这条错误路径只存在于 Flow 文字中。Pipeline Exit Codes 表没有对应的退出码。Mermaid 图中 GEN_CONTRACT → EVAL_C 是直连箭头，没有 schema 验证失败的决策分支。这与 PRD 前置检查（已纳入 Exit Codes 和 Mermaid 图）和 eval 评分失败（已纳入 Mermaid 图的 EVAL_J_SKIP/C_SKIP）的文档化水平不一致。实现者如果只看 Mermaid 图会遗漏这个失败路径。

**What must improve**: 在 Pipeline Exit Codes 表中新增 gen-contracts schema 验证重试耗尽的退出码（建议：退出码 1，retryable）。在 Mermaid 图中 GEN_CONTRACT 节点后增加 SCHEMA_CHECK 决策菱形。

### [blindspot-3] ENV_FAIL 无放弃路径，用户可能无限循环

**Quote**: Mermaid 图: `ENV_FAIL[输出缺失项+修复建议] -->|用户修复后重新检测| ENV_CHECK`

**Issue**: ENV_FAIL → ENV_CHECK 的循环没有退出条件。用户在环境问题无法解决时（如缺少付费的外部服务依赖、缺少管理员权限安装系统级依赖），没有"放弃管线"或"跳过环境检测并降级执行"的选项。PAUSE_J/PAUSE_C 有 3 条恢复路径，FIX_DECIDE 有"否"选项，但 ENV_FAIL 只有"修复后重试"一条路。这个设计假设所有环境问题都能被用户解决，不符合实际场景。

**What must improve**: 为 ENV_FAIL 添加放弃路径（如 `ENV_FAIL -->|用户放弃| END_CANCEL`）或在 Exit Codes 表中新增 ENV_FAIL 用户放弃的退出码（建议：退出码 1，retryable）。

---

## Score Summary

| Dimension | Score | Max |
|-----------|-------|-----|
| 1. Background & Goals | 98 | 100 |
| 2. Flow Diagrams | 150 | 150 |
| 3. Flow Completeness | 182 | 200 |
| 4. User Stories | 194 | 200 |
| 5. Scenario Completeness | 138 | 150 |
| 6. Edge Case Coverage | 93 | 100 |
| 7. Scope Clarity | 85 | 100 |
| **Total** | **940** | **1000** |
