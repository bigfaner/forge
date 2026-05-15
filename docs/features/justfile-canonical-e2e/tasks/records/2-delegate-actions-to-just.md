---
status: "completed"
started: "2026-05-15 00:38"
completed: "2026-05-15 00:48"
time_spent: "~10m"
---

# Task Record: 2 Delegate e2e/actions.go functions to just recipes

## Summary
Replace profile-specific switch/case dispatch in Run/Setup/Compile/Discover with thin exec.Command('just', ...) delegation. Verify() left completely unchanged. Added just-not-on-PATH detection with actionable error message. Extracted runJust() helper to DRY the just invocation pattern.

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/e2e/actions.go
- forge-cli/pkg/e2e/actions_test.go
- forge-cli/scripts/version.txt

### Key Decisions
- Extracted runJust() helper to centralize just invocation, not-found detection, and error formatting for all four action functions
- Feature filtering uses justfile named-argument syntax: 'just test-e2e feature=<name>' per hard rule requirement
- Removed os/exec import from actions.go since subprocess calls now go through the existing ExecRunner interface in exec.go
- Bumped version to 3.10.0 (minor) since this changes e2e execution behavior for all profiles

## Test Results
- **Tests Executed**: Yes
- **Passed**: 40
- **Failed**: 0
- **Coverage**: 88.9%

## Acceptance Criteria
- [x] Run() calls exec.Command('just', 'test-e2e') with optional feature=<name> argument
- [x] Setup() calls exec.Command('just', 'e2e-setup')
- [x] Compile() calls exec.Command('just', 'e2e-compile')
- [x] Discover() calls exec.Command('just', 'e2e-discover')
- [x] Verify() is completely unchanged
- [x] Non-zero exit from just produces a non-nil error; zero exit produces nil error
- [x] Forge CLI exit code matches the just exit code
- [x] When just is not on PATH, error message: 'just' is required but not found on PATH
- [x] All 6 profiles work without per-profile code paths
- [x] go build ./... passes

## Notes
Behavior change: go-test profile Run() filtering changes from directory-scoped (./tests/e2e/features/{name}/...) to justfile recipe test-name-based filtering. This is documented in the task's Implementation Notes. The pre-existing TestSaveIndexAndSignalCompletion_SaveIndexError failure in internal/cmd is unrelated to this change.
