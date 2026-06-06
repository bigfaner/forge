---
status: "completed"
started: "2026-06-06 13:27"
completed: "2026-06-06 13:31"
time_spent: "~4m"
---

# Task Record: 5 Update init-justfile SKILL.md and surface rules

## Summary
Updated init-justfile SKILL.md with surface-key test directory path rules and added Test-Dir-Path guidance to all 5 surface rule files (api, cli, mobile, tui, web)

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/init-justfile/SKILL.md
- plugins/forge/skills/init-justfile/rules/surfaces/api.md
- plugins/forge/skills/init-justfile/rules/surfaces/cli.md
- plugins/forge/skills/init-justfile/rules/surfaces/mobile.md
- plugins/forge/skills/init-justfile/rules/surfaces/tui.md
- plugins/forge/skills/init-justfile/rules/surfaces/web.md

### Key Decisions
无

## Document Metrics
6 files modified, 4 new documentation sections added to SKILL.md, 5 Test-Dir-Path blocks added to surface rules

## Referenced Documents
- docs/proposals/surface-key-test-alignment/proposal.md

## Review Status
final

## Acceptance Criteria
- [x] init-justfile SKILL.md recipe prefix logic aligned with surface-key directory structure
- [x] 5 surface rule files recipe definitions work correctly under tests/<surfaceKey>/<journey>/ paths

## Notes
SKILL.md changes: (1) Added 'Test Directory Path in Recipes' subsection under Surface-Level Targets explaining single vs multi surface path rules. (2) Added 'Test directory path' guidance in Step 3b for recipe body generation. (3) Updated Step 3d group naming from <surface-key-or-type> to <surface-key> with explicit scalar/named guidance. (4) Added test directory path note in Notes section. Surface rule changes: Each of the 5 files (api, cli, mobile, tui, web) received a <Test-Dir-Path> block before the Recipe Template section, specifying the adaptive path rule with a concrete key=type example.
