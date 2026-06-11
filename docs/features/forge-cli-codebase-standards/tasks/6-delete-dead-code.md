---
id: "6"
title: "Delete dead code and unify Debugf"
priority: "P1"
estimated_time: "2h"
complexity: "medium"
dependencies: [5]
surface-key: "."
surface-type: "cli"
breaking: false
type: "coding.cleanup"
mainSession: false
---

# 6: Delete dead code and unify Debugf

## Description
Phase 2a：删除纯粹的死代码——deprecated `Scope` 字段（`pkg/task/frontmatter.go`）、重复的 `Debugf` 定义（保留 `internal/cmd/base/output.go` 中的版本，删除 `internal/cmd/output.go` 中的版本并改为转发）、`.out` 构建产物（`cmd.out`、`cout.out`、`just.out`）。同时确保 `getTaskPhase` 不被误删（有 5 处生产调用）。

## Reference Files
- forge-cli/pkg/task/frontmatter.go:28-30: deprecated `Scope` 字段需删除，检查 `CheckLegacyScope` 是否仍需此字段 (source: proposal.md#Evidence)
- forge-cli/internal/cmd/output.go:13: 重复的 `Debugf` 定义 (source: proposal.md#Evidence)
- forge-cli/internal/cmd/base/output.go:85: 保留此 `Debugf` 定义 (source: proposal.md#Evidence)
- forge-cli/cmd.out, forge-cli/cout.out, forge-cli/just.out: 构建产物需删除并加入 .gitignore (source: proposal.md#Evidence)

## Acceptance Criteria
- [ ] `pkg/task/frontmatter.go` 中 deprecated `Scope` 字段已删除；`CheckLegacyScope` 迁移为不依赖 `Scope` 的实现
- [ ] `internal/cmd/output.go` 中重复的 `Debugf` 已删除，调用点改为导入 `internal/cmd/base` 的版本
- [ ] `cmd.out`、`cout.out`、`just.out` 已从仓库中删除并加入 `.gitignore`
- [ ] `go build ./...` 和 `go test ./...` 全部通过
- [ ] `getTaskPhase`、`checkExistingTaskState`、`compareVersionIDs` 别名未被误删（仅留待 Task 9 处理）

## Hard Rules
- 仅删除以下文件中列出的死代码项，不涉及其他清理
- `getTaskPhase` 不得删除（`validate_index.go` 有 5 处生产调用）

## Implementation Notes
- 删除 `Scope` 字段前需确认 `CheckLegacyScope` 的迁移方案
- Debugf 统一后，检查所有 `internal/cmd/` 下的导入是否正确更新
- `.out` 文件删除后确认 `.gitignore` 包含 `*.out` 规则

### Test Impact
- Affected test suite(s): `forge-cli/internal/cmd/`, `forge-cli/pkg/task/`
- Expected fixture changes: 无
- Risk level: low（纯删除操作，编译验证即可）
