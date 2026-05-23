# QA Evaluation Report — Test Capability v2.0 PRD

**Evaluator**: Senior QA Engineer (Adversarial)
**Iteration**: 3
**Mode**: B (No UI — prd-ui-functions.md absent)
**Date**: 2026-05-23

---

## Phase 1: Reasoning Audit (Pre-Score Anchors)

Before applying the rubric, I traced the document's core argument chain independently:

**P1. Problem → Solution**: The three structural defects (dual-path confusion, shallow testing, limited generality) map to five solution pillars. Iteration 3 improved "limited generality" by expanding scene detection rules to 16 entries covering Go, Node.js, Python, Java, and Rust. The "eval 事实依据 vs. LLM 衍生 Outcome" contradiction is now resolved via the two-category declaration system (事实声明 vs. 合理推理声明) in the rubric framework. The "Mobile 接入成本" Goal now has a quantified metric ("≤ 30 分钟" and "≥ 2 个核心 Journey").

**P2. Solution → Evidence**: Goals metrics are significantly improved. The risk-density Goal has both absolute floor ("≥ 13") and ratio ("≥ 1.5x") with degradation clause. Mobile Goal changed from deliverable list to measurable time bound. Eval gate Goal calibrated via gold standard methodology. Remaining weakness: "消除双路径困惑" metric "所有相关文件完全删除" — "所有" remains unbounded.

**P3. Evidence → Success Criteria**: Delivery phasing gates have concrete thresholds. PAUSE_J/PAUSE_C now have three defined recovery paths (skip, cancel, manual-retry). FIX_DECIDE now has safety constraint ("不降低断言严格度") albeit human-verified only. Pipeline Exit Codes table provides concrete failure state definitions.

**P4. Self-contradiction check**: Iteration 3 resolved several prior contradictions: (a) eval "事实依据" dimension now distinguishes two declaration categories, resolving the LLM-derived Outcome conflict; (b) Fact Table runtime-overwrite strategy now retains static entries as fallback, resolving the coverage formula conflict; (c) gold standard calibration methodology is now defined with scale, process, method, and update frequency. **New contradiction discovered**: Steps 6-7 and Steps 8-9 in Flow Description are near-identical duplicates — steps 6 (gen-contracts) and 8 (gen-contracts) overlap, with step 8 adding schema validation but step 6 lacking it. The PAUSE recovery path is repeated verbatim at lines 153-156 and 160-163. This structural duplication is confusing and suggests incomplete editing.

These anchors inform blindspot attacks where they identify issues outside rubric dimensions.

---

## Previous-Attack Resolution Status (Iteration 2 → 3)

| # | Attack | Status | Evidence |
|---|--------|--------|----------|
| 1 | Goal "降低 Mobile 接入成本" metric 是交付物列表而非量化指标 | **RESOLVED** | prd-spec 第 48 行: "新 Mobile 项目从零到可运行 Maestro 测试 ≤ 30 分钟（含 Convention 草稿审核）；生成 Maestro YAML 骨架 + deep link 测试覆盖 ≥ 2 个核心 Journey" — 量化时间约束 + 数量指标 |
| 2 | "v3.0.0 是重构的最佳窗口"未在 Goal 中体现版本约束 | **NOT RESOLVED** | prd-spec 第 21 行仍为 "v3.0.0 是重构的最佳窗口"，Goals 表格和 Delivery Phasing 无版本绑定 |
| 3 | SCENE_FAIL 到 DETECT 信息流断裂 | **NOT RESOLVED** | 流程图 SCENE_FAIL 箭头标注 "用户选择类型 / 选择存入 session 缓存"，但 session 缓存概念从未在 PRD 中定义 |
| 4 | eval 评分结果到 revise 技能的 Data Flow 缺失 | **RESOLVED** | Data Flow Table 第 189 行新增 EVAL_J/EVAL_C 行: "评分结果（总分 + 各维度得分 + 不通过项明细）→ revise-journey/revise-contract"，格式定义为 `{total, dimensions, failed_dims}` |
| 5 | PAUSE_J/PAUSE_C 后恢复路径未定义 | **RESOLVED** | prd-spec 第 153-163 行定义三种恢复路径: (a) 跳过门禁继续, (b) 放弃管线, (c) 修改后重跑 |
| 6 | Story 7 "不受影响"验证标准模糊 | **RESOLVED** | prd-user-stories 第 120 行: "可接受的差异范围为——测试函数数量变化 ≤ 5%、无新增编译错误、eval 评分偏差 ≤ 30 分" |
| 7 | 场景检测规则 CLI/TUI 仅覆盖 Go 项目 | **RESOLVED** | prd-spec 第 122-140 行扩展至 16 条规则，覆盖 Go、Node.js、Python、Java、Rust 五种语言的 CLI/TUI/API/WebUI/Mobile |
| 8 | gen-test-scripts 输出无质量验证 | **RESOLVED** | prd-spec 第 166 行: "生成后执行语法/可执行性验证：检查 (a) 测试文件语法正确... (b) 导入路径可解析" |
| 9 | gen-contracts 合约 schema 验证失败未覆盖 | **RESOLVED** | prd-spec 第 157 行: "合约生成后执行 schema 验证... 验证失败则记录不符合项明细，自动重新生成一次... 重试仍失败则暂停管线" |
| 10 | eval-contract "事实依据"维度与管线衍生引擎架构矛盾 | **RESOLVED** | prd-spec 第 68 行: eval rubric 维度区分"事实声明"和"合理推理声明"两类，推理声明需有 `required_outcomes` 规则支撑 |
| 11 | Fact Table runtime 覆盖 static 策略可能与覆盖率公式冲突 | **RESOLVED** | prd-spec 第 313 行: "替换时保留 static 条目作为 fallback... 若 runtime 条目的 confidence 不是 confirmed，则回退使用对应的 static 条目计算覆盖率" |
| 12 | PAUSE_J/PAUSE_C 恢复路径定义缺失 | **RESOLVED** | 与 #5 相同 |

