# Pipeline Topology Registry — Eval Report (Iteration 2)

**Reviewer**: CTO adversary (blind to author identity)
**Date**: 2026-05-29
**Rubric**: proposal.md (1000-point scale)
**Pre-revision annotations**: None detected in document. All regions treated as unannotated.

---

## Iteration-1 Issue Resolution Audit

| # | Attack Point | Status | Evidence |
|---|-------------|--------|----------|
| 1 | ResolveFirstTestDep stage-gate/clean-code logic dropped | Resolved | New `injectCleanCodeDep` post-generation step + `ResolveLastRunTestOrBusiness` resolver. Stage-gate explicitly delegated to `GenerateStageGates`. |
| 2 | Drift/consolidate fallback unrepresented | Resolved | `ResolveLastRunTestOrBusiness` resolver handles fallback. Risk table entry added. |
| 3 | SC1 "ONE file" self-contradiction | Resolved | Rewritten to "at most TWO files". |
| 4 | Zero industry references | Resolved | Alternative 2 added (google/wire DI container, Tekton/GitHub Actions pipeline DAG comparison). |
| 5 | Alternative 1 is straw man | Resolved | Trade-off paragraph acknowledges table-driven solves ~80% and explains why full registry is worth the complexity. |
| 6 | Missing scenario: test pipeline disabled | Resolved | Dedicated scenario paragraph added. |
| 7 | GenerateTestTasks signature change undocumented | Resolved | "Breaking change: businessTasks parameter" section added with caller analysis. |
| 8 | T-clean-code reverse dependency missing | Resolved | `injectCleanCodeDep` post-generation step added. |
| 9 | Missing risk: init() panic in CI | Resolved | Risk table entry added with `--no-validate` flag mitigation. |
| 10 | Missing risk: ResolveFirstTestDep regression | Resolved | Risk table entry added (High/High). |
| 11 | "~40% chance" unsupported statistic | Resolved | Rewritten to raw data "3 caused at least one pipeline bug". |
| 12 | injectReviewDocDep undermines declarative purity | Resolved | Escape-hatch protocol defined with 4 rules. |
| 13 | No timeline/effort estimate | Resolved | "2-3 development days" with breakdown. |
| 14 | findHighestGateOrSummary dead code | Resolved | Functions Relationship table marks as "删除（dead code）". |
| 15 | SC7 is code style rule | Resolved | Replaced with SC7 "Dependency chain completeness". |

**Resolution rate**: 15/15 attacks addressed. This is a thorough revision.

---

## Phase 1: Reasoning Audit

### Problem -> Solution Chain

**Problem**: Task IDs and dependency relationships are hardcoded as string literals across 5+ code locations. Adding/removing a task type causes cascading failures.

**Solution**: Pipeline Topology Registry — a single declarative data structure that all consuming code derives from, plus two post-generation injection steps for cross-pipeline dependencies.

**Chain verdict**: The solution directly addresses the stated problem. The registry centralizes scattered definitions. The revision has addressed all three chain gaps from iteration 1:

1. **ResolveFirstTestDep's stage-gate logic**: Now explicitly delegated to `GenerateStageGates` (Out of Scope item 6). Stage-gate generation is independent of pipeline topology.
2. **T-clean-code injection**: Now handled by `injectCleanCodeDep` post-generation step.
3. **Drift fallback**: Now handled by `ResolveLastRunTestOrBusiness` resolver.

**Remaining chain concern**: The dependency chain described in the post-generation section states: `business-tasks -> T-clean-code -> T-review-doc (if exists) -> first-test-task`. However, `injectCleanCodeDep` injects T-clean-code as dependency of the *first test task*, not as an intermediate between business tasks and test pipeline. And `injectReviewDocDep` also injects T-review-doc as dependency of the *first test task*. The resulting chain for the first test task would be: `first-test-task.DependsOn = [..., T-clean-code, T-review-doc]`. This means first-test-task depends on BOTH T-clean-code AND T-review-doc in parallel, not in the serial chain described. The text claims a serial chain but the code produces a parallel fan-in. This is a documentation accuracy issue — the actual behavior may be correct but the description is misleading.

