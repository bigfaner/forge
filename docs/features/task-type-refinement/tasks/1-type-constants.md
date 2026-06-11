---
id: "1"
title: "Add new task type constants and deprecate TypeImplementation"
priority: "P0"
estimated_time: "1h"
dependencies: []
scope: "backend"
breaking: true
type: "implementation"
mainSession: false
---

# 1: Add new task type constants and deprecate TypeImplementation

## Description

Foundation task for the type refinement. Add four new business type constants (`TypeFeature`, `TypeEnhancement`, `TypeCleanup`, `TypeRefactor`) to replace the overly broad `TypeImplementation`. Update the type registry, validation map, and inference logic to support the new types.

## Reference Files
- `docs/proposals/task-type-refinement/proposal.md` — Source proposal
- `forge-cli/pkg/task/types.go` — Type constants, ValidTypes map, TaskTypeRegistry
- `forge-cli/pkg/task/infer.go` — InferType function
- `forge-cli/internal/cmd/validate_index.go` — Uses ValidTypes for validation
- `forge-cli/internal/cmd/list_types.go` — Displays TaskTypeRegistry

## Acceptance Criteria
- [ ] `TypeFeature`, `TypeEnhancement`, `TypeCleanup`, `TypeRefactor` constants defined in `types.go`
- [ ] `TypeImplementation` marked as deprecated (comment + registry deprecation note), NOT removed yet
- [ ] `ValidTypes` map includes all 4 new types
- [ ] `TaskTypeRegistry` includes entries for all 4 new types with category "Core business"
- [ ] `InferType()` in `infer.go` returns empty string for new type IDs (they are explicit, not inferred from patterns)
- [ ] `forge list-types` displays all 4 new types
- [ ] All existing tests pass with deprecated `TypeImplementation` still in ValidTypes

## Hard Rules
- Do NOT remove `TypeImplementation` from `ValidTypes` yet — existing index.json files still use it. Migration (task 5) handles the transition.
- New type constant values must be exact string matches: `"feature"`, `"enhancement"`, `"cleanup"`, `"refactor"`.

## Implementation Notes
- `infer.go` currently has no pattern for business types — non-auto-gen tasks get empty string, callers fall back to `TypeImplementation`. After this change, callers (migrate.go, build.go) should fall back to `TypeFeature` instead. But that change belongs to tasks 2 and 5 — this task only adds the constants.
