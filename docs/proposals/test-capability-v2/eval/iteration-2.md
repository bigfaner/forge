---
created: 2026-05-23
evaluator: CTO-Adversary
iteration: 2
target: N/A
---

# Proposal Evaluation: Test Capability 2.0 — Iteration 2

## Iteration 1 改善追踪

| Iteration 1 Issue | Resolution Status | Evidence |
|-------------------|-------------------|----------|
| Industry Benchmarking 缺乏深度 | **Resolved** | Postman/Newman 增加版本号、功能对标、与 Forge 的差异分析；增量方案扩展为完整论证 |
| "自动衍生"机制是黑箱 | **Resolved** | In Scope 第 6 项详细说明了 LLM prompt 增强策略的技术路线 |
| 关键设计建议悬在"建议"状态 | **Partial** | Run-to-Learn 纳入 In Scope；环境就绪检测和置信度评级仍未归入 |
| Success Criteria 可测试性不足 | **Resolved** | 增加了 Pearson 相关系数、gold standard、inter-rater reliability 定义 |
| 风险等级标记机制未定义 | **Resolved** | In Scope 第 7 项说明了 Risk 字段复用、标记流程、读取时机 |
| gen-test-cases 下游引用未清理 | **Resolved** | Success Criteria 第 2 条增加 run-tasks/run-tests 引用清理验证 |
| Convention 自动生成质量标准缺失 | **Partial** | 增加"结构完整"描述但仍是定性 |
| Mobile 退出条件缺失 | **Unresolved** | 无变更 |
| 管线瓶颈未转化为交付项 | **Partial** | Run-to-Learn 纳入 Scope；其余 4 个瓶颈仍悬空 |
| 缺少集成回归风险 | **Resolved** | 新增第 5 项风险 + 三阶段交付计划 |
| 缺少时间估算 | **Unresolved** | 仍无时间估算，仅增加三阶段划分 |

---

## Phase 1: Reasoning Audit

### Problem -> Solution Chain

| Problem | Solution | Mapping Quality |
|---------|----------|----------------|
| 双路径并行 | 管线统一：退休 gen-test-cases | **Strong** |
| 测试偏重 happy path | 深度增强：LLM prompt 增强策略 + 风险驱动密度 | **Strong** |
| Convention 只覆盖 3 个框架 | 通用扩展：内置更多 Convention + 自动生成 | **Strong** |
| Journey-Contract 缺少评测 | 评测补全：eval-journey + eval-contract | **Strong** |
| 场景差异化停留在模板层面 | 双维度场景差异化 | **Strong** |

**Verdict**: 问题-方案映射完整，iteration 1 中识别的黑箱（自动衍生机制、风险标记机制）均已白箱化。

### Solution -> Evidence Chain

- gen-test-cases 技能目录存在 — 确认（iteration 1 已验证）
- testing-journey-contract.md 中 Risk 字段存在 — iteration 1 已验证，新版提案显式引用
- eval 框架 scorer-gate-revise 模式已成熟 — 确认
- Convention 文件 schema 存在 — 确认

**Verdict**: 新增的解决方案细节（LLM prompt 增强策略、Risk 字段复用）均基于已验证的代码结构。

### Self-Contradiction Check

1. **"三个关键设计建议"定位不一致**（残留）：Run-to-Learn 已纳入 In Scope，但第 2 项（场景特定执行环境就绪检测）和第 3 项（置信度评级系统）仍作为"建议"出现，既不在 In Scope 也不在 Out of Scope。这直接违反了 iteration 1 的改进要求。

2. **管线瓶颈 #2 与设计建议 #2 呼应但不闭合**：瓶颈"执行环境准备缺乏自动化"对应建议"场景特定执行环境就绪检测"，但建议不纳入 Scope。读者无法判断这个瓶颈是否会在 2.0 中被解决。

3. **三阶段交付计划与 "20+ tasks" 估算的对应关系不明**：风险表新增的三阶段计划提供了交付节奏，但没有给出每个阶段的任务数量估算或时间跨度。

---

## Phase 2: Rubric Scoring

### 1. Problem Definition (110 pts)

#### Problem stated clearly (38/40)

三大问题陈述清晰，无歧义。iteration 1 指出的"偏重 happy path"缺乏量化问题在本次仍未解决，但不影响问题的清晰度。

