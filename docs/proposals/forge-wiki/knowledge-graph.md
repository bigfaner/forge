# Forge Wiki 知识图谱

## 什么是知识图谱

知识图谱是一种用图结构表示知识的方式。节点（Node）代表实体或概念，边（Edge）代表它们之间的关系。

```
文档 A ──共享 concern──► concern: auth ◄──共享 concern── 文档 B
文档 A ──替代──────────────────────────────────────────► 文档 C
文档 B ──引用──────────────────────────────────────────► 文档 D
```

传统搜索是"关键词 → 文档列表"的线性检索。知识图谱搜索是"实体 → 关系 → 关联实体"的拓扑检索，能发现关键词搜索无法覆盖的隐含关联。

## 为什么 Forge Wiki 需要知识图谱

### 当前问题

`forge kb search` 基于多维度过滤（`--type`、`--concerns`、`--artifacts`、`--langs`、`--frameworks`）和关键词匹配，能回答"有哪些关于 auth 的知识"。但无法回答：

- "auth 的 JWT 方案为什么从 session 切换过来？" — 需要决策链追溯
- "auth 和 error-handling 之间有什么关联？" — 需要跨域发现
- "这条知识有没有被更新的知识替代？" — 需要生命周期追踪
- "哪些知识是孤立的，没有任何关联？" — 需要结构分析

这些问题的答案藏在文档间的关系中，不在单个文档内。

### 图谱 vs 关键词搜索的能力边界

| 能力 | 关键词搜索 | 知识图谱 |
|------|-----------|---------|
| 找到包含"JWT"的文档 | ✅ | ✅ |
| 找到与 auth 相关但没提 auth 的文档 | ❌ | ✅（通过共享 concern 的邻居） |
| 追溯一个决策的演变历史 | ❌ | ✅（通过 superseded_by 链） |
| 发现两个不相关领域之间的隐含联系 | ❌ | ✅（通过多跳遍历） |
| 识别孤立、无关联的知识 | ❌ | ✅（通过度分析） |

## Forge Wiki 的图谱模型

### 节点类型

Forge Wiki 知识图谱中有两类节点：

**文档节点（Document Node）**

每篇知识文档是一个节点，由 frontmatter `id` 标识。

```yaml
id: auth-jwt-decision
title: "JWT Token 认证决策"
description: "JWT Token 认证方案选型及理由"
type: product
concerns: [auth]
artifacts: [api, web-spa]
langs: [go]
frameworks: [gin]
source_type: decision
status: active
```

节点属性：

| 属性 | 来源 | 用途 |
|------|------|------|
| `id` | frontmatter | 节点唯一标识 |
| `title` | frontmatter | 显示名称 |
| `description` | frontmatter | 搜索结果预览 |
| `type` | frontmatter | 知识层级（language/design/product/business） |
| `concerns` | frontmatter | 与领域节点的连边 |
| `artifacts` | frontmatter | 适用的软件形态 |
| `langs` | frontmatter | 编程语言（仅 language 层级） |
| `frameworks` | frontmatter | 框架（仅 language 层级，可选） |
| `source_type` | frontmatter | 节点类型（decision/lesson/convention/pattern/anti-pattern/reference） |
| `status` | frontmatter | 生命周期状态 |
| `scope` | frontmatter | universal = 匹配一切搜索，不受维度限制 |
| `created` | frontmatter | 时间戳 |
| `source_project` | frontmatter | 来源项目 |

**领域节点（Concern Node）**

每个 `concerns` 值是一个领域节点。

```
concern:auth  concern:jwt  concern:api  concern:error-handling
```

领域节点没有 frontmatter，由文档的 `concerns` 字段隐式定义。`forge kb concerns` 命令输出的是所有领域节点及其连接的文档数量。

### 边类型与构图来源

图谱的边（文档间的关系）从多个来源构建，每个来源产生不同类型的边：

#### 1. 领域共现边（Concern Co-occurrence）

**来源**：frontmatter `concerns` 字段

**机制**：共享同一 concern 的文档之间产生一条边，权重为共享 concern 的数量。

