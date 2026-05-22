---
status: "completed"
started: "2026-05-22 23:52"
completed: "2026-05-23 00:16"
time_spent: "~24m"
---

# Task Record: 2 Migrate test_results.go and journey_isolation.go to pkg/testrunner

## Summary
Migrated test_results.go and journey_isolation.go (plus test) from internal/cmd/ to pkg/testrunner/. Exported previously-unexported symbols to enable cross-package access. Updated all callers in quality_gate.go, test.go, test_promote.go, and integration_test.go to use testrunner.XXX. Kept CLI command registration tests in internal/cmd/ since they reference cmd-internal symbols.

## Changes

### Files Created
- forge-cli/pkg/testrunner/test_results.go
- forge-cli/pkg/testrunner/journey_isolation.go
- forge-cli/pkg/testrunner/journey_isolation_test.go

### Files Modified
- forge-cli/internal/cmd/quality_gate.go
- forge-cli/internal/cmd/test.go
- forge-cli/internal/cmd/test_promote.go
- forge-cli/internal/cmd/integration_test.go

### Key Decisions
- Exported previously-unexported functions (writeUnitTestRawOutput -> WriteUnitTestRawOutput, etc.) to enable cross-package access since internal/cmd callers need to invoke them
- Kept TestTestingRunJourney_CommandRegistered in internal/cmd/journey_isolation_test.go because it references testCmd which is cmd-internal
- Used type aliases (type ContractFailure = testrunner.ContractFailure) during Phase A to keep code compiling at every intermediate step

## Test Results
- **Tests Executed**: Yes
- **Passed**: 678
- **Failed**: 0
- **Coverage**: 80.9%

## Acceptance Criteria
- [x] test_results.go is moved to pkg/testrunner/ with updated package declaration
- [x] journey_isolation.go and journey_isolation_test.go are moved to pkg/testrunner/
- [x] All imports referencing the old location are updated
- [x] go build ./... passes
- [x] go test ./... passes

## Notes
Coverage: 80.9% for internal/cmd, 61.4% for pkg/testrunner. The testrunner coverage is lower because test_results.go functions were previously tested via integration_test.go in cmd package; those tests now call testrunner.WriteUnitTestRawOutput directly.
