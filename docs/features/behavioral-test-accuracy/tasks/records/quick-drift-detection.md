---
status: "completed"
started: "2026-06-08 18:22"
completed: "2026-06-08 18:31"
time_spent: "~9m"
---

# Task Record: T-quick-doc-drift Detect Spec Drift

## Summary
Detected and fixed spec drift: fallbackSortPriority location updated from list.go to list_sort.go in constants.md and enum-constants.md. Regenerated vocabulary index with updated counts (10 decisions, 146 lessons, 18 conventions, 4 business-rules). No other drift found across 14 spec files checked.

## Changes

### Files Created
无

### Files Modified
- docs/conventions/constants.md
- docs/conventions/enum-constants.md
- docs/.vocabulary.md

### Key Decisions
无

## Document Metrics
specs_scanned: 22, drifted: 2 (constants.md, enum-constants.md), fixed: 2, orphaned: 0

## Referenced Documents
- docs/conventions/constants.md
- docs/conventions/enum-constants.md
- docs/conventions/naming.md
- docs/conventions/code-structure.md
- docs/conventions/package-organization.md
- docs/conventions/error-handling.md
- docs/conventions/forge-cli-reference.md
- docs/conventions/dead-code.md
- docs/conventions/dispatcher-quality.md
- docs/conventions/forge-distribution.md
- docs/conventions/prompt-template-hierarchy.md
- docs/conventions/skill-self-containment.md
- docs/conventions/skill-structure.md
- docs/conventions/surface-cli.md
- docs/conventions/surface-rules.md
- docs/business-rules/error-reporting.md
- docs/business-rules/quality-gate.md
- docs/business-rules/surface-orchestration.md
- docs/business-rules/task-lifecycle.md
- docs/conventions/testing/index.md
- docs/conventions/testing/cli/index.md
- docs/conventions/testing/cli/core.md

## Review Status
final

## Acceptance Criteria
- [x] All acceptance criteria met

## Notes
Drift-only mode (no PRD/design files). Used git diff to narrow scope to 14 spec files with domain overlap. Only drift found: fallbackSortPriority file location changed from list.go to list_sort.go during a prior refactor.
