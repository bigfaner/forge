---
created: 2026-05-23
evaluator: CTO-Adversary
iteration: 1
target: N/A
---

# Proposal Evaluation: Test Capability 2.0 — Iteration 1

## Phase 1: Reasoning Audit

### Problem -> Solution Chain

三大结构性缺陷与五大解决方案的映射关系：

| Problem | Solution | Mapping Quality |
|---------|----------|----------------|
| 双路径并行造成困惑和维护负担 | 管线统一：退休 gen-test-cases | **Strong** — 直接消除冗余 |
| 测试偏重 happy path | 深度增强：边界/异常衍生 + 风险驱动密度 | **Strong** — 精准定位瓶颈 |
| Convention 只覆盖 3 个框架 | 通用扩展：内置更多 Convention + 自动生成 | **Strong** — 直接扩展 |
| Journey-Contract 路径缺少评测 | 评测补全：eval-journey + eval-contract | **Strong** — 填补空白 |
| 场景差异化停留在模板层面 | 场景差异化：策略+层级双维度 | **Strong** — 从模板升级到策略 |

**Verdict**: 问题-方案映射清晰完整，无悬空问题或无根方案。

### Solution -> Evidence Chain

提案引用的代码结构经验证：
- `gen-test-cases` 技能目录存在（含 SKILL.md、6个模板、5个类型文件）— 确认
- `eval-test-cases` 命令存在 — 确认
- `testing-journey-contract.md` Convention 已发布且稳定 — 确认
- `test-guide` 技能目录存在 — 确认
- `gen-test-scripts/types/` 包含 5 个场景类型 — 确认
- Convention 文件确为 3 个（go, vitest, ginkgo）— 确认

**Verdict**: 证据可信，核心断言可独立验证。

### Self-Contradiction Check

- "退休 gen-test-cases" vs run-tests SKILL.md 仍引用 gen-test-cases — 提案在 Scope 中列出了"质量门禁更新"但未显式处理 run-tests 中的引用。这不是逻辑矛盾，但是遗漏。
- "排除单元测试" 在整个文档中立场一致，无矛盾。
- "Mobile 尽力而为" vs "各场景类型有独立的测试策略和层级定义文档" — Mobile 仍有独立策略文档，只是策略级别降低，不矛盾。

---

## Phase 2: Rubric Scoring

### 1. Problem Definition (110 pts)

#### Problem stated clearly (35/40)

三大问题陈述清晰："两条并行路径"、"happy path 偏重"、"3 个框架局限"。每个问题都是独立的结构性缺陷，无歧义。

**-5 原因**: 第二个问题"生成的测试偏重 happy path"稍宽泛 — 没有量化"偏重"到什么程度。当前测试覆盖率数据缺失，读者无法判断问题的严重程度。

**Quote**: *"当前合约规范和测试脚本生成主要覆盖 happy path，边界值、异常输入、错误恢复、集成交互等场景需要手动补充"* — "主要覆盖"和"需要手动补充"缺乏量化。

#### Evidence provided (35/40)

5 条证据均为可验证的代码结构事实，经代码库独立验证全部为真。

**-5 原因**: 第四条证据"Journey-Contract 路径上没有质量评测技能"是功能缺失陈述而非用户痛点数据。缺少来自实际用户的困惑反馈或使用数据。

#### Urgency justified (28/30)

v3.0.0 大版本分支窗口论证有力，延迟成本（维护翻倍）具体可量化。

**-2 原因**: 延迟成本"维护成本翻倍"是估算而非实测。没有给出"翻倍"的计算依据（例如每个新框架需适配两条路径 = 2x 工作量）。

**Dimension Total: 98/110**

---

### 2. Solution Clarity (120 pts)

#### Approach is concrete (36/40)

五大方向各含具体行动项：退休指定技能、增强合约规范、内置新 Convention、新增评测技能、按场景差异化。读者可以复述将要构建的内容。

**-4 原因**: "合约规范支持边界/异常场景自动衍生"中的"自动衍生"机制未说明。是通过模板规则？LLM prompt 增强？还是从 contract schema 推导？

**Quote**: *"合约规范支持边界/异常场景自动衍生"* — "自动衍生"的实现路径不透明。

