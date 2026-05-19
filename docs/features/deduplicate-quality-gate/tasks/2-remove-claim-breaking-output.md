---
id: "2"
title: "Remove BREAKING output from task claim"
priority: "P1"
estimated_time: "30min"
dependencies: []
scope: "backend"
breaking: false
type: "coding.cleanup"
mainSession: false
---

# 2: Remove BREAKING Output from Task Claim

## Description

`forge task claim` currently outputs a `BREAKING: true` field when the task has `breaking: true`. This field is only consumed by the dispatchers (`run-tasks.md`, `execute-task.md`) which are being updated in task 5 to no longer use it. The CLI submit gate (task 1) now reads `breaking` directly from the task frontmatter/index, not from claim output.

Remove the BREAKING output from `printTaskDetails()` in `claim.go`.

## Reference Files
- `docs/proposals/deduplicate-quality-gate/proposal.md` — Source proposal (item 3)

## Acceptance Criteria

- [ ] `printTaskDetails()` no longer prints `BREAKING: true`
- [ ] `forge task claim` output no longer contains BREAKING field
- [ ] The `Task.Breaking` field in types.go remains (used by submit.go via index.json)
- [ ] The `TaskState.Breaking` field remains (for state persistence)
- [ ] Existing tests pass; new/updated tests verify BREAKING is absent from output

## Implementation Notes

- `forge-cli/internal/cmd/claim.go` lines 354-356: remove the `if t.Breaking { PrintField("BREAKING", "true") }` block.
- The Breaking field stays in Task and TaskState structs — submit.go reads it from the task loaded via index.json, not from claim output.
- TDD: update existing claim tests to verify BREAKING is absent from output.
