# Journey Evaluation Report — Iteration 1

**Evaluator**: Adversarial QA Engineer
**Date**: 2026-05-29
**Surface Type**: CLI
**Journeys Evaluated**: 4 (status-migration, surface-type-migration, validation-map-consolidation, full-verification)

---

## Phase 1: Reasoning Audit

### status-migration/journey.md

**Problem -> Solution**: The journey addresses migrating Status string literals to typed constants. The solution is sound: define type, migrate, re-export, replace literals, add boundary conversions.

**Solution -> Evidence**: Steps are well-grounded in Go compilation checks (`go build`, `go test`), which is appropriate evidence for a type-safety migration.

**Evidence -> Success Criteria**: Success criteria are clear (compilation passes, tests pass, JSON compatibility). However, the step 1 claim of 7 constants is **factually wrong** -- it lists `StatusCancelled`, `StatusFailed`, `StatusReview` which **do not exist** in the actual codebase.

**Self-contradiction check**: Step 1 lists 7 constants including 3 that don't exist, but the actual `pkg/types/status.go` has 7 different constants (`StatusPending`, `StatusInProgress`, `StatusCompleted`, `StatusBlocked`, `StatusSuspended`, `StatusSkipped`, `StatusRejected`). The constant `StatusInProgress` has value `"in_progress"` (underscore), but the journey's grep in full-verification uses `"in-progress"` (hyphen). Contradiction.

### surface-type-migration/journey.md

**Problem -> Solution**: Sound. SurfaceType migration follows the same pattern.

**Solution -> Evidence**: Steps reference `detect_surface.go`, `execution_order.go`, `detect.go` which actually exist in `pkg/forgeconfig/`.

**Evidence -> Success Criteria**: Reasonable. Claims "~97 occurrences across 6 files" -- the actual codebase shows ~91 typed constant uses across 4 files in `forgeconfig`, with string literals in more files. The count is approximate and unverifiable from the journey alone.

**Self-contradiction check**: No internal contradictions. Step 3 mentions `execution_order.go`, `detect.go` -- both exist.

### validation-map-consolidation/journey.md

**Problem -> Solution**: Sound. Consolidating hardcoded validation maps to use helper functions.

**Solution -> Evidence**: References `validate_index.go` -- but the actual file is at `internal/cmd/task/validate_index.go`, NOT `internal/cmd/validate_index.go`. Also, the actual code **already uses** `types.AllStatuses()` and `types.AllPriorities()` via `buildValidStatusMap()`/`buildValidValidPriorityMap()` -- so this journey describes work that is already done, not a future migration.

**Evidence -> Success Criteria**: Reasonable but the setup precondition states "All Status and Priority magic values already migrated" which is false per the current codebase (many test files still use string literals).

**Self-contradiction check**: Step 1 describes hardcoded `validStatus map[string]bool{"pending": true, ...}` but the actual code uses `buildValidStatusMap()` which calls `types.AllStatuses()`. The journey describes a past state that no longer exists.

### full-verification/journey.md

**Problem -> Solution**: Sound. Final verification of all migrations.

**Solution -> Evidence**: The grep command in Step 1 searches for `"pending"\|"in-progress"\|"completed"\|"blocked"\|"cancelled"\|"failed"\|"review"` -- this has multiple factual errors:
1. The actual `in-progress` value is `"in_progress"` (underscore), not `"in-progress"` (hyphen)
2. `"cancelled"`, `"failed"`, `"review"` are NOT valid Status constants -- these don't exist in `pkg/types/status.go`
3. The grep misses `"suspended"`, `"skipped"`, `"rejected"` which ARE valid constants

**Evidence -> Success Criteria**: The verification grep is fundamentally broken -- it would miss real magic values (suspended, skipped, rejected) while searching for non-existent ones (cancelled, failed, review).

**Self-contradiction check**: Risk level is "Low" but the journey verifies irreversible codebase changes. If verification fails and a magic value is missed, downstream behavior could break. Low risk seems inconsistent with the purpose.

---

## Phase 2: Rubric Scoring

### Dimension 1: Completeness (200 pts)

**Journey metadata (0-50)**:
- All 4 journeys have `feature`, `journey` (kebab-case), `risk_level` (High/Medium/Low), `surface_types`, `sources`, `generated` date. **Full marks.**

**Steps with required fields (0-80)**:
- All journeys have Steps with "User Action" and "Expected Result". However, the edge case steps use "Precondition" + "User Action" + "Expected Result" which is good. But none of the steps have a structured `step name` header separate from the narrative -- the step names ARE the headers, which is acceptable.
- Deduction: Some steps describe developer actions that are not testable by a downstream agent in a CLI context (e.g., "Developer creates `pkg/types/status.go`" -- this is a code editing action, not a CLI invocation). Steps are implementation actions, not CLI user workflows. **-20 pts**.

