---
status: "completed"
started: "2026-05-27 23:17"
completed: "2026-05-27 23:21"
time_spent: "~4m"
---

# Task Record: 2 Add complexity branching and search strategy to prompt templates

## Summary
Added complexity branching (<!-- IF NOT_LOW --> markers) around Step 1.5 and search strategy guidance paragraph to all 4 coding prompt templates (enhancement, feature, refactor, cleanup). coding-fix.md left untouched per Hard Rule.

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/prompt/data/coding-enhancement.md
- forge-cli/pkg/prompt/data/coding-feature.md
- forge-cli/pkg/prompt/data/coding-refactor.md
- forge-cli/pkg/prompt/data/coding-cleanup.md

### Key Decisions
- Search strategy guidance placed as standalone paragraph after <!-- END_IF --> and before Step 2, not inside conditional block
- coding-enhancement.md already had <!-- IF NOT_LOW --> markers from Task 1, only needed search strategy paragraph added

## Test Results
- **Tests Executed**: Yes
- **Passed**: 10
- **Failed**: 0
- **Coverage**: 83.8%

## Acceptance Criteria
- [x] coding-enhancement.md Step 1.5 wrapped in <!-- IF NOT_LOW -->...<!-- END_IF --> markers
- [x] coding-feature.md Step 1.5 wrapped in <!-- IF NOT_LOW -->...<!-- END_IF --> markers
- [x] coding-refactor.md Step 1.5 wrapped in <!-- IF NOT_LOW -->...<!-- END_IF --> markers
- [x] coding-cleanup.md Step 1.5 wrapped in <!-- IF NOT_LOW -->...<!-- END_IF --> markers
- [x] Search strategy guidance paragraph added between Step 1.5 and Step 2 in all 4 templates
- [x] Guidance content matches exact specified text
- [x] coding-fix.md is NOT modified
- [x] Templates remain valid markdown after changes

## Notes
coding-enhancement.md already had the conditional markers from Task 1 (COMPLEXITY placeholder). Only the search strategy paragraph was added to it. The other 3 templates received both changes.
