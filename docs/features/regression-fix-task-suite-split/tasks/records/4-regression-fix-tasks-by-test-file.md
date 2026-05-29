---
status: "completed"
started: "2026-05-29 14:43"
completed: "2026-05-29 14:56"
time_spent: "~13m"
---

# Task Record: 4 实现 addRegressionFixTasks 按测试文件拆分 fix task

## Summary
Implemented addRegressionFixTasks function that splits regression fix tasks by test file, with soft cap of 10+1 overflow, fallback to addFixTask when no test files found, and replaced addFixTask calls in runTestRegressionLegacy and runTestRegressionSurface

## Changes

### Files Created
无

### Files Modified
- forge-cli/internal/cmd/quality_gate.go
- forge-cli/internal/cmd/quality_gate_test.go

### Key Decisions
- Used createFixTask with fixTaskOverride to customize title and description per test file
- addRegressionFixTasks bypasses maxFixTasksPerStep cap by calling createFixTask directly instead of addSingleFixTask
- Fallback to addFixTask when extractFileLineMap returns empty map, with structured log warning

## Test Results
- **Tests Executed**: Yes
- **Passed**: 13
- **Failed**: 0
- **Coverage**: 91.7%

## Acceptance Criteria
- [x] isTestFile and addRegressionFixTasks created, uses extractFileLineMap and createFixTask
- [x] Fix task title format 'fix test: <filename> failure in quality gate' compatible with countFixTasks prefix
- [x] Each fix task description contains test file path and filtered output lines
- [x] Soft cap 10 independent + 1 overflow, total <= 11, not subject to maxFixTasksPerStep
- [x] runTestRegressionLegacy and runTestRegressionSurface call addRegressionFixTasks with fallback warning

## Notes
13 new test cases added covering single file, multiple files, overflow, fallback, empty output, title format, task properties, cap bypass, and description content. Full test suite (go test -count=1 ./internal/cmd/) passes.
