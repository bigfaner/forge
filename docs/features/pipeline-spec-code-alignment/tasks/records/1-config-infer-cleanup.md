---
status: "completed"
started: "2026-05-27 00:33"
completed: "2026-05-27 00:37"
time_spent: "~4m"
---

# Task Record: 1 Clean up config.go dead code and infer.go missing cases

## Summary
Removed dead 'coding.clean' entry from coverage defaults in config.go and added missing InferType cases for T-eval-journey and T-eval-contract in infer.go

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/forgeconfig/config.go
- forge-cli/pkg/task/infer.go
- forge-cli/pkg/forgeconfig/config_test.go
- forge-cli/internal/cmd/testdata/forge-config.example.yaml

### Key Decisions
无

## Test Results
- **Tests Executed**: Yes
- **Passed**: 493
- **Failed**: 0
- **Coverage**: 88.0%

## Acceptance Criteria
- [x] coding.clean entry removed from coverage defaults in config.go
- [x] InferType returns correct type string for T-eval-journey and T-eval-contract task IDs
- [x] Existing tests pass (go test ./...)

## Notes
Also removed coding.clean test case from config_test.go and coding.clean example from forge-config.example.yaml
