---
status: "completed"
started: "2026-05-15 18:08"
completed: "2026-05-15 18:37"
time_spent: "~29m"
---

# Task Record: 2 Rewrite SKILL.md: Scope Inference, split rules, cleanup

## Summary
Rewrote quick-tasks/SKILL.md Step 2 to replace file-path-based Scope Assignment with semantic Scope Inference. Added split-by-functional-steps rule and Reference Files directory-level fill rule. Cleaned Step 3 to remove 'Fill Affected Files' instruction. Updated Template Selection to use description semantics instead of Affected Files. Removed Affected Files check item from Output Checklist.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/quick-tasks/SKILL.md

### Key Decisions
- Template Selection now uses task description semantics instead of Affected Files classification, matching the proposal's information density
- Scope Inference is purely semantic (no file reading) to eliminate over-research of downstream skill files
- Split-by-functional-steps rule added alongside the existing per-bullet rule to handle multi-step In Scope items

## Test Results
- **Tests Executed**: No
- **Passed**: 0
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] Step 2 no longer contains 'Determine affected file paths' instruction or Scope Assignment algorithm
- [x] Step 2 contains Scope Inference rules (frontend/backend/all from task description semantics)
- [x] Step 2 contains split-by-functional-steps rule (multi-step In Scope items -> multiple tasks, total <= 10)
- [x] Step 2 contains Reference Files fill rule (directory-level paths from proposal In Scope)
- [x] Step 3 no longer references 'Fill Affected Files'
- [x] Output Checklist no longer has Affected Files check item
- [x] forge task index still generates valid index.json after changes
- [x] executor agents can work from Description + Reference Files without needing file-level paths

## Notes
Documentation-only task, no testable code changes. Template Selection table also updated to use description semantics instead of Affected Files since the implementation template no longer has that section.
