---
id: "5"
title: "更新活跃文档：移除 gen-and-run 引用"
priority: "P2"
estimated_time: "30min"
dependencies: ["3"]
type: "doc"
mainSession: false
---

# 5: 更新活跃文档

## Description

更新 OVERVIEW.md 和 task-lifecycle 文档中的 gen-and-run 引用，反映当前 staged pipeline 架构。

## Reference Files
- `proposal.md#P2-—-gen-and-run-废弃代码移除` — 活跃文档更新列表
- `proposal.md#Success-Criteria` — grep 零结果验证标准

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `forge-cli/docs/OVERVIEW.md` | 移除 gen-and-run 引用，更新为 staged pipeline 描述 |
| `docs/business-rules/task-lifecycle.md` | 移除 gen-and-run 引用 |

## Acceptance Criteria

- [ ] `grep -r "gen-and-run" forge-cli/docs/ docs/business-rules/` 返回零结果
- [ ] OVERVIEW.md 中 test pipeline 描述反映当前 staged 架构（gen-journeys → gen-contracts → gen-scripts → run → verify-regression）

## Hard Rules

- 仅更新 gen-and-run 相关引用，不重构文档其他部分

## Implementation Notes

OVERVIEW.md 中 `test-pipeline.gen-and-run` 条目需替换或移除。task-lifecycle.md 中的系统类型列表需移除 `test.gen-and-run` 条目。
