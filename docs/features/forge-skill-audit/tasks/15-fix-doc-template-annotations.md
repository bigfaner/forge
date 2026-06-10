---
id: "15"
title: "Add doc template annotations (MINOR-C2, MINOR-D3)"
priority: "P2"
estimated_time: "15m"
dependencies: [10]
type: "doc"
mainSession: false
---

# 15: Add doc template annotations

## Description

breakdown-tasks/templates/task-doc.md 缺少 surface-key/surface-type 不存在的注释说明（quick-tasks 版本有此注释）。同时两个 task-doc.md 都没有 `complexity` 字段但 SKILL.md 的 Complexity 判定未排除 doc 任务，需明确 doc 任务是否需要 complexity。

## Reference Files
- `plugins/forge/skills/breakdown-tasks/templates/task-doc.md`: Add surface-key annotation
- `plugins/forge/skills/quick-tasks/templates/task-doc.md`: Reference for annotation format

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/breakdown-tasks/templates/task-doc.md` | Add surface-key/surface-type absence annotation matching quick-tasks format |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] breakdown-tasks/templates/task-doc.md 包含与 quick-tasks 版本一致的 surface-key/surface-type 缺失注释
- [ ] 两个 task-doc.md 对 complexity 字段的处理方式一致（要么都省略并注释说明，要么都包含）

## Hard Rules
- Before modifying any plugin file, read `docs/conventions/forge-distribution.md`

## Implementation Notes
- quick-tasks task-doc.md 第 9-11 行有参考注释格式
- complexity 字段问题：如果 doc 任务不需要 complexity，应在 task-doc.md frontmatter 注释中说明原因；如果需要，应添加 `complexity: "{{COMPLEXITY}}"` 字段
