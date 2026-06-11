---
created: 2026-05-20
author: "faner"
status: Draft
---

# Proposal: Forge Dashboard — 多项目状态管理与可视化

## Problem

Forge 以单项目模式运作，管理 3-5 个项目时需要频繁 `cd` 到不同目录、执行 `forge task status` 才能了解每个项目的进度。没有全局视野，无法一眼看清所有项目的健康状态、活跃任务和近期活动。

### Evidence

- 3-5 个活跃项目分布在不同目录，每次查看状态需要手动切换
- 现有 Web Dashboard（docs/features/web-dashboard/）仅支持单项目，且未部署
- 项目间的 lessons、decisions、conventions 无法跨项目查看
- 没有统一的"所有项目今天发生了什么"视图

### Urgency

项目数量增长后，全局视野的缺失导致管理成本非线性上升。早期解决可避免信息碎片化加剧。

## Proposed Solution

浏览器端多项目 Dashboard，通过 `forge dashboard` 命令启动本地 API 服务并自动打开浏览器。提供项目状态总览、任务下钻、活动时间线三个核心视图。

**项目发现**：扫描指定目录下所有包含 `.forge/` 的子项目 + 用户手动注册的项目路径。

**命令入口**：`forge dashboard` 一条命令启动服务和浏览器。

### Innovation Highlights

**以项目为一等公民**：现有 Forge 以 feature 为核心组织，Dashboard 以项目为核心组织视图。两个视角互补——Feature 视角用于深入开发，Project 视角用于全局管理。

**零配置启动**：无需手动配置项目列表，目录扫描自动发现。`forge dashboard` 一条命令即可使用。

## Requirements Analysis

### Key Scenarios

1. **全局状态一览**：打开 Dashboard 看到所有项目的 feature 状态、任务进度、健康度
2. **项目下钻**：点击某个项目，查看其活跃 feature、待办任务、近期记录
3. **活动时间线**：按时间倒序查看所有项目的活动记录（任务完成、决策记录、lesson 新增等）
4. **快速跳转**：从 Dashboard 中一键打开某个项目的终端或编辑器

### Non-Functional Requirements

- Dashboard 加载时间 < 3 秒（本地数据，无网络延迟）
- 支持 3-10 个项目的规模
- 浏览器兼容 Chrome/Edge/Firefox 最新版本

### Constraints & Dependencies

- 依赖 Forge CLI 的项目发现能力
- 需要本地 API 服务（HTTP server）
- 前端框架待技术设计阶段确定
- 数据来源为各项目 `docs/features/` 目录的只读访问

## Alternatives & Industry Benchmarking

### Industry Solutions

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零成本 | 多项目管理持续混乱 | Rejected: 问题随项目增长恶化 |
| CLI 全局视图 | `forge overview` 命令 | 无需 Web 服务 | 表格输出无法展示富内容，交互性差 | Partial: 可作为补充，但不是主要方案 |
| 扩展现有 Dashboard | docs/features/web-dashboard/ | 复用已有工作 | 现有 Dashboard 未完成且为单项目设计，扩展架构债务大 | Rejected: 用户决定重新设计 |
| **全新多项目 Dashboard** | — | 干净架构，多项目优先 | 需要重新开发 | **Selected: 架构干净，符合多项目优先的设计理念** |

## Feasibility Assessment

### Technical Feasibility

- 现有 `tmp_api_server.mjs` 可作为 API 层的参考（但需重新设计多项目架构）
- 前端技术成熟，无技术风险
- 项目发现逻辑可通过 CLI 已有的文件系统操作实现

### Resource & Timeline

- API 服务 + 前端界面为中等规模工作量
- 可参考现有 Dashboard 的 UI 设计（Tailwind + 卡片布局）

### Dependency Readiness

- 无外部服务依赖
- 各项目 `docs/features/index.json` 已有结构化数据可直接消费

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| "需要 Web UI 才能管理多项目" | XY Detection | 状态查看也可以用 CLI 解决，但活动时间线和知识浏览确实需要富 UI |
| "应该复用现有 Dashboard" | Occam's Razor | 现有 Dashboard 未完成且架构不同，重新设计比在债务上扩展更简单 |
| "项目发现需要手动注册" | Assumption Flip | 目录扫描可自动发现大部分项目，手动注册只作为补充 |

## Scope

### In Scope

- 项目发现与注册（目录扫描 + 手动路径）
- 本地 API 服务（读取多项目状态数据）
- Dashboard UI — 项目状态总览
- Dashboard UI — 项目下钻（feature、任务、记录）
- Dashboard UI — 活动时间线
- `forge dashboard` CLI 命令（启动服务 + 打开浏览器）

### Out of Scope

- 知识库管理（属于 Forge Wiki 提案）
- 知识库聚合和搜索
- 跨项目任务管理
- 多用户 / 团队协作
- Dashboard 内编辑操作（v1 为只读）
- 原生桌面应用
- 移动端适配

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 项目目录结构差异大，状态解析失败 | M | H | 定义标准化的项目状态接口，异常项目优雅降级 |
| API 服务端口冲突 | L | L | 可配置端口，默认 7300 |
| 前端框架选择增加维护成本 | L | M | 技术设计时评估，选择最小依赖方案 |
| 项目数量多时 Dashboard 性能下降 | L | M | 按需加载数据，懒加载项目详情 |

## Success Criteria

- [ ] `forge dashboard` 一条命令启动并打开浏览器，显示所有已发现项目
- [ ] 项目状态总览页显示每个项目的 feature 数量、任务进度、健康状态
- [ ] 项目下钻页显示活跃 feature 的任务列表和执行记录
- [ ] 活动时间线按时间倒序展示跨项目活动
- [ ] 支持 5 个项目的数据加载，页面响应 < 3 秒

## Next Steps

- Proceed to `/write-prd` to formalize requirements
