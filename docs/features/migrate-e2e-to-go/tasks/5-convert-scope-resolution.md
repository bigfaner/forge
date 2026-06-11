---
id: "5"
title: "Convert scope-resolution tests"
priority: "P1"
estimated_time: "45m"
dependencies: ["1"]
scope: "all"
breaking: false
type: "implementation"
mainSession: false
---

# 5: Convert scope-resolution tests

## Description

Convert `tests/e2e/scope-resolution/scope-resolution.spec.ts` (8 tests) to Go. These tests validate the scope resolution algorithm for quality gate commands (compile, fmt, lint, test).

## Reference Files
- `docs/proposals/migrate-e2e-to-go/proposal.md` — Source proposal
- `tests/e2e/scope-resolution/scope-resolution.spec.ts` — Source file (8 tests)

## Affected Files

### Create
| File | Description |
|------|-------------|
| `forge-cli/tests/e2e/scope_resolution_cli_test.go` | Go test file with 8 test cases |

### Modify
| File | Changes |
|------|---------|

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] All 8 test cases have Go test functions with matching TC numbers
- [ ] Scope resolution assertions correctly verify forge command behavior with different project types
- [ ] `go test ./tests/e2e/... -v -tags=e2e -run TestTC_0` passes for these tests
- [ ] `go build ./...` passes

## Hard Rules

- Build tag `//go:build e2e`
- Package `e2e`
- Preserve TC numbers from source .spec.ts exactly
- TC number collisions with existing tests are OK — the full function name (e.g. `TestTC_011_JustCompilePasses`) differs from existing (e.g. `TestTC_011_ConcurrentSubmitHandlesLockContention`). Go requires unique function names within a package, and the descriptive suffix ensures uniqueness.

## Implementation Notes

- These tests may verify that `forge config get project-type` correctly determines scope for just commands. They may need mock project structures with specific file layouts.
