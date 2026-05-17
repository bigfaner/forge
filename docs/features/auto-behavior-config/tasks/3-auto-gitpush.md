---
id: "3"
title: "Add auto.gitPush step to run-tasks skill"
priority: "P1"
estimated_time: "1h"
dependencies: ["1"]
scope: "all"
breaking: false
type: "enhancement"
mainSession: false
---

# 3: Add auto.gitPush step to run-tasks skill

## Description

Add a post-completion step to `/run-tasks` that reads `auto.gitPush` from `.forge/config.yaml` and runs `git push` after the all-completed hook passes.

## Reference Files
- `docs/proposals/auto-behavior-config/proposal.md` — Source proposal

## Acceptance Criteria
- [ ] `auto.gitPush=true` → `git push` runs after all-completed hook succeeds
- [ ] `auto.gitPush=false` or absent → no push (backward compatible)
- [ ] Push failures produce clear error message (auth, no remote) — no crash
- [ ] Pushes to configured remote (default: origin)

## Hard Rules
- Only push after all-completed hook passes (compile + fmt + lint + test + e2e regression all green)
- Never force push
- Catch and report git push errors gracefully

## Implementation Notes
- File: `plugins/forge/commands/run-tasks.md`
- The run-tasks command is a markdown skill that instructs the agent. Add a post-completion step that reads `auto.gitPush` from `.forge/config.yaml` and conditionally runs `git push`.
- Also add a Go config accessor if needed (check if `profile.ReadConfig` already exposes the `auto` block from Task 1's schema changes).
