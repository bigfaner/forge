---
status: "completed"
started: "2026-06-09 18:30"
completed: "2026-06-09 18:39"
time_spent: "~9m"
---

# Task Record: T-quick-doc-drift Detect Spec Drift

## Summary
Drift-only mode: scanned all project-level spec files (4 business-rules, 15 conventions). Found 1 drifted spec in forge-distribution.md -- test pipeline section did not document CondHasProtocolSurfaceTask conditional skip for gen-contracts/eval-contract. Fixed by updating the pipeline description. Regenerated vocabulary index with updated counts.

## Changes

### Files Created
无

### Files Modified
- docs/conventions/forge-distribution.md
- docs/.vocabulary.md

### Key Decisions
无

## Document Metrics
specs scanned: 19, drifted: 1, orphaned: 0, current: 18, implicit new: 0

## Referenced Documents
- docs/conventions/forge-distribution.md
- docs/business-rules/task-lifecycle.md
- docs/business-rules/quality-gate.md
- docs/business-rules/surface-orchestration.md
- docs/business-rules/error-reporting.md
- docs/conventions/surface-rules.md
- docs/conventions/surface-cli.md
- docs/conventions/naming.md
- docs/conventions/enum-constants.md
- docs/conventions/constants.md
- docs/conventions/code-structure.md
- docs/conventions/dead-code.md
- docs/conventions/package-organization.md
- docs/conventions/dispatcher-quality.md
- docs/conventions/error-handling.md
- docs/conventions/forge-cli-reference.md
- docs/conventions/skill-structure.md
- docs/conventions/skill-self-containment.md
- docs/conventions/prompt-template-hierarchy.md
- docs/proposals/skip-contracts-web-mobile/proposal.md

## Review Status
final

## Acceptance Criteria
- [x] All acceptance criteria met

## Notes
Drift-only mode (no PRD/design files). Scoped by git diff main...HEAD to identify changed files, then validated specs whose domains overlap with changed code. Only drift found: forge-distribution.md test pipeline section missing CondHasProtocolSurfaceTask conditional skip documentation.
