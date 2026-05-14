---
id: "2"
title: "Convert justfile-related Playwright tests to Go"
priority: "P1"
estimated_time: "45m"
dependencies: ["1"]
scope: "backend"
breaking: false
type: "implementation"
mainSession: false
---

# 2: Convert justfile-related Playwright tests to Go

## Description

Convert the 5 justfile-related Playwright test files (86 test cases total) to Go. These tests validate:
- Justfile recipe execution (`justfile-execution`, 9 tests)
- Forge justfile structure and standard recipes (`forge-justfile`, 15 tests)
- Project detection and template assembly (`detection-assembly`, 19 tests)
- Mixed template handling (`mixed-template`, 23 tests)
- Skill/agent file content validation (`cli.spec`, 20 tests)

All tests follow the same pattern: run `forge` CLI commands or read project files, then assert on output/content.

## Reference Files
- `docs/proposals/migrate-e2e-to-go/proposal.md` — Source proposal
- `tests/e2e/justfile-execution/justfile-execution.spec.ts` — 9 tests (TC-001 to TC-009)
- `tests/e2e/justfile-e2e-integration/forge-justfile.spec.ts` — 15 tests (TC-001 to TC-015)
- `tests/e2e/justfile-e2e-integration/detection-assembly.spec.ts` — 19 tests (TC-001 to TC-019)
- `tests/e2e/justfile-e2e-integration/mixed-template.spec.ts` — 23 tests (TC-001 to TC-023)
- `tests/e2e/justfile-e2e-integration/cli.spec.ts` — 20 tests (TC-001 to TC-020)

## Affected Files

### Create
| File | Description |
|------|-------------|
| `forge-cli/tests/e2e/justfile_execution_test.go` | Converted justfile-execution tests (TC-JE-001 to TC-JE-009) |
| `forge-cli/tests/e2e/forge_justfile_test.go` | Converted forge-justfile tests (TC-FJ-001 to TC-FJ-015) |
| `forge-cli/tests/e2e/detection_assembly_test.go` | Converted detection-assembly tests (TC-DA-001 to TC-DA-019) |
| `forge-cli/tests/e2e/mixed_template_test.go` | Converted mixed-template tests (TC-MT-001 to TC-MT-023) |
| `forge-cli/tests/e2e/justfile_skill_agent_test.go` | Converted justfile-e2e-integration/cli tests (TC-SA-001 to TC-SA-020) |

### Modify
| File | Changes |
|------|---------|
| None | |

### Delete
| File | Reason |
|------|--------|
| None | |

## Acceptance Criteria
- [ ] 86 Go test functions matching all Playwright test assertions
- [ ] TC numbers prefixed: JE- (justfile-execution), FJ- (forge-justfile), DA- (detection-assembly), MT- (mixed-template), SA- (skill-agent)
- [ ] `go build -tags=e2e ./tests/e2e/...` compiles without errors
- [ ] `go test ./tests/e2e/... -v -tags=e2e -run "TestTC_JE|TestTC_FJ|TestTC_DA|TestTC_MT|TestTC_SA"` passes

## Hard Rules
- All files MUST use `//go:build e2e` build tag and `package e2e`
- Use `testkit.RunCLIExitCode()` for forge commands, `testkit.ReadProjectFile()` for file reads
- Use `testkit.ProjectRoot` for resolving project-relative paths

## Implementation Notes

### Common patterns across these tests:
1. **CLI execution + output assertion**: `testkit.RunCLIExitCode("task", args...)` then `assert.Contains(t, output, expected)`
2. **File content assertion**: `testkit.ReadProjectFile(t, relPath)` then `assert.Regexp(t, pattern, content)`
3. **Temp fixture setup**: `t.TempDir()` + `os.MkdirAll` + `os.WriteFile` for tests that need temp project structures

### Test-specific notes:
- `forge-justfile` tests verify the 15 standard justfile recipes exist and have correct structure
- `detection-assembly` tests verify project type detection (frontend/backend/mixed) and template selection
- `mixed-template` tests verify the mixed project justfile template with scoped and unscoped recipes
- `cli.spec` (justfile-e2e-integration) tests verify skill/agent/command files use `just <verb>` exclusively
- `justfile-execution` tests verify actual `just compile`/`just build`/`just test` execution