#### User-facing behavior described (38/45)

Key Scenarios 节描述了 5 个用户场景，包含输入输出和预期行为。

**-7 原因**: 用户如何标记"高风险"未说明 — 是 PRD 中的字段？Journey 定义中的标签？命令行参数？这直接影响核心创新"风险驱动密度"的用户体验。

**Quote**: *"用户为安全相关功能标记高风险，管线自动产出更密集的测试矩阵"* — "标记高风险"的操作方式未定义。

#### Technical direction clear (32/35)

Per-Scenario Strategy Summary 表格提供了丰富的技术方向（执行方式、AI 优先侧重、必须衍生的边界 Outcome），足够的实现指导。

**-3 原因**: "三个关键设计建议"中的 Run-to-Learn 机制和置信度评级系统在 Requirements 和 Scope 中均未列入，读者不清楚这些是建议还是交付项。

**Quote**: *"Run-to-Learn 机制：生成骨架测试→运行捕获实际输出→丰富 Fact Table→重新生成精确测试"* — 在 Scope 的 In Scope / Out of Scope 中均未出现。

**Dimension Total: 106/120**

---

### 3. Industry Benchmarking (120 pts)

#### Industry solutions referenced (25/40)

引用了 Cucumber/Gherkin BDD、Postman/Newman、Playwright/TestProject，但每个只有一句话描述。

**-15 原因**: 引用浮于表面。没有版本号、没有具体功能对标、没有说明 Forge 与这些工具的功能重叠度。例如"Postman/Newman API testing: 契约测试模式。Forge 的 Contract 规范更通用，不限于 HTTP API" — "更通用"的判断缺乏证据支撑。

#### At least 3 meaningful alternatives (22/30)

三个替代方案：Do nothing、增量修补、完整 2.0。

**-8 原因**: "增量修补"方案是稻草人 — 只列了"最小改动范围"的优点和"深度和通用性问题延后"的缺点，没有认真探讨增量方案的可行性和收益。真正的增量方案（例如先统一管线再逐步增强）没有被公平评估。

#### Honest trade-off comparison (18/25)

三个方案的比较表过于简略，每项只有一句话。

**-7 原因**: 选中方案（完整 2.0）的 Cons 只有"工作量大，需要 full pipeline"，没有展开具体风险（例如一次性大变更的集成风险、测试覆盖的回归风险等）。

#### Chosen approach justified (20/25)

"结构性问题需要结构性解决方案"的论证逻辑合理，v3.0.0 窗口期加强了论证。

**-5 原因**: 依赖了 v3.0.0 窗口期的假设 — 如果 v3.0.0 发布时间紧迫，"完整 2.0"可能反而不合适。提案没有讨论时间约束与方案粒度的权衡。

**Dimension Total: 85/120**

---

### 4. Requirements Completeness (110 pts)

#### Scenario coverage (33/40)

5 个 Key Scenarios 覆盖了 happy path、冷启动、高风险、多场景、评测门禁。缺少的：

**-7 原因**: 缺少 error/recovery 场景 — 当 eval-journey 评分低于阈值且迭代修正仍不通过时，管线如何处理？自动降级？阻断并报错？这是"评测门禁"场景的关键补充。

**Quote**: *"journey 和 contract 各有独立评测技能，低于阈值自动迭代修正"* — "自动迭代修正"失败的后续处理未定义。

#### Non-functional requirements (35/40)

3 条 NFR 均为实际约束：Convention 加载兼容性、eval 框架复用、退休不影响已有功能。

**-5 原因**: 缺少性能 NFR — "风险驱动密度"可能导致高风险 Journey 生成 20 个测试用例，当 Journey 数量多时生成时间和 token 消耗是否可接受？

#### Constraints & dependencies (25/30)

4 项约束/依赖明确引用了具体文档和已有系统。

**-5 原因**: 缺少一个关键依赖 — 当前 run-tests SKILL.md 和 run-tasks 命令中引用了 gen-test-cases 和 graduate-tests，退休这些技能需要更新这些引用。提案在 Scope 中列了"质量门禁更新"但没有明确列出需要更新的下游引用文件。

