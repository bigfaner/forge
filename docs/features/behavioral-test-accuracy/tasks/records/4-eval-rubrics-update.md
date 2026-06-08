---
status: "completed"
started: "2026-06-08 18:13"
completed: "2026-06-08 18:16"
time_spent: "~3m"
---

# Task Record: 4 eval rubrics 新增 Workflow Coverage + Fixture Specification 评估维度

## Summary
Journey eval rubric 新增 Workflow Coverage 维度 (150pts, 60% min threshold, Golden Path veto); Contract eval rubric 新增 Fixture Specification 维度 (100pts, 60% min threshold, entity completeness veto)。两个维度均含语义验证要求（防 checkbox-compliant），总分和目标分数已同步调整。

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/eval/rubrics/journey.md
- plugins/forge/skills/eval/rubrics/contract.md

### Key Decisions
无

## Document Metrics
Journey: 6→7 dimensions, 1000→1150pts; Contract: 7→8 dimensions, 1000→1100pts. Both new dimensions include veto items and semantic verification requirements.

## Referenced Documents
- docs/proposals/behavioral-test-accuracy/proposal.md

## Review Status
final

## Acceptance Criteria
- [x] Journey eval rubric 新增 Workflow Coverage 维度 (150 分)，含 Golden Path 存在性和多步覆盖度子项
- [x] Workflow Coverage 维度最低通过阈值 ≥90/150 (60%)，Golden Path 存在性一票否决；eval prompt 要求语义验证
- [x] Contract eval rubric 新增 Fixture Specification 维度 (100 分)，含前置数据声明完整性和实体关系覆盖度
- [x] Fixture Specification 维度最低通过阈值 ≥60/100 (60%)，entities 完整性一票否决；eval prompt 要求验证 entity_type 与领域模型一致

## Notes
两个新维度均实现了 backward-compatible 处理：Journey Workflow Coverage 对无 Golden Path 的旧 Journey 通过 veto 机制处理；Contract Fixture Specification 对缺失 fixture_spec 的旧 Contract 免罚。总分调整保持了原有维度的分值比例不变。
