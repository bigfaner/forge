---
created: "2026-05-22"
topic: "Graphify"
mode: "deep-dive"
dimensions: ["概览与定位", "架构与核心概念", "生态与社区", "性能与基准测试"]
---

# Graphify 开源项目研究报告

## Overview

**Graphify 是当前增长最快的 AI 编程助手开源项目之一——2026 年 4 月 3 日发布，不到 50 天即获得 49k+ GitHub Stars（全球排名 #438）。它是一个多模态知识图谱构建工具，让 Claude Code、Codex 等 AI 编程助手能将整个仓库（代码、文档、PDF、图片、视频）转化为可查询的知识图谱，实现平均 71.5 倍的 Token 压缩效率。**

**Research mode:** Deep Dive

**Key question:** 深入了解 Graphify 的定位、架构设计、社区生态和性能表现

---

## Research Background & Objectives

本报告旨在系统性了解 Graphify 开源项目——它是谁、怎么工作的、社区反馈如何、性能表现如何。Graphify 在短短一个多月内从零增长到近 5 万 GitHub Stars，已经成为 AI 编程辅助领域不可忽视的基础设施级项目。

### Research Scope

| 维度 | Value |
|---|---|
| Topic | Graphify (safishamsi/graphify) |
| Mode | Deep Dive |
| Dimensions covered | 概览与定位、架构与核心概念、生态与社区、性能与基准测试 |
| Candidates | — |
| Project adaptation | No |

---

## 概览与定位

### 项目身份

