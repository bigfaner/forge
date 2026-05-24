# Forge

> Claude Code 工作流增强工具集：让 AI 编程从"聊天"变成"工程"

[![Version](https://img.shields.io/badge/Version-5.6.0-blue.svg)](https://github.com/bigfaner/forge)
[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

---

## Forge 解决什么？

| 痛点 | Forge 的解法 |
|------|-------------|
| 方向漂移 | `brainstorm -> PRD -> 设计 -> 任务` 结构化流程 |
| 质量失控 | Quality Gate（compile -> fmt -> lint -> test）+ TDD 工作流 |
| 上下文丢失 | forge CLI 任务追踪 + manifest.md 全链路追溯 |
| 知识不沉淀 | `/learn` 跨会话积累决策与经验 |

---

## 两种工作模式

### 完整模式（复杂功能：>2h, >10 任务）

```
/brainstorm -> /write-prd -> /tech-design -> /breakdown-tasks -> /run-tasks
     |             |            |  -> /ui-design      |              |
 proposal.md   prd/*.{3}   design/*.{2}  ui/    tasks + index.json  自动执行
```

每阶段产出文档，可选通过 `/eval-*` 系列技能迭代评分至达标。

### 快速模式（小功能：1-2h, 1-10 任务）

```
/quick -> /brainstorm -> /quick-tasks -> /run-tasks
```

跳过 PRD 和设计，`proposal.md` 直接驱动任务。纯文档 feature 自动跳过测试，生成文档评估任务。

---

## 安装

### 前置要求

- [Go 1.25+](https://golang.org/dl/)
- [Claude Code](https://docs.anthropic.com/en/docs/claude-code) CLI

### 安装步骤

```bash
# Marketplace 安装（推荐）
/plugin marketplace add git@github.com:bigfaner/forge.git
/plugin install forge@forge --scope project
/init-forge
forge --help
```

或本地安装：`git clone` -> `/plugin marketplace add .` -> `/plugin install forge@forge` -> `/init-forge`

### 5 分钟体验

```bash
# 快速模式
/quick

# 完整模式
/brainstorm -> /write-prd -> /tech-design -> /ui-design -> /breakdown-tasks -> /run-tasks
```

---

## 命令速查

> 与 `forge --help` 一一对应

| 命令 | 用途 |
|------|------|
| `forge init` | 初始化 Forge 项目环境 |
| `forge config` | 管理项目配置（.forge/config.yaml） |
| `forge surfaces` | 查询项目 surfaces 配置 |
| `forge feature` | 设置或显示当前 feature 上下文 |
| `forge task` | 任务生命周期管理 |
| `forge test` | 测试工具集（promote / run-journey / verify） |
| `forge prompt` | 生成 agent 执行提示词 |
| `forge quality-gate` | 检查所有任务完成，运行回归测试 |
| `forge fact` | 管理结构化系统事实表 |
| `forge worktree` | 管理 git worktree 并行开发 |
| `forge forensic` | 分析会话记录，诊断 agent 偏差 |
| `forge verify-task-done` | 提交前验证任务完成状态 |
| `forge cleanup` | 清理已完成任务的状态文件 |
| `forge lesson` | 列出或查看经验详情 |
| `forge proposal` | 列出或查看提案详情 |
| `forge research` | 列出或查看研究报告详情 |
| `forge claude` | 跳过权限检查启动 Claude CLI |
| `forge completion` | 生成指定 shell 的自动补全脚本 |
| `forge help` | 查看任意命令的帮助信息 |

### 常用 task 子命令

| 子命令 | 用途 |
|--------|------|
| `forge task add` | 创建新任务 |
| `forge task list` | 列出当前 feature 所有任务 |
| `forge task claim` | 认领下一个可用任务 |
| `forge task submit` | 提交任务执行结果 |
| `forge task status` | 查询任务状态 |
| `forge task index` | 从任务 Markdown 重建 index.json |
| `forge task list-types` | 列出所有支持的任务类型 |
| `forge task query` | 查询任务信息 |
| `forge task migrate` | 推断并补充 index.json 中的 type 字段 |

---

## Skills 一览（21 个）

> 计数与 `ls plugins/forge/skills/ | wc -l` 一致

### 规划

| Skill | 产出 |
|-------|------|
| `/brainstorm` | 结构化提案 `proposal.md` |
| `/write-prd` | PRD 三件套 + `manifest.md` |
| `/tech-design` | 技术设计文档 |
| `/ui-design` | UI 规格 + 可选 HTML 原型 |
| `/breakdown-tasks` | 任务文件 + `index.json` + `manifest.md` |

### 快速模式

| Skill | 产出 |
|-------|------|
| `/quick` | 启动快速模式流程 |
| `/quick-tasks` | 从提案直接生成任务 |

### 评估（1000 分制，对抗式迭代至达标）

`/eval-prd` / `/eval-design` / `/eval-ui` / `/eval-proposal` / `/eval-journey` / `/eval-contract` / `/eval-consistency` / `/eval`（通用评估，参数化 rubric）

### 测试生命周期

`/gen-sitemap` -> `/gen-journeys` -> `/eval-journey` -> `/gen-contracts` -> `/eval-contract` -> `/gen-test-scripts` -> `/run-tests` -> `forge test promote` -> `/consolidate-specs`

### 执行

| Skill | 用途 |
|-------|------|
| `/execute-task` | 执行单任务（TDD + Quality Gate + record） |
| `/run-tasks` | 自动循环分发 |
| `/submit-task` | 记录完成（必须） |

### 辅助

`/fix-bug` / `/git-commit` / `/git-checkout` / `/learn` / `/consolidate-specs` / `/init-forge` / `/gen-sitemap` / `/extract-design-md` / `/forensic` / `/deep-research` / `/clean-code` / `/test-guide` / `/simplify-skill`

---

## Agents

| Agent | 职责 |
|-------|------|
| **task-executor** | 执行单个任务（TDD + Quality Gate + record） |

---

## 任务类型表（21 种 dot-notation 类型）

> 与 `forge task list-types` 一一对应

### coding（5 种）

| 类型 | 用途 |
|------|------|
| `coding.feature` | 实现新运行时行为 |
| `coding.enhancement` | 增强已有行为 |
| `coding.cleanup` | 移除死代码或修复技术债 |
| `coding.refactor` | 无行为变更的重构 |
| `coding.fix` | 修复 bug 或问题 |

### doc（5 种）

| 类型 | 用途 |
|------|------|
| `doc` | 编写或更新文档 |
| `doc.review` | 对照验收标准审查文档 |
| `doc.summary` | 生成文档摘要 |
| `doc.consolidate` | 合并文档文件 |
| `doc.drift` | 检测并修复规范漂移 |

### test（5 种）

| 类型 | 用途 |
|------|------|
| `test.gen-journeys` | 从规格生成测试旅程 |
| `test.gen-contracts` | 从旅程生成测试契约 |
| `test.gen-scripts` | 生成可执行测试脚本 |
| `test.run` | 运行测试脚本并收集结果 |
| `test.verify-regression` | 晋升后验证回归套件 |

### eval（2 种）

| 类型 | 用途 |
|------|------|
| `eval.journey` | 评估 Journey 质量（rubric 评分） |
| `eval.contract` | 评估 Contract 质量（rubric 评分） |

### validation（2 种）

| 类型 | 用途 |
|------|------|
| `validation.code` | 验证代码质量和正确性 |
| `validation.ux` | 验证用户体验质量 |

### 其他（2 种）

| 类型 | 用途 |
|------|------|
| `gate` | 阶段退出质量门禁 |
| `code-quality.simplify` | 简化和清理代码质量 |

---

## 架构

```
forge/
+-- plugins/forge/           # Forge plugin
|   +-- skills/              # 21 个 Skills
|   +-- commands/            # 18 个 Slash Commands
|   +-- agents/              # 1 个 Subagent (task-executor)
|   +-- hooks/               # Hooks + guide.md
+-- forge-cli/               # Go CLI 源码 (forge binary)
+-- tests/                   # 回归测试套件
+-- docs/                    # 项目文档
```

```
docs/features/<slug>/
+-- manifest.md              # Feature 单一入口（自动维护）
+-- prd/                     # /write-prd 产出
+-- design/                  # /tech-design 产出
+-- ui/                      # /ui-design 产出（可选）
+-- testing/                 # /gen-journeys + /gen-contracts 产出
+-- tasks/
    +-- index.json           # 任务定义
    +-- *.md                 # 任务详情
    +-- records/             # 执行记录
```

详细架构参见 [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md)。

---

## 贡献

```bash
git clone git@github.com:bigfaner/forge.git && cd forge
cd forge-cli && go mod download
go test -race -cover ./...
```

提交遵循 [Conventional Commits](https://www.conventionalcommits.org/)。

---

## 文档索引

| 文档 | 说明 |
|------|------|
| [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) | 核心架构、工作流管道、Agent 协作、Quality Gate |
| [forge-cli/docs/OVERVIEW.md](forge-cli/docs/OVERVIEW.md) | CLI 完整命令参考 |
| [forge-cli/docs/WORKFLOW.md](forge-cli/docs/WORKFLOW.md) | 内部流程图解 |
| [docs/official-references/plugin.md](docs/official-references/plugin.md) | 插件系统技术参考 |
| [docs/official-references/plugin-marketplace.md](docs/official-references/plugin-marketplace.md) | Marketplace 分发指南 |
| [docs/official-references/hooks.md](docs/official-references/hooks.md) | Hooks 技术参考 |

---

## License

[MIT](LICENSE)
