---
created: "2026-05-26"
tags: [testing, architecture]
---

# Breaking Task Must Scope Integration Test Fixtures

## Problem
Task 4（split run-tests per-surface-key）标记为 `breaking: true`，但只覆盖了 `forge-cli/pkg/task/` 下的单元测试。`tests/` 目录下的 8 个集成测试套件（60+ 用例）因 config fixture 缺少 `surfaces` 字段全部失败，由 fix-task 逐个修复，耗时远超预期。

## Root Cause
1. `breaking: true` 标记仅用于质量门路由（触发额外检查），未强制要求评估集成测试影响
2. 任务描述中未提及 `tests/` 目录的 config fixtures 需要同步更新
3. `BuildIndex` 是核心函数，被 `tests/` 下多个集成测试套件间接调用，但任务 scope 只写了直接影响的 `autogen.go`
4. task-executor subagent 逐个发现并修复，缺乏全局视角识别"60+ 失败都是同一个根因"的模式

## Solution
Breaking task 的任务描述应包含：
1. **受影响调用方清单**：不仅列直接调用的文件，还包括集成测试中间接调用的路径
2. **fixture 更新指引**：明确哪些测试 fixture 需要同步修改（如 config.yaml 模板需新增 `surfaces` 字段）
3. **批量修复策略**：如果集成测试失败模式相同（同一根因），应在任务 AC 中标注"所有集成测试 config fixture 统一添加 surfaces 配置"

## Reusable Pattern
当任务的 Hard Rules 中有 `breaking: true` 或函数签名变更时：
- grep 所有测试目录（`tests/` + `forge-cli/`）中调用变更函数的位置
- 检查集成测试的 fixture/config 是否依赖变更的隐式契约（如 config 结构）
- 在 AC 中显式列出 fixture 更新范围，避免 fix-task 发现式修复

## Example
```
# Task AC 应包含
- [ ] tests/ 下所有使用 forgeconfig.ReadConfig 的测试 fixture 添加 surfaces: api
- [ ] grep -rl "BuildIndex\|ReadConfig" tests/ 确认无遗漏
```

## References
- Task 4: `4-split-run-tests-per-surface-key.md`（breaking: true，scope 未覆盖 tests/）
- Fix task: fix-3（修复 60+ 集成测试 fixture）
