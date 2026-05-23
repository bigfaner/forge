---
status: "completed"
started: "2026-05-23 23:52"
completed: "2026-05-24 00:02"
time_spent: "~10m"
---

# Task Record: 1 Rename doc.eval → doc.review and rewrite templates

## Summary
Renamed doc.eval to doc.review across entire Go CLI codebase: constants, functions, task IDs, template mappings, and test references. Rewrote autogen template and agent prompt from broken 8-dimension rubric to AC checklist + direct-fix model. Deleted old doc-eval.md files.

## Changes

### Files Created
- forge-cli/pkg/task/data/doc-review.md
- forge-cli/pkg/prompt/data/doc-review.md

### Files Modified
- forge-cli/pkg/task/types.go
- forge-cli/pkg/task/autogen.go
- forge-cli/pkg/task/build.go
- forge-cli/pkg/task/infer.go
- forge-cli/pkg/prompt/prompt.go
- forge-cli/pkg/task/types_test.go
- forge-cli/pkg/task/autogen_test.go
- forge-cli/pkg/task/build_test.go
- forge-cli/pkg/task/infer_test.go
- forge-cli/pkg/task/category_test.go
- forge-cli/pkg/prompt/prompt_test.go
- forge-cli/internal/cmd/quality_gate_test.go
- forge-cli/internal/cmd/task/validate_index_test.go
- forge-cli/tests/task-type-system/task_type_refinement_test.go

### Key Decisions
- Direct rename without alias phase since total files <= 20 and all changes are simple text replacement
- T-review-doc no longer matches T-eval- prefix in isTestTaskID, handled by explicit check in IsAutoGenTaskID
- Agent prompt changed from 8-dimension 1000-point scoring to 4-step AC checklist with direct fixes

## Test Results
- **Tests Executed**: Yes
- **Passed**: 1346
- **Failed**: 0
- **Coverage**: 90.6%

## Acceptance Criteria
- [x] TypeDocEval constant renamed to TypeDocReview with value doc.review in types.go
- [x] doc.eval removed from ValidTypes and SystemTypes, replaced with doc.review
- [x] TaskTypeRegistry entry updated: doc.review with description review documentation against acceptance criteria
- [x] GetDocEvalTask renamed to GetReviewDocTask with Key review-doc, ID T-review-doc, Title Review Documentation Quality, Type TypeDocReview
- [x] ResolveDocEvalDep renamed to ResolveReviewDocDep in autogen.go
- [x] autogenTypeToFile mapping updated: TypeDocReview to data/doc-review.md
- [x] needsDocEval renamed to needsReviewDoc in build.go
- [x] All call sites updated in build.go
- [x] IsAutoGenTaskID updated: T-eval-doc to T-review-doc
- [x] infer.go type inference updated: T-review-doc to TypeDocReview
- [x] prompt.go template mapping updated: TypeDocReview to data/doc-review.md
- [x] forge-cli/pkg/task/data/doc-eval.md renamed to doc-review.md with AC checklist content
- [x] forge-cli/pkg/prompt/data/doc-eval.md renamed to doc-review.md with AC checklist + direct fix workflow
- [x] go build passes with zero errors after all renames
- [x] No remaining references to doc.eval, DocEval, T-eval-doc, or doc-eval in Go source or templates

## Notes
category.go logic unchanged as per hard rules - strings.HasPrefix(typ, doc) already covers doc.review. All targeted tests pass: pkg/task 90.6%, pkg/prompt 90.9%, internal/cmd/task 71.0%.
