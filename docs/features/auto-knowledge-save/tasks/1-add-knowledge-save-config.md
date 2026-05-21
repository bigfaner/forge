---
id: "1"
title: "Add KnowledgeSave ModeToggle to AutoConfig + JSON Schema + config get"
priority: "P0"
estimated_time: "1h"
dependencies: []
scope: "backend"
breaking: true
type: "coding.feature"
mainSession: false
---

# 1: Add KnowledgeSave ModeToggle to AutoConfig + JSON Schema + config get

## Description

Add `KnowledgeSave ModeToggle` to the Go `AutoConfig` struct with defaults `{Quick: true, Full: false}`, update the JSON Schema, add `auto.knowledgeSave` support to `forge config get`, and update the example config YAML.

## Reference Files
- `docs/proposals/auto-knowledge-save/proposal.md` — Source proposal

## Acceptance Criteria
- [ ] `AutoConfig` struct has `KnowledgeSave ModeToggle` field with yaml tag `"knowledgeSave"`
- [ ] `AutoConfigDefaults()` sets `KnowledgeSave: ModeToggle{Quick: true, Full: false}`
- [ ] `IsZero()` checks the new field
- [ ] `WithDefaults()` / `applyDefaults()` handles the new field
- [ ] `parseAutoRaw()` includes `"knowledgeSave"` in `modeFields`
- [ ] `getAutoKeyValue()` handles `"auto.knowledgeSave"` returning `"quick:<val> full:<val>"` format
- [ ] JSON Schema (`forge-config.schema.json`) has `knowledgeSave` under `auto.properties` with quick/full boolean properties and correct descriptions/defaults
- [ ] Example YAML (`forge-config.example.yaml`) documents `knowledgeSave` with defaults
- [ ] Existing tests pass (`go test ./...`)
- [ ] New unit test: `TestGetConfigValue` case for `"auto.knowledgeSave"` returns correct format

## Hard Rules
- Follow existing `ModeToggle` patterns exactly (see `consolidateSpecs` field for reference)
- `forge config get auto.knowledgeSave` must return `"quick:<val> full:<val>"` format (same as `auto.runTasks`)
- No changes to skill/markdown files in this task

## Implementation Notes
- The `getAutoKeyValue()` function currently only handles `auto.runTasks` and `auto.gitPush`. Add a new `if` block for `auto.knowledgeSave` before the `"auto.gitPush"` check.
- The `modeFields` slice in `parseAutoRaw()` lists field names to track for explicit-set detection — add `"knowledgeSave"`.
- JSON Schema defaults should reflect the Go defaults: `quick: true` (default: true), `full: false` (default: false).
