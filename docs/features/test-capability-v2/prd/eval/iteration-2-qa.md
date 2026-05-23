# QA Evaluation Report — Test Capability v2.0 PRD

**Evaluator**: Senior QA Engineer (Adversarial)
**Iteration**: 2
**Mode**: B (No UI — prd-ui-functions.md absent)
**Date**: 2026-05-23

---

## Phase 1: Reasoning Audit (Pre-Score Anchors)

Before applying the rubric, I traced the document's core argument chain independently:

**P1. Problem → Solution**: The three structural defects (dual-path confusion, shallow testing, limited generality) map well to the five solution pillars. Iteration 2 improved the mapping: the "limited generality" problem now explicitly cites "评测门禁缺失" and "Mobile 接入成本极高" as sub-problems, justifying the "eval补全" and "Mobile 骨架" goals. However, test-guide auto-generation of Convention drafts remains a qualitatively more ambitious capability than the problem scope (only 3 frameworks lacking support) would suggest.

**P2. Solution → Evidence**: Goals metrics improved significantly. The risk-driven density Goal now uses "≥ 13" absolute floor AND "≥ 1.5x" ratio with a degradation clause for single-risk journeys. Eval gate Goal clarified as "评分 ≥ 850/1000" (score threshold, not accuracy) with gold standard calibration method defined in Other Notes. The remaining weakness is "降低 Mobile 接入成本" — the metric ("生成 Maestro YAML 骨架 + deep link 测试") is deliverable-based, not cost-based.

**P3. Evidence → Success Criteria**: Delivery phasing gates improved ("高风险测试数 ≥ 低风险 × 1.5，边界 Outcome 无效比例 < 20%"). The new Pipeline Exit Codes table provides concrete failure state definitions. However, success criteria still primarily test pipeline execution completeness rather than generated test quality.

**P4. Self-contradiction check**: Iteration 2 resolved several iteration-1 contradictions: risk density table lower bound now aligns with Goal (13-20 for High, floor matches ≥ 13); Flow Diagram FIX_DECIDE now has two rollback paths matching Flow Description text; Data Flow Table added resolving data propagation gap. New potential contradiction: eval-contract "事实依据" dimension (requires Fact Table backing) vs. LLM-derived boundary Outcomes (not in Fact Table by definition).

These anchors inform blindspot attacks where they identify issues outside rubric dimensions.

---

## Previous-Attack Resolution Status (Iteration 1 → 2)

| # | Attack | Status | Evidence |
|---|--------|--------|----------|
| 1 | Goal "评分准确率 >= 850/1000" 用词不当且 ground truth 未定义 | **RESOLVED** | prd-spec Goals: "eval-journey/eval-contract 评分 ≥ 850/1000（基于 gold standard 评分集校准）" — now clearly a threshold; gold standard calibration method defined in Other Notes |
| 2 | Mobile 接入成本降低与 Background 三大缺陷无逻辑联系 | **PARTIALLY RESOLVED** | Background Why #3 now explicitly mentions "Mobile 场景无任何测试生成支持，接入成本极高"; but Goal metric still deliverable-based not cost-based |
| 3 | Flow Diagram FIX_DECIDE 缺少 Contract 语义错误回退路径 | **RESOLVED** | Flow Diagram FIX_DECIDE now has two branches: "是: 脚本问题 ≤ 2 次 → GEN_SCRIPTS" and "是: Contract 语义错误 ≤ 2 次 → GEN_CONTRACT" |
| 4 | Run-to-Learn 骨架测试执行失败分支完全缺失 | **RESOLVED** | Flow Diagram added R2L_DEGRADE node; Flow Description added "Run-to-Learn Failure Handling" table with 4 failure scenarios |
| 5 | 数据流传递路径未文档化 | **RESOLVED** | New "Data Flow Table" with 7 rows covering scene type, Convention, Journey, Contract, Fact Table, confidence rating |
| 6 | Story 2 AC 的 1.5x 比较在只有高风险 Journey 时不可验证 | **RESOLVED** | Story 2 AC now includes degradation clause: "若同一功能仅有高风险 Journey（无低风险 variant 可比较），退化为绝对值验证：高风险 Journey 总 Outcome 数 ≥ 13" |
| 7 | Story 5 缺少 When 条件的 Given/Then 块 | **RESOLVED** | Story 5 third GWT block now has When condition: "When 用户查看 Forge 插件的 `conventions/` 目录" |
| 8 | 风险驱动密度表下界数据与 Goal 1.5x 指标数学矛盾 | **RESOLVED** | Density table now shows High 13-20, Low 4-7; floor 13/7 = 1.86 > 1.5; 13 meets absolute floor |
| 9 | gen-journeys 提取失败（PRD 无用户故事）未覆盖 | **RESOLVED** | Flow Description pre-condition now includes "PRD 质量前置检查：PRD 必须包含至少 1 个 User Story... 若 PRD 不存在或质量前置检查未通过，管线在步骤 1 报错并输出缺失项明细" |
| 10 | eval-skipped 绕过质量门禁的降级后果未评估 | **RESOLVED** | eval-skipped degradation strategy now defined with 4 steps: (1) downstream continues, (2) test file marked eval-skipped + LOW, (3) test report lists eval-skipped items, (4) user can manually clear mark |
| 11 | "LLM prompt 增强策略"作为 in-scope 交付物是实现手段而非功能 | **PARTIALLY RESOLVED** | In-scope now says "边界/异常场景自动衍生引擎" instead of "LLM prompt 增强策略", but description still mentions "LLM prompt 增强 + required_outcomes 规则" |
| 12 | 场景类型检测结果传递机制未定义 | **RESOLVED** | Data Flow Table: "SCENE_DETECT → scene_type → 写入 `.forge/session.yaml` 的 `scene_type` 字段" |
| 13 | eval-skipped 降级路径绕过质量门禁核心目的 | **RESOLVED** | eval-skipped now has explicit 4-step degradation strategy |
| 14 | LLM 衍生边界 Outcome 与 eval-contract "事实依据"维度的架构矛盾 | **NOT RESOLVED** | eval-contract rubric "事实依据" dimension says "声明是否基于 Fact Table 已知事实，未知来源是否标注 UNKNOWN" — LLM-derived boundary Outcomes (not-found, already-exists, etc.) are not in Fact Table. This contradiction persists. |
| 15 | Mobile Maestro YAML "骨架"定义完全空白 | **NOT RESOLVED** | "骨架" definition remains absent; no minimum structure example, no complexity criteria |

