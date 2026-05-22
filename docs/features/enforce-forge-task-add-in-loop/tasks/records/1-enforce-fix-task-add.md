---
status: "completed"
started: "2026-05-22 10:49"
completed: "2026-05-22 10:53"
time_spent: "~4m"
---

# Task Record: 1 Enforce explicit forge task add in run-tasks dispatcher

## Summary
Replace all ambiguous spawn fix task instructions in run-tasks.md with explicit forge task add commands using proper flags (--template fix-task, --source-task-id, --block-source, --description)

## Changes

### Files Created
无

### Files Modified
- plugins/forge/commands/run-tasks.md

### Key Decisions
- Each fix task instance uses explicit forge task add with --template fix-task, --source-task-id, --block-source, and --description flags
- Built-in dedup via HasActiveFixTasks() handles duplicate fix task prevention

## Test Results
- **Tests Executed**: Yes
- **Passed**: 23
- **Failed**: 0
- **Coverage**: 83.4%

## Acceptance Criteria
- [x] run-tasks.md no longer contains the phrase spawn fix task
- [x] Step 1.5 (main session failure) uses explicit forge task add with proper flags
- [x] Step 2b (dispatched task blocked) uses explicit forge task add with proper flags
- [x] Step 3 (main session error table) uses explicit forge task add with proper flags
- [x] Each instance includes correct --source-task-id and --block-source flags
- [x] Each instance includes informative --description explaining failure reason
- [x] Failure tracking comment (consecutive_failures) remains accurate

## Notes
Three locations updated: Step 1.5 main session failure, Step 2b dispatched task blocked, Step 3 error table. All use explicit forge task add with --template fix-task, --source-task-id, --block-source, and --description.
