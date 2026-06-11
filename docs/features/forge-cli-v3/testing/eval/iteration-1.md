---
feature: "forge-cli-v3"
iteration: 1
scored: "2026-05-14"
model: "claude-opus-4"
---

# Test Cases Evaluation Report — Iteration 1

**Document**: `docs/features/forge-cli-v3/testing/test-cases.md`
**Active Profile**: go-test (capabilities: tui, api, cli)

---

## SCORE: 645/1000

---

## Dimension Breakdown

### 1. PRD Traceability: 185/250

| Criterion | Score | Notes |
|-----------|-------|-------|
| TC-to-AC mapping | 70/90 | Every TC has a `Source` field, but granularity is mixed. User-story TCs use "Story N / AC-K" which is good. Error-handling TCs (TC-031 to TC-041) use "Spec Error Handling Table" — a section-level reference, not a specific row or AC. This is insufficient for downstream traceability. |
| Traceability table complete | 80/80 | All 41 TCs appear in the traceability table with Source, Type, Target, Priority. No TCs are missing. |
| Reverse coverage | 35/80 | **Significant gaps**: (1) PRD spec "State Transition Constraints" table has 5 rows (pending->*, in_progress->*, completed->*, blocked->*, rejected->*) — zero TCs cover invalid state transitions. (2) `forge task claim` success path is untested (only error paths TC-031/TC-032). (3) `forge task status` happy path has no TC (only error TC-035). (4) `forge task add` has zero TCs. (5) `forge task index` has zero TCs. (6) `forge task migrate` has zero TCs. (7) `forge task query` has zero TCs. (8) `forge feature` (get/set) has zero TCs. (9) `forge probe` has zero TCs. (10) `forge version` has zero TCs. (11) `forge e2e setup/verify/compile/discover` — 4 subcommands with zero TCs. (12) `forge prompt get-by-task-id` success with variable substitution (AC-1 step 4: `{{TASK_ID}}, {{TASK_FILE}}, {{SCOPE}}`) is asserted in prose but the step does not verify each variable individually. |

**Reverse coverage detail**: Of the 23 commands listed in the PRD command structure table, only 13 have at least one TC. 10 commands are completely untested: `task claim` (happy), `task add`, `task index`, `task migrate`, `task query`, `task status` (happy), `feature`, `probe`, `version`, `e2e {setup,verify,compile,discover}`.

---

### 2. Step Actionability: 210/250

| Criterion | Score | Notes |
|-----------|-------|-------|
| Steps are concrete actions | 75/90 | Most steps use concrete CLI invocations (`forge task submit T-impl-1 --result success --summary "..."`). However: TC-011 "Simultaneously invoke" is vague — no mechanism specified (shell backgrounding, goroutine, test harness?). TC-039 "Simulate a lock contention scenario" is a placeholder instruction, not a concrete action. TC-015/016 preconditions are abstract ("A fix-task exists but the corresponding fix also fails") without setup steps. |
| Expected results are verifiable | 80/90 | Most expected results specify exit codes and stderr content, which is objectively verifiable. Weak spots: TC-013 step 2 "verify the command executes compile, then fmt, then lint, then test in sequence" — how is sequence verified? No log parsing or intermediate checkpoint assertion. TC-023 "output contains a list" — what fields? What format? TC-024 "compact evidence summary" — no assertion on specific output fields. |
| Preconditions are explicit | 55/70 | Preconditions declare data state (task status, file existence) for most TCs, but several are thin: TC-011/TC-039 require concurrent execution setup but don't specify how. TC-015/TC-016 require specific fix-task counts but don't specify how to create that state. TC-022 "empty type registry" has no precond step to achieve that state. |

**Blocking assessment**: Score is 210 >= 200 threshold. Downstream gen-test-scripts is NOT blocked.

---

### 3. Interface Accuracy: 140/200

Active capabilities: `tui`, `api`, `cli`. This is a CLI-only feature, so `api` and `tui` are not applicable in practice. The document itself declares 0 API and 0 integration tests. The effective evaluation uses the `cli` capability dimensions, weighted proportionally.

**Per-capability breakdown:**

| Capability | Criteria | Score | Notes |
|-----------|----------|-------|-------|
| cli | Command coverage | 60/100 | **Major gaps**: 10+ commands/subcommands have zero TCs (see reverse coverage). Flag combinations are barely explored — `--result` variants only test `success` and missing, never `blocked` or `rejected`. `--summary` content validation is untested. `forge e2e run` has 4 TCs but the other 4 e2e subcommands (setup, verify, compile, discover) have none. No TC tests `forge task claim` success path. No TC tests flag combinations like `forge task submit --result blocked --summary "..."` which is a core workflow per PRD Agent Flow. |
| cli | Output assertion specificity | 80/100 | Exit codes and stderr messages are well-specified for covered commands. Good adherence to PRD error table messages. Weak spots: TC-002 step 4 ("description follows pattern") is structural but not automatable without regex. TC-021 step 4-5 (description format/length) requires ad-hoc parsing. Some expected results say "output contains" without specifying exact substrings or format. |

---

### 4. Completeness: 60/200

