# CLI 作为 Agent 与知识库的 Transport Layer

> 状态：Proposal
> 日期：2026-05-11

## 核心主张

CLI 不是知识库，不是搜索引擎，不是 RAG 引擎。**CLI 是 Agent 与知识库后端之间的 transport layer**——一根管子，负责认证、路由、格式转换，不做智能决策。

```
Agent ──→ CLI (transport) ──→ 知识库后端 (intelligence)
               │                    │
          只负责：               可以是：
          - 认证                  - RAG
          - 路由                  - 文档推荐服务
          - 格式转换              - 传统搜索 (ES/Meilisearch)
          - 透传                  - LLM 总结服务
```

## 为什么是 CLI 而不是 MCP

| 维度 | MCP Server | CLI Transport |
|------|-----------|---------------|
| Token 开销 | tool description 固定占用 context | `-h` 一次性学习，按需调用 |
| Agent 自主性 | 受限于注册的 tool schema | 自然语言 `ask`，自由度更高 |
| 后端灵活性 | 每个 MCP server 绑定一种后端 | CLI 路由到任意后端 |
| 部署复杂度 | 需要 MCP server 进程常驻 | CLI 按需执行，无守护进程 |
| 生态兼容 | 仅 Claude Code | 任何支持 CLI 的 Agent |
| 推荐能力 | 需要额外 tool 定义 | `suggested_queries` 原生内置 |

MCP 的 tool description 本身就占 context window。注册 10 个 MCP tool 的 description 可能比一次 `kb -h` 还贵。

CLI 的核心优势：**轻、通、透**。不试图做智能，只做可靠的管道。

## 架构

```
┌─────────────────────────────────────────────────────────┐
│  Agent (Claude Code / Cursor / Copilot)                  │
│  ┌─────────────────────────────────────────────────────┐ │
│  │  意图生成：我需要什么 → kb ask "..."                  │ │
│  │  结果消费：拿到回答 → 纳入编码上下文                   │ │
│  └──────────────────────┬──────────────────────────────┘ │
└─────────────────────────┼───────────────────────────────┘
                          │ stdio / exit code
┌─────────────────────────▼───────────────────────────────┐
│  CLI Transport Layer                                     │
│  ┌──────────┬──────────┬──────────┬──────────────────┐  │
│  │ Auth     │ Router   │ Format   │ State            │  │
│  │ Token    │ 后端选择  │ 输入输出  │ Session/Cache    │  │
│  │ Refresh  │ 负载均衡  │ 编解码   │ 去重/合并        │  │
│  └──────────┴──────────┴──────────┴──────────────────┘  │
└──────────┬──────────────┬──────────────┬────────────────┘
           │              │              │
     ┌─────▼─────┐  ┌────▼─────┐  ┌────▼──────┐
     │   RAG     │  │  Rec-Sys │  │ LLM-Sum   │
     │ 向量检索   │  │ 文档推荐  │  │ 知识总结   │
     └───────────┘  └──────────┘  └───────────┘
```

## Agent-CLI 协议

### 命令设计

CLI 暴露给 Agent 的接口极简，复杂度全部下沉到后端：

```bash
kb -h                          # 自学习：我有什么能力
kb ask <query>                 # 核心：自然语言提问
kb ask <query> --context <...> # 带编码上下文的提问
kb related <topic>             # 发现：围绕某个主题还有哪些知识
kb summary <doc-id>            # 深挖：对某个文档要摘要或全文
kb context --auto              # 自动从 cwd/git 推断 context
```

### Agent 自学习闭环

CLI 的 `-h` 机制让 Agent 形成自发的知识探索链：

```
kb -h → 学习能力
  ↓
kb ask "X" → 获取知识 + suggested_queries
  ↓
kb ask "Y" → 发现盲区并补全
  ↓
知识足够，开始编码
```

关键在于 `suggested_queries`——后端不只回答当前问题，还建议下一步该问什么。Agent 自主决定是否继续追问。

### 意图表达的三个层次

```
Session 开始 ──→ 编码中 ──→ 遇到边界情况 ──→ Session 结束
     │               │              │               │
  宽泛探索        精准提问        关联发现        知识沉淀
```

```bash
# 1. 宽泛探索 — Session 启动时
kb ask "当前项目相关的所有技术规范索引"

# 2. 精准提问 — 编码中
kb ask "错误码设计规范" --context "Go microservice, gRPC transport"

# 3. 关联发现 — 遇到边界
kb related "rate-limiting"

# 4. 深挖细节 — 需要完整内容
kb summary "auth-middleware-v2" --section "token-refresh"
```

### Context 自动推断

`--context` 让后端精准匹配。CLI 可从环境自动提取：

```bash
kb context --auto
# Detected: Go, gRPC, microservice, project=payment-service
```

## CLI-Backend 协议

### 后端接口抽象

```typescript
interface KBBackend {
  ask(query: string, context?: QueryContext): Promise<KBResponse>;
  related(topic: string): Promise<RelatedItem[]>;
  get(docId: string, options?: GetOptions): Promise<KBResponse>;
  capabilities(): BackendCapabilities;
}
```

### 适配器模式

CLI 不绑定任何特定后端，通过适配器支持多种后端：

| 适配器 | ask() 转换 | 特点 |
|--------|-----------|------|
| RAGAdapter | embedding(query) → topK | 高精度，需要向量库 |
| RecAdapter | 用户画像 + 内容匹配 | 个性化推荐 |
| LLMAdapter | 检索 + LLM 合成回答 | 最省 token，可能丢细节 |
| SearchAdapter | ES/Meilisearch 关键词匹配 | 简单可靠 |
| HybridAdapter | 编排多个后端 | RAG 粗筛 + LLM 精排 + 推荐补充 |

