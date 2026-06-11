# Forge

> 结构化 AI 编程工具链 — 让 Claude Code 从"聊天"变成"工程"

[![Version](https://img.shields.io/badge/Version-3.0.0-blue.svg)](https://github.com/bigfaner/forge)
[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

你一定经历过这些场景：

- 聊了两小时，AI 把需求理解偏了，代码全白写
- "帮我加个按钮"，三小时后整个模块被重构了
- 昨天修的 bug，今天 AI 又改回来，因为它根本不记得
- 每次开新会话都要重新教 AI 项目规范和决策上下文

**Forge** 是为 Claude Code 打造的结构化工作流工具链。它用 `brainstorm → PRD → 设计 → 任务 → 自动执行` 的工程化流程，把 AI 编程从"凭感觉聊天"变成"按规范交付"。

---

## 竞品对比

| 维度 | Forge | Superpowers | Spec Kit | OpenSpec |
|------|:-----:|:-----------:|:--------:|:--------:|
| 结构化流程 | ✓ | ✓ | ✓ | ✓ |
| 质量门控 (Quality Gate) | ✓ | ✗ | ✗ | ✗ |
| 上下文持久化 (manifest + worktree) | ✓ | ✗ | ✗ | ✗ |
| 知识沉淀 (`/learn`) | ✓ | ✗ | ✗ | ✗ |
| 跨会话连续性 | ✓ | ✗ | ✗ | ✗ |
| Agent 自动编排 (`/run-tasks`) | ✓ | ✓ | ✗ | ✗ |
| 多 Agent / Subagent | ✗ | ✓ | ✗ | ✗ |
| 跨 IDE / Agent 平台 | ✗ | ✓ | ✓ | ✓ |

> 数据来源：截至 2026-06，基于各项目 GitHub README 及文档的实际功能验证。

---

## 核心特性

### 质量门控

每一步都有自动化的质量检查。`compile → fmt → lint → test` 四层门控确保 AI 产出的代码不会"看起来对、跑起来崩"。Quality Gate 在任务提交和阶段退出时自动执行，不合格就打回。

### 上下文持久化

`manifest.md` 追踪 feature 全链路状态，`forge worktree` 隔离并行开发环境，`index.json` 记录任务依赖和进度。AI 不再"失忆"——每个会话都能从上一次的精确断点继续。

### 知识沉淀

`/learn` 把项目中的设计决策、经验教训、技术规范沉淀为可复用的知识。新的会话、新的贡献者都能站在前人的肩膀上，而不是每次从零开始。

### Agent 自动编排

`/run-tasks` 自动认领、分发、执行任务，每个任务由独立的 task-executor agent 以 TDD 方式完成，产出可追溯的执行记录。复杂功能从拆解到交付，无需人工干预每个步骤。

---

## 5 分钟体验

```bash
# 初始化项目
forge init

# 快速模式 — 适合小功能（1-2h，1-10 任务）
/quick

# 完整模式 — 适合复杂功能（>2h，>10 任务）
/brainstorm -> /write-prd -> /tech-design -> /breakdown-tasks -> /run-tasks
```

---

## 安装

### 前置要求

- [Claude Code](https://docs.anthropic.com/en/docs/claude-code) CLI
- curl（macOS/Linux 自带，Windows 10+ 自带）

### 安装步骤

**macOS / Linux：**

```bash
# 1. 安装 forge CLI
curl -fsSL https://github.com/bigfaner/forge/releases/latest/download/install.sh | bash

# 2. 安装 forge Plugin（CLI binary + Plugin 一步到位）
forge upgrade

# 3. 在项目中初始化
cd my-project && forge init
```

**Windows (PowerShell)：**

```powershell
# 1. 安装 forge CLI
irm https://github.com/bigfaner/forge/releases/latest/download/install.ps1 | iex

# 2. 安装 forge Plugin（CLI binary + Plugin 一步到位）
forge upgrade

# 3. 在项目中初始化
cd my-project; forge init
```

本地开发者构建：`git clone` -> `cd forge-cli && bash scripts/install-local.sh` -> `forge upgrade`

---

## 贡献

```bash
git clone git@github.com:bigfaner/forge.git && cd forge
cd forge-cli && go mod download
go test -race -cover ./...
```

提交遵循 [Conventional Commits](https://www.conventionalcommits.org/)。

## License

[MIT](LICENSE)
