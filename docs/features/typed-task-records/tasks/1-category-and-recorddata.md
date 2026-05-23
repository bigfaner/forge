---
id: "1"
title: "Add CategoryForType() and doc-specific RecordData fields"
priority: "P0"
estimated_time: "1h"
dependencies: []
scope: "backend"
breaking: false
type: "coding.feature"
mainSession: false
---

# 1: Add CategoryForType() and doc-specific RecordData fields

## Description

Foundation for type-differentiated records. Two changes:

1. **`CategoryForType(typ string) string`** in `forge-cli/pkg/task/` — maps all 21 task types to categories: `coding`, `doc`, `test`, `validation`, `gate`. This is the canonical categorization used by record rendering and validation.

2. **Extend `RecordData`** in `types.go` with optional doc-specific fields: `ReferencedDocs []string`, `ReviewStatus string`, `DocMetrics string`. These fields carry doc-task metadata that the current uniform template ignores.

## Reference Files
- `docs/proposals/typed-task-records/proposal.md` — Source proposal
- `forge-cli/pkg/task/types.go` — RecordData struct, type constants
- `forge-cli/pkg/task/build.go` — IsTestableType (existing type prefix pattern)

## Acceptance Criteria
- [ ] `CategoryForType()` returns correct category for all 21 types
- [ ] `CategoryForType("")` returns `"coding"` as default
- [ ] `CategoryForType("code-quality.simplify")` returns `"coding"`
- [ ] `RecordData` has new optional fields: `ReferencedDocs`, `ReviewStatus`, `DocMetrics` with `json` tags and `omitempty`
- [ ] Existing `RecordData` JSON deserialization is backward compatible (new fields optional)
- [ ] Unit tests cover all 21 types + empty string + unknown type

## Hard Rules
- Category constants must be exported strings (e.g., `CategoryCoding = "coding"`)
- Follow the existing prefix-matching pattern from `IsTestableType` where possible
- No changes to `fillRecordTemplate()` or `validateRecordData()` in this task

## Implementation Notes
- Place `CategoryForType()` in a new file `forge-cli/pkg/task/category.go` to keep `types.go` focused on data structures
- Use `strings.HasPrefix` for `coding.*`, `doc*`, `test.*`, `validation.*` — fallthrough to explicit matches for `gate` and `code-quality.simplify`
- Phase 1 only needs doc vs non-doc distinction, but define all 5 categories for future use
