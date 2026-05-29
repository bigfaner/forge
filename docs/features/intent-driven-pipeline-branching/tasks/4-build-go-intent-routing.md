---
id: "4"
title: "build.go — intent-driven pipeline routing"
priority: "P0"
estimated_time: "2-3h"
complexity: "high"
dependencies: []
surface-key: "."
surface-type: "cli"
breaking: true
type: "coding.feature"
mainSession: false
---

# 4: build.go — intent-driven pipeline routing

## Description

在 `build.go` 中实现 intent 驱动的 pipeline 路由。核心改动：
1. `BuildIndexOpts` 新增 `Intent` 字段，由 CLI handler 从 `Proposal.Intent` 传入
2. `needsTestPipeline()` 增加 `intent` 参数，当 intent 为 `refactor`/`cleanup` 时直接返回 `false`
3. `detectMode()` 在 intent=`cleanup` 时强制返回 Quick 模式，忽略文档存在性
4. `IsTestableType()` 不修改，保持纯类型判断职责

数据流：CLI handler (`internal/cmd/task/index.go`) 调用 `proposal.FindBySlug()` → 获取 `Proposal.Intent` → 赋值给 `BuildIndexOpts.Intent` → 传入 `BuildIndex()` → `needsTestPipeline(tasks, intent)` 和 `detectMode()` 感知 intent。

## Reference Files
- `forge-cli/pkg/task/build.go`: `BuildIndexOpts` 结构体（L17-23）、`detectMode()` 函数（L452-461）、`needsTestPipeline()` 函数（L526-536）、`IsTestableType()` 函数（L518-521） (source: proposal.md#In-Scope, items 5-6)
- `forge-cli/internal/cmd/task/index.go`: CLI handler，需读取 proposal intent 并传入 BuildIndexOpts（L60-66 当前无 Intent 传递） (source: proposal.md#Feasibility-Assessment, CLI 层第 5 点)
- `forge-cli/pkg/proposal/proposal.go`: `Proposal` 结构体已有 `Intent` 字段（L20），`FindBySlug()` 可直接使用 (source: proposal.md#Constraints-&-Dependencies)
- `docs/proposals/intent-driven-pipeline-branching/proposal.md#Feasibility-Assessment`: 完整数据流定义——CLI handler 从 proposal frontmatter 读取 intent，传入 BuildIndex，不重复解析

## Acceptance Criteria
- [ ] `BuildIndexOpts` 新增 `Intent string` 字段（与 `Mode` 并列），CLI handler 从 `Proposal.Intent` 赋值，空值默认 `new-feature`
- [ ] `needsTestPipeline()` 签名变为 `needsTestPipeline(tasks map[string]Task, intent string) bool`，当 intent 为 `refactor`/`cleanup` 时直接返回 `false`（不遍历 tasks 也不调用 `IsTestableType()`）
- [ ] `detectMode()` 在 intent=`cleanup` 时强制返回 `"quick"`，忽略 PRD 文档存在性
- [ ] `IsTestableType()` 不修改，保持纯类型判断职责（对所有 `coding.*` 返回 `true`）
- [ ] 向后兼容：缺少 intent 字段的已有 proposal（`opts.Intent` 为空）默认走 `new-feature` 管道，行为与当前完全一致

## Hard Rules
- `IsTestableType()` 不修改——保持纯类型判断职责，intent 短路逻辑放在 `needsTestPipeline()` 中

## Implementation Notes

### Test Impact
- Affected test suite(s): `forge-cli/internal/cmd/`, `forge-cli/pkg/task/`
- Expected fixture changes: 新增 intent 相关的 test fixture（proposal.md 文件含不同 intent 值）
- Risk level: medium

- `forge task add` 中的 `BuildIndex` 调用（`internal/cmd/task/add.go:209`）也需要传递 Intent——需要从已有 index 或 proposal 中获取 intent
- `BuildIndex()` 内部不再重复解析 frontmatter，直接使用 `opts.Intent`
- proposal.md 不存在时（用户在未完成 brainstorm 的特征目录上运行 `forge task index`），`proposal.FindBySlug()` 返回空 Proposal，CLI handler 将 `opts.Intent` 设为默认值 `new-feature`
