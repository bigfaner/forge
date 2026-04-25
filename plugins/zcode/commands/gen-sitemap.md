---
name: gen-sitemap
description: Auto-generate and maintain sitemap.json for a web app. Uses agent-browser to explore routes, capture accessibility tree, and discover dynamic states. Preserves element IDs across runs.
argument-hints:
  - name: base-url
    description: 待探索的应用基础 URL（如 http://localhost:3456），config.yaml 存在时可省略
    required: false
---

# /gen-sitemap

自动生成并维护 web 应用的 `docs/sitemap/sitemap.json`。

**核心原则**：sitemap 是 web 应用的完整结构化地图，作为 Playwright locator 生成的唯一原料。元素 ID（E-NNN）是稳定标识，跨生成保持不变。

## Prerequisites

- Web 应用已启动并可访问
- agent-browser 已安装（可选工具，用于自动探索页面结构和动态状态）

**验证 agent-browser 安装**：

```bash
npx agent-browser --version
```

若失败，先安装：

```bash
npx agent-browser install
```

安装完成后再重新运行本命令。

## Config Resolution

在执行 workflow 之前，解析配置源：

1. 检查 `tests/e2e/config.yaml` 是否存在
2. **若不存在**：从模板 `plugins/zcode/references/shared/config.yaml` 复制到 `tests/e2e/config.yaml`，然后**中止并提示用户**：

```
已创建 tests/e2e/config.yaml（模板），请根据实际环境填写配置后重新运行。
  baseUrl: 当前应用的实际地址
  username/password: 测试账号凭据（若应用需要认证）
  loginLocators: 若登录页定位器与默认值不匹配，取消注释并自定义
```

3. **若已存在**：读取 `baseUrl`、`username`、`password` 等字段

**base-url 优先级**：命令行参数 > config.yaml 中的 `baseUrl` > 报错中止

**认证页面探索**：若 config.yaml 中 `username` 和 `password` 非空，在 Step 3 探索每个页面前先执行登录：

```
ab('open <baseUrl>/login')
ab('wait --load networkidle')
// 使用 config.yaml 中的 username/password 填充登录表单
ab('click <login_button>')
ab('wait --load networkidle')
```

这确保 agent-browser 能访问需要认证的页面，避免因未登录导致探索中断。

## Schema

完整示例见 `plugins/zcode/references/shared/sitemap.json`。

**关键字段**:

| 字段 | 说明 |
|------|------|
| `baseUrl` | 应用基础 URL（如 `http://localhost:3456`） |
| `updatedAt` | 最后更新时间（RFC3339 格式） |
| `layout.name` | 布局组件名（如 `AppLayout`） |
| `layout.wraps` | 共享此布局的路由列表 |
| `layout.elements[]` | 布局级共享元素（侧边栏、顶部导航等），ID 格式 `L-NNN` |
| `pages[].elements[].role` | 可访问角色（button, heading 等） |
| `pages[].elements[].name` | 可访问名称 |
| `pages[].elements[].level` | heading 层级（仅 heading 角色） |
| `pages[].elements[].label` | 关联 label 文本（仅 textbox 等表单元素） |
| `pages[].elements[].placeholder` | 占位文本（仅 textbox） |
| `pages[].states[]` | 动态状态（modal、tab panel、dropdown 等） |
| `pages[].states[].trigger` | 触发元素 ID（如 `"E-002"`） |
| `pages[].states[].elements` | 状态内的元素（同样带 E-NNN ID） |

**布局元素 vs 页面元素**：

- `layout.elements`：所有被布局包裹的页面共享的元素（导航栏、侧边栏、页脚），ID 格式 `L-NNN`
- `pages[].elements`：仅属于当前页面的特有元素，ID 格式 `E-NNN`
- 测试脚本生成时，布局元素可用在任何被包裹的页面中

## Workflow

```
1. Load existing sitemap → 2. Analyze layout → 3. Discover routes → 4. Explore pages → 5. Merge & write
```

### Step 1: Load Existing Sitemap

读取 `docs/sitemap/sitemap.json`（如果存在）。

构建 ID 索引：以 `route + role + name` 为 key，映射到已有 ID。

若无现有 sitemap，从 `E-001` 开始编号。

### Step 2: Analyze Layout

<EXTREMELY-IMPORTANT>
**此步骤必须在页面探索前执行。** 目标是识别共享布局，避免在每个页面重复提取布局元素。

**先读代码，再探索。** 不要先无脑探索再回头去重。
</EXTREMELY-IMPORTANT>

**2a. 读取路由定义**：查找应用的路由配置文件（如 `App.tsx`、`router.ts`、`routes.tsx`），识别布局嵌套：

```
<Route element={<AppLayout />}>      ← 共享布局
  <Route path="/" element={<Dashboard />} />
  <Route path="/settings" element={<Settings />} />
  ...
</Route>
<Route path="/login">                ← 无布局（独立页面）
```

记录：
- `layout.name`：布局组件名（如 `AppLayout`）
- `layout.wraps`：被布局包裹的路由列表
- 不被任何布局包裹的路由（如 `/login`）为独立页面