**-2 原因**: "主要覆盖 happy path"仍然缺乏量化 — 当前测试的 happy path vs non-happy-path 比例是多少？读者无法判断问题的严重程度。

#### Evidence provided (37/40)

5 条证据均为可验证的代码结构事实。iteration 1 指出的"缺少用户痛点数据"问题部分改善 — 新版在 Assumptions Challenged 中增加了更多定性证据（"从未执行、无技能实现、有已知缺陷"）。

**-3 原因**: 第四条证据"Journey-Contract 路径上没有质量评测技能"仍是功能缺失陈述，缺乏来自实际使用场景的痛点数据（例如：用户因缺少评测门禁而产出了低质量测试的实际案例）。

#### Urgency justified (28/30)

v3.0.0 窗口期论证有力。iteration 1 指出的"维护成本翻倍"缺乏计算依据的问题未变。

**-2 原因**: "每新增一个框架或场景类型，都要在两条并行路径上分别适配，维护成本翻倍"仍是定性陈述。虽然逻辑成立（2 条路径 × N 个框架 = 2N vs 1N），但"翻倍"没有附带实际成本数据（例如每个新框架适配的工作量人天估算）。

**Dimension Total: 103/110**

---

### 2. Solution Clarity (120 pts)

#### Approach is concrete (39/40)

五大方向各含具体行动项。iteration 1 中最大的弱点 — "自动衍生"机制不透明 — 现已通过 LLM prompt 增强策略的详细说明解决（根据 Input 维度字段类型注入类型特定边界提示词，结合 Convention 中的 required_outcomes 规则）。

**-1 原因**: 风险驱动密度中"gen-journeys 根据 PRD 中的安全/合规关键词自动建议 risk_level"的匹配规则未说明。哪些关键词？如何定义？这影响自动化程度的核心假设。

#### User-facing behavior described (42/45)

5 个 Key Scenarios 描述了输入输出和预期行为。iteration 1 指出的"如何标记高风险"问题已解决 — 明确说明复用 journey 文档中的 risk_level 字段。

**-3 原因**: "新项目冷启动"场景缺少失败路径 — test-guide 自动检测框架失败时（检测到未知框架、信号冲突），用户看到什么？回退到什么行为？Key Scenarios 只描述了 happy path。

#### Technical direction clear (33/35)

Per-Scenario Strategy Summary 表格提供了丰富的技术方向。LLM prompt 增强策略为自动衍生提供了清晰的实现指导。

**-2 原因**: "三个关键设计建议"中"环境就绪检测"和"置信度评级系统"的技术方向出现了但未纳入交付范围。这为读者提供了不可执行的技术方向。

**Quote**: *"场景特定执行环境就绪检测：CLI 检查二进制、WebUI 探测 dev server、API 检查服务+DB 连通性"* — 出现在"关键设计建议"中，但 Scope 的 In/Out 均未包含。

**Dimension Total: 114/120**

---

### 3. Industry Benchmarking (120 pts)

#### Industry solutions referenced (35/40)

显著改善。Postman/Newman 现在有具体版本号（Newman v6+, Collections v2.1 schema）、功能描述（pre-request scripts、test assertions）、以及与 Forge 的差异分析（HTTP-only vs 六维度通用模型）。Cucumber 和 Playwright 也增加了版本号。

**-5 原因**: 仍然缺少这些工具的测试覆盖率方法论对标。例如 Cucumber 的 BDD 场景覆盖策略与 Forge 的"风险驱动密度"有何异同？这是竞品分析的核心维度。

#### At least 3 meaningful alternatives (27/30)

三个替代方案均经过分析。iteration 1 中的稻草人问题（增量方案不公平评估）已改善 — "增量分期"方案现在有了完整的论证："管线统一会触及 gen-test-scripts 核心路径，后续深度增强仍需修改同一区域，存在二次变更成本"。

**-3 原因**: "Do nothing"方案仍是一句话带过（"双路径混乱持续"），没有分析其具体后果（例如对新框架接入的影响时间线、用户投诉增长预测等）。

#### Honest trade-off comparison (22/25)

选中方案的 Cons 从原来的"工作量大"扩展为隐含在增量方案对比中的"一次性大变更风险"。

**-3 原因**: 选中方案的代价分析仍然不够完整。三阶段交付计划（在 Risk 中提出）应该在这里作为代价的一部分被讨论 — 它意味着第一和第二阶段之间有回归验证的停顿期，这影响交付节奏。

