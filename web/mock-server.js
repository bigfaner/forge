#!/usr/bin/env node
// Mock API Server — runs on port 7300, proxied by Vite dev server
import { createServer } from 'http'

const FEATURES = [
  {
    slug: 'auth-login',
    title: '用户认证与登录系统',
    lastUpdated: '2026-04-07T10:23:00Z',
    tasks: [
      { id: 'T-001', title: '设计认证方案', description: '调研 JWT vs Session，确定技术选型，输出 ADR 文档。', phase: 1, priority: 'P0', status: 'completed', estimatedTime: '2h', dependencies: [], files: ['docs/adr/auth-strategy.md'], record: 'rec-001' },
      { id: 'T-002', title: '实现 JWT 签发与验证中间件', description: '封装 issueToken / verifyToken，支持 RS256 算法，过期时间可配置。', phase: 1, priority: 'P0', status: 'completed', estimatedTime: '4h', dependencies: ['T-001'], files: ['pkg/auth/jwt.go', 'pkg/auth/jwt_test.go'], record: 'rec-002' },
      { id: 'T-003', title: '用户注册接口', description: 'POST /api/register，密码 bcrypt 加盐，唯一性校验，返回 201。', phase: 1, priority: 'P0', status: 'completed', estimatedTime: '3h', dependencies: ['T-002'], files: ['internal/handler/register.go'], record: 'rec-003' },
      { id: 'T-004', title: '用户登录接口', description: 'POST /api/login，验证凭据后下发 access_token + refresh_token。', phase: 1, priority: 'P0', status: 'in_progress', estimatedTime: '3h', dependencies: ['T-003'], files: ['internal/handler/login.go'] },
      { id: 'T-005', title: 'Refresh Token 轮换', description: '实现无感刷新，旧 token 使用后立即作废，防止重放攻击。', phase: 2, priority: 'P1', status: 'pending', estimatedTime: '4h', dependencies: ['T-004'], files: ['internal/handler/refresh.go'] },
      { id: 'T-006', title: '登出与 Token 黑名单', description: 'Redis 维护撤销 token 集合，过期自动清理。', phase: 2, priority: 'P1', status: 'pending', estimatedTime: '2h', dependencies: ['T-005'], files: ['internal/handler/logout.go', 'pkg/cache/redis.go'] },
      { id: 'T-007', title: 'OAuth2 GitHub 第三方登录', description: '接入 GitHub OAuth App，回调后绑定或创建本地账户。', phase: 2, priority: 'P2', status: 'pending', estimatedTime: '6h', dependencies: ['T-004'], files: ['internal/handler/oauth_github.go'] },
      { id: 'T-008', title: '频率限制（Rate Limit）', description: '登录接口每 IP 每分钟限 10 次，超限返回 429。', phase: 2, priority: 'P1', status: 'blocked', estimatedTime: '2h', dependencies: ['T-004'], files: ['internal/middleware/ratelimit.go'] },
      { id: 'T-009', title: '前端登录表单', description: 'React 表单，字段校验，错误提示，loading 状态。', phase: 3, priority: 'P1', status: 'pending', estimatedTime: '4h', dependencies: ['T-004'], files: ['web/src/pages/Login.tsx', 'web/src/components/LoginForm.tsx'] },
      { id: 'T-010', title: '端到端测试', description: 'Playwright 覆盖注册→登录→刷新→登出全流程。', phase: 3, priority: 'P1', status: 'skipped', estimatedTime: '5h', dependencies: ['T-009'], files: ['e2e/auth.spec.ts'] },
      { id: 'T-011', title: '安全审计与渗透测试', description: 'OWASP Top10 自查清单，SQL 注入、XSS、CSRF 验证。', phase: 3, priority: 'P0', status: 'pending', estimatedTime: '8h', dependencies: ['T-008', 'T-009'], files: ['docs/security-audit.md'] },
      { id: 'T-012', title: '文档与 Changelog', description: '更新 API 文档，补充认证流程图，发布 v1.0.0 说明。', phase: 3, priority: 'P2', status: 'pending', estimatedTime: '2h', dependencies: ['T-011'], files: ['docs/api.md', 'CHANGELOG.md'] },
    ]
  },
  {
    slug: 'task-dashboard',
    title: '任务管理看板（本项目）',
    lastUpdated: '2026-04-08T08:00:00Z',
    tasks: [
      { id: 'T-001', title: '定义数据模型与类型系统', description: '梳理 Feature / Task / Record / Lesson 数据结构，输出 TypeScript 类型文件。', phase: 1, priority: 'P0', status: 'completed', estimatedTime: '2h', dependencies: [], files: ['web/src/lib/types.ts'], record: 'rec-004' },
      { id: 'T-002', title: '搭建 Vite + React + Tailwind 脚手架', description: '初始化项目，配置路由、React Query、主题系统。', phase: 1, priority: 'P0', status: 'completed', estimatedTime: '3h', dependencies: ['T-001'], files: ['web/package.json', 'web/vite.config.ts', 'web/src/main.tsx'], record: 'rec-005' },
      { id: 'T-003', title: '实现 API 客户端层', description: '封装 fetch wrapper，统一错误处理，支持文本与 JSON 响应。', phase: 1, priority: 'P0', status: 'completed', estimatedTime: '2h', dependencies: ['T-002'], files: ['web/src/lib/api.ts'], record: 'rec-006' },
      { id: 'T-004', title: 'AppLayout 与导航栏', description: '侧边栏导航，响应式折叠，亮/暗主题切换。', phase: 1, priority: 'P1', status: 'completed', estimatedTime: '3h', dependencies: ['T-003'], files: ['web/src/layouts/AppLayout.tsx', 'web/src/components/ThemeToggle.tsx'], record: 'rec-007' },
      { id: 'T-005', title: 'Dashboard 总览页', description: '聚合统计卡片，活跃 / 已完成 Feature 分组展示。', phase: 2, priority: 'P0', status: 'completed', estimatedTime: '3h', dependencies: ['T-004'], files: ['web/src/pages/Dashboard.tsx'], record: 'rec-008' },
      { id: 'T-006', title: 'FeatureList 列表页', description: 'Feature 卡片网格，状态徽章，按 Active / Completed 筛选。', phase: 2, priority: 'P1', status: 'completed', estimatedTime: '2h', dependencies: ['T-004'], files: ['web/src/pages/FeatureList.tsx', 'web/src/components/FeatureCard.tsx'], record: 'rec-009' },
      { id: 'T-007', title: 'FeatureDetail — Kanban / List / DAG 视图', description: '三种任务视图切换，DAG 依赖图用 @xyflow/react 渲染。', phase: 2, priority: 'P0', status: 'in_progress', estimatedTime: '8h', dependencies: ['T-006'], files: ['web/src/pages/FeatureDetail.tsx', 'web/src/components/task-board/KanbanView.tsx', 'web/src/components/task-board/DagView.tsx'] },
      { id: 'T-008', title: 'TaskDetail 详情页', description: '展示任务元信息、文件列表、执行记录，支持状态变更。', phase: 2, priority: 'P1', status: 'pending', estimatedTime: '4h', dependencies: ['T-007'], files: ['web/src/pages/TaskDetail.tsx'] },
      { id: 'T-009', title: 'Records 时间线页', description: '全局执行记录时间轴，按时间倒序，展示 commit / coverage。', phase: 3, priority: 'P1', status: 'pending', estimatedTime: '3h', dependencies: ['T-005'], files: ['web/src/pages/Records.tsx'] },
      { id: 'T-010', title: 'Lessons 知识库页', description: '分类展开列表，全文搜索，Markdown 渲染。', phase: 3, priority: 'P1', status: 'pending', estimatedTime: '4h', dependencies: ['T-005'], files: ['web/src/pages/Lessons.tsx'] },
      { id: 'T-011', title: 'Mock Server 脚本', description: '用 Node.js 内置 http 模块实现完整 mock API，覆盖所有端点。', phase: 3, priority: 'P2', status: 'in_progress', estimatedTime: '2h', dependencies: ['T-003'], files: ['mock-server.js'] },
      { id: 'T-012', title: 'Claim Task 功能', description: '从未领取任务中随机分配一个，写入 in_progress 状态。', phase: 3, priority: 'P2', status: 'blocked', estimatedTime: '2h', dependencies: ['T-008'], files: ['web/src/components/ClaimButton.tsx'] },
      { id: 'T-013', title: '性能优化与懒加载', description: '路由级代码分割，React Query staleTime 调优，虚拟滚动大列表。', phase: 4, priority: 'P2', status: 'pending', estimatedTime: '4h', dependencies: ['T-009', 'T-010'], files: ['web/src/main.tsx'] },
    ]
  },
  {
    slug: 'api-gateway',
    title: 'API 网关与流量治理',
    lastUpdated: '2026-04-05T16:45:00Z',
    tasks: [
      { id: 'T-001', title: '网关架构选型', description: '评估 Kong / Envoy / 自研方案，输出 ADR。', phase: 1, priority: 'P0', status: 'completed', estimatedTime: '4h', dependencies: [], files: ['docs/adr/gateway-arch.md'], record: 'rec-010' },
      { id: 'T-002', title: '路由转发核心', description: '基于 Host / Path 的动态路由，支持前缀剥离与重写。', phase: 1, priority: 'P0', status: 'completed', estimatedTime: '8h', dependencies: ['T-001'], files: ['pkg/gateway/router.go', 'pkg/gateway/router_test.go'], record: 'rec-011' },
      { id: 'T-003', title: '服务注册与健康检查', description: 'etcd 作为注册中心，主动 HTTP 探针，故障自动摘除。', phase: 2, priority: 'P0', status: 'in_progress', estimatedTime: '6h', dependencies: ['T-002'], files: ['pkg/registry/etcd.go', 'pkg/health/probe.go'] },
      { id: 'T-004', title: '负载均衡策略', description: '轮询 / 加权轮询 / 最少连接，可运行时切换。', phase: 2, priority: 'P1', status: 'pending', estimatedTime: '4h', dependencies: ['T-003'], files: ['pkg/lb/balancer.go'] },
      { id: 'T-005', title: '熔断器（Circuit Breaker）', description: 'Hystrix 语义，三态状态机，自动半开探测。', phase: 2, priority: 'P0', status: 'pending', estimatedTime: '6h', dependencies: ['T-004'], files: ['pkg/circuitbreaker/breaker.go'] },
      { id: 'T-006', title: '请求限流（Token Bucket）', description: '全局 + 每服务双层限流，支持 Burst。', phase: 2, priority: 'P1', status: 'blocked', estimatedTime: '4h', dependencies: ['T-003'], files: ['pkg/ratelimit/tokenbucket.go'] },
      { id: 'T-007', title: 'Prometheus 指标暴露', description: 'QPS / 延迟分位 / 错误率，标准 /metrics 端点。', phase: 3, priority: 'P1', status: 'skipped', estimatedTime: '3h', dependencies: ['T-002'], files: ['pkg/metrics/prometheus.go'] },
      { id: 'T-008', title: '灰度发布（Canary）', description: '按权重/Header 路由到新版本，支持逐步切流。', phase: 3, priority: 'P2', status: 'pending', estimatedTime: '8h', dependencies: ['T-004'], files: ['pkg/gateway/canary.go'] },
      { id: 'T-009', title: 'mTLS 服务间认证', description: '自签 CA 签发证书，双向 TLS 握手验证。', phase: 3, priority: 'P0', status: 'pending', estimatedTime: '6h', dependencies: ['T-002'], files: ['pkg/tls/mtls.go'] },
      { id: 'T-010', title: '压测与调优', description: 'k6 模拟 5000 RPS，P99 < 50ms，找出并修复瓶颈。', phase: 4, priority: 'P1', status: 'pending', estimatedTime: '8h', dependencies: ['T-005', 'T-006'], files: ['load-test/gateway.js', 'docs/perf-report.md'] },
      { id: 'T-011', title: '生产部署 Runbook', description: 'K8s Helm Chart，滚动升级，回滚流程文档化。', phase: 4, priority: 'P1', status: 'pending', estimatedTime: '6h', dependencies: ['T-010'], files: ['deploy/helm/gateway/', 'docs/runbook.md'] },
      { id: 'T-012', title: '告警规则配置', description: 'Alertmanager 规则：错误率 > 1% / P99 > 200ms 触发告警。', phase: 4, priority: 'P1', status: 'pending', estimatedTime: '2h', dependencies: ['T-007'], files: ['deploy/alerts/gateway.yaml'] },
      { id: 'T-013', title: '文档站点', description: '用 Docusaurus 构建 API 文档站，集成 OpenAPI Spec。', phase: 4, priority: 'P2', status: 'pending', estimatedTime: '6h', dependencies: ['T-011'], files: ['docs-site/'] },
      { id: 'T-014', title: 'Admin 控制台 UI', description: '实时流量图、路由配置、服务健康状态一览。', phase: 4, priority: 'P2', status: 'pending', estimatedTime: '12h', dependencies: ['T-007', 'T-008'], files: ['admin/'] },
    ]
  }
]

