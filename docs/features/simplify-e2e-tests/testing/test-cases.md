---
feature: "simplify-e2e-tests"
sources:
  - docs/proposals/simplify-e2e-tests/proposal.md
generated: "2026-05-16"
---

# Test Cases: simplify-e2e-tests

## Summary

| Type | Count |
|------|-------|
| UI   | 0   |
| **Integration** | **0** |
| API  | 0  |
| CLI  | 4  |
| **Total** | **4** |

---

## CLI Test Cases

## TC-001: Verify tui-ui-design directory deleted
- **Source**: Proposal Success Criterion 1 — "`tests/e2e/tui-ui-design/` directory deleted"
- **Type**: CLI
- **Target**: cli/e2e-cleanup
- **Test ID**: cli/e2e-cleanup/verify-tui-ui-design-directory-deleted
- **Pre-conditions**: Task 1 (Remove text-verification e2e tests) has been executed
- **Steps**:
  1. Check that the directory `tests/e2e/tui-ui-design/` does not exist on the filesystem
- **Expected**: Directory `tests/e2e/tui-ui-design/` is absent; `os.Stat` returns an error with `os.IsNotExist` evaluating to true
- **Priority**: P0

## TC-002: Verify TC-020 removed from justfile-canonical-e2e
- **Source**: Proposal Success Criterion 2 — "`TestTC_020` removed from `justfile-canonical-e2e/justfile_canonical_e2e_cli_test.go`"
- **Type**: CLI
- **Target**: cli/e2e-cleanup
- **Test ID**: cli/e2e-cleanup/verify-tc020-removed-from-justfile-canonical-e2e
- **Pre-conditions**: Task 1 (Remove text-verification e2e tests) has been executed
- **Steps**:
  1. Read the file `tests/e2e/justfile-canonical-e2e/justfile_canonical_e2e_cli_test.go`
  2. Search for the function name `TestTC_020_AllManifestsContainZeroRunAndGraduateFields`
- **Expected**: The function `TestTC_020_AllManifestsContainZeroRunAndGraduateFields` is not present in the file. All other test functions remain intact.
- **Priority**: P0

## TC-003: Verify e2e test suite compiles
- **Source**: Proposal Success Criterion 3 — "`go test -tags=e2e ./tests/e2e/...` compiles without errors"
- **Type**: CLI
- **Target**: cli/e2e-compilation
- **Test ID**: cli/e2e-compilation/verify-e2e-test-suite-compiles
- **Pre-conditions**: Task 1 (Remove text-verification e2e tests) has been executed; Go toolchain is available
- **Steps**:
  1. Run `go test -tags=e2e ./tests/e2e/... -count=1 -run=^$` (compile-only, no test execution)
  2. Check exit code is 0 and no compilation errors appear in output
- **Expected**: Command exits with code 0. No compilation errors. All remaining test files are valid Go source.
- **Priority**: P0

## TC-004: Verify remaining CLI behavior tests pass
- **Source**: Proposal Success Criterion 4 — "Remaining CLI-focused tests pass"
- **Type**: CLI
- **Target**: cli/e2e-execution
- **Test ID**: cli/e2e-execution/verify-remaining-cli-behavior-tests-pass
- **Pre-conditions**: Task 1 completed; e2e test suite compiles; forge binary is built
- **Steps**:
  1. Run `go test -tags=e2e ./tests/e2e/... -count=1 -timeout 120s`
  2. Check exit code is 0 and all tests pass
- **Expected**: All remaining CLI behavior tests pass with exit code 0. No test failures related to the deleted text-verification tests.
- **Priority**: P0

---

## Traceability

| TC ID | Source | Type | Target | Priority |
|-------|--------|------|--------|----------|
| TC-001 | Proposal Success Criterion 1 | CLI | cli/e2e-cleanup | P0 |
| TC-002 | Proposal Success Criterion 2 | CLI | cli/e2e-cleanup | P0 |
| TC-003 | Proposal Success Criterion 3 | CLI | cli/e2e-compilation | P0 |
| TC-004 | Proposal Success Criterion 4 | CLI | cli/e2e-execution | P0 |
