---
journey: "multi-surface-interleaved"
step: 3
step-action: "Execute the first surface's run-tests task"
generated: "2026-06-08"
sources:
  - docs/features/test-pipeline-interleaved/testing/multi-surface-interleaved/journey.md
skip_eval: true
---

# Contract: multi-surface-interleaved / Step 3: Execute the first surface's run-tests task

> **Note**: Contracts generated without eval-journey verification (SKIP_EVAL_GATE=true). Review with extra scrutiny.

<!-- gen-contracts: do not edit manually. Regenerate via /gen-contracts. -->

## Outcome "success"
- Preconditions: "T-test-run-{surface-1} is pending and its dependency T-test-gen-scripts-{surface-1} is completed"
- Input: "Task executor runs T-test-run-{surface-1}"
- Output: "API tests execute successfully. All tests pass. Task completes. The task does NOT wait for gen-scripts of other surfaces"
- State: "T-test-run-{surface-1} status changed to completed. This unblocks T-test-gen-scripts-{surface-2} (interleaved dependency)"
- Side-effect: "Test execution results recorded. No production code modifications"

## Outcome "api-bug-discovered"
- Preconditions: "API tests execute but one or more tests fail due to a real API bug in production code"
- Input: "Task executor runs T-test-run-{surface-1} which reports test failures"
- Output: "Task follows hardened AC: confirms failure is due to production code (not test script bug). Creates fix task via forge task add"
- State: "T-test-run-{surface-1} status set to blocked. Fix task created and linked to source. Second surface's gen-scripts remains blocked until fix resolves"
- Side-effect: "Fix task created via forge task add --type coding.fix --block-source. Production code NOT modified without confirmation"

## Outcome "test-script-bug"
- Preconditions: "Generated test scripts themselves contain a bug (e.g., wrong assertion, incorrect setup)"
- Input: "Task executor runs T-test-run-{surface-1} which reports failures from test script errors"
- Output: "Agent identifies this as a test script bug (not a production code bug). Fixes the test script directly without modifying production code"
- State: "Test scripts corrected. T-test-run-{surface-1} proceeds to completion or retry. No production code changes"
- Side-effect: "Test script files modified to fix the generation error. This does not count as faking the test"

## Outcome "not-found-missing-scripts"
- Preconditions: "Test script files expected by run-tests task are missing or not found"
- Input: "Task executor runs T-test-run-{surface-1} with missing test scripts"
- Output: "Error message indicating test scripts not found. Task blocked"
- State: "Task status set to blocked. No tests executed"
- Side-effect: "none"

## Journey Invariants

- Every generated test task DAG maintains the interleaving invariant: for surface N>0, T-test-gen-scripts-{surface-N} depends on T-test-run-{surface-N-1}, NOT on T-test-gen-scripts-{surface-N-1}
- Every run-tests task includes AC enforcing real tests (no fake/stub-only tests that always pass)
- Every run-tests task enforces the "confirm before modifying production code" rule -- test script bugs may be fixed directly, but production code modification requires explicit confirmation and fix-task creation
- The pipeline does not regress for any number of configured surfaces
- Error handling follows the task-executor's existing Pause Protocol
