---
status: "completed"
started: "2026-05-24 03:30"
completed: "2026-05-24 03:38"
time_spent: "~8m"
---

# Task Record: 11 Fix validateRecordData os.Exit to return error

## Summary
Refactored validateRecordData from void+os.Exit to return error pattern. Replaced all base.Exit() calls with return error. Updated doSubmit caller to handle returned error. Converted 8 subprocess-pattern tests to direct error assertions. Updated testbridge.go export and integration_test.go caller.

## Changes

### Files Created
无

### Files Modified
- forge-cli/internal/cmd/task/submit.go
- forge-cli/internal/cmd/task/submit_test.go
- forge-cli/internal/cmd/task/testbridge.go
- forge-cli/internal/cmd/integration_test.go

### Key Decisions
- Changed validateRecordData signature from func(...) to func(...) error, replacing base.Exit() with return err
- Converted 8 subprocess-pattern tests (exec.Command+env var) to direct error assertions, eliminating os.Exit dependency
- Simplified PRESERVE tests: replaced stderr-capture-and-check-for-ERROR with direct err != nil checks where no warning output needed

## Test Results
- **Tests Executed**: Yes
- **Passed**: 42
- **Failed**: 0
- **Coverage**: 70.1%

## Acceptance Criteria
- [x] validateRecordData returns error instead of calling os.Exit
- [x] All callers updated to handle returned error
- [x] 30+ test cases adapted to new signature
- [x] 0 non-top-level os.Exit calls in validateRecordData
- [x] go build ./... passes
- [x] go test ./internal/cmd/task/ passes

## Notes
Hard Rules preserved: quality_gate.go os.Exit(0) calls in RunE handlers unchanged. validateQualityGate uses panic (not os.Exit) — out of scope for this task.
