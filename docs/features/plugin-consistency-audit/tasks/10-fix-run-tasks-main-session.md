---
id: "10"
title: "Fix: run-tasks MAIN_SESSION missing submit-task + git-commit"
priority: "P1"
estimated_time: "15min"
dependencies: []
type: "doc"
complexity: "low"
mainSession: false
---

# 10: Fix: run-tasks MAIN_SESSION missing submit-task + git-commit

## Description
`run-tasks` command 的 Step 1.5 (Main Session Routing) 在执行完 MAIN_SESSION 任务后，缺少 submit-task 和 git-commit 步骤。非 MAIN_SESSION 任务由 task-executor agent 内部调用 submit-task + git-commit，但 MAIN_SESSION 路径没有这些步骤，导致 MAIN_SESSION 任务完成后可能缺少 record 和 commit。(Source: CMD-07, Report 05)

## Reference Files
- `docs/features/plugin-consistency-audit/reports/05-commands-agent-hooks.md#CMD-07`: P1 级 INCOMPLETE (source: Report 05)
- `plugins/forge/commands/run-tasks.md`: 需修改的文件，Step 1.5 节 (source: audit finding)
- `plugins/forge/agents/task-executor.md`: 参考 task-executor 的 submit + commit 流程 (source: cross-reference)

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/commands/run-tasks.md` | 在 Step 1.5 的执行后增加 submit-task 和 git-commit 调用指令 |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] Step 1.5 包含执行后的 submit-task 指令：`Skill(skill="forge:submit-task")`
- [ ] Step 1.5 包含 git-commit 指令（条件性：如有未提交的变更）
- [ ] Step 1.5 的流程保持向后兼容：如果任务自身已处理 submit/commit，不重复执行

## Hard Rules
- 仅修改 `plugins/forge/commands/run-tasks.md`

## Implementation Notes
- 参考 task-executor agent 的 Step 4 (submit-task) 和 Step 5 (git-commit) 的调用模式
- 注意不要与 Stop hook 的 commit 逻辑冲突
