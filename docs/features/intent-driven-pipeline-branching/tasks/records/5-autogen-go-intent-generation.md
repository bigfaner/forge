---
status: "completed"
started: "2026-05-29 16:40"
completed: "2026-05-29 17:04"
time_spent: "~24m"
---

# Task Record: 5 autogen.go — intent-driven task generation and wiring

## Summary
Implemented intent-driven task generation and dependency wiring in autogen.go. Modified GetBreakdownTestTasks, GetQuickTestTasks, resolveBreakdownDeps, resolveQuickDeps, ResolveFirstTestDep to accept and respond to intent parameter. Updated GenerateTestTasks in build.go to pass intent through. Added 20 new unit tests covering all 5 wiring scenarios and zero business task edge cases.

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/task/autogen.go
- forge-cli/pkg/task/autogen_test.go
- forge-cli/pkg/task/build.go
- forge-cli/pkg/task/autoconfig_test.go
- forge-cli/pkg/task/build_test.go
- forge-cli/scripts/version.txt

### Key Decisions
- Added isSkipTestIntent() helper to centralize refactor/cleanup intent detection
- GetBreakdownTestTasks/GetQuickTestTasks skip auto.Test.* block for refactor/cleanup but preserve auto.Validation.*, auto.ConsolidateSpecs.*, auto.CleanCode.* tasks
- resolveBreakdownDeps/resolveQuickDeps do not wire downstream tasks to run-test for refactor/cleanup (no lastRunID), leaving wiring to ResolveFirstTestDep
- ResolveFirstTestDep gains intent awareness: for refactor/cleanup it wires validate-code (breakdown) or clean-code (quick) directly to last business task, with zero-business-task protection
- GenerateTestTasks and resolveTestDepsAndInjectReviewDoc in build.go propagate intent to autogen functions

## Test Results
- **Tests Executed**: Yes
- **Passed**: 20
- **Failed**: 0
- **Coverage**: 87.8%

## Acceptance Criteria
- [x] GetBreakdownTestTasks/GetQuickTestTasks skip test tasks (gen-journeys/gen-contracts/gen-scripts/run-tests) for refactor/cleanup but keep validate-code/clean-code/consolidate-specs
- [x] resolveBreakdownDeps/resolveQuickDeps sense intent: refactor/cleanup wires downstream tasks to last business task, new-feature unchanged
- [x] Zero business task protection: refactor/cleanup with empty business tasks produces no downstream tasks, no panic
- [x] intent: new-feature pipeline behavior identical to current (backward compat)
- [x] Unit tests cover all 5 wiring scenarios + zero business task edge case

## Notes
Version bumped from 5.14.2 to 5.15.0 (minor: new feature). All existing tests pass with the new intent parameter (backward compatible with empty string default).
