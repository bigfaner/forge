---
id: "4"
title: "Migrate tests/e2e/ remaining Journeys: quality-gate, task-type-system, e2e-pipeline, command-regression, test-suite-health"
priority: "P1"
estimated_time: "2h"
dependencies: ["1"]
scope: "backend"
breaking: false
type: "coding.refactor"
mainSession: false
---

# 4: Migrate tests/e2e/ remaining Journeys (5 Journeys)

## Description

Migrate remaining test files into 5 smaller Journey directories.

| Journey | Source | Target |
|---------|--------|--------|
| `quality-gate/` | `quality_gate_fix_task_loop_breaker_cli_test.go` | `quality_gate_test.go` |
| `task-type-system/` | `features/task-type-refinement/` | `task_type_refinement_test.go` |
| `e2e-pipeline/` | `justfile-canonical-e2e/` | `justfile_canonical_test.go` |
| `command-regression/` | `features/test-knowledge-convention-driven/removed_commands_cli_test.go` | `removed_commands_test.go` |
| `test-suite-health/` | `e2e_test_quality_cleanup_cli_test.go`, `gen_journeys_skill_cli_test.go`, `simplify_e2e_tests_cli_test.go` | 3 test files |

Each Journey needs: `main_test.go` (via testkit), `contracts/` with spec files.

## Reference Files
- `docs/proposals/forge-cli-test-spec-alignment/proposal.md` — Journey mapping table
- `tests/testkit/` — Shared infrastructure (Task 1)

## Acceptance Criteria
- [ ] 5 Journey directories created with correct package names
- [ ] All test files migrated (total ~8 files across 5 Journeys)
- [ ] Each Journey has `contracts/` with appropriate spec files
- [ ] All tests pass: `go test ./tests/quality-gate/... ./tests/task-type-system/... ./tests/e2e-pipeline/... ./tests/command-regression/... ./tests/test-suite-health/... -tags=e2e -count=1`
- [ ] `justfile-canonical-e2e/` separate go.mod eliminated (consolidated into `tests/go.mod`)

## Hard Rules
- Package names: `qualitygate`, `tasktypesystem`, `e2epipeline`, `commandregression`, `testsuitehealth`
- All test functions use `//go:build e2e` build tag
- Sub-package tests must be flattened into Journey packages
- Contract counts: quality-gate (1), task-type-system (3), e2e-pipeline (2), command-regression (1), test-suite-health (1)

## Implementation Notes
- `justfile-canonical-e2e/` has its own go.mod — must be consolidated into `tests/go.mod`
- `quality_cleanup_test.go` is a meta-test — ensure it still works after directory reorganization
