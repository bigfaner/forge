# Proposal Evaluation Report — Iteration 1

**Proposal**: Behavioral Test Accuracy
**Date**: 2026-06-08
**Scorer**: CTO Adversary

---

## Reasoning Audit (Phase 1)

### Argument Chain Trace

1. **Problem → Solution**: Three-layer root cause (L1: CRUD-only tests, L2: empty seed data, L3: low eval scores bypassed) maps cleanly to three solution pillars (Golden Path journeys, Fixture Specification, assertion depth rules). Mapping is tight — no orphan causes or mismatched solutions.

2. **Solution → Evidence**: The pm-work-tracker milestone map case is a real project failure, not hypothetical. It demonstrates all three root causes simultaneously. Evidence strength: concrete and specific.

3. **Evidence → Success Criteria**: SC-1 through SC-7 cover each solution pillar. SC-7 provides regression anchor tied to the original evidence case.

4. **Self-contradiction check**: Found one tension — see Logical Consistency dimension.

### SC Consistency Deep-Dive

Clustered by affected area:

- **gen-journeys**: In Scope #1, SC-1, SC-5 (journey eval). SC-1 requires Golden Path with "primary user story core domain action sequence" (semantic completeness). SC-5 requires Journey eval "Workflow Coverage" dimension with Golden Path existence as veto item. Directionally consistent.
- **gen-contracts**: In Scope #2, SC-2, SC-6. SC-2 requires fixture_spec in Preconditions. SC-6 requires Contract eval "Fixture Specification" dimension with entity coverage as veto. Consistent.
- **gen-test-scripts**: In Scope #3, SC-3, SC-4. SC-3 requires ≥80% behavioral assertions. SC-4 requires fixture to create ≥N child entities when fixture_spec declares N. Consistent — one governs assertions, the other governs fixture generation.
- **eval rubrics**: In Scope #4, SC-5, SC-6. Both define new dimensions with 60% minimum thresholds. No overlap conflict — one is journey eval, other is contract eval.
- **Regression**: SC-7 is end-to-end validation, references all three areas. No conflict with individual SCs.

One ambiguous pair: SC-3 (≥80% behavioral assertions) applied to **all** generated tests vs. the assertion classification table which excludes health check/readiness endpoints. If a feature generates many health-check tests, the denominator definition matters. The proposal states "计入 ≥80% 分母" column clarifies this — health checks are excluded from the ratio. Acceptable but borderline.

No contradictions found. All pairs satisfiable.

---

## Rubric Scoring (Phase 2)

### Dimension 1: Problem Definition — 102/110

