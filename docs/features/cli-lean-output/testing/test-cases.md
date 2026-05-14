---
feature: "cli-lean-output"
sources:
  - docs/proposals/cli-lean-output/proposal.md
  - docs/features/cli-lean-output/tasks/1-trim-claim-output.md
  - docs/features/cli-lean-output/tasks/2-trim-submit-query-status.md
generated: "2026-05-15"
profile: "go-test"
---

# Test Cases: cli-lean-output

> **WARNING**: sitemap.json not found — Element set to `sitemap-missing`. Run `/gen-sitemap` for precise element references.

## Summary

| Type | Count |
|------|-------|
| UI   | 0    |
| **Integration** | **0** |
| API  | 0   |
| CLI  | 19  |
| **Total** | **19** |

> **Note**: This feature modifies CLI structured output only. No UI or API interfaces are affected. All test cases are CLI type.

---

## CLI Test Cases

### claim command

## TC-001: Claim outputs only essential fields
- **Source**: Proposal "Success Criteria" item 1 + Task 1 AC-1
- **Type**: CLI
- **Target**: cli/claim
- **Test ID**: cli/claim/claim-outputs-only-essential-fields
- **Pre-conditions**: At least one pending task exists in the task index
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge task claim` to claim a task
  2. Capture the structured output block (between `---` delimiters)
  3. Count the number of fields in the output block
- **Expected**: Output contains exactly ACTION, TASK_ID, FEATURE, FILE fields. Optional fields (SCOPE, BREAKING, MAIN_SESSION) may appear only when their conditions are met. No other fields are present.
- **Priority**: P0

## TC-002: Claim output includes ACTION CLAIMED
- **Source**: Task 1 AC-2 — `printNewTask()` wraps with ACTION: CLAIMED
- **Type**: CLI
- **Target**: cli/claim
- **Test ID**: cli/claim/claim-output-includes-action-claimed
- **Pre-conditions**: A pending task is available to claim
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge task claim` to claim a new task
  2. Read the output block
  3. Check for the ACTION field
- **Expected**: ACTION field has value "CLAIMED"
- **Priority**: P0

## TC-003: Claim output includes TASK_ID
- **Source**: Proposal — TASK_ID is consumed "everywhere"
- **Type**: CLI
- **Target**: cli/claim
- **Test ID**: cli/claim/claim-output-includes-task-id
- **Pre-conditions**: A pending task is available to claim
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge task claim`
  2. Read the output block
  3. Check for TASK_ID field
- **Expected**: TASK_ID field is present and matches the claimed task's ID
- **Priority**: P0

## TC-004: Claim output includes FEATURE
- **Source**: Proposal — FEATURE is consumed by "E2E gate"
- **Type**: CLI
- **Target**: cli/claim
- **Test ID**: cli/claim/claim-output-includes-feature
- **Pre-conditions**: A pending task exists within a feature
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge task claim`
  2. Read the output block
  3. Check for FEATURE field
- **Expected**: FEATURE field is present and matches the feature slug of the claimed task
- **Priority**: P0

## TC-005: Claim output includes FILE
- **Source**: Proposal — FILE is consumed by "agent reads task file"
- **Type**: CLI
- **Target**: cli/claim
- **Test ID**: cli/claim/claim-output-includes-file
- **Pre-conditions**: A pending task exists with a valid file path
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge task claim`
  2. Read the output block
  3. Check for FILE field
- **Expected**: FILE field is present and contains the absolute path to the task markdown file
- **Priority**: P0

## TC-006: Claim SCOPE present when non-empty
- **Source**: Task 1 AC-1 — SCOPE only when non-empty
- **Type**: CLI
- **Target**: cli/claim
- **Test ID**: cli/claim/claim-scope-present-when-non-empty
- **Pre-conditions**: A pending task has a non-empty scope (e.g., "backend")
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge task claim` for a task with scope = "backend"
  2. Read the output block
  3. Check for SCOPE field
- **Expected**: SCOPE field is present with value "backend"
- **Priority**: P1

## TC-007: Claim SCOPE absent when empty
- **Source**: Task 1 AC-1 — SCOPE only when non-empty
- **Type**: CLI
- **Target**: cli/claim
- **Test ID**: cli/claim/claim-scope-absent-when-empty
- **Pre-conditions**: A pending task has an empty scope
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge task claim` for a task with empty scope
  2. Read the output block
  3. Check that SCOPE field is not present
- **Expected**: SCOPE field is absent from the output block
- **Priority**: P1

## TC-008: Claim BREAKING present when true
- **Source**: Task 1 AC-5 — BREAKING present with "true" when true
- **Type**: CLI
- **Target**: cli/claim
- **Test ID**: cli/claim/claim-breaking-present-when-true
- **Pre-conditions**: A pending task has breaking = true
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge task claim` for a task with breaking = true
  2. Read the output block
  3. Check for BREAKING field
