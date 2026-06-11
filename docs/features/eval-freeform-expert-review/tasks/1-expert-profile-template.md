---
id: "1"
title: "Expert Profile Template & Inference Prompt"
priority: "P0"
estimated_time: "1h"
dependencies: []
type: "doc"
mainSession: false
---

# 1: Expert Profile Template & Inference Prompt

## Description

创建动态专家档案的模板和推断 prompt。这是 Phase 0 的基础——自由专家评审依赖于此模板生成的专家档案。

根据提案，系统需要分析提案内容（domain、技术栈、复杂度、关键决策），推断最适合评审的专家档案（背景、专业领域、评审风格），并通过 AskUserQuestion 让用户确认（接受 / 修改 / 重新生成，最多 3 轮修改）。

## Reference Files
- `docs/proposals/eval-freeform-expert-review/proposal.md` — Source proposal

## Affected Files

### Create
| File | Description |
|------|-------------|
| `plugins/forge/skills/eval/experts/freeform/expert-template.md` | 专家档案 YAML front matter + Markdown 正文模板 |
| `plugins/forge/skills/eval/experts/freeform/expert-inference.md` | 专家推断 prompt（分析提案 → 生成专家档案） |

### Modify
| File | Changes |
|------|---------|

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria

- [ ] 专家档案模板包含提案规定的 6 个必填字段：domain、background、review_style、generated_for、created_at、review_history
- [ ] 模板格式为 YAML front matter + Markdown 正文，兼容现有 `experts/scorer/*.md` 的 prompt 格式
- [ ] 推断 prompt 指导 LLM 从提案中提取 domain、技术栈、复杂度、关键决策，推断专家背景
- [ ] 推断 prompt 包含用户确认机制（AskUserQuestion 三选一：接受 / 修改 / 重新生成）
- [ ] 推断 prompt 限制最多 3 轮修改，超过后提示接受或跳过
- [ ] 推断 prompt 包含降级逻辑：用户连续 3 次拒绝 → 提示手动输入或跳过

## Hard Rules

- 专家档案必须包含可验证的领域关键词和背景描述（应对 hallucinated expertise 风险）
- 当用户无法判断专家领域胜任度时，推断 prompt 应生成交叉引用（专家关键词 vs 提案技术术语）和 3-5 个自检问题

## Implementation Notes

- 参考现有 `experts/scorer/cto.md` 的格式风格
- 推断 prompt 是子 agent 的指令文档，会被拼接到 agent prompt 中
