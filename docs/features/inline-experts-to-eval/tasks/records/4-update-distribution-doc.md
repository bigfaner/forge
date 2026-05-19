---
status: "completed"
started: "2026-05-19 11:23"
completed: "2026-05-19 11:25"
time_spent: "~2m"
---

# Task Record: 4 Update forge-distribution.md for new expert location

## Summary
Updated forge-distribution.md to reflect expert files moved from agents/experts/ to skills/eval/experts/. Updated directory tree, component table, and section 3 path references.

## Changes

### Files Created
无

### Files Modified
- docs/conventions/forge-distribution.md

### Key Decisions
- Kept experts/ as a generic subtree under <skill-name>/ in the directory tree (not eval-specific), since the convention applies to any skill that may include expert subdirectories in the future

## Test Results
- **Tests Executed**: No
- **Passed**: 0
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] Directory tree: agents/experts/ subtree moved under skills/eval/, removed from agents/
- [x] Component table: remove agents/experts/ from agents row description; add experts/ mention under skills row
- [x] Section 3: update title and all paths to reflect new location under skills/eval/experts/
- [x] No remaining references to agents/experts/ in the document

## Notes
Documentation-only task. grep confirms zero remaining references to agents/experts/.
