---
status: "completed"
started: "2026-05-16 20:36"
completed: "2026-05-16 20:41"
time_spent: "~5m"
---

# Task Record: 1 Fix checkDependenciesMet: add SourceTaskID == selfID check

## Summary
Added self-block check to checkDependenciesMet: when an active fix-task has SourceTaskID matching the candidate task's own ID, the candidate is blocked from being claimed. This covers the --block-source scenario where a fix task targets the source task itself.

## Changes

### Files Created
无

### Files Modified
- forge-cli/internal/cmd/claim.go
- forge-cli/internal/cmd/claim_test.go

### Key Decisions
- Added a separate loop after the dependency-based fix-task check to scan all tasks for SourceTaskID == selfID matches, rather than merging into the existing loop. This keeps the two concerns (dependency blocking vs self-blocking) clearly separated.
- The self-block check does NOT exclude selfID (other.ID != selfID) since we specifically want to match fix tasks targeting the candidate, not the candidate itself.

## Test Results
- **Tests Executed**: Yes
- **Passed**: 912
- **Failed**: 0
- **Coverage**: 80.8%

## Acceptance Criteria
- [x] checkDependenciesMet returns false when an active fix-task has SourceTaskID == selfID
- [x] Existing behavior unchanged for tasks without active fix-tasks targeting them
- [x] Existing tests continue to pass
- [x] New test cases added for the self-block scenario

## Notes
Hard rules verified: no changes to submit.go; only adds blocking conditions, never relaxes existing checks.
