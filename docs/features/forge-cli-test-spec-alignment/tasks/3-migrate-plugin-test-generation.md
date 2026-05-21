---
id: "3"
title: "Migrate tests/e2e/ Journey: test-generation"
priority: "P1"
estimated_time: "2h"
dependencies: ["1"]
scope: "backend"
breaking: false
type: "coding.refactor"
mainSession: false
---

# 3: Migrate tests/e2e/ Journey: test-generation

## Description

Migrate 7 test files into the `test-generation/` Journey directory — the largest single Journey in the plugin test suite.

**Source → Target:**
- `tests/e2e/quick_test_slim_cli_test.go` → `tests/test-generation/quick_test_slim_test.go`
- `tests/e2e/test_scripts_per_type_cli_test.go` → `tests/test-generation/test_scripts_per_type_test.go`
- `tests/e2e/features/test-knowledge-convention-driven/gen_test_scripts_cli_test.go` → `tests/test-generation/gen_test_scripts_test.go`
- `tests/e2e/features/test-knowledge-convention-driven/integration_cli_test.go` → `tests/test-generation/integration_test.go`
- `tests/e2e/features/test-knowledge-convention-driven/forge_commands_cli_test.go` → `tests/test-generation/forge_commands_test.go`
- `tests/e2e/features/test-knowledge-convention-driven/test_guide_cli_test.go` → `tests/test-generation/test_guide_test.go`
- Contracts: `tests/test-generation/contracts/step-1-task-index.md`, `step-2-gen-test-scripts.md`, `step-3-run-tests.md`

## Reference Files
- `docs/proposals/forge-cli-test-spec-alignment/proposal.md` — Journey mapping table
- `tests/testkit/` — Shared infrastructure (Task 1)

## Acceptance Criteria
- [ ] 7 test files migrated with correct package name and imports
- [ ] `contracts/` contains 3 Contract spec files with six-dimension declarations
- [ ] `main_test.go` initializes binary via testkit
- [ ] Tests pass: `go test ./tests/test-generation/... -tags=e2e -count=1`
- [ ] Duplicate test cases removed (check overlap between root and sub-package tests)

## Hard Rules
- Package name: `testgeneration`
- All test functions use `//go:build e2e` build tag
- 5 files from `test-knowledge-convention-driven/` sub-package must be flattened

## Implementation Notes
- Check for shared helpers within the `test-knowledge-convention-driven/` sub-package
- Some tests may overlap with `quick_test_slim_cli_test.go` or `test_scripts_per_type_cli_test.go` — deduplicate
- Contract specs reflect the 3-step pipeline: index → generate → run
