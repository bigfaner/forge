---
created: 2026-05-20
author: "faner"
status: Draft
---

# Proposal: Forge Wiki — 开发者数字孪生的知识底座

## Vision

Forge Wiki 的终极目标是构建开发者的**数字孪生**——一个完整记录你的决策风格、经验教训、技术偏好、成长轨迹的知识系统。当 AI agent 拥有这个数字孪生，它不再是一个通用助手，而是一个完全适配你的个人 AI 搭档。

**三个核心原则**：
- **Agent 优先**：知识库的第一消费者是 AI agent，不是人。设计围绕"agent 如何最高效地检索和应用知识"展开
- **时间感知**：知识有生命周期。2024 年的"最佳实践"在 2026 年可能是反模式。每条知识都有时间维度
- **持续复利**：知识编译一次、持续维护、永久积累。不是每次查询重新推导，而是站在已有认知上继续构建

## Problem

Forge 在每个项目中产出了大量高质量文档（PRDs、designs、decisions、lessons、conventions、business-rules），但这些知识被锁死在各自的项目仓库里，无法跨项目复用。每次启动新项目，都要从零开始重新积累约定、重新踩已知的坑。

更深层的问题：AI agent 每次会话都是"失忆"的——它从 CLAUDE.md 重建上下文，但没有跨项目、跨时间的长期记忆。agent 永远不会"越来越懂你"。

### Evidence

- 3-5 个活跃项目，每个项目都有独立的 `docs/conventions/`、`docs/business-rules/`、`docs/decisions/`、`docs/lessons/`
- 跨项目存在重复的 conventions（如错误处理模式、API 设计规范），但没有机制去发现和复用
- `/learn` 和 `/consolidate-specs` 只在单项目内运作，知识无法流动
- 项目完成后，积累的决策和教训随仓库归档，新项目无法参考
- 同一个 agent 在不同项目中无法复用之前积累的领域知识

### Urgency

AI 时代知识是最重要的资产。每个项目从零开始是巨大的浪费。当前 Forge 的知识管理是"写了就忘"的模式——知识在产出时有用，但没有沉淀机制。随着项目数量增长，这个问题会越来越严重。

## Proposed Solution

基于 Karpathy 的 LLM Wiki 模式，构建一个以 agent 为第一消费者的个人软件开发现知识库——数字孪生的知识底座。各项目产出的文档作为 Raw Sources，AI 自动提炼、交叉引用、持续维护一个结构化的知识 Wiki。

**三层架构**：
- **Raw Sources**：各项目的 docs/ 目录（PRDs、designs、decisions、lessons、conventions 等）— 只读，不可变
- **Wiki**：独立目录下的 AI 维护的结构化 Markdown 知识库 — AI 全权维护
- **Schema**：告诉 AI 如何维护 wiki 的指令文档（CLAUDE.md 或类似机制）

**知识流：Forge Push → Wiki Receive**
- **Forge（生产者）**：在知识产生时主动推送给 Wiki。触发点包括：`/learn` 记录决策/教训、`/consolidate-specs` 提取约定、任务完成时记录模式、Bug 修复时捕获根因
- **Wiki（接收者 + 维护者）**：接收 Forge 推送的知识，自动整合——去重、交叉引用、标记矛盾、标注时间戳和时效状态

**Wiki 三种操作**：
- **Receive**：接收 Forge 推送的知识，整合进 Wiki，更新交叉引用，标记矛盾。这是被动操作，由 Forge 的推送触发
- **Query**：自然语言查询 Wiki，返回综合答案，查询结果可回填为新的 Wiki 页面。**Agent 优先**：Query 接口设计为 agent 友好，返回结构化、可直接注入上下文的知识片段
- **Lint**：定期健康检查——矛盾检测、**过期标记**、孤立页面、缺失交叉引用

**时间维度**：每条知识条目携带时间戳和时效标签。Agent 查询时自动感知知识的新鲜度，优先引用近期知识，对过期知识标记警告。知识演变历史可追溯——"为什么我们现在用 gRPC 而不是 GraphQL"。

**Agent 集成**：Forge 任务执行时，AI agent 自动从 Wiki 检索相关知识注入上下文。不是 agent 去找知识，而是**知识找 agent**——基于当前任务的领域标签自动推送相关知识。

**数字孪生积累**：随着项目经验不断 ingest，Wiki 逐步构建出开发者的能力画像——决策偏好、常见模式、踩过的坑、擅长的领域。这是数字孪生的知识底座。

### Innovation Highlights

**Agent 优先的知识架构**：传统知识库为人设计（文件夹、标签、搜索框）。Forge Wiki 为 agent 设计——知识结构、检索接口、返回格式都优化为 agent 可直接消费的形式。人通过 Obsidian 浏览，agent 通过 CLI/Skill 查询，同一个 Wiki 服务两个消费者。

