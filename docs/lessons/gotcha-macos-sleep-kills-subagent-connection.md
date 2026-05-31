---
created: "2026-05-31"
tags: [architecture, error-handling, local-dev-deployment]
---

# macOS 休眠杀死子代理网络连接，客户端无 read timeout 导致 9.2 小时挂起

## Problem

Task 11（reorganize internal/cmd/）的 task-executor 子代理在 26 分钟实际工作后，挂起 9.2 小时（23:34 → 次日 08:46），直到用户手动中断。

- 133 次工具调用后，子代理在 23:34:30 收到最后一次 Write tool result
- 之后 JSONL 零活动，直到用户 08:46 中断

## Root Cause

因果链（3 层）：

1. **表面现象**：Task 11 执行 9h 38m 11s，子代理 9.2 小时无任何活动
2. **直接原因**：macOS 在 **23:34:49** 进入 Idle Sleep（pmset 日志确认），距子代理最后一次活动仅 19 秒。系统休眠后网络连接中断，Claude Code 的 Node.js 进程被挂起。用户在 08:45:58 按电源键唤醒后手动中断了会话
3. **根因**：**Claude Code 的子代理 API 客户端无 read timeout**——系统休眠杀死 TCP 连接后，进程恢复时没有检测到连接已断，无限等待已关闭的 socket

**pmset 日志证据（关键时间线）：**
```
23:34:49  Sleep    Entering Sleep (Idle Sleep) — 子代理最后活动后 19 秒
23:50:17  DarkWake (maintenance)
23:50:23  Sleep    (back to sleep)
00:06:47  Sleep    (user left)
...每 15 分钟 DarkWake 一次做 maintenance，立刻睡回去...
08:45:58  Wake     (power button — 用户唤醒)
08:46:23  [Request interrupted by user] — 用户手动中断
```

## Solution

1. **Agent tool 增加 timeout 强制执行**：父会话派发子代理时必须设置 timeout（如 1800000ms = 30min），到期自动中断子代理
2. **API 客户端增加 read timeout**：单次 API 请求超过 N 秒（如 300s）无数据返回，判定连接断开并重试
3. **Claude Code 持有 caffeinate assertion**：执行任务期间阻止系统休眠（或至少在 Agent 阻塞等待子代理时阻止休眠）

## Reusable Pattern

- **长时间运行的 CLI 工具必须持有 caffeinate assertion**：macOS 默认几分钟无用户输入就会休眠，后台进程的网络连接会被杀死。Claude Code 在 Agent 阻塞等待子代理时应阻止休眠
- **网络客户端必须有 read timeout**：即使 TCP keepalive 开启（macOS 的 TCPKeepAlive=active），系统休眠期间 keepalive 探测无法发出，对端关闭连接后客户端无法感知
- **根因诊断：先查 pmset 再猜 API**：当出现长时间无活动时，第一步应查 `pmset -g log` 排查系统休眠，而非猜测 API 层问题

## Related Files

- `docs/forensics/task-11-9h-stuck/report.md` — 完整取证报告
- `docs/forensics/task-11-9h-stuck/evidence/subagent-11/evidence.json` — 子代理证据
- `docs/lessons/gotcha-large-output-stall-subagent.md` — 类似案例：大文件输出导致 2h 卡死（输出端问题）
- `docs/lessons/gotcha-task-executor-thinking-overhead.md` — 类似案例：87% 时间消耗在 thinking
- `docs/lessons/gotcha-task-executor-never-returns.md` — 类似案例：executor 不返回主会话
