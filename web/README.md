# ZCode Web Dashboard

ZCode 任务管理系统的前端 Web Dashboard，基于 React + TypeScript + Vite 构建。

---

## 入口文件

| 文件 | 说明 |
|------|------|
| [`index.html`](./index.html) | HTML 入口，Vite 从此文件启动 |
| [`src/main.tsx`](./src/main.tsx) | React 应用挂载点，初始化 QueryClient 和主题 |
| [`src/App.tsx`](./src/App.tsx) | 路由配置，定义 7 个页面路径 |

---

## 快速启动

### 安装依赖

```bash
cd web
npm install --cache /tmp/npm-cache   # 若 npm 缓存有权限问题用此命令
# 或修复权限后正常安装：
sudo chown -R $(whoami) ~/.npm && npm install
```

---

### 方式一：仅查看 UI 效果（无数据）

启动静态预览，API 请求会失败但页面框架可正常浏览：

```bash
cd web
npm run preview -- --port 4173
```

访问：**http://localhost:4173**

> 侧边栏导航、深色/浅色模式切换、路由跳转均可正常使用，数据区域显示空状态。

---

### 方式二：查看完整数据效果（Mock Server）

**终端 1** — 启动 Mock API Server（无需 Go 环境）：

```bash
# 在项目根目录运行
node mock-server.js
# Mock API 运行在 http://localhost:7300
```

**终端 2** — 启动前端开发服务器：

```bash
cd web
npm run dev -- --port 4173
```

访问：**http://localhost:4173**

Mock 数据包含：
- 3 个 Feature（auth-login / task-dashboard / api-gateway）
- 39 个任务，涵盖全部 5 种状态
- 11 条执行记录时间线
- 6 篇知识库文章

---

### 方式三：连接真实后端

启动 Go 后端（需安装 Go 环境）：

```bash
cd task-cli
go build -o task ./cmd/task && ./task serve --port 7300
```

然后启动前端开发服务器，`vite.config.ts` 已配置 `/api` 请求自动代理到 `localhost:7300`：

```bash
cd web
npm run dev -- --port 4173
```

---

## 目录结构

```
web/
├── index.html                    # HTML 入口
├── prd.md                        # 产品需求文档
├── design.md                     # 技术设计文档
├── design-system/
│   └── zcode-dashboard/
│       └── MASTER.md             # UI 设计系统（颜色/字体/风格）
├── src/
│   ├── main.tsx                  # React 挂载 + QueryClient 初始化
│   ├── App.tsx                   # 路由配置
│   ├── index.css                 # 全局样式 + CSS 变量（深/浅色 token）
│   ├── layouts/
│   │   └── AppLayout.tsx         # 侧边栏导航 + Header + ThemeToggle
│   ├── pages/                    # 页面组件（对应路由）
│   │   ├── Dashboard.tsx         # /         全局统计概览
│   │   ├── FeatureList.tsx       # /features Feature 卡片列表
│   │   ├── FeatureDetail.tsx     # /features/:slug  详情（PRD/Design/Tasks）
│   │   ├── TaskDetail.tsx        # /features/:slug/tasks/:id 任务详情
│   │   ├── Records.tsx           # /records  执行记录时间线
│   │   ├── Lessons.tsx           # /lessons  知识库
│   │   └── Settings.tsx          # /settings 服务器配置信息
│   ├── components/               # 可复用组件
│   │   ├── task-board/
│   │   │   ├── KanbanView.tsx    # 看板视图（按状态分列）
│   │   │   ├── ListView.tsx      # 表格视图（可排序）
│   │   │   └── DagView.tsx       # DAG 依赖图（ReactFlow 层次布局）
│   │   ├── ClaimButton.tsx       # 认领任务按钮
│   │   ├── FeatureCard.tsx       # Feature 卡片（含进度条）
│   │   ├── MarkdownViewer.tsx    # Markdown 渲染器
│   │   ├── StatusBadge.tsx       # 状态/优先级徽章
│   │   └── ThemeToggle.tsx       # 深色/浅色切换
│   └── lib/                      # 工具库
│       ├── api.ts                # 所有 API 请求封装
│       ├── types.ts              # TypeScript 类型定义
│       ├── utils.ts              # 工具函数（cn、状态颜色映射等）
│       └── theme.ts              # 主题管理（localStorage 持久化）
├── package.json
├── vite.config.ts                # Vite 配置（含 /api 代理）
├── tailwind.config.ts            # Tailwind 配置（含深色模式）
├── tsconfig.json
└── dist/                         # 构建产物（被 Go embed，不提交 git）
```

---

## API 接口清单

接口封装位于 [`src/lib/api.ts`](./src/lib/api.ts)，所有请求代理到 `http://localhost:7300/api`。

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/health` | 服务器版本、项目路径、当前 Feature |
| GET | `/api/features` | 所有 Feature 列表（含任务统计） |
| GET | `/api/features/:slug` | Feature 详情 + 任务列表 |
| GET | `/api/features/:slug/prd` | PRD Markdown 原文 |
| GET | `/api/features/:slug/design` | Design Markdown 原文 |
| GET | `/api/features/:slug/tasks` | 任务列表 |
| GET | `/api/features/:slug/tasks/:id` | 单个任务详情（含执行记录） |
| GET | `/api/features/:slug/records` | 该 Feature 的执行记录 |
| GET | `/api/records` | 全局执行记录时间线 |
| GET | `/api/lessons` | 知识库文章列表 |
| GET | `/api/lessons/:name` | 单篇知识库 Markdown 原文 |
| POST | `/api/tasks/claim` | 认领下一个可用任务 |
| POST | `/api/features/:slug/tasks/:id/status` | 更新任务状态 |

---

## 构建

```bash
cd web
npm run build
# 产物输出到 web/dist/，由 task-cli Go 二进制 embed
```
