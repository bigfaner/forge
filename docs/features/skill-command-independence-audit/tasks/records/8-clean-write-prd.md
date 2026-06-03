---
status: "completed"
started: "2026-06-04 00:55"
completed: "2026-06-04 00:58"
time_spent: "~3m"
---

# Task Record: 8 Delete Related sections + reduce redundancy in write-prd

## Summary
Delete Related/Integration sections + reduce redundancy in write-prd SKILL.md: removed ## Integration chapter, consolidated 4x Process Flow/Checklist/Output Documents into unified versions with intent-differentiated columns

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/write-prd/SKILL.md

### Key Decisions
无

## Document Metrics
454 lines -> 361 lines (~20% reduction); 3 sections consolidated from 12 to 3; 1 chapter deleted

## Referenced Documents
- docs/proposals/skill-command-independence-audit/proposal.md

## Review Status
final

## Acceptance Criteria
- [x] Related Skills/Integration/References sections deleted
- [x] 4 intent variants consolidated with core differences preserved
- [x] Override Signals table is write-prd-specific (not shared with tech-design)

## Notes
Process Flow: 4 separate flows -> 1 unified flow + intent diff table. Checklist: 4 separate lists -> 1 unified list with conditional annotations. Output Documents: 4 separate tables -> 1 unified table with Generated-For column.
