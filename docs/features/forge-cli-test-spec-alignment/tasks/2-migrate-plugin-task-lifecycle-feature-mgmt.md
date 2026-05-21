---
id: "2"
title: "Migrate tests/e2e/ Journeys: task-lifecycle + feature-management"
priority: "P1"
estimated_time: "2h"
dependencies: ["1"]
scope: "backend"
breaking: false
type: "coding.refactor"
mainSession: false
---

# 2: Migrate tests/e2e/ Journeys: task-lifecycle + feature-management

## Description

Migrate 8 test files from `tests/e2e/` into 2 Journey directories.

**task-lifecycle Journey:**
- `tests/e2e/task_lifecycle_hardening_cli_test.go` + `tests/e2e/features/fix-task-claim-priority/` → `tests/task-lifecycle/task_lifecycle_test.go`
- `tests/e2e/task_record_immutability_cli_test.go` → `tests/task-lifecycle/task_record_test.go`
- Contracts: `tests/task-lifecycle/contracts/step-1-task-claim.md`, `step-2-task-submit.md`

**feature-management Journey:**
- `tests/e2e/feature_set_command_cli_test.go` → `tests/feature-management/feature_set_test.go`
- `tests/e2e/features/cli-list-reverse-chronological/` → `tests/feature-management/cli_list_reverse_chronological_test.go`
- `tests/e2e/features/proposal-status-lifecycle/` → `tests/feature-management/proposal_status_lifecycle_test.go`
- Contracts: `tests/feature-management/contracts/step-1-feature-set.md`, `step-2-feature-list.md`, `step-3-feature-status.md`

## Reference Files
- `docs/proposals/forge-cli-test-spec-alignment/proposal.md` — Journey mapping table
- `tests/testkit/` — Shared infrastructure (Task 1)
- All source files listed above

## Acceptance Criteria
- [ ] Both Journey directories created with correct package names
- [ ] All 8 test files migrated (5 from task-lifecycle, 3 from feature-management)
- [ ] Each Journey has `main_test.go` calling testkit, `contracts/` with spec files
- [ ] Tests pass: `go test ./tests/task-lifecycle/... ./tests/feature-management/... -tags=e2e -count=1`
- [ ] Duplicate test cases removed (fix-task-claim-priority merged into task_lifecycle)

## Hard Rules
- Package names: `tasklifecycle`, `featuremanagement`
- All test functions use `//go:build e2e` build tag
- Sub-package tests must be flattened into their Journey package
- Contract specs: six-dimension declarations extracted from existing test descriptions

## Implementation Notes
- fix-task-claim-priority is a sub-package with its own main_test.go — merge into task-lifecycle
- Deduplicate overlapping test scenarios during merge
