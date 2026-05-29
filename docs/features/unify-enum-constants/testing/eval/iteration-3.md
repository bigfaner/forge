# Journey Eval — Iteration 3

**Date**: 2026-05-29
**Evaluator**: Adversary (strict)
**Journeys evaluated**: 4

---

## Phase 1: Reasoning Audit

### status-migration

1. **Problem -> Solution**: Clear. Migrate 117 Status string literals to typed constants. Steps are logical and sequential.
2. **Solution -> Evidence**: Each step has `go build` or `go test` as verification. Adequate.
3. **Evidence -> Success Criteria**: Compilation success + test pass. Well-defined.
4. **Self-contradiction check**:
   - Step 6 title says "Status Value Not Found" but Precondition describes a task ID not found scenario: `"forge task status 99 when task 99 does not exist"`. This conflates status-value-not-found with task-not-found. The step title and precondition describe different failure modes. **CONTRADICTION.**
   - Invariant says "Every Status field assignment in production code must use `types.StatusXxx` constants, never raw string literals" but `statemachine.go` uses `types.Status("*")` for wildcard matching (line 53, 57, 61, 68, 69, etc.). This is a raw string literal `"*"` cast to Status at the call site. The invariant is violated by the actual code. **CONTRADICTION with reality.**

### surface-type-migration

1. **Problem -> Solution**: Clear. Migrate ~97 Surface Type string literals.
2. **Solution -> Evidence**: Compilation checks at each step.
3. **Evidence -> Success Criteria**: Builds pass.
4. **Self-contradiction check**:
   - Step 3 says "remaining files (`execution_order.go`, `detect.go`, and other forgeconfig files)" but the actual path is `pkg/forgeconfig/` and the main file is `detect_surface.go`. The file names `execution_order.go` and `detect.go` are unverified. The journey says "~97 Surface Type string literals across 6 files" but provides no evidence for this count. **UNVERIFIED FACTUAL CLAIM.**
   - Step 6 Precondition: "`.forge/config.yaml` defines a surface key that resolves to a type already configured" -- this describes an already-exists scenario for config file entries. The expected result is vague: "The already-existing surface type is reported. No duplicate entries created. If validation enforces uniqueness, error is reported on stderr." The "If" conditional is ambiguous -- the outcome is uncertain. **AMBIGUOUS OUTCOME.**

### validation-map-consolidation

1. **Problem -> Solution**: Verify validation maps use helper functions. Clear.
2. **Solution -> Evidence**: Code inspection + unit tests. Good.
3. **Evidence -> Success Criteria**: Tests pass, no hardcoded maps.
4. **Self-contradiction check**:
   - Step 6 "Priority Already Exists" Precondition says "Task index contains two tasks with the same ID but different priorities". Two tasks with the same ID is impossible in a Go map (`map[string]Task`) -- the second would overwrite the first. The precondition describes an impossible state. **IMPOSSIBLE PRECONDITION.**
   - Step 4b "Validation of Case Sensitivity": User Action says "Developer runs `forge task validate` on the config" but the actual CLI command is `forge validate-index` (per `validateIndexCmd.Use`). **FACTUAL ERROR in CLI command name.**

### full-verification

1. **Problem -> Solution**: Verify complete migration. Clear.
2. **Solution -> Evidence**: grep + build + test + CLI smoke test.
3. **Evidence -> Success Criteria**: Zero matches, builds pass, tests pass.
4. **Self-contradiction check**:
   - Step 4b: Claims `assert.Equal(t, "pending", task.Status)` passes because "reflect.DeepEqual compares underlying values." This is correct for Go semantics -- `assert.Equal` uses `reflect.DeepEqual`, and `types.Status` has underlying type `string`, so the comparison succeeds. **CORRECT.** However, the Expected Result says "Test passes because... types.Status has underlying type string and reflect.DeepEqual compares underlying values" -- this is a technical explanation, not a clear pass/fail criterion. The criterion mixes explanation with assertion.
   - Step 5 says "Maintainer runs sample CLI commands (`forge task list`, `forge task status <id>`, etc.)" but Step 5b uses `forge task list --status unknown`. The `--status` flag behavior is not documented anywhere in the codebase as a verified feature. **UNVERIFIED CLI FLAG.**
   - Still missing mandatory not-found and already-exists derived outcomes. Step 5b is a not-found scenario but it is not structured as a derived outcome of an existing step. Iteration 2 flagged this and it was not addressed.

