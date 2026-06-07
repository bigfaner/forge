---
status: "completed"
started: "2026-06-07 21:57"
completed: "2026-06-07 22:02"
time_spent: "~5m"
---

# Task Record: 4 Fix CLI help text for init/forensic/worktree/fact commands

## Summary
Fix CLI help text for 5 commands (C1, C7-C10): updated cobra Long descriptions for init, forensic search, forensic subagents, worktree status, and fact summary

## Changes

### Files Created
无

### Files Modified
- forge-cli/internal/cmd/init.go
- forge-cli/internal/cmd/forensic/commands.go
- forge-cli/internal/cmd/worktree/cmd_status.go
- forge-cli/internal/cmd/fact/summary.go

### Key Decisions
- Init Long: restructured as numbered steps to reflect actual runInit flow including surface detection (step 5/6)
- Forensic search/subagents: wrote Long descriptions by tracing RunE implementations for accurate behavior documentation
- Worktree status: added all 5 output fields (WORKTREE, BRANCH, COMMIT, UNCOMMITTED, UNPUSHED) to Long description
- Fact summary: documented [COVERAGE] section with runtime-confirmed ratio explanation

## Test Results
- **Tests Executed**: Yes
- **Passed**: 42
- **Failed**: 0
- **Coverage**: 87.1%

## Acceptance Criteria
- [x] forge init --help Long description includes surface detection step
- [x] forge forensic search --help has Long description (was empty)
- [x] forge forensic subagents --help has Long description (was empty)
- [x] forge worktree status --help Long includes UNPUSHED field
- [x] forge fact summary --help Long includes COVERAGE indicator

## Notes
All changes are cobra Long string constants only. No logic changes. Compile, fmt, lint, and all targeted tests pass.
