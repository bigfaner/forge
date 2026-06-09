---
id: "3"
title: "Handle `--source-branch`, `--no-launch`, `--interactive` for existing worktrees"
priority: "P1"
estimated_time: "1h"
complexity: "medium"
dependencies: ["2"]
surface-key: ""
surface-type: "cli"
breaking: false
type: "coding.enhancement"
mainSession: false
---

# 3: Handle `--source-branch`, `--no-launch`, `--interactive` for existing worktrees

## Description

处理 worktree 已存在时各 CLI flag 的边界情况：`--source-branch` 应被忽略并输出 warning；`--no-launch` 应仅输出路径；`--interactive` 应正常进入已有 worktree。

## Reference Files
- `docs/proposals/worktree-start-idempotent/proposal.md` — Key Scenarios, Success Criteria
- `forge-cli/internal/cmd/worktree/cmd_start.go`: Handle flag edge cases in idempotent path

## Acceptance Criteria
- [ ] `--source-branch` 在 worktree 已存在时被忽略，stderr 输出 warning: `worktree already exists, ignoring --source-branch`
- [ ] `--no-launch` + worktree 已存在时，仅输出 worktree 路径（stdout），不启动 Claude，退出码 0
- [ ] `--interactive` + worktree 已存在时，正常进入该 worktree 并启动全新 Claude 会话，行为与显式指定 slug 一致

## Implementation Notes
- `--source-branch` 忽略时需明确输出 warning，避免用户误以为分支被切换
- `--no-launch` 的已有 worktree 路径应与新建 worktree 时 `--no-launch` 输出格式一致
- `--interactive` 模式下选择已有 worktree 应走与显式指定 slug 相同的 idempotent 路径