**Outcomes cover happy path plus required derived scenarios (0-70)**:
- Happy path covered. Edge cases present but incomplete:
  - status-migration: No "not-found" or "already-exists" derived outcome
  - surface-type-migration: Has "unknown value" (similar to not-found) but no "already-exists"
  - validation-map-consolidation: No "not-found" or "already-exists"
  - full-verification: No "not-found" or "already-exists"
- For a refactoring/migration journey, "not-found" and "already-exists" are less critical, but the CLI surface rule mandates them. **-25 pts**.

**Score: 125/200**

---

### Dimension 2: Semantic Purity (200 pts)

**Outcome descriptions use natural language (0-80)**:
- Outcomes are well-written in natural language. No regex or code in expected results except for grep commands in full-verification Step 1. The grep pattern `"pending"\|"in-progress"\|"completed"\|"blocked"\|"cancelled"\|"failed"\|"review"` is code, not natural language. **-10 pts**.
- Some expected results reference specific Go compiler output ("Exit code non-zero", "Error message clearly identifies the file and line") which is acceptable natural language.

**Preconditions are declarative (0-60)**:
- Edge case preconditions use "Precondition:" prefix and are declarative. Good. E.g., "A Status value exists in the codebase but is not defined as a constant" -- declarative, not procedural. **Full marks.**

**No implementation coupling in Step descriptions (0-60)**:
- Steps are heavily coupled to specific file paths and Go package structure. Step descriptions reference `statemachine.go`, `add.go`, `state.go`, `deps.go`, `build.go`, `autogen.go`, `types.go` by name. This is tight implementation coupling -- if any file is renamed or restructured, the journey breaks. **-25 pts**.
- The step descriptions describe code editing actions ("Developer replaces all Status string literals in `statemachine.go`"), not observable CLI behavior. This is implementation-level detail, not user-facing behavior. **-15 pts**.

**Score: 110/200**

---

### Dimension 3: Precondition Exclusivity (150 pts)

**Preconditions distinct across Outcomes (0-60)**:
- Within each step, the happy path has no precondition and edge cases have distinct preconditions. Good separation. However, status-migration Step 3b has "A string literal was incorrectly replaced with the wrong typed constant" while Step 1b has "A Status value exists in the codebase but is not defined as a constant" -- these are distinct. **Full marks.**

**Preconditions sufficient to uniquely select an Outcome (0-50)**:
- The happy path outcomes have no preconditions, making them the default. Edge cases have explicit preconditions. This is sufficient for selection. However, some edge cases like Step 5c (status-migration) "Developer writes `task.Status == "pending"`" is both a precondition AND a user action, blurring the boundary. **-10 pts**.

**No missing Preconditions for error/boundary Outcomes (0-40)**:
- surface-type-migration Step 3b: "A config file specifies `surface: "desktop"` which is not a defined constant" -- this precondition assumes "desktop" is invalid, but doesn't establish that "desktop" is NOT one of the existing constants. Weak but acceptable.
- validation-map-consolidation Step 2b: "`AllStatuses()` returns an empty slice due to a bug" -- this precondition describes an internal bug, not an external state. It's testing the helper function's implementation, not the validation map's behavior. **-10 pts**.

**Score: 130/150**

---

### Dimension 4: Fact Alignment (150 pts)

**Factual claims traceable to fact_id or marked UNKNOWN (0-60)**:
- No fact_ids are used anywhere. No UNKNOWN markers. All claims are asserted without traceability. The journeys have `sources` in frontmatter pointing to PRD docs, but individual claims are not linked. **-30 pts**.

**Inferred claims have required_outcomes rule support (0-50)**:
- Many claims are inferred (e.g., "117 Status string literals across 22 files", "~97 Surface Type string literals across 6 files") but none are marked as inferred or have source classification. **-25 pts**.

