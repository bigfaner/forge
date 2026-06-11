---
id: "4"
title: "Fix ARCHITECTURE.md factual errors"
priority: "P0"
estimated_time: "1h"
dependencies: ["1", "2"]
type: "doc"
mainSession: false
---

# 4: Fix ARCHITECTURE.md factual errors

## Description
ARCHITECTURE.md 存在 6 个 Critical 事实性错误：声称 4 agents（实际 1）、引用 PostToolUse hook（不存在）、过时路径、组件计数错误。修正已存在内容，不新增子系统概述（P1.12 范围）。

## Reference Files
- `proposal.md#Problem` — Evidence table: ARCHITECTURE.md 6 Critical + 3 Minor errors
- `proposal.md#Scope` — P0.2 defines ARCHITECTURE.md fix scope (existing content only)
- `proposal.md#Scope` — P1.12 is separate: 9 subsystem overviews (out of scope for this task)

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `docs/ARCHITECTURE.md` | Fix 6 Critical errors (agent count, PostToolUse, paths) + 3 Minor formatting issues |

## Acceptance Criteria
- [ ] Agent 计数与 `ls plugins/forge/agents/ | wc -l` 一致
- [ ] Hook 列表与 `ls plugins/forge/hooks/` 一致
- [ ] Skill 计数与 `ls plugins/forge/skills/ | wc -l` 一致
- [ ] 无 PostToolUse 引用（`grep -c "PostToolUse" docs/ARCHITECTURE.md` = 0）
- [ ] 路径引用指向实际存在的目录

## Hard Rules
- 仅修正事实性声明，不重构文档结构
- 不新增内容（子系统概述属于 Task 11）
- 逐条交叉验证，不改未审计部分

## Implementation Notes
修正后 ARCHITECTURE.md 作为 P1.12（Task 11）的基础，需确保现有内容 100% 准确后再扩展。
