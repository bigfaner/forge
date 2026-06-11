---
id: "2"
title: "Inline contract model into gen-journeys + clean Related + reduce summaries"
priority: "P0"
estimated_time: "2h"
dependencies: []
type: "doc"
complexity: "high"
mainSession: false
---

# 2: Inline contract model into gen-journeys + clean Related + reduce summaries

## Description
gen-journeys 引用 gen-contracts/rules/journey-contract-model.md（3次），需将 Contract 结构定义 + Outcome 语义内联。同时删除 Related Skills/Integration 章节，合并 References 到内联知识，并精简 5 个 per-surface 内联摘要。

## Reference Files
- `docs/proposals/skill-command-independence-audit/proposal.md` — Scope > In Scope, Key Risks
- plugins/forge/skills/gen-journeys/SKILL.md: 需内联 contract model、删除 Related、合并 References、精简 per-surface 摘要 (ref: Scope)
- plugins/forge/skills/gen-contracts/rules/journey-contract-model.md: Contract 结构定义 + Outcome 语义 ~60 行需内联 (ref: Scope)

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| plugins/forge/skills/gen-journeys/SKILL.md | 内联 Contract 结构定义，删除 Related/Integration，合并 References，精简 per-surface 摘要 |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] gen-journeys/SKILL.md 包含内联的 Contract 结构定义 + Outcome 语义（~60 行），带 `<!-- INLINE:origin=... -->` 标记
- [ ] Related Skills、Integration 章节已删除
- [ ] References 中的概念定义已合并到内联知识
- [ ] 5 个 per-surface 内联摘要已精简，保留关键差异信息
- [ ] gen-journeys 不再包含对 gen-contracts 内部文件的跨 skill 引用

## Hard Rules
- 内联段落必须添加 `<!-- INLINE:origin=gen-contracts/rules/journey-contract-model.md -->` 标记以提供可追溯性
- 保留所有 HARD-RULE / HARD-GATE / EXTREMELY-IMPORTANT / PROHIBITIONS 块计数不变

## Implementation Notes
INJECT: Contract 结构定义 + Outcome 语义 ~60 行；SKIP: 代码示例和实现细节。此任务涉及 4 种不同类型的变更（内联、删除、合并、精简），AC=5 导致 high complexity。
