---
id: "2"
title: "提取 createFixTask 共享 helper 函数"
priority: "P0"
estimated_time: "1h"
complexity: "medium"
dependencies: [1]
surface-key: "."
surface-type: "cli"
breaking: true
type: "coding.refactor"
mainSession: false
---

# 2: 提取 createFixTask 共享 helper 函数

## Description

将 `addSingleFixTask` 的 task 创建逻辑（surface inference、template defaults、opts 构造、task 创建、markdown 生成、state 更新）提取为共享 helper 函数 `createFixTask`。`addRegressionFixTasks` 和 `addSingleFixTask` 均调用此 helper，仅 `addSingleFixTask` 执行 cap 检查。共享 helper 须有独立单元测试覆盖，确保两条调用路径行为一致。

## Reference Files
- `forge-cli/internal/cmd/quality_gate.go`: addSingleFixTask（L709-803）提取 L726-803 为 createFixTask helper (source: proposal.md#In-Scope)
- `forge-cli/internal/cmd/quality_gate.go`: inferSurface（L728）需包含在 helper 中 (source: proposal.md#In-Scope)
- `forge-cli/internal/cmd/quality_gate_test.go`: 现有测试作为重构后行为一致性验证基线 (source: proposal.md#Key-Risks)

## Acceptance Criteria
- [ ] 新建 `createFixTask` 函数，封装 task 创建核心逻辑（surface inference、opts 构造、AddTask、CreateTaskMarkdown、EnsureForgeState），接收 title、sourceFiles、output、errorDocPath、step 等参数
- [ ] `addSingleFixTask` 重构为：cap 检查（countFixTasks + maxFixTasksPerStep）+ 调用 `createFixTask`
- [ ] `addSingleFixTask` 原有行为不变（compile/fmt/lint/unit-test 路径的 fix task 创建结果与重构前一致）
- [ ] `createFixTask` 有独立单元测试覆盖 task 字段填充、markdown 生成、state 更新

## Implementation Notes

### Test Impact
- Affected test suite(s): forge-cli/internal/cmd/
- Expected fixture changes: 无，重构不改变外部行为
- Risk level: medium

- 重构先于新功能（Task 4）合并，确保现有 compile/fmt/lint/unit-test 路径不受影响
- `createFixTask` 签名设计需考虑 Task 4 的调用需求（regression 路径不传 sourceFiles 逗号列表，传按文件分组的 description）
- opts 构造中的 Vars（SOURCE_FILES、TEST_SCRIPT、TEST_RESULTS、SOURCE_TASK_ID）需支持自定义覆盖
