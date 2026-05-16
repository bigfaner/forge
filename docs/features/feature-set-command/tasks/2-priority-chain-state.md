---
id: "2"
title: "Adjust GetCurrentFeature() to read state.json first"
priority: "P1"
estimated_time: "1h"
dependencies: ["1"]
scope: "backend"
breaking: true
type: "implementation"
mainSession: false
---

# 2: Adjust GetCurrentFeature() to read state.json first

## Description

Modify `GetCurrentFeature()` to check `.forge/state.json` as the highest-priority source before falling back to git context. This makes explicit feature selection (via `forge feature set`) take precedence over implicit git-derived resolution.

Additionally, add a `GetCurrentFeatureWithSource()` function that returns the resolution source alongside the feature slug, enabling verbose display in Task 3.

## Reference Files
- `docs/proposals/feature-set-command/proposal.md` — Source proposal
- `forge-cli/pkg/feature/feature.go` — `GetCurrentFeature()` and `getFeatureFromFeaturesDir()`
- `forge-cli/pkg/feature/forge_state.go` — `ReadForgeState()` reads state.json
- `forge-cli/internal/cmd/feature.go` — Caller of `GetCurrentFeature()`

## Acceptance Criteria
- [ ] `GetCurrentFeature()` returns the feature from state.json when it exists and the feature directory is valid
- [ ] When state.json is absent (deleted by quality-gate cleanup), behavior falls back to git context — identical to current behavior
- [ ] `GetCurrentFeatureWithSource()` returns `(slug, source, error)` where source is one of: `"state.json"`, `"worktree"`, `"branch"`, `"features-dir"`
- [ ] `GetCurrentFeature()` calls `GetCurrentFeatureWithSource()` internally, preserving existing API
- [ ] All existing callers of `GetCurrentFeature()` continue to work without changes
- [ ] Existing tests pass without modification (state.json absent = identical behavior)

## Hard Rules
- `GetCurrentFeature()` signature must not change — add `GetCurrentFeatureWithSource()` as a new function
- When state.json feature directory doesn't exist, skip to next priority (don't auto-create from state.json)
- State.json read failure (corrupt file) is silently ignored — fall through to git context

## Implementation Notes
- Priority chain becomes: state.json → git context (worktree → branch) → features-dir scanning (state.json → single feature)
- The proposal names the resolution sources as: `state.json`, `worktree`, `branch` (or `git` for combined). For `git.GetFeatureFromGit()`, we may need to distinguish worktree vs branch — check if the existing API exposes this.
- Risk: `set` + `task claim` both write state.json — but with same feature slug, so no conflict (per proposal risk table).
