# Freeform Expert Review: surface-test-ordering

**Expert**: Test Pipeline Architect
**Date**: 2026-05-25
**Proposal**: docs/proposals/surface-test-ordering/proposal.md

## Executive Summary

The proposal identifies a real architectural gap (cross-surface test ordering and gen-journeys semantic inconsistency) and proposes two targeted changes. The gen-journeys merge is well-motivated and architecturally sound. However, the run-test split into per-surface-key serial tasks introduces significant complications that the proposal underestimates: it breaks `InferType`'s intentional exact-match contract for `T-test-run`, requires pervasive changes across `isTestTaskID`/`isAutoGenForDep`/`resolveTestDepsAndInjectReviewDoc`, and collides with the existing `typeSuffixedID` naming convention (which currently only applies to gen-journeys and gen-scripts). The proposal also conflates surface-key (e.g., "backend") with surface-type (e.g., "api") in its naming scheme, which would create ambiguity in the type inference system.

## Architectural Analysis

### 1. Two Naming Schemes Collide: surface-key vs surface-type

The current autogen framework uses **surface-type** suffixes consistently for per-type tasks:
- `T-test-gen-journeys-api`, `T-test-gen-journeys-cli` (deduplicated surface types)
- `T-test-gen-scripts-api`, `T-test-gen-scripts-tui` (deduplicated surface types)

The proposal introduces **surface-key** suffixes for run-test:
- `T-test-run-backend`, `T-test-run-frontend` (map keys from config)

This is a fundamental inconsistency. The existing system passes `capabilities []string` (deduplicated surface types) to both `GetBreakdownTestTasks` and `GetQuickTestTasks` (autogen.go lines 103, 219). The proposal would require changing the function signatures to also accept the full `surfaces map[string]string` (surface-key -> surface-type mapping), since run-test needs surface-keys while gen-journeys/gen-scripts use surface-types.

**Risk**: A project with `surfaces: { auth: api, payments: api }` would generate:
- `T-test-gen-journeys-api` (single, deduplicated) -- correct
- `T-test-gen-scripts-api` (single, deduplicated) -- correct
- `T-test-run-auth`, `T-test-run-payments` (per surface-key) -- this is the new pattern

But `T-test-run-auth` and `T-test-run-payments` both have the same surface-type "api". The proposal's default priority ordering (api > web > cli > tui > mobile) is based on surface-type, not surface-key. With two api surfaces, the proposal says "report conflict, require explicit config." This is workable but means the common multi-api microservice pattern always requires manual configuration, undermining the "zero-config" claim.

### 2. InferType Contract Must Change

Currently `infer.go` line 22 uses **exact match** for `T-test-run`:
```go
case id == "T-test-run":
    return TypeTestRun
```

This is explicitly tested in `infer_test.go` line 53: `{"T-test-run-api", ""}` -- type-suffixed run IDs intentionally return empty string. The proposal's `T-test-run-{surface-key}` pattern requires changing this to use `typeSuffixedID`, which would require:
1. Adding `typeSuffixedID(id, "T-test-run")` to the case
2. Updating the test to expect `TypeTestRun` instead of `""`
3. Updating `ExtractTypeSuffix` support for the new base

The proposal does not mention this change. This is a **critical omission**.

### 3. isTestTaskID Pattern Must Evolve

`build.go` line 470-477:
```go
func isTestTaskID(id string) bool {
    return strings.HasPrefix(id, "T-test-") ||
        strings.HasPrefix(id, "T-quick-") ||
        ...
}
```

The `T-test-` prefix already covers `T-test-run-backend` and `T-test-run-frontend`, so `isTestTaskID` would work without changes. However, `isAutoGenForDep` (autogen.go line 880) delegates to `isTestTaskID`, so new run-test IDs would be correctly excluded from business task dependency resolution. This is **fine as-is**.

### 4. Single-surface Degradation Correctness

For `surfaces: api` (scalar form, stored as `{".": "api"}`), the proposal claims degradation to a single task identical to current behavior. Let's verify:

Current: `T-test-run` (single task, no suffix)
Proposed: `T-test-run-{surface-key}` where surface-key is `.` (dot)

