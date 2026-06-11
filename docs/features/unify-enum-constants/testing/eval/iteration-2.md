# Journey Evaluation Report — Iteration 2

**Evaluator**: Adversarial QA Engineer
**Date**: 2026-05-29
**Surface Type**: CLI
**Journeys Evaluated**: 4 (status-migration, surface-type-migration, validation-map-consolidation, full-verification)
**Iteration 1 Score**: 466/1000

---

## Phase 1: Reasoning Audit

### status-migration/journey.md

**Problem -> Solution**: The journey addresses migrating Status string literals to typed constants. Solution is sound: define type, migrate, re-export, replace literals, add boundary conversions.

**Solution -> Evidence**: Steps use Go compilation checks (`go build`, `go test`) as evidence. Appropriate for a type-safety migration.

**Evidence -> Success Criteria**: Success criteria are clear (compilation passes, tests pass, JSON compatibility). Constant names now correctly list `StatusSuspended`, `StatusSkipped`, `StatusRejected` instead of the hallucinated names from iteration 1. **Fixed.**

**Self-contradiction check**: Step 6 claims "CLI Error -- Status Value Not Found" with `forge task status --filter unknown` but the actual CLI does not have a `--filter` flag for task status. The command `forge task status` takes a task ID, not a status filter. The scenario conflates two different commands. Step 7 claims "Duplicate Status Constant" but Go does NOT prevent duplicate constant values -- the journey correctly notes this, but then claims "Unit test for `AllStatuses()` uniqueness fails" which presupposes such a test exists.

### surface-type-migration/journey.md

**Problem -> Solution**: Sound. SurfaceType migration follows the same pattern.

**Solution -> Evidence**: References `detect_surface.go`, `execution_order.go`, `detect.go` which exist in `pkg/forgeconfig/`. Reasonable.

**Evidence -> Success Criteria**: Claims "~97 Surface Type string literals across 6 files" -- still unverifiable and no source classification.

**Self-contradiction check**: Step 2b references `forge surfaces <path>` command. Step 5 also references `forge surfaces ./unknown-project/`. Step 6 references `forge surfaces detect`. The `forge surfaces` command exists, but the journey uses it in inconsistent contexts -- sometimes as a path-based query, sometimes as a detection subcommand. The scenarios are reasonable but would benefit from consistent command syntax.

### validation-map-consolidation/journey.md

**Problem -> Solution**: The journey describes consolidating hardcoded validation maps. **BUT the actual `validate_index.go` already uses `types.AllStatuses()` and `types.AllPriorities()` via `buildValidStatusMap()`/`buildValidPriorityMap()`**. This journey describes work that is already completed, not pending work.

**Solution -> Evidence**: Step 1 claims: `validStatus map[string]bool{"pending": true, "in-progress": true, ...}`. This is **factually wrong** on two counts:
1. The actual code uses `buildValidStatusMap()` calling `types.AllStatuses()` -- no hardcoded map exists.
2. The value `"in-progress"` (hyphen) is wrong -- the actual constant is `"in_progress"` (underscore). This was flagged in iteration 1 for full-verification but NOT fixed in validation-map-consolidation.

**Evidence -> Success Criteria**: The setup precondition states "`internal/cmd/task/validate_index.go` contains hardcoded `validStatus` and `validPriority` maps" -- this is false. The file already uses `types.AllStatuses()`/`types.AllPriorities()`.

**Self-contradiction check**: The journey's own invariants state "Validation maps must be derived from `types.AllStatuses()`/`types.AllPriorities()` -- no hardcoded enum values in validation logic" -- but the journey's Steps 2 and 3 describe replacing hardcoded maps with helper functions, which is already done. The invariant is already satisfied before the journey begins.

### full-verification/journey.md

**Problem -> Solution**: Sound. Final verification of all migrations.

**Solution -> Evidence**: The grep command in Step 1 now correctly uses `"in_progress"` (underscore) and includes all 7 actual status values. **Fixed from iteration 1.**

**Evidence -> Success Criteria**: Verification grep is now correct. Exit code expectations are properly stated (exit code 1 when grep finds nothing, exit code 0 for successful builds).

**Self-contradiction check**: Step 4b claims `assert.Equal(t, "pending", task.Status)` would fail because `types.StatusPending` is not equal to untyped `"pending"`. This is **factually wrong** in Go -- `types.Status` is defined as `type Status string`, and untyped string constants are assignable to named string types. `assert.Equal` uses `reflect.DeepEqual`, and since `StatusPending = "pending"`, comparing `task.Status` (of type `Status`) to `"pending"` (untyped string constant) would succeed. The test would NOT fail as described.