### Solution -> Evidence Chain

Evidence: 5 bug examples, 6 touch points documented. Code verification from iteration 1 confirms accuracy. "3 of 5" raw data is now used instead of percentage. Chain is solid.

### Evidence -> Success Criteria Chain

SC1 (at most TWO files) is now consistent with the solution description. SC4 (behavioral parity) is supported by the exhaustive coverage of the three ResolveFirstTestDep sub-concerns. SC7 (dependency chain completeness) is a strong proxy for correctness. Chain is sound.

### Self-Contradiction Check

1. **Post-processing dependency chain description vs. code**: The text says `business-tasks -> T-clean-code -> T-review-doc -> first-test-task` (serial chain). But `injectCleanCodeDep` prepends `T-clean-code` to `first-test-task.DependsOn`, and `injectReviewDocDep` prepends `T-review-doc` to `first-test-task.DependsOn`. The result is `first-test-task.DependsOn = [T-review-doc, T-clean-code, ...original deps...]`. This is a parallel dependency (fan-in), not a serial chain. The serial description implies T-review-doc depends on T-clean-code, but T-review-doc's registry `DependsOn` is `[ResolveDocTasks]` — it depends on doc business tasks, not T-clean-code. The serial chain description is factually incorrect.

2. **GenContext progressive population ordering dependency**: The proposal states "Order matters: nodes are processed in declaration order" and `GenContext` is "Populated progressively as nodes are processed in declaration order." T-review-doc is the FIRST node in the registry, but it depends on `ResolveDocTasks` which reads from `ctx.BusinessTasks`. Since `BusinessTasks` is a pre-populated field (not generated during iteration), this is fine. However, `ResolveUpstream` used by T-test-run reads `ctx.UpstreamIDs` which depends on having processed previous nodes. If someone reorders the registry and moves T-test-run before T-test-gen-scripts, `UpstreamIDs` would point to wrong tasks. The ordering constraint is implicit and fragile.

3. **Escape-hatch protocol cap vs. current design**: The protocol caps at 5 escape hatches. Currently there are 2 (`injectCleanCodeDep`, `injectReviewDocDep`). This seems reasonable but the cap is arbitrary — why 5? If 3 more cross-pipeline dependencies emerge, the protocol says "must extend registry expressiveness" but doesn't specify what that extension looks like.

### SC Consistency Deep-Dive

**Cluster: pipeline.go + types.go**
- SC1 (at most TWO files): Adding a task type changes pipeline.go and optionally types.go.
- In Scope item 1: Define PipelineNode, DepRef, PipelineRegistry in pipeline.go.
- **Bidirectional check**: If SC1 is satisfied, In Scope item 1 is necessarily satisfied (the registry entry in pipeline.go is exactly what In Scope item 1 defines). If In Scope item 1 is satisfied, SC1 can be satisfied for new types that don't need a new type constant (only pipeline.go). For types needing a new constant, types.go is also needed — consistent with SC1's "at most TWO files" wording. **No contradiction.**

**Cluster: GenerateTestTasks and callers**
- In Scope item 4: Refactor build.go steps 7/7.5/7.6 to use registry-driven generation.
- In Scope item 7: Two-phase validation.
- SC4 (behavioral parity): forge task index produces identical output.
- SC2 (zero panics): No findTaskIndexOrPanic calls remain.
- **Bidirectional check**: If SC4 is satisfied, the refactoring preserves behavior. If In Scope items 4 and 7 are satisfied, the code is refactored. SC2 is a necessary condition for SC4 (panics would violate identical output). **No contradiction.**