**Resolution Summary**: 9/12 fully resolved, 0/12 partially resolved, 3/12 not resolved.

---

## Phase 2: Rubric Scoring

### 1. Background & Goals — 85/100

| Criterion | Score | Rationale |
|-----------|-------|-----------|
| Background 三要素 (Reason/Target/Users) | 26/30 | Why/What/Who 三段齐全。Iteration 3 在 Why #3 中明确提及"评测门禁缺失"和"Mobile 接入成本极高"，使 Target 的 5 个方向有完整 backing。"定位"声明明确了管线边界。**扣分**: (1) Who 部分仍遗漏"存量用户"——退休旧路径直接影响已使用 gen-test-cases 的用户，但他们未被列为受影响群体。In-scope 明确排除了"已使用 gen-test-cases 项目的迁移工具"，但 Background 未将存量用户列为受影响群体。(2) Target 有 5 个方向但 Why 只有 3 大缺陷，"评测补全"和"信息增强"在 Why 中仅隐含于缺陷 #3 和 #2，未作为独立缺陷列出。 |
| Goals 量化 | 27/30 | 6 个 Goal 均有 Metric 列，显著改进。风险密度 Goal 有绝对下限 (≥ 13) + 比值 (≥ 1.5x) + 退化条款。Mobile Goal 改为"≤ 30 分钟 + ≥ 2 个核心 Journey"，可量化可验证。Eval 门禁 Goal 有 gold standard 校准方法定义。**扣分**: (1) "消除双路径困惑"的 metric "gen-test-cases 及所有相关文件完全删除"——"所有相关文件"的范围边界不清晰（历史 eval 报告中的引用？PRD 本身的引用？其他项目中的文件？）。(2) "提升通用性" Goal "内置 ≥ 3 个新 Convention 文件"——数量指标虽然明确，但未度量"新项目接入成本降低"的程度（如从多少步到多少步）。 |
| 背景与目标逻辑一致性 | 32/40 | 三大问题与 Goals 大致对应。**扣分**: (1) "v3.0.0 是重构的最佳窗口"在 Why 中提出但从未在 Goal 或 Delivery Phasing 中体现为版本约束或时间约束——如果版本窗口是关键动机，应有对应的时间性 Goal。(2) "提升测试信息质量"Goal 引用"Fact Table 覆盖率"但 Goals 表格未定义覆盖率基线——Story 6 AC 中基线为"gen-contracts 静态侦察结果"，Goals 表不自包含。(3) Gold standard 校准方法中"每次新增场景类型或修改 rubric 维度时重新校准"——校准的维护责任（谁来做？）和成本未提及。 |

### 2. Flow Diagrams — 135/150

