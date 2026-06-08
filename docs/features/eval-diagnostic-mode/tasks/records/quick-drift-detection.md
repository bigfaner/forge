---
status: "completed"
started: "2026-06-08 12:54"
completed: "2026-06-08 12:59"
time_spent: "~5m"
---

# Task Record: T-quick-doc-drift Detect Spec Drift

## Summary
Detected and fixed 1 spec drift: prompt-template-hierarchy.md incorrectly stated TASK-CONSTRAINTS was 'not yet used in any shipped template', but 7 prompt templates now use it. Updated description to reflect actual usage. All other spec files (task-lifecycle, quality-gate, dispatcher-quality, forge-distribution, forge-cli-reference) verified consistent with code.

## Changes

### Files Created
无

### Files Modified
- docs/conventions/prompt-template-hierarchy.md

### Key Decisions
无

## Document Metrics
drifts found: 1, drifts fixed: 1, specs verified clean: 7

## Referenced Documents
- docs/business-rules/task-lifecycle.md
- docs/business-rules/quality-gate.md
- docs/conventions/dispatcher-quality.md
- docs/conventions/forge-distribution.md
- docs/conventions/forge-cli-reference.md
- docs/conventions/enum-constants.md
- docs/conventions/constants.md

## Review Status
final

## Acceptance Criteria
- [x] Drift check completed (skipped if no code changes detected)

## Notes
Used git diff to scope code changes, then cross-referenced spec domains to narrow verification. Only prompt-template-hierarchy.md had drift (TASK-CONSTRAINTS usage description outdated). Task-lifecycle SystemTypes list (12 entries including eval.journey, eval.contract) matches code exactly.
