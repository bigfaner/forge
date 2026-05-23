---
title: "Forge CLI 命令参考"
domains: [cli, commands, reference, skills]
---

# Forge CLI 命令参考

本文档记录所有有效的 `forge` CLI 命令，供 skill 作者编写或修改 skill 时参考。

> **权威来源**：`forge-cli/internal/cmd/root.go` 中的 `init()` 函数。本文档必须与此文件保持同步。

## 顶层命令

| 命令 | 用途 | 源文件 |
|------|------|--------|
| `forge init` | 初始化 forge 项目环境 | `init.go` |
| `forge config` | 管理 forge 配置 | `config.go` |
| `forge feature [slug]` | 设置或显示当前 feature | `feature.go` |
| `forge proposal [slug]` | 列出或查看 proposal 详情 | `proposal.go` |
| `forge lesson [name]` | 列出或查看 lesson 详情 | `lesson.go` |
| `forge cleanup` | 清理已完成的任务状态 | `cleanup.go` |
| `forge probe` | ~~HTTP 健康检查~~ — 已移除，功能由 `pkg/e2eprobe` 库提供 | _(已移除)_ |
| `forge quality-gate` | 检查所有任务是否完成，然后运行测试 | `quality_gate.go` |
| `forge verify-task-done` | 在 git commit 前验证任务完成状态 | `verify_task_done.go` |
| `forge version` | 打印 CLI 版本号（隐藏命令，不出现在 --help 中） | `version.go` |
| `forge research [slug]` | 列出或查看 research report 详情 | `research.go` |
| `forge claude` | 跳过权限检查启动 Claude CLI | `claude.go` |

## 命令组

### forge task — 任务生命周期管理

源文件：`task_parent.go`

| 命令 | 用途 | 源文件 |
|------|------|--------|
| `forge task claim` | 认领下一个可用任务 | `claim.go` |
| `forge task submit` | 提交任务执行结果 | `submit.go` |
| `forge task status <task-id>` | 查询任务状态（只读） | `status.go` |
| `forge task query` | 查询任务信息 | `query.go` |
| `forge task check-deps` | 检查任务依赖关系 | `check_deps.go` |
| `forge task validate-index [file]` | 验证 index.json 文件 | `validate_index.go` |
| `forge task add` | 向当前 feature 添加新任务 | `add.go` |
| `forge task index` | 从任务 markdown 文件构建或重建 index.json | `index.go` |
| `forge task migrate` | 为所有任务推断 type 字段并迁移 index.json | `migrate.go` |
| `forge task list-types` | 列出所有支持的任务类型 | `list_types.go` |
| `forge task reopen <task-id>` | 重新打开已拒绝/跳过的任务（恢复为 pending） | `reopen.go` |
| `forge task transition <task-id> <status> --reason` | 手动切换任务状态（操作员覆盖） | `transition.go` |

### forge e2e — 端到端测试管理

源文件：`e2e_parent.go`

| 命令 | 用途 | 源文件 |
|------|------|--------|
| `forge e2e validate-specs` | 验证生成的 spec 文件是否符合结构规则 | `e2e_validate_specs.go` |
| `forge e2e run` | 运行 e2e 测试 | `e2e_run.go` |
| `forge e2e setup` | 安装 e2e 依赖（幂等） | `e2e_setup.go` |
| `forge e2e verify` | 检查未解决的 VERIFY 标记 | `e2e_verify.go` |
| `forge e2e compile` | 编译检查 e2e 测试文件 | `e2e_compile.go` |
| `forge e2e discover` | 列出所有 e2e 测试用例（不运行） | `e2e_discover.go` |

### forge test — 测试工具

源文件：`test.go`

| 命令 | 用途 | 源文件 |
|------|------|--------|
| `forge test promote <journey-name>` | 将 journey 的 @feature 标签提升为 @regression | `test_promote.go` |
| `forge test run-journey <journey-name>` | 在隔离的临时目录中运行单个 journey | `journey_isolation.go` |
| `forge test verify` | 对比 Contract spec 与当前代码，检测契约破坏 | `test_verify.go` |

### forge prompt — Agent 执行提示词管理

