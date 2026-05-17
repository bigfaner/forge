---
status: "completed"
started: "2026-05-17 17:07"
completed: "2026-05-17 17:26"
time_spent: "~19m"
---

# Task Record: 1 Implement worktree command group scaffold + start subcommand

## Summary
Implement forge worktree command group scaffold and start subcommand. Created worktreeCmd parent cobra command (group with no Run) and worktreeStartCmd subcommand. The start subcommand creates a git worktree at ../<slug> with branch <slug> from HEAD, then launches claude --dangerously-skip-permissions in the worktree directory. Handles branch-already-exists (resume context), target-dir-exists (error with hint), claude-not-in-path (pre-flight check), and not-git-repo errors. Uses filepath.Join for all path construction and reuses lookPathFunc/runClaudeFunc patterns from claude.go for testability.

## Changes

### Files Created
- forge-cli/internal/cmd/worktree.go
- forge-cli/internal/cmd/worktree_test.go

### Files Modified
- forge-cli/internal/cmd/root.go
- forge-cli/internal/cmd/root_test.go
- forge-cli/scripts/version.txt

### Key Decisions
- worktreeCmd registered as group parent (6th group) in root.go init(), following existing pattern of task/e2e/forensic/testing/prompt groups
- worktreeStartCmd uses cobra.ExactArgs(1) to enforce single slug argument
- Target directory computed via filepath.Join(projectRoot, '..', slug) then filepath.Abs() for cross-platform safety
- Branch existence checked via git rev-parse --verify to support resume-from-existing-branch flow
- Reuses lookPathFunc and runClaudeFunc from claude.go for testability (hard rule)

## Test Results
- **Tests Executed**: Yes
- **Passed**: 10
- **Failed**: 0
- **Coverage**: 80.0%

## Acceptance Criteria
- [x] forge worktree start <slug> creates a git worktree at ../<slug> with branch <slug> from HEAD
- [x] If branch <slug> already exists, creates worktree from that branch (resume context)
- [x] If worktree directory already exists, errors with hint to use resume
- [x] If target directory name conflicts with existing non-worktree directory, errors with clear message
- [x] Launches claude --dangerously-skip-permissions in the worktree directory after creation
- [x] Detects claude binary availability before attempting launch (like forge claude does)
- [x] forge feature in the created worktree auto-detects the correct feature via GetWorktreeName()

## Notes
Verified GetWorktreeName() auto-detection in TestWorktreeStart_WorktreeNameAutoDetection by creating a real git worktree via the start subcommand and asserting that GetWorktreeName(worktreeDir) returns the slug.
