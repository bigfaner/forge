---
created: "2026-05-29"
tags: [testing, architecture]
---

# Reviser Agent 无超时保护导致长时间运行

## Problem

eval 管道的 Reviser subagent 在处理 9 个 attack points 时运行接近 2 小时仍未完成。虽然 agent 确实在做有效修改（proposal.md 已被改动），但执行效率极低。

## Root Cause

**Level 1 — Edit 工具的 old_string 匹配开销**：Reviser 对 150 行文档执行 9 次 Edit 调用，每次需要找到唯一的 old_string。在多次 edit 后文档内容已变，old_string 可能不再匹配或不再唯一，导致 agent 反复尝试不同的匹配串。

**Level 2 — Quality check 导致过度 re-read**：Reviser protocol 要求 "Every attack point from the scorer has been addressed"，agent 在每次 edit 后可能 re-read 文件验证结果。9 个 attack points = 至少 9 次 edit + 9 次 re-read = 18 次 tool call，加上失败重试可能达到 30+ 次。

**Level 3 — 无 max-duration 约束**：Reviser protocol 只有 "Maximum 3 rounds of self-review" 的软约束，但没有 wall-clock 时间上限。Agent 在遇到 edit 困难时不会主动放弃或简化策略，而是持续尝试直到所有 attack points 都被处理。

**Level 4（根源）— 单 agent 处理所有 attack points 的架构瓶颈**：Reviser 设计为串行处理所有 attack points。当 attack points 数量多（>5）或文档长（>100 行）时，单个 agent 的上下文窗口压力增大，后期 edit 的效率指数级下降。

## Solution

1. **Reviser protocol 添加 max-duration 约束**：在 HARD-RULE 中增加 "如果执行超过 10 分钟仍未完成所有 attack points，优先处理 high-severity 项，跳过 medium/low 项并报告。"
2. **Attack points 批次处理**：当 attack points > 5 时，按 severity 排序，先处理 high，再处理 medium，跳过 low。如果 high 处理完后时间接近上限，停止并报告剩余未处理项。
3. **减少 re-read**：在 protocol 中明确 "Edit 后不要 re-read 验证——Edit 工具本身会在失败时报错。只有在下一个 edit 需要当前文档状态时才 re-read。"

## Reusable Pattern

当 eval 管道的 subagent（Scorer 或 Reviser）运行时间超过预期时：
- 检查 attack points 数量——>5 个时考虑分批或降级
- 检查文档长度——>100 行时 Edit 的 old_string 匹配成本上升
- 检查是否有 max-duration 保护——没有时 agent 会无限重试
- 如果 agent 实际在做有效修改（非卡死），可以等待而非中断，但应设 10-15 分钟上限

## Example

```
# 症状：Reviser 运行 2 小时未完成
# 原因：9 个 attack points × Edit 匹配困难 × 无超时保护
# 处理：手动中断，proposal.md 已包含大部分修改，手动验证后继续
```

## Related Files

- `plugins/forge/skills/eval/experts/protocol/reviser-protocol.md`
- `docs/proposals/intent-driven-pipeline-branching/proposal.md`
