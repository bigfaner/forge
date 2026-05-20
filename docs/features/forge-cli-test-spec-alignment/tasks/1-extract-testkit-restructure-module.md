---
id: "1"
title: "Extract testkit package and restructure Go module"
priority: "P0"
estimated_time: "2h"
dependencies: []
scope: "backend"
breaking: true
type: "coding.refactor"
mainSession: false
---

# 1: Extract testkit package and restructure Go module

## Description

Foundational task: move `tests/e2e/go.mod` up to `tests/go.mod` (module name → `forge-tests`), then extract shared infrastructure into a `testkit` sub-package. This enables all subsequent Journey migrations to import from a single shared package.

**Source files:**
- `tests/e2e/go.mod` → `tests/go.mod`
- `tests/e2e/forge_binary.go` → `tests/testkit/forge_binary.go`
- `tests/e2e/helpers_test.go` → `tests/testkit/helpers.go`
- `tests/e2e/main_test.go` → deleted (each Journey will have its own)

## Reference Files
- `docs/proposals/forge-cli-test-spec-alignment/proposal.md` — Source proposal, Section "1. 提取共享基础设施到 testkit 包"
- `tests/e2e/go.mod` — Current module definition
- `tests/e2e/forge_binary.go` — Current binary initialization
- `tests/e2e/helpers_test.go` — Current shared helpers

## Acceptance Criteria
- [ ] `tests/go.mod` exists with module name `forge-tests`
- [ ] `tests/testkit/forge_binary.go` exports `ForgeBinary` var and `ForgeCmd()` func (public API)
- [ ] `tests/testkit/helpers.go` exports `ParseBlock`, `HasField`, `FieldValue`, `WithRetry` as public functions
- [ ] `tests/testkit/` compiles without errors (`go build ./...`)
- [ ] Existing root-level test files still compile after testkit creation (temporary — they'll be migrated in subsequent tasks)

## Hard Rules
- Do NOT modify any test logic, only restructure package boundaries
- Exported function signatures in testkit must remain compatible with all current callers
- `tests/e2e/main_test.go` is NOT migrated — it will be deleted in cleanup task

## Implementation Notes
- The `testkit` package must be a normal Go package (not `_test` suffix) so other packages can import it
- `forge_binary.go` currently uses `TestMain` to build the binary; in testkit it should use `init()` instead
- Helper functions currently use unexported names; they need to be exported (capitalized) when moved to testkit
- After this task, `tests/e2e/` root files still exist temporarily for incremental migration
