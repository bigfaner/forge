---
status: "completed"
started: "2026-05-26 22:45"
completed: "2026-05-26 22:50"
time_spent: "~5m"
---

# Task Record: T-quick-doc-drift Detect Spec Drift

## Summary
Drift detection: scanned 12 spec files against codebase, found 1 drifted rule (BIZ-task-lifecycle-003 had wrong system type count 15 vs actual 13), auto-fixed and committed

## Changes

### Files Created
无

### Files Modified
- docs/business-rules/task-lifecycle.md
- docs/.vocabulary.md

### Key Decisions
无

## Document Metrics
12 specs checked, 1 drifted, 1 auto-fixed

## Referenced Documents
- docs/business-rules/surface-orchestration.md
- docs/business-rules/task-lifecycle.md
- docs/business-rules/quality-gate.md
- docs/business-rules/error-reporting.md
- docs/conventions/surface-rules.md
- docs/conventions/surface-cli.md
- docs/conventions/code-structure.md
- docs/conventions/error-handling.md
- docs/conventions/forge-cli-reference.md
- docs/conventions/dispatcher-quality.md
- docs/conventions/forge-distribution.md
- docs/conventions/testing/index.md
- docs/reference/test-type-model.md

## Review Status
drift fixed

## Acceptance Criteria
- [x] Scoped drift detection to files changed by surface-test-type-model feature
- [x] All spec files with overlapping domains verified against codebase
- [x] Drifted rule BIZ-task-lifecycle-003 corrected (15 -> 13 base types)
- [x] Changes committed with [auto-specs] tag

## Notes
Drift in BIZ-task-lifecycle-003: spec claimed '15 total, including wildcard patterns' but code has exactly 13 base types in SystemTypes map. Surface-specific variants (e.g. test.gen-scripts.cli) are dynamically recognized by IsSystemType() stripping the last segment. Vocabulary index regenerated with updated lesson count (112 -> 115).
