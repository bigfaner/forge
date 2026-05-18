---
created: "2026-05-18"
tags: [local-dev-deployment, architecture]
---

# Auto-Push Requires Upstream Tracking on Feature Branches

## Problem

`.forge/config.yaml` has `gitPush: true`, but commits on the `eval-adversarial-scorer` worktree branch were never pushed to remote. The config was silently ignored.

## Root Cause

1. **Symptom**: 3 commits ahead of remote, no push happened after task completion.
2. **Direct cause**: The branch `eval-adversarial-scorer` has no upstream tracking set (`git branch -vv` shows no `[origin/...]`).
3. **Root cause**: `gitPush: true` in forge config tells the agent to push after commits, but `git push` without `-u` on an untracked branch fails silently (or is skipped). Worktree branches created by `EnterWorktree` don't automatically set upstream tracking.

## Solution

After creating a worktree branch (or any new feature branch), set upstream tracking immediately:

```bash
git push -u origin <branch-name>
```

Or ensure the task executor / run-tasks flow does this as part of branch setup.

## Reusable Pattern

When `gitPush: true` is set in forge config, verify that the current branch has upstream tracking before assuming pushes will work. If `git branch -vv` shows no `[origin/...]` tracking, push with `-u` first.

## Related Files

- `.forge/config.yaml` — `gitPush: true` setting
- Worktree branch creation flow
