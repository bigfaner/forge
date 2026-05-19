---
created: "2026-05-18"
tags: [local-dev-deployment, architecture]
---

# Auto-Push Requires Upstream Tracking on Feature Branches

## Problem

`.forge/config.yaml` has `gitPush: true`, but commits on the `eval-adversarial-scorer` worktree branch were never pushed to remote. The config was silently ignored.

## Root Cause

1. **Symptom**: Commits ahead of remote, no push happened after task completion.
2. **Direct cause**: New feature branch has no upstream tracking set (`git branch -vv` shows no `[origin/...]`).
3. **Code cause**: `feature_complete.go:247` `gitPush()` runs bare `git push` — fails on untracked branches with "no upstream branch". Error is logged to stderr but non-blocking (hook protocol exits 0), so the failure is invisible.
4. **Why it happens every time**: Worktree branches created by `EnterWorktree` (or `git checkout -b`) don't set upstream tracking. The first `git push` must include `-u`.

## Solution

Fixed `gitPush()` to use `git push -u origin HEAD` instead of bare `git push`. This sets upstream tracking on first push and works for all subsequent pushes.

```bash
# Before (broken for new branches)
git push

# After (works for both new and existing branches)
git push -u origin HEAD
```

## Reusable Pattern

When implementing a `git push` wrapper, always use `git push -u origin HEAD` to handle both new and existing branches. Never assume upstream tracking is set.

## Related Files

- `forge-cli/internal/cmd/feature_complete.go` — `gitPush()` function
- `.forge/config.yaml` — `auto.gitPush: true` setting
