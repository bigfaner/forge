---
status: "completed"
started: "2026-05-31 15:11"
completed: "2026-05-31 15:13"
time_spent: "~2m"
---

# Task Record: 3 Insert task sizing audit step in quick-tasks skill

## Summary
Inserted Task Sizing Audit step (Step 4) in quick-tasks SKILL.md between task file creation and index generation. Renamed all validate-index references to validate. Renumbered steps 4-8 to 5-9 for continuity.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/quick-tasks/SKILL.md

### Key Decisions
无

## Document Metrics
1 audit step inserted, 1 validate-index→validate rename, 5 steps renumbered, flow diagram updated

## Referenced Documents
- docs/proposals/task-sizing-gate/proposal.md
- plugins/forge/skills/breakdown-tasks/SKILL.md

## Review Status
final

## Acceptance Criteria
- [x] quick-tasks SKILL.md contains independent task sizing audit step after task file creation and before forge task index
- [x] audit step instructs LLM to check multi-verb and cross-domain AC, auto-split and output report
- [x] All validate-index references updated to validate
- [x] All step numbers are consecutive with no gaps

## Notes
Audit step format kept consistent with breakdown-tasks SKILL.md Step 5. Split Protocol uses <seq> convention instead of <phase>.<sub> since quick-tasks uses simple integer IDs.
