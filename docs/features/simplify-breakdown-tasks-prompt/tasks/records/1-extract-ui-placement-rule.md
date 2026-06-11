---
status: "completed"
started: "2026-05-19 23:00"
completed: "2026-05-19 23:02"
time_spent: "~2m"
---

# Task Record: 1 Extract ui-placement rule file

## Summary
Extracted all UI-related conditional rules from SKILL.md into a standalone rule file at rules/ui-placement.md. The file consolidates content gated by 5 conditional tags (HAS_UI, NO_UI, UI_ONLY, HAS_PLACEMENT, RULE) into a single self-contained document with a unified load condition (ui/ui-design.md exists OR prd/prd-ui-functions.md exists).

## Changes

### Files Created
- plugins/forge/skills/breakdown-tasks/rules/ui-placement.md

### Files Modified
无

### Key Decisions
- Merged 5 conditional tags into a single load condition using OR logic (ui/ui-design.md exists OR prd/prd-ui-functions.md exists) rather than maintaining separate activation paths
- Used skill-relative file paths throughout the rule file, consistent with forge distribution model
- Preserved placement format note with canonical form documentation and examples

## Test Results
- **Tests Executed**: No
- **Passed**: 0
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] rules/ui-placement.md created under plugins/forge/skills/breakdown-tasks/rules/
- [x] File contains UI-specific element mapping rows from UI_ONLY block
- [x] File contains placement validation procedure from HAS_PLACEMENT block
- [x] File contains UI task split rules from RULE block
- [x] File contains UI reference file requirements for Build, Integration, and Page Assembly tasks
- [x] File contains UI dependency layer rules
- [x] File contains UI prototype reading instruction
- [x] Load condition documented at top of file
- [x] Guard clause included for empty/malformed artifacts
- [x] Maintenance note listing skeleton section dependencies
- [x] File is independently understandable

## Notes
Doc-type task, no test metrics applicable. SKILL.md was NOT modified per hard rules.
