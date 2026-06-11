---
status: "completed"
started: "2026-05-20 00:49"
completed: "2026-05-20 00:50"
time_spent: "~1m"
---

# Task Record: 2 Inject Karpathy principles into coding-fix template, replacing existing rules

## Summary
Replaced the first <IMPORTANT> block (MINIMAL CHANGES + NO REFACTORING) in coding-fix.md with a <CODING_PRINCIPLES> block containing Think Before Coding, Simplicity First, and Surgical Changes principles. New principles provide stronger and more actionable guidance while covering the same constraints as the old rules.

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/prompt/data/coding-fix.md

### Key Decisions
- Adapted principle text for fix context: Think Before Coding emphasizes root cause verification over the generic version's assumption identification
- Simplicity First directly replaces MINIMAL CHANGES with explicit 'no speculative changes' and 'no while I'm here improvements' coverage
- Surgical Changes directly replaces NO REFACTORING with 'scope boundary = failing code path only' — more actionable than the original blanket rule

## Test Results
- **Tests Executed**: No
- **Passed**: 0
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] First <IMPORTANT> block replaced by <CODING_PRINCIPLES> with Think Before Coding + Simplicity First + Surgical Changes
- [x] Semantic coverage: new principles enforce at least same constraints as old rules
- [x] Second <IMPORTANT> block (Hard Rules) preserved unchanged
- [x] <CODING_PRINCIPLES> positioned after role description, before ## Workflow
- [x] Step numbering (1/4, 2/4, 3/4, 4/4) unchanged
- [x] Template placeholders undisturbed

## Notes
无
