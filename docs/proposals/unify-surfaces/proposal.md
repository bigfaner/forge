---
created: 2026-05-24
author: "faner"
status: Draft
---

# Proposal: Unify interfaces and surface into surfaces

## Problem

Forge 测试管线中存在两个功能重叠但机制不同的概念：`interfaces`（项目级数组，控制测试任务生成）和 `surface`（journey 级标量，驱动测试策略）。两者取值相同（web-ui/api/cli/tui/mobile-ui）但设定方式不同（一个手动配置，一个自动检测），存在同步风险。

### Evidence

1. **已知 bug**：`interfaces` 为空时，`forge task index` 静默跳过所有测试任务生成（见 `docs/lessons/gotcha-test-pipeline-no-languages.md`）
2. **命名不一致 bug**：config schema 定义 `web-ui`/`mobile-ui`，但 Go 代码 `uiInterfaces` map 检查 `web`/`mobile`，导致 `T-validate-ux` 任务无法正确生成
3. **surface 字段不存在于 Go Config struct**：SKILL.md 指示 gen-journeys 将检测结果写入 `surface` 字段，但 Go 代码中 `Config` struct 无此字段
4. **forge init 不配置 interfaces**：init 流程故意排除了 interfaces，用户必须手动编辑 config 才能启用测试任务生成

### Urgency

v3.0.0 正在重新设计测试管线。两个重叠概念在 v3 发布前统一，比发布后再迁移成本低得多。如果带着这两个概念发布，每次测试管线变更都要同时维护两套逻辑。

## Proposed Solution

用一个统一的 `surfaces` map 字段（path → surface）替代 `interfaces` + `surface` 两个字段。`forge init` 自动检测每个路径的接口类型并写入 `surfaces`，提供独立 CLI 命令 `forge surfaces` 供查询。

### Innovation Highlights

**路径映射是核心洞察**：init 检测本就是在特定路径发现信号（`frontend/package.json` → web），直接记录路径比丢失路径信息更精确。gen-journeys 通过 `forge surfaces <path>` 一行命令获取 surface 类型，无需独立检测。

**与行业方案的差异化**：Cypress/Turborepo/ESLint 都只解决"路径 → 配置"或"检测 → 配置"的单向问题。Forge 的独特之处在于将**双向数据流**统一到一个 map 中：init 时从代码结构自动推导配置（代码 → config），运行时从配置查询 surface 类型（config → 测试策略）。这个双向闭环消除了大多数工具中"检测结果与配置不一致"的根本原因。

### Config 结构

```yaml
# .forge/config.yaml
surfaces:
  frontend: web
  backend: api
  cli: cli
```

**Go struct 声明**：

```go
// Config struct 变更
// 旧: Interfaces []string `yaml:"interfaces"`
// 新:
Surfaces map[string]string `yaml:"surfaces"` // 注意：不使用 omitempty

// 兼容读取过渡逻辑（v3.0.0 引入，v3.1.0/v4.0.0 移除）
// 规则 1：当 Surfaces 不为空时，直接使用 Surfaces，Interfaces 被完全忽略（无论其值）
// 规则 2：当 Surfaces 为空但旧 Interfaces 不为空时，回退读 Interfaces 并映射旧命名
//   （非 strict 模式）回退读 + 控制台 deprecation 警告
//   （strict 模式，即 --strict 标志或 FORGE_STRICT=1 环境变量）跳过回退，直接报错
```

**为什么不使用 `omitempty`**：空 `surfaces: {}` 如果被 `omitempty` 丢弃，YAML 序列化后字段消失，Go 反序列化得到 nil map。此时 `forge task index` 无法区分"未配置"和"配置为空"，复现提案声称要修复的静默跳过 bug。不使用 `omitempty` 确保空 map 序列化为 `surfaces: {}`，保留显式配置语义。

- **键**：相对于项目根目录的路径（无前导 `./`，无尾随 `/`）
- **值**：单个 surface 类型字符串
- **同路径多 surface**：拆分子路径（如 `app/pages: web` + `app/api: api`）
- **单模块项目**：`".": api`

### 路径规范化与匹配算法

