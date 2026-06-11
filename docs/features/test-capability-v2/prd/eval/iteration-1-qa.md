# QA Evaluation Report — Test Capability v2.0 PRD

**Evaluator**: Senior QA Engineer (Adversarial)
**Iteration**: 1
**Mode**: B (No UI — prd-ui-functions.md absent)
**Date**: 2026-05-23

---

## Phase 1: Reasoning Audit (Pre-Score Anchors)

Before applying the rubric, I traced the document's core argument chain independently:

**P1. Problem → Solution alignment**: The three stated problems (dual-path confusion, shallow testing, limited generality) map well to the solution pillars (pipeline unification, depth enhancement, extensibility). However, the "limited generality" problem cites only 3 frameworks lacking Convention support, while the solution (test-guide auto-generation) is a far more ambitious capability than simply adding 3 new Convention files. The solution overreaches relative to the problem scope — auto-generating Conventions from project signals is a qualitatively different capability than shipping pre-built Convention files.

**P2. Solution → Evidence alignment**: The Goals table provides metrics, but several key metrics are either unmeasurable as stated (eval accuracy >= 850/1000 without defined ground truth) or rely on undefined computation (Fact Table coverage rate without a formula). Evidence does not fully support the solution's claims of measurability.

**P3. Evidence → Success Criteria**: Delivery phasing defines gate standards ("2+ projects run full pipeline without errors"), but these gates test integration completeness, not quality of the generated tests. A pipeline could produce low-quality tests and still pass the gate. The success criteria test proxies, not the actual outcome (test quality).

**P4. Self-contradiction check**: The document claims to eliminate gen-test-cases but retains test.graduate deletion as a separate in-scope item — these are related but the PRD treats them independently without explaining the relationship. Also, "Mobile best-effort" in Per-Scenario Strategy contradicts the Goal "降低 Mobile 接入成本" — if the approach is best-effort, it's unclear how cost is meaningfully lowered.

These anchors will be channeled into blindspot attacks where they identify issues not fully captured by rubric dimensions.

---

## Phase 2: Rubric Scoring

### 1. Background & Goals — 72/100

| Criterion | Score | Rationale |
|-----------|-------|-----------|
| Background 三要素 (Reason/Target/Users) | 27/30 | Why/Target/Who 三段齐全。Why 列出三大结构性缺陷并引用了具体命令名。Target 明确了 5 个升级方向。Who 区分了项目开发者和维护者。**扣分**: Who 部分遗漏了"已使用 gen-test-cases 的存量用户"——退休旧路径直接影响到他们，但他们未被列为受影响群体。此外 Background 只提到 3 大缺陷但 Target 有 5 个升级方向（"评测补全"和"信息增强"在 Why 中只是隐含问题，未被列为独立缺陷），逻辑上存在 Target > Problem 的不对称。 |
| Goals 量化 | 20/30 | 6 个 Goal 均有 Metric 列，量化程度参差不齐。可测量的：`高风险 Journey 平均测试数 >= 8`、`高风险 >= 低风险 x 1.5`、`Fact Table 覆盖率提升 >= 20 个百分点`、`内置 >= 3 个新 Convention`。**扣分点**: (1) `eval-journey/eval-contract 评分 >= 850/1000（基于 gold standard 评分集校准）`——gold standard 评分集本身不存在于当前 PRD 中，校准方法的定义在 Other Notes 中但校准数据集的构建时间点、维护责任不明确；Goal 的"准确率"措辞实际上描述的是"绝对评分值"而非"准确率"，用词不当。 (2) `消除双路径` 的 metric 是"gen-test-cases 及所有相关文件完全删除"——这是二值可验证的，但"所有相关文件"的范围边界不清晰（是否包括历史 PRD/eval 报告中的引用？）。 (3) `降低 Mobile 接入成本` 的 metric 是"生成 Maestro YAML 骨架 + deep link 测试"——这没有任何量化指标，无法验证"成本降低"的程度。扣 10 分。 |
| 背景与目标逻辑一致性 | 25/40 | 背景的三大问题与 Goals 大致对应：双路径 → 消除双路径；测试深度不足 → 风险驱动密度；通用性有限 → Convention 扩充。**扣分点**: (1) 背景中"评测门禁缺失"是 Journey-Contract 路径的一个子问题，但 Goal 将其升级为独立的"建立评测门禁"目标——升级合理但逻辑链不完整，缺少"为什么评测门禁缺失比其他子问题更值得独立解决"的论证。(2) 背景的"信息缺口"问题在 Goal 中被表述为"Run-to-Learn 迭代后 Fact Table 覆盖率提升 >= 20 个百分点"，但覆盖率的计算公式直到 Other Notes 才定义——如果公式的定义使得基线无法确定（首次运行时 Fact Table 的初始内容是什么？），这个 Goal 可能无法测量。(3) Goal "降低 Mobile 接入成本"与背景的联系薄弱——Background 三大缺陷中没有一条专门提到 Mobile。(4) "v3.0.0 是重构的最佳窗口"这一论证出现在 Why 中但从未在 Goal 中体现为交付时间约束或版本绑定。扣 15 分。 |

