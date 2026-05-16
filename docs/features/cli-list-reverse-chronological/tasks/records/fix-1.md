---
status: "completed"
started: "2026-05-16 09:54"
completed: "2026-05-16 10:00"
time_spent: "~6m"
---

# Task Record: fix-1 Fix: T-quick-2 e2e test failures

## Summary
Diagnosed TC-001, TC-004, TC-005 e2e test failures as caused by stale forge binary. The sort implementation in proposal.go and feature.go was already correct but the installed CLI binary was not rebuilt after the code changes. Running go install resolved all 3 failures. No source code changes were needed.

## Changes

### Files Created
无

### Files Modified
无

### Key Decisions
- No code modification required -- root cause was stale binary, not incorrect logic

## Test Results
- **Tests Executed**: Yes
- **Passed**: 20
- **Failed**: 0
- **Coverage**: 90.2%

## Acceptance Criteria
- [x] TC-001 proposal list sorted by created descending
- [x] TC-004 feature list sorted by manifest mtime descending
- [x] TC-005 missing manifest sorts to end

## Notes
All 6 e2e tests pass after reinstall. Unit tests all pass (20/20 packages).
