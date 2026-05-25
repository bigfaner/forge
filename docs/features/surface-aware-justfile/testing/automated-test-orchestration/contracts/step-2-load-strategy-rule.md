---
journey: "automated-test-orchestration"
step: 2
step-action: "load execution strategy rule file"
generated: "2026-05-26"
sources:
  - docs/features/surface-aware-justfile/testing/automated-test-orchestration/journey.md
---

# Contract: automated-test-orchestration / Step 2: Load execution strategy rule file

<!-- gen-contracts: do not edit manually. Regenerate via /gen-contracts. -->

## Outcome "success"
- Preconditions: "surface type is detected from step 1 (frontmatter or forge surfaces CLI fallback)"
- Input: "run-tests loads the execution strategy rule file for the detected surface type"
- Output: "rule file loaded successfully, e.g., skills/run-tests/rules/surfaces/web.md defines orchestration sequence: dev (background) -> probe -> test -> teardown"
- State: "execution strategy loaded, orchestration sequence defined and ready to execute"
- Side-effect: "none"

## Outcome "rule-file-not-found"
- Preconditions: "surface type is detected but the corresponding rule file does not exist"
- Input: "run-tests attempts to load rule file for detected surface type"
- Output: "error message to stderr: surface type with type value execution strategy rule file not found, supported types: web/api/cli/tui/mobile"
- State: "no execution strategy loaded, run-tests aborted"
- Side-effect: "process exits with exit code 2 (blocking)"

## Outcome "not-found-surface-type"
<!-- source: cli-required — surface rule mandates not-found for resource access steps -->
- Preconditions: "forge surfaces CLI returns no match for the queried path"
- Input: "run-tests attempts to determine surface type via forge surfaces CLI for a path that matches no configured surface entry"
- Output: "forge surfaces CLI exits with code 1, stderr contains error message with recovery hint (run forge init to configure surfaces)"
- State: "surface type unavailable, run-tests reports error and aborts"
- Side-effect: "process exits with exit code 2 (blocking)"

## Journey Invariants

- HARD-GATE: After probe failure, no probe retry or dev restart within the same orchestration cycle
- Every error message includes both the specific failure reason and a recovery hint
- Teardown is always idempotent: PID not existing is not an error; kill failure retries once then logs and continues
- Exit code semantics are consistent: 0=success, 1=retryable, 2=blocking
- .forge/test-state.json is always cleaned up regardless of orchestration outcome
- No orphan processes remain after run-tests completes
