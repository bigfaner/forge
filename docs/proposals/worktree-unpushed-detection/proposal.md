---
created: "2026-05-23"
status: Draft
intent: "Add unpushed commit detection to worktree status/remove to prevent accidental data loss"
---

# Proposal: Worktree Unpushed Commit Detection

## Problem

`forge worktree remove` can delete a worktree whose branch contains commits that have never been pushed to remote. Combined with `--hard` flag (which also deletes the local branch), this can cause **irreversible data loss**. The current safety net only detects uncommitted changes (reactively via git's error), but silently drops unpushed commits.

`forge worktree status` also lacks any push tracking information — users cannot see at a glance which worktrees have unpushed work.

## Solution

Two changes:

### 1. `forge worktree status` — Show unpushed commit count

Add an `UNPUSHED` field to the status output:

```
---
WORKTREE:  auth-login
BRANCH:    auth-login
COMMIT:    abc1234 add login form
UNCOMMITTED: (none)
UNPUSHED:  3 commits
---
```

Detection: run `git rev-list --count @{u}..HEAD` inside the worktree. If the branch has no upstream tracking (`@{u}` fails), treat as "no remote tracking" and display accordingly (e.g., `no remote`).

### 2. `forge worktree remove` — Block on unpushed commits

Before executing `git worktree remove`, proactively check for unpushed commits. If found:
- Print: `"error: branch has N unpushed commit(s) -- push first, or use --force to discard"`
- Return error, do not remove

The existing `--force` flag overrides this check (consistent with how it already overrides the dirty-worktree check).

Detection logic:
- Run `git rev-list --count @{u}..HEAD` inside the worktree directory
- If exit code != 0 (no upstream tracking): skip check (branch was never pushed, no baseline to compare against)
- If count > 0: block removal

## Alternatives

| Alternative | Trade-off |
|-------------|-----------|
| **Do nothing** — rely on user discipline | Unpushed commits silently lost on `--hard` remove |
| **Warn only** (don't block) | Users learn to ignore warnings; defeats the purpose |

## Scope

### In Scope
- Add `unpushedCount` to worktree status display (`printWorktreeStatus`)
- Add proactive unpushed check in `runWorktreeRemove` before removal
- `--force` flag overrides the unpushed check
- Unit tests for both paths

### Out of Scope
- Checking for open PRs on the branch
- Checking if branch has been merged into main
- New flags or subcommands
- Changes to `push` or `start` commands

## Risks

| Risk | Severity | Mitigation |
|------|----------|------------|
| False positive: branch with no upstream is treated as "unpushed" | Medium | Skip check when upstream tracking is absent; display "no remote" in status |
| Performance: additional git command per worktree in status | Low | `git rev-list --count` is fast (local-only operation) |
| Breaking existing workflows that rely on `remove` without push | Medium | `--force` provides escape hatch |

## Success Criteria

- [ ] `forge worktree status` displays UNPUSHED count for each worktree
- [ ] `forge worktree status` displays "no remote" for branches without upstream tracking
- [ ] `forge worktree remove` blocks when unpushed commits exist
- [ ] `forge worktree remove --force` overrides the unpushed check
- [ ] Branches without upstream tracking do not block removal
- [ ] Unit tests cover: unpushed commits found, no upstream, clean state
- [ ] Version bump in `scripts/version.txt` (minor: behavior change with escape hatch)
