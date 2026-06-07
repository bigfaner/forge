---
created: 2026-06-06
author: "faner"
status: Draft
intent: "new-feature"
---

# Proposal: Forge Workspace — 多项目过程文档统一管理

## Problem

独立开发者同时推进多个 Forge 项目（4-8 个），但 Forge 是单项目视角设计——所有**过程文档**（proposals、features、PRDs、designs、tasks）锚定在项目目录内。多项目场景下产生具体痛点：

1. **过程文档散落**：每个项目的 `docs/` 下有大量过程文档（proposals、features、tasks），查找时需要先切到项目目录。无法在一个地方看到所有项目的 proposals 或 features
2. **缺乏全局状态视图**：无法一眼看到所有项目的阶段、活跃任务、阻塞项，需要逐个 `cd` + `forge task list`
3. **多项目切换摩擦**：在项目间切换时丢失上下文，需要反复重新加载状态
4. **知识无法复用**：lessons/conventions/decisions 沉淀在各自项目里，跨项目无法共享

> 四个痛点中，1-3 是**过程文档管理**问题（Workspace 解决），4 是**知识管理**问题（Wiki 解决）。本提案聚焦前者。

### Evidence

- 4-8 个活跃 Forge 项目分布在一个统一父目录下
- 查看「所有项目当前状态」需要逐个 `cd` + `forge task list`
- 已有 proposals 和 features 分布在各自项目的 `docs/` 下，没有跨项目的统一视图
- 已有两个 Draft 提案（forge-dashboard、forge-wiki）分别解决可见性和知识管理，但缺乏统一的**过程文档管理层**

### Urgency

每多一个项目，过程文档管理摩擦非线性增长。4-8 个项目的规模已影响日常效率。推迟意味着持续在项目间来回切换来了解全貌。

## Proposed Solution

引入 **Forge Workspace** 概念——项目父目录即工作区，作为多项目**过程文档**的统一管理入口。

### 三模块协作架构

Forge 多项目管理由三个互补模块组成，各有明确职责边界：

```
┌──────────────────────────────────────────────────────────┐
│                                                          │
│  ┌──────────────────┐                                    │
│  │   Workspace      │  过程文档管理                      │
│  │                  │  (proposals + features + tasks)     │
│  │  · 项目发现/注册  │                                    │
│  │  · 跨项目提案     │◄──── proposals 暂存与分配          │
│  │  · feature 总览   │◄──── features 聚合与状态追踪      │
│  │  · 任务聚合       │◄──── tasks 跨项目进度              │
│  └────────┬─────────┘                                    │
│           │                                              │
│           │ 提供结构化数据                                 │
│           ▼                                              │
│  ┌──────────────────┐                                    │
│  │   Dashboard      │  可视化查看                        │
│  │                  │  (文档 + 任务 + 状态)               │
│  │  · 项目状态卡片   │                                    │
│  │  · 活动时间线     │                                    │
│  │  · 文档/任务下钻  │                                    │
│  └──────────────────┘                                    │
│                                                          │
│  ┌──────────────────┐                                    │
│  │   Wiki           │  知识管理                          │
│  │                  │  (lessons/conventions/decisions)    │
│  │  · 知识 push/pull│                                    │
│  │  · 多维搜索      │                                    │
│  │  · 健康检查      │                                    │
│  └──────────────────┘                                    │
│                                                          │
│  .forge-workspace.yaml  (注册表 + 配置 + 模块开关)       │
│                                                          │
└──────────────────────────────────────────────────────────┘
         │              │              │
    ┌────┴────┐    ┌────┴────┐    ┌────┴────┐
    │ my-api  │    │ my-web  │    │cli-tool │  ...
    │ (项目)   │    │ (项目)   │    │ (项目)   │
    └─────────┘    └─────────┘    └─────────┘
```

**模块职责边界**：

| 模块 | 管理什么 | 不管理什么 | 优先级 |
|------|---------|-----------|--------|
| **Workspace** | 过程文档（proposals、features、tasks、PRDs、designs） | 知识（lessons/conventions/decisions） | P0（本次实现） |
| **Dashboard** | 可视化查看（消费 Workspace 的数据，展示文档与任务） | 数据生产、持久化 | P2（后续提案） |
| **Wiki** | 知识（lessons、conventions、decisions、business-rules） | 过程文档、UI 展示 | P1（后续提案，基于 forge-wiki 演进） |

**协作关系**：
- Workspace 提供过程文档的结构化数据 → Dashboard 消费并可视化
- Wiki 独立管理知识 → Workspace 的 feature 流程中可引用 Wiki 知识
- 三者通过 `.forge-workspace.yaml` 共享项目注册信息

