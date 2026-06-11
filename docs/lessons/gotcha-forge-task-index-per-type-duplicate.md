---
created: "2026-05-16"
tags: [architecture, testing]
---

# `forge task index` Re-run Creates Duplicate Per-Type Tasks

## Problem

After T-quick-1 (gen-cases) generates `test-cases.md` with type information, `forge task index` is re-run and creates a per-type variant (e.g., `T-quick-2-cli`) alongside the original generic task (`T-quick-2`). Both tasks do the same work. The original generic task becomes an orphan — completed but nothing downstream depends on it. The per-type variant takes over the dependency chain.

## Root Cause

Causal chain (4 levels):

1. **Symptom**: Two gen-scripts tasks (T-quick-2 and T-quick-2-cli) with identical dependencies on T-quick-1. T-quick-2 completes first, then T-quick-2-cli runs and may overwrite its output.
2. **Direct cause**: `forge task index` was run twice — once before test-cases.md existed (creates generic T-quick-2), once after (creates per-type T-quick-2-cli from the `CLI: 7` summary table).
3. **Root cause**: `forge task index` creates per-type variants when `test-cases.md` has type counts > 0, but does NOT remove or skip the original generic task. The dependency graph is rewritten: T-quick-3 depends on T-quick-2-cli instead of T-quick-2, leaving T-quick-2 orphaned.
4. **Trigger condition**: Quick mode workflow calls `forge task index` in `/quick-tasks` Step 5 (before test cases exist), and the task executor re-runs it after T-quick-1 generates test cases with type information.

## Solution

`forge task index` should handle the case where per-type variants replace a generic task:

1. When creating per-type variants, check if the generic task already exists and is completed → skip per-type creation, keep generic.
2. Or: when creating per-type variants, remove the generic task from the index and delete its .md file.
3. Or: only run `forge task index` AFTER test-cases.md exists, so per-type detection happens in one pass.

## Reusable Pattern

**Beware of two-pass index generation in quick mode.** The first pass (at `/quick-tasks` time) has no type info and creates generic tasks. The second pass (after gen-cases) has type info and creates per-type variants. If both passes create gen-scripts tasks, you get duplicates.

For quick mode specifically: the initial `forge task index` should either:
- Skip gen-scripts creation until test-cases.md is available, OR
- Create per-type tasks based on a placeholder that gets refined later

## Example

```
# First run (no test-cases.md): creates generic T-quick-2
forge task index --feature my-feature

# After gen-cases produces test-cases.md with CLI type:
forge task index --feature my-feature
# → Creates T-quick-2-cli, updates T-quick-3 deps
# → T-quick-2 still exists but is now orphaned
```

## Related Files

- `docs/features/e2e-test-quality-cleanup/tasks/index.json` — Shows both T-quick-2 and T-quick-2-cli
- `docs/features/e2e-test-quality-cleanup/testing/test-cases.md` — Source of CLI type detection
- `plugins/forge/skills/quick-tasks/SKILL.md` — Calls `forge task index` in Step 5
