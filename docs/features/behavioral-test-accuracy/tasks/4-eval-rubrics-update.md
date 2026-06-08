---
id: "4"
title: "eval rubrics 新增 Workflow Coverage + Fixture Specification 评估维度"
priority: "P1"
estimated_time: "1h"
dependencies: [1, 2]
type: "doc"
mainSession: false
---

# 4: eval rubrics 新增 Workflow Coverage + Fixture Specification 评估维度

## Description

当前 Journey eval rubric 不评估 Golden Path 存在性和工作流覆盖度，Contract eval rubric 不评估 Fixture Specification 完整性。本任务在两个 eval rubric 中新增对应评估维度，使行为性测试质量可度量。

Journey eval rubric 新增 "Workflow Coverage" 维度（150 分），包含 Golden Path 存在性（一票否决子项）和多步覆盖度评分。Contract eval rubric 新增 "Fixture Specification" 维度（100 分），包含前置数据声明完整性和实体关系覆盖度（完整性子项一票否决）。两个维度均设 60% 最低通过阈值。

## Reference Files
- `docs/proposals/behavioral-test-accuracy/proposal.md` — Success Criteria (SC-5, SC-6), Key Risks
- `plugins/forge/skills/eval/rubrics/journey.md` — Journey eval rubric，需新增 Workflow Coverage 维度 (ref: Success Criteria)
- `plugins/forge/skills/eval/rubrics/contract.md` — Contract eval rubric，需新增 Fixture Specification 维度 (ref: Success Criteria)

## Affected Files

### Create
| File | Description |
|------|-------------|
| |

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/eval/rubrics/journey.md` | 新增 "Workflow Coverage" 评估维度（150 分） |
| `plugins/forge/skills/eval/rubrics/contract.md` | 新增 "Fixture Specification" 评估维度（100 分） |

### Delete
| File | Reason |
|------|--------|
| |

## Acceptance Criteria
- [ ] Journey eval rubric 新增 "Workflow Coverage" 维度（150 分），评分标准包含 Golden Path 存在性子项和多步覆盖度子项
- [ ] Workflow Coverage 维度最低通过阈值 ≥90/150（60%），Golden Path 存在性子项不得为 0 分（一票否决）；eval prompt 要求评审者验证步骤序列是否对应 PRD/Design 中的具体用户故事
- [ ] Contract eval rubric 新增 "Fixture Specification" 维度（100 分），评分标准包含前置数据声明完整性和实体关系覆盖度
- [ ] Fixture Specification 维度最低通过阈值 ≥60/100（60%），entities 必须包含 Contract 涉及的所有实体类型（完整性子项一票否决）；eval prompt 要求评审者验证 entity_type 是否与 Design 中的领域模型一致

## Implementation Notes

- 新增维度时需同步调整 rubric 总分，确保其他维度分值比例合理
- 防 checkbox-compliant 机制：eval prompt 要求评审者进行语义验证（如检查 Golden Path 步骤是否对应用户故事），而非仅检查字段是否存在或步骤数量
- 两个维度的一票否决子项需在 rubric 中明确标记为 veto item