**Problem stated clearly (38/40)**:
The problem is sharply defined: tests are structural (API doesn't crash, data format correct) rather than behavioral (functionality actually works). The three-layer root cause analysis (L1/L2/L3) is specific and well-structured. Minor deduction: the "rather than" framing presents a binary, but the actual issue is a spectrum — the tests verify *some* behavior, just not the *critical* behavior. The problem statement could acknowledge this nuance. (-2)

**Evidence provided (35/40)**:
The pm-work-tracker example is concrete and specific — a real project where all tests passed but the core feature (milestone map had no milestones) was broken. L1-L3 root cause analysis is thorough. Deduction: the evidence is a single project example. A second independent instance would strengthen the case for "systematic deficiency" claim. The urgency section asserts "any feature with parent-child entity relationships" would hit this, but only provides one data point. (-5)

**Urgency justified (29/30)**:
The urgency is well-justified: it's positioned as a systematic architectural deficiency, not a one-off. The statement "pm-work-tracker 项目已因此产生了虚假的'全部通过'结果" makes the cost of delay concrete. Minor deduction: the proposal doesn't quantify how many current/queued features have parent-child relationships and are thus at risk. (-1)

**Attacks**:
1. [Problem Definition]: Single evidence point for "systematic" claim — "这是管线架构的系统性缺陷" but only pm-work-tracker is cited. A second independent instance would make the systematic claim airtight.
2. [Problem Definition]: L3 root cause (eval scores below target but bypassed) is stated as evidence but the proposal's own solution doesn't directly address eval gate bypass behavior — that's listed as out-of-scope ("eval gate 行为（已由 eval-diagnostic-mode 覆盖）"). L3 is presented as contributing to the problem but the solution only addresses L1+L2. This is a gap in the problem-solution mapping at the root cause level.

---

### Dimension 2: Solution Clarity — 108/120

**Approach is concrete (38/40)**:
Three clear pillars with specific changes to specific skills. The Fixture Specification Schema with YAML examples is particularly concrete. Minor deduction: the "断言分类判据" table is clear for the assertion classification, but the enforcement mechanism (how does the system automatically classify assertions as behavioral vs. structural during generation?) is not explained. Is it a rubric for human eval, or a rule for the generation prompt? This ambiguity matters for implementation. (-2)

**User-facing behavior described (38/45)**:
The "Happy path" scenario describes the pipeline flow well. However, "user-facing behavior" is ambiguous here — is the "user" the developer running the pipeline, or the QA person reading test results? The proposal conflates both. From the developer's perspective: they run gen-journeys and get richer journeys. From the QA perspective: tests now fail when they should. The observable behavior for the **developer** (what they see different when running the pipeline) is only described narratively, not with concrete before/after examples. (-7)

**Technical direction clear (32/35)**:
The YAML schema for fixture_spec, the assertion classification table, and the feature complexity heuristic table all provide concrete technical direction. Deduction: the interaction between "gen-journeys 强制 Golden Path" and the complexity classification is described but the enforcement mechanism is unclear. Is there a validation step in gen-journeys that checks for Golden Path? Is it purely a rubric eval? The boundary between "generation rule" and "eval rubric check" is blurry. (-3)

**Attacks**:
3. [Solution Clarity]: Before/after comparison missing — the happy path describes the new pipeline flow but never shows what the *old* output looked like (empty CRUD test) vs. the *new* output (behavioral test). A single concrete before/after example would make the user-facing behavior crystal clear.
4. [Solution Clarity]: Assertion enforcement mechanism undefined — "断言深度规则" is mentioned as a gen-test-scripts rule, but the proposal doesn't specify whether the 80% threshold is enforced at generation time (prompt rule), at eval time (rubric), or at run time (assertion counting script). This matters significantly for implementation.

---

### Dimension 3: Industry Benchmarking — 95/120

**Industry solutions referenced (30/40)**:
Three industry references: Playwright/Cypress best practices, TestNG/JUnit fixture patterns, Contract Testing (Pact). These are legitimate references. However, Playwright and Cypress are *test execution frameworks*, not test *generation* frameworks. The proposal's challenge (AI-generated tests) is fundamentally different from human-written tests. The references acknowledge this ("手动测试的行业标准，但 Forge 需要将这些原则编码为生成规则") but the gap is significant. No reference to AI test generation tools/papers (e.g., Codium/test generation research, LLM-based test generation benchmarks) which would be more directly relevant. (-10)

**At least 3 meaningful alternatives (22/30)**:
Four alternatives listed including "do nothing". The "仅增强 eval rubric" alternative cites "eval-test-cases 提案" — is this a real proposal or a straw man? It's listed with minimal description and immediately rejected. The "轻量级规则增强" has no source and no concrete description of what it entails. It reads as a straw man: "仅改 gen-test-scripts" with cons that are assumed rather than demonstrated. (-8)

**Honest trade-off comparison (22/25)**:
The comparison table is structured. The selected approach's cons are acknowledged ("改动涉及 3 skill + 2 rubric"). However, the "改动涉及 3 skill + 2 rubric" con is presented as the only downside — the actual risk of coordinated changes across 5 files causing integration issues is understated. (-3)

**Chosen approach justified against benchmarks (21/25)**:
The justification is "唯一能系统性解决三层根因的方案" which is a strong claim. It's partially justified by the root cause analysis. But the proposal doesn't explain why a staged approach (e.g., start with fixture spec only, add Golden Path later) wouldn't also address L1 and L2. The all-or-nothing framing is presented without proving that incremental is inferior. (-4)

**Attacks**:
5. [Industry Benchmarking]: Missing AI test generation benchmarks — the proposal is about AI-generated tests but cites only human-authored test frameworks. No reference to research or tools in the AI test generation domain (e.g., how other AI coding tools handle behavioral test generation).
6. [Industry Benchmarking]: "轻量级规则增强" is a straw man — "仅改 gen-test-scripts" with "上游仍是 CRUD 描述，推断可能不准确" as the only con. This alternative is not given a fair treatment; there's no source, no concrete description, and it exists solely to be rejected.

---

### Dimension 4: Requirements Completeness — 88/110

**Scenario coverage (32/40)**:
Happy path, edge cases (simple feature, single-entity CRUD), and error scenarios (missing PRD workflow, unverifiable entity relationships) are all listed. However, the error scenario "实体关系无法从代码推断" mentions "标记为 unverifiable" but doesn't describe what happens downstream — does the test proceed with minimal data? Does it block the pipeline? The user experience for this error path is undefined. (-4)

A missing error scenario: what happens when gen-journeys generates a Golden Path but gen-contracts can't map it to a valid fixture_spec (e.g., the PRD describes a workflow involving entities not yet modeled in code)? This gap between journey and contract is unaddressed. (-4)

**Non-functional requirements (20/40)**:
This is the weakest area. The proposal doesn't address:
- **Performance**: Will the richer fixture generation and deeper assertions slow down test generation? Test execution time?
- **Compatibility**: The fixture_spec schema is new — what happens with existing projects that already have generated contracts without fixture_spec? Migration path?
- **Maintainability**: The complexity heuristic table adds cognitive load for future contributors — is this documented?
- **Backward compatibility**: Are existing contracts without fixture_spec still valid? The proposal doesn't say. (-20)

**Constraints & dependencies (36/30)**: Capped at 30.
Dependencies are clearly stated: three existing skills + eval rubrics. No external dependencies. No new pipeline stages. Clear and honest. Full marks within cap.

**Attacks**:
7. [Requirements Completeness]: No backward compatibility consideration — existing projects may already have generated contracts without fixture_spec. What happens when the pipeline re-runs? Do old contracts need migration? The proposal is silent on this.
8. [Requirements Completeness]: Performance impact unaddressed — richer fixtures mean longer test setup, deeper assertions mean longer test execution. No acknowledgment of potential regression.

---

### Dimension 5: Solution Creativity — 80/100

**Novelty over industry baseline (35/40)**:
The proposal honestly acknowledges that "golden path" and "rich test data" are standard QA practices. The novelty claim is specifically about making AI generation pipelines infer these requirements automatically, and the Contract-level declarative Fixture Specification. The fixture_spec approach is genuinely innovative — moving data requirements upstream to the contract layer rather than inferring at test generation time. (-5 for overstating novelty slightly — "Contract 级别的声明式 Fixture Specification" is essentially dependency injection for test data, which is a well-known pattern in many frameworks)

**Cross-domain inspiration (28/35)**:
The proposal draws from Contract Testing (Pact) and fixture patterns (TestNG/JUnit). However, no inspiration from domains like type theory (contracts as types), formal verification, or even database migration systems (declarative schema → fixture as "data migration"). The cross-domain inspiration is limited to the testing domain itself. (-7)

**Simplicity of insight (17/25)**:
The core insight ("declare what data you need upstream, consume it downstream") is clean and elegant. However, the overall solution has accumulated significant complexity: assertion classification tables, feature complexity heuristics, evaluation dimension thresholds, veto mechanisms, anti-checkbox-compliance mechanisms. The elegance of the core insight is somewhat buried under implementation detail that belongs in design, not proposal. (-8)

**Attacks**:
9. [Solution Creativity]: Over-engineering in proposal stage — the assertion classification table, feature complexity heuristic table, and eval dimension scoring details (150 points, 100 points, 60% thresholds, veto items, anti-checkbox mechanisms) are design-level details that clutter the proposal. The core creative insight (declarative fixture spec + golden path enforcement) is strong enough to stand alone.

---

### Dimension 6: Feasibility — 78/100

**Technical feasibility (32/40)**:
The proposal correctly identifies that all changes are within existing skill files. However, the claim "不涉及架构变更" undersells the complexity: gen-journeys must now parse PRD/Design to extract "primary user story core domain action sequence" — this requires significant prompt engineering and may be unreliable for ambiguous PRDs. The feature complexity heuristic ("实体关系 > 工作流描述") is a hard rule that may not always be correct. (-8)

**Resource & timeline feasibility (22/30)**:
"每个 skill 的修改是独立的，可以并行" — this is optimistic. The three skills form a data pipeline: gen-journeys output feeds gen-contracts, which feeds gen-test-scripts. Changes to the fixture_spec schema require coordinated updates across all three. True parallelism requires upfront agreement on the fixture_spec schema, which isn't called out as a dependency. (-8)

**Dependency readiness (24/30)**:
No external dependencies, which is good. But the proposal assumes gen-journeys can reliably extract "primary user story core domain action sequence" from PRD/Design documents. This assumes PRD/Design documents are well-structured and contain this information. The error scenario acknowledges this might fail ("PRD 中缺少工作流描述") but the mitigation ("应从用户故事中推断，或要求补充") is vague — who is asked to supplement? The AI? The user? (-6)

**Attacks**:
10. [Feasibility]: "可以并行" claim is misleading — gen-journeys → gen-contracts → gen-test-scripts is a data pipeline where fixture_spec schema changes must be agreed upon first. The three skills are not truly independent; they have a sequential data dependency.

---

### Dimension 7: Scope Definition — 72/80

**In-scope items are concrete (26/30)**:
Four in-scope items, each naming specific files/deliverables (skill rules, templates, rubrics). Minor issue: "gen-test-scripts 规则：断言深度规则 + seed data 丰富度规则" — are these two separate rules or one combined rule? The granularity could be clearer. (-4)

**Out-of-scope explicitly listed (23/25)**:
Five items listed. Good coverage. Minor issue: "共享 fixture 库 / 集中式 fixture 管理" is listed as out-of-scope but the fixture_spec per-contract approach might naturally lead to duplication across contracts. The proposal doesn't acknowledge this potential issue. (-2)

**Scope is bounded (23/25)**:
The scope is well-bounded: modify 3 skills + 2 rubrics, no new skills or pipeline stages. However, the eval rubric changes are described with enough specificity (new dimensions, point allocations, veto mechanisms) that they feel like design specifications rather than scope boundaries. This blurs the line between "what's in scope" and "how to implement it." (-2)

**Attacks**:
11. [Scope Definition]: Fixture duplication across contracts unacknowledged — with per-contract fixture_spec, related contracts (e.g., "create milestone" and "update milestone") may declare overlapping fixture requirements. Whether this leads to duplicated test setup is an open question the scope section ignores.

---

### Dimension 8: Risk Assessment — 75/90

**Risks identified (24/30)**:
Three risks identified. Missing risks:
- **Integration risk**: Coordinated changes across 3 skills and 2 rubrics create integration testing challenges. Not identified.
- **Regression risk**: New assertion rules may cause existing passing tests to be re-evaluated as "insufficient," breaking established workflows. Not identified.
- **PRD quality dependency risk**: The entire solution assumes PRD/Design documents contain enough workflow detail. If they don't (which is common in early-stage features), the Golden Path generation will be unreliable. This is partially addressed in error scenarios but not as a risk. (-6)

**Likelihood + impact rated (25/30)**:
Ratings seem reasonable. The "Golden Path 强制要求导致简单 feature 的 Journey 过度膨胀" risk is rated M/L which is honest. However, "Contract Fixture Specification 增加 Contract 复杂度" impact is rated M — but if this causes existing contracts to fail eval, the impact on project velocity could be H. (-5)

**Mitigations are actionable (26/30)**:
Mitigations reference specific mechanisms (complexity heuristic table, 80% threshold, assertion classification table). This is actionable. Minor deduction: "新 rubric 维度设 60% 最低通过阈值，且关键子项一票否决" — the veto mechanism itself creates a new risk (false negatives blocking valid work) which isn't addressed. (-4)

**Attacks**:
12. [Risk Assessment]: Missing integration risk — three skills form a data pipeline with a shared schema (fixture_spec). Coordinated changes across all three plus two rubrics create real integration risk not captured in the risk table.
13. [Risk Assessment]: Veto mechanism risk unaddressed — "关键子项一票否决" can cause false negatives. What happens when the veto incorrectly fires (e.g., gen-journeys generates a valid Golden Path but the eval rejects it due to PRD ambiguity)? No fallback mechanism described.

---

### Dimension 9: Success Criteria — 68/80

**Criteria are measurable and testable (25/30)**:
SC-1 through SC-7 are generally measurable. SC-1 ("至少生成 1 个 Golden Path Journey") is testable. SC-3 ("≥80% 的断言验证业务结果") is quantified. SC-7 is regression-testable. However:
- SC-1's "核心领域动作序列（语义完整性约束）" is subjective — who determines if the "semantic completeness" constraint is met? The eval rubric (SC-5) provides some checking, but "semantic completeness" is inherently hard to automate. (-3)
- SC-5 and SC-6 specify eval dimensions with point values (150 and 100 respectively) — this presumes the eval rubric total is being restructured, which is a design decision embedded in a success criterion. (-2)

**Coverage is complete (20/25)**:
SCs cover all four in-scope areas. However:
- No SC for the assertion classification table itself — who validates that the classification (e.g., "assert list.length > 0 is behavioral") is correct?
- No SC for the feature complexity heuristic — the complexity classification rule is proposed but there's no criterion to verify it produces correct classifications. (-5)

**SC internal consistency (23/25)**:
Per the deep-dive above, no contradictions found. All SCs are independently satisfiable. Minor ambiguity: SC-3's 80% denominator depends on the assertion classification table's "计入 ≥80% 分母" column, which excludes health checks. But no SC validates that the classification table is correctly applied. (-2)

**Attacks**:
14. [Success Criteria]: "语义完整性约束" is subjective and not automatable — SC-1 requires Golden Path to cover "核心领域动作序列（语义完整性约束）" but this is a judgment call. The anti-checkbox mechanism in SC-5 helps, but the SC itself can't be objectively verified without human judgment.
15. [Success Criteria]: No SC for validating the classification heuristic — the feature complexity classification table and the assertion classification table are core mechanisms, but no success criterion verifies they produce correct results.

---

### Dimension 10: Logical Consistency — 80/90

**Solution addresses the stated problem (32/35)**:
The three-pillar solution directly addresses L1 (CRUD-only → Golden Path with multi-step workflows), L2 (empty seed data → fixture_spec with entity requirements), and the structural assertion problem (→ assertion depth rules with 80% behavioral threshold). One gap: L3 (low eval scores bypassed) is identified as a root cause but the solution treats eval gates as out-of-scope ("eval gate 行为已由 eval-diagnostic-mode 覆盖"). This means one of the three stated root causes is deferred. The proposal should either remove L3 from the root cause analysis or acknowledge it as partially out-of-scope. (-3)

**Scope ↔ Solution ↔ Success Criteria aligned (24/30)**:
In Scope #1 (gen-journeys) ↔ SC-1 (Golden Path) ↔ SC-5 (journey eval): aligned.
In Scope #2 (gen-contracts) ↔ SC-2 (fixture_spec) ↔ SC-6 (contract eval): aligned.
In Scope #3 (gen-test-scripts) ↔ SC-3 (assertions) + SC-4 (fixture count): aligned.
In Scope #4 (eval rubrics) ↔ SC-5 + SC-6: aligned.
SC-7 is regression validation, cross-cutting.

Gap: The "Innovation Highlights" section claims "Contract 级别的声明式 Fixture Specification" eliminates "推断错误的可能性" — but the error scenario for "实体关系无法从代码推断" acknowledges this isn't always solvable. The innovation claim is slightly overstated relative to the solution's actual capability. (-6)

**Requirements ↔ Solution coherent (24/25)**:
Requirements (happy path, edge cases, error scenarios) map cleanly to the three solution pillars. No orphan requirements. One minor gap: the error scenario "PRD 中缺少工作流描述：gen-journeys 应从用户故事中推断" describes a solution behavior that has no corresponding SC. What's the acceptance criteria for "推断" quality? (-1)

**Attacks**:
16. [Logical Consistency]: L3 root cause deferred but presented as addressed — the root cause analysis lists L3 (eval scores below target bypassed) as a contributing factor, but the solution doesn't address eval gate bypass behavior (explicitly out-of-scope). This creates a logical gap: the problem analysis identifies three root causes but only two are addressed by the solution.
17. [Logical Consistency]: Innovation claim overstates capability — "这消除了'推断错误'的可能性" but the error scenarios acknowledge unverifiable entity relationships. The fixture_spec approach reduces inference errors but doesn't eliminate them.

---

## Phase 3: Blindspot Hunt

18. [blindspot]: **No definition of "behavioral test" vs. "structural test" at the concept level** — The proposal dives into assertion classification but never provides a foundational definition. The assertion table gives examples but not a general principle. For a proposal titled "Behavioral Test Accuracy," the absence of a clear, one-sentence definition of what makes a test "behavioral" is a conceptual gap. The closest is: "验证功能真正可用" in the Problem section, but this is informal and not carried through as a formal definition.

19. [blindspot]: **No consideration of test execution failure diagnosis** — The proposal focuses on making tests *fail when they should* (catching real bugs), but doesn't address what happens when these richer tests fail. Currently, an HTTP 500 error is easy to diagnose. A failed behavioral assertion ("assert milestone.map_id == map.id") gives less diagnostic context. The proposal increases test quality at the cost of test failure diagnosability, and this trade-off is not acknowledged.

20. [blindspot]: **The "80% behavioral assertion" metric may incentivize the wrong behavior** — If the system optimizes for the 80% ratio, it may generate trivial behavioral assertions (e.g., asserting that a created entity's name matches the input) rather than meaningful ones (asserting state transitions, side effects). The proposal doesn't distinguish between "shallow behavioral" and "deep behavioral" assertions.

21. [blindspot]: **No success criterion for the actual user outcome** — SC-1 through SC-7 measure pipeline output quality, but none measure whether the *end user* (developer using Forge) actually gets better bug detection. SC-7 comes closest with the pm-work-tracker regression test, but it's a single-feature regression test. There's no criterion like "when applied to a project with known bugs, the new pipeline detects ≥X% of them."

---

## Bias Detection Report

**Annotated regions**: 8 attack points (attacks #1-2 partial, #14, #16, #17 from annotated SC/Innovation sections + #18, #19, #20 which reference annotated content) / 8 annotated paragraphs = density 1.00
**Unannotated regions**: 13 attack points / 12 unannotated paragraphs = density 1.08
**Ratio (annotated/unannotated)**: 0.93

The ratio is close to 1.0, indicating no significant bias toward or against pre-revised regions. Attack density is roughly uniform across both annotated and unannotated content.

---

## Score Summary

| Dimension | Score | Max |
|-----------|-------|-----|
| Problem Definition | 102 | 110 |
| Solution Clarity | 108 | 120 |
| Industry Benchmarking | 95 | 120 |
| Requirements Completeness | 88 | 110 |
| Solution Creativity | 80 | 100 |
| Feasibility | 78 | 100 |
| Scope Definition | 72 | 80 |
| Risk Assessment | 75 | 90 |
| Success Criteria | 68 | 80 |
| Logical Consistency | 80 | 90 |
| **Total** | **846** | **1000** |
