---
status: "completed"
started: "2026-05-23 09:56"
completed: "2026-05-23 10:02"
time_spent: "~6m"
---

# Task Record: 1 Show unpushed commit count in worktree status

## Summary
Added CountUnpushedCommits helper to pkg/git and integrated UNPUSHED field into worktree status output. The helper uses git rev-list --count @{u}..HEAD to detect unpushed commits, returning ErrNoUpstream when no upstream tracking is configured. The status command now displays 'UNPUSHED: N commits', 'UNPUSHED: no remote', or 'UNPUSHED: (none)' accordingly.

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/git/git.go
- forge-cli/pkg/git/git_test.go
- forge-cli/internal/cmd/worktree/worktree.go
- forge-cli/internal/cmd/worktree/worktree_test.go

### Key Decisions
- Used errors.New sentinel ErrNoUpstream rather than a custom error type, matching existing error patterns in the codebase
- Added countUnpushedCommitsFunc as overridable var in worktree.go, consistent with existing testability pattern (gitRunFunc, listWorktreesFunc, etc.)
- Extracted formatUnpushed helper for clean separation of formatting logic

## Test Results
- **Tests Executed**: Yes
- **Passed**: 10
- **Failed**: 0
- **Coverage**: 91.9%

## Acceptance Criteria
- [x] CountUnpushedCommits(dir string) (int, error) added to pkg/git/ package
- [x] Returns (0, ErrNoUpstream) when branch has no upstream tracking
- [x] Returns (count, nil) where count is the number of unpushed commits
- [x] forge worktree status displays UNPUSHED field showing commit count
- [x] Branches without upstream tracking display UNPUSHED: no remote
- [x] Fully pushed branches display UNPUSHED: (none)
- [x] Detection uses git rev-list --count @{u}..HEAD inside the worktree directory

## Notes
10 new tests added: 6 unit tests in pkg/git (ErrNoUpstream sentinel, no-upstream scenarios, with-upstream zero/some unpushed, non-git-dir) and 4 integration tests in worktree (no-remote display, commits display, none display, field ordering after UNCOMMITTED). pkg/git coverage: 91.9%, internal/cmd/worktree coverage: 86.1%.
