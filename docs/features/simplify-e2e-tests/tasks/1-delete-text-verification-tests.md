---
id: "1"
title: "Remove text-verification e2e tests"
priority: "P1"
estimated_time: "30m"
dependencies: []
scope: "backend"
breaking: false
type: "implementation"
mainSession: false
---

# 1: Remove text-verification e2e tests

## Description

Remove 32 text-verification test cases across 2 files that verify static text content of plugin/template files rather than forge-cli runtime behavior. These tests break whenever plugin documentation is reworded and add maintenance burden without testing CLI functionality.

**Action 1 — Delete**: `tests/e2e/tui-ui-design/` directory (31 text-verification tests in `tui_ui_design_cli_test.go` reading SKILL.md, template .md files, and asserting string presence).

**Action 2 — Edit**: Remove `TestTC_020_AllManifestsContainZeroRunAndGraduateFields` from `tests/e2e/justfile-canonical-e2e/justfile_canonical_e2e_cli_test.go` (reads 6 static `manifest.yaml` files and asserts YAML key presence/absence — no CLI invocation). Preserve the remaining 19 CLI behavior tests in this file.

## Reference Files
- `docs/proposals/simplify-e2e-tests/proposal.md` — Source proposal
- `tests/e2e/tui-ui-design/tui_ui_design_cli_test.go` — Target for deletion
- `tests/e2e/justfile-canonical-e2e/justfile_canonical_e2e_cli_test.go` — Remove TC-020 only

## Acceptance Criteria
- [ ] `tests/e2e/tui-ui-design/` directory no longer exists
- [ ] `TestTC_020` removed from `justfile_canonical_e2e_cli_test.go`, all other tests preserved
- [ ] `go test -tags=e2e ./tests/e2e/...` compiles without errors
- [ ] Remaining CLI behavior tests still pass

## Hard Rules
- Do NOT modify any other test functions in `justfile_canonical_e2e_cli_test.go` — only remove TC-020
- Do NOT modify `cli_lean_output_cli_test.go`, `test_scripts_per_type_cli_test.go`, or helper files

## Implementation Notes
- `tui-ui-design/` package is self-contained with local helpers (`readFile`, `fileContains`, `truncate`) — no shared cleanup needed
- TC-020 in justfile-canonical-e2e reads static YAML manifests (`forge-cli/pkg/profile/profiles/*/manifest.yaml`) — removing it loses no CLI behavior coverage
- Quality gate (compile + test) will verify no broken imports or missing dependencies
