---
created: "2026-06-05"
tags: [local-dev-deployment, architecture]
---

# Stale Manifest Blocks Auto-Push via --if-done Guard

## Problem

All coding tasks completed and committed, but `forge feature complete --if-done` (Stop hook) never pushed. Running `git push` manually revealed no upstream. The feature branch was 10+ commits ahead with no remote tracking.

## Root Cause

1. **Symptom**: No auto-push after all tasks completed. Manual `git push` fails with "no upstream branch".
2. **Direct cause**: `forge feature complete --if-done` exited 0 silently — the `--if-done` guard saw manifest status as `tasks` (not `completed`) and bailed out without pushing.
3. **Code cause**: The Stop hook reads `manifest.md` status, not `index.json`. The dispatcher loop updates `index.json` task statuses via CLI and direct edits, but never syncs the manifest's status table.
4. **Why it happens**: `quick-tasks` mode generates tasks that the dispatcher completes, but the manifest is a static artifact from the `quick-tasks` generation step. It is never updated after task execution.

## Solution

After all tasks complete, the dispatcher (or the `run-tasks` skill) should update the manifest status to `completed` before the Stop hook fires. Alternatively, `forge feature complete --if-done` should read from `index.json` as the source of truth rather than the manifest.

As an immediate workaround, manually set upstream before expecting auto-push:

```bash
git push -u origin HEAD
```

## Reusable Pattern

When a guard condition (`--if-done`) reads from a secondary artifact (manifest) while the primary source of truth (`index.json`) is updated independently, the secondary artifact will drift. Either:
1. Make the guard read from the primary source, or
2. Ensure the secondary artifact is updated as part of the same write path.

## Related Files

- `docs/features/<slug>/manifest.md` — status table (secondary artifact)
- `docs/features/<slug>/tasks/index.json` — task statuses (primary source)
- `forge-cli/internal/cmd/feature_complete.go` — `--if-done` guard logic
