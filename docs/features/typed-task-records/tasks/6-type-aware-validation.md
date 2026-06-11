---
id: "6"
title: "Type-aware validation in validateRecordData()"
priority: "P0"
estimated_time: "1h"
dependencies: ["1"]
scope: "backend"
breaking: false
type: "coding.enhancement"
mainSession: false
---

# 6: Type-aware validation in validateRecordData()

## Description

Make `validateRecordData()` type-aware using `CategoryForType()`. Currently, any `completed` task with zero test counts triggers a hard exit — even doc tasks that legitimately have no tests. Per-category validation rules:

| Category | Test Evidence Check | testsFailed Auto-downgrade |
|----------|-------------------|---------------------------|
| coding | Required (unchanged) | Active (unchanged) |
| doc | Skipped entirely | Skipped |
| test | Skipped (test tasks produce cases, not pass/fail) | Skipped |
| validation | Skipped | Skipped |
| gate | Skipped | Skipped |

## Reference Files
- `docs/proposals/typed-task-records/proposal.md` — Source proposal
- `forge-cli/internal/cmd/task/submit.go` — `validateRecordData()` (lines 311-365), `doSubmit()` (lines 103-235)
- `forge-cli/pkg/task/category.go` — `CategoryForType()` (from task 1)

## Acceptance Criteria
- [ ] `validateRecordData` accepts task type parameter: `validateRecordData(rd *task.RecordData, taskType string)`
- [ ] Doc-type tasks with `testsPassed=0, testsFailed=0, coverage=0` pass validation (no `ErrNoTestEvidence`)
- [ ] Coding-type tasks with same values still fail with `ErrNoTestEvidence` (unchanged)
- [ ] Doc-type tasks with `testsFailed > 0` are NOT auto-downgraded to `blocked`
- [ ] Coding-type tasks with `testsFailed > 0` still auto-downgrade (unchanged)
- [ ] `summary` remains hard-required for ALL task types
- [ ] `doSubmit()` updated to pass `t.Type` to `validateRecordData()`
- [ ] Unit tests for validation behavior per category

## Hard Rules
- Do NOT remove test evidence check — only skip for non-coding categories
- The `IsTestableType` check in `doSubmit` (line 133-137) must stay consistent with this change
- `keyDecisions` and `acceptanceCriteria` warnings remain for all completed tasks regardless of category

## Implementation Notes
- Only `coding` category needs test evidence validation. All other categories skip it.
- The simplest implementation: `if CategoryForType(taskType) != CategoryCoding { skip test checks }`
