# 架构概览

> 最后更新：2026-06-11

本文档从用户视角介绍 Forge 的架构设计。理解这些概念后，你将知道 Forge 是如何组织的、各组件做什么、以及数据在你的项目中如何流转。

---

## 目录

- [Forge 是什么](#forge-是什么)
- [插件机制](#插件机制)
- [四大组件](#四大组件)
- [数据流向](#数据流向)
- [状态管理](#状态管理)
- [目录约定](#目录约定)
- [工作模式](#工作模式)

---

## Forge 是什么

Forge 是一个 Claude Code 插件，核心理念是 **Spec Driven Development（SDD）** 和 **Harness Engineering**——将 AI 编程从自由对话变为受控的工程流水线。它通过结构化的工作流管道（规划、评估、执行、验证）来管理功能开发的全生命周期。

简单来说：

| 没有 Forge | 有 Forge |
|-----------|---------|
| 直接让 AI 写代码，方向容易漂移 | 先规划再动手，每一步都有文档产出 |
| 质量靠人工检查 | 对抗式评估 + Quality Gate 双重保障 |
| 跨会话知识容易丢失 | 任务记录和知识库自动沉淀 |

---

## 插件机制

### 安装位置

Forge 作为 Claude Code 插件安装到用户环境中，安装后的位置：

```
~/.claude/plugins/cache/forge/forge/<version>/
```

这个目录包含 Forge 的所有运行时组件（skills、commands、agents、hooks），不包含源码仓库中的文档、测试等开发资源。

### 如何加载

当你启动 Claude Code 时，插件系统会自动加载已安装的插件。Forge 通过 Hooks 在会话启动时注入规范指南（guide.md），确保 AI 助手在每次对话中都能遵循 Forge 的工作流规范。

### 调用方式

所有 Forge 功能通过 Claude Code 的斜杠命令（slash command）调用：

```
/brainstorm      # 规划类
/eval-prd        # 评估类
/run-tasks       # 执行类
```

斜杠命令会触发对应的 Skill 或 Command 组件来完成工作。

---

## 四大组件

Forge 由四类组件协作完成工作。下表列出每个组件的角色和交互方式：

| 组件 | 用途 | 触发方式 | 示例 |
|------|------|---------|------|
| **Skill** | 完成特定任务的技能模块，包含完整的工作流定义 | 斜杠命令 `/skill-name` | `/write-prd`、`/run-tasks`、`/eval-design` |
| **Command** | 轻量级命令入口，提供快速访问功能 | 斜杠命令 `/command-name` | `/git-commit`、`/submit-task` |
| **Agent** | 自主执行的子代理，由 Skill 或 Command 调度 | 由 Skill/Command 自动调度 | `task-executor` |
| **Hook** | 生命周期事件的自动触发器，在关键时刻执行检查 | 会话启动、停止等事件自动触发 | 会话结束时运行 Quality Gate |

### 对抗式评估体系

每种文档配备专属评估器（共 8 种），采用**专家角色 + 1000+ 分制评分表**进行多轮迭代修订，支持跨文档一致性校验（PRD → 设计 → 任务对齐）：

| 评估器 | 目标文档 | 说明 |
|--------|---------|------|
| `/eval-proposal` | proposal.md | 领域专家 + CTO 双层评审 |
| `/eval-prd` | prd/*.{3} | 产品经理视角评审 |
| `/eval-design` | design/*.{2} | 架构师视角评审 |
| `/eval-ui` | ui/*.{2} | UI/UX 专家评审 |
| `/eval-journey` | journey.md | 测试策略评审 |
| `/eval-contract` | contracts/*.md | 行为契约评审 |
| `/eval-consistency` | 跨文档 | PRD → 设计 → 任务一致性校验 |
| `/eval` | 任意文档 | 通用评估入口 |

### Journey-Contract 测试模型

从 PRD 到可执行测试脚本的两层抽象：

```
PRD 用户故事
     │
     ▼
Journey（用户流程 + 风险分级）  ← /gen-journeys
     │
     ▼
Contract（六维度行为契约）      ← /gen-contracts
     │
     ▼
Surface 感知测试脚本            ← /gen-test-scripts
  web / api / cli / tui / mobile
```

### 组件协作关系

```
用户输入斜杠命令
        │
        ▼
   Skill / Command
        │
        ├── 读取模板和规范
        ├── 调用 forge CLI 管理状态
        ├── 调度 Agent 执行开发工作
        └── 返回结果
        │
        ▼
   Agent（task-executor）
        │
        ├── 读取项目知识
        ├── TDD 实现代码
        ├── 运行 Quality Gate
        └── 提交记录
        │
        ▼
   Hook（自动触发）
        │
        ├── 会话启动：注入上下文
        └── 会话结束：运行质量验证
```

---

## 数据流向

以下是从用户输入到文件系统变更的完整数据流：

```
┌──────────────────────────────────────────────────────────────────┐
│                         用户交互层                                │
│                                                                  │
│   用户输入斜杠命令（如 /brainstorm、/run-tasks）                  │
└───────────────────────────┬──────────────────────────────────────┘
                            │
                            ▼
┌──────────────────────────────────────────────────────────────────┐
│                         Forge 处理层                              │
│                                                                  │
│   Skill/Command 接收命令                                         │
│         │                                                        │
│         ├── 加载模板和规范文件                                    │
│         ├── 调用 forge CLI 读写状态                               │
│         │       │                                                │
│         │       ├── forge task add/claim/submit                  │
│         │       ├── forge feature set                            │
│         │       └── forge quality-gate                           │
│         │                                                        │
│         └── 调度 Agent（task-executor）                           │
│                 │                                                │
│                 ├── 读取 docs/business-rules/、docs/conventions/  │
│                 ├── 执行 TDD 开发                                │
│                 └── 运行 Quality Gate（compile → fmt → lint → unit-test）│
│                                                                  │
└───────────────────────────┬──────────────────────────────────────┘
                            │
                            ▼
┌──────────────────────────────────────────────────────────────────┐
│                         文件系统层                                │
│                                                                  │
│   docs/features/<slug>/          # Feature 工作区                │
│       ├── manifest.md            # Feature 入口                  │
│       ├── prd/                   # 需求文档                      │
│       ├── design/                # 设计文档                      │
│       ├── tasks/                 # 任务定义和记录                 │
│       └── testing/               # 测试产出                      │
│                                                                  │
│   docs/business-rules/           # 跨功能业务规则                 │
│   docs/conventions/              # 技术规范                       │
│   docs/decisions/                # 技术决策                       │
│   docs/lessons/                  # 经验教训                       │
│                                                                  │
│   tests/<surface>/               # 测试脚本                      │
│   .forge/config.yaml             # 项目配置                      │
│                                                                  │
└──────────────────────────────────────────────────────────────────┘
```

### 数据流要点

1. **所有状态变更通过 forge CLI**：任务创建、认领、提交都通过 `forge task` 命令完成，确保状态一致性
2. **文档驱动开发**：每个阶段都有对应的文档产出（proposal → PRD → design → tasks），文档是流程推进的依据
3. **Hook 自动保障**：在会话启动时注入规范，结束时运行质量检查，用户无需手动触发

---

## 状态管理

### Feature 状态流转

每个 Feature 通过 `manifest.md` 跟踪状态，状态由各 Skill 在完成对应阶段后自动设置：

```
(none) ──→ prd ──→ design ──→ tasks ──→ in-progress ──→ completed
              │         │         │           │
          /write-prd    │    /breakdown-  首次 task claim
                       │    tasks 完成
                       │
                 /tech-design +
                 /ui-design 完成
```

| 状态 | 含义 | 触发条件 |
|------|------|---------|
| `(none)` | 初始状态，尚无文档 | Feature 创建时 |
| `prd` | PRD 文档就绪 | `/write-prd` 完成 |
| `design` | 设计文档就绪 | `/tech-design` 或 `/ui-design` 完成（后完成者设置） |
| `tasks` | 任务已拆分 | `/breakdown-tasks` 完成 |
| `in-progress` | 任务执行中 | 首次 `forge task claim` |
| `completed` | 全部完成 | 所有任务完成 |

### 任务状态

单个任务在执行过程中有以下状态：

```
pending ──→ in-progress ──→ completed
                │
                ├──→ blocked     # 依赖其他任务或遇到错误
                └──→ skipped     # 被跳过
```

### Quality Gate

所有 coding 任务在提交前必须通过 Quality Gate，按顺序执行验证：

```
compile → fmt → lint → unit-test → test → probe
```

任何一步失败则停止，需要修复后重新通过。根据任务类型，Gate 的严格程度不同：

| 任务类型 | Gate 范围 | 说明 |
|---------|----------|------|
| Breaking 任务 | compile → fmt → lint → unit-test | 快速反馈，含单元测试 |
| Non-breaking 任务 | compile → fmt → lint | 仅静态检查 |
| 全部完成时 | compile → fmt → lint → unit-test → test → probe | 项目级全量验证 |

---

## 目录约定

### 用户项目结构

Forge 在用户项目中维护以下目录结构：

```
<你的项目>/
├── .forge/
│   └── config.yaml              # Forge 项目配置
│
├── docs/
│   ├── ARCHITECTURE.md          # 系统架构（你维护）
│   ├── business-rules/          # 跨功能业务规则（自动生成）
│   ├── conventions/             # 技术规范（自动生成）
│   ├── decisions/               # 技术决策（/learn 生成）
│   ├── lessons/                 # 经验教训（/learn 生成）
│   ├── proposals/               # 改进提案（/brainstorm 生成）
│   ├── sitemap/sitemap.json     # 页面元素映射（/gen-web-sitemap 生成）
│   │
│   └── features/<slug>/         # Feature 工作区（核心）
│       ├── manifest.md          # Feature 入口（自动维护）
│       ├── prd/                 # /write-prd 产出
│       │   ├── prd-spec.md
│       │   ├── prd-user-stories.md
│       │   └── prd-ui-functions.md  # 可选
│       ├── design/              # /tech-design 产出
│       │   ├── tech-design.md
│       │   └── api-handbook.md      # 可选
│       ├── ui/                  # /ui-design 产出（可选）
│       │   ├── ui-design.md
│       │   ├── DESIGN.md            # 可选
│       │   └── prototype/           # 可选
│       ├── testing/             # 测试文档（自动生成）
│       │   └── <journey>/
│       │       ├── journey.md
│       │       └── contracts/
│       ├── tasks/               # 任务（核心）
│       │   ├── index.json       # 任务索引（forge CLI 维护）
│       │   ├── *.md             # 任务详情文件
│       │   ├── records/         # 执行记录（forge task submit 生成）
│       │   └── specs/           # 规范预览（/consolidate-specs 生成）
│       └── eval/                # 评估报告（可选）
│
└── tests/                       # 测试脚本
    └── <surface-key>/
        ├── helpers.ts
        ├── features/<slug>/     # 功能级测试
        └── <module>/            # 回归测试（自动晋升）
```

### 关键目录说明

| 目录 | 生成方式 | 是否提交到 Git |
|------|---------|--------------|
| `.forge/config.yaml` | `forge init` | 是 |
| `docs/features/<slug>/manifest.md` | 各 Skill 自动维护 | 是 |
| `docs/features/<slug>/tasks/index.json` | `/breakdown-tasks` 或 `forge task index` | 是 |
| `docs/features/<slug>/tasks/records/` | `forge task submit` | 是 |
| `docs/business-rules/` | `/consolidate-specs` | 是 |
| `docs/conventions/` | `/consolidate-specs` 或 `/test-guide` | 是 |
| `docs/decisions/` | `/learn` | 是 |
| `docs/lessons/` | `/learn` | 是 |

### 重要约束

- **`tasks/records/`** 由 `forge task submit` 自动生成，不要手动编辑
- **`docs/business-rules/`** 和 **`docs/conventions/`** 由 `/consolidate-specs` 从 feature 文档中提取，agent 在执行任务时会读取这些规范
- **`tests/`** 中的测试通过标签管理（`@feature` → `@regression` 晋升），不使用文件迁移

---

## 工作模式

Forge 提供两种工作模式，根据功能复杂度选择：

### 完整模式

适用于复杂功能（>10 任务），经过完整的文档化流程：

```
/brainstorm ──→ /write-prd ──→ /tech-design ──→ /breakdown-tasks ──→ /run-tasks
     │               │              │  ──→ /ui-design       │              │
 proposal.md    prd/*.{3}     design/*.{2}  ui/    tasks + index.json  自动执行
```

每个阶段都有文档产出，可选通过 `/eval-*` 系列技能迭代评分至达标。

### 快速模式

适用于小功能（1-10 任务），跳过 PRD 和设计：

```
/quick ──→ /brainstorm ──→ /quick-tasks ──→ /run-tasks
                │               │               │
           proposal.md     tasks/*.md       自动执行
```

### 如何选择

| 场景 | 推荐模式 |
|------|---------|
| 全新功能，涉及架构决策 | 完整模式 |
| 需要 UI 设计 | 完整模式 |
| 有明确的验收标准 | 完整模式 |
| Bug 修复 | 快速模式 |
| 小改动、配置调整 | 快速模式 |
| 重构已有功能 | 快速模式 |

### Pipeline 全景

```
┌─────────────────────────────────────────────────────┐
│ 完整模式 Pipeline                                    │
│                                                     │
│ /brainstorm → /write-prd → /ui-design? → /tech-design │
│      ↓            ↓           ↓            ↓        │
│  /eval-proposal /eval-prd  /eval-ui   /eval-design   │
│                                                     │
│ → /breakdown-tasks → /run-tasks → /submit-task       │
│                                                     │
└─────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────┐
│ 快速模式 Pipeline                                    │
│                                                     │
│ /quick → /brainstorm → /quick-tasks → /run-tasks    │
│                                                     │
└─────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────┐
│ 测试 Pipeline（完整模式）                             │
│                                                     │
│ /gen-journeys → /eval-journey → /gen-contracts       │
│ → /eval-contract → /gen-test-scripts → /run-tests    │
│                                                     │
└─────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────┐
│ 辅助 Skill（任意阶段可用）                            │
│                                                     │
│ /consolidate-specs  /learn  /forensic  /deep-research│
│ /clean-code  /gen-web-sitemap  /test-guide           │
│ /init-justfile  /extract-design-md  /eval  /fix-bug  │
│                                                     │
└─────────────────────────────────────────────────────┘
```
