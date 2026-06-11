---
id: "1"
title: "Inline surface detection into gen-contracts + clean Related sections"
priority: "P0"
estimated_time: "1h"
dependencies: []
type: "doc"
complexity: "medium"
mainSession: false
---

# 1: Inline surface detection into gen-contracts + clean Related sections

## Description
gen-contracts 引用 gen-journeys/SKILL.md "Surface Detection" section（反向引用），需将 surface 检测规则内联到 gen-contracts，并清理 Related Skills/Integration/References 章节。消除 gen-journeys→gen-contracts 双向耦合的 gen-contracts 端。

## Reference Files
- `docs/proposals/skill-command-independence-audit/proposal.md` — Scope > In Scope, Key Risks
- plugins/forge/skills/gen-contracts/SKILL.md: 需内联 surface 检测规则并删除 Related/Integration/References 章节 (ref: Scope)
- plugins/forge/skills/gen-journeys/SKILL.md: "Surface Detection" section ~20 行需内联到 gen-contracts (ref: Scope)

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| plugins/forge/skills/gen-contracts/SKILL.md | 内联 surface 检测规则段落，删除 Related/Integration/References 章节，合并 References 概念定义到内联知识 |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] gen-contracts/SKILL.md 包含从 gen-journeys/SKILL.md 内联的 surface 检测规则段落（~20 行），带 `<!-- INLINE:origin=... -->` 标记
- [ ] Related Skills、Integration 章节已删除
- [ ] References 中的概念定义（Contract、Outcome、Semantic Descriptors 等 6 个概念）已合并到内联知识作为定义段落
- [ ] gen-contracts 不再包含对 gen-journeys 内部文件的跨 skill 引用

## Hard Rules
- 内联段落必须添加 `<!-- INLINE:origin=gen-journeys/SKILL.md#Surface Detection -->` 标记以提供可追溯性
- 保留所有 HARD-RULE / HARD-GATE / EXTREMELY-IMPORTANT / PROHIBITIONS 块计数不变

## Implementation Notes
INJECT: surface 检测规则段落 ~20 行；SKIP: surface 生成步骤。此为 gen-journeys→gen-contracts 反向引用的对应端，修复后消除双向耦合。
