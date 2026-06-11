---
status: "completed"
started: "2026-05-28 02:40"
completed: "2026-05-28 02:47"
time_spent: "~7m"
---

# Task Record: T-quick-doc-drift Detect Spec Drift

## Summary
Spec drift detection for unified-template-engine branch: verified 15 spec files against current codebase, found and auto-fixed 2 drifts in forge-distribution.md (outdated version example, stale forge test promote reference)

## Changes

### Files Created
无

### Files Modified
- docs/conventions/forge-distribution.md

### Key Decisions
无

## Document Metrics
15 specs checked, 2 drifted, 2 auto-fixed

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
- [x] Used git diff --name-only main...HEAD to narrow scope
- [x] Verified all spec files whose domains overlap with changed files
- [x] Auto-fixed drifted specs and committed with [auto-specs] tag
- [x] No spec files drifted without fix

## Notes
Drift details: (1) forge-distribution.md version example outdated from 3.0.0-rc.23 to 3.0.0-rc.28; (2) forge-distribution.md referenced removed command forge test promote instead of /run-tests skill. All other 13 spec files verified current.
