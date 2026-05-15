---
id: "8"
title: "Convert init-justfile and plugin-content tests"
priority: "P1"
estimated_time: "45m"
dependencies: ["1"]
scope: "all"
breaking: false
type: "implementation"
mainSession: false
---

# 8: Convert init-justfile and plugin-content tests

## Description

Convert `tests/e2e/init-justfile/init-justfile.spec.ts` (7 tests) and `tests/e2e/plugin-content/skill-content.spec.ts` (1 test) to Go. Small test suites covering justfile initialization and plugin skill content validation.

## Reference Files
- `docs/proposals/migrate-e2e-to-go/proposal.md` — Source proposal
- `tests/e2e/init-justfile/init-justfile.spec.ts` — Source file (7 tests)
- `tests/e2e/plugin-content/skill-content.spec.ts` — Source file (1 test)

## Affected Files

### Create
| File | Description |
|------|-------------|
| `forge-cli/tests/e2e/init_justfile_cli_test.go` | Go test file with 7 init-justfile test cases |
| `forge-cli/tests/e2e/plugin_content_cli_test.go` | Go test file with 1 skill-content test case |

### Modify
| File | Changes |
|------|---------|

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] All 8 test cases (7 + 1) have Go test functions with matching TC numbers
- [ ] `go test ./tests/e2e/... -v -tags=e2e -run TestTC_0` passes for these tests
- [ ] `go build ./...` passes

## Hard Rules

- Build tag `//go:build e2e`
- Package `e2e`
- Preserve TC numbers from source .spec.ts exactly
- TC number collisions with existing tests are OK — the full function name (e.g. `TestTC_011_JustCompilePasses`) differs from existing (e.g. `TestTC_011_ConcurrentSubmitHandlesLockContention`). Go requires unique function names within a package, and the descriptive suffix ensures uniqueness.

## Implementation Notes

- Small task — both source files have simple test patterns. The init-justfile tests verify `forge init-justfile` creates correct justfile structure. The plugin-content test validates skill content file format.
