---
id: "3"
title: "Remove --hard: delete worktree + branch + prune"
priority: "P1"
estimated_time: "1.5h"
dependencies: []
scope: "backend"
breaking: false
type: "coding.feature"
mainSession: false
---

# 3: Remove --hard: delete worktree + branch + prune

## Description

Add `--hard` flag to `forge worktree remove` that performs a thorough cleanup: removes the worktree, deletes the local branch, and runs `git worktree prune` to clean up stale administrative files. Currently remove only deletes the worktree directory and leaves the branch intact.

## Reference Files
- `docs/proposals/worktree-experience/proposal.md` — Source proposal
- `forge-cli/internal/cmd/worktree.go` — Remove command implementation (runWorktreeRemove)
- `forge-cli/internal/cmd/worktree_test.go` — Existing tests

## Acceptance Criteria
- [ ] `forge worktree remove <slug> --hard` performs three steps: worktree remove → branch delete → prune
- [ ] Each step reports its status (success or skip with reason)
- [ ] If worktree removal fails, branch deletion and prune are skipped
- [ ] Before deleting branch, check for uncommitted changes in the worktree (fail with message if found)
- [ ] Before deleting branch, check if branch is merged into source-branch (warn if not, require `--force` to proceed)
- [ ] `forge worktree remove <slug>` (without --hard) behavior is unchanged

## Hard Rules
- Do NOT delete remote branches — only local cleanup
- Branch deletion must refuse if uncommitted changes exist, unless `--force` is also passed
- `--hard` without `--force` should warn but allow deletion of unmerged branches (only block on uncommitted changes)

## Implementation Notes
- Current remove flow: resolve slug → `git worktree remove <dir>` → print message
- Hard flow extends: after worktree remove → `git branch -d <branch>` (or `-D` with --force) → `git worktree prune`
- The existing `runWorktreeRemove` looks up the branch name before removal (line ~68). Reuse that for branch deletion.
- Key risk: `--hard` can cause data loss if the user has unpushed commits. The uncommitted-changes check is critical.
