---
status: "completed"
started: "2026-05-16 17:47"
completed: "2026-05-16 17:48"
time_spent: "~1m"
---

# Task Record: 3 Extend test-cases rubric with antipattern detection dimension

## Summary
Added a 6th scoring dimension 'Test Code Quality' (200 pts) to rubrics/test-cases.md that checks for 6 known antipatterns: recursive test invocation, unconditional t.Skip, vacuous assertions, conditional skip without self-contained fixture, duplicate test function names across packages, and static-file text grep tests. Redistributed points from Completeness (200->150) and Structure & ID Integrity (100->50) to keep total at 1000.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/eval/rubrics/test-cases.md

### Key Decisions
- Added blocking threshold of <150 pts for Test Code Quality dimension (similar to Step Actionability's <200 threshold)
- Redistributed points proportionally: Completeness lost 50 pts (70/70/60 -> 50/50/50), Structure lost 50 pts (40/30/30 -> 20/15/15)
- Referenced lesson documents as authoritative sources for antipattern definitions rather than duplicating explanations

## Test Results
- **Tests Executed**: No
- **Passed**: 0
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] rubrics/test-cases.md includes a Test Code Quality dimension
- [x] The dimension checks for all 6 antipatterns: recursive test invocation, unconditional t.Skip, vacuous assertions, conditional skip without fixture, duplicate test function names, static-file text grep
- [x] Total rubric points remain 1000 (Completeness 200->150, Structure 100->50, Test Code Quality 200)

## Notes
Documentation-only task; no test scope. Used coverage=-1.0 per documentation task convention.
