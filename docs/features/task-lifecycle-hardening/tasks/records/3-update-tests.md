---
status: "completed"
started: "2026-05-16 20:55"
completed: "2026-05-16 21:03"
time_spent: "~8m"
---

# Task Record: 3 Update tests for lazy unblock and self-block scenarios

## Summary
Added table-driven tests for block-source lifecycle (fix done -> claim -> source auto-unblocked) and auto-downgrade unblock scenarios (task blocked -> dep completed -> claim auto-unblocks) to claim_test.go, covering all 6 acceptance criteria from the task definition.

## Changes

### Files Created
无

### Files Modified
- forge-cli/internal/cmd/claim_test.go

### Key Decisions
- Used table-driven subtests following existing patterns (TestClaimNextTask_BlockSourceLifecycle, TestClaimNextTask_AutoDowngradeUnblock)
- Tested both positive (auto-unblock succeeds) and negative (stays blocked) paths for each scenario
- Block-source lifecycle tests cover: fix completed unblocks source, fix active keeps source blocked, multiple fixes all completed, unmet regular deps override fix completion
- Auto-downgrade tests cover: dep completes unblocks, dep pending keeps blocked, no-deps task auto-unblocks immediately

## Test Results
- **Tests Executed**: Yes
- **Passed**: 137
- **Failed**: 0
- **Coverage**: 80.8%

## Acceptance Criteria
- [x] Test: checkDependenciesMet returns false when active fix-task has SourceTaskID == selfID
- [x] Test: checkDependenciesMet returns true when fix-task targeting self is completed
- [x] Test: claimNextTask auto-unblocks blocked task whose dependency just completed
- [x] Test: blocked task stays blocked when fix-task targeting it is still active
- [x] Test: --block-source scenario (fix done -> claim -> source auto-unblocked)
- [x] Test: auto-downgrade scenario (task blocked -> submit passes -> claim auto-unblocks)
- [x] All existing tests pass with updated expectations

## Notes
7 new subtests added across 2 test functions. All existing 130 tests continue to pass. Coverage at 80.8% for internal/cmd package.
