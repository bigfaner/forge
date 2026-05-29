---
id: "4"
title: "实现 addRegressionFixTasks 按测试文件拆分 fix task"
priority: "P0"
estimated_time: "2h"
complexity: "high"
dependencies: [2, 3]
surface-key: "."
surface-type: "cli"
breaking: false
type: "coding.feature"
mainSession: false
---

# 4: 实现 addRegressionFixTasks 按测试文件拆分 fix task

## Description

新建 `addRegressionFixTasks` 函数和 `isTestFile` 函数，替代 `runTestRegression` 中 `addFixTask` 调用。`isTestFile` 识别 Go 测试文件（`*_test.go`），`addRegressionFixTasks` 使用 `extractFileLineMap` 按测试文件分组创建独立 fix task。每个 task 只包含该测试文件相关的输出行。

## Reference Files
- `forge-cli/internal/cmd/quality_gate.go`: runTestRegressionLegacy（L260）和 runTestRegressionSurface（L289）需替换 addFixTask 调用 (source: proposal.md#In-Scope)
- `forge-cli/internal/cmd/quality_gate.go`: countFixTasks（L614-628）title prefix `"fix test:"` 匹配逻辑需验证兼容新 title 格式 (source: proposal.md#In-Scope)
- `forge-cli/internal/cmd/quality_gate.go`: addFixTask（L648-673）作为 fallback 行为保留 (source: proposal.md#In-Scope)
- `docs/lessons/gotcha-quality-gate-fix-task-loop.md`: countFixTasks cap 机制已正常工作的验证记录 (source: proposal.md#Proposed-Solution)

## Acceptance Criteria
- [ ] 新建 `isTestFile(filename string) bool`，MVP 仅匹配 `*_test.go` 命名约定；新建 `addRegressionFixTasks`，调用 `extractFileLineMap` 获取文件→输出行映射，仅为包含直接 `--- FAIL:` 条目的测试文件创建独立 fix task，调用 `createFixTask` helper（Task 2 产物）
- [ ] 每个 fix task 的 title 格式为 `"fix test: <filename> failure in quality gate"`（如 `"fix test: handler_test.go failure in quality gate"`），使 `countFixTasks` 的 title prefix `"fix test:"` 匹配仍能正确计数
- [ ] 每个 fix task 的 description 包含该测试文件路径 + 筛选后的相关输出行（含上下文窗口）
- [ ] Regression 专用软上限 10 个独立 task + 1 个 overflow task（第 11 至 N 个文件的所有输出行合并为 1 个 task，title 为 `"fix test: regression overflow (N-10 files)"`），总 task 数 ≤ 11，不受 `maxFixTasksPerStep` 硬上限限制
- [ ] `runTestRegressionLegacy`（L260）和 `runTestRegressionSurface`（L289）中的 `addFixTask` 调用替换为 `addRegressionFixTasks`；当 `isTestFile` 返回零匹配时 fallback 到现有 `addFixTask` 行为，并输出结构化日志警告（`WARNING: isTestFile returned zero matches for output, falling back to directory-grouped fix task`）

## Implementation Notes

### Test Impact
- Affected test suite(s): forge-cli/internal/cmd/
- Expected fixture changes: 需构造包含多文件失败的 regression output 样本作为测试数据
- Risk level: medium

- 仅修改以下调用点：`runTestRegressionLegacy` L260 和 `runTestRegressionSurface` L289
- Overflow task 不使用 `groupFilesByDir`（该函数对单目录返回 nil），直接拼接剩余文件输出行
- 非测试文件路径的输出行和未归属行归入 overflow task（若无 overflow 则归入 fallback task）
- `countFixTasks` 的 title prefix `"fix test:"` 需验证覆盖新 title 格式 `"fix test: <filename> failure in quality gate"`：当前 prefix 为 `"fix " + step + ":"` 即 `"fix test:"`，新 title 以 `"fix test:"` 开头，匹配正确
