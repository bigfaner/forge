---
id: "1"
title: "Extend testkit + convert gen-test-scripts and plugin-content tests"
priority: "P1"
estimated_time: "30m"
dependencies: []
scope: "backend"
breaking: false
type: "implementation"
mainSession: false
---

# 1: Extend testkit + convert gen-test-scripts and plugin-content tests

## Description

Foundation task: extend the Go testkit with helpers needed by all converted tests, then convert the first two Playwright test files as a validation of the approach.

The Playwright tests use helpers from `tests/e2e/helpers.ts`:
- `runCli(cmd)` → runs any shell command with project root as cwd
- `readProjectFile(relPath)` → reads file relative to project root
- `projectFileExists(relPath)` → checks file existence
- `PROJECT_ROOT` → project root path constant

The Go testkit (`forge-cli/tests/e2e/testkit/helpers.go`) already has `RunCLI*` for forge commands but lacks:
- Generic shell command execution (some tests run `node validate-specs.mjs`, not `forge`)
- File reading relative to project root
- File existence checking relative to project root
- Project root constant

## Reference Files
- `docs/proposals/migrate-e2e-to-go/proposal.md` — Source proposal
- `tests/e2e/helpers.ts` — Source helpers (Node.js)
- `tests/e2e/gen-test-scripts/cli.spec.ts` — Source test (7 test cases, TC-001 through TC-007)
- `tests/e2e/plugin-content/skill-content.spec.ts` — Source test (1 test case, TC-001)
- `forge-cli/tests/e2e/testkit/helpers.go` — Target testkit
- `forge-cli/tests/e2e/helpers_test.go` — Existing local helpers (reference for patterns)

## Affected Files

### Create
| File | Description |
|------|-------------|
| `forge-cli/tests/e2e/testkit/project.go` | Project root constant, ReadProjectFile, ProjectFileExists, RunShell |
| `forge-cli/tests/e2e/gen_test_scripts_cli_test.go` | Converted gen-test-scripts tests (TC-GTS-001 through TC-GTS-007) |
| `forge-cli/tests/e2e/plugin_content_test.go` | Converted plugin-content test (TC-PC-001) |

### Modify
| File | Changes |
|------|---------|
| None | |

### Delete
| File | Reason |
|------|--------|
| None | |

## Acceptance Criteria
- [ ] `testkit/project.go` exports: `ProjectRoot` (string), `ReadProjectFile(t, relPath)`, `ProjectFileExists(relPath)`, `RunShell(t, cmd, cwd)`
- [ ] `ProjectRoot` resolves to the forge project root (parent of `forge-cli/`)
- [ ] `RunShell` executes arbitrary shell commands and returns `(stdout, stderr, exitCode)`
- [ ] gen-test-scripts tests: 7 Go test functions matching TC-001 through TC-007 assertions
- [ ] plugin-content test: 1 Go test function matching TC-001 assertions
- [ ] TC numbers prefixed to avoid collision with existing Go tests (e.g., TC-GTS-001, TC-PC-001)
- [ ] `go build -tags=e2e ./tests/e2e/...` compiles without errors
- [ ] `go test ./tests/e2e/... -v -tags=e2e -run "TestTC_GTS|TestTC_PC"` passes

## Hard Rules
- All new files MUST use `//go:build e2e` build tag
- All test files MUST be `package e2e`
- Use `testkit` package for helpers, NOT local `helpers_test.go`
- TC numbering: use prefix `GTS-` for gen-test-scripts, `PC-` for plugin-content

## Implementation Notes

### testkit/project.go
```go
// ProjectRoot resolves to the forge project root (parent of forge-cli/).
// Detects by walking up from testkit/ dir looking for forge-cli/ marker.
var ProjectRoot = resolveProjectRoot()

func resolveProjectRoot() string {
    // Start from testkit/ directory, walk up until finding go.mod with module containing "forge-cli"
    // Return parent of that directory
}

func ReadProjectFile(t *testing.T, relPath string) string { ... }
func ProjectFileExists(relPath string) bool { ... }
func RunShell(t *testing.T, cmd string, cwd string) (stdout, stderr string, exitCode int) { ... }
```

### gen-test-scripts tests
- TC-GTS-001 through TC-GTS-004: Create temp fixture dirs, run `node validate-specs.mjs`, parse JSON output, assert error/warning rules
- TC-GTS-005 through TC-GTS-007: Read SKILL.md files, assert content patterns with regex
- Use `t.TempDir()` for fixture creation, `testkit.RunShell()` for node execution, `testkit.ReadProjectFile()` for file reading

### plugin-content test
- TC-PC-001: Iterate over skill files, check for forbidden raw commands, assert zero violations
- Uses `testkit.ReadProjectFile()` for reading skill files
