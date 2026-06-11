---
status: "completed"
started: "2026-05-29 11:39"
completed: "2026-05-29 11:46"
time_spent: "~7m"
---

# Task Record: T-quick-doc-drift Detect Spec Drift

## Summary
Detected and fixed spec drift in ARCHITECTURE.md: removed references to deleted /eval-consistency and /eval-harness skills, removed task-cli as standalone subsystem (merged into forge CLI), updated /record-decision and /learn-lesson references to /learn

## Changes

### Files Created
无

### Files Modified
- docs/ARCHITECTURE.md

### Key Decisions
无

## Document Metrics
4 drift items found and fixed: 2 deleted skills, 1 removed subsystem, 2 outdated command references

## Referenced Documents
- docs/business-rules/quality-gate.md
- docs/conventions/enum-constants.md
- docs/conventions/forge-cli-reference.md
- docs/conventions/forge-distribution.md

## Review Status
final

## Acceptance Criteria
- [x] Run git diff to identify changed files
- [x] Only verify specs whose domains overlap with changed files
- [x] Skip specs with no overlap
- [x] Auto-fix drifted specs
- [x] Commit with [auto-specs] tag

## Notes
Scanned docs/business-rules/ (4 files) and docs/conventions/ (12 files). Only ARCHITECTURE.md had drift: eval-consistency and eval-harness skills deleted but still listed in evaluation table; task-cli removed as separate project but still listed as subsystem; /record-decision and /learn-lesson replaced by /learn but still referenced in directory annotations. Other spec files had no drift.
