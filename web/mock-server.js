/**
 * ZCode Dashboard — Mock API Server
 * Run: node mock-server.js
 * Serves mock data on http://localhost:7300/api/...
 */

const http = require('http')

// ─── Mock Data ──────────────────────────────────────────────────────────────

const FEATURES = [
  {
    slug: 'auth-login',
    title: '用户登录与认证',
    stats: { total: 12, pending: 3, in_progress: 1, completed: 7, blocked: 1, skipped: 0 },
    lastUpdated: '2026-04-07T14:23:00Z',
  },
  {
    slug: 'task-dashboard',
    title: 'Web Dashboard 前端',
    stats: { total: 18, pending: 8, in_progress: 2, completed: 6, blocked: 0, skipped: 2 },
    lastUpdated: '2026-04-08T09:45:00Z',
  },
  {
    slug: 'api-gateway',
    title: 'API Gateway 集成',
    stats: { total: 9, pending: 0, in_progress: 0, completed: 9, blocked: 0, skipped: 0 },
    lastUpdated: '2026-04-05T16:00:00Z',
  },
]

const TASKS = {
  'auth-login': [
    {
      id: '1.1', title: '定义用户数据模型', description: '设计 User struct，包含 ID、email、passwordHash、createdAt 字段',
      phase: 1, priority: 'P0', status: 'completed', estimatedTime: '2h',
      dependencies: [], files: ['pkg/model/user.go', 'pkg/model/user_test.go'],
      record: { summary: '实现了 User 数据模型，包含完整的 JSON 序列化支持', filesCreated: ['pkg/model/user.go'], filesModified: ['pkg/model/user_test.go'], decisions: ['使用 bcrypt 存储密码哈希', 'ID 使用 UUID v4'], testResults: 'PASS 8/8', coverage: '94.2%', commitHash: 'a3f9c12', completedAt: '2026-04-06T10:30:00Z' }
    },
    {
      id: '1.2', title: '实现登录接口', description: '实现 POST /api/auth/login，验证凭据并返回 JWT token',
      phase: 1, priority: 'P0', status: 'completed', estimatedTime: '3h',
      dependencies: ['1.1'], files: ['internal/handler/auth.go', 'internal/handler/auth_test.go'],
      record: { summary: '实现了 JWT 登录接口，token 有效期 24h', filesCreated: ['internal/handler/auth.go'], filesModified: ['internal/handler/auth_test.go'], decisions: ['JWT 有效期设为 24h', '使用 RS256 算法签名'], testResults: 'PASS 12/12', coverage: '87.3%', commitHash: 'b7e2d45', completedAt: '2026-04-06T14:15:00Z' }
    },
    {
      id: '1.3', title: '实现 Token 刷新逻辑', description: '实现 POST /api/auth/refresh，使用 refresh token 换取新的 access token',
      phase: 1, priority: 'P1', status: 'completed', estimatedTime: '2h',
      dependencies: ['1.2'], files: ['internal/handler/auth.go'],
      record: { summary: 'Refresh token 机制实现完成，有效期 7 天', filesCreated: [], filesModified: ['internal/handler/auth.go'], decisions: ['Refresh token 存入 Redis，支持撤销'], testResults: 'PASS 6/6', coverage: '91.0%', commitHash: 'c9f3a78', completedAt: '2026-04-06T16:00:00Z' }
    },
    {
      id: '2.1', title: '实现登出接口', description: '实现 POST /api/auth/logout，将 token 加入黑名单',
      phase: 2, priority: 'P1', status: 'completed', estimatedTime: '1h',
      dependencies: ['1.3'], files: ['internal/handler/auth.go'],
      record: { summary: '登出接口实现，token 加入 Redis 黑名单', filesCreated: [], filesModified: ['internal/handler/auth.go'], decisions: ['使用 Redis TTL 自动清理过期黑名单'], testResults: 'PASS 4/4', coverage: '88.5%', commitHash: 'd2b1e56', completedAt: '2026-04-07T09:00:00Z' }
    },
    {
      id: '2.2', title: '中间件：JWT 鉴权', description: '实现 JWT 验证中间件，保护需要认证的路由',
      phase: 2, priority: 'P0', status: 'completed', estimatedTime: '2h',
      dependencies: ['1.2'], files: ['internal/middleware/auth.go', 'internal/middleware/auth_test.go'],
      record: { summary: 'JWT 鉴权中间件实现，支持 Bearer token 和 Cookie', filesCreated: ['internal/middleware/auth.go'], filesModified: ['internal/middleware/auth_test.go'], decisions: ['支持 Authorization header 和 httpOnly Cookie 两种方式'], testResults: 'PASS 10/10', coverage: '92.1%', commitHash: 'e4c7f89', completedAt: '2026-04-07T11:30:00Z' }
    },
    {
      id: '2.3', title: '集成测试：完整认证流程', description: '编写端到端测试，覆盖登录→使用→刷新→登出完整流程',
      phase: 2, priority: 'P1', status: 'completed', estimatedTime: '3h',
      dependencies: ['2.1', '2.2'], files: ['test/auth_integration_test.go'],
      record: { summary: '集成测试覆盖完整认证流程，发现并修复了 token 并发刷新的竞态问题', filesCreated: ['test/auth_integration_test.go'], filesModified: ['internal/handler/auth.go'], decisions: ['增加 token 刷新的分布式锁'], testResults: 'PASS 15/15', coverage: '89.7%', commitHash: 'f5d8a01', completedAt: '2026-04-07T15:00:00Z' }
    },
    {
      id: '3.1', title: '密码重置功能', description: '实现忘记密码→发送邮件→重置密码完整流程',
      phase: 3, priority: 'P1', status: 'in_progress', estimatedTime: '4h',
      dependencies: ['2.x'], files: ['internal/handler/password.go', 'internal/service/email.go'],
      record: undefined
    },
    {
      id: '3.2', title: 'OAuth2 第三方登录', description: '接入 GitHub OAuth2，支持第三方账户绑定',
      phase: 3, priority: 'P2', status: 'blocked', estimatedTime: '6h',
      dependencies: ['3.1'], files: ['internal/handler/oauth.go', 'pkg/oauth/github.go'],
      record: undefined
    },
    {
      id: '3.3', title: '登录频率限制', description: '实现登录失败次数限制，防止暴力破解',
      phase: 3, priority: 'P1', status: 'pending', estimatedTime: '2h',
      dependencies: ['2.2'], files: ['internal/middleware/ratelimit.go'],
      record: undefined
    },
    {
      id: '3.4', title: '用户注册接口', description: '实现 POST /api/auth/register，含邮箱验证',
      phase: 3, priority: 'P0', status: 'pending', estimatedTime: '3h',
      dependencies: ['1.1'], files: ['internal/handler/register.go'],
      record: undefined
    },
    {
      id: '3.5', title: '二步验证 (2FA)', description: '实现 TOTP 二步验证，兼容 Google Authenticator',
      phase: 3, priority: 'P2', status: 'pending', estimatedTime: '5h',
      dependencies: ['3.4'], files: ['internal/handler/totp.go', 'pkg/totp/totp.go'],
      record: undefined
    },
    {
      id: '3.6', title: '用户会话管理', description: '实现多设备会话列表查看与主动踢出',
      phase: 3, priority: 'P2', status: 'pending', estimatedTime: '3h',
      dependencies: ['2.1'], files: ['internal/handler/session.go'],
      record: undefined
    },
  ],
  'task-dashboard': [
    {
      id: '1.1', title: 'task serve 命令', description: '实现 cobra task serve --port 7300 命令',
      phase: 1, priority: 'P0', status: 'completed', estimatedTime: '1h',
      dependencies: [], files: ['internal/cmd/serve.go'],
      record: { summary: 'task serve 命令实现完成，支持 --port 参数', filesCreated: ['internal/cmd/serve.go'], filesModified: ['internal/cmd/root.go'], decisions: ['默认端口 7300'], testResults: 'PASS 3/3', coverage: '85.0%', commitHash: 'g6e9b23', completedAt: '2026-04-07T10:00:00Z' }
    },
    {
      id: '1.2', title: 'HTTP Server + embed', description: '实现 net/http 路由，embed 前端静态文件',
      phase: 1, priority: 'P0', status: 'completed', estimatedTime: '2h',
      dependencies: ['1.1'], files: ['internal/server/server.go', 'web/web.go'],
      record: { summary: 'HTTP Server 实现，前端 dist 通过 go:embed 内嵌', filesCreated: ['internal/server/server.go', 'web/web.go'], filesModified: [], decisions: ['SPA fallback 返回 index.html'], testResults: 'PASS 5/5', coverage: '80.0%', commitHash: 'h7f0c34', completedAt: '2026-04-07T13:00:00Z' }
    },
    {
      id: '1.3', title: 'Feature 列表 API', description: '实现 GET /api/features，扫描 docs/features/*',
      phase: 1, priority: 'P0', status: 'completed', estimatedTime: '2h',
      dependencies: ['1.2'], files: ['internal/server/handlers/features.go'],
      record: { summary: 'Feature 列表 API 实现，含任务统计', filesCreated: ['internal/server/handlers/features.go'], filesModified: [], decisions: ['lastUpdated 取 index.json 文件修改时间'], testResults: 'PASS 7/7', coverage: '88.0%', commitHash: 'i8a1d45', completedAt: '2026-04-07T16:00:00Z' }
    },
    {
      id: '2.1', title: 'Vite + React 脚手架', description: '初始化前端项目，配置 TypeScript + Tailwind',
      phase: 2, priority: 'P0', status: 'completed', estimatedTime: '1h',
      dependencies: ['1.x'], files: ['web/package.json', 'web/vite.config.ts', 'web/tailwind.config.ts'],
      record: { summary: 'Vite React TypeScript 项目初始化完成', filesCreated: ['web/package.json', 'web/vite.config.ts'], filesModified: [], decisions: ['使用 Tailwind v3 + shadcn/ui'], testResults: 'N/A', coverage: 'N/A', commitHash: 'j9b2e56', completedAt: '2026-04-08T08:00:00Z' }
    },
    {
      id: '2.2', title: 'AppLayout + 路由', description: '实现侧边栏导航 + React Router v6',
      phase: 2, priority: 'P0', status: 'completed', estimatedTime: '2h',
      dependencies: ['2.1'], files: ['web/src/layouts/AppLayout.tsx', 'web/src/App.tsx'],
      record: { summary: 'AppLayout 侧边栏实现，7 个路由配置完成', filesCreated: ['web/src/layouts/AppLayout.tsx'], filesModified: ['web/src/App.tsx'], decisions: ['深色模式默认开启，支持切换'], testResults: 'N/A', coverage: 'N/A', commitHash: 'k0c3f67', completedAt: '2026-04-08T10:00:00Z' }
    },
    {
      id: '2.3', title: 'Dashboard 首页', description: '实现全局统计 + Feature 卡片网格',
      phase: 2, priority: 'P0', status: 'completed', estimatedTime: '3h',
      dependencies: ['2.2'], files: ['web/src/pages/Dashboard.tsx', 'web/src/components/FeatureCard.tsx'],
      record: { summary: 'Dashboard 页面实现，含统计卡片和 Feature 列表', filesCreated: ['web/src/pages/Dashboard.tsx'], filesModified: [], decisions: ['30s 轮询刷新数据'], testResults: 'N/A', coverage: 'N/A', commitHash: 'l1d4g78', completedAt: '2026-04-08T14:00:00Z' }
    },
    {
      id: '3.1', title: 'Feature 详情页（Tasks Tab）', description: '实现 Kanban / List / DAG 三视图',
      phase: 3, priority: 'P0', status: 'in_progress', estimatedTime: '4h',
      dependencies: ['2.x'], files: ['web/src/pages/FeatureDetail.tsx', 'web/src/components/task-board/KanbanView.tsx'],
      record: undefined
    },
    {
      id: '3.2', title: 'DAG 依赖图', description: '使用 ReactFlow 实现层次布局任务依赖图',
      phase: 3, priority: 'P1', status: 'in_progress', estimatedTime: '3h',
      dependencies: ['3.1'], files: ['web/src/components/task-board/DagView.tsx'],
      record: undefined
    },
    {
      id: '3.3', title: 'Claim Task 按钮', description: '实现一键认领任务，调用 POST /api/tasks/claim',
      phase: 3, priority: 'P0', status: 'pending', estimatedTime: '1h',
      dependencies: ['3.1'], files: ['web/src/components/ClaimButton.tsx'],
      record: undefined
    },
    {
      id: '3.4', title: '任务详情页', description: '展示任务信息、依赖、状态变更、执行记录',
      phase: 3, priority: 'P1', status: 'pending', estimatedTime: '3h',
      dependencies: ['3.1'], files: ['web/src/pages/TaskDetail.tsx'],
      record: undefined
    },
    {
      id: '3.5', title: 'PRD / Design Tab', description: 'Markdown 渲染 prd.md 和 design.md',
      phase: 3, priority: 'P1', status: 'pending', estimatedTime: '1h',
      dependencies: ['3.1'], files: ['web/src/pages/PrdTab.tsx', 'web/src/pages/DesignTab.tsx'],
      record: undefined
    },
    {
      id: '4.1', title: '执行记录时间线', description: '全局执行记录时间线页面',
      phase: 4, priority: 'P1', status: 'pending', estimatedTime: '2h',
      dependencies: ['3.x'], files: ['web/src/pages/Records.tsx'],
      record: undefined
    },
    {
      id: '4.2', title: '知识库页面', description: 'Lessons 列表 + 分类筛选 + 搜索',
      phase: 4, priority: 'P2', status: 'skipped', estimatedTime: '2h',
      dependencies: ['3.x'], files: ['web/src/pages/Lessons.tsx'],
      record: undefined
    },
    {
      id: '4.3', title: '设置页面', description: 'Health info + 服务器配置',
      phase: 4, priority: 'P2', status: 'skipped', estimatedTime: '1h',
      dependencies: [], files: ['web/src/pages/Settings.tsx'],
      record: undefined
    },
    {
      id: '4.4', title: 'Go embed 集成', description: '将前端 dist 打包进 task-cli 二进制',
      phase: 4, priority: 'P0', status: 'pending', estimatedTime: '1h',
      dependencies: ['4.1', '4.2', '4.3'], files: ['task-cli/web/web.go', 'Makefile'],
      record: undefined
    },
    {
      id: '4.5', title: 'E2E 验收测试', description: '端到端验证 Dashboard 完整功能',
      phase: 4, priority: 'P1', status: 'pending', estimatedTime: '3h',
      dependencies: ['4.4'], files: ['test/dashboard_e2e_test.go'],
      record: undefined
    },
    {
      id: '4.6', title: 'Makefile 构建脚本', description: '统一 make web + make build 构建流程',
      phase: 4, priority: 'P1', status: 'pending', estimatedTime: '1h',
      dependencies: ['4.4'], files: ['Makefile'],
      record: undefined
    },
  ],
  'api-gateway': [
    {
      id: '1.1', title: '路由注册模块', description: '实现动态路由注册，支持中间件链',
      phase: 1, priority: 'P0', status: 'completed', estimatedTime: '3h',
      dependencies: [], files: ['pkg/router/router.go'],
      record: { summary: '路由模块完成，支持 RESTful 路径参数', filesCreated: ['pkg/router/router.go'], filesModified: [], decisions: ['使用 trie 树匹配路由'], testResults: 'PASS 18/18', coverage: '93.5%', commitHash: 'm2e5h89', completedAt: '2026-04-04T14:00:00Z' }
    },
    {
      id: '1.2', title: '请求限流中间件', description: '令牌桶算法实现 API 限流',
      phase: 1, priority: 'P0', status: 'completed', estimatedTime: '2h',
      dependencies: [], files: ['pkg/middleware/ratelimit.go'],
      record: { summary: '令牌桶限流实现，支持按 IP 和 API Key 分别限制', filesCreated: ['pkg/middleware/ratelimit.go'], filesModified: [], decisions: ['默认 100 req/s per IP'], testResults: 'PASS 9/9', coverage: '90.2%', commitHash: 'n3f6i90', completedAt: '2026-04-04T16:30:00Z' }
    },
  ],
}

