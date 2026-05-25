---
journey: "automated-test-orchestration"
step: 5
step-action: "execute test"
generated: "2026-05-26"
sources:
  - docs/features/surface-aware-justfile/testing/automated-test-orchestration/journey.md
---

# Contract: automated-test-orchestration / Step 5: Execute test

<!-- gen-contracts: do not edit manually. Regenerate via /gen-contracts. -->

## Outcome "success"
- Preconditions: "probe completed successfully (or surface type has no probe step), test recipe exists in justfile"
- Input: "run-tests executes just test"
- Output: "test suite runs, all tests pass (exit code 0)"
- State: "tests completed successfully, proceeds to teardown"
- Side-effect: "none"

## Outcome "test-failure"
- Preconditions: "test recipe exists but test execution fails (exit code 1 for retryable or exit code 2 for blocking)"
- Input: "run-tests executes just test, process exits with non-zero code"
- Output: "test output shows failure details. run-tests executes teardown and exits with the same exit code (1 for retryable, 2 for blocking). For exit code 2, prompts user that test environment is abnormal and suggests retry"
- State: "tests failed, teardown executed for cleanup, .forge/test-state.json cleaned up"
- Side-effect: "teardown executed, exit code propagated to caller"

## Journey Invariants

- HARD-GATE: After probe failure, no probe retry or dev restart within the same orchestration cycle
- Every error message includes both the specific failure reason and a recovery hint
- Teardown is always idempotent: PID not existing is not an error; kill failure retries once then logs and continues
- Exit code semantics are consistent: 0=success, 1=retryable, 2=blocking
- .forge/test-state.json is always cleaned up regardless of orchestration outcome
- No orphan processes remain after run-tests completes
