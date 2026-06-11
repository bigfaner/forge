---
status: "completed"
started: "2026-05-20 00:14"
completed: "2026-05-20 00:16"
time_spent: "~2m"
---

# Task Record: 4 Submit-task SKILL.md cleanup — remove CLI-enforced validation rules

## Summary
Removed CLI-enforced validation rules from submit-task SKILL.md to eliminate duplication with submit.go. Removed 'Validation Rules (enforced by CLI)' section (quality gate pre-check + data validation table), reduced 'What forge task submit Does' to a two-line summary, removed just test [scope] example from metrics collection, and cleaned all noTest references from Fields table and coverage rules.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/submit-task/SKILL.md

### Key Decisions
- Kept only agent-unique instructions in SKILL.md: metrics collection workflow, JSON data format, type reclassification, forbidden operations, and recovery steps
- Replaced 'just test [scope]' example with instruction to capture metrics from targeted test runs during development
- Coverage auto-set to -1.0 now described solely in terms of non-coding.* types, with no noTest mention

## Test Results
- **Tests Executed**: No
- **Passed**: 0
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] 'Validation Rules (enforced by CLI)' section removed
- [x] 'What forge task submit Does' section reduced to single line summary
- [x] Quality Gate Pre-check subsection removed
- [x] Data Validation subsection removed
- [x] Metrics Collection updated: no just test [scope] command, instructs to capture from targeted runs
- [x] Coverage rules updated: -1.0 auto-set for non-coding.* types only, no noTest mention
- [x] Remaining sections intact: File Locations, JSON Data Format, Fields table, Metrics Collection, Type Reclassification, Usage, Forbidden Operations, Recovery

## Notes
Doc-type task. No code changes, no tests needed. SKILL.md reduced from 207 lines to 157 lines by removing CLI-duplicated content.
