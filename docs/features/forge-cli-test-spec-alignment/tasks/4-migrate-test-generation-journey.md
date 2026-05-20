---
id: "4"
title: "Migrate test-generation Journey"
priority: "P1"
estimated_time: "2h"
dependencies: ["1"]
scope: "backend"
breaking: false
type: "coding.refactor"
mainSession: false
---

# 4: Migrate test-generation Journey

## Description

Migrate 7 test files into the `test-generation/` Journey directory. This is the largest Journey, covering the test generation pipeline (task index → gen-test-scripts → run tests).

**Source → Target mapping:**
- `tests/e2e/quick_test_slim_cli_test.go` → `tests/test-generation/quick_test_slim_test.go`
- `tests/e2e/test_scripts_per_type_cli_test.go` → `tests/test-generation/test_scripts_per_type_test.go`
- `tests/e2e/features/test-knowledge-convention-driven/gen_test_scripts_cli_test.go` → `tests/test-generation/gen_test_scripts_test.go`
- `tests/e2e/features/test-knowledge-convention-driven/integration_cli_test.go` → `tests/test-generation/integration_test.go`
- `tests/e2e/features/test-knowledge-convention-driven/forge_commands_cli_test.go` → `tests/test-generation/forge_commands_test.go`
- `tests/e2e/features/test-knowledge-convention-driven/test_guide_cli_test.go` → `tests/test-generation/test_guide_test.go`
- New: `tests/test-generation/contracts/step-1-task-index.md`
- New: `tests/test-generation/contracts/step-2-gen-test-scripts.md`
- New: `tests/test-generation/contracts/step-3-run-tests.md`
- New: `tests/test-generation/main_test.go`

## Reference Files
- `docs/proposals/forge-cli-test-spec-alignment/proposal.md` — Journey mapping table
- `tests/testkit/` — Shared infrastructure
- `tests/e2e/quick_test_slim_cli_test.go` — Source
- `tests/e2e/test_scripts_per_type_cli_test.go` — Source
- `tests/e2e/features/test-knowledge-convention-driven/` — Source (5 files)

## Acceptance Criteria
- [ ] All 7 test files migrated with correct package name and imports
- [ ] `tests/test-generation/contracts/` contains 3 Contract spec files
- [ ] `tests/test-generation/main_test.go` initializes binary via testkit
- [ ] All tests pass: `go test ./tests/test-generation/... -tags=e2e -count=1`
- [ ] Duplicate test cases removed (check test-knowledge-convention-driven for overlap with root-level tests)

## Hard Rules
- Package name: `testgeneration`
- All test functions use `//go:build e2e` build tag
- The 5 files from `test-knowledge-convention-driven/` sub-package must be flattened into the Journey package

## Implementation Notes
- This is the largest migration (7 files). The `test-knowledge-convention-driven/` sub-package has 5 test files — carefully check for shared helpers within that sub-package
- Some tests in test-knowledge-convention-driven may overlap with `quick_test_slim_cli_test.go` or `test_scripts_per_type_cli_test.go` — deduplicate
- Contract specs should reflect the 3-step generation pipeline: index → generate → run