**Dimension Total: 93/110**

---

### 5. Solution Creativity (100 pts)

#### Novelty over industry baseline (32/40)

风险驱动密度（高/中/低三级测试密度差异化）在 AI 测试生成领域有一定新意，不同于常见的统一覆盖率策略。

**-8 原因**: "三级分类"本身并不新颖 — 测试领域中基于风险的优先级划分（Risk-Based Testing）是成熟实践（ISO 29119 标准、ISTQB 大纲均有覆盖）。创新点在于将风险等级映射到自动生成的测试密度，而非风险分类本身。

#### Cross-domain inspiration (25/35)

"Convention 自动生成"借用了 IDE 中"从项目结构推导配置"的思路（类似于 VS Code 的自动设置推荐、Spring Boot 的 auto-configuration）。Run-to-Learn 借用了编译器"多趟处理"的概念。

**-10 原因**: 跨域灵感没有显式说明来源。提案中"Convention 自动生成"和"Run-to-Learn"的灵感来源未标注，读者无法判断这些是创新还是借用。

#### Simplicity of insight (20/25)

"退休 gen-test-cases 统一到 Journey-Contract"是简洁有力的结构性简化。Per-Scenario Strategy Summary 表格以统一的维度矩阵呈现 5 种场景，清晰且实用。

**-5 原因**: "双维度场景差异化"增加了系统复杂度（策略层面 + 层级层面），对于 CLI 和 API 这种确定性较高的场景类型，是否真的需要两维度的区分？提案没有证明单维度不够用的场景。

**Dimension Total: 77/100**

---

### 6. Feasibility (100 pts)

#### Technical feasibility (35/40)

所有技术基础已验证：Journey-Contract 模型已建立、eval 框架已成熟、Convention 系统已有 schema、gen-test-scripts 已有分区结构。增强是在已有基础上扩展，不需要新建基础设施。

**-5 原因**: "边界/异常场景自动衍生"的技术可行性未充分论证 — 如何从语义描述符自动推导出"not-found"或"already-exists"等边界 Outcome？这需要理解业务语义，LLM 的推导准确性如何保证？

#### Resource & timeline feasibility (22/30)

"预计 20+ coding tasks"是粗略估计，"需要 full pipeline"是定性描述。

**-8 原因**: 没有给出时间估算（周/月/冲刺），没有评估团队并行度（几人同时开发？），没有关键路径分析。20+ tasks 可能是 2 周也可能是 2 个月。

**Quote**: *"需要 full pipeline（PRD → Design → Tasks），预计 20+ coding tasks"* — 缺乏时间估算和资源配置。

#### Dependency readiness (28/30)

所有上游依赖都已就位。

**-2 原因**: 新增 pytest、JUnit、Rust/cargo test Convention 文件需要领域知识（框架特定的测试约定、断言模式、目录结构），提案没有评估团队是否具备这些框架的足够经验。

**Dimension Total: 85/100**

---

### 7. Scope Definition (80 pts)

#### In-scope items concrete (27/30)

13 项 In Scope 条目均为可交付物，大多数有明确的产出物定义。

**-3 原因**: "合约规范增强：支持边界/异常场景自动衍生描述"中的"增强"是一个能力而非交付物 — 具体改哪个文件？gen-contracts 技能？合约规范文档？

#### Out-of-scope explicitly listed (23/25)

9 项 Out of Scope 清晰列出，包含单元测试、性能测试、安全测试等。

**-2 原因**: "三个关键设计建议"中的 Run-to-Learn 机制和置信度评级系统在 In Scope 和 Out of Scope 中均未出现，处于范围灰色地带。

#### Scope is bounded (20/25)

"预计 20+ coding tasks"提供了一个粗略的范围边界。

**-5 原因**: 没有明确的时间框（deadline 或 sprint 规划）。"v3.0.0 窗口"是一个定性时间约束，没有转化为具体的里程碑或截止日期。

**Dimension Total: 70/80**

---

### 8. Risk Assessment (90 pts)

#### Risks identified (24/30)

4 项风险覆盖了退休影响、准确率、复杂度和阈值定义。

