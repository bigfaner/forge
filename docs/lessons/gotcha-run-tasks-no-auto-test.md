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

如果规则配置为 Claude Code Hook（settings.json），Hook 由 harness 执行，不受 Iron Laws 约束，
但需要确认：
- Hook 绑定的是 `Stop` 事件（不是 `PostToolUse`）
- Hook 的触发条件（matcher）能匹配 dispatcher 的最终输出
- Hook 命令本身执行无误

## Solution
所有任务完成后，手动运行：
```
/run-e2e-tests
```

如需真正自动化，有两个可行方案：
1. **修改 `/run-tasks` skill 本身**，在循环结束后加显式测试触发步骤（需移除对应 Iron Law）
2. **配置 `Stop` hook**，在 dispatcher 输出包含 `"All tasks completed"` 时触发测试脚本

## Key Takeaway
`/run-tasks` 是纯粹的任务编排器，其 Iron Laws 优先级高于任何外部配置的规则或 hook。
测试执行永远是用户的显式动作，不会被 dispatcher 代劳。
想要真正自动化，必须修改 skill 本身或使用正确事件的 Hook，而不是依赖对话级别的规则。
