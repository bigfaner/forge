---
id: "7"
title: "Cross-module dependency audit and CI check"
priority: "P0"
estimated_time: "1.5h"
complexity: "medium"
dependencies: []
surface-key: "."
surface-type: "cli"
breaking: false
type: "coding.enhancement"
mainSession: false
---

# 7: Cross-module dependency audit and CI check

## Description
Phase 2a 前置条件：审计 monorepo 内是否存在其他 Go 模块 import `forge-cli` 的 `internal/` 或 `pkg/` 包。审计完成后在 `Makefile` 中添加 `check-cross-module-deps` target 作为 CI check，为 Phase 2c 的"不保留兼容层"决策提供持续保护。

## Reference Files
- forge-cli/go.mod: module path 为 `forge-cli`，审计须搜索此路径被引用的情况 (source: proposal.md#Dependency-Readiness)
- Makefile: 需添加 `check-cross-module-deps` target (source: proposal.md#Scope item 15)

## Acceptance Criteria
- [ ] 审计已完成：`go list -m` 和 `go mod graph` 检测 Go module 级依赖，`grep -rn 'forge-cli/internal\|forge-cli/pkg' --include='*.go'` 搜索 monorepo 根目录
- [ ] 审计结果记录在 `docs/features/forge-cli-codebase-standards/cross-module-audit.md`
- [ ] `Makefile` 新增 `check-cross-module-deps` target：`grep -rn '"forge-cli/internal'` 搜索 monorepo 排除 `forge-cli/` 自身，返回非零则构建失败
- [ ] fallback 评估已完成：若无跨模块依赖则 Phase 2c 可完整执行；若有则记录影响范围

## Implementation Notes
- 审计方法须包含两层：Go module 级（`go list -m`）和文本级（`grep`）
- CI check 应能在 monorepo 根目录下运行，不限于 `forge-cli/` 目录
