---
id: "3"
title: "Update quick-tasks and breakdown-tasks SKILL.md type tables"
priority: "P2"
estimated_time: "20min"
dependencies: ["1"]
type: "doc"
mainSession: false
---

# 3: Update quick-tasks and breakdown-tasks SKILL.md type tables

## Description

更新 `quick-tasks/SKILL.md` 和 `breakdown-tasks/SKILL.md` 的类型分配表：移除 `gate`（系统类型），添加 `doc.consolidate` 和 `doc.drift` 作为合法业务类型并附使用指导。

## Reference Files
- `docs/proposals/system-type-exclusion/proposal.md` — Source proposal
- `plugins/forge/skills/quick-tasks/SKILL.md` — quick-tasks 类型表（~line 88-99）
- `plugins/forge/skills/breakdown-tasks/SKILL.md` — breakdown-tasks 类型表（~line 94-102）

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/quick-tasks/SKILL.md` | 移除 gate 行，添加 doc.consolidate 和 doc.drift 行 |
| `plugins/forge/skills/breakdown-tasks/SKILL.md` | 移除 gate 行，添加 doc.consolidate 和 doc.drift 行 |

## Acceptance Criteria

- [ ] quick-tasks/SKILL.md 类型表中不含 `gate`
- [ ] breakdown-tasks/SKILL.md 类型表中不含 `gate`
- [ ] 两个类型表均包含 `doc.consolidate` 并附使用场景说明
- [ ] 两个类型表均包含 `doc.drift` 并附使用场景说明

## Hard Rules

- 修改 SKILL.md 前，必须先加载 `docs/conventions/forge-distribution.md` 了解分发模型约束

## Implementation Notes

- `doc.consolidate` 使用场景：用户为老项目手动创建 consolidate 任务（将分散的规范文件合并到 docs/business-rules/ 或 docs/conventions/）
- `doc.drift` 使用场景：用户手动创建 drift 审计任务（检测现有规范与代码的不一致）
- 两个 SKILL.md 的类型表结构一致，修改需同步
