---
status: "completed"
started: "2026-06-06 16:43"
completed: "2026-06-06 17:03"
time_spent: "~20m"
---

# Task Record: 4 拆分 pipeline.go 提取校验逻辑

## Summary
拆分 pipeline.go (1103 行) 为 pipeline.go (332 行) + pipeline_validate.go (484 行)。提取校验逻辑、PipelineRegistry、Dependency Resolvers 和 ID 匹配函数到 pipeline_validate.go，pipeline.go 保留核心类型定义、Gate/Condition 函数和 GenerateTestTasks 生成逻辑。零行为变更，所有测试通过。

## Changes

### Files Created
- forge-cli/pkg/task/pipeline_validate.go

### Files Modified
- forge-cli/pkg/task/pipeline.go

### Key Decisions
- 将 Dependency Resolvers 移至 pipeline_validate.go，因为它们主要为 PipelineRegistry 服务且被校验逻辑引用
- 将 ID 匹配函数 (matchRegistryID, matchTypeSuffixedID, matchSurfaceKeyID) 与校验逻辑同文件放置，因为 idExistsInRegistry 被 matchIDToTemplate 引用
- PipelineRegistry 移至 pipeline_validate.go 与 init() 校验同文件，保持注册表+校验的紧密内聚

## Test Results
- **Tests Executed**: Yes
- **Passed**: 513
- **Failed**: 0
- **Coverage**: 86.5%

## Acceptance Criteria
- [x] pipeline.go 和 pipeline_validate.go 各 <= 500 行
- [x] var 块和类型定义集中放置，不穿插在函数间打断阅读流
- [x] go test ./... 全绿，零行为变更
- [x] 所有函数 <= 80 行

## Notes
原始 pipeline.go 1103 行拆分为 pipeline.go 332 行 + pipeline_validate.go 484 行。文件按职责划分：pipeline.go = 类型+Gate+Condition+生成逻辑，pipeline_validate.go = Resolvers+Registry+匹配+校验。
