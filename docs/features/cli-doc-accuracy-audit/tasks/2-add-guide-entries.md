---
id: "2"
title: "Add new commands and flags to guide.md"
priority: "P1"
estimated_time: "30m"
dependencies: [1]
type: "doc"
mainSession: false
---

# 2: Add new commands and flags to guide.md

## Description
在 guide.md CLI 参考部分补充 4 处缺失的命令和标志（G5-G8）。这些命令已在 CLI 中实现但 guide.md 未收录，导致 agent 会话中不知道这些能力。

## Reference Files
- `docs/proposals/cli-doc-accuracy-audit/proposal.md` — In Scope (G5-G8), Success Criteria
- `plugins/forge/hooks/guide.md`: CLI 参考部分新增命令条目 (ref: In Scope)

## Affected Files

### Create
| File | Description |
|------|-------------|
| _(none)_ |

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/hooks/guide.md` | G5: 新增 task query; G6: 新增 task check-deps; G7: 新增 feature list; G8: task list 补充 --tree |

### Delete
| File | Reason |
|------|--------|
| _(none)_ |

## Acceptance Criteria
- [ ] guide.md 新增 `forge task query <id-or-key>` 命令描述，与 `forge task query --help` 输出一致
- [ ] guide.md 新增 `forge task check-deps` 命令描述，与 `forge task check-deps --help` 输出一致
- [ ] guide.md 新增 `forge feature list` 命令描述，与 `forge feature list --help` 输出一致
- [ ] guide.md 中 `forge task list` 描述补充 `--tree` 标志

## Implementation Notes
- 先运行 `forge task query --help`、`forge task check-deps --help`、`forge feature list --help`、`forge task list --help` 获取准确的命令行为
- 新增条目的格式应与 guide.md 中现有命令条目格式保持一致