### 2. Flow Diagrams — 105/150

| Criterion | Score | Rationale |
|-----------|-------|-----------|
| Mermaid 图存在 | 50/50 | 存在完整的 mermaid flowchart，使用 START/END 节点、决策菱形、处理节点，覆盖了管线的主要阶段。 |
| 主路径完整 (start → end) | 25/50 | 主路径存在：START → SCENE_DETECT → DETECT → GEN_JOURNEY → EVAL_J → GEN_CONTRACT → EVAL_C → GEN_SCRIPTS → R2L_CHOICE → ENV_CHECK → CONFIDENCE → RUN_TESTS → REPORT → END。**扣分点**: (1) 流程图缺少 PRD 作为输入源——GEN_JOURNEY 的输入是 PRD 用户故事，但 START 节点后直接进入场景检测，PRD 的存在性检查未在图中体现（虽然 Flow Description 中提到了）。(2) SCENE_FAIL 的恢复路径标注"选择存入 session 缓存"回到 DETECT，但 DETECT 只检查 Convention 文件是否存在——场景类型检测结果没有在图中作为状态传递到后续步骤。(3) test-guide 的"用户审核确认"是人工交互节点，在自动化管线中如何实现（暂停等待？CLI prompt？）未在图中体现。(4) FIX_DECIDE 的"是"分支回到 GEN_SCRIPTS 但标注"<= 2 次"——循环上限在图中以注释形式存在但未以决策菱形表达。(5) CONFIDENCE → RUN_TESTS 没有决策点——即使全部 LOW 置信度也直接执行测试，缺少"LOW 置信度是否需要人工确认"的判断。(6) Report 中"全部通过"和"测试失败"分支缺少对 FIX_DECIDE 的"Contract 语义错误"回退路径——图中只画了回到 GEN_SCRIPTS，但 Flow Description 文字提到还应回退到 gen-contracts。扣 25 分。 |
| 决策点 + 错误分支 | 30/50 | 评测门禁的迭代和 3 轮耗尽分支（PAUSE_J/PAUSE_C）覆盖良好。R2L 可选分支合理。**扣分点**: (1) 缺少 SCENE_DETECT 检测失败的分支——虽然文字提到"无法匹配或匹配到多个类型 → 暂停管线"，但图中 SCENE_DETECT 的失败分支（未知类型）没有连接到任何处理节点。(2) PAUSE_J/PAUSE_C 后用户决定"继续"或"放弃"的分支未画出——它们是终止节点但文字说"用户决定后续操作"。 (3) Run-to-Learn 迭代中骨架测试执行失败（编译错误、运行时崩溃、脏数据输出）的分支完全缺失——图只画了成功路径（"<= 3 轮或覆盖率达标"→ ENV_CHECK）。(4) ENV_FAIL 有"输出缺失项+修复建议"和"用户修复后重新检测"的回路（在图中有箭头回到 ENV_CHECK），但 R2L 失败没有类似回路或降级路径。扣 20 分。 |

### 3. Flow Completeness — 125/200

