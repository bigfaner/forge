# ZCode Web Dashboard — 技术设计文档

---

## 元数据

| 字段 | 内容 |
|------|------|
| 功能名称 | ZCode Web Dashboard |
| 文档版本 | v0.1.0 |
| 创建日期 | 2026-04-08 |
| 状态 | 已确认 |
| 依赖文档 | [prd.md](./prd.md) |

---

## 一、整体架构

```
┌─────────────────────────────────────────────────────┐
│                   浏览器                              │
│  React + TypeScript + TanStack Query                 │
│  shadcn/ui + Tailwind + ReactFlow                    │
└───────────────────┬─────────────────────────────────┘
                    │ HTTP  (localhost:7300)
┌───────────────────▼─────────────────────────────────┐
│              task serve  (Go)                        │
│  ┌─────────────────┐   ┌──────────────────────────┐ │
│  │  Static Handler │   │     API Router            │ │
│  │  embed web/dist │   │  /api/**  → handlers/     │ │
│  └─────────────────┘   └──────────┬───────────────┘ │
│                                   │                  │
│  ┌────────────────────────────────▼───────────────┐ │
│  │            pkg/  (复用现有逻辑)                  │ │
│  │  pkg/task   pkg/feature   pkg/project  pkg/git  │ │
│  └────────────────────────────────────────────────┘ │
│                                   │                  │
│                        本地文件系统                    │
│              docs/features/*/tasks/index.json        │
│              docs/features/*/prd.md                  │
│              docs/features/*/design.md               │
│              docs/features/*/tasks/records/*.md      │
│              docs/lessons/*.md                       │
└─────────────────────────────────────────────────────┘
```

---

## 二、后端设计

### 2.1 新增目录结构

```
task-cli/
├── internal/
│   ├── cmd/
│   │   └── serve.go              ← 新增：注册 task serve 命令
│   └── server/
│       ├── server.go             ← HTTP 路由注册、embed 静态文件
│       └── handlers/
│           ├── features.go       ← Feature 相关接口
│           ├── tasks.go          ← Task 操作接口
│           ├── records.go        ← 执行记录接口
│           ├── lessons.go        ← 知识库接口
│           └── health.go         ← 健康检查接口
└── web/dist/                     ← 前端构建产物（被 embed）
```

### 2.2 `serve.go` 命令定义

```go
var serveCmd = &cobra.Command{
    Use:   "serve",
    Short: "Start the web dashboard",
    RunE:  runServe,
}

func init() {
    serveCmd.Flags().IntP("port", "p", 7300, "Port to listen on")
    rootCmd.AddCommand(serveCmd)
}

func runServe(cmd *cobra.Command, args []string) error {
    port, _ := cmd.Flags().GetInt("port")
    root, err := project.FindProjectRoot()
    if err != nil {
        return err
    }
    fmt.Printf("Dashboard running at http://localhost:%d\n", port)
    return server.Start(root, port)
}
```

### 2.3 `server.go` 路由注册

```go
//go:embed ../web/dist
var staticFiles embed.FS

func Start(projectRoot string, port int) error {
    mux := http.NewServeMux()

    // API 路由
    mux.Handle("/api/features",            handlers.Features(projectRoot))
    mux.Handle("/api/features/",           handlers.FeaturesDetail(projectRoot))
    mux.Handle("/api/tasks/claim",         handlers.ClaimTask(projectRoot))
    mux.Handle("/api/tasks/",              handlers.TasksOps(projectRoot))
    mux.Handle("/api/lessons",             handlers.Lessons(projectRoot))
    mux.Handle("/api/lessons/",            handlers.LessonDetail(projectRoot))
    mux.Handle("/api/health",              handlers.Health(projectRoot))

    // 静态文件（SPA fallback：所有非 /api 路径返回 index.html）
    mux.Handle("/", spaHandler(staticFiles))

    return http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
}
```

### 2.4 API 响应结构

#### `GET /api/features`

