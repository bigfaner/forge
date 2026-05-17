---
status: "completed"
started: "2026-05-17 13:31"
completed: "2026-05-17 13:34"
time_spent: "~3m"
---

# Task Record: 2 Update prompt templates with discovery instruction

## Summary
Replaced vague 'Read relevant project knowledge files' instruction in all 5 prompt templates (fix.md, cleanup.md, enhancement.md, feature.md, refactor.md) with the discovery instruction that uses YAML frontmatter domains field for relevance matching.

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/prompt/data/fix.md
- forge-cli/pkg/prompt/data/cleanup.md
- forge-cli/pkg/prompt/data/enhancement.md
- forge-cli/pkg/prompt/data/feature.md
- forge-cli/pkg/prompt/data/refactor.md

### Key Decisions
- fix.md had slightly different wording ('based on the affected files and error context') compared to the other 4 templates ('based on the task domain'), but both were replaced with the same discovery instruction per proposal spec
- For the 4 templates that combined knowledge reading with task file reading in one sentence, the discovery instruction was inserted as a separate block before 'Then read the task file' to maintain the logical flow

## Test Results
- **Tests Executed**: Yes
- **Passed**: 20
- **Failed**: 0
- **Coverage**: 90.6%

## Acceptance Criteria
- [x] All 5 files contain the discovery instruction
- [x] No file contains the text 'Read relevant project knowledge files'
- [x] The discovery instruction matches the proposal's specification

## Notes
All 5 prompt templates now have identical discovery instruction text matching the proposal spec exactly. grep confirmed zero occurrences of the old vague instruction remain.
