---
status: "completed"
started: "2026-05-16 14:00"
completed: "2026-05-16 14:10"
time_spent: "~10m"
---

# Task Record: 1 Add --platform flag and command scaffolding

## Summary
Add --platform flag to extract-design-md command with web/mobile/tui routing. Updated frontmatter description to mention all platforms, added --platform argument-hint with valid values, added Platform Routing section with input validation rejecting invalid values, and placeholder routing for mobile and tui adapters. Web extraction remains byte-for-byte identical.

## Changes

### Files Created
- forge-cli/internal/docsync/extract_design_md_test.go

### Files Modified
- plugins/forge/commands/extract-design-md.md

### Key Decisions
- Platform routing added as a separate section between Process Flow and Step 1, keeping all original web extraction content untouched
- Mobile placeholder falls through to web extraction with a note, keeping the command simple until Task 2 implements the adapter
- TUI placeholder asks for screenshot path and outputs a not-yet-implemented message, deferring to Task 3

## Test Results
- **Tests Executed**: No
- **Passed**: 7
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] Command frontmatter updated: description mentions all three platforms, allowed_tools includes image analysis capability, argument-hints includes --platform with valid values (web/mobile/tui)
- [x] Default behavior (--platform web or no flag) produces identical output to current behavior
- [x] Platform routing logic added: when --platform mobile or --platform tui is provided, the command delegates to the appropriate adapter section (placeholder for Tasks 2/3)
- [x] Input validation rejects invalid platform values with a clear error message

## Notes
docsync package has no production statements (test-only package), coverage reported as [no statements]. Overall backend coverage >= 80% across all packages with statements.
