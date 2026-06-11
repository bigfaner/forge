---
id: "1"
title: "Create install.sh and install.ps1 scripts"
priority: "P0"
estimated_time: "2h"
complexity: "medium"
dependencies: []
surface-key: "."
surface-type: "cli"
breaking: false
type: "coding.feature"
mainSession: false
---

# 1: Create install.sh and install.ps1 scripts

## Description

Create the curl-pipe-bash installation script (`install.sh`) for macOS/Linux and the PowerShell equivalent (`install.ps1`) for Windows. These scripts are the primary entry point for users installing forge CLI for the first time — they download a pre-compiled binary from GitHub Releases instead of building from source.

## Reference Files

- `forge-cli/scripts/install-local.sh`: existing local build+install script — reuse platform detection, PATH management, and atomic replace patterns (source: proposal.md#Implementation-1)
- `forge-cli/scripts/install-local.ps1`: existing Windows PowerShell install script — reuse Windows-specific patterns (source: proposal.md#Implementation-1)
- `forge-cli/scripts/version.txt`: version format reference, tag = `forge-cli/v{version}`, binary = `forge-{version}-{os}-{arch}` (source: proposal.md#Implementation-1)

## Acceptance Criteria

- [ ] install.sh detects OS (darwin/linux) and architecture (amd64/arm64) from the running system
- [ ] install.sh fetches latest version from GitHub Release API, constructs download URL using tag format `forge-cli/v{version}` and binary name `forge-{version}-{os}-{arch}`
- [ ] install.sh downloads binary to `~/.forge/bin/forge.new` and atomically replaces (`mv` to `~/.forge/bin/forge`)
- [ ] install.sh adds `~/.forge/bin/` to PATH in shell RC files (`.bashrc`, `.zshrc`, `.profile`)
- [ ] install.ps1 handles Windows (amd64/arm64), downloads to `%USERPROFILE%\.forge\bin\`, updates User PATH via `[Environment]::SetEnvironmentVariable`
- [ ] Both scripts output verification instructions after successful installation

## Hard Rules

- Tag uses `v` prefix (e.g. `forge-cli/v5.17.0`), binary filename does NOT use `v` prefix (e.g. `forge-5.17.0-darwin-arm64`)

## Implementation Notes

- Reuse the platform detection, PATH management, and atomic replace logic from `install-local.sh` and `install-local.ps1`
- The new scripts differ from install-local: they download from GitHub Releases instead of building locally
- install.ps1 does NOT need to handle the "running exe" rename dance — that's only needed in the `forge upgrade` command (Task 2) because the running process holds a file lock on its own binary