**Resolution Summary**: 10/15 fully resolved, 2/15 partially resolved, 3/15 not resolved.

---

## Phase 2: Rubric Scoring

### 1. Background & Goals — 82/100

| Criterion | Score | Rationale |
|-----------|-------|-----------|
| Background 三要素 (Reason/Target/Users) | 25/30 | Why/Target/Who 三段齐全。Iteration 2 在 Why #3 中补充了"评测门禁缺失"和"Mobile 接入成本极高"两个子问题，使 Target 的 5 个升级方向有了更完整的 backing。定位声明"管线只生成开发者手动编写成本高的复杂测试"澄清了范围。**扣分**: (1) Who 部分仍遗漏"存量用户"——退休旧路径直接影响已使用 gen-test-cases 的用户，但他们未被列为受影响群体。(2) Target 有 5 个方向但 Why 只有 3 大缺陷，"评测补全"和"信息增强"在 Why 中只是隐含在缺陷 #3 和 #2 中，未作为独立缺陷列出。 |
| Goals 量化 | 25/30 | 6 个 Goal 均有 Metric 列。Iteration 2 改进: "建立评测门禁"Goal 改为"评分 ≥ 850/1000（基于 gold standard 评分集校准）"，用词从"准确率"改为"评分"，gold standard 校准方法在 Other Notes 中定义。风险密度 Goal 改为"高风险 Journey 平均测试数 ≥ 13，且高风险旅程测试数 ≥ 低风险旅程 × 1.5"，双重约束更严格。**扣分**: (1) "降低 Mobile 接入成本"的 metric "生成 Maestro YAML 骨架 + deep link 测试"是交付物列表而非量化指标——无法验证"成本降低"的程度。应有对照度量（如从零到可审阅测试的步骤数）。(2) "消除双路径困惑"的 metric "gen-test-cases 及所有相关文件完全删除"——"所有相关文件"的范围边界不清晰（历史 eval 报告中的引用是否需要删除？PRD 本身的引用呢？）。 |
| 背景与目标逻辑一致性 | 32/40 | 三大问题与 Goals 大致对应，iteration 2 改善了逻辑链。**扣分**: (1) Goal "降低 Mobile 接入成本"与 Background 的联系虽然建立了（Background #3 提到"接入成本极高"），但 Goal 的 metric 没有度量"成本降低"——只是列出了交付物。(2) "v3.0.0 是重构的最佳窗口"出现在 Why 中但从未在 Goal 中体现为交付时间约束或版本绑定。(3) Goal "提升测试信息质量"的覆盖率基线在 Story 6 AC 中定义为"gen-contracts 静态侦察结果"——但 Goals 表格未引用此定义，不自包含。(4) "建立评测门禁"Goal 引用了"gold standard 评分集校准"——校准方法的定义在 Other Notes 中，但校准数据集的构建时间点（在交付阶段二之前？之后？）和维护责任未明确。 |

