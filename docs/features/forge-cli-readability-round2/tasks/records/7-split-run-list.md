---
status: "completed"
started: "2026-06-06 17:24"
completed: "2026-06-06 17:34"
time_spent: "~10m"
---

# Task Record: 7 拆分 runList 217 行函数

## Summary
Extracted 217-line runList into 11 named sub-functions (parseSortMode, resolveFeatureArgs, loadTaskIndex, checkLegacyScope, isTreeMode, handleTreeMode, sortTaskIDs, newDisplayIDFunc, computeColumnWidths, boolStr, printTaskTable). Also extracted sorting logic into list_sort.go. All functions <= 80 lines, nesting <= 4 levels, file <= 500 lines.

## Changes

### Files Created
- forge-cli/internal/cmd/task/list_sort.go

### Files Modified
- forge-cli/internal/cmd/task/list.go

### Key Decisions
- Extracted colWidths struct to pass column widths instead of 6 individual return values
- handleTreeMode returns (bool, error) to distinguish TUI-launched from fallback-to-table
- loadTaskIndex uses sentinel error pattern to preserve nil-return-on-no-tasks behavior
- Moved naturalSortTaskIDs, sortKey, idSortKey, fallbackSortPriority to list_sort.go to keep list.go under 500 lines

## Test Results
- **Tests Executed**: Yes
- **Passed**: 172
- **Failed**: 0
- **Coverage**: 75.6%

## Acceptance Criteria
- [x] runList and all extracted sub-functions <= 80 lines
- [x] All function nesting <= 4 levels
- [x] go test ./... all green, zero behavior change
- [x] File <= 500 lines

## Notes
无
