# PipelineNode struct fields must not be hardcoded in expand variants

## Problem

All auto-generated pipeline tasks had `breaking=true` in `forge task list` output, but most pipeline tasks (doc review, test gen, drift detection) are verification/validation tasks that should not trigger the full unit-test quality gate on submit.

## Root Cause

`PipelineNode` struct had `MainSession` as a configurable per-node bool field, but `Breaking` was not a struct field at all — it was hardcoded to `true` in all 4 `expandNode()` code paths (default, single-surface, multi-surface, per-surface-type). This asymmetry was a design oversight: when one behavioral flag is configurable, its sibling flag should also be configurable rather than hardcoded.

## Solution

Added `Breaking bool` to `PipelineNode`. Go's zero value (`false`) naturally provides the correct default — auto-gen tasks are non-breaking by default. The 4 hardcoded `Breaking: true` sites were replaced with `node.Breaking`. Individual registry nodes that need breaking behavior can opt in explicitly.

## Reusable Pattern

When a struct has multiple boolean behavioral flags, check whether ALL of them are exposed as configurable fields. If one is configurable (like `MainSession`), its siblings (like `Breaking`) should follow the same pattern rather than being hardcoded at expansion sites. This avoids the "configurable vs hardcoded" asymmetry that makes the system misleading in display output.

## References

- `forge-cli/pkg/task/pipeline.go` — PipelineNode struct and expandNode variants
- `forge-cli/pkg/task/pipeline_validate.go` — PipelineRegistry (all 12 nodes unaffected, Go zero value = false)
