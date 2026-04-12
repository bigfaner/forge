# ZCode Web Dashboard — 产品需求文档（PRD）

---

## 元数据

| 字段 | 内容 |
|------|------|
| 功能名称 | ZCode Web Dashboard |
| 文档版本 | v0.1.0 |
| 创建日期 | 2026-04-08 |
| 状态 | 已确认 |
| 负责人 | — |

---

## 一、背景与动机

ZCode 当前完全依赖 CLI 和 Markdown 文件驱动工作流（PRD → Design → Tasks → Execution）。虽然该方案对 AI Agent 友好，但存在以下痛点：

1. **可视化缺失**：任务依赖关系、执行进度、覆盖率等信息分散在多个 JSON / Markdown 文件中，人工阅读成本高
2. **操作门槛**：执行 `task claim` / `task status` 等操作需要记忆命令，不直观
3. **跨 Feature 全局视图缺失**：无法一眼看到所有 Feature 的整体进度
4. **知识沉淀不可检索**：`docs/lessons/` 中的经验文件只能逐一打开阅读

### 目标

提供一个**本地 Web Dashboard**，在不破坏现有 CLI 工作流的前提下，为人类开发者提供可视化管理界面，支持：
- 浏览和管理所有 Feature 的任务状态
- 可视化任务依赖关系（DAG）
- 直接触发核心 CLI 操作（claim、status 变更、record 提交）
- 检索经验知识库

---

## 二、目标用户

| 用户 | 场景 |
|------|------|
| 使用 ZCode 的个人开发者 | 本地运行 `task serve`，浏览器管理自己的任务 |
| 使用 ZCode 的小型团队 | 共享本地网络下的 Dashboard，协同查看任务进度 |

---

## 三、范围

### In Scope

- Feature 列表与进度概览
- Feature 详情（PRD 文档、技术设计文档、任务看板）
- 任务看板（Kanban 视图、列表视图、DAG 依赖图）
- 任务详情与执行记录查看
- 核心操作：Claim Task、更新任务状态、提交执行记录
- 执行记录全局时间线
- 经验知识库（Lessons）浏览与搜索
- 本地运行，通过 `task serve` 命令启动

### Out of Scope（本版本不包含）

- 用户认证 / 权限管理
- 多机器远程访问
- PRD / Design 文档的在线编辑
- 实时推送（WebSocket / SSE）
- 移动端适配

---

## 四、用户故事

### 功能组 A：全局概览

**A1** 作为开发者，我想在 Dashboard 首页看到所有 Feature 的状态摘要（进行中/已完成/阻塞），以便快速了解整体进度。

**A2** 作为开发者，我想看到当前 in_progress 的任务入口，以便快速恢复上次中断的工作。

**A3** 作为开发者，我想看到最近的执行记录 feed，以便了解近期产出。

### 功能组 B：Feature 管理

**B1** 作为开发者，我想浏览所有 Feature 列表，查看每个 Feature 的任务总数和完成进度。

**B2** 作为开发者，我想查看某个 Feature 的 PRD 和技术设计文档（Markdown 渲染），以便了解需求背景。

**B3** 作为开发者，我想切换 Feature 的任务视图（看板 / 列表 / DAG 依赖图），以便从不同角度理解任务结构。

### 功能组 C���任务操作

**C1** 作为开发者，我想点击"Claim Task"按钮认领下一个可用任务，而不需要手动执行 `task claim`。

**C2** 作为开发者，我想在任务详情页直接修改任务状态（blocked / skipped），而不需要记忆命令格式。

**C3** 作为开发者，我想查看任务的执行记录（summary、decisions、coverage、文件变更），了解历史执行情况。

**C4** 作为开发者，我想在 DAG 视图中清晰看到任务的依赖关系，以便理解执行顺序。

### 功能组 D：知识库

**D1** 作为开发者，我想浏览所有 Lessons 文件，按分类（debug / arch / tool / pattern / gotcha）筛选。

**D2** 作为开发者，我想通过关键词搜索 Lessons，快速定位相关经验。

---

## 五、功能需求

### 5.1 `task serve` 命令

