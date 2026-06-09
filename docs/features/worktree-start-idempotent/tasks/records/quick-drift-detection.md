---
status: "completed"
started: "2026-06-09 16:29"
completed: "2026-06-09 16:44"
time_spent: "~15m"
---

# Task Record: T-quick-doc-drift Detect Spec Drift

## Summary
Detected and fixed 2 spec drifts: (1) skill-structure.md line counts stale (gen-test-scripts 489->536, init-justfile 451->490, gen-journeys 428->454); (2) forge-cli-reference.md worktree start description missing idempotent behavior and --no-launch/-b flags. Regenerated vocabulary index.

## Changes

### Files Created
无

### Files Modified
- docs/conventions/skill-structure.md
- docs/conventions/forge-cli-reference.md
- docs/.vocabulary.md

### Key Decisions
无

## Document Metrics
19 spec files scanned, 2 drifts detected, 2 drifts fixed, 0 orphaned, 0 implicit new rules

## Referenced Documents
- docs/conventions/skill-structure.md
- docs/conventions/forge-cli-reference.md
- docs/conventions/forge-distribution.md
- docs/conventions/skill-self-containment.md
- docs/conventions/code-structure.md
- docs/conventions/constants.md
- docs/conventions/dead-code.md
- docs/conventions/dispatcher-quality.md
- docs/conventions/enum-constants.md
- docs/conventions/error-handling.md
- docs/conventions/naming.md
- docs/conventions/package-organization.md
- docs/conventions/prompt-template-hierarchy.md
- docs/conventions/surface-cli.md
- docs/conventions/surface-rules.md
- docs/business-rules/task-lifecycle.md
- docs/business-rules/quality-gate.md
- docs/business-rules/error-reporting.md
- docs/business-rules/surface-orchestration.md

## Review Status
final

## Acceptance Criteria
- [x] All acceptance criteria met

## Notes
Drift-only mode (no PRD/design files). Used git diff to narrow scope to relevant spec files. All 19 spec files validated. Only skill-structure.md and forge-cli-reference.md had drift.