| Criterion | Score | Rationale |
|-----------|-------|-----------|
| Mermaid 图存在 | 50/50 | 完整的 mermaid flowchart，使用 START/END、决策菱形、处理节点，覆盖四个阶段。新增 PRD_CHECK 节点解决了前置条件缺失问题。 |
| 主路径完整 (start → end) | 42/50 | START → PRD_CHECK → SCENE_DETECT → DETECT → GEN_JOURNEY → EVAL_J → GEN_CONTRACT → EVAL_C → GEN_SCRIPTS → R2L_CHOICE → ENV_CHECK → CONFIDENCE → RUN_TESTS → REPORT → END。新增 PAUSE_J/PAUSE_C 三种恢复分支（跳过/放弃/修改后重跑）。**扣分**: (1) SCENE_FAIL → DETECT 箭头标注"选择存入 session 缓存"——session 缓存概念在 PRD 中从未定义（生命周期、作用域、过期机制），Data Flow Table 中也没有对应行。(2) R2L "覆盖率达标"阈值未在图中标注——R2L 分支条件 "≤ 3 轮或覆盖率达标"，但"达标"阈值是多少？结合覆盖率公式，这个阈值直接影响迭代行为。 |
| 决策点 + 错误分支 | 43/50 | 评测门禁迭代（REVISE_J/REVISE_C）、3 轮耗尽（PAUSE_J/PAUSE_C + 三种恢复路径）、eval-skipped 降级、R2L 失败降级（R2L_DEGRADE）、PRD_CHECK 失败、SCENE_FAIL、ENV_FAIL 回路、FIX_DECIDE 双分支（脚本问题/Contract 语义错误）覆盖良好。**扣分**: (1) CONFIDENCE 节点没有 Fact Table 输入来源——置信度评级依赖 Fact Table 中 confirmed 事实占比计算，但流程图中 Fact Table 从未作为数据源出现在 CONFIDENCE 节点的输入中。(2) TEST_GUIDE → GEN_JOURNEY 的箭头只标注"用户审核确认"——但 test-guide 检测到的框架信息如何传递到后续步骤？Convention 文件是 gen-test-scripts 的输入，不是 GEN_JOURNEY 的输入。数据流有中间环节缺失。 |

### 3. Flow Completeness — 155/200

| Criterion | Score | Rationale |
|-----------|-------|-----------|
| 流程步骤描述完整业务过程 | 55/70 | 四阶段覆盖完整生命周期。场景检测规则扩展至 16 条（覆盖 5 种语言），风险分级三级定义，PRD 质量前置检查，PAUSE 恢复路径（3 种选项），FIX_DECIDE 失败类型区分 + 安全约束。**扣分**: (1) **严重结构错误**：步骤 6-7 和步骤 8-9 是重复内容——步骤 6 (gen-contracts) 和步骤 8 (gen-contracts) 描述同一操作，但步骤 8 增加了 schema 验证而步骤 6 没有。PAUSE 恢复路径也在第 153-156 行和第 160-163 行完全重复。这使下游实现者无法确定哪组步骤是权威描述。(-5) (2) SCENE_FAIL 用户确认场景类型后无验证——用户可能选错类型（如将 API 项目标为 CLI），管线继续执行但后续测试不匹配。(3) 步骤 15 测试报告的格式、输出位置、内容结构未说明。 |
| 数据流文档 | 55/70 | Data Flow Table 有 9 行覆盖关键数据传递路径（场景类型、Convention、Journey、Contract、静态/运行时 Fact Table、置信度、eval 评分结果、测试代码）。新增 EVAL_J/EVAL_C → revise 行，包含结构化格式定义。**扣分**: (1) Convention 文件在流程中的角色不清晰——TEST_GUIDE 产出 Convention 草稿后，流程图显示直接到 GEN_JOURNEY，但 Convention 是 gen-test-scripts 的输入，中间的 Journey/Contract 生成步骤不消费 Convention。Data Flow Table 说 Convention 写入 `conventions/` 目录供 gen-test-scripts 读取，但流程图上 TEST_GUIDE → GEN_JOURNEY 的箭头暗示 Convention 在 GEN_JOURNEY 之前就被消费了。(-5) (2) PRD 文档作为 GEN_JOURNEY 的输入不在 Data Flow Table 中——但 Flow Description 前置条件明确要求 PRD 存在。(-3) (3) session 缓存（SCENE_FAIL 写入）不在 Data Flow Table 中。(-2) |
| 异常处理与边界情况 | 45/60 | 重大改进: gen-contracts schema 验证（步骤 8/157 行），gen-test-scripts 语法/可执行性验证（步骤 10/166 行），PAUSE 恢复三种路径，FIX_DECIDE 失败类型区分 + 安全约束（"不降低断言严格度"），Run-to-Learn Failure Handling 4 种失败场景，eval-skipped 降级策略，Pipeline Exit Codes 6 种退出码。**未覆盖**: (1) SCENE_FAIL 用户确认场景类型后——用户选择的类型可能不正确，无验证机制。(-3) (2) Convention 草稿重试 2 次用尽后的处理——流程图中 TEST_GUIDE 有"用户拒绝, 重试 ≤ 2"的循环，但 2 次后仍未通过怎么办？管线是暂停还是回退？(-3) (3) PRD 中的步骤 6 (gen-contracts 无 schema 验证) 和步骤 8 (gen-contracts 有 schema 验证) 共存——哪个是正确的流程？如果步骤 6 是正确的（无 schema 验证），则 schema 验证覆盖无效；如果步骤 8 是正确的，则步骤 6 的描述不完整。(-4) |