const RECORDS = [
  { featureSlug: 'auth-login', taskId: 'T-001', taskTitle: '设计认证方案', coverage: '—', filesChanged: 1, commitHash: 'a1b2c3d', completedAt: '2026-04-01T09:15:00Z' },
  { featureSlug: 'auth-login', taskId: 'T-002', taskTitle: '实现 JWT 签发与验证中间件', coverage: '94.2%', filesChanged: 2, commitHash: 'e4f5g6h', completedAt: '2026-04-02T14:30:00Z' },
  { featureSlug: 'auth-login', taskId: 'T-003', taskTitle: '用户注册接口', coverage: '88.5%', filesChanged: 1, commitHash: 'i7j8k9l', completedAt: '2026-04-03T11:00:00Z' },
  { featureSlug: 'task-dashboard', taskId: 'T-001', taskTitle: '定义数据模型与类型系统', coverage: '—', filesChanged: 1, commitHash: 'm0n1o2p', completedAt: '2026-04-04T10:00:00Z' },
  { featureSlug: 'task-dashboard', taskId: 'T-002', taskTitle: '搭建 Vite + React + Tailwind 脚手架', coverage: '—', filesChanged: 4, commitHash: 'q3r4s5t', completedAt: '2026-04-04T16:00:00Z' },
  { featureSlug: 'task-dashboard', taskId: 'T-003', taskTitle: '实现 API 客户端层', coverage: '—', filesChanged: 1, commitHash: 'u6v7w8x', completedAt: '2026-04-05T09:30:00Z' },
  { featureSlug: 'task-dashboard', taskId: 'T-004', taskTitle: 'AppLayout 与导航栏', coverage: '—', filesChanged: 2, commitHash: 'y9z0a1b', completedAt: '2026-04-05T14:00:00Z' },
  { featureSlug: 'task-dashboard', taskId: 'T-005', taskTitle: 'Dashboard 总览页', coverage: '—', filesChanged: 1, commitHash: 'c2d3e4f', completedAt: '2026-04-06T10:00:00Z' },
  { featureSlug: 'task-dashboard', taskId: 'T-006', taskTitle: 'FeatureList 列表页', coverage: '—', filesChanged: 2, commitHash: 'g5h6i7j', completedAt: '2026-04-07T09:00:00Z' },
  { featureSlug: 'api-gateway', taskId: 'T-001', taskTitle: '网关架构选型', coverage: '—', filesChanged: 1, commitHash: 'k8l9m0n', completedAt: '2026-04-02T11:00:00Z' },
  { featureSlug: 'api-gateway', taskId: 'T-002', taskTitle: '路由转发核心', coverage: '91.8%', filesChanged: 2, commitHash: 'o1p2q3r', completedAt: '2026-04-05T17:00:00Z' },
]

