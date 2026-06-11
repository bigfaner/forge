---
id: "2"
title: "Integrate topo sort into forge task list + add --sort flag"
priority: "P0"
estimated_time: "2h"
dependencies: ["1"]
surface-key: ""
surface-type: "cli"
breaking: false
type: "coding.enhancement"
mainSession: false
---

# 2: Integrate topo sort into forge task list + add --sort flag

## Description

Change `forge task list` default sorting from natural ID order to topological order using the algorithm from Task 1. Add `--sort` flag to allow users to restore the old natural ID ordering. Add `[cycle]` and `[missing: <id>]` markers in the table output. Ensure non-TTY environments suppress color output for these markers.

## Reference Files
- `proposal.md#Proposed-Solution` — defines default topo sort, `--sort id` flag, and claim alignment
- `proposal.md#Scope` — In Scope: default topo sort table, `--sort id` flag, cycle/missing markers, non-TTY color disable
- `proposal.md#Key-Risks` — breaking change risk: existing scripts depending on ID order must have `--sort id` migration path
- `forge-cli/internal/cmd/task/list.go` — current implementation using `naturalSortTaskIDs`, the target for integration

## Acceptance Criteria

- [ ] `forge task list` defaults to topological sort ordering
- [ ] `forge task list --sort id` restores natural ID ordering (current behavior)
- [ ] Cycle nodes display `[cycle]` marker in the table row
- [ ] Missing deps display `[missing: <id>]` marker in the table row
- [ ] Non-TTY environments (piped output) do not emit ANSI color codes for markers
- [ ] Empty feature still shows "no tasks found"
- [ ] Existing `list_test.go` tests pass after integration
- [ ] New tests cover: topo sort default, `--sort id` fallback, cycle marker, missing dep marker, pipe mode color suppression

## Hard Rules

- `--sort` flag values: `topo` (default when flag omitted) and `id` — no other values
- Pipeline output (non-TTY) must remain machine-parseable: markers are plain text, no ANSI escape codes
- Column alignment must not break with marker text (account for marker width in dynamic column calculation)

## Implementation Notes

- Current `naturalSortTaskIDs` at list.go:205-233 is the sorting function to replace as default
- The `--sort` flag can use `cobra.Flags().String("sort", "topo", ...)` with validation
- For color: use `base.IsTerminal()` or `golang.org/x/term` to detect TTY; markers in non-TTY are plain `[cycle]`/`[missing: id]`
- The cycle/missing info comes from the `TopologicalSort` return values — inject markers into the ID or a new column