| Criterion | Score | Rationale |
|-----------|-------|-----------|
| 流程步骤描述完整业务过程 | 55/70 | 阶段一至四的 13 个步骤覆盖了从管线准备到执行报告的完整过程。每阶段有明确输入输出。场景类型检测规则有详细表格。**扣分点**: (1) 步骤 2 的场景类型检测表格有 10 条规则但没有优先级——如果项目同时匹配 CLI（`main.go` + `cobra.Command`）和 API（`main.go` + `http.Handler`），按什么规则判定？表格最后一条"无法匹配或匹配到多个类型 → 暂停管线"提供了兜底，但缺少判定优先级。(2) 步骤 9 Run-to-Learn 的迭代终止条件不完整——说"<= 3 轮或覆盖率达标"，但覆盖率达标阈值未在 Flow Description 中定义（虽然在 Other Notes 中有公式）。(3) 步骤 13 输出报告的格式、输出位置、内容结构未说明。(4) 步骤 14 的失败修复策略区分了"脚本问题"和"Contract 语义错误"两种回退路径，但没有说明如何区分这两种失败类型——是靠错误码？靠 LLM 分析？扣 15 分。 |
| 数据流文档 | 35/70 | 没有专门的数据流表。隐式数据流 PRD → Journey → Contract → Test Scripts 可推断。**扣分点**: (1) Convention 文件在管线中的消费路径不明确——哪些步骤读取 Convention？格式是什么？test-guide 生成草稿后如何被后续步骤发现和加载？(2) Fact Table 的数据模型直到 Other Notes 才定义，但 Flow Description 中大量引用"覆盖率"概念——应将核心数据模型前置到 Flow Description 或至少引用。(3) eval-journey/eval-contract 的评分结果数据结构未说明——下游 agent 如何消费评分反馈来修正 Journey/Contract？(4) 置信度评级的输入（什么数据源决定 HIGH/MEDIUM/LOW？）和输出（标注在哪里？测试文件的注释？报告字段？）未说明。(5) 场景类型检测结果如何在后续步骤间传递（session 缓存的具体机制？）未说明。扣 35 分。 |
| 异常处理与边界情况 | 35/60 | **覆盖的**: 评测门禁迭代耗尽后暂停、环境就绪检测失败、场景类型检测失败（混合类型暂停管线）、PRD 不存在时步骤 1 报错、FIX_DECIDE 修复耗尽。**缺失的**: (1) gen-journeys 提取失败（PRD 格式异常、无用户故事可提取）的处理未说明——Flow Description 前置条件只说"项目必须已有 PRD"，没说 PRD 质量要求。(2) gen-contracts 生成合约后 schema 验证失败的处理未说明。(3) Run-to-Learn 骨架测试编译失败/运行时崩溃/脏数据输出——虽然 Other Notes 的 Run-to-Learn Failure Handling 表格详细定义了 4 种失败场景的处理策略，但这部分信息不在 Flow Description 中，应至少引用。(4) Convention 草稿自动生成失败（无法检测到任何已知测试框架）的处理未说明。扣 25 分。 |

### 4. User Stories — 148/200

| Criterion | Score | Rationale |
|-----------|-------|-----------|
| 覆盖率：每类目标用户一个 Story | 35/50 | Background 定义了两类用户：项目开发者和维护者。7 个 Story 中 5 个面向项目开发者，2 个面向维护者（Story 4 评测门禁、Story 7 可扩展场景类型系统）。**扣分点**: (1) "Forge 用户（项目开发者）"是一个过于宽泛的角色——它同时是 Story 1-3, 5-6 的主角，但没有区分"新用户"（首次使用管线）和"存量用户"（已使用 gen-test-cases）的需求差异。存量用户受退休旧路径影响最大但没有对应 Story。(2) 维护者只有 2 个 Story，但 Background 说维护者需要"清晰的管线架构和可扩展的场景类型系统"——管线架构的清晰性（如 eval rubric 维度设计、评分阈值设定）没有对应 Story。扣 15 分。 |
| 格式正确 (As a / I want / So that) | 48/50 | 7 个 Story 均使用 As a / I want / So that 格式。**扣分点**: Story 1 的 "I want to" 包含括号说明"不需要在 gen-test-cases 和 Journey-Contract 之间选择"——这是对旧状态的反面描述而非对期望行为的正面描述，略显间接但可接受。Story 3 的 "I want to" 列举了 CLI/WebUI/Mobile 三种场景的具体策略——这更像是实现细节而非用户意图。扣 2 分。 |
| AC per Story (Given/When/Then) | 35/50 | 每个 Story 都有 AC。**扣分点**: (1) 格式不完全统一——Story 3 有两个独立的 Given/When/Then 块但未编号区分（"Given 一个 CLI 类型的项目" vs "Given 一个 Mobile 类型的项目"是两个独立场景但缺乏编号）。(2) Story 5 有三个 Given/When/Then 块，第三个块（"Given 内置 Convention 库 / Then 包含 pytest、JUnit、Rust/cargo test 共 >= 3 个新增 Convention 文件"）缺少 When 条件。(3) Story 6 的 AC 中 `And` 子句过多（4 个 And），模糊了核心 Then 断言与补充断言的层级。(4) Story 7 的 AC "无需修改任何管线技能代码"是一个否定式断言——难以正面验证。扣 15 分。 |
| AC 可验证性与边界覆盖 | 30/50 | **可验证性问题**: (1) Story 5 AC "用户审核修改量评估方式为：diff --stat 统计用户修改行数占草稿总行数的比例；目标 <= 20%（此指标为人工审核参考，标记为 **human-verified**）"——尽管 PRD 已经将其标注为 human-verified 并明确了 diff --stat 作为度量方式，但 "目标 <= 20%"作为 AC 的通过/失败判定标准仍然依赖人工判断"修改是否合理"，diff --stat 只能量化"修改了多少"而不能量化"修改是否因为草稿质量差"。 (2) Story 2 AC "同一功能的 Journey 比较高风险 vs 低风险 variant 的总 Outcome 数量，高风险 >= 低风险 x 1.5"——如果某个功能只有高风险 Journey（没有低风险 variant），这个 AC 无法验证。需要说明退化为绝对值的情况。(3) Story 6 AC "经过 <= 3 轮迭代后，Fact Table 覆盖率从初始值（gen-contracts 静态侦察结果）提升 >= 20 个百分点"——覆盖率公式在 Other Notes 中定义了，但"初始值"的获取方式（gen-contracts 静态侦察的结果如何映射到 Fact Table？）不明确。(4) Story 7 AC "已有场景类型（CLI/TUI/WebUI/Mobile/API）的测试生成结果不受影响（回归验证）"——"不受影响"如何定义？完全相同？语义等价？允许合理的随机差异？(5) Story 1 AC 要求删除 gen-test-cases 相关内容并列出了具体清单，这是可验证的——但"全局搜索 gen-test-cases 关键词（除 PRD/历史文档外）无匹配"的范围判定标准（什么是"历史文档"？eval 目录下的报告算吗？）不精确。扣 20 分。 |