const TASK_RECORDS = {
  'rec-001': { summary: '完成 JWT vs Session 对比分析，选定 JWT + Redis 黑名单方案，输出 ADR 文档。', filesCreated: ['docs/adr/auth-strategy.md'], filesModified: [], decisions: ['采用 RS256 非对称签名保障多服务场景安全', '过期时间 access=15m，refresh=7d'], testResults: '无单元测试（纯文档任务）', coverage: '—', commitHash: 'a1b2c3d', completedAt: '2026-04-01T09:15:00Z' },
  'rec-002': { summary: '封装 issueToken / verifyToken，RS256 签名，支持 claims 扩展，编写完整单元测试。', filesCreated: ['pkg/auth/jwt.go', 'pkg/auth/jwt_test.go'], filesModified: [], decisions: ['密钥从环境变量注入，不落磁盘', '提供 ParseUnverified 用于日志追踪'], testResults: '18/18 passed', coverage: '94.2%', commitHash: 'e4f5g6h', completedAt: '2026-04-02T14:30:00Z' },
  'rec-003': { summary: '实现 POST /api/register，bcrypt cost=12，邮箱唯一性 DB 约束，400/409 错误结构统一。', filesCreated: ['internal/handler/register.go'], filesModified: ['internal/router/router.go'], decisions: ['返回 201 而非 200 符合 REST 语义', '不在响应中返回密码 hash'], testResults: '12/12 passed', coverage: '88.5%', commitHash: 'i7j8k9l', completedAt: '2026-04-03T11:00:00Z' },
  'rec-004': { summary: '梳理所有实体关系，输出 TypeScript 类型文件，与后端 JSON schema 对齐。', filesCreated: ['web/src/lib/types.ts'], filesModified: [], decisions: ['TaskRecord 支持 string | object 兼容旧格式', 'RecordEntry 独立于 Task 以支持全局时间线'], testResults: '—', coverage: '—', commitHash: 'm0n1o2p', completedAt: '2026-04-04T10:00:00Z' },
  'rec-005': { summary: '初始化 Vite 5 + React 18 + Tailwind 3，配置 path alias，React Query v5 Provider，Router v6。', filesCreated: ['web/package.json', 'web/vite.config.ts', 'web/src/main.tsx', 'web/src/App.tsx'], filesModified: [], decisions: ['使用 @tanstack/react-query v5 代替 SWR，缓存控制更精细', 'Tailwind 优先于 CSS-in-JS，减少运行时开销'], testResults: '—', coverage: '—', commitHash: 'q3r4s5t', completedAt: '2026-04-04T16:00:00Z' },
  'rec-006': { summary: '封装 fetch wrapper，统一错误抛出，根据 Content-Type 自动判断 JSON/文本响应。', filesCreated: ['web/src/lib/api.ts'], filesModified: [], decisions: ['不使用 axios，减少包体积', 'text/markdown 与 text/plain 统一走 res.text()'], testResults: '—', coverage: '—', commitHash: 'u6v7w8x', completedAt: '2026-04-05T09:30:00Z' },
  'rec-007': { summary: '实现侧边栏导航，移动端折叠，亮/暗主题系统（localStorage 持久化 + prefers-color-scheme 初始值）。', filesCreated: ['web/src/layouts/AppLayout.tsx', 'web/src/components/ThemeToggle.tsx', 'web/src/lib/theme.ts'], filesModified: ['web/src/index.css'], decisions: ['主题用 data-theme attribute 而非 class，避免与 Tailwind dark: 冲突', '侧边栏宽度用 CSS variable 便于动画'], testResults: '—', coverage: '—', commitHash: 'y9z0a1b', completedAt: '2026-04-05T14:00:00Z' },
  'rec-008': { summary: '总览页聚合 4 个统计卡片，分 Active / Completed 两区展示 Feature，30s 自动刷新。', filesCreated: ['web/src/pages/Dashboard.tsx'], filesModified: [], decisions: ['统计卡片数字用 CSS counter 动画提升视觉', 'refetchInterval 设 30000ms 避免服务器压力过大'], testResults: '—', coverage: '—', commitHash: 'c2d3e4f', completedAt: '2026-04-06T10:00:00Z' },
  'rec-009': { summary: 'Feature 卡片网格，All / Active / Completed 三种筛选，StatusBadge 组件封装。', filesCreated: ['web/src/pages/FeatureList.tsx', 'web/src/components/FeatureCard.tsx', 'web/src/components/StatusBadge.tsx'], filesModified: [], decisions: ['筛选逻辑在前端做，避免额外 API 调用', 'FeatureCard 复用于 Dashboard'], testResults: '—', coverage: '—', commitHash: 'g5h6i7j', completedAt: '2026-04-07T09:00:00Z' },
  'rec-010': { summary: '评估 Kong / Envoy / 自研三方案，综合考量团队能力与运维复杂度，选定自研轻量网关。', filesCreated: ['docs/adr/gateway-arch.md'], filesModified: [], decisions: ['Kong 功能强大但引入 PostgreSQL 依赖复杂度过高', 'Envoy xDS 学习曲线陡峭，自研可控性更强'], testResults: '—', coverage: '—', commitHash: 'k8l9m0n', completedAt: '2026-04-02T11:00:00Z' },
  'rec-011': { summary: '实现基于 radix tree 的路由匹配，支持 path 前缀剥离、Header 转发、超时配置，benchmark 1.2M rps。', filesCreated: ['pkg/gateway/router.go', 'pkg/gateway/router_test.go'], filesModified: ['pkg/gateway/config.go'], decisions: ['radix tree 比 map 在路由数 >50 时性能提升 3x', 'timeout 默认 30s，可每路由覆盖'], testResults: '34/34 passed', coverage: '91.8%', commitHash: 'o1p2q3r', completedAt: '2026-04-05T17:00:00Z' },
}