#### Chosen approach justified (22/25)

"结构性问题需要结构性解决方案"的论证逻辑合理。"管线统一和深度增强紧密耦合，拆期会导致二次重构"是新增的有力论证。

**-3 原因**: 仍依赖 v3.0.0 窗口期假设。如果 v3.0.0 发布时间线紧迫（例如 2 周内必须合并），"完整 2.0 + 三阶段交付"是否仍然可行？提案没有给出时间约束下的降级方案。

**Dimension Total: 106/120**

---

### 4. Requirements Completeness (110 pts)

#### Scenario coverage (36/40)

5 个 Key Scenarios 覆盖了主要用户路径。iteration 1 指出的"error/recovery 场景缺失"部分改善 — 新版在风险表中增加了三阶段回归验证，但 Key Scenarios 中仍缺少失败路径。

**-4 原因**: "评测门禁"场景仍缺少"自动迭代修正失败的后续处理"。Key Scenarios 说"低于阈值自动迭代修正"，但修正 N 次仍不通过时怎么办？阻断？降级？标记为需人工审查？这是评测门禁的关键用户路径。

#### Non-functional requirements (37/40)

3 条 NFR 均为实际约束。iteration 1 指出的"缺少性能 NFR"问题仍未解决。

**-3 原因**: "风险驱动密度"可能导致高风险 Journey 生成 20 个测试用例（密度表中估算为 10-20）。当项目有 10+ 个 Journey 时，总测试数可达 200+。LLM 生成和 token 消耗是否可接受？提案在密度表中给出了"总测试数估算"但没有讨论性能边界。

#### Constraints & dependencies (27/30)

4 项约束/依赖明确。iteration 1 指出的"下游引用文件"问题已部分改善 — Success Criteria 增加了"run-tasks/run-tests 中的引用已清理"。

**-3 原因**: 新增 3 个 Convention 文件（pytest、JUnit、Rust/cargo test）需要对应框架的领域知识。JUnit 5 的测试约定与 JUnit 4 有显著差异（@ParameterizedTest、@Nested 等），提案没有说明目标版本或框架特性覆盖范围。这直接影响 Convention 文件的实用价值。

**Dimension Total: 100/110**

---

### 5. Solution Creativity (100 pts)

#### Novelty over industry baseline (34/40)

风险驱动密度（风险等级 -> 测试密度差异化）在 AI 测试生成领域有一定新意。LLM prompt 增强策略（根据字段类型注入边界提示词）是一个简洁的实现方案。

**-6 原因**: iteration 1 指出的"风险分级本身不新颖"问题仍然成立。新版本的技术路线说明虽然澄清了实现方式，但创新的核心仍是"风险 -> 密度"的映射，这仍然是 Risk-Based Testing 的直接应用。真正的新意在于将这一映射自动化到 AI 生成管线中，但提案没有强调这一点。

#### Cross-domain inspiration (28/35)

LLM prompt 增强策略借用了"类型系统 -> 约束推导"的编译器思路（string -> empty/overflow/unicode 是类型推导的边界枚举）。Run-to-Learn 借用了编译器多趟处理概念。

**-7 原因**: 跨域灵感仍未显式标注来源。读者需要自行推断这些思路的灵感来源。提案应说明"类似于编译器的类型推导"或"类似于 IDE 的 auto-configuration"。

#### Simplicity of insight (22/25)

"退休 gen-test-cases 统一到 Journey-Contract"是简洁有力的结构性简化。LLM prompt 增强策略避免了构建独立的规则引擎，选择在生成时注入提示词，是最简实现路径。

**-3 原因**: "双维度场景差异化"增加了系统复杂度。iteration 1 提出的"单维度是否够用"质疑未回应。对于 CLI 和 API 这种确定性较高的场景，策略层面和层级层面的差异化是否有必要独立为两个维度？

**Dimension Total: 84/100**

---

### 6. Feasibility (100 pts)

#### Technical feasibility (37/40)

所有技术基础已验证。LLM prompt 增强策略的可行性高于 iteration 1 的黑箱方案 — 它不需要独立的推导引擎，只需在 gen-contracts 的 prompt 中注入边界提示词。

