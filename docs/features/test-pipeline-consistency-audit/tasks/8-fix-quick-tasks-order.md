---
id: "8"
title: "修正 quick-tasks 执行顺序和术语"
priority: "P1"
estimated_time: "30m"
dependencies: []
type: "doc"
mainSession: false
---

# 8: 修正 quick-tasks 执行顺序和术语

## Description
修正 `quick-tasks/SKILL.md` 中测试任务执行顺序从 `gen-journeys → gen-contracts → gen-test-scripts` 改为 `gen-journeys → run-test`（`breakdown-tasks/SKILL.md` 的顺序正确，不动）。替换 "Integration Test Impact Assessment" 为 "Test Impact Assessment"（在 `breakdown-tasks` 和 `quick-tasks` 两处）。

## Reference Files
- `proposal.md#Layer-2-Skill-文档层术语统一` — 第 10 项定义了执行顺序修正和术语替换
- `proposal.md#Success-Criteria` — 验证条件：术语替换完成

## Acceptance Criteria
- [ ] `quick-tasks/SKILL.md` 执行顺序已修正为 `gen-journeys → run-test`
- [ ] `breakdown-tasks/SKILL.md` 和 `quick-tasks/SKILL.md` 中 "Integration Test Impact Assessment" → "Test Impact Assessment"
- [ ] `breakdown-tasks/SKILL.md` 执行顺序未被修改（已经正确）

## Implementation Notes
- `breakdown-tasks/SKILL.md` 顺序正确不动，仅修改术语
- `gen-test-scripts/types/ui.md` 中的 "Integration Test" 是 UI 组件集成测试的专门概念，保留

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/quick-tasks/SKILL.md` | 执行顺序修正 + 术语替换 |
| `plugins/forge/skills/breakdown-tasks/SKILL.md` | "Integration Test Impact Assessment" → "Test Impact Assessment" |

### Delete
| File | Reason |
|------|--------|
