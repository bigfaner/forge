---
id: "1"
title: "Add forge task list subcommand"
priority: "P1"
estimated_time: "1-2h"
dependencies: []
scope: "backend"
breaking: false
type: "coding.feature"
mainSession: false
---

# 1: Add forge task list subcommand

## Description

Add `forge task list` subcommand that displays all tasks for the current feature in a table format. This fills a gap where users must manually open `index.json` to see task IDs, titles, and statuses — critical for using `forge task claim` or `forge task transition` which require task IDs.

## Reference Files
- `docs/proposals/list-tasks/proposal.md` — Source proposal
- `forge-cli/internal/cmd/task/status.go` — Existing subcommand pattern
- `forge-cli/internal/cmd/task/list_types.go` — Simpler subcommand example
- `forge-cli/internal/cmd/task/register.go` — Command registration
- `forge-cli/internal/cmd/base/output.go` — Output formatting helpers
- `forge-cli/pkg/task/` — Task data types (LoadIndex, TaskIndex, Task)
- `forge-cli/pkg/feature/` — Feature resolution (RequireFeature, GetFeatureIndexFile)
- `scripts/version.txt` — Version bump

## Acceptance Criteria

1. `forge task list` resolves current feature via `feature.RequireFeature()`, loads `index.json`, and displays all tasks in a table
2. Table columns: ID, Type, Title (truncated if too long), Status
3. Header shows total count and feature slug (e.g. `8 found  (feature: feature-set-command)`)
4. Tasks sorted by ID in natural order: numeric IDs first (1, 2, 3...), then test/gate IDs (T-1, T-2...)
5. Features with no `index.json` or empty task list print a clear "no tasks found" message
6. Unit tests cover: normal list output, empty feature, sorted output order
7. Version bumped in `scripts/version.txt` (minor: new command)

## Hard Rules

- Follow existing cobra subcommand pattern: package-level `var xxxCmd` + `runXxx` function
- Register `listCmd` in `register.go` alongside existing subcommands
- No extra flags for initial version (no `--status` filter)

## Implementation Notes

- Natural sort: extract numeric prefix from ID for primary sort, fall back to string comparison for secondary. Test/gate IDs (T-1, T-2) sort after business IDs.
- Title truncation: cap at ~30 chars with ellipsis to keep table scannable for large lists.
- Use `base.PrintBlockStart()` / `PrintBlockEnd()` for consistent output framing, or use the table format from the proposal mockup.
- The `query.go` subcommand has flag examples if needed later.
