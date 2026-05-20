---
status: "completed"
started: "2026-05-20 22:21"
completed: "2026-05-20 22:23"
time_spent: "~2m"
---

# Task Record: 3 Update knowledge extraction in fix-bug, write-prd, tech-design

## Summary
Added auto-save configuration check (Step 4.5) to 3 knowledge extraction files: fix-bug.md, write-prd/rules/knowledge-extraction.md, tech-design/rules/knowledge-extraction.md. When auto.knowledgeSave is true for the current mode, Step 5 (AskUserQuestion) is skipped and all candidates auto-save to Step 6.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/commands/fix-bug.md
- plugins/forge/skills/write-prd/rules/knowledge-extraction.md
- plugins/forge/skills/tech-design/rules/knowledge-extraction.md

### Key Decisions
- Inserted config check as Step 4.5 (between Step 4 silent exit and Step 5 confirmation) to minimize disruption to existing flow numbering
- Used identical branching pattern across all 3 files for consistency, matching the auto.runTasks pattern from quick.md

## Test Results
- **Tests Executed**: No
- **Passed**: 0
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] All 3 files read auto.knowledgeSave via forge config get auto.knowledgeSave before Step 5
- [x] When mode is true, Step 5 (AskUserQuestion) is skipped and all candidates auto-saved to Step 6
- [x] When mode is false, existing AskUserQuestion flow is preserved unchanged
- [x] The config read and branching logic is consistent across all 3 files

## Notes
无
