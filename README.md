# ZCode

> Claude Code 工作流增强工具集：让 AI 编程从"聊天"变成"工程"

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

---

## 为什么选择 ZCode？

### 问题

用 Claude Code 开发时，你是否遇到过：

| 问题 | 表现 |
|------|------|
| **方向漂移** | 做着做着发现偏离了原始需求 |
| **重复错误** | 同样的坑踩多次，没有知识沉淀 |
| **质量失控** | AI 生成的代码缺乏测试覆盖 |
| **上下文丢失** | 长会话中 Claude 忘记之前的决策 |
| **协作困难** | 多人开发时任务状态不透明 |

### ZCode 的解法

```
┌─────────────────────────────────────────────────────────────────┐
│  结构化工作流 + 任务管理 + 知识沉淀 = 可预测的交付质量           │
└─────────────────────────────────────────────────────────────────┘
```

| ZCode 提供 | 解决的问题 |
|-----------|-----------|
| **PRD → 设计 → 任务** 流程 | 防止方向漂移 |
| **task-cli** 任务声明与追踪 | 上下文不丢失 |
| **强制 record-task** | 知识沉淀，可回溯 |
| **TDD workflow** | 质量保证（覆盖率 80%+） |
| **Git 分支自动识别** | 多人协作任务隔离 |

### 与其他方案对比

| 方案 | 优点 | 缺点 |
|------|------|------|
| **直接用 Claude Code** | 灵活、快速 | 无结构、难追溯、质量不稳定 |
| **GitHub Projects/Jira** | 成熟的任务管理 | 与 AI 开发脱节，需要人工同步 |
| **Cursor Rules** | 轻量级指导 | 仅有规则，无执行流程 |
| **ZCode** | AI 原生工作流 + 任务管理 + 知识沉淀 | 需要 Go 环境，学习成本 |

---

## 核心概念

```
/write-prd → /design-tech → /breakdown-tasks → /run-tasks
     ↓              ↓              ↓               ↓
   prd.md      design.md      index.json       自动执行
```

| 阶段 | 命令 | 产出物 | 说明 |
|------|------|--------|------|
| 需求 | `/write-prd` | `prd.md` | 协作式澄清需求，防止 XY 问题 |
| 设计 | `/design-tech` | `design.md` | 技术决策，避免实现时返工 |
| 拆分 | `/breakdown-tasks` | `index.json` | 任务粒度 1-4 小时 |
| 执行 | `/run-tasks` | 代码 + 记录 | 自动 TDD 循环 |

---

## 快速开始

### 前提条件

