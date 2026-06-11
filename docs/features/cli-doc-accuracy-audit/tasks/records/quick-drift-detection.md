---
status: "completed"
started: "2026-06-07 22:18"
completed: "2026-06-07 22:21"
time_spent: "~3m"
---

# Task Record: T-quick-doc-drift Detect Spec Drift

## Summary
Detected and fixed 2 spec drifts in forge-cli-reference.md: (1) forge cleanup description only mentioned 'completed' but code handles completed/blocked/suspended/rejected; (2) forge worktree status description omitted the UNPUSHED field. Verified 19 spec files, narrowed scope to 2 overlapping specs via git diff + domains matching.

## Changes

### Files Created
无

### Files Modified
- docs/conventions/forge-cli-reference.md

### Key Decisions
无

## Document Metrics
specs_scanned: 19, specs_overlapping: 2, drifts_found: 2, drifts_fixed: 2

## Referenced Documents
- docs/conventions/forge-cli-reference.md
- docs/conventions/forge-distribution.md
- plugins/forge/hooks/guide.md

## Review Status
final

## Acceptance Criteria
- [x] git diff --name-only main...HEAD executed, changed file scope identified
- [x] Only specs with domain overlap checked, no omissions, no extra scans
- [x] Drifts auto-fixed and committed, or confirmed no drift

## Notes
forge-distribution.md was checked but guide.md content changes did not affect distribution structure. Other specs (surface-cli, testing/*) had no domain overlap with changed Go source files.
