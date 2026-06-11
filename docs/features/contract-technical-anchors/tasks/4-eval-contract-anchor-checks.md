---
id: "4"
title: "eval-contract 增加锚点完整性和一致性检查"
priority: "P1"
estimated_time: "1.5h"
dependencies: [1]
type: "doc"
mainSession: false
---

# 4: eval-contract 增加锚点完整性和一致性检查

## Description
eval-contract 评分规则增加技术锚点完整性检查（当 handbook 存在时，Contract 必须包含对应锚点字段）和 handbook 内部一致性检查（设计阶段检测 endpoint 冲突），从评分层面保证锚点质量。

## Reference Files
- `docs/proposals/contract-technical-anchors/proposal.md` — Scope, Success Criteria
- `plugins/forge/skills/eval/SKILL.md`: 增加锚点评分逻辑 (ref: Scope)
- `plugins/forge/skills/eval/rubrics/`: 评分标准文件，增加锚点完整性规则 (ref: Success Criteria)

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/eval/SKILL.md` | 增加锚点完整性评分逻辑 |
| `plugins/forge/skills/eval/rubrics/` | 增加锚点评分标准 |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] eval-contract 评分包含锚点字段完整性检查：当 handbook 存在时，Contract 缺少对应锚点字段扣分
- [ ] handbook 内部一致性检查：检测同一 endpoint 的 method 冲突、路径冲突
- [ ] 评分结果报告明确列出缺失的锚点字段（按 surface 类型分组）

## Implementation Notes
- 完整性检查仅在 handbook 存在时生效，handbook 不存在时不扣分（向后兼容）
- 内部一致性检查覆盖 api-handbook 的 endpoint 冲突检测
