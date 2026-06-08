# Freeform Expert Review: Behavioral Test Accuracy Proposal

**Reviewer**: Behavioral Test Quality Architect
**Date**: 2026-06-08
**Document**: `docs/proposals/behavioral-test-accuracy/proposal.md`

---

## Background Assessment

This proposal addresses a problem I have witnessed destroy engineering trust in test automation more times than I care to count: a test pipeline that reports full green status while the actual product is functionally broken. The pm-work-tracker milestone map incident — where "完整管线（gen-contracts → gen-test-scripts → run-test-backend → run-test-frontend）全部通过" but the map contained zero milestones — is a textbook case of what I call "vacuous test coverage." The proposal correctly identifies the three-layer root cause: CRUD-only test semantics (L1), empty seed data (L2), and insufficient eval gates (L3).

The proposed solution architecture traces the information flow chain in the right direction: from Journey (workflow semantics) through Contract (fixture requirements) to Test Script (rich assertions). This is the correct information flow — behavioral intent must originate upstream, because no amount of downstream cleverness can reconstruct workflow semantics that were never captured. The proposal's innovation highlight of "Contract 级别的声明式 Fixture Specification" addresses a genuine gap: currently, fixture needs are inferred at test generation time with insufficient domain context.

However, the proposal also reveals some important tensions and ambiguities that, if left unresolved, could undermine the entire effort. The following sections examine the most critical risks and offer concrete improvement paths.

---

## Key Risks

风险：Golden Path Journey 的定义存在可被钻空子的模糊性。提案要求"每个 feature 必须至少包含一个跨越多步操作的 Golden Path Journey"，但"多步操作"本身是一个可以被打折的概念。一个 create-then-read 序列技术上就是两步操作，但它并不构成真正的用户工作流。提案在 Success Criteria 中进一步要求"跨越 3+ 步骤的完整工作流"，但这个"3+"阈值是如何得出的？为什么不是 4 或 5？更重要的是，没有任何规则阻止 AI 生成一个技术上 3 步但语义上仍然是孤立操作的序列——例如"创建地图 → 读取地图 → 删除地图"满足 3 步要求，却完全不验证"在地图中添加里程碑"这个核心行为。这种"checkbox compliance"式的 Golden Path 比没有 Golden Path 更危险，因为它给人虚假的质量保证感。提案在 Key Risks 中承认"Golden Path 强制要求导致简单 feature 的 Journey 过度膨胀"，但关注点放错了方向——真正的风险不是过度膨胀，而是表面合规但实质空洞。

问题：Fixture Specification 的"声明式"性质缺乏可审计的具体 schema。提案声明这是"关键创新点"，但它仅在 Requirements Analysis 中描述为"每个 Contract 的 Preconditions 必须声明前置数据状态（需要哪些实体、实体间关系、最小数据量）"。这个描述停留在概念层面，没有给出具体的声明格式或 schema 示例。一个声明为"需要里程碑实体，关系为 belongs_to map，最小数量 1"的 Fixture Specification 是可审计的——评审者可以验证这三个维度是否完整。但如果生成的 Specification 写的是"创建必要的测试数据以覆盖里程碑场景"，那么"声明式"就退化成了模糊意图声明，和当前状态没有本质区别。提案没有定义什么构成"完整"的 Fixture Specification，也没有提供评审者可以用来判断完整性的判据。

风险：80% 业务结果断言阈值的操作性定义缺失。提案在 SC-3 中设定"≥80% 的断言验证业务结果（实体存在、状态正确、关系完整），而非仅 HTTP 状态码"。但"业务结果"的边界在哪里？验证一个 API 返回的 list 长度等于预期数量——这是结构性断言（检查 array length）还是行为性断言（验证 N 个子实体存在）？验证 response body 包含特定 field value ——这是结构性还是行为性？提案没有提供任何分类判据或边界案例。在 AI 生成管线的语境下，这意味着 AI 可以通过微妙的措辞调整将结构性断言包装为"行为性"断言，而 eval rubric 可能无法有效区分。80% 这个阈值看似给了 20% 的宽容度，但在没有清晰分类标准的情况下，这个比例根本无法被可靠测量。

问题：提案对"简单 feature"和"复杂 feature"的区分策略不够明确。提案在 Key Scenarios 中提到"简单 feature（无父子实体关系）：Golden Path Journey 仍然适用，但 fixture specification 可以声明为最小数据集"以及"单实体 CRUD feature：Golden Path 可能就是完整的 CRUD 循环"。这些描述暗示了一个重要的设计决策——不是所有 feature 都需要同等深度的 Golden Path——但提案没有定义区分的判据。谁来判断一个 feature 是否"简单"？是 gen-journeys skill 自动判断，还是需要人工标注？如果是自动判断，判断的依据是什么——PRD 中的实体关系描述？Design 文档的复杂度？如果判断失误，简单 feature 可能被过度工程化，复杂 feature 可能被错误简化。提案在 Key Risks 的 mitigation 中说"规则应区分'简单 feature'和'复杂 feature'的期望"，但"应"不是规范性的，它将关键设计决策推迟到了实现阶段。

