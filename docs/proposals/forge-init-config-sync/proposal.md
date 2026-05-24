---
name: forge-init-config-sync
status: Superseded
superseded-by: unify-surfaces
date: 2026-05-20
---

# Proposal: Sync forge init with Latest config.yaml Structure

## Problem

After commit `3bb6125c` intentionally removed the profile-based project-type/languages/interfaces selection from `forge init`, the config generation step fell behind the latest `.forge/config.yaml` schema. Specifically:

1. **`auto.validation`** (ModeToggle with quick/full) was added to `AutoConfig` but `askAutoBehavior()` in `forge init` and the auto section in `forge config init` never got the corresponding prompts.
2. **`worktree`** config (source-branch + copy-files) exists in `forge config init` (stdin version) but is missing from `forge init` (huh TUI version).
3. **`test-command`** in config.yaml contradicts Forge's design principle of "just as the unified test abstraction layer". It's consumed by `RunProjectTests` (as a bypass override — the default path already uses `just test`) and `journey_isolation.go` (runs raw command in isolated temp dir). Both should go through `just` instead.

This means users running `forge init` get an incomplete config that doesn't reflect the current schema, and the config contains a field (`test-command`) that violates the project's architectural principle.

## Solution

### Part 1: Add missing interactive sections to both init commands

#### forge init (init.go — huh TUI)

1. Add `auto.validation` quick/full confirm prompts in `askAutoBehavior()`, between cleanCode and gitPush.
2. Add worktree config section after auto behavior: optional source-branch input and copy-files multi-select.

#### forge config init (config.go — stdin)

1. Add `auto.validation` quick/full prompts in the auto-behavior section, between cleanCode and gitPush.
2. Worktree prompts already exist — no change needed.

### Part 2: Remove test-command from config

`test-command` violates Forge's "just as unified abstraction" principle. Two consumers exist, both should use `just` instead:

| Consumer | Current behavior | Refactored behavior |
|----------|-----------------|-------------------|
| `RunProjectTests` (quality_gate.go) | Uses `testCommand` as override; empty → falls through to `just test` | Always use the fallback chain (already works when `testCommand` is empty) |
| `executeJourneyInIsolation` (journey_isolation.go) | Runs raw `testCommand` in isolated temp dir | Run `just e2e-test` from project root with journey filter |

Changes required:

1. Remove `TestCommand` field from `forgeconfig.Config` struct
2. Remove `test-command` case from `forgeconfig.GetConfigValue`
3. Remove `TestCommand` field from `task.TaskInfo` / `task.TaskIndex` (pkg/task/types.go)
4. Refactor `journey_isolation.go`: remove `readTestCommand`, change `executeJourneyInIsolation` to use `just e2e-test` from project root
5. Refactor `quality_gate.go` / `RunProjectTests`: remove `testCommand` parameter
6. Update related tests
7. Remove `test-command` from example YAML and JSON schema
8. Version bump

### Field Classification (deliberately excluded)

| Field | Reason for exclusion |
|-------|---------------------|
| `project-type` | Intentionally removed in v3 refactor, not in Go Config struct |
| `languages` | Intentionally removed from init (was profile system) |
| `interfaces` | Intentionally removed from init (was profile system) |
| `test-framework` | Not consumed by any business logic |
| `coverage` | Complex config with built-in defaults, advanced users can edit manually |

## Alternatives

### Do nothing
Users manually edit config.yaml after init. Works but defeats the purpose of interactive init.

### Full coverage (all fields)
Re-add project-type, languages, interfaces, test-framework, test-command, coverage to init. Rejected — project-type/languages/interfaces were intentionally removed, and the rest are advanced overrides or design violations.

## Scope

### In Scope
- Add `auto.validation` ModeToggle prompts to `forge init` (huh TUI)
- Add worktree config section to `forge init` (huh TUI)
- Add `auto.validation` prompts to `forge config init` (stdin)
- Remove `test-command` from `forgeconfig.Config` and all consumers
- Refactor `journey_isolation.go` to use `just` instead of raw test command
- Refactor `RunProjectTests` to remove `testCommand` parameter
- Update existing tests for all affected commands
- Update example YAML and JSON schema to remove `test-command`
- Bump version in scripts/version.txt

### Out of Scope
- Re-adding removed profile fields (project-type, languages, interfaces)
- Adding test-framework/coverage prompts
- Removing `test-framework` from Config struct (separate concern — it's unused but not a design violation like test-command)

## Risks

| Risk | Impact | Mitigation |
|------|--------|-----------|
| Increased init prompt count (4 more screens) | Users may find init slower | Worktree section is skippable (press Enter to skip both fields) |
| huh TUI worktree prompts untested | Tests mock configInitFunc | Add test coverage for new prompt paths |
| journey isolation refactor breaks run-journey | Tests in journey_isolation_test.go cover this | Run existing e2e tests for run-journey |
| RunProjectTests signature change affects callers | Only 2 callers exist (quality_gate + test) | Update both callers |

## Success Criteria

- [ ] `forge init` generates config with `auto.validation` quick/full values
- [ ] `forge init` generates config with worktree section (when user provides values)
- [ ] `forge init` allows skipping worktree (empty source-branch + no copy-files = no worktree block)
- [ ] `forge config init` includes validation quick/full prompts
- [ ] `test-command` field removed from `forgeconfig.Config` struct
- [ ] No code references `test-command` or `TestCommand` in config context
- [ ] `journey_isolation.go` runs tests via `just` instead of raw command
- [ ] `RunProjectTests` signature simplified (no testCommand parameter)
- [ ] All existing tests pass
- [ ] New tests cover validation, worktree prompt paths, and test-command removal
