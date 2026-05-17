---
status: "completed"
started: "2026-05-17 17:27"
completed: "2026-05-17 17:42"
time_spent: "~15m"
---

# Task Record: 2 Implement forge worktree list subcommand

## Summary
Implement forge worktree list subcommand that displays all git worktrees with name, branch, and path. Worktrees matching a feature slug in docs/features/ are marked as forge-managed. Main worktree is distinguished with [main] tag. Uses git worktree list --porcelain for reliable parsing.

## Changes

### Files Created
- forge-cli/pkg/git/worktree_list.go
- forge-cli/pkg/git/worktree_list_test.go

### Files Modified
- forge-cli/internal/cmd/worktree.go
- forge-cli/internal/cmd/worktree_test.go
- forge-cli/internal/cmd/root.go

### Key Decisions
- ParsePorcelainWorktrees is a pure function in pkg/git for testable porcelain output parsing separate from git binary execution
- listWorktreesFunc variable in cmd package follows existing lookPathFunc/runClaudeFunc pattern for testability
- WorktreeEntry.Name() derives name from directory basename per Hard Rules requirement
- listForgeFeatures scans docs/features/ directories to identify forge-managed worktrees

## Test Results
- **Tests Executed**: Yes
- **Passed**: 17
- **Failed**: 0
- **Coverage**: 92.6%

## Acceptance Criteria
- [x] forge worktree list displays all git worktrees (from git worktree list)
- [x] Each entry shows: worktree name, branch, path
- [x] Entries whose worktree name matches a directory in docs/features/ are marked as forge-managed
- [x] Main worktree (current project) is included but distinguished from feature worktrees
- [x] No worktrees prints 'No worktrees found' message

## Notes
ParsePorcelainWorktrees handles bare repos (skips them), detached HEAD (empty branch shown as '(detached)'), and blank-line separated blocks. Uses tabwriter for aligned column output.
