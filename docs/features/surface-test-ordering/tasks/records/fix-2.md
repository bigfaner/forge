---
status: "completed"
started: "2026-05-26 13:04"
completed: "2026-05-26 13:15"
time_spent: "~11m"
---

# Task Record: fix-2 Fix: TestSurfacesTypes and TestSurfacesJSONTypes fail with empty output

## Summary
Fix map-iteration-order flakiness in TestGetBreakdownTestTasks_PerType_TwoTypes by adding sort.Strings to SurfaceTypes() in forgeconfig/detect.go, and updating test assertions to match the new deterministic alphabetical order (api before tui)

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/forgeconfig/detect.go
- forge-cli/pkg/task/autogen_test.go

### Key Decisions
- Sort SurfaceTypes() output alphabetically (sort.Strings) to guarantee deterministic gen-test-scripts task ordering regardless of Go map iteration order
- Updated test assertions to expect api before tui (alphabetical) rather than the previous unstable order

## Test Results
- **Tests Executed**: Yes
- **Passed**: 18
- **Failed**: 0
- **Coverage**: 82.5%

## Acceptance Criteria
- [x] TestGetBreakdownTestTasks_PerType_TwoTypes passes reliably (10/10 runs)
- [x] Full test suite passes with 0 failures
- [x] go vet and go build pass cleanly

## Notes
Root cause: SurfaceTypes() iterated a Go map without sorting, producing non-deterministic gen-test-scripts task ordering. The two originally reported tests (TestSurfacesTypes, TestSurfacesJSONTypes) were already passing; the actual flaky failure was TestGetBreakdownTestTasks_PerType_TwoTypes in pkg/task.
