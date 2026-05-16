---
status: "completed"
started: "2026-05-16 20:42"
completed: "2026-05-16 20:54"
time_spent: "~12m"
---

# Task Record: 2 Add lazy unblock scan in claimNextTask

## Summary
Add lazy unblock scan at the start of claimNextTask that iterates all blocked tasks, checks if their dependencies are now met via checkDependenciesMet, and auto-transitions eligible ones to pending. The scan runs before the existing hasPending check so newly-unblocked tasks are visible for claiming. Auto-unblocked tasks are logged with fmt.Printf. Tasks with active fix tasks targeting them (self-block rule) remain blocked.

## Changes

### Files Created
无

### Files Modified
- forge-cli/internal/cmd/claim.go
- forge-cli/internal/cmd/claim_test.go
- forge-cli/internal/cmd/integration_test.go

### Key Decisions
- Scan is O(n) where n = tasks in the feature (typically <20), no performance concern
- Lazy scan reuses existing checkDependenciesMet which already includes selfID fix and fix-task blocking logic
- Auto-unblock is a claim-time pull pattern replacing fragile submit-time cascade

## Test Results
- **Tests Executed**: Yes
- **Passed**: 7
- **Failed**: 0
- **Coverage**: 80.8%

## Acceptance Criteria
- [x] Blocked task auto-transitions to pending when checkDependenciesMet returns true on claim
- [x] Auto-unblocked tasks are logged
- [x] The scan runs before the existing hasPending check so newly-unblocked tasks are visible
- [x] Fix-task in progress targeting the task keeps it blocked
- [x] Existing claim behavior unchanged when no blocked tasks exist

## Notes
Updated TestRunRecord_AutoDowngrade_ThenClaim to give the downgraded task an unmet dependency so it stays blocked under the new lazy scan. Updated TestExecuteClaim_BlockedTaskClearsStateAndProceeds to expect auto-unblock + claim since the blocked task has no deps.
