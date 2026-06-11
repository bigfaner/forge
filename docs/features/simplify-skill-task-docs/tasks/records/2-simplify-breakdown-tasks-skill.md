---
status: "completed"
started: "2026-05-15 20:10"
completed: "2026-05-15 20:17"
time_spent: "~7m"
---

# Task Record: 2 Simplify breakdown-tasks SKILL.md

## Summary
Simplified breakdown-tasks SKILL.md by removing CLI auto-generated task descriptions: merged Steps 4b/4c into HARD-RULE block, reduced Step 4d (now 4b) to brief statement + fix-task command, removed 13 task types inference table, condensed Scope Assignment examples, Hard Rules section, Existing-Code Modification Split, and Template Selection examples. Line count reduced from 478 to 355 (25.7% reduction, target >=25%).

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/breakdown-tasks/SKILL.md

### Key Decisions
- Merged stage-gate auto-generation info into HARD-RULE block to avoid duplication between 4b and the rule block
- Removed separate 4b section entirely, renumbered 4d to 4b
- Condensed Scope Assignment by removing examples table (9 lines) and non-mixed projects note
- Condensed Hard Rules from 5 bullets to single sentence
- Condensed Existing-Code Modification Split by removing example block (8 lines) and merging sub-items into single lines

## Test Results
- **Tests Executed**: Yes
- **Passed**: 0
- **Failed**: 0
- **Coverage**: 0.0%

## Acceptance Criteria
- [x] Steps 4b/4c merged into brief statement about auto-generation
- [x] Step 4d reduced to brief statement + fix-task command reference
- [x] 13 task types inference table removed from Template Selection
- [x] Fix-task forge task add one-liner kept
- [x] Step 0 Profile Resolution unchanged
- [x] Steps 1-3, 4a, 5-7 unchanged (core skill logic preserved)
- [x] Output Checklist items about test tasks and gate files simplified
- [x] Total line count reduced by >=25%

## Notes
Documentation-only task. No tests to run. Coverage set to -1.0.