### 5. Scenario Completeness — 90/150

| Criterion | Score | Rationale |
|-----------|-------|-----------|
| 端到端场景覆盖 | 38/60 | Per-Scenario Strategy 表格提供了 5 种场景类型的差异化策略。风险驱动密度表格给出了三级分类。Flow Description 的 4 阶段 13 步覆盖了管线全生命周期。**扣分点**: (1) Mobile "尽力而为"策略过于模糊——"只生成 Maestro YAML 骨架 + deep link 测试"但"骨架"的具体定义（包含哪些操作？导航到几层？）和"复杂场景标记 manual-only"的判定标准未说明。(2) TUI 场景的测试执行环境准备方式完全空白——TUI 应用需要终端模拟器或 headless 模式，但 Per-Scenario 表格中 TUI 行只有 `timeout` 边界，环境就绪检测标准缺失（CLI 有"二进制编译"、WebUI 有"dev server"、API 有"服务+DB"，TUI 有什么？）。(3) API 场景的"平衡 50/50"——Contract 测试和 Journey 烟测试各覆盖什么范围？如果 Contract 测试已覆盖 API 端点的全部参数组合，Journey 烟测试的增量价值是什么？(4) 场景类型检测的端到端场景未完整覆盖——检测成功（单类型）、检测失败（混合类型）有描述，但检测到类型后的"确认"环节如何与用户交互未说明（自动确认？用户确认？）。扣 22 分。 |
| 隐式假设暴露 | 18/40 | **未暴露的隐式假设**: (1) 假设 PRD 中存在可提取的用户故事——如果 PRD 的 User Stories 部分为空或质量差，gen-journeys 的输入是什么？(2) 假设项目的场景类型是单一互斥的——但实际项目可能同时是 CLI + API（如 Go CLI 工具附带 HTTP API），表格有兜底但策略未定义。(3) 假设 Convention 文件有明确的 schema——Story 5 AC 引用了"4 个必需 section"但 PRD spec 的 Functional Specs 中才定义了这些 section，前面的 Flow Description 已经隐式依赖了它们。(4) 假设 eval 评分结果可以被 revise 技能消费——但评分结果到修正指令的映射机制未说明。(5) 假设 LLM 生成的边界 Outcome 在语义上有意义——"LLM prompt 增强策略"是否能产出有实际测试价值的边界场景？R2L 迭代如果证明 LLM 衍生的边界是无效的（与实际行为不符），是否有反馈机制？(6) 假设骨架测试可以在隔离环境中运行——但 Write 端点的副作用处理（生成回滚语句）暗示测试环境可能有状态污染。扣 22 分。 |
| 业务规则一致性 | 34/50 | 主体一致。**矛盾/不一致**: (1) **BIZ-task-lifecycle-003** 规定 `test.graduate` 是系统保留类型，非自动生成任务不得使用。PRD 将"删除 test.graduate 任务类型"列为 in-scope——这意味着从系统注册中移除一个保留类型。但该业务规则列出 13 个保留类型，PRD 只提议删除其中一个，是否需要同步清理保留类型列表的文档或代码？未说明。(2) **BIZ-quality-gate-001** 定义了质量门禁的三个阶段（compile → test → e2e），PRD 的"质量门禁更新以反映新管线"是 in-scope 但未说明如何更新——是修改质量门禁技能还是修改配置？(3) PRD 的 Functional Specs 中 Convention Schema 必需 Section 定义表在"Related Changes"下方，但这个定义实际上是新 Convention 文件和 test-guide 的核心约束——它应该在 Scope 或 Core Concepts 中定义而非藏在 Related Changes 中。(4) Story 2 AC 要求"高风险 >= 低风险 x 1.5"，Goals 表格也有相同指标，但风险驱动密度表显示 High 总测试数估算为 10-20，Low 为 4-8——10/8 = 1.25，20/4 = 5.0，范围太大，下界 1.25 不满足 1.5 倍要求。扣 16 分。 |

