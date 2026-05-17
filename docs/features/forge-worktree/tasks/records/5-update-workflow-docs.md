---
status: "completed"
started: "2026-05-17 18:13"
completed: "2026-05-17 18:15"
time_spent: "~2m"
---

# Task Record: 5 Update WORKFLOW.md with forge worktree commands

## Summary
Updated Section 9 Option 2 in both WORKFLOW.md and WORKFLOW.zh.md to document forge worktree start/list/resume/remove commands as the primary workflow, with manual git worktree steps retained as fallback reference.

## Changes

### Files Created
无

### Files Modified
- forge-cli/docs/WORKFLOW.md
- forge-cli/docs/WORKFLOW.zh.md

### Key Decisions
- Structured Option 2 with sub-sections: Commands table, Typical Workflow, Multiple Features in Parallel, Manual Git Worktree (Fallback)
- Kept existing document structure and formatting style (headings, code blocks, tables)
- Both EN and ZH versions updated with equivalent content

## Test Results
- **Tests Executed**: No
- **Passed**: 0
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] Section 9 Option 2 updated to show forge worktree start <slug> as primary workflow
- [x] All four commands documented with examples (start, list, resume, remove)
- [x] Manual git worktree steps kept as fallback reference
- [x] Both WORKFLOW.md and WORKFLOW.zh.md updated

## Notes
Documentation-only task. noTest applies.
