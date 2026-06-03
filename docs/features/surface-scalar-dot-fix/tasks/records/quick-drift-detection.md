---
status: "completed"
started: "2026-06-03 22:53"
completed: "2026-06-03 22:58"
time_spent: "~5m"
---

# Task Record: T-quick-doc-drift Detect Spec Drift

## Summary
Drift-only mode: validated 5 overlapping spec files against surface-scalar-dot-fix changes. Found 3 drifted rules, 0 orphaned. Fixed TECH-surface-rules-002 (recipe naming scalar vs named), TECH-surface-rules-003 (data propagation text mode), BIZ-surface-orchestration-006 (naming constraint scalar acknowledgment). Updated domains frontmatter for both files.

## Changes

### Files Created
无

### Files Modified
- docs/conventions/surface-rules.md
- docs/business-rules/surface-orchestration.md

### Key Decisions
无

## Document Metrics
5 spec files checked, 3 rules drifted and fixed, 4 rules current, 0 orphaned

## Referenced Documents
- docs/proposals/surface-scalar-dot-fix/proposal.md
- docs/conventions/surface-cli.md
- docs/conventions/skill-self-containment.md
- docs/conventions/skill-structure.md

## Review Status
final

## Acceptance Criteria
- [x] Only specs with domains overlapping changed files were checked
- [x] Drift report produced and drift fixed

## Notes
Narrowed scope to 5 overlapping spec files using git diff domain-matching per task discovery strategy. Non-overlapping specs (12 files) skipped. No implicit new rules discovered.
