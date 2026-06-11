---
id: "5"
title: "清理旧包中的重复代码"
priority: "P2"
estimated_time: "1h"
dependencies: ["4"]
scope: "backend"
breaking: false
type: "coding.cleanup"
mainSession: false
---

# 5: 清理旧包中的重复代码

## Description

三个命令全部迁移完成后，清理 `pkg/research/`、`pkg/proposal/`、`pkg/lesson/` 中不再需要的重复代码。确保这些包仅保留数据模型定义和调用 `pkg/infocmd/` 的薄胶水层。

## Reference Files
- `docs/proposals/infocmd-shared-pkg/proposal.md` — Source proposal
- `forge-cli/pkg/research/research.go` — 待清理
- `forge-cli/pkg/proposal/proposal.go` — 待清理
- `forge-cli/pkg/lesson/lesson.go` — 待清理

## Acceptance Criteria
- [ ] `parseFrontmatter` 仅存在于 `pkg/infocmd/`，三个旧包中不再有副本
- [ ] 三个旧包中不再有手动实现的 `Discover()` 和 `FindByXxx()` 函数
- [ ] 旧包仅保留：exported struct 定义 + `ScanConfig` 构造 + 薄包装函数
- [ ] `go vet ./...` 通过
- [ ] 所有测试通过
- [ ] 新增 info-command 时只需定义 struct + 列配置 + ~30 行胶水代码（验证标准）

## Hard Rules
- 不删除任何 exported struct 或 field（保持 API 兼容）
- 不改变任何 CLI 输出

## Implementation Notes

清理目标：
1. 删除三个包中的 `parseFrontmatter` 函数
2. 删除三个包中的内部 `metadata` struct（如果 parseFrontmatter 迁移后不再需要）
3. 删除 `proposal/` 中的 `fileExists` helper（如果已移入 infocmd）
4. 删除 `lesson/` 中的 `lessonWithMeta` / research 中的 `reportWithMeta` 内部 struct（排序逻辑已统一到 infocmd）
5. 确保每个旧包的代码量显著减少（目标：每个包 < 50 行有效代码）
