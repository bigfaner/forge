---
status: "completed"
started: "2026-06-01 20:59"
completed: "2026-06-01 21:08"
time_spent: "~9m"
---

# Task Record: 1 Create install.sh and install.ps1 scripts

## Summary
Created install.sh (macOS/Linux) and install.ps1 (Windows) scripts that download pre-compiled forge CLI binaries from GitHub Releases, with platform detection, atomic binary replacement, and PATH management.

## Changes

### Files Created
- forge-cli/scripts/install.sh
- forge-cli/scripts/install.ps1
- forge-cli/scripts/install_test.go

### Files Modified
无

### Key Decisions
- Reused platform detection, PATH management, and atomic replace patterns from install-local.sh and install-local.ps1
- install.sh iterates all three RC files (.bashrc, .zshrc, .profile) rather than only the detected shell, for better multi-shell support
- install.ps1 does NOT handle the running-exe rename dance per task spec — that is only needed in forge upgrade (Task 2)
- Tag uses v prefix (forge-cli/v{version}), binary filename does NOT use v prefix (forge-{version}-{os}-{arch}) per Hard Rule

## Test Results
- **Tests Executed**: Yes
- **Passed**: 22
- **Failed**: 0
- **Coverage**: 0.0%

## Acceptance Criteria
- [x] install.sh detects OS (darwin/linux) and architecture (amd64/arm64)
- [x] install.sh fetches latest version from GitHub Release API, constructs download URL with tag format forge-cli/v{version} and binary name forge-{version}-{os}-{arch}
- [x] install.sh downloads binary to ~/.forge/bin/forge.new and atomically replaces (mv to ~/.forge/bin/forge)
- [x] install.sh adds ~/.forge/bin/ to PATH in shell RC files (.bashrc, .zshrc, .profile)
- [x] install.ps1 handles Windows (amd64/arm64), downloads to %USERPROFILE%\.forge\bin\, updates User PATH via [Environment]::SetEnvironmentVariable
- [x] Both scripts output verification instructions after successful installation

## Notes
22 Go tests verify script contents against all AC items and Hard Rules. No Go source code in scripts/ directory, so coverage shows 'no statements'. All static checks (compile, fmt, lint) pass.