**No hallucinated unclassified claims (0-40)**:
- **CRITICAL**: status-migration Step 1 lists constants `StatusCancelled`, `StatusFailed`, `StatusReview` which DO NOT EXIST in the actual codebase. The real constants are `StatusSuspended`, `StatusSkipped`, `StatusRejected`. This is a **hallucinated claim**. **-40 pts (per rule: -30 per instance = -90, capped at -40 for this subdimension)**.
- full-verification Step 1 grep searches for `"cancelled"`, `"failed"`, `"review"` which don't exist as status values. Also misses `"suspended"`, `"skipped"`, `"rejected"` which DO exist. **Another hallucinated claim.**
- full-verification Step 1 grep searches for `"in-progress"` (hyphen) but the actual value is `"in_progress"` (underscore). **Factual error.**
- status-migration Setup claims "117 Status string literals across 22 files" -- the actual codebase has far more occurrences (240+ in pkg/ alone, 843 in internal/) across more than 22 files. **Inaccurate claim.**
- validation-map-consolidation references `internal/cmd/validate_index.go` but the actual path is `internal/cmd/task/validate_index.go`. **Factual error.**
- validation-map-consolidation Step 1 describes hardcoded `validStatus map[string]bool{"pending": true, ...}` but the actual code uses `buildValidStatusMap()` calling `types.AllStatuses()`. **Describes a state that doesn't exist.**

**Score: 5/150** (Multiple hallucinated claims about non-existent constants, wrong file paths, wrong values)

---

### Dimension 5: Surface Fitness (150 pts)

**Mandatory derived Outcomes (not-found + already-exists) (0-60)**:
- None of the 4 journeys include "not-found" or "already-exists" derived outcomes. The CLI surface rule explicitly mandates these: "Mandatory derived Outcomes (must be considered for every CLI Journey): not-found, already-exists". **0 pts.**
- This is a **surface type violation** per the rubric rules: -25 pts per instance x 4 journeys = -100 pts. Capped at 0 for this subdimension.

**Test strategy proportions (Contract 80% / Journey 20%) (0-50)**:
- The journeys don't mention test strategy proportions or contract vs. journey test split. No explicit acknowledgment of the CLI guidance. **-25 pts**.

**CLI-specific environment assumptions (0-40)**:
- The journeys reference `go build`, `go test`, `grep` commands which are subprocess-invokable. But the steps describe "Developer creates" and "Developer replaces" which are code-editing actions, NOT CLI command invocations. The journeys are written as developer code-editing workflows, not as CLI subprocess invocations. This contradicts the CLI surface model which requires subprocess isolation and exit code verification. **-30 pts**.
- The only steps that resemble CLI testing are in full-verification (grep, go build, go test, forge task list). Even these are described as manual developer actions, not as automated subprocess tests.

**Score: 10/150**

---

### Dimension 6: Internal Consistency (150 pts)

**Invariants hold in every Step (0-60)**:
- status-migration Invariant: "All typed constant values must exactly match the original string literals (zero behavior change)". But Step 1 lists `StatusCancelled`, `StatusFailed`, `StatusReview` as constants -- if these were created, they would NOT match any original string literals since they don't exist in the codebase. **Invariant violation.** **-40 pts.**
- full-verification Invariant: "Zero Status/SurfaceType/Priority string literals in production code". But the verification grep in Step 1 searches for wrong values, meaning this invariant cannot be properly verified with the given approach. **-10 pts.**

**Cross-Step references consistent (0-50)**:
- status-migration: Step 1 defines 7 constants. Step 3 replaces literals in statemachine.go. Step 4 replaces in "remaining files". The progression is logical. But Step 1 lists wrong constants, which propagates inconsistency to all subsequent steps. **-20 pts.**
- full-verification Step 4b references `assert.Equal(t, "pending", task.Status)` but the actual constant value is `"in_progress"` for `StatusInProgress` -- the invariant about test assertions is reasonable but the specific example uses a value that may not trigger the described failure (untyped `"pending"` IS assignable to `types.Status`). **-10 pts.**

**Risk level consistent with Journey content (0-40)**:
- status-migration: High risk. Content involves state machine mutation and type system changes. **Consistent.**
- surface-type-migration: High risk. ~97 occurrences across 6 files. **Consistent.**
- validation-map-consolidation: Medium risk. Multi-step without irreversible side effects. **Consistent.**
- full-verification: Low risk. Read-only verification. **Consistent.**
- **Full marks.**

**Score: 60/150**

---

## Phase 3: Blindspot Hunt

### [blindspot-1] Wrong Status Constant Names — Critical
status-migration Step 1 lists `StatusCancelled`, `StatusFailed`, `StatusReview`. The actual constants are `StatusSuspended`, `StatusSkipped`, `StatusRejected`. This error propagates to full-verification's grep command which also searches for non-existent values and misses real ones. A downstream agent following these journeys would create wrong constants and fail to detect real magic values.

### [blindspot-2] Wrong String Value for InProgress
The actual `StatusInProgress` constant has value `"in_progress"` (underscore), but full-verification Step 1 greps for `"in-progress"` (hyphen). This means the verification step would never catch `in-progress` magic values if they existed, and would incorrectly search for a non-existent variant.

