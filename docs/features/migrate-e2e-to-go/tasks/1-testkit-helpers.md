---
id: "1"
title: "Add file assertion helpers to testkit"
priority: "P0"
estimated_time: "30m"
dependencies: []
scope: "all"
breaking: false
type: "implementation"
mainSession: false
---

# 1: Add file assertion helpers to testkit

## Description

The Playwright tests in `tests/e2e/helpers.ts` export `readProjectFile()` and `projectFileExists()` helpers that are used by nearly every test file. The Go `testkit` package (`forge-cli/tests/e2e/testkit/helpers.go`) currently only has CLI execution helpers (`RunCLI`, `RunCLIExitCode`, `RunCLIWithResult`, `WithRetry`). Add equivalent file assertion helpers so subsequent conversion tasks can use them.

## Reference Files
- `docs/proposals/migrate-e2e-to-go/proposal.md` — Source proposal
- `tests/e2e/helpers.ts` — Source helpers to port
- `forge-cli/tests/e2e/testkit/helpers.go` — Target file to extend
- `forge-cli/tests/e2e/helpers_test.go` — Existing unexported helpers (reference only)

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| `forge-cli/tests/e2e/testkit/helpers.go` | Add `FileContains`, `FileNotContains`, `ReadProjectFile`, `ProjectFileExists` functions |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] `testkit.FileContains(t, filePath, substring)` passes test if file contains substring
- [ ] `testkit.FileNotContains(t, filePath, substring)` passes test if file does NOT contain substring
- [ ] `testkit.ReadProjectFile(relPath)` reads file relative to project root
- [ ] `testkit.ProjectFileExists(relPath)` returns bool for file existence
- [ ] `go build ./...` passes
- [ ] `go test ./tests/e2e/... -tags=e2e -run TestNothing` compiles without errors (no existing tests break)

## Hard Rules

- All functions must use `t.Helper()` pattern consistent with existing testkit functions
- Functions must resolve project root using runtime detection (same approach as helpers.ts `PROJECT_ROOT`)
- File paths must be OS-agnostic (use `filepath.Join`)

## Implementation Notes

- The existing `helpers_test.go` has unexported versions (`runCLI`, `runCLIWithResult`, etc.) that duplicate testkit. The new helpers should be exported in `testkit` so they can be used by all test files including feature-scoped ones in different packages.
- Consider also adding `WriteFile` and `MkdirAll` helpers for test fixture setup, as many Playwright tests create temporary files/directories.
