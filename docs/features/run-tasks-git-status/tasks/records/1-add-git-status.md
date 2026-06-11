---
status: "completed"
started: "2026-05-23 10:54"
completed: "2026-05-23 10:54"
time_spent: ""
---

# Task Record: 1 Add git status to run-tasks Post-Completion

## Summary
Added Git Status Summary subsection to run-tasks Post-Completion section. Instructs dispatcher to display branch name + ahead/behind relative to main and working tree changes (git status --short) after loop ends. All git commands wrapped in error handling with silent skip on failure. Existing Post-Completion content preserved unchanged.

## Changes

### Files Created
None

### Files Modified
- plugins/forge/commands/run-tasks.md

### Key Decisions
None

## Test Results
- **Tests Executed**: No
- **Passed**: 0
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] Post-Completion section instructs dispatcher to run git commands showing: current branch name, ahead/behind relative to main, and changed/untracked files (git status --short)
- [x] Git commands wrapped in error handling — skip silently on failure
- [x] Existing Post-Completion content preserved (test task message, artifact commit prohibition)

## Notes
Added as ### Git Status Summary subsection under ## Post-Completion. Branch info via git branch --show-current + git rev-list --left-right --count main...HEAD. Working tree via git status --short. No diff statistics per Hard Rules.
