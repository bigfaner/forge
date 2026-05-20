---
id: "3"
title: "Migrate feature-management Journey"
priority: "P1"
estimated_time: "2h"
dependencies: ["1"]
scope: "backend"
breaking: false
type: "coding.refactor"
mainSession: false
---

# 3: Migrate feature-management Journey

## Description

Migrate 3 test sources into the `feature-management/` Journey directory. This Journey covers feature CRUD operations, listing, and status lifecycle.

**Source → Target mapping:**
- `tests/e2e/feature_set_command_cli_test.go` → `tests/feature-management/feature_set_test.go`
- `tests/e2e/features/cli-list-reverse-chronological/` → `tests/feature-management/cli_list_reverse_chronological_test.go`
- `tests/e2e/features/proposal-status-lifecycle/` → `tests/feature-management/proposal_status_lifecycle_test.go`
- New: `tests/feature-management/contracts/step-1-feature-set.md`
- New: `tests/feature-management/contracts/step-2-feature-list.md`
- New: `tests/feature-management/contracts/step-3-feature-status.md`
- New: `tests/feature-management/main_test.go`

## Reference Files
- `docs/proposals/forge-cli-test-spec-alignment/proposal.md` — Journey mapping table
- `tests/testkit/` — Shared infrastructure
- `tests/e2e/feature_set_command_cli_test.go` — Source: feature set command tests
- `tests/e2e/features/cli-list-reverse-chronological/` — Source: CLI list ordering tests
- `tests/e2e/features/proposal-status-lifecycle/` — Source: proposal status tests

## Acceptance Criteria
- [ ] 3 test files migrated with correct package name and imports
- [ ] `tests/feature-management/contracts/` contains 3 Contract spec files
- [ ] `tests/feature-management/main_test.go` initializes binary via testkit
- [ ] All tests pass: `go test ./tests/feature-management/... -tags=e2e -count=1`

## Hard Rules
- Package name: `featuremanagement`
- All test functions use `//go:build e2e` build tag
- Sub-package tests (cli-list-reverse-chronological, proposal-status-lifecycle) must be flattened into the Journey package

## Implementation Notes
- Sub-packages may have their own helpers or setup — ensure all are brought into the single Journey package
- Check for and remove any duplicate test cases during migration
