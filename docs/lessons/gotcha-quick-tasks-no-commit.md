---
created: "2026-05-18"
tags: [architecture, local-dev-deployment]
---

# quick-tasks generates planning docs but never commits them

## Problem

After running `/quick-tasks` followed by `/run-tasks`, task definition files (`tasks/*.md`, `index.json`, `manifest.md`) remain untracked in git. Only the code changes from task executors get committed. Discovered in auto-consolidate-specs feature where 2 of 3 task files were left untracked.

## Root Cause

1. **Surface**: Task .md files are untracked after the full quick-tasks → run-tasks pipeline completes
2. **Immediate cause**: `run-tasks` task executor agents only commit files they modified (code changes), not pre-existing planning docs
3. **Structural cause**: `quick-tasks` skill has Steps 0-7 but no commit step — it generates task files, runs `forge task index`, creates manifest, validates, then stops. The skill assumes someone else will commit the planning artifacts
4. **Root gap**: No skill in the pipeline takes responsibility for committing planning docs. quick-tasks creates them, run-tasks consumes them, but neither commits them

## Solution

quick-tasks should add a Step 8 after validation: commit all generated planning artifacts (task .md files, index.json, manifest.md) as a single commit before handing off to `/run-tasks`.

## Reusable Pattern

Any skill that generates planning artifacts (quick-tasks, breakdown-tasks, write-prd, tech-design) should commit its output before the pipeline continues. The principle: **the skill that creates artifacts is responsible for persisting them**. Downstream consumers (run-tasks, task executors) should not be expected to commit files they didn't create.

## Example

```bash
# Should be Step 8 in quick-tasks, after forge task validate-index
git add docs/features/<slug>/tasks/*.md docs/features/<slug>/tasks/index.json docs/features/<slug>/manifest.md
git commit -m "chore(<slug>): add task definitions for quick pipeline"
```

## Related Files

- `plugins/forge/skills/quick-tasks/SKILL.md` — missing commit step
- `plugins/forge/skills/breakdown-tasks/SKILL.md` — likely same gap (verify)
- `plugins/forge/skills/run-tasks/SKILL.md` — consumer, should not need to commit planning docs

## References

- auto-consolidate-specs feature: 2 of 3 task files left untracked after full pipeline run
