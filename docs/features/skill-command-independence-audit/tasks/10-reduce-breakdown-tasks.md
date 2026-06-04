---
id: "10"
title: "Reduce shared content in breakdown-tasks"
priority: "P1"
estimated_time: "1h"
dependencies: []
type: "doc"
complexity: "medium"
mainSession: false
---

# 10: Reduce shared content in breakdown-tasks

## Description
breakdown-tasks 与 quick-tasks 共享 ~150 行内容，需精简共享部分，保留 breakdown-tasks 专属逻辑。

## Reference Files
- `docs/proposals/skill-command-independence-audit/proposal.md` — Scope > In Scope, Key Risks
- plugins/forge/skills/breakdown-tasks/SKILL.md: 精简与 quick-tasks 共享的内容 (ref: Scope)
- plugins/forge/skills/quick-tasks/SKILL.md: 共享内容参考 (ref: Scope)

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| plugins/forge/skills/breakdown-tasks/SKILL.md | 精简与 quick-tasks 共享的 ~150 行内容 |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] 与 quick-tasks 共享的 ~150 行内容已精简，保留 breakdown-tasks 专属逻辑
- [ ] breakdown-tasks 仍能完整指导 AI agent 执行 task breakdown

## Implementation Notes
breakdown-tasks 有完整的 PRD+Design pipeline 上下文，与 quick-tasks 的共享内容主要集中在 task 结构定义和 priority 分类。
