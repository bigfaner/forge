# Forge

> Claude Code 工作流增强工具集：让 AI 编程从"聊天"变成"工程"

[![Version](https://img.shields.io/badge/Version-2.16.1-blue.svg)](https://github.com/bigfaner/forge)
[![Go Version](https://img.shields.io/badge/Go-1.26.1+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

---

## Forge 解决什么？

| 痛点 | Forge 的解法 |
|------|-------------|
| 方向漂移 | `brainstorm → PRD → 设计 → 任务` 结构化流程 |
| 质量失控 | Quality Gate（compile → fmt → lint → test）+ TDD 工作流 |
| 上下文丢失 | task-cli 任务追踪 + manifest.md 全链路追溯 |
| 知识不沉淀 | `/learn-lesson` + `/record-decision` 跨会话积累 |

---

## 两种工作模式

### 完整模式（复杂功能：>2h, >10 任务）

```
/brainstorm → /write-prd → /tech-design → /breakdown-tasks → /run-tasks
     ↓             ↓            ↓ ↘ /ui-design      ↓              ↓
 proposal.md   prd/*.{3}   design/*.{2}  ui/    tasks + index.json  自动执行
```

每阶段产出文档，可选通过 `/eval-*` 系列技能迭代评分至达标（默认可选，非强制）。

### 快速模式（小功能：1-2h, 1-10 任务）

```
/quick → /brainstorm → /quick-tasks → /run-tasks
```

跳过 PRD 和设计，`proposal.md` 直接驱动任务。纯文档 feature 自动跳过测试，生成文档评估任务。

---

## 快速开始

### 安装

```bash
# Marketplace 安装（推荐）
/plugin marketplace add git@github.com:bigfaner/forge.git
/plugin install forge@forge --scope project
/init-forge
task --version
```

或本地安装：`git clone` → `/plugin marketplace add .` → `/plugin install forge@forge` → `/init-forge`

### 5 分钟体验

```bash
# 快速模式
/quick

# 完整模式
/brainstorm → /write-prd → /tech-design → /ui-design → /breakdown-tasks → /run-tasks
```

---

## Skills 一览

### 规划

| Skill | 产出 |
|-------|------|
| `/brainstorm` | 结构化提案 `proposal.md` |
| `/write-prd` | PRD 三件套 + `manifest.md` |
| `/tech-design` | 技术设计文档 |
| `/ui-design` | UI 规格 + 可选 HTML 原型（5 种内置风格） |
| `/breakdown-tasks` | 任务文件 + `index.json` + `manifest.md` |

### 快速模式

| Skill | 产出 |
|-------|------|
| `/quick` | 启动快速模式流程 |
| `/quick-tasks` | 从提案直接生成任务 |

### 评估（1000 分制，对抗式迭代至达标；`/eval-harness` 例外，使用 100 分制）

`/eval-prd` · `/eval-design` · `/eval-ui` · `/eval-proposal` · `/eval-test-cases` · `/eval-consistency` · `/eval-harness`

### 测试生命周期

`/gen-sitemap` → `/gen-test-cases` → `/eval-test-cases` → `/gen-test-scripts` → `/run-e2e-tests` → `/graduate-tests` → `verify-regression` → `/consolidate-specs`

### 执行

| Skill | 用途 |
|-------|------|
| `/execute-task` | 执行单任务 |
| `/run-tasks` | 自动循环分发 |
| `/submit-task` | 记录完成（必须） |

### 辅助

`/fix-bug` · `/git-commit` · `/git-checkout` · `/learn-lesson` · `/record-decision` · `/consolidate-specs` · `/init-justfile` · `/init-forge` · `/gen-sitemap` · `/extract-design-md` · `/simplify-skill` · `/forensic` · `/improve-harness`

---

## Agents

| Agent | 职责 |
|-------|------|
| **task-executor** | 执行单个任务（TDD + Quality Gate + record） |
| **doc-scorer** | 按评分标准对文档打分 |
| **doc-reviser** | 根据评分报告修订文档 |

---

## 项目结构

```
forge/
├── plugins/forge/          # Forge plugin
│   ├── skills/             # 23 个 Skills
│   ├── commands/           # 11 个 Slash Commands
│   └── agents/             # 3 个 Subagents
├── task-cli/               # Go CLI 工具源码
├── tests/e2e/              # Playwright E2E 回归测试
├── docs/                   # 项目文档
└── web/                    # Web 看板（Vite + React）
```

```
docs/features/<slug>/
├── manifest.md             # Feature 单一入口（自动维护）
├── prd/                    # /write-prd 产出
├── design/                 # /tech-design 产出
├── ui/                     # /ui-design 产出（可选）
├── testing/                # /gen-test-cases 产出
└── tasks/
    ├── index.json          # 任务定义
    ├── *.md                # 任务详情
    └── records/            # 执行记录
```

---

## 贡献

```bash
git clone git@github.com:bigfaner/forge.git && cd forge
cd task-cli && go mod download
go test -race -cover ./...
```

提交遵循 [Conventional Commits](https://www.conventionalcommits.org/)，测试覆盖率 >= 80%。

---

## 文档索引

| 文档 | 说明 |
|------|------|
| [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) | 核心架构、工作流管道、Agent 协作、Quality Gate |
| [task-cli/docs/OVERVIEW.md](task-cli/docs/OVERVIEW.md) | CLI 完整命令参考 |
| [task-cli/docs/WORKFLOW.md](task-cli/docs/WORKFLOW.md) | 内部流程图解 |
| [docs/official-references/plugin.md](docs/official-references/plugin.md) | 插件系统技术参考 |
| [docs/official-references/plugin-marketplace.md](docs/official-references/plugin-marketplace.md) | Marketplace 分发指南 |
| [docs/official-references/hooks.md](docs/official-references/hooks.md) | Hooks 技术参考 |

---

## License

[MIT](LICENSE)