### 2. Flow Diagrams — 130/150

| Criterion | Score | Rationale |
|-----------|-------|-----------|
| Mermaid 图存在 | 50/50 | 完整的 mermaid flowchart，使用 START/END、决策菱形、处理节点，覆盖四个阶段。 |
| 主路径完整 (start → end) | 40/50 | START → SCENE_DETECT → DETECT → GEN_JOURNEY → EVAL_J → GEN_CONTRACT → EVAL_C → GEN_SCRIPTS → R2L_CHOICE → ENV_CHECK → CONFIDENCE → RUN_TESTS → REPORT → END。Iteration 2 改进: FIX_DECIDE 现在有两个回退分支（GEN_SCRIPTS 和 GEN_CONTRACT），与 Flow Description 文字一致。**扣分**: (1) SCENE_FAIL 的箭头标注"选择存入 session 缓存"回到 DETECT——但 DETECT 只检查 Convention 文件存在性，场景类型信息如何传递到 GEN_JOURNEY 及后续步骤？Data Flow Table 说"写入 `.forge/session.yaml` 的 `scene_type` 字段"，但流程图中未体现这个写入操作。(2) TEST_GUIDE 到 GEN_JOURNEY 的箭头标注"用户审核确认"——但 test-guide 检测到的框架信息和 DETECT 检测到的 Convention 存在性是两个独立路径，信息汇聚点不清晰。(3) R2L_CHOICE 的退出条件"≤ 3 轮或覆盖率达标"——覆盖率达标阈值未在流程图中标注。 |
| 决策点 + 错误分支 | 40/50 | 评测门禁迭代（REVISE_J/REVISE_C）和 3 轮耗尽（PAUSE_J/PAUSE_C）覆盖良好。Iteration 2 新增: eval-skipped 降级路径（LLM 解析失败 → eval-skipped → LOW）、R2L 失败降级路径（R2L_DEGRADE）。**扣分**: (1) PAUSE_J/PAUSE_C 后用户决定的分支仍未画出——它们连接到 END 但 Flow Description 说"用户决定后续操作"，暗示可能有继续选项。(2) Run-to-Learn Failure Handling 表格定义了 4 种失败场景，但流程图中只有 R2L_DEGRADE 一个节点，没有区分不同失败类型的处理。(3) CONFIDENCE → RUN_TESTS 没有决策点——即使全部 LOW 置信度也直接执行测试，缺少"是否需要人工确认"的判断。 |

### 3. Flow Completeness — 150/200

| Criterion | Score | Rationale |
|-----------|-------|-----------|
| 流程步骤描述完整业务过程 | 55/70 | 四阶段 13 步覆盖完整生命周期。Iteration 2 新增: 场景类型检测规则表格（10 条规则含兜底），风险分级判定规则（High/Medium/Low 三级定义），PRD 质量前置检查，Pipeline Exit Codes 表格。**扣分**: (1) 场景类型检测表格有 10 条规则但缺少优先级——多规则同时匹配时如何判定？表格最后一条"无法匹配或匹配到多个类型 → 暂停管线"提供了兜底，但在实际项目中 CLI+API 混合（如 Go CLI 工具附带 HTTP API）是常见场景，不应仅靠兜底处理。(2) 步骤 13 输出报告的格式、输出位置、内容结构未说明。(3) 步骤 14 的两种失败类型区分方式仍不明确——靠错误码？靠 LLM 分析？靠测试输出模式匹配？ |
| 数据流文档 | 55/70 | Iteration 2 新增 Data Flow Table，覆盖 7 行数据传递路径（场景类型、Convention、Journey、Contract、静态 Fact Table、运行时 Fact Table、置信度评级），包括源步骤、产出数据、消费步骤、传递方式四列。这是 iteration 1 重大缺失的重大改进。**扣分**: (1) eval 评分结果的数据结构未说明——下游 revise 技能如何消费评分反馈来修正 Journey/Contract？评分结果包含哪些字段（分数、维度明细、具体扣分项）？(2) Convention 草稿被拒绝后的"基于用户反馈重新生成"循环在 Data Flow 中没有对应行——用户反馈的数据格式是什么？(3) PRD 文档作为 GEN_JOURNEY 的输入不在 Data Flow Table 中。 |
| 异常处理与边界情况 | 40/60 | Iteration 2 重大改进: Run-to-Learn Failure Handling 表格（4 种失败场景含检测方式和处理策略），eval-skipped 降级策略（4 步处理），Pipeline Exit Codes 表格（6 种退出码），PRD 质量前置检查。**未覆盖**: (1) gen-contracts 合约 schema 验证失败的处理。(2) gen-test-scripts 代码生成失败（语法错误、框架不兼容）。(3) Convention 草稿自动生成失败（无法检测到任何已知测试框架）。(4) PAUSE_J/PAUSE_C 后的恢复流程未定义——用户可以选择什么（继续？放弃？修改后重跑？）。(5) run-tests 执行中测试运行器本身崩溃（非测试失败，而是执行环境异常）。 |

