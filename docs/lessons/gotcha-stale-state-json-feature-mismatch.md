---
created: "2026-05-16"
tags: [local-dev-deployment, testing]
---

# Stale `.forge/state.json` Causes Task Claim to Use Wrong Feature

## Problem

`forge task claim` returns "No pending tasks available" even though `index.json` has 10 pending tasks. The CLI was reading a stale `.forge/state.json` that pointed to a previous feature (`feature-set-command`) whose tasks were all completed.

## Root Cause

Causal chain (3 levels):

1. **Symptom**: `forge task claim` → "No pending tasks available" despite 10 pending tasks in `docs/features/e2e-test-quality-cleanup/tasks/index.json`.
2. **Direct cause**: `forge task claim` reads `.forge/state.json` to determine the active feature. The file contained `{"feature": "feature-set-command"}` — a previous feature whose tasks were already completed.
3. **Root cause**: The run-tasks workflow (`/run-tasks`) and its dispatcher do not automatically update `state.json` when starting a new feature. The previous feature's state persists across sessions.
4. **Trigger condition**: Starting a new feature's task execution without explicitly setting the active feature first.

## Solution

Before claiming tasks for a new feature, always set the active feature:

```bash
forge feature set <slug>
```

This writes the correct feature slug to `.forge/state.json`.

## Reusable Pattern

**When starting task execution for any feature, the first action must be `forge feature set <slug>`.** The `/run-tasks` dispatcher and `forge task claim` rely on `.forge/state.json` to determine which feature's tasks to claim. If the state is stale from a previous feature, claims will fail or target the wrong feature.

This applies to:
- Starting `/run-tasks` after `/quick-tasks` generates tasks
- Resuming work on a different feature in a new session
- Any workflow that calls `forge task claim`

## Example

```bash
# WRONG: state.json still points to old feature
forge task claim  # → "No pending tasks available"

# RIGHT: set feature first
forge feature set e2e-test-quality-cleanup
forge task claim  # → "ACTION: CLAIMED" with correct task
```

## Related Files

- `.forge/state.json` — Active feature state
- `docs/features/<slug>/tasks/index.json` — Task index per feature
