---
created: "2026-05-22"
tags: [architecture, testing]
---

# 预存语法错误导致任务执行器卡住

## Problem

Task 3.1 (RunE migration) 的子代理在派发后卡住，用户手动中断后报告 "卡住了"。子代理未能产生可见进展。

## Root Cause

因果链（5 层）：

1. **表层**：Task 3.1 的子代理派发后长时间无输出，用户中断
2. **第 1 层**：`go build ./...` 在 submit.go 和 test_promote.go 中存在真实语法错误（`unexpected ) at end of statement`），导致整个代码库不可编译。子代理的每一步（读文件、改代码、编译验证）都被阻塞

3. **第 2 层**：Task 3.1 是全局 Run→RunE 迁移，预估 3h，横跨 internal/cmd/ 下所有命令文件。子代理需要大量时间读取文件、理解范围、逐个迁移——但代码不编译时，任何验证步骤都会失败

4. **第 3 层**：Task 2.4 (submit-refactor) 的 breaking 质量门声称编译通过，但实际在 submit.go 中遗留了语法错误。后续的 fix.1 修复了部分编译错误但漏掉了 submit.go 和 test_promote.go 的语法错误

5. **第 4 层**：Dispatcher 协议在派发 coding.* 任务前不做预编译检查。仅依赖 `forge task submit` 时的质量门，这意味着一旦某次质量门漏检，后续所有任务都在损坏的代码库上执行

## Solution

**调度层修复**：在派发 coding.* 任务前，执行预编译检查：

1. `go build ./...` 失败 → 自动创建 fix 任务修复编译错误，当前任务延期
2. `go build ./...` 成功 → 正常派发

**任务粒度修复**：预估 >2h 的大型重构任务（全局 Run→RunE 迁移）应拆分为更小的子任务（每个文件或每组文件一个），避免子代理在单个调用中超时或卡住。

## Reusable Pattern

**"不可编译的代码库是调度器的死锁陷阱"**：任何时候代码库处于不可编译状态，都不应派发新的 coding.* 任务。必须有一个前置的编译健康门控，优先排队编译修复任务，再派发业务任务。

## References

- `submit.go:128` — `unexpected ) at end of statement` 语法错误
- `test_promote.go:46,67` — 更多语法错误
- `docs/lessons/gotcha-dispatcher-ignores-compilation-diagnostics.md` — 相关的质量门漏检问题