---
id: "2"
title: "Make `forge worktree start` idempotent — core behavior"
priority: "P0"
estimated_time: "1.5h"
complexity: "high"
dependencies: ["1"]
surface-key: ""
surface-type: "cli"
breaking: false
type: "coding.enhancement"
mainSession: false
---

# 2: Make `forge worktree start` idempotent — core behavior

## Description

修改 `forge worktree start <slug>` 的核心逻辑，使其在 worktree 已存在时不再报错，而是跳过创建并启动全新 Claude 会话。复用 `cmd_resume.go` 中的 worktree 验证逻辑（symlink 解析 + `.git` 文件检查），确保已有 worktree 的有效性。

**核心变更**：将 `cmd_start.go` 中 `os.Stat(targetDir)` 已存在的分支从「报错退出」改为「验证 → 启动新会话」。

## Reference Files
- `docs/proposals/worktree-start-idempotent/proposal.md` — Proposed Solution, Key Scenarios, Key Risks, Success Criteria
- `forge-cli/internal/cmd/worktree/cmd_start.go`: Modify "already exists" branch from error to idempotent path
- `forge-cli/internal/cmd/worktree/cmd_resume.go`: Reuse worktree validation (symlink resolve + `.git` check, lines 58-71)
- `forge-cli/internal/cmd/worktree/helpers.go`: `copyIncludesToWorktree` / `validateIncludes` helpers

## Acceptance Criteria
- [ ] `forge worktree start <slug>` 在 worktree 已存在时不再报错，而是跳过创建并启动全新 Claude 会话
- [ ] 已存在时 stderr 输出包含关键词 `entering existing worktree`；新建时 stderr 输出包含关键词 `created new worktree`
- [ ] `includes` 文件复制在 worktree 已存在时被跳过（文件已在首次创建时复制）
- [ ] worktree 目录存在但 `.git` 文件缺失或损坏时，`start` 报错退出并提示建议 `forge worktree remove` 后重建
- [ ] worktree 不存在时的行为与当前完全一致（回归测试通过）

## Implementation Notes
- 复用 `cmd_resume.go` 的验证模式：`filepath.EvalSymlinks` + `os.Stat(.git)`
- 已存在路径中跳过 `git worktree add` 和 `copyIncludesToWorktree` 两个步骤
- 已存在时不进入交互式配置流程（配置已在首次创建时完成）
- 参考 `cmd_resume.go` lines 58-71 的验证代码结构
