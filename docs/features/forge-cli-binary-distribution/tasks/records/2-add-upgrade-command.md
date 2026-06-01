---
status: "completed"
started: "2026-06-01 21:09"
completed: "2026-06-01 21:26"
time_spent: "~17m"
---

# Task Record: 2 Implement forge upgrade CLI subcommand

## Summary
Implemented forge upgrade CLI subcommand with two-phase upgrade: CLI binary self-update from GitHub Releases API and Plugin install/update via claude CLI. Includes Windows rename dance for safe in-place binary replacement, version comparison logic, and unified summary output.

## Changes

### Files Created
- forge-cli/internal/cmd/upgrade.go
- forge-cli/internal/cmd/upgrade_test.go

### Files Modified
- forge-cli/internal/cmd/root.go
- forge-cli/internal/cmd/root_test.go

### Key Decisions
- Two independent phases (CLI binary + Plugin) that succeed or fail independently
- Function variables for all external dependencies (exec.LookPath, http.Get, base.RunClaude) for testability
- Windows rename dance: forge.exe -> forge.old, write new, delete forge.old per Hard Rules
- Unix atomic replace via forge.new -> forge rename

## Test Results
- **Tests Executed**: Yes
- **Passed**: 27
- **Failed**: 0
- **Coverage**: 81.5%

## Acceptance Criteria
- [x] New upgrade subcommand registered in Cobra command tree with prerequisite check (claude CLI in PATH)
- [x] CLI binary upgrade: compare current version with latest from GitHub Release API (parse tag forge-cli/v{version}), skip if same version
- [x] Binary download + atomic replace at ~/.forge/bin/forge; Windows special handling: rename old binary to forge.old before write, delete forge.old after
- [x] Plugin management: detect if marketplace added -> claude plugin marketplace add if missing; detect plugin installed -> claude plugin install forge or claude plugin update forge
- [x] Unified output showing results for both CLI and Plugin operations
- [x] Unit tests for version comparison logic and download URL construction

## Notes
replaceUnixBinary has 0% coverage on Windows (expected - tested via TestReplaceWindowsBinary). Plugin management functions have partial coverage because they delegate to external claude CLI binary.
