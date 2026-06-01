---
id: "2"
title: "write-prd: Conditionalize sitemap references by surface type"
priority: "P1"
estimated_time: "1.5h"
dependencies: [1]
type: "doc"
mainSession: false
complexity: "medium"
---

# 2: write-prd: Conditionalize sitemap references by surface type

## Description

write-prd skill 有 3 个位置无条件引用 sitemap.json：SKILL.md Step 1（line 226）、ui-functions.md（line 7）、self-check.md（lines 19-20）。需在每个读取点前增加 web surface 检查，非 web 项目跳过 sitemap 相关操作。

## Reference Files
- `plugins/forge/skills/write-prd/SKILL.md`: Step 1 line 226 "Read `docs/sitemap/sitemap.json` if it exists" 需改为 surface 条件化读取 (source: proposal.md#Scope-In-Scope)
- `plugins/forge/skills/write-prd/rules/ui-functions.md`: Placement Rules 第 1 条 "Read `docs/sitemap/sitemap.json`" 需增加 surface 前置条件 (source: proposal.md#Scope-In-Scope)
- `plugins/forge/skills/write-prd/rules/self-check.md`: new-feature intent 表格的 Placement consistency 和 Sitemap availability 两行需增加 surface 前置条件 (source: proposal.md#Scope-In-Scope)

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/write-prd/SKILL.md` | Step 1 中 sitemap.json 读取增加 web surface 检查 |
| `plugins/forge/skills/write-prd/rules/ui-functions.md` | Placement Rules 第 1 条增加 web surface 前置条件 |
| `plugins/forge/skills/write-prd/rules/self-check.md` | Placement consistency 和 Sitemap availability 检查增加 web surface 前置条件 |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] SKILL.md Step 1 读取 sitemap.json 前，先检查项目是否有 web surface（通过 `forge surfaces --json`），无则跳过读取
- [ ] ui-functions.md Placement Rules 读取 sitemap 前，增加 web surface 前置条件检查
- [ ] self-check.md 的 Placement consistency 检查增加 web surface 前置条件；Sitemap availability 检查同样受 surface 条件守卫

## Implementation Notes

3 个文件的守卫逻辑保持一致：先 `forge surfaces --json` 检测 web surface，有则继续原有逻辑，无则跳过。注意 write-prd 可能在同一会话中被多种项目类型调用，守卫应在读取时而非 skill 入口处执行，确保 monorepo（含 web surface）正常工作。
