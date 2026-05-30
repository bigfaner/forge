---
status: "completed"
started: "2026-05-30 01:50"
completed: "2026-05-30 01:54"
time_spent: "~4m"
---

# Task Record: T-quick-doc-drift Detect Spec Drift

## Summary
Detected and fixed spec drift in docs/conventions/forge-cli-reference.md. Drift found in source file path references for feature, worktree, and forensic command groups which were refactored into subdirectories. All other 15 spec files (30+ rules) validated as current.

## Changes

### Files Created
无

### Files Modified
- docs/conventions/forge-cli-reference.md

### Key Decisions
无

## Document Metrics
1 drifted rule fixed out of 30+ total rules across 16 spec files; drift rate: ~3%

## Referenced Documents
- docs/conventions/forge-cli-reference.md
- docs/business-rules/error-reporting.md
- docs/business-rules/quality-gate.md
- docs/business-rules/surface-orchestration.md
- docs/business-rules/task-lifecycle.md
- docs/conventions/code-structure.md
- docs/conventions/dispatcher-quality.md
- docs/conventions/enum-constants.md
- docs/conventions/error-handling.md
- docs/conventions/forge-distribution.md
- docs/conventions/prompt-template-hierarchy.md
- docs/conventions/skill-self-containment.md
- docs/conventions/skill-structure.md
- docs/conventions/surface-cli.md
- docs/conventions/surface-rules.md

## Review Status
final

## Acceptance Criteria
- [x] Git diff used to identify changed files for scope narrowing
- [x] All spec files in docs/business-rules/ and docs/conventions/ validated against current code
- [x] Drifted source file paths in forge-cli-reference.md corrected
- [x] Commit includes [auto-specs] tag for traceability

## Notes
Drift-only mode (no PRD/design docs). Only forge-cli-reference.md had drift: feature.go, worktree.go, and forensic/forensic.go source paths were stale after refactoring to subdirectory structures. No orphaned or new implicit rules discovered.
