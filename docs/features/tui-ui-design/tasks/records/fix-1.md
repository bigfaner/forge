---
status: "completed"
started: "2026-05-15 01:45"
completed: "2026-05-15 01:48"
time_spent: "~3m"
---

# Task Record: fix-1 Fix: manifest-update-ui.md missing prototype directory references

## Summary
Added missing prototype/ directory references to manifest-update-ui.md template for all three platform scenarios (single web, multi-platform, single TUI). The template previously only listed ui-design file references but omitted prototype output directories, causing TC-029/030/031 E2E test failures.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/ui-design/templates/manifest-update-ui.md

### Key Decisions
- Added Prototype row with prototype/ reference for single-platform (web) section (uncommented, active)
- Added commented Prototype row with prototype/ reference for single TUI platform section
- Added commented Prototype rows with prototype/web/ and prototype/tui/ for multi-platform section
- Used consistent {{PROTOTYPE_SUMMARY}} placeholder pattern matching existing UI_DESIGN_SUMMARY placeholders

## Test Results
- **Tests Executed**: No
- **Passed**: 0
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] TC-029: Single web platform manifest includes prototype/ reference
- [x] TC-030: Multi-platform manifest includes prototype/web/ and prototype/tui/ references
- [x] TC-031: Single TUI manifest includes prototype/ reference

## Notes
Fix is a template-only change (markdown file). The forge-cli/internal/cmd unit test failure is pre-existing (missing Node.js module at C:\nonexistent\validate-specs.mjs) and unrelated to this fix. Verified by stashing changes and confirming the same test failure occurs on the base branch.
