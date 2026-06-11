---
status: "completed"
started: "2026-05-20 00:53"
completed: "2026-05-20 00:54"
time_spent: "~1m"
---

# Task Record: 5 Inject Karpathy principles into fix-bug command

## Summary
Added CODING_PRINCIPLES block with all 4 Karpathy principles (Think Before Coding, Simplicity First, Surgical Changes, Goal-Driven Execution) to fix-bug command. Positioned after Core principle line before Parameters. Merged overlapping rules: removed redundant 'Fix only what failing tests require / no scope creep / no refactoring / no improvements' from EXTREMELY-IMPORTANT and Step 4 Fix section, as these are now covered by Simplicity First + Surgical Changes. Preserved atomic commit requirement, HARD-GATE at Step 2, Knowledge Review, workflow step numbering (1-6), and Common Pitfalls table.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/commands/fix-bug.md

### Key Decisions
- Tailored principle text to bug-fix context: Think Before Coding references 'restating bug and suspected root cause', Goal-Driven Execution references 'tests that fail before fix and pass after'
- Replaced Step 4 Fix Principles list with a single guideline referencing CODING_PRINCIPLES, keeping the unique dependency root-cause guidance
- Removed 'Fix only what failing tests require — no scope creep, no refactoring, no improvements' from EXTREMELY-IMPORTANT since Simplicity First + Surgical Changes cover this semantically

## Test Results
- **Tests Executed**: No
- **Passed**: 0
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] Contains CODING_PRINCIPLES block with Think Before Coding + Simplicity First + Surgical Changes + Goal-Driven Execution
- [x] CODING_PRINCIPLES positioned after Core principle line, before Parameters
- [x] EXTREMELY-IMPORTANT overlapping rules merged into principles
- [x] No scope creep / no refactoring / no improvements rule absorbed by Simplicity First + Surgical Changes
- [x] Knowledge Review section untouched
- [x] Workflow step numbering (Steps 1-6) unchanged
- [x] Common Pitfalls table preserved

## Notes
Atomic commit requirement and HARD-GATE at Step 2 both preserved per Hard Rules.
