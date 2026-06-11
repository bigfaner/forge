# Forge Wiki 知识分类体系

## 设计原则

### 第一性原理

知识检索的本质是回答两个问题：
1. **这条知识适用于哪类产出物？** — 我在做的是什么东西？
2. **这条知识解决什么层面的问题？** — 我遇到的是哪类问题？

在此基础上，对于代码级知识还需回答第三个问题：
3. **这条知识依赖什么技术栈？** — 我用的是不是同一套技术？

### 核心决策

| 决策 | 选择 | 理由 |
|------|------|------|
| 多少个维度 | 三个维度 + 一个层级声明 | 层级决定维度的作用方式，维度之间尽量正交 |
| 语言/框架是维度还是属性 | 维度（language 层级） | Go 错误处理和 Java 异常处理是硬边界，不是排序权重 |
| 跨语言知识如何处理 | langs 为空 | 空 langs = 跨语言，匹配所有搜索 |
| 非技术知识如何处理 | type 声明层级 | product/business 知识不填 langs/frameworks |
| concern 是否分技术/业务 | 统一为单一层次化列表 | 技术/业务二分法导致同一概念跨两个值域（如 auth/security），Agent 推断困难 |

## 知识层级（type）

type 声明知识的抽象层级，决定哪些维度参与匹配。

| 值 | 语义 | 举例 | langs/frameworks | artifact |
|---|------|------|-------|----------|
| `language` | 语言/框架特定的实现知识 | Go error wrapping、React hooks 规范、Spring Bean 配置 | 必填 | 必填 |
| `design` | 跨语言的架构与设计知识 | RESTful 设计原则、微服务拆分策略、Strategy 模式 | 空（跨语言） | 必填 |
| `product` | 产品决策与需求知识 | 为什么选 gRPC 而非 GraphQL、用户故事、功能取舍 | 空 | 可选 |
| `business` | 业务规则与领域知识 | 支付对账规则、数据权限模型、合规要求 | 空 | 通常不填 |

**层级与维度的关系**：

| 维度 | language | design | product | business |
|------|------|--------|---------|----------|
| artifact | 必填，参与匹配 | 必填，参与匹配 | 可选 | 通常不填 |
| concern | 统一值域，参与匹配 | 统一值域，参与匹配 | 统一值域，参与匹配 | 统一值域，参与匹配 |
| langs + frameworks | 必填，参与匹配 | 空（跨语言） | 空 | 空 |

> **关键变化**：concern 不再根据 type 使用不同值域。所有层级共用一个层次化的 concern 列表。type 只决定 langs/frameworks 是否参与匹配。

## 维度一：Artifact（软件形态）

描述知识适用的软件产出物类型。

| 值 | 语义 | 判定信号 |
|---|------|---------|
| `cli` | 命令行工具 | `cmd/` 目录、依赖 cobra/urfave/cli、有 `main` 函数接受 flags |
| `tui` | 终端用户界面 | 依赖 bubbletea/tview/rich、全屏终端渲染 |
| `web-spa` | Web 单页应用 | 依赖 React/Vue/Svelte、`src/components/` 目录 |
| `web-ssr` | Web 服务端渲染 | 依赖 Next.js/Nuxt.js、`pages/` + `getServerSideProps` |
| `api` | API 服务 | `handlers/`/`routes/` 目录、依赖 gin/express/spring-web |
| `mobile` | 移动应用 | `android/`/`ios/` 目录、依赖 React Native/Flutter |
| `library` | 库/包 | 无 `main` 函数、`export` 导出、发布到 pkg.dev/npm |
| `daemon` | 守护进程/后台服务 | systemd 配置、常驻进程、无直接用户交互 |
| `desktop` | 桌面应用 | 依赖 Electron/Tauri/Qt |
| `embedded` | 嵌入式/固件 | 跨平台编译目标、依赖 embedded-hal/Arduino |
| `batch` | 批处理/定时任务 | cron 配置、ETL 管道、一次性脚本 |
| `plugin` | 插件/扩展 | Forge skill、VSCode 扩展、Chrome 扩展 |
| `infra` | 基础设施即代码 | Terraform/Pulumi/Ansible |

