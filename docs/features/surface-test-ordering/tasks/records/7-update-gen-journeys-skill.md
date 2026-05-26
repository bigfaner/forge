---
status: "completed"
started: "2026-05-26 14:10"
completed: "2026-05-26 14:15"
time_spent: "~5m"
---

# Task Record: 7 Update gen-journeys SKILL.md for multi-surface support

## Summary
Updated gen-journeys SKILL.md for multi-surface support: added Multi-Surface Rules Loading section with per-surface rule guidance, surface_types frontmatter annotation for Journey files, and surface coverage validation

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/gen-journeys/SKILL.md
- plugins/forge/skills/gen-journeys/templates/journey.md

### Key Decisions
无

## Document Metrics
2 files modified, 1 new section added (Multi-Surface Rules Loading ~60 lines), 3 existing sections updated (Surface Detection, Step 4, Step 5 validation)

## Referenced Documents
- docs/proposals/surface-test-ordering/proposal.md
- docs/conventions/forge-distribution.md

## Review Status
completed

## Acceptance Criteria
- [x] SKILL.md contains multi-surface rule loading guidance organized by surface type sections
- [x] Output format requires each Journey to annotate covered surface type set (e.g. [web, api])
- [x] All configured surface types must be covered by at least one Journey

## Notes
Surface Detection section updated from single-surface to multi-surface detection flow. Journey template updated with surface_types frontmatter field. Step 5 validation table extended with surface coverage checks.
