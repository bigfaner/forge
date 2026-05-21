---
id: "5"
title: "Migrate forge-cli/tests/e2e/ Journeys: task-lifecycle + quality-gate"
priority: "P1"
estimated_time: "2h"
dependencies: []
scope: "backend"
breaking: false
type: "coding.refactor"
mainSession: false
---

# 5: Migrate forge-cli/tests/e2e/ Journeys: task-lifecycle + quality-gate

## Description

Migrate forge-cli integration tests for task lifecycle and quality gate into Journey directories. forge-cli testkit already exists at `forge-cli/tests/e2e/testkit/`.

**task-lifecycle Journey:**
- `lifecycle_cli_test.go` (cleanup + quality gate tests) → `forge-cli/tests/task-lifecycle/lifecycle_test.go`
- `submit_cli_test.go` → `forge-cli/tests/task-lifecycle/submit_test.go`
- `features/fix-task-claim-priority/` → merge into `forge-cli/tests/task-lifecycle/fix_task_claim_priority_test.go`
- `features/task-stage-gates/` → merge into `forge-cli/tests/task-lifecycle/task_stage_gates_test.go`
- Contracts: `forge-cli/tests/task-lifecycle/contracts/step-1-task-claim.md`, `step-2-task-submit.md`, `step-3-quality-gate.md`

**quality-gate Journey** (split from lifecycle_cli_test.go if quality gate tests are separable):
- If quality gate tests are tightly coupled with lifecycle, keep in task-lifecycle Journey and skip separate quality-gate Journey.

## Reference Files
- `docs/proposals/forge-cli-test-spec-alignment/proposal.md` — Reference pattern
- `forge-cli/tests/e2e/testkit/` — Already existing testkit (use as-is)
- `forge-cli/tests/e2e/lifecycle_cli_test.go` — Source
- `forge-cli/tests/e2e/submit_cli_test.go` — Source
- `forge-cli/tests/e2e/features/fix-task-claim-priority/` — Source
- `forge-cli/tests/e2e/features/task-stage-gates/` — Source

## Acceptance Criteria
- [ ] `forge-cli/tests/task-lifecycle/` created with migrated test files
- [ ] `contracts/` contains spec files with six-dimension declarations
- [ ] `main_test.go` calls `testkit.SetForgeBinary()` and builds binary
- [ ] Tests pass: `go test ./forge-cli/tests/task-lifecycle/... -tags=e2e -count=1`
- [ ] Sub-package tests (fix-task-claim-priority, task-stage-gates) flattened

## Hard Rules
- Package name: `tasklifecycle`
- All test functions use `//go:build e2e` build tag
- Import testkit from `forge-cli/tests/e2e/testkit` initially (path will be updated in cleanup task)
- Do NOT modify test logic, only restructure

## Implementation Notes
- forge-cli tests are part of `forge-cli` Go module (no separate go.mod)
- `testkit` is at `forge-cli/tests/e2e/testkit/` — import path is `github.com/user/forge-cli/tests/e2e/testkit`
- When moving to `forge-cli/tests/task-lifecycle/`, the import path changes to `github.com/user/forge-cli/tests/e2e/testkit` (testkit stays at e2e level during migration, moves up in cleanup)
- lifecycle_cli_test.go has 5 tests (TC_012-016) covering cleanup and quality gate
