---
status: "completed"
started: "2026-05-27 19:44"
completed: "2026-05-27 19:47"
time_spent: "~3m"
---

# Task Record: 6 更新 prompt/task 模板中的旧术语

## Summary
将 pkg/prompt/data/ 和 pkg/task/data/ 中测试相关模板的旧术语 'profile'/'active profile'/'profile resolution' 替换为 'surface' 术语

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/prompt/data/test-gen-scripts.md
- forge-cli/pkg/prompt/data/test-run.md
- forge-cli/pkg/task/data/test-gen-scripts.md

### Key Decisions
- 'active profile' -> 'surface': 'active profile' 指当前激活的 profile，对应新模型中的 surface
- 'profile resolution' -> 'surface resolution': profile 解析概念映射到 surface 解析
- 'active test profile' -> 'surface': 模板中 'the framework from the active test profile' 简化为 'the framework from the surface'

## Test Results
- **Tests Executed**: Yes
- **Passed**: 472
- **Failed**: 0
- **Coverage**: 82.7%

## Acceptance Criteria
- [x] pkg/prompt/data/test-gen-scripts.md 中旧术语已替换
- [x] pkg/prompt/data/test-run.md 中旧术语已替换
- [x] pkg/task/data/test-gen-scripts.md 中旧术语已替换
- [x] grep -rn 'profile' forge-cli/pkg/prompt/data/ forge-cli/pkg/task/data/ --include='*.md' 返回 0 结果

## Notes
仅修改了 3 个 markdown 模板文件（共 4 处替换），未改动 Go 代码。目标包测试全部通过（472 passed, 0 failed），compile/fmt/lint 全部通过。
