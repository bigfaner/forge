---
status: "completed"
started: "2026-05-22 22:20"
completed: "2026-05-22 22:22"
time_spent: "~2m"
---

# Task Record: 1 Add Impact Declaration to coding-refactor.md template

## Summary
在 coding-refactor.md 模板的 Step 2 末尾新增 Impact Declaration 子步骤（子步骤 6），要求执行器在动手前对受影响测试进行 PRESERVE/EVOLVE 分类并输出结构化声明。同步更新 Step 3 Universal constraints 和 Step 4 targeted test 失败处理，区分 EVOLVE（更新断言）和 PRESERVE/未声明（BEHAVIOR_CHANGE_DETECTED + 暂停）。

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/prompt/data/coding-refactor.md

### Key Decisions
- Impact Declaration 作为 Step 2 的子步骤 6 追加，不修改现有 Impact Mapping 逻辑（步骤 1-5）
- EVOLVE 条目必须同时包含 reason 和 expected_change，否则降级为 PRESERVE，防止过度声明
- 声明格式包含具体示例（TestAddCmd_BlockSource），降低执行器遵循门槛

## Test Results
- **Tests Executed**: No
- **Passed**: 0
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] Step 2 末尾新增 Impact Declaration 子步骤，要求执行器输出 IMPACT_DECLARATION 结构化声明
- [x] 声明格式包含 test name、classification (PRESERVE/EVOLVE)、reason、expected_change
- [x] EVOLVE 条目缺少 reason 或 expected_change 时，重新分类为 PRESERVE
- [x] Step 3 Universal constraints 更新：EVOLVE 测试失败时更新断言而非输出 BEHAVIOR_CHANGE_DETECTED
- [x] Step 4 targeted test 失败处理更新：区分 EVOLVE（更新）、PRESERVE/未声明（暂停）
- [x] 声明输出在执行记录中可见（Step 2 Output 行包含 impact_declaration 摘要）
- [x] go build ./... 通过

## Notes
模板变更不影响编译。BEHAVIOR_CHANGE_DETECTED 机制完整保留用于 PRESERVE/未声明测试。声明格式包含具体示例，符合 Hard Rules 要求。