**路径规范化规则**（所有输入路径在匹配前必须规范化）：

1. 去除前导 `./`（`./frontend` → `frontend`）
2. 去除尾随 `/`（`frontend/` → `frontend`）
3. 统一使用 `/` 分隔符（Windows 环境下 `\` 转为 `/`）
4. 包含 `..` 的路径视为非法输入，返回错误（不做路径解析，防止安全风险）
5. 不解析符号链接（按字面路径匹配）

**匹配算法：路径段前缀匹配**（不是字符前缀匹配）：

将查询路径和配置键都按 `/` 分割为 segment 数组，然后按 segment 逐段匹配。匹配的 segment 数最多者胜出。

```
配置: { "frontend": "web", "frontend/api": "api" }

查询 "frontend/api/routes" → segments: ["frontend","api","routes"]
  匹配 "frontend" (1 segment) 和 "frontend/api" (2 segments)
  → 最长匹配: "frontend/api" → api

查询 "frontend-new" → segments: ["frontend-new"]
  不匹配 "frontend" (segment "frontend" ≠ "frontend-new")
  → 无匹配，报错
```

**为什么用路径段而非字符前缀**：字符前缀匹配会导致 `frontend` 错误匹配 `frontend-new`，这在实际项目（重命名目录、新旧并存）中是常见场景。

### CLI 命令与退出码契约

新增独立命令 `forge surfaces`，不挂在 `forge config` 下，保持关注点分离。

**退出码定义**：

| 命令 | 成功（exit 0） | 未匹配（exit 1） | 输出流 |
|------|---------------|-----------------|--------|
| `forge surfaces` | 每行一个 `path=surface` | 不适用（空 map 也返回 0） | stdout |
| `forge surfaces <path>` | 单个 surface 类型字符串（无额外格式化） | stderr 输出错误信息（含手动配置提示） | stdout / stderr |
| `forge surfaces --types` | 空格分隔的去重类型列表 | 不适用 | stdout |

**gen-journeys skill 调用契约**：skill 通过检查退出码区分"成功"和"未找到"——退出码 0 时解析 stdout 获取 surface 类型，退出码 1 时提示用户手动配置 surfaces。

```bash
# 查看所有 surface 映射
forge surfaces
# frontend=web
# backend=api
# cli=cli

# 按路径查询（gen-journeys 等 skill 调用）
forge surfaces frontend/src
# 输出: web (exit 0)
# 或: Error: no surface found for path "unknown-dir". Run `forge init` to configure surfaces. (exit 1, stderr)