```json
{
  "features": [
    {
      "slug": "auth-login",
      "title": "用户登录功能",
      "stats": {
        "total": 12,
        "pending": 3,
        "in_progress": 1,
        "completed": 7,
        "blocked": 1,
        "skipped": 0
      },
      "lastUpdated": "2026-04-08T10:30:00Z"
    }
  ]
}
```

#### `GET /api/features/:slug/tasks`

```json
{
  "feature": "auth-login",
  "title": "用户登录功能",
  "tasks": [
    {
      "id": "1.1",
      "title": "定义用户数据模型",
      "description": "...",
      "phase": 1,
      "priority": "P0",
      "status": "completed",
      "estimatedTime": "2h",
      "dependencies": [],
      "files": ["pkg/model/user.go"],
      "record": "tasks/records/1.1-user-model.md"
    }
  ]
}
```

#### `GET /api/features/:slug/tasks/:id`

```json
{
  "id": "1.2",
  "title": "实现登录接口",
  "description": "...",
  "phase": 1,
  "priority": "P0",
  "status": "completed",
  "estimatedTime": "3h",
  "dependencies": ["1.1"],
  "files": ["internal/handler/auth.go", "internal/handler/auth_test.go"],
  "record": {
    "summary": "实现了 POST /api/auth/login 接口",
    "filesCreated": ["internal/handler/auth.go"],
    "filesModified": ["internal/handler/auth_test.go"],
    "decisions": ["使用 JWT，有效期 24h"],
    "testResults": "PASS 12/12",
    "coverage": "87.3%",
    "commitHash": "a3f9c12"
  }
}
```

#### `POST /api/tasks/claim` 响应

```json
{
  "taskId": "2.1",
  "key": "auth-login",
  "title": "实现 Token 刷新逻辑",
  "file": "docs/features/auth-login/tasks/2.1-token-refresh.md"
}
```

#### `GET /api/lessons`

```json
{
  "lessons": [
    {
      "name": "debug-go-embed-path",
      "category": "debug",
      "title": "Go embed 路径问题排查",
      "excerpt": "使用 embed.FS 时路径不包含前缀目录..."
    }
  ]
}
```

#### `GET /api/health`

```json
{
  "version": "1.2.0",
  "projectRoot": "/Users/nasuki/my-project",
  "currentFeature": "auth-login"
}
```

---

## 三、前端设计

### 3.1 路由结构

```
/                              → Dashboard（首页）
/features                      → Feature 列表
/features/:slug                → Feature 详情（默认 Tasks Tab）
/features/:slug?tab=prd        → Feature PRD Tab
/features/:slug?tab=design     → Feature Design Tab
/features/:slug?tab=tasks      → Feature Tasks Tab
/features/:slug/tasks/:id      → 任务详情
/records                       → 执行记录时间线
/lessons                       → 知识库
/settings                      → 设置
```

### 3.2 组件树

```
App
└── AppLayout
    ├── Sidebar
    │   ├── NavItem (Dashboard)
    │   ├── NavItem (Features)
    │   ├── NavItem (Records)
    │   ├── NavItem (Lessons)
    │   └── NavItem (Settings)
    ├── Header
    │   ├── FeatureSwitcher        ← 当前 Feature 切换下拉框
    │   └── ThemeToggle            ← 深色/浅色切换
    └── <Outlet>                   ← 页面内容区
        │
        ├── Dashboard
        │   ├── StatsGrid          ← Feature 状态卡片 (4格)
        │   ├── GlobalProgressRing ← 全局任务环形图
        │   ├── InProgressCard     ← 当前任务快捷入口
        │   └── RecentRecordsFeed  ← 最近10条执行记录
        │
        ├── FeatureList
        │   ├── FilterBar          ← all/active/completed 筛选
        │   └── FeatureCard[]      ← 每个 Feature 卡片
        │       ├── ProgressBar
        │       └── StatsRow
        │
        ├── FeatureDetail
        │   ├── TabNav             ← PRD / Design / Tasks
        │   ├── PrdTab
        │   │   └── MarkdownViewer
        │   ├── DesignTab
        │   │   └── MarkdownViewer
        │   └── TasksTab
        │       ├── TasksToolbar   ← 视图切换 + Claim Task 按钮
        │       ├── KanbanView
        │       │   └── KanbanColumn[] → TaskCard[]
        │       ├── ListView
        │       │   └── TaskTable  ← 可排序表格
        │       └── DagView
        │           └── ReactFlow  ← 层次布局有向图
        │
        ├── TaskDetail
        │   ├── TaskMeta           ← title/phase/priority/files
        │   ├── DependencyList     ← 可跳转依赖项
        │   ├── StatusSelector     ← blocked/skipped 下拉操作
        │   └── RecordPanel        ← 执行记录（只读）
        │
        ├── Records
        │   ├── RecordsFilter      ← Feature/日期范围筛选
        │   └── RecordTimeline     ← 时间线列表
        │
        ├── Lessons
        │   ├── CategoryTabs       ← all/debug/arch/tool/pattern/gotcha
        │   ├── SearchBar          ← 前端本地过滤
        │   └── LessonList
        │       └── LessonItem     ← 可展开 MarkdownViewer
        │
        └── Settings
            ├── HealthInfo         ← 版本/路径/Feature
            └── PollingControl     ← 轮询间隔设置
```