const LESSONS = [
  { name: 'debug-go-embed-path', category: 'debug', title: 'Go embed 路径不能包含 .. 的问题', excerpt: '使用 //go:embed 时，路径不能跨越模块边界（不能包含 ..），需要将 embed 包放在 dist 目录的同级或父目录。' },
  { name: 'arch-jwt-refresh-race', category: 'arch', title: 'JWT Refresh Token 并发刷新竞态问题', excerpt: '多个请求同时触发 token 刷新时，会导致旧 token 被多次使用。解决方案：使用分布式锁保证同一时刻只有一个刷新操作。' },
  { name: 'pattern-handler-closure', category: 'pattern', title: 'Go HTTP Handler 用闭包注入依赖', excerpt: '避免使用全局变量，使用 func NewHandler(deps) http.HandlerFunc 的闭包模式注入依赖，便于测试。' },
  { name: 'tool-tanstack-query-invalidate', category: 'tool', title: 'TanStack Query mutation 后的缓存失效策略', excerpt: 'mutation 成功后调用 queryClient.invalidateQueries() 时，要精确指定 queryKey 的层级，避免过度刷新导致性能问题。' },
  { name: 'gotcha-vite-proxy-trailing-slash', category: 'gotcha', title: 'Vite proxy 配置对 trailing slash 敏感', excerpt: 'Vite proxy 的 target 末尾不要加 /，否则会导致双斜杠路径问题（/api//features）。' },
  { name: 'debug-tailwind-dark-mode', category: 'debug', title: 'Tailwind dark: 变体不生效的排查步骤', excerpt: '需确认 tailwind.config 中 darkMode: "class" 已设置，且 document.documentElement 上有 dark class。' },
]

