---
created: "2026-06-06"
tags: [completion-guard, agent-feedback, architecture]
---

# Rejected Tasks Silently Block Feature Completion

## Problem

All meaningful coding tasks completed and committed, but `forge feature complete --if-done` (Stop hook) never pushed. The feature merged without a completion commit, and manifest status stayed at `tasks` indefinitely.

## Root Cause

1. **Symptom**: No auto-push after all coding tasks completed. No completion commit generated.
2. **Direct cause**: `checkFeatureCompletion()` guard requires every task in `index.json` to be `completed` or `skipped`. Two tasks had status `rejected` — the guard silently returned `nil`, the completion pipeline never ran, and the Stop hook exited 0 with no output.
3. **Why silent**: When the guard fails, it returns `nil` with no diagnostics. The agent (dispatcher) sees a clean exit and assumes everything is fine, with no clue that tasks were left in a terminal-but-unrecognized state.

The lesson originally misattributed this to "stale manifest blocking the guard" — but manifest.md is a **write target** of `completeFeature()`, never a read target of the guard. The manifest stayed stale *because* the guard never passed, not the other way around.

## Solution

When the guard finds incomplete tasks, it should print a summary of which tasks are blocking completion and their current statuses. This gives the agent (dispatcher / Stop hook) actionable information to investigate and recover:

- **Rejected but out-of-scope** → agent can `forge task transition <id> skipped --reason "..."`
- **Rejected due to fix not needed** → agent can `forge task transition <id> skipped --reason "..."`
- **Unexpected rejection** → agent can `forge task reopen <id>` and retry

Implementation: in `checkFeatureCompletion()`, when a non-completed/non-skipped task is found, collect the full list and log it before returning `nil`.

```go
// Pseudocode
var blocking []string
for _, t := range index.TasksMap() {
    if t.Status != types.StatusCompleted && t.Status != types.StatusSkipped {
        blocking = append(blocking, fmt.Sprintf("  %s (%s): %s", t.ID, t.Status, t.Title))
    }
}
if len(blocking) > 0 {
    fmt.Fprintf(os.Stderr, "feature not complete — %d task(s) not done:\n%s\n", len(blocking), strings.Join(blocking, "\n"))
    return nil
}
```

## Reusable Pattern

A guard condition that **silently exits 0 on failure** is indistinguishable from success. When a guard rejects, it must produce output describing *what* blocked and *what statuses* were found — otherwise autonomous agents have no signal to act on.

This applies broadly: any hook or CLI command that bails out early should explain why, so downstream consumers (agents, CI pipelines, humans) can diagnose and recover.

## Related Files

- `forge-cli/internal/cmd/feature/feature_complete.go` — guard logic in `checkFeatureCompletion()`
- `docs/features/<slug>/tasks/index.json` — task statuses (source of truth for the guard)