### 3.3 TypeScript 类型定义（`src/lib/types.ts`）

```typescript
// --- Feature ---
export interface FeatureStats {
  total: number
  pending: number
  in_progress: number
  completed: number
  blocked: number
  skipped: number
}

export interface FeatureSummary {
  slug: string
  title: string
  stats: FeatureStats
  lastUpdated: string
}

// --- Task ---
export type TaskStatus = 'pending' | 'in_progress' | 'completed' | 'blocked' | 'skipped'
export type TaskPriority = 'P0' | 'P1' | 'P2'

export interface Task {
  id: string
  title: string
  description: string
  phase: number
  priority: TaskPriority
  status: TaskStatus
  estimatedTime: string
  dependencies: string[]
  files: string[]
  record?: string
}

export interface TaskDetail extends Task {
  record?: TaskRecord
}

// --- Record ---
export interface TaskRecord {
  summary: string
  filesCreated: string[]
  filesModified: string[]
  decisions: string[]
  testResults: string
  coverage: string
  commitHash: string
  completedAt: string
}

// --- Lesson ---
export type LessonCategory = 'debug' | 'arch' | 'tool' | 'pattern' | 'gotcha'

export interface LessonMeta {
  name: string
  category: LessonCategory
  title: string
  excerpt: string
}

// --- Health ---
export interface HealthInfo {
  version: string
  projectRoot: string
  currentFeature: string
}
```

### 3.4 API 请求封装（`src/lib/api.ts`）

```typescript
const BASE = '/api'

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(`${BASE}${path}`, init)
  if (!res.ok) throw new Error(await res.text())
  return res.json()
}

export const api = {
  features: {
    list: () => request<{ features: FeatureSummary[] }>('/features'),
    tasks: (slug: string) => request<{ tasks: Task[] }>(`/features/${slug}/tasks`),
    task: (slug: string, id: string) => request<TaskDetail>(`/features/${slug}/tasks/${id}`),
    prd: (slug: string) => request<string>(`/features/${slug}/prd`),
    design: (slug: string) => request<string>(`/features/${slug}/design`),
    records: (slug: string) => request<{ records: TaskRecord[] }>(`/features/${slug}/records`),
  },
  tasks: {
    claim: () => request<ClaimResult>('/tasks/claim', { method: 'POST' }),
    setStatus: (id: string, status: TaskStatus) =>
      request(`/tasks/${id}/status`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ status }),
      }),
  },
  lessons: {
    list: () => request<{ lessons: LessonMeta[] }>('/lessons'),
    get: (name: string) => request<string>(`/lessons/${name}`),
  },
  health: () => request<HealthInfo>('/health'),
}
```

### 3.5 状态管理策略

使用 **TanStack Query v5**，遵循以下规则：