const LESSONS = [
  { name: 'jwt-rs256-pitfall', category: 'gotcha', title: 'JWT RS256 验签：公钥路径错误导致静默失败', excerpt: '使用 RS256 时若公钥文件路径错误，某些库不会抛出而是返回 nil error，所有 token 均通过验证。' },
  { name: 'react-query-v5-migration', category: 'tool', title: 'React Query v5 Breaking Changes 踩坑记录', excerpt: 'v5 移除了 onSuccess/onError 回调，useQuery 返回值结构变更，升级时需全量替换。' },
  { name: 'radix-tree-routing', category: 'arch', title: '基于 Radix Tree 的高性能路由设计', excerpt: '当路由规则超过 50 条时，Radix Tree 比 map 前缀匹配性能提升约 3 倍，内存占用减少 40%。' },
  { name: 'vite-proxy-cors', category: 'debug', title: 'Vite Dev Proxy 与 CORS 的常见误区', excerpt: 'Vite proxy 只在 dev server 生效，production 构建后需在 Nginx/网关层配置 CORS，不能依赖 changeOrigin。' },
  { name: 'bcrypt-cost-tradeoff', category: 'pattern', title: 'bcrypt cost 参数的性能与安全权衡', excerpt: 'cost=10 约 100ms，cost=12 约 400ms。注册/登录接口可接受，但批量导入用户时需降低或改用离线处理。' },
  { name: 'circuit-breaker-halfopen', category: 'arch', title: '熔断器半开状态的探测策略选择', excerpt: '半开时放行单个请求探测比放行百分比更可预测；探测成功后建议渐进恢复流量而非立即全量放开。' },
]

