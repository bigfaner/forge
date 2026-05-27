---
status: "completed"
started: "2026-05-28 02:15"
completed: "2026-05-28 02:29"
time_spent: "~14m"
---

# Task Record: 8 目录重组：合并、拆分、重命名

## Summary
Reorganized template directories: merged pkg/template/ into pkg/task/, split pkg/task/data/ into pkg/task/templates/ (14 files) and pkg/task/records/ (6 files without record- prefix), renamed pkg/prompt/data/ to pkg/prompt/templates/. Updated all go:embed directives, template path functions, and import references.

## Changes

### Files Created
- forge-cli/pkg/task/tasktemplate.go
- forge-cli/pkg/task/tasktemplate_test.go
- forge-cli/pkg/task/templates/code-quality-simplify.md
- forge-cli/pkg/task/templates/coding.cleanup.md
- forge-cli/pkg/task/templates/coding.fix.md
- forge-cli/pkg/task/templates/doc-consolidate.md
- forge-cli/pkg/task/templates/doc-drift.md
- forge-cli/pkg/task/templates/doc-review.md
- forge-cli/pkg/task/templates/eval-contract.md
- forge-cli/pkg/task/templates/eval-journey.md
- forge-cli/pkg/task/templates/test-gen-contracts.md
- forge-cli/pkg/task/templates/test-gen-journeys.md
- forge-cli/pkg/task/templates/test-gen-scripts.md
- forge-cli/pkg/task/templates/test-run.md
- forge-cli/pkg/task/templates/validation-code.md
- forge-cli/pkg/task/templates/validation-ux.md
- forge-cli/pkg/task/records/coding.md
- forge-cli/pkg/task/records/doc.md
- forge-cli/pkg/task/records/eval.md
- forge-cli/pkg/task/records/gate.md
- forge-cli/pkg/task/records/test.md
- forge-cli/pkg/task/records/validation.md
- forge-cli/pkg/prompt/templates/code-quality-simplify.md
- forge-cli/pkg/prompt/templates/coding-cleanup.md
- forge-cli/pkg/prompt/templates/coding-enhancement.md
- forge-cli/pkg/prompt/templates/coding-feature.md
- forge-cli/pkg/prompt/templates/coding-fix.md
- forge-cli/pkg/prompt/templates/coding-refactor.md
- forge-cli/pkg/prompt/templates/doc-consolidate.md
- forge-cli/pkg/prompt/templates/doc-drift.md
- forge-cli/pkg/prompt/templates/doc-review.md
- forge-cli/pkg/prompt/templates/doc-summary.md
- forge-cli/pkg/prompt/templates/doc.md
- forge-cli/pkg/prompt/templates/eval-contract.md
- forge-cli/pkg/prompt/templates/eval-journey.md
- forge-cli/pkg/prompt/templates/fix-record-missed.md
- forge-cli/pkg/prompt/templates/gate.md
- forge-cli/pkg/prompt/templates/test-gen-contracts.md
- forge-cli/pkg/prompt/templates/test-gen-journeys.md
- forge-cli/pkg/prompt/templates/test-gen-scripts.md
- forge-cli/pkg/prompt/templates/test-run.md
- forge-cli/pkg/prompt/templates/validation-code.md
- forge-cli/pkg/prompt/templates/validation-ux.md

### Files Modified
- forge-cli/pkg/task/add.go
- forge-cli/pkg/task/autogen.go
- forge-cli/pkg/task/autogen_test.go
- forge-cli/pkg/task/record.go
- forge-cli/pkg/prompt/prompt.go
- forge-cli/pkg/prompt/prompt_test.go
- forge-cli/internal/cmd/quality_gate.go
- forge-cli/internal/cmd/task/add.go

### Key Decisions
- Merged pkg/template/ code into pkg/task/tasktemplate.go using autogenTemplateFS instead of separate embed FS
- Renamed TaskTemplateData to TemplateData to avoid stutter (task.TaskTemplateData -> task.TemplateData)
- Record templates renamed without record- prefix (e.g. record-coding.md -> coding.md in records/)

## Test Results
- **Tests Executed**: Yes
- **Passed**: 10
- **Failed**: 0
- **Coverage**: 86.1%

## Acceptance Criteria
- [x] pkg/template/ package deleted, templates and code merged into pkg/task/
- [x] pkg/task/data/ split into templates/ (14 files) and records/ (6 files, no record- prefix)
- [x] pkg/prompt/data/ renamed to pkg/prompt/templates/
- [x] All go:embed paths, templatePath() functions and import paths updated
- [x] go build ./... passes with no compile errors

## Notes
无
