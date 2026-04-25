---
name: gen-test-scripts
description: Generate executable TypeScript e2e test scripts from test cases. Uses Playwright for UI (semantic locators from sitemap.json), fetch for API, child_process for CLI. Zero external test frameworks (node:test + node:assert).
---

# Gen Test Scripts

从测试用例生成可执行的 TypeScript e2e 测试脚本。

**核心原则**：每个生成的 `test()` 必须可独立运行、可重复执行、有明确断言。测试用例是输入，脚本是产出。

<HARD-GATE>
本 skill 只生成测试脚本（testing/scripts/），不执行测试。
测试执行由 `/run-e2e-tests` skill 负责。
</HARD-GATE>

## Prerequisites

检查上一阶段产物，缺失则中止并提示用户：

| 产物                                     | 缺失时提示                                    |
| ---------------------------------------- | --------------------------------------------- |
| `testing/test-cases.md`                | 先执行 `/gen-test-cases`  |
| `docs/sitemap/sitemap.json`（仅 UI 测试需要） | 先执行 `/gen-sitemap` |
| `tests/e2e/config.yaml` | 先执行 `/gen-sitemap` 或手动创建 |

**注意**：本 skill 既可以手动调用，也可以作为 `/breakdown-tasks` 追加的标准任务 T-test-2 被 agent 调用。

```bash
ls docs/features/<slug>/testing/test-cases.md
ls docs/sitemap/sitemap.json  # 仅 UI 测试需要（项目级文件，非 feature 级）
```

**注意**：`<slug>` 为当前 feature 名称，通过 `task feature` 命令获取。`docs/sitemap/sitemap.json` 是项目级文件（一个应用只有一份），不是按 feature 隔离的。

### Sitemap

UI 测试的 locator 原料来自 `docs/sitemap/sitemap.json`。由 `/gen-sitemap` 命令自动生成和维护。

**位置**: `docs/sitemap/sitemap.json`
**完整示例**: `plugins/zcode/references/shared/sitemap.json`

**关键 locator 字段**:

| 字段                          | 用途                                            |
| ----------------------------- | ----------------------------------------------- |
| `layout.elements[].role` + `name` | → 布局级共享元素（导航栏等），可用于所有被 `layout.wraps` 包裹的页面 |
| `elements[].role` + `name`    | → `page.getByRole(role, { name })` (优先级 0)  |
| `elements[].level`            | → `getByRole('heading', { level })` 精确定位    |
| `elements[].label`            | → `page.getByLabel(label)` (优先级 1)           |
| `elements[].placeholder`      | → `page.getByPlaceholder(placeholder)` (优先级 2) |
| `states[].trigger`            | → 点击触发元素进入动态状态                      |
| `states[].elements`           | → 动态状态内的元素 locator                      |

测试用例中通过 `Route` 字段匹配 sitemap 页面。`Element` 字段（引用元素 ID 如 E-NNN / L-NNN）为可选——存在时精确映射，不存在时使用页面全部元素。

## When to Use

**Trigger:**

- User asks to "generate test scripts" or "create e2e scripts"
- User provides `/gen-test-scripts` command
- After `/gen-test-cases` has produced `testing/test-cases.md`

**Skip:**

- `testing/test-cases.md` doesn't exist (run `/gen-test-cases` first)

## Workflow

```
1. Read test cases → 2. Resolve sitemap → 3. Map locators → 4. Generate spec files → 5. Generate config files
```

### Step 1: Read Test Cases

Read `docs/features/<slug>/testing/test-cases.md`.

Parse each test case:

- Extract TC ID, title, type, route, feature, pre-conditions, steps, expected result, priority
- Group by type: UI / API / CLI
- Count each group

#### Auth Classification

对每个测试用例，按认证需求分为三类：

| 分类 | 检测规则 | 生成策略 |
|------|---------|---------|
| **login-test** | `Target` 匹配 `ui/login`, `ui/auth`, `ui/signin`, `api/auth`, `api/login`, `api/token` | 无共享认证。UI 用独立 `loginPage`，API 用原始 `curl()` |
| **auth-required-test** | `Pre-conditions` 含 "authenticated", "logged in", "已登录", "已认证"，或 `Target` 暗示受保护资源 | 使用共享认证 |
| **public-test** | 无认证相关前置条件，目标是公开资源 | 无需认证 |

