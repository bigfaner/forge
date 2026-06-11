---
status: "completed"
started: "2026-05-15 20:17"
completed: "2026-05-15 20:22"
time_spent: "~5m"
---

# Task Record: 3 Add Task Types reference section to README.md

## Summary
Add 'Task Types & Pipeline 参考' section to README.md with 13 task types table, Quick pipeline (T-quick-1~5), Full pipeline (T-test-1~5), fix-task command template with --block-source, profile-suffix convention, and gate/summary auto-generation rules.

## Changes

### Files Created
无

### Files Modified
- README.md

### Key Decisions
- Used tables throughout for compactness per acceptance criteria
- Captured all content removed from quick-tasks and breakdown-tasks SKILL.md files in tasks 1 and 2
- Added disclaimer note: '以下内容由 forge task index 自动生成，以 CLI 行为准'

## Test Results
- **Tests Executed**: No
- **Passed**: 0
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] New section 'Task Types & Pipeline 参考' added before '文档索引'
- [x] Contains 13 task types table (type + who generates + purpose)
- [x] Contains Quick pipeline responsibility chain (T-quick-1~5 with brief descriptions)
- [x] Contains Full pipeline responsibility chain (T-test-1~5 with brief descriptions)
- [x] Contains fix-task command template with --block-source explanation
- [x] Contains profile-suffix convention (single vs multiple profiles)
- [x] Contains gate/summary auto-generation rules (phases with >=2 business tasks)
- [x] No content lost from the SKILL.md files — everything removed in tasks 1 and 2 is present here
- [x] Section is compact (tables, not prose) and scannable

## Notes
Documentation-only task, no test execution needed. All 13 task types verified against TaskTypeRegistry in forge-cli/pkg/task/types.go.
