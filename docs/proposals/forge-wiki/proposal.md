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

### 架构：Hub-Spoke 知识中心

```
Project A ──push──►┌─────────────┐◄──push── Project B
                   │  forge-wiki  │
Project C ◄──pull──└─────────────┘──pull──► New Project D
                   (集中知识库)
```

- **Push**：项目将本地产出的知识（decisions、lessons、conventions 等）推送到 forge-wiki 集中存储
- **Pull**：项目从 forge-wiki 检索公共的、必要的规范和领域知识

**三层架构**：
- **Raw Sources**：各项目的 docs/ 目录（PRDs、designs、decisions、lessons、conventions 等）— 只读，不可变
- **Wiki**：独立目录下的 AI 维护的结构化 Markdown 知识库 — AI 全权维护
- **Schema**：告诉 AI 如何维护 wiki 的指令文档（CLAUDE.md 或类似机制）

### CLI 接口：`forge kb` 命令族

知识库的第一消费者是 AI Agent，CLI 接口围绕 Agent 的交互模式设计。

| 命令 | 方向 | 作用 |
|------|------|------|
| `forge kb push` | Push | 项目知识推送到集中 wiki |
| `forge kb search --domains <domains> <query>` | Pull | 搜索相关文档，返回元数据列表（名称 + 标签 + 描述） |
| `forge kb get <name>` | Pull | 获取完整文档内容 |
| `forge kb get <name> --outline` | Pull | 仅获取文档大纲（标题结构） |
| `forge kb get <name> --section <section>` | Pull | 仅获取文档的指定章节 |
| `forge kb domains` | Meta | 列出 wiki 中所有可用知识领域及文档数量 |

### Agentic Search：两步检索模式

Agent 的检索遵循"轻量发现 → 精准获取"模式，避免一次性加载大量无关内容：

**Step 1 — 搜索（便宜）**：
```bash
forge kb search --domains auth "认证流程"
# → [{name: "auth-flow", tags: ["auth", "jwt"], desc: "..."},
#    {name: "session-mgmt", tags: ["auth", "session"], desc: "..."}]
```
只返回元数据（名称、标签、描述），Agent 自行判断哪些文档相关。

**Step 2 — 按需获取（精准）**：
```bash
# 先看大纲，决定需要哪个章节
forge kb get auth-flow --outline
# → 1.概述 2.JWT配置 3.错误处理 4.刷新机制

# 只拉相关章节
forge kb get auth-flow --section "错误处理"
```

"Agentic" 的本质不是搜索算法多智能，而是 **Agent 在中间做了判断** —— 第一步只给元数据（节省 token），Agent 推理后决定深入哪篇文档的哪个章节（精准消费）。

### Domains 双向利用

知识文档使用 frontmatter 中的 `domains` 字段作为领域标签。搜索时反向利用 `domains` 进行双重过滤：

1. **`--domains` 精确缩小候选集**：字符串匹配，cheap
2. **query 在候选集内做关键词匹配**：在缩小的范围内进一步筛选

Agent 获取 `--domains` 的成本几乎为零 —— 任务本身的 `domains` frontmatter 就有这个值。

### Agent 上下文注入

Agent 不需要"知道" forge-wiki 的存在——通过 Forge 插件现有的 `session-start` hook 自动注入。

在 `session-start` 脚本中追加 wiki 领域目录：
```bash
WIKI_DOMAINS=$(forge kb domains 2>/dev/null)
if [ -n "$WIKI_DOMAINS" ]; then
  wiki_context="<forge-wiki>\n可用知识领域: ${WIKI_DOMAINS}\n检索: forge kb search --domains <domains> <query>\n</forge-wiki>"
  # 拼接到 guide_content 后一起注入
fi
```

与 guide.md 注入机制完全一致，通过 `hookSpecificOutput.additionalContext` 输出，用户零配置。

**完整的三层感知链路**：

| 层 | 机制 | Agent 感知成本 |
|---|---|---|
| **发现** | session-start hook 注入 domains | 零（自动） |
| **检索** | `forge kb search --domains` | 一次 tool call，返回轻量元数据 |
| **获取** | `forge kb get --section` | 按需，只拉相关章节 |

### 知识流：Forge Push → Wiki Receive

- **Forge（生产者）**：在知识产生时通过 `forge kb push` 推送到 Wiki。触发点包括：`/learn` 记录决策/教训、`/consolidate-specs` 提取约定、任务完成时记录模式、Bug 修复时捕获根因
- **Wiki（接收者 + 维护者）**：接收 Forge 推送的知识，自动整合——去重、交叉引用、标记矛盾、标注时间戳和时效状态

### Wiki 维护操作

- **Receive**：接收 Forge 推送的知识，整合进 Wiki，更新交叉引用，标记矛盾。这是被动操作，由 Forge 的推送触发
- **Lint**：定期健康检查——矛盾检测、**过期标记**、孤立页面、缺失交叉引用

### 时间维度

每条知识条目携带时间戳和时效标签。Agent 查询时自动感知知识的新鲜度，优先引用近期知识，对过期知识标记警告。知识演变历史可追溯——"为什么我们现在用 gRPC 而不是 GraphQL"。