**-3 原因**: "gen-journeys 根据 PRD 中的安全/合规关键词自动建议 risk_level"的实现可行性未充分论证。NLP 关键词匹配的准确性如何保证？误判（将低风险标记为高风险）的后果是测试密度浪费，漏判（将高风险标记为低风险）的后果是测试覆盖不足。后者的影响更严重。

#### Resource & timeline feasibility (25/30)

三阶段交付计划提供了比 iteration 1 更清晰的节奏。但仍然缺少具体的时间估算。

**-5 原因**: 三阶段划分是改进，但没有给出每阶段的任务数量、持续时间、或人力需求。20+ tasks 分配到 3 个阶段是 7/7/6 还是 3/10/7？关键路径是什么？这直接影响项目管理的可操作性。

**Quote**: *"需要 full pipeline（PRD → Design → Tasks），预计 20+ coding tasks"* — 仍是唯一的资源估算，无时间维度。

#### Dependency readiness (28/30)

所有上游依赖已就位。新增 Convention 文件的框架知识需求未评估。

**-2 原因**: JUnit Convention 需要区分 JUnit 4 vs JUnit 5 的约定差异，pytest 需要了解 fixture/parametrize 模式，Rust 需要 cargo test 的特性。团队是否具备这些框架的足够经验未说明。

**Dimension Total: 90/100**

---

### 7. Scope Definition (80 pts)

#### In-scope items concrete (28/30)

13 项 In Scope 条目中大多数有明确的产出物定义。iteration 1 指出的"合约规范增强是能力而非交付物"问题已改善 — 现在附带了完整的技术路线说明。

**-2 原因**: "test-guide 增强：自动扫描项目信号检测测试框架并生成 Convention 草稿"的"自动扫描"机制未说明。扫描哪些文件？package.json、go.mod、Cargo.toml？这影响交付物的边界。

#### Out-of-scope explicitly listed (22/25)

9 项 Out of Scope 清晰。iteration 1 指出的灰色地带问题部分改善 — Run-to-Learn 已纳入 In Scope。

**-3 原因**: "三个关键设计建议"中的 2 项（环境就绪检测、置信度评级）仍处于 In Scope 和 Out of Scope 的灰色地带。iteration 1 明确要求"纳入 Scope 或显式排除"，但新版只纳入了 1/3。

#### Scope is bounded (22/25)

三阶段交付计划为范围提供了更好的边界。

**-3 原因**: 仍无时间框（deadline 或 sprint 规划）。三阶段计划是节奏而非时间约束 — 没有"第一阶段需在 X 周内完成"的约束，阶段之间的边界是功能性的而非时间性的。

**Dimension Total: 72/80**

---

### 8. Risk Assessment (90 pts)

#### Risks identified (28/30)

5 项风险覆盖了退休影响、准确率、复杂度、阈值定义、集成回归。iteration 1 指出的"缺少集成回归风险"已补充（第 5 项风险）。

**-2 原因**: 缺少"LLM 生成边界用例的准确性风险" — LLM prompt 增强策略依赖 LLM 正确理解字段类型并生成合理的边界值。如果 LLM 对特定类型（如 enum 的 invalid_value）的推导不准确，产出的测试用例可能是无效的。这是新方案引入的核心技术风险，但在 Risk 表中缺失。

#### Likelihood + impact rated (27/30)

5 项风险的 L/I 评级合理。第 5 项风险的 H/H 评级是对 iteration 1 "偏保守"批评的回应。

**-3 原因**: "Convention 自动生成准确率不足"标注为 M/M，但如果自动生成的草稿质量很差（用户需要大量修改），Impact 应为 H — 因为这直接影响了"将新框架接入时间从手动编写降到审核微调"的核心价值主张。当前评级低估了此风险的影响。

#### Mitigations actionable (27/30)

5 项 Mitigation 均为可执行的行动。第 5 项的三阶段交付 + 回归验证是对 iteration 1 的显著改善。

**-3 原因**: 第 5 项 Mitigation 说"每阶段完成后对已交付功能做回归验证"，但回归验证的具体方式未定义。是手动运行端到端测试？自动化 CI 门禁？还是 eval 技能重新评分？没有定义验证标准，"回归验证"就是一句空话。

**Dimension Total: 82/90**

---

### 9. Success Criteria (80 pts)

#### Measurable and testable (47/55)

显著改善。评分准确率现在有明确的度量方法（Pearson 相关系数 >= 0.85、gold standard 标注流程、inter-rater reliability <= 100 分差异）。iteration 1 中"准确率 850/1000 度量方法不明"的问题已解决。

