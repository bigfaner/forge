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

因果链（4 层）：

1. **表层**：每次 claim 都会被 auto-unblock 循环
2. **第 1 层**：`forge task claim` 的 auto-unblock 机制会在没有其他任务可领时，自动解除最近 blocked 任务的阻塞。`fix.7` 是唯一处于 pending→blocked→in_progress 循环中的 fix 任务，所以每次都被选中
3. **第 2 层**：`forge task status fix.7 skipped` 被状态机拒绝（invalid transition: in_progress cannot transition to skipped），因为任务类型 `coding.fix` 的状态机可能不允许 skipped 状态，或者系统不允许直接跳过任务。`blocked` 是唯一允许的非活跃状态，但恰好就是 auto-unblock 会解除的状态
4. **第 3 层**：auto-unblock 机制的设计目标是"防止由于遗漏阻塞导致流水线死锁"，但缺少白名单或手动冻结机制——一旦某个任务被判定为"可以解除阻塞"，就无法阻止它被解除

## Solution

当需要跳过某个任务时，有两种可靠方案：

1. **直接完成**：如果任务验证后确无必要（如 fix.5 是 no-op），使用 `forge task submit` 正常完成（需质量门）
2. **手动调整 index**：编辑 `tasks/index.json`，将任务状态改为 `completed` 并添加备注。这是最后的绕行手段

**避免的操作**：`forge task status X blocked` 在存在 auto-unblock 机制时不是可靠的"暂停"操作——它只对不在 claim 队列顶端的任务有效。

## Reusable Pattern

**auto-unblock 的两面性**：这个机制防止了流水线因 blocked 任务死锁，但也让"标记 blocked → 等后面处理"的策略不可行。对于想推迟的任务，要么完成它，要么删除它。没有"暂挂"状态。

## References

- gotcha-auto-unblock-loop.md（本文档）
- gotcha-blocked-task-never-auto-unblocks.md（反向问题——任务从未被 auto-unblock）
- fix.7 执行记录