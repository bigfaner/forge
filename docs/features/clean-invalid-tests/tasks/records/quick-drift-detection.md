---
status: "completed"
started: "2026-05-27 11:16"
completed: "2026-05-27 11:30"
time_spent: "~14m"
---

# Task Record: T-quick-doc-drift Detect Spec Drift

## Summary
Detected and auto-fixed spec drift in 2 files referencing deleted tests/e2e/ directory

## Changes

### Files Created
无

### Files Modified
- docs/conventions/testing/go.md
- docs/conventions/forge-distribution.md

### Key Decisions
无

## Document Metrics
12 specs checked, 2 drifted, 2 auto-fixed

## Referenced Documents
- docs/conventions/testing/go.md
- docs/conventions/testing/ginkgo.md
- docs/conventions/testing/vitest.md
- docs/conventions/forge-distribution.md
- docs/business-rules/error-reporting.md
- docs/business-rules/quality-gate.md
- docs/business-rules/surface-orchestration.md
- docs/business-rules/task-lifecycle.md
- docs/conventions/dispatcher-quality.md
- docs/conventions/error-handling.md
- docs/conventions/surface-cli.md
- docs/conventions/surface-rules.md

## Review Status
drift fixed

## Acceptance Criteria
- [x] All spec files in docs/business-rules/ and docs/conventions/ checked for drift against codebase
- [x] Drifted specs auto-fixed and committed with [auto-specs] tag
- [x] No orphaned rules found

## Notes
Used git diff --name-only main...HEAD to narrow scope. Drift limited to tests/e2e/ path references after directory removal. vitest.md and ginkgo.md kept unchanged (template conventions, not project-specific paths).
