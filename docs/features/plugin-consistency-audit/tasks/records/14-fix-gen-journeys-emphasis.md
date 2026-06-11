---
status: "completed"
started: "2026-05-30 06:07"
completed: "2026-05-30 06:08"
time_spent: "~1m"
---

# Task Record: 14 Fix: gen-journeys test level emphasis mismatch

## Summary
Fixed gen-journeys SKILL.md Per-Surface Rule Application section: updated API, Web, TUI test level emphasis to match corresponding rules/surface-*.md files (C-15 P1 CONFLICT)

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/gen-journeys/SKILL.md

### Key Decisions
无

## Document Metrics
3/5 surface emphasis lines corrected, CLI and Mobile unchanged

## Referenced Documents
- docs/features/plugin-consistency-audit/reports/04-skills-batch-c.md
- plugins/forge/skills/gen-journeys/rules/surface-api.md
- plugins/forge/skills/gen-journeys/rules/surface-web.md
- plugins/forge/skills/gen-journeys/rules/surface-tui.md
- plugins/forge/skills/gen-journeys/rules/surface-cli.md
- plugins/forge/skills/gen-journeys/rules/surface-mobile.md

## Review Status
final

## Acceptance Criteria
- [x] API surface test level emphasis matches rules/surface-api.md (Balanced 50/50)
- [x] Web surface test level emphasis matches rules/surface-web.md (Balanced 50/50)
- [x] TUI surface test level emphasis matches rules/surface-tui.md (Contract 80% / Journey smoke 20%)
- [x] CLI and Mobile surfaces unchanged (already consistent)

## Notes
Audit finding C-15: 3 of 5 surfaces had mismatched test level emphasis descriptions in SKILL.md vs rules files. Fixed API (integration-heavy -> Balanced 50/50), Web (e2e-heavy -> Balanced 50/50), TUI (integration-heavy -> Contract 80% / Journey smoke 20%). Rules files treated as authoritative per Hard Rules.
