---
created: "2026-05-19"
tags: [architecture, testing]
---

# quick-tasks does not auto-chain to run-tasks

## Problem
After `/quick-tasks` generates planning artifacts and commits them, the user expects execution to start automatically but the session stops. This creates friction — the user must manually invoke `/run-tasks` to begin work.

## Root Cause
1. `/quick-tasks` is a standalone planning skill, not a pipeline orchestrator
2. The separation is intentional: generated tasks need human review before execution starts
3. `/quick` is the full-pipeline skill that chains brainstorm → quick-tasks → run-tasks
4. Individual skills (quick-tasks, run-tasks) are designed to be composable, not monolithic

## Solution
Use `/quick` for the full automated pipeline (brainstorm → tasks → execute). Use `/quick-tasks` standalone when you want to review/edit tasks before execution.

## Reusable Pattern
- Full pipeline (no review): `/quick`
- Plan-only with review gate: `/quick-tasks` → review → `/run-tasks`
- The commit step at the end of `/quick-tasks` is a natural stopping point for review

## Related Files
- `plugins/forge/skills/quick-tasks/SKILL.md`
- `plugins/forge/skills/quick/SKILL.md`
- `plugins/forge/skills/run-tasks/SKILL.md`
