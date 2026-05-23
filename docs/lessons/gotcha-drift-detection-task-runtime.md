---
created: "2026-05-23"
tags: [architecture, testing]
---

# doc.drift 任务的 task-executor 执行时间过长

## Problem

`doc.drift` 类型的 drift detection 任务由 task-executor 子代理执行时，耗时可达 4+ 小时，严重阻塞任务调度循环。

## Root Cause

1. **表层**: drift detection 任务执行时间远超其他任务类型，4+ 小时仍未完成
2. **中层**: `doc.drift` 任务的模板内容极其空洞（仅一句 "Execute this test pipeline task"），task-executor 无法从任务文件获得任何搜索范围指引，只能做全代码库扫描
3. **根因**: `forge task index` 生成的 drift detection 任务是通用模板，没有注入 feature 上下文（哪些 spec 文件存在、哪些代码路径受影响）。task-executor 对所有任务类型使用相同的执行流程（探索→执行→验证），对于空洞的 doc.drift 任务会触发无方向的全面探索

## Solution

1. **跳过自动执行**: 对于纯代码重构类 feature（无 spec 变更），drift detection 可以直接 skip，因为没有 spec 漂移的可能
2. **在模板中注入上下文**: `forge task index` 生成 drift detection 任务时，扫描 feature 实际涉及的 spec 文件和代码路径，写入任务内容
3. **手动执行**: drift detection 不是关键路径任务，可放在所有 coding 任务完成后由用户手动触发

## Reusable Pattern

当 auto-generated 的任务模板内容过于空洞（无具体文件路径、无搜索范围、无上下文）时：
- 不应交给 task-executor 自动执行，空洞指令会导致无方向的全面探索
- 应在模板生成阶段注入足够的上下文，或标记为手动执行
- 对于 `doc.*` 类型任务，评估 feature 性质：纯代码重构无 spec 变更 → 直接 skip
