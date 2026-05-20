---
status: "completed"
started: "2026-05-20 17:03"
completed: "2026-05-20 17:16"
time_spent: "~13m"
---

# Task Record: 1 Refactor start: branch-first creation + --no-launch flag

## Summary
Refactored forge worktree start: (1) Branch-first creation using `git branch <slug> <source>` then `git worktree add <dir> <slug>` instead of atomic `git worktree add -b <slug>`. Cleanup on worktree add failure removes the created branch. (2) Added --no-launch flag to skip claude launch after worktree creation, printing worktree path and exiting. Claude pre-flight check also skipped when --no-launch is set.

## Changes

### Files Created
无

### Files Modified
- forge-cli/internal/cmd/worktree.go
- forge-cli/internal/cmd/worktree_test.go

### Key Decisions
- Branch-first creation uses `git branch <slug> <source>` instead of `git checkout -b` because we don't want to switch the current working tree's HEAD -- only create the branch ref
- Remote branch resolution (Layer 2) also uses branch-first: `git branch <slug> origin/<slug>` then `git worktree add <dir> <slug>`, replacing the old `git worktree add -b <slug> <dir> origin/<slug>`
- Cleanup on worktree add failure: if `git worktree add` fails after branch creation, the branch is deleted via `git branch -D <slug>` to avoid orphan branches
- --no-launch skips both the claude pre-flight check (lookPath) and the claude launch, making it suitable for CI/batch use
- --no-launch prints 'worktree created at <path>' to stdout for scripting consumption

## Test Results
- **Tests Executed**: Yes
- **Passed**: 606
- **Failed**: 0
- **Coverage**: 81.8%

## Acceptance Criteria
- [x] Branch creation uses git checkout -b before git worktree add, decoupling the two steps
- [x] forge worktree start <slug> --no-launch creates the worktree and exits without launching claude
- [x] Without --no-launch, behavior is unchanged (launches claude as before)
- [x] Branch resolution logic (local existing / remote existing / new from source) still works correctly
- [x] Copy-files still applied after worktree creation
- [x] All existing tests pass; new behavior has unit tests

## Notes
7 new tests added: NoLaunchFlagRegistered, NoLaunch_CreatesWorktreeWithoutLaunchingClaude, NoLaunch_WithSourceBranch, WithoutNoLaunch_LaunchesClaude, BranchFirstCreation_NewBranchFromHead, BranchFirstCreation_SkipsCheckoutWhenBranchExists, BranchFirstCreation_WithSourceBranch, BranchFirstCreation_CleansUpBranchOnWorktreeFailure. Existing mock-based test (RemoteBranchResolution) updated to match new branch-first flow.
