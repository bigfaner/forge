# Pipeline Topology Registry -- Freeform Narrative Review

**Reviewer**: Go Pipeline Integration & Type System Engineer
**Date**: 2026-05-29
**Document**: `docs/proposals/pipeline-topology-registry/proposal.md`

---

## Background Assessment

The proposal addresses a real and well-documented fragility in the Forge CLI's auto-task generation pipeline. The core problem -- that adding or removing a task type requires coordinated changes across 6+ code locations -- is validated by concrete bug evidence. I verified the current codebase and confirmed that `findTaskIndexOrPanic` appears 13 times in `autogen.go` (with 4 more for the prefix variant), the dependency wiring is indeed scattered across `resolveBreakdownDeps`, `resolveQuickDeps`, `wireRunTestChain`, `wireQuickRunTestChain`, and `ResolveDriftFallbackDep`, and the type inference switch in `infer.go` contains 15+ cases with mixed matching strategies (exact match, prefix match, suffix match, surface-key lookup).

The proposed solution -- a declarative `PipelineRegistry` that serves as the single source of truth for all auto-generated tasks -- is architecturally sound. The registry pattern eliminates the root cause: instead of N scattered definitions that must be manually kept in sync, there is one ordered list of `PipelineNode` entries from which all consuming logic derives. The approach of using `ConfigGateFunc`, `GenerateCondFunc`, and `DepResolveFunc` function types to capture mode-dependent and context-dependent behavior is idiomatic Go and preserves the flexibility that the current procedural code provides.

The proposal's scope is appropriately bounded. It correctly excludes runtime task creation (fix-*, disc-*), the task state machine, prompt templates, and stage gate generation. The user-facing behavior is explicitly stated as unchanged, making this a pure refactoring with a clear behavioral parity contract.

---

## Key Risk Identification

问题：
The proposal's registry entries for downstream tasks (validation, consolidation, drift, clean-code) all lack an explicit `GenerateCondition` field. The document states: "GenerateCondition: nil defaults to CondHasTestableTasks." But in the current codebase, `GetBreakdownTestTasks` and `GetQuickTestTasks` generate validation, consolidation, and clean-code tasks **unconditionally** once their config gate passes -- they do not check whether testable business tasks exist. The only condition applied is `auto.Validation.Full`, `auto.ConsolidateSpecs.Full`, etc. A doc-only feature with `auto.Validation.Full = true` would not generate `T-validate-code` in the current code because `needsTestPipeline` returns false and the entire step 7.5 block is skipped -- but under the registry's default `CondHasTestableTasks`, these nodes would never pass the gate either, so the behavior happens to match. However, this is a coincidence, not a deliberate design. If a future change makes `CondHasTestableTasks` the default for nodes that should actually use `CondAlways` (or no condition at all), the behavior will silently diverge. The proposal does not address this semantic mismatch and does not explain why `CondHasTestableTasks` is the correct default for nodes whose config gates already provide the only necessary gating.

风险：
The `T-clean-code` entry in the registry has `DependsOn: []DepRef{}` with the comment "resolved by caller: depends on last business task." This represents a fundamental escape hatch from the registry's promise of being the single source of truth. In the current codebase, `ResolveFirstTestDep` handles T-clean-code by: (1) setting its dependency to the highest gate or last business task, and (2) inserting T-clean-code as a dependency of the first test task. This bidirectional wiring logic -- where T-clean-code both depends on business tasks and is depended upon by test tasks -- is nowhere in the registry. The `DependsOn: []DepRef{}` entry silently delegates this to unspecified "caller" logic, which means the registry cannot fully validate its own dependency graph. Any init-time referential integrity check would mark this empty dependency set as valid when it is actually incomplete. The proposal should either: (a) encode this wiring pattern in the registry using a new resolver function (e.g., `ResolveLastBusinessTask`), or (b) explicitly document that `ResolveFirstTestDep` remains outside the registry and explain why it cannot be captured.

风险：
The registry completely omits the **single-surface degenerate case** for run-test tasks. In the current codebase, when `isSingleSurface(surfaces)` is true, the run-test task gets `ID: "T-test-run"` (no suffix) and `Key: "run-test"`. But the registry entry uses `ID: "T-test-run-{surface-key}"` with `Expansion: "per-surface-key"`. The proposal never explains how expansion handles the single-surface case. Does the `{surface-key}` placeholder resolve to an empty string? Does the expansion skip and emit `T-test-run` as a plain ID? This matters because `InferType` in the current code explicitly matches `id == "T-test-run"` as a separate case from `testRunSurfaceKeyMatch(id, surfaces)`. If the registry's ID pattern matching does not handle this degenerate case, the behavioral parity guarantee fails for single-surface projects -- which is the most common project configuration.

