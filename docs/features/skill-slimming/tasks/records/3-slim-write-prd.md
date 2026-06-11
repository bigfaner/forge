---
status: "completed"
started: "2026-05-20 13:51"
completed: "2026-05-20 13:56"
time_spent: "~5m"
---

# Task Record: 3 Slim write-prd (407→≤350 lines)

## Summary
将 write-prd/SKILL.md 从 407 行拆分至 231 行，抽取知识提取规则、UI Functions 规则和自检规则到 rules/ 子目录

## Changes

### Files Created
- plugins/forge/skills/write-prd/rules/knowledge-extraction.md
- plugins/forge/skills/write-prd/rules/ui-functions.md
- plugins/forge/skills/write-prd/rules/self-check.md

### Files Modified
- plugins/forge/skills/write-prd/SKILL.md

### Key Decisions
- 将 Step 12 Knowledge Review 的完整知识提取流程（~142 行）移至 rules/knowledge-extraction.md，SKILL.md 仅保留一行引用
- 将 Step 8 UI Functions 的 Placement 规则、Navigation Architecture 和下游影响规则移至 rules/ui-functions.md
- 将 Step 9.5 Self-Check 的检查表格移至 rules/self-check.md
- 与 Task 1 (consolidate-specs) 和 Task 2 (tech-design) 保持一致的拆分风格：rules/ 子目录存放规则细节

## Test Results
- **Tests Executed**: No
- **Passed**: 0
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] SKILL.md 行数 ≤ 350 行
- [x] 所有步骤编号及其描述保留
- [x] 条件分支逻辑和 I/O 契约保留
- [x] 引用的 rules/templates 路径均存在可读
- [x] 拆分风格与 Task 1、2 保持一致

## Notes
无