配置化切换：

```yaml
# kb.config.yaml
backend:
  primary: rag
  fallback: search
  rag:
    endpoint: http://rag.internal:8080
    model: text-embedding-3-small
    top_k: 5
  search:
    endpoint: http://es.internal:9200
    index: tech-specs
```

## Session-Level Cache

CLI 层唯一值得有的智能——**去重和增量合并**。

### 问题

Agent 在同一个 session 中可能多次查询相关知识：

```
T1: kb ask "错误码规范"        → 返回 200 tokens
T3: kb ask "错误码HTTP映射"    → 和 T1 有 60% 内容重叠
T7: kb ask "错误码gRPC映射"    → 和 T1 有 40% 内容重叠
```

三次查询浪费了重复 token。

### 解法

```
~/.kb/session/
  query_cache:     "错误码规范" → {response, topics, ts}
  fetched_docs:    "error-code-v2" → sections_fetched
  dedup_tracker:   content_hash → first_seen_query
```

重复查询时只返回增量部分，标注 `[cached: 错误码规范 §2.1]`。

## 输出格式

### 默认：精简 Markdown

```bash
$ kb ask "错误码规范"
## 错误码设计规范 v2.1
- 格式：`{service}-{module}-{code}` (例：`auth-token-001`)
- HTTP映射：4xx 业务错误，5xx 系统错误
- gRPC映射：用 status.Code + 自定义 metadata
[来源: docs/error-code-v2.md §3.1-3.3] [置信度: 0.92]
```

### 结构化：JSON

```bash
$ kb ask "错误码规范" --format json
{
  "answer": "...",
  "sources": [{"doc": "error-code-v2.md", "sections": ["3.1","3.2","3.3"]}],
  "confidence": 0.92,
  "related": ["rate-limiting", "error-middleware"],
  "suggested_queries": ["gRPC status code映射", "错误码注册流程"]
}
```

每个响应都附带：
- **来源**：可追溯，Agent 可判断可信度
- **置信度**：后端对回答质量的自我评估
- **suggested_queries**：推荐下一步探索方向

## 安全模型

```
Agent                    CLI                      Backend
  │                       │                          │
  │  kb ask "..."         │                          │
  │──────────────────────►│  1. Token 注入           │
  │                       │     (env / keychain)     │
  │                       │─────────────────────────►│
  │                       │                          │
  │                       │  2. 权限检查              │
  │                       │◄─────────────────────────│
  │                       │     (用户能访问哪些文档)    │
  │                       │                          │
  │                       │  3. 审计日志              │
  │                       │──────────────────────────→ audit-log
```

CLI 层安全职责：

- **认证**：Token 管理、刷新、安全存储
- **授权**：透传用户身份，后端做权限判断
- **脱敏**：输出过滤（检测到密钥/证书内容自动截断）
- **审计**：记录每次查询，可追溯

## 可靠性：Fail-Safe 设计

CLI 是 Agent 和知识库之间的唯一通道，必须 fail-safe：

```bash
# 后端不可用 — 降级到本地缓存
$ kb ask "错误码规范"
# [WARN] RAG backend unreachable, falling back to local cache (updated 2h ago)
# (cached response...)

# 完全不可用 — 明确告知 Agent
$ kb ask "错误码规范"
# [ERROR] Knowledge base unavailable. Response confidence: NONE.
# Do NOT use any speculative knowledge. Ask the human for guidance.
```

**核心原则：宁可让 Agent 知道"我不知道"，也不能让它拿到错误答案后自信地写出违规代码。**

## 完整交互示例

Agent 编码 session 中与知识库的完整交互流程：

```bash
# Session 启动，Agent 自学习
$ kb -h
# Commands: ask, related, summary, context
# Backends: rag (primary), search (fallback)

# Agent 自动推断上下文
$ kb context --auto
# Detected: Go, gRPC, microservice, project=payment-service

# 任务：写支付回调接口
$ kb ask "支付回调接口设计规范"
# ## 支付回调 v3
# - 必须验签（RSA-SHA256）
# - 幂等设计：用 callback_id 去重
# - 超时重试：指数退避，最多5次
# [来源: payment/callback-spec-v3.md] [confidence: 0.95]
# suggested: ["回调签名验证实现", "幂等设计方案", "支付状态机"]

$ kb ask "回调签名验证实现" --context "Go, crypto/rsa"
# ## RSA-SHA256 验签 (Go)
# - 公钥从配置中心获取，不硬编码
# - signature 放在 HTTP header `X-Signature`
# [来源: payment/callback-spec-v3.md §4.2] [confidence: 0.94]
# suggested: ["配置中心公钥管理", "签名错误处理"]

$ kb related "幂等设计"
# Related topics:
# 1. 分布式锁方案 (relevance: 0.88)
# 2. 数据库唯一约束 (relevance: 0.82)
# 3. 消息队列幂等 (relevance: 0.79)

# Agent 判断已有足够知识，开始编码
```

## 设计原则总结

1. **CLI 是管子，不是大脑**——智能留给后端，决策留给 Agent
2. **`-h` 驱动自学习**——Agent 一次学习能力，按需调用
3. **`suggested_queries` 驱动探索**——后端建议，Agent 自主追问
4. **输出可追溯**——每个回答附带来源和置信度
5. **Fail-safe 而非 fail-open**——不可用时明确告知，不返回猜测
6. **后端可替换**——适配器模式，CLI 层零改动切换后端
7. **Session cache 去重**——唯一值得在 CLI 层做的智能
