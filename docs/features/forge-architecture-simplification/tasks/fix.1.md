---
id: "fix.1"
title: "Fix TestRunStatus silent failure in internal/cmd"
priority: "P0"
dependencies: []
status: 
type: "coding.fix"
breaking: true
---

# fix.1: Fix TestRunStatus silent failure in internal/cmd

## Problem
`go test ./internal/cmd/` fails silently — `TestRunStatus` runs but produces no error output and exits with code 1. Likely caused by `os.Exit()` being called during `rootCmd.Execute()` for `task status 1.1 blocked` (2-arg mutation mode).

## Scope
- `forge-cli/internal/cmd/feature_test.go` — TestRunStatus uses 2-arg status mutation
- `forge-cli/internal/cmd/status.go` — may call Exit/OS.Exit in error path
- `forge-cli/internal/cmd/errors.go` — Exit() function added by task 2.2

## Acceptance Criteria
- [ ] `go test ./internal/cmd/ -run TestRunStatus` passes with clear output
- [ ] `go test ./internal/cmd/...` all pass (0 failures)
- [ ] No other test regressions across `go test ./...`
