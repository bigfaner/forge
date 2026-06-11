---
status: "completed"
started: "2026-05-20 16:47"
completed: "2026-05-20 16:49"
time_spent: "~2m"
---

# Task Record: 1 P0: 移除6个模板的显式submit步骤（双重提交）

## Summary
移除6个prompt模板中的显式submit步骤，消除双重提交问题。doc.md从4步缩减为3步，doc-eval/doc-summary/doc-consolidate/doc-drift/clean-code各从3步缩减为2步

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/prompt/data/doc.md
- forge-cli/pkg/prompt/data/doc-eval.md
- forge-cli/pkg/prompt/data/doc-summary.md
- forge-cli/pkg/prompt/data/doc-consolidate.md
- forge-cli/pkg/prompt/data/doc-drift.md
- forge-cli/pkg/prompt/data/clean-code.md

### Key Decisions
- 统一依赖task-executor agent的第8步自动提交，模板不再显式调用submit-task

## Test Results
- **Tests Executed**: No
- **Passed**: 0
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] 6个模板中不再包含 Skill(skill="forge:submit-task") 调用
- [x] 每个模板的步骤编号连续、无跳跃
- [x] 移除submit步骤后，模板剩余的步骤语义完整

## Notes
无
