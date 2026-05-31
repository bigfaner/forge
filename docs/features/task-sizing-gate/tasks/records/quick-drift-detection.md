---
status: "completed"
started: "2026-05-31 15:13"
completed: "2026-05-31 15:22"
time_spent: "~9m"
---

# Task Record: T-quick-doc-drift Detect Spec Drift

## Summary
Detected and fixed spec drift in 4 files: task-lifecycle.md (validate-index→validate), code-structure.md (CS-2/3/4 marked fixed), dead-code.md (DC-2/3/4 marked fixed), forge-cli-reference.md (validate-index→validate). Added 2 implicit rules: BIZ-task-lifecycle-005 (Task Sizing Constraints) and BIZ-task-lifecycle-006 (Task Complexity Classification). Regenerated vocabulary index.

## Changes

### Files Created
无

### Files Modified
- docs/business-rules/task-lifecycle.md
- docs/conventions/code-structure.md
- docs/conventions/dead-code.md
- docs/conventions/forge-cli-reference.md
- docs/.vocabulary.md

### Key Decisions
无

## Document Metrics
4 drifted specs fixed, 2 implicit rules added, 5 domains updated, vocabulary regenerated

## Referenced Documents
- docs/proposals/task-sizing-gate/proposal.md
- docs/features/task-sizing-gate/tasks/quick-drift-detection.md
- forge-cli/pkg/task/types.go
- forge-cli/internal/cmd/task/validate.go
- forge-cli/internal/cmd/task/claim.go

## Review Status
final

## Acceptance Criteria
- [x] Run git diff to identify files changed by task-sizing-gate feature
- [x] List all spec files in docs/business-rules/ and docs/conventions/
- [x] Only verify specs whose domains overlap with changed files
- [x] Auto-fix drifted specs and commit with [auto-specs] tag

## Notes
Drift findings: (1) validate-index.go renamed to validate.go + AC count validation added -- forge-cli-reference.md and task-lifecycle.md updated. (2) All 3 test-bridge aliases (getTaskPhase, checkExistingTaskState, compareVersionIDs) eliminated from claim.go -- dead-code.md DC-2/3/4 and code-structure.md CS-2/3/4 marked as fixed. (3) Task sizing rules (AC 1-6, audit step, complexity field) existed only in skill files, not in project-level specs -- promoted to BIZ-task-lifecycle-005/006.
