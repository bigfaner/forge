---
status: "completed"
started: "2026-05-20 17:56"
completed: "2026-05-20 18:06"
time_spent: "~10m"
---

# Task Record: 5 Add worktree status subcommand

## Summary
Added `forge worktree status [<slug>]` subcommand that displays worktree status including branch name, latest commit (hash + message), and uncommitted files list. When no slug is provided, shows status for all forge-managed worktrees. Uses structured output format (PrintBlockStart/PrintField/PrintBlockEnd pattern). The command is strictly read-only.

## Changes

### Files Created
无

### Files Modified
- forge-cli/internal/cmd/worktree.go
- forge-cli/internal/cmd/worktree_test.go
- forge-cli/internal/cmd/root.go
- forge-cli/scripts/version.txt

### Key Decisions
- Used gitRunFunc (existing injectable function) for git log/status calls in printWorktreeStatus to maintain testability without adding new injection points
- Reused listForgeFeatures and listWorktreesFunc for slug resolution and forge-managed detection
- For 'no slug' mode, filters to only forge-managed worktrees (non-main + matching feature slug) to avoid noise from the main worktree
- Structured output uses --- delimiters with KEY: VALUE format consistent with project's output.go pattern

## Test Results
- **Tests Executed**: Yes
- **Passed**: 10
- **Failed**: 0
- **Coverage**: 94.7%

## Acceptance Criteria
- [x] forge worktree status <slug> shows: branch name, latest commit (hash + message), uncommitted files list
- [x] forge worktree status (no slug) shows status for all forge-managed worktrees
- [x] Output uses structured format (consistent with project's output.go pattern)
- [x] Non-existent slug prints clear error message
- [x] Command is strictly read-only — never modifies any file

## Notes
Bumped version from 4.6.0 to 4.7.0 (minor: new command). Coverage measured at function level for status-related functions: runWorktreeStatus 84.6%, showWorktreeStatus 100%, showAllWorktreeStatus 100%, printWorktreeStatus 94.7%.
