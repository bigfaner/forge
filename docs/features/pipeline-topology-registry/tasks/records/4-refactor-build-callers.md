---
status: "blocked"
started: "2026-05-29 23:26"
completed: "N/A"
time_spent: ""
---

# Task Record: 4 Refactor build.go callers and delete legacy functions

## Summary
Deleted 5 legacy functions from autogen.go: ResolveFirstTestDep, GetReviewDocTask, ResolveReviewDocDep, findTaskIndexOrPanic, ResolveDriftFallbackDep. go build ./... passes. Blocked because forge task submit runs go vet (which compiles test files) and tests reference deleted functions — requires task 6 (test updates) to resolve.

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/task/autogen.go

### Key Decisions
- Kept GetBreakdownTestTasks, GetQuickTestTasks, findTaskIndex, findTaskIndexByPrefix as bridge functions since they are used by test files (deferred to task 6)
- No changes to build.go — steps 7/7.5/7.6 were already refactored by tasks 1-3

## Test Results
- **Tests Executed**: Yes
- **Passed**: 0
- **Failed**: 0
- **Coverage**: 0.0%

## Acceptance Criteria
- [x] Step 7/7.5 unified: single call to registry-driven GenerateTestTasks with businessTasks + existingTasks
- [x] T-review-doc generated as part of registry pipeline, gated by CondHasDocTasks
- [x] Step 7.6 (ResolveDriftFallbackDep) removed
- [x] All functions in deletion table deleted
- [x] needsTestPipeline preserved for stage-gate generation control
- [x] needsReviewDoc preserved for doc task criteria extraction
- [x] go build ./... passes

## Notes
All 7 acceptance criteria met. Blocked on forge task submit quality gate: go vet compiles test files which reference deleted functions. Task 6 must update autogen_test.go and build_test.go to remove calls to ResolveFirstTestDep, GetReviewDocTask, ResolveReviewDocDep, findTaskIndexOrPanic, ResolveDriftFallbackDep.