---

## Phase 2: Rubric Scoring

### Dimension 1: Completeness (200 pts)

**Journey metadata (0-50)**:
- All 4 journeys have `feature`, `journey` (kebab-case), `risk_level` (High/Medium/Low), `surface_types`, `sources`, `generated` date. **Full marks: 50/50**

**Steps with required fields (0-80)**:
- All journeys have Steps with "User Action" and "Expected Result" headers. Edge cases have "Precondition" + "User Action" + "Expected Result". Steps have named headers.
- However, steps describe developer code-editing actions, not CLI invocations. "Developer creates `pkg/types/status.go`" is a code-writing step, not a testable CLI command. A downstream agent cannot execute "Developer creates a file" as a CLI test step. **-20 pts**.
- status-migration Step 5: "Developer adds `types.Status(stringVal)` conversions at CLI flag parsing boundaries" -- vague. Which CLI flags? Which entry points? Not actionable. **-10 pts**.
- surface-type-migration Step 4: Same vagueness -- "at config parsing and surface detection output boundaries" -- no specific entry points named. **-10 pts**.

**Score: 40/80**

**Outcomes cover happy path plus required derived scenarios (0-70)**:
- Happy path covered in all journeys.
- CLI error boundary outcomes added since iteration 1:
  - status-migration: Steps 6 ("Status Value Not Found") and 7 ("Duplicate Status Constant"). Step 6 addresses not-found. Step 7 partially addresses already-exists but for constants, not resources.
  - surface-type-migration: Steps 5 ("Surface Not Found") and 6 ("Duplicate Surface Type in Config"). Addresses not-found and already-exists.
  - validation-map-consolidation: Steps 5 ("Task Status Not Found") and 6 ("Priority Already Exists"). Addresses not-found and already-exists.
  - full-verification: NO CLI error boundary outcomes. Only happy path + edge cases. **Missing both not-found and already-exists.** **-15 pts**.
- The "not-found" and "already-exists" outcomes in the first 3 journeys are not structured as derived Outcomes within existing Steps. They are separate Steps entirely, which means they are not true "derived Outcomes" as the CLI surface rule envisions (alternative outcomes for the same action under different preconditions). They are separate actions. **-10 pts**.

**Score: 45/70**

**Completeness Total: 135/200** (above threshold 120)

---

### Dimension 2: Semantic Purity (200 pts)

**Outcome descriptions use natural language (0-80)**:
- Most outcomes are in natural language. "Exit code non-zero. stderr contains descriptive error message listing valid status values." -- good.
- full-verification Step 1 contains a grep command as a User Action: `grep -r '"pending"\|"in_progress"\|"completed"\|"blocked"\|"suspended"\|"skipped"\|"rejected"' --include='*.go'`. This is code, not natural language. The Expected Result references "Exit code 1 (grep found nothing)" which mixes CLI behavior with a specific tool. **-10 pts**.
- validation-map-consolidation Step 1 contains code: `validStatus map[string]bool{"pending": true, "in-progress": true, ...}`. This is a Go map literal in an Expected Result field. **-5 pts**.

**Score: 65/80**

**Preconditions are declarative (0-60)**:
- Edge case preconditions use "Precondition:" and are declarative. "A Status value exists in the codebase but is not defined as a constant" -- good. "AllStatuses() returns an empty slice due to a bug" -- declarative state, good.
- However, some preconditions are procedural in disguise: "Developer writes `task.Status == "pending"` (untyped string comparison)" in status-migration Step 5c. This prescribes an action within a precondition. **-10 pts**.

**Score: 50/60**

**No implementation coupling in Step descriptions (0-60)**:
- Steps are heavily coupled to specific file paths: `statemachine.go`, `add.go`, `state.go`, `deps.go`, `build.go`, `autogen.go`, `types.go`, `detect_surface.go`, `execution_order.go`, `detect.go`. If any file is renamed or restructured, the journey breaks. **-20 pts**.
- Steps describe code-editing actions ("Developer replaces all Status string literals in `statemachine.go`"), not observable CLI behavior. This is implementation-level detail, not user-facing behavior. **-15 pts**.

**Score: 25/60**

**Semantic Purity Total: 140/200** (above threshold 120)

---

### Dimension 3: Precondition Exclusivity (150 pts)