### 4. User Stories — 160/200

| Criterion | Score | Rationale |
|-----------|-------|-----------|
| 覆盖率：每类目标用户一个 Story | 40/50 | 两类用户均有 Story: 项目开发者（Story 1-3, 5-6，共 5 个），维护者（Story 4, 7，共 2 个）。Iteration 2 新增 Story 7（可扩展场景类型系统）覆盖维护者扩展需求。**扣分**: (1) "存量用户"（已使用 gen-test-cases 的用户）无独立 Story——退休旧路径直接影响他们，但 Story 1 以"管线统一"角度描述而非"迁移"角度。In-scope 明确排除了"已使用 gen-test-cases 项目的迁移工具"，但存量用户在 Background 中未被识别为受影响群体。(2) 维护者 Story 4 只覆盖评测门禁、Story 7 只覆盖场景扩展，但 Background 说维护者需要"清晰的管线架构"——管线架构的清晰性（如 eval rubric 维度设计、评分阈值设定）没有对应 Story。 |
| 格式正确 (As a / I want / So that) | 45/50 | 7 个 Story 均使用 As a / I want / So that 格式。**扣分**: (1) Story 3 的 "I want to" 列举了 CLI/WebUI/Mobile 三种场景的具体策略——这是实现细节而非用户意图。(2) Story 1 的 "I want to" 包含括号说明"不需要在 gen-test-cases 和 Journey-Contract 之间选择"——这是对旧状态的反面描述而非期望行为的正面描述。 |
| AC per Story (Given/When/Then) | 40/50 | 每个 Story 都有 Given/When/Then AC。**扣分**: (1) Story 3 有两个独立的 Given/When/Then 块（CLI 和 Mobile）但未编号区分。(2) Story 5 有三个 Given/When/Then 块，虽然第三个块现在有 When 条件（"When 用户查看 Forge 插件的 `conventions/` 目录"），但这更像是陈述句而非触发动作。(3) Story 6 的 AC 中 `And` 子句过多（4 个 And），模糊了核心 Then 断言与补充断言的层级。(4) Story 4 有两个 Given/When/Then 块，分别覆盖 eval-journey 和 eval-contract，但未编号区分。 |
| AC 可验证性与边界覆盖 | 35/50 | **可验证性问题**: (1) Story 2 AC 的 1.5x 比较现在有退化子句（"若同一功能仅有高风险 Journey，退化为绝对值验证：高风险 Journey 总 Outcome 数 ≥ 13"），这是重大改进。但 Story 2 AC 还包含 "API 场景的每个认证端点自动衍生 unauthorized Outcome；TUI 场景的每个异步 Cmd 自动衍生 timeout Outcome"等——这些"自动衍生"如何验证？是否检查生成的 Outcome 列表中确实包含了 `unauthorized`/`timeout`？如果是，AC 应明确验证方式。(2) Story 5 "diff --stat 统计用户修改行数占草稿总行数的比例；目标 ≤ 20%"——human-verified 标注诚实，但"修改是否因为草稿质量差"无法从 diff --stat 判断。(3) Story 7 "已有场景类型的测试生成结果不受影响（回归验证）"——"不受影响"的验证标准仍为模糊（完全相同？语义等价？允许 LLM 随机差异？）。(4) Story 1 "全局搜索 gen-test-cases 关键词（除 PRD/历史文档外）无匹配"——"历史文档"的范围判定标准不精确（eval 目录下的报告算吗？）。 |

### 5. Scenario Completeness — 105/150

