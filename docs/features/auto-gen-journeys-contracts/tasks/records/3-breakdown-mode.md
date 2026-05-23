---
status: "completed"
started: "2026-05-24 00:06"
completed: "2026-05-24 00:19"
time_spent: "~13m"
---

# Task Record: 3 autogen.go Breakdown 模式：插入 gen-journeys/gen-contracts 并重写依赖解析

## Summary
Modified GetBreakdownTestTasks() to insert gen-journeys (per-type) and gen-contracts tasks at pipeline head, and rewrote resolveBreakdownDeps() from hardcoded arithmetic indices to findTaskIndex/findTaskIndexByPrefix-based ID lookup with panic on missing tasks.

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/task/autogen.go
- forge-cli/pkg/task/autogen_test.go
- forge-cli/pkg/task/autoconfig_test.go

### Key Decisions
- New pipeline order: gen-journeys-per-type -> eval-journey -> gen-contracts -> eval-contract -> gen-scripts-per-type -> run -> verify-regression
- gen-journeys uses StrategyKind=interface, TypeTestGenJourneys, TestType per interface type
- gen-contracts is a single shared task (not per-type) with TypeTestGenContracts
- Added findTaskIndexOrPanic that panics with missing task ID and all current task IDs for debugging
- ResolveFirstTestDep updated to find gen-journeys as first test task (fallback to eval-journey then gen-scripts)
- eval-journey depends on all gen-journeys tasks (fan-in), matching existing run-depends-on-all-gen-scripts pattern

## Test Results
- **Tests Executed**: Yes
- **Passed**: 10
- **Failed**: 0
- **Coverage**: 90.6%

## Acceptance Criteria
- [x] GetBreakdownTestTasks() generates T-test-gen-journeys-{type} per interface type with correct Type/TestType/StrategyKind
- [x] GetBreakdownTestTasks() generates T-test-gen-contracts with Type=TypeTestGenContracts
- [x] gen-journeys before eval-journey, gen-contracts after eval-journey and before eval-contract
- [x] New tasks use embed template via autogenTypeToFile mapping
- [x] resolveBreakdownDeps() no longer uses hardcoded arithmetic indices
- [x] All dependencies resolved via findTaskIndex or findTaskIndexByPrefix
- [x] Full dependency chain: gen-journeys -> eval-journey -> gen-contracts -> eval-contract -> gen-scripts -> run -> verify-regression
- [x] findTaskIndexOrPanic panics with missing task ID prefix and all current task IDs
- [x] All existing Breakdown mode tests pass after structural changes

## Notes
Breaking change: GetBreakdownTestTasks return value structure changed (task count increased by 2 for single-type, more for multi-type). Updated 4 existing tests and 1 autoconfig test to match new structure.