**互斥性**：一个产出物在同一时刻不可能既是 `cli` 又是 `api`。但一条知识可适用于多种形态（用数组表示），一个项目也可包含多种形态（全栈 = `web-spa` + `api`）。

## 维度二：Concern（关注点）

描述知识解决什么问题。**统一为单一层次化列表**，技术关注点和业务领域关注点共享同一个值域。

### 层次化 Concern 列表

用 `/` 表示父子关系。父级是宽泛的领域，子级是具体的关注点。

```
# ─── 核心技术关注点 ───
structure               # 项目结构与模块组织
  ├── directory-layout  #   目录布局
  ├── module-design     #   模块划分
  └── dependency-mgmt   #   依赖管理

io                      # 输入输出与数据流
display                 # 界面呈现
state                   # 状态管理与数据模型
control                 # 流程控制与业务逻辑

error                   # 错误处理与异常管理

testing                 # 测试策略与方法
  ├── unit              #   单元测试
  ├── integration       #   集成测试
  └── e2e               #   端到端测试

security                # 安全
  ├── auth              #   认证授权（登录、权限、SSO、JWT/Session 选型）
  ├── encryption        #   加密（TLS、密钥管理、数据加密）
  ├── input-validation  #   输入校验（XSS、SQL 注入防护）
  └── vulnerability     #   漏洞防护（CVE、依赖审计）

perf                    # 性能与资源管理

deploy                  # 构建与部署
  ├── ci-cd             #   CI/CD 流水线
  ├── containerization  #   容器化（Docker、K8s）
  └── env-config        #   环境配置

dx                      # 开发者体验

# ─── 业务领域关注点 ───
payment                 # 支付领域
compliance              # 合规与审计
data-privacy            # 数据隐私
user-management         # 用户管理
inventory               # 库存管理
logistics               # 物流
*自定义*                # 用户自行扩展
```

### 层次匹配规则

搜索父级 concern 时，自动匹配其子级；搜索子级时，不匹配父级和兄弟：

```
搜索 concerns: [security]
  → 匹配 concerns: [security]           ✅（精确）
  → 匹配 concerns: [security/auth]      ✅（子级）
  → 匹配 concerns: [security/encryption] ✅（子级）
  → 不匹配 concerns: [error]            ❌（无关）

搜索 concerns: [security/auth]
  → 匹配 concerns: [security/auth]      ✅（精确）
  → 不匹配 concerns: [security]         ❌（父级不匹配子级搜索）
  → 不匹配 concerns: [security/encryption] ❌（兄弟不互匹配）
```

> 搜索范围越宽泛（用父级），召回越多；搜索越精确（用子级），结果越聚焦。

### 为什么不再区分技术/业务 concern

旧设计中 `auth` 是业务 concern，`security` 是技术 concern，导致：
- Agent 难以判断"添加 JWT 认证"该搜 `auth` 还是 `security`
- 同一概念被拆到两个值域，违反分面互斥性

新设计中 `security/auth` 统一了两个视角：
- `type=language` + `concerns: [security/auth]` → 认证的技术实现知识
- `type=product` + `concerns: [security/auth]` → 认证方案的产品决策

concern 描述"解决什么问题"，type 描述"知识的抽象层级"，两者正交。

### 扩展规范

业务领域关注点是开放的——用户可根据项目领域自行定义新值。

**命名规范**：
- 使用英文小写，连字符分隔：`data-privacy` 而非 `dataPrivacy`
- 使用领域通用名称，避免项目特定名称：`payment` 而非 `stripe-integration`
- 大领域用父级，子领域用 `/`：`payment/refund`、`payment/settlement`
- 参考：如果团队已有领域词汇表（如 DDD 中的 bounded context 名称），优先复用

## 维度三：Language（编程语言）

描述知识依赖的编程语言。**仅 language 层级适用**。

### 值域

`go`、`typescript`、`java`、`python`、`rust`、`csharp`、`kotlin`、`swift`、`cpp`、`ruby`、`php`、`...

> 语言值是开放的，用户可自行扩展。

### 匹配规则

精确匹配，交集非空即可：

```
文档 langs: [go]             搜索 langs: [go]           → 匹配 ✅
文档 langs: [go]             搜索 langs: [java]         → 不匹配 ❌
文档 langs: [typescript]     搜索 langs: [typescript]   → 匹配 ✅
```

## 维度四：Framework（框架）

描述知识依赖的具体框架。**仅 language 层级适用，可选**。

### 值域

`gin`、`fiber`、`cobra`、`react`、`vue`、`express`、`spring`、`fastapi`、`django`、`axum`、`...

