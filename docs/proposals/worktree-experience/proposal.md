---
created: 2026-05-20
author: fanhuifeng
status: Draft
---

# Proposal: Worktree Experience Enhancement

## Problem

`forge worktree` 的四个子命令（start/list/remove/resume）缺少关键交互能力，导致并行开发工作流效率低下：用户需要手动输入 slug、手动 git push、无法快速查看 worktree 状态、无法批量清理分支和残留记录。

### Evidence

1. 无 shell 补全——每次操作都需要手动输入完整 slug，容易打错
2. 无交互选择——需要先跑 `forge proposal` 或 `forge feature list` 查看 slug 再回来操作
3. `remove` 只删目录不删分支，也不 prune，残留记录累积
4. 无 push 命令——创建 worktree 后需手动到目录里 git push
5. 无 status 命令——无法快速查看某个 worktree 有无未提交代码
6. `start` 始终启动 claude，无法仅创建环境用于脚本或批量操作
7. `resume` 只重新启动 claude，无法恢复之前的会话上下文

### Urgency

worktree 是并行开发的核心工具，当前体验阻碍了多 feature 并行的效率。每次操作多出的手动步骤在频繁使用时累积显著。

## Proposed Solution

对 `forge worktree` 进行全面增强：新增 3 个子命令（push/status 交互模式）、增强 2 个现有命令（remove --hard、start --no-launch）、添加 shell 补全、集成 Claude Code 会话管理。

### Innovation Highlights

- **Cobra ValidArgsFunction 动态补全**：根据上下文提供不同补全源（start → 未完成 proposal/feature；remove/resume → 现存 worktree）
- **Claude -c 会话恢复**：利用 Claude Code 的 `-c <sessionName>` 能力，将 slug 作为会话标识，实现 resume 时恢复之前的对话上下文
- **分支优先创建**：先 `git checkout -b` 再 `git worktree add`，解耦分支与 worktree 的创建，使中间步骤可观测

## Requirements Analysis

### Key Scenarios

1. **交互创建**：用户运行 `forge worktree start -i`，从下拉列表选择未完成的 proposal 或 feature，自动创建 worktree
2. **仅创建环境**：CI 脚本或批量操作时 `forge worktree start <slug> --no-launch`，只创建分支和 worktree
3. **彻底清理**：`forge worktree remove <slug> --hard`，删除 worktree + 本地分支 + prune 残留
4. **推送分支**：`forge worktree push` 在 worktree 目录下推送当前分支到远端
5. **查看状态**：`forge worktree status <slug>` 显示分支、提交状态、未提交文件
6. **恢复会话**：`forge worktree resume <slug>` 用 `claude -c <slug>` 恢复之前的 Claude Code 会话
7. **Tab 补全**：start 补全未完成 proposal/feature；remove/resume 补全现存 worktree 列表

### Non-Functional Requirements

- 补全响应时间 < 200ms（避免用户感知延迟）
- status 命令不修改任何文件系统状态（只读操作）

### Constraints & Dependencies

- 依赖 Cobra 的 shell completion 机制（ValidArgsFunction）
- 依赖 Claude Code 的 `-c <sessionName>` 会话恢复能力
- `--hard` 删除分支前需确认无未提交更改

## Alternatives & Industry Benchmarking

### Industry Solutions

Git worktree 是标准 Git 功能，大多数工具（如 gh、git-town）通过 alias 和 hook 简化操作。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零成本 | 效率持续低下 | Rejected: 痛点明确且累积 |
| 外部 alias/hook | git-town 等 | 不改 forge | 功能分散，补全不统一 | Rejected: 体验不连贯 |
| **forge worktree 全面增强** | 本提案 | 统一体验，补全集成 | 需改动 forge-cli | **Selected: 最佳 ROI** |

## Feasibility Assessment

### Technical Feasibility

完全可行。所有改动基于现有 Cobra 框架和 Git CLI，无外部依赖。Claude Code `-c` 是已有能力。

### Resource & Timeline

预计 7-8 个 coding 任务，适合 quick mode 单次完成。

### Dependency Readiness

所有依赖已就绪：Cobra 补全 API、Git CLI、Claude Code CLI。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| 用户需要 tmux 管理会话 | XY Detection | Overturned：用户实际需求是 Claude -c 会话恢复，不需要 tmux |
| start 应自动 push 新分支 | Assumption Flip | Overturned：用户明确选择不自动 push，保持手动控制 |
| remove 删除远端分支 | Stress Test | Overturned：只需清理本地，远端分支保留以便回溯 |

## Scope

### In Scope

- `forge worktree start -i/--interactive`：交互选择未完成 proposal/feature
- `forge worktree start --no-launch`：仅创建 worktree 不启动 claude
- `forge worktree remove --hard`：删除 worktree + 本地分支 + prune
- `forge worktree push`：推送当前 worktree 分支到远端
- `forge worktree status <slug>`：显示分支、提交状态、未提交文件
- `forge worktree resume`：用 `claude -c <slug>` 恢复会话
- `start` 分支创建顺序调整：先创建分支再创建 worktree
- Shell 补全：start/remove/resume 子命令的 slug 动态补全

### Out of Scope

- tmux 会话管理
- 退出会话自动清理 worktree
- `--hard` 删除远端分支
- `start` 自动 push 到远端
- worktree 级别的配置文件隔离

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Claude -c 的 sessionName 格式限制 | M | M | 先验证 Claude Code 会话命名机制，不支持时回退到当前行为 |
| `--hard` 误删有价值分支 | M | H | 删除前检查未提交更改，有更改时拒绝或要求 --force |
| 补全在大型 repo 中延迟 | L | L | proposal/feature 列表通常 < 50 项，性能可控 |

## Success Criteria

- [ ] `forge worktree start -i` 能列出并选择未完成 proposal/feature
- [ ] `forge worktree start <slug> --no-launch` 创建 worktree 但不启动 claude
- [ ] `forge worktree remove <slug> --hard` 删除 worktree + 分支 + prune，3 步全部成功
- [ ] `forge worktree push` 能将当前分支 push 到 origin
- [ ] `forge worktree status <slug>` 显示分支名、最新提交、未提交文件列表
- [ ] `forge worktree resume <slug>` 使用 `claude -c <slug>` 恢复会话
- [ ] Tab 补全在 start/remove/resume 子命令上正常工作
- [ ] 所有新命令有对应的单元测试

## Next Steps

- Proceed to `/quick-tasks` to generate implementation tasks
