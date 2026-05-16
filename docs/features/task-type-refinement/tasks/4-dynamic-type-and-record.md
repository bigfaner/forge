---
id: "4"
title: "Dynamic fix task type by failure step and Type Reclassification in records"
priority: "P1"
estimated_time: "1h"
dependencies: ["1"]
scope: "backend"
breaking: false
type: "implementation"
mainSession: false
---

# 4: Dynamic fix task type by failure step and Type Reclassification in records

## Description

Two related changes to the task lifecycle: (1) When quality gate fails and a fix task is dynamically created, the type is deterministically set based on which step failed — compile/test → `fix`, fmt/lint → `cleanup`. (2) The record format gains a Type Reclassification block that documents when an executor changes a task's type during execution.

## Reference Files
- `docs/proposals/task-type-refinement/proposal.md` — Source proposal (D4, D5)
- `forge-cli/internal/cmd/quality_gate.go` — `addFixTask()` function (lines 324-401)
- `forge-cli/internal/cmd/submit.go` — `RecordData` struct, `fillRecordTemplate()` function
- `forge-cli/pkg/task/types.go` — `RecordData` struct (lines 201-214)

## Acceptance Criteria
- [ ] `addFixTask()` sets type to `fix` when failure is from compile or test step
- [ ] `addFixTask()` sets type to `cleanup` when failure is from fmt or lint step
- [ ] `RecordData` struct has optional `TypeReclassification` field (original type, actual type, reason)
- [ ] `fillRecordTemplate()` renders "## Type Reclassification" block only when `TypeReclassification` is non-nil
- [ ] Type mapping follows proposal D4 table exactly

## Hard Rules
- The failure step → type mapping must be deterministic, no heuristics:
  - compile failure → `TypeFix`
  - fmt failure → `TypeCleanup`
  - lint failure → `TypeCleanup`
  - unit test failure → `TypeFix`
  - e2e test failure → `TypeFix`

## Implementation Notes
- `addFixTask()` receives the step info from the calling context. Trace how the failure step name flows into `addFixTask()` to determine where to inject the type decision.
- For Type Reclassification in records: this is a rendering concern in `fillRecordTemplate()`. The struct field should be optional (pointer or zero-value check) so existing records are unaffected.
