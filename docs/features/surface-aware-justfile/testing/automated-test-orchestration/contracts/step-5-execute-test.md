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

## Outcome "test-validation-error"
<!-- source: inferred -->
<!-- reasoning: Web surface rule mandates validation-error. Test recipe may have configuration errors (invalid test command, missing test framework setup). This is the validation-error analog for the test execution step. -->
<!-- required_outcomes: web-validation-error -->
- Preconditions: "test recipe configuration is invalid (missing test command, malformed test configuration, or test framework not properly initialized)"
- Input: "run-tests executes just test but the test recipe itself is misconfigured"
- Output: "immediate test failure with descriptive error indicating the configuration issue, with recovery hint to check the test recipe and test framework setup"
- State: "no test execution, teardown executed"
- Side-effect: "process exits with exit code 2 (blocking), teardown executed"

## Outcome "test-execution-timeout"
<!-- source: inferred -->
<!-- reasoning: TUI surface rule mandates timeout Outcome for async operations. Test execution may hang (deadlock, infinite loop in tests). This is the timeout analog for the test execution step. -->
<!-- required_outcomes: tui-timeout -->
- Preconditions: "test execution has been running beyond the expected completion time"
- Input: "run-tests executes just test and the process does not complete within the test timeout period"
- Output: "timeout error indicating test execution exceeded the expected time limit, with elapsed time and recovery hint (check for deadlocked tests or resource contention)"
- State: "test process may still be running, teardown executed to clean up"
- Side-effect: "test process terminated, teardown executed"

## Journey Invariants

- HARD-GATE: After probe failure, no probe retry or dev restart within the same orchestration cycle
- Every error message includes both the specific failure reason and a recovery hint
- Teardown is always idempotent: PID not existing is not an error; kill failure retries once then logs and continues
- Exit code semantics are consistent: 0=success, 1=retryable, 2=blocking
- .forge/test-state.json is always cleaned up regardless of orchestration outcome
- No orphan processes remain after run-tests completes
- Step-specific: test execution always follows successful probe (or no-probe surface type), never skipped
