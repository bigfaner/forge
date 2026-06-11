---
status: "completed"
started: "2026-05-20 17:09"
completed: "2026-05-20 17:11"
time_spent: "~2m"
---

# Task Record: 7 P3: coding-refactor fmt处理简化 + coding-enhancement同步注释

## Summary
Simplified coding-refactor.md Step 4 fmt failure handling by removing complex git stash flow, replaced with clear decision logic. Added SYNC NOTICE comment to coding-enhancement.md header to remind maintainers to sync changes with coding-feature.md.

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/prompt/data/coding-refactor.md
- forge-cli/pkg/prompt/data/coding-enhancement.md

### Key Decisions
- Simplified fmt failure handling from git stash-based baseline comparison to a simple ownership check: if refactor-touched files have fmt issues, fix them; if only pre-existing files changed, continue
- Used HTML comment for sync notice in coding-enhancement.md so it does not interfere with template rendering at runtime

## Test Results
- **Tests Executed**: No
- **Passed**: 0
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] coding-refactor.md no longer contains git stash && just fmt && git diff --name-only && git stash pop complex git operations
- [x] coding-enhancement.md header includes sync maintenance comment
- [x] Simplified fmt handling still covers key scenarios (strategy for when fmt fails)

## Notes
无