| 数据类型 | staleTime | refetchInterval | 说明 |
|---------|-----------|-----------------|------|
| Feature 列表 | 10s | 30s | 轮询感知任务变化 |
| 任务列表 | 5s | 30s | 执行中频繁变化 |
| PRD / Design Markdown | 60s | 不轮询 | 文档变化少 |
| 执行记录 | 10s | 30s | 新记录需感知 |
| Lessons | 60s | 不轮询 | 几乎不变 |
| Health | 30s | 60s | 低频 |

**Query Key 规范**：
```typescript
['features']                          // Feature 列表
['features', slug, 'tasks']           // 某 Feature 的任务列表
['features', slug, 'tasks', id]       // 单个任务详情
['features', slug, 'prd']             // PRD 内容
['features', slug, 'design']          // Design 内容
['lessons']                           // Lesson 列表
['health']                            // 健康信息
```

**Mutation 后的缓存失效**：
- `claim` 成功 → invalidate `['features', slug, 'tasks']`
- `setStatus` 成功 → invalidate `['features', slug, 'tasks']` + `['features', slug, 'tasks', id]`

### 3.6 深色模式

使用 **Tailwind CSS `dark:` 变体 + shadcn/ui 内置主题**：

```typescript
// src/lib/theme.ts
type Theme = 'light' | 'dark' | 'system'

export function setTheme(theme: Theme) {
  const root = document.documentElement
  if (theme === 'dark' || (theme === 'system' && prefersDark())) {
    root.classList.add('dark')
  } else {
    root.classList.remove('dark')
  }
  localStorage.setItem('theme', theme)
}
```

- 默认跟随系统（`system`）
- `ThemeToggle` 组件持久化到 `localStorage`
- Tailwind 配置：`darkMode: 'class'`

### 3.7 DAG 视图设计

使用 **ReactFlow** + **Dagre** 层次布局算法：

```
Phase 1        Phase 2        Phase 3
┌─────┐        ┌─────┐        ┌─────┐
│ 1.1 │───────▶│ 2.1 │───────▶│ 3.1 │
└─────┘        └─────┘        └─────┘
               ┌─────┐
               │ 2.2 │
               └─────┘
```

节点着色规则：

| 状态 | 颜色（浅色） | 颜色（深色） |
|------|------------|------------|
| pending | gray-200 | gray-700 |
| in_progress | blue-200 | blue-800 |
| completed | green-200 | green-800 |
| blocked | red-200 | red-800 |
| skipped | yellow-200 | yellow-800 |

---

## 四、构建与集成

### 4.1 目录约定

```
task-cli/web/          ← 前端工程（与 task-cli 同仓库）
task-cli/web/dist/     ← 构建产物（被 Go embed，不提交到 git）
```

`.gitignore` 新增：
```
task-cli/web/dist/
task-cli/web/node_modules/
```

### 4.2 Makefile 构建流程

```makefile
# 前端构建
.PHONY: web
web:
	cd web && npm install && npm run build

# Go 构建（依赖前端产物）
.PHONY: build
build: web
	go build -o task ./cmd/task

# 仅开发时（前端 dev server 代理到 Go Server）
.PHONY: dev
dev:
	task serve &
	cd web && npm run dev
```

### 4.3 前端开发代理

`vite.config.ts` 配置代理，开发时不需要 embed：

```typescript
export default defineConfig({
  server: {
    proxy: {
      '/api': 'http://localhost:7300',
    },
  },
})
```

---

## 五、文件结构汇总

