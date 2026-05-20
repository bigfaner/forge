---
id: "2"
title: "Isolate justfile-canonical-e2e/ to TestMain auto-build"
priority: "P1"
estimated_time: "1h"
dependencies: []
scope: "backend"
breaking: false
type: "coding.enhancement"
mainSession: false
---

# 2: Isolate justfile-canonical-e2e/ to TestMain auto-build

## Description

`tests/e2e/justfile-canonical-e2e/helpers_test.go` uses a hardcoded relative path `../../../forge-cli/bin/forge` with a lazy build-if-missing fallback. This writes to a shared location that can become stale across branches. Convert to TestMain auto-build: compile into a temp directory, clean up after tests.

## Reference Files
- `docs/proposals/e2e-isolated-build-artifact/proposal.md` — Source proposal
- `tests/e2e/main_test.go` — Reference implementation of TestMain build pattern

## Acceptance Criteria
- `tests/e2e/justfile-canonical-e2e/main_test.go` added with TestMain that builds forge into temp dir
- `forgeBinaryPath` package variable set to the temp binary path
- `helpers_test.go` uses `forgeBinaryPath` instead of `../../../forge-cli/bin/forge`
- No reference to `forge-cli/bin/forge` in this module
- Temp directory cleaned up on exit

## Hard Rules
- This is a separate Go module (`tests/e2e/justfile-canonical-e2e/go.mod`), build path must resolve independently

## Implementation Notes
- The build command dir should walk up from the test module to find `forge-cli/cmd/forge`
- Follow the pattern from `tests/e2e/main_test.go`
