# CTO Scorer Report — Proposal: Skip gen-contracts for Interaction-Only Features

**Iteration**: 2
**Date**: 2026-06-09
**Scorer**: CTO adversary

---

## Pre-Score Anchors (Phase 1: Reasoning Audit)

### Problem → Solution trace
The stated problem is: gen-contracts produces no output for pure Web/Mobile journeys, gen-test-scripts silently skips them, pipeline reports no error, resulting in 71% test coverage gap. The proposed solution addresses both symptoms: (1) PipelineRegistry skip eliminates useless gen-contracts execution, (2) gen-test-scripts direct path ensures web/mobile journeys get test scripts. The mapping is direct and honest. The "Skill 层是权威决策者" principle (line 24) clarifies the decision hierarchy.

### Solution → Evidence trace
The "protocol vs. interaction execution model" distinction is used to justify the dual-path solution. The "经验验证承诺" (line 45) is a significant addition: the proposal now explicitly acknowledges that "journey.md 包含生成所需全部信息" is an unvalidated hypothesis and commits to verifying it on 3 existing web-only journeys before implementation. This transforms an assertion into a testable assumption. However, the mapping itself — "步骤段落 → step-action, 前置条件描述 → fixture_spec, 预期结果描述 → Outcome" (line 108) — is described but not demonstrated on an actual journey. The validation commitment is procedural, not evidential.

### Evidence → Success Criteria trace
SC-1 through SC-8 now cover more ground: SC-5 requires "用户动作调用和至少一个可视化断言（非空、非骨架）" (line 143), a positive quality criterion rather than a negative existence check. SC-8 directly validates types/web.md output quality. SC-6 includes type-matching. This is stronger than iteration 1. Remaining gap: no SC validates the "经验验证承诺" — the pre-implementation validation of journey.md mapping completeness. The validation is described in Innovation Highlights but has no corresponding SC. If the validation fails, the proposal says "回退到补充结构化字段方案" (line 45), but this alternative has no SC, no scope items, and no feasibility assessment.

### Self-contradiction check
No direct contradictions found. The fallback plan (line 45: "失败则回退到补充结构化字段方案") is mentioned but not scoped — it's a named alternative without requirements, risks, or success criteria. This is not a contradiction but an incompleteness: the proposal acknowledges a failure path without committing to its details.

### SC Consistency Deep-Dive

**Cluster 1: PipelineRegistry** (SC-1, SC-2, SC-3, SC-4 + InScope-1, InScope-2)
- SC-1 + SC-3: Pure web feature skips gen-contracts; mixed feature generates gen-contracts. Different features in same project. **No contradiction.**
- SC-4 + SC-1: SC-4 is project-level context for the same behavior SC-1 tests at feature level. **No contradiction, SC-4 adds project-level context.**
- InScope-2 (dependency chain adjustment) ↔ SC-3: Mixed feature keeps original chain. **Compatible.**
- InScope-1 conservative strategy ↔ SC-1: SC-1 says "不生成 T-test-gen-contracts" for pure web. InScope-1 says "缺失、为空、或未知值...不触发跳过". Pure web has valid surface-type, so conservative strategy doesn't interfere. **Compatible.**

**Cluster 2: gen-test-scripts** (SC-5, SC-6, SC-7, SC-8 + InScope-3, InScope-4)
- SC-5 + SC-7: Direct path for web/mobile; API/CLI/TUI unchanged. Different surface types. **No contradiction.**
- SC-5 + SC-8: SC-5 validates script content (actions + assertions); SC-8 validates types/web.md output quality. These overlap — SC-5 tests output quality directly while SC-8 tests the rules that produce that output. **Redundancy, not contradiction.** SC-8 could be subsumed by SC-5.
- SC-6 + InScope-4: SC-6 requires type-matching; InScope-4 defines the mapping. InScope-4 provides the mechanism SC-6 tests. **Compatible.**
- InScope-3 (直达映射) ↔ SC-7: InScope-3 modifies Step 2 routing. SC-7 requires API/CLI/TUI path unchanged. If routing is a clean if-else branch, compatible. The InScope-3 description specifies "仅含 web/mobile 时跳过 contract" which implies the else-branch is unchanged. **Compatible.**

