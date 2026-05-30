---
id: "6"
title: "Consolidated audit report and effectiveness validation"
priority: "P0"
estimated_time: "1.5h"
dependencies: [2, 3, 4, 5]
type: "doc"
complexity: "high"
mainSession: true
---

# 6: Consolidated audit report and effectiveness validation

## Description

合并 Task 2-5 的所有审计发现，生成最终结构化问题报告。执行有效性验证（确认已知 run-tests 问题被复现）、五类分类覆盖率检查、P0-P3 严重等级分级、误报率抽检方案。报告标注基准 commit hash 和 AI 模型信息。

## Reference Files
- `docs/features/plugin-consistency-audit/reports/02-skills-batch-a.md`: Task 2 审计结果 (source: Task 2)
- `docs/features/plugin-consistency-audit/reports/03-skills-batch-b.md`: Task 3 审计结果 (source: Task 3)
- `docs/features/plugin-consistency-audit/reports/04-skills-batch-c.md`: Task 4 审计结果 (source: Task 4)
- `docs/features/plugin-consistency-audit/reports/05-commands-agent-hooks.md`: Task 5 审计结果 (source: Task 5)
- `docs/proposals/plugin-consistency-audit/proposal.md#Success-Criteria`: 有效性验证、误报率抽检、五类分类覆盖要求 (source: proposal.md#Success-Criteria)

## Affected Files

### Create
| File | Description |
|------|-------------|
| `docs/features/plugin-consistency-audit/reports/06-consolidated-report.md` | 最终合并审计报告 |

### Modify
| File | Changes |
|------|---------|

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] 所有审计发现已合并去重，按 P0→P3 排序输出
- [ ] **有效性验证通过**: run-tests 的 `rules/env-check.md` Playwright 硬编码问题在最终报告中作为 P1 级 CONFLICT 出现
- [ ] 五类分类（CONFLICT/REDUNDANT/TIMING/REFERENCE/INCOMPLETE）均至少有 1 个实例；若某类为 0，列出所有含多步骤流程的组件清单并确认已逐一验证
- [ ] 报告包含基准 commit hash、AI 模型版本、审计参数（temperature 等）
- [ ] 误报率抽检方案已定义：随机抽取 ≥20% 的 P0/P1 问题清单，标注待人工验证

## Implementation Notes
- 合并时去重逻辑：同一 (component, file_path, description) 三元组视为重复
- P0-P3 分级依据 proposal 的 Severity Level Definitions：P0 运行时错误/流程卡死，P1 行为偏差，P2 信息冗余，P3 风格/措辞
- 误报率抽检：列出 P0/P1 问题清单并随机抽取 ≥20%，标注为"待人工验证"，不执行实际人工复核
- 报告格式遵循 proposal schema：每条问题包含 `{component, file_path, layer, category, severity, description, fix_suggestion, confidence}`
