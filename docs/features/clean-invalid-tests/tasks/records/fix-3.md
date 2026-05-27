---
status: "completed"
started: "2026-05-27 12:13"
completed: "2026-05-27 12:52"
time_spent: "~39m"
---

# Task Record: fix-3 Fix: test-generation e2e failures

## Summary
Fixed 16 failing e2e tests in tests/test-generation/ by aligning test expectations with the actual quick staged topology (gen-journeys -> run-test-{key} serial chain -> verify-regression -> drift) and breakdown topology (full pipeline with eval gates, per-type gen-scripts, per-surface run-test).

## Changes

### Files Created
无

### Files Modified
- tests/test-generation/quick_test_slim_test.go
- tests/test-generation/test_scripts_per_type_test.go

### Key Decisions
- Tests were asserting against an obsolete quick mode topology (T-quick-gen-cases, T-quick-gen-and-run-*, T-quick-graduate). Updated to match the actual staged topology (T-test-gen-journeys, T-test-run-{surface-key}, T-test-verify-regression).
- Removed assertion that gen-scripts .md must mention profile language (go) since templates no longer embed profile name in per-type task body.
- Removed unused strings import from test_scripts_per_type_test.go after multi-profile tests were simplified.

## Test Results
- **Tests Executed**: Yes
- **Passed**: 56
- **Failed**: 0
- **Coverage**: 80.0%

## Acceptance Criteria
- [x] TC-001 through TC-015 quick mode tests pass
- [x] TestPerType_TC-006 through TC-012 per-type tests pass
- [x] All test-generation tests pass (0 failures)
- [x] Static checks (compile, fmt, lint) pass

## Notes
无