### 4. User Stories — 178/200

| Criterion | Score | Rationale |
|-----------|-------|-----------|
| 覆盖率：每类目标用户一个 Story | 40/50 | 两类用户均有 Story: 项目开发者（Story 1-3, 5-6，共 5 个），维护者（Story 4, 7，共 2 个）。**扣分**: (1) "存量用户"（已使用 gen-test-cases 的用户）无独立 Story——退休旧路径直接影响他们。In-scope 排除了"迁移工具"，但 Background 未识别存量用户为受影响群体。(2) 维护者 Story 4 覆盖评测门禁、Story 7 覆盖场景扩展，但管线架构的清晰性（如 eval rubric 维度设计、评分阈值设定、gold standard 校准）没有对应 Story。 |
| 格式正确 (As a / I want / So that) | 43/50 | 7 个 Story 均使用 As a / I want / So that 格式。**扣分**: (1) Story 1 的 "I want to" 包含括号说明"不需要在 gen-test-cases 和 Journey-Contract 之间选择"——这是对旧状态的反面描述而非期望行为的正面描述。(2) Story 3 的 "I want to" 列举了三种场景的具体策略（subprocess 断言、浏览器自动化、Maestro YAML）——这是实现细节而非用户意图。(3) Story 6 的 "So that" 过长，包含了静态侦察的细节描述。 |
| AC per Story (Given/When/Then) | 48/50 | 每个 Story 都有 Given/When/Then AC。Story 4 有两组 AC（eval-journey + eval-contract），Story 5 有三组 AC（框架检测 + 草稿拒绝 + 内置库）。**扣分**: (1) Story 3 有两个独立的 GWT 块（CLI 和 Mobile）但未编号区分。(-1) (2) Story 4 两组 GWT 块也未编号区分。(-1) |
| AC 可验证性与边界覆盖 | 47/50 | Story 2 AC 的 1.5x 比较有退化子句（绝对值 ≥ 13）。Story 3 AC 明确了"Contract 测试占比"的计算分母："分母：生成的测试函数总数，按 `Contract 测试函数 / (Contract 测试函数 + Journey 烟测试函数)` 计算"。Story 7 AC 明确了回归验证的量化标准（函数数量变化 ≤ 5%、无新增编译错误、eval 评分偏差 ≤ 30 分）。Story 5 `diff --stat` 标记为 human-verified。**扣分**: (1) Story 1 "全局搜索 gen-test-cases 关键词（除 PRD/历史文档外）无匹配"——"历史文档"的范围判定标准不精确（eval 目录下的报告算吗？PRD 文档本身算吗？）。(-1) (2) Story 5 `diff --stat` 的 "≤ 20%" — 什么构成"一行修改"（空行？注释？格式？）未定义，不同工具的 diff 统计粒度可能不一致。(-1) (3) Story 6 "边界/异常 Outcome 占比 ≥ 30%（基线：初始静态侦察时的占比，标记为 human-verified）"——"占比"的分母是 Outcome 总数还是测试函数总数？基线度量时点不精确（第 1 轮 R2L 前？gen-contracts 完成后？）。(-1) |

### 5. Scenario Completeness — 118/150