### 数字孪生积累

随着项目经验不断 push，Wiki 逐步构建出开发者的能力画像——决策偏好、常见模式、踩过的坑、擅长的领域。这是数字孪生的知识底座。

### Innovation Highlights

**Agent 优先的知识架构**：传统知识库为人设计（文件夹、标签、搜索框）。Forge Wiki 为 agent 设计——知识结构、检索接口、返回格式都优化为 agent 可直接消费的形式。人通过 Obsidian 浏览，agent 通过 `forge kb` CLI 命令检索，同一个 Wiki 服务两个消费者。

**时间感知的知识生命周期**：知识不是静态快照，而是有时间维度的活文档。2024-2026 年间对同一问题的不同决策形成演变轨迹。Agent 在建议方案时会说"我们 2024 年用 REST，2025 年试了 GraphQL，2026 年回归 gRPC，原因是复杂度管理"。这让 agent 的建议不是基于最新的文档，而是基于你的完整决策历史。

**LLM Wiki 模式应用于软件开发**：Karpathy 的 LLM Wiki 主要用于研究和个人笔记管理，我们将其应用于软件开发现知识管理——一个更结构化、更可验证的领域，并在此基础上构建数字孪生的知识底座。

**知识与工具的闭环**：Forge 是知识的源头（通过 `/learn`、`/consolidate-specs`、任务执行产出），在知识产生时通过 `forge kb push` 推送到 Wiki。Forge 同时也是知识的消费者（agent 通过 `forge kb search` / `forge kb get` 检索）。形成"产出 → push → 维护 → search → 消费 → 新产出"的正循环。

**AI 维护而非人工维护**：传统知识库的死因是维护成本随规模非线性增长。LLM Wiki 的核心洞察是让 AI 承担所有维护工作——交叉引用、摘要更新、冲突检测、一致性维护。

## Requirements Analysis

### Key Scenarios

1. **知识推送**：任务执行中 `/learn` 记录了一条教训 → `forge kb push` 推送到 Wiki → Wiki 接收并整合（去重、交叉引用、标注时间戳）
2. **Agent 自动发现**：Agent 开启会话时，`session-start` hook 自动注入 `forge kb domains` 的输出。Agent 在上下文中直接看到可用知识领域，零额外操作
3. **两步检索**：Agent 执行 auth 相关任务 → `forge kb search --domains auth "jwt token"` 发现相关文档 → `forge kb get auth-flow --section "刷新机制"` 只拉需要的章节
4. **查询知识演变**：`forge kb search "API 设计风格"` 返回相关文档列表，文档内包含时间线——"2024 REST → 2025 GraphQL → 2026 gRPC"，每次转变的原因可追溯
5. **Wiki 健康检查**：`forge kb lint` 检测矛盾、过期条目（>6个月未验证）、孤立页面、缺失交叉引用
6. **跨项目知识发现**：多项目 `forge kb push` 后，Wiki 自动发现重复出现的模式，提炼为通用知识条目

### Non-Functional Requirements

- Wiki 存储为纯 Markdown 文件，可被任何编辑器打开（Obsidian、VS Code、终端）
- 知识库独立于 Forge 安装目录，用户指定位置
- Agent Query 返回结构化格式（不依赖 LLM 二次解析），可直接注入 agent 上下文
- `forge kb search` 返回轻量元数据（名称 + domains + 描述），Agent 按需通过 `forge kb get --section` 获取具体内容，避免加载全文浪费 token
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
- `forge kb push`：项目知识推送到 Wiki，自动整合（去重、交叉引用、标注来源和时间戳）
- `forge kb search`：两步 agentic search（元数据列表 → 按需获取），支持 `--domains` 过滤
- `forge kb get`：获取文档全文或指定章节（`--outline` / `--section`），节省 token
- `forge kb domains`：列出可用知识领域
- `forge kb lint`：Wiki 健康检查（含过期检测）
- Agent 上下文注入：通过 session-start hook 自动注入 domains 目录，与 guide.md 注入机制一致
- Forge Push 集成：在 `/learn`、`/consolidate-specs`、任务完成等知识产生点触发 `forge kb push`

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

- [ ] `forge kb push` 能将项目知识推送到 Wiki，至少覆盖 3 个项目的知识产出
- [ ] `forge kb search --domains auth "API 设计"` 能返回相关文档的元数据列表（名称 + domains + 描述），支持跨项目检索
- [ ] `forge kb get <name> --section <section>` 能返回指定章节内容，token 开销远低于全文加载
- [ ] Agent 会话启动时自动收到 `forge kb domains` 的输出作为上下文（通过 session-start hook）
- [ ] 任务执行时 agent 能基于注入的 domains 信息主动调用 `forge kb search` 检索相关知识
- [ ] Wiki 通过 `forge kb lint` 检查，无矛盾条目，过期条目被正确标记
- [ ] Wiki 可在 Obsidian 中正常浏览，graph view 可显示知识关联

## Next Steps

- Proceed to `/write-prd` to formalize requirements
