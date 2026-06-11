---
status: "completed"
started: "2026-05-24 16:25"
completed: "2026-05-24 16:30"
time_spent: "~5m"
---

# Task Record: 7 gen-journeys skill adaptation for surfaces

## Summary
Updated gen-journeys SKILL.md to use `forge surfaces <path>` CLI query instead of independent project scanning. Renamed surface-webui.md to surface-web.md and updated all naming references (WebUI -> web, Mobile -> mobile) across all rule files.

## Changes

### Files Created
- plugins/forge/skills/gen-journeys/rules/surface-web.md

### Files Modified
- plugins/forge/skills/gen-journeys/SKILL.md
- plugins/forge/skills/gen-journeys/rules/surface-api.md
- plugins/forge/skills/gen-journeys/rules/surface-mobile.md

### Key Decisions
无

## Document Metrics
1 file created (rename), 3 files modified, 1 file deleted. SKILL.md Surface Detection section fully rewritten (~65 lines replaced). All 6 AC items passed.

## Referenced Documents
- docs/proposals/unify-surfaces/proposal.md

## Review Status
completed

## Acceptance Criteria
- [x] SKILL.md Surface Detection section rewritten to use forge surfaces <path> CLI call
- [x] SKILL.md uses exit code contract: exit 0 -> parse stdout, exit 1 -> prompt user
- [x] surface-webui.md renamed to surface-web.md
- [x] surface-mobileui.md renamed to surface-mobile.md (if exists with old naming)
- [x] Internal naming references updated: webui -> web, mobileui -> mobile
- [x] Rule file loading logic updated to match new filenames

## Notes
surface-mobileui.md did not exist (already named surface-mobile.md). SKILL.md HARD-RULE updated to explicitly forbid independent project file scanning. Rule file content structure preserved as required by Hard Rules -- only naming references changed.
