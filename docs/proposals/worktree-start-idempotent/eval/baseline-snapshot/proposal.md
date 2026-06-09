---
created: "2026-06-09"
author: fanhuifeng
status: Draft
intent: "enhancement"  # 修改 start 命令的现有语义，使其幂等；重命名配置项
---

# Proposal: Make `forge worktree start` Idempotent & Rename `copy-files` → `includes`

## Problem

两个问题：

1. `forge worktree start <slug>` 在 worktree 已存在时报错退出，缺少一条命令能「进入已有 worktree 并开启全新 Claude 会话」
2. `worktree.copy-files` 配置命名描述的是实现行为（"复制文件"）而非意图（"包含什么"），与 Claude Code 的 `.worktreeinclude` 命名风格不一致

### Evidence

1. `start` 在 `.forge/worktrees/<slug>` 已存在时直接报错：`worktree already exists, use "resume" instead`
2. `resume` 使用 `claude -c` 恢复上一次会话，无法开启全新会话
3. 手动 `cd .forge/worktrees/<slug> && claude --dangerously-skip-permissions` 可行但破坏了统一的 CLI 入口
4. `copy-files` 命名在功能扩展时显得局限（如果未来支持 glob 模式或目录，"copy-files" 不再准确）

### Urgency

日常高频操作——每次需要在新会话中继续某个 feature 工作时都会遇到。Cost of delay：持续的体验摩擦，用户要么每次手动 cd，要么接受 resume 旧会话的上下文噪音。

## Proposed Solution

将 `forge worktree start <slug>` 改为幂等操作，并将 `worktree.copy-files` 重命名为 `worktree.includes`：

**幂等 start：**
- **worktree 不存在** → 创建 worktree + 启动全新 Claude 会话（现有行为不变）
- **worktree 已存在** → 跳过创建，直接启动全新 Claude 会话（新行为）

**配置重命名：**
- `worktree.copy-files` → `worktree.includes`（对齐 Claude Code 的 `.worktreeinclude` 命名风格）
- 直接替换，不保留任何旧字段兼容逻辑

`start` 成为「全新会话」的统一入口，`resume` 保持为「恢复旧会话」的专用命令。

### Innovation Highlights

幂等设计，不是业界首创，但符合 CLI 最佳实践（`mkdir -p`、`touch` 都是幂等操作）。核心洞察是 `start` 的语义应聚焦于「启动新会话」而非「创建 worktree」——worktree 创建是实现细节，用户真正关心的是开始工作。

## Requirements Analysis

### Key Scenarios

1. **Happy path: worktree 不存在** → 创建 + 启动新会话（与当前行为一致）
2. **Happy path: worktree 已存在** → 跳过创建 + 启动新会话
3. **已有 worktree 但用了 `--source-branch`** → 忽略该 flag（worktree 已存在，分支已确定），输出提示信息
4. **已有 worktree 但配置了 `includes`** → 跳过文件复制（文件应已在首次创建时复制）
5. **`--no-launch` + worktree 已存在** → 仅验证 worktree 存在并输出路径，不启动 Claude
6. **`--interactive` + worktree 已存在** → 正常进入，因为 slug 来自交互选择

### Non-Functional Requirements

- 向后兼容：`start` 在 worktree 不存在时的行为完全不变
- 性能：无影响（只是跳过了创建步骤）

### Constraints & Dependencies

- 修改 `forge-cli/internal/cmd/worktree/cmd_start.go`（幂等逻辑）
- 修改 `forge-cli/pkg/forgeconfig/` 中的配置结构体（`CopyFiles` → `Includes` + 兼容读取）
- 依赖现有的 worktree 验证逻辑（resume 已有实现可复用）

## Alternatives & Industry Benchmarking

### Industry Solutions

幂等 CLI 命令是常见模式：`kubectl apply`、`terraform plan` 都采用「已存在则跳过」的策略。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零成本 | 持续的体验摩擦，手动 cd 或被迫 resume | Rejected: 高频痛点 |
| 新增 `open` 子命令 | — | 职责分离清晰 | 增加命令数量，用户需记住 start vs open 的区别 | Rejected: 认知负担大于收益 |
| 给 `resume` 加 `--fresh` flag | — | 不增加命令 | 语义矛盾（resume + fresh 自相矛盾） | Rejected: 命名不直观 |
| **幂等 `start`** | `mkdir -p` 模式 | 零新命令、直觉行为、幂等 | 轻微改变 start 的现有语义 | **Selected: 最符合「最小惊讶原则」** |

## Feasibility Assessment

### Technical Feasibility

完全可行。实现上只需修改 `cmd_start.go` 中的「已存在则报错」分支，将其改为「已存在则跳过创建并启动 Claude」。可复用 `cmd_resume.go` 中的 worktree 验证逻辑。

### Resource & Timeline

改动范围较小（2-3 个文件，约 40 行变更），预计 1 小时内完成。

### Dependency Readiness

无外部依赖。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| `start` 的语义是「创建 worktree」 | Assumption Flip | Overturned: 用户视角的语义是「开始工作」，创建 worktree 只是实现手段 |
| 已有 worktree 时 `start` 报错是安全行为 | XY Detection | Refined: 报错保护的是「防止重复创建」，而非「防止进入」。幂等跳过同样安全 |
| 需要一个新命令来解决 | Occam's Razor | Overturned: 修改现有命令的语义比新增命令更简单 |
| `copy-files` 是好名字 | Assumption Flip | Overturned: 命名描述了实现而非意图。Claude Code 的 `.worktreeinclude` 更好——`includes` 表达「包含什么」 |

## Scope

### In Scope

- 修改 `forge worktree start <slug>` 在 worktree 已存在时的行为：跳过创建，启动全新 Claude 会话
- 当 worktree 已存在时忽略 `--source-branch`，输出提示信息
- 当 worktree 已存在时跳过 `includes` 文件复制
- 输出区分性日志，让用户知道是「创建了新 worktree」还是「进入了已有 worktree」
- 重命名配置项 `worktree.copy-files` → `worktree.includes`（直接替换，不保留旧字段）

### Out of Scope

- 修改 `resume` 命令的行为
- 添加新的 CLI 子命令
- 修改 worktree 的创建/验证逻辑本身
- `start` 的 `--interactive` 模式变更
- 支持 glob 模式或目录包含（未来可扩展）
- 采用 `.worktreeinclude` 独立文件格式（保持 YAML 配置统一）

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 用户脚本依赖「已存在则报错」的行为 | L | M | 在 release notes 中标注行为变更；日志明确输出"已存在，跳过创建" |
| `--source-branch` 被静默忽略导致困惑 | M | L | 输出 warning: "worktree already exists, ignoring --source-branch" |
| `includes` 被跳过但用户预期更新 | L | L | 首次创建时已复制，后续不应再覆盖（会丢失 worktree 内的修改） |

## Success Criteria

- [ ] `forge worktree start <slug>` 在 worktree 已存在时不再报错，而是直接启动全新 Claude 会话
- [ ] 已存在时输出明确的提示信息，区分「新建」和「进入」两种路径
- [ ] `--source-branch` 在 worktree 已存在时被忽略并输出 warning
- [ ] `includes` 在 worktree 已存在时被跳过
- [ ] worktree 不存在时的行为与当前完全一致（回归测试通过）
- [ ] `worktree.includes` 配置项正常工作
- [ ] 代码中不存在任何 `copy-files` / `CopyFiles` 兼容逻辑

## Next Steps

- Proceed to `/write-prd` to formalize requirements