统计每类数量，决定是否启用共享认证（当存在 `auth-required-test` 时启用）。

<HARD-RULE>
Login test 和 authenticated test 不得混在同一个 `describe` 块中。
</HARD-RULE>

### Step 2: Resolve Sitemap

**仅当存在 UI 类型测试用例时执行。**

读取 `docs/sitemap/sitemap.json`。对测试用例中涉及的每个路由：

1. 按 `Route` 字段匹配 `sitemap.json` 中的 `pages[].route`（动态路由如 `/tasks/:id` 匹配具体路径 `/tasks/123`）
2. 收集匹配页面的全部元素数据（base elements + 相关 states）
3. **若测试用例含 `Element` 字段**（引用 sitemap 元素 ID 如 E-NNN / L-NNN）：仅使用指定的元素，精确映射测试步骤到 sitemap 元素
4. **若测试用例无 `Element` 字段**：使用该页面的所有可用元素构建 locator 映射，agent 根据测试步骤描述自行匹配最合适的元素
5. **若 sitemap 包含 `layout` 字段**：对每个被 `layout.wraps` 包裹的路由，将 `layout.elements` 合并为可用元素（布局元素在该页面的 locator 映射中也可用）

若测试用例引用的 route 在 sitemap 中不存在，报告缺失并建议重新运行 `/gen-sitemap`。

### Step 3: Map Locators

基于 Step 2 收集的 sitemap 元素数据，按优先级为每个元素生成 Playwright locator：

| 优先级      | 条件                                             | 生成代码                                               |
| ----------- | ------------------------------------------------ | ------------------------------------------------------ |
| 0（最稳定） | `role` ∈ {button,link,heading,...} + `name` 非空 | `page.getByRole('button', { name: 'Submit' })`         |
| 0+          | heading + `level` 非空                           | `page.getByRole('heading', { name: 'Dashboard', level: 1 })` |
| 1           | `label` 非空                                     | `page.getByLabel('Email address')`                     |
| 2           | `placeholder` 非空，label 为空                   | `page.getByPlaceholder('Search...')`                   |
| 3           | 静态文本节点                                     | `page.getByText('No results found')`                   |
| 4           | `data-testid` 可见                               | `page.getByTestId('user-avatar')`                      |
| 5（降级）   | 以上均不满足                                     | `page.locator('.btn') // UNSTABLE: no semantic anchor` |

对动态状态内的测试步骤：先用触发元素的 locator 点击，再对状态内元素映射 locator。

构建内存映射表供 Step 4 使用。

### Step 4: Generate Spec Files

For each type group, generate a spec file from the corresponding template.

**模板使用方式**：模板中包含 `CONDITIONAL` 注释块，标记了按 auth 分类需要启用/禁用的代码段。根据 Step 1 的 auth 分类结果，**取消注释**对应的 CONDITIONAL 块，移除不匹配的块，然后填充测试数据。不要从头重写模板结构。

<EXTREMELY-IMPORTANT>
**UI 测试使用 Playwright Locator API**（仅作浏览器驱动库），禁止使用 agent-browser 在生成的 spec 文件中。

**API 测试使用 Node.js 内置 `fetch`**，禁止引入 axios、supertest 等外部 HTTP 库。

**CLI 测试使用 `child_process.execSync`**，通过 `runCli()` helper 调用。

**所有测试使用 `node:test` + `node:assert`**，禁止引入 jest、mocha、vitest 等测试框架。

**`@eN` ref 不得出现在任何生成的 spec 文件中。**
</EXTREMELY-IMPORTANT>

**UI tests (`testing/scripts/ui.spec.ts`)**:

- Read template: `plugins/zcode/skills/gen-test-scripts/templates/playwright-ui.spec.ts`
- **Auth setup**（顶层 `before` 钩子）：
  - 若存在 `auth-required-test`：在 `setupBrowser()` 之后调用 `await loginViaUI(page)`
  - 若仅有 `public-test` 和/或 `login-test`：不调用 `loginViaUI()`
