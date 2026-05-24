---
id: "2"
title: "Split eval SKILL.md (488→≤350 lines)"
priority: "P0"
estimated_time: "1.5h"
dependencies: []
type: "doc"
mainSession: false
---

# 2: Split eval SKILL.md (488→≤350 lines)

## Description
eval SKILL.md 当前 488 行，超标 39%（上限 350 行）。提取 freeform pipeline 逻辑到独立 rules/ 文件，使 SKILL.md 回归约束。eval 已有 9 个 rules 文件，新增 freeform pipeline rule 需与现有结构一致。

## Reference Files
- `proposal.md#Proposed-Solution` — defines eval SKILL.md split as part of P0.4
- `proposal.md#Scope` — P0.4 specifies freeform pipeline extraction target
- `proposal.md#Key-Risks` — SKILL.md split breaking reference chain, rollback via git stash

## Affected Files

### Create
| File | Description |
|------|-------------|
| `plugins/forge/skills/eval/rules/freeform-pipeline.md` | Freeform evaluation pipeline logic extracted from SKILL.md |

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/eval/SKILL.md` | Remove freeform pipeline inline content, add Load directive for new rule |

## Acceptance Criteria
- [ ] SKILL.md ≤ 350 行
- [ ] 新 rule 被 SKILL.md 通过 Load 引用（入度 ≥ 1）
- [ ] 拆分后 SKILL.md 流程完整，无断裂引用
- [ ] `wc -l plugins/forge/skills/eval/SKILL.md` ≤ 350

## Hard Rules
- 新 rule 文件需符合 skill-self-containment.md 规范
- freeform injection 和 freeform expert persistence rules 已存在，新 rule 需与它们协作不冲突
- 回滚方案：`git stash` 回归失败则降级为 P1

## Implementation Notes
eval 含 9 个现有 rules，freeform 相关已有 freeform-expert-persistence.md 和 freeform-injection.md。新 freeform-pipeline.md 承载流程编排逻辑，与现有两文件互补不重叠。