**2b. 探索布局元素**：对布局包裹的第一个页面做 snapshot，提取**仅属于布局层**的元素：

```
ab('open <baseUrl><first_wrapped_route>')
ab('wait --load networkidle')
snapshot = abJson('snapshot -i')
```

从 snapshot 中识别布局区域（通常为 `role=navigation`、`role=banner`、侧边栏容器等），提取其中的元素归入 `layout.elements`。

**若无法读取路由代码**（纯 HTML 项目或无源码访问）：跳过此步骤，所有元素按页面级处理。在 Step 4 探索前两个页面后，对比两次 snapshot 中 `role + name` 完全相同的元素作为布局候选项，归入 `layout.elements`，后续页面探索时过滤掉这些候选项。

### Step 3: Discover Routes

使用 agent-browser 导航到用户提供的 `base-url`，提取页面中所有链接：

```
ab('open <base-url>')
ab('wait --load networkidle')
links = abJson('snapshot -i')  // 提取所有 role=link 节点的 href
```

1. 过滤为同源路径（排除外部链接、`mailto:`、`javascript:`）
2. 去重得到路由列表
3. 对每个新路由递归提取链接（广度优先，最大深度 3）
4. 合并现有 sitemap 中手动添加的路由

**动态路由处理**：带参数的路由（如 `/tasks/123`）记录为模板形式 `/tasks/:id`。参数化规则：

| URL 段模式 | 替换为 | 示例 |
|-----------|--------|------|
| 纯数字 | `:id` | `/tasks/42` → `/tasks/:id` |
| UUID 格式 | `:uuid` | `/orders/550e8400-...` → `/orders/:uuid` |
| 32 位 hex | `:hash` | `/files/a1b2c3d4e5f6...` → `/files/:hash` |

去重时，模板相同的路由只保留一个条目。`layout.wraps` 中也使用模板形式。

### Step 4: Explore Pages

对每个路由逐一用 agent-browser 探索：

```
ab('open <baseUrl><route>')
ab('wait --load networkidle')
snapshot = abJson('snapshot -i')
```

#### 布局元素过滤

<HARD-RULE>
若 Step 2 识别了共享布局，此步骤**必须跳过布局元素**。
对比 snapshot 与 `layout.elements`，过滤掉 `role + name` 匹配的元素。
只有页面特有的内容区元素才归入 `pages[].elements`。
</HARD-RULE>

#### 基础元素提取

1. 获取页面 title
2. 从 snapshot 提取元素，过滤条件：
   - 排除已归入 `layout.elements` 的元素（按 `role + name` 匹配）
   - `role` ∈ {button, link, heading, textbox, checkbox, radio, combobox, tab, dialog, alert, navigation, search, form, menuitem, switch}
   - `name` 非空
3. 对每个元素记录完整属性：
   - 通用：`{ role, name }`
   - heading：额外记录 `level`
   - textbox/combobox：额外记录 `label`（关联 label 文本）和 `placeholder`

#### 动态状态探索

对基础元素中 role=button/tab/disclosure 且 name 非空的触发元素：

```
ab('click @eN')
ab('wait --load networkidle')
state_snapshot = abJson('snapshot -i')
// 提取新增元素（对比基础 snapshot）
ab('press Escape')  // 或 ab('click @close_btn') 重置
```

1. 比较状态 snapshot 与基础 snapshot，提取新增元素
2. 记录为 `states` 条目：`{ name, trigger: "<元素ID>", elements: [...] }`
3. `trigger` 引用触发元素的 E-NNN ID（如 `"E-002"`）
4. 状态内元素同样分配 E-NNN ID

> **注意**：`@eN` 是 agent-browser CLI 的元素引用语法，仅在 sitemap 生成过程中使用。生成的测试脚本（`*.spec.ts`）中禁止使用 `@eN`，必须使用 Playwright Locator API。

```
ab('close')
```

### Step 5: Merge & Write

对每个元素（含 layout 和 states 内元素），用 `route + role + name` 三元组匹配现有 sitemap：

- **匹配成功** → 保留原有 ID
- **无匹配** → 分配新 ID（当前最大 ID + 1）
- **现有 ID 无匹配** → 元素已移除，从 sitemap 中删除

写入 `docs/sitemap/sitemap.json`。

**报告变更**：

```
Sitemap updated: docs/sitemap/sitemap.json
  Layout: AppLayout (5 shared elements, L-001..L-005)
  3 pages, 12 elements (page-specific) + 8 elements (states)
  +4 new elements (E-015..E-018)
  -2 removed elements (previously on /settings)
  17 unchanged
```

## Element ID 分配规则

- **布局元素**：格式 `L-NNN`，全局唯一，独立编号空间
- **页面元素**：格式 `E-NNN`，全局唯一（base 和 states 共享 ID 空间）
- 首次生成：分别从 `L-001` 和 `E-001` 顺序分配
- 增量更新：新元素从当前最大 ID + 1 开始
- ID 永不重复使用（删除的 ID 不回收）
