---
created: "2026-05-30"
author: "faner"
status: Approved
intent: "cleanup"
---

# Proposal: 用户手册与架构说明文档体系

## Problem

Forge 缺少面向终端用户的完整文档体系。当前文档分散在多个文件中，信息存在重复和缺口：

- **环境配置**：README.md 仅 3 行前置要求 + 6 行安装命令，缺少 OS 要求、Claude Code 版本兼容性、Go 环境验证、安装后验证步骤
- **初始化流程**：`forge init` 的完整流程、`.forge/config.yaml` 全字段含义、Surface 检测机制等关键初始化知识没有独立文档
- **架构概览**：`docs/ARCHITECTURE.md`（555行）面向开发者，包含内部实现细节（Go 包结构、CLI 内部机制），终端用户难以从中提取"我需要了解什么"
- **使用指南**：README.md 的"5分钟体验"仅两条命令，缺少 Full Mode / Quick Mode 的端到端实战示例和常见问题排查

新用户需要从至少 3 个不同文件中拼凑信息才能完成从安装到使用的全过程。

### Evidence

- README.md 安装部分（第 46-71 行）仅 25 行，覆盖 marketplace 和 local 两种安装方式，但无环境验证、无排错指引
- `docs/ARCHITECTURE.md` 面向开发者，包含 `forge-cli/internal/cmd/` 包引用、`ResolveScope` 实现细节等内部信息
- `docs/reference/` 目录仅包含 `test-type-model.md`，无环境配置或初始化文档
- 多个用户反馈（项目 memory 中记录）表明现有文档不完整

### Urgency

Forge v3.0.0 即将发布，新版本引入了 CLI 命令分组、Surface 自动检测、Worktree 管理等大量新功能，用户需要一个完整的文档体系来理解和使用这些能力。发布前补充文档是最佳时机。

## Proposed Solution

创建 `docs/user-guide/` 目录，包含四个面向终端用户的独立文档：

1. **`environment-setup.md`** — 前置条件、安装方式、环境验证、常见安装问题
2. **`initialization.md`** — `forge init` 完整流程、config.yaml 全字段参考、Surface 检测、首个项目设置
3. **`architecture-overview.md`** — 用户视角的插件机制、组件角色（skill/command/agent/hook）、数据流向、状态管理、目录约定
4. **`usage-guide.md`** — Full Mode / Quick Mode 端到端实战、单命令场景、常见问题与排错

同时在 README.md 的文档索引表中添加用户手册链接。

### Innovation Highlights

方案本身是标准的多文件结构化文档实践，无特殊创新。选择此方案的核心原因是职责分离原则：将用户文档与开发者文档（ARCHITECTURE.md）物理隔离，各自面向不同读者群体独立维护。

## Requirements Analysis

### Key Scenarios

- 新用户首次安装 Forge，从零开始配置环境并完成第一个 feature
- 已安装用户升级到 v3.0.0，需要了解新的 config 选项和 Surface 检测
- 用户遇到问题（安装失败、配置错误、工作流异常），需要排错指引
- 用户想了解 Forge 内部如何运作（插件机制、数据流向），以便更好地使用

### Non-Functional Requirements

- 所有文档使用中文，与项目 CLAUDE.md 要求一致
- 每个文档独立可读，不依赖其他文档的上下文
- 代码示例可直接复制执行，无需额外修改

### Constraints & Dependencies

- 文档内容来源于现有代码和配置，需要确保与 v3.0.0 代码一致
- architecture-overview.md 不包含开发者内部实现细节（Go 包结构、CLI 内部命令注册等）

## Alternatives & Industry Benchmarking

### Industry Solutions

主流开源项目通常采用以下文档结构：
- **Next.js**：`docs/` 下按主题分目录（getting-started, app, pages），每目录多个 mdx 文件
- **Astro**：类似结构，独立 getting-started 指南
- **Tauri**：单页快速上手 + 分专题指南

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零成本 | 新用户困惑，文档缺口持续 | Rejected: v3.0.0 发布前需要完整文档 |
| 扩展 README.md | 内部 | 集中在一个文件 | README 过长（当前已 452 行），两种读者混合 | Rejected: 违反职责分离 |
| 扩展 ARCHITECTURE.md | 内部 | 不新增文件 | 面向开发者的文档混入用户内容 | Rejected: 读者群体不匹配 |
| **多文件 user-guide/** | 行业标准 | 职责单一、按需查阅、独立维护 | 需要维护更多文件 | **Selected: 最佳职责分离** |

## Feasibility Assessment

### Technical Feasibility

纯文档任务，无技术风险。所有信息来源于现有代码和配置文件。

### Resource & Timeline

4 个文档文件 + README 更新，预计可在单次 /quick 流水线中完成。

### Dependency Readiness

所有信息已存在于代码库中，无外部依赖。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| "ARCHITECTURE.md 对用户也适用" | XY Detection | Overturned: ARCHITECTURE.md 面向开发者，包含 Go 包结构等内部信息，用户需要的是"怎么用"而非"怎么实现" |
| "README.md 已经够用" | 5 Whys | Refined: README 覆盖了安装和命令参考，但缺少初始化详解、排错指引和端到端实战 |

## Scope

### In Scope

- `docs/user-guide/environment-setup.md` — 前置条件、安装方式、环境验证、常见安装问题
- `docs/user-guide/initialization.md` — forge init 流程、config.yaml 全字段参考、Surface 检测、首个项目设置
- `docs/user-guide/architecture-overview.md` — 用户视角架构（插件机制、组件角色、数据流、状态管理、目录约定）
- `docs/user-guide/usage-guide.md` — Full Mode / Quick Mode 端到端实战、单命令场景、常见问题与排错
- README.md 文档索引表添加用户手册链接

### Out of Scope

- 修改现有 `docs/ARCHITECTURE.md`（保持开发者视角不变）
- Forge CLI 内部实现文档（Go 包结构、内部 API）
- 贡献者指南（README.md 已有）
- API 文档（`forge -h` 已覆盖）
- 国际化（仅中文）

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 文档内容与 v3.0.0 代码不同步 | L | M | 内容来源于代码和配置文件的实际读取，非凭记忆 |
| 文档过时（后续版本变更后未更新） | M | L | 每个文档顶部标注"最后更新"日期和对应版本 |

## Success Criteria

- [ ] `docs/user-guide/` 目录包含 4 个独立的 Markdown 文件
- [ ] `environment-setup.md` 覆盖 3 种安装方式（marketplace / local / 开发模式）和安装后验证命令
- [ ] `initialization.md` 包含 config.yaml 全字段表格（至少 8 个配置项）和 Surface 检测说明
- [ ] `architecture-overview.md` 包含插件机制图解和数据流向说明，不包含 Go 包结构等开发者内部信息
- [ ] `usage-guide.md` 包含 Full Mode 和 Quick Mode 各至少一个端到端示例，以及 5 条以上常见问题
- [ ] README.md 文档索引表包含 `docs/user-guide/` 的 4 个文件链接

## Next Steps

- Proceed to `/quick-tasks` to generate implementation tasks
