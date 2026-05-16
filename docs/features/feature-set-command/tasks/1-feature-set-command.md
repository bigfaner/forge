---
id: "1"
title: "Add forge feature set subcommand"
priority: "P1"
estimated_time: "1h"
dependencies: []
scope: "backend"
breaking: false
type: "implementation"
mainSession: false
---

# 1: Add forge feature set subcommand

## Description

Add an explicit `forge feature set <slug>` subcommand that writes the feature slug to `.forge/state.json` via `EnsureForgeState()` and ensures the feature directory structure exists via `EnsureFeatureDir()`. This provides an explicit override for feature resolution, complementing the existing implicit resolution from git context.

Currently `forge feature <slug>` (positional arg) only calls `SetFeature()` → `EnsureFeatureDir()`, which creates directories but does NOT persist the selection to state.json. The new `set` subcommand does both.

## Reference Files
- `docs/proposals/feature-set-command/proposal.md` — Source proposal
- `forge-cli/internal/cmd/feature.go` — Feature command definitions (add new subcommand here)
- `forge-cli/pkg/feature/forge_state.go` — `EnsureForgeState()` writes state.json with `allCompleted=false`
- `forge-cli/pkg/feature/feature.go` — `EnsureFeatureDir()` creates feature directory structure

## Acceptance Criteria
- [ ] `forge feature set <slug>` creates the feature directory structure under `docs/features/<slug>/`
- [ ] `forge feature set <slug>` writes `.forge/state.json` with `feature=<slug>` and `allCompleted=false`
- [ ] `forge feature set <slug>` returns an error when slug is empty
- [ ] `forge feature set <slug>` prints the feature slug to stdout on success
- [ ] Existing `forge feature <slug>` positional arg behavior unchanged (backward compatible)

## Hard Rules
- Reuse existing `EnsureForgeState()` and `EnsureFeatureDir()` — no new state-writing functions
- Register `featureSetCmd` in `init()` alongside existing subcommands
- Validate slug is non-empty before any filesystem operations

## Implementation Notes
- The existing `runFeature()` handles positional arg as "set" — do NOT modify it. The new `set` subcommand is a separate Cobra command.
- `EnsureForgeState()` already handles directory creation (`.forge/`) — no need to call `EnsureForgeDir()` separately.
- Follow the pattern of `featureListCmd` / `featureStatusCmd` for subcommand registration.
