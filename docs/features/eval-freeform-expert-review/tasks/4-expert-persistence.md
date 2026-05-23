---
id: "4"
title: "Expert Persistence, Reuse & Deprecation"
priority: "P1"
estimated_time: "1h"
dependencies: ["1"]
type: "doc"
mainSession: false
---

# 4: Expert Persistence, Reuse & Deprecation

## Description

定义专家档案的持久化、复用匹配和质量追踪机制。动态专家档案保存到 `docs/experts/` 全局目录，后续评审可复用已有专家，并追踪专家有效性。

## Reference Files
- `docs/proposals/eval-freeform-expert-review/proposal.md` — Source proposal

## Affected Files

### Create
| File | Description |
|------|-------------|
| `plugins/forge/skills/eval/rules/freeform-expert-persistence.md` | 专家持久化、复用匹配和废弃追踪规则 |

### Modify
| File | Changes |
|------|---------|

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria

- [ ] 定义 `docs/experts/` 目录结构：每个专家一个文件，文件名基于 domain 命名（如 `distributed-systems-architect.md`）
- [ ] 定义 YAML front matter 必填字段：domain、background、review_style、generated_for、created_at、review_history、deprecated（默认 false）
- [ ] 定义复用匹配逻辑：读取 `docs/experts/` 中所有非 deprecated 专家 → 比较专家 domain 与提案 domain 的关键词重叠度 → 匹配度最高的专家通过 AskUserQuestion 呈现给用户（接受复用 / 拒绝生成新专家）
- [ ] 定义质量追踪机制：每次使用专家后，在 review_history 中记录「是否有实质变化」（rubric 评分差异 ≥ 15 分或 attack points 变动）
- [ ] 定义废弃逻辑：连续 3 次无实质变化 → 自动标记 `deprecated: true` → 后续匹配时跳过
- [ ] 定义用户手动废弃：用户可手动将专家的 deprecated 字段设为 true

## Hard Rules

- 专家档案写入 `docs/experts/` 时必须遵守项目的 docs/ 目录约定
- 复用匹配的语义必须足够具体，使开发者能直接实现（定义关键词提取和重叠度计算方法）

## Implementation Notes

- 复用匹配是轻量级的关键词匹配，不需要向量搜索或 embedding
- 质量追踪数据存储在专家档案的 review_history 数组中
