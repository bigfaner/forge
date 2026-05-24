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

### Config 结构

```yaml
# .forge/config.yaml
surfaces:
  frontend: web
  backend: api
  cli: cli
```

- **键**：相对于项目根目录的路径（无前导 `./`，无尾随 `/`）
- **值**：单个 surface 类型字符串
- **同路径多 surface**：拆分子路径（如 `app/pages: web` + `app/api: api`）
- **单模块项目**：`".": api`

### CLI 命令

新增独立命令 `forge surfaces`，不挂在 `forge config` 下，保持关注点分离：

```bash
# 查看所有 surface 映射
forge surfaces
# frontend=web  backend=api  cli=cli

# 按路径查询（gen-journeys 等 skill 调用）
forge surfaces <path>
# 输出: web
# 最长前缀匹配，无匹配则报错提示手动指定

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

### Non-Functional Requirements

- **检测速度**：init 中的 surface 检测应在 5 秒内完成（文件扫描 + 依赖解析）
- **向后兼容**：旧的 `interfaces` 字段应被识别并迁移，不是直接报错

### Constraints & Dependencies

- `forge init` 是 Go TUI 命令，检测逻辑必须用 Go 实现（不能依赖 LLM）
- gen-journeys 是 LLM 驱动的 skill，通过 `forge surfaces <path>` CLI 命令查询 surface
- `forge task index` 内部直接读 config struct 提取去重类型列表，不需要走 CLI
- 依赖现有的 `forge-init-config-sync` proposal（已 Approved），本 proposal 超集覆盖

## Alternatives & Industry Benchmarking

### Industry Solutions

测试框架通常根据项目类型自动推断测试策略，不需要用户手动声明接口类型。例如 Cypress 自动检测 Web 项目、Postman 自动推断 API 测试。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零成本 | 同步风险持续存在，命名 bug 不修复 | Rejected: 已有实际 bug |
| 保留两个概念，interfaces 自动检测 | 最小改动 | 改动量最小 | 概念重叠持续，仍需维护两套逻辑 | Rejected: 不够彻底 |
| surfaces 数组 | 本 proposal v1 | 简单 | 丢失路径信息，gen-journeys 仍需额外推断 | Rejected: 信息损失 |
| surfaces map + forge config 查询 | 本 proposal v2 | 统一 | 污染 config 命令，混合关注点 | Rejected: 关注点耦合 |
| **surfaces map + 独立 forge surfaces 命令** | 本 proposal v3 | 路径级精度 + 关注点分离 | 同路径多 surface 需拆子路径 | **Selected: 精确、简洁、职责清晰** |

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

### Resource & Timeline

Go 检测逻辑约 150-250 行代码（文件扫描 + 依赖解析），`forge surfaces` 命令约 50-100 行（路径匹配 + 格式化输出），加上 TUI 确认界面修改。skill 文档更新涉及 gen-journeys SKILL.md 和相关 rule 文件。工作量可控。

### Dependency Readiness

无外部依赖。`forge-init-config-sync` proposal 已 Approved 但未实现，本 proposal 可替代它。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| interfaces 和 surface 是不同层次的概念 | XY Detection | Confirmed: 它们确实作用于不同粒度（项目级 vs journey 级），但取值和用途高度重叠，统一后通过"路径映射 + CLI 查询"即可覆盖 |
| 删除 interfaces 就够了 | 5 Whys | Overturned: monorepo 存在多接口场景，需要路径级映射而非标量 |
| 简单数组就够了 | Assumption Flip | Overturned: 数组丢失路径信息，gen-journeys 仍需独立推断。map 保留检测时的路径上下文，消除推断需求 |
| surface 查询应复用 forge config | 用户反馈 | Overturned: 混合关注点增加 config 命令复杂度。独立 `forge surfaces` 命令更清晰 |
| init 中可以做完整检测 | Assumption Flip | Refined: init 只能做基于硬性规则的简单检测，LLM 驱动的精细检测仍在 gen-journeys 中 |

## Scope

### In Scope

- config schema 变更：删除 `interfaces` 和 `surface`，新增 `surfaces` map 字段（path → surface）
- 命名规范统一：`web-ui` → `web`，`mobile-ui` → `mobile`
- `forge init` Go 代码：新增 surface 自动检测逻辑（文件扫描 + 依赖解析）+ TUI 确认界面
- `forge surfaces` CLI 命令：独立命令，支持全量查看、路径查询、类型列表
- `forge task index` Go 代码：从读 `interfaces` 改为读 `surfaces`（提取去重 surface 类型列表）
- gen-journeys skill：从独立 surface 检测改为调用 `forge surfaces <path>` 查询
- gen-journeys surface rule 文件：更新引用以支持新命名

### Out of Scope

- 下游 skill 全面适配（gen-contracts, gen-test-scripts, eval-journey, eval-contract, run-tests）— 可在后续迭代中更新引用
- 旧 config 自动迁移工具 — init 重新检测会覆盖，旧字段可忽略
- 文档全面更新（ARCHITECTURE.md, conventions/ 等）— 配合下游 skill 适配一起做

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Go 检测规则无法覆盖边缘项目类型 | M | L | init 确认界面允许用户手动编辑路径和 surface |
| 同路径多 surface 拆分不直观 | M | M | 文档提供拆分示例（Next.js），init 检测覆盖常见模式 |
| gen-journeys 路径匹配与用户预期不一致 | L | M | `forge surfaces <path>` 用最长前缀匹配 + 无匹配时明确报错 |
| 下游 skill 引用旧字段名导致运行时错误 | M | H | out-of-scope skill 暂时保留对旧字段名的兼容读取 |

## Success Criteria

- [ ] `forge init` 能自动检测至少 3 种 surface 类型（web, api, cli）
- [ ] 检测结果以 path → surface map 形式在 TUI 中展示并允许用户确认/编辑
- [ ] `surfaces` 写入 `.forge/config.yaml`，格式为 map（path → surface string）
- [ ] `forge surfaces` 命令支持全量查看、路径查询（最长前缀匹配）、类型列表
- [ ] `forge task index` 从 `surfaces` 提取去重类型列表并正确生成对应测试任务
- [ ] 旧 `interfaces` 字段的命名不一致 bug 被修复
- [ ] gen-journeys 通过 `forge surfaces <path>` 查询 surface，不再独立检测

## Next Steps

- Proceed to `/write-prd` to formalize requirements