**Cluster 3: types/web.md, types/mobile.md** (InScope-5 + SC-8)
- SC-8 now covers types/web.md/types/mobile.md output quality. **Previously identified gap addressed.**

**Result**: No formal contradictions. One observation: SC-8 is arguably redundant with SC-5 (both test output quality at different levels of indirection).

---

## Dimension Scores (Phase 2: Rubric Scoring)

### 1. Problem Definition (98/110)

**Problem stated clearly (39/40)**:
The core problem is precise: "gen-contracts 只能为协议级 surface（API/CLI/TUI）生成有意义的契约，纯交互级 journey 无法产出契约文件，导致 gen-test-scripts 无输入可用，最终静默跳过这些 journey 的测试生成" (line 12). The "结构性断层" framing is unambiguous. The "协议级" vs. "交互级" terminology is used consistently and clarified in Innovation Highlights. Two readers would converge on the same interpretation. Minor deduction: the Problem section does not briefly define what "协议级" and "交互级" mean before using them — readers must proceed to Innovation Highlights for the definition.

**Evidence provided (38/40)**:
Concrete and quantified: "5 个纯 Web 旅程（71%）零测试覆盖——且流水线未报告任何异常" (line 16). This is verifiable against the milestone-map feature. The single data point concern from iteration 1 is noted, but the evidence is specific, measurable, and compelling for demonstrating the problem exists. Deduction: still only one feature's data provided. A second feature exhibiting the same gap would demonstrate systemic nature.

**Urgency justified (21/30)**:
"每新增一个 Web-only feature 都会重复此问题" and "测试覆盖率缺口对质量门禁不可见" (line 20) establish compounding cost and hidden risk. However, the urgency section still lacks concrete timeline pressure: no sprint/quarter context, no quantified cost of delay, no indication of how many Web features are planned. The argument is directionally correct but remains qualitative.

### 2. Solution Clarity (108/120)

**Approach is concrete (39/40)**:
The two-layer solution with explicit decision priority ("Skill 层是权威决策者", line 24) is clearly described with specific component names: (1) `CondHasProtocolSurfaceTask` in PipelineRegistry, (2) direct path in gen-test-scripts via surface_types routing. The InScope-3 mapping ("步骤段落 → step-action, 前置条件描述 → fixture_spec, 预期结果描述 → Outcome", line 108) provides concrete transformation logic. A reader can explain back what will be built. Minor deduction: the mapping template format within types/web.md is not specified — is it a Jinja template, a prompt instruction, or a rule description?

**User-facing behavior described (40/45)**:
The Developer-Facing Observability section (lines 31-36) is a significant addition that directly addresses iteration 1's largest gap. Five observable developer experiences are described: task list changes, skip logs, direct path identification, coverage reports, and failure fallback. These are specific, observable, and actionable. Deduction: the coverage report format ("Coverage: web 3/3 journeys, api 2/2 journeys") is described as a single-line example but the full output format (what file, what location, when triggered) is not specified. Also, the interaction between skip log and direct path log could be confusing if both appear for the same feature — the proposal doesn't describe the log sequence for a complete pipeline run.

**Technical direction clear (29/35)**:
InScope-3 (line 108) provides concrete transformation mapping. InScope-4 (line 109) defines the Surface → Test Type mapping table explicitly. The modification points in gen-test-scripts Step 2 are specified. However: the actual template structure within types/web.md that performs the mapping from journey.md narrative steps to structured test actions remains unspecified. The proposal says "types/web.md 定义映射模板" (line 108) but does not show or describe the template structure. This is the core implementation detail for the direct path.

### 3. Industry Benchmarking (90/120)

