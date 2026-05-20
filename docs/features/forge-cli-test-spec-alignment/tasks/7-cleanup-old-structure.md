---
id: "7"
title: "Cleanup old structure and update justfile"
priority: "P0"
estimated_time: "1h"
dependencies: ["2", "3", "4", "5", "6"]
scope: "backend"
breaking: true
type: "coding.cleanup"
mainSession: false
---

# 7: Cleanup old structure and update justfile

## Description

Final cleanup task: remove the old `tests/e2e/` directory structure, update justfile to point to new Journey paths, and verify the complete E2E test suite runs end-to-end.

**Files to delete:**
- `tests/e2e/main_test.go`
- `tests/e2e/forge_binary.go`
- `tests/e2e/helpers_test.go`
- `tests/e2e/*.cli_test.go` (all root-level test files)
- `tests/e2e/features/` (entire directory)
- `tests/e2e/justfile-canonical-e2e/` (entire directory)
- `tests/e2e/go.mod`, `tests/e2e/go.sum`
- `tests/e2e/results/` (if empty)

**Files to modify:**
- `justfile` — update `test-e2e` and related commands from `./tests/e2e/...` to `./tests/...`
- `.graduated/` markers — update paths if marker files reference old locations

## Reference Files
- `docs/proposals/forge-cli-test-spec-alignment/proposal.md` — Section "4. 更新构建配置"
- `justfile` — E2E test commands
- `tests/e2e/.graduated/` — Graduated test markers

## Acceptance Criteria
- [ ] `tests/e2e/` directory fully removed (no remnant files)
- [ ] `justfile` E2E commands updated to new paths (`./tests/...` instead of `./tests/e2e/...`)
- [ ] `just test-e2e` runs and passes all Journey tests
- [ ] No broken imports or references to old `tests/e2e/` paths
- [ ] `.graduated/` markers relocated to `tests/.graduated/` (if applicable)

## Hard Rules
- Run `just test-e2e` as the final validation step before marking complete
- Do NOT delete `tests/testkit/` or any Journey directories
- Ensure `tests/go.mod` is the only go.mod under `tests/` (no nested modules)

## Implementation Notes
- Verify the justfile doesn't have hardcoded `tests/e2e/` paths in other commands (e.g., `e2e-setup`, `probe`)
- Check `.graduated/` for marker files and update their internal path references
- After deletion, run `go build ./tests/...` to verify no compilation errors across all Journey packages
- The `results/` directory may be a runtime artifact — check if it needs to be created elsewhere or if it's auto-generated
