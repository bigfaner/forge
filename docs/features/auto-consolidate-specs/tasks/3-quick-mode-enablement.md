---
id: "3"
title: "Enable consolidate-specs in quick mode"
priority: "P1"
estimated_time: "30min"
dependencies: []
scope: "backend"
breaking: false
type: "feature"
mainSession: false
---

# 3: Enable consolidate-specs in quick mode

## Description

Enable the consolidate-specs task (T-quick-specs-1) in quick mode by changing config defaults and verifying the quick-tasks template supports the slot. Currently `auto.consolidateSpecs.quick: false` disables it.

## Reference Files
- `docs/proposals/auto-consolidate-specs/proposal.md` — Source proposal
- `forge-cli/pkg/profile/config.go` — Config defaults (`AutoConfigDefaults`)
- `.forge/config.yaml` — Project config
- `plugins/forge/skills/quick-tasks/templates/index.json` — Quick-tasks template
- `forge-cli/pkg/task/testgen.go` — Test task generation (`GetQuickTestTasks`)

## Acceptance Criteria

- [ ] `config.go`: Default for `ConsolidateSpecs.Quick` changed from `false` to `true`
- [ ] `forge-config.example.yaml`: Example reflects new default (`quick: true`)
- [ ] `forge-config.schema.json`: Schema description updated if needed
- [ ] `index.json` quick-tasks template: includes T-quick-specs-1 slot placeholder
- [ ] `forge task index --feature <slug>` generates T-quick-specs-1 when config is default
- [ ] `.forge/config.yaml` updated to reflect new default (`quick: true`)

## Hard Rules

- The Go code in `testgen.go` (`GetQuickTestTasks`) already has the generation logic — only config defaults need changing
- Must not break full-mode behavior (`auto.consolidateSpecs.full: true` stays default)

## Implementation Notes

- The quick-tasks template may not need changes if `forge task index` dynamically generates T-quick-specs-1 — verify this during implementation
- `testgen.go` line 161-168: already generates T-quick-specs-1 when `auto.ConsolidateSpecs.Quick` is true
- The config schema and example should be updated to match the new default
