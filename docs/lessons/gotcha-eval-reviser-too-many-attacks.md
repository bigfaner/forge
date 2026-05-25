---
created: "2026-05-25"
tags: [architecture, testing]
---

# Eval Reviser 处理大量 Attack Points 时耗时过长

## Problem

`/eval --type proposal` 流程中，Iteration 1 reviser 处理 10 个 attack points（4 high + 6 medium）对 ~650 行提案文档的修订，执行超过 1 小时仍未完成，最终被用户中断。

## Root Cause

多个因素叠加导致 reviser 单轮执行时间失控：

1. **Attack points 数量过多**：10 个 findings 全部交给单轮 reviser 处理，每个 finding 需要读取文档、定位编辑点、生成新内容、验证上下文一致性
2. **新内容生成而非微调**：多个 findings 要求新增完整段落（如 surface-orchestration.yaml 生命周期管理、scope 迁移顺序约束），不是简单的措辞修改，LLM 生成长文本耗时显著
3. **Edit 工具唯一性约束**：文档较长时，Edit 需要包含足够的上下文才能唯一定位 old_string，多次定位失败后触发重试
4. **Self-review 轮次放大**：reviser 协议允许最多 3 轮 self-review，每轮重新读取文档验证修改，与 attack points 数量乘积效应明显
5. **顺序处理无并行**：所有 attack points 在单个 agent 中顺序处理，无法利用并行能力

## Solution

**短期缓解**：减少单轮 reviser 处理的 attack points 数量。仅处理 high severity findings，medium 留待后续迭代或降低 MAX_ITERATIONS 从 3 到 2（减少总迭代次数）。

**结构性改进方向**：将 reviser 改为 batch 模式——每次 subagent 调用只处理 2-3 个 attack points，主会话按批次串行调度。代价是增加 Agent tool 调用次数，但单次调用耗时可控。

## Reusable Pattern

当 eval 管线的 reviser 需要处理大量 attack points 时：

- **预判耗时**：attack points > 5 且文档 > 500 行时，预期单轮 reviser 耗时 > 30 分钟
- **分批策略**：按 severity 分批——high 优先，medium 按相关性分组。每批 3-4 个 findings
- **迭代预算**：MAX_ITERATIONS = 2 足以覆盖 pre-revision + 1 轮正式修订。3 轮迭代在复杂提案上代价过高
- **监控中断**：reviser 执行超过 20 分钟时考虑中断并评估已完成的修改是否足够推进下一轮 scoring

## Related Files

- `plugins/forge/skills/eval/rules/reviser-composition.md`
- `plugins/forge/skills/eval/rubrics/proposal.md`
