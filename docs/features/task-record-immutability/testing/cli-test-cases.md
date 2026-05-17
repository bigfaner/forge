# CLI Test Cases: task-record-immutability

Profile: go-test

## TC-001: Submit record succeeds when no record exists
- **Source**: Proposal / Key Scenario 1, Success Criteria item 1 (implicit happy path)
- **Type**: CLI
- **Target**: cli/task-submit
- **Test ID**: cli/task-submit/submit-record-succeeds-when-no-record-exists
- **Pre-conditions**: A valid task exists in `tasks/index.json` with status `pending` or `in-progress`. No record file exists at `tasks/records/<task>.md`. `CLAUDE_PROJECT_DIR` is set to an isolated fixture directory containing a valid `.forge/` project structure with `tasks/index.json`.
- **Steps**:
  1. Run `forge task submit <task-id>` with a valid submission payload
- **Expected**: Exit code 0. Record file created at `tasks/records/<task>.md` containing the submission content.
- **Priority**: P0

## TC-002: Submit record blocked when record already exists
- **Source**: Proposal / Key Scenario 2, Success Criteria item 1
- **Type**: CLI
- **Target**: cli/task-submit
- **Test ID**: cli/task-submit/submit-record-blocked-when-record-already-exists
- **Pre-conditions**: A valid task exists in `tasks/index.json`. A record file already exists at `tasks/records/<task>.md`. `CLAUDE_PROJECT_DIR` is set to an isolated fixture directory.
- **Steps**:
  1. Run `forge task submit <task-id>` (without `--force`)
- **Expected**: Exit code 1. stderr contains `Record for task <task-id> already exists`. stderr contains `Use --force to overwrite, or create a fix task instead.` Existing record file is NOT modified.
- **Priority**: P0

## TC-003: Submit record with --force overwrites existing record
- **Source**: Proposal / Key Scenario 3, Success Criteria item 2
- **Type**: CLI
- **Target**: cli/task-submit
- **Test ID**: cli/task-submit/submit-with-force-overwrites-existing-record
- **Pre-conditions**: A valid task exists in `tasks/index.json`. A record file already exists at `tasks/records/<task>.md` with known content. `CLAUDE_PROJECT_DIR` is set to an isolated fixture directory.
- **Steps**:
  1. Run `forge task submit <task-id> --force` with new submission content
- **Expected**: Exit code 0. stderr contains `WARNING: Overwriting existing record`. Record file at `tasks/records/<task>.md` is replaced with the new submission content.
- **Priority**: P0

## TC-004: Default query shows 4 fields unchanged
- **Source**: Proposal / Key Scenario 4, Success Criteria item 3
- **Type**: CLI
- **Target**: cli/task-query
- **Test ID**: cli/task-query/default-query-shows-4-fields-unchanged
- **Pre-conditions**: A valid task exists in `tasks/index.json` with known values for TASK_ID, STATUS, SCOPE, and BREAKING fields. `CLAUDE_PROJECT_DIR` is set to an isolated fixture directory.
- **Steps**:
  1. Run `forge task query <task-id>` (no flags)
- **Expected**: Exit code 0. stdout contains exactly the 4 fields: TASK_ID, STATUS, SCOPE (if set), BREAKING (if true). stdout does NOT contain TITLE, PRIORITY, TYPE, DEPENDENCIES, TASK_FILE, RECORD_FILE, KEY, or RELATED_FIXES.
- **Priority**: P1

## TC-005: Verbose query displays all task fields
- **Source**: Proposal / Key Scenario 5, Success Criteria item 4
- **Type**: CLI
- **Target**: cli/task-query
- **Test ID**: cli/task-query/verbose-query-displays-all-task-fields
- **Pre-conditions**: A valid task exists in `tasks/index.json` with populated fields: KEY, TASK_ID, TITLE, STATUS, PRIORITY, TYPE, SCOPE, DEPENDENCIES, TASK_FILE, RECORD_FILE. `CLAUDE_PROJECT_DIR` is set to an isolated fixture directory.
- **Steps**:
  1. Run `forge task query <task-id> --verbose`
- **Expected**: Exit code 0. stdout contains all fields: KEY, TASK_ID, TITLE, STATUS, PRIORITY, TYPE, SCOPE, DEPENDENCIES, TASK_FILE, RECORD_FILE. Output is bounded by `>>>` and `<<<` markers.
- **Priority**: P0

## TC-006: Verbose query shows RELATED_FIXES for tasks with fix records
- **Source**: Proposal / Success Criteria item 5
- **Type**: CLI
- **Target**: cli/task-query
- **Test ID**: cli/task-query/verbose-query-shows-related-fixes
- **Pre-conditions**: A valid task (e.g., task "2") exists in `tasks/index.json`. A fix task exists in `tasks/index.json` with `sourceTaskID` set to "2". `CLAUDE_PROJECT_DIR` is set to an isolated fixture directory.
- **Steps**:
  1. Run `forge task query 2 --verbose`
- **Expected**: Exit code 0. stdout contains `RELATED_FIXES` field. RELATED_FIXES shows each fix as `<id> [<status>] <title>`, one per line. Fix task ID, status, and title match the values from `tasks/index.json`.
- **Priority**: P0

