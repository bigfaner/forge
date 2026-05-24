---
id: "4"
title: "forge task index: migrate from Interfaces to Surfaces"
priority: "P1"
estimated_time: "1.5h"
dependencies: ["1"]
scope: "backend"
breaking: false
type: "coding.enhancement"
mainSession: false
---

# 4: forge task index: migrate from Interfaces to Surfaces

## Description

Switch `forge task index` (BuildIndex) from reading `Interfaces []string` to reading `Surfaces map[string]string`. Extract deduplicated surface types from the map. Handle empty surfaces with explicit error. Unify naming: remove `uiInterfaces` map that checks `web`/`mobile` separately, use the new unified type names directly.

## Reference Files
- `proposal.md#统一命名规范` — old web-ui/mobile-ui → new web/mobile mapping table
- `proposal.md#Non-Functional-Requirements` — empty surfaces error behavior
- `proposal.md#未知类型处理策略` — unknown types ignored + warn log, not passed downstream
- `proposal.md#Success-Criteria` — task index generates correct test tasks, empty surfaces errors, naming unification

## Acceptance Criteria

- [ ] `BuildIndex` reads `Surfaces` map instead of `Interfaces` slice via `forgeconfig.ReadSurfaces()`
- [ ] Deduplicated surface types extracted from map values (e.g., `{frontend: web, backend: api}` → `[web, api]`)
- [ ] Empty surfaces (nil or 0 entries) → error exit with "run forge init" guidance
- [ ] Unknown surface types logged as warning and ignored (not passed to test task generation)
- [ ] `uiInterfaces` map in `autogen.go` removed — test task generation uses unified type names directly
- [ ] Old naming values (`web-ui`, `mobile-ui`) no longer appear in Go code
- [ ] `GetBreakdownTestTasks` and `GetQuickTestTasks` accept `[]string` surface types (unchanged signature, just different source)

## Hard Rules

- Unknown types MUST be ignored with `log.Warn`, not error — prevents breaking existing configs with typos
- Empty surfaces MUST error (exit 1) — never silently skip, that's the original bug

## Implementation Notes

- Key files: `forge-cli/pkg/task/build.go` (line 66: `ReadInterfaces` → `ReadSurfaces`), `forge-cli/pkg/task/autogen.go` (lines 34-50: remove `uiInterfaces` map and `hasUIInterface`), `forge-cli/pkg/task/extract.go` (line 15: update body context)
- The `BodyContext.Interfaces` field should be renamed to `BodyContext.SurfaceTypes` or similar
