---
created: "2026-05-23"
tags: [architecture, testing]
---

# Freeform Expert Findings 是二阶影响——不直接触发修订

## Problem

在 eval-proposal 流程中，Phase 0 自由专家评审产出了 20 个结构化发现（5 high + 13 medium + 2 low），但这些发现并未直接用于修订 proposal。Reviser 只看到 Scorer 合并后的攻击点，无法区分哪些来自 rubric、哪些来自 freeform。

## Root Cause

### Level 0: 架构设计——Freeform findings 通过 Scorer 间接传递

Eval 流程的信息传递链是：

```
Freeform Review → 提取 findings → 注入 Scorer Prompt
                                        ↓
                                  Scorer 结合 rubric + findings → 生成 ATTACKS
                                                                  ↓
                                                            Reviser 修订文档
```

Freeform findings 被 `<injected-freeform-findings>` 块包裹，注入到 Scorer prompt 中。Scorer 的职责是将这些 findings 映射到 rubric 维度（作为 dimension attacks）或标为 `[beyond-rubric]`。Reviser 只接收合并后的 ATTACKS 列表。

### Level 1: 映射决策权在 Scorer，不在 Reviser

Scorer 决定每个 freeform finding 的命运：
- 映射到 rubric 维度 → 成为标准 attack point，Reviser 必须处理
- 标为 `[beyond-rubric]` → 仍然出现在 ATTACKS 中，但 Reviser 可能将其优先级低于 rubric 维度内的攻击点
- Scorer 忽略（认为不重要或已由 rubric 覆盖）→ 完全丢失

### Level 2: 实际案例中的信息损失

在 spec-authority-enforcement 的 eval 中：
- Freeform 评审发现"标记稀释效应"和"Agent 层职责混淆"这两个核心问题
- Iteration 1 的 Scorer 将"标记稀释"标为 `[beyond-rubric]`，将"职责混淆"也标为 `[beyond-rubric]`
- Reviser 处理了这两个 `[beyond-rubric]` 攻击点（说明 Reviser 确实会处理 beyond-rubric findings）
- 但理论上，如果 Scorer 选择忽略某些 freeform findings，这些发现将完全丢失

### Level 3: 设计权衡——间接传递是有意为之

这不是 bug，而是架构权衡：
- **优点**：Scorer 统一过滤和映射，避免 freeform findings 与 rubric 评分产生矛盾
- **缺点**：freeform findings 的信息量在 Scorer 映射过程中可能被压缩或丢失

## Solution

当前架构已能工作，但有两个改进方向：

1. **Reviser 可见性**：在 Reviser prompt 中附注"以下 ATTACKS 中部分源自自由专家评审（标记为 beyond-rubric），请同等对待"
2. **信息保真验证**：在 Phase 0 注入后，比较 freeform findings 的覆盖率和 scorer ATTACKS 中的对应率，如果 < 80% 则在 eval report 中标注

## Reusable Pattern

**评估流程中的信息传递链**：当评审流程包含多个阶段（如 freeform + rubric），下游修订者只接收中间层的合并结果时，需要：
- 意识到中间层（Scorer）有权忽略上游输入（freeform findings）
- 如果上游输入是高价值的（如领域专家评审），考虑让修订者直接访问上游输出
- 或在中间层的输出中保留上游信息的溯源标记

**反面模式**：假设多阶段评审的每个阶段都完整保留了前一阶段的所有发现。

## Related Files

- `plugins/forge/skills/eval/SKILL.md` — eval 主流程定义
- `plugins/forge/skills/eval/rules/freeform-injection.md` — freeform findings 注入规则
- `plugins/forge/skills/eval/rules/scorer-composition.md` — scorer prompt 组合规则
- `plugins/forge/skills/eval/experts/protocol/reviser-protocol.md` — reviser 协议
- `docs/proposals/spec-authority-enforcement/eval/freeform-review.md` — 实际案例的自由评审
- `docs/proposals/spec-authority-enforcement/eval/iteration-1.md` — scorer 输出（含 beyond-rubric 标记）
