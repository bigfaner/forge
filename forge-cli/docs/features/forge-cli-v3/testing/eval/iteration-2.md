---
date: "2026-05-14"
doc_dir: docs/features/forge-cli-v3/testing/
iteration: 2
target_score: 800
evaluator: Claude (automated, adversarial)
---

# Test Cases Eval — Iteration 2

**Score: 845/1000** (target: 800)

```
┌─────────────────────────────────────────────────────────────────┐
│                  TEST CASES QUALITY SCORECARD                     │
├──────────────────────────────────────────────────────────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 1. PRD Traceability          │  225     │  250     │ ✅         │
│    TC-to-AC mapping          │  80/90   │          │            │
│    Traceability table        │  80/80   │          │            │
│    Reverse coverage          │  65/80   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 2. Step Actionability        │  225     │  250     │ ✅         │
│    Steps concrete            │  80/90   │          │            │
│    Expected results          │  85/90   │          │            │
│    Preconditions explicit    │  60/70   │          │            │
├──────────────────────────────┼──────────┼──────────┼──────────┤
│ 3. Interface Accuracy        │  170     │  200     │ ✅         │
│    Command coverage          │  85/100  │          │            │
│    Output assertion specific │  85/100  │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 4. Completeness              │  125     │  200     │ ⚠️         │
│    Type coverage             │  60/70   │          │            │
│    Boundary cases            │  25/70   │          │            │
│    Integration scenarios     │  40/60   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 5. Structure & ID Integrity  │  100     │  100     │ ✅         │
│    IDs sequential/unique     │  40/40   │          │            │
│    Classification correct    │  30/30   │          │            │
│    Summary matches actual    │  30/30   │          │            │
├──────────────────────────────┼──────────┼──────────┼──────────┼──────────┤
│ TOTAL                        │  845     │  1000    │            │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

---

## Dimension Details

### 1. PRD Traceability: 225/250

| Criterion | Score | Notes |
|-----------|-------|-------|
| TC-to-AC mapping | 80/90 | All 71 TCs have Source fields. User-story TCs use "Story N / AC-K" format. State transition TCs use "Spec State Transition Constraints Table (row)". New workflow TCs reference "Spec Agent Task Execution Flow" and "Spec CI/Hook Flow" with sufficient specificity. **Remaining issue**: TC-031 through TC-041 still use "Spec Error Handling Table" as a section-level reference rather than pointing to a specific row. The error handling table has 22 rows with distinct (command, scenario) pairs, and "Spec Error Handling Table" alone is ambiguous. Deduction: -10. |
| Traceability table complete | 80/80 | All 71 TCs appear in the traceability table with Source, Type, Target, and Priority. No missing entries. The table grew from 41 to 71 rows correctly. |
| Reverse coverage | 65/80 | Major improvement from iteration 1 (35/80). All 23 PRD commands now have at least one TC. State transition matrix is fully covered (5 states, all valid transitions, 3 invalid transitions, 2 terminal-state rejections). All 8 user stories have mapped TCs. **Remaining gaps**: (1) PRD "Naming Change Table" has 10 rows mapping old-to-new commands — no TC explicitly verifies renaming equivalence (e.g., that `forge task submit` behaves identically to the old `task record`). This is a behavioral equivalence contract that a refactoring must validate. (2) PRD "Related Changes" table lists 5 change points (hooks.json, justfile, skills, docs, Go tests) — no TC verifies these reference updates are complete. (3) `forge task check-deps` has only an error-path TC (TC-033) — no happy-path TC where dependencies are met and the command succeeds. (4) `forge task validate-index` has only an error-path TC (TC-034) — no happy-path TC. Deduction: -15. |

---

### 2. Step Actionability: 225/250

| Criterion | Score | Notes |
|-----------|-------|-------|
| Steps are concrete actions | 80/90 | Most steps use concrete CLI invocations with specific flags and arguments. New TCs (TC-042 through TC-055) use `jq` commands to verify index.json state, which is automatable. **Remaining issues**: (1) TC-011 "Simultaneously invoke ... from two processes" — still vague on mechanism (shell backgrounding with `&`? Go test goroutine? expect script?). (2) TC-039 "Simulate a lock contention scenario" — still a placeholder instruction, not a concrete action. (3) TC-015/TC-016 preconditions are abstract ("A fix-task exists but the corresponding fix also fails") without concrete setup steps (which file to corrupt? which compile error to inject?). (4) TC-050 step 4 "Stop the service" — no concrete mechanism specified. Deduction: -10. |
| Expected results are verifiable | 85/90 | Good improvement. New TCs specify `jq` queries for state verification (e.g., TC-042 step 4, TC-044 step 4-5). Exit codes and stderr content are well-specified. **Remaining issues**: (1) TC-002 step 4 "verify description follows pattern" — structural assertion without regex. (2) TC-013 step 2 "verify the command executes compile, then fmt, then lint, then test in sequence" — sequence verification mechanism unspecified. (3) TC-024 "compact evidence summary" — no assertion on specific output fields. (4) TC-025 "summary information" — unspecified format. Deduction: -5. |
| Preconditions explicit | 60/70 | Improvement from iteration 1. New TCs declare preconditions more clearly (e.g., TC-057 through TC-066 specify exact initial status for each task). TC-067 through TC-071 integration TCs declare feature names and task IDs. **Remaining issues**: (1) TC-011/TC-039 require concurrent execution setup but don't specify how to achieve it. (2) TC-015/TC-016 require fix-task setup but don't specify how to create that state. (3) TC-022 "empty type registry" has no precondition step to achieve that state (delete all type definitions?). (4) TC-045 "older schema version format" — no guidance on what the old format looks like or how to create a fixture. Deduction: -10. |

**Blocking assessment**: Score is 225 >= 200 threshold. Downstream gen-test-scripts is NOT blocked.

---

### 3. Interface Accuracy: 170/200

Active capabilities: `tui`, `api`, `cli`. This is a CLI-only feature, so `api` and `tui` are not applicable in practice. The effective evaluation uses the `cli` capability dimensions, weighted proportionally (200 pts allocated to cli command coverage and output assertion specificity).

| Criterion | Score | Notes |
|-----------|-------|-------|
| Command coverage | 85/100 | Major improvement. All 23 PRD commands now have at least one TC. `--result` variants now cover `success`, `blocked`, and `rejected`. State transition matrix is fully covered with 10 TCs. **Remaining gaps**: (1) `forge task submit --result rejected` is covered by TC-058 but only in the state transition section — there is no standalone `--result rejected` TC verifying record creation and index.json update independently. (2) `forge task check-deps` and `forge task validate-index` have only error-path TCs — no happy paths. (3) `forge task query` TC-046 tests `--status pending` and `--status completed` but not other potential flags (e.g., `--type`, `--feature`). The PRD spec lists `query` as a subcommand but does not define its flags — the TC covers the reasonable assumption. (4) No TC tests `forge task add` with missing required flags (e.g., omitting `--title` or `--type`). (5) No TC tests `forge e2e run` with multiple features or runs without `--feature` flag. Deduction: -15. |
| Output assertion specificity | 85/100 | Exit codes and stderr messages are well-specified for covered commands. New TCs use `jq` for structured JSON verification. **Remaining issues**: (1) TC-048 "stdout contains the current feature name/slug" — no specific expected value. (2) TC-050 "stdout contains HTTP status code (e.g., '200 OK')" — the "e.g." is a suggestion, not an assertion. (3) TC-052/053/054/055 output assertions are thin ("profile name displayed", "verification pass message", "compilation success message", "stdout lists all test files") — no exact strings or format specified. (4) TC-069 step 6 is a template `<id>` that would need to be resolved at test generation time — acceptable but noted. (5) TC-071 step 1 output "fix-compile-1" is an exact assertion — good pattern. Deduction: -15. |

---

### 4. Completeness: 125/200

| Criterion | Score | Notes |
|-----------|-------|-------|
| Type coverage | 60/70 | CLI is well-covered with 66 TCs. Integration has 5 TCs covering all 3 PRD flows (Agent, Developer, Hook). PRD declares "CLI-only, no UI" so 0 UI TCs is correct. **Remaining gap**: The PRD spec describes internal Go APIs (profile detection, state machine transitions, index.json schema validation) that could warrant API-level unit tests. However, since the PRD scope is CLI refactoring and the test profile is `go-test`, the CLI + Integration classification is acceptable for the declared scope. Deduction: -10 for the internal API gap. |
| Boundary and edge cases | 25/70 | Some improvement. TC-022 (empty registry), TC-014 (no terminal tasks), TC-036 (no search results), TC-062/063/064 (invalid state transitions) are boundary cases. **Significant remaining gaps**: (1) No TC for very long `--summary` values (max length boundary). (2) No TC for special characters in task IDs (spaces, unicode, SQL injection patterns). (3) No TC for empty `--summary` string. (4) No TC for `--result` with invalid value (e.g., `--result invalid_state`). (5) No TC for `forge task add` with duplicate title. (6) No TC for `forge e2e run` with multiple features simultaneously. (7) No TC for `forge task index` when no task markdown files exist (empty directory). (8) No TC for `forge quality-gate` when individual steps have mixed results (compile passes, lint fails). (9) No TC for concurrent `forge task claim` by two agents on the same task. Deduction: -45. |
| Integration scenarios | 40/60 | Major improvement from 0/60 in iteration 1. Five integration TCs now cover: (1) Agent Flow claim-execute-submit lifecycle (TC-067). (2) Agent Flow quality gate failure with fix-task escalation (TC-068). (3) Hook Flow cleanup-then-quality-gate-then-verify sequence (TC-069). (4) Developer Flow profile-detect-set-e2e-run (TC-070). (5) Quality gate retry loop with max-3 fix-tasks (TC-071). **Remaining gaps**: (1) No integration TC for the "quality-gate failure -> fix-task creation -> fix execution -> quality-gate retry -> success" full recovery loop. TC-068 stops at "blocked" escalation but doesn't test the recovery branch where the fix succeeds. (2) No integration TC for the PreToolUse hook (verify-task-done) as a gate in the Agent Flow. TC-069 tests it standalone but not as part of the agent execution flow. (3) No integration TC testing the `forge cleanup` after a quality-gate failure + fix-task cycle (cleanup should remove fix-task state files after fix completes). Deduction: -20. |

---

### 5. Structure & ID Integrity: 100/100

| Criterion | Score | Notes |
|-----------|-------|-------|
| TC IDs sequential and unique | 40/40 | TC-001 through TC-071, sequential, no gaps, no duplicates. |
| Classification correct | 30/30 | 66 TCs classified as CLI, 5 as Integration. No TCs in wrong sections. Document has sections for CLI Test Cases and Integration Test Cases with a note explaining 0 UI/API TCs. All classifications match the PRD scope. |
| Summary matches actual | 30/30 | Summary table shows: UI=0, Integration=5, API=0, CLI=66, Total=71. Manual count confirms: TC-067 through TC-071 are Integration (5), TC-001 through TC-066 are CLI (66). Match is exact. |

**Deductions applied:**
- TC-001 through TC-041 still have `Element: sitemap-missing` placeholder. TC-042 through TC-071 use `Element: N/A`. For a CLI-only feature, Route/Element fields are genuinely N/A. The `sitemap-missing` placeholder in the first 41 TCs is a template artifact. However, since these TCs are clearly CLI (no UI routes), and the feature is CLI-only per PRD, the practical impact is zero. No deduction applied — the document now uses `N/A` for all new TCs, and the legacy `sitemap-missing` is cosmetic.

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| TC-031 to TC-041 (11 TCs) | Source field "Spec Error Handling Table" is section-level, not row-level | -10 (in Traceability) |
| TC-011, TC-039 | "Simultaneously invoke" / "Simulate lock contention" — vague mechanism | -10 (in Step Actionability) |
| TC-015, TC-016 | Abstract preconditions without concrete setup steps | -5 (in Step Actionability) |
| TC-002, TC-013, TC-024, TC-025 | Expected results lack specific assertion values | -5 (in Step Actionability) |
| TC-022, TC-045 | Preconditions cannot be achieved without guidance | -5 (in Step Actionability) |
| TC-033, TC-034 | Error-only TCs; no happy-path counterparts for check-deps and validate-index | -15 (in Interface Accuracy) |
| TC-048, TC-050, TC-052-055 | Output assertions use vague descriptors instead of exact strings | -15 (in Interface Accuracy) |
| Boundary cases | Missing 9+ boundary/edge case categories | -45 (in Completeness) |
| Integration scenarios | Missing recovery loop, PreToolUse-in-Agent-Flow, post-fix cleanup | -20 (in Completeness) |
| Internal APIs | No TCs for internal Go API boundaries (profile detection, state machine) | -10 (in Completeness) |

---

## Attack Points

### Attack 1: Boundary and Edge Case Coverage is Thin (Completeness: 25/70)

**Where**: The entire test suite has 71 TCs but only ~4 are genuine boundary/edge cases (TC-014, TC-022, TC-036, TC-062/063/064 for invalid transitions).
**Why it's weak**: The PRD error handling table defines 22 error scenarios. The test suite covers most of them. However, it misses critical input boundaries: (1) what happens with an empty `--summary` string? (2) what happens with `--result invalid_value`? (3) what happens with very long task IDs or special characters? (4) what happens when `forge task index` runs on an empty directory? (5) what happens when `forge quality-gate` encounters a partial failure (compile passes but lint fails)? These are not exotic edge cases — they are common operational scenarios for a CLI tool. A refactoring must preserve behavior at these boundaries.
**What must improve**: Add 5-8 boundary TCs: (1) `--result` with invalid value, (2) empty `--summary`, (3) special characters in task ID, (4) `forge task index` with no markdown files, (5) partial quality-gate failure (compile pass, lint fail), (6) `forge task add` with missing flags, (7) `forge e2e run` without `--feature` flag, (8) concurrent `forge task claim` on same task.

### Attack 2: Integration Recovery Loop Untested (Completeness: 40/60)

**Where**: TC-068 tests quality-gate failure -> fix-task creation -> blocked escalation, but stops at the "blocked" terminal. TC-071 tests the retry loop up to max-3. Neither tests the successful recovery path.
**Why it's weak**: The PRD Agent Task Execution Flow explicitly shows a recovery branch: "Quality Gate fails -> Fix -> Quality Gate retries -> Succeeds -> Submit". This is the most common happy-path for the quality gate — most failures get fixed and succeed on retry. The test suite has no TC for this. TC-068 goes straight to "blocked" escalation (failure), and TC-071 goes to "max fix-tasks reached" (catastrophic failure). The successful middle path is untested.
**What must improve**: Add 1 integration TC: "Agent Flow — quality gate failure, fix succeeds, quality gate passes, task submits". This is the primary recovery path and arguably the most important integration test.

### Attack 3: Error-Only Commands Missing Happy Paths (Interface Accuracy: 85/100)

**Where**: `forge task check-deps` has only TC-033 (unmet dependency error). `forge task validate-index` has only TC-034 (invalid schema error). Neither has a happy-path TC.
**Why it's weak**: These commands are invoked in the PRD Agent Task Execution Flow and Hook Flow. A downstream test script generator will produce tests that only verify error behavior for these commands. The happy paths (dependencies met -> success, valid index -> success) are never verified. For a refactoring that must preserve behavioral equivalence, this is a risky gap — the success behavior is the primary contract.
**What must improve**: Add 2 TCs: (1) `forge task check-deps` when all dependencies are completed -> exit code 0 with success message, (2) `forge task validate-index` when index.json is valid -> exit code 0 with validation pass message.

---

## Previous Issues Check

| Previous Attack | Addressed? | Evidence |
|----------------|------------|----------|
| Attack 1: Missing Commands — 10+ PRD commands with zero TCs | ✅ Fully addressed | All 23 PRD commands now have at least one TC. Added TC-042 (task claim happy), TC-043 (task add), TC-044 (task index), TC-045 (task migrate), TC-046 (task query), TC-047 (task status happy), TC-048/TC-049 (feature get/set), TC-050 (probe), TC-051 (version), TC-052-055 (e2e setup/verify/compile/discover). |
| Attack 2: Zero End-to-End Workflow Coverage | ✅ Fully addressed | Added 5 integration TCs: TC-067 (Agent Flow lifecycle), TC-068 (quality gate failure + escalation), TC-069 (Hook Flow sequence), TC-070 (Developer Flow profile setup + e2e), TC-071 (quality gate retry loop). |
| Attack 3: State Machine Transitions Untested | ✅ Fully addressed | Added 10 state transition TCs: TC-057-TC-061 (5 valid transitions), TC-062-TC-064 (3 invalid transitions), TC-065-TC-066 (2 terminal state rejections). Full coverage of the PRD state transition constraints table. |

---

## Verdict

- **Score**: 845/1000
- **Target**: 800/1000
- **Gap**: +45 points above target
- **Step Actionability**: 225/250 (not blocking)
- **Action**: Target reached. The document improved from 645/1000 (iteration 1) to 845/1000 (iteration 2), a +200 point gain. All three iteration-1 attacks were fully addressed. The remaining weaknesses (boundary cases, recovery loop integration, happy-path gaps) are P1/P2 improvements that can be addressed during test script generation without blocking downstream work.
