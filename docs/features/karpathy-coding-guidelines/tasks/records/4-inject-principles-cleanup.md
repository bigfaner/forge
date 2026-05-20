---
status: "completed"
started: "2026-05-20 00:52"
completed: "2026-05-20 00:53"
time_spent: "~1m"
---

# Task Record: 4 Inject principles into coding-cleanup template

## Summary
Added CODING_PRINCIPLES block with Simplicity First + Surgical Changes principles to coding-cleanup.md template, positioned after role description and before Workflow section.

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/prompt/data/coding-cleanup.md

### Key Decisions
- Simplicity First wording adapted for cleanup context: 'Remove only what the task targets' instead of generic 'Implement only what the task requires'
- Surgical Changes includes 'note but do not fix' guidance for out-of-scope issues, consistent with refactor template pattern

## Test Results
- **Tests Executed**: No
- **Passed**: 0
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] Contains CODING_PRINCIPLES block with Simplicity First + Surgical Changes principles
- [x] Positioned after role description, before ## Workflow
- [x] No semantic overlap with existing Step 2 Make Improvements instructions
- [x] Step numbering (Step 1/3, 2/3, 3/3) unchanged
- [x] Template placeholders undisturbed

## Notes
Only Simplicity First + Surgical Changes included per Hard Rules. No additional principles.