# 列出去重类型列表（调试用）
forge surfaces --types
# web api cli
```

### 统一命名规范

统一使用**无连字符**的短名称，同时修复现有的命名不一致 bug：

| 旧值 (config schema) | 旧值 (Go code) | 新统一值 |
|----------------------|----------------|---------|
| `web-ui` | `web` | `web` |
| `mobile-ui` | `mobile` | `mobile` |
| `api` | `api` | `api` |
| `cli` | `cli` | `cli` |
| `tui` | `tui` | `tui` |

## Requirements Analysis

### Key Scenarios

1. **单接口项目**：`forge init` 检测到 `".": api`，gen-journeys 所有 journey 都用 api 策略
2. **monorepo 多接口**：`forge init` 检测到 `frontend: web` + `backend: api`，gen-journeys 通过 `forge surfaces <path>` 查询
3. **检测失败**：`forge init` 无法确定 surface 类型，提示用户手动输入路径和 surface
4. **检测多选**：`forge init` 检测到多个候选路径，展示结果供用户确认或编辑
5. **用户覆盖**：检测正确但用户想添加/删除/修改某个映射，可在确认时编辑
6. **Next.js 全栈**：用户手动拆分为 `app/pages: web` + `app/api: api`
7. **信号冲突**：同一个 `package.json` 同时包含 `react`（web 信号）和 `express`（api 信号）
8. **已有项目过渡**：用户不重新运行 init，旧 `interfaces` 字段仍需被识别
9. **CI 环境（strict 模式）**：CI 中启用 `FORGE_STRICT=1` 时，若 `surfaces` 不存在或为空，即使 `interfaces` 有值也应报错退出（exit 1），而非静默回退到 `interfaces`
10. **空 surfaces**：init 后用户清空了所有 surface 映射，`forge task index` 应明确报告"无 surface 配置"而非静默跳过

### 信号冲突消歧规则

当同一个 manifest 文件（如根目录 `package.json`）同时匹配多个 surface 类型信号时，按以下优先级自动消歧：

| 优先级 | Surface 类型 | 理由 |
|--------|------------|------|
| 1（最高） | `web` | 前端应用通常内含后端 API，但用户面向的是前端 surface |
| 2 | `mobile` | 移动端应用可能共享 web 组件（react-native + react 共存），但移动端是独立部署目标 |
| 3 | `api` | 纯后端服务，无前端交互 |
| 4 | `cli` | 命令行工具 |
| 5（最低） | `tui` | 终端 UI 应用 |

**典型冲突场景**：
- `react` + `express` → web 优先（前端 fullstack 框架，API 是内部实现）
- `react-native` + `react` → mobile 优先（移动端独立部署，web 共享组件库）
- `react-native` + `express` → mobile 优先（移动端应用包含 API 调用，但用户面向移动端）

**冲突处理流程**：
1. init 检测到同一 manifest 存在信号冲突
2. 按优先级表自动选择最高优先级的 surface 类型作为默认值
3. TUI 确认界面中**标注冲突信号**（如"检测到 web + api 信号，已按优先级选择 web"），用户可手动覆盖

### 未知类型处理策略

`surfaces` map 中的值必须是已知的 surface 类型（`web`/`api`/`cli`/`tui`/`mobile`）。当 `forge task index` 从 surfaces 提取去重类型列表时遇到未知类型：
- **不报错，不中断**：未知类型被忽略，不影响已知类型的正常处理
- **输出 deprecation 级别日志**：`log.Warn("unknown surface type ignored", "type", "unknown-type", "path", "frontend")`
- **不传透给下游**：未知类型不出现在 `forge surfaces --types` 的输出中，不生成对应的测试任务

这保持了与当前 `hasUIInterface` 函数对未知类型返回 false 的一致行为。

### Non-Functional Requirements

- **检测速度**：init 中的 surface 检测应在 5 秒内完成（文件扫描 + 依赖解析，深度限制 1-10 层保证不会遍历过深）
- **向后兼容（兼容读取过渡期）**：v3.0.0 引入 `surfaces` 字段的同时，Go 代码增加兼容读取过渡逻辑，遵循以下规则：
  - **规则 1（surfaces 非空优先）**：当 `surfaces` map 不为空时，直接使用 `surfaces`，`interfaces` 字段被**完全忽略**（无论其值是什么）。即使 `interfaces` 包含 `surfaces` 中不存在的类型，也不合并。
  - **规则 2（空 surfaces 回退）**：当 `surfaces` 为空但 `interfaces` 不为空时，行为取决于运行模式：
    - **普通模式**（默认）：回退读取 `interfaces` 字段（同时将旧值映射到新命名：`web-ui` → `web`，`mobile-ui` → `mobile`），并在控制台输出 deprecation 警告（`"interfaces field is deprecated, run forge init to migrate to surfaces"`）。
    - **strict 模式**（通过 `--strict` 命令行标志或 `FORGE_STRICT=1` 环境变量启用）：跳过回退，直接报错退出（`"surfaces is empty; in strict mode, interfaces fallback is disabled. Run forge init to configure surfaces"`，exit 1）。
  - strict 模式适用于 CI 环境，确保配置问题被显式暴露而非静默兼容。此过渡逻辑计划在 v3.1.0 或 v4.0.0 中移除。
- **路径规范化性能**：路径规范化（去除前导/尾随字符、分隔符转换）不应成为性能瓶颈，实现应避免不必要的字符串分配
- **Windows 兼容性**：所有路径处理必须正确处理 `\` 分隔符，统一转为 `/` 后存储和匹配
- **YAML 序列化一致性**：`surfaces` 字段的 YAML tag **不使用** `omitempty`（详见 Config struct 声明），避免空 map 被丢弃后复现静默跳过 bug

### Constraints & Dependencies

- `forge init` 是 Go TUI 命令，检测逻辑必须用 Go 实现（不能依赖 LLM）
- gen-journeys 是 LLM 驱动的 skill，通过 `forge surfaces <path>` CLI 命令查询 surface
- `forge task index` 内部直接读 config struct 提取去重类型列表，不需要走 CLI
- 路径分隔符跨平台约束：所有路径在存储和匹配时统一使用 `/`，Windows 环境下的 `\` 在输入时转为 `/`
- `forge-init-config-sync` proposal（已 Approved 但未实现）被本 proposal 超集覆盖，需标记为 Superseded

## Alternatives & Industry Benchmarking

### Industry Solutions

测试框架和构建工具通常根据项目结构自动推断测试/构建策略，无需用户手动声明接口类型。以下是具体工具的检测机制分析：

**1. Cypress（v13+）Component Testing**
- **检测机制**：Cypress 读取项目根目录的 `package.json`，检查 `devDependencies` 中是否包含 `react`/`vue`/`svelte` 等前端框架依赖，自动决定使用哪个组件测试适配器（cypress/react、cypress/vue 等）。
- **配置方式**：`cypress.config.ts` 中的 `component` 字段。如果检测到前端框架但用户未配置，Cypress 会在首次启动时引导用户完成配置（interactive setup）。
- **与 Forge 的差异**：Cypress 只检测单个项目根目录，不支持 monorepo 子目录级别的自动检测。Forge 的 `surfaces` map 需要处理 monorepo 多目录场景，这是 Cypress 不涉及的。

**2. Turborepo Pipeline Configuration**
- **检测机制**：Turborepo 不做自动检测，而是通过 `turbo.json` 的 `pipeline` 配置声明每个 workspace package 的构建和测试任务。依赖关系通过 workspace 内部的 `package.json` 依赖自动推断执行拓扑。
- **配置方式**：`turbo.json` 的 pipeline 定义任务间依赖（`dependsOn`），workspace package 间通过 `internalDependencies` 自动追踪。
- **与 Forge 的差异**：Turborepo 将"做什么"留给用户声明，"执行顺序"自动化。Forge 的方向相反——"做什么"（surface 类型）通过检测自动推断，用户只需确认。Turborepo 的声明式模式更确定但配置成本高，Forge 的检测模式更智能但需要处理消歧。我们选择检测模式是因为 Forge 的目标用户是开发者个人，减少手动配置是核心价值。

**3. ESLint Override by Path（flat config）**
- **检测机制**：ESLint flat config（`eslint.config.js`）通过 `files` glob 模式匹配文件路径，为不同路径应用不同的 lint 规则集。这是"路径 → 配置"映射的行业验证模式。
- **配置方式**：`files: ["frontend/**/*.{js,jsx}"]` + 对应规则数组。支持 glob 而非精确路径匹配。
- **与 Forge 的参考价值**：ESLint 证明了"按路径区分配置策略"是行业成熟模式。Forge 的 `surfaces` map 借鉴了这一思路，但使用精确路径而非 glob（因为 surface 类型在目录级别确定，不需要文件级粒度）。ESLint 的 glob 匹配比 Forge 的路径段前缀匹配更复杂，我们选择更简单的方案以降低实现和调试成本。

**4. Jest Projects Configuration**
- **检测机制**：Jest 通过 `projects` 字段（数组）为 monorepo 中每个子项目定义独立的测试配置。不自动检测——用户必须手动声明项目列表和配置。
- **配置方式**：`projects: ["apps/*"]` + 每个项目的 `jest.config.js`。支持 glob 匹配目录。
- **与 Forge 的差异**：Jest 的 `projects` 是声明式数组，不含类型信息（不区分 web/api）。Forge 的 `surfaces` map 同时编码了路径和类型信息，这是 Forge 特有的需求（因为测试策略依赖 surface 类型）。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零成本 | 同步风险持续存在，命名 bug 不修复 | Rejected: 已有实际 bug |
| 保留两个概念，interfaces 自动检测（类似 Jest projects 声明式模式） | 最小改动 | 改动量最小 | 概念重叠持续，仍需维护两套逻辑 | Rejected: 不够彻底 |
| **纯声明式模式**（用户手动配置 surfaces，不做自动检测） | Turborepo `turbo.json` pipeline 模式 | 确定性高、可调试性好、CI 友好（无检测不确定性） | 配置成本高，违背 Forge"减少手动配置"的核心价值；Turborepo 的声明式模式适合团队协作但 Forge 面向个人开发者 | Rejected: 配置成本过高，不适合 Forge 目标用户 |
| **交互式引导模式**（检测 + 首次运行时引导用户完成配置，类似 Cypress 首次启动引导） | Cypress v13+ 首次启动 interactive setup | 用户教育成本低（边用边学），配置结果确定性高（用户最终确认），已有行业验证 | 需要交互式环境（CI 不友好），引导流程设计成本高，Forge 的 init 已经承担了这个角色 | Rejected: Forge init 已涵盖此模式的"检测+确认"流程，无需额外的运行时引导 |
| **Glob 模式匹配**（类似 ESLint flat config 的 `files` glob） | ESLint `eslint.config.js` | 灵活性极高，支持文件级粒度 | 实现复杂度高（glob 引擎），调试困难；surface 类型在目录级确定，不需要文件级粒度，glob 能力过剩 | Rejected: 复杂度过高，收益不匹配 |
| **项目数组 + 独立配置**（类似 Jest projects，每个子项目有独立 surface 配置） | Jest `projects` field | 成熟模式，生态验证充分 | 不含类型信息（Jest 不区分 web/api），需额外抽象层；数组无路径映射能力，gen-journeys 仍需推断 | Rejected: 信息密度不足 |
| surfaces 数组 | 本 proposal v1 | 简单 | 丢失路径信息，gen-journeys 仍需额外推断 | Rejected: 信息损失 |
| surfaces map + forge config 查询 | 本 proposal v2 | 统一 | 污染 config 命令，混合关注点 | Rejected: 关注点耦合 |
| **surfaces map + 独立 forge surfaces 命令** | 本 proposal v3（综合 ESLint 路径配置 + Turborepo workspace 隔离 + Cypress 检测引导思路） | 路径级精度 + 自动检测减少配置 + 关注点分离 | 同路径多 surface 需拆子路径；strict 模式增加 CLI 复杂度 | **Selected: 精确、简洁、职责清晰** |

### Trade-off 深度分析

**"同路径多 surface 需拆子路径"的实际影响**：
此约束仅在 fullstack 框架（Next.js、Nuxt）中出现。全栈框架在 Forge 的目标用户群（中小型项目）中占比约 15-20%。拆分子路径（`app/pages: web` + `app/api: api`）的用户体验可接受——init 检测可自动建议拆分方案，用户只需确认。对于纯 API 或纯前端项目（占比 80%+），不涉及此约束。

## Feasibility Assessment

### Technical Feasibility

Go 代码实现 surface 检测是完全可行的。检测逻辑基于文件模式匹配：

| 路径信号 | 依赖/内容检测 | Surface |
|---------|-------------|---------|
| `package.json` | react/vue/svelte + DOM entry | `web` |
| `package.json` | express/fastify/koa（无前端框架） | `api` |
| `package.json` | commander/yargs/oclif（无前端框架） | `cli` |
| `package.json` | blessed/ink/neo-blessed | `tui` |
| `go.mod` | gin/echo/chi/net-http | `api` |
| `go.mod` | cobra/urfave | `cli` |
| `go.mod` | bubbletea/tview | `tui` |
| `Cargo.toml` | actix/axum/rocket | `api` |
| `Cargo.toml` | clap/structopt | `cli` |
| `Cargo.toml` | ratatui/cursive | `tui` |
| `AndroidManifest.xml` | — | `mobile` |
| `*.xcodeproj` / `pubspec.yaml`+flutter | — | `mobile` |
| `pyproject.toml`/`setup.py` | flask/fastapi/django | `api` |
| `pyproject.toml`/`setup.py` | click/typer/argparse | `cli` |

检测输出直接就是 path → surface 的 map，无需额外转换。

### 检测遍历策略（Detection Traversal Strategy）

在 monorepo（pnpm workspace / yarn workspaces / npm workspaces）中，根目录的 `package.json` 通常只声明 workspace 配置，实际依赖在子目录。检测算法需遵循以下策略：

**遍历规则**：
1. **深度限制**：默认最多遍历 3 层子目录。选择依据：npm/yarn/pnpm workspaces 的典型结构为 `packages/<scope>/<pkg>` 或 `apps/<pkg>`，2-3 层覆盖绝大多数场景。3 层覆盖 `apps/web/client`、`packages/team-a/frontend` 等常见结构，4+ 层在实际项目中极罕见（<1%，且通常可简化目录结构）。深度限制可通过 `FORGE_DETECT_DEPTH` 环境变量覆盖（有效值范围 1-10；设为 0 不合法，init 报错提示有效范围）。不提供"无限制"选项是因为超大项目中无限制遍历可能导致挂起或 OOM，与 5 秒检测速度要求冲突。
2. **排除目录**：跳过 `node_modules`、`.git`、`vendor`、`dist`、`build`、`__pycache__`、`.next`、`target`。排除目录列表不提供配置能力（这些目录不可能包含有效的 surface 信号）。
3. **Workspace manifest 处理**：当检测到 `pnpm-workspace.yaml` 或 `package.json` 中包含 `workspaces` 字段时，跳过根目录的依赖检测（根 `package.json` 的 dependencies 不作为 surface 信号），只检测子目录的 manifest 文件
4. **非 workspace 项目**：根目录的 `package.json` / `go.mod` / `Cargo.toml` 正常检测，路径记录为 `"."`

**遍历示例**：
```
monorepo/
├── package.json (workspaces: ["apps/*", "packages/*"])  → 跳过根检测
├── pnpm-workspace.yaml                                    → 跳过根检测
├── apps/
│   ├── web/package.json (react)                           → apps/web: web
│   └── api/package.json (express)                         → apps/api: api
└── packages/
    └── cli/package.json (commander)                       → packages/cli: cli
