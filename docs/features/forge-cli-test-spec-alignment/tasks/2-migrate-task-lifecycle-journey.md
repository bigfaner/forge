---
id: "2"
title: "Migrate task-lifecycle Journey"
priority: "P1"
estimated_time: "2h"
dependencies: ["1"]
scope: "backend"
breaking: false
type: "coding.refactor"
mainSession: false
---

# 2: Migrate task-lifecycle Journey

## Description

Migrate 3 existing test files into the `task-lifecycle/` Journey directory. This Journey covers the full task lifecycle: claim → submit → verify immutability.

**Source → Target mapping:**
- `tests/e2e/task_lifecycle_hardening_cli_test.go` + `tests/e2e/features/fix-task-claim-priority/` → `tests/task-lifecycle/task_lifecycle_test.go`
- `tests/e2e/task_record_immutability_cli_test.go` → `tests/task-lifecycle/task_record_test.go`
- New: `tests/task-lifecycle/contracts/step-1-task-claim.md`
- New: `tests/task-lifecycle/contracts/step-2-task-submit.md`
- New: `tests/task-lifecycle/main_test.go`

## Reference Files
- `docs/proposals/forge-cli-test-spec-alignment/proposal.md` — Journey mapping table
- `tests/testkit/` — Shared infrastructure (from Task 1)
- `tests/e2e/task_lifecycle_hardening_cli_test.go` — Source: task lifecycle tests
- `tests/e2e/task_record_immutability_cli_test.go` — Source: record immutability tests

## Acceptance Criteria
- [ ] `tests/task-lifecycle/task_lifecycle_test.go` contains merged tests from task_lifecycle_hardening + fix-task-claim-priority
- [ ] `tests/task-lifecycle/task_record_test.go` contains migrated task_record_immutability tests
- [ ] `tests/task-lifecycle/contracts/` contains 2 Contract spec files with six-dimension declarations
- [ ] `tests/task-lifecycle/main_test.go` initializes binary via testkit
- [ ] All tests pass: `go test ./tests/task-lifecycle/... -tags=e2e -count=1`
- [ ] No duplicate test cases remain (identical command + scenario tested only once)

## Hard Rules
- Merge strategy: task_lifecycle_hardening + fix-task-claim-priority → single file (related lifecycle stages)
- Package name: `tasklifecycle`
- All test functions use `//go:build e2e` build tag

## Implementation Notes
- The fix-task-claim-priority tests are in a sub-package (`tests/e2e/features/fix-task-claim-priority/`) — check for additional helper functions or setup that needs migration
- When merging, de-duplicate any overlapping test scenarios (same command, same expected outcome)
- Contract specs should be extracted from existing test descriptions/assertions, not written from scratch
