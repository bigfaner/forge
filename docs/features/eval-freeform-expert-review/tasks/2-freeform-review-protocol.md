---
id: "2"
title: "Freeform Review Protocol & Agent Prompt"
priority: "P0"
estimated_time: "1h"
dependencies: ["1"]
type: "doc"
mainSession: false
---

# 2: Freeform Review Protocol & Agent Prompt

## Description

创建自由叙事评审的协议和子 agent prompt。该专家以纯叙事形式对提案进行深度评审——无 rubric、无评分、无预设维度，完全由专家自主决定关注什么。

根据提案，评审产出为纯叙事格式，不包含任何 rubric 维度或评分。评审结果用于后续提取 key findings。

## Reference Files
- `docs/proposals/eval-freeform-expert-review/proposal.md` — Source proposal

## Affected Files

### Create
| File | Description |
|------|-------------|
| `plugins/forge/skills/eval/experts/freeform/freeform-review-protocol.md` | 自由评审协议（约束、产出格式、质量要求） |
| `plugins/forge/skills/eval/experts/freeform/freeform-reviewer.md` | 自由评审子 agent prompt（组合协议 + 动态专家档案） |

### Modify
| File | Changes |
|------|---------|

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria

- [ ] 自由评审协议明确定义：纯叙事格式、无 rubric、无评分、无预设维度
- [ ] 协议要求评审产出中的显式风险点使用明确语言标注（便于后续提取）
- [ ] 协议包含结构化段落框架（背景评估、关键风险识别、改进建议），确保 Jaccard 相似度 ≥ 0.6 的确定性 NFR
- [ ] 子 agent prompt 使用低 temperature（0.3）减少随机性
- [ ] 子 agent prompt 组合方式：协议 + 动态专家档案（来自 Task 1 的模板）
- [ ] 评审产出保存路径定义：`<doc_dir>/eval/freeform-review.md`

## Hard Rules

- 自由评审协议不得引用或暗示任何 rubric 维度（防止锚定效应）
- 评审产出中风险点必须使用明确的语言标注（如「风险：」「问题：」「建议：」），便于提取

## Implementation Notes

- 参考现有 `experts/protocol/scorer-protocol.md` 的三阶段结构，但简化为纯叙事流程
- 子 agent prompt 是 general-purpose agent 的指令，model: sonnet
