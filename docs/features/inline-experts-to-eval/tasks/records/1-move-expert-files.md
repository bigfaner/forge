---
status: "completed"
started: "2026-05-19 11:07"
completed: "2026-05-19 11:11"
time_spent: "~4m"
---

# Task Record: 1 Move expert files from agents/ to skills/eval/experts/

## Summary
Moved 11 expert/protocol files from agents/experts/ to skills/eval/experts/ using git mv. All files relocated with identical content (0 insertions/0 deletions). Removed empty agents/experts/ directory. task-executor.md remains in agents/.

## Changes

### Files Created
- plugins/forge/skills/eval/experts/protocol/scorer-protocol.md
- plugins/forge/skills/eval/experts/protocol/reviser-protocol.md
- plugins/forge/skills/eval/experts/scorer/architect.md
- plugins/forge/skills/eval/experts/scorer/code-reviewer.md
- plugins/forge/skills/eval/experts/scorer/cto.md
- plugins/forge/skills/eval/experts/scorer/editor.md
- plugins/forge/skills/eval/experts/scorer/harness-engineer.md
- plugins/forge/skills/eval/experts/scorer/pm.md
- plugins/forge/skills/eval/experts/scorer/qa.md
- plugins/forge/skills/eval/experts/scorer/ux-auditor.md
- plugins/forge/skills/eval/experts/scorer/ux-engineer.md

### Files Modified
无

### Key Decisions
- Used git mv to preserve clean history tracking
- Removed only agents/experts/ subdirectories, preserving agents/task-executor.md

## Test Results
- **Tests Executed**: No
- **Passed**: 0
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] All 11 files exist under skills/eval/experts/ with identical content
- [x] agents/experts/ directory no longer exists
- [x] No forge:experts:* agents appear in /agents panel after plugin cache refresh

## Notes
This is Part A only (file move) from the proposal. Part B (SKILL.md path updates and mechanism fixes) is a separate task.
