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

`/fix-bug` · `/git-commit` · `/git-checkout` · `/learn-lesson` · `/record-decision` · `/consolidate-specs` · `/init-justfile` · `/init-forge` · `/gen-sitemap` · `/extract-design-md` · `/forensic` · `/improve-harness`

---

## Agents

| Agent | 职责 |
|-------|------|
| **task-executor** | 执行单个任务（TDD + Quality Gate + record） |

### Eval Experts（评估协议 + 专家角色）

评估通过 protocol + expert 组合执行，不使用独立 agent 定义文件。参见 `agents/experts/`。

---

## 项目结构

```
forge/
├── plugins/forge/          # Forge plugin
│   ├── skills/             # 17 个 Skills
│   ├── commands/           # 17 个 Slash Commands
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

## Task Types & Pipeline 参考

> 以下内容由 `forge task index` 自动生成，以 CLI 行为为准。

### 13 种任务类型

| Type | 谁生成 | 用途 |
|------|--------|------|
| `implementation` | Skill agent | 实现功能代码 |
| `documentation` | Skill agent | 编写或更新文档 |
| `doc-evaluation` | `forge task index`（docs-only） | 评估文档质量（T-eval-doc） |
| `doc-generation.summary` | `forge task index` | 生成阶段摘要（`N.summary`） |
| `doc-generation.consolidate` | `forge task index` | 合并规格文档（T-test-5） |
| `test-pipeline.gen-cases` | `forge task index` | 生成测试用例（T-test-1 / T-quick-1） |
| `test-pipeline.eval-cases` | `forge task index` | 评估测试用例质量（T-test-1b） |
| `test-pipeline.gen-scripts` | `forge task index` | 生成可执行测试脚本（T-test-2 / T-quick-2） |
| `test-pipeline.run` | `forge task index` | 运行测试并收集结果（T-test-3 / T-quick-3） |
| `test-pipeline.graduate` | `forge task index` | 将测试晋升到回归套件（T-test-4 / T-quick-4） |
| `test-pipeline.verify-regression` | `forge task index` | 验证完整回归套件（T-test-4.5 / T-quick-5） |
| `fix` | `forge task add` | 修复失败的测试或质量门禁（`fix-*` / `disc-*`） |
| `gate` | `forge task index` | 阶段退出质量门禁（`N.gate`） |

### Quick Pipeline 责任链（T-quick-1~5）

| ID | 职责 |
|----|------|
| T-quick-1 | 从 proposal Success Criteria 生成测试用例文档（无 sitemap、无 eval） |
| T-quick-2 | 从测试用例生成 e2e 测试脚本 |
| T-quick-3 | 执行 feature e2e 测试；失败则标记 blocked 并添加 fix task（P0） |
| T-quick-4 | 将测试脚本晋升到回归套件 `tests/e2e/` |
| T-quick-5 | 运行完整回归套件；失败则标记 blocked 并添加 fix task（P0） |

依赖链：每个 profile 的 T-quick-1~4 串行，T-quick-5 依赖所有 T-quick-4。

### Full Pipeline 责任链（T-test-1~5）

| ID | 职责 |
|----|------|
| T-test-1 | 生成测试用例文档（先调 `/gen-sitemap` 如果 sitemap.json 缺失，再调 `/gen-test-cases`） |
| T-test-1b | 评估测试用例的下游可执行性（main session 任务，调用 `/eval-test-cases`） |
| T-test-2 | 从评估后的测试用例生成测试脚本（调用 `/gen-test-scripts`） |
| T-test-3 | 执行 feature e2e 测试；失败则标记 blocked 并添加 fix task（P0），修复后重跑 |
| T-test-4 | 验证 e2e 通过（检查 `latest.md`），然后晋升脚本到 `tests/e2e/` |
| T-test-4.5 | 运行完整回归套件；失败则标记 blocked 并添加 fix task（P0），修复后重跑 |
| T-test-5 | 提取业务规则和技术规格，用户确认后合并（调用 `/consolidate-specs`） |

依赖链：T-test-1 → T-test-1b → 每个 profile 的 T-test-2~4 串行 → T-test-4.5（依赖所有 T-test-4） → T-test-5。

### Fix-task 命令

```bash
forge task add --template fix-task --title "Fix: <描述>" \
  --source-task-id <源任务ID> \
  --block-source \
  --var SOURCE_FILES="<受影响文件>" \
  --var TEST_SCRIPT="<失败的测试>" \
  --var TEST_RESULTS="<结果路径>" \
  --description "<根因分析>"
```

| 参数 | 说明 |
|------|------|
| `--block-source` | 原子操作：在 fix-task 创建前将源任务设为 blocked |
| 去重 | `task add` 自动去重：输出 `ACTION: ADDED`（新建）或 `ACTION: SKIPPED`（已有活跃 fix-task） |
| 自动恢复 | fix-task 完成后，`task record` 自动将源任务恢复为 pending（需所有依赖完成） |
| 嵌套 | fix-task 本身失败时，`--source-task-id` 指向失败的 fix-task（非原始源），最多 3 层 |

### Profile-suffix 规则

| 场景 | ID 格式 | 示例 |
|------|---------|------|
| 单 profile（默认） | 无后缀 | `T-test-2`, `T-quick-1` |
| 多 profile（2+） | 字母后缀（a, b, c...） | `T-test-2a`, `T-test-2b`, `T-quick-1a` |

后缀规则：第一个 profile 为 `a`，第二个为 `b`，以此类推。共享任务（T-test-1、T-test-1b、T-test-4.5、T-test-5、T-quick-5）无后缀。

### Gate/Summary 自动生成规则

`forge task index` 自动检测阶段（phase）并生成 stage-gate 文件：

- **检测条件**：任务 ID 匹配 `<数字>.<数字>` 格式（如 `1.1`, `2.3`），排除 `T-test-*`、`T-quick-*`、`.summary`、`.gate`
- **生成条件**：同一 phase 内有 >= 2 个业务任务时，生成 `<N>.summary.md` 和 `<N>.gate.md`
- **依赖关系**：`N.summary` 依赖 phase N 所有业务任务；`N.gate` 依赖 `N.summary`；下一 phase 的任务依赖 `N.gate`
- **幂等性**：已存在的文件不会覆盖

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
