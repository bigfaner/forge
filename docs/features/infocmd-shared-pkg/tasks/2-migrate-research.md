---
id: "2"
title: "迁移 research 命令使用 pkg/infocmd/"
priority: "P1"
estimated_time: "1.5h"
dependencies: ["1"]
scope: "backend"
breaking: false
type: "coding.refactor"
mainSession: false
---

# 2: 迁移 research 命令使用 pkg/infocmd/

## Description

将 `pkg/research/research.go` 和 `internal/cmd/research.go` 改造为使用 `pkg/infocmd/` 共享工具包。research 是最简单的 info-command（无跨引用、flat file 扫描、Discover 内部排序），适合作为第一个迁移目标验证框架设计。

## Reference Files
- `docs/proposals/infocmd-shared-pkg/proposal.md` — Source proposal
- `forge-cli/pkg/research/research.go` — 当前数据层实现
- `forge-cli/internal/cmd/research.go` — 当前命令层实现
- `forge-cli/internal/cmd/base/output.go` — 输出工具函数

## Acceptance Criteria
- [ ] `pkg/research/research.go` 使用 `infocmd.Discover` 和 `infocmd.FindByID` 替代手动实现
- [ ] `internal/cmd/research.go` 的列表/详情渲染保持原有输出格式
- [ ] `forge research` 命令输出与重构前逐字节一致
- [ ] 现有 research 相关测试全部通过
- [ ] `pkg/research/` 中不再有 `parseFrontmatter` 的独立副本

## Hard Rules
- CLI 输出格式零行为变更
- 不改变 `cobra.Command` 的注册方式

## Implementation Notes

### research 的特殊性
- Flat file 扫描：`docs/research/*.md`
- Slug 从文件名派生
- Discover 内部已有排序
- 无跨引用逻辑

### 迁移策略
1. 在 `pkg/research/` 中定义 `Report` 的 `ScanConfig`，调用 `infocmd.Discover`
2. `FindBySlug` 委托给 `infocmd.FindByID`，`IDKey = func(r Report) string { return r.Slug }`
3. `internal/cmd/research.go` 保持原有渲染逻辑，仅数据获取方式改变
