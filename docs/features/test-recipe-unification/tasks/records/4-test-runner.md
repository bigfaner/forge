---
status: "completed"
started: "2026-05-24 21:51"
completed: "2026-05-24 21:55"
time_spent: "~4m"
---

# Task Record: 4 Update test runner probe chain and journey isolation

## Summary
Updated RunProjectTests probe chain to unit-test -> test -> go test with fallback; migrated journey_isolation.go from `just e2e-test` to `just test <journey>`

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/testrunner/testrunner.go
- forge-cli/pkg/testrunner/journey_isolation.go
- forge-cli/pkg/testrunner/testrunner_test.go
- forge-cli/pkg/testrunner/journey_isolation_test.go

### Key Decisions
- Probe chain: unit-test -> test -> go test -> npm -> pytest (removed Makefile branch per proposal spec)
- RunProjectTests retains fallback: first matching recipe wins
- ExecuteJourneyInIsolation passes journey name as positional arg to `just test` matching recipe signature `test journey=''`

## Test Results
- **Tests Executed**: Yes
- **Passed**: 27
- **Failed**: 0
- **Coverage**: 74.3%

## Acceptance Criteria
- [x] RunProjectTests() probe chain: unit-test -> test -> go test (HasRecipe fallback)
- [x] journey_isolation.go calls `just test <journeyName>` instead of `just e2e-test`
- [x] RunProjectTests retains fallback mechanism (first matching recipe wins)
- [x] Gate sequence path and RunProjectTests path behave independently

## Notes
Removed Makefile probe branch from RunProjectTests per proposal spec which only defines unit-test -> test -> go test chain. All existing tests pass.