**-8 原因**:

1. "test-guide 能从项目文件信号自动检测 >= 5 种测试框架，并能为检测到的框架生成结构完整的 Convention 草稿" — "结构完整"仍是定性描述。什么算"结构完整"？覆盖 Convention schema 的所有必填字段？能生成可执行的测试？缺少质量量化标准。

2. "各场景类型有独立的测试策略和层级定义文档" — "有"是二元判断。iteration 1 指出的"一份空文档也满足此条件"问题未解决。应增加质量标准（例如：定义文档覆盖 Per-Scenario Strategy Summary 表格中的所有维度）。

#### Coverage complete (22/25)

9 条标准覆盖了 5 大方向。iteration 1 指出的"缺少 gen-test-cases 迁移验证"问题在 Success Criteria 第 1 条中间接覆盖（"所有相关文件已完全删除"）。

**-3 原因**: 缺少对"三阶段交付计划"的验证标准。Risk 表中提出了三阶段交付，但没有 Success Criteria 确认每阶段的交付完整性。例如：第一阶段完成后，应验证"退休 gen-test-cases 后，现有 Journey-Contract 管线仍可端到端运行" — 这在第 1 条和第 9 条中间接覆盖，但没有显式的阶段门禁标准。

**Dimension Total: 69/80**

---

### 10. Logical Consistency (90 pts)

#### Solution addresses problem (33/35)

五大解决方案与三大问题 + 两个附加问题的映射完整。

**-2 原因**: iteration 1 指出的"从模板层面到策略层面的跨越是否必要"的质疑未回应。文档仍然没有证明单维度（仅代码模板差异化）无法满足需求。

#### Scope <-> Solution <-> Criteria aligned (27/30)

In Scope 的 13 项与 Solution 的 5 大方向对齐。iteration 1 指出的"test-guide 检测 vs 生成层次不一致"问题已改善 — Success Criteria 现在包含"检测"和"生成"两个层次。

**-3 原因**: In Scope 第 11 项"Run-to-Learn 机制"在 Success Criteria 中没有对应条目。一个 In Scope 交付项没有验证标准，这意味着"完成"的定义不明确。Run-to-Learn 的成功标准是什么？Fact Table 丰富度提升 X%？测试准确率提升？

#### Requirements <-> Solution coherent (23/25)

Key Scenarios 与 Solution 方向一致。iteration 1 指出的"三个关键设计建议定位不一致"问题部分改善 — 1/3 已纳入 Scope。

**-2 原因**: 管线瓶颈 #2（执行环境准备缺乏自动化）和 #4/#5（失败诊断和测试数据管理策略）在 Requirements Analysis 中被识别为关键瓶颈，但在 Solution 和 Scope 中均无对应条目。识别了问题但不解决，造成文档内部的不一致 — 如果是未来工作，应标注为 Out of Scope。

**Dimension Total: 83/90**

---

## Phase 3: Blindspot Hunt

### [blindspot] 环境就绪检测和置信度评级的 Scope 归属仍是灰色地带

**Quote**: *"场景特定执行环境就绪检测：CLI 检查二进制、WebUI 探测 dev server、API 检查服务+DB 连通性"* 和 *"置信度评级系统：HIGH/MEDIUM/LOW 替代二元 pass/fail，区分可信赖测试和需审查测试"*

这两项出现在"三个关键设计建议"中，但在 In Scope 和 Out of Scope 中均未出现。iteration 1 的评估明确要求"纳入 Scope 或显式排除"，新版只解决了 3 项中的 1 项（Run-to-Learn 纳入 In Scope）。这不只是遗漏 — 它意味着读者无法判断这 2 项能力是否会在 2.0 中交付，直接影响项目计划的完整性。

### [blindspot] 三阶段交付计划缺少阶段门禁标准

**Quote**: *"分三阶段交付并设置回归验证点：(1) 管线统一...；(2) 深度增强...；(3) 通用扩展...。每阶段完成后对已交付功能做回归验证，防止后续阶段破坏前序成果"*

阶段划分清晰，但缺少阶段门禁标准：每阶段"完成"的判断依据是什么？如果第一阶段"退休 gen-test-cases"完成后回归验证发现 Journey-Contract 管线不可运行，是回滚还是修复？修复的时限是多少？这直接影响项目风险管理。