风险：eval rubric 新增维度的评分标准可能导致新的虚假通过路径。提案在 SC-5 和 SC-6 中分别要求 Journey eval 新增 "Workflow Coverage"（150 分）和 Contract eval 新增 "Fixture Specification"（100 分）。但提案没有讨论这些维度的最低通过阈值。如果阈值设得太低——例如只要有一个 multi-step journey 就通过 Workflow Coverage——那么它对质量的提升微乎其微。如果阈值设得合理，但评分标准本身可以被"钻空子"（如前所述的 checkbox-compliant Golden Path），那么新增维度反而可能制造更多的虚假通过信号。提案在 Assumptions Challenged 中已经承认"eval gate 是必要条件但不是充分条件"，但新方案似乎仍然依赖 eval gate 作为主要质量保障机制，而没有解决 eval rubric 本身可能被形式化地满足这一根本问题。

问题：提案缺少对现有失败案例的回归验证计划。pm-work-tracker 里程碑地图是唯一的 motivating example，但提案的 Success Criteria 没有包含"对 pm-work-tracker 重新运行管线并验证里程碑地图测试能发现空地图问题"这样的回归测试条款。一个针对真实失败案例的端到端验证是检验整个提案有效性的最直接方式。没有这个验证，我们无法知道三层根因是否真的被同时解决了——可能 L1 被解决了但 L2 仍然存在，或者 L1 和 L2 被解决但 eval 仍然通过了一个实质上有缺陷的 Journey。

---

## Improvement Suggestions

建议：为 Golden Path Journey 增加语义完整性约束，而不仅仅是步骤数量要求。当前提案在 SC-1 中要求"至少生成 1 个 Golden Path Journey（跨越 3+ 步骤的完整工作流）"，但应该进一步要求 Golden Path 必须覆盖 feature 的核心领域动作——即 PRD 或 Design 文档中描述的 primary user story 的完整路径。例如，对于里程碑地图 feature，Golden Path 不仅要"跨越 3+ 步骤"，还必须包含"在地图中创建里程碑"这一核心行为。这可以通过要求 gen-journeys skill 从 PRD 的 primary user story 中提取"关键动作列表"，并验证 Golden Path 覆盖了列表中的所有动作来实现。

建议：为 Fixture Specification 定义一个明确的 schema，包含必需字段和可选字段。提案声明这是"Contract 级别的声明式 Fixture Specification"，应该为这个声明提供一个具体的格式规范。例如，每个 Fixture Specification 条目必须包含：entity_type（必需）、relationship_type（可选，对简单 feature）、min_count（必需，默认为 1）、field_constraints（可选，如"status 必须为 active"）。有了这样的 schema，eval rubric 就可以基于字段完整性来评分，评审者也可以明确判断一个 Specification 是否"完整"——而不是依赖主观判断。提案可以参考 TestNG/JUnit fixture patterns（已在 Industry Benchmarking 中提到）的声明式风格，将其形式化为 Forge 的专用 schema。

建议：为"业务结果断言"提供分类判据和边界案例示例。提案在 SC-3 中提到"验证业务结果（实体存在、状态正确、关系完整）"，应该在 Success Criteria 或单独的附录中提供分类判据。例如："验证 API 返回的子实体列表长度 >= N"属于行为性断言（因为它验证了 fixture 数据是否被正确持久化和检索）；"验证 response status code 为 200"属于结构性断言；"验证 response body 中某个字段不为 null"的归类取决于该字段是否承载业务含义——如果是 entity_id，则是行为性的；如果是 created_at timestamp，则可能不是。提供 5-10 个明确的分类示例，以及 3-5 个边界案例和它们的判定结果，可以让 80% 阈值变得可操作、可审计。

建议：增加一个基于 pm-work-tracker 的端到端回归验证成功标准。提案应该增加一条 SC："对 pm-work-tracker 里程碑地图 feature 重新运行完整管线（gen-journeys → gen-contracts → gen-test-scripts → run-tests），验证：(a) 生成的 Golden Path Journey 包含'在地图中创建里程碑'步骤；(b) Contract 的 Fixture Specification 声明了 map 和 milestone 的实体关系；(c) 生成的测试在空地图 fixture 上应失败（而不是通过）。" 这个回归验证不仅是确认修复效果的最直接手段，也能作为整个提案的"smoke test"——如果连这个已知失败案例都无法正确处理，那么提案的架构就有根本性问题。

建议：在 Proposal 中明确简单/复杂 feature 的区分机制，而不是推迟到实现阶段。提案可以定义一个启发式规则，例如：如果 PRD 或 Design 文档中描述了多于一种实体类型，并且这些实体之间存在关系描述（如"belongs_to"、"has_many"、"contains"），则该 feature 被归类为"复杂 feature"，需要完整的 Golden Path 和详细的 Fixture Specification。否则，feature 被归类为"简单 feature"，Golden Path 可以是完整的 CRUD 循环，Fixture Specification 可以使用最小数据集。这个启发式规则可以内嵌到 gen-journeys skill 的规则中，作为自动化的 feature 分类机制，而不是依赖人工判断或 AI 的自由裁量。
