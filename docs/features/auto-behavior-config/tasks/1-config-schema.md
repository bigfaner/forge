---
id: "1"
title: "Add auto config block to schema and example"
priority: "P1"
estimated_time: "45m"
dependencies: []
scope: "backend"
breaking: false
type: "enhancement"
mainSession: false
---

# 1: Add auto config block to schema and example

## Description

Extend the forge config schema with an `auto` object containing mode-scoped config fields (`e2eTest`, `consolidateSpecs`, `cleanCode` with `quick`/`full` sub-keys) and a global `gitPush` flag. Update the example config to document these options.

## Reference Files
- `docs/proposals/auto-behavior-config/proposal.md` — Source proposal

## Acceptance Criteria
- [ ] `forge-config.schema.json` defines `auto` object with `e2eTest`, `consolidateSpecs`, `cleanCode` (each with `quick`/`full` bool) and `gitPush` (bool)
- [ ] `additionalProperties: false` preserved on all objects
- [ ] `forge-config.example.yaml` documents all 7 fields with comments
- [ ] Existing configs without `auto` block continue to work (backward compatible)

## Hard Rules
- Keep `additionalProperties: false` everywhere. Add explicit field definitions.
- Defaults: `e2eTest.{quick,true; full,true}`, `consolidateSpecs.{quick,true; full,true}`, `cleanCode.{quick,false; full,false}`, `gitPush: false`

## Implementation Notes
- Schema: `plugins/forge/references/shared/forge-config.schema.json`
- Example: `plugins/forge/references/shared/forge-config.example.yaml`
- Go struct (in Task 2): `AutoConfig` with `ModeToggle` sub-struct for mode-scoped fields
- Task 1 only defines the schema + example. Go config reading happens in Task 2.
