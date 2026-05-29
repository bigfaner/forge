---
status: "completed"
started: "2026-05-29 11:31"
completed: "2026-05-29 11:32"
time_spent: "~1m"
---

# Task Record: doc-fix-1 Fix: TYPE not listed as extractable field in claim output docs

## Summary
Added missing TYPE field to extractable fields list in execute-task.md and run-tasks.md claim output documentation

## Changes

### Files Created
无

### Files Modified
- plugins/forge/commands/execute-task.md
- plugins/forge/commands/run-tasks.md

### Key Decisions
无

## Document Metrics
N/A

## Referenced Documents
- forge-cli/internal/cmd/task/claim.go

## Review Status
final

## Acceptance Criteria
- [x] TYPE field listed in execute-task.md Extract from claim output section
- [x] TYPE field listed in run-tasks.md Extract section
- [x] No code files modified (doc-only fix)

## Notes
TYPE is output by forge task claim (claim.go:288) but was missing from both command docs. Added with description referencing the full type string format.
