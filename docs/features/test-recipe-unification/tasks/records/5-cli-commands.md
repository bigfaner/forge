---
status: "completed"
started: "2026-05-24 21:56"
completed: "2026-05-24 22:00"
time_spent: "~4m"
---

# Task Record: 5 Update CLI command text and help messages

## Summary
Updated CLI command text and help messages: init.go wizard prompts changed from 'e2e tests' to 'advanced tests', test.go run-journey description updated from 'just e2e-test' to 'just test'

## Changes

### Files Created
- forge-cli/internal/cmd/test/test_cli_test.go

### Files Modified
- forge-cli/internal/cmd/init.go
- forge-cli/internal/cmd/init_test.go
- forge-cli/internal/cmd/test/test.go

### Key Decisions
- Used 'advanced tests' terminology in wizard prompts to reflect surface-level test model
- Used 'surface-level advanced tests' in descriptions to explain the two-layer model distinction

## Test Results
- **Tests Executed**: Yes
- **Passed**: 5
- **Failed**: 0
- **Coverage**: 4.5%

## Acceptance Criteria
- [x] init.go wizard prompts reference 'test' not 'e2e-test'
- [x] test/test.go all e2e-test references updated to test
- [x] Help text accurately describes the two-layer test model

## Notes
Text-only changes. Coverage percentage reflects targeted tests on CLI help text only; package-level coverage is low because this task only touched string literals in command definitions.
