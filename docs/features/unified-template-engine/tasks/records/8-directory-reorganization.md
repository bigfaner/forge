---
status: "completed"
started: "2026-05-28 02:35"
completed: "2026-05-28 02:40"
time_spent: "~5m"
---

# Task Record: 8 目录重组：合并、拆分、重命名

## Summary
Directory reorganization: merged pkg/template/ into pkg/task/, split pkg/task/data/ into pkg/task/templates/ (14 files) and pkg/task/records/ (6 files with record- prefix removed), renamed pkg/prompt/data/ to pkg/prompt/templates/, updated all //go:embed paths and import references, deleted empty pkg/template/ package.

## Changes

### Files Created
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
- forge-cli/pkg/task/tasktemplate.go
- forge-cli/pkg/task/tasktemplate_test.go

### Files Modified
- forge-cli/internal/cmd/quality_gate.go
- forge-cli/internal/cmd/task/add.go
- forge-cli/pkg/prompt/prompt.go
- forge-cli/pkg/prompt/prompt_test.go
- forge-cli/pkg/task/add.go
- forge-cli/pkg/task/autogen.go
- forge-cli/pkg/task/autogen_test.go
- forge-cli/pkg/task/record.go

### Key Decisions
- Merged pkg/template/ code logic into pkg/task/tasktemplate.go rather than keeping a separate package
- Renamed record files by dropping record- prefix (e.g. record-coding.md -> coding.md) as specified in task

## Test Results
- **Tests Executed**: Yes
- **Passed**: 28
- **Failed**: 0
- **Coverage**: 0.0%

## Acceptance Criteria
- [x] pkg/template/ package deleted, templates moved to pkg/task/templates/, code merged to pkg/task/
- [x] pkg/task/data/ split into pkg/task/templates/ (14 files) and pkg/task/records/ (6 files, no record- prefix)
- [x] pkg/prompt/data/ renamed to pkg/prompt/templates/
- [x] All //go:embed paths, templatePath() functions and import paths updated
- [x] go build ./... passes with no compile errors

## Notes
Fix-record recovery task. Original implementation was in commit a494e3df. All 28 test packages pass, lint clean, fmt clean.
