---
status: "completed"
started: "2026-06-09 22:08"
completed: "2026-06-09 22:27"
time_spent: "~19m"
---

# Task Record: 1 Scaffold 核心框架 + cli/tui 模板

## Summary
Implemented `forge justfile scaffold` Cobra sub-command with core framework: command registration (`forge justfile scaffold --type/--key`), parameter validation (unknown type, scalar+key, named without key), placeholder injection using `<<...>>` syntax, dual-platform `[unix]`/`[windows]` variants, and `# user-customized` boundary markers. Scaffolded cli and tui simple surface types with shared recipe set (test + teardown + compile + fmt + lint + unit-test).

## Changes

### Files Created
- forge-cli/internal/cmd/scaffold/types.go
- forge-cli/internal/cmd/scaffold/generate.go
- forge-cli/internal/cmd/scaffold/register.go
- forge-cli/internal/cmd/scaffold/generate_test.go
- forge-cli/internal/cmd/justfile.go

### Files Modified
- forge-cli/pkg/types/surface.go

### Key Decisions
- Used table-driven SurfaceSpec pattern: surfaceSpecs map indexed by type, cli/tui share simpleRecipes()
- ValidateArgs checks all 5 known surface types for arg correctness, then surfaceSpecs for generation support — allows api/web/mobile validation without full recipe support yet
- Added AllSurfaceTypesSet() to pkg/types for O(1) membership check without importing forgeconfig (avoids cycles)
- recipeName() handles scalar (no prefix) vs named (key-verb prefix) naming convention

## Test Results
- **Tests Executed**: Yes
- **Passed**: 20
- **Failed**: 0
- **Coverage**: 81.4%

## Acceptance Criteria
- [x] `forge justfile scaffold --type cli` outputs test + teardown + compile + fmt + lint + unit-test with <<...>> placeholders
- [x] Parameter validation: unknown type errors; scalar+key errors; named without key errors
- [x] cli and tui generate identical recipe sets (test+teardown+quality only, no dev/probe)
- [x] All lifecycle/quality recipes marked # user-customized; scalar surfaces have no prefix
- [x] All recipes include [unix] and [windows] dual-platform variants

## Notes
Task 1 scope is cli/tui only. api/web/mobile surface types are recognized for validation but return 'not yet supported' error for scaffold generation — future tasks will add full specs. All new files in forge-cli/internal/cmd/scaffold/ per Hard Rules.
