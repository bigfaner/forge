---
id: "2"
title: "Block worktree remove on unpushed commits"
priority: "P1"
estimated_time: "1h"
dependencies: ["1"]
scope: "backend"
breaking: false
type: "coding.feature"
mainSession: false
---

# 2: Block worktree remove on unpushed commits

## Description

Add a proactive unpushed commit check before `git worktree remove` executes. When unpushed commits exist, the command prints an error and aborts removal. The existing `--force` flag overrides this check, consistent with how it already overrides the dirty-worktree check.

## Reference Files
- `docs/proposals/worktree-unpushed-detection/proposal.md` — Source proposal
- `forge-cli/pkg/git/` — `CountUnpushedCommits` helper (created in task 1)
- `forge-cli/internal/cmd/worktree/worktree.go` — `runWorktreeRemove` (line 174)
- `forge-cli/scripts/version.txt` — Version file to bump

## Acceptance Criteria
- [ ] `forge worktree remove <slug>` blocks when unpushed commits exist on the worktree's branch
- [ ] Error message: `"error: branch has N unpushed commit(s) — push first, or use --force to discard"`
- [ ] `forge worktree remove --force <slug>` overrides the unpushed check and proceeds with removal
- [ ] Branches without upstream tracking do not block removal (skip check when `ErrNoUpstream`)
- [ ] Version bumped in `forge-cli/scripts/version.txt` (minor: new behavior with escape hatch)

## Hard Rules
- Check must happen before `git worktree remove` is invoked (not after)
- Use `git.CountUnpushedCommits` from `pkg/git/` — do not inline detection logic
- `--force` flag already exists — no new CLI flags

## Implementation Notes
- Insert the unpushed check after resolving the worktree path and branch name, but before building the `git worktree remove` args
- Call `git.CountUnpushedCommits(targetDir)` and handle:
  - `ErrNoUpstream` → skip check, proceed with removal
  - `count > 0 && !force` → print error and return
  - `count > 0 && force` → proceed (force overrides)
- The `--force` flag is already read at line 182 — check it to decide whether to skip or block
- Version bump: minor (5.2.1 → 5.3.0) since this adds user-facing safety behavior
