---
status: "completed"
started: "2026-05-23 00:17"
completed: "2026-05-23 01:51"
time_spent: "~1h 34m"
---

# Task Record: 3 Move forge task group to task/ subdirectory

## Summary
Moved all forge task subcommand files (~12 source + 12 test files) from forge-cli/internal/cmd/ to new forge-cli/internal/cmd/task/ subdirectory. Created base/ sub-package to break circular dependency (cmd <-> task), updated root.go to use taskpkg.Register(), created testbridge.go for cross-package test access, and updated all test files to work with the new package structure.

## Changes

### Files Created
- forge-cli/internal/cmd/base/errors.go
- forge-cli/internal/cmd/base/output.go
- forge-cli/internal/cmd/task/register.go
- forge-cli/internal/cmd/task/testbridge.go
- forge-cli/internal/cmd/task/testing_helpers_test.go
- forge-cli/internal/cmd/task/testmain_test.go

### Files Modified
- forge-cli/internal/cmd/root.go
- forge-cli/internal/cmd/errors.go
- forge-cli/internal/cmd/output.go
- forge-cli/internal/cmd/integration_test.go
- forge-cli/internal/cmd/characterization_test.go
- forge-cli/internal/cmd/feature_test.go
- forge-cli/internal/cmd/output_contract_test.go
- forge-cli/internal/cmd/root_test.go
- forge-cli/internal/cmd/task/add.go
- forge-cli/internal/cmd/task/check_deps.go
- forge-cli/internal/cmd/task/claim.go
- forge-cli/internal/cmd/task/index.go
- forge-cli/internal/cmd/task/list_types.go
- forge-cli/internal/cmd/task/migrate.go
- forge-cli/internal/cmd/task/query.go
- forge-cli/internal/cmd/task/reopen.go
- forge-cli/internal/cmd/task/status.go
- forge-cli/internal/cmd/task/submit.go
- forge-cli/internal/cmd/task/task_parent.go
- forge-cli/internal/cmd/task/transition.go
- forge-cli/internal/cmd/task/validate_index.go
- forge-cli/internal/cmd/task/add_cmd_test.go
- forge-cli/internal/cmd/task/claim_test.go
- forge-cli/internal/cmd/task/claim_integration_test.go
- forge-cli/internal/cmd/task/check_deps_test.go
- forge-cli/internal/cmd/task/index_test.go
- forge-cli/internal/cmd/task/list_types_test.go
- forge-cli/internal/cmd/task/migrate_test.go
- forge-cli/internal/cmd/task/query_test.go
- forge-cli/internal/cmd/task/reopen_test.go
- forge-cli/internal/cmd/task/status_test.go
- forge-cli/internal/cmd/task/submit_test.go
- forge-cli/internal/cmd/task/validate_index_test.go

### Key Decisions
- Created internal/cmd/base/ package to break circular dependency between cmd and task packages
- Used re-export pattern in errors.go and output.go (type aliases + var assignments) for backward compatibility
- Inlined Debugf in output.go instead of re-exporting to preserve variadic function semantics
- Created testbridge.go (non-_test.go file) to export internal symbols for cross-package test access
- Used Cmd.SetArgs(['subcommand', ...]) pattern for cobra command testing in task package tests
- Duplicated test helpers (setupFullProject, captureOutput, etc.) in task/testing_helpers_test.go rather than creating shared test package

## Test Results
- **Tests Executed**: Yes
- **Passed**: 140
- **Failed**: 0
- **Coverage**: 70.7%

## Acceptance Criteria
- [x] All task subcommand files moved to cmd/task/
- [x] No circular dependencies between cmd and task packages
- [x] go build ./... passes
- [x] go vet ./... passes
- [x] All existing tests pass

## Notes
The testbridge.go approach exports internal symbols for cross-package testing. This is a pragmatic tradeoff -- the alternative would have been to split integration_test.go into separate packages or export all internal functions. The base package extraction is minimal (errors.go + output.go) and could be expanded if other packages need similar utilities.
