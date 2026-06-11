---
status: "completed"
started: "2026-05-19 01:10"
completed: "2026-05-19 01:12"
time_spent: "~2m"
---

# Task Record: 3 Strengthen docs-only classification in quick-tasks and breakdown-tasks

## Summary
Strengthened docs-only classification in quick-tasks and breakdown-tasks SKILL.md files. Added explicit 'classify by output artifact, not by intent' rule with cross-reference to type-assignment.md classification table in both skills' Type Assignment sections.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/quick-tasks/SKILL.md
- plugins/forge/skills/breakdown-tasks/SKILL.md

### Key Decisions
- Kept changes minimal — added the rule text after the existing type-assignment.md read instruction rather than restructuring the section, to avoid duplicating the classification table

## Test Results
- **Tests Executed**: Yes
- **Passed**: 0
- **Failed**: 0
- **Coverage**: 0.0%

## Acceptance Criteria
- [x] quick-tasks SKILL.md Type Assignment section explicitly states 'classify by output artifact, not intent' with cross-reference to type-assignment.md classification rule
- [x] breakdown-tasks SKILL.md Type Assignment section explicitly states the same rule
- [x] Both skills reference the new classification table from type-assignment.md

## Notes
Documentation-only task. Both files now have identical Type Assignment section text that makes the 'classify by output artifact' rule prominent alongside the existing read instruction for type-assignment.md.
