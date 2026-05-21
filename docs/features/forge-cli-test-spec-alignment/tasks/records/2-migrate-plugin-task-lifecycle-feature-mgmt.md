---
status: "completed"
started: "2026-05-20 22:13"
completed: "2026-05-20 22:24"
time_spent: "~11m"
---

# Task Record: 2 Migrate tests/e2e/ Journeys: task-lifecycle + feature-management

## Summary
Migrated 8 test files from tests/e2e/ into 2 Journey directories (task-lifecycle, feature-management). Created unified helpers using testkit package, merged fix-task-claim-priority into task-lifecycle, created contract specs with six-dimension declarations.

## Changes

### Files Created
- tests/task-lifecycle/main_test.go
- tests/task-lifecycle/task_lifecycle_test.go
- tests/task-lifecycle/task_record_test.go
- tests/task-lifecycle/contracts/step-1-task-claim.md
- tests/task-lifecycle/contracts/step-2-task-submit.md
- tests/feature-management/main_test.go
- tests/feature-management/feature_set_test.go
- tests/feature-management/cli_list_reverse_chronological_test.go
- tests/feature-management/proposal_status_lifecycle_test.go
- tests/feature-management/contracts/step-1-feature-set.md
- tests/feature-management/contracts/step-2-feature-list.md
- tests/feature-management/contracts/step-3-feature-status.md

### Files Modified
无

### Key Decisions
- Unified setupFeatureFixture helper with featureSlug parameter to replace per-feature tlhSetupFeatureFixture and setupFeatureFixture duplicates
- Used testkit.ParseBlock/HasField/FieldValue instead of per-file parseBlock/hasFieldWithValue/getFieldValue helpers
- Renumbered TC IDs in fix-task-claim-priority tests (TC-001..TC-006 -> TC-015..TC-020) to avoid collision with task-lifecycle-hardening tests (TC-001..TC-014)
- Renumbered TC IDs in task_record tests (TC-001..TC-012 -> TC-021..TC-032) for same reason
- Package names follow Hard Rules: tasklifecycle, featuremanagement

## Test Results
- **Tests Executed**: No
- **Passed**: 0
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] Both Journey directories created with correct package names
- [x] All 8 test files migrated (5 task-lifecycle, 3 feature-management)
- [x] Each Journey has main_test.go calling testkit, contracts/ with spec files
- [x] Tests compile: go test -tags=e2e -run=^$ ./task-lifecycle/... ./feature-management/...
- [x] Duplicate test cases removed (fix-task-claim-priority merged into task-lifecycle)

## Notes
This is a migration-only task. Source files in tests/e2e/ are NOT deleted (that would be a separate cleanup task). All new Journey files use testkit.ForgeBinary instead of local e2e.ForgeBinary, aligning with the new shared infrastructure from Task 1.