**Preconditions distinct across Outcomes (0-60)**:
- Within each step, happy path has no precondition (default) and edge cases have distinct preconditions. Good separation.
- status-migration: Step 1b ("A Status value exists in the codebase but is not defined as a constant") vs Step 3b ("A string literal was incorrectly replaced with the wrong typed constant") -- distinct. Good.
- surface-type-migration: Step 2b ("encounters a project with no recognizable surface signals") vs Step 3b ("config file specifies `surface: "desktop"` which is not a defined constant") -- distinct. Good.

**Score: 55/60**

**Preconditions sufficient to uniquely select an Outcome (0-50)**:
- Happy path outcomes have no preconditions, making them the default. Edge cases have explicit preconditions. This is sufficient for selection.
- validation-map-consolidation Step 2b: "`AllStatuses()` returns an empty slice due to a bug" -- this is the ONLY precondition that leads to "all statuses rejected." Sufficient.
- status-migration Step 5c precondition is "Developer writes `task.Status == "pending"`" -- this is both a precondition AND an action, blurring the boundary. An agent cannot tell if this is the starting state or the action to perform. **-10 pts**.

**Score: 40/50**

**No missing Preconditions for error/boundary Outcomes (0-40)**:
- CLI error boundary steps have preconditions:
  - status-migration Step 6: "CLI command references a status value that does not match any typed constant" -- good.
  - status-migration Step 7: "Developer accidentally defines two Status constants with the same string value" -- good.
  - surface-type-migration Step 5: "CLI command references a surface type that does not match any typed constant" -- good.
  - surface-type-migration Step 6: "`.forge/config.yaml` defines a surface key that resolves to a type already configured" -- good.
  - validation-map-consolidation Step 5: "forge task validate is run on an index with a task whose status is not in AllStatuses() output" -- good.
  - validation-map-consolidation Step 6: "Task index contains two tasks with the same ID but different priorities" -- good.
- full-verification has no CLI error boundary steps, so this is N/A for missing preconditions on non-existent steps. However, Step 5b ("A CLI command receives a status filter that does not match any defined constant") is present but lacks a clear precondition header -- the "Precondition" line is embedded in the step name. **-5 pts**.

**Score: 35/40**

**Precondition Exclusivity Total: 130/150** (above threshold 90)

---

### Dimension 4: Fact Alignment (150 pts)

**Factual claims traceable to fact_id or marked UNKNOWN (0-60)**:
- No fact_ids are used anywhere. No UNKNOWN markers. All claims are asserted without traceability. The journeys have `sources` in frontmatter pointing to PRD docs, but individual claims within steps are not linked to specific facts. Same issue as iteration 1. **-30 pts**.

**Score: 30/60**

**Inferred claims have required_outcomes rule support and source: inferred (0-50)**:
- Multiple claims are inferred but not classified:
  - "117 Status string literals across 22 files" (status-migration Setup) -- unverifiable, no source classification.
  - "~97 Surface Type string literals across 6 files" (surface-type-migration Setup) -- unverifiable, no source classification.
  - "JSON output unchanged -- `type Status string` marshals/unmarshals identically to plain `string`" (status-migration Step 5b) -- this is an inferred claim about Go behavior, correct but not classified.
  - None are marked as `source: inferred`. **-25 pts**.

**Score: 25/50**

**No hallucinated claims without classification (0-40)**:
- **HALLUCINATION 1**: validation-map-consolidation Step 1: `validStatus map[string]bool{"pending": true, "in-progress": true, ...}`. The actual code uses `buildValidStatusMap()` calling `types.AllStatuses()`. No hardcoded map exists. AND `"in-progress"` is wrong -- the actual value is `"in_progress"`. This is a **dual hallucination** (wrong code state + wrong constant value). **-30 pts** (per rule: -30 per instance).

- **HALLUCINATION 2**: validation-map-consolidation Setup: "`internal/cmd/task/validate_index.go` contains hardcoded `validStatus` and `validPriority` maps". The file DOES NOT contain hardcoded maps. It uses `buildValidStatusMap()` which calls `types.AllStatuses()`. The entire journey describes already-completed work. **-30 pts** (per rule: -30 per instance, capped at subdimension minimum).

- **FACTUAL ERROR**: full-verification Step 4b: "Test fails because `types.StatusPending` is not equal to untyped `"pending"` in strict equality." This is **factually wrong** in Go. `type Status string` makes `Status` a named type with `string` as underlying type. `reflect.DeepEqual(string("pending"), Status("pending"))` returns `true` because both have the same underlying string value. `assert.Equal` uses `reflect.DeepEqual` under the hood. The test would NOT fail as described. However, this is not a hallucination of a constant -- it is a wrong behavioral claim about Go semantics. **Not a -30 deduction but a factual error that undermines the step's validity.**

