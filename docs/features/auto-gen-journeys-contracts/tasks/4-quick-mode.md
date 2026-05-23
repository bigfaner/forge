---
id: "4"
title: "autogen.go Quick 模式：替换 gen-and-run 为 staged across types 拓扑"
priority: "P0"
estimated_time: "2h"
dependencies: ["1", "2"]
scope: "backend"
breaking: true
type: "coding.feature"
mainSession: false
---

# 4: autogen.go Quick 模式：替换 gen-and-run 为 staged across types 拓扑

## Description

修改 `GetQuickTestTasks()` 将 gen-and-run 替换为 gen-journeys + gen-contracts + gen-scripts 的拆分任务，并重写 `resolveQuickDeps()` 为 staged across types 拓扑。

## Reference Files
- `docs/proposals/auto-gen-journeys-contracts/proposal.md` — Source proposal
- `forge-cli/pkg/task/autogen.go` — GetQuickTestTasks (L176-237), resolveQuickDeps (L470-502)
- `forge-cli/pkg/task/types.go` — TypeTestGenJourneys, TypeTestGenContracts (Task 1 新增)

## Acceptance Criteria

- [ ] `GetQuickTestTasks()` 不再生成 gen-and-run 任务（TypeTestGenAndRun 不再出现在 Quick 模式输出中）
- [ ] 为每个 interface type 生成 `T-test-gen-journeys-{type}` 任务
- [ ] 生成一个 `T-test-gen-contracts` 任务
- [ ] 为每个 interface type 生成 `T-test-gen-scripts-{type}` 任务
- [ ] 生成 T-test-run 和 T-test-verify-regression（或 Quick 模式前缀变体）
- [ ] gen-journeys 模板 body 包含 `AUTO_COMMIT=true` 条件指令
- [ ] gen-contracts 模板 body 包含 `SKIP_EVAL_GATE=true` 条件指令
- [ ] `resolveQuickDeps()` 实现 staged across types 拓扑：
  - Stage 1: 所有 gen-journeys 并行（无依赖）
  - Stage 2: gen-contracts 依赖所有 gen-journeys
  - Stage 3: 所有 gen-scripts 依赖 gen-contracts
  - Stage 4: run 依赖所有 gen-scripts
  - Stage 5: verify-regression 依赖 run
- [ ] 所有依赖查找使用 `findTaskIndex` / `findTaskIndexByPrefix`（无算术索引）
- [ ] findTaskIndex 返回 -1 时 panic 并输出明确错误信息
- [ ] 所有现有 Quick 模式单测通过

## Hard Rules

- Quick 模式任务 ID 前缀保持与现有约定一致（检查 Quick 模式是用 T- 还是 T-quick- 前缀）
- 不删除 TypeTestGenAndRun 类型定义（保留用于历史 index.json 兼容）
- breaking change：修改了 GetQuickTestTasks 的返回值结构

## Implementation Notes

- proposal.md 的 staged across types 策略意味着：gen-journeys 各 type 间无依赖，可并行执行。gen-contracts 需要所有 Journey 完成后才能执行代码侦察
- 可以考虑提取 Breakdown/Quick 共享的依赖链逻辑（参数化 includeEval），但这是可选优化
- gen-journeys 任务的 BodyContext 中 Mode="quick"，模板据此选择 proposal.md 输入路径
- 现有 resolveQuickDeps 中 T-quick-validate-code、T-quick-doc-drift、T-clean-code 的依赖应保持不变（依赖 verify-regression）
