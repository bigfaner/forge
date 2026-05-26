---
status: "completed"
started: "2026-05-26 22:59"
completed: "2026-05-26 23:01"
time_spent: "~2m"
---

# Task Record: 2 Update plugin components to use --type coding.fix

## Summary
Replaced all 10 occurrences of --template fix-task with --type coding.fix across 6 plugin files (1 agent, 2 commands, 3 skills). No --template cleanup-task occurrences found.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/agents/task-executor.md
- plugins/forge/commands/execute-task.md
- plugins/forge/commands/run-tasks.md
- plugins/forge/skills/breakdown-tasks/SKILL.md
- plugins/forge/skills/quick-tasks/SKILL.md
- plugins/forge/skills/submit-task/SKILL.md

### Key Decisions
无

## Document Metrics
6 files modified, 10 occurrences replaced, 0 remaining old syntax

## Referenced Documents
- docs/proposals/unify-task-add-type-param/proposal.md

## Review Status
final

## Acceptance Criteria
- [x] No plugin markdown file contains --template fix-task or --template cleanup-task
- [x] All 6 files use --type coding.fix
- [x] No changes to template variable flags (--var, SOURCE_FILES, TEST_SCRIPT, TEST_RESULTS)

## Notes
Mechanical search-and-replace only. No cleanup-task occurrences existed in plugin files.