| Criterion | Score | Rationale |
|-----------|-------|-----------|
| 端到端场景覆盖 | 45/60 | Iteration 2 新增: 场景检测规则表格（10 条），风险分级判定规则（三级定义），PRD 质量前置检查。Per-Scenario Strategy 表格覆盖 5 种场景类型。TUI 的 ENV_CHECK 在流程图中标注了"二进制+stdin pipe"。**扣分**: (1) Mobile "尽力而为"策略仍然模糊——"骨架"的具体定义（Maestro YAML 最小结构？包含哪些操作？）和"复杂场景标记 manual-only"的判定标准未补充。(2) 场景检测规则表格 10 条规则中 CLI 和 TUI 的检测全部依赖 `main.go`（Go 项目信号），Python/Node.js CLI 工具的检测规则缺失。(3) API 场景的"平衡 50/50"——Contract 测试和 Journey 烟测试各覆盖什么范围？增量价值未说明。(4) 场景类型检测后的"确认"环节如何与用户交互未说明（自动确认？用户确认？）。 |
| 隐式假设暴露 | 22/40 | **已解决的隐式假设**: (1) PRD 质量前置检查解决了"假设 PRD 中存在可提取的用户故事"。**未暴露的隐式假设**: (1) 假设 eval 评分结果可以被 revise 技能消费——但评分结果到修正指令的映射机制未说明。(2) 假设 LLM 衍生的边界 Outcome 在语义上有意义——"LLM prompt 增强"是否能产出有实际测试价值的边界场景？(3) 假设骨架测试可以在隔离环境中运行——但 API Write 端点的副作用处理暗示测试环境可能有状态污染。(4) 假设 Fact Table 的 runtime 替换 static 策略总是提升信息质量——但 runtime fact 的 confidence 可能是 assumed（LLM 生成），替换 static 的 inferred 可能反而降低质量。(5) 假设 Convention Schema 的 4 个 section 定义是稳定不变的——但不同测试框架可能有不同的 section 需求。 |
| 业务规则一致性 | 38/50 | 主体一致。Iteration 2 新增 Pipeline Exit Codes 表格与 BIZ-error-reporting-001/002 对齐（exit 1 = retryable, exit 2 = blocking, exit 0 = success）。**矛盾/不一致**: (1) BIZ-quality-gate-001 定义了三阶段质量门禁（compile → test → e2e），PRD 的"质量门禁更新以反映新管线"是 in-scope 但未说明如何更新——是修改质量门禁技能还是修改配置？(2) BIZ-task-lifecycle-003 列出 13 个保留类型，PRD 提议删除 `test.graduate`——但保留类型列表的同步更新（文档或代码）未在 Scope 中提及。(3) PRD 在 Scope 中新增了 Eval Rubric 评分维度框架，但"一致性"维度（0-150）与 Story 4 AC 的维度阈值表（一致性 ≥ 90）对齐——框架定义一致。 |

### 6. Edge Case Coverage — 63/100

| Criterion | Score | Rationale |
|-----------|-------|-----------|
| 错误路径文档 | 25/40 | Iteration 2 重大改进: Run-to-Learn Failure Handling 表格（4 种失败场景含检测方式和处理策略），eval-skipped 降级策略（4 步处理），Pipeline Exit Codes 表格（6 种退出码），PRD 质量前置检查。**未覆盖**: (1) gen-contracts 合约 schema 验证失败的处理。(2) gen-test-scripts 代码生成失败（语法错误、框架不兼容）。(3) Convention 草稿自动生成失败（无法检测到任何已知测试框架）。(4) run-tests 执行中测试运行器本身崩溃。(5) gen-test-scripts 输出的质量验证缺失——生成的代码没有语法/可执行性检查就直接进入 Run-to-Learn 或 run-tests。 |
| 边界条件 | 20/35 | **已覆盖**: 评测迭代上限 3 轮、Run-to-Learn 迭代 ≤ 3 轮、风险等级三级分类、FIX_DECIDE 修复 ≤ 2 次、eval 评分解析失败重试 1 次、PRD 质量前置检查最低要求。**未覆盖**: (1) eval 评分恰好在维度阈值边界（如完整性恰好 120/200）的通过/不通过判定——阈值是 ≥ 还是 >？(2) 0 个 Journey 的极端情况（PRD 存在且通过质量前置检查但无用户故事可提取——前置检查要求"至少 1 个 User Story"所以理论上不会发生，但前置检查本身的"User Story"定义可能不够严格）。(3) 单个 Journey 包含大量 Step（如 50+）时的性能和超时影响。(4) 全部 HIGH 和全部 LOW 置信度的处理策略差异。(5) Convention 文件格式损坏或版本不兼容时的容错。(6) 场景检测规则多规则同时匹配（非兜底的"混合类型"场景，而是规则优先级冲突）。 |
| 失败恢复描述 | 18/25 | **已覆盖**: ENV_FAIL → 用户修复后重新检测（图中有回路）。eval-skipped → 4 步降级策略（标记、报告、用户审核、手动清除）。Run-to-Learn 失败 → 兜底原则 + R2L_DEGRADE。FIX_DECIDE → 修复耗尽后输出报告。Pipeline Exit Codes 表格明确了每种终止点的语义。**未覆盖**: (1) PAUSE_J/PAUSE_C 后的恢复流程——用户选择"继续"意味着什么？接受当前低分文档？手动修改后重新评分？从中断点继续？(2) Convention 草稿被拒绝后"基于用户反馈重新生成"——"用户反馈"的数据格式和"最多重试 2 次"后的最终失败处理未定义。(3) SCENE_FAIL 用户确认场景类型后的验证——用户选择的类型可能不正确。 |

