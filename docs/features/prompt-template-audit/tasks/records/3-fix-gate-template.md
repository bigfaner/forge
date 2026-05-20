---
status: "completed"
started: "2026-05-20 16:53"
completed: "2026-05-20 16:55"
time_spent: "~2m"
---

# Task Record: 3 P1: 修复 gate 模板缺失 conventions 和 PHASE_SUMMARY

## Summary
修复 gate.md 模板，添加 conventions 加载步骤（docs/conventions/ 和 docs/business-rules/，按 domains 字段过滤）和 PHASE_SUMMARY 条件加载语句，使其与 coding-feature.md 和 validation-code.md 的写法保持一致

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/prompt/data/gate.md

### Key Decisions
- conventions 加载写法完全对齐 coding-feature.md 和 validation-code.md 的 Step 1 格式
- PHASE_SUMMARY header 声明已存在（第4行），仅需补充 Step 1 中的条件加载语句

## Test Results
- **Tests Executed**: No
- **Passed**: 0
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] gate.md header 包含 PHASE_SUMMARY 声明
- [x] Step 1 包含 PHASE_SUMMARY 条件加载语句
- [x] Step 1 包含 docs/conventions/ 和 docs/business-rules/ 的加载指令（按 domains 字段过滤相关性）
- [x] conventions 加载方式与其他模板一致（读 frontmatter domains 字段判断相关性）

## Notes
无
