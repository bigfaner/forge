---
status: "completed"
started: "2026-05-17 21:56"
completed: "2026-05-17 22:03"
time_spent: "~7m"
---

# Task Record: 7 Add auto-extract trigger to fix-bug

## Summary
Added Knowledge Review section to fix-bug command after Step 6 (Commit) and before Output Summary. The section triggers knowledge auto-extraction by referencing the shared extraction routine, scanning root cause analysis and fix approach for notable knowledge (non-obvious root causes, debugging patterns, gotchas). Silent for trivial fixes.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/commands/fix-bug.md

### Key Decisions
- Followed the same pattern as Task 6 (run-tasks trigger) for consistency across trigger points
- Included extraction routine by reference (read knowledge-extraction.md) per Hard Rules, not by copying content

## Test Results
- **Tests Executed**: Yes
- **Passed**: 20
- **Failed**: 0
- **Coverage**: 83.4%

## Acceptance Criteria
- [x] After Step 6 (Commit), a new knowledge review step runs before Output Summary
- [x] Step reads plugins/forge/references/shared/knowledge-extraction.md for extraction logic
- [x] Scans root cause analysis and fix approach for notable knowledge
- [x] Looks for: non-obvious root causes, debugging patterns, gotchas
- [x] Silent when the fix was trivial (typo, simple config change)
- [x] Presents extracted knowledge via AskUserQuestion for user confirmation
- [x] Writes confirmed knowledge to appropriate directories using shared formats
- [x] Does not interfere with existing fix-bug workflow (Steps 1-6)

## Notes
Enhancement task modifying markdown command file only, no code tests applicable. Go test suite runs as quality gate.