**-6 原因**: 缺少"集成回归风险" — 一次性变更 20+ tasks 可能导致管线整体退化，尤其是退休 gen-test-cases 后整个测试生成流程的改变。这是大规模重构的标准风险。

#### Likelihood + impact rated (25/30)

4 项风险的 L/I 评级合理，但偏保守 — 所有 Likelihood 都是 M，没有 H-L 或 L-H 的差异化分布。

**-5 原因**: "风险驱动密度的阈值难以定义"标注为 H/M 但 Mitigation 只有"初始版本用简单的三级分类"。如果阈值定义错误（高风险标记为低风险），影响的是测试覆盖率，Impact 可能被低估。

#### Mitigations actionable (25/30)

4 项 Mitigation 均为可执行的行动。

**-5 原因**: "先全面搜索确认无外部技能/agent 依赖 gen-test-cases 输出格式" — 这个搜索应该是一个前置条件而非 Mitigation。如果搜索结果发现存在外部依赖，Mitigation 是什么？提案没有给出备选方案。

**Dimension Total: 74/90**

---

### 9. Success Criteria (80 pts)

#### Measurable and testable (42/55)

9 条标准中 7 条可量化验证（删除确认、评分阈值 ≥850、数量差 ≥50%、≥3 个 Convention、检测 ≥5 种框架）。2 条偏定性。

**-13 原因**:
1. "对 journey 文档评分准确率 ≥ 850/1000" — "准确率"未定义。是人工标注对比？是 inter-rater reliability？评分准确率 850/1000 的度量方法不明。
2. "各场景类型有独立的测试策略和层级定义文档" — "有"是二元判断，没有质量标准。一份空文档也满足此条件。

#### Coverage complete (20/25)

覆盖了管线统一、深度增强、通用扩展、评测补全、场景差异化五大方向。

**-5 原因**: 缺少对"已使用 gen-test-cases 项目的迁移"的验证标准 — Out of Scope 中明确排除了迁移工具，但没有标准确认"不提供迁移工具"是否可接受（例如：确认无外部用户使用 gen-test-cases）。

**Dimension Total: 62/80**

---

### 10. Logical Consistency (90 pts)

#### Solution addresses problem (32/35)

五大解决方案与三大问题 + 两个附加问题的映射完整且无遗漏。

**-3 原因**: "场景差异化"同时出现在 Problem（问题 4 的隐含）和 Solution 中，但 Problem 中只说"差异化仅停留在代码模板层面"，Solution 则提出了"策略+层级双维度"。从"模板层面"到"策略层面"的跨越是否真的必要？文档没有证明模板层面无法满足需求。

#### Scope <-> Solution <-> Criteria aligned (25/30)

In Scope 的 13 项与 Solution 的 5 大方向基本对齐。Success Criteria 的 9 条覆盖了大部分 In Scope 项。

**-5 原因**: In Scope 第 10 项"test-guide 增强：自动扫描项目信号生成 Convention 草稿"在 Success Criteria 中有对应（"test-guide 能从项目文件信号自动检测 ≥ 5 种测试框架"），但 In Scope 说的是"生成 Convention 草稿"而 Criteria 说的是"检测框架" — 生成和检测是不同层次的能力。

#### Requirements <-> Solution coherent (22/25)

Key Scenarios 与 Solution 方向一致，Constraints 引用了具体文档。

**-3 原因**: "三个关键设计建议"（Run-to-Learn、环境就绪检测、置信度评级）在 Requirements Analysis 中作为"关键瓶颈"和"设计建议"出现，但在 Solution 和 Scope 中没有被吸纳或排除。这造成文档内部的不一致 — 如果是建议，应在 Scope 中标注为"未来考虑"；如果是交付项，应在 In Scope 中列出。

**Dimension Total: 79/90**

---

## Phase 3: Blindspot Hunt

### [blindspot] gen-test-cases 的实际使用痕迹未清理

**Quote**: *"退休 gen-test-cases 技能及相关评测能力"*

