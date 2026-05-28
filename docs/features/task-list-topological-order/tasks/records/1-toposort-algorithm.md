---
status: "completed"
started: "2026-05-28 00:21"
completed: "2026-05-28 00:26"
time_spent: "~5m"
---

# Task Record: 1 Implement Kahn's algorithm topological sort

## Summary
Implemented Kahn's algorithm topological sort in TopologicalSort() with wildcard expansion, cycle detection, and missing dep identification. All 16 tests pass with 89.1% package coverage.

## Changes

### Files Created
- forge-cli/pkg/task/toposort.go
- forge-cli/pkg/task/toposort_test.go

### Files Modified
无

### Key Decisions
- Used min-heap (idQueue) with CompareVersionIDs for deterministic natural ID ordering at same topological level
- Self-referencing deps create self-edges so Kahn's naturally detects them as cycle members (in-degree never reaches 0)
- Wildcard expansion reuses ResolveWildcardDep; results deduplicated against exact IDs via seen-set before building adjacency list
- Missing deps collected as a set to avoid duplicates from multiple tasks referencing the same missing ID

## Test Results
- **Tests Executed**: Yes
- **Passed**: 16
- **Failed**: 0
- **Coverage**: 89.1%

## Acceptance Criteria
- [x] TopologicalSort(tasks TaskIndex) ([]string, []string, []string) returns ordered IDs, cycle IDs, and missing dep IDs
- [x] Kahn's algorithm produces correct ordering: if B depends on A, A appears before B
- [x] Same-level tasks (no dependency relationship) sorted by natural ID for determinism
- [x] Wildcard deps (1.x) expanded via ResolveWildcardDep before building adjacency list; expansion is stable across runs
- [x] Cycle detection: remaining nodes after Kahn's completes are identified as cycle members
- [x] Missing deps: dependencies referencing non-existent task IDs are identified
- [x] Empty task list returns empty results (no panic)
- [x] Single task with no deps returns [that-task-id]
- [x] Self-referencing dependency detected as cycle
- [x] Test coverage for: linear chain, diamond deps, disconnected components, wildcards, mixed exact+wildcard deps, cycle scenarios

## Notes
All Hard Rules satisfied: ResolveWildcardDep reused (not reimplemented), O(V+E) via Kahn's + min-heap, no side effects on input TaskIndex.
