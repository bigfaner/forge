---
status: "completed"
started: "2026-06-05 17:36"
completed: "2026-06-05 17:41"
time_spent: "~5m"
---

# Task Record: 4 eval-contract 增加锚点完整性和一致性检查

## Summary
为 eval-contract 增加第7维度 Anchor Integrity (锚点完整性)，包含锚点字段完整性检查、锚点值匹配检查、handbook 内部一致性检查（endpoint 冲突检测）。同步更新 SKILL.md、pre-processing.md、scorer-composition.md、report-format.md 以支持新维度的数据流和报告格式。

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/eval/rubrics/contract.md
- plugins/forge/skills/eval/SKILL.md
- plugins/forge/skills/eval/rules/pre-processing.md
- plugins/forge/skills/eval/rules/scorer-composition.md
- plugins/forge/skills/eval/rules/report-format.md

### Key Decisions
无

## Document Metrics
7 dimensions, 1000 pts total, Anchor Integrity dimension 100 pts (3 criteria), backward-compatible when handbook absent

## Referenced Documents
- docs/proposals/contract-technical-anchors/proposal.md
- plugins/forge/skills/eval/SKILL.md
- plugins/forge/skills/eval/rubrics/contract.md

## Review Status
final

## Acceptance Criteria
- [x] eval-contract 评分包含锚点字段完整性检查：当 handbook 存在时，Contract 缺少对应锚点字段扣分
- [x] handbook 内部一致性检查：检测同一 endpoint 的 method 冲突、路径冲突
- [x] 评分结果报告明确列出缺失的锚点字段（按 surface 类型分组）

## Notes
点数重新分配：Completeness 200->150, Surface Fitness 150->100，腾出100点给新维度 Anchor Integrity。向后兼容：handbook 不存在时该维度满分。
