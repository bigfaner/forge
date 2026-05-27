---
created: "2026-05-27"
tags: [architecture]
---

# Prompt 模板按 type 选择，忽略任务复杂度差异

## Problem

Task 3（quality_gate 路径替换，6 条 AC，纯文本替换）执行耗时 25 分钟。大量时间浪费在读取无关文档（conventions、business-rules、proposal 全文）和做 5 维度 spec-code conflict scan 上。

## Root Cause

因果链（3 层）：

1. **表面现象**：task-executor 对简单任务做了过重的探索（Step 1 读 conventions + business-rules + task file + 所有 Reference Files，Step 1.5 做 5 维度 scan），25 分钟中探索占比过高
2. **直接原因**：prompt 模板 `coding-enhancement.md` 对所有 `coding.enhancement` 任务强制执行完整的 Step 1 + Step 1.5 流程，不区分任务复杂度。Task 3 的 Reference Files 指向 `proposal.md#Section-Title`，导致 executor 还读了 proposal
3. **根因**：prompt 模板只按 type（enhancement/feature/refactor/cleanup/fix）选择，但同一 type 内任务复杂度差异巨大。简单机械替换（6 条 grep + sed 级别的 AC）和复杂功能增强共享同一个重探索模板。模板缺少复杂度分支机制

## Solution

在 prompt 模板中增加复杂度分支：

1. **轻量判定条件**：如果 Hard Rules 为空 且 AC ≤ 3 条 且 Reference Files ≤ 1 个 → 跳过 Step 1.5 spec-code scan，简化 Step 1（只读 task file + Reference Files，跳过 conventions/business-rules）
2. **替代方案**：在 task frontmatter 加 `complexity: low/medium/high`，quick-tasks 生成时根据 AC 数量和 Reference Files 数量自动判定，模板根据 complexity 调整探索深度

## Reusable Pattern

- **模板粒度 ≠ type 粒度**：按 type 选模板（5 种 coding type）太粗。真正影响探索深度的是任务复杂度，不是 type。一个 enhancement 可能是单行路径替换，也可能是跨模块功能重构
- **复杂度快速判定**：AC 数量 + Hard Rules 是否为空 + Reference Files 数量 → 三者组合足以判定任务是否需要重探索
- **探索成本的隐性代价**：Step 1.5 的 5 维度 scan 对简单任务是纯浪费。每次 scan 都要读代码文件 → 对比 spec → 输出 checklist，即使结论全是 "N/A"

## Related Files

- `forge-cli/pkg/prompt/data/coding-enhancement.md` — 当前模板，对所有 enhancement 任务强制完整探索
- `forge-cli/pkg/prompt/data/coding-cleanup.md` — cleanup 类型模板，可能存在同样问题
- [[gotcha-quick-tasks-merge-threshold]] — 任务拆分粒度问题（合并标准错误导致任务 scope 过大）
- [[gotcha-task-executor-thinking-overhead]] — task-executor 87% 时间在 thinking 的根因分析