```

### Resource & Timeline

Go 检测逻辑约 150-250 行代码（文件扫描 + 依赖解析），`forge surfaces` 命令约 50-100 行（路径段前缀匹配 + 格式化输出 + 退出码处理），兼容读取过渡逻辑约 30-50 行，加上 TUI 确认界面修改。gen-journeys skill 适配工作量：SKILL.md 指令更新 + rule 文件重命名和内容更新（约 5-8 个文件）。总计可控。

### Dependency Readiness

无外部依赖。`forge-init-config-sync` proposal 已 Approved 但未实现。本 proposal 超集覆盖其功能（config schema 变更 + init 集成），实施时需将 `forge-init-config-sync` 标记为 Superseded-by 本 proposal，避免两个 proposal 同时实施时产生冲突的代码变更。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| interfaces 和 surface 是不同层次的概念 | XY Detection | Confirmed: 它们确实作用于不同粒度（项目级 vs journey 级），但取值和用途高度重叠，统一后通过"路径映射 + CLI 查询"即可覆盖 |
| 删除 interfaces 就够了 | 5 Whys | Overturned: monorepo 存在多接口场景，需要路径级映射而非标量 |
| 简单数组就够了 | Assumption Flip | Overturned: 数组丢失路径信息，gen-journeys 仍需独立推断。map 保留检测时的路径上下文，消除推断需求 |
| surface 查询应复用 forge config | 用户反馈 | Overturned: 混合关注点增加 config 命令复杂度。独立 `forge surfaces` 命令更清晰 |
| init 中可以做完整检测 | Assumption Flip | Refined: init 只能做基于硬性规则的简单检测，LLM 驱动的精细检测仍在 gen-journeys 中 |
| 新字段用 omitempty 是安全的 | Bug Class Analysis | Overturned: `interfaces` 的静默跳过 bug 根因就是 omitempty 导致空值被丢弃。新 `surfaces` 字段必须不使用 omitempty，否则复现同类 bug |
| 不迁移旧字段是安全的 | Edge Case Analysis | Overturned: 已有项目不重新 init 时旧字段被完全忽略，forge task index 静默停止生成任务。需要兼容读取过渡期 |
| 字符前缀匹配足够精确 | Boundary Test | Overturned: `frontend` 会错误匹配 `frontend-new`。必须使用路径段前缀匹配 |
| 检测深度无限制是合理的 | Performance Analysis | Overturned: 无限制遍历在大型 monorepo 中可能挂起或 OOM，与 5 秒检测速度要求矛盾。限制为 1-10 并拒绝 0 值 |

## Scope

### In Scope

- config schema 变更：删除 `interfaces` 和 `surface`，新增 `surfaces` map 字段（path → surface）
- Go Config struct：`Surfaces map[string]string \`yaml:"surfaces"\``（不使用 omitempty），兼容读取过渡逻辑
- 命名规范统一：`web-ui` → `web`，`mobile-ui` → `mobile`
- `forge init` Go 代码：新增 surface 自动检测逻辑（文件扫描 + 依赖解析）+ 信号冲突消歧 + TUI 确认界面
- `forge surfaces` CLI 命令：独立命令，支持全量查看、路径查询（路径段前缀匹配）、类型列表，退出码契约
- `forge task index` Go 代码：从读 `interfaces` 改为读 `surfaces`（提取去重 surface 类型列表，含兼容读取过渡 + strict 模式支持）
- strict 模式：`--strict` 命令行标志或 `FORGE_STRICT=1` 环境变量，跳过 `interfaces` 兼容回退，空 surfaces 直接报错（适用于 CI 环境）
- gen-journeys skill 适配（核心消费者，必须同步更新）：
  - 从独立 surface 检测改为调用 `forge surfaces <path>` 查询
  - SKILL.md 中 `surface` 字段引用更新为 `surfaces` map
  - surface rule 文件重命名：`surface-webui.md` → `surface-web.md`，`surface-mobileui.md` → `surface-mobile.md`
  - rule 文件中旧命名值引用更新（`webui` → `web`，`mobileui` → `mobile`）
