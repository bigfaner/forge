---
feature: "quality-gate-fix-task-loop-breaker"
sources:
  - docs/proposals/quality-gate-fix-task-loop-breaker/proposal.md
generated: "2026-05-16"
---

# Test Cases: quality-gate-fix-task-loop-breaker

## Summary

| Type | Count |
|------|-------|
| UI   | 0   |
| TUI  | 0  |
| Mobile | 0 |
| API  | 0  |
| CLI  | 11  |
| **Total** | **11** |

---

## UI Test Cases

_None. Project does not expose a web-UI interface._

---

## TUI Test Cases

_None. Project does not expose a terminal-UI (full-screen rendering) interface._

---

## Mobile Test Cases

_None. Project does not expose a mobile interface._

---

## API Test Cases

_None. Project does not expose HTTP API endpoints._

---

## CLI Test Cases

### TC-001: addFixTask creates tasks with step-scoped SourceTaskID sentinel

- **Source**: Proposal Success Criteria #1 — "`addFixTask` creates tasks with `SourceTaskID: \"quality-gate:<step>\"` (verified by test per step)"
- **Type**: CLI
- **Target**: cli/quality-gate
- **Test ID**: cli/quality-gate/addfixtask-creates-step-scoped-sourcetaskid-sentinel
- **Pre-conditions**: A feature with at least one completed task and forge state set; a valid fix-task template exists in the project. Create a temp project dir with `features/<slug>/tasks/index.json` containing a completed task, `.forge/state.json` with the feature slug, and a fix-task template.
- **Steps**:
  1. Run `forge quality-gate` (or invoke `addFixTask` directly) with a step name such as "compile"
  2. Load the updated `index.json` from disk
  3. Inspect the newly created fix task's `SourceTaskID` field
- **Expected**: The fix task's `SourceTaskID` equals `"quality-gate:compile"` — the step-scoped sentinel format. The `Vars["SOURCE_TASK_ID"]` field equals `"N/A (project-wide gate)"` for template rendering (intentionally diverges from the struct field).
- **Priority**: P0

### TC-002: countFixTasks counts fix tasks cumulatively regardless of status

- **Source**: Proposal Success Criteria #2 — "`countFixTasks` counts fix tasks regardless of status (completed + active + blocked)"
- **Type**: CLI
- **Target**: cli/quality-gate
- **Test ID**: cli/quality-gate/countfixtasks-counts-cumulatively-regardless-of-status
- **Pre-conditions**: A feature task index containing fix tasks for the same step with various statuses (one completed, one active, one blocked, one skipped). Create a temp project dir with `features/<slug>/tasks/index.json` containing these tasks, all with `SourceTaskID: "quality-gate:unit-test"` and titles prefixed with "fix unit-test:".
- **Steps**:
  1. Load the task index from disk
  2. Call `countFixTasks(index, "unit-test")`
  3. Verify the returned count includes ALL statuses
- **Expected**: The count is 4 (completed + active + blocked + skipped). No status is excluded.
- **Priority**: P0

### TC-003: Unit tests are retried once before creating a fix task

- **Source**: Proposal Success Criteria #3 — "When unit tests fail, they are retried once before creating a fix task"
- **Type**: CLI
- **Target**: cli/quality-gate
- **Test ID**: cli/quality-gate/unit-tests-retried-once-before-fix-task
- **Pre-conditions**: A feature with all tasks completed and forge state set. A mock test runner that fails on the first call and also fails on the second call. Create a temp project dir with completed tasks, `.forge/state.json`, and inject a failing test runner stub.
- **Steps**:
  1. Invoke the quality gate unit-test step (or `runUnitTestStep` directly) with a test runner that fails both times
  2. Count how many times the test runner was invoked
  3. Verify a fix task was created after the second failure
- **Expected**: The test runner is invoked exactly 2 times (initial + retry). A fix task is created after both attempts fail.
- **Priority**: P0

### TC-004: Retry passes — warning logged, no fix task created

