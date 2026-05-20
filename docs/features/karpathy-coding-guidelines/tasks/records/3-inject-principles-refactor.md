---
status: "completed"
started: "2026-05-20 00:51"
completed: "2026-05-20 00:52"
time_spent: "~1m"
---

# Task Record: 3 Inject Surgical Changes principle into coding-refactor template

## Summary
Added CODING_PRINCIPLES block with Surgical Changes principle to coding-refactor.md template, positioned after External behavior definition and before Pre-check section

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/prompt/data/coding-refactor.md

### Key Decisions
- Included only Surgical Changes principle (not Think Before Coding) since Impact Mapping step already covers the 'think first' concern
- Principle text references Impact Map (Step 2) to create cross-link between behavior constraint and scope mapping
- Kept principle concise (3 bullet points) to avoid inflating the already complex refactor template

## Test Results
- **Tests Executed**: No
- **Passed**: 0
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] Contains CODING_PRINCIPLES block with only the Surgical Changes principle
- [x] Positioned after External behavior definition block, before Pre-check
- [x] No overlap with Impact Mapping step content
- [x] Pre-check section preserved
- [x] Step numbering unchanged (1/4, 2/4, 3/4, 4/4)
- [x] Template placeholders undisturbed
- [x] ONLY Surgical Changes principle included, not Think Before Coding

## Notes
无
