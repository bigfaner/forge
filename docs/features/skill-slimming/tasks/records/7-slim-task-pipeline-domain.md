---
status: "completed"
started: "2026-05-20 14:40"
completed: "2026-05-20 14:48"
time_spent: "~8m"
---

# Task Record: 7 Slim task pipeline domain (breakdown-tasks + quick-tasks + submit-task)

## Summary
Slim task pipeline domain: breakdown-tasks (144→144 lines, extracted scope-assignment rules, fixed doc type inconsistency), quick-tasks (208→173 lines, streamlined Type Assignment/Intent Propagation/Commit sections, fixed type terminology), submit-task (156→156 lines, replaced doc* ambiguous wildcard with precise doc type). All 3 files well under 350-line limit.

## Changes

### Files Created
- plugins/forge/skills/breakdown-tasks/rules/scope-assignment.md

### Files Modified
- plugins/forge/skills/breakdown-tasks/SKILL.md
- plugins/forge/skills/quick-tasks/SKILL.md
- plugins/forge/skills/submit-task/SKILL.md

### Key Decisions
- Unified type: 'documentation' to 'doc' across breakdown-tasks and quick-tasks for consistency
- Extracted Scope Assignment path classification rules from breakdown-tasks into rules/scope-assignment.md (Splitting Heuristic: >5 lines of rule definitions)
- Replaced doc* ambiguous wildcard in submit-task with precise 'doc' type reference
- Streamlined quick-tasks Type Assignment from ~40 lines to ~15 lines by removing redundant explanation
- Condensed quick-tasks Step 8 Commit section from ~24 lines to ~10 lines while preserving all constraints

## Test Results
- **Tests Executed**: No
- **Passed**: 0
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] Each SKILL.md line count <= 350
- [x] No ambiguous terms remain (noTest, doc* replaced with precise references)
- [x] All referenced auxiliary file paths exist and are readable

## Notes
Doc-type task. All 3 files were already small (144-208 lines), so focus was on disambiguation and redundancy trimming rather than splitting.