**Cluster: Dependency chain**
- SC7 (dependency chain completeness): No orphaned tasks, no dangling references.
- In Scope item 7 (Phase 2 validation): Validates resolver results exist in generated set.
- **Bidirectional check**: If SC7 is satisfied, Phase 2 validation would pass. If Phase 2 validation is implemented, it provides a mechanism to enforce SC7. **No contradiction.**

**Cluster: Cross-concern**
- SC6 (GenerateCondition correctness): T-review-doc in all three scenarios.
- T-review-doc registry entry: `ConfigGate: nil`, `GenerateCondition: CondHasDocTasks`.
- **Bidirectional check**: If the registry entry is as declared, SC6's three scenarios are logically derivable from CondHasDocTasks and the three scenarios table. **No contradiction.**

**Re-verification of revised SC entries**: The revised SC7 (dependency chain completeness) is internally satisfiable and does not contradict any In Scope item or other SC entry. No new contradictions introduced by revision.

---

## Phase 2: Rubric Scoring

### 1. Problem Definition: 95/110

| Criterion | Score | Rationale |
|-----------|-------|-----------|
| Problem stated clearly | 38/40 | Core problem is unambiguous: scattered string literals cause cascading failures. The "5+ code locations" and 6-item code smell list make this concrete. Minor gap: "5+ code locations" in Problem opening vs "6 touch points" in Code smell — inconsistency in the count. |
| Evidence provided | 38/40 | 5 concrete bugs with root causes and symptoms. Code-smell audit of 6 touch points. Raw data "3 of 5" is used instead of percentage. Solid evidence base. Only gap: no user feedback or telemetry cited — all evidence is developer-internal. |
| Urgency justified | 19/30 | "Cost of inaction" states raw data (3 of 5 caused bugs) and mentions developer velocity constraint. But still no quantified impact: how many hours per incident? How often are task types added? What is the blast radius when bugs escape to production CI? The urgency is implied by frequency, not demonstrated by impact. |

### 2. Solution Clarity: 102/120

| Criterion | Score | Rationale |
|-----------|-------|-----------|
| Approach is concrete | 38/40 | Full Go struct definitions, registry entries, resolver functions, post-generation injection code. A reader can explain back exactly what will be built. Minor gap: the generation algorithm (how `GenerateTestTasks` iterates, filters, expands, and populates `GenContext`) is described in prose without pseudocode — but the registry structure makes the algorithm derivable. |
| User-facing behavior described | 45/45 | "No change. The pipeline generates the same tasks with the same dependencies. The refactoring is purely internal architecture." Clear, honest, and complete. |
| Technical direction clear | 19/35 | The registry data structure and resolver pattern are detailed. The post-generation injection mechanism is concrete (actual Go code). However, the core generation algorithm — the `GenerateTestTasks` function that filters by Mode/ConfigGate/GenerateCondition/UISurfaceOnly, expands per-surface nodes, and progressively populates GenContext — is still described in a 4-step bullet list. The expansion algorithm (how template placeholders are substituted, how GenContext fields are updated per-node) is implicit. A reader must infer the algorithm from the data structure rather than reading it directly. The post-processing ordering dependency (injectCleanCodeDep before injectReviewDocDep) is stated but not justified — why this order matters is only explained by the serial chain description which is itself inaccurate (see Phase 1 finding). |

### 3. Industry Benchmarking: 78/120

