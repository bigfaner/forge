---
id: "3"
title: "breakdown-tasks/ui-placement: Add surface guard before route validation"
priority: "P1"
estimated_time: "0.5h"
dependencies: [1]
type: "doc"
mainSession: false
complexity: "low"
---

# 3: breakdown-tasks/ui-placement: Add surface guard before route validation

## Description

breakdown-tasks 的 ui-placement.md Placement Validation 段落（lines 26-28）通过检查 sitemap.json 存在性来决定是否验证 route，但未检查 surface 类型。TUI/Mobile 项目可能存在遗留的 sitemap.json 导致误校验。需在 route 验证前增加 web surface 检查。

## Reference Files
- `plugins/forge/skills/breakdown-tasks/rules/ui-placement.md`: Placement Validation 段落 step 2 "Check if `docs/sitemap/sitemap.json` exists" 需增加 surface 类型前置检查 (source: proposal.md#Scope-In-Scope)

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/breakdown-tasks/rules/ui-placement.md` | Placement Validation 段落增加 web surface 检查 |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] Placement Validation 段落在检查 sitemap.json 存在性之前，先检查项目是否有 web surface
- [ ] 无 web surface 时，跳过 route 验证并输出适当提示（而非警告 sitemap 缺失）

## Implementation Notes

当前 step 2 的 WARN 消息 "sitemap.json not found — cannot verify existing-page routes" 对非 web 项目有误导性。非 web 项目不应暗示用户需要运行 /gen-web-sitemap，而应直接跳过。
