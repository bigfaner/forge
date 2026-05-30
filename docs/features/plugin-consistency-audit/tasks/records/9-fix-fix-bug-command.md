---
status: "completed"
started: "2026-05-30 05:59"
completed: "2026-05-30 06:00"
time_spent: "~1m"
---

# Task Record: 9 Fix: fix-bug command - add AskUserQuestion, replace Playwright, add mobile/tui

## Summary
Fix fix-bug command: added AskUserQuestion to allowed-tools, replaced Playwright hardcode with generic test profile description, added mobile and tui rows to Bug surface table

## Changes

### Files Created
无

### Files Modified
- plugins/forge/commands/fix-bug.md

### Key Decisions
无

## Document Metrics
3 fixes applied: 1 allowed-tools addition, 1 hardcode replacement, 2 table rows added

## Referenced Documents
- docs/features/plugin-consistency-audit/reports/05-commands-agent-hooks.md
- docs/reference/test-type-model.md
- plugins/forge/skills/run-tests/rules/surfaces/mobile.md
- plugins/forge/skills/run-tests/rules/surfaces/tui.md

## Review Status
final

## Acceptance Criteria
- [x] allowed-tools contains AskUserQuestion
- [x] Bug surface table replaces Playwright with generic description
- [x] Bug surface table includes mobile row (Maestro YAML / mobile test profile)
- [x] Bug surface table includes tui row (child_process + stdin pipe / tui test profile)

## Notes
Runner descriptions aligned with test-type-model.md execution model column. Mobile uses Maestro YAML, TUI uses child process + stdin pipe, matching the authoritative spec.