const LESSON_CONTENTS = {
  'jwt-rs256-pitfall': `# JWT RS256 验签：公钥路径错误导致静默失败\n\n**分类：** gotcha | **发现于：** auth-login T-002\n\n## 问题描述\n\n使用 RS256 算法时，若传入 \`ParseRSAPublicKeyFromPEM\` 的字节为空（公钥文件路径错误），\n某些 Go JWT 库不会抛出而是返回 \`nil error\`，导致**所有 token 均通过验证**。\n\n## 复现步骤\n\n\`\`\`go\n// 错误示范\npubKeyBytes, _ := os.ReadFile(os.Getenv("PUBLIC_KEY_PATH"))\nkey, _ := jwt.ParseRSAPublicKeyFromPEM(pubKeyBytes) // key == nil, err == nil\n\`\`\`\n\n\`\`\`go\n// 正确做法\npubKeyBytes, err := os.ReadFile(os.Getenv("PUBLIC_KEY_PATH"))\nif err != nil || len(pubKeyBytes) == 0 {\n    log.Fatal("public key not found:", err)\n}\n\`\`\`\n\n## 预防措施\n\n启动时做 key pair 自检（用测试 token 走一遍签发→验签流程）。`,
  'react-query-v5-migration': `# React Query v5 Breaking Changes 踩坑记录\n\n**分类：** tool | **发现于：** task-dashboard T-002\n\n## 主要变更\n\n### 1. 移除 onSuccess / onError 回调\n\n\`\`\`tsx\n// ❌ v4（已移除）\nuseQuery({ queryKey, queryFn, onSuccess: (data) => toast(data) })\n\n// ✅ v5\nconst { data } = useQuery({ queryKey, queryFn })\nuseEffect(() => { if (data) toast(data) }, [data])\n\`\`\`\n\n### 2. isLoading → isPending\n\n### 3. cacheTime → gcTime`,
  'radix-tree-routing': `# 基于 Radix Tree 的高性能路由设计\n\n**分类：** arch | **发现于：** api-gateway T-002\n\n## 背景\n\n初版路由用 \`map[string]Handler\` 实现，100 条规则时 P99 延迟 0.8ms。\n切换 Radix Tree 后降至 0.27ms，内存从 2.1MB 降至 1.2MB。\n\n## 何时不需要\n\n路由规则 < 20 条时，简单 slice 线性扫描更易维护，差异可忽略不计。`,
  'vite-proxy-cors': `# Vite Dev Proxy 与 CORS 的常见误区\n\n**分类：** debug | **发现于：** task-dashboard T-003\n\n## 误区\n\n开发时配了 \`vite.config.ts\` proxy，CORS 正常，部署后报错。\n\n## 原因\n\nVite proxy 是**开发服务器**功能，build 后的静态文件里没有 proxy 逻辑。\n\n## 生产环境做法\n\n用 Nginx 反向代理将 \`/api/*\` 代理到后端，或后端添加 CORS 中间件。`,
  'bcrypt-cost-tradeoff': `# bcrypt cost 参数的性能与安全权衡\n\n**分类：** pattern | **发现于：** auth-login T-003\n\n| cost | 耗时 | 适用场景 |\n|------|------|----------|\n| 10 | ~90ms | 高并发低安全 |\n| 12 | ~380ms | **推荐** |\n| 14 | ~1500ms | 高安全场景 |\n\n**注意**：批量导入用户时需异步队列处理，避免超时。`,
  'circuit-breaker-halfopen': `# 熔断器半开状态的探测策略选择\n\n**分类：** arch | **发现于：** api-gateway T-005（设计阶段）\n\n## 推荐：单请求探测\n\n半开后只放行 1 个请求，成功则关闭熔断，失败则重新打开。\n\n**优点**：可预测，不会对下游造成额外压力。\n\n探测成功后不要立即恢复 100% 流量，建议在 30s 内线性恢复。`,
}

