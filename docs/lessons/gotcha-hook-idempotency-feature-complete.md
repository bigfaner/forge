---
created: "2026-05-19"
tags: [error-handling, architecture]
---

# Hook must be idempotent: feature complete created 34 duplicate commits

## Problem

The `forge feature complete --if-done` Stop hook fired on every session end and created duplicate `mark feature completed` commits. Over the course of a day, 34 redundant commits accumulated across 5 features (worktree-remote-branch-reuse: 14, eval-adversarial-scorer: 11, knowledge-accumulation-loop: 5, worktree-source-branch: 2, contract-journey-test-model: 1, skill-ecosystem-audit: 1). Each commit only added a blank line to manifest.md frontmatter.

## Root Cause

Three independent failures compounded:

1. **No idempotency in `updateFileStatus`** (feature_complete.go): The function rewrote the file on every call, even when `status` was already set to the target value. The string reconstruction (`"---\n" + strings.Join(lines, "\n") + "\n---" + body`) introduced a whitespace diff each time, so git saw a change even when the logical content was identical.

2. **No completion state check in `checkFeatureCompletion`**: The guard only verified that all tasks in `index.json` were completed/skipped. It never checked whether the feature completion hook had already run successfully. Since task statuses don't revert, the guard always passed.

3. **`state.json` allCompleted field was unreliable**: `WriteForgeState` (called by submit.go) sets `allCompleted: true`, but `EnsureForgeState` (called by add.go, claim.go, quality_gate.go) overwrites it to `false`. The field was never intended as a hook completion marker but was the only candidate — and it had already been reset.

## Solution

Two-layer defense:

1. **state.json `completedAt` field** (primary): `checkFeatureCompletion` reads the new `completedAt` field from state.json. If non-empty, the entire completion flow is skipped at the entry point. `completeFeature` writes this timestamp only after a successful commit.

2. **`updateFileStatus` idempotency** (secondary): Before rewriting, check if the status line already matches the target value. If so, return nil without touching the file.

```go
// Primary guard in checkFeatureCompletion
if state := feature.ReadForgeState(projectRoot); state != nil && state.CompletedAt != "" {
    return nil
}

// Secondary guard in updateFileStatus
if strings.TrimSpace(line) == "status: "+value {
    return nil
}
```

## Reusable Pattern

Any hook that modifies tracked files and commits must be idempotent at two levels:
- **Entry guard**: Check a persistent marker (state file, db flag) before entering the main logic. Don't rely on inspecting the output file — the marker is the source of truth.
- **Write guard**: Before writing a file, verify the change is actually needed. String reconstruction almost always introduces invisible diffs.

Pattern: **hook = read marker → if set, exit → do work → write marker**. The marker must be written by the same code path that does the work, not by a separate caller that may be overridden.

## Related Files

- forge-cli/internal/cmd/feature_complete.go
- forge-cli/pkg/feature/forge_state.go
- plugins/forge/hooks/hooks.json (Stop hook config)
