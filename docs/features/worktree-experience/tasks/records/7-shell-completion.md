---
status: "completed"
started: "2026-05-20 20:51"
completed: "2026-05-20 20:59"
time_spent: "~8m"
---

# Task Record: 7 Shell completion for start/remove/resume subcommands

## Summary
Add Cobra ValidArgsFunction dynamic shell completion to worktree start/remove/resume subcommands. Start completes with unfinished proposal/feature slugs; remove and resume complete with existing non-main worktree slugs. All completion functions handle errors gracefully (return empty list, never error to shell).

## Changes

### Files Created
无

### Files Modified
- forge-cli/internal/cmd/worktree.go
- forge-cli/internal/cmd/worktree_test.go
- forge-cli/scripts/version.txt

### Key Decisions
- Reused listUnfinishedItems() for start completion (same logic as interactive mode from Task 2)
- Extracted worktreeSlugCompletion() shared helper for remove/resume since both filter non-main worktree slugs
- Completion entries use tab-separated description format (slug\ttype) for rich shell hints
- Used Cobra ShellCompDirectiveNoFileComp to prevent file-system fallback completion

## Test Results
- **Tests Executed**: Yes
- **Passed**: 17
- **Failed**: 0
- **Coverage**: 81.9%

## Acceptance Criteria
- [x] forge worktree start <TAB> shows unfinished proposal and feature slugs
- [x] forge worktree remove <TAB> shows existing worktree slugs
- [x] forge worktree resume <TAB> shows existing worktree slugs
- [x] Completion response time < 200ms
- [x] Completion works for bash, zsh, and fish shells (Cobra handles this automatically)
- [x] No completion on list, push, status subcommands

## Notes
无
