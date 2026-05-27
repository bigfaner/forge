---
created: "2026-05-27"
tags: [architecture]
---

# Task-executor 不可见 thinking 时间导致 52 分钟执行简单任务

## Problem

Task 3（quality_gate 路径替换 + mobile-test-setup 集成，6 条 AC）执行了 52 分 26 秒。forensic 工具报告 4.2 分钟，与实际 wall-clock 时间差距 12 倍。48 分钟的模型 thinking 时间完全不可见。

## Root Cause

因果链（3 层）：

1. **表面现象**：task-executor 执行 Task 3 耗时 52 分钟，远超预期的 5-10 分钟。forensic 工具 `forge forensic extract` 报告 duration 为 4.2 分钟，与实际 wall-clock 52 分钟严重不符
2. **直接原因**：executor 在 48 次 tool call 之间产生了大量 thinking 时间。forensic 工具通过 JSONL entry 时间戳差值计算 thinking 时间，但只捕获到 251.8 秒（4.2 分钟）。剩余 48 分钟的 thinking 发生在 API 调用过程中（模型在生成 tool call 前的 extended thinking），这些时间在 JSONL 中没有对应 entry
3. **根因**：两层问题叠加：
   - **Forensic 盲区**：`forge forensic extract` 只计算 JSONL entry 间的时间差。模型的 extended thinking 发生在 API 请求-响应周期内，这个时间在 JSONL 中体现为一个 entry 的写入时间，无法被差值计算捕获。52 分钟中只有 4.2 分钟（8%）被 forensic 看见
   - **Sonnet 模型的 extended thinking 开销**：sonnet 模型每次 API 调用都进行 extended thinking（��部推理链），每次调用耗时 2-10 秒的 thinking。48 次 tool call × 平均 60 秒 thinking = 48 分钟。这个 per-call thinking 开销是 task-executor 使用 sonnet 模型的结构性成本

## Solution

1. **Forensic 修复**：用 raw JSONL 的首尾 entry 时间戳计算 wall-clock duration，而非从 thinking turns 加总
2. **模型选择优化**：对于简单任务（AC ≤ 3，纯文本替换），考虑使用 haiku 模型减少 per-call thinking 开销
3. **Prompt 精简**：减少不必要的探索步骤，降低 tool call 数量（当前 26 个 Bash 中大部分是 grep/find/ls，可用更少的调用完成）

## Reusable Pattern

- **Forensic 工具有盲区**：`forge forensic extract` 报告的 duration 可能严重低估实际时间。检查方法：对比 raw JSONL 首尾时间戳 vs forensic 报告的 duration，如果差距 >2 倍，说明存在大量不可见 thinking
- **Tool call 数量 × per-call thinking 时间 = 总不可见时间**：48 calls × ~60s/call ≈ 48 分钟。减少 tool call 数量直接减少不可见时间
- **简单任务用重模型是浪费**：路径替换任务不需要 deep reasoning，但 sonnet 的 extended thinking 每次调用都触发。按任务复杂度选模型比统一用 sonnet 高效得多

## Related Files

- `forge-cli/internal/cmd/forensic.go` — forensic extract 的 duration 计算
- [[gotcha-task-executor-thinking-overhead]] — 87% thinking 时间的根因分析
- [[gotcha-prompt-template-complexity-agnostic]] — prompt 模板不区分任务复杂度导致过度探索