- **Expected**: BREAKING field is present with value "true"
- **Priority**: P1

## TC-009: Claim BREAKING absent when false
- **Source**: Task 1 AC-5 — Boolean fields absent when false
- **Type**: CLI
- **Target**: cli/claim
- **Test ID**: cli/claim/claim-breaking-absent-when-false
- **Pre-conditions**: A pending task has breaking = false
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge task claim` for a task with breaking = false
  2. Read the output block
  3. Check that BREAKING field is not present
- **Expected**: BREAKING field is absent from the output block
- **Priority**: P1

## TC-010: Claim MAIN_SESSION present when true
- **Source**: Task 1 AC-5 — MAIN_SESSION present with "true" when true
- **Type**: CLI
- **Target**: cli/claim
- **Test ID**: cli/claim/claim-main-session-present-when-true
- **Pre-conditions**: A pending task has mainSession = true
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge task claim` for a task with mainSession = true
  2. Read the output block
  3. Check for MAIN_SESSION field
- **Expected**: MAIN_SESSION field is present with value "true"
- **Priority**: P1

## TC-011: Claim MAIN_SESSION absent when false
- **Source**: Task 1 AC-5 — Boolean fields absent when false
- **Type**: CLI
- **Target**: cli/claim
- **Test ID**: cli/claim/claim-main-session-absent-when-false
- **Pre-conditions**: A pending task has mainSession = false
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge task claim` for a task with mainSession = false
  2. Read the output block
  3. Check that MAIN_SESSION field is not present
- **Expected**: MAIN_SESSION field is absent from the output block
- **Priority**: P1

## TC-012: Claim removed fields not present
- **Source**: Task 1 AC-4 — Removed fields no longer appear
- **Type**: CLI
- **Target**: cli/claim
- **Test ID**: cli/claim/claim-removed-fields-not-present
- **Pre-conditions**: A pending task is available to claim
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge task claim`
  2. Read the output block
  3. Check for each removed field: KEY, TITLE, PRIORITY, STATUS, ESTIMATED_TIME, DEPENDENCIES, TYPE, PROFILE, NO_TEST, RECORD
- **Expected**: None of the 10 removed fields appear in the output block
- **Priority**: P0

## TC-013: Claim CONTINUE wraps with ACTION CONTINUE and STARTED_AT
- **Source**: Task 1 AC-3 — `printContinueTask()` wraps with ACTION: CONTINUE + trimmed fields + STARTED_AT
- **Type**: CLI
- **Target**: cli/claim
- **Test ID**: cli/claim/claim-continue-wraps-with-action-continue-and-started-at
- **Pre-conditions**: A task is in "in_progress" state (previously claimed)
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge task claim` when an in-progress task exists
  2. Read the output block
  3. Check for ACTION: CONTINUE and STARTED_AT fields
- **Expected**: ACTION field has value "CONTINUE", STARTED_AT is present, and trimmed essential fields (TASK_ID, FEATURE, FILE) are present
- **Priority**: P0

## TC-014: Claim field order matches specification
- **Source**: Task 1 Implementation Notes — field order: ACTION, TASK_ID, FEATURE, FILE, SCOPE, BREAKING, MAIN_SESSION
- **Type**: CLI
- **Target**: cli/claim
- **Test ID**: cli/claim/claim-field-order-matches-specification
- **Pre-conditions**: A pending task is available with all conditional fields active (SCOPE non-empty, BREAKING true, MAIN_SESSION true)
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge task claim` for a task with all conditional fields active
  2. Read the output block
  3. Extract field names in order
- **Expected**: Fields appear in order: ACTION, TASK_ID, FEATURE, FILE, SCOPE, BREAKING, MAIN_SESSION. No KEY or TYPE fields appear.
- **Priority**: P1

### submit command

## TC-015: Submit outputs only STATUS field
- **Source**: Task 2 AC-1 + Proposal "Success Criteria" item 2
- **Type**: CLI
- **Target**: cli/submit
- **Test ID**: cli/submit/submit-outputs-only-status-field
- **Pre-conditions**: A task is in "in_progress" state (claimed)
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge task submit <task-id>` for an in-progress task (non-JSON, non-quiet mode)
  2. Read the output block
  3. Count the fields
- **Expected**: Output block contains exactly 1 field: STATUS. TASK_ID and RECORD_FILE are not present.
- **Priority**: P0

## TC-016: Submit JSON mode unchanged
- **Source**: Task 2 AC-5 — JSON mode (`--json`) in submit is NOT changed
- **Type**: CLI
- **Target**: cli/submit
- **Test ID**: cli/submit/submit-json-mode-unchanged
- **Pre-conditions**: A task is in "in_progress" state
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge task submit <task-id> --json`
  2. Capture the JSON output
  3. Compare with expected JSON structure
