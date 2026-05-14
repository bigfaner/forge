---
status: "completed"
started: "2026-05-15 01:01"
completed: "2026-05-15 01:07"
time_spent: "~6m"
---

# Task Record: 3 ui-design SKILL.md and template TUI support

## Summary
Modified ui-design SKILL.md to add TUI platform detection, TUI theme selection (Modern Dark / Minimal ASCII / DESIGN.md custom), TUI-specific panel design requirements with 5 mandatory structural items, multi-platform file splitting (web+tui), and TUI prototype generation rules. Modified ui-design.md template to add TUI component template section with all 5 structural requirements from the lesson (ASCII Layout Mockup, Dimensions, Character Palette, Color Mapping, Edge Cases) plus States, Key Bindings, and Data Binding.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/ui-design/SKILL.md
- plugins/forge/skills/ui-design/templates/ui-design.md

### Key Decisions
- TUI platform branches off at Step 3 (Select Design Style) with its own theme selection separate from web/mobile styles
- TUI uses the same template file (ui-design.md) but with an additional TUI Component section appended for TUI panels
- Multi-platform output naming convention: single TUI produces ui-design-tui.md, multi-platform produces ui-design-web.md + ui-design-tui.md
- 5 structural requirements enforced via HARD-RULE tag in SKILL.md to prevent agents from skipping visual specs
- TUI prototype generates HTML simulating terminal window with simulated key buttons for panel switching

## Test Results
- **Tests Executed**: Yes
- **Passed**: 0
- **Failed**: 0
- **Coverage**: 0.0%

## Acceptance Criteria
- [x] SKILL.md detects platform=tui from PRD and enters TUI branch
- [x] TUI branch presents theme selection: Modern Dark, Minimal ASCII, or DESIGN.md custom
- [x] TUI branch uses platforms/tui.md for navigation rules and selected theme for visual style
- [x] ui-design.md template includes TUI component template with all 5 structural requirements from lesson
- [x] Multi-platform features produce separate files: ui-design-web.md + ui-design-tui.md
- [x] Single TUI feature produces ui-design-tui.md
- [x] Existing web/mobile behavior unchanged

## Notes
Pre-existing test failure in forge-cli/internal/cmd is unrelated to these changes (confirmed by stashing and re-running). compile, fmt, lint all pass cleanly. This task modifies Markdown skill/template files only, no Go code.