- **Login test cases**：放入嵌套 `describe('Login', ...)` 块
  - 通过 `page.context().newPage()` 创建独立 `loginPage`（无认证 cookie）
  - 块结束后 `loginPage.close()`
  - 不使用共享认证
- **Authenticated test cases** + **public-test cases**：放在顶层，使用共享 `page`（已登录状态）
- Use Step 3 的 locator 映射表替换模板中的示例 locator
- 每个 `test()` 使用 `await page.xxx.waitFor()` 替代 `assert.ok(snapshotContains())`
- 保留 `screenshot(page, 'TC-NNN')` 调用
- 降级 locator 附注释 `// UNSTABLE: no semantic anchor`
- Each test includes a comment: `// Traceability: TC-NNN → {Source}`

**API tests (`testing/scripts/api.spec.ts`)**:

- Read template: `plugins/zcode/skills/gen-test-scripts/templates/api.spec.ts`
- **Auth setup**（顶层 `before` 钩子）：
  - 若存在 `auth-required-test`：调用 `const token = await getApiToken(apiBaseUrl)` 并创建 `authCurl = createAuthCurl(apiBaseUrl, token)`
  - 若仅有 `public-test` 和/或 `login-test`：不设置 auth
- **Login/auth API test cases**：使用原始 `curl()`（无 Bearer header）
- **Authenticated test cases**：使用 `authCurl(method, path)`（自动注入 Bearer token）
- **Public test cases**：使用 `curl()`（无认证）
- For each API test case, generate a `test()` block:
  - Assert on status code, response body, headers
  - Each test includes a traceability comment

**CLI tests (`testing/scripts/cli.spec.ts`)**:

- Read template: `plugins/zcode/skills/gen-test-scripts/templates/cli.spec.ts`
- For each CLI test case, generate a `test()` block:
  - Use `runCli()` helper to execute commands
  - Assert on stdout, stderr, exit code
  - Each test includes a traceability comment

<HARD-RULE>
**Traceability**: 每个 `test()` 必须包含追溯注释 `// Traceability: TC-NNN → {PRD Source}`，确保测试脚本与 PRD 的可追溯性链路完整。

**Skip empty groups**: If no test cases of a given type exist, skip generating that spec file.
</HARD-RULE>

### Step 5: Generate Config Files

Write `testing/scripts/package.json` from template:

- `tsx` + `playwright` as devDependencies
- Scripts: `"test:ui"`, `"test:api"`, `"test:cli"`, `"test:all"`

Write `testing/scripts/tsconfig.json` from template:

- Target ES2022, module NodeNext, strict mode

## Output Files

All files go to `docs/features/<slug>/testing/scripts/`:

```
scripts/
  helpers.ts       # Shared utilities (browser lifecycle, curl, auth helpers, runCli)
  ui.spec.ts       # UI tests via Playwright (shared auth via loginViaUI)
  api.spec.ts      # API tests via fetch (shared auth via authCurl)
  cli.spec.ts      # CLI tests via child_process
  package.json     # tsx + playwright
  tsconfig.json    # ES2022 + NodeNext
```

## Locator Priority Reference

When generating UI test code, use Playwright locators following this priority:

| Priority | Method                    | Example                                              |
| -------- | ------------------------- | ---------------------------------------------------- |
| 0        | `page.getByRole()`        | `page.getByRole('button', { name: 'Submit' })`       |
| 0+       | `page.getByRole()`        | `page.getByRole('heading', { name: 'Dashboard', level: 1 })` |
| 1        | `page.getByLabel()`       | `page.getByLabel('Email address')`                   |
| 2        | `page.getByPlaceholder()` | `page.getByPlaceholder('Search...')`                 |
| 3        | `page.getByText()`        | `page.getByText('No results found')`                 |
| 4        | `page.getByTestId()`      | `page.getByTestId('user-avatar')`                    |
| 5        | `page.locator()`          | `page.locator('.btn') // UNSTABLE`                   |

## Related Skills

| Skill             | Usage                                   |
| ----------------- | --------------------------------------- |
| `/gen-test-cases` | Generate test cases from PRD            |
| `/gen-sitemap`    | Generate and maintain sitemap.json      |
| `/run-e2e-tests`  | Execute test scripts and report results |
