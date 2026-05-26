---
id: "4"
title: "Clean tests/ invalid tests and tui-ui-design directory"
priority: "P0"
estimated_time: "1h"
dependencies: []
surface-key: "."
surface-type: "cli"
breaking: false
type: "coding.cleanup"
mainSession: false
---

# 4: Clean tests/ invalid tests and tui-ui-design directory

## Description
Delete five categories of invalid tests in `tests/`:

1. **Recursive `go test` tests** (2): `simplify_e2e_tests` TC-003/TC-004 spawn `exec.Command("go", "test")` recursively. On Windows, this creates 126+ orphan processes consuming 6GB+ RAM — machine-freezing bug.
2. **Permanent skip tests** (2): `feature_set_command` TC-016/TC-017 with `t.Skip("requires real git worktree")` — never execute.
3. **Hollow assertion tests** (8): `cli_lean_output` TC-006~011, TC-013, TC-018 use `if cond { assert }` where `cond` is always false — assertions never fire.
4. **Entire `cli_lean_output_cli_test.go`** (19 tests): 8 hollow + 11 conditional skips with no fixture — entire file is dead.
5. **Entire `tui-ui-design/` directory** (31 text-verification tests): all read static files and grep for text.

## Reference Files
- `docs/proposals/clean-invalid-tests/proposal.md#Evidence` — table rows for recursive tests, permanent skips, hollow assertions
- `docs/proposals/clean-invalid-tests/proposal.md#Scope` — In Scope items 9-13
- `docs/proposals/clean-invalid-tests/proposal.md#Urgency` — "递归测试在 Windows 上可导致机器卡死" — P0 safety issue
- `docs/proposals/clean-invalid-tests/proposal.md#Assumptions-Challenged` — hollow assertions never fire, removing them loses nothing

## Acceptance Criteria
- [ ] Zero `exec.Command("go", "test"` calls in `tests/` (verify: `grep -rn '"go",.*"test"' tests/ --include='*_test.go'`)
- [ ] `feature_set_command` skip tests deleted
- [ ] `cli_lean_output_cli_test.go` deleted entirely
- [ ] `tests/.graduated/tui-ui-design/` directory deleted entirely (if exists in root tests/, also delete there)
- [ ] `tests/test-suite-health/simplify_e2e_tests_test.go` — TC-003/TC-004 test functions deleted (or entire file if all tests are recursive)
- [ ] Empty files and directories cleaned up
- [ ] `go build -tags=e2e ./tests/...` compiles successfully

## Hard Rules
- The recursive `go test` tests are DANGEROUS — delete them first before any other work
- Do NOT modify `tests/test-suite-health/` meta-test assertions — only delete the problematic test cases within them
- Verify `tui-ui-design/` exists before deleting (may already be only in `.graduated/`)

## Implementation Notes
- `simplify_e2e_tests_test.go` is in `tests/test-suite-health/` — delete only TC-003 and TC-004 test functions, not the entire file (other TCs may be valid)
- For `feature_set_command`: likely `tests/.graduated/feature-set-command/` or in `tests/feature-management/feature_set_test.go` — check both locations
- For `cli_lean_output`: check both `tests/` root and `.graduated/` directories
- After all deletions, run `go build -tags=e2e ./tests/...` to verify compilation