### 7. Scope Clarity — 83/100

| Criterion | Score | Rationale |
|-----------|-------|-----------|
| In-scope 项目是具体可交付物 | 28/35 | Iteration 2 改进: "合约规范增强"重新表述为"边界/异常场景自动衍生引擎"。Eval Rubric 评分维度框架作为 in-scope 项列出了 6 维度表格。**模糊项**: (1) "Run-to-Learn 机制：骨架测试 → 运行捕获输出 → 丰富 Fact Table → 重新生成"仍是流程描述而非交付物清单。(2) "质量门禁更新以反映新管线：将现有单一门禁替换为多阶段门禁"——更新什么技能/配置？更新后的行为与 BIZ-quality-gate-001 的关系？(3) "场景特定执行环境就绪检测（CLI/TUI/WebUI/API）"——交付物是什么？配置文件？检测脚本？检测规则已在流程图中标注但未作为独立交付物列出。 |
| Out-of-scope 明确列出延迟项 | 25/30 | 11 项 out-of-scope 清晰列出。两个注释（Maestro YAML vs 编译/lint、required_outcomes vs schema 变更）澄清了边界模糊问题。**扣分**: (1) "已使用 gen-test-cases 项目的迁移工具"被排除——存量用户需手动处理，但 Background 未将存量用户列为受影响群体。(2) "执行环境自动准备与配置（仅做就绪检测）"的边界——ENV_FAIL 后修复建议的深度未定义（只说"缺少 X"还是提供安装命令？）。(3) "失败诊断场景特定策略"被排除但 FIX_DECIDE 区分"脚本问题"和"Contract 语义错误"——这本身是一种诊断策略，边界模糊。 |
| 范围与功能规格和用户故事一致 | 30/35 | 主体一致。In-scope 15 项、Functional Specs 8 个 Change Points、7 个 Story 之间映射基本完整。Iteration 2 新增 Eval Rubric 维度框架对应 Story 4。**不一致处**: (1) Story 7（可扩展场景类型系统）的核心是"通过添加配置文件接入新场景类型"，In-scope 中最接近的项是"场景差异化：CLI/TUI/WebUI/API 核心支持 + Mobile 尽力而为"——但"差异化"和"可扩展性"是不同概念。(2) Story 6 的置信度评级（HIGH/MEDIUM/LOW）在 In-scope 中有"置信度评级系统"对应项，也在 Change Point #6 中提及，但 Change Point #6 合并了"环境就绪检测 + 置信度评级"——是否应拆分？ |

---

## Phase 3: Cross-Dimension Coherence Check

**CD-1 (Goals vs. Scenario Completeness)**: Goals "高风险 Journey 平均测试数 ≥ 13"与 Per-Scenario Strategy 表格 High 总测试数 13-20 一致（下界匹配）。Story 2 AC 退化条件"≥ 13"也对齐。Iteration 1 的数学矛盾已修复。

**CD-2 (User Stories vs. Scope)**: Story 7（可扩展场景类型系统）在 In-scope 中缺少直接对应的"场景类型可扩展配置"交付物项——最接近的是"场景差异化"项，但概念不完全匹配。

**CD-3 (Flow Diagrams vs. Flow Description)**: FIX_DECIDE 两个回退路径在图和文字中均一致。Iteration 1 的 CD-3 已修复。

**CD-4 (Flow Completeness vs. Edge Case Coverage)**: Flow Description 的 eval-skipped 降级策略（4 步处理）与 Edge Case 的 eval-skipped 覆盖一致。Pipeline Exit Codes 与 Flow Description 的终止点描述对齐。

**CD-5 (Data Flow Table vs. Flow Description)**: Data Flow Table 覆盖了 Flow Description 中大部分数据传递路径，但 eval 评分结果→revise 技能的传递路径缺失——这在 Flow Description 的 eval 迭代中是关键数据流。

---

## Phase 4: Blindspot Hunt

### [blindspot-1] 场景检测规则的覆盖不平衡——非 Go 项目的 CLI/TUI 检测缺失

