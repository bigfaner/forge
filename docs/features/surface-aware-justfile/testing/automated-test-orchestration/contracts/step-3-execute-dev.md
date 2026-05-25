---
journey: "automated-test-orchestration"
step: 3
step-action: "execute dev (background start)"
generated: "2026-05-26"
sources:
  - docs/features/surface-aware-justfile/testing/automated-test-orchestration/journey.md
---

# Contract: automated-test-orchestration / Step 3: Execute dev (background start)

<!-- gen-contracts: do not edit manually. Regenerate via /gen-contracts. -->

## Outcome "success"
- Preconditions: "execution strategy loaded, web or api surface type, dev recipe exists in justfile, port available for dev server"
- Input: "run-tests executes just dev (or surface-key-prefixed dev recipe)"
- Output: "dev server starts in the background, no blocking of run-tests process"
- State: "dev server process running, PID recorded in .forge/test-state.json"
- Side-effect: "background process spawned, .forge/test-state.json updated with PID"

## Outcome "dev-failure"
<!-- source: inferred -->
<!-- reasoning: Journey edge case 6b describes dev server crashing during probe. Dev step itself can also fail at startup (dependency missing, port occupied) or crash immediately. Tech design interface 2 defines exit code 1 for dev as startup failure. Fact Table indicates PID-based process tracking via .forge/test-state.json. -->
- Preconditions: "dev recipe execution fails (exit code 1 or crashes) due to dependency missing, port occupied, configuration error, or immediate crash"
- Input: "run-tests executes just dev, process exits with non-zero code or terminates unexpectedly"
- Output: "error output from dev recipe on stderr indicating startup failure reason, or process termination detected"
- State: "no valid background process to track, .forge/test-state.json not updated or reflects failed dev step"
- Side-effect: "run-tests proceeds to teardown and exits with appropriate exit code (1 for retryable, 2 for blocking)"

## Journey Invariants

- HARD-GATE: After probe failure, no probe retry or dev restart within the same orchestration cycle
- Every error message includes both the specific failure reason and a recovery hint
- Teardown is always idempotent: PID not existing is not an error; kill failure retries once then logs and continues
- Exit code semantics are consistent: 0=success, 1=retryable, 2=blocking
- .forge/test-state.json is always cleaned up regardless of orchestration outcome
- No orphan processes remain after run-tests completes
