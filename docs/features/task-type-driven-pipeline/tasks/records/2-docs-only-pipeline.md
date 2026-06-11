---
status: "completed"
started: "2026-05-14 23:01"
completed: "2026-05-14 23:19"
time_spent: "~18m"
---

# Task Record: 2 Implement docs-only detection and conditional pipeline in BuildIndex

## Summary
Implement docs-only detection and conditional pipeline in BuildIndex. Added isDocsOnlyFeature() to detect when all business tasks are non-implementation/non-fix. For docs-only features: skip stage-gate generation, skip test task generation, generate T-eval-doc evaluation task instead. Added hard error when business tasks from .md files have missing/empty type field after InferType. Added GetDocEvalTask() and ResolveDocEvalDep() to testgen.go. Updated existing tests to include type fields in task .md files.

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/task/build.go
- forge-cli/pkg/task/testgen.go
- forge-cli/pkg/task/build_test.go
- forge-cli/pkg/task/testgen_test.go
- forge-cli/scripts/version.txt

### Key Decisions
- isDocsOnlyFeature() excludes auto-generated tasks (gates, summaries, T-test/T-quick, T-eval-doc) from the docs-only check, only examining business tasks
- Hard error for missing type only applies to tasks read from .md files (existingKeys), not auto-generated tasks that get type from InferType
- fix- and disc- prefixed tasks are treated as business tasks (not auto-generated) since they modify code and should trigger the full pipeline
- ResolveDocEvalDep uses lexicographic ID comparison to find the last business task, handling both simple numeric IDs and dotted phase IDs
- T-eval-doc is only generated when there is at least one business task (hasBusinessTasks check prevents empty-dir edge case)
- writeTaskMD test helper updated to auto-set type field via InferType with fallback to TypeImplementation, ensuring all existing tests produce valid task files

## Test Results
- **Tests Executed**: Yes
- **Passed**: 163
- **Failed**: 0
- **Coverage**: 89.8%

## Acceptance Criteria
- [x] isDocsOnlyFeature(tasks) returns true only when ALL business tasks have type != implementation AND type != fix
- [x] Business task with empty type (after InferType) causes BuildIndex to return a hard error naming the specific file
- [x] Docs-only features skip stage-gate generation (step 6.5 in BuildIndex)
- [x] Docs-only features skip test task generation (step 7 in BuildIndex)
- [x] Docs-only features generate exactly one T-eval-doc task with: ID T-eval-doc, type doc-evaluation, noTest true, dependency on last business task
- [x] Features with any implementation or fix tasks behave identically to current behavior (gates + tests generated)
- [x] Mixed features (implementation + documentation) are treated as code features (full pipeline)
- [x] Table-driven tests in build_test.go cover: docs-only, code feature, mixed feature, missing type error
- [x] Table-driven tests in testgen_test.go cover GetDocEvalTask() output

## Notes
Version bumped from 3.6.0 to 3.7.0 (minor: new feature). Pre-existing test failures in pkg/project and internal/cmd are unrelated to this change.
