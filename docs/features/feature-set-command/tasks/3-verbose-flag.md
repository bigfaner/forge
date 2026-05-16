---
id: "3"
title: "Add verbose flag to forge feature command"
priority: "P2"
estimated_time: "1h"
dependencies: ["2"]
scope: "backend"
breaking: false
type: "implementation"
mainSession: false
---

# 3: Add verbose flag to forge feature command

## Description

Add a `-v` flag to the `forge feature` command (bare invocation, no subcommand) that displays the feature name along with its resolution source. This helps users understand how the current feature was determined — whether from explicit `set`, git context, or directory scanning.

## Reference Files
- `docs/proposals/feature-set-command/proposal.md` — Source proposal
- `forge-cli/internal/cmd/feature.go` — `runFeature()` function and `featureCmd` definition
- `forge-cli/pkg/feature/feature.go` — `GetCurrentFeatureWithSource()` (added in Task 2)

## Acceptance Criteria
- [ ] `forge feature -v` shows `FEATURE: my-feature (from: state.json)` when resolved from state.json
- [ ] `forge feature -v` shows `FEATURE: my-feature (from: worktree)` when resolved from git worktree
- [ ] `forge feature -v` shows `FEATURE: my-feature (from: branch)` when resolved from git branch
- [ ] `forge feature -v` shows `FEATURE: my-feature (from: features-dir)` when resolved from directory scanning
- [ ] `forge feature -v` shows `FEATURE: (none)` when no feature is set
- [ ] `forge feature` (without `-v`) behavior unchanged — shows only the slug
- [ ] `-v` flag only applies to bare `forge feature` — does not affect `set`, `list`, or `status` subcommands

## Hard Rules
- Register `-v` as a local flag on `featureCmd`, not a persistent flag — avoid leaking to subcommands
- Use `GetCurrentFeatureWithSource()` from Task 2; do NOT duplicate resolution logic in the cmd layer

## Implementation Notes
- Use `featureCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "show resolution source")`
- The flag only affects the `len(args) == 0` branch in `runFeature()` — when displaying the current feature.
- When no feature is found, `-v` still shows `(none)` (same as non-verbose, no additional info to show).
