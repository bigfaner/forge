---
status: "completed"
started: "2026-05-24 21:36"
completed: "2026-05-24 21:50"
time_spent: "~14m"
---

# Task Record: 3 Update quality gate steps and failure handling

## Summary
Updated quality gate steps and failure handling: migrated runE2ERegression to runTestRegression (e2e-setup->test-setup, e2e-test->test), updated RunGate to error on missing required recipes with init-justfile hint (no fallback), updated cobra command descriptions and warning messages

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/just/just.go
- forge-cli/pkg/just/just_test.go
- forge-cli/internal/cmd/quality_gate.go
- forge-cli/internal/cmd/quality_gate_test.go
- forge-cli/internal/cmd/integration_test.go

### Key Decisions
- RunGate now calls onFail and returns false when a required recipe is missing, instead of printing a WARNING and skipping
- Added hasRecipeWithArg helper to handle scoped recipes that require arguments
- Renamed runE2ERegression to runTestRegression with all e2e recipe references updated to test/test-setup
- Updated cobra Long description and warning messages to remove e2e references

## Test Results
- **Tests Executed**: Yes
- **Passed**: 58
- **Failed**: 0
- **Coverage**: 60.1%

## Acceptance Criteria
- [x] Step 2 (unit test step) runs just unit-test
- [x] Step 3 (advanced test step) runs just test
- [x] addFixTask uses generic rule step -> 'just ' + step, no hardcoded recipe name mapping
- [x] handleGateFailure guide/label map has 'test' key instead of 'e2e-test'
- [x] runE2ERegression migrated: e2e-setup -> test-setup, e2e-test -> test, function renamed to runTestRegression
- [x] When unit-test recipe missing, quality gate reports error suggesting run init-justfile (no fallback)

## Notes
Integration tests in integration_test.go updated to use unit-test recipe instead of test for breaking gate scenarios
