---
status: "completed"
started: "2026-05-24 19:16"
completed: "2026-05-24 19:18"
time_spent: "~2m"
---

# Task Record: 5 更新活跃文档：移除 gen-and-run 引用

## Summary
Removed gen-and-run references from OVERVIEW.md and task-lifecycle.md; updated OVERVIEW.md test pipeline types to reflect current staged architecture (gen-journeys -> gen-contracts -> gen-scripts -> run -> verify-regression); updated task-lifecycle.md system types list to remove gen-and-run and replace deprecated types with current ones

## Changes

### Files Created
无

### Files Modified
- forge-cli/docs/OVERVIEW.md
- docs/business-rules/task-lifecycle.md

### Key Decisions
无

## Document Metrics
2 files updated, 2 gen-and-run references removed, 1 type table rewritten with 5 staged pipeline types

## Referenced Documents
- docs/proposals/pipeline-integration-stitch/proposal.md

## Review Status
completed

## Acceptance Criteria
- [x] grep -r 'gen-and-run' forge-cli/docs/ docs/business-rules/ returns zero results
- [x] OVERVIEW.md test pipeline description reflects staged architecture (gen-journeys -> gen-contracts -> gen-scripts -> run -> verify-regression)

## Notes
Also corrected outdated test-pipeline.* prefix to test.* in OVERVIEW.md type table, and updated task-lifecycle.md system types list to reflect current types (replaced test.gen-cases/test.eval-cases/test.gen-and-run with test.gen-journeys/test.gen-contracts). Count updated from 12 to 11 system types.
