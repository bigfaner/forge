---
status: "completed"
started: "2026-05-15 01:08"
completed: "2026-05-15 01:14"
time_spent: "~6m"
---

# Task Record: 4 TUI HTML prototype simulation rules

## Summary
Add TUI-specific prototype generation rules to prototype.md: single-file index.html with terminal-window div, simulated key buttons for panel switching, monospace dark-theme CSS, panel layout matching ASCII mockups, and output path rules for single-TUI vs multi-platform features.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/ui-design/templates/prototype.md

### Key Decisions
- TUI prototypes are single index.html (not multi-file like web/mobile), following proposal D6
- Terminal CSS uses #1e1e1e background with monospace font stack as terminal approximation
- Simulated key buttons ([Tab], [1]-[9], [q], [:command]) at bottom of terminal window for panel switching
- ASCII mockup from ui-design.md rendered inside <pre> blocks to preserve spacing
- Panel focus uses CSS class toggling (.focused) to simulate keyboard-driven panel highlighting
- Output path: prototype/index.html for single-TUI, prototype/tui/index.html for multi-platform

## Test Results
- **Tests Executed**: No
- **Passed**: 0
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] prototype.md includes TUI-specific prototype generation rules
- [x] TUI prototype is a single index.html with all panels rendered inside a black terminal window div (proposal D6)
- [x] Includes simulated key buttons at the bottom: [Tab], [1], [q], [:command] to switch between panels
- [x] Uses monospace font and dark background to approximate terminal appearance
- [x] Panel layout in HTML matches the ASCII mockup from ui-design.md
- [x] TUI prototypes output to prototype/tui/ (multi-platform) or prototype/ (single TUI feature)

## Notes
Template-only change (no Go code). Pre-existing test failure in forge-cli/internal/cmd (TestSaveIndexAndSignalCompletion_SaveIndexError) confirmed unrelated to this change.