- 路径规范化与匹配：路径段前缀匹配实现 + 路径规范化规则

### Out of Scope

- 下游 skill 全面适配（gen-contracts, gen-test-scripts, eval-journey, eval-contract, run-tests）— gen-journeys 以外的 skill 可在后续迭代中更新引用，因它们不直接读写 `surfaces` 字段
- 旧 config 自动迁移工具（`interfaces` → `surfaces` 的一键转换）— 通过兼容读取过渡期解决，无需独立迁移工具
- 文档全面更新（ARCHITECTURE.md, conventions/ 等）— 配合下游 skill 适配一起做

### 版本计划

- **v3.0.0**：`surfaces` map + 检测 + CLI 命令 + gen-journeys 适配 + 兼容读取过渡逻辑（含 `interfaces` 回退 + strict 模式）
- **v3.1.0 / v4.0.0**：移除 `interfaces` 兼容读取过渡逻辑，全面切换到 `surfaces`

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Go 检测规则无法覆盖边缘项目类型 | M | L | init 确认界面允许用户手动编辑路径和 surface |
| 同路径多 surface 拆分不直观 | M | H | 文档提供拆分示例（Next.js），init 检测覆盖常见模式。Next.js 等 fullstack 框架常见，impact 调高为 H |
| gen-journeys 路径匹配与用户预期不一致 | L | M | `forge surfaces <path>` 用路径段前缀匹配（非字符前缀）+ 无匹配时明确报错 |
| 下游 skill 引用旧字段名导致运行时错误 | M | H | out-of-scope skill 暂时保留对旧字段名的兼容读取，gen-journeys 已纳入 In Scope 同步更新 |
| **旧 `interfaces` 字段被忽略导致静默信息丢失** | M | H | v3.0.0 引入兼容读取过渡逻辑：surfaces 为空时回退读 interfaces + deprecation 警告；CI 环境通过 `FORGE_STRICT=1` 跳过回退直接报错。确保未重新 init 的项目在普通模式不丢失配置，CI 环境显式暴露配置问题 |
| **surfaces 空 map 被 omitempty 丢弃复现静默跳过 bug** | H | H | Config struct 的 YAML tag 显式不使用 `omitempty`：`yaml:"surfaces"`。空 map 序列化为 `surfaces: {}` 保留显式语义 |
| **gen-journeys 新旧字段并存导致同步问题** | M | H | gen-journeys 的 surface 字段引用和 rule 文件重命名已纳入 In Scope，确保 v3.0.0 同步发布。兼容期内新旧字段不冲突（Go 读 surfaces，gen-journeys 通过 CLI 查询 surfaces） |

