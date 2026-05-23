---
id: "2"
title: "Refactor fillRecordTemplate() for type-differentiated rendering"
priority: "P0"
estimated_time: "1.5h"
dependencies: ["1"]
scope: "backend"
breaking: false
type: "coding.refactor"
mainSession: false
---

# 2: Refactor fillRecordTemplate() for type-differentiated rendering

## Description

Replace the uniform `fillRecordTemplate()` with category-aware conditional rendering. Doc tasks produce "Document Metrics" / "Referenced Documents" / "Review Status" sections instead of "Test Results" / "Coverage". Coding tasks produce identical output to current format (backward compatible).

## Reference Files
- `docs/proposals/typed-task-records/proposal.md` — Source proposal
- `forge-cli/internal/cmd/task/submit.go` — `fillRecordTemplate()` (lines 367-448)
- `forge-cli/pkg/task/types.go` — `RecordData` struct
- `forge-cli/pkg/task/category.go` — `CategoryForType()` (from task 1)

## Acceptance Criteria
- [ ] Doc task records contain zero test-related sections (no "Test Results", no "Coverage")
- [ ] Doc task records include "Document Metrics", "Referenced Documents", "Review Status" sections
- [ ] Coding task records are byte-identical to current format (backward compatible)
- [ ] `CategoryForType()` is used to determine rendering branch
- [ ] Test/validation/gate task records render with "Test Results" section (same as coding, Phase 1 conservative)
- [ ] Unit tests for each category's rendered output

## Hard Rules
- Phase 1 scope: only doc vs non-doc rendering split (2 branches). Do NOT add test/validation/gate-specific templates yet.
- Must call `CategoryForType(t.Type)` — no inline type-checking
- `fillRecordTemplate` signature unchanged: `(t *task.Task, rd *task.RecordData, startedTime string) string`

## Implementation Notes
- The proposal's Phase 2 will add per-category templates via `text/template`. Phase 1 uses conditional branches in existing `fillRecordTemplate()` — simpler change, same effect for doc types.
- Doc-specific sections should use data from the new `RecordData` fields (`ReferencedDocs`, `ReviewStatus`, `DocMetrics`), falling back to "N/A" when empty.
