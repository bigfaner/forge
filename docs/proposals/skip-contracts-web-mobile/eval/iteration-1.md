# CTO Scorer Report — Proposal: Skip gen-contracts for Interaction-Only Features

**Iteration**: 1
**Date**: 2026-06-09
**Scorer**: CTO adversary

---

## Pre-Score Anchors (Phase 1: Reasoning Audit)

### Problem → Solution trace
The stated problem is: gen-contracts produces no output for pure Web/Mobile journeys, gen-test-scripts silently skips them, pipeline reports no error. The proposed solution addresses both symptoms: (1) PipelineRegistry skip eliminates useless gen-contracts execution, (2) gen-test-scripts direct path ensures web/mobile journeys get test scripts. The mapping is direct and honest.

### Solution → Evidence trace
The "protocol vs. interaction execution model" distinction is used to justify the dual-path solution. Evidence supporting this distinction: the Innovation Highlights section provides a conceptual argument (structured data I/O vs. user actions → visual states) but no empirical evidence that journey.md alone produces equivalent-quality test scripts. The claim "journey.md 已包含生成所需全部信息" remains an assertion, not a demonstrated fact.

### Evidence → Success Criteria trace
SC-1 through SC-7 test the observable outcomes (tasks present/absent, scripts generated, coverage gap detected). They do NOT test the quality of generated test scripts — only their existence and type correctness. This is a proxy: the real goal is "Web features have useful automated tests," but SC only measures "Web features have test script files."

### Self-contradiction check
No direct contradictions found. However, the pre-revised annotation on line 24 introduces an architectural principle ("Skill 层是权威决策者") that the baseline version did not have. This is an improvement but creates a new tension: if the Skill layer is the authoritative decision-maker, then the Pipeline layer skip is purely optimization. Yet SC-1 and SC-2 test the Pipeline layer's skip behavior as if it were a guarantee, not an optimization.

### SC Consistency Deep-Dive

**Cluster 1: PipelineRegistry** (SC-1, SC-2, SC-3, SC-4 + InScope-1, InScope-2)
- SC-1 + SC-3: Assume pure web feature (SC-1 satisfied). Now assume mixed feature with same project (SC-3 satisfied). Both can coexist — different features in the same project. **No contradiction.**
- SC-4 + SC-1: Assume SC-4 satisfied (frontend-only feature skips gen-contracts). SC-1 also satisfied (pure web feature). These are the same scenario described at different granularity. **No contradiction, but SC-4 is subsumed by SC-1/SC-2 — it adds project-level context that SC-1 doesn't test.**
- InScope-2 (dependency chain adjustment) ↔ SC-3: If gen-contracts is NOT skipped (mixed feature), the dependency chain stays as-is. This is compatible. **No contradiction.**

**Cluster 2: gen-test-scripts** (SC-5, SC-6, SC-7 + InScope-3, InScope-4)
- SC-5 + SC-7: Assume SC-5 satisfied (direct path generates scripts for web journey). SC-7 satisfied (API/CLI/TUI pipeline unchanged). These target different surface types. **No contradiction.**
- SC-6 + SC-5: Assume SC-5 satisfied (scripts generated). SC-6 satisfied (coverage self-check passes). Compatible. But: what if SC-5 generates scripts for the WRONG journey? SC-6 checks count and type match, which should catch this. **No contradiction, but SC-6's "测试类型与 surface type 匹配" verification depends on a mapping that is not specified in the proposal.**
- InScope-3 (前置条件路由) ↔ SC-7: The prerequisite routing change modifies gen-test-scripts behavior. SC-7 requires "已有 API/CLI/TUI 流水线行为不变." If the routing logic is implemented as an if-else (check surface_type → branch), the API/CLI/TUI path is the unchanged branch. **Compatible, but the pre-revised InScope-3 description now references "Step 2" and "Step 2.5" modification — this implies touching shared code that both paths use, creating regression risk that SC-7 must guard against.**

**Cluster 3: types/web.md, types/mobile.md** (InScope-5)
- No dedicated SC for this In Scope item. **Gap, not contradiction.**

**Result**: No formal contradictions. Two gaps: (1) SC-4 is arguably subsumed by SC-1/SC-2, (2) types/web.md changes lack SC coverage.