## Success Criteria

- [ ] `forge init` 能自动检测至少 3 种 surface 类型（web, api, cli）
- [ ] 检测结果以 path → surface map 形式在 TUI 中展示，TUI 验收标准如下：
  - **确认按钮**：存在明确的确认操作（如 "Confirm" 按钮或 Enter 快捷键），确认后将 surfaces 写入 config
  - **编辑入口**：每个 path=surface 行可通过特定操作（如按 `e` 键或点击行）进入编辑模式，修改路径或 surface 类型
  - **冲突信号标注**：当检测到信号冲突时，该行显示格式为 `path: surface (冲突信号: web + api，已按优先级选择 web)`，冲突类型以高亮或颜色区分
  - **添加/删除**：提供添加新映射（空白行输入）和删除已有映射（选中行后按 `d` 键）的操作入口
- [ ] `surfaces` 写入 `.forge/config.yaml`，格式为 map（path → surface string），空 map 写入 `surfaces: {}` 而非字段缺失
- [ ] `forge surfaces` 命令支持全量查看（每行 `path=surface`）、路径查询（路径段前缀匹配，未匹配时 exit 1 + stderr 错误信息）、类型列表
- [ ] `forge task index` 从 `surfaces` 提取去重类型列表并正确生成对应测试任务
- [ ] 兼容读取过渡验证：
  - **规则 1 验证**：当 `surfaces` 非空且 `interfaces` 也非空时，`forge task index` 只使用 `surfaces`，`interfaces` 中的值被完全忽略（验证方法：config 中 `surfaces: {frontend: web}` + `interfaces: [api, cli]`，`forge surfaces --types` 输出仅包含 `web`，不包含 `api` 或 `cli`）
  - **规则 2 普通模式验证**：当 `surfaces` 为空但 `interfaces` 不为空时，`forge task index` 回退读 `interfaces`（含旧命名映射 `web-ui` → `web`）并输出 deprecation 警告
  - **规则 2 strict 模式验证**：设置 `FORGE_STRICT=1` 后，当 `surfaces` 为空时，即使 `interfaces` 非空也报错退出（exit 1 + stderr 错误信息）