- **Expected**: JSON output structure is unchanged from the pre-feature version (all original JSON fields present)
- **Priority**: P0

### query command

## TC-017: Query outputs essential fields with conditional SCOPE and BREAKING
- **Source**: Task 2 AC-2 + Proposal "Success Criteria" item 3
- **Type**: CLI
- **Target**: cli/query
- **Test ID**: cli/query/query-outputs-essential-fields-with-conditional-scope-and-breaking
- **Pre-conditions**: A task exists with known state
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge task query <task-id>` for a task with scope = "backend" and breaking = true
  2. Read the output block
  3. Check for TASK_ID, STATUS, SCOPE, BREAKING fields
- **Expected**: Output contains TASK_ID, STATUS, SCOPE ("backend"), BREAKING ("true"). No KEY, TITLE, PRIORITY, ESTIMATED_TIME, DEPENDENCIES, FILE, or RECORD fields.
- **Priority**: P0

## TC-018: Query omits SCOPE when empty and BREAKING when false
- **Source**: Task 2 AC-2 — SCOPE (when non-empty), BREAKING (when true)
- **Type**: CLI
- **Target**: cli/query
- **Test ID**: cli/query/query-omits-scope-when-empty-and-breaking-when-false
- **Pre-conditions**: A task exists with empty scope and breaking = false
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge task query <task-id>` for a task with empty scope and breaking = false
  2. Read the output block
  3. Check that SCOPE and BREAKING are absent
- **Expected**: Output contains only TASK_ID and STATUS. SCOPE and BREAKING fields are absent.
- **Priority**: P1

### status command

## TC-019: Status outputs only TASK_ID and STATUS
- **Source**: Task 2 AC-3, AC-4 — status outputs TASK_ID + STATUS (both query and update modes)
- **Type**: CLI
- **Target**: cli/status
- **Test ID**: cli/status/status-outputs-only-task-id-and-status
- **Pre-conditions**: A task exists in the index
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge task status <task-id>` in query mode
  2. Read the output block
  3. Check that only TASK_ID and STATUS are present
  4. Run `forge task status <task-id> completed` in update mode
  5. Read the output block
  6. Check that only TASK_ID and STATUS are present
- **Expected**: Both modes output exactly TASK_ID and STATUS. KEY, TITLE, and DEPENDENCIES are absent.
- **Priority**: P0

---

## Traceability

| TC ID | Source | Type | Target | Priority |
|-------|--------|------|--------|----------|
| TC-001 | Proposal Success Criteria #1, Task 1 AC-1 | CLI | cli/claim | P0 |
| TC-002 | Task 1 AC-2 | CLI | cli/claim | P0 |
| TC-003 | Proposal — TASK_ID consumed everywhere | CLI | cli/claim | P0 |
| TC-004 | Proposal — FEATURE consumed by E2E gate | CLI | cli/claim | P0 |
| TC-005 | Proposal — FILE consumed by agent | CLI | cli/claim | P0 |
| TC-006 | Task 1 AC-1 — SCOPE conditional (present) | CLI | cli/claim | P1 |
| TC-007 | Task 1 AC-1 — SCOPE conditional (absent) | CLI | cli/claim | P1 |
| TC-008 | Task 1 AC-5 — BREAKING present when true | CLI | cli/claim | P1 |
| TC-009 | Task 1 AC-5 — BREAKING absent when false | CLI | cli/claim | P1 |
| TC-010 | Task 1 AC-5 — MAIN_SESSION present when true | CLI | cli/claim | P1 |
| TC-011 | Task 1 AC-5 — MAIN_SESSION absent when false | CLI | cli/claim | P1 |
| TC-012 | Task 1 AC-4 — Removed fields absent | CLI | cli/claim | P0 |
| TC-013 | Task 1 AC-3 — CONTINUE mode | CLI | cli/claim | P0 |
| TC-014 | Task 1 Implementation Notes — field order | CLI | cli/claim | P1 |
| TC-015 | Task 2 AC-1, Proposal Success Criteria #2 | CLI | cli/submit | P0 |
| TC-016 | Task 2 AC-5 — JSON mode unchanged | CLI | cli/submit | P0 |
| TC-017 | Task 2 AC-2, Proposal Success Criteria #3 | CLI | cli/query | P0 |
| TC-018 | Task 2 AC-2 — conditional fields absent | CLI | cli/query | P1 |
| TC-019 | Task 2 AC-3, AC-4 — status lean output | CLI | cli/status | P0 |