const PRD = {
  'auth-login': `# PRD：用户认证与登录系统\n\n**版本：** v1.2 | **状态：** In Progress\n\n## 用户故事\n\n| ID | 需求 | 优先级 |\n|----|------|--------|\n| US-1 | ��过邮箱+密码注册 | P0 |\n| US-2 | 登录获取 JWT | P0 |\n| US-3 | 无感刷新 token | P1 |\n| US-4 | GitHub 一键登录 | P2 |\n\n## 非功能需求\n\n- 登录接口 P99 < 500ms\n- 日志脱敏：不记录明文密码\n- 连续失败 5 次锁定 15 分钟`,
  'task-dashboard': `# PRD：任务管理看板\n\n**版本：** v0.1 | **状态：** Active Development\n\n## 核心功能\n\n1. **总览页**：所有 Feature 进度聚合\n2. **Feature 详情**：Kanban / List / DAG 三视图\n3. **任务详情**：执行记录 + 决策 + commit 链接\n4. **知识库**：沉淀工程经验，可搜索分类\n5. **时间线**：跨 Feature 执行历史\n\n## 成功指标\n\n查看任务状态时间从 3min 降至 30s。`,
  'api-gateway': `# PRD：API 网关与流量治理\n\n**版本：** v2.0 | **状态：** Phase 2\n\n## 核心能力\n\n| 能力 | 目标 | 状态 |\n|------|------|------|\n| 路由转发 | P99 < 1ms | ✅ 完成 |\n| 服务注册 | 故障摘除 < 5s | 🔄 进行中 |\n| 熔断限流 | 错误率 < 0.1% | ⏳ 待开始 |\n| 灰度发布 | 支持 1% 切流 | ⏳ 待开始 |\n\n## SLA\n\n可用性 99.95%，单机 5000 RPS，P99 < 50ms。`,
}