- [ ] 命名规范统一验证：config 中 `web-ui` / `mobile-ui` 值不再出现，Go 代码和 config 使用一致的 `web` / `mobile` 短名称
- [ ] gen-journeys 通过 `forge surfaces <path>` 查询 surface（检查退出码 0），不再独立检测
- [ ] gen-journeys rule 文件已重命名（`surface-webui.md` → `surface-web.md`），文件内引用的命名值已更新
- [ ] 路径规范化边界验证：
  - `..` 路径报错：`forge surfaces ../etc` 返回 exit 1 + stderr 错误信息 "path contains '..'"
  - Windows `\` 转换：`forge surfaces frontend\src` 在 Windows 上等同于 `forge surfaces frontend/src`，返回相同结果
  - 符号链接不解析：如果 `frontend` 是符号链接指向 `packages/web`，`forge surfaces frontend` 按字面路径 `frontend` 匹配，不尝试解析为 `packages/web`
- [ ] 路径匹配边界验证：`forge surfaces frontend-new` 不匹配配置键 `frontend`（路径段匹配 vs 字符前缀匹配）
- [ ] `FORGE_DETECT_DEPTH` 无效值拒绝：设为 0 或负数时 init 报错提示有效范围（1-10），不静默忽略

## Next Steps

- Proceed to `/write-prd` to formalize requirements
