# Pipeline Topology Registry — Eval Report (Iteration 1)

**Reviewer**: CTO adversary (blind to author identity)
**Date**: 2026-05-29
**Rubric**: proposal.md (1000-point scale)
**Pre-revision annotations**: None detected in document. All regions treated as unannotated.

---

## Phase 1: Reasoning Audit

### Problem -> Solution Chain

**Problem**: Task IDs and dependency relationships are hardcoded as string literals across 6+ code locations. Adding/removing a task type causes cascading failures.

**Solution**: Pipeline Topology Registry — a single declarative data structure that all consuming code derives from.

**Chain verdict**: The solution directly addresses the stated problem. The registry centralizes scattered definitions. However, the chain has two gaps:

1. **ResolveFirstTestDep's stage-gate logic is silently dropped.** The proposal states `ResolveFirstTestDep` is "replaced by injectReviewDocDep's first-task lookup" (Functions Relationship table), but the current `ResolveFirstTestDep` does three things: (a) resolve first-test-task dep to highest gate or last business task in breakdown mode, (b) inject T-clean-code as intermediate dependency between business tasks and first test task, (c) inject T-review-doc. The proposal's replacement (`injectReviewDocDep`) only covers (c). Items (a) and (b) are not represented in the registry or any other mechanism.

2. **ResolveDriftFallbackDep's fallback is not represented.** When test pipeline is disabled, drift/consolidate tasks currently fall back to depending on the last business task. The registry's `DependsOn: []DepRef{{Resolve: ResolveLastRunTest}}` returns nil when no run-test tasks exist. The fallback path is unaddressed.

### Solution -> Evidence Chain

Evidence: 5 bug examples, 6 touch points documented, ~40% failure rate claimed. Code verification confirms:
- 13 `findTaskIndexOrPanic` calls in autogen.go (confirmed)
- 6-prefix `isTestTaskID` (confirmed)
- 4 files contain scattered references (autogen.go, build.go, infer.go, extract.go) — matches claim

### Evidence -> SC Chain

SC1 claims "exactly ONE file (pipeline.go)" but immediately contradicts with "plus optionally the type constant in types.go". If a new task type requires both a registry entry AND a type constant, that is TWO files, not ONE. The word "optionally" is misleading since every new auto-generated task type would need a new type constant.

### Self-Contradiction Check

- SC1 ("ONE file") vs. body ("plus optionally the type constant in types.go") — minor contradiction.
- Functions Relationship table says `ResolveFirstTestDep` is "deleted, logic replaced by injectReviewDocDep" — but ResolveFirstTestDep does far more than review-doc injection. The table entry is factually incorrect.
- Registry declares `T-clean-code` depends on `ResolveLastBusinessTask`, but current code has first-test-task depending on T-clean-code. This reverse dependency is not captured anywhere in the proposal.

---

## Phase 2: Rubric Scoring

### 1. Problem Definition: 82/110

| Criterion | Score | Rationale |
|-----------|-------|-----------|
| Problem stated clearly | 35/40 | Core problem is unambiguous: scattered string literals cause cascading failures. Clear enough for any reader to understand. |
| Evidence provided | 32/40 | 5 concrete bugs with root causes. Code-smell audit of 6 touch points. However, the "~40% chance" statistic is stated without methodology — "based on the last 5 additions (3 caused issues)" is a sample of 5, too small for a percentage claim. No user feedback or telemetry data cited. |
| Urgency justified | 15/30 | "Cost of inaction" section mentions developer velocity and manual verification burden, but no quantified impact (hours lost per incident, time to add a task type, number of affected users). The urgency is implied, not demonstrated. |

### 2. Solution Clarity: 85/120

