---
status: "completed"
started: "2026-05-22 16:50"
completed: "2026-05-22 17:03"
time_spent: "~13m"
---

# Task Record: fix-1 fix unit-test: just test failure in quality gate

## Summary
Fix 15 pre-existing unit test failures in forge-cli/internal/cmd caused by AIError/state machine refactoring — subprocess tests were swallowing errors with `_ = func()` instead of propagating them through Exit()

## Changes

### Files Created
无

### Files Modified
- forge-cli/internal/cmd/characterization_test.go
- forge-cli/internal/cmd/forensic_test.go
- forge-cli/internal/cmd/integration_test.go
- forge-cli/internal/cmd/reopen_test.go

### Key Decisions
- All failing tests used subprocess pattern where errors were ignored with `_ = func()`; fix propagates errors via Exit() for non-zero exit
- TestValidateQualityGate_FailingCompileGate used panic() instead of Exit(); wrapped in recover/defer to convert panic to proper Exit() output
- TestReopen_WithLock_SaveIndexError assertion had casing mismatch ('failed to acquire lock' vs 'Failed to acquire lock')

## Test Results
- **Tests Executed**: Yes
- **Passed**: 1300
- **Failed**: 0
- **Coverage**: 80.5%

## Acceptance Criteria
- [x] All previously failing tests pass
- [x] No regressions in existing tests

## Notes
无