| Criterion | Score | Rationale |
|-----------|-------|-----------|
| Industry solutions referenced | 22/40 | Alternative 2 references google/wire (DI container) and mentions Tekton Pipelines and GitHub Actions as pipeline DAGs. This is an improvement from iteration 1. However, the references are brief — the comparison is "DI doesn't fit our domain" without deeper analysis of how Tekton/GitHub Actions solve similar topology definition problems. No discussion of DAG validation libraries (e.g., hashicorp/go-tfe), pipeline DSL patterns, or rule engine approaches. |
| At least 3 meaningful alternatives | 25/30 | 4 alternatives: Table-Driven, DI Container, Do Nothing, Code Generation. All are genuinely different approaches. DI Container is industry-validated (google/wire). "Do Nothing" is included per rubric requirement. Strong. Minor gap: Alternative 3 (Do Nothing) gets only two sentences of rejection — it's present but thin. |
| Honest trade-off comparison | 15/25 | Alternative 1's trade-off is now honest (acknowledges 80% benefit). Alternative 2's comparison against pipeline DAGs is reasonable but shallow — "The registry's DepResolveFunc is closer to a pipeline DAG" is asserted without showing how Tekton/GitHub Actions differ in practice. No comparison table with weighted criteria. Alternative 4 (Code Generation) is dismissed with "no runtime benefit over native Go registry" — but code generation could provide compile-time validation that the registry's init-time validation replicates at runtime. The comparison misses this angle. |
| Chosen approach justified against benchmarks | 16/25 | The justification is that the registry provides "equivalent static validation in a structure that maps directly to task generation" vs. DI. This is reasonable but the comparison is asymmetric — DI is analyzed for misfit but pipeline DAGs (Tekton, GitHub Actions) are mentioned without being seriously evaluated. The proposal should explain why a native Go registry is superior to adopting a pipeline DAG pattern (which is what Tekton uses) beyond "maps directly to task generation." |

### 4. Requirements Completeness: 85/110

| Criterion | Score | Rationale |
|-----------|-------|-----------|
| Scenario coverage | 32/40 | Three T-review-doc scenarios documented. Test pipeline disabled / all ConfigGates off scenario added. Single-surface degeneration for per-surface-key expansion documented. Missing scenarios: (1) Empty business tasks list — what happens when `businessTasks` is empty? `ResolveLastBusinessTask` returns nil, `CondHasTestableTasks` returns false, `CondHasDocTasks` returns false. Only `CondAlways` nodes with passing ConfigGates would generate — but their dependency resolvers may return nil. (2) What happens when a `DepResolveFunc` returns nil — is the task skipped or generated with empty DependsOn? (3) Mixed mode scenario — what if `GenerateTestTasks` is called in an unexpected mode string? |
| Non-functional requirements | 32/40 | Performance addressed (~12 entries). Breaking change documented with caller analysis. Effort estimate provided. Missing: (1) Memory footprint of registry at scale — trivial but not stated. (2) Init-time validation latency — the validation runs at CLI startup but no mention of how long it takes or whether it could slow startup. (3) Backward compatibility with existing index.json files — the proposal says "identical output" but doesn't explicitly confirm that existing index.json files remain valid. |
| Constraints & dependencies | 21/30 | `forgeconfig.AutoConfig` dependency is implicit. `GenContext` progressive population ordering constraint is mentioned ("Order matters"). Go version requirements not mentioned. `CategoryForType` and `IsTestableType` stability not discussed. The `--no-validate` flag adds a new CLI interface requirement not discussed in constraints. |

### 5. Solution Creativity: 55/100

| Criterion | Score | Rationale |
|-----------|-------|-----------|
| Novelty over industry baseline | 25/40 | The `DepResolveFunc` with progressive `GenContext` is a clean design — not groundbreaking but well-suited to the domain. The escape-hatch protocol with explicit rules and caps is a thoughtful addition. The differentiation from industry patterns is now articulated (vs. DI, vs. pipeline DAGs). However, the core pattern (declarative registry + resolvers) is a well-established approach. |
| Cross-domain inspiration | 15/35 | The DI and pipeline DAG references show awareness of industry approaches. But no borrowing from compiler IR pipelines, database query planners, or rule engines that face similar problems. The escape-hatch protocol could have been inspired by compiler optimization pass managers (which have similar "pass ordering with escape hatches for cross-pass concerns"). |
| Simplicity of insight | 15/25 | "Centralize scattered definitions into a registry" is straightforward and clean. The ConfigGate/GenerateCondition/DepResolveFunc trichotomy is well-structured. However, the two post-generation injection steps with ordering dependencies, plus the escape-hatch protocol, add complexity that partially undermines the simplicity. The `GenContext` progressive population contract is a subtle ordering dependency that could confuse future maintainers. |

