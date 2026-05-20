---
id: "1"
title: "Extract testkit + restructure tests/e2e/ Go module"
priority: "P0"
estimated_time: "2h"
dependencies: []
scope: "backend"
breaking: true
type: "coding.refactor"
mainSession: false
---

# 1: Extract testkit + restructure tests/e2e/ Go module

## Description

Foundational task for `tests/e2e/` (plugin E2E tests): move `go.mod` up to `tests/go.mod` (module → `forge-tests`), extract shared infrastructure into `testkit/` sub-package.

**Source files:**
- `tests/e2e/go.mod` → `tests/go.mod`
- `tests/e2e/forge_binary.go` → `tests/testkit/forge_binary.go`
- `tests/e2e/helpers_test.go` → `tests/testkit/helpers.go`
- `tests/e2e/main_test.go` → deleted (each Journey gets its own)

## Reference Files
- `docs/proposals/forge-cli-test-spec-alignment/proposal.md` — Section "1. 提取共享基础设施到 testkit 包"
- `tests/e2e/go.mod` — Current module definition
- `tests/e2e/forge_binary.go` — Binary initialization
- `tests/e2e/helpers_test.go` — Shared helpers
- `forge-cli/tests/e2e/testkit/helpers.go` — Reference: forge-cli already has a testkit

## Acceptance Criteria
- [ ] `tests/go.mod` exists with module name `forge-tests`
- [ ] `tests/testkit/forge_binary.go` exports `ForgeBinary` var and `ForgeCmd()` func
- [ ] `tests/testkit/helpers.go` exports `ParseBlock`, `HasField`, `FieldValue`, `WithRetry` as public functions
- [ ] `tests/testkit/` compiles: `go build ./tests/testkit/...`
- [ ] Existing root-level test files still compile (temporary, migrated in subsequent tasks)

## Hard Rules
- Do NOT modify any test logic, only restructure package boundaries
- Exported function signatures must remain compatible with all current callers
- Follow the same export pattern as `forge-cli/tests/e2e/testkit/helpers.go` for consistency

## Implementation Notes
- `testkit` must be a normal Go package (not `_test` suffix) so other packages can import it
- `forge_binary.go` uses `TestMain` → change to `init()` in testkit
- Helper names need capitalization (unexported → exported) when moved to testkit
- After this task, `tests/e2e/` root files remain temporarily for incremental migration
