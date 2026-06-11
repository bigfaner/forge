---
id: "16"
title: "Add config commands to guide.md (residual L-11)"
priority: "P2"
estimated_time: "15m"
dependencies: [10]
type: "doc"
mainSession: false
---

# 16: Add config commands to guide.md

## Description

hooks/guide.md 缺少 `forge config get` 和 `forge config set` 命令文档。这是 skill 中最高频使用的配置读取机制（20+ 处引用），但 guide.md 的 CLI 命令参考部分未列出该命令，新 skill 开发者无法了解 config 读取的标准模式。

## Reference Files
- `plugins/forge/hooks/guide.md`: Add forge config get/set documentation

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/hooks/guide.md` | Add Config Management section with `forge config get <key>` and `forge config set <key> <value>` documentation |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] guide.md 包含 Config Management 部分，列出 `forge config get <key>` 和 `forge config set <key> <value>` 命令说明

## Hard Rules
- Before modifying any plugin file, read `docs/conventions/forge-distribution.md`

## Implementation Notes
- guide.md 当前涵盖 Task Management、Feature Management、Pipeline Utilities 三类命令
- 新增 Config Management 部分应放在 Pipeline Utilities 之前或之后，保持格式一致
