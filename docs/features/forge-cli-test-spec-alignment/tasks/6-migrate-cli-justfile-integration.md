---
id: "6"
title: "Migrate forge-cli/tests/e2e/ Journey: justfile-integration"
priority: "P1"
estimated_time: "2h"
dependencies: []
scope: "backend"
breaking: false
type: "coding.refactor"
mainSession: false
---

# 6: Migrate forge-cli/tests/e2e/ Journey: justfile-integration

## Description

Migrate 4 justfile-related test files into the `justfile-integration/` Journey directory.

**Source → Target:**
- `init_justfile_cli_test.go` → `forge-cli/tests/justfile-integration/init_justfile_test.go`
- `justfile_execution_cli_test.go` → `forge-cli/tests/justfile-integration/execution_test.go`
- `justfile_forge_detection_cli_test.go` → `forge-cli/tests/justfile-integration/forge_detection_test.go`
- `justfile_mixed_cli_cli_test.go` → `forge-cli/tests/justfile-integration/mixed_cli_test.go`
- Contracts: `forge-cli/tests/justfile-integration/contracts/step-1-init.md`, `step-2-detection.md`, `step-3-execution.md`, `step-4-mixed-scope.md`

## Reference Files
- `forge-cli/tests/e2e/testkit/` — Already existing testkit
- All source files listed above

## Acceptance Criteria
- [ ] 4 test files migrated with correct package name and imports
- [ ] `contracts/` contains 4 Contract spec files
- [ ] `main_test.go` initializes binary via testkit
- [ ] Tests pass: `go test ./forge-cli/tests/justfile-integration/... -tags=e2e -count=1`

## Hard Rules
- Package name: `justfileintegration`
- All test functions use `//go:build e2e` build tag

## Implementation Notes
- These are the largest test files (justfile_forge_detection has 20 tests, justfile_mixed_cli has 20 tests)
- The justfile tests create temporary project directories — ensure testkit paths resolve correctly in new package location
