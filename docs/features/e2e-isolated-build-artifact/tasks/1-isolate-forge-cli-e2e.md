---
id: "1"
title: "Isolate forge-cli/tests/e2e/ to TestMain auto-build"
priority: "P1"
estimated_time: "1-2h"
dependencies: []
scope: "backend"
breaking: false
type: "coding.enhancement"
mainSession: false
---

# 1: Isolate forge-cli/tests/e2e/ to TestMain auto-build

## Description

`forge-cli/tests/e2e/` uses `exec.Command("forge", ...)` which resolves the binary from system PATH (`~/.forge/bin/forge`). This tests the installed version, not the current source. Convert to TestMain auto-build pattern: compile forge binary from source into a temp directory in TestMain, pass the path to tests via a package variable, clean up on exit.

## Reference Files
- `docs/proposals/e2e-isolated-build-artifact/proposal.md` — Source proposal
- `tests/e2e/main_test.go` — Reference implementation of TestMain build pattern

## Acceptance Criteria
- `forge-cli/tests/e2e/` has a `main_test.go` with TestMain that builds forge binary from `../../cmd/forge` (relative to forge-cli module) into a temp directory
- `forgeBinary` package variable holds the temp binary path
- `helpers_test.go` and `testkit/helpers.go` use the `forgeBinary` variable instead of bare `"forge"`
- Temp directory cleaned up in TestMain cleanup
- `grep -r 'exec.Command("forge"' forge-cli/tests/e2e/` returns zero results

## Hard Rules
- Do not modify the test assertions or test logic, only the binary resolution
- Keep testkit/helpers.go's public API unchanged (callers outside this module may exist)

## Implementation Notes
- Follow the pattern from `tests/e2e/main_test.go` but adapt the build path to the forge-cli module structure
- The build command dir should walk up to find the forge source root
- `testkit/helpers.go` needs a way to receive the binary path — consider a package-level var or init function
