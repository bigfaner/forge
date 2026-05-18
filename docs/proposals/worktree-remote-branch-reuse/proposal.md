---
name: worktree-remote-branch-reuse
status: Draft
date: 2026-05-18
---

# Problem

`forge worktree start <slug>` only checks local branches when deciding whether to reuse an existing branch. If a branch exists on remote (`origin/<slug>`) but not locally, the command creates a new branch from `source-branch` instead of reusing the remote branch. This causes divergence from existing work that was already pushed.

## Evidence

- Current code at `worktree.go:268` uses `git rev-parse --verify <slug>` which only resolves local refs
- Remote-only branches are invisible to this check
- User has experienced this when resuming work on a different machine or after a branch was pushed but the local worktree was removed

## Urgency

Medium — causes silent divergence when it happens. Workaround is manual (`git fetch origin && git checkout <slug>`) but the tool should handle this automatically.

# Solution

Extend the branch resolution logic to three layers:

1. **Local branch** `<slug>` exists → use directly (current behavior, unchanged)
2. **Remote branch** `origin/<slug>` exists → create worktree from remote, establishing local tracking branch
3. **Neither** → create from `source-branch` (current behavior, unchanged)

Auto-fetch from origin before checking, with graceful fallback if network is unavailable.

## Success Criteria

- [ ] `forge worktree start <slug>` detects remote-only `origin/<slug>` and creates worktree from it
- [ ] Auto-fetch runs before branch check; fetch failure degrades gracefully (skip remote check, proceed with local/source-branch)
- [ ] Informative output message when using remote branch: "creating worktree from remote branch origin/<slug>"
- [ ] `--source-branch` flag is ignored when a branch (local or remote) with the slug name already exists, consistent with current local-branch behavior
- [ ] Existing tests pass; new tests cover remote-branch detection and fetch-failure fallback

# Alternatives

| Approach | Pros | Cons |
|----------|------|------|
| **Three-layer resolution** (proposed) | Automatic, consistent with local-branch behavior, minimal code change | Adds network latency from fetch |
| **Do nothing** | Zero risk | Users must manually fetch + checkout before running worktree start |
| **`--reuse-remote` flag** | Explicit opt-in | Extra flag to remember, easy to forget |

# Scope

## In Scope

- Auto-fetch + remote branch detection in `runWorktreeStart`
- Git command for creating worktree from remote branch (`git worktree add -b <slug> <path> origin/<slug>`)
- User-facing output messages for remote branch reuse
- Test cases in `worktree_test.go`

## Out of Scope

- Prefix matching (`feature/<slug>`, etc.)
- Changes to `resume`/`list`/`remove` subcommands
- Any other `forge worktree start` behavior changes

# Risks

| Risk | Mitigation |
|------|-----------|
| Fetch adds latency on slow networks | Graceful degradation: if fetch fails, skip remote check and proceed with local/source-branch |
| Remote branch diverged significantly from expected state | This is the same risk as local branch reuse — user chose the slug, they own the branch state |
| `git worktree add -b <slug> <path> origin/<slug>` fails if local ref somehow exists | Local branch check happens first (step 1), so step 2 only runs when local branch is absent |