---

## Phase 2: Rubric Scoring

### Dimension 1: Completeness (200 pts)

**Journey metadata (0-50)**:
- All 4 journeys have correct kebab-case names, valid risk_level (High/Medium/Low), surface_types: ["cli"], and generated date. **45/50** (dates are all identical "2026-05-29", suggesting bulk generation rather than individual consideration, but acceptable).

**Steps complete with required fields (0-80)**:
- All 4 journeys have Steps with User Action and Expected Result. No missing required fields. **70/80** (no explicit "step name" field -- steps use markdown headers as names, which is acceptable but less machine-parseable).

**Outcomes cover happy path plus required derived scenarios (0-70)**:
- Happy path covered in all 4 journeys.
- Edge cases present in all journeys.
- full-verification still lacks CLI error boundary outcomes as derived scenarios (not-found/already-exists). It has Step 5b/5c but they are not "derived outcomes" of the verification step.
- status-migration Step 7 "Duplicate Status Constant" Expected Result says "Go allows duplicate constant values" but doesn't verify what happens in practice.
- **50/70** (full-verification missing derived scenarios; impossible precondition in validation-map Step 6).

**Dimension 1 Total: 165/200**

---

### Dimension 2: Semantic Purity (200 pts)

**Outcome descriptions use natural language, not code/regex (0-80)**:
- Most outcomes use natural language. Some use code artifacts:
  - status-migration Step 5c Expected Result: "Compiler allows the comparison (untyped string constant is assignable to `types.Status`), but linter/review should flag it." -- mixes Go type system terminology with normative guidance.
  - full-verification Step 4b Expected Result: "Test passes because `types.Status` has underlying type `string` and `reflect.DeepEqual` compares underlying values." -- deeply technical implementation explanation rather than behavioral description.
  - full-verification Step 1: `grep -r '"pending"\|"in_progress"\|"completed"...` -- this is literally a regex command in the User Action, not a natural language description.
- **55/80** (mix of natural language and code/implementation details across journeys).

**Preconditions are declarative, not procedural (0-60)**:
- Most preconditions are declarative: "A Status value exists in the codebase but is not defined as a constant"
- Some are procedural:
  - surface-type-migration Step 4b Precondition: "Two packages compare Surface Type values using different representations (one typed, one string)" -- describes a procedural scenario rather than a state.
  - full-verification Step 4b Precondition: "A test uses `assert.Equal(t, "pending", task.Status)` with untyped string" -- describes specific code rather than a declarative state.
- **45/60**

**No implementation coupling in Step descriptions (0-60)**:
- status-migration Step 3: "Developer replaces all Status string literals in `statemachine.go`" -- tightly coupled to specific file name.
- surface-type-migration Step 2: "Developer replaces Surface Type string literals in `pkg/forgeconfig/detect_surface.go`" -- specific file path.
- validation-map-consolidation Step 1: "Developer inspects `internal/cmd/task/validate_index.go`" -- specific file path.
- full-verification Step 1: literal grep command in User Action.
- These are migration journeys, so some file-level specificity is expected, but the coupling to exact file paths and exact commands violates semantic purity.
- **35/60**

**Dimension 2 Total: 135/200**

---

### Dimension 3: Precondition Exclusivity (150 pts)

**Preconditions distinct across Outcomes within each Step (0-60)**:
- Most steps have only one outcome (happy path) with edge cases as separate steps (Step Xb, Xc). The journey format does not structure multiple outcomes within a single step. This means preconditions are trivially "distinct" because there's no ambiguity -- but also no real exclusivity test because the format avoids the problem rather than solving it.
- **40/60** (format avoids the problem rather than demonstrating distinct preconditions).

