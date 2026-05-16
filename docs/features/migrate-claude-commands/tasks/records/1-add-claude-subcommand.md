---
status: "completed"
started: "2026-05-16 23:25"
completed: "2026-05-16 23:35"
time_spent: "~10m"
---

# Task Record: 1 Add forge claude subcommand with arg passthrough

## Summary
Add forge claude subcommand that launches Claude CLI with --dangerously-skip-permissions always injected. Uses DisableFlagParsing for transparent arg passthrough, exec.LookPath for pre-flight binary validation, and syscall.Exec for process replacement.

## Changes

### Files Created
- forge-cli/internal/cmd/claude.go
- forge-cli/internal/cmd/claude_test.go

### Files Modified
- forge-cli/internal/cmd/root.go
- forge-cli/internal/cmd/root_test.go
- forge-cli/scripts/version.txt

### Key Decisions
- Used DisableFlagParsing: true on cobra command to avoid flag parsing conflicts -- all user args pass through transparently without -- separator
- Injected --dangerously-skip-permissions as first arg -- always prepended, not configurable per hard rules
- Used syscall.Exec for process replacement instead of exec.Command -- avoids spawning a child process
- Extracted lookPathFunc and runClaudeFunc as testable variables for dependency injection in tests

## Test Results
- **Tests Executed**: Yes
- **Passed**: 9
- **Failed**: 0
- **Coverage**: 81.2%

## Acceptance Criteria
- [x] forge claude launches Claude CLI with --dangerously-skip-permissions
- [x] forge claude -c continues the last conversation
- [x] forge claude -w <name> opens a worktree session
- [x] Any Claude CLI flag passes through: forge claude --model opus -p "prompt"
- [x] Clear error when claude binary is not in PATH
- [x] Unit tests for: PATH validation, arg passthrough, flag injection

## Notes
Version bumped from 3.16.0 to 3.17.0 (minor: new command). Updated root_test.go command counts to reflect the new claude subcommand.