- **Source**: Proposal Success Criteria #4 — "If retry passes, a warning is logged and no fix task is created"
- **Type**: CLI
- **Target**: cli/quality-gate
- **Test ID**: cli/quality-gate/retry-passes-warning-logged-no-fix-task
- **Pre-conditions**: A feature with all tasks completed and forge state set. A mock test runner that fails on the first call and passes on the second call. Create a temp project dir with completed tasks and `.forge/state.json`.
- **Steps**:
  1. Invoke `runUnitTestStep` with a test runner that fails first, passes on retry
  2. Capture stderr output
  3. Check the task index for any new fix tasks
- **Expected**: `runUnitTestStep` returns `(true, "", nil)` (passed, no fix ID, no error). Stderr contains "WARNING: unit tests passed on retry (transient failure)". No fix task exists in the task index.
- **Priority**: P0

### TC-005: Retry fails — fix-task description mentions "retried once, both attempts failed"

- **Source**: Proposal Success Criteria #5 — "If retry fails, fix-task description mentions 'retried once, both attempts failed'"
- **Type**: CLI
- **Target**: cli/quality-gate
- **Test ID**: cli/quality-gate/retry-fail-description-mentions-retried-once
- **Pre-conditions**: A feature with all tasks completed, forge state set, and a fix-task template. A mock test runner that fails both times with distinct error output per attempt. Create a temp project dir with the required files.
- **Steps**:
  1. Invoke `runUnitTestStep` with a test runner that fails both times, returning "FAIL: attempt-1" and "FAIL: attempt-2" respectively
  2. Read the fix task markdown file from disk
  3. Inspect the task description content
- **Expected**: The written `tests/results/unit-raw-output.txt` contains "retried once, both attempts failed", "=== First attempt ===", and "=== Retry attempt ===" sections with the respective error output. The fix task is created.
- **Priority**: P0

### TC-006: Cumulative cap stops fix-task creation after 3 per step

- **Source**: Proposal Success Criteria #6 — "After 3 cumulative fix tasks for a step, no more are created regardless of status"
- **Type**: CLI
- **Target**: cli/quality-gate
- **Test ID**: cli/quality-gate/cumulative-cap-stops-fix-task-after-3
- **Pre-conditions**: A feature task index with 3 existing fix tasks for "unit-test" step (all with `SourceTaskID: "quality-gate:unit-test"` and title prefix "fix unit-test:"), each with different statuses (completed, skipped, active). Create a temp project dir with this pre-populated index, a completed main task, `.forge/state.json`, and a fix-task template.
- **Steps**:
  1. Invoke `addFixTask` for the "unit-test" step
  2. Check the return value
- **Expected**: `addFixTask` returns `("", ErrMaxFixTasks)`. Stderr contains "max fix-tasks reached for unit-test, manual intervention required". No new fix task is created in the index.
- **Priority**: P0

### TC-007: Cross-step independence — fix for step A does not block step B

- **Source**: Proposal Success Criteria #7 — "Pending fix for step A does NOT block fix creation for step B (cross-step independence)"
- **Type**: CLI
- **Target**: cli/quality-gate
- **Test ID**: cli/quality-gate/cross-step-independence-fix-a-does-not-block-b
- **Pre-conditions**: A feature task index with 3 existing fix tasks for "compile" step (all completed/skipped, capped). Create a temp project dir with this index, a completed main task, `.forge/state.json`, and a fix-task template.
- **Steps**:
  1. Invoke `addFixTask` for the "unit-test" step (different step from the capped one)
  2. Check the return value and the task index
- **Expected**: `addFixTask` succeeds — returns a non-empty task ID and no error. The new fix task has `SourceTaskID: "quality-gate:unit-test"` (not "quality-gate:compile"). The cap on "compile" does not affect "unit-test".
- **Priority**: P0

### TC-008: addFixTask returns explicit errors on template-not-found