| Criterion | Score | Rationale |
|-----------|-------|-----------|
| Approach is concrete | 35/40 | Full Go struct definitions, registry entries, resolver functions. A reader can explain back what will be built. Strong. |
| User-facing behavior described | 40/45 | Explicitly stated: "No change. The pipeline generates the same tasks with the same dependencies. The refactoring is purely internal architecture." Clear and honest. |
| Technical direction clear | 10/35 | Registry data structure is detailed, but the generation algorithm is vague. How does `GenerateTestTasks` produce output from the registry? The description says "filters the registry by Mode, ConfigGate, GenerateCondition, UISurfaceOnly, then expands per-surface nodes and resolves dependency references" — but the expansion algorithm, dependency resolution ordering, and how `GenContext` is progressively populated are described in prose without pseudocode or flow. The critical `ResolveFirstTestDep` replacement logic is missing entirely (see Phase 1 gap 1). |

### 3. Industry Benchmarking: 45/120

| Criterion | Score | Rationale |
|-----------|-------|-----------|
| Industry solutions referenced | 10/40 | No industry solutions, open-source projects, or published patterns cited. No mention of DAG libraries, pipeline frameworks (Airflow, Tekton), or even Go-specific approaches (wire, dig). The proposal is entirely self-invented. |
| At least 3 meaningful alternatives | 18/30 | Three alternatives presented: Table-Driven, Do Nothing, Code Generation. "Do Nothing" is explicitly required by the rubric. However, none reference industry-validated solutions. Alternative 1 (Table-Driven) is essentially the same approach with less ambition — a straw man. Alternative 3 (Code Generation) is dismissed with "adds a build step" without substantiation. |
| Honest trade-off comparison | 10/25 | Trade-offs are cherry-picked. Alternative 1 is rejected because it "doesn't solve the root cause" — but it could solve 80% of the problem with 20% of the effort. Alternative 3's "build step" objection is weak when `go generate` is standard practice in Go. No comparison table with weighted criteria. |
| Chosen approach justified against benchmarks | 7/25 | No benchmark to justify against. The justification is internal logic only (centralization = good). No analysis of whether other projects with similar pipeline complexity use registries, tables, or code generation. |

### 4. Requirements Completeness: 62/110

| Criterion | Score | Rationale |
|-----------|-------|-----------|
| Scenario coverage | 22/40 | Three scenarios for T-review-doc are documented. Missing scenarios: (1) What happens when ALL ConfigGates are disabled? (2) Single-surface vs multi-surface behavior for each node — only T-test-run's single-surface degeneration is described. (3) Empty business tasks list. (4) Business tasks with only non-testable types when test pipeline is enabled. (5) The drift fallback when test pipeline is disabled — this is an existing scenario not covered by the registry. |
| Non-functional requirements | 25/40 | Performance addressed ("~12 entries; iteration cost is negligible"). No mention of: memory footprint of registry at scale, init-time validation latency, backward compatibility with existing index.json files, migration path for in-flight features. Compatibility is critical — the proposal changes the `GenerateTestTasks` signature (adds `businessTasks` parameter) without noting this as a breaking API change. |
| Constraints & dependencies | 15/30 | External dependency on `forgeconfig.AutoConfig` is implicit but not called out. The `GenContext` struct has a progressive population contract that depends on declaration order — this ordering constraint is mentioned but not validated. No mention of Go version requirements or dependency on existing type system. |

### 5. Solution Creativity: 45/100

| Criterion | Score | Rationale |
|-----------|-------|-----------|
| Novelty over industry baseline | 20/40 | The registry pattern is standard — not novel. The `DepResolveFunc` with progressive `GenContext` is a reasonable design but not innovative. The proposal does not articulate what differentiates it from standard DAG/pipeline patterns. |
| Cross-domain inspiration | 10/35 | No cross-domain references. No mention of compiler IR pipelines, database query planners, or CI/CD DAGs that face similar problems. |
| Simplicity of insight | 15/25 | The core insight ("centralize scattered definitions into a registry") is straightforward and clean. The `ConfigGateFunc`/`GenerateCondFunc`/`DepResolveFunc` trichotomy is well-structured. However, the `injectReviewDocDep` post-generation step undermines the declarative purity — it's an escape hatch that admits the registry cannot express cross-pipeline dependencies. |

### 6. Feasibility: 65/100