提案要退休 gen-test-cases，但 run-tests SKILL.md 的 Related Skills 表中仍有 `/gen-test-cases` 条目，run-tasks 命令中引用了 `T-test-graduate`。虽然提案在 In Scope 中列了"删除 test.graduate 任务类型和相关任务文件"和"质量门禁更新"，但没有明确列出需要更新 `run-tests/SKILL.md` 和 `run-tasks.md` 中的引用。退休一个被多个下游引用的技能，需要一个完整的"影响地图"。

### [blindspot] 边界/异常衍生的"自动"机制是黑箱

**Quote**: *"合约规范支持边界/异常场景自动衍生描述"*

这是提案的核心创新之一，但"自动衍生"的实现机制完全是黑箱。是从语义描述符通过规则推导？是 LLM 通过 prompt engineering 生成？还是预定义的边界模板匹配？不同的实现方式对准确性、可调试性和维护成本有根本性影响。提案应至少说明采用哪种技术路线。

### [blindspot] 风险等级由谁标记、何时标记

**Quote**: *"用户为安全相关功能标记高风险，管线自动产出更密集的测试矩阵"*

风险标记是风险驱动密度的前提，但提案未定义：
1. 风险标记在哪里？Journey 定义的 Risk 字段？PRD 的元数据？
2. 谁来标记？用户手动？自动推断？
3. 何时标记？Journey 生成时？Contract 生成时？

经验证，testing-journey-contract.md 中 Journey 的 Risk 字段确实存在（High/Medium/Low），但提案应该显式引用这个已有字段并说明如何利用它，而非让读者自行猜测。

### [blindspot] Convention 自动生成的质量标准缺失

**Quote**: *"test-guide 从项目文件信号自动推导 Convention 文件草稿，将新框架接入时间从'手动编写'降到'审核微调'"*

"审核微调"意味着生成质量需要足够高才能降低而非增加总工作量。但没有定义质量标准 — 如果自动生成的草稿需要大量修改，用户可能还不如从头写。提案应定义最低质量要求（例如生成的 Convention 覆盖 X% 的框架特性）。

### [blindspot] Mobile "尽力而为"策略缺乏退出条件

**Quote**: *"Mobile 场景的'尽力而为'策略：只生成 Maestro YAML 骨架和 deep link 测试"*

"尽力而为"是合理的降级，但缺乏：1) 升级到"核心"级别的条件（例如 AI Agent 适用性提升到 3/5？Maestro 生态成熟？）；2) 用户反馈收集机制 — 用户使用 Mobile 测试后的满意度如何度量？

### [blindspot] 管线瓶颈的优先级未转化为交付计划

**Quote**: *"管线关键瓶颈（优先级排序）：1. 语义描述符→regex 转换断裂..."*

5 个关键瓶颈和 3 个关键设计建议被识别并排序，但这些发现没有被转化为 Scope 中的交付项或 Success Criteria 中的验证标准。这意味着最有价值的改进点可能被淹没在"合约规范增强"这样的宽泛条目中。

---

## Score Summary

| Dimension | Score | Max |
|-----------|-------|-----|
| Problem Definition | 98 | 110 |
| Solution Clarity | 106 | 120 |
| Industry Benchmarking | 85 | 120 |
| Requirements Completeness | 93 | 110 |
| Solution Creativity | 77 | 100 |
| Feasibility | 85 | 100 |
| Scope Definition | 70 | 80 |
| Risk Assessment | 74 | 90 |
| Success Criteria | 62 | 80 |
| Logical Consistency | 79 | 90 |
| **Total** | **829** | **1000** |

## Top Priority Improvements

1. **Industry Benchmarking 需要实质性增强** — 当前替代方案分析是最大的弱点。需要：(a) 引用具体版本和功能；(b) 公平评估增量方案；(c) 展开选中方案的代价分析。
2. **"自动衍生"机制需要白箱化** — 边界/异常场景的自动衍生是核心创新，必须说明技术路线。
3. **关键设计建议需要纳入 Scope 或显式排除** — Run-to-Learn、环境就绪检测、置信度评级不能悬在"建议"状态。
4. **Success Criteria 需要更强的可测试性** — "评分准确率 ≥ 850/1000" 需要定义度量方法。
5. **风险等级标记机制需要显式定义** — 利用 testing-journey-contract.md 中已有的 Risk 字段，并说明标记流程。
