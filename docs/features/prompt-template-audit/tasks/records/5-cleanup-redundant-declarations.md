---
status: "completed"
started: "2026-05-20 16:57"
completed: "2026-05-20 17:00"
time_spent: "~3m"
---

# Task Record: 5 P1: 清理未使用占位符声明 + coding-refactor 格式统一

## Summary
清理模板中声明但未使用的占位符变量：doc.md 和 doc-eval.md 移除 FEATURE_SLUG 声明，doc-consolidate.md 和 doc-drift.md 移除 SCOPE 声明；将 coding-refactor.md 的 CODING_PRINCIPLES 从 ### 标题格式统一为 - 列表格式（与 coding-feature.md 一致）

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/prompt/data/doc.md
- forge-cli/pkg/prompt/data/doc-eval.md
- forge-cli/pkg/prompt/data/doc-consolidate.md
- forge-cli/pkg/prompt/data/doc-drift.md
- forge-cli/pkg/prompt/data/coding-refactor.md

### Key Decisions
- FEATURE_SLUG 仅在 doc-summary.md 中有实际使用（构建 records 路径），其余 doc 模板仅声明未消费，安全移除
- SCOPE 在 coding.* 和 validation.* 模板中有实际使用，但在 doc-consolidate 和 doc-drift 中仅声明未消费，安全移除
- coding-refactor.md CODING_PRINCIPLES 原有三个 bullet 点合并为两个 list 条目（Surgical Changes + Scope Limits），保持语义完整

## Test Results
- **Tests Executed**: No
- **Passed**: 0
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] doc.md 和 doc-eval.md 的 header 不再包含 FEATURE_SLUG 行
- [x] doc-consolidate.md 和 doc-drift.md 的 header 不再包含 SCOPE 行
- [x] coding-refactor.md 的 CODING_PRINCIPLES 使用 - 列表格式（与 coding-feature 一致）
- [x] 所有修改不影响模板的实际输出（移除的是声明但未使用的变量）

## Notes
无
