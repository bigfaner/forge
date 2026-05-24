---
id: "1"
title: "Enhance forge task list with slug parameter and worktree-aware reading"
priority: "P0"
estimated_time: "1-2h"
dependencies: []
scope: "backend"
breaking: false
type: "coding.enhancement"
mainSession: false
---

# 1: Enhance forge task list with slug parameter and worktree-aware reading

## Description

`forge task list` currently relies entirely on auto-detection (`.forge/state.json` → worktree name → branch name → features dir scan) to determine which feature to display. When a feature's tasks are being worked on in a worktree, the main repo's `index.json` copy may be stale. Users must `cd` into the worktree to see latest task status.

Enhance the command with an optional positional slug parameter and `--local` flag so users can query any feature's tasks from any working directory, with automatic worktree-aware index reading.

## Reference Files

- `docs/proposals/task-list-slug-worktree/proposal.md` — Source proposal
- `forge-cli/internal/cmd/task/list.go` — Target file (current command implementation)
- `forge-cli/pkg/git/worktree_list.go` — `ListWorktrees()`, `WorktreeEntry` for worktree detection
- `forge-cli/pkg/git/git.go` — `GetWorktreeName()` for current worktree detection
- `forge-cli/internal/cmd/worktree/cmd_push.go` — `resolveWorktreeDir()` pattern for slug-to-path resolution
- `forge-cli/pkg/feature/feature.go` — `RequireFeature()`, `GetCurrentFeature()`
- `forge-cli/pkg/feature/paths.go` — `GetFeatureIndexFile()` for index path construction
- `forge-cli/pkg/project/root.go` — `FindProjectRoot()`
- `docs/conventions/forge-cli-reference.md` — CLI command conventions

## Acceptance Criteria

- `forge task list my-feature` shows tasks for the specified feature, reading from worktree if one exists for that slug
- `forge task list my-feature --local` reads from main repo's `index.json` regardless of worktree existence
- `forge task list` (no args) behaves exactly as before — existing auto-detection logic unchanged
- Clear error message when slug doesn't match any `docs/features/<slug>/` directory
- "no tasks found" message when feature exists but has no `index.json`
- Unit tests cover: slug resolution, worktree detection, `--local` override, backward compatibility (no-arg path)
- Version bumped in `scripts/version.txt` per semver (minor: new user-facing capability)
- `docs/conventions/forge-cli-reference.md` updated to reflect new usage

## Hard Rules

- Must not break existing `forge task list` behavior (no-arg path is backward compatible)
- Worktree detection must use existing `git.ListWorktrees()` or `resolveWorktreeDir()` pattern — do not invent new path resolution
- Follow CLI table rendering conventions from `docs/conventions/forge-cli-reference.md`

## Implementation Notes

- `resolveWorktreeDir()` in `internal/cmd/worktree/cmd_push.go` constructs `.forge/worktrees/<slug>` path — reuse this pattern or extract to shared utility
- Worktree detection priority: check if `.forge/worktrees/<slug>` directory exists (fast path), fallback to `git.ListWorktrees()` if needed
- `cobra.NoArgs` must change to `cobra.MaximumNArgs(1)` to accept optional positional slug
- Add `--local` bool flag via `listCmd.Flags().Bool("local", false, "...")`
- When slug is provided: bypass `feature.RequireFeature()`, directly construct index path from slug
- Risk: worktree path resolution differs by OS — mitigated by using existing utilities that handle OS differences
