---
status: "completed"
started: "2026-05-20 00:48"
completed: "2026-05-20 00:49"
time_spent: "~1m"
---

# Task Record: 1 Inject Karpathy principles into coding-feature and coding-enhancement templates

## Summary
Added Karpathy's 4 coding principles (Think Before Coding, Simplicity First, Surgical Changes, Goal-Driven Execution) wrapped in CODING_PRINCIPLES XML tags to both coding-feature.md and coding-enhancement.md templates, positioned after role description and before Workflow section.

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/prompt/data/coding-feature.md
- forge-cli/pkg/prompt/data/coding-enhancement.md

### Key Decisions
无

## Test Results
- **Tests Executed**: No
- **Passed**: 0
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] Both files contain <CODING_PRINCIPLES> block positioned after role description, before ## Workflow
- [x] Block includes all 4 principles: Think Before Coding, Simplicity First, Surgical Changes, Goal-Driven Execution
- [x] No timing conflict with Step 1 — Think Before Coding guides Step 1 behavior, does not insert a new step
- [x] No semantic overlap with existing <IMPORTANT> block — principles complement, not duplicate
- [x] Step numbering (Step 1/3, 2/3, 3/3) unchanged
- [x] Template placeholders ({{TASK_ID}}, {{TASK_FILE}}, {{SCOPE}}, {{PHASE_SUMMARY}}) undisturbed

## Notes
无
