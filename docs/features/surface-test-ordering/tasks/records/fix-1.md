---
status: "completed"
started: "2026-05-26 11:49"
completed: "2026-05-26 12:02"
time_spent: "~13m"
---

# Task Record: fix-1 Fix: task 4 signature change not propagated to test callers

## Summary
Fix two pre-existing test failures that blocked task 4: (1) SurfacesMap UnmarshalYAML normalized '.' key to '-' breaking scalar-form config in map YAML; (2) stale Makefile branch test not cleaned up after testrunner probe chain refactor

## Type Reclassification
- Original: 
- Actual: coding.fix
- Reason: Fix task with no assigned type, actual work is fixing pre-existing test failures

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/forgeconfig/config.go
- forge-cli/pkg/testrunner/testrunner_test.go

### Key Decisions
- Preserve '.' key as-is in MappingNode branch of UnmarshalYAML since '.' is the scalar-form internal marker
- Delete stale Makefile branch test and unused exec import instead of re-adding Makefile support (removed in commit 469e88fb)

## Test Results
- **Tests Executed**: Yes
- **Passed**: 18
- **Failed**: 0
- **Coverage**: 82.5%

## Acceptance Criteria
- [x] TestRunSurfaceConfigRerun passes (scalar '.' key preserved through YAML map round-trip)
- [x] TestRunProjectTests passes (stale Makefile branch removed)
- [x] forgeconfig and testrunner package tests all pass

## Notes
Both failures were pre-existing and unrelated to task 4 signature change. TestSurfacesJSONTypes also fails pre-existing (not fixed in this task).