**Preconditions sufficient to uniquely select an Outcome (0-50)**:
- Same issue: with one outcome per step, selection is trivially unique.
- For multi-precondition edge cases (e.g., status-migration Step 6 vs Step 7), the preconditions are distinct enough.
- **40/50**

**No missing Preconditions for error/boundary Outcomes (0-40)**:
- validation-map-consolidation Step 6: Precondition "Task index contains two tasks with the same ID but different priorities" -- IMPOSSIBLE. Two entries with the same key in a Go map cannot coexist. Missing realistic precondition for this scenario.
- status-migration Step 6: Precondition conflates "status value not found" with "task not found" (title says status, precondition says task 99 doesn't exist).
- full-verification: No precondition for what happens when `go test ./...` encounters a test binary compilation failure (vs test failure). Step 4 only considers test failures.
- **20/40**

**Dimension 3 Total: 100/150**

---

### Dimension 4: Fact Alignment (150 pts)

**Factual claims traceable to fact_id or marked UNKNOWN (0-60)**:
- No fact_ids are used anywhere in any journey. None of the journeys reference a fact table or fact source system. The claims are stated as bare assertions.
- Verified claims against actual code:
  - CORRECT: `buildValidStatusMap()` calls `types.AllStatuses()` (confirmed in validate_index.go line 44)
  - CORRECT: `buildValidPriorityMap()` calls `types.AllPriorities()` (confirmed in validate_index.go line 53)
  - CORRECT: `type Status string` with 7 constants (confirmed in status.go)
  - CORRECT: `type SurfaceType string` with 5 constants (confirmed in surface.go)
  - CORRECT: `type Priority string` with 3 constants (confirmed in priority.go)
  - CORRECT: `IsTerminalStatus` includes completed, rejected, skipped (confirmed in status.go line 36)
  - CORRECT: `pkg/feature/constants.go` re-exports via `type Status = types.Status` (confirmed line 75)
  - INCORRECT: validation-map-consolidation Step 4b User Action says `forge task validate` but actual CLI command is `forge validate-index` (per validateIndexCmd.Use = "validate-index [file]")
- **30/60** (no fact_id system used; one factual error found).

**Inferred claims have required_outcomes rule support and source: inferred (0-50)**:
- status-migration Setup claims "117 Status string literals across 22 files" -- no source, no evidence, no "inferred" tag. This is an unclassified factual claim.
- surface-type-migration Overview claims "~97 Surface Type string literals across 6 files" -- same issue.
- No "source: inferred" markings anywhere.
- **15/50**

**No hallucinated claims without classification (0-40)**:
- status-migration "117 Status string literals across 22 files" -- unverifiable, no classification. **-30 pts.**
- surface-type-migration "~97 Surface Type string literals across 6 files" -- unverifiable, no classification. **-30 pts.**
- validation-map-consolidation Step 4b: `forge task validate` command does not exist. The actual command is `forge validate-index`. This is a hallucinated CLI command. **-30 pts.**
- surface-type-migration Step 3 mentions "execution_order.go", "detect.go" as file names. `detect_surface.go` exists but "execution_order.go" and "detect.go" are not verified. **-30 pts.** (partial -- detect_surface.go exists but the other names are unverified).
- Total hallucination deductions: -120 pts, but this sub-dimension is floored at 0.
- **0/40**

**Dimension 4 Total: 45/150** (floored; 30 + 15 + 0 = 45, with -120 in hallucination deductions absorbed)

---

### Dimension 5: Surface Fitness (150 pts)

**Mandatory derived Outcomes from CLI surface `required_outcomes` (not-found + already-exists) are present (0-60)**:
- status-migration: Step 6 "Status Value Not Found" = not-found (present). Step 7 "Duplicate Status Constant" = partial already-exists for compile-time duplicates (present but semantically different from CLI resource already-exists). **Present with caveat.**
- surface-type-migration: Step 5 "Surface Not Found" = not-found (present). Step 6 "Duplicate Surface Type in Config" = already-exists (present).
- validation-map-consolidation: Step 5 "Task Status Not Found" = not-found (present). Step 6 "Priority Already Exists" = already-exists (present but impossible precondition).
- full-verification: **MISSING both.** Step 5b is a not-found scenario but is not a derived outcome. No already-exists scenario at all. Iteration 2 flagged this explicitly: "full-verification missing mandatory `not-found` and `already-exists` derived outcomes." This was not addressed in iteration 3.
- All journeys structure these as separate Steps rather than derived Outcomes within existing Steps (alternative outcomes for the same action). The CLI surface rule says "Mandatory derived Outcomes (must be considered for every CLI Journey)." This structural mismatch was flagged in iteration 2 and not fixed.
- **25/60** (3/4 journeys have approximate coverage; full-verification missing; structural mismatch persists).

**Test strategy proportions match CLI guidance (Contract 80% / Journey 20%) (0-50)**:
- No test strategy section in any journey. No mention of contract vs journey test proportion. The CLI surface rule says "Contract 80% / Journey smoke 20%" but the journeys make no reference to this allocation.
- The journeys themselves are journey-level smoke tests. There is no indication of how contract-level testing would cover the individual steps.
- **10/50**

**CLI-specific environment assumptions realistic (subprocess isolation, exit code checks) (0-40)**:
- status-migration Step 6: "Exit code non-zero" -- correct exit code assertion.
- surface-type-migration Step 2b: "Exit code 1" -- correct.
- full-verification Step 1: "Exit code 1 (grep found nothing)" -- correct (grep returns 1 when no matches).
- However, no journey mentions subprocess isolation, temporary directories, or environment variable isolation -- all required by the CLI surface rule ("Subprocess isolation" section).
- No journey mentions compiling a dedicated test binary or using `go test -c`.
- Most "CLI commands" in the journeys are conceptual actions by a developer, not automated test invocations. The journeys describe manual developer workflows, not automated CLI tests.
- **15/40**

**Dimension 5 Total: 50/150**

---

### Dimension 6: Internal Consistency (150 pts)

**Invariants hold in every Step (0-60)**:
- status-migration Invariant: "Every Status field assignment in production code must use `types.StatusXxx` constants, never raw string literals." VIOLATED by statemachine.go which uses `types.Status("*")` -- this is a string literal "*" wrapped in a type conversion. The wildcard `"*"` is not a named constant like `StatusPending`. The invariant as stated forbids this pattern. **-40 pts.**
- status-migration Invariant: "go build ./... must pass after each migration step (incremental correctness)." Step 1b Expected Result says "Compiler reports type mismatch error for the missing constant. Exit code non-zero." -- this means `go build` does NOT pass in this step, violating the invariant. The invariant should say "except during intentional error demonstration" or the step should be framed differently. **CONTRADICTION.**
- surface-type-migration Invariant: "All typed constant values must exactly match the original string literals (zero behavior change)." Reasonable and verified.
- validation-map-consolidation Invariant: "Validation behavior must be identical before and after consolidation." Reasonable.
- full-verification Invariant: "Zero Status/SurfaceType/Priority string literals in production code." Same issue as status-migration -- `types.Status("*")` contains a string literal.
- **20/60**

**Cross-Step references are consistent (0-50)**:
- status-migration Step 2 says "moves Status constant definitions from `pkg/feature/constants.go` to `pkg/types/status.go`, and adds re-export" -- but the actual code shows `pkg/feature/constants.go` already has re-exports via type alias (`type Status = types.Status`). The step describes creating the re-export, but in the actual code the re-exports already exist. The journey describes a migration that has already happened.
- validation-map-consolidation references "prerequisite: status-migration and surface-type-migration journeys" -- consistent.
- full-verification references "All enum migration journeys are completed" -- consistent.
- **35/50**

**Risk level consistent with Journey content (0-40)**:
- status-migration: High risk. Content involves data migration across 22 files. Justified.
- surface-type-migration: High risk. Content involves ~97 replacements. Justified.
- validation-map-consolidation: Medium risk. Verification-only. Justified.
- full-verification: Low risk. Read-only verification. Justified.
- **35/40**

**Dimension 6 Total: 90/150**

---

## Phase 3: Blindspot Hunt

[blindspot-1] **Unreachable precondition**: validation-map-consolidation Step 6 says "Task index contains two tasks with the same ID but different priorities." This is impossible in Go -- a `map[string]Task` cannot have two entries with the same key. The test case cannot be triggered. No agent could execute this step. This was present in iteration 2 and not fixed.

[blindspot-2] **CLI command hallucination**: validation-map-consolidation Step 4b uses `forge task validate` but the actual CLI command is `forge validate-index`. Any agent executing this step would get a "unknown command" error. This was not caught in iterations 1 or 2.

[blindspot-3] **Wildcard status constant violates invariant**: `statemachine.go` uses `types.Status("*")` extensively for wildcard matching in transition rules. This is a raw string literal `"*"` not represented as a named constant. The invariant "Every Status field assignment in production code must use `types.StatusXxx` constants, never raw string literals" is violated by the actual implementation. Either the invariant is wrong (it should carve out wildcard matching) or the implementation is incomplete.

[blindspot-4] **No test strategy alignment**: The CLI surface rule mandates "Contract 80% / Journey smoke 20%" but none of the 4 journeys reference this allocation or describe how contract tests would cover individual steps. The journeys are pure journey-level documents with no contract test scaffolding.

[blindspot-5] **full-verification Step 4b accuracy nuance**: While it is correct that `assert.Equal(t, "pending", task.Status)` passes (because `reflect.DeepEqual` compares underlying string values), the Expected Result says "the test should be updated to use `types.StatusPending` for type safety." This is a normative recommendation, not an observable outcome. An agent executing this step cannot verify "should be updated" -- it can only verify the test passes.

[blindspot-6] **Missing `go vet` step**: full-verification runs `go build` and `go test` but never runs `go vet` or `golangci-lint`. The forge-cli CLAUDE.md explicitly lists `go vet ./...` and `golangci-lint run ./...` as common commands. A verification journey should include lint checks.

[blindspot-7] **No race condition test**: The forge-cli CLAUDE.md specifies `go test -race -cover ./...` as the standard test command. No journey mentions `-race` flag. For a migration that changes type signatures, race conditions could be introduced.

---

## Final Scores

| Dimension | Score | Min Threshold | Pass? |
|-----------|-------|---------------|-------|
| 1. Completeness | 165/200 | 120 | YES |
| 2. Semantic Purity | 135/200 | 120 | YES |
| 3. Precondition Exclusivity | 100/150 | 90 | YES |
| 4. Fact Alignment | 45/150 | 90 | **NO** |
| 5. Surface Fitness | 50/150 | 90 | **NO** |
| 6. Internal Consistency | 90/150 | 90 | YES |
| **TOTAL** | **585/1000** | 850 | **NO** |

---

## Attack Summary

1. **Fact Alignment**: Three unclassified hallucinated claims -- "117 Status string literals across 22 files" (status-migration), "~97 Surface Type string literals across 6 files" (surface-type-migration), and `forge task validate` (validation-map-consolidation Step 4b, actual command is `forge validate-index`). Each costs -30 pts. File names "execution_order.go" and "detect.go" in surface-type-migration Step 3 are unverified.

2. **Surface Fitness**: full-verification still missing mandatory not-found and already-exists derived outcomes (flagged in iteration 2, not fixed). All journeys structure these as separate Steps rather than derived Outcomes within existing Steps (flagged in iteration 2, not fixed). No test strategy proportion guidance. No subprocess isolation or environment setup guidance.

3. **Internal Consistency**: Invariant "never raw string literals" is violated by `types.Status("*")` in statemachine.go. Invariant "go build passes after each step" is contradicted by Step 1b where build intentionally fails.

4. **Precondition Exclusivity**: validation-map-consolidation Step 6 has an impossible precondition (two tasks with same ID in a Go map). status-migration Step 6 conflates "status value not found" (title) with "task not found" (precondition).

5. **Semantic Purity**: full-verification Step 1 uses literal grep regex in User Action. Step 4b gives implementation explanation rather than behavioral outcome. Multiple journeys use exact file paths and CLI commands in steps.
