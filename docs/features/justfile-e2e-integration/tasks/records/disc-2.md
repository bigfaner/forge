---
status: "completed"
started: "2026-04-29 20:13"
completed: "2026-04-29 20:14"
time_spent: "~1m"
---

# Task Record: disc-2 Fix: add justfile prerequisite check to run-e2e-tests SKILL.md

## Summary
Added Justfile prerequisite check to run-e2e-tests SKILL.md. The Prerequisites section now includes a row for 'Justfile with e2e-setup recipe' and a bash check that verifies the Justfile exists and contains the e2e-setup recipe, with a prompt to run /init-justfile if missing.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/run-e2e-tests/SKILL.md

### Key Decisions
- Added Justfile existence and e2e-setup recipe check before other prerequisite checks so users are directed to /init-justfile early in the workflow

## Test Results
- **Passed**: 0
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] Prerequisites section references Justfile existence check
- [x] Prerequisites section references e2e-setup recipe check
- [x] User is prompted to run /init-justfile when Justfile or e2e-setup recipe is missing

## Notes
无