**Industry solutions referenced (30/40)**:
Significantly improved from iteration 1. Three real tools are now described with specific mechanisms: "Playwright Codegen：Browser Context Recording 捕获用户交互，直接生成 page.click()/page.fill() 等脚本" (line 72), "Maestro：声明式 YAML...- tapOn: 'Login' 即测试步骤" (line 73), "Cypress Studio：从交互录制生成 cy.get().click() 代码" (line 74). The common pattern extraction ("Action-Driven Testing", line 76) is a useful synthesis. Deduction: no links to documentation, no description of how these tools handle assertion generation (which is a key concern in Forge's direct path — SC-5 requires "至少一个可视化断言"). The industry benchmarking describes test generation but not assertion quality.

**At least 3 meaningful alternatives (26/30)**:
Five alternatives including "do nothing." The comparison table is clear. "仅流水线跳过" is genuinely different. Deduction: "泛化 Contract 模型" is dismissed with "Web Contract 退化为 DOM 操作+视觉断言列表，与测试脚本 1:1 重合，维护等价文档无抽象价值" (line 83) — this is the proposal's own thesis ("中间层无信息增益") presented as the reason for rejection, which is somewhat circular. "仅加覆盖率审计" is dismissed as "只诊断不治疗" (line 84) — a fair characterization. The alternatives are more honestly evaluated than in iteration 1 but the strongest alternative (泛化 Contract 模型) still receives primarily the proposal's own framing rather than independent evaluation.

**Honest trade-off comparison (17/25)**:
The Cons for the selected approach ("双路径维护") is acknowledged. However: the Cons for "泛化 Contract 模型" ("中间层无信息增益") is the proposal's central thesis presented as an objective measurement — this is still circular reasoning. The "经验验证承诺" in Innovation Highlights (line 45) actually undermines this dismissal: if the mapping completeness needs empirical validation, then the claim that Contract provides "无信息增益" is not yet proven. The trade-off comparison for the selected approach's "双路径维护" Con is described in one phrase without analysis of what maintenance cost looks like in practice (e.g., when a new surface type is added, both paths need updating).

**Chosen approach justified against benchmarks (17/25)**:
The "Action-Driven Testing" pattern synthesis (line 76) provides conceptual alignment with industry practice. "Forge 的 gen-contracts 为协议级设计，对交互级属过度设计" (line 76) is a clear positioning statement. However: no industry tool was given a genuinely deep comparison. The proposal does not analyze how Playwright's approach would map onto Forge's architecture, or what Forge could learn from Maestro's flow-as-spec pattern for its types/web.md templates. The benchmarking is descriptive but not analytical.

### 4. Requirements Completeness (95/110)

**Scenario coverage (36/40)**:
S1-S4 cover the primary scenarios. InScope-1 now explicitly addresses edge cases: "缺失、为空、或未知值（如 'desktop'）采用保守策略" (line 106). This resolves iteration 1's gap about unrecognized surface-type values. SC-5 now requires positive quality ("用户动作调用和至少一个可视化断言", line 143) rather than just existence. Deduction: one remaining gap — what happens when the direct path generates scripts but they contain syntactically invalid test code? No error scenario covers partial or broken generation output.

**Non-functional requirements (34/40)**:
Significantly improved. The NFR section (lines 63-66) now covers: observability, performance (O(n) check), backward compatibility (SC-7 reference), and reversibility (feature flag reference). The Developer-Facing Observability section provides specific UX descriptions. Deduction: the performance NFR ("surface-type 条件检查 O(n) 遍历业务任务列表") addresses only the PipelineRegistry check, not the direct path generation cost. No NFR for test script generation time or resource consumption. Security and accessibility NFRs are absent but may not be relevant in a CLI tool context.

**Constraints & dependencies (25/30)**:
Two clear constraints with evidence ("breakdown-tasks / quick-tasks 已支持", line 58). The "经验验证承诺" (line 45) effectively adds a validation prerequisite. Deduction: the dependency on `surface-type` field correctness is acknowledged (Risk 1) but the proposal doesn't specify how the field is populated — is it manual entry by the developer? If so, the "已支持" claim means the field exists but not that it's reliably filled.

### 5. Solution Creativity (75/100)

**Novelty over industry baseline (32/40)**:
The core insight — distinguishing protocol-level from interaction-level execution models and routing accordingly — remains a genuine conceptual contribution. The "无信息增益" argument is the creative leap, and the "经验验证承诺" (line 45) shows intellectual honesty about this assumption. The architectural principle ("Skill 层是权威决策者", line 24) adds clarity. The proposal is not copying an industry solution but rather explaining WHY industry tools skip intermediary layers for UI tests, then applying that reasoning to Forge's pipeline. Deduction: the two-layer implementation (Pipeline skip + Skill direct path) adds complexity that a purely skill-driven approach (always route in gen-test-scripts, remove gen-contracts from pipeline entirely) might avoid.

**Cross-domain inspiration (18/35)**:
No references to how other domains handle similar "skip unnecessary transformation stages" decisions. The insight is entirely self-generated within Forge's architecture. This is not a flaw — the solution is internally coherent — but it means the proposal doesn't leverage patterns from compiler optimization (dead code elimination passes), data pipeline routing (early filtering), or network protocol stacks (skip layers when unnecessary).

**Simplicity of insight (25/25)**:
The "why didn't I think of that" quality is strong: if journey.md already contains all information needed, why force it through an intermediary? The "Action-Driven Testing" pattern name (line 76) provides a clean abstraction. The "经验验证承诺" (line 45) adds intellectual honesty without sacrificing simplicity.

### 6. Feasibility (86/100)

**Technical feasibility (34/40)**:
Three existing mechanisms cited: `UISurfaceOnly` pattern, `surface-type` field, existing types/web.md rules. InScope-3 (line 108) provides the specific mapping: steps → step-action, preconditions → fixture_spec, expected results → Outcome. The "经验验证承诺" (line 45) commits to validating this mapping before implementation, which de-risks the core technical unknown. Deduction: the mapping template format within types/web.md remains unspecified — this is the key implementation artifact and its feasibility depends on its design, which is deferred.

**Resource & timeline feasibility (25/30)**:
No timeline or resource estimate provided. The 5 In Scope items suggest manageable scope. The "经验验证承诺" adds a pre-implementation validation step (3 web-only journeys), which is prudent but adds schedule uncertainty — if validation fails, the "补充结构化字段方案" alternative may require significantly more work than the original scope.

**Dependency readiness (27/30)**:
Both dependencies stated as "已支持" or "已有...需扩展" with evidence. The feature flag mechanism is described in the Rollback Mechanism section. Deduction: "types/web.md 已有 Web 规则，需扩展" (line 92) — the extent of "需扩展" is not quantified. Current rules are designed around Contract-driven generation; the "扩展" may be more substantial than "只需" implies.

### 7. Scope Definition (75/80)

**In-scope items are concrete (27/30)**:
Five In Scope items naming specific components. InScope-1 specifies conservative strategy for edge cases. InScope-3 provides detailed mapping. InScope-4 defines the coverage check with explicit type mapping table. These are deliverable-level. Deduction: "types/web.md 定义映射模板" (line 108) — the template itself is not specified as a deliverable artifact (it's implied by InScope-5 but not explicitly named as a template document). The "经验验证承诺" (3 journey validation) is described in Innovation Highlights but not listed as an In Scope item, creating an implicit scope item that isn't tracked.

**Out-of-scope explicitly listed (23/25)**:
Five Out of Scope items named clearly. TUI surface handling exclusion is explicit. Deduction: the "经验验证承诺" fallback ("回退到补充结构化字段方案", line 45) is an implicitly scoped alternative that is neither in-scope nor out-of-scope — it's a contingency without scope definition.

**Scope is bounded (25/25)**:
The 5+5 split suggests a bounded effort. The feature flag mechanism (Rollback Mechanism section) provides a safety valve. The "经验验证承诺" is a gated prerequisite, not an open-ended investigation.

### 8. Risk Assessment (82/90)

**Risks identified (28/30)**:
Four meaningful risks. Risk 1 (surface-type filling, now rated H/H with justification "字段新引入，历史 feature 可能未标注") is honestly assessed. Risk 2 (journey.md structured information insufficiency) directly addresses the core assumption, with "实施前验证 3 个 web-only journey 映射完整性；SC-8 验证断言质量；失败则补充结构化字段" (line 133). This is a genuine risk that was a blindspot in earlier iterations. Deduction: still missing regression risk from modifying gen-test-scripts Step 2 — this shared code path handles both contract-driven and direct-path generation, and changes to the routing logic could introduce regressions in the contract path. SC-7 guards against this but the risk itself is not identified.

**Likelihood + impact rated (25/30)**:
Ratings provided (H/H, M/H, M/M, L/M). Risk 1's upgrade from M to H is justified. Risk 2's Medium likelihood for the core assumption failure is debatable — this is the foundational assumption of the entire direct path, and it's rated Medium because the proposal argues it's "合理推断" (line 45). This is somewhat optimistic given the validation hasn't been performed yet. Risk 4 (L/M for maintenance cost) may underestimate long-term maintenance — every new surface type or pipeline change must consider both paths.

**Mitigations are actionable (29/30)**:
- Risk 1: Dual safety net (Pipeline conservative strategy + Skill coverage self-check with type matching). Concrete and actionable.
- Risk 2: "实施前验证 3 个 web-only journey" + "SC-8 验证断言质量" + "失败则补充结构化字段". Specific validation plan with fallback. Strong.
- Risk 3: "types/web.md 补充充分生成规则" — "充分" is still somewhat vague but SC-5 and SC-8 provide quality gates.
- Risk 4: "分支逻辑按执行模型分流" — adequate for the stated risk level.
Deduction: Risk 2's fallback ("补充结构化字段方案") is an unscoped alternative. If the fallback is triggered, the team has no plan, no scope, no estimate.

### 9. Success Criteria (76/80)

**Criteria are measurable and testable (28/30)**:
SC-1 through SC-4 are testable via task list inspection. SC-5 now requires positive quality: "用户动作调用和至少一个可视化断言（非空、非骨架）" (line 143) — this is a meaningful improvement from iteration 1's negative criterion. SC-6 has explicit type matching with defined mapping. SC-8 validates types/web.md output quality. Deduction: SC-5's "非空、非骨架" qualification for assertions is somewhat subjective — what qualifies as a "skeleton" assertion? A more specific criterion (e.g., "assertion references a specific page element or text content") would be less ambiguous.

**Coverage is complete (24/25)**:
SC entries now cover both layers, main scenarios, and types/web.md changes (SC-8). The SC→InScope traceability is strong. Deduction: the "经验验证承诺" (3 journey validation before implementation) has no corresponding SC. If this is a gated prerequisite, it should have a verifiable exit criterion. Also, no SC validates the feature flag mechanism described in Rollback Mechanism (line 121).

**SC internal consistency (24/25)**:
The SC Consistency Deep-Dive (Phase 1) found no contradictions. SC-8 is arguably redundant with SC-5 (both test output quality at different levels), but redundancy is not contradiction. SC-6's type-matching is consistent with InScope-4's mapping. The consistency_check_result: pass with 21 pairs is plausible. Deduction: SC-4 is arguably subsumed by SC-1/SC-2 (same behavior at project-level vs. feature-level context).

### 10. Logical Consistency (84/90)

**Solution addresses the stated problem (33/35)**:
The two-layer solution directly addresses both symptoms: missing gen-contracts output and missing test scripts for web/mobile journeys. The Developer-Facing Observability section (lines 31-36) addresses the "silent failure" concern from iteration 1 — the pipeline now communicates skip decisions via logs and coverage reports. The "Skill 层是权威决策者" principle ensures the direct path functions even if Pipeline skip fails. Deduction: the solution addresses the symptoms comprehensively, but the "经验验证承诺" introduces a conditional: the solution only works IF journey.md contains sufficient information. The proposal honestly acknowledges this condition but the logical chain has a contingent link that is not yet proven.

**Scope ↔ Solution ↔ Success Criteria aligned (26/30)**:
InScope items map cleanly to solution components. SC entries cover all InScope items (including SC-8 for InScope-5). The Developer-Facing Observability section bridges Solution and InScope. Deduction:
- The "经验验证承诺" (line 45) is described in Innovation Highlights but not reflected in Scope — it's an implicit scope item
- The Rollback Mechanism (lines 119-127) introduces a feature flag and monitoring plan that have no corresponding InScope items or SC entries
- The fallback "补充结构化字段方案" (line 45) is named but has no scope, SC, or feasibility assessment

**Requirements ↔ Solution coherent (25/25)**:
S1-S4 map cleanly to the two-layer solution. NFRs are addressed by specific sections (Observability, Performance, Backward Compatibility, Reversibility). No orphan requirements or unrequired solution features. The mapping is tight and coherent.

---

## Blindspot Hunt (Phase 3)

### [blindspot] B1: Fallback plan is named but not scoped
The proposal explicitly names a fallback: "失败则回退到补充结构化字段方案" (line 45). This fallback is referenced in Risk 2's mitigation (line 133) and in Innovation Highlights. However, this alternative has: no scope definition, no feasibility assessment, no implementation estimate, no success criteria. If the "经验验证承诺" fails (3 journeys cannot be mapped), the team must either (a) design and implement a "补充结构化字段" approach with no prior analysis, or (b) abandon the feature. The proposal acknowledges a critical decision point but provides no preparation for the unfavorable branch. This is a planning gap that no rubric dimension fully captures — Risk Assessment scores the risk, Feasibility scores the primary approach, but neither scores the fallback's readiness.

### [blindspot] B2: Rollback Mechanism introduces untracked scope
The Rollback Mechanism section (lines 119-127) introduces: a feature flag `forge.skip_contracts.enabled`, a monitoring plan tracking "直达路径脚本通过率", and a triggering threshold ("显著低于 contract 路径时触发评估"). None of these appear in the In Scope section. The feature flag requires: flag definition, flag checking in both Pipeline and Skill layers, default value configuration, and documentation. The monitoring plan requires: metric collection, comparison logic, threshold definition ("显著低于" is vague), and alert/notification mechanism. These are not trivial — they represent additional implementation work that is not tracked, estimated, or covered by SC entries. Reasoning audit flagged this independently of dimension scoring.

### [blindspot] B3: Direct path assertion quality has no baseline comparison
SC-5 requires "至少一个可视化断言（非空、非骨架）" (line 143) and SC-8 requires "有意义断言的测试脚本" (line 146). Risk 3 acknowledges "直达路径测试质量低于 Contract 驱动" (line 134). But the proposal provides no baseline for comparison: what quality do Contract-driven test scripts achieve? What is the current assertion quality for API/CLI/TUI journeys? Without a baseline, "有意义断言" and "非骨架" are subjective — the team cannot objectively determine whether direct path quality is acceptable. The monitoring plan in Rollback Mechanism (line 126) mentions "跟踪直达路径脚本通过率，显著低于 contract 路径时触发评估" but this monitors runtime pass rate, not assertion quality at generation time. Two different quality dimensions are conflated.

---

## Bias Detection Report

Pre-revised annotations from iteration 1 are no longer present — the proposal appears to be a clean revision without marked annotations. All sections are unmarked.

- **Total attack points**: 10 dimension attacks + 3 blindspot attacks = 13
- **Distribution**: Evenly spread across sections; no concentration bias detected
- **Annotated regions**: N/A (no pre-revised markers in this iteration)

No `conflict-with-pre-revision` tags needed — this is a clean document without annotation markers.

---

## Summary

| Dimension | Score | Max |
|-----------|-------|-----|
| Problem Definition | 98 | 110 |
| Solution Clarity | 108 | 120 |
| Industry Benchmarking | 90 | 120 |
| Requirements Completeness | 95 | 110 |
| Solution Creativity | 75 | 100 |
| Feasibility | 86 | 100 |
| Scope Definition | 75 | 80 |
| Risk Assessment | 82 | 90 |
| Success Criteria | 76 | 80 |
| Logical Consistency | 84 | 90 |
| **Total** | **869** | **1000** |
