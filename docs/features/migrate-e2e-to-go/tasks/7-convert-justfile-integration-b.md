---
id: "7"
title: "Convert justfile-e2e-integration tests (mixed-template + cli)"
priority: "P1"
estimated_time: "1.5h"
dependencies: ["1"]
scope: "all"
breaking: false
type: "implementation"
mainSession: false
---

# 7: Convert justfile-e2e-integration tests (mixed-template + cli)

## Description

Convert `tests/e2e/justfile-e2e-integration/mixed-template.spec.ts` (23 tests) and `tests/e2e/justfile-e2e-integration/cli.spec.ts` (20 tests) to Go. These tests validate justfile mixed template handling and CLI integration commands.

## Reference Files
- `docs/proposals/migrate-e2e-to-go/proposal.md` — Source proposal
- `tests/e2e/justfile-e2e-integration/mixed-template.spec.ts` — Source file (23 tests)
- `tests/e2e/justfile-e2e-integration/cli.spec.ts` — Source file (20 tests)

## Affected Files

### Create
| File | Description |
|------|-------------|
| `forge-cli/tests/e2e/justfile_mixed_cli_cli_test.go` | Go test file with 43 test cases from both source files |

### Modify
| File | Changes |
|------|---------|

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] All 43 test cases (23 + 20) have Go test functions with matching TC numbers
- [ ] Mixed template and CLI integration assertions work correctly
- [ ] `go test ./tests/e2e/... -v -tags=e2e -run TestTC_0` passes for these tests
- [ ] `go build ./...` passes

## Hard Rules

- Build tag `//go:build e2e`
- Package `e2e`
- Preserve TC numbers from source .spec.ts exactly
- TC number collisions with existing tests are OK — the full function name (e.g. `TestTC_011_JustCompilePasses`) differs from existing (e.g. `TestTC_011_ConcurrentSubmitHandlesLockContention`). Go requires unique function names within a package, and the descriptive suffix ensures uniqueness.

## Implementation Notes

- This is the largest conversion task (43 tests). The mixed-template tests verify that justfiles with mixed recipe types are handled correctly. CLI tests verify justfile-related CLI commands.
- Both source files are from the same `justfile-e2e-integration/` directory — combine into one Go file to keep related tests together.
