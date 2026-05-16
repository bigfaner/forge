---
status: "completed"
started: "2026-05-16 21:30"
completed: "2026-05-16 21:32"
time_spent: "~2m"
---

# Task Record: 6 Update skill type assignment tables for new task types

## Summary
Updated type assignment tables in quick-tasks and breakdown-tasks SKILL.md files to replace deprecated 'implementation' type with the new taxonomy (feature, enhancement, cleanup, refactor). Added Intent Propagation sections to both skills describing how to read proposal 'intent' field as default task type with per-task override. Changed default type in both task.md templates from 'implementation' to 'feature'. Updated template selection logic to reflect 'feature' as the new default for code tasks.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/quick-tasks/SKILL.md
- plugins/forge/skills/breakdown-tasks/SKILL.md
- plugins/forge/skills/quick-tasks/templates/task.md
- plugins/forge/skills/breakdown-tasks/templates/task.md

### Key Decisions
- Used 'feature' as fallback default (matching migration logic in task 5) rather than introducing a new default
- Added deprecation blockquote for 'implementation' in both SKILL.md files with migration guidance, per Hard Rules
- Intent Propagation section is identical in structure across both skills, differing only in quick-tasks using 'documentation' vs breakdown-tasks using 'doc-generation' for the doc type

## Test Results
- **Tests Executed**: Yes
- **Passed**: 0
- **Failed**: 0
- **Coverage**: 0.0%

## Acceptance Criteria
- [x] quick-tasks/SKILL.md Type Assignment table lists: feature, enhancement, cleanup, refactor, documentation, gate
- [x] breakdown-tasks/SKILL.md Type Assignment table lists: feature, enhancement, cleanup, refactor, doc-generation, gate, test-pipeline
- [x] Type Assignment tables include clear 'When to assign' guidance for each type
- [x] Template selection logic updated: tasks with testable runtime behavior -> task.md, doc-only -> task-doc.md (unchanged), but default type in task.md frontmatter is 'feature'
- [x] Skills describe how to read proposal 'intent' field and use it as default type, with per-task override

## Notes
Documentation-only task. No runtime code changes, no test execution needed. Hard Rule satisfied: 'implementation' does not appear in the new type tables, only in a deprecation blockquote.