### 6. Feasibility: 82/100

| Criterion | Score | Rationale |
|-----------|-------|-----------|
| Technical feasibility | 35/40 | Go structs, closures, and iteration are standard. No showstopper dependencies. The ResolveFirstTestDep replacement is now fully addressed with post-generation steps. The init-time validation is straightforward. Remaining concern: `init()` validation of `DependsOn.Ref` strings that reference expanded IDs (e.g., `"T-test-gen-scripts-{surface-type}"`) — the proposal says "all DependsOn.Ref strings reference existing node IDs" but expanded IDs are templates, not concrete IDs. The validation logic must understand templates, which adds non-trivial complexity. |
| Resource & timeline feasibility | 27/30 | "2-3 development days" with itemized breakdown. This is realistic for the scope given the proposal author's familiarity with the codebase. Single PR, single review cycle is achievable. The "Update all existing tests" item is bounded by the 0.5-1 day estimate, which may be optimistic if tests require significant restructuring. |
| Dependency readiness | 20/30 | `forgeconfig.AutoConfig` is available. But `CategoryForType` and `IsTestableType` stability is still not confirmed. The `GenContext` progressive population contract requires all registry consumers to understand ordering semantics — this is an implicit knowledge dependency that is not documented as a constraint. |

### 7. Scope Definition: 72/80

| Criterion | Score | Rationale |
|-----------|-------|-----------|
| In-scope items are concrete | 27/30 | 9 numbered items, each targeting specific files and functions. Effort estimate with breakdown provided. The Functions Relationship table precisely documents each function's fate. Strong. Minor gap: Item 7 (two-phase validation) includes `init()` panic behavior — the trade-off of panic-in-CI is now in the risk table, but the scope item doesn't mention the `--no-validate` flag as a deliverable. |
| Out-of-scope explicitly listed | 22/25 | 6 items explicitly listed. Clear. Stage gate generation is properly scoped out with rationale. `findHighestGateOrSummary` is now correctly handled (marked as dead code to be deleted). |
| Scope is bounded | 23/25 | "2-3 development days" with breakdown. Single PR, single review cycle. Well-bounded. The "Update all existing tests" item is bounded by the 0.5-1 day estimate, which provides a constraint. |

### 8. Risk Assessment: 78/90

| Criterion | Score | Rationale |
|-----------|-------|-----------|
| Risks identified | 27/30 | 8 risks listed, including the two critical ones from iteration 1 (init() panic in CI, ResolveFirstTestDep regression). All major failure modes are now covered. Missing: (1) GenContext ordering dependency — if registry is reordered, progressive population breaks. (2) `DepResolveFunc` returning nil — the behavior when a resolver returns nil is unspecified in the risk table. |
| Likelihood + impact rated | 26/30 | Ratings are honest and reasonable. ResolveFirstTestDep regression at High/High is appropriately frank. init() panic at Medium/High acknowledges the trade-off. Drift fallback at Medium/Medium is fair. Minor issue: "Expanding per-surface nodes creates ID conflicts" at Low likelihood — the init-time validation mitigates this, so the rating is justified. |
| Mitigations are actionable | 25/30 | Snapshot tests are actionable. `--no-validate` flag is actionable. `forge validate-pipeline` CI step is actionable. Exhaustive snapshot comparison is specific. Minor gap: "Port all existing gotcha lessons as test cases" is still vague — which lessons? How many? Where documented? |

### 9. Success Criteria: 72/80

