---
id: "8"
title: "Delete Related sections + reduce redundancy in write-prd"
priority: "P1"
estimated_time: "1h"
dependencies: []
type: "doc"
complexity: "medium"
mainSession: false
---

# 8: Delete Related sections + reduce redundancy in write-prd

## Description
write-prd 的 Related Skills/Integration/References 章节需删除。同时精简 4 种 intent 变体的重复展开。

## Reference Files
- `docs/proposals/skill-command-independence-audit/proposal.md` — Scope > In Scope, Key Risks, Success Criteria
- plugins/forge/skills/write-prd/SKILL.md: 删除 Related 章节并精简 intent 变体 (ref: Scope)

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| plugins/forge/skills/write-prd/SKILL.md | 删除 Related Skills/Integration/References，精简 intent 变体 |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] Related Skills、Integration、References 章节已删除
- [ ] 4 种 intent 变体的重复展开已精简，保留核心差异
- [ ] 与 tech-design 共享的 Override Signals 表已精简为 write-prd 专属版本

## Implementation Notes
write-prd 与 tech-design 共享的 Override Signals 表需各自保留与自身 skill 相关的信号。
