---
id: "4"
title: "Add knowledgeSave and runTasks to forge init guided config"
priority: "P1"
estimated_time: "45m"
dependencies: ["1"]
scope: "backend"
breaking: false
type: "coding.enhancement"
mainSession: false
---

# 4: Add knowledgeSave and runTasks to forge init guided config

## Description

Both `forge init` (TUI via `huh`) and `forge config init` (CLI via stdin) currently ask about e2eTest, consolidateSpecs, cleanCode, validation, and gitPush — but NOT about `runTasks` (existing field, missing from init) or `knowledgeSave` (new field from Task 1). Add guided configuration for both fields to both init paths.

## Reference Files
- `docs/proposals/auto-knowledge-save/proposal.md` — Source proposal

## Acceptance Criteria
- [ ] `askAutoBehavior()` in `init.go` includes `runTasks` questions (quick/full) with correct defaults
- [ ] `askAutoBehavior()` in `init.go` includes `knowledgeSave` questions (quick/full) with correct defaults
- [ ] `runConfigInit()` in `config.go` includes `runTasks` prompts (quick/full) with correct defaults
- [ ] `runConfigInit()` in `config.go` includes `knowledgeSave` prompts (quick/full) with correct defaults
- [ ] Config struct construction in both files includes the new fields
- [ ] Question ordering follows the logical grouping of auto config fields
- [ ] Existing init tests pass; new test cases cover the added questions

## Hard Rules
- Follow existing question patterns exactly (see `consolidateSpecs` questions for reference)
- Defaults must match `AutoConfigDefaults()`: runTasks `{Quick: true, Full: false}`, knowledgeSave `{Quick: true, Full: false}`
- TUI questions use `huh.NewConfirm`; CLI prompts use `readBool(reader, default)`

## Implementation Notes
- `askAutoBehavior()` is in `forge-cli/internal/cmd/init.go` (~lines 243-338)
- `runConfigInit()` is in `forge-cli/internal/cmd/config.go` (~lines 87-170)
- Add questions after the existing validation questions and before gitPush
- Suggested order: e2eTest → consolidateSpecs → cleanCode → validation → runTasks → knowledgeSave → gitPush
