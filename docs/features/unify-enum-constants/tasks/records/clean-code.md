---
status: "completed"
started: "2026-05-29 01:27"
completed: "2026-05-29 01:30"
time_spent: "~3m"
---

# Task Record: T-clean-code Simplify and Clean Code

## Summary
Removed duplicate terminal status representations in pkg/task: inlined isTerminalStatus wrapper and replaced terminalStatuses map with types.IsTerminalStatus()

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/task/add.go
- forge-cli/pkg/task/statemachine.go

### Key Decisions
- Replaced local terminalStatuses map in add.go with types.IsTerminalStatus() to eliminate duplicate terminal-status definitions
- Inlined isTerminalStatus() wrapper in statemachine.go by calling types.IsTerminalStatus() directly

## Test Results
- **Tests Executed**: Yes
- **Passed**: 875
- **Failed**: 0
- **Coverage**: 0.0%

## Acceptance Criteria
- [x] No duplicate terminal-status representations remain in pkg/task
- [x] All existing tests pass after cleanup
- [x] golangci-lint reports 0 issues

## Notes
go build ./..., go vet ./..., and golangci-lint all clean. No tests needed creation since this was a pure simplification of existing code.
