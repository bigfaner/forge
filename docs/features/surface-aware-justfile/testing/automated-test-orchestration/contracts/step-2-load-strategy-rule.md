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
- Preconditions: "surface type is detected but the corresponding execution strategy rule file does not exist on disk"
- Input: "run-tests attempts to load the rule file for the detected surface type but the file is missing"
- Output: "error message to stderr: execution strategy rule file not found for the detected surface type, listing supported types (web/api/cli/tui/mobile)"
- State: "no execution strategy loaded, run-tests aborted"
- Side-effect: "process exits with exit code 2 (blocking)"

## Outcome "rule-file-malformed"
<!-- source: inferred -->
<!-- reasoning: Surface rule file may exist but contain invalid structure (missing orchestration sequence, unknown step names). This is the not-found analog for the rule loading step — resource exists but is unusable. -->
<!-- required_outcomes: cli-not-found -->
- Preconditions: "execution strategy rule file exists but has invalid structure (missing required sections, unparseable content, or unknown step definitions)"
- Input: "run-tests loads the rule file but cannot parse it into a valid execution strategy"
- Output: "error message to stderr indicating the rule file is malformed with specific parsing failure, recovery hint to check the rule file format or regenerate it"
- State: "no valid execution strategy loaded, run-tests aborted"
- Side-effect: "process exits with exit code 2 (blocking)"

## Journey Invariants

- HARD-GATE: After probe failure, no probe retry or dev restart within the same orchestration cycle
- Every error message includes both the specific failure reason and a recovery hint
- Teardown is always idempotent: PID not existing is not an error; kill failure retries once then logs and continues
- Exit code semantics are consistent: 0=success, 1=retryable, 2=blocking
- .forge/test-state.json is always cleaned up regardless of orchestration outcome
- No orphan processes remain after run-tests completes