- **FACTUAL ERROR**: status-migration Step 6: `forge task status --filter unknown`. The actual CLI does not have a `--filter` flag on the `status` subcommand. `forge task status` takes a task ID positional argument, not a filter flag. The command syntax is invented.

- Constants in status-migration are NOW CORRECT: `StatusSuspended`, `StatusSkipped`, `StatusRejected` -- matching actual code. **Fixed from iteration 1.**

- full-verification grep NOW CORRECT: uses `"in_progress"` (underscore) and includes all 7 actual status values. **Fixed from iteration 1.**

- File path in validation-map-consolidation Setup NOW CORRECT: `internal/cmd/task/validate_index.go`. **Fixed from iteration 1.** (But the Setup claim about its contents is still wrong.)

**Score: 0/40** (two hallucinated unclassified claims: -30 each = -60, capped at 0)

**Fact Alignment Total: 55/150** (BELOW threshold 90 -- **FAIL**)

---

### Dimension 5: Surface Fitness (150 pts)

**Mandatory derived Outcomes (not-found + already-exists) (0-60)**:
- status-migration: Step 6 = "Status Value Not Found" (not-found), Step 7 = "Duplicate Status Constant" (partial already-exists). **Present.**
- surface-type-migration: Step 5 = "Surface Not Found" (not-found), Step 6 = "Duplicate Surface Type in Config" (already-exists). **Present.**
- validation-map-consolidation: Step 5 = "Task Status Not Found" (not-found), Step 6 = "Priority Already Exists" (already-exists). **Present.**
- full-verification: **MISSING both not-found and already-exists derived outcomes.** No CLI error boundary steps. The journey is "Low risk / read-only" but this does not exempt it from the CLI surface rule which states "Mandatory derived Outcomes (must be considered for every CLI Journey)". **-20 pts** for full-verification missing both.
- Additionally, the "not-found" and "already-exists" outcomes in the first 3 journeys are implemented as **separate Steps** rather than **alternative Outcomes within existing Steps**. The surface-cli.md says these are "derived Outcomes" -- meaning they should be alternative results of the same action under different preconditions, not separate actions entirely. For example, status-migration Step 4 (replace magic values) should have a derived outcome for "status value not found in constants" if the replacement fails. Instead, Step 6 is an entirely different action ("Developer runs the CLI command with an unrecognized status value"). **-15 pts** for structural mismatch with CLI surface model.

**Score: 25/60**

**Test strategy proportions match CLI guidance (Contract 80% / Journey 20%) (0-50)**:
- None of the journeys mention test strategy proportions or the Contract/Journey test split. The CLI surface rule states: "Test Level Emphasis: Contract 80% / Journey smoke 20%". No journey acknowledges this. No journey indicates which steps would be Contract tests vs Journey smoke tests. **-25 pts**.
- The journeys are structured as sequential workflows (Journey-level) but lack any indication of which individual command behaviors would be tested at the Contract level. Without this split, a downstream agent cannot determine the correct test architecture.

**Score: 25/50**

**CLI-specific environment assumptions realistic (subprocess isolation, exit code checks) (0-40)**:
- Steps reference `go build`, `go test`, `grep` commands with exit codes. Good.
- Steps also describe code-editing actions: "Developer creates `pkg/types/status.go`", "Developer replaces all Status string literals in `statemachine.go`". These are NOT subprocess-invokable commands. A CLI test agent cannot execute these as subprocess tests. **-20 pts**.
- Only full-verification Step 5 mentions actual CLI commands (`forge task list`, `forge task status <id>`). The other journeys have no subprocess-isolated CLI command invocations except the newly added error boundary steps. **-10 pts**.

**Score: 10/40**

**Surface Fitness Total: 60/150** (BELOW threshold 90 -- **FAIL**)

---

### Dimension 6: Internal Consistency (150 pts)

