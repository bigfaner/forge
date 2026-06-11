---
status: "completed"
started: "2026-05-31 15:09"
completed: "2026-05-31 15:11"
time_spent: "~2m"
---

# Task Record: 2 Insert task sizing audit step in breakdown-tasks skill

## Summary
Inserted task sizing audit step (Step 5) in breakdown-tasks SKILL.md between task file creation and index generation. Updated all validate-index references to validate. Renumbered subsequent steps (5->6, 6->7, 7->8, 8->9).

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/breakdown-tasks/SKILL.md

### Key Decisions
无

## Document Metrics
audit step ~35 lines, 3 check items, split protocol with report format

## Referenced Documents
- docs/proposals/task-sizing-gate/proposal.md
- docs/conventions/forge-distribution.md

## Review Status
final

## Acceptance Criteria
- [x] breakdown-tasks SKILL.md has independent task sizing audit step between write-task-files and forge-task-index
- [x] audit step instructs LLM to check multi-verb and cross-domain AC, auto-split and output report
- [x] all validate-index references updated to validate
- [x] all step numbers are sequential with no gaps

## Notes
Steps now run 0-9 (was 0-8). Audit step includes multi-verb detection, AC cross-domain detection, and operational ceiling checks with split protocol and hard gate.
