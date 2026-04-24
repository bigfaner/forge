# /run-tasks 不会自动触发测试

## Problem
配置了"所有任务完成后自动执行测试"的规则，但 `/run-tasks` 结束后测试没有触发。

## Root Cause
`/run-tasks` dispatcher 有硬编码的 Iron Law：**"NO running tests directly"**。
Skill 提示词明确写道：
> Do NOT run e2e tests automatically — the dispatcher must not execute tests. Only suggest.

Dispatcher 的职责只有三个动作：claim → dispatch → verify record。
任何测试执行都被设计为用户手动触发，不受外部规则影响。

两套规则系统的优先级关系：
```
Skill Iron Laws > CLAUDE.md > 用户对话指令
```

## Solution（已实现）

使用 `task all-completed` + Claude Code `Stop` hook 实现真正的自动化：

### 1. `task all-completed` 命令
检查当前 feature 所有任务是否均为 `completed` 或 `skipped`：
- 若未全部完成 → 静默退出 exit 1（无输出）
- 若全部完成 → 依次运行：
  1. Feature e2e 测试（`docs/features/{slug}/testing/scripts/`，若存在）
  2. 项目级测试（自动检测：`go.mod` → `package.json` → `Makefile` → `pytest`，或 `index.json` 中的 `testCommand` 字段）

### 2. Stop hook 配置（`.claude/settings.local.json`）
```json
"hooks": {
  "Stop": [{
    "hooks": [{"type": "command", "command": "task all-completed"}]
  }]
}
```

Hook 由 harness 执行，不受 Iron Laws 约束。每次 Claude stop 时触发，
`task all-completed` 内部判断是否需要运行测试，未完成时静默退出，无性能损耗。

### 3. 可选：在 `index.json` 中指定测试命令
```json
{
  "testCommand": "make test"
}
```

## Key Takeaway
`/run-tasks` 是纯粹的任务编排器，其 Iron Laws 优先级高于任何外部配置的规则。
真正自动化的正确方式：**修改 task-cli 命令本身 + Stop hook**，而不是依赖对话级别的规则。
Hook 执行在 Iron Laws 之外，是唯一可靠的自动化入口。
