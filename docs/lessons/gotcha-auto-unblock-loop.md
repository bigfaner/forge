---
created: "2026-05-22"
tags: [architecture]
---

# Auto-unblock 循环导致无法跳过任务

## Problem

将任务 `fix.7` 设为 blocked 后，每次执行 `forge task claim` 都会被 auto-unblock 重新领取，形成无限循环：

```
$ forge task status fix.7 blocked   # 手动阻塞
$ forge task claim
Previous task '...' is blocked. Claiming new task...
Auto-unblocked task fix.7           # 自动解除阻塞！
ACTION: CLAIMED → fix.7             # 再次被领取
```

## Root Cause

因果链（3 层）：

1. **表层**：`forge task claim` 的 lazy unblock scan 遍历所有 `status == "blocked"` 的任务，deps 满足就自动解除阻塞
2. **第 1 层**：状态模型只有一个 `"blocked"` 状态，无法区分"系统阻塞"和"手动挂起"——前者应该被 auto-unblock，后者不应该
3. **第 2 层**：设计上缺少"暂停"语义——`blocked` 被同时用于系统降级（auto-downgrade）、依赖等待（block-source）和手动挂起，三种意图共用一个状态

## Solution

引入 `suspended` 状态，与 `blocked` 分离：

| 状态 | 触发者 | auto-unblock | 恢复方式 |
|------|--------|-------------|---------|
| `blocked` | 系统（auto-downgrade, block-source） | 自动解除 | deps 满足后 auto-unblock |
| `suspended` | operator（`forge task transition`） | 不触碰 | 手动 `forge task transition X pending` |

操作示例：

```bash
# 挂起任务（替代旧的 forge task status X blocked）
$ forge task transition fix.7 suspended --reason "waiting on external team"

# 恢复任务
$ forge task transition fix.7 pending --reason "external team ready"

# 或者直接决定跳过
$ forge task transition fix.7 skipped --reason "no longer needed"
```

## Reusable Pattern

**状态语义分离**：当多个不同意图共用在同一个状态上时，会出现行为冲突。解决方案是为每种语义引入独立状态，而非在现有状态上叠加布尔字段。本例中 `ManualBlock bool` 是临时 workaround，`suspended` 状态才是正确的设计。

## References

- gotcha-auto-unblock-loop.md（本文档）
- gotcha-blocked-task-never-auto-unblocks.md（反向问题——任务从未被 auto-unblock）
- feat(task): add suspended status to replace ManualBlock field（commit 3a54430f）
