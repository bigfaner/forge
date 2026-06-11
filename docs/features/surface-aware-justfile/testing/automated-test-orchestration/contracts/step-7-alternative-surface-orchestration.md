---
journey: "automated-test-orchestration"
step: 7
step-action: "alternative surface orchestration (CLI/TUI/Mobile)"
generated: "2026-05-26"
sources:
  - docs/features/surface-aware-justfile/testing/automated-test-orchestration/journey.md
---

# Contract: automated-test-orchestration / Step 7: Alternative surface orchestration (CLI/TUI/Mobile)

<!-- gen-contracts: do not edit manually. Regenerate via /gen-contracts. -->

## Outcome "cli-tui-success"
- Preconditions: "task has surface-type: cli or surface-type: tui in frontmatter"
- Input: "user executes run-tests for a CLI/TUI surface task"
- Output: "orchestration sequence is build -> dev -> test. No probe step, no teardown step. Dev is not a background server but a build/compile step"
- State: "tests run directly after build, no background server management, no state file with PID tracking"
- Side-effect: "none"

## Outcome "mobile-success"
- Preconditions: "task has surface-type: mobile in frontmatter, emulator tooling available"
- Input: "user executes run-tests for a mobile surface task"
- Output: "orchestration sequence is test-setup -> dev -> test -> teardown. test-setup prepares the emulator. teardown cleans up emulator and processes"
- State: "emulator launched, tests run against emulator, emulator and processes cleaned up"
- Side-effect: "emulator process spawned and later terminated, .forge/test-state.json tracks emulator and dev server PIDs"

## Outcome "alternative-surface-failure"
<!-- source: inferred -->
<!-- reasoning: Alternative surface types (CLI/TUI/Mobile) can fail at build/setup step. CLI/TUI build failures block test execution. Mobile emulator setup failures are a common failure point (SDK missing, hardware acceleration unavailable). -->
- Preconditions: "alternative surface type detected but build or setup step fails (compilation error for CLI/TUI, emulator failure for mobile)"
- Input: "run-tests executes build or test-setup step and it fails"
- Output: "build/setup failure output on stderr, run-tests skips test step and exits with exit code 1 (retryable) or 2 (blocking)"
- State: "no test execution, no background processes to clean up (or partial emulator resources for mobile)"
- Side-effect: "none for CLI/TUI; partial emulator cleanup may be needed for mobile"

## Journey Invariants

- HARD-GATE: After probe failure, no probe retry or dev restart within the same orchestration cycle
- Every error message includes both the specific failure reason and a recovery hint
- Teardown is always idempotent: PID not existing is not an error; kill failure retries once then logs and continues
- Exit code semantics are consistent: 0=success, 1=retryable, 2=blocking
- .forge/test-state.json is always cleaned up regardless of orchestration outcome
- No orphan processes remain after run-tests completes
