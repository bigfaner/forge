---
name: run-tasks-git-status
status: Approved
type: feature
created: 2026-05-23
---

# Run-Tasks Post-Completion Git Status

## Problem

After `run-tasks` completes all tasks, the user has zero visibility into the current git state. The only output is a static message about test tasks. The user must manually run `git status` / `git log` to understand what changed — branch position, uncommitted artifacts, commits made during execution.

## Solution

In the Post-Completion section of `run-tasks.md`, add instructions to display a concise git summary after the existing completion message (test task notification and artifact commit prohibition):

1. **Branch info**: current branch name + ahead/behind relative to main
2. **Working tree changes**: list of modified/untracked files (from `git status --short`)

## Alternatives

| # | Approach | Trade-off |
|---|----------|-----------|
| 1 | **Post-completion git summary** (recommended) | Zero user effort; always shown; ~3 lines of bash in the command file |
| 2 | Do nothing | User runs `git status` manually; negligible cost but easy to forget after long execution |
| 3 | Full `git diff --stat` | More detail than needed; adds noise for simple workflows |

## Scope

### In Scope
- Add git status display to run-tasks Post-Completion section
- Show: branch name, ahead/behind commits, changed/untracked files

### Out of Scope
- Structured summary of per-task outcomes (separate concern)
- Diff statistics (overkill for this use case)
- Changes to hooks or other commands

## Risks

| Risk | Mitigation |
|------|------------|
| Git commands fail in non-git context | Wrap in error handling; skip silently on failure |

## Success Criteria

- [ ] After run-tasks loop ends, user sees current branch + ahead/behind + file changes
- [ ] Graceful degradation if git commands fail
