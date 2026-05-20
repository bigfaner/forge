---
status: "completed"
started: "2026-05-20 13:18"
completed: "2026-05-20 13:24"
time_spent: "~6m"
---

# Task Record: 3 Update coding-type prompt templates with coverage directives

## Summary
更新所有 5 个 coding 类型 prompt 模板，将裸 COVERAGE_STRATEGY/COVERAGE_TARGET 标签替换为测试步骤中 <IMPORTANT> 块内的可操作覆盖率指令。feature/enhancement 使用 percentage 策略加'达到目标后停止补充测试'；fix 增加针对性修复测试指令；refactor 增加 maintain 策略、不新增测试指令和增量编译策略（修改文件→just compile→通过则继续→失败则修复当前文件）；cleanup 增加 maintain 策略和不新增测试指令。新增 5 个测试验证各模板指令可操作性，全部现有测试保持通过。

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/prompt/data/coding-feature.md
- forge-cli/pkg/prompt/data/coding-enhancement.md
- forge-cli/pkg/prompt/data/coding-fix.md
- forge-cli/pkg/prompt/data/coding-refactor.md
- forge-cli/pkg/prompt/data/coding-cleanup.md
- forge-cli/pkg/prompt/prompt_test.go

### Key Decisions
- 将裸 COVERAGE_STRATEGY/COVERAGE_TARGET 标签从模板顶部移除，改为嵌入各模板测试步骤的 <IMPORTANT> 块中，确保 agent 优先看到
- 在 <IMPORTANT> 块中同时包含策略标识符和目标文本，保持与现有 Task 2 测试的兼容性
- refactor 模板的增量编译策略格式化为：修改文件 → just compile → 通过则继续 → 失败则立即修复当前文件，不要继续修改其他文件

## Test Results
- **Tests Executed**: Yes
- **Passed**: 34
- **Failed**: 0
- **Coverage**: 89.4%

## Acceptance Criteria
- [x] coding-feature.md 模板：测试步骤包含 {{COVERAGE_TARGET}} 指令，agent 知道达到目标后停止补充测试
- [x] coding-enhancement.md 模板：同 feature，使用 {{COVERAGE_TARGET}}
- [x] coding-fix.md 模板：测试步骤包含 {{COVERAGE_TARGET}}，agent 写针对性修复测试而非追求高覆盖
- [x] coding-refactor.md 模板：使用 maintain 策略指令 + 增量编译策略
- [x] coding-cleanup.md 模板：使用 maintain 策略指令，不新增测试
- [x] 所有模板更新后，forge prompt get-by-task-id 能正确渲染覆盖率指令
- [x] 现有测试通过

## Notes
无