源文件：`prompt/register.go`, `prompt/prompt_get.go`

| 命令 | 用途 | 源文件 |
|------|------|--------|
| `forge prompt get-by-task-id <id>` | 为任务合成 agent 提示词 | `prompt/prompt_get.go` |

### forge worktree — Git Worktree 管理

源文件：`worktree.go`

| 命令 | 用途 | 源文件 |
|------|------|--------|
| `forge worktree start [slug]` | 创建 worktree 并在其中启动 Claude（slug 可选，支持 `-i` 交互选择） | `worktree.go` |
| `forge worktree list` | 列出所有 git worktree | `worktree.go` |
| `forge worktree remove <slug>` | 移除 worktree（`--hard` 删除分支，`--force` 强制） | `worktree.go` |
| `forge worktree resume <slug>` | 在已有 worktree 中重新启动 Claude | `worktree.go` |
| `forge worktree push` | 推送当前 worktree 分支到远程并设置 upstream 跟踪 | `worktree.go` |
| `forge worktree status [<slug>]` | 显示 worktree 状态（分支、提交、未提交文件列表）；无参数时显示所有 worktree | `worktree.go` |

### forge config — 配置管理

源文件：`config.go`

| 命令 | 用途 | 源文件 |
|------|------|--------|
| `forge config get <key>` | 获取配置值（纯文本输出） | `config.go` |
| `forge config set <key> <value>` | 设置配置值（支持 dot-notation 嵌套键） | `config.go` |
| `forge config init` | 交互式初始化 .forge/config.yaml | `config.go` |

### forge feature — Feature 管理

源文件：`feature.go`

| 命令 | 用途 | 源文件 |
|------|------|--------|
| `forge feature list` | 列出所有 feature | `feature.go` |
| `forge feature status <slug>` | 显示 feature 状态详情 | `feature.go` |
| `forge feature set <slug>` | 显式设置当前 feature | `feature.go` |
| `forge feature complete` | 所有任务完成后执行 feature 收尾 | `feature_complete.go` |

### forge forensic — 会话取证分析

源文件：`forensic/register.go`, `forensic/forensic.go`

| 命令 | 用途 | 源文件 |
|------|------|--------|
| `forge forensic search [project-path]` | 在 history.jsonl 中搜索匹配的会话 | `forensic/forensic.go` |
| `forge forensic extract <session-jsonl-path>` | 从会话记录中提取紧凑证据 | `forensic/forensic.go` |
| `forge forensic subagents <session-dir-path>` | 列出会话的子 agent 记录 | `forensic/forensic.go` |

## CLI 表格渲染约定

列表命令（`forge feature list`、`forge lesson`、`forge proposal`）的 slug/name 列使用动态宽度：

- **宽度计算**: `clamp(max(30, maxSlugLen + 2), 60)` — 最小 30 字符，最大 60 字符
- **实现**: `calcSlugColWidth()` 计算宽度，`padRight()` 对齐，`truncateSlug()` 截断超长值
- **新增列表命令时**必须遵循此模式（常量 `slugColMinWidth` / `slugColMaxWidth` 在 `proposal.go` 中定义）

## 列表排序约定

所有列表命令按 frontmatter `created` 字段（`YYYY-MM-DD` 格式）降序排列。缺少 `created` 时 fallback 到文件 mtime。git clone 会重置 mtime，因此 `created` frontmatter 是唯一可靠的创建时间源。

## 已移除的命令（禁止使用）

以下命令已从 CLI 中移除，**不得在任何 skill 或文档中引用**：

| 已移除命令 | 原所属组 | 说明 |
|-----------|---------|------|
| `forge test detect` | test | 已移除。通过读取项目文件（package.json / go.mod / Cargo.toml / pyproject.toml）推断测试语言 |
| `forge test interfaces` | test | 已移除。通过读取项目结构和 `docs/conventions/` 推断接口类型 |
| `forge test framework` | test | 已移除 |
| `forge test get` | test | 已移除 |

> **维护提醒**：如果在 skill 编写中需要检测项目信息，请直接读取项目文件，不要依赖 CLI 命令。