| Criterion | Score | Rationale |
|-----------|-------|-----------|
| Technical feasibility | 30/40 | Go structs, closures, and iteration are standard. No showstopper dependencies. However, the `ResolveFirstTestDep` replacement gap (Phase 1 finding) means the implementation will discover missing logic mid-refactor. The `init()` validation for `DependsOn.Ref` strings is non-trivial — it requires resolving template placeholders to verify cross-references, which the proposal glosses over. |
| Resource & timeline feasibility | 20/30 | No timeline estimate given. The scope (9 in-scope items touching 4+ files with test updates) is substantial. "Update all existing tests" is a significant effort item with no breakdown. |
| Dependency readiness | 15/30 | `forgeconfig.AutoConfig` is an existing dependency — available. But the proposal does not address whether `CategoryForType` and `IsTestableType` (used by condition functions) are stable APIs or subject to change. The `GenContext` progressive population contract requires that all registry consumers understand the ordering semantics. |

### 7. Scope Definition: 58/80

| Criterion | Score | Rationale |
|-----------|-------|-----------|
| In-scope items are concrete | 22/30 | 9 numbered items, each targeting specific files and functions. Strong. However, item 7 (two-phase validation) is underspecified — "panics on failure" for init-time validation in a CLI tool that may be used in CI pipelines is a questionable design choice not discussed. |
| Out-of-scope explicitly listed | 20/25 | 6 items explicitly listed as out-of-scope. Clear. But `ResolveFirstTestDep`'s stage-gate resolution logic (which uses `findHighestGateOrSummary`) is listed as "Out of Scope" under `findHighestGateOrSummary`, yet this function is called FROM `ResolveFirstTestDep` which is listed as "deleted". If you delete the caller but keep the callee, the callee is dead code. If you keep the callee's functionality somewhere else, that's not documented. |
| Scope is bounded | 16/25 | No timeline, no sprint allocation, no effort estimate. "Update all existing tests" is an unbounded item. The scope could range from 2 days to 2 weeks depending on test complexity. |

### 8. Risk Assessment: 55/90

| Criterion | Score | Rationale |
|-----------|-------|-----------|
| Risks identified | 18/30 | 5 risks listed. Missing risks: (1) ResolveFirstTestDep logic gap — stage-gate resolution and clean-code injection not covered by registry. (2) Drift fallback gap when test pipeline disabled. (3) `GenerateTestTasks` signature change breaking callers. (4) Init-time panic in CI environments. (5) Progressive `GenContext` population bugs when registry order is changed. |
| Likelihood + impact rated | 20/30 | Ratings are reasonable. "Regression in task generation output" at Medium/High is honest. However, "Expanding per-surface nodes creates ID conflicts" at Low likelihood seems optimistic — this is exactly the kind of edge case that bit the current implementation. |
| Mitigations are actionable | 17/30 | "Snapshot tests" is actionable. "Port all existing gotcha lessons as test cases" is vague — which lessons? Where are they documented? "Benchmark if concerned" is not a mitigation, it's a deferral. |

### 9. Success Criteria: 50/80

| Criterion | Score | Rationale |
|-----------|-------|-----------|
| Criteria are measurable and testable | 18/30 | SC2 (zero panics), SC8 (80% coverage), SC9 (all tests pass) are measurable. SC1 ("exactly ONE file") is self-contradictory (also mentions types.go). SC7 ("zero magic values") is measurable but extremely granular — every string literal must be a constant, which is a style rule, not a success criterion for a refactoring proposal. SC4 ("identical output") is measurable via snapshot testing. |
| Coverage is complete | 12/25 | Critical gaps: No SC for "ResolveFirstTestDep's stage-gate logic is preserved". No SC for "drift fallback when test pipeline is disabled". No SC for "T-clean-code is injected as intermediate dependency between business tasks and first test task". No SC for backward compatibility of GenerateTestTasks signature change. |
| SC internal consistency | 20/25 | SC1 ("ONE file") contradicts itself. SC4 (behavioral parity) conflicts with the missing ResolveFirstTestDep coverage — if the implementation drops the stage-gate resolution, behavioral parity cannot be achieved. SC7 (zero magic values) and SC4 (behavioral parity) could conflict if existing behavior relies on string matching that breaks with constant-based comparisons. Ambiguous but not definitively contradictory. |