---

## Dimension Scores (Phase 2: Rubric Scoring)

### 1. Problem Definition (93/110)

**Problem stated clearly (38/40)**:
The core problem is precise: gen-contracts produces no output for pure Web/Mobile journeys, gen-test-scripts silently skips them, and no error is reported. The "结构性断层" framing is unambiguous. Two readers would converge on the same interpretation. Minor deduction: "交互级" vs. "协议级" terminology is introduced in the Problem section but only fully defined later in Innovation Highlights. A brief parenthetical in the Problem section itself would eliminate ambiguity.

**Evidence provided (35/40)**:
Concrete and quantified: "milestone-map feature 定义了 7 个旅程，其中 5 个配置为 surface_types: ["web"]...5 个纯 Web 旅程（71%）零测试覆盖——且流水线未报告任何异常." This is verifiable against a real feature. Deduction: only one data point is provided. A second feature exhibiting the same gap would strengthen the evidence base and demonstrate the problem is systemic, not feature-specific.

**Urgency justified (20/30)**:
"每新增一个 Web-only feature 都会重复此问题" establishes the compounding cost. "测试覆盖率缺口对质量门禁不可见" highlights the hidden risk. However, the urgency section lacks concrete timeline pressure: no indication of how many Web features are planned, what sprint/quarter they fall in, or what the cost of delay is in quantitative terms. The argument is directionally correct but not quantified.

### 2. Solution Clarity (95/120)

**Approach is concrete (38/40)**:
The two-layer solution is clearly described with specific component names: (1) `CondHasProtocolSurfaceTask` in PipelineRegistry, (2) direct path in gen-test-scripts via surface_types routing. The pre-revised annotation on line 24 adds an important architectural principle ("Skill 层是权威决策者") that clarifies the relationship between layers. A reader can explain back what will be built. Deduction: the direct path's internal transformation logic (how journey.md steps map to test actions without Contract's structured fields) remains at a high level.

**User-facing behavior described (28/45)**:
The proposal describes pipeline behavior (tasks generated/skipped) and system outcomes (test scripts exist), but NOT the developer's observable experience. What does the developer see differently? A changed task list? A log message explaining the skip? A coverage report? The pre-revised InScope-3 and InScope-4 describe internal mechanics (Step 2/2.5 routing, coverage self-check logic) but not the developer-facing UX. SC entries are verification criteria, not user-facing behavior descriptions. This remains the largest gap in Solution Clarity.

**Technical direction clear (29/35)**:
Pre-revised InScope-3 significantly improves technical direction: "Step 2 原逻辑'读取 contract 文件作为输入'，改为先检查 journey 的 surface_types" and "Step 2.5 的 contract 解析步骤对 interaction-level journey 完全跳过，后续步骤使用 journey.md 中的步骤描述替代 contract 的输入输出定义." This gives implementers enough to proceed. Deduction: the replacement data source for Contract's structured fields (step-action, fixture_spec, Outcome) is mentioned at a high level ("journey.md 中的步骤描述替代 contract 的输入输出定义") but the mapping between Contract's structured schema and journey.md's narrative format is not specified. This is a design detail that may require resolution before implementation.

### 3. Industry Benchmarking (80/120)

**Industry solutions referenced (20/40)**:
"多数测试框架（Playwright、Cypress、Maestro）直接从用户场景描述生成测试脚本，不经过中间契约层." Three real tools are named but their approaches are summarized in a single sentence. No description of how Playwright Codegen records user interactions, how Maestro defines flows, or how Cypress handles test generation from specs. No documentation links, no pattern names (e.g., Model-Based Testing, Keyword-Driven Testing). This is name-dropping, not industry benchmarking.

**At least 3 meaningful alternatives (24/30)**:
Five alternatives including "do nothing." The comparison table is clear. The "仅流水线跳过" alternative is genuinely different from the selected approach. Deduction: "泛化 Contract 模型" ("Web Contract = 重写测试脚本，中间层无信息增益") and "仅加覆盖率审计" ("只诊断不治疗") are dismissed in single phrases that read as straw men — they are presented primarily to be rejected rather than genuinely evaluated.

