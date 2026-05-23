---
id: "3"
title: "Extraction Prompt Template & Injection Mechanism"
priority: "P0"
estimated_time: "1h"
dependencies: ["2"]
type: "doc"
mainSession: false
---

# 3: Extraction Prompt Template & Injection Mechanism

## Description

创建从自由评审叙事中提取结构化 key findings 的 prompt 模板，以及将发现注入 rubric scorer 的机制。这是自由评审与 rubric 评审之间的桥梁。

根据提案，提取流程为：自由评审叙事 + 提取 prompt → LLM 提取 → JSON 校验 → 注入 rubric scorer prompt。包含完整提取失败（降级）和部分提取失败（命中率告警）的处理逻辑。

## Reference Files
- `docs/proposals/eval-freeform-expert-review/proposal.md` — Source proposal

## Affected Files

### Create
| File | Description |
|------|-------------|
| `plugins/forge/skills/eval/experts/freeform/extraction-prompt.md` | 提取 prompt 模板（含 System/User 角色和完整指令） |
| `plugins/forge/skills/eval/rules/freeform-injection.md` | 注入机制规则（如何将 findings 追加到 scorer prompt） |

### Modify
| File | Changes |
|------|---------|

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria

- [ ] 提取 prompt 完整模板包含：System 角色、User 角色、输出格式（JSON 数组）、字段定义（summary/severity/quote）、规则（仅显式风险、不合并、severity 基于语气、quote 逐字引用）
- [ ] 提取 prompt 中使用 `{{FREEFORM_REVIEW}}` 占位符
- [ ] 注入规则定义：key findings 以 attack points 列表形式追加到 scorer prompt 末尾，scorer 被要求在评分时回应这些 attack points
- [ ] 注入规则定义 `[beyond-rubric]` 标签：无法映射到 rubric 维度的发现以 `[beyond-rubric]: [finding]` 格式记录在 ATTACKS 列表末尾
- [ ] 注入规则定义矛盾标记：自由评审发现与 rubric 方向矛盾时，scorer 在评分报告中标注「自由评审与 rubric 存在分歧」
- [ ] 定义 JSON 校验规则：提取产出非空、JSON 格式合法、每个元素 summary/severity/quote 三字段均非空、severity 枚举值为 high/medium/low
- [ ] 定义降级逻辑：提取产出为空或格式校验失败 → 跳过注入，降级为标准 rubric 流程，告知用户
- [ ] 定义部分提取失败处理：命中率（成功提取数 / 关键词段落数）< 50% → 报告中标注「提取命中率低」+ 附加完整叙事
- [ ] 定义命中率的粗粒度启发式估算方法及其局限性说明

## Hard Rules

- 注入的内容必须明确标识为「来自自由专家评审」以区分 rubric 自身的 attack points
- 提取 prompt 不得推断隐含风险——仅提取评审者明确表述的风险

## Implementation Notes

- 提取 prompt 模板参考提案中已定义的完整模板
- 注入规则修改 `rules/scorer-composition.md` 的 prompt 组合逻辑
