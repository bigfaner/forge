---
status: "completed"
started: "2026-05-23 13:05"
completed: "2026-05-23 13:16"
time_spent: "~11m"
---

# Task Record: 6 Type-aware validation in validateRecordData()

## Summary
Make validateRecordData() type-aware using CategoryForType(). Added taskType parameter to validateRecordData, skip test evidence checks and testsFailed auto-downgrade for non-coding categories (doc, test, validation, gate). Updated doSubmit() to pass t.Type. Summary, AC, and recommended field warnings remain for all categories.

## Changes

### Files Created
无

### Files Modified
- forge-cli/internal/cmd/task/submit.go
- forge-cli/internal/cmd/task/submit_test.go
- forge-cli/internal/cmd/integration_test.go

### Key Decisions
- Used CategoryForType(taskType) == CategoryCoding as the single gate for test-related validation, consistent with task 1's category system
- Kept IsTestableType check in doSubmit (line 132) as redundant safety net for auto-setting coverage=-1.0
- Unmet AC validation applies to ALL categories (not just coding) as per Hard Rules
- keyDecisions and acceptanceCriteria warnings remain for all completed tasks regardless of category

## Test Results
- **Tests Executed**: Yes
- **Passed**: 24
- **Failed**: 0
- **Coverage**: 4.7%

## Acceptance Criteria
- [x] validateRecordData accepts task type parameter
- [x] Doc-type tasks with testsPassed=0, testsFailed=0, coverage=0 pass validation (no ErrNoTestEvidence)
- [x] Coding-type tasks with same values still fail with ErrNoTestEvidence
- [x] Doc-type tasks with testsFailed > 0 are NOT auto-downgraded to blocked
- [x] Coding-type tasks with testsFailed > 0 still auto-downgrade
- [x] summary remains hard-required for ALL task types
- [x] doSubmit() updated to pass t.Type to validateRecordData()
- [x] Unit tests for validation behavior per category

## Notes
12 new type-aware test cases added in TestValidateRecordData_TypeAware. All existing tests updated to pass taskType parameter. Integration test call site updated.
