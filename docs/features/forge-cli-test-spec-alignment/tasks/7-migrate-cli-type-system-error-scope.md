---
id: "7"
title: "Migrate forge-cli/tests/e2e/ Journeys: task-type-system + error-handling + scope-resolution"
priority: "P1"
estimated_time: "2h"
dependencies: []
scope: "backend"
breaking: false
type: "coding.refactor"
mainSession: false
---

# 7: Migrate forge-cli/tests/e2e/ Journeys: task-type-system + error-handling + scope-resolution

## Description

Migrate 6 test files into 3 Journey directories.

**task-type-system Journey:**
- `task_types_cli_test.go` → `forge-cli/tests/task-type-system/task_types_test.go`
- `task_types_dispatch_cli_test.go` → `forge-cli/tests/task-type-system/task_types_dispatch_test.go`
- `features/task-type-refinement/` → `forge-cli/tests/task-type-system/task_type_refinement_test.go`
- Contracts: `step-1-list-types.md`, `step-2-validate.md`, `step-3-dispatch.md`

**error-handling Journey:**
- `error_handling_cli_test.go` → `forge-cli/tests/error-handling/error_handling_test.go`
- Contracts: `step-1-task-errors.md`, `step-2-forensic-errors.md`, `step-3-submit-errors.md`

**scope-resolution Journey:**
- `scope_resolution_cli_test.go` → `forge-cli/tests/scope-resolution/scope_resolution_test.go`
- Contracts: `step-1-scope-inference.md`, `step-2-scope-dispatch.md`

## Reference Files
- `forge-cli/tests/e2e/testkit/` — Already existing testkit
- All source files listed above

## Acceptance Criteria
- [ ] 3 Journey directories created with 6 migrated test files
- [ ] Each Journey has `contracts/` with spec files
- [ ] Each Journey has `main_test.go` via testkit
- [ ] Tests pass: `go test ./forge-cli/tests/task-type-system/... ./forge-cli/tests/error-handling/... ./forge-cli/tests/scope-resolution/... -tags=e2e -count=1`

## Hard Rules
- Package names: `tasktypesystem`, `errorhandling`, `scoperesolution`
- All test functions use `//go:build e2e` build tag
- task_types_dispatch_cli_test.go has 20 tests — verify all pass after migration

## Implementation Notes
- task_types_dispatch tests use task prompt generation — ensure forge binary path resolves correctly
- error_handling tests create corrupted indexes and invalid states — ensure cleanup works in new package location
