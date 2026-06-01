---
status: "completed"
started: "2026-06-01 21:35"
completed: "2026-06-01 21:40"
time_spent: "~5m"
---

# Task Record: T-quick-doc-drift Detect Spec Drift

## Summary
Detected and fixed spec drift in 2 project-level spec files: forge-cli-reference.md (added forge upgrade command) and forge-distribution.md (updated distribution model to reflect CLI binary as entry point)

## Changes

### Files Created
无

### Files Modified
- docs/conventions/forge-cli-reference.md
- docs/conventions/forge-distribution.md

### Key Decisions
无

## Document Metrics
2 drifted specs fixed, 0 orphaned, 0 implicit new rules. Scoped to forge-cli-reference.md and forge-distribution.md (only specs whose domains overlap with feature changes)

## Referenced Documents
- docs/proposals/forge-cli-binary-distribution/proposal.md
- docs/conventions/forge-cli-reference.md
- docs/conventions/forge-distribution.md

## Review Status
final

## Acceptance Criteria
- [x] git diff --name-only main...HEAD executed to identify changed files
- [x] Only spec files whose domains overlap with changed files are verified
- [x] Drift detected -> auto-fixed and committed with [auto-specs] tag

## Notes
Drift-only mode (no PRD/design files in feature). Used git diff to narrow scope to 2 relevant spec files out of 20 total. forge-cli-reference.md was missing forge upgrade command. forge-distribution.md still described old model (Plugin as entry point) instead of new model (CLI binary as entry point).
