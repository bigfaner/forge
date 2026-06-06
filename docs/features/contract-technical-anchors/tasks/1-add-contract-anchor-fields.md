---
id: "1"
title: "Contract 模板增加技术锚点字段"
priority: "P0"
estimated_time: "1.5h"
dependencies: []
type: "doc"
mainSession: false
---

# 1: Contract 模板增加技术锚点字段

## Description
为 Contract 模板 frontmatter 增加技术锚点字段，建立 Contract 与技术实现的桥梁。Forge 测试管道生成的 Contract 规格缺少技术锚点（API endpoint、CLI command 等），导致 gen-test-scripts 只能依赖 LLM 推断技术细节。本任务在 Contract 模板中定义各 surface 类型的锚点字段，作为后续锚点填充和交叉验证的基础。

## Reference Files
- `docs/proposals/contract-technical-anchors/proposal.md` — Proposed Solution, Anchor Field Schema
- `plugins/forge/skills/gen-contracts/templates/contract.md`: 增加 frontmatter 锚点字段 (ref: Anchor Field Schema)

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/gen-contracts/templates/contract.md` | 增加 API/CLI/TUI/Web/Mobile 锚点字段及 `last_anchor_sync` |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] Contract 模板 frontmatter 包含 API 锚点字段：`endpoint` (string)、`method` (string)，及可选字段 `content_type`、`auth_required`
- [ ] Contract 模板 frontmatter 包含 CLI/TUI 锚点字段：`command` (string)，及可选字段 `subcommand`、`flags`、`aliases`
- [ ] Contract 模板 frontmatter 包含 Web 锚点字段：`page` (string)，及可选字段 `route`、`requires_auth`、`layout`
- [ ] Contract 模板 frontmatter 包含 Mobile 锚点字段：`screen` (string)，及可选字段 `navigation_path`、`deeplink`、`platform`
- [ ] `last_anchor_sync` 时间戳字段包含在 frontmatter 中

## Implementation Notes
- 锚点字段定义参照 proposal 的 Anchor Field Schema 表格
- 字段应为可选（handbook 不存在时 Contract 仍可正常生成）
- 现有 Contract 不受影响——缺少锚点字段时管道降级为现有行为
