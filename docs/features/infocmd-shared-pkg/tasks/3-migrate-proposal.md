---
id: "3"
title: "迁移 proposal 命令使用 pkg/infocmd/"
priority: "P1"
estimated_time: "2h"
dependencies: ["2"]
scope: "backend"
breaking: false
type: "coding.refactor"
mainSession: false
---

# 3: 迁移 proposal 命令使用 pkg/infocmd/

## Description

将 `pkg/proposal/proposal.go` 和 `internal/cmd/proposal.go` 改造为使用 `pkg/infocmd/` 共享工具包。proposal 是最复杂的 info-command（subdirectory 扫描、跨引用检查 PRD/Feature 状态、排序在 command 层），需要保留所有特有逻辑。

## Reference Files
- `docs/proposals/infocmd-shared-pkg/proposal.md` — Source proposal
- `forge-cli/pkg/proposal/proposal.go` — 当前数据层实现（含 fileExists、跨引用逻辑）
- `forge-cli/internal/cmd/proposal.go` — 当前命令层实现（排序在 cmd 层、Extra columns）
- `forge-cli/internal/cmd/base/output.go` — 输出工具函数

## Acceptance Criteria
- [ ] `pkg/proposal/proposal.go` 使用 `infocmd.Discover` 和 `infocmd.FindByID` 替代手动实现
- [ ] proposal 详情中的 PRD（yes/no）和 Feature 状态字段保留
- [ ] proposal 列表排序（Created 降序）保持不变
- [ ] `forge proposal` 命令输出与重构前逐字节一致
- [ ] 现有 proposal 相关测试全部通过
- [ ] `pkg/proposal/` 中不再有 `parseFrontmatter` 的独立副本

## Hard Rules
- CLI 输出格式零行为变更
- 保留 proposal 的跨引用逻辑（PRD 检查、Feature 状态查询）
- proposal 的错误处理风格（`Exit()` 调用）不在此次重构范围内

## Implementation Notes

### proposal 的特殊性
- Subdirectory 扫描：`docs/proposals/*/proposal.md`
- 使用 `feature.ProposalBaseDir` 和 `feature.ProposalFileName` 常量
- Discover 无排序，排序在 command 层完成
- 跨引用：检查 `features/<slug>/prd/spec.md` 是否存在（HasPRD）和 `features/<slug>/manifest.yaml` 的 status（FeatureStatus）
- Created 有 mtime fallback（直接写入 Proposal 结构体）

### 迁移策略
1. `ScanConfig.IsSubdir = true`，`FileName = "proposal.md"`
2. `ParseEntry` 中处理跨引用逻辑（检查 PRD 文件和 manifest）
3. 排序可移入 Discover（统一行为），或在 cmd 层继续手动排序——取决于 task 1 的设计
4. `FindBySlug` 委托给 `infocmd.FindByID`
