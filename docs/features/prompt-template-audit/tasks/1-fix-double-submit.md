---
id: "1"
title: "P0: 移除6个模板的显式submit步骤（双重提交）"
priority: "P0"
estimated_time: "30m"
dependencies: []
type: "doc"
mainSession: false
---

# 1: P0: 移除6个模板的显式submit步骤（双重提交）

## Description
task-executor agent 定义的第 8 步会自动调用 submit-task，但 6 个模板在步骤末尾也显式要求调用 `Skill(skill="forge:submit-task")`，导致双重提交。需要从这 6 个模板中移除显式 submit 步骤，统一依赖 task-executor 的自动提交。

## Reference Files
- `docs/proposals/prompt-template-audit/proposal.md` — Source proposal (Section 1.1)

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `forge-cli/pkg/prompt/data/doc.md` | 移除 Step 4 submit 步骤，将步骤从 4 步缩减为 3 步 |
| `forge-cli/pkg/prompt/data/doc-eval.md` | 移除 Step 3 submit 步骤，将步骤从 3 步缩减为 2 步 |
| `forge-cli/pkg/prompt/data/doc-summary.md` | 移除 Step 3 submit 步骤，将步骤从 3 步缩减为 2 步 |
| `forge-cli/pkg/prompt/data/doc-consolidate.md` | 移除 Step 3 submit 步骤，将步骤从 3 步缩减为 2 步 |
| `forge-cli/pkg/prompt/data/doc-drift.md` | 移除 Step 3 submit 步骤，将步骤从 3 步缩减为 2 步 |
| `forge-cli/pkg/prompt/data/clean-code.md` | 移除 Step 3 submit 步骤，将步骤从 3 步缩减为 2 步 |

## Acceptance Criteria
- [ ] 6 个模板中不再包含 `Skill(skill="forge:submit-task")` 调用
- [ ] 每个模板的步骤编号连续、无跳跃
- [ ] 移除 submit 步骤后，模板剩余的步骤语义完整

## Implementation Notes
- 优先验证 submit-task 的幂等性（方案 A 的前提）
- 移除 submit 步骤后需重新编号后续步骤
- doc.md 原为 4 步（Read → Execute → Self-Check → Submit），移除 Submit 后变为 3 步
- doc-eval/doc-summary/doc-consolidate/doc-drift/clean-code 原为 3 步，移除 Submit 后变为 2 步