### [blindspot-3] Validation Map Already Migrated
validation-map-consolidation describes replacing hardcoded `validStatus map[string]bool{"pending": true, ...}` with `types.AllStatuses()`. But the actual `validate_index.go` ALREADY uses `buildValidStatusMap()` calling `types.AllStatuses()`. This journey describes completed work, not pending work.

### [blindspot-4] Wrong File Path
validation-map-consolidation references `internal/cmd/validate_index.go` but the actual path is `internal/cmd/task/validate_index.go`. A downstream agent would fail to find the file.

### [blindspot-5] Journeys Describe Code Editing, Not CLI Behavior
All 4 journeys describe developer code-editing actions ("Developer creates", "Developer replaces", "Developer adds") rather than CLI command invocations. The CLI surface type requires subprocess-isolated tests with exit code verification. These journeys would be more appropriate for a "code" or "refactoring" surface type.

### [blindspot-6] No Error Path Coverage for CLI Commands
No journey tests what happens when CLI commands fail. For example: what if `forge task validate-index` is run on a malformed JSON file? What if `go build` fails due to a network issue? The CLI surface requires exit code and stderr verification for error paths.

---

## Per-Journey Dimension Breakdown

| Dimension | status-migration | surface-type-migration | validation-map-consolidation | full-verification |
|-----------|:-:|:-:|:-:|:-:|
| Completeness | 30/50 | 35/50 | 30/50 | 30/50 |
| Steps complete | 60/80 | 65/80 | 60/80 | 55/80 |
| Outcomes | 45/70 | 45/70 | 45/70 | 40/70 |
| **Completeness Total** | **135/200** | **145/200** | **135/200** | **125/200** |
| Semantic Purity | 100/200 | 130/200 | 115/200 | 100/200 |
| Precondition Exclusivity | 135/150 | 140/150 | 120/150 | 125/150 |
| Fact Alignment | 5/150 | 50/150 | 15/150 | 5/150 |
| Surface Fitness | 10/150 | 15/150 | 10/150 | 10/150 |
| Internal Consistency | 55/150 | 70/150 | 60/150 | 55/150 |

## Overall Averaged Scores

| Dimension | Average Score |
|-----------|:---:|
| 1. Completeness | **135/200** |
| 2. Semantic Purity | **111/200** |
| 3. Precondition Exclusivity | **130/150** |
| 4. Fact Alignment | **19/150** |
| 5. Surface Fitness | **11/150** |
| 6. Internal Consistency | **60/150** |
| **Total** | **466/1000** |

---

## Attack Points Summary

1. **Fact Alignment**: Hallucinated constants `StatusCancelled`, `StatusFailed`, `StatusReview` in status-migration Step 1 -- "7 constants (`StatusPending`, `StatusInProgress`, `StatusCompleted`, `StatusBlocked`, `StatusCancelled`, `StatusFailed`, `StatusReview`)" -- must be corrected to `StatusSuspended`, `StatusSkipped`, `StatusRejected`.

2. **Fact Alignment**: Wrong string value in full-verification grep -- `"in-progress"` (hyphen) vs actual `"in_progress"` (underscore) -- must use correct value.

3. **Fact Alignment**: Missing real constants in full-verification grep -- `"suspended"`, `"skipped"`, `"rejected"` are not searched, meaning real magic values would be missed -- must include all 7 actual status values.

4. **Surface Fitness**: Missing mandatory CLI derived outcomes (`not-found`, `already-exists`) in all 4 journeys -- "Mandatory derived Outcomes (must be considered for every CLI Journey)" from surface-cli.md -- must add these outcomes to each journey.

5. **Surface Fitness**: Steps describe code-editing actions, not CLI subprocess invocations -- "Developer creates `pkg/types/status.go`" -- must reframe as CLI commands with exit code and output assertions.

6. **Fact Alignment**: Wrong file path in validation-map-consolidation -- "internal/cmd/validate_index.go" -- must be corrected to `internal/cmd/task/validate_index.go`.

7. **Fact Alignment**: validation-map-consolidation describes already-completed work -- "validStatus map[string]bool{\"pending\": true, ...}" -- must reflect the actual current state where `buildValidStatusMap()` already calls `types.AllStatuses()`.

8. **Semantic Purity**: Heavy implementation coupling in step descriptions -- "Developer replaces all Status string literals in `statemachine.go`" -- must describe observable behavior rather than implementation details.

9. **Internal Consistency**: Invariant violation -- "All typed constant values must exactly match the original string literals" but Step 1 proposes constants that don't match any originals -- must align constants with actual codebase values.

10. **Completeness**: Inaccurate occurrence counts -- "117 Status string literals across 22 files" is unverifiable and likely outdated -- must mark as approximate or update with actual counts.