const DESIGN = {
  'auth-login': `# 技术设计：用户认证\n\n## 认证流程\n\n\`\`\`\nClient → POST /login → 验证凭据 → 签发 JWT → 返回 token\nClient → GET /api/* (Bearer token) → 验证 JWT → 响应\n\`\`\`\n\n## 数据库表\n\n\`\`\`sql\nCREATE TABLE users (\n  id UUID PRIMARY KEY,\n  email VARCHAR(255) UNIQUE NOT NULL,\n  password_hash VARCHAR(60) NOT NULL\n);\n\nCREATE TABLE refresh_tokens (\n  token_hash CHAR(64) PRIMARY KEY,\n  user_id UUID REFERENCES users(id),\n  expires_at TIMESTAMPTZ NOT NULL\n);\n\`\`\``,
  'task-dashboard': `# 技术设计：任务管理看板\n\n## 组件树\n\n\`\`\`\nApp\n├── AppLayout\n│   ├── Dashboard\n│   ├── FeatureList → FeatureDetail\n│   │   └── KanbanView / ListView / DagView\n│   ├── TaskDetail\n│   ├── Records\n│   ├── Lessons\n│   └── Settings\n\`\`\`\n\n## 主题系统\n\nCSS custom properties + \`data-theme\` attribute，初始值从 localStorage 读取。`,
  'api-gateway': `# 技术设计：API 网关\n\n## 架构\n\n\`\`\`\nInternet → [Gateway] → Router → Auth → RateLimit → CircuitBreaker → Upstream\n\`\`\`\n\n## 熔断状态机\n\n\`\`\`\nClosed ──(错误率>50%)──▶ Open\n  ▲                         │ (30s后)\n  │ 探测成功            Half-Open\n  └──────────────────────────┘\n\`\`\``,
}