| 需求 | 说明 |
|------|------|
| 启动本地 HTTP Server | 默认端口 7300，支持 `--port` 参数覆盖 |
| 自动检测项目根目录 | 复用 `project.FindProjectRoot()` 逻辑 |
| 静态文件服务 | 将前端构建产物 embed 进二进制，`/` 路由服务前端 |
| 启动提示 | 打印访问地址：`Dashboard running at http://localhost:7300` |

### 5.2 REST API

#### Feature 接口

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/features` | 返回所有 Feature 列表，含任务统计（total / completed / in_progress / blocked） |
| GET | `/api/features/:slug` | 返回 Feature 元数据 + 任务列表 |
| GET | `/api/features/:slug/prd` | 返回 PRD Markdown 原文 |
| GET | `/api/features/:slug/design` | 返回 Design Markdown 原文 |
| GET | `/api/features/:slug/tasks` | 返回任务列表（含状态、优先级、依赖项） |
| GET | `/api/features/:slug/tasks/:id` | 返回单个任务详情 |
| GET | `/api/features/:slug/records` | 返回该 Feature 所有执行记录列表 |

#### 任务操作接口

| 方法 | 路径 | Body | 说明 |
|------|------|------|------|
| POST | `/api/tasks/claim` | — | 认领下一个可用任务，返回 TaskState |
| POST | `/api/tasks/:id/status` | `{ "status": "blocked" }` | 更新任务状态 |
| POST | `/api/tasks/:id/record` | RecordData JSON | 提交执行记录，更新状态为 completed |

#### 知识库接口

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/lessons` | 返回所有 Lesson 文件元数据列表（名称、分类、摘要） |
| GET | `/api/lessons/:name` | 返回单篇 Lesson Markdown 原文 |