| Criterion | Score | Rationale |
|-----------|-------|-----------|
| Criteria are measurable and testable | 26/30 | SC1 (at most TWO files) is measurable. SC2 (zero panics) is measurable via grep. SC4 (identical output) is measurable via snapshot. SC7 (dependency chain completeness) is measurable via graph validation. SC8 (80% coverage) is measurable. Minor gap: SC6 (GenerateCondition correctness) says "T-review-doc is generated in all three scenarios" — but the three scenarios are defined in a table in the Solution section, not in the SC itself. The SC is testable but requires cross-referencing. |
| Coverage is complete | 24/25 | SC7 now covers dependency chain completeness (including orphaned tasks and dangling references). SC4 covers behavioral parity. SC2 covers panic elimination. The revised SC set covers all in-scope items. Gap: No SC for "escape-hatch count stays at or below 5" — the protocol defines a cap but no SC enforces it. |
| SC internal consistency | 22/25 | SC1 and SC4 are compatible (changing only pipeline.go/types.go while maintaining behavioral parity). SC7 and SC4 are compatible (dependency chain completeness is a necessary condition for behavioral parity). Remaining concern: SC4 (behavioral parity) requires "identical output" but the proposal introduces `injectCleanCodeDep` which adds T-clean-code as dependency of first-test-task. If the current code already does this via ResolveFirstTestDep, behavioral parity is preserved. But if the dependency chain order differs (serial vs parallel, as noted in Phase 1), the output may not be "identical" at the dependency level. This is an ambiguity, not a definitive contradiction. |

### 10. Logical Consistency: 75/90

| Criterion | Score | Rationale |
|-----------|-------|-----------|
| Solution addresses the stated problem | 33/35 | The registry centralizes scattered definitions — directly addresses the stated problem. The three chain gaps from iteration 1 are now addressed via post-generation steps and resolvers. Remaining minor gap: the serial chain description in the post-processing section is inaccurate (describes serial chain, code produces parallel fan-in). This is a documentation issue, not a logical gap — the actual behavior may be correct. |
| Scope <-> Solution <-> SC aligned | 22/30 | Scope says delete ResolveFirstTestDep, Solution provides injectCleanCodeDep + injectReviewDocDep + GenerateStageGates as replacements. SC4 (behavioral parity) is achievable if all three sub-concerns are correctly implemented. Alignment is now strong. Remaining gap: the escape-hatch protocol is defined in the Solution but has no corresponding In Scope item or SC enforcement mechanism. The protocol says "escape hatch count <= 5" but no SC enforces this. |
| Requirements <-> Solution coherent | 20/25 | The trichotomy (ConfigGate/GenerateCondition/Mode) covers most requirements. The drift fallback requirement has a corresponding resolver. The clean-code injection has a post-generation step. Remaining gap: the requirement for `DepResolveFunc` returning nil (what happens when a resolver cannot resolve) has no corresponding solution behavior described. The proposal says resolvers "return nil when the reference cannot be resolved" but doesn't specify what the generation algorithm does with nil. |

---

## Phase 3: Blindspot Hunt

### What the rubric missed:

1. **[blindspot] Serial chain description is factually inaccurate**: The post-processing section states the dependency chain is `business-tasks -> T-clean-code -> T-review-doc (if exists) -> first-test-task`. But both `injectCleanCodeDep` and `injectReviewDocDep` prepend to `first-test-task.DependsOn`. T-review-doc's own `DependsOn` in the registry is `[ResolveDocTasks]`, not `[T-clean-code]`. The actual chain is: `first-test-task` depends on `[T-review-doc, T-clean-code, ...upstream deps...]` in parallel. T-review-doc depends on `[doc business tasks]`. T-clean-code depends on `[last business task]`. The serial chain description implies T-review-doc waits for T-clean-code, but the code does not enforce this ordering. If the serial chain is the intended behavior, the code is wrong. If parallel is intended, the description is wrong. — Quote: "确保依赖链：`business-tasks -> T-clean-code -> T-review-doc（若存在）-> first-test-task`" vs. the actual `append([]string{reviewDocID}, generated[i].DependsOn...)` and `append([]string{cleanCodeID}, generated[i].DependsOn...)`.

