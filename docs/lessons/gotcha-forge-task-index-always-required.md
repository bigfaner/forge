---
created: "2026-05-17"
tags: [architecture, testing]
---

# `forge task index` Is Always Required — Even for Docs-Only Features

## Problem

When generating tasks for a docs-only feature (all tasks produce `.md` files, no compilation or test execution), the agent skipped `forge task index` reasoning that since there are no test tasks (Step 4) and no profile needed (Step 0), index generation could also be skipped.

When `forge task claim` was subsequently invoked, it failed:

```
ERROR_CODE: NOT_FOUND
ERROR: Failed to load task index
CAUSE: failed to read index: open ...\tasks\index.json: The system cannot find the file specified.
HINT: cat ...\tasks\index.json
```

## Root Cause

**Causal chain** (3 levels):

1. **Symptom**: `forge task claim` fails with NOT_FOUND because index.json doesn't exist
2. **Direct cause**: Agent skipped `forge task index` (quick-tasks Step 5) based on incorrect inference that docs-only fast path covers index generation
3. **Root cause**: The quick-tasks SKILL.md Docs-Only Fast Path explicitly lists only Step 0 and Step 4 as skippable. Step 5 (`forge task index`) generates `index.json` which is the **operational contract** between task files and the CLI — it is required for ALL features regardless of type. The agent over-generalized the "skip" scope.

**Secondary issue**: `forge task claim` error HINT says `cat index.json` (which doesn't exist) instead of suggesting the fix command `forge task index --feature <slug>`. This makes the error self-unresolvable — the user/agent must already know the fix.

## Solution

1. Always run `forge task index --feature <slug>` after creating task `.md` files, even for docs-only features
2. The docs-only fast path skips only: Step 0 (profile resolution) and Step 4 (test task generation). Step 5 (index generation) is never skippable.

For the CLI HINT: `forge task claim` should suggest `forge task index --feature <slug>` when index.json is missing, not `cat <path>`.

## Reusable Pattern

**`forge task index` is never optional.** It transforms `.md` task files into the `index.json` that `forge task claim`, `forge task status`, and `forge task validate-index` all depend on. The docs-only fast path in quick-tasks only skips profile resolution and test task generation — index generation remains mandatory.

**Rule**: If you wrote task `.md` files, you must run `forge task index --feature <slug>` before any `forge task claim`.

## Example

```bash
# After writing task .md files (Step 3 of quick-tasks)
forge task index --feature my-feature

# For docs-only features, skip --test-profiles is fine since no test tasks
# are generated, but the index.json itself is still created:
forge task index --feature my-feature

# Then validate:
forge task validate-index docs/features/my-feature/tasks/index.json

# Then claim works:
forge task claim  # succeeds
```

## Related Files

- `plugins/forge/skills/quick-tasks/SKILL.md` — Docs-Only Fast Path section
- `forge-cli/internal/cmd/task.go` — `forge task claim` error handling
- `forge-cli/pkg/task/index.go` — index.json generation

## References

- `plugins/forge/skills/quick-tasks/SKILL.md` Step 5: "After all business task `.md` files (Step 3) are written, run `forge task index --feature <slug>`"
