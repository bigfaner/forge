---
id: "4"
title: "Unit and integration tests for ensureJust flow"
priority: "P2"
estimated_time: "2h"
dependencies: ["3"]
scope: "backend"
breaking: false
type: "implementation"
mainSession: false
---

# 4: Unit and integration tests for ensureJust flow

## Description

Write comprehensive tests covering the full ensureJust flow: detection, version parsing, package manager dispatch, embedded binary fallback, and integration with `forge init`. Tests should cover success paths, failure paths, and edge cases.

## Reference Files
- `docs/proposals/forge-init-install-just/proposal.md` — Source proposal
- `forge-cli/pkg/just/ensure.go` — ensureJust logic (from task 2)
- `forge-cli/pkg/just/ensure_test.go` — Partial unit tests (from task 2, extend here)
- `forge-cli/internal/cmd/init_test.go` — Init command tests (from task 3, extend here)

## Affected Files

### Create
| File | Description |
|------|-------------|
| `forge-cli/pkg/just/ensure_integration_test.go` | Integration tests for the full ensureJust flow |

### Modify
| File | Changes |
|------|---------|
| `forge-cli/pkg/just/ensure_test.go` | Extend with additional edge case tests |
| `forge-cli/internal/cmd/init_test.go` | Extend with ensureJust step integration tests |

### Delete
| File | Reason |
|------|--------|
| (none) | |

## Acceptance Criteria

- [ ] Unit tests cover: `DetectJust` (found/not found), `ParseJustVersion` (valid/invalid formats), `IsMinimumVersion` (equal/above/below/edge cases)
- [ ] Unit tests cover: package manager dispatch logic per OS (mock exec commands)
- [ ] Unit tests cover: embedded binary extraction to `~/.forge/bin/` (with temp dirs)
- [ ] Integration test: `forge init` with `--skip-just` skips ensureJust step
- [ ] Integration test: `forge init` with just already installed reports SKIPPED
- [ ] Integration test: `forge init` without just triggers installation attempt
- [ ] Edge cases: Windows `.exe` extension, empty version output, permission denied on extraction
- [ ] All tests use table-driven patterns where applicable
- [ ] Test coverage for `pkg/just/ensure.go` >= 80%

## Hard Rules

- Mock external command execution (`exec.Command`) — do not require actual `just` or package managers on CI
- Use `t.TempDir()` for filesystem tests — never write to real `~/.forge/bin/`
- Integration tests may use build tags (`//go:build !e2e`) if they require real `just` binary

## Implementation Notes

- Use a command executor interface to allow mocking in tests: `type Executor interface { LookPath(name string) (string, error); Run(name string, args ...string) (string, error) }`
- Test version parsing with real-world `just --version` output formats: `just 1.40.0`, `just 1.37.0`
- Test binary extraction by creating a fake "binary" (byte slice) and verifying file contents + permissions
- For init integration tests, follow the existing pattern in `init_test.go` (buffered stdin/stdout)
