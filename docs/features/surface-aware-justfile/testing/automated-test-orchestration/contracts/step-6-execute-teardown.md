---
journey: "automated-test-orchestration"
step: 6
step-action: "execute teardown"
generated: "2026-05-26"
sources:
  - docs/features/surface-aware-justfile/testing/automated-test-orchestration/journey.md
---

# Contract: automated-test-orchestration / Step 6: Execute teardown

<!-- gen-contracts: do not edit manually. Regenerate via /gen-contracts. -->

## Outcome "success"
- Preconditions: "test step completed (pass or fail), background processes tracked in .forge/test-state.json"
- Input: "run-tests executes just test-teardown"
- Output: "background processes terminated (PID-based), .forge/test-state.json cleaned up (deleted or marked completed)"
- State: "all tracked processes terminated, .forge/test-state.json removed or finalized"
- Side-effect: "process termination via PID-based kill, .forge/test-state.json deleted"

## Outcome "teardown-kill-failure"
- Preconditions: "teardown attempts to kill a PID but the kill syscall fails"
- Input: "run-tests executes just test-teardown, kill syscall returns error"
- Output: "teardown retries the kill once. If it still fails, logs the process info and continues cleanup"
- State: ".forge/test-state.json is cleaned up (deleted or marked completed) regardless of kill outcome"
- Side-effect: "failed kill logged with process details for manual investigation"

## Outcome "stale-state-cleanup"
<!-- source: cli-required (already-exists) + journey edge case 8b — merged stale state and conflicting state scenarios -->
- Preconditions: "a previous run-tests was interrupted or left .forge/test-state.json with active/stale state (e.g., user Ctrl+C, incomplete previous run)"
- Input: "run-tests detects existing .forge/test-state.json from a previous incomplete or interrupted run"
- Output: "stale state detected and cleaned up, residual processes terminated, previous state info logged. Fresh orchestration cycle begins"
- State: "old state file removed, new clean orchestration cycle initialized"
- Side-effect: "residual processes from previous run are terminated before new cycle starts"

## Journey Invariants

- HARD-GATE: After probe failure, no probe retry or dev restart within the same orchestration cycle
- Every error message includes both the specific failure reason and a recovery hint
- Teardown is always idempotent: PID not existing is not an error; kill failure retries once then logs and continues
- Exit code semantics are consistent: 0=success, 1=retryable, 2=blocking
- .forge/test-state.json is always cleaned up regardless of orchestration outcome
- No orphan processes remain after run-tests completes
