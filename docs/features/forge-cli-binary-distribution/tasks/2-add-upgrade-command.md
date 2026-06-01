---
id: "2"
title: "Implement forge upgrade CLI subcommand"
priority: "P0"
estimated_time: "2h"
complexity: "high"
dependencies: []
surface-key: "."
surface-type: "cli"
breaking: false
type: "coding.feature"
mainSession: false
---

# 2: Implement forge upgrade CLI subcommand

## Description

Add a new `forge upgrade` CLI subcommand that handles both CLI binary self-update and Plugin installation/upgrade. This is the unified upgrade path for all users after initial installation. Prerequisite: `claude` CLI must be in PATH.

## Reference Files

- `forge-cli/internal/cmd/root.go`: Cobra root command registration — add `upgrade` subcommand here (source: proposal.md#Implementation-2)
- `forge-cli/pkg/types`: `Version` variable injected via `-ldflags`, used to compare current vs latest (source: proposal.md#Implementation-3)
- `forge-cli/scripts/version.txt`: version format, current value `5.16.0` (source: proposal.md#Implementation-2)
- `forge-cli/internal/cmd/init.go`: reference for command registration patterns and output styling (source: proposal.md#Implementation-2)

## Acceptance Criteria

- [ ] New `upgrade` subcommand registered in Cobra command tree with prerequisite check (`claude` CLI in PATH)
- [ ] CLI binary upgrade: compare current version with latest from GitHub Release API (parse tag `forge-cli/v{version}`), skip if same version
- [ ] Binary download + atomic replace at `~/.forge/bin/forge`; Windows special handling: rename old binary to `forge.old` before write, delete `forge.old` after (Windows cannot replace a running exe)
- [ ] Plugin management: detect if marketplace added → `claude plugin marketplace add` if missing; detect plugin installed → `claude plugin install forge` or `claude plugin update forge`
- [ ] Unified output showing results for both CLI and Plugin operations
- [ ] Unit tests for version comparison logic and download URL construction

## Hard Rules

- Windows rename dance: `forge.old` → write new → delete `forge.old`. This is the ONLY safe pattern on Windows when the running process is the binary being replaced.

## Implementation Notes

- The upgrade flow has two independent phases: (1) CLI binary update, (2) Plugin install/update. Each phase should succeed or fail independently — a Plugin failure should not roll back a successful CLI update.
- GitHub Release API URL: `https://api.github.com/repos/bigfaner/forge/releases/latest`
- Download URL pattern: `https://github.com/bigfaner/forge/releases/download/forge-cli/v{version}/forge-{version}-{os}-{arch}`
- Plugin marketplace add command: `claude plugin marketplace add https://github.com/bigfaner/forge.git --sparse .claude-plugin plugins`
