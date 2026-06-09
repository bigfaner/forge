# CTO Scorer Report — Proposal: Skip gen-contracts for Interaction-Only Features

**Iteration**: 1 (baseline)
**Date**: 2026-06-09
**Scorer**: CTO adversary

---

## Dimension Scores

### 1. Problem Definition (92/110)

**Problem stated clearly (37/40)**:
The core problem is unambiguous: gen-contracts produces no usable output for pure Web/Mobile journeys, gen-test-scripts silently skips them, and the pipeline reports no errors. The "structural gap" framing is precise — two readers would interpret the problem the same way. Minor deduction: the term "交互级 journey" is introduced without definition in the Problem section itself (it's clarified later in Innovation Highlights). A reader unfamiliar with Forge internals might not immediately grasp the distinction between "协议级" and "交互级" from the Problem section alone.

**Evidence provided (35/40)**:
Concrete evidence: "milestone-map feature 定义了 7 个旅程，其中 5 个配置为 surface_types: ["web"]...只有 2 个双 surface（web+api）旅程获得了测试脚本，5 个纯 Web 旅程（71%）零测试覆盖". This is quantified, verifiable, and specific to a real feature. Deduction: only one data point (milestone-map) is provided. A second example from a different feature would strengthen the evidence base.

**Urgency justified (20/30)**:
"每新增一个 Web-only feature 都会重复此问题" and "测试覆盖率缺口对质量门禁不可见" establish the growing cost. However, the urgency section is thin — it does not quantify the timeline (how many Web features are planned? what's the current sprint velocity?) or compare the cost of fixing now vs. after N more features accumulate. The urgency argument is directionally correct but lacks concrete timeline pressure.

### 2. Solution Clarity (88/120)

**Approach is concrete (35/40)**:
Two-layer solution is clearly described: (1) PipelineRegistry skip condition via `CondHasProtocolSurfaceTask`, (2) gen-test-scripts direct path for web/mobile journeys. A reader can explain back what will be built. Deduction: the direct path's internal mechanics are vague — "直接从 journey.md + types/web.md 生成测试脚本" doesn't describe how journey.md steps map to test actions without the Contract intermediary's structured fields.

**User-facing behavior described (30/45)**:
The proposal describes pipeline behavior (tasks generated/skipped) and system outcomes (test scripts exist/missing) but not end-user observable behavior. What does the developer see differently? A changed task list? A new log message? The proposal focuses on internal mechanism rather than the developer's experience. SC-1 through SC-7 are verification criteria, not user-facing behavior descriptions.

**Technical direction clear (23/35)**:
`CondHasProtocolSurfaceTask` as a `GenerateCondition` function and `ResolveUpstream` dependency adjustment provide reasonable technical hints. However, the gen-test-scripts internal modification is underspecified: the current prerequisite check ("必须有 Contract 文件") needs to become conditional, and Step 2 (Read Contract Specifications) needs an alternative data source in the direct path. The proposal acknowledges this in In Scope ("前置条件路由") but doesn't give enough technical direction for the Step 2 replacement. The freeform review correctly identifies this gap.

### 3. Industry Benchmarking (78/120)

**Industry solutions referenced (20/40)**:
"多数测试框架（Playwright、Cypress、Maestro）直接从用户场景描述生成测试脚本，不经过中间契约层" — three real tools are named but their approaches are summarized in a single sentence with no depth. No specific pattern name, no reference to documentation, no description of how Playwright codegen or Maestro flow recording actually works. This is name-dropping, not industry benchmarking.

**At least 3 meaningful alternatives (22/30)**:
Five alternatives are presented including "do nothing". The comparison table covers distinct approaches. Deduction: "泛化 Contract 模型" and "仅加覆盖率审计" read as straw men — each has a single-sentence dismissal ("Web Contract = 重写测试脚本，中间层无信息增益", "只诊断不治疗"). The "仅流水线跳过" alternative is the most genuine non-selected option, as it represents a partial implementation of the selected solution.

**Honest trade-off comparison (18/25)**:
The Cons column for the selected approach ("双路径维护") is honest. However, other alternatives' cons are stated without evidence: "泛化 Contract 模型" cons is "中间层无信息增益" which is the proposal's thesis, not an objectively measured trade-off. The comparison is somewhat circular — alternatives are rejected based on the proposal's own premise.

**Chosen approach justified against benchmarks (18/25)**:
The proposal justifies the selected approach as the only one covering all scenarios. The innovation highlight section provides the conceptual justification (protocol vs. interaction execution models). However, no industry-validated alternative was given serious consideration — the proposal essentially compares against self-generated options. An industry-validated alternative (e.g., how Cypress/Playwright handles test generation from user stories in a CI/CD pipeline) would strengthen the comparison.

### 4. Requirements Completeness (78/110)

**Scenario coverage (30/40)**:
Four scenarios (S1-S4) cover pure web, pure mobile, mixed surface, and multi-surface project with frontend-only change. These are the primary happy paths. Deduction: missing edge cases:
- What happens when `surface-type` field is missing or empty on some tasks?
- What happens when a task has an unrecognized surface-type value?
- Error scenarios: what if gen-contracts is incorrectly skipped but later stages need it?

**Non-functional requirements (22/40)**:
No NFRs are mentioned. For a pipeline change: performance impact of the new condition check (negligible but unstated), backward compatibility (SC-7 addresses regression but not as an NFR), observability (how will developers know the skip happened?), and test execution time impact of the direct path vs. contract path. The absence of any NFR section is a significant gap.

**Constraints & dependencies (26/30)**:
Two constraints are clearly stated: (1) `surface-type` field must exist on business tasks, (2) types/web.md and types/mobile.md need direct generation rules. Both are concrete. Deduction: the dependency on `surface-type` being correctly filled is acknowledged but the proposal doesn't describe what "correctly" means — is there a schema validation? What values are valid?

### 5. Solution Creativity (68/100)

**Novelty over industry baseline (28/40)**:
The core insight — distinguishing protocol-level (structured data I/O) from interaction-level (user action → visual state) execution models and routing accordingly — is a genuine contribution. It's not copying an industry solution but rather explaining WHY industry tools don't use an intermediary layer for UI tests. The "无信息增益" argument is the creative leap.

**Cross-domain inspiration (20/35)**:
The proposal draws on Forge's own internal architecture (PipelineRegistry conditions, surface types) rather than cross-domain inspiration. No references to how other domains (compilers, network protocols, data pipelines) handle similar routing decisions based on data characteristics.

**Simplicity of insight (20/25)**:
The "why didn't I think of that" quality is present: if journey.md already contains all information needed, why force it through an intermediary? The dual-path solution is straightforward once the insight is accepted. Deduction: the two-layer solution (pipeline skip + skill direct path) adds complexity that a single-layer solution might avoid. The freeform review's suggestion to make the skill layer self-sufficient and the pipeline skip purely optimization would be simpler.

### 6. Feasibility (82/100)

**Technical feasibility (32/40)**:
Three existing mechanisms are cited as evidence: `UISurfaceOnly` pattern in PipelineRegistry, `surface-type` field on tasks, and existing types/web.md rules. These are concrete and verifiable. Deduction: the gen-test-scripts direct path requires modifying Step 2 (Read Contract Specifications) to work without contracts, and the proposal doesn't address how journey.md's format maps to the structured data that Step 2 currently extracts from contracts. This is an implementation unknown, not necessarily a showstopper.

**Resource & timeline feasibility (25/30)**:
No timeline or resource estimate is provided. The scope is bounded (5 In Scope items) which suggests a manageable effort. Deduction: without explicit timeline, this is inferred rather than demonstrated.

**Dependency readiness (25/30)**:
Both dependencies are stated as "已就绪" with specific evidence. This is strong. Deduction: "gen-test-scripts 的 types/web.md 已有 Web 特定规则，只需扩展" — "只需扩展" is vague about how much extension is needed. The freeform review notes that the current types/web.md rules are designed around Contract-driven generation, so the extension may be non-trivial.

### 7. Scope Definition (62/80)

**In-scope items are concrete (24/30)**:
Five In Scope items, each naming a specific file or component (PipelineRegistry, gen-test-scripts SKILL.md, types/web.md, types/mobile.md). Most are deliverable-level. Deduction: "gen-test-scripts 覆盖率自检" is described as a check but doesn't specify where the check lives (in SKILL.md? a new script?). "前置条件路由" is a behavior description, not a deliverable.

**Out-of-scope explicitly listed (20/25)**:
Five Out of Scope items are named, including gen-contracts skill, Contract format, gen-journeys, run-tests, and TUI surface handling. Clear and specific. Deduction: "TUI surface 处理方式变更（保持 Contract 路径）" is interesting — it implies the proposal considered TUI and decided to exclude it, but doesn't explain the exclusion rationale in the Out of Scope section.

**Scope is bounded (18/25)**:
The 5 In Scope + 5 Out of Scope items suggest a bounded effort. However, no timeline estimate or sprint allocation is provided. The scope feels achievable but the proposal doesn't demonstrate it's bounded by time or resources.

### 8. Risk Assessment (55/90)

**Risks identified (20/30)**:
Three risks are listed: surface-type field errors, direct path test quality, and dual-path maintenance cost. These are meaningful risks. Deduction: missing risks include:
- Regression risk from changing ResolveUpstream dependency chain
- Risk of types/web.md scope creep (identified in freeform review)
- Risk of coverage self-check false positives in multi-surface scenarios (identified in freeform review)

**Likelihood + impact rated (18/30)**:
Likelihood ratings (M, M, L) and impact ratings (H, M, M) are provided. The assessment is reasonable but not deeply justified — why is "业务任务 surface-type 未填充或填写错误" only Medium likelihood when the proposal depends entirely on this field? If the field is user-provided, the likelihood of errors seems higher.

**Mitigations are actionable (17/30)**:
- Risk 1 mitigation: "覆盖率自检硬失败兜底" — actionable but the freeform review correctly identifies this only catches total absence, not wrong type. The mitigation addresses a subset of the risk.
- Risk 2 mitigation: "types/web.md 补充充分的生成规则" — "充分的" is vague. What constitutes "sufficient"?
- Risk 3 mitigation: "两种路径的分支逻辑清晰" — this is a description, not a mitigation action. No specific action to prevent maintenance burden.

### 9. Success Criteria (55/80)

**Criteria are measurable and testable (22/30)**:
SC-1 through SC-4 are testable via task list inspection (presence/absence of T-test-gen-contracts). SC-5 is testable by checking test script generation. SC-6 is testable by verifying the FAIL behavior. SC-7 is testable via regression. Deduction: SC-5's "不报错、不跳过" is a negative criterion (absence of behavior) — it should be complemented with a positive criterion (what DOES the generated script contain?). SC-6's "输出缺口列表" is measurable but the format of the gap list is unspecified.

**Coverage is complete (18/25)**:
SC entries cover the main scenarios (S1-S4) and both layers (pipeline skip + skill direct path). Coverage self-check has its own SC (SC-6). Regression has SC-7. Deduction: No SC for the types/web.md and types/mobile.md changes mentioned in In Scope. No SC for the dependency chain adjustment. No SC for the "前置条件路由" in In Scope. These In Scope items lack corresponding verification criteria.

**SC internal consistency (15/25)**:
No direct contradictions found between SC entries. However, the consistency check is superficial because several SC entries are imprecise:
- SC-3 says "流水线正常生成 T-test-gen-contracts" for mixed features but doesn't specify behavior for the web-only journeys WITHIN that mixed feature
- SC-4 and SC-3 interact (same project, different features) but the proposal doesn't clarify whether these are tested independently or in combination
- The proposal's own consistency_check_result claims "pass" with 21 pairs checked, but this self-assessment cannot be independently verified

### 10. Logical Consistency (75/90)

**Solution addresses the stated problem (30/35)**:
The two-layer solution directly addresses both symptoms: (1) pipeline skip eliminates the useless gen-contracts execution, (2) direct path ensures web/mobile journeys get test scripts. The "71% coverage gap" would be fully addressed. Deduction: the solution addresses the "no test scripts" symptom but not the "silent failure" aspect — the pipeline still doesn't report that it's skipping gen-contracts. A developer running the pipeline should ideally see a log message explaining why gen-contracts was skipped.

**Scope <-> Solution <-> Success Criteria aligned (22/30)**:
General alignment exists: In Scope items map to solution layers, SC entries verify In Scope items. However:
- In Scope "前置条件路由" has no dedicated SC
- In Scope "types/web.md / types/mobile.md：补充 journey.md 直达生成规则" has no dedicated SC
- SC-4 (multi-surface project) tests project-level behavior but the scope is defined at feature level

**Requirements <-> Solution coherent (23/25)**:
S1-S4 map cleanly to the two-layer solution. No orphan requirements or solution features without requirements. The mapping is tight and well-structured.

---

## Blindspot Hunt

### [blindspot] B1: Missing rollback plan

The proposal has no rollback strategy. If the direct path generates low-quality test scripts in production, how does a team revert? The PipelineRegistry skip is a conditional feature that presumably can be disabled, but the gen-test-scripts direct path is a permanent code change. Quote: "gen-test-scripts SKILL.md：前置条件路由（protocol-level 需 Contract，interaction-level 跳过）" — this routing is a one-way door in the skill's logic. There's no feature flag, no gradual rollout, no kill switch described.

### [blindspot] B2: Silent behavior change for developers

The proposal changes pipeline behavior without discussing developer communication. Today, developers see gen-contracts in the task list for all features. After this change, some features won't have gen-contracts tasks. Quote: "纯 Web feature（所有业务任务 surface-type=web）流水线不生成 T-test-gen-contracts 和 T-eval-contract 任务" — developers accustomed to seeing these tasks may be confused when they disappear. No migration guide, no changelog entry, no log message is specified.

### [blindspot] B3: Assumption that journey.md contains sufficient information for test generation

The proposal's core assumption is "journey.md 已包含生成所需全部信息" but this is stated as fact without evidence. Quote from Innovation Highlights: "输入是用户动作、输出是视觉状态，journey.md 已包含生成所需全部信息，中间层无信息增益". The freeform review correctly identifies that journey.md's format may not contain structured fields equivalent to Contract's step-action, fixture_spec, and Outcome definitions. This unverified assumption underpins the entire direct path.

---

## Summary

| Dimension | Score | Max |
|-----------|-------|-----|
| Problem Definition | 92 | 110 |
| Solution Clarity | 88 | 120 |
| Industry Benchmarking | 78 | 120 |
| Requirements Completeness | 78 | 110 |
| Solution Creativity | 68 | 100 |
| Feasibility | 82 | 100 |
| Scope Definition | 62 | 80 |
| Risk Assessment | 55 | 90 |
| Success Criteria | 55 | 80 |
| Logical Consistency | 75 | 90 |
| **Total** | **733** | **1000** |