**Invariants hold in every Step (0-60)**:
- status-migration Invariant: "All typed constant values must exactly match the original string literals ... zero behavior change". The constants listed in Step 1 now match the actual codebase. **Invariant holds. Fixed from iteration 1.**
- validation-map-consolidation Invariant: "Validation maps must be derived from `types.AllStatuses()`/`types.AllPriorities()` -- no hardcoded enum values". But Step 1 describes finding hardcoded maps and Steps 2-3 describe replacing them. If the invariant already holds (which it does in the actual code), then Steps 2-3 are unnecessary. **The invariant contradicts the need for the journey.** **-20 pts**.
- full-verification Invariant: "Zero Status/SurfaceType/Priority string literals in production code". Step 4b describes a test using `"pending"` as a string literal -- but test files are excluded from the grep in Step 1, so this is not an invariant violation. Consistent. However, the claim in Step 4b that the test would FAIL is incorrect (see Fact Alignment). The invariant is stated correctly but the evidence supporting it (the test failure) is fabricated. **-10 pts**.

**Score: 30/60**

**Cross-Step references consistent (0-50)**:
- status-migration: Step 1 defines 7 constants. Steps 3-4 reference migration of those constants. Step 5 adds boundary conversions. Progression is logical and consistent. **Good.**
- surface-type-migration: Step 1 defines 5 constants. Steps 2-3 replace literals. Step 4 adds boundary conversions. Consistent. **Good.**
- validation-map-consolidation: Step 1 identifies hardcoded maps. Steps 2-3 replace them. Step 4 verifies. But the Setup says "All Status and Priority magic values already migrated (prerequisite: journeys 1 and status-migration)" which contradicts Step 1's claim that hardcoded maps still exist. If magic values are already migrated, the validation maps should already use the helpers. **-15 pts**.

**Score: 35/50**

**Risk level consistent with Journey content (0-40)**:
- status-migration: High risk. Involves state machine mutation and type system changes. **Consistent.**
- surface-type-migration: High risk. ~97 occurrences across 6 files. **Consistent.**
- validation-map-consolidation: Medium risk. Multi-step without irreversible side effects. **Consistent.**
- full-verification: Low risk. Read-only verification. **Consistent.**
- **Full marks: 40/40**

**Internal Consistency Total: 105/150** (above threshold 90)

---

## Phase 3: Blindspot Hunt

### [blindspot-1] validation-map-consolidation Describes Already-Completed Work
The journey's Setup claims "contains hardcoded `validStatus` and `validPriority` maps" but the actual `validate_index.go` (lines 38-56) already uses `buildValidStatusMap()` and `buildValidPriorityMap()` which call `types.AllStatuses()` and `types.AllPriorities()`. The entire journey is a ghost -- it describes work that is already done. A downstream agent would waste time looking for hardcoded maps that don't exist.

### [blindspot-2] "in-progress" Still Present in validation-map-consolidation
Iteration 1 flagged `"in-progress"` (hyphen) in full-verification and it was fixed. But the same error persists in validation-map-consolidation Step 1: `"pending": true, "in-progress": true, ...}`. The actual constant value is `"in_progress"` (underscore). This was NOT fixed.

### [blindspot-3] full-verification Step 4b Fabricates Go Language Behavior
The step claims `assert.Equal(t, "pending", task.Status)` would fail because "types.StatusPending is not equal to untyped "pending" in strict equality". In Go, `reflect.DeepEqual` (used by `assert.Equal`) compares underlying values, not types. `Status("pending")` and `"pending"` would be equal. The entire edge case is based on a wrong understanding of Go semantics.

### [blindspot-4] status-migration Step 6 Uses Non-Existent CLI Flag
The step uses `forge task status --filter unknown`. The actual `forge task status` command takes a positional task ID argument, not a `--filter` flag. The correct scenario would be running `forge task status 99` for a non-existent task, or `forge task list --status unknown` for an invalid status filter. The invented flag syntax makes the step unexecutable.

### [blindspot-5] No Subprocess Isolation Model
The CLI surface rule requires: "The CLI under test must run as an isolated subprocess" with "temporary directories for all file I/O". None of the journeys mention subprocess isolation, temp directories, or binary compilation. Steps describe manual developer actions, not automated subprocess invocations.

### [blindspot-6] full-verification Missing Mandatory CLI Outcomes
While the first 3 journeys now have "not-found" and "already-exists" steps (though structured as separate Steps rather than derived Outcomes), full-verification has NO CLI error boundary outcomes at all. This is a gap in the coverage net -- the verification journey should verify that error handling works for invalid inputs.

---

## Per-Journey Dimension Breakdown

