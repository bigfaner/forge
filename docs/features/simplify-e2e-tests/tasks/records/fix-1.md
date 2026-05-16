---
status: "completed"
started: "2026-05-16 00:48"
completed: "2026-05-16 00:48"
time_spent: ""
---

# Task Record: fix-1 Fix: TC-003/TC-004 go.mod path mismatch

## Summary
Fixed TC-003/TC-004 go.mod path mismatch: changed working directory from project root to tests/e2e/ (the Go module root) and package pattern from ./tests/e2e/... to ./... . Also fixed pre-existing bug in resolveQuickDeps where T-quick-6 dependency was assigned to wrong task index in per-type mode, and updated test expectations to account for T-quick-6 (drift-detection) task.

## Changes

### Files Created
无

### Files Modified
- tests/e2e/features/simplify-e2e-tests/simplify_e2e_tests_cli_test.go
- forge-cli/pkg/task/testgen.go
- forge-cli/pkg/task/testgen_test.go

### Key Decisions
- Used e2eRoot(t) instead of projectRoot(t) for cmd.Dir since go.mod lives at tests/e2e/
- Changed go test pattern from ./tests/e2e/... to ./... to match the module-relative path
- Fixed resolveQuickDeps T-quick-6 dependency by searching for ID instead of using hardcoded index formula that was wrong in per-type mode

## Test Results
- **Tests Executed**: Yes
- **Passed**: 21
- **Failed**: 0
- **Coverage**: 0.0%

## Acceptance Criteria
- [x] TC-003 compiles e2e test suite from correct module root
- [x] TC-004 runs e2e tests from correct module root
- [x] All unit tests pass (including pre-existing failures now fixed)

## Notes
Pre-existing bug in resolveQuickDeps caused T-quick-6 dependency to overwrite T-quick-5 dependency in per-type mode because len(profiles)*4+1 was wrong index. Fixed by searching for T-quick-6 ID instead of computing index.
