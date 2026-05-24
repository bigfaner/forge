---
status: "completed"
started: "2026-05-24 18:36"
completed: "2026-05-24 18:48"
time_spent: "~12m"
---

# Task Record: 2 加固依赖注入：合并 resolveTestDepsAndInjectReviewDoc + 更新 findFirstTestTaskIdx

## Summary
Merged resolveTestDepsAndInjectReviewDoc combining ResolveFirstTestDep + T-review-doc prepend into a single atomic operation. Updated findFirstTestTaskIdx to use findTaskIndexByPrefix('T-test-gen-journeys') instead of matching deprecated T-quick-gen-and-run prefix.

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/task/build.go
- forge-cli/pkg/task/build_test.go
- forge-cli/pkg/task/autogen_test.go

### Key Decisions
- resolveTestDepsAndInjectReviewDoc delegates to ResolveFirstTestDep internally, then conditionally prepends T-review-doc when needsEval=true
- findFirstTestTaskIdx simplified to single findTaskIndexByPrefix call, removing mode-specific branches and deprecated T-quick-gen-and-run matching
- Updated existing breakdown-mode integration test to assert on T-test-gen-journeys (first staged pipeline task) instead of T-eval-journey

## Test Results
- **Tests Executed**: Yes
- **Passed**: 10
- **Failed**: 0
- **Coverage**: 55.3%

## Acceptance Criteria
- [x] resolveTestDepsAndInjectReviewDoc(testTasks, idx, 'quick', true) returns deps containing T-review-doc
- [x] resolveTestDepsAndInjectReviewDoc(testTasks, idx, 'quick', false) returns deps without T-review-doc matching old ResolveFirstTestDep output
- [x] findFirstTestTaskIdx quick-mode branch uses findTaskIndexByPrefix(tasks, 'T-test-gen-journeys')
- [x] BuildIndex no longer has independent T-review-doc prepend operation
- [x] Integration test covers Quick mode full dependency chain
- [x] All existing tests pass

## Notes
无