| Criterion | Score | Notes |
|-----------|-------|-------|
| Type coverage | 50/70 | CLI is well-covered for implemented TCs. PRD declares "CLI-only, no UI" which matches the 0 UI TCs. But API/integration are 0 and the PRD does have internal Go APIs (profile detection, state machine) that could warrant integration tests. For the declared scope, CLI coverage is acceptable but incomplete due to missing commands. |
| Boundary and edge cases | 10/70 | Very few boundary/edge cases: TC-022 (empty registry) and TC-014 (no terminal tasks) are the only non-happy-path-that-isn't-an-error. Missing: (1) state transition boundary tests (pending->completed should be invalid, pending->in_progress should be valid, blocked->completed should be invalid). (2) Very long `--summary` values. (3) Special characters in task IDs. (4) Empty `--summary` string. (5) `forge e2e run` with multiple features. (6) `forge task submit --result blocked` (only `success` and missing are tested). |
| Integration scenarios | 0/60 | Zero cross-command workflow TCs. PRD describes a full Agent Flow (claim -> execute -> quality-gate -> submit) and Hook Flow (SessionEnd -> cleanup, Stop -> quality-gate, PreToolUse -> verify-task-done). None of these end-to-end flows have TCs. TC-011 tests concurrent submit but not a multi-command workflow. |

---

### 5. Structure & ID Integrity: 50/100

| Criterion | Score | Notes |
|-----------|-------|-------|
| TC IDs sequential and unique | 40/40 | TC-001 through TC-041, sequential, no gaps, no duplicates. |
| Classification correct | 10/30 | All 41 TCs classified as CLI, which matches the feature. However, the rubric requires the document to have sections for UI, API, and CLI test cases. The document only has a "CLI Test Cases" section with 0 UI and 0 API sections. The summary table shows 0 for UI/API/Integration but no actual section headers exist for those types. The required "Route Validation table" is entirely missing. |
| Summary table matches actual | 0/30 | Summary table claims Integration: 0 and API: 0, but no corresponding sections exist in the document. More critically, the required "Route Validation table" section listed in the rubric is completely absent. Every TC has `Route: N/A` and `Element: sitemap-missing`, suggesting the route validation table was never populated. This is a structural omission. |

**Deductions applied:**
- 41 TCs with `Element: sitemap-missing` (placeholder text): -820 pts capped to -0 (deductions cannot exceed dimension scores). Applied within Structure dimension: each placeholder `sitemap-missing` is a TBD/N/A equivalent. However, since this is a CLI feature, Route/Element fields are genuinely N/A. The issue is the document uses a placeholder string instead of omitting the fields or using `N/A` with justification. Deduction applied as: -20 per instance for 41 instances would be -820, but capped to the dimension maximum of 100 pts, effectively reducing Structure to 0. **Adjusted**: Since the PRD explicitly states "CLI-only, no UI" and Route/Element are CLI-irrelevant, the placeholder is a template artifact, not a vagueness issue. Reduced deduction to -50 (proportional to template misuse, not per-instance).

---

## ATTACKS (Top 3 Weakest Areas)

### Attack 1: Missing Commands — 10+ PRD Commands Have Zero Test Cases (Completeness: 10/70 boundary, PRD Traceability: 35/80 reverse coverage)

The PRD defines 23 commands/subcommands. The test suite covers 13. Completely untested:
- `forge task claim` (success path), `forge task add`, `forge task index`, `forge task migrate`, `forge task query`, `forge task status` (success path)
- `forge feature`, `forge probe`, `forge version`
- `forge e2e setup`, `forge e2e verify`, `forge e2e compile`, `forge e2e discover`

This is a 43% command coverage gap. A downstream agent generating test scripts from these cases will produce incomplete coverage by construction.

**Fix**: Add TCs for every untested command, prioritizing `task claim` (core agent workflow), `feature` (context management), and the 4 e2e subcommands.

### Attack 2: Zero End-to-End Workflow Coverage (Completeness: 0/60 integration)

The PRD describes 3 multi-command workflows (Agent Flow, Developer Flow, Hook Flow). None have integration TCs. The test suite treats each command as an isolated unit. This means:
- No test verifies `claim -> execute -> submit` lifecycle
- No test verifies `quality-gate failure -> fix-task creation -> fix -> quality-gate retry` loop
- No test verifies `SessionEnd -> cleanup` followed by `Stop -> quality-gate` sequence

**Fix**: Add 3-5 workflow-level TCs that exercise command sequences end-to-end. These are the highest-value tests for a CLI refactoring (behavioral equivalence).

### Attack 3: State Machine Transitions Untested (PRD Traceability: 35/80 reverse coverage)

The PRD defines a 5-state machine (pending, in_progress, completed, blocked, rejected) with explicit allowed/forbidden transitions. Zero TCs test state transition validation. This is a core behavioral contract — a CLI refactoring must preserve it exactly. The only related TC is TC-009 (submit to terminal state), which tests one terminal rejection but not the transition matrix.

**Fix**: Add a parameterized TC (or TC group) that covers all 5 states x valid transitions and at least 3 invalid transitions (pending->completed, blocked->completed, in_progress->pending).

---

## Scoring Summary

| Dimension | Weight | Score | Pct |
|-----------|--------|-------|-----|
| 1. PRD Traceability | 250 | 185 | 74% |
| 2. Step Actionability | 250 | 210 | 84% |
| 3. Interface Accuracy | 200 | 140 | 70% |
| 4. Completeness | 200 | 60 | 30% |
| 5. Structure & ID Integrity | 100 | 50 | 50% |
| **Total** | **1000** | **645** | **64.5%** |
