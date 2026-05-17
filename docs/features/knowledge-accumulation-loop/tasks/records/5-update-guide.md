---
status: "completed"
started: "2026-05-17 21:46"
completed: "2026-05-17 21:49"
time_spent: "~3m"
---

# Task Record: 5 Update guide.md — replace old skills with /learn

## Summary
Updated plugins/forge/hooks/guide.md to replace /record-decision and /learn-lesson references with unified /learn skill. Added Knowledge Accumulation section documenting /learn as the primary manual entry point, auto-extract triggers at 4 pipeline completion points (run-tasks, fix-bug, write-prd, tech-design), and /consolidate-specs for bulk extraction + drift detection.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/hooks/guide.md

### Key Decisions
- Placed Knowledge Accumulation section between Evaluation Parameter Exceptions and Other Auxiliary Skills to maintain guide flow
- Split auxiliary skills into knowledge-related (new dedicated section) and other auxiliary skills (existing table)

## Test Results
- **Tests Executed**: Yes
- **Passed**: 0
- **Failed**: 0
- **Coverage**: 0.0%

## Acceptance Criteria
- [x] decisions/ directory description updated: references /learn instead of /record-decision
- [x] lessons/ directory description updated: references /learn instead of /learn-lesson
- [x] No remaining references to /record-decision or /learn-lesson as active skills
- [x] /learn documented as the primary knowledge accumulation entry point
- [x] Auto-extract flow documented: triggers at run-tasks, fix-bug, write-prd, tech-design completion
- [x] /consolidate-specs still documented for bulk extraction + drift detection (unchanged role)

## Notes
Minimal changes per hard rules — only updated knowledge-related references. All other guide sections preserved unchanged.