**设计原则**：
- **职责正交**：过程文档（Workspace）与知识（Wiki）严格分离，不交叉
- **CLI 优先**：所有能力先有 CLI 入口，Dashboard 是可视化层
- **渐进采用**：`forge workspace init` 一条命令即可开始
- **项目不变**：现有项目结构完全不变，Workspace 是纯叠加层

### Innovation Highlights

**过程文档与知识的正交分离**：传统工具把所有文档混在一起管理。Forge Workspace 明确区分「过程文档」（proposals、features、tasks——有生命周期、有状态机）和「知识」（lessons、conventions、decisions——长期沉淀、跨项目复用）。两类文档的消费者、生命周期、管理方式完全不同，分开处理更自然。

**以过程文档为核心的项目管理**：Workspace 不试图成为通用项目管理工具——它只管理 Forge 管线产生的过程文档（brainstorm → PRD → design → tasks → execute）。这个限定让 Workspace 与 Forge 管线深度集成，而非泛化的项目管理。

**跨项目提案空间**：独立开发者的想法经常还没确定属于哪个项目。Workspace 级别的 proposals 提供了一个「创意暂存区」——先沉淀想法，后续再分配到具体项目。

## Requirements Analysis

### Key Scenarios

1. **每日全局视图**：运行 `forge workspace status`，一张表看到所有项目的阶段、任务进度、阻塞项
2. **跨项目提案**：想统一所有项目的错误处理规范 → `forge workspace propose "统一错误处理规范"` → 在 workspace 级别创建 proposal → brainstorm 技能在 workspace 上下文运行
3. **提案落地**：workspace 级 proposal 审批后 → `forge workspace assign <proposal> <project>` → 将 proposal 转为指定项目的 feature
4. **Feature 总览**：`forge workspace features` 跨项目列出所有 features 及其状态（哪个阶段、多少任务完成）
5. **项目注册**：新建项目后 → `forge workspace register ./new-project` → 加入 workspace 管理
6. **单项目下钻**：`forge workspace status <project>` 查看某项目的活跃 features、待办任务、近期记录

### Non-Functional Requirements

- `forge workspace status` 响应时间 < 2 秒（扫描 8 个项目）
- 项目发现兼容任意子目录命名
- 过程文档格式与现有项目 docs/ 保持一致（Markdown + frontmatter）

### Constraints & Dependencies

- 项目父目录下所有子目录平铺（非嵌套），每个项目包含 `.forge/config.yaml`
- Forge CLI 已有的路径解析和 manifest 读取能力
- 过程文档格式与现有技能（brainstorm、write-prd、tech-design、breakdown-tasks）产出的格式兼容

## Alternatives & Industry Benchmarking

### Industry Solutions

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零成本 | 过程文档管理摩擦持续，随项目增长恶化 | Rejected: 成本已在累积 |
| CLI 状态脚本 | shell alias / 自定义脚本 | 极简实现 | 无跨项目提案，无 feature 聚合，难维护 | Rejected: 不解决根本问题 |
| forge-dashboard (现有 Draft) | Web UI 多项目可视化 | 富交互 UI | 仅解决可见性，无过程文档管理，需常驻 HTTP 进程 | Superseded: Dashboard 降级为 Workspace 的可视化模块 |
| forge-wiki (现有 Draft) | Hub-Spoke 知识管理 | 完整知识架构 | 管理知识而非过程文档，无提案/feature 管理 | Complementary: Wiki 作为独立知识模块 |
| Monorepo | 所有项目合一仓库 | 天然共享文档 | 语义耦合、CI 复杂、违反项目独立性 | Rejected: 破坏现有工作流 |
| **Forge Workspace** | 过程文档统一管理 | CLI 优先、职责正交、与 Forge 管线深度集成 | v1 不含可视化和知识管理 | **Selected: 匹配过程文档管理的核心需求** |

## Feasibility Assessment

### Technical Feasibility

- 项目发现：扫描子目录 + 检测 `.forge/config.yaml`，纯文件系统操作
- 状态聚合：读取各项目 `docs/features/*/manifest.md` 和 `tasks/index.json`，现有格式可直接消费
- 跨项目 proposals：复用现有 proposal 模板，路径从项目级提升到 workspace 级
- Feature 总览：读取各项目 manifest.md 中的 status 字段，纯聚合操作

### Resource & Timeline

- Workspace 核心（项目发现 + 状态 + proposals + features 聚合）为中等规模工作量
- 所有实现基于现有 Forge CLI 框架，无需新基础设施

### Dependency Readiness

