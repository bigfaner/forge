---
id: "T-test-4.5"
title: "Verify Full E2E Regression"
priority: "P1"
estimated_time: "15-30min"
dependencies: ["T-test-4"]
status: pending
noTest: false
mainSession: false
---

# Verify Full E2E Regression

## Description

Run the full e2e regression suite to verify graduated specs integrate cleanly
with existing tests.

## Reference Files

- `tests/e2e/` — Full regression suite
- `tests/e2e/.graduated/<slug>` — Graduation marker from T-test-4

## Acceptance Criteria

- [ ] `just test-e2e` passes (full suite, no --feature flag)
- [ ] All graduated and existing specs pass

## User Stories

No direct user story mapping. This is a standard regression verification task.

## Implementation Notes

1. Run `just e2e-setup` (idempotent — skips if already set up)
2. Run: `just test-e2e`
3. On success: mark completed, continue to T-test-5

**On failure**:
- Read Playwright output for failure details (check `tests/e2e/test-results/` and terminal output)
- Analyze each failure: is it a code bug, test script issue, or environment issue?
- Run `task template fix-task` to view the fix-task template and required variables
- For each distinct root cause, create a fix task:
  ```bash
  task add --template fix-task \
           --title "Fix: <concise description>" \
           --source-task-id T-test-4.5 \
           --block-source \
           --var SOURCE_FILES="<affected source file paths>" \
           --var TEST_SCRIPT="tests/e2e/<failing-spec>.spec.ts" \
           --var TEST_RESULTS="tests/e2e/test-results/" \
           --description "<root cause and context>"
  ```
  `task add` automatically deduplicates — check output: `ACTION: ADDED` (new fix task) or `ACTION: SKIPPED` (active fix already exists).
- When you record this task, the CLI will auto-downgrade `status: "completed"` to `status: "blocked"` if `testsFailed > 0` — do NOT manually set status to blocked
- Fix tasks (P0) will be claimed before T-test-5 (P2)
- After fix tasks complete, T-test-4.5 is auto-restored to pending and re-claimed for re-run

**Do NOT** attempt to fix failures inline — create fix tasks and let the dispatcher handle them.
