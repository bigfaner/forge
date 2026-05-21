---
status: "completed"
started: "2026-05-20 22:23"
completed: "2026-05-20 22:32"
time_spent: "~9m"
---

# Task Record: 4 Add knowledgeSave and runTasks to forge init guided config

## Summary
Added runTasks and knowledgeSave guided configuration to both forge init (TUI via huh) and forge config init (CLI via stdin). Questions follow existing patterns with correct defaults from AutoConfigDefaults().

## Changes

### Files Created
无

### Files Modified
- forge-cli/internal/cmd/config.go
- forge-cli/internal/cmd/config_test.go
- forge-cli/internal/cmd/init.go

### Key Decisions
- Placed runTasks and knowledgeSave questions after validation and before gitPush, following the suggested order from the task spec
- runTasks defaults: Quick=true, Full=false; knowledgeSave defaults: Quick=true, Full=false (matching AutoConfigDefaults)
- TUI uses askConfirm helper (huh.NewConfirm), CLI uses readBool with matching defaults

## Test Results
- **Tests Executed**: Yes
- **Passed**: 12
- **Failed**: 0
- **Coverage**: 80.3%

## Acceptance Criteria
- [x] askAutoBehavior() in init.go includes runTasks questions (quick/full) with correct defaults
- [x] askAutoBehavior() in init.go includes knowledgeSave questions (quick/full) with correct defaults
- [x] runConfigInit() in config.go includes runTasks prompts (quick/full) with correct defaults
- [x] runConfigInit() in config.go includes knowledgeSave prompts (quick/full) with correct defaults
- [x] Config struct construction in both files includes the new fields
- [x] Question ordering follows the logical grouping of auto config fields
- [x] Existing init tests pass; new test cases cover the added questions

## Notes
无
