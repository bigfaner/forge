---
status: "completed"
started: "2026-05-20 21:59"
completed: "2026-05-20 22:12"
time_spent: "~13m"
---

# Task Record: 1 Extract testkit + restructure tests/e2e/ Go module

## Summary
Extract shared test infrastructure into testkit package: moved go.mod from tests/e2e/ to tests/ (module renamed to forge-tests), created tests/testkit/ with exported ForgeBinary var, ForgeCmd(), ParseBlock, HasField, HasNoField, FieldIndex, FieldValue, RunCLI, RunCLIRaw, WithRetry. Deleted tests/e2e/main_test.go and helpers_test.go. Updated all import paths from e2e-tests to forge-tests/e2e. Replaced forgeBinary (unexported) with ForgeBinary (exported) in root test files. All packages compile successfully.

## Changes

### Files Created
- tests/go.mod
- tests/go.sum
- tests/testkit/forge_binary.go
- tests/testkit/helpers.go

### Files Modified
- tests/e2e/forge_binary.go
- tests/e2e/feature_set_command_cli_test.go
- tests/e2e/quality_gate_fix_task_loop_breaker_cli_test.go
- tests/e2e/quick_test_slim_cli_test.go
- tests/e2e/task_lifecycle_hardening_cli_test.go
- tests/e2e/task_record_immutability_cli_test.go
- tests/e2e/test_scripts_per_type_cli_test.go
- tests/e2e/features/cli-list-reverse-chronological/cli_list_reverse_chronological_cli_test.go
- tests/e2e/features/fix-task-claim-priority/fix_task_claim_priority_cli_test.go
- tests/e2e/features/proposal-status-lifecycle/proposal_status_lifecycle_cli_test.go
- tests/e2e/features/task-type-refinement/task_type_refinement_cli_test.go
- tests/e2e/features/test-knowledge-convention-driven/test_helpers_test.go

### Key Decisions
- Kept tests/e2e/forge_binary.go in place (provides ForgeBinary to root package e2e) alongside tests/testkit/forge_binary.go (for future journey packages to import)
- Import path for features sub-packages changed from e2e-tests to forge-tests/e2e because module root moved from tests/e2e/ to tests/
- Replaced forgeBinary (unexported, was set in deleted TestMain) with ForgeBinary (exported, set in init()) in root test files -- no logic change, same runtime value
- Deleted helpers_test.go because all its functions were copied to testkit/helpers.go with exported names, and no root test file directly used its unexported functions

## Test Results
- **Tests Executed**: No
- **Passed**: 0
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] tests/go.mod exists with module name forge-tests
- [x] tests/testkit/forge_binary.go exports ForgeBinary var and ForgeCmd() func
- [x] tests/testkit/helpers.go exports ParseBlock, HasField, FieldValue, WithRetry as public functions
- [x] tests/testkit/ compiles: go build ./tests/testkit/...
- [x] Existing root-level test files still compile

## Notes
Refactoring task -- e2e tests require forge binary to actually run. Verified all packages compile with go vet -tags=e2e and go test -tags=e2e dry-run. No behavior changes.