2. **[blindspot] DepResolveFunc nil return behavior is unspecified**: Multiple resolvers explicitly return nil when they cannot resolve (`ResolveLastRunTest`, `ResolveUpstream`, `ResolveDocTasks`, `ResolveLastBusinessTask`). The proposal does not specify what `GenerateTestTasks` does when a resolver returns nil. Options: (a) skip the dependency (DependsOn omits the nil result), (b) skip the entire node, (c) error. This is a critical algorithmic gap that affects all scenarios. — Quote: "Returns concrete task IDs. Returns nil when the reference cannot be resolved."

3. **[blindspot] Registry ordering fragility is an unmanaged risk**: The proposal states "Order matters: nodes are processed in declaration order" and GenContext is "Populated progressively as nodes are processed in declaration order." If someone reorders registry entries (e.g., moving T-test-run before T-test-gen-scripts), `ResolveUpstream` would return wrong IDs. This ordering constraint is implicit — there is no validation that the ordering is correct, and the init-time validation does not check ordering invariants. The risk is not listed in the risk table. — Quote: "Order matters: nodes are processed in declaration order."

4. **[blindspot] Escape-hatch protocol cap has no enforcement mechanism**: The protocol states "escape hatch 总数上限 5 个" but no In Scope item implements enforcement, and no SC verifies it. The cap is an unenforceable declaration. — Quote: "escape hatch 总数上限 5 个，达上限时必须扩展 registry 表达力"

5. **[blindspot] `--no-validate` flag introduces a safety bypass without guardrails**: The risk mitigation for init() panic adds a `--no-validate` flag for emergency bypass. But if this flag is used, the registry runs without validation — exactly the scenario the validation was designed to prevent. There is no discussion of what happens when `--no-validate` is used with a broken registry (silent wrong task generation), or who is authorized to use this flag, or whether it should be logged/alerted. — Quote: "`--no-validate` flag for emergency bypass."

---

## Bias Detection Report

- Annotated regions: 0 attack points / 0 paragraphs = density N/A (no pre-revision markers detected)
- Unannotated regions: 10 attack points (5 blindspot + 5 rubric-derived) / ~50 paragraphs = density 0.20
- Ratio: N/A

---

## ATTACKS

1. [Solution Clarity]: Post-processing dependency chain description is factually inaccurate — "确保依赖链：`business-tasks → T-clean-code → T-review-doc（若存在）→ first-test-task`" describes a serial chain, but the code prepends both T-clean-code and T-review-doc to first-test-task.DependsOn in parallel; T-review-doc does NOT depend on T-clean-code — Fix the description to accurately reflect the parallel fan-in, or fix the code to implement the serial chain.

2. [Solution Clarity]: DepResolveFunc nil return behavior is unspecified — "Returns concrete task IDs. Returns nil when the reference cannot be resolved" — Specify what GenerateTestTasks does when a resolver returns nil: skip the dependency, skip the node, or error?

3. [Solution Clarity]: Core generation algorithm is under-specified — "filters the registry by Mode, ConfigGate, GenerateCondition, UISurfaceOnly, then expands per-surface nodes and resolves dependency references" — The filter-expand-resolve pipeline is described in 4 bullet points without pseudocode. How GenContext is progressively populated per-node during expansion is implicit. Add a step-by-step algorithm description or pseudocode.

4. [Industry Benchmarking]: Pipeline DAG comparison is shallow — "The registry's DepResolveFunc is closer to a pipeline DAG (Tekton Pipelines, GitHub Actions)" — The comparison asserts similarity but does not analyze how Tekton/GitHub Actions define their topology or what lessons apply. Deepen the comparison with specific pattern analysis.

5. [Requirements Completeness]: Empty business tasks scenario unaddressed — When `businessTasks` is empty, `CondHasTestableTasks` and `CondHasDocTasks` return false, only `CondAlways` nodes generate. But their resolvers (`ResolveLastRunTestOrBusiness`, `ResolveLastBusinessTask`) return nil. The proposal does not describe this scenario.