| Criterion | Score | Rationale |
|-----------|-------|-----------|
| 端到端场景覆盖 | 48/60 | 场景检测规则从 10 条扩展到 16 条，覆盖 Go/Node.js/Python/Java/Rust 五种语言的 CLI/TUI/API/WebUI/Mobile。Per-Scenario Strategy 表格覆盖 5 种场景类型。PRD 质量前置检查确保 PRD 可用性。**扣分**: (1) Mobile "尽力而为"策略仍然模糊——"骨架"包含哪些内容？"app lifecycle" 的步骤粒度？"navigation" 的范围？"Maestro YAML 骨架"在 Per-Scenario Strategy 中为"Journey 骨架 + deep link"，但 Journey 骨架与 gen-journeys 产出的 Journey 如何关联未说明。(-4) (2) 场景检测 16 条规则中仍有"无法匹配或匹配到多个类型 → 暂停管线"的兜底——但确认后的策略未说明：用户选择单一类型后，其他类型的信号如何处理？管线后续步骤是否忽略非选定类型的信号？(-3) (3) 场景类型检测后的确认环节——流程图 SCENE_FAIL 输出信号"用户确认场景类型"，但文档未说明确认方式（自动确认？用户确认？命令行参数？配置文件？）。(-3) (4) API 场景的"平衡 50/50"——Contract 测试和 Journey 烟测试各覆盖什么范围？增量价值未说明。(-2) |
| 隐式假设暴露 | 28/40 | **已解决的隐式假设**: (1) PRD 质量前置检查解决了"假设 PRD 中存在可提取的用户故事"。(2) Fact Table runtime/static fallback 策略解决了"runtime 总是优于 static"的假设。**未暴露的隐式假设**: (1) 假设 eval 评分结果可以被 revise 技能消费——Data Flow Table 定义了格式但未说明 revise 技能如何将"不通过项明细"转化为具体的修正指令（修改哪个字段？修改方向是什么？）。(-3) (2) 假设 LLM 衍生的边界 Outcome 在语义上有意义——"LLM prompt 增强 + required_outcomes 规则"是否能产出有实际测试价值的边界场景？required_outcomes 是硬编码的场景类型规则，但实际项目的边界场景可能超出规则覆盖。(-3) (3) 假设 Convention Schema 的 4 个 section 定义对所有测试框架足够——但不同测试框架可能有不同的 section 需求（如 Maestro 的 YAML 结构与 framework/discovery/structure/assertions 模型不匹配）。(-3) (4) 假设 session 缓存（SCENE_FAIL 后写入）是可靠的——但 session 概念从未定义（生命周期？跨运行持久化？清理机制？）。(-3) |
| 业务规则一致性 | 42/50 | 主体一致。Pipeline Exit Codes 与 BIZ-error-reporting-001/002 对齐（exit 0/1/2 语义正确，每条错误含 failure reason + recovery hint）。质量门禁集成关系在第 83-87 行明确区分了 BIZ-quality-gate-001（项目源代码）和 eval 门禁（测试管线产物）。**矛盾/不一致**: (1) **Flow Description 步骤重复**: 步骤 6 (gen-contracts) 和步骤 8 (gen-contracts) 描述同一操作但内容不同（步骤 8 多了 schema 验证）。PAUSE 恢复路径重复出现。这不直接违反业务规则，但造成歧义——下游实现者无法确定哪组描述是权威版本。(-4) (2) BIZ-task-lifecycle-003 列出 `test.graduate` 为系统保留类型，PRD 提议删除——但保留类型列表的同步更新（文档或代码中移除）未在 Scope 或 Story 中提及。(-2) (3) Out of Scope "失败诊断场景特定策略"与 Run-to-Learn Failure Handling 表格内容部分重叠——表格中的 API 写操作回滚策略实质上是"场景特定策略"。(-2) |

### 6. Edge Case Coverage — 78/100

| Criterion | Score | Rationale |
|-----------|-------|-----------|
| 错误路径文档 | 30/40 | 已覆盖: gen-contracts schema 验证失败（自动重试一次 → 暂停管线），gen-test-scripts 语法验证失败（自动重试 → 标记 gen-failed），PAUSE 恢复三种路径，FIX_DECIDE 区分脚本问题/Contract 语义错误 + 安全约束，Run-to-Learn 4 种失败场景 + 兜底原则，eval-skipped 降级策略 4 步，Pipeline Exit Codes 6 种退出码。**未覆盖**: (1) Convention 草稿自动生成失败——test-guide 无法检测到任何已知测试框架时如何处理？(-3) (2) Convention 草稿拒绝 2 次重试用尽后——管线是暂停还是回退到手动 Convention？(-3) (3) run-tests 执行中测试运行器本身崩溃（非测试失败，而是执行环境异常——如 OOM、超时）。(-2) (4) FIX_DECIDE 的"不降低断言严格度"约束标记为 human-verified——这是诚实的，但意味着自动化层面没有防止"通过但无意义"测试的安全网。(-2) |
| 边界条件 | 28/35 | **已覆盖**: eval 迭代上限 3 轮、R2L 迭代 ≤ 3 轮、风险等级三级分类、FIX_DECIDE 修复 ≤ 2 次、eval 评分解析失败重试 1 次、PRD 质量前置检查最低要求、风险密度绝对下限 ≥ 13。**未覆盖**: (1) eval 评分恰好在维度阈值边界（如完整性恰好 120/200）的通过/不通过判定——阈值是 ≥ 还是 >？（维度表格写"≥ 120"，应是一致的，但 Flow Description 步骤 5 写"未达阈值"——"未达"是否包含等于？）(-2) (2) 全部 HIGH 和全部 LOW 置信度时的处理策略差异——CONFIDENCE 评级后直接进入 RUN_TESTS，无差异化行为（如全 LOW 是否需要人工确认）。(-2) (3) 场景检测规则多规则同时匹配时的优先级——16 条规则中可能有交叉（如 `Cargo.toml + clap` 同时匹配 CLI 规则和 Rust 测试规则）。(-3) |
| 失败恢复描述 | 20/25 | **已覆盖**: ENV_FAIL → 用户修复后重新检测（回路）。eval-skipped → 4 步降级策略。Run-to-Learn 失败 → 兜底原则 + R2L_DEGRADE。FIX_DECIDE → 失败类型区分 + 2 次上限 + 安全约束。PAUSE_J/PAUSE_C → 3 种恢复路径。**未覆盖**: (1) Convention 草稿重试 2 次用尽后的最终失败处理——用户仍不满意草稿，管线如何继续？(-2) (2) SCENE_FAIL 用户确认场景类型后——选择的类型可能不正确，无验证机制。(-1) (2) gen-test-scripts 语法验证重试仍失败后标记 `gen-failed`——但"跳过"不阻塞其余测试执行，如果所有测试文件都标记 gen-failed 怎么办？(-2) |