**Honest trade-off comparison (18/25)**:
The Cons column for the selected approach ("双路径维护") is honest. However, the Cons for "泛化 Contract 模型" ("中间层无信息增益") is the proposal's own thesis presented as an objective measurement. The trade-off comparison is somewhat circular — alternatives are evaluated using the proposal's own framework rather than independent criteria.

**Chosen approach justified against benchmarks (18/25)**:
The proposal justifies the selected approach as the only one covering all scenarios. The Innovation Highlights provide the conceptual justification. However, no industry-validated alternative was given serious, deep consideration. A genuine comparison with how Playwright/Cypress handles the pipeline from user story to automated test would strengthen this section.

### 4. Requirements Completeness (85/110)

**Scenario coverage (34/40)**:
S1-S4 cover the primary scenarios well. The pre-revised InScope-1 explicitly addresses the edge case of surface-type field missing/empty ("采用保守策略——视为'可能有协议级'，不触发跳过"), which was a gap in the baseline. This is a genuine improvement. Deduction: missing edge cases still remain:
- What happens when a task has an unrecognized/invalid surface-type value (e.g., typo "desktop" instead of "web")?
- What happens when gen-contracts is correctly skipped but the direct path fails mid-generation?

**Non-functional requirements (25/40)**:
No dedicated NFR section. Backward compatibility is implicitly addressed by SC-7 but not as an NFR. Observability (how developers know gen-contracts was skipped) is not discussed. Performance, accessibility, and security are not mentioned. The absence of NFRs is a persistent gap.

**Constraints & dependencies (26/30)**:
Two clear constraints. The pre-revised InScope-1 adds explicit handling for missing fields ("只有明确的 web/mobile 才触发跳过"), which strengthens the constraint definition. Deduction: the dependency on `surface-type` being correctly filled is acknowledged but the proposal doesn't specify what valid values are (is there a schema? an enum?).

### 5. Solution Creativity (72/100)

**Novelty over industry baseline (30/40)**:
The core insight — distinguishing protocol-level from interaction-level execution models and routing accordingly — is a genuine conceptual contribution. The "无信息增益" argument is the creative leap. The pre-revised annotation on line 24 strengthens this by clarifying the architectural principle (Skill layer as authority). This is not copying an industry solution but rather explaining WHY industry tools skip intermediary layers for UI tests.

**Cross-domain inspiration (18/35)**:
The proposal draws on Forge's internal architecture rather than external domains. No references to how other domains (compilers with IR optimization passes, network protocol stacks, data pipeline routing) handle similar "skip unnecessary transformation stages" decisions. The insight is internally generated, not cross-pollinated.

**Simplicity of insight (24/25)**:
The "why didn't I think of that" quality is strong: if journey.md already contains all information needed, why force it through an intermediary? The pre-revised architectural principle (Skill layer as authority, Pipeline skip as optimization) simplifies the mental model. One deduction: the two-layer implementation still adds branching complexity that a purely skill-driven approach might avoid.

### 6. Feasibility (84/100)

**Technical feasibility (33/40)**:
Three existing mechanisms cited: `UISurfaceOnly` pattern, `surface-type` field, existing types/web.md rules. The pre-revised InScope-3 adds concrete modification points (Step 2, Step 2.5), which improves feasibility assessment. Deduction: the mapping from journey.md's narrative format to the structured data that test generation needs (replacing Contract's step-action, fixture_spec, Outcome) remains an implementation unknown.

**Resource & timeline feasibility (25/30)**:
No timeline or resource estimate provided. The 5 In Scope items suggest manageable scope. Deduction: without explicit timeline, feasibility is inferred rather than demonstrated.

**Dependency readiness (26/30)**:
Both dependencies stated as "已就绪" with evidence. Strong. Deduction: "gen-test-scripts 的 types/web.md 已有 Web 特定规则，只需扩展" — the freeform review correctly notes that current types/web.md rules are designed around Contract-driven generation, so "只需扩展" may understate the required changes.

### 7. Scope Definition (68/80)

**In-scope items are concrete (26/30)**:
Five In Scope items naming specific components. The pre-revised InScope-1, InScope-3, and InScope-4 are significantly more detailed than the baseline, providing specific modification points (Step 2, Step 2.5, surface-type field handling, coverage self-check logic). Most are deliverable-level. Deduction: "覆盖率自检" still doesn't specify where the check logic lives (in SKILL.md? a separate validation step?).

