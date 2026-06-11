---
status: "completed"
started: "2026-06-06 16:15"
completed: "2026-06-06 16:28"
time_spent: "~13m"
---

# Task Record: 2 拆分 BuildIndex 390 行上帝函数

## Summary
将 BuildIndex 390 行上帝函数拆分为 12 个命名步骤函数，分布到 build.go (347行) 和 build_steps.go (425行) 两个文件。BuildIndex 本体从 390 行降至 54 行。所有函数 <= 80 行，所有文件 <= 500 行。导出 API 签名不变，同包拆分零行为变更。

## Changes

### Files Created
- forge-cli/pkg/task/build_steps.go

### Files Modified
- forge-cli/pkg/task/build.go

### Key Decisions
- 引入 buildContext 结构体封装步骤间共享状态，步骤函数作为其方法
- 拆分为 build.go（导出API + 辅助函数）和 build_steps.go（buildContext + 步骤函数）两个文件
- 提取 upsertTaskFromFile 和 upsertAutoEntry 消除 stage-gate 索引中的代码重复
- 提取 validateDocTaskAC 使 detectPipelineNeedsAndAC 更简洁

## Test Results
- **Tests Executed**: Yes
- **Passed**: 513
- **Failed**: 0
- **Coverage**: 86.2%

## Acceptance Criteria
- [x] BuildIndex 及其所有提取出的子函数均 <= 80 行
- [x] build.go 文件 <= 500 行
- [x] go test ./... 全绿，零行为变更
- [x] 导出 API 签名不变（同包拆分，不改包名或导入路径）

## Notes
BuildIndex 从 390 行降至 54 行。build_steps.go 包含 12 个步骤方法 + buildContext 结构体。go test -race ./... 全部通过。
