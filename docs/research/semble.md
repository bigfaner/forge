---
created: "2026-05-24"
topic: "Semble"
mode: "deep-dive"
dimensions: [overview, architecture, performance]
---

# Semble Research Report

## Overview

Semble 是一个专为 AI Agent 设计的本地代码搜索库，通过静态嵌入（Model2Vec）+ BM25 + 代码感知重排序的混合架构，在 CPU 上实现毫秒级检索，达到 137M 参数 Transformer 模型 99% 的检索质量，同时减少 ~98% 的 token 消耗。技术上设计精巧、工程完成度高，是当前 Agent 代码搜索领域的标杆实现。

**Research mode:** Deep Dive

**Key question:** Semble 的技术原理是什么？其架构设计和性能表现如何？

## Research Background & Objectives

本调研旨在深入理解 Semble 的核心技术原理，包括其混合检索架构、静态嵌入模型、代码感知重排序机制，以及在真实场景中的性能表现和社区反馈。调研不结合当前项目，纯技术研究。

### Research Scope

| Dimension | Value |
|---|---|
| Topic | Semble (MinishLab/semble) |
| Mode | Deep Dive |
| Dimensions covered | Overview, Architecture, Performance |
| Project adaptation | No |

---

## Overview

### 项目定位

Semble 是 MinishLab 开发的面向 AI 编码 Agent 的代码搜索库。核心价值主张：

- **速度**：平均仓库索引 ~250ms，查询 ~1.5ms，全部在 CPU 上运行
- **准确**：NDCG@10 达到 0.854，与 137M 参数的 CodeRankEmbed Hybrid（0.862）几乎持平
- **省 token**：相比 Agent 典型的 grep+read 工作流，减少 ~98% token 消耗
- **零配置**：无需 API Key、GPU 或外部服务

### 核心团队

