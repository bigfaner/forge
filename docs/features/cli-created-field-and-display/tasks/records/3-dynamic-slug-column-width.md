---
status: "completed"
started: "2026-05-21 01:03"
completed: "2026-05-21 01:09"
time_spent: "~6m"
---

# Task Record: 3 列表命令 slug 列动态宽度

## Summary
将 feature list、lesson、proposal 三个命令的 slug/name 列从固定宽度改为动态宽度。新增 calcSlugColWidth 函数实现 clamp(max(30, maxSlugLen+2), 60) 逻辑，新增 padRight 辅助函数，修改三个列表命令的表头、分隔线、数据行使用动态宽度。42 字符 slug 现可完整显示。

## Changes

### Files Created
- forge-cli/internal/cmd/slug_width_test.go

### Files Modified
- forge-cli/internal/cmd/proposal.go
- forge-cli/internal/cmd/feature.go
- forge-cli/internal/cmd/lesson.go
- forge-cli/scripts/version.txt

### Key Decisions
- calcSlugColWidth 接受 []int（slug长度列表）而非泛型接口，保持简单
- padRight 替代 fmt %-*s 动态宽度格式化（Go 的 %-Ns 中 N 不能是变量）
- mapper 函数（mapProposalsToSlugLens 等）就近定义在各命令文件中，避免跨包依赖

## Test Results
- **Tests Executed**: Yes
- **Passed**: 31
- **Failed**: 0
- **Coverage**: 8.9%

## Acceptance Criteria
- [x] forge feature list 输出中无 slug 被截断（42 字符 slug 完整显示）
- [x] forge lesson 输出中无 name 被截断
- [x] forge proposal 输出中无 slug 被截断
- [x] 列宽 = clamp(max(30, maxSlugLen + 2), 60)
- [x] truncateSlug 函数根据动态宽度截断

## Notes
无