### 7. Scope Clarity — 88/100

| Criterion | Score | Rationale |
|-----------|-------|-----------|
| In-scope 项目是具体可交付物 | 30/35 | Eval Rubric 评分维度框架以 6 维度表格列出。质量门禁集成关系明确区分了 BIZ-quality-gate-001 和 eval 门禁。**模糊项**: (1) "Run-to-Learn 机制：骨架测试 → 运行捕获输出 → 丰富 Fact Table → 重新生成"仍是流程描述而非交付物清单。(-2) (2) "场景特定执行环境就绪检测（CLI/TUI/WebUI/API）"——交付物是什么？检测规则已在流程图中标注，但未作为独立交付物列出。(-1) (3) "边界/异常场景自动衍生引擎：基于场景类型 required_outcomes 规则 + 项目 Fact Table 自动生成边界和异常 Outcome"——描述清晰但仍偏向功能描述而非交付物规格。(-2) |
| Out-of-scope 明确列出延迟项 | 26/30 | 11 项 out-of-scope 清晰列出。注释（Maestro YAML vs 编译/lint、required_outcomes vs schema 变更）澄清了边界模糊问题。**扣分**: (1) "已使用 gen-test-cases 项目的迁移工具"被排除——存量用户需手动处理，但 Background 未将存量用户列为受影响群体。(-2) (2) "失败诊断场景特定策略"被排除但 Run-to-Learn Failure Handling 包含场景特定处理策略（编译错误记录、运行崩溃标记、API 写操作回滚）——边界模糊。(-2) |
| 范围与功能规格和用户故事一致 | 32/35 | 主体一致。In-scope 15 项、Functional Specs 8 个 Change Points、7 个 Story 之间映射基本完整。Eval Rubric 维度框架对应 Story 4。质量门禁集成关系对应 Change Point #6。**不一致处**: (1) Story 7（可扩展场景类型系统）的核心是"通过添加配置文件接入新场景类型"，In-scope 中最接近的项是"场景差异化：CLI/TUI/WebUI/API 核心支持 + Mobile 尽力而为"——但"差异化"和"可扩展性"是不同概念。Story 7 的配置文件格式定义未作为 In-scope 条目列出。(-3) |

---

## Phase 3: Cross-Dimension Coherence Check

**CD-1 (Goals vs. Scenario Completeness)**: Goals "高风险 Journey 平均测试数 ≥ 13"与 Per-Scenario Strategy 表格 High 总测试数 13-20 一致（下界匹配）。Story 2 AC 退化条件"≥ 13"也对齐。Iteration 1 的数学矛盾已修复。

**CD-2 (User Stories vs. Scope)**: Story 7（可扩展场景类型系统）在 In-scope 中缺少直接对应的"场景类型可扩展配置文件格式定义"交付物项——最接近的是"场景差异化"项，但概念不匹配。

**CD-3 (Flow Diagrams vs. Flow Description)**: Flow Description 存在步骤重复（步骤 6-7 vs. 步骤 8-9），Flow Diagram 只有一套 GEN_CONTRACT → EVAL_C 流程——图和文字的对齐关系因文字重复而模糊。

**CD-4 (Data Flow Table vs. Flow Description)**: Data Flow Table 新增 EVAL 结果行覆盖了 eval → revise 的数据传递。但 session 缓存（SCENE_FAIL 写入）不在 Data Flow Table 中，与 Flow Description 和流程图中 SCENE_FAIL 的描述不一致。

**CD-5 (FIX_DECIDE safety constraint)**: Flow Description 步骤 16 明确区分失败类型（脚本问题 vs. Contract 语义错误）并添加"不降低断言严格度"约束。Flow Diagram FIX_DECIDE 也有两个回退分支。两者一致。但约束标记为 human-verified，Edge Case 维度需反映此限制。

---

## Phase 4: Blindspot Hunt

### [blindspot-1] Flow Description 步骤 6-9 重复导致流程描述歧义

**Quote**: prd-spec 步骤 6 (第 150 行): "gen-contracts 从 Journey 生成 6 维度合约规范，自动衍生边界/异常 Outcome"
**Quote**: prd-spec 步骤 8 (第 157 行): "gen-contracts 从 Journey 生成 6 维度合约规范，自动衍生边界/异常 Outcome。合约生成后执行 schema 验证（6 维度结构完整性 + Outcome Preconditions 互斥性检查）；验证失败则..."