### 6. Edge Case Coverage — 48/100

| Criterion | Score | Rationale |
|-----------|-------|-----------|
| 错误路径文档 | 18/40 | **已覆盖**: 评测门禁迭代耗尽（PAUSE_J/PAUSE_C）、环境就绪检测失败（ENV_FAIL）、场景类型检测混合/未知（SCENE_FAIL）、FIX_DECIDE 修复耗尽、PRD 不存在。Run-to-Learn Failure Handling 表格在 Other Notes 中详细定义了 4 种失败场景（编译失败、运行时崩溃、脏数据输出、API 写操作副作用），这是加分项。**未覆盖**: (1) gen-journeys 提取失败——PRD 存在但无用户故事、用户故事格式异常。(2) gen-contracts 合约 schema 验证失败。(3) gen-test-scripts 代码生成失败（语法错误、框架不兼容）。(4) eval 评分结果无法解析——Flow Description 中提到"若 eval 评分因 LLM 输出无法解析而失败，记录错误日志并重试评分一次；重试仍失败则跳过门禁，标记该 Journey/Contract 为 eval-skipped"，这是一个好的错误处理设计，但它定义了一种绕过门禁的降级路径——eval-skipped 的文档质量没有保障，对下游 gen-test-scripts 的影响未评估。(5) Convention 草稿自动生成失败（无法检测到任何已知测试框架）。(6) run-tests 执行中测试运行器本身崩溃（非测试失败，而是执行环境异常）。扣 22 分。 |
| 边界条件 | 17/35 | **已覆盖**: 评测迭代上限 3 轮、Run-to-Learn 迭代上限 3 轮、风险等级三级分类、FIX_DECIDE 修复上限 2 次、Journey 3-5 步的测试密度估算。**未覆盖**: (1) 0 个 Journey 的极端情况——PRD 存在但无用户故事可提取。(2) 单个 Journey 包含大量 Step（如 50+）时的性能和超时影响。(3) eval 评分恰好在维度阈值边界（如完整性恰好 120/200）的通过/不通过判定。(4) 场景类型检测为混合型（CLI + API 同时存在）时的策略选择——表格有兜底但无策略。(5) 置信度评级的边界——全部 HIGH 和全部 LOW 的处理策略差异。(6) Convention 文件格式损坏或版本不兼容时的容错。(7) 0 条运行时 Fact（Run-to-Learn 首轮骨架测试全部失败）时的覆盖率基线计算。扣 18 分。 |
| 失败恢复描述 | 13/25 | **已覆盖**: ENV_FAIL → 用户修复后重新检测（图中有回路）。eval-skipped → 标记后由用户手动审核。Run-to-Learn 失败 → 兜底原则"不应阻塞管线"，使用静态信息继续。FIX_DECIDE → 修复耗尽后输出报告。(2) Run-to-Learn 编译失败 → 记录到 Fact Table，跳过本轮运行。**未覆盖**: (1) PAUSE_J/PAUSE_C 后用户选择"继续"的后续流程——是接受当前低分文档继续？还是手动修改后重新评分？(2) eval-skipped 的 Journey/Contract 被下游 gen-test-scripts 消费后生成的测试质量无保障——是否应有特殊标记或降级处理？(3) SCENE_FAIL 用户确认场景类型后的验证——用户选择的类型可能不正确（如选择了 CLI 但实际是 API 项目），后续步骤如何处理这种不匹配？扣 12 分。 |

### 7. Scope Clarity — 78/100

