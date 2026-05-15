---
id: "2"
title: "Convert gen-test-scripts and forge-testing-optimization tests"
priority: "P1"
estimated_time: "1h"
dependencies: ["1"]
scope: "all"
breaking: false
type: "implementation"
mainSession: false
---

# 2: Convert gen-test-scripts and forge-testing-optimization tests

## Description

Convert `tests/e2e/gen-test-scripts/cli.spec.ts` (16 tests) and `tests/e2e/features/forge-testing-optimization/cli.spec.ts` (16 tests, duplicate of gen-test-scripts) to Go. The proposal identifies these as duplicates — merge into a single Go test file. All tests validate the `gen-test-scripts` skill's spec validation logic.

## Reference Files
- `docs/proposals/migrate-e2e-to-go/proposal.md` — Source proposal
- `tests/e2e/gen-test-scripts/cli.spec.ts` — Primary source (16 tests)
- `tests/e2e/features/forge-testing-optimization/cli.spec.ts` — Duplicate source (merge into same Go file)

## Affected Files

### Create
| File | Description |
|------|-------------|
| `forge-cli/tests/e2e/gen_test_scripts_cli_test.go` | Merged Go test file with all unique test cases |

### Modify
| File | Changes |
|------|---------|

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] All unique test cases from both .spec.ts files have Go test functions with `TestTC_NNN_Description` naming
- [ ] Test fixture setup (creating temp spec files with deliberate violations) maps to `t.TempDir()` + `os.WriteFile`
- [ ] `runCli()` calls map to `testkit.RunCLIExitCode()` / `testkit.RunCLIWithResult()`
- [ ] File content assertions use new `testkit.FileContains()` / `testkit.FileNotContains()` helpers
- [ ] `go test ./tests/e2e/... -v -tags=e2e -run TestTC_0` passes for these tests
- [ ] `go build ./...` passes

## Hard Rules

- Build tag `//go:build e2e` on every test file
- Package `e2e` (matches existing test files)
- Preserve TC numbers from source .spec.ts files exactly

## Implementation Notes

- These tests validate the `validate-specs.mjs` script — they create fixture spec files with E1-E4 violations and run validation. The Go equivalent will need to create temp directories with the same fixture structure.
- The tests use `writeFileSync` to create fixture files — map to `os.WriteFile` or `testkit` helper.
- Deduplicate: if both files have identical TC numbers, keep only one copy.
