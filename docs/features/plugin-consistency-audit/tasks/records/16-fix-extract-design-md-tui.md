---
status: "completed"
started: "2026-05-30 06:09"
completed: "2026-05-30 06:10"
time_spent: "~1m"
---

# Task Record: 16 Fix: extract-design-md missing TUI match strategy in SKILL.md

## Summary
Added TUI match strategy to SKILL.md Step 3, distinguishing web/mobile (5 built-in styles) from TUI (2 built-in themes) with reference to rules/platform-routing.md

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/extract-design-md/SKILL.md

### Key Decisions
无

## Document Metrics
Step 3 restructured: 1 section -> 2 sub-sections (Web/Mobile + TUI), TUI option table added, cross-reference to rules/platform-routing.md section 4

## Referenced Documents
- docs/features/plugin-consistency-audit/reports/04-skills-batch-c.md
- plugins/forge/skills/extract-design-md/rules/platform-routing.md
- plugins/forge/skills/extract-design-md/rules/match-strategy.md

## Review Status
final

## Acceptance Criteria
- [x] Step 3 distinguishes web/mobile using 5 web built-in styles
- [x] Step 3 distinguishes TUI using 2 TUI themes (modern-dark-tui / minimal-ascii-tui)
- [x] TUI match strategy references rules/platform-routing.md for complete definition

## Notes
Resolved audit finding C-05 (P1 CONFLICT). Step 3 now has ### Web / Mobile and ### TUI sub-sections with distinct match strategies.