**Out-of-scope explicitly listed (20/25)**:
Five Out of Scope items named clearly. TUI surface handling exclusion is explicit. Deduction: "TUI surface 处理方式变更（保持 Contract 路径）" excludes a change without explaining the rationale in this section (it's explained in Innovation Highlights).

**Scope is bounded (22/25)**:
The 5+5 split suggests a bounded effort. The pre-revised additions provide more detail without expanding scope. No timeline estimate, but the scope feels achievable based on the specific modifications described.

### 8. Risk Assessment (72/90)

**Risks identified (25/30)**:
Three meaningful risks. The pre-revised Risk 1 now includes "双重兜底" — both Pipeline conservative strategy AND Skill-layer coverage self-check with type matching. This is significantly improved over the baseline's single mitigation. Deduction: missing risks:
- Regression risk from modifying gen-test-scripts Step 2/2.5 (shared code path)
- Risk of journey.md format being insufficient for quality test generation (the core assumption)

**Likelihood + impact rated (20/30)**:
Ratings provided (M/H, M/M, L/M). The pre-revised Risk 1 mitigation is more detailed, but the likelihood rating for "业务任务 surface-type 未填充或填写错误" remains Medium — this is the proposal's critical dependency, and if the field is user-provided, a Medium likelihood seems optimistic.

**Mitigations are actionable (27/30)**:
- Risk 1: Pre-revised mitigation is now highly actionable: "(1) Pipeline 层对缺失/空字段采用保守策略不跳过；(2) Skill 层覆盖率自检按 surface type 分别验证...还验证测试类型与 surface type 匹配." Two concrete actions. **Significant improvement.**
- Risk 2: "types/web.md 补充充分的生成规则" — "充分的" remains vague but is partially addressed by the pre-revised InScope-3 which specifies what types/web.md needs to contain.
- Risk 3: "两种路径的分支逻辑清晰" — still more of a description than a mitigation action, but the pre-revised architectural principle on line 24 ("Skill 层是权威决策者") provides a stronger mitigation rationale.

### 9. Success Criteria (68/80)

**Criteria are measurable and testable (26/30)**:
SC-1 through SC-4 are testable via task list inspection. SC-5 is testable by checking script generation. The pre-revised SC-6 now includes "并验证测试脚本类型与 surface type 匹配（如 web journey 必须对应 Web E2E Test 而非 API Functional Test）" — this is a significant improvement, making the criterion more precise and testable. Deduction: SC-5's "不报错、不跳过" remains a negative criterion (absence of behavior) without a positive complement (what DOES the generated script contain?).

**Coverage is complete (18/25)**:
SC entries cover both layers and main scenarios. The pre-revised InScope-2 adds "需在 SC-3 和 SC-7 中验证这两种路径均能正确执行," strengthening the SC↔Scope linkage. Deduction: No SC for types/web.md/types/mobile.md changes. No SC for the dependency chain adjustment behavior itself (only the downstream effect).

**SC internal consistency (24/25)**:
The SC Consistency Deep-Dive (Phase 1) found no contradictions. SC-6's pre-revised type-matching check is internally consistent with SC-5's direct path generation. The self-assessed "consistency_check_result: pass" with 21 pairs is plausible given the analysis. Deduction: SC-4 is arguably subsumed by SC-1/SC-2 (same behavior at project-level vs. feature-level context), which means SC-4 may not add independent verification value.

### 10. Logical Consistency (80/90)

**Solution addresses the stated problem (32/35)**:
The two-layer solution directly addresses both symptoms: missing gen-contracts output and missing test scripts for web/mobile journeys. The pre-revised architectural principle ("Skill 层是权威决策者") strengthens the logical foundation. The "71% coverage gap" would be fully addressed. Deduction: the solution addresses "no test scripts" but not "silent failure" — the pipeline still doesn't communicate to the developer WHY gen-contracts was skipped. The observability gap persists.

**Scope ↔ Solution ↔ Success Criteria aligned (24/30)**:
Pre-revised InScope-1/3/4 significantly improve alignment by adding specific modification points and verification references. The SC→InScope traceability is stronger. Deduction:
- InScope-5 (types/web.md/types/mobile.md) still lacks dedicated SC
- InScope-2 (dependency chain) references SC-3/SC-7 but no independent SC validates the dependency resolution behavior itself

**Requirements ↔ Solution coherent (24/25)**:
S1-S4 map cleanly to the two-layer solution. No orphan requirements or unrequired solution features. The mapping is tight. Minor deduction: the "Assumptions Challenged" section is strong methodologically but doesn't feed back into the requirements — overturned assumptions should generate new requirements or constraints.

---

## Blindspot Hunt (Phase 3)

### [blindspot] B1: Missing rollback plan
The proposal has no rollback strategy. If the direct path generates low-quality test scripts in production use, how does a team revert? The gen-test-scripts direct path is a permanent code change to SKILL.md's Step 2/2.5 logic. Quote: "gen-test-scripts SKILL.md：前置条件路由修改——Step 2 原逻辑'读取 contract 文件作为输入'，改为先检查 journey 的 surface_types." This routing change is a one-way door — there is no feature flag, gradual rollout, or kill switch described. If the direct path produces poor results in practice, reverting requires a code change to SKILL.md, not a configuration toggle.

### [blindspot] B2: Core assumption treated as fact
The proposal's foundational claim — "journey.md 已包含生成所需全部信息" — is stated as fact without evidence. Quote from Innovation Highlights: "输入是用户动作、输出是视觉状态，journey.md 已包含生成所需全部信息，中间层无信息增益." This assertion is logically argued (the distinction between structured data I/O and user action → visual state is sound) but never empirically validated. The freeform review correctly identifies that journey.md's format may not contain structured fields equivalent to Contract's step-action, fixture_spec, and Outcome definitions. The pre-revised InScope-3 partially addresses this by specifying that "后续步骤使用 journey.md 中的步骤描述替代 contract 的输入输出定义," but this describes an intention, not a validated capability. The entire direct path depends on this assumption being correct.

### [blindspot] B3: SC-6 type-matching validation depends on unspecified mapping
The pre-revised SC-6 adds "验证测试脚本类型与 surface type 匹配（如 web journey 必须对应 Web E2E Test 而非 API Functional Test）." Quote: "如 web journey 必须对应 Web E2E Test 而非 API Functional Test." This validation requires a surface-type → test-type mapping table, but the proposal does not specify this mapping. Where is it defined? Is it in gen-test-scripts SKILL.md? In a configuration file? If the mapping is hardcoded in the self-check logic, it becomes a maintenance burden; if it's configurable, the configuration needs to be specified. The SC is testable in principle but the mechanism is undefined.

---

## Bias Detection Report

Pre-revised annotations appear on: line 24 (Proposed Solution), line 92 (InScope-1), line 93 (InScope-2), line 94 (InScope-3), line 95 (InScope-4), line 110 (Risk 1). These represent 6 annotated paragraphs.

Unannotated paragraphs: ~26 paragraphs across all sections.

- **Annotated regions**: 4 attack points / 6 paragraphs = density 0.67
- **Unannotated regions**: 12 attack points / 26 paragraphs = density 0.46
- **Ratio (annotated/unannotated)**: 1.46

The annotated regions show higher attack density, which is expected — these regions were revised to address specific issues, and the revisions introduced new claims and details that invite scrutiny. The density ratio of 1.46 does not indicate severe bias; annotated regions legitimately contain more novel claims that require verification. No `conflict-with-pre-revision` tags needed — all attacks on annotated regions align with the pre-revision direction (the revisions improved the document but did not fully resolve the underlying concerns).

---

## Summary

| Dimension | Score | Max |
|-----------|-------|-----|
| Problem Definition | 93 | 110 |
| Solution Clarity | 95 | 120 |
| Industry Benchmarking | 80 | 120 |
| Requirements Completeness | 85 | 110 |
| Solution Creativity | 72 | 100 |
| Feasibility | 84 | 100 |
| Scope Definition | 68 | 80 |
| Risk Assessment | 72 | 90 |
| Success Criteria | 68 | 80 |
| Logical Consistency | 80 | 90 |
| **Total** | **797** | **1000** |
