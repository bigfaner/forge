---
id: "3"
title: "Improve init summary to show actual detected surface types"
priority: "P1"
estimated_time: "1h"
dependencies: ["1"]
scope: "backend"
breaking: false
type: "coding.enhancement"
mainSession: false
---

# 3: Improve init summary to show actual detected surface types

## Description
Change the `forge init` summary format from opaque `"N mappings"` to actual detected surface types with source annotations. Show `CREATED surfaces cli (inferred:cmd-dir)` instead of `CREATED surfaces (1 mappings)`. This makes the init output immediately understandable without requiring the user to inspect the config file.

## Reference Files
- `proposal.md#Proposed-Solution` — init summary format before/after examples
- `proposal.md#Requirements-Analysis` — Key Scenario 4 for subdir detection display, Key Scenario 5 for existing signal detection
- `proposal.md#Success-Criteria` — init summary format criteria (#6-7)
- `forge-cli/internal/cmd/init.go` — runSurfaceConfig (L483-528), detail string construction at L524-527, printInitSummary at L427

## Acceptance Criteria
- [ ] Scalar form: `CREATED surfaces cli` when detected from dependency; `CREATED surfaces cli (inferred:cmd-dir)` when inferred
- [ ] Map form: `CREATED surfaces forge-cli=cli (inferred:cmd-dir)` per entry instead of `(N mappings)`
- [ ] Subdir detection: `forge init` at root with `forge-cli/cli` subdir → summary shows `CREATED surfaces forge-cli=cli (inferred:cmd-dir)`, not `(1 mappings)`
- [ ] Existing signal detection: cobra detected → `CREATED surfaces cli (from cobra)` or `CREATED surfaces cli` — format consistent with inference annotation
- [ ] Re-run skip: `SKIPPED surfaces (already configured)` when user confirms existing config

## Hard Rules
- Summary format change must not break scripted consumers — verify no programmatic parser exists in codebase via grep

## Implementation Notes
- Modify `runSurfaceConfig` at L524-527 where detail string is built
- Current logic: scalar → `surfaces["."]`, map → `fmt.Sprintf("%d mappings", len(surfaces))`
- New logic: iterate surfaces map, build detail string from each entry with source annotation from `DetectResult.Sources`
- Need Sources info to reach `runSurfaceConfig` — either change `askSurfaceConfirmation` return type or pass Sources through `initAction`
- `initAction` struct at init.go may need a Sources field, or the detail string can be constructed inside `runSurfaceConfig` before returning the action
- Consider: `runSurfaceConfig` already has access to `DetectResult` if `askSurfaceConfirmation` is modified to return it (Task 2)
