---
id: "1"
title: "Add CategoryForType() and extend RecordData for all categories"
priority: "P0"
estimated_time: "1.5h"
dependencies: []
scope: "backend"
breaking: false
type: "coding.feature"
mainSession: false
---

# 1: Add CategoryForType() and extend RecordData for all categories

## Description

Foundation for the entire typed-task-records feature. Two deliverables:

1. **`CategoryForType(typ string) string`** — maps all 21 task types to 5 categories: `coding`, `doc`, `test`, `validation`, `gate`. Exported category constants for use by templates, validation, and prompts.

2. **Full `RecordData` extension** — add optional field groups for all categories:

| Category | New Fields | Type |
|----------|-----------|------|
| doc | `ReferencedDocs`, `ReviewStatus`, `DocMetrics` | `[]string`, `string`, `string` |
| test | `CasesGenerated`, `CasesEvaluated`, `ScriptsCreated`, `TestResults` | `int`, `int`, `[]string`, `string` |
| validation | `ValidationPassed`, `IssuesFound` | `bool`, `[]string` |
| gate | `GatePassed`, `GateChecks` | `bool`, `[]string` |

All new fields are optional (`omitempty`) — backward compatible with existing record.json data.

## Reference Files
- `docs/proposals/typed-task-records/proposal.md` — Source proposal
- `forge-cli/pkg/task/types.go` — RecordData struct, type constants (lines 276-289)
- `forge-cli/pkg/task/build.go` — IsTestableType (existing prefix pattern, line 433)

## Acceptance Criteria
- [ ] `CategoryForType()` returns correct category for all 21 types
- [ ] `CategoryForType("")` returns `"coding"` as default
- [ ] `CategoryForType("code-quality.simplify")` returns `"coding"`
- [ ] Category constants exported: `CategoryCoding`, `CategoryDoc`, `CategoryTest`, `CategoryValidation`, `CategoryGate`
- [ ] `RecordData` has all 11 new optional fields with `json` tags and `omitempty`
- [ ] Existing `RecordData` JSON deserialization is backward compatible
- [ ] Unit tests: CategoryForType covers all 21 types + empty string + unknown type; RecordData JSON round-trip for old and new fields

## Hard Rules
- Place `CategoryForType()` in new file `forge-cli/pkg/task/category.go`
- Use `strings.HasPrefix` for `coding.*`, `doc*`, `test.*`, `validation.*` — explicit match for `gate` and `code-quality.simplify`
- No changes to `fillRecordTemplate()` or `validateRecordData()` in this task

## Implementation Notes
- The 5-category model enables Phase 2's per-category template files. Even though Phase 1 only needed doc vs non-doc, we define all categories upfront.
- `TestResults string` is free-text (not structured) to keep it simple; the template decides how to render it.
