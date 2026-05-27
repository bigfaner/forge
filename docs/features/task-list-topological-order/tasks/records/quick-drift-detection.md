---
status: "completed"
started: "2026-05-28 01:03"
completed: "2026-05-28 01:07"
time_spent: "~4m"
---

# Task Record: T-quick-doc-drift Detect Spec Drift

## Summary
Detected and fixed spec drift across project-level spec files. 1 drifted rule auto-fixed (forge-cli-reference.md missing --sort/--tree flags), 1 implicit new rule added (BIZ-task-lifecycle-004: topological ordering). Vocabulary index regenerated.

## Changes

### Files Created
无

### Files Modified
- docs/conventions/forge-cli-reference.md
- docs/business-rules/task-lifecycle.md
- docs/.vocabulary.md

### Key Decisions
无

## Document Metrics
15 specs checked, 1 drifted (auto-fixed), 1 implicit new rule added, 0 orphaned

## Referenced Documents
- docs/business-rules/task-lifecycle.md
- docs/business-rules/quality-gate.md
- docs/business-rules/surface-orchestration.md
- docs/business-rules/error-reporting.md
- docs/conventions/forge-cli-reference.md
- docs/conventions/forge-distribution.md
- docs/conventions/error-handling.md
- docs/conventions/surface-cli.md
- docs/conventions/surface-rules.md
- docs/conventions/dispatcher-quality.md
- docs/conventions/code-structure.md
- docs/conventions/skill-self-containment.md
- docs/conventions/skill-structure.md
- docs/conventions/prompt-template-hierarchy.md
- docs/conventions/testing/ginkgo.md
- docs/conventions/testing/go.md
- docs/conventions/testing/index.md
- docs/conventions/testing/junit.md
- docs/conventions/testing/pytest.md
- docs/conventions/testing/rust.md
- docs/conventions/testing/vitest.md

## Review Status
drift fixed

## Acceptance Criteria
- [x] All project-level specs checked for drift against current codebase
- [x] Drifted rules auto-fixed with preserved IDs
- [x] Changes committed with [auto-specs] tag
- [x] Vocabulary index regenerated

## Notes
Drift-only mode (no PRD/design files). Used git diff to narrow scope; all 15 spec files validated. Only forge-cli-reference.md had drifted (missing --sort/--tree flags for forge task list). Topological ordering (BIZ-task-lifecycle-004) was an implicit new rule discovered during drift scanning.