| Dimension | status-migration | surface-type-migration | validation-map-consolidation | full-verification |
|-----------|:-:|:-:|:-:|:-:|
| Completeness | 35/50 | 35/50 | 35/50 | 30/50 |
| Steps complete | 60/80 | 55/80 | 60/80 | 45/80 |
| Outcomes | 55/70 | 55/70 | 55/70 | 40/70 |
| **Completeness** | **150/200** | **145/200** | **150/200** | **115/200** |
| Semantic Purity | 140/200 | 145/200 | 130/200 | 145/200 |
| Precondition Exclusivity | 135/150 | 135/150 | 130/150 | 120/150 |
| Fact Alignment | 80/150 | 90/150 | 20/150 | 30/150 |
| Surface Fitness | 70/150 | 75/150 | 55/150 | 40/150 |
| Internal Consistency | 115/150 | 120/150 | 80/150 | 105/150 |

## Overall Averaged Scores

| Dimension | Average Score | Threshold | Status |
|-----------|:---:|:---:|:---:|
| 1. Completeness | **140/200** | 120 | PASS |
| 2. Semantic Purity | **140/200** | 120 | PASS |
| 3. Precondition Exclusivity | **130/150** | 90 | PASS |
| 4. Fact Alignment | **55/150** | 90 | **FAIL** |
| 5. Surface Fitness | **60/150** | 90 | **FAIL** |
| 6. Internal Consistency | **105/150** | 90 | PASS |
| **Total** | **630/1000** | 850 | **FAIL** |

---

## Attack Points Summary

1. **Fact Alignment**: validation-map-consolidation Step 1 hallucinates `validStatus map[string]bool{"pending": true, "in-progress": true, ...}` -- actual code uses `buildValidStatusMap()` calling `types.AllStatuses()`. AND `"in-progress"` should be `"in_progress"`. Must reflect actual code state or mark journey as already-completed.

2. **Fact Alignment**: validation-map-consolidation Setup claims "contains hardcoded `validStatus` and `validPriority` maps" in `validate_index.go` -- actual file already uses `types.AllStatuses()`/`types.AllPriorities()`. Must update Setup to reflect reality or remove journey.

3. **Fact Alignment**: full-verification Step 4b fabricates Go behavior -- `assert.Equal(t, "pending", task.Status)` would NOT fail in Go; `reflect.DeepEqual` compares underlying values. Must remove or rewrite this edge case with a correct Go semantic.

4. **Fact Alignment**: status-migration Step 6 uses `forge task status --filter unknown` -- the `--filter` flag does not exist on this command. Must use correct CLI syntax (e.g., `forge task status 99` or `forge task list --status unknown`).

5. **Surface Fitness**: full-verification missing mandatory `not-found` and `already-exists` derived outcomes. CLI surface rule states "must be considered for every CLI Journey". Must add error boundary steps.

6. **Surface Fitness**: All journeys structure "not-found"/"already-exists" as separate Steps rather than derived Outcomes within existing Steps. Must restructure as alternative outcomes with preconditions on existing actions.

7. **Surface Fitness**: No journeys mention subprocess isolation, temp directories, or binary compilation as required by CLI surface rule. Steps describe developer actions, not subprocess invocations. Must add subprocess isolation model.

8. **Surface Fitness**: No journeys acknowledge the Contract 80% / Journey 20% test strategy split. Must annotate which steps are Contract vs Journey smoke tests.

9. **Internal Consistency**: validation-map-consolidation invariant ("Validation maps must be derived from types.AllStatuses()") contradicts the journey's Steps which describe replacing hardcoded maps -- the invariant already holds in the current code. Must resolve the contradiction.

10. **Semantic Purity**: Steps heavily coupled to specific filenames (`statemachine.go`, `add.go`, etc.) and describe code-editing actions ("Developer replaces"). Must describe observable behavior instead of implementation details.

---

## Verdict

**TOTAL: 630/1000 -- FAIL**

Two dimensions below minimum threshold:
- Fact Alignment: 55/150 (threshold 90)
- Surface Fitness: 60/150 (threshold 90)

Key improvements from iteration 1:
- Status constant names corrected (StatusSuspended, StatusSkipped, StatusRejected)
- full-verification grep uses correct `"in_progress"` (underscore) and all 7 values
- File path corrected to `internal/cmd/task/validate_index.go`
- CLI error boundary steps added to 3 of 4 journeys

Remaining critical issues:
- validation-map-consolidation describes already-completed work (the whole journey is a ghost)
- `"in-progress"` (hyphen) persists in validation-map-consolidation Step 1
- full-verification missing mandatory CLI derived outcomes
- Steps describe code-editing, not CLI subprocess behavior
