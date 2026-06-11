---
status: "completed"
started: "2026-05-16 17:49"
completed: "2026-05-16 17:53"
time_spent: "~4m"
---

# Task Record: 4 Harden gen-test-scripts SKILL.md to forbid antipatterns

## Summary
Added 'Antipattern Guard (Mandatory Pre-Emission Check)' section to gen-test-scripts SKILL.md. The section enumerates 6 forbidden patterns (recursive test invocation, unconditional t.Skip, vacuous assertions, conditional skip without fixture, duplicate test functions, static-file text grep) with what/why/instead columns, a concrete validation procedure, and authoritative references to lesson documents.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/gen-test-scripts/SKILL.md

### Key Decisions
- Placed the Antipattern Guard section within Step 4 (Generate Spec Files) rather than as a standalone step, preserving the existing generation flow structure
- Used a table format matching the rubric dimension added in Task 3 for consistency
- Included a concrete 6-step validation procedure after the table to give agents an executable checklist

## Test Results
- **Tests Executed**: No
- **Passed**: 20
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] gen-test-scripts/SKILL.md has a 'Forbidden Patterns' or 'Antipattern Guard' section
- [x] Each of the 6 antipatterns is listed with: what it is, why it's harmful, what to do instead
- [x] The section is referenced in the generation flow (e.g., 'Before writing each test, verify it does not match any forbidden pattern')
- [x] References the lesson documents as authoritative sources

## Notes
Documentation-only task. No code changes. Coverage set to -1.0 as this is a documentation task with no testable code output.
