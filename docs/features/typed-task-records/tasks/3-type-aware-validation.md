---
id: "3"
title: "Type-aware validation in validateRecordData()"
priority: "P0"
estimated_time: "1h"
dependencies: ["1"]
scope: "backend"
breaking: false
type: "coding.enhancement"
mainSession: false
---

# 3: Type-aware validation in validateRecordData()

## Description

Make `validateRecordData()` in `submit.go` type-aware so doc tasks are not penalized for missing test evidence. Currently, any `completed` task with `coverage >= 0` and zero test counts triggers a hard exit — even doc tasks that legitimately have no tests.

Changes:
- Use `CategoryForType()` to determine validation tier
- Doc tasks: skip test evidence check (`ErrNoTestEvidence`) entirely
- Doc tasks: skip `testsFailed` auto-downgrade (already handled by `IsTestableType` in `doSubmit`, but validateRecordData should also be consistent)
- Coding tasks: validation unchanged

## Reference Files
- `docs/proposals/typed-task-records/proposal.md` — Source proposal
- `forge-cli/internal/cmd/task/submit.go` — `validateRecordData()` (lines 311-365), `doSubmit()` (lines 103-235)
- `forge-cli/pkg/task/category.go` — `CategoryForType()` (from task 1)

## Acceptance Criteria
- [ ] `forge task submit` for a doc-type task with `testsPassed=0, testsFailed=0, coverage=0` succeeds (no `ErrNoTestEvidence`)
- [ ] `forge task submit` for a coding-type task with `testsPassed=0, testsFailed=0, coverage=0` still fails with `ErrNoTestEvidence`
- [ ] Doc tasks with `testsFailed > 0` are NOT auto-downgraded to `blocked` (irrelevant metric for doc tasks)
- [ ] Coding tasks with `testsFailed > 0` still auto-downgrade to `blocked` (unchanged behavior)
- [ ] `summary` remains hard-required for ALL task types
- [ ] Unit tests for validation behavior per category

## Hard Rules
- Do NOT remove the test evidence check — only skip it for doc-category tasks
- The `IsTestableType` check in `doSubmit` (line 133-137) and this validation must stay consistent
- Must accept the task type as a parameter: change `validateRecordData(rd)` to `validateRecordData(rd, taskType)` or pass the full `Task` struct

## Implementation Notes
- The function signature needs the task type. Simplest: add `taskType string` parameter, call as `validateRecordData(rd, t.Type)` from `doSubmit`.
- `doSubmit` already has access to the task type via `t.Type`.
- The `IsTestableType` guard at line 133-137 auto-sets `coverage=-1.0` for non-coding tasks, which happens before `validateRecordData`. But the validation function should also independently check the type for clarity and correctness.
