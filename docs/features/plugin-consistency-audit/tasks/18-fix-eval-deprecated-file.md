---
id: "18"
title: "Fix: move eval deprecated freeform-injection.md"
priority: "P1"
estimated_time: "15min"
dependencies: []
type: "doc"
complexity: "low"
mainSession: false
---

# 18: Fix: move eval deprecated freeform-injection.md

## Description
eval skill 的 `rules/freeform-injection.md` 已标记为 `status: deprecated`（deprecated-by: eval-freeform-pre-revision），SKILL.md 也标注为 deprecated。但文件仍位于 rules/ 目录中，存在被意外加载的风险。应移至 `_deprecated/` 子目录或添加 `_` 前缀。

## Reference Files
- `docs/features/plugin-consistency-audit/reports/02-skills-batch-a.md#EV-01`: P1 CONFLICT, deprecated 文件仍在 rules/ (source: Report 02)
- `plugins/forge/skills/eval/rules/freeform-injection.md`: 需移动的文件 (source: audit finding)

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/eval/rules/freeform-injection.md` | 重命名为 `_deprecated/freeform-injection.md` 或添加 `_` 前缀 |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] 文件不再位于 `rules/` 根目录（已移至子目录或重命名）
- [ ] eval SKILL.md 中对该文件的 deprecated 引用路径已更新
- [ ] `rules/freeform-pipeline.md` 中对该文件的引用（如有直接路径）已更新

## Hard Rules
- 仅移动/重命名 `plugins/forge/skills/eval/rules/freeform-injection.md` 及更新相关引用