**Quote**: 场景检测规则表格: "main.go + cobra.Command / urfave/cli → CLI", "main.go + tea.Program / tview.Application → TUI", "pyproject.toml + pytest → API"

**Problem**: 表格 10 条规则中 CLI 和 TUI 的检测全部依赖 `main.go`（Go 项目信号）。Python CLI 工具（Click/Typer）、Node.js CLI 工具（Commander/Yargs/Ink）、Rust CLI 工具（clap）完全不在表格中。同时，`pyproject.toml + pytest` 被归类为 API——如果一个 Python 项目使用 Click 构建 CLI 工具，检测结果会错误地选择 API 策略。这意味着非 Go 项目的 CLI/TUI 开发者在使用 Forge 时，管线会错误地选择测试策略。

**Must improve**: 为 CLI/TUI 检测补充非 Go 语言的信号规则，或在表格中明确声明"当前仅支持 Go 项目的 CLI/TUI 场景检测"。

### [blindspot-2] eval-contract "事实依据"维度与 LLM 衍生边界 Outcome 的架构矛盾

**Quote**: Eval Rubric 评分维度框架: "事实依据（Fact Alignment） 0–150 | 声明是否基于 Fact Table 已知事实，未知来源是否标注 UNKNOWN"
**Quote**: Per-Scenario Strategy 必须衍生的边界 Outcome: "CLI: not-found + already-exists; TUI: timeout; WebUI: validation-error + session-expired; API: unauthorized"

**Problem**: 管线的衍生引擎要求为每个场景类型生成特定的边界 Outcome（如 CLI 的 `not-found`）。这些 Outcome 是基于场景类型规则生成的，其声明内容（如"当资源不存在时返回 404"）可能不在 Fact Table 中——因为 Fact Table 记录的是被测系统的已知事实，而"资源不存在"是一个测试预设条件，不是系统行为事实。如果 eval-contract 严格执行"事实依据"维度，要求所有声明基于 Fact Table，则 LLM 衍生的边界 Outcome 会被标记为 UNKNOWN，降低评分。这意味着 eval 门禁可能与管线自身的衍生引擎产生矛盾——生成引擎产出的合理 Outcome 反而被 eval 拒绝。这是迭代 1 的 blindspot-4，在迭代 2 中仍未解决。

**Must improve**: 在 eval-contract rubric 的"事实依据"维度中明确区分"有事实依据的 Outcome"和"基于合理推理的 Outcome"两类声明。对于管线要求的必须衍生 Outcome（required_outcomes），eval 应评估其推理合理性而非要求 Fact Table 事实依据。

### [blindspot-3] PAUSE_J/PAUSE_C 后的恢复路径定义缺失

**Quote**: Flow Description: "若 eval 评分因 LLM 输出无法解析而失败... 3 轮后仍未达阈值（总分或任一维度），暂停管线并输出当前评分 + 未通过项明细，由用户决定后续操作"
**Quote**: Pipeline Exit Codes: "PAUSE_J / PAUSE_C（eval 3 轮耗尽） | 1 | retryable — 用户决策后可继续"

**Problem**: Pipeline Exit Codes 说"retryable — 用户决策后可继续"，但 PRD 从未定义用户可以选择什么恢复路径。至少有四种合理选项：(1) 接受当前低分文档继续下游步骤；(2) 放弃整个管线；(3) 用户手动修改 Journey/Contract 后重新评分；(4) 调整 eval 阈值后重新评分。这四种选项的后果完全不同——选项 1 绕过门禁，选项 2/3 需要重新运行管线，选项 4 改变了质量标准。如果下游 agent 需要实现"PAUSE"后的恢复流程，它没有任何可依据的规格。

**Must improve**: 为 PAUSE_J/PAUSE_C 定义至少 2-3 种恢复路径及其对应的管线行为（如：接受并继续 → 标记 eval-paused → 继续下游；放弃 → exit 0 + 报告；手动修改并重试 → 从当前阶段重新开始）。

### [blindspot-4] Fact Table 更新策略与覆盖率公式的潜在冲突

**Quote**: Data Flow Table: "R2L | 运行时 Fact Table（骨架测试捕获） | gen-test-scripts（重新生成时引用） | 追加/更新 .forge/fact-table.json，source: runtime，覆盖相同 subject+kind 的 static 条目"
**Quote**: Fact Table Core Model: "Fact Table 覆盖率 = (Contract 中引用 confirmed/runtime 事实的 Outcome 数) / (Contract 总 Outcome 数) × 100%"

