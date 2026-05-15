---
id: "6"
title: "Convert justfile-e2e-integration tests (forge-justfile + detection-assembly)"
priority: "P1"
estimated_time: "1.5h"
dependencies: ["1"]
scope: "all"
breaking: false
type: "implementation"
mainSession: false
---

# 6: Convert justfile-e2e-integration tests (forge-justfile + detection-assembly)

## Description

Convert `tests/e2e/justfile-e2e-integration/forge-justfile.spec.ts` (15 tests) and `tests/e2e/justfile-e2e-integration/detection-assembly.spec.ts` (19 tests) to Go. These tests validate forge-justfile integration: justfile detection, assembly, and recipe execution.

## Reference Files
- `docs/proposals/migrate-e2e-to-go/proposal.md` — Source proposal
- `tests/e2e/justfile-e2e-integration/forge-justfile.spec.ts` — Source file (15 tests)
- `tests/e2e/justfile-e2e-integration/detection-assembly.spec.ts` — Source file (19 tests)

## Affected Files

### Create
| File | Description |
|------|-------------|
| `forge-cli/tests/e2e/justfile_forge_detection_cli_test.go` | Go test file with 34 test cases from both source files |

### Modify
| File | Changes |
|------|---------|

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] All 34 test cases (15 + 19) have Go test functions with matching TC numbers
- [ ] Justfile detection and assembly assertions work correctly
- [ ] `go test ./tests/e2e/... -v -tags=e2e -run TestTC_0` passes for these tests
- [ ] `go build ./...` passes

## Hard Rules

- Build tag `//go:build e2e`
- Package `e2e`
- Preserve TC numbers from source .spec.ts exactly
- TC number collisions with existing tests are OK — the full function name (e.g. `TestTC_011_JustCompilePasses`) differs from existing (e.g. `TestTC_011_ConcurrentSubmitHandlesLockContention`). Go requires unique function names within a package, and the descriptive suffix ensures uniqueness.

## Implementation Notes

- These tests create justfile structures and verify forge correctly detects/assembles them. Map fixture creation to `t.TempDir()` + `os.WriteFile`.
- Detection-assembly tests may verify that `init:justfile` correctly sets up justfile templates.
