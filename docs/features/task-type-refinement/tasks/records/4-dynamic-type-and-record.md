---
status: "completed"
started: "2026-05-16 21:05"
completed: "2026-05-16 21:18"
time_spent: "~13m"
---

# Task Record: 4 Dynamic fix task type by failure step and Type Reclassification in records

## Summary
Implement dynamic fix task type based on quality gate failure step (D4) and Type Reclassification block in records (D5). addFixTask() now deterministically sets type: compile/test -> fix, fmt/lint -> cleanup. RecordData gains optional TypeReclassification field rendered in fillRecordTemplate().

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/task/add.go
- forge-cli/pkg/task/types.go
- forge-cli/internal/cmd/quality_gate.go
- forge-cli/internal/cmd/submit.go
- forge-cli/internal/cmd/quality_gate_test.go
- forge-cli/internal/cmd/submit_test.go
- forge-cli/docs/WORKFLOW.md
- forge-cli/docs/WORKFLOW.zh.md

### Key Decisions
- fixTypeFromStep() uses a deterministic switch statement matching proposal D4 exactly: compile/unit-test/test-e2e -> TypeFix, fmt/lint -> TypeCleanup, default -> TypeFix
- TypeReclassification is a pointer field (nil = no reclassification) so existing records are unaffected
- TypeReclassification block renders between Summary and Changes sections in record template

## Test Results
- **Tests Executed**: Yes
- **Passed**: 20
- **Failed**: 0
- **Coverage**: 80.9%

## Acceptance Criteria
- [x] addFixTask() sets type to fix when failure is from compile or test step
- [x] addFixTask() sets type to cleanup when failure is from fmt or lint step
- [x] RecordData struct has optional TypeReclassification field (original type, actual type, reason)
- [x] fillRecordTemplate() renders Type Reclassification block only when TypeReclassification is non-nil
- [x] Type mapping follows proposal D4 table exactly

## Notes
Added Type field to AddTaskOpts to propagate task type through AddTask -> Task struct. Updated WORKFLOW.md and WORKFLOW.zh.md doc sync tests.