> 框架值是开放的，用户可自行扩展。

### 匹配规则

- 文档 `frameworks` 为空 → 语言级通用知识，匹配所有框架
- 文档 `frameworks` 非空 → 框架特定知识，搜索必须匹配

```
文档 frameworks: []          搜索 frameworks: [gin]     → 匹配 ✅（Go 通用知识）
文档 frameworks: [gin]       搜索 frameworks: []        → 不匹配 ❌（Gin 特定 ≠ 无框架）
文档 frameworks: [gin]       搜索 frameworks: [gin]     → 匹配 ✅（精确匹配）
文档 frameworks: [gin]       搜索 frameworks: [fiber]   → 不匹配 ❌（框架不同）
文档 frameworks: [gin,gorm]  搜索 frameworks: [gin]     → 匹配 ✅（交集非空）
```

**设计逻辑**：
- Go 通用知识（`frameworks: []`）自动适用于 Gin、Fiber、Echo 等所有 Go 项目
- Gin 特定知识（`frameworks: [gin]`）不会污染 Fiber 或 Echo 项目
- 一个项目可用多个框架（`frameworks: [gin, gorm]`），知识只需匹配其中一个

## Frontmatter 设计

### 完整字段定义

```yaml
---
# === 必填 ===
id: go-api-error-handling          # 唯一标识
title: "Go API 错误处理规范"         # 显示名称
description: "Go 后端 API 的错误码设计、错误包装、HTTP 状态码映射"  # 搜索结果预览
type: language                          # language | design | product | business
concerns: [error, io]              # 至少一个，支持层次语法（如 security/auth）

# === 条件必填 ===
artifacts: [api]                   # language/design 必填，product 可选，business 通常不填
langs: [go]                        # 仅 language 必填，其他层级为空

# === 可选 ===
frameworks: [gin]                  # 仅 language 适用，空 = 语言级通用知识
tags: [jwt, token-refresh, exit-code]  # 自由标签，补充搜索
created: 2026-05-21                # 时间戳
updated: 2026-05-21
source_project: my-backend         # 来源项目
source_type: convention            # decision | lesson | convention | pattern | anti-pattern | reference
status: active                     # active | deprecated | superseded
superseded_by: null                # 被哪篇文档替代

# === 通用知识标记 ===
scope: universal                   # 显式声明通用知识，匹配一切搜索
---
```

### scope 规则

```
scope: universal              → 通用知识，匹配一切搜索，不受维度限制
scope: specific（默认，可省略） → 受 artifact/concern/langs/frameworks 维度过滤
```

scope 是显式声明，不是隐式推断。缺少 scope 时默认为 `specific`。

### 各类型知识文档示例