```
文档 A [auth, jwt, api] ──权重3── 文档 B [auth, jwt, api]
文档 A [auth, jwt, api] ──权重1── 文档 C [auth, error-handling]
文档 B [auth, jwt, api] ──权重1── 文档 C [auth, error-handling]
```

**实现**：二部图投影。文档和 concern 构成二部图 G(V_doc, V_concern, E)，投影到文档空间后，共享 concern 的文档对 (A, B) 的边权重 = |concerns(A) ∩ concerns(B)|。

**v1 状态**：已通过 `--concerns` 过滤隐式实现。`forge kb search --concerns auth` 等价于从 concern:auth 节点出发找所有相连文档。

**v2 优化**：可计算领域共现权重用于搜索排序，权重高的文档对在搜索结果中互相提升排名（参考 llm_wiki 的 Source Overlap 信号，权重 ×4.0）。

#### 2. 决策链边（Decision Chain）

**来源**：frontmatter `superseded_by` 字段

**机制**：有向边，从被替代的文档指向替代它的文档。

```
auth-session-decision ──superseded_by──► auth-jwt-decision ──superseded_by──► auth-oauth-decision
(2024 session)                          (2025 jwt)                            (2026 oauth)
```

**用途**：
- 时间线查询：追溯某个领域决策的演变历史
- 矛盾检测：同主题的活跃文档之间不应有矛盾
- 时效标记：被替代的文档自动标记为 `superseded`

**v1 状态**：frontmatter 已定义 `superseded_by` 字段，`forge kb lint` 可检测断裂的链。

**v2 优化**：构建有向无环图（DAG），支持链式遍历查询。Agent 可问"auth 认证方案经历了哪些变化"。

#### 3. 内容引用边（Content Reference / Wikilinks）

**来源**：文档正文中的 `[[wikilinks]]`

**机制**：文档 A 的正文中出现 `[[auth-jwt-decision]]`，产生一条从 A 到 auth-jwt-decision 的边。

```
lesson-auth-timeout ──引用──► auth-jwt-decision
convention-error-handling ──引用──► auth-jwt-decision
```

**用途**：
- 显式关联：作者/AI 明确标记的文档间关系
- 导航：在 Obsidian 中点击 wikilink 跳转
- 图谱可视化：Obsidian graph view 的核心

**v1 状态**：Obsidian 原生支持 wikilinks，无需额外实现。

**v2 优化**：解析 wikilinks 构建显式引用图，与领域共现图合并为完整知识图谱。

#### 4. 同源项目边（Source Project Clustering）

**来源**：frontmatter `source_project` 字段

**机制**：来自同一项目的文档之间产生弱关联边。

```
文档 A (project: my-backend) ──同源── 文档 B (project: my-backend)
文档 A (project: my-backend) ──无──  文档 C (project: my-frontend)
```

**用途**：
- 项目级知识群组：了解某个项目产生了哪些知识
- 跨项目发现：不同项目在同一领域下的知识对比

**权重**：低于领域共现（同项目不意味着知识相关），可设为 0.5。

**v1 状态**：frontmatter 已定义 `source_project` 字段。

**v2 优化**：作为图谱的辅助信号，在领域共现的基础上微调相关性。

#### 5. 时间邻近边（Temporal Proximity）

**来源**：frontmatter `created` 字段

**机制**：同一 concern 下时间相邻的文档产生弱关联，表示"同一时期的思考"。

```
auth-jwt-decision (2025-03) ──时间邻近── auth-api-gateway (2025-04)
```

**用途**：
- 理解知识产生的时代背景
- 发现同一时期多个相关决策之间的关联

**权重**：时间衰减，越近越强。公式参考：`w = e^(-Δt/τ)`，τ 为半衰期（如 180 天）。

**v1 状态**：frontmatter 已定义 `created` 字段。

**v2 优化**：作为图谱的辅助信号，补充领域共现无法捕获的时序关联。

### 完整图谱模型

```
                    concern:auth
                   ╱     |      ╲
    co-occurrence ╱      |       ╲ co-occurrence
               ╱        |        ╲
    doc-A ────────── doc-B ────── doc-C
     │    co-occ     │              │
     │               │              │
  superseded_by    wikilink      source_project
     │               │              │
     ▼               ▼              ▼
    doc-D          doc-E          doc-F
```