问题：
The runtime task coordination section lists only `fix-` and `disc-` prefixes: "Runtime tasks (fix-*, disc-*) are generated by quality-gate and run-tasks dispatcher." But `InferType` in the current code also handles `doc-fix-` prefix (`case strings.HasPrefix(id, "doc-fix-"): return TypeDocFix`). This prefix is not mentioned in the runtime task coordination contract. The omission suggests either: (a) `doc-fix-*` tasks were overlooked and the new `InferType` fallback will not match them, causing `TypeDocFix` inference to break; or (b) `doc-fix-*` tasks are handled differently and the proposal should explain why they are excluded from the coordination contract.

问题：
The `T-review-doc` injection logic in `build.go` -- implemented by `resolveTestDepsAndInjectReviewDoc` -- prepends `"T-review-doc"` to the first test task's dependencies when both doc and coding tasks exist. This is the "mixed" scenario described in the three-scenario table. However, the registry's `T-review-doc` entry has `DependsOn: []DepRef{{Resolve: ResolveDocTasks}}` and the test pipeline nodes have their own `DependsOn` chains, but **nowhere in the registry is the reverse injection captured** -- that T-review-doc must be prepended to the first test task's dependency list. This is the same class of problem as T-clean-code: the registry describes each node's outgoing dependencies but not the dynamic reverse dependencies that `build.go` injects. The proposal lists "T-review-doc injection into test pipeline missed" as a Low-likelihood risk but does not provide a concrete mechanism for how the registry prevents this.

风险：
The `PipelineNode` struct is missing fields that `AutoGenTaskDef` currently uses and that the `.md` file generation depends on: `Key` (the map key in index.json, e.g., "gen-scripts", "gen-scripts-api"), `FileName` (derived from key), `Breaking` (propagated to `Task.Breaking`), and `StrategyContent` (resolved by caller from convention files). The `Key` field is particularly critical because it determines the filename of the generated `.md` file and the lookup key in `index.json`. In the current code, different tasks have different key derivation logic (e.g., `"gen-test-scripts-" + typ` for per-type tasks, `"run-test-" + key` for per-surface tasks, `"quick-drift-detection"` vs `"consolidate-specs"` for the same type in different modes). The proposal does not explain how `Key` is derived from the registry entry, and without this, the `GenerateTestTaskMD` and `TaskFromFile` functions cannot be driven by the registry alone.

问题：
The proposal defines `CondAlways` as a predefined condition function but never uses it in any registry entry. Every non-T-review-doc entry relies on the `CondHasTestableTasks` default. This raises the question: when would `CondAlways` be used? If validation/clean-code tasks should be generated regardless of business task composition (only gated by config), they should use `CondAlways`. If they should require testable tasks, they should use `CondHasTestableTasks`. The proposal is silent on which is correct, and the `CondAlways` function appears to be dead-on-arrival code that obscures the design intent.

风险：
The proposal states "all existing tests pass" as a success criterion and lists 5 test files totaling ~7,300 lines, including `claim_test.go` in `internal/cmd/task/`. But `claim.go` references `task.CategoryForType`, `task.TypeCodingFix`, and `task.GetTaskPhase` -- none of which are being modified. The inclusion of `claim_test.go` in scope suggests these tests call functions like `isTestTaskID` or `IsAutoGenTaskID` indirectly, but my analysis shows they do not. The scope may be overstated, or there may be indirect dependencies through `BuildIndex` integration tests. The proposal should clarify the actual impact on `claim_test.go` rather than blanket-including it.

风险：
The `InferType` refactoring proposes to "iterate the registry and match ID patterns (with wildcard support for `{surface-key}`/`{surface-type}` placeholders)." But the current `InferType` has specialized matching logic that cannot be expressed by simple wildcard patterns: `testRunSurfaceKeyMatch` requires looking up the suffix in a `surfaces` map passed at call time. The registry-based `InferType` would need access to this same runtime state. The proposal does not explain how the registry iteration receives the `surfaces` map or how `GenContext` is populated for inference (which currently is a stateless operation). This is a design gap that will surface during implementation.