### 10. Logical Consistency: 52/90

| Criterion | Score | Rationale |
|-----------|-------|-----------|
| Solution addresses the stated problem | 25/35 | The registry centralizes scattered definitions — directly addresses the stated problem. However, it does NOT fully replace all 6 touch points because: (a) ResolveFirstTestDep's stage-gate logic is lost, (b) clean-code injection is lost, (c) drift fallback is lost. These are gaps in the solution, not failures of the concept. |
| Scope <-> Solution <-> SC aligned | 12/30 | Misalignment: Scope says "delete ResolveFirstTestDep" but Solution doesn't provide replacement for its stage-gate logic. Scope says "delete findMaxBusinessTaskID" (replaced by ResolveLastBusinessTask resolver) but the resolver doesn't handle the clean-code intermediate injection pattern. SC4 (behavioral parity) cannot be met if these deletions occur without replacement. |
| Requirements <-> Solution coherent | 15/25 | The GenerateCondition/ConfigGate/Mode trichotomy covers most requirements. But the requirement "drift tasks depend on last business task when test pipeline disabled" (from the Problem section's bug table) has no corresponding solution mechanism. The requirement "clean-code runs after all business tasks AND before first test task" (from current code) is only half-represented (clean-code depends on last-biz, but first-test-task doesn't depend on clean-code). |

---

## Phase 3: Blindspot Hunt

### What the rubric missed:

1. **Escape hatch proliferation risk**: The `injectReviewDocDep` post-generation step is an escape hatch that admits the declarative registry cannot express cross-pipeline dependencies. This sets a precedent — every future cross-concern dependency will need its own escape hatch, gradually recreating the procedural mess the registry was designed to eliminate. The proposal should either (a) design the registry to handle cross-pipeline deps, or (b) define a bounded escape-hatch protocol with clear rules.

2. **Testing strategy gap**: The proposal mentions "snapshot tests" and "port all existing gotcha lessons" but does not describe a testing strategy for the registry itself. How do you test that `GenContext` progressive population is correct for all mode/config/surface combinations? A combinatorial test matrix is needed.

3. **No rollback plan**: If the refactoring introduces a regression in production CI, there is no rollback strategy described. Since this is a purely internal refactoring, the rollback is trivially "git revert", but this should be stated explicitly.

4. **Versioning and migration**: The proposal changes the `GenerateTestTasks` function signature (adds `businessTasks`). Any caller must be updated. The proposal does not enumerate all callers or describe the migration strategy.

5. **`init()` panic in CI**: Phase 1 validation runs at CLI startup via `init()` and panics on failure. If a registry bug is deployed, every CI pipeline using `forge` will panic immediately. This is by design (fail fast) but should be explicitly justified as a trade-off.

---

## ATTACKS

1. **[Logical Consistency]** ResolveFirstTestDep's stage-gate and clean-code injection logic is silently dropped — "被 injectReviewDocDep 的首任务查找替代 | 删除" (Functions Relationship table) — The proposal must either (a) represent stage-gate resolution and clean-code injection in the registry, or (b) add a post-generation step for these alongside `injectReviewDocDep`, or (c) explicitly document that this logic is being removed and explain why it is safe.

2. **[Logical Consistency]** Drift/consolidate fallback dependency when test pipeline is disabled is unrepresented — Registry defines `DependsOn: []DepRef{{Resolve: ResolveLastRunTest}}` for drift/consolidate nodes, but `ResolveLastRunTest` returns nil when test pipeline is disabled — The proposal must add a fallback resolver or post-generation step that sets drift/consolidate deps to `lastBusinessTask` when `ResolveLastRunTest` returns nil.

3. **[Success Criteria]** SC1 self-contradiction: "exactly ONE file (pipeline.go)" immediately followed by "plus optionally the type constant in types.go" — This is TWO files. Rewrite SC1 to say "at most TWO files: pipeline.go for the registry entry and types.go for the type constant, with no other file changes required."