| Criterion | Score | Rationale |
|-----------|-------|-----------|
| In-scope 项目是具体可交付物 | 25/35 | 15 个 in-scope 项目大部分是具体的。**模糊项**: (1) "合约规范增强：支持边界/异常场景自动衍生描述（LLM prompt 增强策略 + 场景类型 required_outcomes 规则）"——"LLM prompt 增强策略"不是可交付物，是实现手段。(2) "质量门禁更新以反映新管线"——更新什么、在哪里更新、更新后的行为是什么，均未具体化。(3) "Run-to-Learn 机制：骨架测试 → 运行捕获输出 → 丰富 Fact Table → 重新生成"——这是一个流程描述而非交付物清单，交付物应该是"骨架测试生成器"、"Fact Table 更新器"等。(4) "场景特定执行环境就绪检测（CLI/TUI/WebUI/API）"——交付物是什么？配置文件？检测脚本？扣 10 分。 |
| Out-of-scope 明确列出延迟项 | 24/30 | 11 个 out-of-scope 项明确列出。**扣分点**: (1) "已使用 gen-test-cases 项目的迁移工具"被排除——存量用户需要手动处理，但 Background 中未将存量用户列为受影响群体，也未说明手动迁移的步骤。(2) "跨场景组合编排"被排除但未说明是否在后续版本计划中。(3) "执行环境自动准备与配置（仅做就绪检测）"——这个排除项与 in-scope 的"场景特定执行环境就绪检测"的边界在哪里？就绪检测如果发现环境不满足，用户是自行准备还是管线提供安装脚本？(4) "失败诊断场景特定策略"被排除，但 FIX_DECIDE 步骤区分"脚本问题"和"Contract 语义错误"——这本身就是一种失败诊断策略。边界模糊。扣 6 分。 |
| 范围与功能规格和用户故事一致 | 29/35 | 主体一致。**不一致处**: (1) Functional Specs 表格列了 8 个 Change Points，In-scope 列表有 15 项——两者的映射关系不清晰。例如 In-scope 有"置信度评级系统（HIGH/MEDIUM/LOW）"但 Change Points 中没有对应条目。(2) Story 7（可扩展场景类型系统）在 In-scope 列表中没有直接对应项——最接近的是"场景差异化：CLI/TUI/WebUI/API 核心支持 + Mobile 尽力而为"，但 Story 7 的核心是"通过添加配置文件接入新场景类型"，这不仅仅是"差异化"而是"可扩展性"。(3) Out-of-scope 说"合约 6 维度模型（schema）修改"排除，但 In-scope 有"合约规范增强：支持边界/异常场景自动衍生描述"——如果衍生逻辑需要在 schema 中添加 metadata 字段（如 `risk_level`、`required_outcomes`），这是否算 schema 修改？Scope 边界模糊。PRD 的注释已经预见到了这个问题（"注：required_outcomes 是按场景类型配置的实例数据，不属于 schema 变更"），但这只是声明而非论证——下游 agent 需要判断具体实现是否越界时，缺乏判定标准。扣 6 分。 |

---

## Phase 3: Cross-Dimension Coherence Check

**CD-1. Goals vs. Scenario Completeness (内部数据不一致)**:
Goals 表格说"高风险 Journey 平均测试数 >= 8"且"高风险 >= 低风险 x 1.5"。Per-Scenario Strategy 的风险驱动密度表显示 Low 风险总测试数 4-8。如果 Low 取上界 8，High 需要达到 12 才满足 1.5x。但 High 的范围是 10-20，下界 10 不满足 12。更关键的是 High 的下界 10 不满足"平均 >= 8"的可靠性——如果大部分 High Journey 落在 10 附近，平均可能接近 8 但不保证超过。这个数据不一致分散在 Goals 和 Flow Description 两个位置。

**CD-2. User Stories vs. Scope (交付物缺位)**:
Story 6（Run-to-Learn）和 Story 7（可扩展场景类型系统）描述了用户需求，但 In-scope 列表中 Story 7 没有直接对应的交付物项。Story 6 的"置信度评级系统"在 In-scope 中有对应项，但 Fact Table 的数据模型定义（Other Notes 中）和 Convention Schema 的定义（Functional Specs 中）分散在不同位置，缺少统一的"数据模型交付物"项。

**CD-3. Flow Diagrams vs. Flow Description (覆盖范围不一致)**:
Flow Description 步骤 14 描述了两种修复回退路径（脚本问题 → gen-test-scripts，Contract 语义错误 → gen-contracts），但 Flow Diagram 的 FIX_DECIDE 只画了回到 GEN_SCRIPTS 的路径，缺少回到 GEN_CONTRACT 的路径。文字描述比图更完整。

---

## Phase 4: Blindspot Hunt

### [blindspot-1] 场景类型检测结果在后续步骤间的传递机制完全未定义

**Quote**: "SCENE_DETECT -->|单一类型| DETECT{Convention 文件存在?}"

**Problem**: 场景类型检测的结果（CLI/TUI/WebUI/Mobile/API）是后续所有差异化策略的基础——gen-journeys 需要知道场景类型来决定 Journey 格式，gen-contracts 需要知道场景类型来决定必须衍生的 Outcome，gen-test-scripts 需要知道场景类型来选择测试框架和执行模型。但这个检测结果如何在步骤间传递？写入文件？session 变量？CLI 参数？如果用户在阶段二中途切换了项目的文件结构（如添加了 package.json），场景类型是否需要重新检测？这个"状态传播"问题是管线正确运行的先决条件，但 PRD 只在 SCENE_FAIL 的注释中提到"选择存入 session 缓存"，没有定义完整的传递机制。

**Must improve**: 定义场景类型检测结果的存储方式（文件路径、session 键名）、传递方式（每步读取？管线启动时一次性注入？）、以及变更检测策略（是否在每步重新检测？）。

### [blindspot-2] eval-skipped 降级路径绕过了质量门禁的核心目的

