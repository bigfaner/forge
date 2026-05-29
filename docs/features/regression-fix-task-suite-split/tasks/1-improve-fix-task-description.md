---
id: "1"
title: "Phase 0: 改进 fix task description 信息完整性"
priority: "P1"
estimated_time: "1h"
complexity: "low"
dependencies: []
surface-key: "."
surface-type: "cli"
breaking: false
type: "coding.enhancement"
mainSession: false
---

# 1: Phase 0: 改进 fix task description 信息完整性

## Description

当前 `addSingleFixTask` 使用 `just.ExtractConciseError(output, 10)` 取 output 尾部作为 concise error 放入 task description。当 regression 失败数量多且输出长时，concise error 只展示尾部，agent 看不到完整的 `--- FAIL:` 列表，需要额外读取 raw-output.txt 才能定位所有失败。

改为从 output 中提取所有 `--- FAIL:` 行列表替代 tail 截断，使 agent 无需读取 raw-output.txt 即可看到全部失败信息。

## Reference Files
- `forge-cli/internal/cmd/quality_gate.go`: addSingleFixTask（L709-803）中 description 生成逻辑（L736-741），需将 `just.ExtractConciseError(output, 10)` 替换为 `--- FAIL:` 行提取 (source: proposal.md#Phase-0)
- `forge-cli/pkg/just/just.go`: ExtractConciseError 函数，理解当前 tail 截断行为，作为 fallback 保留 (source: proposal.md#Phase-0)

## Acceptance Criteria
- [ ] `addSingleFixTask` 的 description 中 "Concise error" 部分替换为从 output 中提取的所有 `--- FAIL:` 行列表（包含完整失败名称和持续时间）
- [ ] 当 output 中无 `--- FAIL:` 行时（如 compile/fmt/lint 步骤），fallback 到现有 `just.ExtractConciseError` tail 行为，不产生空 description
- [ ] 新增单元测试覆盖：含 `--- FAIL:` 行的 output 提取、无 `--- FAIL:` 行的 output fallback、空 output 处理

## Implementation Notes
- 提取逻辑：遍历 output 按行匹配 `--- FAIL:` 前缀，收集所有匹配行，拼接为 description 中的 concise error 部分
- 非测试步骤（compile/fmt/lint）的输出不含 `--- FAIL:` 行，自动 fallback 到 tail 行为，零回归
