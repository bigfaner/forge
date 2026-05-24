---
id: "11"
title: "Add 9 subsystem overviews to ARCHITECTURE.md"
priority: "P1"
estimated_time: "2h"
dependencies: ["4"]
type: "doc"
mainSession: false
---

# 11: Add 9 subsystem overviews to ARCHITECTURE.md

## Description
ARCHITECTURE.md 缺失 v3.0.0 新增的 9 个子系统概述：surface detection、worktree、Convention、forensic、deep-research、clean-code、extract-design-md、test-guide、learn。每个子系统添加概述+架构角色+SKILL.md 链接。总计 ≤180 行新增内容（每子系统 ≤20 行）。

## Reference Files
- `proposal.md#Scope` — P1.12: defines 9 subsystems and ≤180 line budget
- `proposal.md#Key-Risks` — P1.12 scope creep risk M/M, mitigation via per-subsystem ≤20 line cap

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `docs/ARCHITECTURE.md` | Add 9 subsystem overview sections after existing content |

## Acceptance Criteria
- [ ] 9 个子系统各有独立概述段落
- [ ] 每段落含：架构角色描述 + SKILL.md 链接
- [ ] 新增内容 ≤ 180 行
- [ ] 每子系统 ≤ 20 行

## Hard Rules
- 严格 ≤ 180 行新增，≤ 20 行/子系统
- 不修改 Task 4 已修正的内容
- 概述为描述性文本，不含实现细节

## Implementation Notes
需先读取各子系统 SKILL.md 的前 20 行获取简要描述，再浓缩为架构角色概述。关键约束是范围控制——避免每个子系统写成完整文档。
