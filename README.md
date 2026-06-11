# Forge

> **Spec Driven Development** 工具链 — 让 Claude Code 从"聊天"变成"工程"

[![Version](https://img.shields.io/badge/Version-3.0.0-blue.svg)](https://github.com/bigfaner/forge)
[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

你一定经历过这些场景：

- 聊了两小时，AI 把需求理解偏了，代码全白写
- "帮我加个按钮"，三小时后整个模块被重构了
- 昨天修的 bug，今天 AI 又改回来，因为它根本不记得
- 每次开新会话都要重新教 AI 项目规范和决策上下文

**Forge** 是为 Claude Code 打造的 **Spec Driven Development（SDD）** 工具链，核心理念是 **Harness Engineering**——将 AI 编程从自由对话变为受控的工程流水线。它用 `brainstorm → PRD → 设计 → 任务 → 自动执行` 的工程化流程，把 AI 编程从"凭感觉聊天"变成"按规范交付"。

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

### 结构化 Pipeline

两种模式覆盖不同规模的需求：

| 模式 | 流程 | 适用场景 |
|------|------|----------|
| 完整模式 | brainstorm → PRD → 技术设计 → 任务拆分 → 自动执行 | 复杂特性（>10 个任务） |
| 快速模式 | brainstorm → 任务直接生成 → 自动执行 | 小特性（1-10 个任务） |

### 对抗式评估体系

8 种专属评估器（PRD、技术设计、UI 设计、Proposal、Journey、Contract 等），采用**专家角色 + 1000+ 分制评分表**进行多轮迭代修订，支持跨文档一致性校验（PRD → 设计 → 任务对齐），从源头保障文档质量。

### 自主任务执行引擎

`/run-tasks` 自动循环分发任务至 Subagent，每个 Subagent 严格遵循 **TDD 协议**（RED → GREEN → REFACTOR）。支持**动态修复链**：失败任务自动创建修复任务，修复完成后恢复原任务。7 种任务状态机管控（pending → in_progress → completed / blocked / suspended 等），每步都有 Quality Gate（`compile → fmt → lint → test`）自动拦截。

### Journey-Contract 测试模型

从 PRD 用户故事提取 **Journey**（用户流程 + 风险分级），从 Journey 生成六维度 **Contract**（行为契约），从 Contract 生成可执行的 **Surface 感知测试脚本**（web / api / cli / tui / mobile）。

### 上下文持久化

`manifest.md` 追踪 feature 全链路状态，`forge worktree` 隔离并行开发环境，`index.json` 记录任务依赖和进度。AI 不再"失忆"——每个会话都能从上一次的精确断点继续。

### 知识沉淀

`/learn` 把项目中的设计决策、经验教训、技术规范沉淀为可复用的知识。`/consolidate-specs` 自动检测规范与代码的漂移并修复。新的会话、新的贡献者都能站在前人的肩膀上，而不是每次从零开始。

---

## 工程规模

**21 个 Skill** · **16 个 Slash Command** · **1 个 Subagent** · **20 种任务类型** · **Go CLI**（19 个命令） + **Claude Code Plugin** · Hooks 系统实现会话级自动上下文注入与清理

---

## 5 分钟体验

```bash
# 初始化项目
forge init

# 快速模式 — 流程更短，跳过 PRD/设计/eval
/quick

# 完整模式 — 完整的 brainstorm → PRD → 设计 → 任务 → 执行 流程
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

## 深入了解

- [架构概览](docs/user-guide/architecture-overview.md) — 插件机制、四大组件、数据流与状态管理
- [使用指南](docs/user-guide/usage-guide.md) — Full Mode / Quick Mode 端到端实战、单命令场景、常见问题排错
- [项目初始化](docs/user-guide/initialization.md) — `forge init` 完整流程、配置字段参考、Surface 检测
- [环境配置](docs/user-guide/environment-setup.md) — 从零开始配置 Forge 开发环境

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
