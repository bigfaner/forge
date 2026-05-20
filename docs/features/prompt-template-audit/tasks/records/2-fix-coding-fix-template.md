---
status: "completed"
started: "2026-05-20 16:50"
completed: "2026-05-20 16:52"
time_spent: "~2m"
---

# Task Record: 2 P1: 修复 coding-fix 模板缺失和顺序问题

## Summary
修复 coding-fix.md 模板：(1) 在 header 添加 {{PHASE_SUMMARY}} 声明；(2) 在 Step 1 中添加 PHASE_SUMMARY 条件加载语句；(3) 将 conventions 加载移到 task 文件读取之前，与 coding-feature/enhancement/cleanup 模板保持一致的加载顺序

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/prompt/data/coding-fix.md

### Key Decisions
- 采用与 coding-feature/enhancement/cleanup 完全一致的 conventions -> task file -> PHASE_SUMMARY 条件加载三段式顺序

## Test Results
- **Tests Executed**: No
- **Passed**: 0
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] coding-fix.md header 包含 PHASE_SUMMARY 声明
- [x] Step 1 中包含 PHASE_SUMMARY 条件加载语句
- [x] Step 1 中 conventions 加载在 task 文件读取之前
- [x] 模板结构与 coding-feature/enhancement/cleanup 保持一致的 conventions 加载模式

## Notes
doc 类型任务，无代码变更需要测试