4. **[Industry Benchmarking]** Zero industry references — The Alternatives section names 3 alternatives, none of which cite external projects, patterns, or published approaches. No mention of DAG libraries, pipeline orchestration patterns, or even how other Go CLI tools handle similar pipeline definition problems — Add at least one industry-validated pattern (e.g., dependency injection containers, pipeline DAGs, rule engines) and compare the registry approach against it.

5. **[Industry Benchmarking]** Alternative 1 (Table-Driven) is a straw man — "Doesn't solve the root cause (scattered definitions). Just makes the strings slightly more discoverable." — This dismisses a valid incremental approach that could solve 80% of the problem with minimal risk. The rejection should acknowledge that table-driven is a legitimate intermediate step and explain why the full registry is worth the additional complexity and risk.

6. **[Requirements Completeness]** Missing scenario: test pipeline disabled, all ConfigGates off — The registry nodes for drift/consolidate/validation all use `ResolveLastRunTest` which returns nil when no test pipeline exists. The current `ResolveDriftFallbackDep` handles this, but the proposal deletes it without replacement — Add this scenario to the requirements or document how it is handled.

7. **[Requirements Completeness]** GenerateTestTasks signature change is undocumented breaking change — "GenerateTestTasks(mode, surfaces, executionOrder, auto, businessTasks)" adds a `businessTasks` parameter — Current signature is `GenerateTestTasks(mode string, surfaces map[string]string, executionOrder []string, auto forgeconfig.AutoConfig)`. The proposal must enumerate all callers and describe the migration.

8. **[Solution Clarity]** T-clean-code reverse dependency is missing — Current code: `ResolveFirstTestDep` sets first-test-task to depend on T-clean-code. Registry: T-clean-code depends on ResolveLastBusinessTask, but no node declares dependency on T-clean-code — The proposal must show how first-test-task comes to depend on T-clean-code, either via registry declaration or post-generation step.

9. **[Risk Assessment]** Missing risk: init() panic in CI environments — "Phase 1 (static, init-time): Run at CLI startup via init(). Panics on failure." — A malformed registry entry would crash every CI pipeline. This risk should be listed with a mitigation (e.g., registry validation runs in a separate CLI subcommand during development, not at every invocation).

10. **[Risk Assessment]** Missing risk: ResolveFirstTestDep logic regression — The function performs stage-gate resolution, clean-code injection, and review-doc injection. Deleting it without full replacement is the highest-risk item in the proposal. Not listed in the risk table — Add this risk with High likelihood and High impact.

11. **[Problem Definition]** "~40% chance" statistic is unsupported — "Each new task type has a ~40% chance of introducing at least one pipeline bug, based on the last 5 additions (3 caused issues)" — A sample size of 5 is statistically meaningless for a percentage claim. State the raw data (3 of 5 recent additions caused issues) without the percentage.

12. **[Solution Creativity]** injectReviewDocDef undermines declarative purity — The `injectReviewDocDep` post-generation step exists because the registry cannot express "when T-review-doc exists AND test pipeline is active, prepend T-review-doc as dependency of first test task" — This is a cross-concern dependency. The proposal should acknowledge this limitation and define a protocol for future escape hatches to prevent gradual re-proceduralization.

13. **[Feasibility]** No timeline or effort estimate — The scope includes 9 items touching 4+ files with full test suite updates. No estimate of development time, review cycles, or deployment strategy — Add a rough estimate (e.g., "2-3 development days") to bound the scope.

14. **[Scope Definition]** findHighestGateOrSummary is listed as "保留不变 (Out of Scope)" but its only caller (ResolveFirstTestDep) is "删除" — If you delete the caller, the function becomes dead code. The proposal must either (a) find a new home for stage-gate resolution logic, or (b) acknowledge that this function will be deleted as dead code after the refactoring.

15. **[Success Criteria]** SC7 (zero magic values) is a style rule, not a success criterion — "All string literals in pipeline.go and refactored code must use typed enum constants" — This is a code quality rule that should be enforced by linters, not a proposal success criterion. It inflates the SC count without testing the actual refactoring goal.

