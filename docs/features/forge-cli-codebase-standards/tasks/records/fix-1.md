---
status: "completed"
started: "2026-05-31 10:29"
completed: "2026-05-31 10:42"
time_spent: "~13m"
---

# Task Record: fix-1 fix test: just test failure in quality gate

## Summary
Fixed two pre-existing test failures in forge task index: TC_009 (docs-only features getting test pipeline tasks) and TC_010 (idempotent re-run rejecting auto-gen tasks with system types).

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/task/build.go
- tests/task-type-system/task_type_refinement_test.go

### Key Decisions
- Docs-only features (needsEval && !needsTest) now pass nil surfaces to GenerateTestTasks, suppressing surface-dependent test tasks while keeping T-review-doc
- Added isAutoGenRegistryPrefix() to recognize per-surface-key auto-gen IDs without surfaces map, fixing system-type validation on re-run
- Fixed TC_009 test assertion: changed T-quick- prefix check to T-quick-run- to avoid catching T-quick-doc-drift (a doc task, not test pipeline)

## Test Results
- **Tests Executed**: Yes
- **Passed**: 42
- **Failed**: 0
- **Coverage**: 0.0%

## Acceptance Criteria
- [x] TC_009 passes: docs-only features get T-review-doc but not test pipeline tasks
- [x] TC_010 passes: forge task index is idempotent on second run
- [x] All 9 integration test suites pass

## Notes
无
