---
status: "completed"
started: "2026-05-30 22:45"
completed: "2026-05-30 22:51"
time_spent: "~6m"
---

# Task Record: T-quick-doc-drift Detect Spec Drift

## Summary
Spec drift detection for forge-cli-codebase-standards feature. Verified all 19 spec files against codebase. Found and fixed 3 drifts: (1) package-organization.md subpackage count 7->8, (2) code-structure.md CS-1/CS-5 deviation tables stale after fixes, (3) dead-code.md DC-1/DC-5 deviation tables stale after fixes.

## Changes

### Files Created
无

### Files Modified
- docs/conventions/package-organization.md
- docs/conventions/code-structure.md
- docs/conventions/dead-code.md

### Key Decisions
无

## Document Metrics
19 specs verified, 3 drifts found and auto-fixed, 0 specs require manual review

## Referenced Documents
- docs/conventions/forge-cli-reference.md
- docs/conventions/enum-constants.md
- docs/conventions/constants.md
- docs/conventions/naming.md
- docs/conventions/forge-distribution.md
- docs/conventions/skill-structure.md
- docs/conventions/skill-self-containment.md
- docs/conventions/dispatcher-quality.md
- docs/conventions/error-handling.md
- docs/conventions/surface-cli.md
- docs/conventions/surface-rules.md
- docs/conventions/prompt-template-hierarchy.md
- docs/business-rules/task-lifecycle.md
- docs/business-rules/quality-gate.md
- docs/business-rules/error-reporting.md
- docs/business-rules/surface-orchestration.md

## Review Status
final

## Acceptance Criteria
- [x] Run git diff --name-only main...HEAD to identify changed files
- [x] List all spec files in docs/business-rules/ and docs/conventions/
- [x] For each spec file, read its domains frontmatter
- [x] Only verify specs whose domains overlap with changed files
- [x] Auto-fix drifted specs and commit with [auto-specs] tag

## Notes
Drift details: (1) package-organization.md listed 7 subpackages but actual count is 8 (docs/ was omitted from count despite being listed). (2) CS-1 Debugf dedup and CS-5 Scope removal were completed in earlier tasks but deviation tables not updated. (3) Same DC-1/DC-5 stale entries in dead-code.md. All other specs (forge-cli-reference.md commands, enum-constants.md SystemTypes, task-lifecycle.md state machine, quality-gate.md pipeline, error-reporting.md exit codes, surface-orchestration.md surface types, constants.md deviation line numbers, naming.md conventions) verified accurate against current code.
