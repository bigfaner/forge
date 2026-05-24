---
status: "completed"
started: "2026-05-24 01:17"
completed: "2026-05-24 01:22"
time_spent: "~5m"
---

# Task Record: fix-2 Fix: pre-existing pkg/task test failures blocking fix-1

## Summary
Fixed 3 pre-existing test failures in pkg/task (TestBuildIndex_MixedFeature_*). Root causes: (1) tests expected T-quick-gen-and-run-* task IDs but GetQuickTestTasks was refactored to generate T-test-gen-journeys-* pipeline; (2) build.go mixed-feature logic only wrote testTasks[0] deps back to index even when firstTestIdx pointed to a different task (breakdown mode).

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/task/build.go
- forge-cli/pkg/task/build_test.go

### Key Decisions
- Updated findFirstTestTaskIdx result to be computed once and reused for both T-review-doc injection and index write-back, fixing the stale-deps bug in breakdown mode
- Updated tests to match current pipeline task ID format (T-test-gen-journeys-*) instead of legacy T-quick-gen-and-run-* format

## Test Results
- **Tests Executed**: Yes
- **Passed**: 750
- **Failed**: 0
- **Coverage**: 91.9%

## Acceptance Criteria
- [x] TestBuildIndex_MixedFeature_GeneratesBothPipelines passes
- [x] TestBuildIndex_MixedFeature_ReviewDocBeforeTestPipeline passes
- [x] TestBuildIndex_MixedFeature_BreakdownMode_GenJourneysDependsOnReviewDoc passes
- [x] No regression in full pkg/task test suite (race-safe)

## Notes
Two distinct root causes: (1) test assertions out of sync with GetQuickTestTasks refactor; (2) build.go only persisted testTasks[0] deps to index, but findFirstTestTaskIdx in breakdown mode returns index of T-eval-journey (not 0). Fixed by computing firstTestIdx once and using it consistently for both injection and index write-back.
