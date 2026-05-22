---
status: "completed"
started: "2026-05-22 09:13"
completed: "2026-05-22 09:26"
time_spent: "~13m"
---

# Task Record: fix.1 Fix TestRunStatus silent failure in internal/cmd

## Summary
Fix TestRunStatus silent failure by supporting 2-arg mutation mode in task status command using state machine validation

## Changes

### Files Created
无

### Files Modified
- forge-cli/internal/cmd/status.go
- forge-cli/internal/cmd/integration_test.go
- forge-cli/internal/cmd/characterization_test.go
- forge-cli/internal/cmd/reopen_test.go

### Key Decisions
- Status command now supports 2-arg mutation (<task-id> <status>) using state machine validation via ValidateTransition with empty role
- Uses indexPkg.WithLock for concurrent-safe writes, matching the pattern from reopen.go
- Non-terminal, non-completed transitions are allowed (pending -> blocked, in_progress -> pending, etc.)
- Terminal state transitions (from completed/rejected) and completed-target transitions require dedicated commands (submit/reopen)
- Characterization test TestStatus_AllowsMutation updated to reflect correct state machine behavior
- Reopen test TestStatus_ReadOnly_AnyStatusArgument updated to allow valid non-terminal transitions

## Test Results
- **Tests Executed**: Yes
- **Passed**: 702
- **Failed**: 0
- **Coverage**: 81.0%

## Acceptance Criteria
- [x] go test ./internal/cmd/ -run TestRunStatus passes with clear output
- [x] go test ./internal/cmd/... all pass (0 failures)
- [x] No other test regressions across go test ./...

## Notes
无
