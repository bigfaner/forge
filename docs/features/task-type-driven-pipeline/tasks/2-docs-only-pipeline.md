---
id: "2"
title: "Implement docs-only detection and conditional pipeline in BuildIndex"
priority: "P1"
estimated_time: "2h"
dependencies: ["1"]
scope: "backend"
breaking: false
type: "implementation"
mainSession: false
---

# 2: Implement docs-only detection and conditional pipeline in BuildIndex

## Description
Core logic for the task-type-driven pipeline. After scanning business tasks, detect whether the feature is docs-only (all tasks are `documentation`/`doc-evaluation` type, none are `implementation`/`fix`). For docs-only features: skip stage-gate generation, skip test task generation, and generate `T-eval-doc` instead. Also add a hard error when any business task has a missing/empty type field.

## Reference Files
- `docs/proposals/task-type-driven-pipeline/proposal.md` — Source proposal (D1, D3, D5)

## Affected Files

### Create
| File | Description |
|------|-------------|
| (none) | |

### Modify
| File | Changes |
|------|---------|
| `forge-cli/pkg/task/build.go` | Add `isDocsOnlyFeature()` func; add hard error on missing type for business tasks; conditionally skip stage-gate generation for docs-only; conditionally skip test tasks for docs-only; generate T-eval-doc for docs-only features |
| `forge-cli/pkg/task/testgen.go` | Add `GetDocEvalTask()` function returning a `TestTaskDef` for `T-eval-doc`; add `ResolveDocEvalDep()` to set dependency on last business task |

### Delete
| File | Reason |
|------|--------|
| (none) | |

## Acceptance Criteria
- [ ] `isDocsOnlyFeature(tasks)` returns true only when ALL business tasks have type != `implementation` AND type != `fix`
- [ ] Business task with empty `type` (after `InferType`) causes `BuildIndex` to return a hard error naming the specific file
- [ ] Docs-only features skip stage-gate generation (step 6.5 in BuildIndex)
- [ ] Docs-only features skip test task generation (step 7 in BuildIndex)
- [ ] Docs-only features generate exactly one `T-eval-doc` task with: ID `T-eval-doc`, type `doc-evaluation`, noTest true, dependency on last business task
- [ ] Features with any `implementation` or `fix` tasks behave identically to current behavior (gates + tests generated)
- [ ] Mixed features (implementation + documentation) are treated as code features (full pipeline)
- [ ] Table-driven tests in `build_test.go` cover: docs-only, code feature, mixed feature, missing type error
- [ ] Table-driven tests in `testgen_test.go` cover `GetDocEvalTask()` output

## Implementation Notes
- Detection timing: step 5.5 — after scanning .md files (step 5), before stage-gate generation (step 6.5). This means `isDocsOnlyFeature()` operates on the `existingKeys` / indexed tasks from step 5.
- Business task identification: a task is "business" if its ID does NOT match test-task patterns (`T-` prefix), gate/summary suffixes, or fix-task prefix. Reuse the `isTestTaskID` helper and add a check for gate/summary/fix prefixes.
- The hard error for missing type should only apply to business tasks read from .md files, not to auto-generated tasks (gates, summaries, test tasks) which get their type from `InferType`.
- `GetDocEvalTask()` should follow the same pattern as `GetBreakdownTestTasks` / `GetQuickTestTasks` — return a `TestTaskDef` that `BuildIndex` can use to generate the .md file and index entry.
