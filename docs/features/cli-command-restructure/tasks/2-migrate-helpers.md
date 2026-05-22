---
id: "2"
title: "Migrate test_results.go and journey_isolation.go to pkg/testrunner"
priority: "P0"
estimated_time: "1h"
dependencies: ["1"]
scope: "backend"
breaking: false
type: "coding.refactor"
mainSession: false
---

# 2: Migrate test_results.go and journey_isolation.go to pkg/testrunner

## Description

Move `test_results.go` and `journey_isolation.go` (plus its test) from `internal/cmd/` to `pkg/testrunner/`. These files are helper utilities closely related to testrunner functionality and belong in the shared package.

## Reference Files
- `docs/proposals/cli-command-restructure/proposal.md` — Source proposal

## Acceptance Criteria

- `test_results.go` is moved to `pkg/testrunner/` with updated package declaration
- `journey_isolation.go` and `journey_isolation_test.go` are moved to `pkg/testrunner/`
- All imports referencing the old location are updated
- `go build ./...` passes
- `go test ./...` passes

## Hard Rules

- Update package declaration from `cmd` to `testrunner` in moved files
- All exported symbols remain exported (no visibility changes)
- Update all import paths in files that reference the moved code

## Implementation Notes

Files to move:
- `forge-cli/internal/cmd/test_results.go` → `forge-cli/pkg/testrunner/test_results.go`
- `forge-cli/internal/cmd/journey_isolation.go` → `forge-cli/pkg/testrunner/journey_isolation.go`
- `forge-cli/internal/cmd/journey_isolation_test.go` → `forge-cli/pkg/testrunner/journey_isolation_test.go`

After moving:
1. Change `package cmd` to `package testrunner` in each file
2. Update internal imports: remove old `cmd` package references, add `testrunner` import
3. Grep for all references to the moved symbols and update import paths
4. Files likely affected: test_promote.go, test_verify.go, and other cmd files that use these helpers
