---
status: "completed"
started: "2026-06-06 17:04"
completed: "2026-06-06 17:17"
time_spent: "~13m"
---

# Task Record: 5 拆分 detect_surface.go 提取信号表

## Summary
将 detect_surface.go (962行) 拆分为 detect_surface.go (484行) + detect_surface_signals.go (485行)。信号映射表、信号检测函数、依赖源追踪函数集中到 signals 文件；编排逻辑、推断函数保留在主文件。

## Changes

### Files Created
- forge-cli/pkg/forgeconfig/detect_surface_signals.go

### Files Modified
- forge-cli/pkg/forgeconfig/detect_surface.go

### Key Decisions
- 同包文件拆分，不涉及跨包移动，外部调用方零影响
- 信号映射表 + 信号检测函数 + 冲突解决 + 辅助函数归入 detect_surface_signals.go
- 编排逻辑 + 推断函数归入 detect_surface.go
- 压缩 DetectSurfacesWithConflicts 的注释使其 <= 80 行

## Test Results
- **Tests Executed**: Yes
- **Passed**: 103
- **Failed**: 0
- **Coverage**: 85.6%

## Acceptance Criteria
- [x] detect_surface.go 和 detect_surface_signals.go 各 <= 500 行
- [x] 信号映射表集中在 detect_surface_signals.go
- [x] go test ./... 全绿，零行为变更
- [x] 所有函数 <= 80 行

## Notes
无
