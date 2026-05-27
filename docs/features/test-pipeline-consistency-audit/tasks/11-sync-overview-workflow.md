---
id: "11"
title: "同步更新 OVERVIEW 和 WORKFLOW 文档（含中文版）"
priority: "P2"
estimated_time: "1h"
dependencies: [10]
type: "doc"
mainSession: false
---

# 11: 同步更新 OVERVIEW 和 WORKFLOW 文档（含中文版）

## Description
替换 `OVERVIEW.md`、`WORKFLOW.md` 及其中文版（`OVERVIEW.zh.md`、`WORKFLOW.zh.md`）中所有 "e2e" 泛用（约 15+10 处）为正确的 surface-specific 术语，更新 graduation/staging 描述为 tag-based promotion，移除 "profile" 旧术语引用。

## Reference Files
- `proposal.md#Layer-2-Skill-文档层术语统一` — 第 13 项定义了文档替换范围
- `proposal.md#Success-Criteria` — 验证条件：中文版文档 grep 返回 0

## Acceptance Criteria
- [ ] `OVERVIEW.md` 和 `OVERVIEW.zh.md` 中 "e2e" 泛用替换为 surface-specific 术语
- [ ] `WORKFLOW.md` 和 `WORKFLOW.zh.md` 中 "e2e" 泛用替换为 surface-specific 术语
- [ ] graduation/staging 描述已更新为 tag-based promotion
- [ ] "profile" 旧术语引用已移除
- [ ] `grep -rn "tests/e2e" forge-cli/docs/OVERVIEW.zh.md forge-cli/docs/WORKFLOW.zh.md` 返回 0 结果

## Implementation Notes
- 中文版术语需保持与英文版一致（surface、Convention 等专有名词可保留英文）
- 约 25 处修改，批量替换为主

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| `docs/OVERVIEW.md` | e2e → surface-specific 术语 |
| `docs/WORKFLOW.md` | e2e → surface-specific 术语 |
| `docs/OVERVIEW.zh.md` | 中文版同步更新 |
| `docs/WORKFLOW.zh.md` | 中文版同步更新 |

### Delete
| File | Reason |
|------|--------|