This means the task would be `T-test-run-.` (with a literal dot), which is invalid. The scalar form `{".": "api"}` uses `"."` as a sentinel key (config.go line 158). The proposal must address this edge case -- when surfaces has a single entry with the `"."` key, the run-test should use the unsuffixed `T-test-run` form.

This is a **critical bug** in the proposal as written. The success criterion "single surface project tasks and dependency chain identical to before" cannot be met without special-casing the `"."` key.

## Dependency Chain Analysis

### resolveBreakdownDeps Changes

Current chain (autogen.go line 527):
```
gen-journeys-per-type -> eval-journey -> gen-contracts -> eval-contract -> gen-scripts-per-type -> run -> verify-regression
```

With the proposal's two changes:

**gen-journeys merge**: The eval-journey block (lines 535-540) currently collects deps from all `T-test-gen-journeys-{typ}` tasks. If gen-journeys becomes a single task (e.g., `T-test-gen-journeys`), eval-journey would depend on just that one task. This simplifies the dep chain. However, the proposal says gen-journeys should be "单任务" (single task), but doesn't specify the ID. Would it be `T-test-gen-journeys` (no suffix)? If so, `InferType` already handles exact match for this ID (infer.go line 18).

**run-test split**: The current `runIdx` lookup (line 531: `findTaskIndexOrPanic(tasks, "T-test-run")`) would need to be replaced with a loop over surface-keys. The verify-regression task (line 563) would need to depend on ALL `T-test-run-{key}` tasks instead of a single `T-test-run`. The proposal mentions this ("T-test-verify-regression depends on all run-test sub-tasks") but doesn't show the concrete code change.

The serial ordering between run-test tasks (e.g., `T-test-run-backend` before `T-test-run-frontend`) must be wired as:
```
T-test-run-backend depends on [all gen-scripts]
T-test-run-frontend depends on [T-test-run-backend]  // serial, not parallel
```

This is a **fundamental topology change**. Currently run-test is a fan-in node (all gen-scripts -> single run). The proposal makes it a serial chain (all gen-scripts -> first run -> second run -> ... -> verify-regression). The proposal must specify how `resolveBreakdownDeps` constructs this chain given the `execution-order` config.

### resolveQuickDeps Changes

Same topology change applies. Current (line 615): `run depends on all gen-scripts`. Proposed: first run-test depends on all gen-scripts, subsequent run-tests form a serial chain.

### ResolveFirstTestDep Impact

`ResolveFirstTestDep` (autogen.go line 693) uses `findTaskIndexByPrefix(tasks, "T-test-gen-journeys")` to find the first test task. If gen-journeys merges to a single `T-test-gen-journeys` (no suffix), `findTaskIndexByPrefix` still works. If it keeps some suffix pattern, this needs updating.

More critically, `resolveTestDepsAndInjectReviewDoc` (build.go line 535) injects `T-review-doc` as the first dependency of the first test pipeline task. With the gen-journeys merge, the first test task would be `T-test-gen-journeys` (single task), and this logic would still work. No issue here.

## Backward Compatibility Assessment

### Scenario 1: Existing single-surface project with `T-test-run` in index.json

An existing index.json contains `"run-test": {ID: "T-test-run", ...}`. After the change, the autogen would generate `T-test-run-.` (with the dot key) or some other ID. The `PreserveRuntimeFields` call in BuildIndex (build.go line 364) uses `index.ByID(td.ID)` to match -- if the new ID differs from `T-test-run`, the old task becomes an orphan (deleted in step 6 of BuildIndex, line 228). The user loses runtime state (status, blocked-reason) for the run-test task.

**Severity: HIGH**. The proposal must define a migration path: either recognize `T-test-run` as equivalent to the new per-surface-key task, or provide a rename mechanism.

### Scenario 2: Existing multi-surface project with multiple gen-journeys tasks

Before: `T-test-gen-journeys-api`, `T-test-gen-journeys-web` (two tasks).
After: `T-test-gen-journeys` (single task).

Same orphan problem: the old per-type gen-journeys tasks would be deleted, losing runtime state.

### Scenario 3: InferType breaks for `T-test-run-api`

