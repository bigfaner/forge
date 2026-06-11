---
id: "1"
title: "Implement three-layer branch resolution for worktree start"
priority: "P0"
estimated_time: "1.5h"
dependencies: []
scope: "backend"
breaking: false
type: "enhancement"
mainSession: false
---

# 1: Implement three-layer branch resolution for worktree start

## Description

Extend `runWorktreeStart` in `forge-cli/internal/cmd/worktree.go` to detect remote-only branches (`origin/<slug>`) in addition to local branches. When a remote branch exists but no local branch does, create the worktree from the remote branch instead of from `source-branch`.

Current code only checks local branches via `git rev-parse --verify <slug>`. Need to add auto-fetch and remote branch detection before falling back to source-branch.

## Reference Files
- `docs/proposals/worktree-remote-branch-reuse/proposal.md` — Source proposal
- `forge-cli/internal/cmd/worktree.go` — Main implementation file (lines 258-292)
- `forge-cli/pkg/git/git.go` — Git utilities (may need `Run` for fetch)

## Acceptance Criteria

- `git fetch origin` runs before branch existence check; fetch failure degrades gracefully (log warning, skip remote check, proceed with local/source-branch)
- After local branch check (existing), add check for `origin/<slug>` via `git rev-parse --verify origin/<slug>` or similar
- When remote branch exists and local does not: run `git worktree add -b <slug> <targetDir> origin/<slug>` to create local tracking branch from remote
- Output message to stdout when using remote branch: "creating worktree from remote branch origin/<slug>"
- `--source-branch` flag is ignored when branch exists (local or remote), consistent with current local-branch behavior
- Existing local-branch reuse behavior is unchanged

## Hard Rules

- Fetch failure must NOT block worktree creation — degrade to current behavior (local check + source-branch)
- Do NOT add new exported functions to `pkg/git/` for this — use existing `git.Run()` directly

## Implementation Notes

- The three-layer check replaces the current two-path branch (line 266-292):
  1. Local `<slug>` exists → `git worktree add <targetDir> <slug>` (unchanged)
  2. Remote `origin/<slug>` exists → `git worktree add -b <slug> <targetDir> origin/<slug>` (new)
  3. Neither → create from source-branch (unchanged)
- Fetch should use `git.Run(projectRoot, "fetch", "origin")` — ignore errors (network down, no remote, etc.)
- Remote check: `git.Run(projectRoot, "rev-parse", "--verify", "remotes/origin/"+slug)` — using `remotes/origin/` prefix avoids ambiguity with tags