- Forge CLI 已有完整的文件读取、manifest 解析、config 加载能力
- 过程文档格式与现有技能产出格式兼容
- 无外部服务或 API 依赖

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| "过程文档和知识是同一类东西，应该一起管理" | Assumption Flip | 过程文档有生命周期（Draft → Approved → Done），知识是长期沉淀。消费者不同（过程文档由开发者追踪进度，知识由 agent 跨项目检索）。分开管理更自然：Workspace 管过程，Wiki 管知识 |
| "需要 Web Dashboard 来看项目状态" | XY Detection | CLI 表格输出对 4-8 个项目足够。Dashboard 是锦上添花，不是必需品。v1 做 CLI，Dashboard 留作后续模块 |
| "Proposals 必须属于具体项目" | 5 Whys | 独立开发者的想法经常在项目之前产生。先有想法再有项目是正常流程，workspace 级 proposals 匹配真实心智模型 |
| "项目文档需要迁移到 workspace" | Stress Test | 迁移成本高、风险大、破坏 git 历史。正确做法：项目不变，workspace 是纯叠加层 |
| "应该复用 forge-wiki 的设计作为知识共享" | Occam's Razor | forge-wiki 的 agentic search、session cache、降级策略对 v1 过重。知识管理独立演进为 Wiki 模块，与 Workspace 职责正交 |
| "独立开发者需要团队级工具" | Assumption Flip | 团队工具（权限、协作、审计）对独立开发者是噪音。聚焦个人效率：看见、整理、推进 |

## Scope

### In Scope

**v1 — Workspace 核心（过程文档管理）**：

1. **项目发现与注册**
   - `forge workspace init` — 自动扫描子目录，注册含 `.forge/config.yaml` 的项目
   - `forge workspace register <path>` — 手动注册项目
   - `.forge-workspace.yaml` — workspace 配置文件（项目列表 + 模块开关）

2. **跨项目状态总览**
   - `forge workspace status` — 表格输出：项目名、当前 phase、任务进度（done/total）、阻塞数
   - `forge workspace status <project>` — 单项目详情（活跃 features + tasks + 近期记录）

3. **Feature 总览**
   - `forge workspace features` — 跨项目列出所有 features 及其状态（阶段、任务进度）
   - `forge workspace features <project>` — 单项目 feature 列表

4. **Workspace 级 Proposals**
   - `docs/proposals/` 目录 — workspace 级提案存储
   - `forge workspace propose "<title>"` — 创建 workspace 级提案
   - brainstorm/eval 技能在 workspace 上下文运行
   - `forge workspace assign <proposal> <project>` — 将提案分配到具体项目，启动 feature 流程

**架构蓝图（定义接口，不实现）**：

5. **Wiki 模块接口** — 定义知识管理模块与 Workspace 的交互点（独立提案）
6. **Dashboard 模块接口** — 定义状态数据的标准输出格式（JSON），为可视化模块消费做准备

### Out of Scope

- **知识管理**（lessons、conventions、decisions 的 push/pull/search）→ Wiki 模块独立提案
- **Dashboard Web UI**（卡片视图、活动时间线、下钻）→ Dashboard 模块独立提案
- 跨项目任务调度（全局任务队列、项目间依赖）
- 多用户 / 团队协作
- 数字孪生 / 能力画像

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 项目结构变化导致发现失败 | L | M | 优雅降级：发现失败的项目在 status 中标记为「unhealthy」，不影响其他项目 |
| 跨项目 proposal 落地时上下文丢失 | M | M | assign 时自动从 workspace 级 proposal 继承关键上下文到目标项目的 PRD |
| 模块间耦合导致后续 Wiki/Dashboard 难独立演进 | L | H | 蓝图中定义清晰的模块接口；职责正交（过程文档 vs 知识 vs 可视化）确保边界清晰 |
| Feature 聚合时项目间 manifest 格式不一致 | M | M | 定义容错解析策略：缺失字段用默认值填充，解析失败不阻塞其他项目 |
| Workspace 级 proposals 与项目级 proposals 的认知混淆 | M | L | CLI 命令明确区分：workspace 级用 `forge workspace propose`，项目级用 `/brainstorm`（在项目内运行） |

## Success Criteria

- [ ] `forge workspace init` 自动发现父目录下所有 Forge 项目（含 `.forge/config.yaml` 的子目录），生成 `.forge-workspace.yaml`
- [ ] `forge workspace status` 输出包含项目名、当前 phase、任务进度（done/total）、阻塞数的表格，8 个项目扫描 < 2 秒
- [ ] `forge workspace status <project>` 显示单项目的活跃 features、待办任务和近期执行记录
- [ ] `forge workspace features` 跨项目列出所有 features 及其状态（阶段、任务进度）
- [ ] `forge workspace propose "统一错误处理"` 在 workspace 级 `docs/proposals/` 创建提案，brainstorm 技能可在 workspace 上下文运行
- [ ] `forge workspace assign <proposal-slug> <project>` 将 workspace 提案关联到目标项目，项目内可继承提案上下文启动 feature 流程
- [ ] 项目发现失败（目录缺失、manifest 损坏）时在 status 中标记为「unhealthy」，不阻塞其他项目
