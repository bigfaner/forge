---
status: "completed"
started: "2026-05-28 22:50"
completed: "2026-05-28 22:53"
time_spent: "~3m"
---

# Task Record: 1 Delete CLI behavior descriptions from pipeline skills/commands

## Summary
Deleted CLI behavior descriptions from 5 pipeline files: submit-task/SKILL.md (removed 'What forge task submit Does' section), breakdown-tasks/SKILL.md and quick-tasks/SKILL.md (removed 'Auto-generated tasks by forge task index' blocks), execute-task.md (removed field semantic explanations and subagent internal behavior), run-tasks.md (removed field semantic explanations, dedup explanations, and subagent internal behavior). All output contract field names and exit code contracts preserved.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/submit-task/SKILL.md
- plugins/forge/skills/breakdown-tasks/SKILL.md
- plugins/forge/skills/quick-tasks/SKILL.md
- plugins/forge/commands/execute-task.md
- plugins/forge/commands/run-tasks.md

### Key Decisions
无

## Document Metrics
38 lines deleted, 0 lines added; 5 files modified

## Referenced Documents
- docs/proposals/skill-instruction-audit/proposal.md

## Review Status
final

## Acceptance Criteria
- [x] submit-task/SKILL.md has no 'What .* Does' section; forge task submit command and exit code contract remain
- [x] breakdown-tasks/SKILL.md and quick-tasks/SKILL.md have no 'Auto-generated tasks by forge task index' block; Step 5 command remains
- [x] execute-task.md retains field name list but removes per-field semantic explanations and example values
- [x] run-tasks.md retains field name list; no 'Subagent calls forge prompt get-by-task-id internally'
- [x] No output contract field names lost (SURFACE_KEY, SURFACE_TYPE, TASK_ID, FILE, MAIN_SESSION remain)

## Notes
Applied three-category boundary rule from proposal: kept imperative instructions and output contracts, deleted behavioral explanations. Only modified the 5 files listed in Hard Rules.
