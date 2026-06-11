---
id: "11"
title: "Reduce proposal-only features in eval"
priority: "P1"
estimated_time: "30m"
dependencies: []
type: "doc"
complexity: "low"
mainSession: false
---

# 11: Reduce proposal-only features in eval

## Description
eval SKILL.md 334 行中约 100 行为 proposal-only 特性描述，需精简。

## Reference Files
- `docs/proposals/skill-command-independence-audit/proposal.md` — Scope > In Scope, Key Risks
- plugins/forge/skills/eval/SKILL.md: 精简 proposal-only 特性描述 (ref: Scope)

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| plugins/forge/skills/eval/SKILL.md | 精简 proposal-only 特性描述 ~100 行 |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] proposal-only 特性描述已精简，保留通用 eval 功能说明
- [ ] 所有 eval 子类型（17 种 rubric）的触发条件和执行流程完整保留

## Implementation Notes
eval 是通用评估 skill，proposal-only 特性指仅适用于 eval-proposal 子类型的展开内容。精简时保留 freeform 评估的动态专家生成逻辑。
