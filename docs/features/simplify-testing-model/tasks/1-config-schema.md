---
id: "1"
title: "Config schema: replace TestProfiles/Capabilities with Interfaces/Languages"
priority: "P0"
estimated_time: "2h"
dependencies: []
scope: "backend"
breaking: true
type: "refactor"
mainSession: false
---

# 1: Config schema: replace TestProfiles/Capabilities with Interfaces/Languages

## Description
Remove `TestProfiles []string` and `Capabilities []string` fields from `ForgeConfig` struct in `config.go`. Add `Interfaces []string` (replacing capabilities — system-facing interface types) and `Languages []string` (optional override of auto-detected language). Replace `ReadTestProfiles()` with `ReadLanguages()` that returns the `languages` config field if set, otherwise calls `DetectLanguages()`. Global rename `capabilities` → `interfaces` in all Go exported symbols, CLI help text, and internal variable names.

## Reference Files
- `docs/proposals/simplify-testing-model/proposal.md` — Source proposal
- `forge-cli/pkg/profile/config.go` — ForgeConfig struct and config reading
- `forge-cli/pkg/profile/embed.go` — Strategy file embedding and capability lookup

## Acceptance Criteria
- `ForgeConfig` struct has `Interfaces []string` and `Languages []string` fields; `TestProfiles` and `Capabilities` fields removed
- `ReadLanguages()` function exists: returns `config.Languages` if set, otherwise calls `DetectLanguages()`
- `ReadInterfaces()` function exists: returns `config.Interfaces` if set, otherwise defaults to union of all detected languages' capabilities
- Zero Go exported symbols contain "capability" or "Capability" (grep `forge-cli/pkg/` and `forge-cli/internal/`)
- Config YAML field names: `interfaces`, `languages` (snake_case matching existing `project-type` convention)
- `go build ./...` passes

## Hard Rules
- Do not change detection logic (Task 3) or embed paths (Task 2) in this task
- Do not rename `profiles/` directory in this task
- `interfaces` valid values: web-ui, tui, mobile-ui, api, cli (same set as v2 capabilities)

## Implementation Notes
- The `ReadLanguages()` function is the new primary entry point for language resolution, replacing the 3-step fallback chain in `runProfileResolve` (config → detect → "none")
- `Languages` field is optional: when omitted, auto-detection is the primary path (not a fallback)
- The `languageCapabilities` hardcoded map (defined in Task 2) provides the default interfaces per language — `ReadInterfaces()` uses this when `config.Interfaces` is empty
