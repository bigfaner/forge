---
id: "4"
title: "Convert task-cli typed-task-dispatch tests"
priority: "P1"
estimated_time: "1.5h"
dependencies: ["1"]
scope: "all"
breaking: false
type: "implementation"
mainSession: false
---

# 4: Convert task-cli typed-task-dispatch tests

## Description

Convert `tests/e2e/task-cli/typed-task-dispatch.spec.ts` (20 tests) to Go. These tests validate the task type dispatch system — ensuring that tasks with different `type` fields (implementation, fix, etc.) are dispatched to the correct handler.

## Reference Files
- `docs/proposals/migrate-e2e-to-go/proposal.md` — Source proposal
- `tests/e2e/task-cli/typed-task-dispatch.spec.ts` — Source file (20 tests)

## Affected Files

### Create
| File | Description |
|------|-------------|
| `forge-cli/tests/e2e/task_types_dispatch_cli_test.go` | Go test file with 20 test cases |

### Modify
| File | Changes |
|------|---------|

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] All 20 test cases have Go test functions with matching TC numbers
- [ ] Task dispatch assertions correctly verify exit codes and output
- [ ] `go test ./tests/e2e/... -v -tags=e2e -run TestTC_0` passes for these tests
- [ ] `go build ./...` passes

## Hard Rules

- Build tag `//go:build e2e`
- Package `e2e`
- Preserve TC numbers from source .spec.ts exactly
- TC number collisions with existing tests are OK — the full function name (e.g. `TestTC_011_JustCompilePasses`) differs from existing (e.g. `TestTC_011_ConcurrentSubmitHandlesLockContention`). Go requires unique function names within a package, and the descriptive suffix ensures uniqueness.
- Do NOT modify existing `task_types_cli_test.go` — this is a new file for different test cases

## Implementation Notes

- This is the largest single conversion. Some tests may need feature fixture setup (creating task markdown files with specific frontmatter).
- The existing `task_types_cli_test.go` already covers some task type tests — check for TC number overlaps and avoid duplication.