- [Go 1.26.1+](https://golang.org/dl/)
- [Claude Code CLI](https://github.com/anthropics/claude-code)
- Git

### 安装

#### 方式一：通过 Plugin Marketplace 安装（推荐）

ZCode 可作为 Claude Code plugin marketplace 使用，支持标准化安装和自动更新：

```bash
# 1. 添加 marketplace
/plugin marketplace git@github.com:bigfaner/zcode.git

# 2. 安装 zcode plugin（--scope project 与团队共享）
/plugin install zcode@zcode-marketplace --scope project

# 3. 编译安装 task-cli
/init-zcode

# 4. 验证
task --version
```

**安装范围说明**：

| 范围 | 设置文件 | 用例 |
|------|----------|------|
| `user` | `~/.claude/settings.json` | 个人使用（默认） |
| `project` | `.claude/settings.json` | 团队共享（推荐） |
| `local` | `.claude/settings.local.json` | 本地开发，gitignored |

#### 方式二：本地安装

```bash
# 1. 克隆仓库
git clone git@github.com:bigfaner/zcode.git
cd zcode

# 2. 添加本地 marketplace 并安装
/plugin marketplace add .
/plugin install zcode@zcode-marketplace

# 3. 编译安装 task-cli
/init-zcode

# 4. 验证
task --version
```

### 5 分钟体验

```bash
# 1. 创建 feature 目录
mkdir -p docs/features/demo/{tasks,records}

# 2. 开始写 PRD（Claude 会引导你）
/write-prd

# 3. 技术设计
/design-tech

# 4. 拆分任务
/breakdown-tasks

# 5. 自动执行
/run-tasks
```

---

## 详细指南

### 目录结构

```
project-root/
├── docs/
│   ├── features/<slug>/           # Feature 工作区
│   │   ├── prd.md                 # 需求文档
│   │   ├── design.md              # 设计文档
│   │   └── tasks/
│   │       ├── index.json         # 任务定义（核心）
│   │       ├── process/           # 运行时状态（不提交）
│   │       │   ├── state.json     # 当前任务状态
│   │       │   └── record.json    # 进行中的记录
│   │       ├── 1.1-<title>.md     # 任务详情
│   │       └── records/               # 执行记录
│   │           └── 1.1-<title>.md
│   └── lessons/                   # 经验教训
├── task-cli/                      # CLI 工具源码
└── plugins/zcode/                 # 插件定义
```

### task-cli 命令

#### 核心操作

```bash
task claim              # 声明下一个任务（自动排序）
task record 1.1 -data record.json   # 记录任务完成
task status 1.1         # 查询状态
task status 1.1 completed           # 更新状态
task query 1.1          # 查询详情
task feature auth       # 切换/设置 feature
```

#### 验证与维护

```bash
task validate           # 验证 index.json
task check              # 检查依赖完整性
task verifyCompletion   # 验证任务完成（hook 用）
task cleanup            # 清理状态文件（hook 用）
```

### 任务状态流转

```
pending → in_progress → completed
                 ↓
              blocked → in_progress
                 ↓
              skipped
```

### 依赖语法

| 语法 | 示例 | 含义 |
|------|------|------|
| 精确 | `1.1`, `2.3` | 特定任务完成后才可开始 |
| 通配 | `1.x` | Phase 1 所有任务完成后才可开始 |

### Feature 自动识别

优先级从高到低：

1. **Git Worktree 名称** — `worktrees/auth-login` → `auth-login`
2. **Git 分支名称** — `feature/auth-login` → `auth-login`
3. **State 文件存在** — 有 `tasks/process/state.json` 的 feature
4. **唯一 feature** — 只有一个 feature 有 `index.json`

分支前缀自动剥离：`feature/`, `feat/`, `fix/`, `bugfix/`, `hotfix/`, `chore/`

---

## Skills 参考

### 规划阶段

| Skill | 用途 | 何时使用 |
|-------|------|----------|
| `/write-prd` | 产出 PRD | 用户描述需求但细节不明 |
| `/design-tech` | 产出设计文档 | PRD 已完成，需要技术方案 |
| `/breakdown-tasks` | 产出任务列表 | 设计完成，需要拆分执行 |

### 执行阶段

| Skill | 用途 | 何时使用 |
|-------|------|----------|
| `/claim-task` | 声明任务 | 手动领取下一个任务 |
| `/execute-task` | 执行单任务 | 手动执行一个任务 |
| `/run-tasks` | 自动循环执行 | 自动化批量执行 |
| `/record-task` | 记录完成 | 任务完成后必须调用 |
| `/set-task-status` | 更新状态 | 手动修改任务状态 |

### 辅助工具

| Skill | 用途 |
|-------|------|
| `/git-commit` | Conventional Commits 规范提交 |
| `/learn-lesson` | 提取可复用知识到 `docs/lessons/` |
| `/simplify-skill` | 重构 skill 文件 |
| `/init-zcode` | 编译安装 task-cli |

---

## Agents

| Agent | 触发方式 | 职责 |
|-------|----------|------|
| **task-executor** | `/run-tasks` 分发 | 执行单个任务（5 步 TDD） |
| **error-fixer** | 执行失败时分发 | 修复编译/测试错误 |

### task-executor 工作流

```
Step 1: Read task definition
Step 2: TDD (RED → GREEN → REFACTOR)
Step 3: Full verification (build + test + coverage >= 80%)
Step 4: Record task (MANDATORY)
Step 5: Git commit
```

**铁律**：
- 无 record-task = 任务未完成
- 验证不通过 = 不提交
- 最多 3 次 subagent 调用

---

## 数据模型

### index.json 示例

```json
{
  "feature": "auth-login",
  "title": "用户登录功能",
  "tasks": {
    "1.1": {
      "id": "1.1",
      "title": "定义认证接口",
      "description": "定义 AuthService 接口",
      "phase": 1,
      "priority": "P0",
      "status": "pending",
      "dependencies": [],
      "files": ["internal/auth/service.go"]
    },
    "1.2": {
      "id": "1.2",
      "title": "实现 JWT 签发",
      "phase": 1,
      "priority": "P0",
      "status": "pending",
      "dependencies": ["1.1"],
      "files": ["internal/auth/jwt.go"]
    }
  }
}
```

### record.json 示例

```json
{
  "taskId":"3.3.1",
  "status": "completed",
  "summary": "实现 AuthService 接口定义",
  "filesCreated": ["internal/auth/service.go"],
  "filesModified": [],
  "keyDecisions": ["使用 interface{} 支持多种认证方式"],
  "testsPassed": 5,
  "testsFailed": 0,
  "coverage": 92.5,
  "acceptanceCriteria": [
    { "criterion": "接口定义完整", "met": true },
    { "criterion": "单元测试覆盖", "met": true }
  ]
}
```

---

## Hooks

| 事件 | 触发 | 作用 |
|------|------|------|
| `PostToolUse` (Edit/Write `index.json`) | `validate-index.sh` | 自动验证任务定义 |
| `SessionEnd` | `task cleanup` | 清理已完成任务状态 |
| `SubagentStop` | `task cleanup` | 子代理结束时清理 |

---

## 故障排除

### 常见问题

| 问题 | 原因 | 解决方案 |
|------|------|----------|
| `task: command not found` | 未安装或 PATH 未刷新 | 重新打开终端或 `source ~/.bashrc` |
| `feature not found` | 目录不存在或未设置 | `task feature <slug>` 手动设置 |
| `no available task` | 所有任务已完成或依赖未满足 | 检查 `task check` 输出 |
| `index.json 验证失败` | JSON 语法错误或依赖引用无效 | 检查 JSON 格式和任务 ID |
| `verifyCompletion 失败` | 任务未完成或缺少 record 文件 | 确保 `/record-task` 已执行 |

### 调试命令

```bash
# 查看当前 feature
task feature

# 查看所有任务状态
cat docs/features/<slug>/tasks/index.json | grep status

# 查看运行时状态
cat docs/features/<slug>/tasks/process/state.json

# 强制重置状态
rm docs/features/<slug>/tasks/process/state.json
```

---

## Plugin Marketplace 架构

ZCode 遵循 [Claude Code plugin 规范](docs/official-references/plugin.md)，可作为 marketplace 分发。

### 目录结构

```
zcode/                              # Marketplace 根目录
├── .claude-plugin/
│   └── marketplace.json            # Marketplace 清单
├── plugins/
│   └── zcode/                      # ZCode plugin
│       ├── .claude-plugin/
│       │   └── plugin.json         # Plugin 清单
│       ├── skills/                 # Skills（命令）
│       ├── agents/                 # Subagents
│       └── hooks/                  # Hooks 配置
└── task-cli/                       # CLI 工具源码
```

### marketplace.json

```json
{
  "name": "zcode",
  "owner": { "name": "bigfaner" },
  "plugins": [
    {
      "name": "zcode",
      "source": "./plugins/zcode",
      "description": "Claude Code 工作流增强工具集"
    }
  ]
}
```

### Plugin 组件

ZCode plugin 包含以下组件：

| 组件 | 目录 | 用途 |
|------|------|------|
| **Skills** | `skills/` | `/write-prd`、`/run-tasks` 等命令 |
| **Agents** | `agents/` | `task-executor`、`error-fixer` 子代理 |
| **Hooks** | `hooks/hooks.json` | 自动验证和清理 |

### 相关 CLI 命令

```bash
# Marketplace 管理
/plugin marketplace add <owner/repo>   # 添加 marketplace
/plugin marketplace list               # 列出已添加的 marketplaces
/plugin marketplace update zcode       # 更新 marketplace

# Plugin 管理
/plugin install zcode@zcode            # 安装 plugin
/plugin update zcode@zcode             # 更新到最新版本
/plugin uninstall zcode@zcode          # 卸载
/plugin enable zcode@zcode             # 启用
/plugin disable zcode@zcode            # 禁用

# 验证
/plugin validate .                     # 验证本地 marketplace
```

---

## 架构约束（团队内部）

```
依赖方向: cmd → internal → pkg (严禁反向)
模块交互: 通过接口/类型定义，不直接依赖内部实现
测试覆盖: >= 80%
TDD 流程: RED → GREEN → REFACTOR
```

---

## 贡献指南

### 开发环境

```bash
# 克隆
git clone git@github.com:bigfaner/zcode.git
cd zcode

# 安装依赖
cd task-cli && go mod download

# 运行测试
go test -race -cover ./...

# Lint
golangci-lint run ./...
```

### 提交规范

使用 Conventional Commits：

```
feat(cli): add parallel task execution
fix(claim): handle circular dependencies
docs(readme): add troubleshooting section
```

### PR 检查清单

- [ ] 测试通过 (`go test -race -cover ./...`)
- [ ] 覆盖率 >= 80%
- [ ] Lint 通过 (`golangci-lint run`)
- [ ] 文档已更新（如适用）
- [ ] 提交消息符合规范

---

## 文档索引

| 文档 | 说明 |
|------|------|
| [task-cli/docs/OVERVIEW.md](task-cli/docs/OVERVIEW.md) | CLI 功能详解 |
| [task-cli/docs/WORKFLOW.md](task-cli/docs/WORKFLOW.md) | 内部流程图解 |
| [docs/official-references/plugin.md](docs/official-references/plugin.md) | Claude Code 插件系统完整技术参考 |
| [docs/official-references/plugin-marketplace.md](docs/official-references/plugin-marketplace.md) | Plugin marketplace 创建与分发指南 |

---

## License

[MIT](LICENSE)