**边类型汇总**：

| 边类型 | 方向 | 权重建议 | 构图来源 | v1 | v2 |
|--------|------|---------|---------|----|----|
| 领域共现 | 无向 | 共享 concern 数量 | `concerns` | 隐式（`--concerns` 过滤） | 显式图 |
| 决策链 | 有向（旧→新） | 1.0（确定性） | `superseded_by` | frontmatter 预留 | DAG 遍历 |
| 内容引用 | 有向 | 1.0（确定性） | `[[wikilinks]]` | Obsidian 原生 | 显式图 |
| 同源项目 | 无向 | 0.5（弱关联） | `source_project` | frontmatter 预留 | 聚类辅助 |
| 时间邻近 | 无向 | e^(-Δt/τ) | `created` | frontmatter 预留 | 时序辅助 |

### scope: universal 在图谱中的行为

`scope: universal` 的文档在图谱中参与构图（可与其他文档产生 concern 共现边、wikilink 边等），但搜索时不受任何维度过滤限制。在 v2 图谱扩展中，universal 文档作为"通用知识锚点"，连接不同领域的知识。

## 图谱在搜索中的应用

### v1：隐式图谱（当前设计）

v1 不构建显式图结构。多维度过滤等价于从 concern 节点出发的一跳查询：

```
forge kb search --concerns auth --type language --langs go "token 刷新"
         │
         ▼
    concern:auth ──1跳──► [doc-A, doc-B, doc-C]
                              │
                         --type/--langs 过滤
                              │
                              ▼
                         [doc-B, doc-C]
                              │
                         关键词 "token 刷新" 过滤
                              │
                              ▼
                         [doc-B]  ← 返回元数据
```

Agent 自主决定是否深入 doc-B。

### v2：显式图谱搜索

v2 构建显式图后，搜索流程升级为：

```
forge kb search --concerns auth --type language --langs go "token 刷新"
         │
         ▼
    Step 1: concerns 过滤 → 候选集 [doc-A, doc-B, doc-C]
         │
         ▼
    Step 2: 关键词匹配 → 命中 [doc-B]
         │
         ▼
    Step 3: 图谱扩展 → doc-B 的关联文档
         │        ├── 领域共现: doc-A (权重3), doc-C (权重1)
         │        ├── 决策链: doc-D (superseded_by)
         │        └── 内容引用: doc-E (wikilink)
         │
         ▼
    Step 4: 相关性排序 → 返回 [doc-B, doc-A, doc-D]
```

**与 llm_wiki / graphify 的对比**：

| 维度 | llm_wiki | graphify | Forge Wiki v2 |
|------|---------|----------|---------------|
| 扩展算法 | 4 信号模型（阈值 2.0，最多扩 3 邻居） | BFS/DFS（depth 2-3） | 领域共现 + 决策链 + 引用，按权重排序 |
| 控制权 | 系统决定扩展哪些 | 系统决定遍历深度 | Agent 决定是否查看扩展结果 |
| 结果注入 | 直接注入 LLM 上下文 | 渲染为文本截断 | 返回元数据，Agent 自主判断 |

### 搜索排序信号（v2 参考）

参考 llm_wiki 的 4 信号模型，为 Forge Wiki v2 设计适配的信号：

| 信号 | 计算方式 | 权重 | 灵感来源 |
|------|---------|------|---------|
| Concerns Overlap | 共享 concerns 数量 | ×4.0 | llm_wiki Source Overlap (×4.0) |
| Superseded Chain | 在同一决策链上 | ×3.0 | llm_wiki Direct Link (×3.0) |
| Wikilink | 正文显式引用 | ×2.5 | llm_wiki Direct Link (×3.0)，降低权重因为引用不等于相关 |
| Source Project | 同源项目 | ×0.5 | graphify 社区检测 |
| Temporal Proximity | e^(-Δt/τ) | ×1.0 | 新增，利用时间维度 |

> **与 v1 排序的关系**：v1 使用 knowledge-classification.md 中定义的维度匹配排序（concern 重叠、artifact 特异性、language/framework 精确度等）。v2 在 v1 排序基础上叠加图谱信号（共现权重、决策链、引用），两者是递进关系而非替代关系。

