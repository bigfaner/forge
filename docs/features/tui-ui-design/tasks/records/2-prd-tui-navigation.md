---
status: "completed"
started: "2026-05-15 00:53"
completed: "N/A"
time_spent: ""
---

# Task Record: 2 PRD TUI navigation template and write-prd awareness

## Summary
Added TUI Navigation Architecture section to prd-ui-functions.md template with Keymap, Panel Layout, Modes, and Navigation Rules tables, conditionally rendered when platform=tui. Updated write-prd SKILL.md Step 8 with platform-aware navigation handling instructions.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/write-prd/templates/prd-ui-functions.md
- plugins/forge/skills/write-prd/SKILL.md

### Key Decisions
- TUI navigation uses a separate section (Keymap + Panel Layout + Modes) rather than extending the existing web/mobile tables, since TUI is keyboard-driven vs pointer-driven
- Both navigation sections coexist in the template with HTML comments marking conditional rendering boundaries -- the agent renders exactly one based on platform value

## Test Results
- **Tests Executed**: No
- **Passed**: 18
- **Failed**: 1
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] prd-ui-functions.md includes TUI Navigation Architecture section with Keymap, Panel Layout, Modes, Navigation Rules tables
- [x] TUI navigation section is conditionally rendered -- only appears when platform=tui, does not affect web/mobile templates
- [x] write-prd SKILL.md references the TUI navigation template and triggers it when platform=tui
- [x] Existing web/mobile PRD generation behavior unchanged

## Notes
Pre-existing test failure in forge-cli/internal/cmd (unrelated to this change, verified on clean state). This is a template/documentation task with no code changes. Coverage set to -1.0 as no new testable code was added.
