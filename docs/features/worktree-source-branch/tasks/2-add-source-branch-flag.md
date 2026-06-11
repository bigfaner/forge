---
id: "2"
title: "Add --source-branch flag to worktree start command"
priority: "P0"
estimated_time: "45m"
dependencies: ["1"]
scope: "backend"
breaking: false
type: "feature"
mainSession: false
---

# 2: Add --source-branch flag to worktree start command

## Description
Add a `--source-branch` / `-b` CLI flag to `forge worktree start`. Implement source resolution with priority: flag > config > HEAD. Update the worktree creation logic to use the resolved source branch as the start point for `git worktree add`.

## Reference Files
- `docs/proposals/worktree-source-branch/proposal.md` — Source proposal
- `forge-cli/internal/cmd/worktree.go` — worktree start command and `runWorktreeStart`
- `forge-cli/pkg/profile/config.go` — ForgeConfig with new WorktreeConfig

## Acceptance Criteria
- [ ] `forge worktree start <slug> --source-branch develop` creates worktree from `develop` branch
- [ ] `forge worktree start <slug> -b v3.0.0` creates worktree from `v3.0.0` branch
- [ ] Without flag or config, behavior is identical to current version (creates from HEAD)
- [ ] `worktree.source-branch` in config.yaml sets default source branch
- [ ] Flag overrides config default
- [ ] Clear error message when specified branch does not exist locally or remotely
- [ ] CLI help text updated to show `--source-branch` flag and usage
- [ ] `Long` description of `worktreeStartCmd` updated to reflect source-branch support

## Hard Rules
- Pre-validate branch exists (via `git rev-parse --verify`) before worktree creation
- Source branch applies to the "new branch" path only (when `branchExists == false`). For existing branches, the worktree is created from the existing branch as today.

## Implementation Notes
- The `git worktree add -b <slug> <path> <start-point>` command already supports specifying a start point — just pass the resolved source branch as the last argument
- Read config via `profile.ReadConfig(projectRoot)` and check `cfg.Worktree.SourceBranch`
- Flag registration: `worktreeStartCmd.Flags().StringP("source-branch", "b", "", ...)`