**Quote**: prd-spec 第 153-156 行和第 160-163 行: PAUSE_J/PAUSE_C 恢复路径完全相同的内容出现两次

**Problem**: 步骤 6-7 和步骤 8-9 描述了相同的操作（gen-contracts + eval-contract），但步骤 8 包含了步骤 6 没有的 schema 验证逻辑。PAUSE 恢复路径也完整重复。这导致：(a) 下游实现者无法确定哪组步骤是权威描述——步骤 6 是否也应该有 schema 验证？还是步骤 8 是步骤 6 的替代版？(b) 如果两组都有效，gen-contracts 被执行了两次，这在逻辑上不合理。(c) 步骤编号跳跃（阶段二有步骤 4-9 共 6 步，但实际操作只有 4 个：gen-journeys、eval-journey、gen-contracts、eval-contract），暗示这是编辑残留。这是文档质量问题而非功能缺陷，但对实现者的理解有重大影响。

**Must improve**: 合并步骤 6 和步骤 8（gen-contracts），保留包含 schema 验证的完整版本。合并步骤 7 和步骤 9（eval-contract）。删除重复的 PAUSE 恢复路径，保留一处即可。调整步骤编号使阶段二为步骤 4-7（4 步）。

### [blindspot-2] Mobile 场景与 test-guide Convention 检测断裂

**Quote**: prd-spec 第 74 行 In-scope: "场景差异化：CLI/TUI/WebUI/API 核心支持 + Mobile 尽力而为"
**Quote**: prd-spec 第 75 行 In-scope: "内置 Convention 文件扩充（pytest、JUnit、Rust/cargo test）"
**Quote**: prd-spec 场景检测规则表第 128 行: "AndroidManifest.xml 或 *.xcodeproj + UI 框架依赖 → Mobile"
**Quote**: prd-spec Per-Scenario Strategy 表格 Mobile 列: "Journey 骨架 + deep link"

**Problem**: Mobile 场景需要 Maestro 测试框架，但 Convention 扩充计划（pytest/JUnit/Rust）不包含 Maestro。场景检测规则表能检测到 Mobile 项目（AndroidManifest.xml 等），但 test-guide 的自动框架检测如何处理 Mobile 项目？如果检测不到 Maestro Convention，test-guide 会无法生成合适的草稿。如果 Mobile 跳过 Convention 检测步骤，则 test-guide 的通用性（Story 5 "在一个全新项目中首次运行测试生成时... 自动检测测试框架并生成 Convention 草稿"）在 Mobile 场景下失效。PRD 未说明这一断裂。In-scope 列出了"Mobile 尽力而为"但"尽力而为"不包含 Convention 支持——这意味着 Mobile 项目无法走完整管线流程。

**Must improve**: 明确 Mobile 场景的 Convention 策略：(a) 内置 Maestro Convention 草稿，或 (b) Mobile 场景跳过 Convention 检测步骤直接进入 gen-journeys，或 (c) 在 Out-of-scope 中明确说明"Mobile 场景的 Convention 自动生成"。

### [blindspot-3] session 缓存概念在 PRD 中无定义

**Quote**: Flow Diagram SCENE_FAIL 箭头标注: "用户选择类型 / 选择存入 session 缓存"
**Quote**: Data Flow Table SCENE_DETECT 行: "写入 `.forge/session.yaml` 的 `scene_type` 字段"

**Problem**: Data Flow Table 提到了 `.forge/session.yaml` 作为 scene_type 的传递方式，但 Flow Diagram SCENE_FAIL 说的"session 缓存"暗示了一种缓存机制——用户选择后存入缓存，下次运行时自动使用。然而：(a) "session" 的生命周期未定义——是单次运行有效？跨运行持久化？过期时间？(b) "缓存"暗示读取优先级——如果缓存存在，是否跳过 SCENE_DETECT？(c) 用户选择场景类型后"存入 session 缓存"和 SCENE_DETECT "写入 session.yaml" 是同一操作还是两个不同操作？(d) 缓存何时失效——项目结构变更后？版本升级后？这个概念在 PRD 中只在流程图箭头标注中出现一次，无任何其他解释。

**Must improve**: 在 Core Concepts 或 Flow Description 中定义 session 缓存的生命周期、作用域、过期机制和读取优先级。或者将"session 缓存"简化为"写入 session.yaml"以避免引入未定义概念。

### [blindspot-4] eval rubric "合理推理声明"的可验证性标准不足

