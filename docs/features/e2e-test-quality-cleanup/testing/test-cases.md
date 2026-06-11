---
feature: "e2e-test-quality-cleanup"
sources:
  - docs/proposals/e2e-test-quality-cleanup/proposal.md
generated: "2026-05-16"
---

# Test Cases: e2e-test-quality-cleanup

## Summary

| Type | Count |
|------|-------|
| UI   | 0   |
| **Integration** | **0** |
| API  | 0  |
| CLI  | 7  |
| **Total** | **7** |

---

## CLI Test Cases

## TC-001: Deleted test files do not exist
- **Source**: Proposal In Scope #1, #2, #3, #4
- **Type**: CLI
- **Target**: cli/test-cleanup
- **Test ID**: cli/test-cleanup/deleted-test-files-do-not-exist
- **Pre-conditions**: Cleanup tasks 1-2 completed
- **Steps**:
  1. Check that `tests/e2e/extract_design_md_platform_adapters_cli_test.go` does not exist
  2. Check that `tests/e2e/cli_list_reverse_chronological_cli_test.go` (root copy) does not exist
  3. Check that `tests/e2e/fix_task_claim_priority_cli_test.go` (root copy) does not exist
  4. Check that `tests/e2e/cli_lean_output_cli_test.go` does not exist
- **Expected**: All four files are absent from the filesystem
- **Priority**: P0

## TC-002: Deleted test functions do not exist
- **Source**: Proposal In Scope #5, #6, #7
- **Type**: CLI
- **Target**: cli/test-cleanup
- **Test ID**: cli/test-cleanup/deleted-test-functions-do-not-exist
- **Pre-conditions**: Cleanup task 2 completed
- **Steps**:
  1. Read `tests/e2e/simplify_e2e_tests_cli_test.go` and verify `TestTC_003_*` and `TestTC_004_*` function names are absent
  2. Read `tests/e2e/feature_set_command_cli_test.go` and verify `TestTC_016_*` and `TestTC_017_*` function names are absent
  3. Read `tests/e2e/quick_test_slim_cli_test.go` and verify `TestTC_003_*`, `TestTC_009_*`, `TestTC_010_*`, `TestTC_013_*`, `TestTC_016_*` function names are absent
- **Expected**: All specified test function names are absent from their respective files
- **Priority**: P0

## TC-003: E2E test suite compiles successfully
- **Source**: Proposal Success Criterion "just test-e2e compiles and all pass" + Key Risk "compilation failure"
- **Type**: CLI
- **Target**: cli/e2e-compile
- **Test ID**: cli/e2e-compile/e2e-test-suite-compiles-successfully
- **Pre-conditions**: All cleanup tasks completed
- **Steps**:
  1. Run `just e2e-compile` (or equivalent `go test -tags=e2e -c ./tests/e2e/...`)
  2. Verify exit code is 0
- **Expected**: Compilation succeeds with no errors
- **Priority**: P0

## TC-004: Zero unconditional t.Skip in test files
- **Source**: Proposal Success Criterion "zero t.Skip unconditional skips"
- **Type**: CLI
- **Target**: cli/antipattern-scan
- **Test ID**: cli/antipattern-scan/zero-unconditional-t-skip
- **Pre-conditions**: All cleanup tasks completed
- **Steps**:
  1. Grep all `_test.go` files in `tests/e2e/` for `t.Skip(` calls
  2. Exclude any `t.Skip(` calls that are inside conditional branches (if/else)
  3. Verify no unconditional `t.Skip(` remains
- **Expected**: No unconditional `t.Skip(` calls exist in any e2e test file
- **Priority**: P0

## TC-005: Zero recursive go test invocations
- **Source**: Proposal Success Criterion "zero recursive exec.Command go test calls"
- **Type**: CLI
- **Target**: cli/antipattern-scan
- **Test ID**: cli/antipattern-scan/zero-recursive-go-test-invocations
- **Pre-conditions**: All cleanup tasks completed
- **Steps**:
  1. Grep all `_test.go` files in `tests/e2e/` for `exec.Command("go", "test"`
  2. Verify no matches found
- **Expected**: No `exec.Command("go", "test"` calls exist in any e2e test file
- **Priority**: P0

## TC-006: No static file text-grep tests remain
- **Source**: Proposal Success Criterion "zero static source file text check tests"
- **Type**: CLI
- **Target**: cli/antipattern-scan
- **Test ID**: cli/antipattern-scan/no-static-file-text-grep-tests
- **Pre-conditions**: All cleanup tasks completed
- **Steps**:
  1. Grep all `_test.go` files in `tests/e2e/` for patterns that read `.md` or `.go` source files and assert on text content
  2. Specifically check for `os.ReadFile` followed by `assert.Contains` on static source files (not test fixtures)
  3. Verify no such patterns exist
- **Expected**: No tests that read static source files and check for specific text strings
- **Priority**: P1

## TC-007: No duplicate test files between root and features directory
- **Source**: Proposal Success Criterion "no duplicate files between tests/e2e/ and tests/e2e/features/"
- **Type**: CLI
- **Target**: cli/antipattern-scan
- **Test ID**: cli/antipattern-scan/no-duplicate-test-files-root-and-features
- **Pre-conditions**: All cleanup tasks completed
- **Steps**:
  1. List all `_test.go` files in `tests/e2e/` (root level)
  2. List all `_test.go` files in `tests/e2e/features/*/`
  3. Compare filenames: no file in root should have an identical file (by content hash) in any features/ subdirectory
- **Expected**: No duplicate test files exist between root and features/ directories
- **Priority**: P1

---

## Traceability

| TC ID | Source | Type | Target | Priority |
|-------|--------|------|--------|----------|
| TC-001 | Proposal In Scope #1,#2,#3,#4 | CLI | cli/test-cleanup | P0 |
| TC-002 | Proposal In Scope #5,#6,#7 | CLI | cli/test-cleanup | P0 |
| TC-003 | Proposal Success Criterion #1 + Key Risk | CLI | cli/e2e-compile | P0 |
| TC-004 | Proposal Success Criterion #2 | CLI | cli/antipattern-scan | P0 |
| TC-005 | Proposal Success Criterion #3 | CLI | cli/antipattern-scan | P0 |
| TC-006 | Proposal Success Criterion #4 | CLI | cli/antipattern-scan | P1 |
| TC-007 | Proposal Success Criterion #5 | CLI | cli/antipattern-scan | P1 |
