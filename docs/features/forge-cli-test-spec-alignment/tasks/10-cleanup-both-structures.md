---
id: "10"
title: "Cleanup both old structures + update justfiles"
priority: "P0"
estimated_time: "2h"
dependencies: ["2", "3", "4", "5", "6", "7", "8", "9"]
scope: "backend"
breaking: true
type: "coding.cleanup"
mainSession: false
---

# 10: Cleanup both old structures + update justfiles

## Description

Final cleanup for both test directories: remove old files, move testkit, update import paths and justfiles.

**tests/e2e/ (plugin tests):**
- Delete `tests/e2e/main_test.go`, `forge_binary.go`, `helpers_test.go`
- Delete all root-level `*.cli_test.go` files
- Delete `tests/e2e/features/` directory
- Delete `tests/e2e/justfile-canonical-e2e/` directory
- Delete `tests/e2e/go.mod`, `tests/e2e/go.sum`
- Delete `tests/e2e/results/` (if empty)
- Move `.graduated/` markers to `tests/.graduated/`

**forge-cli/tests/e2e/ (CLI tests):**
- Delete `forge-cli/tests/e2e/main_test.go`, `helpers_test.go`
- Delete all root-level `*_test.go` files
- Delete `forge-cli/tests/e2e/features/` directory
- Move `forge-cli/tests/e2e/testkit/` → `forge-cli/tests/testkit/`
- Move `forge-cli/tests/e2e/.graduated/` → `forge-cli/tests/.graduated/`
- Update all Journey import paths from `tests/e2e/testkit` to `tests/testkit`

**Justfile updates:**
- Update E2E test commands: `./tests/e2e/...` → `./tests/...`
- Update forge-cli E2E paths: `./forge-cli/tests/e2e/...` → `./forge-cli/tests/...`

## Reference Files
- `docs/proposals/forge-cli-test-spec-alignment/proposal.md` — Section "4. 更新构建配置"
- `justfile` — E2E test commands
- `forge-cli/justfile` or root `justfile` — CLI test commands

## Acceptance Criteria
- [ ] `tests/e2e/` directory fully removed
- [ ] `forge-cli/tests/e2e/` directory fully removed (testkit moved to `forge-cli/tests/testkit/`)
- [ ] All Journey import paths updated to new testkit location
- [ ] Justfile E2E commands updated for both test suites
- [ ] `just test-e2e` runs and passes all Journey tests
- [ ] No broken imports or references to old paths
- [ ] `.graduated/` markers relocated correctly

## Hard Rules
- Run `just test-e2e` as final validation before marking complete
- Do NOT delete `tests/testkit/` or any Journey directories
- Ensure only `tests/go.mod` exists under `tests/` (no nested go.mod)
- Update ALL import paths in forge-cli Journey test files when testkit moves

## Implementation Notes
- Verify justfile doesn't have hardcoded old paths in other commands (e2e-setup, probe, etc.)
- After testkit move, all forge-cli Journey files need import path update: search-replace `tests/e2e/testkit` → `tests/testkit`
- The `results/` directories may be runtime artifacts — check if auto-generated
- Verify `go build ./tests/... ./forge-cli/tests/...` compiles without errors