```yaml
# ─── Language 知识：Go API 错误处理 ───
id: go-api-error-handling
type: language
artifacts: [api]
concerns: [error, io]
langs: [go]
tags: [error-code, fmt.Errorf, http-status]

# ─── Language 知识：React Hooks 规范 ───
id: react-hooks-convention
type: language
artifacts: [web-spa]
concerns: [state, display]
langs: [typescript]
frameworks: [react]
tags: [useEffect, useState, rules-of-hooks]

# ─── Language 知识：Go 通用错误处理（跨框架）───
id: go-error-wrapping
type: language
artifacts: [api, cli, library]      # Go 错误包装适用于多种产出物
concerns: [error]
langs: [go]                         # 无 frameworks = 语言级通用，适用于 gin/fiber/echo
tags: [fmt.Errorf, errors.Is, errors.As]

# ─── Language 知识：Gin 中间件规范 ───
id: gin-middleware-convention
type: language
artifacts: [api]
concerns: [io, security/auth]       # 层次化 concern：认证授权
langs: [go]
frameworks: [gin]                   # Gin 特定知识
tags: [middleware, context, handler]

# ─── Language 知识：JWT 实现 ───
id: go-jwt-implementation
type: language
artifacts: [api]
concerns: [security/auth]           # 认证的技术实现
langs: [go]
frameworks: [gin]
tags: [jwt, token-refresh, rsa]

# ─── Design 知识：RESTful API 设计原则 ───
id: restful-api-design
type: design
artifacts: [api]
concerns: [io, error]
tags: [rest, http-methods, versioning]

# ─── Design 知识：Strategy 模式 ───
id: strategy-pattern
type: design
artifacts: [api, cli, library]      # 设计模式适用于多种产出物
concerns: [structure, control]

# ─── Product 知识：认证方案决策 ───
id: auth-decision-jwt-vs-session
type: product
artifacts: [api, web-spa]           # 前后端都相关
concerns: [security/auth]           # 同一个 concern，不同层级
tags: [jwt, session, oauth, decision]

# ─── Business 知识：支付对账规则 ───
id: payment-reconciliation-rules
type: business
artifacts: []                       # 不限软件形态
concerns: [payment, compliance]     # 业务领域关注点
tags: [reconciliation, settlement, audit]

# ─── 通用知识：Git 提交规范 ───
id: git-commit-convention
scope: universal
concerns: [dx]
tags: [conventional-commits, git]
```

## 匹配算法

```python
def concern_ancestor_match(doc_concerns, search_concerns):
    """层次匹配：搜索父级时匹配子级，搜索子级时不匹配父级"""
    for sc in search_concerns:
        for dc in doc_concerns:
            if dc == sc:                    # 精确匹配
                return True
            if dc.startswith(sc + '/'):     # doc 是 search 的子级
                return True
    return False

def matches(doc, search):
    # 通用知识直接匹配
    if doc.scope == 'universal':
        return True

    # type 必须匹配
    if doc.type != search.type:
        return False

    # concern 层次匹配
    if not concern_ancestor_match(doc.concerns, search.concerns):
        return False

    # artifact 匹配（空 = 不限）
    if doc.artifacts:
        if not set(search.artifacts) & set(doc.artifacts):
            return False

    # language 匹配（仅 language 类型，精确匹配）
    if doc.type == 'language' and doc.langs:
        if not set(search.langs) & set(doc.langs):
            return False

    # framework 匹配（仅 language 类型，空 = 语言级通用知识）
    if doc.type == 'language' and doc.frameworks:
        if not set(search.frameworks) & set(doc.frameworks):
            return False

    return True
```

**匹配示例**：搜索 `type=language, artifacts=[api], concerns=[security], langs=[go], frameworks=[gin]`

| 知识文档 | concerns | 结果 |
|---------|----------|------|
| Go Gin 认证中间件 | security/auth | ✅ security 匹配 security/auth（子级） |
| Go API 错误处理 | error | ❌ concern 不匹配 |
| Go Gin 错误中间件 | error, io | ❌ concern 不匹配 |
| 认证方案决策（product） | security/auth | ❌ type=product ≠ language |
| Git 提交规范 | dx | ✅ 通用知识 |

## 排序算法

```python
def concern_match_depth(doc_concern, search_concern):
    """层次匹配的深度得分：精确匹配 > 近层子级 > 远层子级"""
    if doc_concern == search_concern:
        return 3.0    # 精确匹配，最高分
    if doc_concern.startswith(search_concern + '/'):
        depth = doc_concern.count('/') - search_concern.count('/')
        return max(3.0 - depth * 0.5, 1.0)  # 每深一层减 0.5，最低 1.0
    return 0.0

def rank(doc, search):
    score = 0.0

    # 1. Concern 层次匹配深度（所有类型都适用）
    concern_score = 0.0
    for dc in doc.concerns:
        for sc in search.concerns:
            concern_score += concern_match_depth(dc, sc)
    score += concern_score

    # 2. Artifact 特异性（有 artifact 限制 = 更具体）
    if doc.artifacts:
        overlap = len(set(search.artifacts) & set(doc.artifacts))
        score += overlap * 2.0

    # 3. Language 精确度（仅 language 类型）
    if doc.type == 'language' and doc.langs:
        lang_overlap = len(set(doc.langs) & set(search.langs))
        score += lang_overlap * 3.0

    # 4. Framework 精确度（仅 language 类型，空 = 通用知识加分）
    if doc.type == 'language':
        if doc.frameworks:
            fw_overlap = len(set(doc.frameworks) & set(search.frameworks))
            score += fw_overlap * 2.5
        else:
            score += 1.0  # 语言级通用知识，额外加分（适用面广）

    # 5. 关键词匹配（query 在 title/description/tags 中的命中）
    score += keyword_match(search.query, doc) * 1.0

    # 6. 时间衰减（180 天半衰期）
    days_old = (now - doc.created).days
    score += math.exp(-days_old / 180) * 0.5

    # 7. Source type 权重
    weights = {'convention': 1.3, 'pattern': 1.2, 'decision': 1.0,
               'lesson': 1.0, 'anti-pattern': 1.1, 'reference': 0.8}
    score *= weights.get(doc.source_type, 1.0)

    # 8. Universal 知识惩罚（降低噪音，不过度占据搜索结果）
    if doc.scope == 'universal':
        score *= 0.7

    return score
```

