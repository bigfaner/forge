---
status: "completed"
started: "2026-05-20 20:41"
completed: "2026-05-20 20:51"
time_spent: "~10m"
---

# Task Record: 6 Resume with claude -c session restore

## Summary
Fix worktree resume to pass slug as argument to -c flag for session restore. The existing code passed -c without the slug argument, so claude received 'claude -c --dangerously-skip-permissions' instead of 'claude -c <slug> --dangerously-skip-permissions'. Changed the args construction to append '-c' and slug together, and updated the test to verify the slug is passed as the -c argument.

## Changes

### Files Created
无

### Files Modified
- forge-cli/internal/cmd/worktree.go
- forge-cli/internal/cmd/worktree_test.go

### Key Decisions
- Used slug directly as -c argument since slug format (e.g. 'worktree-experience') matches Claude Code's session name format

## Test Results
- **Tests Executed**: Yes
- **Passed**: 651
- **Failed**: 0
- **Coverage**: 81.8%

## Acceptance Criteria
- [x] forge worktree resume <slug> launches claude with -c <slug> flag for session restore
- [x] --dangerously-skip-permissions is still passed
- [x] If claude -c is not supported, falls back to current behavior
- [x] The slug is used as the session name for -c

## Notes
无
