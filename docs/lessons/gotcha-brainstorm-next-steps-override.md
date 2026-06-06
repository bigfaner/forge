---
created: "2026-06-06"
tags: [architecture, testing]
---

# brainstorm skill 生成的提案 Next Steps 被 agent 自行替换

## Problem

brainstorm skill 生成的 proposal.md 中，Next Steps 应为 `Proceed to /write-prd to formalize requirements`（模板静态文本），但 agent 将其替换为 `Proceed to /tech-design for implementation details`，跳过了 PRD 阶段。

## Root Cause

1. **L1**: agent 在写提案时自行判断"这是 enhancement，需求已足够明确，不需要 PRD"，然后用自身判断覆盖了模板的 Next Steps 静态文本
2. **L2**: brainstorm skill 的 Step 5 仅说"Save to `docs/proposals/<slug>/proposal.md` using `templates/proposal.md`"，对 Next Steps 没有任何约束（无 HARD-RULE、无 GATE），agent 可以自由替换模板中的任何静态内容
3. **L3**: 模板 `templates/proposal.md` 的 Next Steps 是固定文本（非 `{{placeholder}}`），但 skill 没有声明"模板中的 Next Steps 不可修改"或"必须遵循标准 pipeline 流程"

## Solution

修正了当前提案的 Next Steps 为 `/write-prd`。

brainstorm skill 的 Step 5 应增加规则，确保 Next Steps 遵循标准 pipeline 流程，不被 agent 自行替换。

## Reusable Pattern

当 skill 模板中存在表示流程顺序的静态文本（如 Next Steps、Pipeline Position），应视为不可变指引，agent 不得基于自身判断替换。如果 skill 需要灵活处理，应使用 `{{placeholder}}` 并在 skill 中提供推导逻辑，而非让 agent 自行发挥。

## References

- brainstorm skill Step 5（`plugins/forge/skills/brainstorm/SKILL.md`）
- 提案模板 `templates/proposal.md` 第 101-102 行
- 标准 pipeline 流程：brainstorm → proposal → eval → write-prd → tech-design → breakdown-tasks