const RECORDS = [
  { featureSlug: 'auth-login', taskId: '2.3', taskTitle: '集成测试：完整认证流程', coverage: '89.7%', filesChanged: 2, commitHash: 'f5d8a01', completedAt: '2026-04-07T15:00:00Z' },
  { featureSlug: 'auth-login', taskId: '2.2', taskTitle: '中间件：JWT 鉴权', coverage: '92.1%', filesChanged: 2, commitHash: 'e4c7f89', completedAt: '2026-04-07T11:30:00Z' },
  { featureSlug: 'auth-login', taskId: '2.1', taskTitle: '实现登出接口', coverage: '88.5%', filesChanged: 1, commitHash: 'd2b1e56', completedAt: '2026-04-07T09:00:00Z' },
  { featureSlug: 'task-dashboard', taskId: '2.3', taskTitle: 'Dashboard 首页', coverage: 'N/A', filesChanged: 2, commitHash: 'l1d4g78', completedAt: '2026-04-08T14:00:00Z' },
  { featureSlug: 'task-dashboard', taskId: '2.2', taskTitle: 'AppLayout + 路由', coverage: 'N/A', filesChanged: 2, commitHash: 'k0c3f67', completedAt: '2026-04-08T10:00:00Z' },
  { featureSlug: 'task-dashboard', taskId: '2.1', taskTitle: 'Vite + React 脚手架', coverage: 'N/A', filesChanged: 3, commitHash: 'j9b2e56', completedAt: '2026-04-08T08:00:00Z' },
  { featureSlug: 'auth-login', taskId: '1.3', taskTitle: '实现 Token 刷新逻辑', coverage: '91.0%', filesChanged: 1, commitHash: 'c9f3a78', completedAt: '2026-04-06T16:00:00Z' },
  { featureSlug: 'auth-login', taskId: '1.2', taskTitle: '实现登录接口', coverage: '87.3%', filesChanged: 2, commitHash: 'b7e2d45', completedAt: '2026-04-06T14:15:00Z' },
  { featureSlug: 'auth-login', taskId: '1.1', taskTitle: '定义用户数据模型', coverage: '94.2%', filesChanged: 2, commitHash: 'a3f9c12', completedAt: '2026-04-06T10:30:00Z' },
  { featureSlug: 'api-gateway', taskId: '1.2', taskTitle: '请求限流中间件', coverage: '90.2%', filesChanged: 1, commitHash: 'n3f6i90', completedAt: '2026-04-04T16:30:00Z' },
  { featureSlug: 'api-gateway', taskId: '1.1', taskTitle: '路由注册模块', coverage: '93.5%', filesChanged: 1, commitHash: 'm2e5h89', completedAt: '2026-04-04T14:00:00Z' },
]