### [blindspot] gen-journeys 自动建议 risk_level 的准确性验证缺失

**Quote**: *"通过 gen-journeys 根据 PRD 中的安全/合规关键词自动建议 risk_level"*

这是风险标记机制的关键补充（降低用户手动标记的负担），但提案没有为这一自动化能力设定任何验证标准。Success Criteria 中"高风险旅程自动生成的测试用例数量比低风险旅程多 >= 50%"验证的是密度的差异化，而非 risk_level 自动建议的准确性。如果自动建议总是建议 Medium（最安全的默认值），密度差异化可能形同虚设。

### [blindspot] Convention 草稿"结构完整"缺少可操作定义

**Quote**: *"test-guide 能从项目文件信号自动检测 >= 5 种测试框架，并能为检测到的框架生成结构完整的 Convention 草稿"*

iteration 1 要求定义 Convention 自动生成的质量标准，新版以"结构完整"回应。但"结构完整"缺少可操作定义 — 是覆盖 Convention schema 的所有必填字段？是能生成至少 1 个可执行的测试？还是与手写 Convention 的 diff < 30%？不同的定义导致完全不同的交付标准。

### [blindspot] Run-to-Learn 机制的迭代终止条件未定义

**Quote**: *"Run-to-Learn 机制：生成骨架测试→运行捕获实际输出→丰富 Fact Table→重新生成精确测试，作为管线内置的迭代增强环节"*

Run-to-Learn 是一个迭代增强机制，但缺少终止条件：迭代多少次？Fact Table 丰富度达到什么阈值算"足够"？LLM 重新生成的测试与骨架测试的 diff 小于某个阈值时停止？没有终止条件，这个机制可能无限迭代（浪费 token）或过早终止（质量不足）。

---

## Score Summary

| Dimension | Score | Max |
|-----------|-------|-----|
| Problem Definition | 103 | 110 |
| Solution Clarity | 114 | 120 |
| Industry Benchmarking | 106 | 120 |
| Requirements Completeness | 100 | 110 |
| Solution Creativity | 84 | 100 |
| Feasibility | 90 | 100 |
| Scope Definition | 72 | 80 |
| Risk Assessment | 82 | 90 |
| Success Criteria | 69 | 80 |
| Logical Consistency | 83 | 90 |
| **Total** | **903** | **1000** |

## Top Priority Improvements

1. **环境就绪检测和置信度评级必须归入 Scope 或显式排除** — 这是 iteration 1 就提出的改进要求，iteration 2 只解决了 1/3。悬在"建议"状态的交付项是项目管理的不确定因素。

2. **Run-to-Learn 缺少 Success Criteria 和终止条件** — 已纳入 In Scope 但没有验证标准。一个没有"完成"定义的交付项等于没有交付项。需要定义迭代终止条件和质量度量。

3. **三阶段交付计划需要阶段门禁标准** — 每阶段的"完成"判断依据、回归验证方式、失败时的处理策略需要明确。当前的"回归验证"是一句空话。

4. **Convention 草稿"结构完整"需要可操作定义** — 用可量化的标准替代定性描述（例如：覆盖 Convention schema 必填字段的 100%、生成的测试可执行的比率 >= X%）。

5. **gen-journeys 自动建议 risk_level 需要准确性验证** — 自动建议是风险标记的关键补充，但没有验证标准。建议增加 Success Criterion：自动建议与人工标注的一致率 >= X%。

## Overall Assessment

Iteration 2 相比 iteration 1 有实质性改善（829 -> 903，+74 分），核心改善来自：(1) Industry Benchmarking 的深度增强（+21 分）；(2) 自动衍生机制的白箱化对 Solution Clarity 的提升（+8 分）；(3) Success Criteria 度量方法的细化（+7 分）；(4) 集成回归风险的补充对 Risk Assessment 的改善（+8 分）。

剩余的核心问题是 **Scope 的灰色地带**（2 项设计建议未归入）和 **交付标准的完整性**（Run-to-Learn 无 Success Criteria、三阶段无门禁标准）。这两个问题是相互关联的 — 未归入 Scope 的能力自然无法有验证标准。解决路径是：要么将剩余 2 项建议显式排除出 Out of Scope 并注明"未来迭代"，要么纳入 In Scope 并配以对应的 Success Criteria。
