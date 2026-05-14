---
id: "3"
title: "Convert task-cli, scope-resolution, and init-justfile tests to Go"
priority: "P1"
estimated_time: "30m"
dependencies: ["1"]
scope: "backend"
breaking: false
type: "implementation"
mainSession: false
---

# 3: Convert task-cli, scope-resolution, and init-justfile tests to Go

## Description

Convert the remaining 3 Playwright test files (35 test cases total) to Go:
- `task-cli/typed-task-dispatch.spec.ts` (20 tests) — Task CLI type routing, templates, prompts, migration
- `scope-resolution/scope-resolution.spec.ts` (8 tests) — Scope field in index.json, mixed/frontend/backend behavior
- `init-justfile/init-justfile.spec.ts` (7 tests) — Init-justfile project detection, standard commands, boundary markers

These tests validate the forge CLI's task management and justfile initialization features.

## Reference Files
- `docs/proposals/migrate-e2e-to-go/proposal.md` — Source proposal
- `tests/e2e/task-cli/typed-task-dispatch.spec.ts` — 20 tests (TC-001 to TC-020)
- `tests/e2e/scope-resolution/scope-resolution.spec.ts` — 8 tests (TC-001 to TC-008)
- `tests/e2e/init-justfile/init-justfile.spec.ts` — 7 tests (TC-001 to TC-007)

## Affected Files

### Create
| File | Description |
|------|-------------|
| `forge-cli/tests/e2e/typed_task_dispatch_test.go` | Converted task-cli tests (TC-TD-001 to TC-TD-020) |
| `forge-cli/tests/e2e/scope_resolution_test.go` | Converted scope-resolution tests (TC-SR-001 to TC-SR-008) |
| `forge-cli/tests/e2e/init_justfile_test.go` | Converted init-justfile tests (TC-IJ-001 to TC-IJ-007) |

### Modify
| File | Changes |
|------|---------|
| None | |

### Delete
| File | Reason |
|------|--------|
| None | |

## Acceptance Criteria
- [ ] 35 Go test functions matching all Playwright test assertions
- [ ] TC numbers prefixed: TD- (task-dispatch), SR- (scope-resolution), IJ- (init-justfile)
- [ ] `go build -tags=e2e ./tests/e2e/...` compiles without errors
- [ ] `go test ./tests/e2e/... -v -tags=e2e -run "TestTC_TD|TestTC_SR|TestTC_IJ"` passes

## Hard Rules
- All files MUST use `//go:build e2e` build tag and `package e2e`
- Use `testkit.RunCLIExitCode()` for forge commands, `testkit.ReadProjectFile()` for file reads

## Implementation Notes

### task-cli tests (20 tests):
- Test task type routing: `forge task add --type implementation/fix/...`
- Test new type templates: correct YAML frontmatter generated
- Test task prompt generation: `forge task prompt <id>`
- Test task migration: `forge task migrate`
- Test breakdown-tasks and execute-task integration
- Test error-fixer behavior
- Test task validation: `forge task validate-index`
- Test phase detection
- Uses index backup/restore pattern for tests that modify shared state

### scope-resolution tests (8 tests):
- Test mixed project scope dispatch: `just <verb> frontend/backend`
- Test frontend-only scope: no scope argument
- Test cross-scope: all tasks regardless of scope
- Test fallback behavior when project-type is unknown
- Tests create temp directories with `.forge/` config and verify `forge config get project-type` output

### init-justfile tests (7 tests):
- Test frontend/backend/mixed project detection
- Test standard commands generation
- Test boundary markers in generated justfiles
- Test error handling for invalid configurations
