---
name: gen-sitemap
description: Auto-generate and maintain sitemap.json for a web app. Uses agent-browser to explore routes, capture accessibility tree, and discover dynamic states. Preserves element IDs across runs.
argument-hints:
  - name: base-url
    description: 待探索的应用基础 URL（如 http://localhost:5173）
    required: true
---

# /gen-sitemap

自动生成并维护 web 应用的 `docs/sitemap/sitemap.json`。

**核心原则**：sitemap 是 web 应用的完整结构化地图，作为 Playwright locator 生成的唯一原料。元素 ID（E-NNN）是稳定标识，跨生成保持不变。

## Prerequisites

- Web 应用已启动并可访问（通过 `base-url` 参数指定）
- agent-browser 已安装（`npx agent-browser install`）

## Schema

完整示例见 `plugins/zcode/references/shared/sitemap.json`。

**关键字段**:

| 字段 | 说明 |
|------|------|
| `pages[].elements[].role` | 可访问角色（button, heading 等） |
| `pages[].elements[].name` | 可访问名称 |
| `pages[].elements[].level` | heading 层级（仅 heading 角色） |
| `pages[].elements[].label` | 关联 label 文本（仅 textbox 等表单元素） |
| `pages[].elements[].placeholder` | 占位文本（仅 textbox） |
| `pages[].states[]` | 动态状态（modal、tab panel、dropdown 等） |
| `pages[].states[].trigger` | 触发元素 ID（如 `"E-002"`） |
| `pages[].states[].elements` | 状态内的元素（同样带 E-NNN ID） |

## Workflow

```
1. Load existing sitemap → 2. Discover routes → 3. Explore pages → 4. Merge & write
```

### Step 1: Load Existing Sitemap

读取 `docs/sitemap/sitemap.json`（如果存在）。

构建 ID 索引：以 `route + role + name` 为 key，映射到已有 ID。

若无现有 sitemap，从 `E-001` 开始编号。

### Step 2: Discover Routes

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

**动态路由处理**：带参数的路由（如 `/tasks/123`）记录为模板形式 `/tasks/:id`（去除 URL 中的数字和 UUID 段）。

### Step 3: Explore Pages

对每个路由逐一用 agent-browser 探索：

```
ab('open <baseUrl><route>')
ab('wait --load networkidle')
snapshot = abJson('snapshot -i')
```

#### 基础元素提取

1. 获取页面 title
2. 从 snapshot 提取元素，过滤条件：
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

```
ab('close')
```

### Step 4: Merge & Write

对每个元素（含 states 内元素），用 `route + role + name` 三元组匹配现有 sitemap：

- **匹配成功** → 保留原有 ID
- **无匹配** → 分配新 ID（当前最大 ID + 1）
- **现有 ID 无匹配** → 元素已移除，从 sitemap 中删除

写入 `docs/sitemap/sitemap.json`。

**报告变更**：

```
Sitemap updated: docs/sitemap/sitemap.json
  3 pages, 18 elements (base) + 8 elements (states)
  +4 new elements (E-015..E-018)
  -2 removed elements (previously on /settings)
  22 unchanged
```

## Element ID 分配规则

- 格式：`E-NNN`（三位数字），全局唯一（base 和 states 共享 ID 空间）
- 首次生成：从 `E-001` 顺序分配
- 增量更新：新元素从当前最大 ID + 1 开始
- ID 永不重复使用（删除的 ID 不回收）
