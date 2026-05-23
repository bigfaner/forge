---
status: "completed"
started: "2026-05-24 00:03"
completed: "2026-05-24 00:11"
time_spent: "~8m"
---

# Task Record: 2 Modify build.go routing for mixed features

## Summary
Modified build.go routing logic so mixed features (containing both doc and coding.* tasks) generate both T-review-doc AND test pipeline tasks. Changed needsReviewDoc() from ALL-doc to ANY-doc detection using CategoryForType. Added dependency injection so T-review-doc executes before the first test pipeline task for mixed features.

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/task/build.go
- forge-cli/pkg/task/build_test.go

### Key Decisions
- Used CategoryForType(t.Type) == CategoryDoc instead of exact TypeDoc match to cover doc.consolidate and doc.drift subtypes
- Added findFirstTestTaskIdx helper to locate the first test pipeline task for dependency injection across both breakdown and quick modes
- Injected T-review-doc as first dependency of first test task (prepended, not replaced) to preserve existing gate-based dependency chain

## Test Results
- **Tests Executed**: Yes
- **Passed**: 6
- **Failed**: 0
- **Coverage**: 90.6%

## Acceptance Criteria
- [x] needsReviewDoc() returns true when ANY non-auto-gen task has type doc
- [x] needsTestPipeline() unchanged
- [x] Both needsReviewDoc and needsTestPipeline can return true simultaneously for mixed features
- [x] Pure doc features: only T-review-doc generated
- [x] Pure code features: only test pipeline generated
- [x] Mixed features: both T-review-doc and test pipeline tasks generated
- [x] T-review-doc is dependency of first test pipeline task for mixed features
- [x] forge task index produces correct index.json for all three scenarios
- [x] Existing stage-gate logic remains unchanged
- [x] go build ./... and existing tests pass

## Notes
6 new tests added, 4 existing test cases updated in TestNeedsReviewDoc to reflect new ANY-doc behavior. The old 'mixed doc+feature does NOT need doc review' test case was inverted to match the new behavior.
