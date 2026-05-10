---
id: "T-quick-5"
title: "Verify Quick E2E Regression"
priority: "P1"
estimated_time: "15min"
dependencies: ["T-quick-4"]
status: pending
mainSession: false
---

# Verify Quick E2E Regression

## Description

Run the full e2e regression suite to verify graduated specs integrate cleanly
with existing tests.

## Reference Files

- `tests/e2e/` — Full regression suite
- `tests/e2e/.graduated/<slug>` — Graduation marker from T-quick-4

## Acceptance Criteria

- [ ] `just test-e2e` passes (full suite, no --feature flag)
- [ ] All graduated and existing specs pass

## Implementation Notes

1. Run `just e2e-setup` (idempotent — skips if already set up)
2. Run: `just test-e2e`
3. On success: mark completed

**On failure**:
- Read Playwright output for failure details (check `tests/e2e/test-results/` and terminal output)
- Analyze each failure: is it a code bug, test script issue, or environment issue?
- Run `task template fix-task` to view the fix-task template and required variables
- For each distinct root cause, create a fix task:
  ```bash
  task add --template fix-task \
           --title "Fix: <concise description>" \
           --source-task-id T-quick-5 \
           --block-source \
           --var SOURCE_FILES="<affected source file paths>" \
           --var TEST_SCRIPT="tests/e2e/<failing-spec>.spec.ts" \
           --var TEST_RESULTS="tests/e2e/test-results/" \
           --description "<root cause and context>"
  ```
  `task add` automatically deduplicates — check output: `ACTION: ADDED` (new fix task) or `ACTION: SKIPPED` (active fix already exists).
- Fix tasks (P0) will be claimed before other P1/P2 tasks
- After fix tasks complete, T-quick-5 is unblocked and re-claimed for re-run

**Do NOT** attempt to fix failures inline — create fix tasks and let the dispatcher handle them.

## Execution Workflow

1. Set up e2e environment (idempotent).
   - Command: `just e2e-setup`
   - Success: exit 0 (skips if already set up).
   - Failure: fix environment setup, retry this step.
2. Run full e2e regression suite.
   - Command: `just test-e2e` (no `--feature` flag — runs all specs).
   - Success: all tests pass (exit 0).
   - Failure: proceed to step 3.
3. Classify and create fix tasks (only on failure).
   - Action: read Playwright output and `tests/e2e/test-results/` for failure details.
   - For each distinct root cause, run `task add --template fix-task --source-task-id T-quick-5 --block-source ...`.
   - Success: fix tasks created.
   - Failure: if `task add` fails, set task status to `blocked` and stop.
4. Record results.
   - On all pass: record with status `completed`.
   - On failures with fix tasks created: record with testsFailed > 0 (CLI auto-downgrades to `blocked`).
5. Stop. Proceed to Step 3 (Record).