#### 系统接口

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/health` | 返回 task-cli 版本、项目根目录、当前 Feature |

### 5.3 前端页面

#### 页面 1：Dashboard 首页（`/`）

- Feature 状态卡片网格（进行中 / 已完成 / 阻塞 数量统计）
- 全局任务进度环形图
- 当前 in_progress 任务快捷入口卡片
- 最近 10 条执行记录 feed（任务名 + 时间 + coverage）

#### 页面 2：Feature 列表（`/features`）

- Feature 卡片列表（slug、标题、任务进度条、最后活跃时间）
- 按状态筛选（all / active / completed）

#### 页面 3：Feature 详情（`/features/:slug`）

三个 Tab：

- **PRD Tab**：Markdown 渲染 `prd.md`，顶部显示文档状态标记
- **Design Tab**：Markdown 渲染 `design.md`，顶部显示文档状态标记
- **Tasks Tab**：含视图切换
  - Kanban 视图：按状态分 5 列（pending / in_progress / completed / blocked / skipped）
  - 列表视图：表格，支持按 phase / priority / status 排序和筛选
  - DAG 视图：ReactFlow 渲染任务依赖有向无环图，节点着色区分状态
  - 顶部操作栏："Claim Task" 按钮（触发 POST /api/tasks/claim）

#### 页面 4：任务详情（`/features/:slug/tasks/:id`）

- 任务基本信息（title、description、phase、priority、estimated_time、files）
- 依赖项列表（可点击跳转）
- 状态操作（下拉菜单：blocked / skipped，保存后调用 status API）
- 执行记录面板（summary、decisions、test results、coverage、filesCreated、filesModified、commit hash）

#### 页面 5：执行记录（`/records`）

- 跨所有 Feature 的时间线视图
- 每条记录显示：Feature slug、任务 ID、任务标题、coverage、文件变更数、commit hash、完成时间
- 按 Feature / 日期范围筛选

#### 页面 6：知识库（`/lessons`）

- 文章列表（名称、分类徽章、第一行摘要）
- 分类 Tab 筛选（all / debug / arch / tool / pattern / gotcha）
- 关键词搜索（前端本地过滤）
- 点击展开 Markdown 全文

#### 页面 7：设置（`/settings`）

- task-cli 版本信息（来自 `/api/health`）
- 项目根目录显示
- 当前 Feature 上下文显示
- Server 端口信息

---

## 六、非功能需求

| 类别 | 要求 |
|------|------|
| 性能 | 页面首屏加载 < 1s（本地文件读取） |
| 依赖 | 前端构建产物 embed 进 task-cli 二进制，零额外运行时依赖 |
| 兼容性 | 支持 Chrome / Safari / Firefox 最新版 |
| 安全 | 仅监听 localhost，不对外暴露（无需认证） |
| UI 主题 | 支持深色 / 浅色模式切换（shadcn/ui 内置支持） |
| 数据刷新 | 前端定时轮询（默认每 30s），无需 WebSocket |

---

## 七、技术栈

### 后端（Go）

| 组件 | 方案 |
|------|------|
| HTTP Server | `net/http` 标准库（零新依赖） |
| 静态文件 | `embed.FS` 内嵌 `web/dist/` |
| 业务逻辑 | 直接复用 `pkg/task`、`pkg/feature`、`pkg/project` |
| 新增目录 | `task-cli/internal/server/`（handlers + routing） |

### 前端

| 组件 | 方案 | 理由 |
|------|------|------|
| 框架 | React 18 + TypeScript | 生态成熟 |
| 构建工具 | Vite | 快速，产物体积小 |
| UI 组件 | shadcn/ui + Tailwind CSS | 无运行时依赖 |
| 路由 | React Router v6 | 标准选择 |
| 服务端状态 | TanStack Query v5 | 自动缓存 + 刷新，适配 REST |
| DAG 可视化 | ReactFlow | 任务依赖图渲染 |
| Markdown 渲染 | react-markdown + remark-gfm | 渲染 PRD/Design 文档 |

---

## 八、目录结构规划

```
web/                              ← 本文件夹，前端工程根目录
├── prd.md                        ← 本文档
├── design.md                     ← 技术设计文档（待写）
├── src/
│   ├── main.tsx
│   ├── App.tsx                   ← 路由配置
│   ├── layouts/
│   │   └── AppLayout.tsx         ← 侧边栏 + Header
│   ├── pages/
│   │   ├── Dashboard.tsx
│   │   ├── FeatureList.tsx
│   │   ├── FeatureDetail.tsx
│   │   │   ├── PrdTab.tsx
│   │   │   ├── DesignTab.tsx
│   │   │   └── TasksTab.tsx
│   │   ├── TaskDetail.tsx
│   │   ├── Records.tsx
│   │   ├── Lessons.tsx
│   │   └── Settings.tsx
│   ├── components/
│   │   ├── TaskBoard/
│   │   │   ├── KanbanView.tsx
│   │   │   ├── ListView.tsx
│   │   │   └── DagView.tsx
│   │   ├── TaskCard.tsx
│   │   ├── MarkdownViewer.tsx
│   │   ├── FeatureCard.tsx
│   │   ├── RecordTimeline.tsx
│   │   └── StatusBadge.tsx
│   └── lib/
│       ├── api.ts                ← API 请求封装
│       └── types.ts              ← TypeScript 类型定义
├── public/
├── package.json
├── vite.config.ts
├── tsconfig.json
└── dist/                         ← 构建产物（被 Go embed）
```

---

## 九、依赖关系与里程碑

```
里程碑 1：后端 API（Go）
  └─ task serve 命令 + HTTP Server + 所有 API handlers

里程碑 2：前端基础框架
  └─ Vite 脚手架 + Layout + 路由 + API 封装 + TypeScript 类型

里程碑 3：核心页面
  └─ Dashboard + Feature 列表 + Feature 详情（PRD/Design/Tasks Tab）

里程碑 4：任务交互
  └─ Claim Task + 状态变更 + 任务详情 + 执行记录

里程碑 5：辅助页面
  └─ 全局执行记录 + 知识库 + 设置页
```

---

## 十、已确认决策

| # | 问题 | 决策 |
|---|------|------|
| Q1 | DAG 视图的布局算法？ | ✅ 层次布局（按 phase 分层） |
| Q2 | 执行记录是否支持在 Dashboard 手动填写表单提交？ | ✅ 只读，不支持——记录只能通过 CLI `task record` 提交 |
| Q3 | 是否需要深色模式？ | ✅ 需要——支持深色 / 浅色切换 |
| Q4 | `task serve` 是否需要实时刷新？ | ✅ 前端定时轮询（无需 WebSocket / SSE） |
