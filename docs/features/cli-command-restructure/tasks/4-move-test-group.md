---
id: "4"
title: "Move forge test group to test/ subdirectory"
priority: "P1"
estimated_time: "1h"
dependencies: ["2"]
scope: "backend"
breaking: false
type: "coding.refactor"
mainSession: false
---

# 4: Move forge test group to test/ subdirectory

## Description

Move all forge test subcommand files from `internal/cmd/` to `internal/cmd/test/` subdirectory. Create a `Register()` function for the test command group.

## Reference Files
- `docs/proposals/cli-command-restructure/proposal.md` — Source proposal

## Acceptance Criteria

- All test subcommand files are in `internal/cmd/test/`
- New package exports a `Register()` function
- root.go updated to use the new package
- `go build ./...` passes
- `go test ./...` passes
- `forge test promote/run-journey/verify` behavior unchanged

## Hard Rules

- The test sub-package must NOT import `internal/cmd` (no circular deps)
- All test-related tests move to the subdirectory

## Implementation Notes

Files to move to `internal/cmd/test/`:
- test.go, test_test.go
- test_promote.go, test_promote_test.go
- test_verify.go, test_verify_test.go

Note: `test_results.go` was already moved to `pkg/testrunner/` in task 2. Verify no references to it remain in cmd package.

The test package will import from `pkg/testrunner/` for results handling.