```
web/
├── prd.md                        ← 产品需求文档
├── design.md                     ← 本文档
├── src/
│   ├── main.tsx                  ← 入口，挂载 QueryClient + Router
│   ├── App.tsx                   ← 路由配置
│   ├── layouts/
│   │   └── AppLayout.tsx         ← 侧边栏 + Header + Outlet
│   ├── pages/
│   │   ├── Dashboard.tsx
│   │   ├── FeatureList.tsx
│   │   ├── FeatureDetail.tsx     ← Tab 路由容器
│   │   ├── PrdTab.tsx
│   │   ├── DesignTab.tsx
│   │   ├── TasksTab.tsx          ← 视图切换容器
│   │   ├── TaskDetail.tsx
│   │   ├── Records.tsx
│   │   ├── Lessons.tsx
│   │   └── Settings.tsx
│   ├── components/
│   │   ├── task-board/
│   │   │   ├── KanbanView.tsx
│   │   │   ├── ListView.tsx
│   │   │   └── DagView.tsx
│   │   ├── TaskCard.tsx
│   │   ├── MarkdownViewer.tsx
│   │   ├── FeatureCard.tsx
│   │   ├── RecordTimeline.tsx
│   │   ├── StatusBadge.tsx
│   │   └── ThemeToggle.tsx
│   └── lib/
│       ├── api.ts
│       ├── types.ts
│       └── theme.ts
├── public/
│   └── favicon.ico
├── package.json
├── vite.config.ts
├── tsconfig.json
├── tailwind.config.ts
└── components.json               ← shadcn/ui 配置
```

---

## 六、里程碑与任务拆分

### 里程碑 1：后端 API（Go）

| 任务 | 文件 | 优先级 |
|------|------|--------|
| 实现 `task serve` 命令 | `internal/cmd/serve.go` | P0 |
| 实现 HTTP Server + embed | `internal/server/server.go` | P0 |
| Feature 列表接口 | `handlers/features.go` | P0 |
| Feature 任务列表 + 详情接口 | `handlers/tasks.go` | P0 |
| PRD / Design 文档接口 | `handlers/features.go` | P1 |
| Claim / Status 操作接口 | `handlers/tasks.go` | P1 |
| 执行记录接口 | `handlers/records.go` | P1 |
| Lessons 接口 | `handlers/lessons.go` | P2 |
| Health 接口 | `handlers/health.go` | P2 |

### 里程碑 2：前端基础框架

| 任务 | 文件 | 优先级 |
|------|------|--------|
| Vite + React + TS 脚手架 | `package.json` / `vite.config.ts` | P0 |
| shadcn/ui + Tailwind 初始化 | `tailwind.config.ts` / `components.json` | P0 |
| AppLayout（侧边栏 + Header） | `layouts/AppLayout.tsx` | P0 |
| 路由配置 | `App.tsx` | P0 |
| API 封装 + TypeScript 类型 | `lib/api.ts` / `lib/types.ts` | P0 |
| TanStack Query 初始化 | `main.tsx` | P0 |
| 深色模式 ThemeToggle | `components/ThemeToggle.tsx` / `lib/theme.ts` | P1 |

### 里程碑 3：核心页面

| 任务 | 文件 | 优先级 |
|------|------|--------|
| Dashboard 首页 | `pages/Dashboard.tsx` | P0 |
| Feature 列表页 | `pages/FeatureList.tsx` | P0 |
| Feature 详情 Tab 容器 | `pages/FeatureDetail.tsx` | P0 |
| PRD / Design Tab | `pages/PrdTab.tsx` / `DesignTab.tsx` | P0 |
| Tasks Kanban 视图 | `components/task-board/KanbanView.tsx` | P0 |
| Tasks 列表视图 | `components/task-board/ListView.tsx` | P1 |
| Tasks DAG 视图 | `components/task-board/DagView.tsx` | P1 |

### 里程碑 4：任务交互

| 任务 | 文件 | 优先级 |
|------|------|--------|
| Claim Task 按钮逻辑 | `pages/TasksTab.tsx` | P0 |
| 任务详情页 | `pages/TaskDetail.tsx` | P0 |
| 状态变更操作 | `pages/TaskDetail.tsx` | P1 |
| 执行记录面板（只读） | `pages/TaskDetail.tsx` | P1 |

### 里程碑 5：辅助页面

| 任务 | 文件 | 优先级 |
|------|------|--------|
| 执行记录全局时间线 | `pages/Records.tsx` | P1 |
| 知识库页 + 搜索 | `pages/Lessons.tsx` | P1 |
| 设置页 | `pages/Settings.tsx` | P2 |