- **Source**: Proposal Success Criteria #8 — "`addFixTask` returns explicit errors on template-not-found, task-add-failure, and markdown-creation-failure"
- **Type**: CLI
- **Target**: cli/quality-gate
- **Test ID**: cli/quality-gate/addfixtask-returns-explicit-errors-template-not-found
- **Pre-conditions**: A feature with completed tasks and forge state set, but no "fix-task" template registered. Create a temp project dir with completed tasks and `.forge/state.json`, but no templates directory or an empty templates directory.
- **Steps**:
  1. Invoke `addFixTask` with a step name (e.g., "compile")
  2. Inspect the error return value
- **Expected**: `addFixTask` returns `("", error)` where the error message contains "template \"fix-task\" not found". No fix task is created.
- **Priority**: P1

### TC-009: addFixTask returns explicit errors on task-add failure

- **Source**: Proposal Success Criteria #8 — "`addFixTask` returns explicit errors on ... task-add-failure"
- **Type**: CLI
- **Target**: cli/quality-gate
- **Test ID**: cli/quality-gate/addfixtask-returns-explicit-errors-task-add-failure
- **Pre-conditions**: A feature with a valid fix-task template but an invalid/malformed task index that will cause `AddTask` to fail (e.g., index.json contains invalid JSON or is read-only). Create a temp project dir with the fix-task template but a corrupted index.json.
- **Steps**:
  1. Invoke `addFixTask` with a step name
  2. Inspect the error return value
- **Expected**: `addFixTask` returns `("", error)` where the error message contains "failed to add fix task". No fix task is created.
- **Priority**: P1

### TC-010: addFixTask returns explicit errors on markdown creation failure

- **Source**: Proposal Success Criteria #8 — "`addFixTask` returns explicit errors on ... markdown-creation-failure"
- **Type**: CLI
- **Target**: cli/quality-gate
- **Test ID**: cli/quality-gate/addfixtask-returns-explicit-errors-markdown-creation-failure
- **Pre-conditions**: A feature with a valid fix-task template and valid task index, but the tasks directory is read-only (or does not exist and cannot be created). Create a temp project dir with a valid template, writable index, but a read-only tasks directory.
- **Steps**:
  1. Invoke `addFixTask` with a step name
  2. Inspect the error return value
- **Expected**: `addFixTask` returns `("", error)` where the error message contains "failed to create fix task file". The task is added to the index (in memory) but the markdown file is not created on disk.
- **Priority**: P1

### TC-011: Version bumped in scripts/version.txt

- **Source**: Proposal Success Criteria #9 — "Version bumped in `scripts/version.txt`"
- **Type**: CLI
- **Target**: cli/quality-gate
- **Test ID**: cli/quality-gate/version-bumped-in-scripts-version-txt
- **Pre-conditions**: A working copy of the repository at the current branch with the quality-gate changes applied. The file `scripts/version.txt` exists.
- **Steps**:
  1. Read `scripts/version.txt` from disk
  2. Compare the version against the known pre-change version
- **Expected**: The version has been incremented (patch bump at minimum). The version string is valid semver.
- **Priority**: P2

---

## Traceability

| TC ID | Source | Type | Target | Priority |
|-------|--------|------|--------|----------|
| TC-001 | Proposal SC #1 | CLI | cli/quality-gate | P0 |
| TC-002 | Proposal SC #2 | CLI | cli/quality-gate | P0 |
| TC-003 | Proposal SC #3 | CLI | cli/quality-gate | P0 |
| TC-004 | Proposal SC #4 | CLI | cli/quality-gate | P0 |
| TC-005 | Proposal SC #5 | CLI | cli/quality-gate | P0 |
| TC-006 | Proposal SC #6 | CLI | cli/quality-gate | P0 |
| TC-007 | Proposal SC #7 | CLI | cli/quality-gate | P0 |
| TC-008 | Proposal SC #8 (template-not-found) | CLI | cli/quality-gate | P1 |
| TC-009 | Proposal SC #8 (task-add-failure) | CLI | cli/quality-gate | P1 |
| TC-010 | Proposal SC #8 (markdown-creation-failure) | CLI | cli/quality-gate | P1 |
| TC-011 | Proposal SC #9 | CLI | cli/quality-gate | P2 |
