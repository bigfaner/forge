---
status: "completed"
started: "2026-05-28 01:23"
completed: "2026-05-28 01:43"
time_spent: "~20m"
---

# Task Record: 5 迁移 pkg/task/autogen 渲染引擎并清理 scope

## Summary
Migrated pkg/task/autogen renderBody() from strings.ReplaceAll to text/template.Execute(), defined autogenTemplateData struct, migrated 12 template files to {{.X}} syntax, removed BodyContext.Scope field, removed removeLineContaining()/removeSection() post-processing functions, updated ValidateAutogenTemplates() with missingkey=error + zero-value execution

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/task/autogen.go
- forge-cli/pkg/task/autogen_test.go
- forge-cli/pkg/task/extract.go
- forge-cli/pkg/task/extract_test.go
- forge-cli/pkg/task/data/test-gen-contracts.md
- forge-cli/pkg/task/data/test-gen-journeys.md
- forge-cli/pkg/task/data/doc-consolidate.md
- forge-cli/pkg/task/data/test-run.md
- forge-cli/pkg/task/data/code-quality-simplify.md
- forge-cli/pkg/task/data/doc-drift.md
- forge-cli/pkg/task/data/doc-review.md
- forge-cli/pkg/task/data/eval-contract.md
- forge-cli/pkg/task/data/eval-journey.md
- forge-cli/pkg/task/data/test-gen-scripts.md
- forge-cli/pkg/task/data/validation-code.md
- forge-cli/pkg/task/data/validation-ux.md

### Key Decisions
- autogenTemplateData pre-formats all fields (SurfaceTypes, AcceptanceCriteria, DocTaskCriteria) as strings before template rendering, keeping template logic simple
- {{SCOPE}} paragraph-level templates (test-gen-contracts, test-gen-journeys) use {{if .SurfaceKey}} to wrap entire Scope section; inline templates (doc-consolidate, test-run) use {{if .SurfaceKey}} on the Feature Context line
- extractScope() function retained in extract.go for potential future use even though BodyContext.Scope is removed

## Test Results
- **Tests Executed**: Yes
- **Passed**: 87
- **Failed**: 0
- **Coverage**: 86.9%

## Acceptance Criteria
- [x] autogenTemplateData struct contains all proposal fields, renderBody() uses text/template.Execute()
- [x] 12 non-record template files migrated to {{.X}} format, no {{SCOPE}} remnants in pkg/task/data/
- [x] BodyContext struct Scope field deleted, semantics covered by SurfaceKey/SurfaceType
- [x] removeLineContaining() and removeSection() removed from autogen.go
- [x] ValidateAutogenTemplates() uses missingkey=error + zero-value struct Execute validation

## Notes
All 17 pkg test suites pass. No {{SCOPE}} or old-style {{PLACEHOLDER}} format remaining in template files.
