---
status: "blocked"
started: "2026-05-16 01:18"
completed: "N/A"
time_spent: ""
---

# Task Record: T-quick-5 Verify Quick E2E Regression

## Summary
Verified e2e regression suite and fixed missing justfile recipes. Added test-e2e, e2e-compile, e2e-discover, and e2e-setup recipes to justfile for go-test profile (these were removed during Playwright-to-Go migration but never replaced). Fixed pre-existing TC-002 test failure in justfile-canonical-e2e by creating the missing feature directory fixture. All 19 justfile-canonical-e2e tests now pass. 20 pre-existing cli_lean_output failures remain (environment-dependent: forge task claim requires pending tasks).

## Changes

### Files Created
无

### Files Modified
- justfile
- tests/e2e/justfile-canonical-e2e/justfile_canonical_e2e_cli_test.go

### Key Decisions
- Added go-test-specific e2e recipes to justfile instead of restoring Playwright versions, matching the project's current go-test profile
- Fixed TC-002 by creating feature directory fixture rather than relaxing test assertions, preserving the test's intent of verifying just delegation

## Test Results
- **Tests Executed**: No
- **Passed**: 35
- **Failed**: 20
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] just test-e2e recipe exists and runs go-test e2e suite
- [x] justfile-canonical-e2e tests all pass (19/19)
- [x] Pre-existing cli_lean_output failures not introduced by this task
- [x] simplify-e2e-tests feature tests pass (3/3 direct + 12/12 per-type)

## Notes
BLOCKED by 20 pre-existing cli_lean_output test failures (verified against pre-change code). These tests require a clean fixture with pending tasks -- they fail in the current project state where all tasks are in_progress or completed. Not caused by this task's changes. All feature-specific tests pass. To unblock: either fix cli_lean_output tests with proper test fixtures, or exclude them from the regression suite.
