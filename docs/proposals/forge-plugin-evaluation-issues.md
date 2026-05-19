---
title: "Forge Plugin Evaluation Issues"
domains: [forge, plugin, evaluation, backlog]
---

# Forge Plugin Evaluation Issues

Evaluation date: 2026-05-19. Two rounds: plugin consistency + task scheduling mechanism.

## Task Scheduling — P0

### TS-P0-1: `claim` 和 `status` 缺少文件锁

`claim.go:103` 和 `status.go:105` 直接调用 `task.SaveIndex()` 无锁无原子写入。`submit.go:89-98` 正确使用了 `LockFile` + `SaveIndexAtomic`。并发 claim 可能导致 index.json 损坏。

### TS-P0-2: `execute-task` 缺少 `forge feature set` 步骤

`run-tasks.md` Step 0 处理了 feature 设置，`execute-task.md` 直接调用 `forge task claim`，未设置 feature 时会失败。

### TS-P0-3: `task-executor` 输出格式在 blocked 时假设有 commit hash

`task-executor.md:53` 输出格式包含 `<commit-hash>` 字段，但 blocked 任务跳过 git-commit，字段为空。

## Task Scheduling — P1

### TS-P1-1: `run-tasks` 缺少 e2e gate

`execute-task` 有 Step 3b Feature E2E Gate（检查 test-e2e recipe 和 e2e spec 文件），`run-tasks` 无等价机制。通过 `/run-tasks` 调度的任务跳过了 per-task e2e 验证。

### TS-P1-2: `execute-task` Step 1.5.4 fix-task 创建缺少 `--block-source`

line 39 仅说 "spawn fix task (same as Step 2 verify logic)"，未显式提供命令。LLM 可能遗漏 `--block-source` flag。

### TS-P1-3: CLI 依赖解析逻辑重复

`claim.go` 的 `checkDependenciesMet` 和 `status.go` 的 `checkUnmetDeps` 实现几乎相同。前者额外检查 active fix-task，后者不检查——语义差异脆弱且未文档化。

### TS-P1-4: `isBusinessTask` 三处重复定义

`validate_index.go:253`、`prompt/prompt.go:204`、`task/add.go:397` 各有一份。

### TS-P1-5: 索引保存模式不一致

三种保存机制混用：`task.SaveIndex()`（非原子无锁）用于 claim/status，`SaveIndexAtomic()`（原子无锁）用于 submit，`BuildIndex()` 用于 add/index。claim 和 status 至少应使用原子保存。

### TS-P1-6: `quality_gate.go` 的 `addFixTask` 复制了 `add.go` 的 `executeAdd` 逻辑

注释承认 "mirrors executeAdd() in add.go"，共享逻辑应提取为内部函数。

## Task Scheduling — P2 设计缺口

### TS-P2-1: 无单次任务重试计数

fix-task 链有 `--source-task-id` 自动追溯到 root，但 dispatcher 不跟踪重试次数，无深度上限。

### TS-P2-2: 无过期 `in_progress` 任务恢复机制

进程崩溃时 task 卡在 `in_progress`，无文档化的恢复流程。`forge task claim` 有 `ACTION: CONTINUE` 但 dispatcher 未解释 resume 机制。

### TS-P2-3: 状态机转换未在 plugin 端文档化

有效转换由 CLI 代码强制，但 plugin 端没有枚举完整转换图的文件。

### TS-P2-4: 无 dispatcher 级别审计日志

调度决策不持久化，只能从 subagent 对话记录中重建。

## Plugin Consistency — 已修复

以下问题已于 2026-05-19 修复：

- 类型系统统一为 `coding.*` 前缀（breakdown-tasks, quick-tasks）
- fix-bug 缺少 `forge:` 命名空间
- run-e2e-tests 测试路径对齐到 `tests/<journey>/`
- gen-contracts typo `/gen-jneys` → `/gen-journeys`
- extract-design-md frontmatter 字段名修正
- extract-design-md 重构为 skill（模板提取到 templates/）
- gen-sitemap 重构为 skill（schema 提取到 templates/）
- clean-code command description 对齐 skill
- guide.md 精简 autoConfig 部分
- 新增 `docs/conventions/skill-self-containment.md`

## Pipeline 文档补充（待办）

完整 pipeline 路径（brainstorm → write-prd → tech-design → breakdown-tasks → run-tasks → submit-task）缺乏单文件文档。计划后续补充到 README.md。
