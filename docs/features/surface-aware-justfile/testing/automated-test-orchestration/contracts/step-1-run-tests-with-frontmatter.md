---
journey: "automated-test-orchestration"
step: 1
step-action: "run run-tests with task frontmatter surface-type"
generated: "2026-05-26"
sources:
  - docs/features/surface-aware-justfile/testing/automated-test-orchestration/journey.md
---

# Contract: automated-test-orchestration / Step 1: Run run-tests with task frontmatter surface-type

<!-- gen-contracts: do not edit manually. Regenerate via /gen-contracts. -->

## Outcome "success"
- Preconditions: "task file exists with surface-type: web (or other valid type) in its frontmatter, project justfile contains surface-specific recipes"
- Input: "user executes run-tests for a task that has surface-type in its frontmatter"
- Output: "run-tests reads the surface-type from the task frontmatter and proceeds to load the corresponding execution strategy rule file. No fallback to forge surfaces CLI is needed"
- State: "surface type determined from frontmatter, ready to load execution strategy"
- Side-effect: "none"

## Outcome "frontmatter-missing-fallback-cli"
- Preconditions: "task frontmatter has no surface-type field, forge surfaces CLI is available and functional"
- Input: "user executes run-tests for a task without surface-type in frontmatter"
- Output: "run-tests falls back to forge surfaces CLI with the file path to query surface type via longest-prefix-match. On success, proceeds to load rule file"
- State: "surface type determined from forge surfaces CLI fallback, ready to load execution strategy"
- Side-effect: "forge surfaces CLI invoked as subprocess"

## Outcome "surface-type-unavailable"
- Preconditions: "task frontmatter has no surface-type field AND forge surfaces CLI fails, returns no result, or returns malformed output (not valid JSON or missing expected fields)"
- Input: "user executes run-tests with no usable surface-type from either source"
- Output: "error message to stderr with the failure reason and a recovery hint (configure surfaces in .forge/config.yaml or specify surface-type in task frontmatter, or verify forge CLI version)"
- State: "no surface type determined, run-tests aborted"
- Side-effect: "process exits with exit code 2 (blocking) when both sources unavailable, or exit code 1 (retryable) when CLI output is malformed"

## Outcome "session-expired-during-detection"
<!-- source: inferred -->
<!-- reasoning: Web surface rule mandates session-expired for authenticated flows. run-tests may invoke forge CLI which requires valid project context. If the project configuration was invalidated mid-session (config.yaml deleted or corrupted), surface detection fails with a session-like expiry. -->
<!-- required_outcomes: web-session-expired -->
- Preconditions: "project configuration was invalidated after run-tests started (config.yaml deleted, corrupted, or permissions changed)"
- Input: "run-tests attempts to read surface-type but project context is no longer valid"
- Output: "error message indicating project context is no longer valid, with recovery hint to verify project configuration and re-initialize if needed"
- State: "no surface type determined, run-tests aborted"
- Side-effect: "process exits with exit code 2 (blocking)"

## Journey Invariants

- HARD-GATE: After probe failure, no probe retry or dev restart within the same orchestration cycle
- Every error message includes both the specific failure reason and a recovery hint
- Teardown is always idempotent: PID not existing is not an error; kill failure retries once then logs and continues
- Exit code semantics are consistent: 0=success, 1=retryable, 2=blocking
- .forge/test-state.json is always cleaned up regardless of orchestration outcome
- No orphan processes remain after run-tests completes
- Step-specific: surface-type resolution from frontmatter takes priority over forge surfaces CLI fallback