**与 llm_wiki 的关键区别**：llm_wiki 的搜索算法决定最终结果，Agent 被动接收。Forge Wiki v2 的图谱扩展只用于丰富搜索结果的元数据（如附加"相关文档"列表），Agent 仍然自主决定是否查看。

## 图谱维护

### 图的构建时机

| 事件 | 触发操作 | v1 | v2 |
|------|---------|----|----|
| `forge kb push` | 更新 concern 索引，更新 `superseded_by` 链 | ✅ | ✅ + 重建受影响的局部图 |
| `forge kb lint` | 检测断裂的 `superseded_by` 链、孤立节点 | ✅ | ✅ + 完整图谱健康检查 |
| 文档内容变更 | 重新解析 `[[wikilinks]]` | — | ✅ |

### 图的存储

**v1**：无独立图存储。所有图信息分散在文档 frontmatter 中：
- `concerns` → 领域共现边（运行时计算）
- `superseded_by` → 决策链边
- `[[wikilinks]]` → 内容引用边（Obsidian 原生解析）

`forge kb concerns` 命令扫描所有文档的 `concerns` 字段，构建临时领域索引。

**v2**：可选的独立图索引文件，避免每次查询都全量扫描：

```
forge-wiki/
  wiki/
    auth-jwt-decision.md
    auth-session-decision.md
    ...
  index/
    graph.json        # 邻接表 + 边权重
    concerns.json     # concern → [doc_ids] 映射
    embeddings/       # 向量索引（可选，参考 llm_wiki LanceDB）
```

### 图的增量更新

参考 graphify 的 `--update` 模式和 llm_wiki 的 SHA256 缓存：

1. 文档内容 SHA256 哈希缓存
2. `forge kb push` 时只处理新增或修改的文档
3. 局部重建：只更新受影响节点的邻接关系
4. 全量重建：`forge kb lint --rebuild-graph`（低频，如每周）

## 可视化

### Obsidian Graph View

Forge Wiki 的文档目录本身就是 Obsidian vault。Obsidian 的 graph view 可以直接可视化：

- `[[wikilinks]]` → 文档间的连线
- 文件夹结构 → 聚类
- 搜索 → 高亮子图

无需额外实现，用户在 Obsidian 中打开 wiki 目录即可获得交互式图谱可视化。

### v2 增强可视化

如果需要展示 Obsidian 不支持的图谱特性（领域共现权重、决策链方向、时间线）：

- 生成 `graph.json`（NetworkX 格式），用 D3.js / sigma.js 渲染
- 参考 llm_wiki 的 ForceAtlas2 + sigma.js 实现
- 参考 graphify 的 `graph.html` 交互式输出

## 参考实现对比

| 维度 | llm_wiki | graphify | Forge Wiki v2（规划） |
|------|---------|----------|---------------------|
| 图类型 | Wiki 页面 wikilink 图 | 代码实体结构-语义图 | 知识文档多关系图 |
| 节点数 | 数百~数千页面 | 数百~数千实体 | 数百页面 |
| 边来源 | wikilinks + shared sources | AST + LLM 提取 | concerns + superseded_by + wikilinks + source_project + temporal |
| 社区检测 | Louvain | Leiden | Leiden（偏代码图谱场景 Louvain 更轻量） |
| 向量搜索 | LanceDB（可选） | 无 | v2 可选 |
| Agent 集成 | 系统级注入 | MCP + Skill | Agent 自主判断 |
| 核心差异 | 多信号融合，系统决策 | 图遍历，结构理解 | 多来源构图，Agent 决策 |

## 总结

Forge Wiki 的知识图谱策略是**渐进式**的：

1. **v1**：`concerns` 隐式图谱 + frontmatter 预留关系字段。零额外基础设施，搜索靠 `--concerns` 等多维度过滤 + 关键词匹配
2. **v2**：构建显式多关系图谱（5 种边类型），支持图谱扩展搜索、决策链追溯、跨域发现。新增图索引文件，增量更新

核心设计原则：**frontmatter 中每个字段既服务 v1 的直接需求，又为 v2 的图谱构建提供数据基础**。不存任何只为未来预留的字段——所有字段在 v1 都有明确的消费者。
