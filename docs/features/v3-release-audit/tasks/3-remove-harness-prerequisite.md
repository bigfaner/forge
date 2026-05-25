---
id: "3"
title: "Remove harness rubric type from Prerequisites tables"
priority: "P0"
estimated_time: "30m"
dependencies: ["1", "2"]
type: "doc"
mainSession: false
---

# 3: Remove harness rubric type from Prerequisites tables

## Description
部分 SKILL.md 的 Prerequisites 表引用了 `harness` 类型 rubric，但 harness rubric 文件缺失。按提案推荐方案 (c)：从 Prerequisites 表移除 harness 类型引用，保留其他 rubric 类型不变。不创建 harness.md（超范围），不修改运行时行为。

## Reference Files
- `proposal.md#Scope` — P0.5 defines harness decision: recommend option (c) remove from Prerequisites
- `proposal.md#Key-Risks` — harness rubric non-compliance risk rated L/L

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/eval/SKILL.md` | Remove harness type from Prerequisites table |
| `plugins/forge/skills/eval/rules/rubric-reference.md` | Remove harness type references |
| `plugins/forge/skills/eval/rules/pre-processing.md` | Remove harness type references |
| `plugins/forge/skills/forensic/SKILL.md` | Remove harness type from Prerequisites if present |

## Acceptance Criteria
- [ ] 所有 SKILL.md Prerequisites 表不含 harness 类型
- [ ] harness 相关文件保留不删除
- [ ] `grep -ri "harness" plugins/forge/skills/*/SKILL.md` 仅剩非 Prerequisites 上下文引用（如有）

## Hard Rules
- 纯删除操作，不引入新内容
- 保留 rubric、guide 等其他类型的 Prerequisites 引用

## Implementation Notes
方案 (c) 是最安全路径：纯删除，不涉及运行时变更或新文件创建。
