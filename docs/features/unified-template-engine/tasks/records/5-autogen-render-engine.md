---
status: "completed"
started: "2026-05-28 01:45"
completed: "2026-05-28 01:48"
time_spent: "~3m"
---

# Task Record: 5 迁移 pkg/task/autogen 渲染引擎并清理 scope

## Summary
Migrated pkg/task/autogen renderBody() from strings.ReplaceAll to text/template.Execute(), defined autogenTemplateData struct, migrated 12 non-record template placeholders to {{.X}} format, unified {{SCOPE}} handling with conditional blocks and inline replacement, removed BodyContext.Scope field, removed removeLineContaining()/removeSection() post-processing functions, added ValidateAutogenTemplates() with missingkey=error validation.

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/task/autogen.go
- forge-cli/pkg/task/autogen_test.go
- forge-cli/pkg/task/extract.go
- forge-cli/pkg/task/extract_test.go
- forge-cli/pkg/task/data/code-quality-simplify.md
- forge-cli/pkg/task/data/doc-consolidate.md
- forge-cli/pkg/task/data/doc-drift.md
- forge-cli/pkg/task/data/doc-review.md
- forge-cli/pkg/task/data/eval-contract.md
- forge-cli/pkg/task/data/eval-journey.md
- forge-cli/pkg/task/data/test-gen-contracts.md
- forge-cli/pkg/task/data/test-gen-journeys.md
- forge-cli/pkg/task/data/test-gen-scripts.md
- forge-cli/pkg/task/data/test-run.md
- forge-cli/pkg/task/data/validation-code.md
- forge-cli/pkg/task/data/validation-ux.md

### Key Decisions
- Used text/template with autogenTemplateData struct replacing string-based ReplaceAll approach
- Paragraph-level {{SCOPE}} in test-gen-contracts/test-gen-journeys wrapped with {{if .SurfaceKey}} conditional blocks
- Inline {{SCOPE}} in doc-consolidate/test-run replaced directly with {{.SurfaceKey}}
- BodyContext.Scope ([]string) removed entirely, unified to SurfaceKey (string)

## Test Results
- **Tests Executed**: Yes
- **Passed**: 845
- **Failed**: 0
- **Coverage**: 0.0%

## Acceptance Criteria
- [x] autogenTemplateData struct with all proposal fields, renderBody() uses text/template.Execute()
- [x] 12 non-record template placeholders migrated to {{.X}} format, no {{SCOPE}} residue in pkg/task/data/
- [x] BodyContext.Scope field deleted, semantics covered by SurfaceKey/SurfaceType
- [x] removeLineContaining() and removeSection() removed from autogen.go
- [x] ValidateAutogenTemplates() with missingkey=error + zero-value struct Execute validation passes

## Notes
Fix-record recovery task. Implementation was already committed (c6bcfd15). All verification steps (compile, fmt, lint, unit-test) passed.
