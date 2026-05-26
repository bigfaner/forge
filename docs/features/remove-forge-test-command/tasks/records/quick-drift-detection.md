---
status: "completed"
started: "2026-05-26 22:45"
completed: "2026-05-26 22:48"
time_spent: "~3m"
---

# Task Record: T-quick-doc-drift Detect Spec Drift

## Summary
Detected and fixed spec drift in forge-cli-reference.md: added 3 missing removed commands (forge test promote, run-journey, verify) to the removal table

## Changes

### Files Created
无

### Files Modified
- docs/conventions/forge-cli-reference.md

### Key Decisions
无

## Document Metrics
14 specs checked, 1 drifted, 1 auto-fixed (3 missing command entries added)

## Referenced Documents
- docs/conventions/forge-cli-reference.md
- docs/business-rules/error-reporting.md
- docs/business-rules/quality-gate.md
- docs/business-rules/surface-orchestration.md
- docs/business-rules/task-lifecycle.md
- docs/conventions/code-structure.md
- docs/conventions/dispatcher-quality.md
- docs/conventions/error-handling.md
- docs/conventions/forge-distribution.md
- docs/conventions/prompt-template-hierarchy.md
- docs/conventions/skill-self-containment.md
- docs/conventions/skill-structure.md
- docs/conventions/surface-cli.md
- docs/conventions/surface-rules.md

## Review Status
drift fixed

## Acceptance Criteria
- [x] All project-level spec files checked for drift against current codebase
- [x] Drifted specs auto-fixed and committed with [auto-specs] tag
- [x] No orphaned or stale rules remain in project-level specs

## Notes
Drift detection narrowed scope using git diff main...HEAD per task discovery strategy. Only forge-cli-reference.md had drift: missing entries for forge test promote/run-journey/verify in the removed commands table. All other 13 spec files verified current.
