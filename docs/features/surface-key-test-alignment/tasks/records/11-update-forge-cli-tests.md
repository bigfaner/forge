---
status: "completed"
started: "2026-06-06 13:40"
completed: "2026-06-06 13:49"
time_spent: "~9m"
---

# Task Record: 11 Update forge-cli test files for surface-key naming

## Summary
Updated forge-cli test files to reflect per-surface-key naming for gen-test-scripts expansion: renamed PerType->PerKey test names, added key!=type test cases verifying gen-scripts uses surface-key (backend, frontend) not surface-type (api, web), updated comments and prompt test for key!=type scenario

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/task/autogen_test.go
- forge-cli/pkg/task/pipeline_test.go
- forge-cli/pkg/task/build_test.go
- forge-cli/pkg/prompt/prompt_test.go

### Key Decisions
- Renamed 8 test functions from PerType to PerKey to reflect actual expansion semantics
- Added TestGetBreakdownTestTasks_KeyNotType_GenScriptsUsesKeyNaming to verify gen-scripts IDs use surface-key when key!=type
- Added TestSynthesize_GenScripts_KeyNotType_UsesTypeArg to verify prompt template correctly uses surface-key in ID and surface-type in --type arg
- Kept matchTypeSuffixedID test but added comment noting gen-test-scripts no longer uses per-surface-type

## Test Results
- **Tests Executed**: Yes
- **Passed**: 162
- **Failed**: 0
- **Coverage**: 86.3%

## Acceptance Criteria
- [x] All test cases referencing gen-test-scripts-{type} naming updated to {key} naming
- [x] Test fixtures and expectations updated to per-surface-key expansion
- [x] go test ./... all passing

## Notes
Pipeline registry already changed to per-surface-key in task 1. This task aligned test naming and added key!=type regression coverage.
