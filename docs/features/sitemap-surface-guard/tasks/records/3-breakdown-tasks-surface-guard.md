---
status: "completed"
started: "2026-06-01 21:27"
completed: "2026-06-01 21:30"
time_spent: "~3m"
---

# Task Record: 3 breakdown-tasks/ui-placement: Add surface guard before route validation

## Summary
Added web surface type guard to breakdown-tasks/rules/ui-placement.md Placement Validation section. Step 2 now checks forge surfaces --json for web surface before accessing sitemap.json. Non-web projects skip route verification with an INFO message instead of a misleading WARN.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/breakdown-tasks/rules/ui-placement.md

### Key Decisions
无

## Document Metrics
1 file modified, 6 lines changed in Placement Validation section

## Referenced Documents
- docs/proposals/sitemap-surface-guard/proposal.md
- plugins/forge/skills/gen-test-scripts/types/ui.md

## Review Status
final

## Acceptance Criteria
- [x] Placement Validation section checks web surface before sitemap.json existence
- [x] No web surface skips route verification with appropriate INFO message (not WARN)

## Notes
Inserted new step 2 (Surface check) before original step 2 (sitemap.json check). Steps renumbered to 1-6. Guard pattern consistent with gen-test-scripts/types/ui.md. Uses forge surfaces --json per proposal.
