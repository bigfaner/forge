---
created: "2026-06-07"
tags: [architecture, error-handling]
---

# Dual State Stores Require Cleanup on All Mutation Paths

## Problem

`forge task transition` updated `index.json` but did not clear `process/state.json`. When a claimed task was transitioned away from `in_progress`, subsequent `forge task claim` failed with "Task data integrity issues detected — Task has unexpected status: pending".

## Root Cause

Forge stores claimed task state in two files:
- `tasks/index.json` — task status field
- `tasks/process/state.json` — claimed task snapshot with `startedTime`

The `transition` command only mutated `index.json`. When `claim` ran `CheckExistingTaskState()`, it found `state.json` pointing to a task whose index status was now `pending` (not in the handled cases: `in_progress`, `completed`, `blocked`, `suspended`, `rejected`), fell through to `default`, and returned a data integrity error.

Structural cause: new mutation paths added to the system (transition, reopen) did not follow the same cleanup protocol as existing ones (submit, cleanup hook).

## Solution

Both `doTransition` and `doReopen` now check whether `process/state.json` references the task being mutated, and delete it if so. This follows the same pattern used by `cleanup.go` and `CheckExistingTaskState()`.

The `reopen` fix is particularly important for `skipped` tasks: `CheckExistingTaskState` handles `rejected` via auto-delete but does **not** handle `skipped` — it falls through to `default` and returns a data integrity error.

## Reusable Pattern

When a system maintains **multiple state stores** (e.g., index + process state), every mutation path must update **all** stores consistently. Audit checklist for new commands:

1. Identify all state files that reference the entity being mutated
2. For each mutation path (claim, submit, transition, reopen, cleanup), verify that stale references are cleaned up
3. Write a test that exercises the mutation path and then calls the dependent command (e.g., transition → claim) to catch orphaned state

## Related Files

- `forge-cli/internal/cmd/task/transition.go` — transition fix
- `forge-cli/internal/cmd/task/reopen.go` — reopen fix
- `forge-cli/pkg/task/state.go` — `CheckExistingTaskState()`, `DeleteState()`
- `forge-cli/internal/cmd/cleanup.go` — existing cleanup pattern
