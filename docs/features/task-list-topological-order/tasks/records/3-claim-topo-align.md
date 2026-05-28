---
status: "completed"
started: "2026-05-28 00:39"
completed: "2026-05-28 00:45"
time_spent: "~6m"
---

# Task Record: 3 Align claimNextTask with topological ordering

## Summary
Aligned claimNextTask with topological ordering: replaced priority+version sort with topo-depth-first ordering, keeping priority as tiebreaker within same depth level

## Changes

### Files Created
无

### Files Modified
- forge-cli/internal/cmd/task/claim.go
- forge-cli/internal/cmd/task/claim_test.go

### Key Decisions
- Added computeTopoDepths function using Kahn's algorithm BFS to calculate topological depth for all tasks in index
- Sort eligible tasks by (topo_depth, priority, natural_ID) instead of (priority, natural_ID)
- Topo depth is computed over the full task graph so tasks with unmet deps still get correct depth values

## Test Results
- **Tests Executed**: Yes
- **Passed**: 6
- **Failed**: 0
- **Coverage**: 71.7%

## Acceptance Criteria
- [x] claimNextTask selects eligible tasks in topological order
- [x] Priority still acts as tiebreaker among tasks at same topological level
- [x] Existing claim_test.go tests pass after modification
- [x] New test: verify topo-ordered claiming for multi-dependency task graph
- [x] New test: verify priority tiebreaker within same topological level

## Notes
computeTopoDepths uses task ID (not map key) as lookup key to match the sort comparator. All existing tests pass without modification.
