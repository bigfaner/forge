---
status: "completed"
started: "2026-05-20 21:23"
completed: "2026-05-20 21:24"
time_spent: "~1m"
---

# Task Record: 3 Update forge guide for runTasks

## Summary
Updated forge guide Automation Config section to add runTasks configuration item with quick/full sub-keys and default values

## Changes

### Files Created
无

### Files Modified
- plugins/forge/hooks/guide.md

### Key Decisions
- Kept existing single-paragraph style for Automation Config section, adding runTasks inline with existing settings

## Test Results
- **Tests Executed**: No
- **Passed**: 0
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] Automation Config section lists runTasks config item with quick/full sub-keys
- [x] Default values documented: quick: true, full: false
- [x] Purpose explained: controls whether /quick skips confirmation gate and auto-executes

## Notes
Task type is doc, no quality gate or test metrics required