**Problem**: 运行时事实覆盖相同 `subject`+`kind` 的 static 事实后，如果运行时事实的 `confidence` 不是 `confirmed`（例如运行时崩溃导致 confidence 为 `assumed`），而原来 static 事实的 confidence 是 `inferred`，则覆盖率公式只计算 `source = runtime` 且 `confidence = confirmed` 的事实。这意味着替换后，如果 runtime fact 的 confidence 不满足 confirmed 条件，该 subject+kind 的覆盖率贡献变为 0——原本 static 的 inferred 不算（公式要求 confirmed），runtime 的 assumed 也不算。覆盖率可能因 Run-to-Learn 反而下降。这个更新策略与覆盖率公式存在潜在冲突。

**Must improve**: 明确 Fact Table 更新策略: runtime 事实覆盖 static 事实时，是替换还是共存？如果替换，覆盖率公式如何处理 confidence 不满足的 runtime 事实？是否应该保留 static 事实作为 fallback？

---

## Score Summary

| Dimension | Score | Max |
|-----------|-------|-----|
| 1. Background & Goals | 82 | 100 |
| 2. Flow Diagrams | 130 | 150 |
| 3. Flow Completeness | 150 | 200 |
| 4. User Stories | 160 | 200 |
| 5. Scenario Completeness | 105 | 150 |
| 6. Edge Case Coverage | 63 | 100 |
| 7. Scope Clarity | 83 | 100 |
| **Total** | **773** | **1000** |

---

## Attacks Summary

1. [Background & Goals]: Goal "降低 Mobile 接入成本"的 metric 是交付物列表而非量化指标 — "生成 Maestro YAML 骨架 + deep link 测试" — 补充对照度量（如从零到可审阅测试的步骤数/时间）
2. [Background & Goals]: "v3.0.0 是重构的最佳窗口"在 Why 中提出但未在 Goal 中体现为版本约束或时间约束 — "正处于大版本分支上，可以做大范围变更而不破坏已发布版本" — 如果版本窗口是关键动机，应在 Goals 或 Delivery Phasing 中绑定版本约束
3. [Flow Diagrams]: SCENE_FAIL 到 DETECT 的信息流断裂 — "SCENE_FAIL -->|用户选择类型 / 选择存入 session 缓存| DETECT" — 在 Data Flow Table 中补充 SCENE_FAIL 后 scene_type 的写入路径，或在流程图中体现写入 session 的操作
4. [Flow Completeness]: eval 评分结果到 revise 技能的数据传递路径在 Data Flow Table 中缺失 — Data Flow Table 有 7 行但不含 eval → revise 的数据格式 — 添加 EVAL → REVISE 行，说明评分结果的数据结构和传递方式
5. [Flow Completeness]: PAUSE_J/PAUSE_C 后的恢复路径未定义 — "暂停管线并输出当前评分 + 未通过项明细，由用户决定后续操作" — 定义至少 2-3 种恢复路径（接受继续、放弃、手动修改重试）及其管线行为
6. [User Stories]: Story 7 "不受影响"的验证标准模糊 — "已有场景类型的测试生成结果不受影响（回归验证）" — 定义可接受的差异范围（如输出结构等价、允许 LLM 随机差异但不允许结构变更）
7. [Scenario Completeness]: 场景检测规则 CLI/TUI 检测仅覆盖 Go 项目 — 表格中 CLI 行只有 `main.go + cobra.Command` — 补充非 Go 语言的 CLI/TUI 检测规则或明确声明当前覆盖范围
8. [Edge Case Coverage]: gen-test-scripts 输出无质量验证 — Flow Diagram `GEN_SCRIPTS --> R2L_CHOICE` 直接进入下游 — 添加生成后的语法/可执行性检查步骤或至少定义 gen-test-scripts 输出的质量标准
9. [Edge Case Coverage]: gen-contracts 合约 schema 验证失败未覆盖 — 合约是管线核心中间产物 — 定义 schema 验证失败的处理（重试？降级？报错？）
10. [blindspot]: eval-contract "事实依据"维度与管线衍生引擎的必须 Outcome 架构矛盾 — LLM 衍生的边界 Outcome 不在 Fact Table 中，eval 可能拒绝管线自身产出的合理 Outcome — 在 eval rubric 中区分"事实依据"和"合理推理"两类声明
11. [blindspot]: Fact Table runtime 覆盖 static 策略可能与覆盖率公式冲突 — runtime 替换后 confidence 可能不是 confirmed，覆盖率反而下降 — 明确替换策略: 是否保留 static 作为 fallback
12. [blindspot]: PAUSE_J/PAUSE_C 恢复路径定义缺失 — "由用户决定后续操作"没有定义可选操作 — 定义 2-3 种恢复路径及对应管线行为