6. [Risk Assessment]: Registry ordering fragility is an unmanaged risk — "Order matters: nodes are processed in declaration order" — Reordering the registry silently breaks progressive GenContext population. Not listed in risk table. Add risk or add init-time ordering validation.

7. [Success Criteria]: Escape-hatch protocol cap has no SC enforcement — "escape hatch 总数上限 5 个" — No In Scope item enforces this cap and no SC verifies it. Either add an SC or acknowledge the cap is aspirational.

8. [Logical Consistency]: `--no-validate` safety bypass creates unanalyzed failure mode — "`--no-validate` flag for emergency bypass" — Using this flag with a broken registry produces silently wrong task generation, which is exactly what validation prevents. Document what happens in this mode and whether usage should be logged.

9. [Solution Creativity]: GenContext progressive population is an implicit ordering contract that increases maintenance burden — Future developers must understand that registry order affects behavior, but this constraint is enforced only by convention. Consider adding init-time ordering invariants (e.g., verify that nodes using ResolveUpstream appear after the node they expect as upstream).

10. [Feasibility]: Init-time validation of template-based DependsOn.Ref strings is non-trivial — "Validates: all DependsOn.Ref strings reference existing node IDs" — But DependsOn.Ref for T-test-gen-journeys references "T-eval-journey" (a concrete ID), while DependsOn.Ref for T-test-gen-scripts references "T-eval-contract" (concrete). No DependsOn.Ref currently uses template placeholders. But the validation logic must handle both cases, and this complexity is not acknowledged.

---

## Score Summary

```
SCORE: 794/1000
DIMENSIONS:
  1. Problem Definition: 95/110
  2. Solution Clarity: 102/120
  3. Industry Benchmarking: 78/120
  4. Requirements Completeness: 85/110
  5. Solution Creativity: 55/100
  6. Feasibility: 82/100
  7. Scope Definition: 72/80
  8. Risk Assessment: 78/90
  9. Success Criteria: 72/80
  10. Logical Consistency: 75/90
ATTACKS:
1. [Solution Clarity]: Serial chain description vs parallel fan-in code — "确保依赖链：business-tasks → T-clean-code → T-review-doc（若存在）→ first-test-task" — Fix description or code to be consistent.
2. [Solution Clarity]: DepResolveFunc nil return behavior unspecified — "Returns nil when the reference cannot be resolved" — Specify GenerateTestTasks behavior on nil resolver output.
3. [Solution Clarity]: Core generation algorithm under-specified — "filters the registry by... then expands per-surface nodes and resolves dependency references" — Add pseudocode for filter-expand-resolve pipeline.
4. [Industry Benchmarking]: Pipeline DAG comparison shallow — "closer to a pipeline DAG (Tekton Pipelines, GitHub Actions)" — Deepen with specific pattern analysis.
5. [Requirements Completeness]: Empty business tasks scenario — Resolvers return nil, CondHas* return false — Add scenario description and handling.
6. [Risk Assessment]: Registry ordering fragility — "Order matters: nodes are processed in declaration order" — Add to risk table or add init-time ordering validation.
7. [Success Criteria]: Escape-hatch cap unenforced — "escape hatch 总数上限 5 个" — Add SC or acknowledge as aspirational.
8. [Logical Consistency]: --no-validate bypass creates unanalyzed failure mode — "--no-validate flag for emergency bypass" — Document silent-wrong-output risk.
9. [Solution Creativity]: GenContext ordering is implicit convention — "Order matters" — Add init-time ordering invariant validation.
10. [Feasibility]: Template-based Ref validation complexity unacknowledged — "Validates: all DependsOn.Ref strings reference existing node IDs" — Acknowledge complexity of validating template placeholders.
```

**Attack density**: 10 attacks across 7 of 10 dimensions. Unannotated regions: 10/10 (no pre-revision markers detected).
