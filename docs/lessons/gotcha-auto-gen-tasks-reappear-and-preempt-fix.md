---
created: "2026-05-22"
tags: [architecture, testing]
---

# 自动生成的测试任务反复出现导致 Dispatcher 领取错误任务

## Problem

Dispatcher 循环执行 Phase 3 任务时突然领取了 `T-test-gen-cases`（自动生成的 e2e 测试任务），而不是应该领取的 `fix.8`（修复预存测试失败）或 `3.4`（被 blocked 的配置任务）。用户发现 dispatcher 卡住，询问原因。

## Root Cause

因果链（5 层）：

1. **表层**：Dispatcher 领取了 T-test-gen-cases，用户觉得「卡住了」

2. **第 1 层**：`forge task claim` 的 auto-unblock 机制自动将 `blocked` 状态的 3.4 重新激活为 `pending`。同时 index 中存在 T-test-gen-cases（P1, pending），而 fix.8 也是 P1（pending）。排序规则是 P1 > P2，同优先级按 semver ID 排序

3. **第 2 层**：P1 任务中 ID 排序 `compareVersionIDs("T-test-gen-cases", "fix.8")`：
   - Segment 0: 两者都是字母段，持平
   - Segment 1: "T-test-gen-cases" 无此段（返回 -1），"8" 存在（返回 8）
   - 结果：-1 < 8，T-test-gen-cases 排序在 fix.8 之前

4. **第 3 层**：T-test-gen-cases 的 `.md` 文件在多次删除后仍然存在。原因是 `forge task index` 的 Step 7.5 在配置 `e2eTest.full: false` 时不应生成测试任务，但 PreserveRuntimeFields/cleanup 逻辑的交互导致已存在的任务未被清理

5. **第 4 层**：`auto-unblock` 机制意图是「自动激活因上游 fix 任务而 blocked 的正常任务」，但它无法区分「因等待 fix 而被 blocked」和「因拆分而被手动 blocked」两种情况。导致被手动 blocked 的 3.4 （已拆分为 3.1.1-3.1.3）被错误激活

## Solution

**短期**：清理所有自动生成测试任务的 .md 文件并设置 `e2eTest.full: false`、`validation.full: false`、`cleanCode.full: false`、`consolidateSpecs.full: false` 以彻底禁用自动生成。

**中期**（claim 排序修复）：P0/P1/P2 排序基础上，相同优先级时优先领取 fix 任务（`fix-*` ID 前缀），然后才是 feature 任务和测试任务。避免自动生成测试任务抢占 fix 任务。

**长期**（auto-unblock 区分）：blocked 任务应区分 `blockedReason`：
- `waiting_for_fix`：上游 fix 任务未完成 → 可自动激活
- `manually_blocked`（如拆分、弃用）→ 不可自动激活

## Reusable Pattern

**"自动生成任务不应抢占业务修复任务"**：任何自动生成的任务（T-test-*, T-clean-*, T-specs-*）在 Dispatcher 领取排序中应排在 fix 任务之后。当 fix 任务与自动任务同级时，fix 任务应优先。

## References

- `claim.go:242-247` — 任务领取排序逻辑（priority → semver ID）
- `claim.go:79` — `checkExistingTaskState` auto-unblock 逻辑
- `autogen.go:32-121` — `GetBreakdownTestTasks` 测试任务生成
- `build.go:164-170` — 孤儿清理逻辑
- `docs/lessons/gotcha-auto-unblock-loop.md` — auto-unblock 已知问题