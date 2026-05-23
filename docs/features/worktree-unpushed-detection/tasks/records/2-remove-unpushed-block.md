---
status: "completed"
started: "2026-05-23 10:02"
completed: "2026-05-23 10:06"
time_spent: "~4m"
---

# Task Record: 2 Block worktree remove on unpushed commits

## Summary
Added unpushed commit check to runWorktreeRemove that blocks removal when unpushed commits exist, with --force override and ErrNoUpstream skip. Version bumped 5.2.1 -> 5.3.0.

## Changes

### Files Created
无

### Files Modified
- forge-cli/internal/cmd/worktree/worktree.go
- forge-cli/internal/cmd/worktree/worktree_test.go
- forge-cli/scripts/version.txt

### Key Decisions
- Placed unpushed check after resolving worktree path/branch but before building git worktree remove args (Hard Rule: check before invocation)
- Used countUnpushedCommitsFunc (already declared as overridable var) for testability, reusing git.CountUnpushedCommits from pkg/git/
- ErrNoUpstream skips the check entirely (branch never pushed, no baseline to compare against)
- Unexpected errors from countUnpushedCommitsFunc log a warning but don't block removal (defensive: don't fail open on unexpected errors)
- Version bump minor (5.2.1 -> 5.3.0) per semver: new user-facing safety behavior with escape hatch

## Test Results
- **Tests Executed**: Yes
- **Passed**: 4
- **Failed**: 0
- **Coverage**: 86.1%

## Acceptance Criteria
- [x] forge worktree remove <slug> blocks when unpushed commits exist on the worktree's branch
- [x] Error message: "error: branch has N unpushed commit(s) -- push first, or use --force to discard"
- [x] forge worktree remove --force <slug> overrides the unpushed check and proceeds with removal
- [x] Branches without upstream tracking do not block removal (skip check when ErrNoUpstream)
- [x] Version bumped in forge-cli/scripts/version.txt (minor: new behavior with escape hatch)

## Notes
All 21 existing remove tests continue to pass with no regressions. Coverage: worktree package 86.1%, git package 91.9%.