**时间感知的知识生命周期**：知识不是静态快照，而是有时间维度的活文档。2024-2026 年间对同一问题的不同决策形成演变轨迹。Agent 在建议方案时会说"我们 2024 年用 REST，2025 年试了 GraphQL，2026 年回归 gRPC，原因是复杂度管理"。这让 agent 的建议不是基于最新的文档，而是基于你的完整决策历史。

**LLM Wiki 模式应用于软件开发**：Karpathy 的 LLM Wiki 主要用于研究和个人笔记管理，我们将其应用于软件开发现知识管理——一个更结构化、更可验证的领域，并在此基础上构建数字孪生的知识底座。

**知识与工具的闭环**：Forge 是知识的源头（通过 `/learn`、`/consolidate-specs`、任务执行产出），在知识产生的瞬间主动推送给 Wiki。Forge 同时也是知识的消费者（agent 执行时从 Wiki 检索）。形成"产出 → 推送 → 维护 → 检索 → 消费 → 新产出"的正循环。知识流动是事件驱动的，不需要人工触发。

**AI 维护而非人工维护**：传统知识库的死因是维护成本随规模非线性增长。LLM Wiki 的核心洞察是让 AI 承担所有维护工作——交叉引用、摘要更新、冲突检测、一致性维护。

## Requirements Analysis

### Key Scenarios

1. **知识自动推送**：任务执行中 `/learn` 记录了一条教训 → Forge 自动推送给 Wiki → Wiki 接收并整合（去重、交叉引用、标注时间戳）。知识在产生的瞬间就进入 Wiki，无需手动 ingest
2. **Agent 自动参考**：任务执行时，agent 基于当前任务的领域标签自动从 Wiki 检索相关知识。不需要显式查询，知识自动注入上下文
3. **查询知识演变**：`forge wiki query "API 设计风格的演变"` 返回时间线视图——"2024 REST → 2025 GraphQL → 2026 gRPC"，包含每次转变的原因
4. **Wiki 健康检查**：`forge wiki lint` 检测矛盾、过期条目（>6个月未验证）、孤立页面、缺失交叉引用
5. **跨项目知识发现**：发现多个项目中重复出现的模式，提炼为通用知识条目

### Non-Functional Requirements

- Wiki 存储为纯 Markdown 文件，可被任何编辑器打开（Obsidian、VS Code、终端）
- 知识库独立于 Forge 安装目录，用户指定位置
- Agent Query 返回结构化格式（不依赖 LLM 二次解析），可直接注入 agent 上下文
- 每条知识条目必须包含：时间戳、来源项目、领域标签、时效状态
- Wiki 规模达到数百页面时仍可有效检索

### Constraints & Dependencies

- 依赖 Forge CLI 的项目发现能力（项目注册表或目录扫描）
- 依赖 LLM（Claude）进行知识提炼和维护操作
- Wiki 目录结构需要兼容 Obsidian 的 wiki link 和 graph view
- Raw Sources 只读，Wiki 层才能写入

## Alternatives & Industry Benchmarking

### Industry Solutions

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零成本 | 知识持续断裂，每个项目从零开始 | Rejected: 核心问题会随项目增长恶化 |
| RAG 知识库 | NotebookLM / ChatGPT file upload | 查询方便 | 无持久积累，每次重新推导，无交叉引用 | Rejected: 不解决知识沉淀问题 |
| 传统 Wiki | Confluence / Notion | 人工可编辑 | 维护成本高，人不愿意维护 | Rejected: 维护负担是核心痛点 |
| 通用 LLM Wiki | nashsu/llm_wiki | 功能完善（图谱、向量搜索、桌面应用） | 通用工具，无 Forge 集成，无 Push 模型，无时间维度 | Rejected: 无法形成 Forge 知识闭环 |
| **Forge Wiki** | Karpathy's pattern + Forge 生态 | Push 集成、Agent 自动注入、时间感知、开发领域特化 | 前期需自建，通用能力需逐步补齐 | **Selected: 知识闭环是任何通用工具无法提供的独特价值** |

### 开源 LLM Wiki 对比

| 维度 | nashsu/llm_wiki | SamurAIGPT/llm-wiki-agent | lucasastorian/llmwiki | **Forge Wiki** |
|------|-----------------|---------------------------|----------------------|----------------|
| 形态 | Tauri 桌面应用 | Agent Skill（CLAUDE.md） | Web + MCP Server | Forge Skill + CLI |
| 知识来源 | 通用文档 | 通用文档 | 通用文档 | **Forge 项目文档** |
| Ingest | 手动拖入 | 手动命令 | 手动上传 | **Forge 自动推送** |
| Agent 集成 | HTTP API | 原生嵌入 | MCP Server | **任务执行时自动注入** |
| 时间维度 | 无 | 无 | 无 | **有** |
| Obsidian | 原生兼容 | Symlink 兼容 | 无 | 原生兼容 |

### 借鉴点

