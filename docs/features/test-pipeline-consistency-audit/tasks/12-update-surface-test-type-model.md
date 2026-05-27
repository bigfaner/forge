---
id: "12"
title: "更新 surface-test-type-model 提案的 recipe 命名"
priority: "P2"
estimated_time: "30m"
dependencies: [5]
type: "doc"
mainSession: false
---

# 12: 更新 surface-test-type-model 提案的 recipe 命名

## Description
更新已批准的 `surface-test-type-model/proposal.md`：第 73 行和第 85 行的 recipe 命名从 `test-<surface-type>-<scope>`（如 `test-cli-functional`）改为 `<surface-key>-test`（如 `cli-test`）；第 107 行的多 surface recipe 命名同步更新；第 73 行的 alias 过渡方案不再适用（v3.0.0 大版本允许破坏性变更，alias 直接删除），该提案 NFR1 的向后兼容要求被本提案覆盖。

## Reference Files
- `proposal.md#Layer-2-Skill-文档层术语统一` — 第 14 项定义了 supersede 范围
- `proposal.md#Scope` — In Scope 第 12 项
- `proposal.md#Risks` — surface-test-type-model 被部分 supersede 的风险

## Acceptance Criteria
- [ ] 第 73 行 recipe 命名已更新为 `<surface-key>-test`
- [ ] 第 85 行 recipe 命名已更新为 `<surface-key>-test`
- [ ] 第 107 行多 surface recipe 命名已同步更新
- [ ] NFR1 向后兼容要求标记为 v3.0.0 已覆盖
- [ ] 仅修改 recipe 命名部分，测试类型映射和术语定义不变

## Hard Rules
- 仅 supersede recipe 命名部分，不修改测试类型映射和术语定义

## Implementation Notes
- 这是已批准提案的部分更新，scope 严格限定在 recipe 命名

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| `docs/proposals/surface-test-type-model/proposal.md` | 第 73、85、107 行 recipe 命名更新 + NFR1 标记 |

### Delete
| File | Reason |
|------|--------|
