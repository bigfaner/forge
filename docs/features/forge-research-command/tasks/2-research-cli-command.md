---
id: "2"
title: "Add forge research CLI command"
priority: "P0"
estimated_time: "30min-1h"
dependencies: ["1"]
scope: "backend"
breaking: false
type: "coding.feature"
mainSession: false
---

# 2: Add forge research CLI command

## Description

Create `cmd/research.go` with list and detail modes following the established pattern from `cmd/proposal.go` and `cmd/lesson.go`. Register the command in `cmd/root.go`. Bump version in `scripts/version.txt` (minor: new command).

## Reference Files
- `docs/proposals/forge-research-command/proposal.md` — Source proposal (command behavior examples)
- `forge-cli/internal/cmd/proposal.go` — Pattern reference for list/detail command
- `forge-cli/internal/cmd/lesson.go` — Pattern reference for list/detail command
- `forge-cli/internal/cmd/root.go` — Command registration point
- `forge-cli/scripts/version.txt` — Current version: 5.2.1 → bump to 5.3.0

## Acceptance Criteria

- [ ] `forge research` (no args) lists all reports in table format: SLUG, CREATED, TOPIC, MODE
- [ ] `forge research <slug>` shows detail view with: SLUG, TOPIC, CREATED, MODE, DIMENSIONS, FILE
- [ ] "no research found" message when no reports exist (matching proposal/lesson pattern)
- [ ] Dynamic slug column width via `CalcSlugColWidth` (matching proposal/lesson pattern)
- [ ] Command registered in `root.go` init() alongside proposalCmd and lessonCmd
- [ ] Version bumped from 5.2.1 to 5.3.0 in `scripts/version.txt`
- [ ] Uses `pkg/project.FindProjectRoot()` for root resolution
- [ ] Uses `PrintBlockStart/End` and `PrintField` helpers for consistent output formatting
- [ ] Proper AIError wrapping for discovery and not-found errors

## Hard Rules

- Use `cobra.MaximumNArgs(1)` for argument validation (matching proposal/lesson)
- Use `RunE` for error propagation
- Follow dependency direction: `cmd → pkg` (import `pkg/research`, not the reverse)

## Implementation Notes

- List mode output should match the proposal's ASCII art format
- Detail mode: DIMENSIONS should be comma-separated string, FILE should show relative path
- Use `mapReportsToSlugLens` helper for column width calculation (follow naming convention from proposal/lesson)