---

## Score Summary

```
SCORE: 599/1000
DIMENSIONS:
  1. Problem Definition: 82/110
  2. Solution Clarity: 85/120
  3. Industry Benchmarking: 45/120
  4. Requirements Completeness: 62/110
  5. Solution Creativity: 45/100
  6. Feasibility: 65/100
  7. Scope Definition: 58/80
  8. Risk Assessment: 55/90
  9. Success Criteria: 50/80
  10. Logical Consistency: 52/90
ATTACKS:
1. [Logical Consistency]: ResolveFirstTestDep's stage-gate and clean-code injection logic silently dropped — "被 injectReviewDocDep 的首任务查找替代 | 删除" — Must represent this logic in registry or post-generation step; current replacement only covers review-doc injection.
2. [Logical Consistency]: Drift/consolidate fallback when test pipeline disabled unrepresented — `DependsOn: []DepRef{{Resolve: ResolveLastRunTest}}` returns nil when no run-test tasks exist — Must add fallback resolver or post-generation step for this scenario.
3. [Success Criteria]: SC1 self-contradiction — "exactly ONE file (pipeline.go) — the registry definition plus optionally the type constant in types.go" — "exactly ONE" and "plus optionally" are contradictory; rewrite to say "at most TWO files".
4. [Industry Benchmarking]: Zero industry references — Alternatives section cites no external projects, patterns, or published approaches — Cite at least one industry-validated pipeline/DAG pattern and compare against it.
5. [Industry Benchmarking]: Alternative 1 (Table-Driven) is a straw man — "Doesn't solve the root cause (scattered definitions). Just makes the strings slightly more discoverable." — Acknowledge it solves most of the problem at lower risk; justify why full registry is worth the additional complexity.
6. [Requirements Completeness]: Missing scenario: test pipeline disabled — "ResolveDriftFallbackDep" is deleted but its scenario (drift tasks with no test pipeline) is not covered by any registry mechanism — Add scenario and handling.
7. [Requirements Completeness]: GenerateTestTasks signature change undocumented — Proposal adds `businessTasks` parameter to `GenerateTestTasks(mode, surfaces, executionOrder, auto, businessTasks)` — Enumerate all callers and describe migration.
8. [Solution Clarity]: T-clean-code reverse dependency missing — Registry has T-clean-code depending on ResolveLastBusinessTask, but no node depends on T-clean-code — Must show how first-test-task comes to depend on T-clean-code (currently done by ResolveFirstTestDep).
9. [Risk Assessment]: Missing risk: init() panic in CI — "Phase 1 (static, init-time): Run at CLI startup via init(). Panics on failure." — A registry bug crashes all CI pipelines; add to risk table with mitigation.
10. [Risk Assessment]: Missing risk: ResolveFirstTestDep logic regression — Highest-risk deletion in the proposal; not listed in risk table — Add with High likelihood, High impact.
11. [Problem Definition]: "~40% chance" unsupported — "based on the last 5 additions (3 caused issues)" — Sample of 5 is too small for a percentage claim; state raw numbers.
12. [Solution Creativity]: injectReviewDocDep undermines declarative purity — Registry cannot express cross-pipeline deps; escape hatch sets precedent for gradual re-proceduralization — Acknowledge limitation and define escape-hatch protocol.
13. [Feasibility]: No timeline or effort estimate — 9 scope items, 4+ files, full test suite update — Add rough estimate to bound scope.
14. [Scope Definition]: findHighestGateOrSummary listed as "保留不变" but only caller (ResolveFirstTestDep) is "删除" — Function becomes dead code; document its fate explicitly.
15. [Success Criteria]: SC7 (zero magic values) is a code style rule, not a refactoring success criterion — Move to code quality enforcement mechanism; replace with a criterion that tests the refactoring's actual goal.
```

**Attack density**: 15 attacks across 8 of 10 dimensions. Unannotated regions: 15/15 (no pre-revision markers detected in document).
