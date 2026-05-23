---
id: "1"
title: "Show unpushed commit count in worktree status"
priority: "P1"
estimated_time: "1h"
dependencies: []
scope: "backend"
breaking: false
type: "coding.feature"
mainSession: false
---

# 1: Show unpushed commit count in worktree status

## Description

Add an `UNPUSHED` field to the `forge worktree status` output so users can see at a glance which worktrees have unpushed commits. This addresses the problem that users currently have no visibility into push state, leading to accidental data loss when removing worktrees with `--hard`.

## Reference Files
- `docs/proposals/worktree-unpushed-detection/proposal.md` — Source proposal
- `forge-cli/internal/cmd/worktree/worktree.go` — `printWorktreeStatus` (line 464)
- `forge-cli/pkg/git/` — Git helper package

## Acceptance Criteria
- [ ] `forge worktree status` displays an `UNPUSHED` field showing commit count (e.g., `UNPUSHED: 3 commits`)
- [ ] Branches without upstream tracking display `UNPUSHED: no remote`
- [ ] Fully pushed branches display `UNPUSHED: (none)`
- [ ] Detection uses `git rev-list --count @{u}..HEAD` inside the worktree directory

## Hard Rules
- Use `git rev-list --count @{u}..HEAD` for detection — exit code != 0 means no upstream tracking
- Keep the field after `UNCOMMITTED` in the output block to maintain visual consistency

## Implementation Notes
- Add the unpushed check after the existing uncommitted files check in `printWorktreeStatus`
- Extract a reusable helper (e.g., `countUnpushedCommits(worktreePath string) (int, error)`) so the remove command can reuse it in task 2
- No upstream tracking → return a sentinel error; caller formats as "no remote"
