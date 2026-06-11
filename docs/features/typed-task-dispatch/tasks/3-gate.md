---
id: "3.gate"
title: "Phase 3 Gate: Agent 与命令更新验证"
priority: "P0"
estimated_time: "1h"
dependencies: ["3.summary"]
status: pending
breaking: true
type: "gate"
---

# 3.gate: Phase 3 Gate — Agent 与命令更新验证

## Description

Exit verification gate for Phase 3. Confirms that task-executor, run-tasks, and execute-task are all updated and routing works correctly before proceeding to Phase 4 (Cleanup).

## Verification Checklist

- [ ] task-executor.md ≤ 50 行，只含 Hard Constraints
- [ ] run-tasks.md Step 2 调用 `task prompt <id>`，不再传 TASK_FILE/NO_TEST
- [ ] run-tasks.md 不含任何对 `forge:error-fixer` 的引用
- [ ] execute-task.md 路由与 run-tasks.md 一致
- [ ] eval-cases 类型在主会话执行（不 dispatch subagent）
- [ ] record 缺失恢复使用 `task prompt --fix-record-missed`，不使用 error-fixer
- [ ] `task prompt <id>` exit 非零时任务被标记 blocked，blockedReason 有内容

## Acceptance Criteria

- [ ] 所有上述检查项通过
- [ ] `task validate docs/features/typed-task-dispatch/tasks/index.json` 无报错
