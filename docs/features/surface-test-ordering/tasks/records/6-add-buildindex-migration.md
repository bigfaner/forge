---
status: "completed"
started: "2026-05-26 13:42"
completed: "2026-05-26 14:08"
time_spent: "~26m"
---

# Task Record: 6 Add BuildIndex migration for existing fix-tasks

## Summary
Added BuildIndex migration logic for fix-tasks with SourceTaskID 'T-test-run' in multi-surface projects. Migration remaps fix-tasks to T-test-run-{first-surface-key}, copies status/blockedReason to the first run-test task, and removes the old T-test-run entry. Single-surface projects are unaffected.

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/task/build.go
- forge-cli/pkg/task/build_test.go

### Key Decisions
- Two-phase migration: phase 1 (step 5.9) detects and saves old T-test-run state before orphan cleanup, phase 2 (step 7.5) applies saved state to first run-test task after test pipeline generation
- Migration runs before orphan cleanup (step 6) so old T-test-run entry is still available for state extraction
- Only applies to multi-surface projects (isSingleSurface guard); single-surface projects skip migration entirely

## Test Results
- **Tests Executed**: Yes
- **Passed**: 6
- **Failed**: 0
- **Coverage**: 87.7%

## Acceptance Criteria
- [x] Multi-surface fix-tasks with SourceTaskID 'T-test-run' auto-remap to T-test-run-{first-surface-key}
- [x] Single-surface project (surfaces: api) has no migration behavior, T-test-run unchanged
- [x] Old T-test-run entry removed from index.json after migration

## Notes
All 30 packages pass with zero regressions. Full coverage at 87.7% for pkg/task.