| 来源 | 借鉴 | 原因 |
|------|------|------|
| nashsu/llm_wiki | Two-Step Ingest（先分析再生成） | 提高提炼质量 |
| nashsu/llm_wiki | Review System（异步人工审核队列） | 知识准确性需要人工校验 |
| SamurAIGPT/llm-wiki-agent | 确定性 + 语义边双重建图 | 知识图谱同时包含显式和隐含关系 |
| lucasastorian/llmwiki | 文件系统为 source of truth | Wiki 就是 Markdown 文件，简单可靠 |

## Feasibility Assessment

### Technical Feasibility

- Wiki 本质是 Markdown 文件 + CLAUDE.md 指令，无需额外基础设施
- Ingest/Query/Lint 操作可封装为 Forge CLI 子命令或 Claude Code Skills
- 项目文档已经是结构化 Markdown，LLM 可直接读取和提炼

### Resource & Timeline

- Schema 设计和 Wiki 目录结构需要 1-2 轮迭代
- CLI 命令封装可复用 Forge 现有框架
- 预计中等规模工作量

### Dependency Readiness

- Forge CLI 已有文件读取和项目发现能力
- Claude Code 已有文件读写和目录操作能力
- 无外部 API 或服务依赖

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| "需要一个 Web Dashboard 来管理知识" | XY Detection | 真正的问题是知识断裂而非缺少 UI。LLM Wiki 的 CLI/Skill 模式更适合与 Forge 集成 |
| "知识需要手动编辑和整理" | Assumption Flip | 如果 AI 能维护，就不需要手动编辑。Karpathy 的模式证明 AI 维护的 Wiki 是可行的 |
| "所有项目知识都值得保留" | Stress Test | 大部分项目文档是过程性的，只有提炼后的知识（patterns、anti-patterns、decisions）才值得跨项目保留。Ingest 操作需要筛选机制 |
| "知识库应该为人设计" | Assumption Flip | 知识的第一消费者是 agent，不是人。传统知识库为人优化（文件夹、标签、搜索框），Forge Wiki 应该为 agent 优化（结构化输出、领域标签匹配、上下文片段） |
| "知识是静态的" | 5 Whys | 知识会过时。没有时间维度的知识库会给 agent 过期的建议。时间戳和时效状态是必需的，不是可选的 |

## Scope

### In Scope

- Wiki 目录结构设计（兼容 Obsidian，agent 优先）
- Schema 指令文档（AI 维护规则）
- **时间维度**：每条知识的时间戳、时效标签、演变历史
- Receive 操作：接收 Forge 推送的知识并自动整合（去重、交叉引用、标注来源和时间戳）
- Query 操作：自然语言查询 Wiki（agent 友好的结构化返回）
- Lint 操作：Wiki 健康检查（含过期检测）
- Agent 集成：任务执行时基于领域标签自动推送相关知识
- Forge Push 集成：在 `/learn`、`/consolidate-specs`、任务完成等知识产生点自动推送给 Wiki

### Out of Scope

- Web UI / Dashboard 可视化（属于 Forge Dashboard 提案）
- 多用户 / 团队协作
- 知识的手动编辑界面（v1 由 AI 全权维护）
- 向量检索 / RAG 基础设施（Wiki 规模较小时 index.md 足够）
- 跨项目任务管理
- 数字孪生的高级能力（能力画像、风格指纹、决策模型）— 作为远期愿景，v1 只积累原始数据

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| AI 提炼质量不稳定，Wiki 内容不准确 | M | H | Lint 操作定期检查；Raw Sources 可溯源；Ingest 时用户参与确认 |
| Wiki 规模增长后检索效率下降 | M | M | 初期用 index.md 索引，规模增大后引入轻量搜索工具（如 qmd） |
| Schema 指令不够精确，AI 维护行为不一致 | M | M | Schema 迭代优化，像 CLAUDE.md 一样持续调整 |
| 项目文档格式不统一，影响 Ingest 质量 | L | M | 利用现有 consolidate-specs 的标准化格式作为 Raw Sources |
| Agent 自动推送的知识不相关或过多 | M | M | 基于领域标签精准匹配，设置相关度阈值；agent 可忽略不相关知识 |
| 时间维度增加 Wiki 维护复杂度 | L | M | Lint 操作自动检测过期条目；时效标签由 AI 在 Ingest 时自动标注 |

## Success Criteria

- [ ] Forge 在 `/learn`、`/consolidate-specs`、任务完成时自动推送知识到 Wiki，至少覆盖 3 个项目的知识产出
- [ ] 执行 `forge wiki query "API 设计"` 能返回跨项目的综合知识（引用来自 2+ 个项目的知识），并包含时间演变信息
- [ ] 任务执行时 agent 能基于领域标签自动引用 Wiki 中相关知识（至少在 50% 的相关场景中被引用）
- [ ] Wiki 通过 `forge wiki lint` 检查，无矛盾条目，过期条目被正确标记
- [ ] Agent Query 返回的知识片段格式为可直接注入上下文的结构化格式，无需 LLM 二次解析
- [ ] Wiki 可在 Obsidian 中正常浏览，graph view 可显示知识关联

## Next Steps

- Proceed to `/write-prd` to formalize requirements
