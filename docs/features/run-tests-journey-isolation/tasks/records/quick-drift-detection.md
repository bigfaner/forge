---
status: "completed"
started: "2026-05-27 00:45"
completed: "2026-05-27 00:49"
time_spent: "~4m"
---

# Task Record: T-quick-doc-drift Detect Spec Drift

## Summary
Detected and fixed spec drift in surface-orchestration rules. 3 rules drifted (BIZ-surface-orchestration-001/002/003), all auto-fixed. Remaining specs verified current.

## Changes

### Files Created
无

### Files Modified
- docs/business-rules/surface-orchestration.md

### Key Decisions
无

## Document Metrics
3 specs drifted, 3 auto-fixed; 20 specs verified current

## Referenced Documents
- docs/business-rules/error-reporting.md
- docs/business-rules/quality-gate.md
- docs/business-rules/surface-orchestration.md
- docs/business-rules/task-lifecycle.md
- docs/conventions/code-structure.md
- docs/conventions/dispatcher-quality.md
- docs/conventions/error-handling.md
- docs/conventions/forge-cli-reference.md
- docs/conventions/forge-distribution.md
- docs/conventions/prompt-template-hierarchy.md
- docs/conventions/skill-self-containment.md
- docs/conventions/skill-structure.md
- docs/conventions/surface-cli.md
- docs/conventions/surface-rules.md
- docs/conventions/testing/index.md

## Review Status
drift fixed

## Acceptance Criteria
- [x] Detect spec drift via git diff scope narrowing
- [x] Fix drifted specs and commit with [auto-specs] tag

## Notes
Drift caused by journey isolation changes: CLI/TUI no longer have build/dev steps; probe retry interval changed from 30s to 5s; sequence table updated to reflect per-journey test loops.
