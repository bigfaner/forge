---
journey: "automated-test-orchestration"
step: 4
step-action: "execute probe (retry polling)"
generated: "2026-05-26"
sources:
  - docs/features/surface-aware-justfile/testing/automated-test-orchestration/journey.md
---

# Contract: automated-test-orchestration / Step 4: Execute probe (retry polling)

<!-- gen-contracts: do not edit manually. Regenerate via /gen-contracts. -->

## Outcome "success"
- Preconditions: "dev server is running in background (PID valid), service is becoming ready"
- Input: "run-tests executes just probe with up to 3 retries at 30-second intervals (total timeout 90s)"
- Output: "probe checks service readiness, logs format: [probe] [retry current/max] url -- reason. On success (exit 0), proceeds to test step"
- State: "service confirmed ready, .forge/test-state.json retains PID for later teardown"
- Side-effect: "none"

## Outcome "probe-retryable-failure"
- Preconditions: "dev server started but probe detects transient errors (connection timeout, DNS failure) across all retry attempts"
- Input: "run-tests executes just probe, all 3 consecutive probe attempts fail with retryable error (exit 1)"
- Output: "after 3 consecutive probe failures, run-tests executes teardown, then exits with exit code 1 (retryable)"
- State: "HARD-GATE enforced: no probe retry or dev restart within the same orchestration cycle"
- Side-effect: "teardown executed, .forge/test-state.json cleaned up"

## Outcome "probe-blocking-failure"
- Preconditions: "probe detects a non-transient error (port occupied, auth failure, config error)"
- Input: "run-tests executes just probe, probe fails immediately with blocking error (exit 2)"
- Output: "run-tests executes teardown, then exits with exit code 2 (blocking). HARD-GATE enforced: no retry within this cycle"
- State: "no retry attempted, teardown executed"
- Side-effect: "teardown executed, .forge/test-state.json cleaned up, upper-level scheduler distinguishes retryable(1) vs blocking(2)"

## Outcome "dev-crash-during-probe"
- Preconditions: "dev server process terminates unexpectedly while probe is retrying"
- Input: "probe retries continue but dev server PID is no longer valid"
- Output: "probe retries exhaust (3 attempts, 90s total). run-tests executes teardown (PID already gone, idempotent skip)"
- State: "no orphan process remains, .forge/test-state.json cleaned up"
- Side-effect: "run-tests exits with exit code 1 (retryable)"

## Outcome "probe-validation-error"
<!-- source: inferred -->
<!-- reasoning: Web surface rule mandates validation-error for user-facing flows. Probe configuration may have invalid parameters (malformed URL, incorrect health check path). Although probe is automated, misconfigured probe is analogous to validation error. -->
<!-- required_outcomes: web-validation-error -->
- Preconditions: "probe configuration is invalid (malformed health check URL, missing or invalid endpoint path in justfile recipe)"
- Input: "run-tests executes just probe with an invalid or misconfigured probe recipe"
- Output: "probe fails immediately with a descriptive error indicating the configuration issue, with recovery hint to check the probe recipe in the justfile"
- State: "no health check performed, teardown executed"
- Side-effect: "process exits with exit code 2 (blocking), teardown executed"

## Journey Invariants

- HARD-GATE: After probe failure, no probe retry or dev restart within the same orchestration cycle
- Every error message includes both the specific failure reason and a recovery hint
- Teardown is always idempotent: PID not existing is not an error; kill failure retries once then logs and continues
- Exit code semantics are consistent: 0=success, 1=retryable, 2=blocking
- .forge/test-state.json is always cleaned up regardless of orchestration outcome
- No orphan processes remain after run-tests completes
- Step-specific: probe retries are bounded at 3 attempts with 30-second intervals, total timeout 90s
