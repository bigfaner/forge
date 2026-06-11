---
id: "7"
title: "Delete Related sections + reduce redundancy in tech-design"
priority: "P1"
estimated_time: "1h"
dependencies: []
type: "doc"
complexity: "medium"
mainSession: false
---

# 7: Delete Related sections + reduce redundancy in tech-design

## Description
tech-design 的 Related Skills/Integration/References 章节内容可从正文隐含推断，需删除。同时精简 4 种 intent 变体的重复展开和 Override Signals 表。

## Reference Files
- `docs/proposals/skill-command-independence-audit/proposal.md` — Scope > In Scope, Key Risks, Success Criteria
- plugins/forge/skills/tech-design/SKILL.md: 删除 Related 章节并精简 intent 变体和 Override Signals (ref: Scope)

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| plugins/forge/skills/tech-design/SKILL.md | 删除 Related Skills/Integration/References，精简 intent 变体和 Override Signals 表 |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] Related Skills、Integration、References 章节已删除，且内容可从正文隐含推断
- [ ] 4 种 intent 变体的重复展开已精简，保留核心差异
- [ ] Override Signals 表已精简，去除与 write-prd 的重复内容

## Implementation Notes
tech-design 当前 445 行中 intent 变体膨胀是主要冗余源。精简时保留所有硬规则和决策表，只压缩描述性文字。
