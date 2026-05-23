---
id: "4"
title: "迁移 lesson 命令使用 pkg/infocmd/"
priority: "P1"
estimated_time: "1.5h"
dependencies: ["3"]
scope: "backend"
breaking: false
type: "coding.refactor"
mainSession: false
---

# 4: 迁移 lesson 命令使用 pkg/infocmd/

## Description

将 `pkg/lesson/lesson.go` 和 `internal/cmd/lesson.go` 改造为使用 `pkg/infocmd/` 共享工具包。lesson 的特殊性在于使用 `Name` 而非 `Slug` 作为标识符，且有 `Category`（从文件名前缀推断）和 `Tags` 字段。

## Reference Files
- `docs/proposals/infocmd-shared-pkg/proposal.md` — Source proposal
- `forge-cli/pkg/lesson/lesson.go` — 当前数据层实现（含 inferCategory、categoryPrefixes）
- `forge-cli/internal/cmd/lesson.go` — 当前命令层实现
- `forge-cli/internal/cmd/base/output.go` — 输出工具函数

## Acceptance Criteria
- [ ] `pkg/lesson/lesson.go` 使用 `infocmd.Discover` 和 `infocmd.FindByID` 替代手动实现
- [ ] lesson 的 `Name` 标识符正确映射（而非 Slug）
- [ ] lesson 的 `Category` 推断逻辑保留（从文件名前缀）
- [ ] lesson 的 `Date` → `Created` fallback 保留
- [ ] `forge lesson` 命令输出与重构前逐字节一致
- [ ] 现有 lesson 相关测试全部通过
- [ ] `pkg/lesson/` 中不再有 `parseFrontmatter` 的独立副本

## Hard Rules
- CLI 输出格式零行为变更
- 不改变 `cobra.Command` 的注册方式

## Implementation Notes

### lesson 的特殊性
- Flat file 扫描：`docs/lessons/*.md`
- 使用 `Name` 而非 `Slug`（但语义相同：filename minus `.md`）
- `inferCategory` 从文件名前缀推断分类（如 `feedback_*` → `feedback`）
- Created 有双重 fallback：先 `meta.Date`，再 mtime
- metadata 中有 `Tags []string` 和 `Severity string`

### 迁移策略
1. `IDKey = func(l Lesson) string { return l.Name }`
2. `ParseEntry` 中处理 Category 推断和 Date fallback
3. `FindByName` 委托给 `infocmd.FindByID`
