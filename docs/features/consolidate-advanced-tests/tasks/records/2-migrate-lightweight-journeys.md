---
status: "completed"
started: "2026-06-06 23:56"
completed: "2026-06-06 23:58"
time_spent: "~2m"
---

# Task Record: 2 迁移 error-handling 和 scope-resolution journeys

## Summary
Migrated error-handling and scope-resolution journeys from forge-cli/tests/ to tests/. Both journeys now use ForgeBinary init mode with correct import paths (forge-tests/testkit). All 10 tests pass (2 error-handling + 8 scope-resolution). Previous TC_008 blocking issue resolved by fix-1.

## Changes

### Files Created
- tests/error-handling/main_test.go
- tests/error-handling/error_handling_test.go
- tests/error-handling/contracts/step-1-task-errors.md
- tests/error-handling/contracts/step-2-forensic-errors.md
- tests/error-handling/contracts/step-3-submit-errors.md
- tests/scope-resolution/main_test.go
- tests/scope-resolution/scope_resolution_test.go
- tests/scope-resolution/contracts/step-1-scope-inference.md
- tests/scope-resolution/contracts/step-2-scope-dispatch.md

### Files Modified
无

### Key Decisions
- Used ForgeBinary init mode (_ = testkit.ForgeBinary) matching existing task-lifecycle pattern
- Contracts copied as-is (markdown documents, no modification needed)

## Test Results
- **Tests Executed**: Yes
- **Passed**: 10
- **Failed**: 0
- **Coverage**: 0.0%

## Acceptance Criteria
- [x] tests/error-handling/ contains migrated test files with import testkit "forge-tests/testkit"
- [x] tests/scope-resolution/ contains migrated test files with import testkit "forge-tests/testkit"
- [x] Both journeys' main_test.go use ForgeBinary init mode (_ = testkit.ForgeBinary)
- [x] just test includes both journeys and all tests pass

## Notes
Migration was completed in a previous session. This run verified all AC pass after fix-1 resolved the pre-existing TC_008 failure. compile, fmt, lint all clean.