const PRD = `# 用户登录与认证 — PRD

## 背景
系统需要安全的用户认证机制，支持邮箱密码登录、第三方 OAuth 登录，以及 JWT token 管理。

## 目标
- 实现安全的用户认证系统
- 支持 access token + refresh token 双 token 机制
- 防止暴力破解攻击
- 支持多设备会话管理

## 用户故事
- 用户可以使用邮箱和密码登录
- 登录后获得 JWT access token（24h 有效期）
- access token 过期后可使用 refresh token 续期
- 用户可以主动登出，使 token 立即失效
- 管理员可以强制踢出指定用户的所有会话

## 验收标准
- [ ] 登录接口响应时间 < 200ms
- [ ] 密码错误 5 次后锁定 30 分钟
- [ ] Refresh token 支持轮换（rotation）
- [ ] 所有认证接口有完整的审计日志
`

const DESIGN = `# 技术设计文档 — 用户登录与认证

## 架构概览

\`\`\`
Client → API Gateway → Auth Handler → Auth Service → Database
                                   → JWT Service → Redis
\`\`\`

## 核心接口

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/auth/login | 登录，返回 access + refresh token |
| POST | /api/auth/refresh | 使用 refresh token 换新 token |
| POST | /api/auth/logout | 登出，撤销 token |
| POST | /api/auth/register | 注册新用户 |

## 数据模型

\`\`\`go
type User struct {
    ID           string    \`json:"id"\`
    Email        string    \`json:"email"\`
    PasswordHash string    \`json:"-"\`
    CreatedAt    time.Time \`json:"createdAt"\`
}

type TokenPair struct {
    AccessToken  string \`json:"accessToken"\`
    RefreshToken string \`json:"refreshToken"\`
    ExpiresIn    int    \`json:"expiresIn"\`
}
\`\`\`

## 安全考虑
- ��码使用 bcrypt（cost=12）存储
- JWT 使用 RS256 非对称加密
- Refresh token 存入 Redis，支持主动撤销
- 登录失败次数限制（5次/30分钟）
`

