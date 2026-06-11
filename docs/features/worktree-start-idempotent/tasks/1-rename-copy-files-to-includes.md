---
id: "1"
title: "Rename `copy-files` → `includes` across codebase"
priority: "P0"
estimated_time: "1h"
complexity: "medium"
dependencies: []
surface-key: ""
surface-type: "cli"
breaking: false
type: "coding.cleanup"
mainSession: false
---

# 1: Rename `copy-files` → `includes` across codebase

## Description

将 `worktree.copy-files` 配置项重命名为 `worktree.includes`，对齐 Claude Code 的 `.worktreeinclude` 命名风格。此为纯机械替换，不改变任何运行时行为。

重命名范围：Go 结构体字段（`CopyFiles` → `Includes`）、YAML tag（`copy-files` → `includes`）、函数名（`validateCopyFiles` → `validateIncludes`、`copyFilesToWorktree` → `copyIncludesToWorktree`）、JSON schema、示例 YAML、以及所有测试文件中的引用。

## Reference Files
- `docs/proposals/worktree-start-idempotent/proposal.md` — Proposed Solution, Constraints & Dependencies, Success Criteria
- `forge-cli/pkg/forgeconfig/config.go`: Rename `CopyFiles` field to `Includes`, YAML tag `copy-files` → `includes`
- `forge-cli/internal/cmd/worktree/helpers.go`: Rename `validateCopyFiles` → `validateIncludes`, `copyFilesToWorktree` → `copyIncludesToWorktree`
- `forge-cli/internal/cmd/worktree/cmd_start.go`: Update all `copyFiles` variable references
- `forge-cli/internal/cmd/init_config.go`: Update config init UI labels and field references

## Acceptance Criteria
- [ ] `WorktreeConfig.CopyFiles` renamed to `WorktreeConfig.Includes`, YAML tag changed to `includes`
- [ ] `grep -r "copy-files\|CopyFiles\|copy_files" forge-cli/ --include="*.go"` returns zero results
- [ ] JSON schema and example YAML updated to use `includes`
- [ ] All existing tests pass after rename

## Implementation Notes
- 直接替换，不保留 `copy-files` 旧字段的兼容逻辑
- 全量搜索替换覆盖约 62 处引用（10 个文件）
- 注意 `copyFilesToWorktree` / `validateCopyFiles` 函数名及函数变量也要重命名
