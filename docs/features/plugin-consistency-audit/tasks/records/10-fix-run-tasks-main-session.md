---
status: "completed"
started: "2026-05-30 06:00"
completed: "2026-05-30 06:01"
time_spent: "~1m"
---

# Task Record: 10 Fix: run-tasks MAIN_SESSION missing submit-task + git-commit

## Summary
Added submit-task and git-commit steps to run-tasks Step 1.5 MAIN_SESSION path, matching task-executor agent's Step 4-5 pattern with backward-compatible skip-if-already-handled clauses

## Changes

### Files Created
无

### Files Modified
- plugins/forge/commands/run-tasks.md

### Key Decisions
无

## Document Metrics
2 new steps (submit + commit), 1 mermaid diagram update

## Referenced Documents
- docs/features/plugin-consistency-audit/reports/05-commands-agent-hooks.md
- plugins/forge/agents/task-executor.md

## Review Status
final

## Acceptance Criteria
- [x] Step 1.5 includes submit-task directive: Skill(skill="forge:submit-task")
- [x] Step 1.5 includes conditional git-commit directive
- [x] Backward compatible: skip if task already handled submit/commit

## Notes
Steps 5-6 mirror task-executor agent Step 4-5 pattern. Mermaid diagram updated with '1.5.5 Submit + Commit' node to reflect the new flow.