// ─── Router ──────────────────────────────────────────────────────────────────

function json(res, data, status = 200) {
  res.writeHead(status, { 'Content-Type': 'application/json', 'Access-Control-Allow-Origin': '*' })
  res.end(JSON.stringify(data))
}

function text(res, data) {
  res.writeHead(200, { 'Content-Type': 'text/plain; charset=utf-8', 'Access-Control-Allow-Origin': '*' })
  res.end(data)
}

const server = http.createServer((req, res) => {
  const url = req.url || '/'
  const parts = url.replace('/api/', '').split('?')[0].split('/')

  // CORS preflight
  if (req.method === 'OPTIONS') {
    res.writeHead(204, { 'Access-Control-Allow-Origin': '*', 'Access-Control-Allow-Methods': 'GET,POST', 'Access-Control-Allow-Headers': 'Content-Type' })
    return res.end()
  }

  // GET /api/health
  if (url === '/api/health') {
    return json(res, { version: '1.2.0', projectRoot: '/Users/nasuki/Downloads/zcode-main', currentFeature: 'task-dashboard' })
  }

  // GET /api/features
  if (url === '/api/features') {
    return json(res, { features: FEATURES })
  }

  // GET /api/records (global)
  if (url === '/api/records') {
    return json(res, { records: RECORDS })
  }

  // GET /api/lessons
  if (url === '/api/lessons') {
    return json(res, { lessons: LESSONS })
  }

  // GET /api/lessons/:name
  if (parts[0] === 'lessons' && parts[1]) {
    const l = LESSONS.find(x => x.name === parts[1])
    if (!l) return json(res, { error: 'not found' }, 404)
    return text(res, `# ${l.title}\n\n${l.excerpt}\n\n## 详情\n\n这里是该经验的完整内容，包含代码示例和详细说明。\n\n\`\`\`go\n// 示例代码\nfunc example() error {\n    return nil\n}\n\`\`\`\n`)
  }

  // GET /api/features/:slug
  if (parts[0] === 'features' && parts[1] && !parts[2]) {
    const f = FEATURES.find(x => x.slug === parts[1])
    if (!f) return json(res, { error: 'not found' }, 404)
    return json(res, { ...f, tasks: TASKS[parts[1]] || [] })
  }

  // GET /api/features/:slug/prd
  if (parts[0] === 'features' && parts[2] === 'prd') {
    return text(res, PRD)
  }

  // GET /api/features/:slug/design
  if (parts[0] === 'features' && parts[2] === 'design') {
    return text(res, DESIGN)
  }

  // GET /api/features/:slug/tasks
  if (parts[0] === 'features' && parts[2] === 'tasks' && !parts[3]) {
    const tasks = TASKS[parts[1]] || []
    return json(res, { tasks })
  }

  // GET /api/features/:slug/tasks/:id
  if (parts[0] === 'features' && parts[2] === 'tasks' && parts[3]) {
    const tasks = TASKS[parts[1]] || []
    const t = tasks.find(x => x.id === decodeURIComponent(parts[3]))
    if (!t) return json(res, { error: 'not found' }, 404)
    return json(res, t)
  }

  // GET /api/features/:slug/records
  if (parts[0] === 'features' && parts[2] === 'records') {
    const slug = parts[1]
    return json(res, { records: RECORDS.filter(r => r.featureSlug === slug) })
  }

  // POST /api/tasks/claim
  if (req.method === 'POST' && url === '/api/tasks/claim') {
    return json(res, { taskId: '3.3', key: 'task-dashboard', title: 'Claim Task 按钮', file: 'docs/features/task-dashboard/tasks/3.3-claim-task.md' })
  }

  // POST /api/features/:slug/tasks/:id/status
  if (req.method === 'POST' && parts[0] === 'features' && parts[2] === 'tasks' && parts[4] === 'status') {
    return json(res, { success: true })
  }

  json(res, { error: 'not found' }, 404)
})

const PORT = 7300
server.listen(PORT, () => {
  console.log(`\n  Mock API Server running at http://localhost:${PORT}`)
  console.log('  Start frontend: cd web && npm run dev\n')
})
