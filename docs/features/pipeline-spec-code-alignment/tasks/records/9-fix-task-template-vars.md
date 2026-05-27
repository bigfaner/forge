---
status: "completed"
started: "2026-05-27 01:19"
completed: "2026-05-27 01:22"
time_spent: "~3m"
---

# Task Record: 9 Add fix-task template variables to all creation points

## Summary
Added --var SOURCE_FILES/TEST_SCRIPT/TEST_RESULTS to all fix-task creation points (10 locations across 4 files), added IT impact assessment guidance to quick-tasks and breakdown-tasks SKILL.md

## Changes

### Files Created
无

### Files Modified
- plugins/forge/commands/run-tasks.md
- plugins/forge/agents/task-executor.md
- plugins/forge/commands/execute-task.md
- plugins/forge/skills/submit-task/SKILL.md
- plugins/forge/skills/quick-tasks/SKILL.md
- plugins/forge/skills/breakdown-tasks/SKILL.md

### Key Decisions
无

## Document Metrics
6 files modified, 10 fix-task creation points updated with --var, 2 new IT impact assessment sections added

## Referenced Documents
- docs/proposals/pipeline-spec-code-alignment/proposal.md

## Review Status
completed

## Acceptance Criteria
- [x] Every fix-task creation point in run-tasks.md includes --var SOURCE_FILES/TEST_SCRIPT/TEST_RESULTS
- [x] task-executor.md fix-task creation includes all three --var parameters
- [x] execute-task.md fix-task creation includes all three --var + --description
- [x] submit-task/SKILL.md recovery includes all three --var
- [x] quick-tasks and breakdown-tasks SKILL.md have breaking task IT impact assessment guidance
- [x] Fix-task grouping guidance specifies by test suite (directory), not problem type

## Notes
Found 10 fix-task creation points total (run-tasks: 5, task-executor: 1, execute-task: 4, submit-task: 1). execute-task.md main session fails entry also gained missing --description parameter.
