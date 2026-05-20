---
status: "completed"
started: "2026-05-20 16:55"
completed: "2026-05-20 16:57"
time_spent: "~2m"
---

# Task Record: 4 P1: test.* 模板 HARD-RULE 标签重命名为 TASK-CONSTRAINTS

## Summary
将 7 个 test.* 模板中的 HARD-RULE 标签重命名为 TASK-CONSTRAINTS，消除与 task 文件中 Hard Rules 概念的歧义；同时从 test-eval-cases.md 中移除 MAIN_SESSION 的 EXTREMELY-IMPORTANT 声明

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/prompt/data/test-gen-cases.md
- forge-cli/pkg/prompt/data/test-eval-cases.md
- forge-cli/pkg/prompt/data/test-gen-scripts.md
- forge-cli/pkg/prompt/data/test-run.md
- forge-cli/pkg/prompt/data/test-gen-and-run.md
- forge-cli/pkg/prompt/data/test-graduate.md
- forge-cli/pkg/prompt/data/test-verify-regression.md

### Key Decisions
- section 标题从 'Hard Rules' 同步改为 'Task Constraints' 保持上下文一致
- test-eval-cases.md 中整个 'Task-Specific Rules' 节和 EXTREMELY-IMPORTANT 块一并移除，因为 MAIN_SESSION 由 dispatcher 在分派时决定

## Test Results
- **Tests Executed**: No
- **Passed**: 0
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] 所有 7 个 test.* 模板中 HARD-RULE 替换为 TASK-CONSTRAINTS
- [x] test-eval-cases.md 中 MAIN_SESSION 相关的 EXTREMELY-IMPORTANT 声明已移除
- [x] 标签内的约束内容不变，仅标签名变更

## Notes
无
