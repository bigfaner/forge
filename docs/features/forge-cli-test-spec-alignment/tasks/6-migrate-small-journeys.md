---
id: "6"
title: "Migrate small Journeys (quality-gate, task-type-system, e2e-pipeline, command-regression)"
priority: "P1"
estimated_time: "2h"
dependencies: ["1"]
scope: "backend"
breaking: false
type: "coding.refactor"
mainSession: false
---

# 6: Migrate small Journeys (quality-gate, task-type-system, e2e-pipeline, command-regression)

## Description

Migrate 4 single-file Journeys. Each has exactly 1 source test file and 1 Contract spec to create.

**Source → Target mapping:**

| Journey | Source | Target |
|---------|--------|--------|
| `quality-gate/` | `tests/e2e/quality_gate_fix_task_loop_breaker_cli_test.go` | `tests/quality-gate/quality_gate_test.go` |
| `task-type-system/` | `tests/e2e/features/task-type-refinement/` | `tests/task-type-system/task_type_refinement_test.go` |
| `e2e-pipeline/` | `tests/e2e/justfile-canonical-e2e/` | `tests/e2e-pipeline/justfile_canonical_test.go` |
| `command-regression/` | `tests/e2e/features/test-knowledge-convention-driven/removed_commands_cli_test.go` | `tests/command-regression/removed_commands_test.go` |

Each Journey also needs:
- `contracts/` directory with 1+ Contract spec files
- `main_test.go` initializing binary via testkit

## Reference Files
- `docs/proposals/forge-cli-test-spec-alignment/proposal.md` — Journey mapping table
- `tests/testkit/` — Shared infrastructure
- `tests/e2e/quality_gate_fix_task_loop_breaker_cli_test.go` — Source
- `tests/e2e/features/task-type-refinement/` — Source
- `tests/e2e/justfile-canonical-e2e/` — Source
- `tests/e2e/features/test-knowledge-convention-driven/removed_commands_cli_test.go` — Source

## Acceptance Criteria
- [ ] 4 Journey directories created, each with migrated test file(s)
- [ ] Each Journey has `contracts/` with appropriate Contract spec(s)
- [ ] Each Journey has `main_test.go` via testkit
- [ ] All 4 Journeys pass: `go test ./tests/quality-gate/... ./tests/task-type-system/... ./tests/e2e-pipeline/... ./tests/command-regression/... -tags=e2e -count=1`
- [ ] `e2e-pipeline/` Journey's separate `go.mod` is eliminated (consolidated into `tests/go.mod`)

## Hard Rules
- Package names: `qualitygate`, `tasktypesystem`, `e2epipeline`, `commandregression`
- All test functions use `//go:build e2e` build tag
- `justfile-canonical-e2e/` has its own `go.mod` — must be consolidated into `tests/go.mod`
- Sub-package tests (task-type-refinement) must be flattened into their Journey package

## Implementation Notes
- `justfile-canonical-e2e/` is a separate Go module with its own `go.mod` — the most complex migration in this group. Verify all dependencies are available in `tests/go.mod`
- `quality_gate_fix_task_loop_breaker_cli_test.go` (16KB) tests the loop-breaker mechanism — ensure the test still correctly references the fix-task workflow
- Contract specs for each Journey: quality-gate (1 contract), task-type-system (3 contracts: list-types, validate-index, migrate), e2e-pipeline (2 contracts: setup, run), command-regression (1 contract)
