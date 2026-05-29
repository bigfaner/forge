---
status: "completed"
started: "2026-05-30 06:13"
completed: "2026-05-30 06:14"
time_spent: "~1m"
---

# Task Record: 19 Fix: gen-sitemap Step 2b/4 overlap handling

## Summary
Added Step 2b/Step 4 overlap handling to gen-sitemap SKILL.md: reuse Step 2b snapshots for already-visited routes, extract page-specific elements only, and explore remaining routes with full workflow.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/gen-sitemap/SKILL.md

### Key Decisions
无

## Document Metrics
1 section added (~20 lines), 2 AC satisfied

## Referenced Documents
- docs/features/plugin-consistency-audit/reports/04-skills-batch-c.md

## Review Status
final

## Acceptance Criteria
- [x] Step 4 explicitly states how Step 2b-explored pages are handled (reuse results or skip layout comparison)
- [x] No redundant exploration between Step 4 and Step 2b

## Notes
Fix addresses T-03 (P1 TIMING) from Report 04. Added a new subsection 'Handling Pages Already Explored in Step 2b' with 3 numbered rules and pseudocode example.
