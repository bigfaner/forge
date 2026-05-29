---
status: "completed"
started: "2026-05-29 11:14"
completed: "2026-05-29 11:20"
time_spent: "~6m"
---

# Task Record: 3 Update skill files with fix-type derivation rule

## Summary
Replaced all 13 hardcoded --type coding.fix in error-handling instructions across 7 plugin files with a category-based derivation rule. Added TASK_CATEGORY as extractable field in execute-task.md and run-tasks.md. Added doc.fix to valid types in record-format-coding.md. Documented derivation table in run-tasks.md and execute-task.md as canonical locations.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/agents/task-executor.md
- plugins/forge/commands/execute-task.md
- plugins/forge/commands/run-tasks.md
- plugins/forge/skills/submit-task/SKILL.md
- plugins/forge/skills/breakdown-tasks/SKILL.md
- plugins/forge/skills/quick-tasks/SKILL.md
- plugins/forge/skills/submit-task/data/record-format-coding.md

### Key Decisions
无

## Document Metrics
13 occurrences replaced, 7 files modified, 2 canonical derivation tables added

## Referenced Documents
- docs/proposals/derive-fix-task-type/proposal.md

## Review Status
final

## Acceptance Criteria
- [x] Error-handling instructions in task-executor.md, execute-task.md, run-tasks.md, submit-task/SKILL.md use derivation rule
- [x] Derivation rule table documented in at least one canonical location
- [x] TASK_CATEGORY documented as extractable field from claim output
- [x] grep for --type coding.fix returns zero matches in error-handling contexts

## Notes
Informational mentions of coding.fix in type reclassification tables and Type Assignment tables were kept as-is per task Implementation Notes distinction between error-handling and informational.
