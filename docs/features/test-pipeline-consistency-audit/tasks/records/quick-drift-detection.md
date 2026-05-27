---
status: "completed"
started: "2026-05-27 20:24"
completed: "2026-05-27 20:29"
time_spent: "~5m"
---

# Task Record: T-quick-doc-drift Detect Spec Drift

## Summary
Detected and fixed spec drift: checked 15 spec files across docs/business-rules/ and docs/conventions/, found 1 drifted rule (BIZ-task-lifecycle-003 stated 13 system types but code has 12 — removed non-existent test.verify-regression), auto-fixed and committed with [auto-specs] tag

## Changes

### Files Created
无

### Files Modified
- docs/business-rules/task-lifecycle.md

### Key Decisions
无

## Document Metrics
15 specs checked, 1 drifted, 1 auto-fixed

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
- [x] Run git diff to identify files changed by this feature
- [x] Read all spec file frontmatter domains to scope drift check
- [x] Validate spec rules against current codebase
- [x] Auto-fix drifted rules and commit with [auto-specs] tag

## Notes
Drift-only mode (no PRD/design files). Validated all 15 spec files. Only drift found: BIZ-task-lifecycle-003 system type count (13 -> 12) and removal of non-existent test.verify-regression type. All other specs verified current against codebase.
