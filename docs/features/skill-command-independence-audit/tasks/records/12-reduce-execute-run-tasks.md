---
status: "completed"
started: "2026-06-04 01:13"
completed: "2026-06-04 01:16"
time_spent: "~3m"
---

# Task Record: 12 Reduce overlapping logic in execute-task and run-tasks commands

## Summary
Reduced overlapping logic in execute-task and run-tasks commands: extracted fix task template into Error Handling section, simplified claim output format, removed inline duplicated forge task add commands. execute-task: 150->128 lines (-14.7%), run-tasks: 165->148 lines (-10.3%).

## Changes

### Files Created
无

### Files Modified
- plugins/forge/commands/execute-task.md
- plugins/forge/commands/run-tasks.md

### Key Decisions
无

## Document Metrics
execute-task: 150->128 lines (-14.7%), run-tasks: 165->148 lines (-10.3%), total reduction: 39 lines (-12.4%)

## Referenced Documents
- docs/proposals/skill-command-independence-audit/proposal.md

## Review Status
final

## Acceptance Criteria
- [x] claim format description in both commands is independent and concise
- [x] fix-type table in both commands is independent and concise
- [x] both commands are fully executable with no dangling references

## Notes
Extracted fix task template into Error Handling section of each command. All inline forge task add commands now reference the template instead of duplicating it. Hard rules, decision tables, and step sequences preserved.
