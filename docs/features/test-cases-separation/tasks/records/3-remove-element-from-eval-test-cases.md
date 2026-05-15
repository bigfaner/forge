---
status: "completed"
started: "2026-05-15 19:46"
completed: "2026-05-15 19:48"
time_spent: "~2m"
---

# Task Record: 3 Remove Element evaluation from eval-test-cases skill and rubric

## Summary
Removed Element field evaluation from eval-test-cases rubric and SKILL.md. Renamed Dimension 3 web-ui from 'Route & Element Accuracy' to 'Route Accuracy'. Removed 'Elements are identifiable' evaluation item. Redistributed web-ui points from 70+70+60 to 120+80 (total 200pts preserved).

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/eval-test-cases/templates/rubric.md

### Key Decisions
- Redistributed web-ui points as 120 (Routes valid and specific) + 80 (Route consistency) to maintain 200pts total after removing Elements are identifiable (70pts) and Route/Element consistency (60pts)
- SKILL.md had no Element-specific instructions to remove; only rubric.md required changes
- Added 'No Route field contains implementation details (testid, selector, CSS)' to Route consistency criterion to reinforce the separation of concerns from the proposal

## Test Results
- **Tests Executed**: No
- **Passed**: 0
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] eval-test-cases SKILL.md no longer instructs agents to evaluate Element field quality
- [x] rubric.md Dimension 3 web-ui title is 'Route Accuracy' (not 'Route & Element Accuracy')
- [x] 'Elements are identifiable' evaluation item is removed from rubric
- [x] Total scoring for web-ui dimension remains 200pts (points redistributed to remaining items)

## Notes
Documentation-only task. SKILL.md verified clean of Element references. No test execution needed.