**Quote**: Flow Description: "若 eval 评分因 LLM 输出无法解析而失败，记录错误日志并重试评分一次；重试仍失败则跳过门禁，标记该 Journey/Contract 为 eval-skipped，由用户手动审核"

**Problem**: 质量门禁的目的是"保障下游 gen-test-scripts 收到的输入质量"（Story 4 的 So that 子句）。但 eval-skipped 路径允许未经质量验证的文档通过门禁，直接进入下游。这重新引入了 Story 4 旨在消除的风险——如果 eval 的可靠性不够（LLM 输出无法解析），门禁反而成为了可选步骤而非必经步骤。这不是边界情况——LLM 输出的结构化解析在生产环境中是常见的失败点。

**Must improve**: 为 eval-skipped 路径定义降级策略——例如自动降低该批次测试的置信度评级为 LOW，或在报告中醒目标注"未经评测门禁验证"。当前的"由用户手动审核"过于依赖人工介入，与管线的自动化目标矛盾。

### [blindspot-3] 风险驱动密度表的数据与 Goal 指标存在数学不一致

**Quote**: Goals: "高风险 Journey 平均测试数 >= 8，且高风险旅程测试数 >= 低风险旅程 x 1.5"; Per-Scenario Strategy: Low 总测试数 4-8，High 总测试数 10-20

**Problem**: 如果 Low Journey 的测试数在上界（8），则 High Journey 需要至少 12（8 x 1.5）才满足 Goal。但 High 的下界是 10，这意味着存在 High Journey（测试数 10-11）和 Low Journey（测试数 8）的组合不满足 1.5x 要求。同时"平均 >= 8"与"High 范围 10-20"之间存在宽松度——如果 3 个 High Journey 分别为 10, 10, 10，平均 10 >= 8 通过，但如果有一个 Low Journey 为 8，则 10 < 8 x 1.5 = 12 不通过。Goal 中的两个指标可能互相矛盾。Reasoning audit flagged this independently of dimension scoring.

**Must improve**: 调整风险驱动密度表的下界确保一致性——例如将 High 下界提升到 12，或将 Low 上界降低到 6，或在 Goal 中将 1.5x 限定为"同功能内对比"而非跨功能对比。

### [blindspot-4] "LLM prompt 增强策略"的交付物定义缺失——这是核心能力还是实现细节？

**Quote**: In-scope: "合约规范增强：支持边界/异常场景自动衍生描述（LLM prompt 增强策略 + 场景类型 required_outcomes 规则）"

**Problem**: "LLM prompt 增强策略"被列为 in-scope 交付物，但它实际上是一种实现手段而非可交付的功能。如果 Forge 更换了 LLM 模型（如从 Claude 切换到 GPT），prompt 策略可能完全失效。更重要的是，PRD 没有定义"增强"成功的标准——衍生出的边界 Outcome 什么算"有效"？什么算"幻觉"？eval-contract 的"事实依据"维度（90/150 最低阈值）只检查是否基于 Fact Table，但 LLM 衍生的边界（如"空输入"、"超长输入"、"特殊字符输入"）可能不在 Fact Table 中——这些是 LLM 的推理产物，不是事实。如果 eval-contract 拒绝所有不在 Fact Table 中的声明，则 LLM 衍生出的边界 Outcome 会被标记为 UNKNOWN，反而降低了评分。这是一个根本性的架构矛盾。

**Must improve**: 明确 LLM 衍生的边界 Outcome 在 eval-contract 评分中的处理方式——是否允许基于推理（非事实）的 Outcome？如果是，eval-contract 的"事实依据"维度需要区分"有事实依据的 Outcome"和"基于合理推理的 Outcome"两种类别。

### [blindspot-5] "生成 Maestro YAML 骨架"的"骨架"定义完全空白

**Quote**: Story 3 AC: "输出 Maestro YAML 骨架（app lifecycle + navigation）和 deep link 测试，复杂场景标记 manual-only"

**Problem**: "骨架"在此语境下是 Mobile 场景的核心交付物，但其定义完全空白。一个 Maestro YAML 文件的"骨架"至少包含：appId、launchApp、基本导航流——但"app lifecycle"包含哪些生命周期事件？"navigation"导航到哪个层级？deep link 测试需要覆盖哪些 URI scheme？如果下游 agent 需要实现这个功能，它没有任何可依据的规格。"复杂场景标记 manual-only"中的"复杂"判定标准也未定义。

**Must improve**: 在 Per-Scenario Strategy 表格中补充 Mobile 的骨架定义，至少包含：(1) Maestro YAML 的最小结构示例，(2) "app lifecycle"至少包含启动+关闭，(3) "navigation"至少包含首屏导航，(4) "复杂场景"的判定标准。

---

## Score Summary

