---
status: "completed"
started: "2026-06-09 17:13"
completed: "2026-06-09 17:18"
time_spent: "~5m"
---

# Task Record: fix-5 fix test: just test failure in quality gate

## Summary
Fixed 5 pre-existing test failures in tests/test-generation/ by updating quick mode pipeline topology expectations to match current forge task index behavior (4 tasks -> 7 tasks with gen-contracts and per-surface gen-scripts).

## Changes

### Files Created
无

### Files Modified
- tests/test-generation/quick_test_slim_test.go
- tests/test-generation/test_scripts_per_type_test.go

### Key Decisions
- Updated test expectations from old 4-task topology (gen-journeys -> run-test -> drift) to new 7-task topology (gen-journeys -> gen-contracts -> gen-scripts-per-surface -> run-test-per-surface -> drift)
- Fixed TC-005 byID lookup keys to use task IDs (T-test-gen-contracts) instead of map keys (gen-contracts)

## Test Results
- **Tests Executed**: Yes
- **Passed**: 1937
- **Failed**: 0
- **Coverage**: 0.0%

## Acceptance Criteria
- [x] TestTC_001_QuickModeSingleProfileTaskCount passes with correct count
- [x] TestTC_005_QuickModeDependencyChainCorrect passes with new topology
- [x] TestTC_006_QuickModeRunTestSerialChainFanIn passes
- [x] TestTC_012_QuickModeSingleProfileProducesCorrectTaskCount passes
- [x] TestPerType_TC_006_TaskIndexRunDependsOnAllPerTypeGenTasks passes
- [x] All other tests still pass (no regression)

## Notes
These were pre-existing failures on main branch. The forge task index command now generates an expanded pipeline (gen-contracts + per-surface gen-scripts) that the tests hadn't been updated for.
