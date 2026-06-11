---
id: "2"
title: "Harden test-run auto-gen template with strict AC"
priority: "P1"
estimated_time: "30min"
dependencies: []
type: "doc"
mainSession: false
---

# 2: Harden test-run auto-gen template with strict AC

## Description

test-run 任务的 auto-gen 模板（`pkg/task/templates/test-run.md`）缺少严格的验收标准，导致 agent 可能生成假测试或未经验证的测试。需要在模板的 AcceptanceCriteria 部分添加硬性 AC，确保生成的测试任务要求真实通过的测试。

## Reference Files
- `docs/proposals/test-pipeline-interleaved/proposal.md` — Proposed Solution, Scope > In Scope, Success Criteria
- `forge-cli/pkg/task/templates/test-run.md`: add hardened AC to auto-gen template (ref: Scope > In Scope)

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| `forge-cli/pkg/task/templates/test-run.md` | 在 AcceptanceCriteria 部分添加硬性 AC 条目 |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria

- [ ] 模板 AcceptanceCriteria 部分包含 AC：所有测试用例必须通过（不能有 skip 或预期失败）
- [ ] 模板 AcceptanceCriteria 部分包含 AC：必须是真实测试（验证实际功能行为），不能是占位符或 always-pass 假测试

## Implementation Notes

当前模板内容较短（29行），AcceptanceCriteria 字段为空。需要在 `context` 部分的 `AcceptanceCriteria` 后追加具体的 AC 条目。AC 应使用模板变量保持与现有模板风格一致。