| Dimension | Score | Max |
|-----------|-------|-----|
| 1. Background & Goals | 72 | 100 |
| 2. Flow Diagrams | 105 | 150 |
| 3. Flow Completeness | 125 | 200 |
| 4. User Stories | 148 | 200 |
| 5. Scenario Completeness | 90 | 150 |
| 6. Edge Case Coverage | 48 | 100 |
| 7. Scope Clarity | 78 | 100 |
| **Total** | **666** | **1000** |

---

## Attacks Summary

1. [Background & Goals]: Goal "评分准确率 >= 850/1000" 用词不当且 ground truth 未定义 — "eval-journey/eval-contract 评分 >= 850/1000（基于 gold standard 评分集校准）" — 明确区分"评分值"与"准确率"，定义 gold standard 数据集的构建时间点和维护责任
2. [Background & Goals]: Mobile 接入成本降低与 Background 三大缺陷无逻辑联系 — "降低 Mobile 接入成本：生成 Maestro YAML 骨架 + deep link 测试" — 要么在 Background Why 中补充 Mobile 问题陈述，要么将此 Goal 标注为独立演进方向
3. [Flow Diagrams]: FIX_DECIDE 缺少 Contract 语义错误回退到 gen-contracts 的路径 — Flow Description 步骤 14 文字描述了两种回退路径但图中只画了一种 — 在 Flow Diagram 中为 FIX_DECIDE 添加区分"脚本问题"和"Contract 语义错误"的决策分支
4. [Flow Diagrams]: Run-to-Learn 骨架测试执行失败分支完全缺失 — "Run-to-Learn 迭代 <= 3 轮或覆盖率达标 --> ENV_CHECK" 只有成功路径 — 为 R2L 节点添加骨架测试失败（编译/运行/脏数据）的分支和降级路径
5. [Flow Completeness]: 数据流中 Convention、Fact Table、置信度评级的数据传递路径未文档化 — 多个步骤隐式消费这些数据但传递机制不明 — 添加数据流表或至少在 Flow Description 中为每个步骤明确输入/输出数据
6. [User Stories]: Story 2 AC 的 1.5x 比较在只有高风险 Journey 时不可验证 — "同一功能的 Journey 比较高风险 vs 低风险 variant 的总 Outcome 数量" — 添加退化为绝对值时的备选 AC（如"高风险 Journey 总 Outcome 数 >= 8"）
7. [User Stories]: Story 5 缺少 When 条件的 Given/Then 块 — "Given 内置 Convention 库 / Then 包含 pytest、JUnit、Rust/cargo test 共 >= 3 个新增 Convention 文件" — 补充 When 条件或合并到前一个 Given/When/Then 块
8. [Scenario Completeness]: 风险驱动密度表下界数据与 Goal 1.5x 指标数学矛盾 — High 下界 10 vs Low 上界 8，10/8 = 1.25 < 1.5 — 调整密度表下界或 Goal 乘数使二者一致
9. [Edge Case Coverage]: gen-journeys 提取失败（PRD 无用户故事）未覆盖 — Flow Description 只说"项目必须已有 PRD"但未定义 PRD 质量要求 — 添加 PRD 质量前置检查或 gen-journeys 的降级处理
10. [Edge Case Coverage]: eval-skipped 绕过质量门禁的降级后果未评估 — "标记该 Journey/Contract 为 eval-skipped，由用户手动审核" — 定义 eval-skipped 文档进入下游时的降级策略（如自动降低置信度）
11. [Scope Clarity]: "LLM prompt 增强策略"作为 in-scope 交付物是实现手段而非功能 — "合约规范增强：支持边界/异常场景自动衍生描述（LLM prompt 增强策略 + 场景类型 required_outcomes 规则）" — 将"LLM prompt 增强策略"从 in-scope 列表中移除，替换为功能描述如"边界/异常场景自动衍生引擎"
12. [blindspot]: 场景类型检测结果在后续步骤间的传递机制未定义 — "选择存入 session 缓存"是唯一提到传递的地方 — 定义完整的场景类型状态传播机制（存储格式、传递路径、变更检测）
13. [blindspot]: eval-skipped 降级路径绕过质量门禁核心目的 — "跳过门禁，标记该 Journey/Contract 为 eval-skipped" — 为 eval-skipped 路径定义自动降级策略，降低管线对人工审核的依赖
14. [blindspot]: LLM 衍生边界 Outcome 与 eval-contract "事实依据"维度的架构矛盾 — eval-contract 要求声明基于 Fact Table，但 LLM 衍生的边界不在 Fact Table 中 — 在 eval-contract 维度中区分"事实依据"和"合理推理"两类声明
15. [blindspot]: Mobile Maestro YAML "骨架"定义完全空白 — "输出 Maestro YAML 骨架（app lifecycle + navigation）" — 补充骨架的最小结构定义和复杂场景判定标准
