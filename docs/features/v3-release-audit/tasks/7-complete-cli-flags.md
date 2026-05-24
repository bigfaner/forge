---
id: "7"
title: "Complete CLI flags documentation (8 missing)"
priority: "P1"
estimated_time: "30m"
dependencies: ["6"]
type: "doc"
mainSession: false
---

# 7: Complete CLI flags documentation (8 missing)

## Description
CLI flags 遗漏 8 处：worktree 子命令 flags、feature -v、task query -v 等。这些遗漏导致 `--help` 输出与文档不一致。补全文档中的 flags 描述。

## Reference Files
- `proposal.md#Scope` — P1.6: CLI flags completion, 8 missing entries
- `proposal.md#Success-Criteria` — "CLI flags 补全，`--help` 与文档一致"

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `README.md` | Add missing CLI flags to command reference |
| `docs/ARCHITECTURE.md` | Update CLI section if flag descriptions exist |

## Acceptance Criteria
- [ ] 所有 `forge <subcommand> --help` 输出的 flags 在文档中有对应描述
- [ ] 无文档中存在但 `--help` 不输出的幽灵 flags

## Hard Rules
- 以 `forge --help` 和各子命令 `--help` 的实际输出为准
- 不修改 Go 源码

## Implementation Notes
需逐一运行 `forge <subcommand> --help` 记录实际 flags，与文档比对补全。
