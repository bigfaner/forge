---
status: "completed"
started: "2026-05-26 13:16"
completed: "2026-05-26 13:22"
time_spent: "~6m"
---

# Task Record: 4 Split run-tests into per-surface-key serial tasks

## Summary
Split run-tests into per-surface-key serial tasks: wireRunTestChain centralizes serial chain wiring for both breakdown and quick modes. Single surface degenerates to no-suffix T-test-run. T-test-verify-regression depends on chain tail. InferType supports surface-key prefix matching.

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/task/autogen.go
- forge-cli/pkg/task/autogen_test.go
- forge-cli/pkg/task/build.go
- forge-cli/pkg/task/build_test.go
- forge-cli/pkg/task/autoconfig_test.go
- forge-cli/pkg/forgeconfig/detect.go

### Key Decisions
- wireRunTestChain helper centralizes serial chain wiring for both breakdown and quick modes
- Single surface detected via isSingleSurface helper for backward-compatible degradation
- Map-iteration-dependent tests use byID lookups instead of positional assertions to avoid flakiness

## Test Results
- **Tests Executed**: Yes
- **Passed**: 8
- **Failed**: 0
- **Coverage**: 87.4%

## Acceptance Criteria
- [x] AC1: surfaces { frontend: web, backend: api } with no execution-order produces T-test-run-backend before T-test-run-frontend
- [x] AC2: T-test-run-backend fails -> T-test-run-frontend blocked (serial dependency chain)
- [x] AC3: Single surface (surfaces: api) degrades to T-test-run (no suffix), same task ID and deps
- [x] AC4: Quick mode: T-test-gen-journeys is direct upstream of T-test-run-*
- [x] AC5: T-test-verify-regression depends on last run-test in execution-order

## Notes
All 5 acceptance criteria verified with dedicated test cases. Coverage at 87.4% (target 80%). No new code changes needed - implementation was already complete from prior attempt.