The test at `infer_test.go` line 53 explicitly asserts that `T-test-run-api` returns `""`. After the change, this test would fail. The proposal must enumerate all test files that need updating.

### Scenario 4: `execution-order` references surface-keys that don't exist in `surfaces`

The proposal mentions validation but doesn't specify where this validation occurs. Should it be in `forgeconfig.ReadConfig` (config.go), in `BuildIndex` (build.go), or in a dedicated validator? If validation is in BuildIndex, it would only trigger on `forge task build`, not on config load.

## Edge Cases & Missing Details

### 1. Surface-key with special characters

Config allows arbitrary keys in `surfaces: { my-service/v2: api }`. The key would produce `T-test-run-my-service/v2`, which contains a `/` character. This is invalid for task IDs and filenames. The proposal must specify key normalization rules.

### 2. Empty execution-order with single surface

If `surfaces: { backend: api }` and `execution-order: []`, should the system use the single surface-key as-is or fall back to the unsuffixed `T-test-run`? The proposal doesn't specify.

### 3. execution-order contains surface-keys not in surfaces

The proposal says "report error" but doesn't specify the error message format, severity (warning vs. fatal), or timing (config load vs. build time).

### 4. Quick mode with the same changes

The proposal says "update resolveQuickDeps" but the quick mode pipeline has a different structure (no eval-journey, no eval-contract). The gen-journeys merge affects quick mode equally -- currently quick mode also generates per-type gen-journeys (autogen.go line 229). The proposal doesn't explicitly state that quick mode gen-journeys should also merge, though it's implied.

### 5. gen-journeys single task: what is the TestType field?

Currently gen-journeys tasks have `TestType: typ` (autogen.go line 119). If merged to a single task, what should `TestType` be? Empty string? A comma-separated list? The `renderBody` function (autogen.go line 317) omits the `{{TEST_TYPE}}` line when `TestType` is empty, which would change the generated .md content.

### 6. Task .md filename convention

The `FileName` field is derived from `Key` via `d.Key + ".md"` (autogen.go line 831). The proposal's new task keys (e.g., `run-test-backend`) would produce `run-test-backend.md`. The embed template lookup uses `autogenTemplatePath(def.Type)` which maps `test.run` to `data/test-run.md`. This template lookup is type-based, not key-based, so it would still work. But the generated file would be different from the template file, and `ValidateAutogenTemplates` (autogen.go line 27) only validates types, not keys. This is fine.

### 7. SourceTaskID interaction

When a run-test task fails and spawns fix-tasks, the fix-task records `SourceTaskID: "T-test-run"`. After splitting, the SourceTaskID would be `T-test-run-backend` etc. The `autoRestoreSourceTask` function (referenced in submit_test.go) does `index.ByID(fix.SourceTaskID)` to find the source. This should work as long as the new IDs are in the index. But migration of existing fix-tasks with `SourceTaskID: "T-test-run"` would break -- they'd reference a non-existent task.

## Findings

1. **[CRITICAL] Single-surface dot-key sentinel**: Scalar form `surfaces: api` produces `{".": "api"}`. Using surface-key as suffix would create `T-test-run-.`, which is invalid. Must special-case the `"."` sentinel key to produce unsuffixed `T-test-run` for single-surface projects. (config.go line 158, autogen.go line 155)

2. **[CRITICAL] InferType exact-match contract**: `T-test-run` uses exact match in `InferType` (infer.go line 22). The proposal's per-surface-key tasks require adding `typeSuffixedID` support. The existing test (infer_test.go line 53) explicitly asserts `T-test-run-api` returns empty string. This must be updated.

3. **[CRITICAL] Naming scheme inconsistency**: The proposal uses surface-key suffixes (`T-test-run-backend`) while the existing system uses surface-type suffixes (`T-test-gen-journeys-api`). This creates two different naming conventions in the same task set. The autogen functions currently only receive `capabilities []string` (deduplicated types), not the full surfaces map.

4. **[HIGH] Backward compatibility - index.json orphan risk**: Existing `T-test-run` entries in index.json would become orphans when replaced with `T-test-run-{key}`. Runtime state (status, blocked-reason) would be lost. Same for `T-test-gen-journeys-{type}` -> single `T-test-gen-journeys`. (build.go line 228)