**排序效果示例**：搜索 `type=language, artifacts=[api], concerns=[security], langs=[go], frameworks=[gin]`

| 知识文档 | concerns | 排名因素 | 排名 |
|---------|----------|---------|------|
| Go Gin 认证中间件 | security/auth | concern-精确(3)+artifact(2)+lang(3)+fw(2.5) = 10.5 | **1** |
| Go 安全通用规范 | security | concern-精确(3)+artifact(2)+lang(3)+通用(1) = 9 | **2** |
| Go 加密实现 | security/encryption | concern-子级(2.5)+artifact(2)+lang(3)+fw(2.5) = 10 | **1*** |
| Git 提交规范 | dx | concern(0)+universal 惩罚 | 低 |

> *注：加密实现虽然 concern 是子级匹配（2.5 < 3），但有 framework 精确匹配加分，可能超过通用安全规范。排序会根据具体文档的维度匹配度动态调整。

## 项目配置

### .forge/config.yaml（v3.0.0）

```yaml
# 软件形态（必填，数组形式与知识文档 frontmatter 一致）
artifacts: [api]              # cli | tui | web-spa | web-ssr | api | mobile | library | ...
                              # 全栈项目可配多个: [web-spa, api]

# 语言（自动检测或手动配置）
langs: [go]

# 框架（自动检测或手动配置，可选）
frameworks: [gin]

# 旧字段兼容：project-type 已废弃，由 artifacts 替代
# project-type: backend  →  artifacts: [api]
```

### 自动检测

| 项目文件 | 检测规则 | langs | frameworks |
|---------|---------|-------|-----------|
| `go.mod` | 存在即 Go | `go` | — |
| `go.mod` + require `gin-gonic` | Go + Gin | `go` | `gin` |
| `go.mod` + require `gofiber` | Go + Fiber | `go` | `fiber` |
| `go.mod` + require `cobra` | Go + Cobra | `go` | `cobra` |
| `package.json` + deps `react` | TS + React | `typescript` | `react` |
| `package.json` + deps `vue` | TS + Vue | `typescript` | `vue` |
| `pom.xml` | Java | `java` | — |
| `pom.xml` + dep `spring-boot` | Java + Spring | `java` | `spring` |
| `Cargo.toml` | Rust | `rust` | — |
| `pyproject.toml` + dep `fastapi` | Python + FastAPI | `python` | `fastapi` |
| 无法判断 | — | 提示用户手动配置 | — |

### Concern 推断

Concern 不在项目配置中静态声明，由 CLI 从任务关键词推断（简化模式）或由 Agent 指定（精确模式）：

| 任务描述 | 推断的 concerns | 推断的 type |
|---------|----------------|------------|
| "修复 API 端点的错误返回格式" | error, io | language |
| "添加用户认证功能" | security/auth | language |
| "重构目录结构为 Clean Architecture" | structure | design |
| "评估 REST vs gRPC 的选型" | io, control | product |
| "梳理支付对账的业务规则" | payment, compliance | business |
| "添加 E2E 测试" | testing/e2e | language |
| "优化首屏加载速度" | perf, display | language |

## 搜索命令

### 简化模式（默认）

Agent 只传关键词，CLI 自动填充维度参数：

