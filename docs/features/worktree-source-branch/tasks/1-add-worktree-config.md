---
id: "1"
title: "Add WorktreeConfig to ForgeConfig and update schema"
priority: "P0"
estimated_time: "30m"
dependencies: []
scope: "backend"
breaking: true
type: "feature"
mainSession: false
---

# 1: Add WorktreeConfig to ForgeConfig and update schema

## Description
Add a `WorktreeConfig` struct to `ForgeConfig` in `forge-cli/pkg/profile/config.go` with `source-branch` and `copy-files` fields. Update the JSON schema and example YAML to document the new section. Add config key accessors for `worktree.source-branch` and `worktree.copy-files` so `forge config get` can read them.

## Reference Files
- `docs/proposals/worktree-source-branch/proposal.md` — Source proposal
- `forge-cli/pkg/profile/config.go` — ForgeConfig struct and config key accessors
- `plugins/forge/references/shared/forge-config.schema.json` — JSON schema
- `plugins/forge/references/shared/forge-config.example.yaml` — Example YAML

## Acceptance Criteria
- [ ] `WorktreeConfig` struct exists with `SourceBranch string` (`yaml:"source-branch"`) and `CopyFiles []string` (`yaml:"copy-files"`)
- [ ] `ForgeConfig` has `Worktree *WorktreeConfig` field (`yaml:"worktree,omitempty"`)
- [ ] `forge config get worktree.source-branch` returns the configured value
- [ ] `forge config get worktree.copy-files` returns the configured list
- [ ] `forge-config.schema.json` has a `worktree` property with `source-branch` (string) and `copy-files` (array of strings)
- [ ] `forge-config.example.yaml` includes a commented-out `worktree` section
- [ ] Existing config without `worktree` section still loads correctly (backward compatible)
- [ ] `additionalProperties: false` preserved at root and worktree level

## Hard Rules
- Do NOT add worktree to `required` in schema — must be optional
- Follow existing config accessor pattern in `configKeyAccessors` map

## Implementation Notes
- The `Worktree` field should be a pointer (`*WorktreeConfig`) so nil = not configured, matching the `Auto` field pattern
- Schema must enforce `additionalProperties: false` on the worktree object for forward compatibility
