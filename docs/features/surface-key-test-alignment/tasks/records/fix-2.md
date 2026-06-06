---
status: "completed"
started: "2026-06-06 14:57"
completed: "2026-06-06 15:05"
time_spent: "~8m"
---

# Task Record: fix-2 fix unit-test: just unit-test failure in quality gate

## Summary
Fixed flaky TestPersistentPreRun_InitsWithProjectRoot test by clearing CLAUDE_PROJECT_DIR and PROJECT_ROOT env vars in all three PersistentPreRun tests. Root cause: GetProjectRootFromEnv() returns env var path instead of temp dir when these are set by CI/hook environments.

## Changes

### Files Created
无

### Files Modified
- forge-cli/internal/cmd/root_test.go

### Key Decisions
- Added t.Setenv to clear CLAUDE_PROJECT_DIR and PROJECT_ROOT in 3 tests that depend on FindProjectRoot resolving from temp dir

## Test Results
- **Tests Executed**: Yes
- **Passed**: 3
- **Failed**: 0
- **Coverage**: 0.0%

## Acceptance Criteria
- [x] Root cause identified: env var override in GetProjectRootFromEnv
- [x] Fix applied to all 3 affected tests
- [x] Tests pass with CLAUDE_PROJECT_DIR set

## Notes
The flaky test was not related to surface-key alignment changes. It was a pre-existing Windows-specific issue exposed by quality-gate hook setting CLAUDE_PROJECT_DIR.