问题：
The proposal mentions init-time validation: "the registry is checked for referential integrity (all DependsOn references resolve to existing nodes or are valid special references)." But `DepRef.Resolve` functions are dynamic -- they produce different IDs depending on runtime state (e.g., `ResolveUpstream` returns IDs of previously generated nodes). Static init-time validation cannot check dynamic references. The proposal conflates static reference validation (checking that `Ref: "T-test-gen-journeys"` matches a registry entry's ID) with dynamic resolver validation (which cannot be done without runtime state). This should be split into two distinct checks with different guarantees.

问题：
The existing `ValidateAutogenTemplates` function (called at startup in `cmd/forge/run.go`) iterates `ValidTypes` to verify template file existence and validity. The proposal does not mention whether this validation is updated to use the registry, whether it remains separate, or how it interacts with the new init-time registry validation. Two separate validation passes at startup would be wasteful and confusing; they should be unified or their relationship should be explicitly documented.

风险：
The proposal's scope includes refactoring `build.go` steps 7/7.5/7.6 but does not address `resolveTestDepsAndInjectReviewDoc`, `ResolveFirstTestDep`, `findHighestGateOrSummary`, or `findMaxBusinessTaskID`. These functions in `build.go` implement the critical "insert T-clean-code between business and test tasks" and "inject T-review-doc as first test dependency" logic. If they remain as-is while the rest of `autogen.go` is refactored to use the registry, the codebase will have a hybrid architecture where some dependency resolution is registry-driven and some is still procedural, defeating the stated goal of a "single source of truth."

---

## Improvement Suggestions

建议：
Define an explicit `KeyDerivation` strategy for `PipelineNode`. The current `AutoGenTaskDef.Key` field has mode-dependent and expansion-dependent derivation logic (e.g., `"run-test"` vs `"run-test-" + key`, `"quick-drift-detection"` vs `"consolidate-specs"` for the same type in different modes). Add a `Key` field or `KeyFunc` to `PipelineNode` that produces the index.json map key, making the registry truly self-contained for `.md` file generation.

建议：
Replace `DependsOn: []DepRef{}` on the T-clean-code entry with a dedicated resolver such as `ResolveLastBusinessTask`, and add a new field (e.g., `PrependToFirstTestDep bool`) or a `ReverseDep` mechanism that captures the bidirectional wiring. This eliminates the "resolved by caller" escape hatch and makes the registry's dependency graph fully self-describing. The init-time validation can then check that every node's dependency set is non-empty (for nodes that are not the first in the pipeline).

建议：
Add an `ExpansionRules` sub-section to the design that explicitly covers: (a) how single-surface projects emit `T-test-run` without a suffix; (b) how expanded IDs are validated for uniqueness; (c) how the `Key` field is generated for expanded nodes. This is currently the largest gap between the registry's declarative model and the runtime behavior that the existing code handles.

建议：
Include `doc-fix-*` in the runtime task coordination section alongside `fix-*` and `disc-*`. Update the proposed `InferType` fallback logic to explicitly list all three runtime prefixes: `fix-`, `disc-`, and `doc-fix-`. This prevents `TypeDocFix` inference from being silently broken.

建议：
Unify the init-time validation story. Instead of proposing a separate "init-time referential integrity check," extend the existing `ValidateAutogenTemplates` to also validate the registry (or rename it to `ValidatePipeline` and have it do both). This avoids two separate validation passes and ensures that template-file-to-registry-node correspondence is checked atomically. Document in the proposal how `cmd/forge/run.go` should be updated.

建议：
Add `GenerateCondition: CondAlways` explicitly to validation, consolidation, drift, and clean-code entries if they should be generated based solely on config gating. Remove the "nil defaults to CondHasTestableTasks" default rule entirely -- requiring every entry to specify its condition explicitly makes the registry self-documenting and prevents accidental behavioral divergence from silent defaults.

建议：
Address the `surfaces` map dependency for `InferType` refactoring. Either: (a) make `InferType` accept a `GenContext` instead of a bare `surfaces` map, threading the context through from the registry; or (b) keep `InferType`'s current two-pass approach (registry iteration for static IDs, then prefix/surface-key fallback) and document the boundary clearly. The current proposal's "iterate the registry and match ID patterns" is underspecified for the surface-key lookup case.

建议：
Document the relationship between the registry and `ResolveFirstTestDep`/`resolveTestDepsAndInjectReviewDoc` explicitly. If these functions remain outside the registry, explain why (e.g., they depend on `existingTasks` which is not available at registry definition time) and add a clear contract for how they interact with registry-generated tasks. If they should be absorbed, describe the mechanism (e.g., a post-generation hook in the `GenContext`).
