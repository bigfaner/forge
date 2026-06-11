---
status: "completed"
started: "2026-05-20 13:17"
completed: "2026-05-20 13:18"
time_spent: "~1m"
---

# Task Record: 3 Update quick-tasks and breakdown-tasks SKILL.md type tables

## Summary
更新 quick-tasks/SKILL.md 和 breakdown-tasks/SKILL.md 的类型分配表：移除 gate（系统类型），添加 doc.consolidate 和 doc.drift 作为合法业务类型并附使用场景说明。同步更新了质量门分类表，将 Meta/gate 分类移除，doc.consolidate 和 doc.drift 归入 Doc 分类。

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/quick-tasks/SKILL.md
- plugins/forge/skills/breakdown-tasks/SKILL.md

### Key Decisions
- gate 是系统类型，由 forge task index 自动生成，不应出现在业务任务的类型分配表中
- doc.consolidate 和 doc.drift 是双重身份类型（既可自动生成也可手动创建），作为合法业务类型保留在表中
- 质量门分类表中移除了 Meta/gate 分类行，Doc 分类扩展为包含 doc、doc.consolidate、doc.drift 三种类型

## Test Results
- **Tests Executed**: No
- **Passed**: 0
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] quick-tasks/SKILL.md 类型表中不含 gate
- [x] breakdown-tasks/SKILL.md 类型表中不含 gate
- [x] 两个类型表均包含 doc.consolidate 并附使用场景说明
- [x] 两个类型表均包含 doc.drift 并附使用场景说明

## Notes
无
