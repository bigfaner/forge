---
status: "completed"
started: "2026-05-30 06:23"
completed: "2026-05-30 06:24"
time_spent: "~1m"
---

# Task Record: 26 Fix: ui-design TUI requirements triplication

## Summary
Replaced TUI structural requirements triplication with single-authority references. tui.md now references rules/tui-panel-requirements.md instead of duplicating 47 lines of structural requirements. ui-design.md TUI Component section now references rules/tui-panel-requirements.md as the authoritative source while retaining template tables for fill-in use.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/ui-design/templates/platforms/tui.md
- plugins/forge/skills/ui-design/templates/ui-design.md

### Key Decisions
无

## Document Metrics
reduced redundancy from 3-file triplication to 1 authoritative source + 2 references

## Referenced Documents
- docs/features/plugin-consistency-audit/reports/03-skills-batch-b.md
- plugins/forge/skills/ui-design/rules/tui-panel-requirements.md

## Review Status
final

## Acceptance Criteria
- [x] templates/platforms/tui.md no longer verbatim-repeats 5 structural requirements, references rules/tui-panel-requirements.md
- [x] templates/ui-design.md TUI Component section references rules/tui-panel-requirements.md
- [x] rules/tui-panel-requirements.md unchanged (authoritative source)

## Notes
ui-design.md retains template tables (Dimensions, Character Palette, Color Mapping, Edge Cases) as fill-in scaffolding since it is a template file, but the section heading and comment now point to rules/tui-panel-requirements.md as the authority.
