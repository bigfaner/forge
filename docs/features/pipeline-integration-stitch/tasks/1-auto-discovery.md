---
id: "1"
title: "Auto-discovery: 替换 prompt.go 和 autogen.go 手写 map + init-time 校验 + clean-code.md 重命名"
priority: "P0"
estimated_time: "2h"
dependencies: []
scope: "backend"
breaking: true
type: "coding.refactor"
mainSession: false
---

# 1: Auto-discovery 机制重构

## Description

将 `prompt.go` 的 `typeToTemplate` 和 `autogen.go` 的 `autogenTypeToFile` 两个手写 map 替换为基于命名约定的自动发现。约定：`strings.ReplaceAll(typeName, ".", "-") + ".md"`。重命名 `prompt/data/clean-code.md` → `code-quality-simplify.md` 消除唯一例外，实现零 override。在 CLI 入口（`main()` 函数）添加 init-time 校验：遍历所有已注册类型常量，验证模板文件存在且映射唯一，缺失/碰撞时 fatal exit。

## Reference Files
- `docs/proposals/pipeline-integration-stitch/proposal.md` — Source proposal
- `forge-cli/pkg/prompt/prompt.go` — 当前手写 map
- `forge-cli/pkg/task/autogen.go` — 当前手写 map
- `forge-cli/pkg/task/types.go` — 类型常量定义（自动发现的数据源）

## Acceptance Criteria
- [ ] `prompt.go` 中 `typeToTemplate` map 被移除，`Synthesize()` 通过命名约定自动推导模板路径
- [ ] `autogen.go` 中 `autogenTypeToFile` map 被移除，模板路径通过命名约定自动推导
- [ ] `prompt/data/clean-code.md` 重命名为 `code-quality-simplify.md`
- [ ] CLI 入口添加 `ValidateTemplateConventions()` 函数，遍历 `ValidTypes` 中所有类型，验证 embed.FS 中模板文件存在且映射唯一
- [ ] 缺失模板或映射碰撞时 CLI 启动失败并报告具体类型名
- [ ] `grep -r "clean-code" forge-cli/` 仅匹配重命名后的 `code-quality-simplify`，无残留 `clean-code.md` 引用
- [ ] 所有现有测试通过

## Hard Rules
- init-time 校验必须在 CLI 入口（`main()` 函数启动路径），不得使用 `init()` 函数
- 重命名 clean-code.md 后必须 grep 验证零残留引用

## Implementation Notes
- 自动发现核心：`"data/" + strings.ReplaceAll(typeName, ".", "-") + ".md"`
- init-time 校验需同时检查：文件存在（ReadFile）+ 映射唯一性（无重复文件名）
- 选择 CLI 入口而非 init() 的理由：(1) init() 在 go test 时也执行；(2) 测试可以独立运行；(3) CI 环境解耦
