---
status: "completed"
started: "2026-05-17 13:22"
completed: "2026-05-17 13:25"
time_spent: "~3m"
---

# Task Record: 4 Create types/tui.md instruction file

## Summary
Created plugins/forge/skills/gen-test-scripts/types/tui.md with TUI-specific generation instructions: reconnaissance strategy (grep patterns for bubbletea, tview, tcell, termbox, ratatui, textual, urwid), Fact Table required keys (TUI_BINARY, TUI_ENTRY_POINT, TUI_KEYBIND_*), verification method (terminal framework imports + alternate screen detection), generation patterns (non-interactive stdin pipe, key sequence encoding, screen state assertions, exit code + output combination), and TUI antipattern guards (sleep-based waits, source code inspection, interactive prompts without stdin pipe). Includes HARD-RULE tags enforcing non-interactive execution and manual-only marking for real-terminal tests.

## Changes

### Files Created
- plugins/forge/skills/gen-test-scripts/types/tui.md

### Files Modified
无

### Key Decisions
- Modeled structure after existing cli.md and api.md type files for consistency
- Included Key Sequence Encoding table with stdin escape codes for arrow keys, escape, tab, ctrl+c
- Added Tests Requiring Real Terminal Interaction section with HARD-RULE for manual-only marking, per task Hard Rules
- Used TUI_BINARY + TUI_KEYBIND_* as minimum Fact Table keys (parallels CLI_BINARY + CLI_COMMAND_* pattern)

## Test Results
- **Tests Executed**: No
- **Passed**: 0
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] plugins/forge/skills/gen-test-scripts/types/tui.md exists
- [x] Frontmatter declares type: tui and conventions: [testing-tui.md]
- [x] Contains Reconnaissance Strategy section with TUI-specific search patterns
- [x] Contains Fact Table Required Keys section listing minimum keys for TUI type
- [x] Contains Verification Method section describing how to confirm TUI interface
- [x] Contains Generation Patterns section for TUI test case translation
- [x] Contains TUI Antipattern Guards section
- [x] At least 3 section headings are unique to this file
- [x] TUI vs CLI disambiguation documented in reconnaissance strategy

## Notes
Documentation task -- no code changes, no test execution required. 7 unique section headings (Key Sequence Encoding, Screen State Assertions, Non-Interactive Execution Model, Tests Requiring Real Terminal Interaction, Sleep-Based Waits, Source Code Inspection, Interactive Prompts Without Stdin Pipe).