**Quote**: prd-spec 第 68 行 Eval Rubric 维度定义: "合理推理声明（LLM 衍生的边界 Outcome 不在 Fact Table 中，但基于场景类型 required_outcomes 规则衍生，需标注推理依据和 source: inferred）。评分标准：推理声明需有 required_outcomes 规则支撑"

**Quote**: prd-spec Per-Scenario Strategy 必须衍生边界 Outcome: "CLI: not-found + already-exists; TUI: timeout; WebUI: validation-error + session-expired; API: unauthorized"

**Problem**: 评分标准"推理声明需有 required_outcomes 规则支撑"意味着只要声明的 Outcome 名称在 required_outcomes 列表中就能通过——但这只验证了名称匹配，不验证声明的具体内容是否合理。例如，CLI 场景的 `not-found` Outcome 声明内容可能是"当文件不存在时返回 404"（合理），也可能是"当用户输入错误命令时系统崩溃"（不合理但名称匹配 required_outcomes）。eval-contract 需要评估的是声明的语义正确性，而不仅仅是名称匹配。当前评分标准过于宽松，可能导致名称匹配但语义错误的 Outcome 通过 eval 门禁。

**Must improve**: 在 eval-contract rubric 的"合理推理声明"评分标准中增加语义验证要求——不仅检查 Outcome 名称是否在 required_outcomes 列表中，还需验证声明的具体内容（Preconditions、Expected Behavior）是否符合场景类型的语义预期。

---

## Score Summary

| Dimension | Score | Max |
|-----------|-------|-----|
| 1. Background & Goals | 85 | 100 |
| 2. Flow Diagrams | 135 | 150 |
| 3. Flow Completeness | 155 | 200 |
| 4. User Stories | 178 | 200 |
| 5. Scenario Completeness | 118 | 150 |
| 6. Edge Case Coverage | 78 | 100 |
| 7. Scope Clarity | 88 | 100 |
| **Total** | **837** | **1000** |

---

## Attacks Summary

1. [Background & Goals]: "v3.0.0 是重构的最佳窗口"在 Why 中提出但 Goals 和 Delivery Phasing 无版本约束 — "正处于大版本分支上，可以做大范围变更而不破坏已发布版本" — 如果版本窗口是关键动机，应在 Goals 或 Delivery Phasing 中绑定版本约束
2. [Background & Goals]: "消除双路径困惑"metric "所有相关文件完全删除"范围边界不清晰 — "gen-test-cases 及所有相关文件完全删除" — 定义"相关文件"的精确范围（包含 eval 历史报告？PRD 中的引用？）
3. [Flow Diagrams]: SCENE_FAIL "session 缓存"概念未定义 — 流程图标注"选择存入 session 缓存" — 在 Core Concepts 或 Data Flow Table 中定义 session 缓存的生命周期和语义
4. [Flow Diagrams]: R2L "覆盖率达标"阈值未在图中标注 — "≤ 3 轮或覆盖率达标" — 明确"达标"的具体覆盖率阈值（如 ≥ 80%）
5. [Flow Completeness]: **Flow Description 步骤 6-7 和步骤 8-9 重复** — 步骤 6 (gen-contracts) 和步骤 8 (gen-contracts) 描述同一操作但内容不同 — 合并为单组步骤，删除重复的 PAUSE 恢复路径
6. [Flow Completeness]: Convention 文件在流程中的角色不清 — TEST_GUIDE → GEN_JOURNEY 箭头暗示 Convention 在 GEN_JOURNEY 前被消费 — 澄清 Convention 的消费时点（gen-test-scripts）并在流程图中体现
7. [User Stories]: Story 7 在 Scope 中缺少对应交付物 — Story 7 核心是"通过添加配置文件接入新场景类型" — 在 In-scope 中添加"场景类型配置文件格式定义"交付物
8. [Scenario Completeness]: Mobile 场景与 test-guide Convention 检测断裂 — Convention 扩充不含 Maestro — 明确 Mobile 场景的 Convention 策略
9. [Edge Case Coverage]: Convention 草稿重试 2 次用尽后无后续处理 — Story 5 "最多重试 2 次" — 定义 2 次后管线行为（暂停？回退？使用空白 Convention？）
10. [Edge Case Coverage]: FIX_DECIDE "不降低断言严格度"仅 human-verified 无自动化检测 — 缺少自动化安全网意味着 agent 可能绕过此约束 — 考虑至少定义可自动检测的最低标准（如断言数量不减少）
11. [blindspot]: eval "合理推理声明"评分标准过于宽松 — 只检查 Outcome 名称匹配 required_outcomes — 增加语义验证要求，确保声明内容符合场景类型预期
12. [blindspot]: session 缓存概念在 PRD 中无定义 — 仅在流程图箭头标注中出现一次 — 定义生命周期、作用域、过期机制，或简化为"写入 session.yaml"