Graphify 是由 **Safi Shamsi** 创建并维护的开源项目，托管于 [github.com/safishamsi/graphify](https://github.com/safishamsi/graphify)。采用 **MIT 许可证**，主要语言为 Python，要求 Python 3.10+。

- **创建时间**: 2026 年 4 月 3 日
- **当前 Stars**: ~49k（全球 GitHub 排名 #438）
- **Forks**: 5.3k
- **贡献者**: 41 人
- **PyPI 包名**: `graphifyy`（双 y，CLI 命令仍为 `graphify`）

### 解决的核心问题

AI 编程助手面临一个根本性矛盾：**上下文窗口限制**。一个完整的代码库加上文档、RFC、论文和图表，无法塞进单个 Prompt 中。传统 RAG 将一切分块后通过 embedding 相似度检索，但这种方法丢失了**结构信息**——谁调用了谁、哪个模块依赖哪个、创建某个函数的提交信息里记录了什么样的设计理由。

Graphify 的解决方案：**用知识图谱替代向量检索**。

- 节点 = 概念（类、函数、设计决策、论文段落、图表）
- 边 = 关系（`calls`、`imports`、`rationale_for`、`semantically_similar_to`）
- 编程助手遍历边而非检索分块

### 定位与差异化

| 对比项目 | 定位 | Graphify 的优势 |
|---------|------|----------------|
| Sourcegraph | 跨仓库代码搜索 | 不是知识图谱，缺少设计语义 |
| Code2Vec | 函数级 embedding | 无图谱结构，不支持多模态输入 |
| Neo4j | 通用图数据库 | 不能从代码自动生成图谱 |

Graphify 既不与通用图数据库竞争，也不与代码搜索工具竞争——它在两者之间创造了新品类：**从代码库自动构建的知识图谱层**，专为 AI 编程助手设计。

### 商业生态

Graphify 项目本身是 MIT 开源免费的。其作者正在开发基于 Graphify 的商业产品 **Penpax**（graphifylabs.ai）——始终在线层，将图谱方法扩展到会议、浏览器历史、邮件、文件和代码，持续在后台更新，全本地化运行，无云依赖。

**Sources:**
- 官方网站 — https://graphify.net/
- GitHub 仓库 — https://github.com/safishamsi/graphify
- Star History — https://www.star-history.com/safishamsi/graphify

---

## 架构与核心概念

### 七阶段流水线

Graphify 是一个多阶段流水线架构。每个阶段是独立的函数/模块，通过普通 Python dict 和 NetworkX 图传递数据，**无共享状态、无副作用**（`graphify-out/` 以外）。

```
detect() → extract() → build_graph() → cluster() → analyze() → report() → export()
```

### 模块职责

| 模块 | 功能 | 输入 → 输出 |
|------|------|-------------|
| `detect.py` | `collect_files(root)` | 目录 → `[Path]` 过滤后的文件列表 |
| `extract.py` | `extract(path)` | 文件路径 → `{nodes, edges}` 字典 |
| `build.py` | `build_graph(extractions)` | 提取结果列表 → `nx.Graph` |
| `cluster.py` | `cluster(G)` | 图 → 带 `community` 属性的图 |
| `analyze.py` | `analyze(G)` | 图 → 分析结果（god nodes、surprises、questions） |
| `report.py` | `render_report(G, analysis)` | 图 + 分析 → GRAPH_REPORT.md 字符串 |
| `export.py` | `export(G, out_dir, ...)` | 图 → Obsidian vault / graph.json / graph.html / graph.svg |

### 提取架构（核心创新）

Graphify 的提取机制分为两路：

**1. 代码文件（本地，零 API 调用）**
- 使用 **Tree-sitter**（确定性解析器，类似编译器前端）解析 31 种编程语言的 AST
- 提取：类、函数、调用关系、导入依赖
- 纯本地执行，代码永不离开机器

**2. 非代码文件（通过 AI 模型 API）**
- 文档/PDF/图片/视频 → 通过已配置的 AI 模型（Claude、Gemini、OpenAI 等）进行语义提取
- 仅发送语义描述（非原始代码）

### 提取输出 Schema

```json
{
  "nodes": [
    {"id": "unique_string", "label": "human name", "source_file": "path", "source_location": "L42"}
  ],
  "edges": [
    {"source": "id_a", "target": "id_b", "relation": "calls|imports|uses|...",
     "confidence": "EXTRACTED|INFERRED|AMBIGUOUS"}
  ]
}
```

### 置信度标签体系

| 标签 | 含义 |
|------|------|
| `EXTRACTED` | 关系在源码中明确声明（如 import 语句、直接调用） |
| `INFERRED` | 关系是合理推断（如调用图二次传递、同现上下文） |
| `AMBIGUOUS` | 关系不确定，在 GRAPH_REPORT.md 中标记供人工审查 |

### 聚类算法

使用 **Leiden 算法**进行语义社区检测——不需要向量 embedding。Leiden 是 Louvain 算法的改进版，能发现更精细的社区结构。

### God Nodes 与 Surprising Connections

- **God Nodes**：度数最高的节点——系统中最核心的概念
- **Surprising Connections**：跨文件或跨域的意外连接，按意外程度排序
- 例如在 httpx 示例中，`DigestAuth → Response` 被标记为意外边

### MCP 服务支持

Graphify 支持以 **MCP（Model Context Protocol）** 标准协议暴露图谱，提供结构化查询：`query_graph`、`get_node`、`get_neighbors`、`shortest_path` 等。

### 安全设计

- URL 仅允许 http/https
- 文件下载有大小限制和超时
- 输出路径进行 containment 校验
- 节点标签进行 HTML 转义
- 无遥测，无使用追踪

### 支持的文件类型

| 类型 | 格式 |
|------|------|
| 代码（31 种语言） | `.py .ts .js .jsx .tsx .go .rs .java .c .cpp .rb .cs .kt .scala .php .swift .lua .zig .ps1 .ex .exs .m .mm .jl .vue .svelte .astro .groovy .dart .v .sv .sql .f90 .pas .sh .json` 等 |
| 文档 | `.md .mdx .html .txt .rst .yaml .yml` |
| Office（可选） | `.docx .xlsx` |
| PDF | `.pdf` |
| 图片 | `.png .jpg .webp .gif` |
| 视频/音频（可选） | `.mp4 .mov .mp3 .wav` 及 YouTube URL |
| Google Workspace（可选） | `.gdoc .gsheet .gslides` |

**Sources:**
- ARCHITECTURE.md — https://github.com/safishamsi/graphify/blob/master/ARCHITECTURE.md
- 官方架构说明 — https://graphify.net/
- 安全审计 Issue — https://github.com/safishamsi/graphify/issues/63

---

## 生态与社区

### GitHub 增长数据（截至 2026 年 5 月 22 日）

| 指标 | 数值 |
|------|------|
| Stars | ~49,000 |
| 全球排名 | #438 |
| Forks | ~5,300 |
| 贡献者 | 41 人 |
| 周增速 | +189 stars / 周 |
| 创建时间 | 2026 年 4 月 3 日 |
| 年龄 | ~50 天 |

### 增长轨迹

根据创作者在 Reddit 上的分享（26 天时）：**450k+ PyPI 下载量、~40k GitHub Stars**。这是一个极其陡峭的增长曲线——几乎是最快达到 40k Stars 的 Python 项目之一。

### 平台支持生态

Graphify 已适配 16+ AI 编程助手平台：

Claude Code、OpenAI Codex、OpenCode、Cursor、Gemini CLI、GitHub Copilot CLI、VS Code Copilot Chat、Aider、OpenClaw、Factory Droid、Trae、Trae CN、Hermes、Kimi Code、Kiro、Pi、Google Antigravity

每个平台都有专门的 skill 文件（`skill-*.md`），也可通过 CLI 一键注册。

### 可选功能扩展生态

| Extra | 用途 |
|-------|------|
| `pdf` | PDF 提取 |
| `office` | .docx / .xlsx 支持 |
| `video` | 视频/音频转录（faster-whisper + yt-dlp） |
| `mcp` | MCP 协议服务 |
| `neo4j` | Neo4j 导出 |
| `leiden` | Leiden 社区检测 |
| `ollama` | Ollama 本地推理 |
| `openai` | OpenAI API |
| `gemini` | Google Gemini API |
| `bedrock` | AWS Bedrock（IAM 认证） |
| `svg` | SVG 图导出 |

### 国际化

README 已翻译为 **26 种语言**，包括简体中文、日语、韩语、德语、法语、西班牙语、阿拉伯语等。

### 社区文章与媒体

已有多篇技术博客、YouTube 视频、LinkedIn 帖子覆盖 Graphify。知名内容包括：
- Analytics Vidhya（2026 年 4 月）的深度指南
- Medium 上的多篇教程和评测
- GopenAI 博客的技术解析
- YouTube 上的使用教程

### 商业支撑

- **The Memory Layer** — Safi Shamsi 撰写的 GitBook 书籍
- **Penpax** — 基于 Graphify 的商业产品（graphifylabs.ai）
- **GitHub Sponsors** — 已开通赞助通道

**Sources:**
- Star History — https://www.star-history.com/safishamsi/graphify
- Reddit (创作者分享) — https://www.reddit.com/r/ClaudeAI/comments/1t18eeh/
- PyPI Downloads Dashboard — https://clickpy.clickhouse.com/dashboard/graphifyy
- Analytics Vidhya 文章 — https://www.analyticsvidhya.com/blog/2026/04/graphify-guide/

---

## 性能与基准测试

### 核心指标：Token 压缩比

Graphify 的核心价值主张：**在保持结构信息的前提下，极大减少 AI 编程助手理解代码库所需的 Token 数量**。

官方基准测试结果（Karpathy 混合语料，52 个文件，~92k 词）：

| 指标 | 数值 |
|------|------|
| 原始语料大小 | ~92k 词 → ~123k tokens（原始方式） |
| 平均查询 Token 消耗 | ~1,700 tokens（使用 Graphify） |
| **压缩比** | **71.5×** |
| 图谱规模 | 285 节点, 340 边, 53 个社区 |

### 不同规模下的表现

| 语料 | 文件数 | 节点 | 边 | 社区 | 说明 |
|------|--------|------|-----|------|------|
| httpx | 6 个 Python 文件 | 144 | 330 | 6 | 小型 HTTP 传输层库 |
| Karpathy 混合 | ~52 文件（~92k 词） | 285 | 340 | 53 | 代码 + 论文 + 图表 |
| 大规模（500k 词） | — | — | — | — | BFS 子图查询 ~2k tokens vs ~670k 原始 |

### 扩展性

- Tree-sitter 解析和 NetworkX 构建随代码规模**线性扩展**
- 在大规模语料（~500k 词）上，BFS 子图查询依然保持在 ~2k tokens（vs ~670k 原始）
- 支持 `--max-workers` 控制 AST 并行线程数
- 支持增量更新（`--update` 仅重新提取变更文件）

### 与替代方案对比

| 维度 | Graphify | 传统 RAG（Vector Store） |
|------|----------|------------------------|
| 结构保留 | 是（图结构天然保留关系） | 否（chunk 打散后丢失结构） |
| 多模态 | 原生支持（代码+文档+PDF+图片+视频） | 通常文本 only |
| 置信度标记 | 每边标记 EXTRACTED/INFERRED/AMBIGUOUS | 无 |
| 跨文件关系 | 自动识别 | 需要额外处理 |
| 查询 Token 效率 | ~1.7k 查询 vs ~123k 原始（71.5×） | 取决于 chunk 大小 |
| 初始构建成本 | 一次性（随文件量线性） | 持续 embedding |
| 代码隐私 | 代码文件纯本地处理 | 取决于部署方式 |

### 实测示例

在 httpx 工作示例中：
- 图谱生成结果：144 节点、330 边、6 个社区
- God Nodes: `Client`、`AsyncClient`、`Response`、`Request`
- Surprise Edge: `DigestAuth → Response`

**Sources:**
- 基准测试模块 — https://github.com/safishamsi/graphify/blob/master/graphify/benchmark.py
- 官方文档 How It Works — https://graphify.net/
- 工作示例 — https://github.com/safishamsi/graphify/tree/master/worked

---

## Risks & Caveats

| Risk | Severity | Mitigation |
|------|----------|------------|
| 项目极新（不足 2 个月），长期维护性未知 | Medium | MIT 许可证，代码可 fork；已有 41 位贡献者 |
| LLM 语义提取依赖外部 API（非代码文件需用 API Key） | Low | 仅发语义描述不发原始代码；支持 Ollama 本地推理彻底避免外传 |
| 代码隐私风险（语义提取阶段） | Low | 代码文件由 Tree-sitter 本地处理，从不离开机器 |
| 与 Claude Code 深度绑定（核心体验最优） | Medium | 已适配 16 个平台，但 Claude Code 体验最好 |
| Leiden 算法支持受限（Python < 3.13 才支持） | Low | 不影响基础功能，仅影响社区检测精度 |
| 大型代码库的 HTML 可视化可能过重 | Low | 可通过 `--no-viz` 跳过，直接使用 JSON |
| 商业可持续性依赖 Penpax 产品销售 | Medium | 核心功能 MIT 开源，商业补充不影响开源版本 |

---

## Recommendation

Graphify 是 **2026 年最值得关注的 AI 编程辅助基础设施项目之一**。它在 50 天内从零增长到 49k Stars，解决了 AI 编程助手理解代码库时的根本性痛点——上下文窗口限制。

**核心优势：**
- 结构优先而非相似度优先的代码理解方式
- 71.5× Token 压缩效率极为可观
- 纯本地处理代码文件，隐私友好
- 模块化架构，易于扩展和定制
- 16+ 平台兼容，生态正在快速形成

**适用场景：**
- 使用 AI 编程助手（Claude Code、Codex 等）处理中大型代码库的开发者
- 需要让 AI 助手快速理解跨文件、跨模块架构关系的团队
- 希望在保持代码隐私的同时获得深度代码理解的场景

**Confidence level:** High — 项目质量、高速增长、架构设计都有充分公开证据支撑

---

## Sources

| Source | URL | Used for |
|--------|-----|----------|
| 官方网站 | https://graphify.net/ | 概览、架构、核心特性 |
| GitHub 仓库 | https://github.com/safishamsi/graphify | 项目结构、代码、文档 |
| ARCHITECTURE.md | https://github.com/safishamsi/graphify/blob/master/ARCHITECTURE.md | 架构与模块职责 |
| Star History | https://www.star-history.com/safishamsi/graphify | GitHub 增长数据、排名 |
| Reddit 创作者分享 | https://www.reddit.com/r/ClaudeAI/comments/1t18eeh/ | 社区增长、下载量数据 |
| PyPI Downloads | https://clickpy.clickhouse.com/dashboard/graphifyy | 下载量统计 |
| Analytics Vidhya 文章 | https://www.analyticsvidhya.com/blog/2026/04/graphify-guide/ | 第三方评测 |
| Knowledge Graphs for AI | https://graphify.net/knowledge-graph-for-ai-coding-assistants.html | 知识图谱技术原理 |
| GitHub Issue #63 | https://github.com/safishamsi/graphify/issues/63 | 安全审查 |
| Benchmark 模块 | https://github.com/safishamsi/graphify/blob/master/graphify/benchmark.py | 性能基准测试代码 |
| Knightli 技术博客 | https://www.knightli.com/en/2026/05/21/safishamsi-graphify-ai-code-knowledge-graph/ | 第三方技术分析 |
| GopenAI Blog | https://blog.gopenai.com/graphify-build-a-knowledge-graph-from-your-entire-codebase-without-sending-your-code-to-anyone-1b6924474b50 | 技术深度解析 |