- **Thomas van Dongen** 和 **Stephan Tulkens**（MinishLab）
- 同一团队也是 [Model2Vec](https://github.com/MinishLab/model2vec) 的作者，在静态嵌入领域有深厚积累
- 2026 年发布，MIT 协议，已发布 Zenodo DOI

### 集成方式

Semble 提供三种使用方式：

1. **MCP Server**：一行命令接入 Claude Code、Cursor、Codex、OpenCode、VS Code、GitHub Copilot CLI、Windsurf、Gemini CLI、Kiro、Zed 等几乎所有主流 Agent
2. **CLI / AGENTS.md**：通过 Bash 调用，适用于子 Agent（子 Agent 无法直接调用 MCP 工具）
3. **Python Library**：`from semble import SembleIndex` 直接编程调用

**Sources:**
- GitHub README — https://github.com/MinishLab/semble
- Official Docs — https://minish.ai/packages/semble/introduction/

---

## Architecture

### 系统架构总览

```
输入查询 (NL 或 Symbol)
       ↓
┌──────────────────────────────────────┐
│           Index 阶段 (离线)            │
│                                      │
│  文件 → tree-sitter 分块 (Chonkie)    │
│    ↓                ↓                │
│  BM25 索引     静态嵌入 (potion-code-16M) │
│  (词法匹配)     (语义相似度, 无前向传播)   │
└──────────────────────────────────────┘
       ↓
┌──────────────────────────────────────┐
│           Search 阶段 (在线)           │
│                                      │
│  BM25 检索  +  语义检索                │
│       ↓            ↓                 │
│     RRF (Reciprocal Rank Fusion)     │
│       ↓                              │
│  α 加权融合 (自适应查询类型)            │
│       ↓                              │
│  代码感知重排序 (5 组信号)              │
│       ↓                              │
│  返回 top-k chunks                   │
└──────────────────────────────────────┘
```

### 源码结构

```
src/semble/
├── chunking/          # 代码分块 (tree-sitter via Chonkie)
├── index/
│   ├── create.py      # 索引创建流程
│   ├── dense.py       # 语义向量索引 (Vicinity)
│   ├── sparse.py      # BM25 索引 (bm25s)
│   ├── files.py       # 文件类型过滤
│   └── file_walker.py # 文件遍历
├── ranking/
│   ├── boosting.py    # 定义提升、stem 匹配、文件连贯性
│   ├── penalties.py   # 噪声惩罚 (测试文件、compat 等)
│   └── weighting.py   # α 权重自适应
├── search.py          # 核心搜索逻辑 (hybrid/semantic/bm25)
├── mcp.py             # MCP Server
├── cli.py             # CLI 入口
└── types.py           # 类型定义
```

### 关键技术细节

#### 1. 代码分块 (Chunking)

使用 Chonkie 库（基于 tree-sitter）将源文件按语法结构分块，而非按行数切割。这意味着每个 chunk 是完整的语法单元（函数、类、方法等），上下文边界更合理。

支持 tree-sitter-language-pack 覆盖的所有语言（19+ 种），并通过 `ContentType` 枚举区分：
- `CODE`：源代码（默认）
- `DOCS`：文档（markdown, rst 等）
- `CONFIG`：配置文件（yaml, toml, terraform 等）

#### 2. 双路检索 (Dual Retrieval)

**语义检索 (Semantic)**：
- 模型：`potion-code-16M`（16M 参数的代码专用静态嵌入模型）
- 技术：Model2Vec — 将 sentence transformer 蒸馏为静态查找表，**无 transformer 前向传播**
- 运行时只需一次查表操作，所以速度极快
- 向量索引使用 Vicinity 库（余弦距离）

**词法检索 (BM25)**：
- 使用 `bm25s` 库（纯 NumPy 实现，无 Cython 依赖）
- 对标识符和 API 名称做精确匹配
- 自定义 tokenizer 处理驼峰/下划线分割

**融合 (RRF)**：
```python
# RRF 公式: score = 1 / (k + rank), k=60
# 语义和 BM25 各自排位后转换为 RRF 分数
# 然后按 α 权重线性组合
combined = α × semantic_rrf + (1-α) × bm25_rrf
```

#### 3. 自适应权重 (Adaptive Weighting)

```python
# 查询类型检测
if is_symbol_query(query):  # 如 "Foo::bar", "getUserById", "_private"
    α = 0.3  # 偏向 BM25（精确匹配）
else:                        # 自然语言查询，如 "authentication flow"
    α = 0.5  # 语义和词法均衡
```

符号查询检测使用正则匹配：命名空间限定符（`::`、`\`、`->`、`.`）、前导下划线、包含大写字母/下划线的标识符。

#### 4. 代码感知重排序 (Code-Aware Reranking)

5 组重排序信号：

| 信号 | 机制 | 代码位置 |
|---|---|---|
| **Definition Boost** | 包含 `class`/`def`/`func` 等定义关键字的 chunk 排名提升，支持 20+ 种定义关键字和 SQL DDL | `boosting.py: _boost_symbol_definitions` |
| **Identifier Stems** | 查询词的 stem 匹配 chunk 中的标识符（`parse config` → `parseConfig`/`ConfigParser`/`config_parser`） | `boosting.py: _boost_stem_matches` |
| **File Coherence** | 同一文件多个 chunk 命中时，提升该文件最高排名的 chunk（20% 分位加成） | `boosting.py: boost_multi_chunk_files` |
| **Noise Penalties** | 测试文件、`compat/`/`legacy/` shim、示例代码、`.d.ts` 声明文件降权 | `penalties.py` |
| **Embedded Symbols** | NL 查询中嵌入的驼峰标识符（如 "how does StateManager work"）半额提升 | `boosting.py: _boost_embedded_symbols` |

#### 5. 索引与缓存

- 本地路径：创建索引后可序列化到磁盘（`semble index -o <path>`），后续查询复用
- Git URL：按需 `git clone --depth 1` 到临时目录，索引内容保留在内存
- MCP 模式下，索引按 session 生命周期缓存；本地路径自动监听文件变更并重新索引

**Sources:**
- `src/semble/search.py` — https://github.com/MinishLab/semble/blob/main/src/semble/search.py
- `src/semble/ranking/boosting.py` — https://github.com/MinishLab/semble/blob/main/src/semble/ranking/boosting.py
- `src/semble/ranking/weighting.py` — https://github.com/MinishLab/semble/blob/main/src/semble/ranking/weighting.py
- `src/semble/index/index.py` — https://github.com/MinishLab/semble/blob/main/src/semble/index/index.py
- potion-code-16M — https://huggingface.co/minishlab/potion-code-16M

---

## Performance

### Benchmark 方法论

- **数据集**：~1,250 查询 × 63 仓库 × 19 语言
- **查询类别**：语义 (711)、架构 (343)、符号 (204)
- **标注**：Claude Sonnet 4.6 生成查询和相关性标注，同模型做质量校验
- **硬件**：全部 CPU 运行

### 质量与速度

| 方法 | NDCG@10 | 索引时间 | 查询 p50 |
|---|---|---|---|
| CodeRankEmbed Hybrid (137M) | 0.862 | 57s | 16ms |
| **semble** | **0.854** | **263ms** | **1.5ms** |
| CodeRankEmbed (dense only) | 0.765 | 57s | 16ms |
| ColGREP | 0.693 | 5.8s | 124ms |
| BM25 only | 0.673 | 263ms | 0.02ms |
| grepai | 0.561 | 35s | 48ms |
| probe | 0.387 | — | 207ms |
| ripgrep | 0.126 | — | 12ms |

**关键对比**：semble 达到 CodeRankEmbed Hybrid 99.1% 的检索质量，但索引快 218x、查询快 10.7x。

### Token 效率

| 方法 | 每查询预期 token | 节省 |
|---|---|---|
| ripgrep + read file | 45,692 | baseline |
| **semble** | **566** | **98.8%** |

召回率 vs token 预算：

| Token 预算 | semble | ripgrep+read |
|---|---|---|
| 500 | 0.685 | 0.001 |
| 1k | 0.849 | 0.008 |
| 2k | **0.938** | 0.037 |
| 4k | 0.976 | 0.088 |
| 8k | 0.991 | 0.212 |

semble 仅用 2k token 就达到 94% 召回率，而 ripgrep+read 即使在 32k token 也只有 58%。

### 消融实验

| 检索组合 | Raw (无重排序) | + semble 重排序 |
|---|---|---|
| BM25 only | 0.675 | 0.834 (+0.159) |
| potion-code-16M only | 0.650 | 0.821 (+0.171) |
| **BM25 + potion-code-16M** | — | **0.854** |

重排序栈带来 +0.16~0.17 的 NDCG 提升，是 semble 竞争力的关键。

按查询类别分析：

| 模式 | 架构查询 | 语义查询 | 符号查询 |
|---|---|---|---|
| semble hybrid | **0.802** | **0.846** | **0.958** |

符号查询（函数/类/变量名查找）是 semble 最强的场景，NDCG 达到 0.958。

### 按语言表现

semble 在 19 种语言中有 6 种取得最高分（C++、Ruby、Zig、Go、Java、Rust），其余多数与 CodeRankEmbed Hybrid 差距在 0.01-0.03 以内。TypeScript 是最大弱项（0.706 vs 0.708，差距极小）。

### 真实场景社区反馈

来自 Hacker News 讨论的实测反馈：

- **正面**：用户报告在大型代码库（如 rust-lang/rust）上 semble 26 秒完成索引 vs 竞品数小时；在 84K 行 C 项目上 indexing 几乎无感
- **谨慎**：多位用户指出 token 节省数字虽然真实，但检索质量 ≠ Agent 端到端性能；有些 Agent 会不信任搜索结果而重复 grep
- **已知问题**：MCP 连接偶发失败（`Connection closed` 错误）；Codex CLI 通过 MCP 调用时会挂起；不同模型对搜索结果的信任度不同（Sonnet 4.6 > Opus 4.7）
- **作者回应**：承认端到端 Agent 评测是下一步工作，目前 benchmark 只测量检索质量

**Sources:**
- Benchmarks README — https://github.com/MinishLab/semble/blob/main/benchmarks/README.md
- Hacker News Discussion — https://news.ycombinator.com/item?id=48169874
- Reddit r/ClaudeAI — https://www.reddit.com/r/ClaudeAI/comments/1szvo7t/

---

## Risks & Caveats

| Risk | Severity | Mitigation |
|---|---|---|
| 端到端 Agent 效果未验证 — benchmark 仅测量检索质量，未证明 Agent 整体效率提升 | High | 作者已承认并计划做 Agent-level eval；社区实测结果波动大（25k-95k token），取决于模型和任务 |
| Agent 不信任搜索结果 — 部分模型（尤其 Opus 4.7）被 RL 训练偏好 grep，会用 semble 后仍重复 grep | Medium | 通过优化 prompt 和工具描述可缓解；作者正在研究中 |
| MCP 连接稳定性 — Codex CLI 挂起、偶发 Connection closed 错误 | Medium | 可改用 CLI/AGENTS.md 模式绕过 MCP；项目活跃维护中 |
| 静态嵌入的语义上限 — potion-code-16M 是静态模型，对复杂语义查询的理解能力有天花板 | Low | 通过 BM25 + 重排序补偿；实际 benchmark 表现接近大模型 |
| Python 实现 — 非 Rust/Go，极端场景下可能成为瓶颈 | Low | 实测已经足够快（263ms 索引），Python 生态对 ML 模型加载更友好 |
| 供应链安全 — transitive dependency 理论上可能被攻击 | Low | 使用 PyPI trusted publishing + uv `--exclude-newer`，最小化依赖 |

---

## Recommendation

Semble 是一个工程完成度很高的代码搜索库，其「静态嵌入 + BM25 + 代码感知重排序」的混合架构设计精巧，在不依赖 GPU/API 的前提下达到了接近大型 Transformer 模型的检索质量。对于理解 Agent 代码搜索的技术方案而言，这是一个值得深入学习的参考实现。

需要注意：它的 benchmark 测量的是检索质量而非端到端 Agent 性能提升。在实际 Agent 集成中，搜索结果的质量并不直接转化为 Agent 工作效率的提升——这取决于模型是否信任搜索结果、prompt 设计、以及任务类型。

**Confidence level:** High — 基于官方文档、源码分析、benchmark 数据和 HN 社区讨论的交叉验证。

---

## Sources

| Source | URL | Used for |
|---|---|---|
| GitHub README | https://github.com/MinishLab/semble | Overview, Architecture, Integration |
| Official Docs | https://minish.ai/packages/semble/introduction/ | Overview, How it works |
| Benchmarks README | https://github.com/MinishLab/semble/blob/main/benchmarks/README.md | Performance, Ablations, Per-language |
| source: search.py | https://github.com/MinishLab/semble/blob/main/src/semble/search.py | RRF fusion, hybrid search |
| source: ranking/boosting.py | https://github.com/MinishLab/semble/blob/main/src/semble/ranking/boosting.py | Code-aware reranking signals |
| source: ranking/weighting.py | https://github.com/MinishLab/semble/blob/main/src/semble/ranking/weighting.py | Adaptive α weighting |
| source: index/index.py | https://github.com/MinishLab/semble/blob/main/src/semble/index/index.py | SembleIndex API, from_path/from_git |
| potion-code-16M | https://huggingface.co/minishlab/potion-code-16M | Static embedding model |
| Hacker News Discussion | https://news.ycombinator.com/item?id=48169874 | Community feedback, real-world usage |
| Reddit r/ClaudeAI | https://www.reddit.com/r/ClaudeAI/comments/1szvo7t/ | Community feedback |
| OSSInsight | https://ossinsight.io/analyze/MinishLab/semble | Repository analytics |