5. **[HIGH] Backward compatibility - SourceTaskID migration**: Existing fix-tasks with `SourceTaskID: "T-test-run"` would reference a non-existent task after the rename. The auto-restore logic would silently fail. (submit_test.go lines 864-1039)

6. **[HIGH] Surface-key normalization not specified**: Arbitrary YAML map keys may contain characters invalid for task IDs/filenames (e.g., `/`, spaces, uppercase). No normalization rules are defined.

7. **[MEDIUM] gen-journeys TestType field after merge**: Single gen-journeys task has no meaningful TestType. The `renderBody` function's behavior with empty TestType differs from per-type tasks. Template content would change. (autogen.go line 117-119)

8. **[MEDIUM] execution-order validation timing unspecified**: Where and when should invalid surface-key references be caught? Config load time (fail fast) or build time (contextual)? The proposal says "forge task build" but earlier detection is better.

9. **[MEDIUM] Quick mode gen-journeys merge not explicit**: The proposal focuses on breakdown mode but quick mode has the same per-type gen-journeys loop (autogen.go line 229). The merge must apply to both modes.

10. **[LOW] Default priority ordering incomplete**: The proposed order (api > web > cli > tui > mobile) doesn't cover all pairwise combinations. What about projects with tui + cli surfaces? The proposal acknowledges this as low-risk but doesn't specify fallback behavior (alphabetical? config order?).

11. **[LOW] Missing `Config` struct field for execution-order**: The proposal mentions adding `ExecutionOrder []string` to config but the current `Config` struct (config.go line 201-210) doesn't have it. The proposal should show the exact field definition and YAML key name.

12. **[LOW] T-test-verify-regression dependency change scope**: The proposal states verify-regression depends on "all run-test sub-tasks" but doesn't show whether this is a fan-in (parallel) or serial chain attachment. Given that run-test tasks are serial, verify-regression should only depend on the LAST run-test task, not all of them, for correct scheduling semantics.

## Recommendations

1. **[P1] Fix the dot-key sentinel**: When `SurfacesMap` has a single entry with key `"."`, use unsuffixed `T-test-run` (current behavior). Only apply per-surface-key split when `len(surfaces) > 1` and no key is `"."`. This preserves backward compatibility for the common case.

2. **[P1] Align naming with surface-type, not surface-key**: Use `T-test-run-api` (surface-type) instead of `T-test-run-backend` (surface-key). For multi-same-type scenarios, use an indexed suffix: `T-test-run-api-1`, `T-test-run-api-2`. This keeps the naming convention consistent with gen-journeys/gen-scripts and avoids the surface-key normalization problem entirely.

3. **[P1] Update InferType before implementation**: Add `typeSuffixedID(id, "T-test-run")` to the type inference switch. Update the test. Document this as a prerequisite change.

4. **[P1] Define migration path for existing index.json**: Either (a) add a rename mapping from `T-test-run` to the new ID during BuildIndex, or (b) keep `T-test-run` as a deprecated alias that redirects to the new task, or (c) provide a `forge task migrate` command.

5. **[P2] Move execution-order validation to config load time**: Add validation in `ReadConfig` or a dedicated `ValidateExecutionOrder(surfaces, order)` function that fails fast with a clear error message. Do not defer to build time.

6. **[P2] Specify dependency topology explicitly**: Show the concrete chain for a 3-surface project (e.g., api + web + cli) in both breakdown and quick modes. Use a diagram or pseudocode to make the serial insertion logic unambiguous.

7. **[P2] Address SourceTaskID migration**: When renaming `T-test-run` to `T-test-run-{key}`, scan existing fix-tasks in the index and update their `SourceTaskID` field. Document this as a migration step.

8. **[P3] Consider alternative: execution-order as a run-time concern only**: Instead of splitting `T-test-run` into multiple tasks, keep it as a single task but pass the ordered surface list to the task executor. The SKILL.md for run-test would execute surfaces in order and stop on first failure. This avoids all the InferType/naming/migration complexity while achieving the same fail-fast behavior. The proposal rejected "internal ordering" but the rejection reason ("invisible to scheduler") may not justify the architectural cost shown in this review.
