---
id: "3"
title: "Convert justfile-execution tests"
priority: "P1"
estimated_time: "45m"
dependencies: ["1"]
scope: "all"
breaking: false
type: "implementation"
mainSession: false
---

# 3: Convert justfile-execution tests

## Description

Convert `tests/e2e/justfile-execution/justfile-execution.spec.ts` (9 tests) to Go. These tests validate justfile recipe execution through the forge CLI, covering task compilation, fmt, lint, and test commands.

## Reference Files
- `docs/proposals/migrate-e2e-to-go/proposal.md` — Source proposal
- `tests/e2e/justfile-execution/justfile-execution.spec.ts` — Source file (9 tests)

## Affected Files

### Create
| File | Description |
|------|-------------|
| `forge-cli/tests/e2e/justfile_execution_cli_test.go` | Go test file with 9 test cases |

### Modify
| File | Changes |
|------|---------|

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] All 9 test cases have Go test functions with matching TC numbers
- [ ] `runCli()` calls map to `testkit.RunCLIExitCode()` / `testkit.RunCLIWithResult()`
- [ ] File content assertions use `testkit.FileContains()` / `testkit.FileNotContains()`
- [ ] `go test ./tests/e2e/... -v -tags=e2e -run TestTC_0` passes for these tests
- [ ] `go build ./...` passes

## Hard Rules

- Build tag `//go:build e2e`
- Package `e2e`
- Preserve TC numbers from source .spec.ts exactly
- TC number collisions with existing tests are OK — the full function name (e.g. `TestTC_011_JustCompilePasses`) differs from existing (e.g. `TestTC_011_ConcurrentSubmitHandlesLockContention`). Go requires unique function names within a package, and the descriptive suffix ensures uniqueness.

## Implementation Notes

- These tests likely create temporary justfiles and run forge commands. Map `beforeAll`/`afterAll` to `t.TempDir()` + cleanup.
- If tests depend on justfile existence in specific paths, ensure the Go version uses the same fixture approach.