```bash
forge kb search "token refresh"
```

CLI 内部自动完成：
1. 从 `.forge/config.yaml` 读取 `artifacts`、`langs`、`frameworks`
2. 从 query 关键词推断 `concerns`（如 "token refresh" → `security/auth`）
3. 在所有 `type` 层级中搜索，按匹配度排序

### 精确模式（可选）

需要精确控制时，手动指定维度参数：

```bash
forge kb search \
  --type language \                  # 知识层级
  --artifacts api \              # 软件形态（从项目配置）
  --concerns security/auth \    # 关注点（支持层次语法）
  --langs go \                   # 语言（从项目配置/自动检测）
  --frameworks gin \             # 框架（从项目配置/自动检测）
  "token refresh"                # 关键词
```

| 参数 | 来源 | 简化模式 | 精确模式 |
|------|------|---------|---------|
| `--type` | Agent 从任务推断 | 省略（搜索所有层级） | 不可省略 |
| `--artifacts` | 项目配置 | 自动填充 | 可省略 |
| `--concerns` | Agent/CLI 推断 | CLI 从 query 推断 | 不可省略 |
| `--langs` | 项目配置 | 自动填充 | 非 language 类型省略 |
| `--frameworks` | 项目配置 | 自动填充 | 非 language 类型省略 |

省略的参数 = 不限制该维度，扩大搜索范围。信息越全，结果越精准；信息不全，降级为更宽泛的搜索。

### 多类型搜索

简化模式下 CLI 自动搜索所有层级，无需手动多次调用。精确模式下 Agent 可以连续搜索不同层级：

```bash
# 精确模式：先查代码级知识（最精准）
forge kb search --type language --artifacts api --concerns security/auth --langs go "JWT 认证"

# 再查设计级知识（跨语言参考）
forge kb search --type design --artifacts api --concerns security/auth "认证方案设计"

# 还可以查产品决策背景
forge kb search --type product --concerns security/auth "认证方案选型"
```

### 搜索降级策略

当精确搜索结果不足时，CLI 按以下梯度自动降级：

```
精确搜索 → 去 frameworks → 去 langs → 去 artifacts → 去 type → 只留 concerns + 关键词
```

降级结果分组返回，标注降级了哪些维度：

```
## 精确匹配（language / api / go / gin / security/auth）
## go-gin-auth-middleware
- type: language | concerns: [security/auth] | langs: [go]
- Gin 中间件认证规范

## 降级搜索（已去掉 frameworks 限制）
## go-auth-common
- type: language | concerns: [security/auth] | langs: [go]
- Go 通用认证工具库使用规范
```

降级梯度原则：从最具体的维度开始去掉，保留最能区分知识相关性的维度。`concerns` 最后保留，因为它表达"解决什么问题"，是最核心的匹配信号。

降级到只剩 concerns 后，还可通过 concern 层次关系进一步扩展（搜索 `security/auth` → 扩展到 `security` → 找到通用安全知识）。

## 全景分类示例

8 种项目类型的完整分类：

| 项目 | artifact | langs | frameworks | 搜索时匹配的知识举例 |
|------|----------|-------|-----------|---------------------|
| Go CLI（forge CLI） | cli | go | cobra | CLI 参数解析、exit code、Go 错误包装 |
| Go TUI（终端仪表盘） | tui | go | bubbletea | 终端渲染、键盘事件处理、TUI 组件 |
| React Web 前端 | web-spa | typescript | react | React hooks、组件设计、状态管理 |
| Go API 后端 | api | go | gin | HTTP 路由、中间件、Go 错误码 |
| 全栈（React + Go） | web-spa + api | typescript, go | react, gin | 前端组件 + 后端 API + 跨层通信 |
| React Native 移动 | mobile | typescript | react-native | 原生桥接、设备 API、离线存储 |
| Go 通用库（日志库） | library | go | — | 接口设计、零分配、公共 API |
| Java Spring 后端 | api | java | spring | Spring Bean、JPA、REST 控制器 |

每种项目都有唯一的 artifact + langs + frameworks 组合，搜索结果天然精确。跨项目通用知识（design/product/business 层级）通过 langs 为空或 scope: universal 自动匹配所有项目。
