---
status: "completed"
started: "2026-05-29 14:57"
completed: "2026-05-29 15:05"
time_spent: "~8m"
---

# Task Record: T-quick-doc-drift Detect Spec Drift

## Summary
Detected and fixed 2 spec drifts: (1) SystemTypes comment in pkg/task/types.go said '13 total' but map has 12 entries (after doc.fix removal), (2) dispatcher-quality.md hardcoded coding.fix without category-based derivation rule introduced by derive-fix-task-type feature

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/task/types.go
- docs/conventions/dispatcher-quality.md

### Key Decisions
无

## Document Metrics
2 drifts found and fixed across 2 files, 5 spec files reviewed

## Referenced Documents
- docs/business-rules/task-lifecycle.md
- docs/business-rules/quality-gate.md
- docs/conventions/dispatcher-quality.md
- docs/conventions/enum-constants.md
- docs/features/derive-fix-task-type/tasks/3-update-skill-files-derivation.md

## Review Status
final

## Acceptance Criteria
- [x] git diff --name-only main...HEAD used to narrow scope to relevant specs
- [x] All spec files with domains overlapping changed files are verified against code
- [x] Drifts auto-fixed with correct content

## Notes
Drift 1: SystemTypes comment in types.go line 145 said '13 total' but commit 8e931b4c removed doc.fix leaving 12 entries. Drift 2: dispatcher-quality.md still hardcoded 'coding.fix' for all fix tasks; derive-fix-task-type feature introduced category-based derivation (doc/eval -> doc.fix, coding/test/validation/gate -> coding.fix) in plugin files but this convention was not updated. Quality-gate.md is correct as-is because quality-gate hook only deals with code-level failures (compile/fmt/lint/test).