## TC-007: Verbose query omits RELATED_FIXES when no fixes exist
- **Source**: Proposal / Success Criteria item 6
- **Type**: CLI
- **Target**: cli/task-query
- **Test ID**: cli/task-query/verbose-query-omits-related-fixes-when-none-exist
- **Pre-conditions**: A valid task exists in `tasks/index.json`. No fix tasks in `tasks/index.json` have `sourceTaskID` matching this task. `CLAUDE_PROJECT_DIR` is set to an isolated fixture directory.
- **Steps**:
  1. Run `forge task query <task-id> --verbose`
- **Expected**: Exit code 0. stdout does NOT contain `RELATED_FIXES` field. All other verbose fields (KEY, TASK_ID, TITLE, STATUS, PRIORITY, TYPE, SCOPE, DEPENDENCIES, TASK_FILE, RECORD_FILE) are present.
- **Priority**: P1

## TC-008: Status command behavior unchanged
- **Source**: Proposal / Success Criteria item 7
- **Type**: CLI
- **Target**: cli/task-status
- **Test ID**: cli/task-status/status-command-behavior-unchanged
- **Pre-conditions**: A valid task exists in `tasks/index.json` with a known status. `CLAUDE_PROJECT_DIR` is set to an isolated fixture directory.
- **Steps**:
  1. Run `forge task status <task-id>`
- **Expected**: Exit code 0. Output matches the existing `forge task status` behavior (displays the task status). Output is identical to what the command produced before the write-once and verbose query changes.
- **Priority**: P1

## TC-009: Verbose query with short flag -v
- **Source**: Proposal / Proposed Solution 2 ("--verbose / -v flag")
- **Type**: CLI
- **Target**: cli/task-query
- **Test ID**: cli/task-query/verbose-query-with-short-flag
- **Pre-conditions**: A valid task exists in `tasks/index.json` with populated fields. `CLAUDE_PROJECT_DIR` is set to an isolated fixture directory.
- **Steps**:
  1. Run `forge task query <task-id> -v`
- **Expected**: Exit code 0. stdout output is identical to `forge task query <task-id> --verbose` (same fields, same format).
- **Priority**: P2

## TC-010: Verbose query omits SCOPE when empty
- **Source**: Proposal / Proposed Solution 2 ("SCOPE (omit if empty)")
- **Type**: CLI
- **Target**: cli/task-query
- **Test ID**: cli/task-query/verbose-query-omits-scope-when-empty
- **Pre-conditions**: A valid task exists in `tasks/index.json` with an empty SCOPE field. `CLAUDE_PROJECT_DIR` is set to an isolated fixture directory.
- **Steps**:
  1. Run `forge task query <task-id> --verbose`
- **Expected**: Exit code 0. stdout does NOT contain `SCOPE` field. All other verbose fields are present.
- **Priority**: P2

## TC-011: Verbose query omits BREAKING when false
- **Source**: Proposal / Proposed Solution 2 (default mode: "BREAKING (if true)")
- **Type**: CLI
- **Target**: cli/task-query
- **Test ID**: cli/task-query/verbose-query-omits-breaking-when-false
- **Pre-conditions**: A valid task exists in `tasks/index.json` with BREAKING set to false or unset. `CLAUDE_PROJECT_DIR` is set to an isolated fixture directory.
- **Steps**:
  1. Run `forge task query <task-id> --verbose`
- **Expected**: Exit code 0. stdout does NOT contain `BREAKING` field. All other verbose fields are present.
- **Priority**: P2

## TC-012: Verbose query displays multi-line DEPENDENCIES
- **Source**: Proposal / Proposed Solution 2 ("DEPENDENCIES (multi-line if multiple)")
- **Type**: CLI
- **Target**: cli/task-query
- **Test ID**: cli/task-query/verbose-query-displays-multi-line-dependencies
- **Pre-conditions**: A valid task exists in `tasks/index.json` with multiple dependencies (e.g., ["1", "3"]). `CLAUDE_PROJECT_DIR` is set to an isolated fixture directory.
- **Steps**:
  1. Run `forge task query <task-id> --verbose`
- **Expected**: Exit code 0. stdout contains `DEPENDENCIES` field with each dependency displayed on a separate line.
- **Priority**: P2

## Traceability

| TC ID | Source | Type | Target | Priority |
|-------|--------|------|--------|----------|
| TC-001 | Proposal / Key Scenario 1, SC-1 | CLI | cli/task-submit | P0 |
| TC-002 | Proposal / Key Scenario 2, SC-1 | CLI | cli/task-submit | P0 |
| TC-003 | Proposal / Key Scenario 3, SC-2 | CLI | cli/task-submit | P0 |
| TC-004 | Proposal / Key Scenario 4, SC-3 | CLI | cli/task-query | P1 |
| TC-005 | Proposal / Key Scenario 5, SC-4 | CLI | cli/task-query | P0 |
| TC-006 | Proposal / SC-5 | CLI | cli/task-query | P0 |
| TC-007 | Proposal / SC-6 | CLI | cli/task-query | P1 |
| TC-008 | Proposal / SC-7 | CLI | cli/task-status | P1 |
| TC-009 | Proposal / Proposed Solution 2 | CLI | cli/task-query | P2 |
| TC-010 | Proposal / Proposed Solution 2 | CLI | cli/task-query | P2 |
| TC-011 | Proposal / Proposed Solution 2 | CLI | cli/task-query | P2 |
| TC-012 | Proposal / Proposed Solution 2 | CLI | cli/task-query | P2 |
