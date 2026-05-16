---
status: "completed"
started: "2026-05-16 21:52"
completed: "2026-05-16 22:10"
time_spent: "~18m"
---

# Task Record: fix-2 fix test-e2e: just test-e2e failure in quality gate

## Summary
Fix e2e test helpers to use local forge binary (forgeBin) and override CLAUDE_PROJECT_DIR in setup functions, preventing stale-binary and wrong-project-root issues in quality-gate hook context.

## Changes

### Files Created
无

### Files Modified
- tests/e2e/task_lifecycle_hardening_cli_test.go
- tests/e2e/test_scripts_per_type_cli_test.go
- tests/e2e/quick_test_slim_cli_test.go

### Key Decisions
- Used forgeBin(t) in tlhForgeClaim instead of system forge binary
- Added t.Setenv in both setupFeatureProject and quickSlimSetupProject to isolate test project roots from inherited env

## Test Results
- **Tests Executed**: Yes
- **Passed**: 69
- **Failed**: 0
- **Coverage**: 0.0%

## Acceptance Criteria
- [x] All e2e tests pass including task-lifecycle-hardening feature tests
- [x] Tests pass when CLAUDE_PROJECT_DIR is set to real project root
- [x] Feature e2e tests use locally built binary with latest changes

## Notes
无
