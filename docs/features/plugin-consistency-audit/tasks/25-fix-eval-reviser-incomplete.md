---
id: "25"
title: "Fix: eval reviser-composition incomplete entries"
priority: "P2"
estimated_time: "15min"
dependencies: []
type: "doc"
complexity: "low"
mainSession: false
---

# 25: Fix: eval reviser-composition incomplete entries

## Description
eval skill 的 `rules/reviser-composition.md` 第 29-31 行 "Reviser Type-Specific Constraints" 中，`journey`、`contract` 和 `consistency` 条目有不完整句子："After reviser completes:" 后无后续内容。

## Reference Files
- `docs/features/plugin-consistency-audit/reports/02-skills-batch-a.md#EV-03`: P2 INCOMPLETE, reviser type-specific constraints 不完整 (source: Report 02)
- `plugins/forge/skills/eval/rules/reviser-composition.md`: 需修改的 L29-31 (source: audit finding)

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/eval/rules/reviser-composition.md` | 完成 journey/contract/consistency 的 "After reviser completes" 描述 |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] journey 条目 "After reviser completes" 后有完整描述（如 "compare journey structure against original workflow"）
- [ ] contract 条目有完整描述
- [ ] consistency 条目有完整描述
- [ ] 无悬空句子

## Hard Rules
- 仅修改 `plugins/forge/skills/eval/rules/reviser-composition.md`
