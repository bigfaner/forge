---
id: "1"
title: "Delete dead code: forge e2e group, forge probe, pkg/e2e"
priority: "P0"
estimated_time: "1h"
dependencies: []
scope: "backend"
breaking: false
type: "coding.cleanup"
mainSession: false
---

# 1: Delete dead code: forge e2e group, forge probe, pkg/e2e

## Description

Remove all dead code that has been superseded by justfile delegation (forge e2e) or is only used in tests (forge probe). This includes 7 e2e subcommand files, the probe command, and the entire pkg/e2e package. Also remove their registrations from root.go.

## Reference Files
- `docs/proposals/cli-command-restructure/proposal.md` — Source proposal

## Acceptance Criteria

- `forge e2e` and all its subcommands (run, setup, compile, discover, validate-specs, verify) are removed
- `forge probe` is removed
- `pkg/e2e/` package is fully deleted
- root.go no longer registers e2e or probe commands
- `go build ./...` passes
- `go test ./...` passes
- `forge --help` does not show e2e or probe

## Hard Rules

- Do NOT modify any surviving command's behavior
- Remove registrations from root.go init() in the same commit
- Delete corresponding test files for removed commands

## Implementation Notes

Files to delete from `forge-cli/internal/cmd/`:
- e2e_parent.go, e2e_run.go, e2e_setup.go, e2e_validate_specs.go, e2e_verify.go, e2e_compile.go, e2e_discover.go
- e2e_subcommands_test.go, e2e_validate_specs_test.go
- probe.go, probe_test.go

Package to delete: `forge-cli/pkg/e2e/` (actions.go, actions_test.go, e2e.go, e2e_test.go, exec.go, exec_test.go)

In root.go init(), remove:
- `rootCmd.AddCommand(e2eCmd)` and `rootCmd.AddCommand(probeCmd)`
- The entire "E2E group subcommands" block (6 AddCommand calls)

Search for any remaining references to deleted symbols (`e2eCmd`, `probeCmd`, `e2e.*`, `probe.*`) across the codebase and clean up.