function computeStats(tasks) {
  const s = { total: tasks.length, pending: 0, in_progress: 0, completed: 0, blocked: 0, skipped: 0 }
  for (const t of tasks) s[t.status]++
  return s
}

function json(res, data, status = 200) {
  res.writeHead(status, { 'Content-Type': 'application/json', 'Access-Control-Allow-Origin': '*' })
  res.end(JSON.stringify(data))
}

function txt(res, content) {
  res.writeHead(200, { 'Content-Type': 'text/markdown; charset=utf-8', 'Access-Control-Allow-Origin': '*' })
  res.end(content)
}

function notFound(res) {
  res.writeHead(404, { 'Content-Type': 'application/json', 'Access-Control-Allow-Origin': '*' })
  res.end(JSON.stringify({ error: 'not found' }))
}

createServer((req, res) => {
  if (req.method === 'OPTIONS') {
    res.writeHead(204, { 'Access-Control-Allow-Origin': '*', 'Access-Control-Allow-Methods': 'GET,POST', 'Access-Control-Allow-Headers': 'Content-Type' })
    return res.end()
  }

  const url = req.url.split('?')[0]

  if (url === '/api/health')
    return json(res, { version: '1.4.2', projectRoot: '/Users/dev/zcode', currentFeature: 'task-dashboard' })

  if (url === '/api/features')
    return json(res, { features: FEATURES.map(f => ({ slug: f.slug, title: f.title, stats: computeStats(f.tasks), lastUpdated: f.lastUpdated })) })

  if (url === '/api/records')
    return json(res, { records: [...RECORDS].sort((a, b) => b.completedAt.localeCompare(a.completedAt)) })

  if (url === '/api/lessons')
    return json(res, { lessons: LESSONS })

  if (url === '/api/tasks/claim' && req.method === 'POST')
    return json(res, { taskId: 'T-005', key: 'auth-login/T-005', title: 'Refresh Token 轮换', file: 'features/auth-login/tasks/T-005.yaml' })

  let m
  if ((m = url.match(/^\/api\/features\/([^/]+)$/)) && req.method === 'GET') {
    const f = FEATURES.find(x => x.slug === m[1]); if (!f) return notFound(res)
    return json(res, { slug: f.slug, title: f.title, stats: computeStats(f.tasks), lastUpdated: f.lastUpdated, tasks: f.tasks })
  }
  if ((m = url.match(/^\/api\/features\/([^/]+)\/tasks$/))) {
    const f = FEATURES.find(x => x.slug === m[1]); if (!f) return notFound(res)
    return json(res, { tasks: f.tasks })
  }
  if ((m = url.match(/^\/api\/features\/([^/]+)\/tasks\/([^/]+)\/status$/)) && req.method === 'POST')
    return json(res, { ok: true })

  if ((m = url.match(/^\/api\/features\/([^/]+)\/tasks\/([^/]+)$/))) {
    const f = FEATURES.find(x => x.slug === m[1])
    const t = f?.tasks.find(x => x.id === decodeURIComponent(m[2])); if (!t) return notFound(res)
    const detail = { ...t }
    if (typeof detail.record === 'string') detail.record = TASK_RECORDS[detail.record]
    return json(res, detail)
  }
  if ((m = url.match(/^\/api\/features\/([^/]+)\/prd$/)))
    return PRD[m[1]] ? txt(res, PRD[m[1]]) : notFound(res)

  if ((m = url.match(/^\/api\/features\/([^/]+)\/design$/)))
    return DESIGN[m[1]] ? txt(res, DESIGN[m[1]]) : notFound(res)

  if ((m = url.match(/^\/api\/features\/([^/]+)\/records$/)))
    return json(res, { records: RECORDS.filter(r => r.featureSlug === m[1]) })

  if ((m = url.match(/^\/api\/lessons\/([^/]+)$/)))
    return LESSON_CONTENTS[m[1]] ? txt(res, LESSON_CONTENTS[m[1]]) : notFound(res)

  notFound(res)
}).listen(7300, () => {
  console.log('Mock API running → http://localhost:7300')
  console.log('Features: auth-login, task-dashboard, api-gateway (39 tasks, 11 records, 6 lessons)')
})
