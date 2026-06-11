---
id: "12"
title: "Reduce overlapping logic in execute-task and run-tasks commands"
priority: "P1"
estimated_time: "1h"
dependencies: []
type: "doc"
complexity: "medium"
mainSession: false
---

# 12: Reduce overlapping logic in execute-task and run-tasks commands

## Description
execute-task 和 run-tasks commands 共享约 20-30 行接口契约（claim 格式 + fix-type 表），核心逻辑各自独立。需精简重叠部分，使每个 command 各自自洽。

## Reference Files
- `docs/proposals/skill-command-independence-audit/proposal.md` — Scope > In Scope, Key Risks, Success Criteria
- plugins/forge/commands/execute-task.md: 精简与 run-tasks 重叠的接口契约 (ref: Scope)
- plugins/forge/commands/run-tasks.md: 精简与 execute-task 重叠的接口契约 (ref: Scope)

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| plugins/forge/commands/execute-task.md | 精简与 run-tasks 重叠的 claim 格式和 fix-type 表 |
| plugins/forge/commands/run-tasks.md | 精简与 execute-task 重叠的 claim 格式和 fix-type 表 |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] claim 格式描述在两个 command 中各自独立且简洁
- [ ] fix-type 表在两个 command 中各自独立且简洁
- [ ] 两个 command 各自完整可执行，无悬挂引用

## Hard Rules
- 仅修改以下文件：plugins/forge/commands/execute-task.md、plugins/forge/commands/run-tasks.md

## Implementation Notes
两个 command 的核心逻辑不同（execute-task 是单任务执行，run-tasks 是自动调度），只是接口契约有重叠。
