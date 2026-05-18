---
id: "1"
title: "SKILL.md non-interactive auto-integration + [auto-specs] commit"
priority: "P0"
estimated_time: "1-2h"
dependencies: []
scope: "backend"
breaking: false
type: "enhancement"
mainSession: false
---

# 1: SKILL.md non-interactive auto-integration + [auto-specs] commit

## Description

Modify `plugins/forge/skills/consolidate-specs/SKILL.md` to enable fully automated spec integration in non-interactive (pipeline) mode. Currently Step 6 blocks on `[CROSS]` items and Step 11 lacks a commit tag for auto-integrated changes. The core insight: spec errors have no runtime risk and git revert provides perfect rollback, so auto-integration is safe.

## Reference Files
- `docs/proposals/auto-consolidate-specs/proposal.md` — Source proposal
- `plugins/forge/skills/consolidate-specs/SKILL.md` — Primary target
- `plugins/forge/hooks/guide.md` — Documents auto-behavior config table

## Acceptance Criteria

- [ ] Step 6: In non-interactive mode, all `[CROSS]` items are auto-integrated without blocking (no `blocked` status)
- [ ] Step 6: `[CROSS]` items with >50% overlap still auto-merge, but commit message includes `[auto-specs]` + warning note
- [ ] Step 11: Auto-integrated commits include `[auto-specs]` tag in commit message
- [ ] `git log --grep="[auto-specs]"` finds all auto-integrated commits
- [ ] Manual `/consolidate-specs` interactive behavior is unchanged (CROSS items still prompt user)
- [ ] Drift-only path (Steps 9-11) also uses `[auto-specs]` commit tag

## Hard Rules

- Do NOT remove the interactive code path — only add auto-mode behavior when running in non-interactive (pipeline) context
- `[auto-specs]` commits must be separate from code change commits

## Implementation Notes

- Key risk: auto-integrating wrong rules → mitigated by `[auto-specs]` tag enabling easy revert
- The skill should detect non-interactive mode context (pipeline execution) vs manual invocation
- Consider using a `--non-interactive` flag or detecting the execution context from the task template
