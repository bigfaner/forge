---
id: "17"
title: "Fix: submit-task AC vs testsFailed rule conflict"
priority: "P1"
estimated_time: "15min"
dependencies: []
type: "doc"
complexity: "low"
mainSession: false
---

# 17: Fix: submit-task AC vs testsFailed rule conflict

## Description
submit-task SKILL.md "Common Rules" 说 "`acceptanceCriteria` with any `met: false` → rejected for `completed`"。但 `data/record-format-coding.md` 说 "`testsFailed > 0` with `completed` → auto-downgrade to `blocked`"。两条规则触发条件不同：coding tasks 可能 AC 全 pass 但 testsFailed > 0。需在 SKILL.md 中明确优先级。

## Reference Files
- `docs/features/plugin-consistency-audit/reports/04-skills-batch-c.md#C-30`: P1 CONFLICT, AC vs testsFailed 规则冲突 (source: Report 04)
- `plugins/forge/skills/submit-task/SKILL.md`: 需修改的 Common Rules 节 (source: audit finding)
- `plugins/forge/skills/submit-task/data/record-format-coding.md`: coding category 的 testsFailed 规则 (source: audit finding)

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/submit-task/SKILL.md` | 在 Common Rules 中明确 category-specific rules 与 common rules 的优先级关系 |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] Common Rules 明确说明：category-specific rules 与 common rules 冲突时以 category-specific 为准
- [ ] 或将两条规则合并为一条统一的 `completed` 准入条件

## Hard Rules
- 仅修改 `plugins/forge/skills/submit-task/SKILL.md`
