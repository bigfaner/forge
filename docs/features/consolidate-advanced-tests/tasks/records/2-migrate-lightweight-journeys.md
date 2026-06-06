---
status: "blocked"
started: "2026-06-06 23:44"
completed: "N/A"
time_spent: ""
---

# Task Record: 2 迁移 error-handling 和 scope-resolution journeys

## Summary
Migrated error-handling and scope-resolution journeys from forge-cli/tests/ to tests/. Updated import paths from forge-cli/tests/testkit to forge-tests/testkit, rewrote main_test.go to use ForgeBinary init pattern. All error-handling tests pass (2/2). Scope-resolution has 7/8 passing with 1 pre-existing failure (TC_008 fails identically at source location).

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
- Used ForgeBinary init pattern (_ = testkit.ForgeBinary) instead of manual binary build in main_test.go
- Verified target testkit already contains RunCLIExitCode, ProjectRoot, ReadProjectFile functions - no testkit changes needed
- TC_008 pre-existing failure confirmed at source (forge-cli/tests/) - not introduced by migration

## Test Results
- **Tests Executed**: Yes
- **Passed**: 9
- **Failed**: 1
- **Coverage**: 0.0%

## Acceptance Criteria
- [x] tests/error-handling/ contains migrated test files with import testkit 'forge-tests/testkit'
- [x] tests/scope-resolution/ contains migrated test files with import testkit 'forge-tests/testkit'
- [x] Both journeys' main_test.go use ForgeBinary init pattern (_ = testkit.ForgeBinary)
- [ ] just test includes both journeys and all pass

## Notes
TC_008 (TestTC_008_FrontendOnlyTaskScopeMarkedAsFrontend) fails in BOTH source and target locations. Verified by running the same test in forge-cli/tests/scope-resolution/ where it also fails. This is a pre-existing issue in the breakdown-tasks SKILL.md content, not a migration regression.
