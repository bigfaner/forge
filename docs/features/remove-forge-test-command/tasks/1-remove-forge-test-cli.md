---
id: "1"
title: "Delete forge test CLI commands and Go dependencies"
priority: "P0"
estimated_time: "1.5h"
dependencies: []
surface-key: ""
surface-type: ""
breaking: false
type: "coding.cleanup"
mainSession: false
---

# 1: Delete forge test CLI commands and Go dependencies

## Description
Remove the entire `forge test` command group (promote, run-journey, verify) and all Go package code that exclusively serves it. This includes the command implementation directory, the contract package, journey isolation functions in testrunner, integration tests, and the e2e regression test suite. The testrunner shared functions (`RunProjectTests`, `WriteUnitTestRawOutput`, `WriteRegressionRawOutput`, `Capitalize`, `PrintHookJSON`) used by `quality-gate` and `feature complete` must be preserved.

## Reference Files
- `proposal.md#Proposed-Solution` — defines deletion scope and testrunner preservation boundary
- `proposal.md#Constraints-&-Dependencies` — lists 5 functions that must be preserved and packages safe to delete
- `proposal.md#Scope` — exact In Scope items this task implements
- `proposal.md#Success-Criteria` — compilation and zero-residue verification

## Acceptance Criteria
- [ ] `forge test` execution returns "unknown command" error
- [ ] `go build ./...` compiles with zero errors
- [ ] `go test ./...` all pass
- [ ] `forge quality-gate` and `forge feature complete` commands return exit code 0

## Hard Rules
- **DO NOT delete** any file in `forge-cli/pkg/testrunner/` except `journey_isolation.go` and `journey_isolation_test.go`. The shared functions in `testrunner.go` and `test_results.go` are used by `quality-gate` and `feature complete`.
- Delete directories/files in this exact order to avoid partial states: first delete `cmd/test/` and `pkg/contract/`, then clean `root.go`/`root_test.go`/`quality_gate.go`, then delete `journey_isolation*.go`, then delete `tests/command-regression/`.
- Run `go build ./...` after all deletions to catch broken imports before proceeding.

## Implementation Notes
- Files to delete:
  - `forge-cli/internal/cmd/test/` (entire directory: 7 files)
  - `forge-cli/pkg/contract/` (entire directory: 7 files)
  - `forge-cli/internal/cmd/test_test.go`
  - `forge-cli/internal/cmd/test_verify_test.go`
  - `forge-cli/pkg/testrunner/journey_isolation.go`
  - `forge-cli/pkg/testrunner/journey_isolation_test.go`
  - `tests/command-regression/` (entire directory: 3 files)
- Files to modify:
  - `forge-cli/internal/cmd/root.go` — remove `testpkg` import and `testpkg.Register()` / `rootCmd.AddCommand(testpkg.Cmd)` calls
  - `forge-cli/internal/cmd/root_test.go` — remove `testpkg` import (line 10) and `TestRootCmd_TestGroupHasSubcommands` function (references `testpkg.Cmd`)
  - `forge-cli/internal/cmd/quality_gate.go` — update `forge test promote` hint text on line 152
