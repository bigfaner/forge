---
id: "2"
title: "创建 forge CLI 命令参考文档"
priority: "P2"
estimated_time: "30m"
dependencies: []
type: "doc"
mainSession: false
---

# 2: 创建 forge CLI 命令参考文档

## Description

新增 `docs/conventions/forge-cli-reference.md`，记录所有有效的 forge CLI 命令及其用途。供 skill 作者编写或修改 skill 时参考，防止再次引入不存在的命令引用。

## Reference Files
- `docs/proposals/stale-cli-commands-cleanup/proposal.md` — Source proposal
- `forge-cli/internal/cmd/root.go` — CLI 命令注册（权威来源）

## Affected Files

### Create
| File | Description |
|------|-------------|
| `docs/conventions/forge-cli-reference.md` | 完整的 forge CLI 命令清单 |

### Modify
| File | Changes |
|------|---------|
| (none) | |

### Delete
| File | Reason |
|------|--------|
| (none) | |

## Acceptance Criteria

- [ ] 文档覆盖所有有效 CLI 命令（约 46 个），包括命令组、子命令和简要说明
- [ ] 每个命令标注其所在的源文件路径（便于维护）
- [ ] 文档包含 frontmatter `domains` 字段用于 consolidate-specs 自动加载
- [ ] 明确标注已移除的命令（detect、interfaces、framework、get），防止误用

## Hard Rules

- 命令清单必须与 `forge-cli/internal/cmd/root.go` 中的注册完全一致
- 不要包含 `forge version`（hidden command）
- 不要描述命令的详细用法，只记录名称、用途和源文件

## Implementation Notes

参考 `forge-cli/internal/cmd/root.go` 中的 `init()` 函数获取完整的命令注册树。每个命令组的子命令分散在各自的 `.go` 文件中。
