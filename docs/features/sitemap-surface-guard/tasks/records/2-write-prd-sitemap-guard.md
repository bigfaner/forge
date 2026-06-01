---
status: "completed"
started: "2026-06-01 21:24"
completed: "2026-06-01 21:26"
time_spent: "~2m"
---

# Task Record: 2 write-prd: Conditionalize sitemap references by surface type

## Summary
Added web surface guard to 3 write-prd files (SKILL.md, ui-functions.md, self-check.md) to conditionalize sitemap.json references by surface type

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/write-prd/SKILL.md
- plugins/forge/skills/write-prd/rules/ui-functions.md
- plugins/forge/skills/write-prd/rules/self-check.md

### Key Decisions
无

## Document Metrics
3 files modified, 3 guard points added, consistent forge surfaces --json pattern

## Referenced Documents
- docs/proposals/sitemap-surface-guard/proposal.md
- plugins/forge/skills/gen-test-scripts/types/ui.md

## Review Status
final

## Acceptance Criteria
- [x] SKILL.md Step 1 reads sitemap.json only after checking web surface via forge surfaces --json
- [x] ui-functions.md Placement Rules adds web surface precondition before sitemap read
- [x] self-check.md Placement consistency and Sitemap availability checks guarded by web surface condition

## Notes
Guard pattern follows gen-test-scripts/types/ui.md as reference template. All 3 files use consistent forge surfaces --json check for web surface detection.
