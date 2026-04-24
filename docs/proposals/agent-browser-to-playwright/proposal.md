---
title: agent-browser 探索 → Playwright 稳定脚本
status: proposed
date: 2026-04-24
---

# 提案：agent-browser 探索 → Playwright 稳定脚本

## 背景

当前 `gen-test-scripts` skill 使用 `agent-browser` CLI 生成 UI 测试脚本。核心问题是 `@eN` ref 不稳定——ref 是 snapshot 时的临时 ID，DOM 变化后即失效，导致测试脆弱。

本提案将 agent-browser 的角色限定为**探索工具**，利用其 Accessibility Tree 输出作为语义桥梁，生成基于 Playwright 稳定 Locator API 的测试脚本。

## 核心洞察

agent-browser `snapshot` 输出的 Accessibility Tree 与 Playwright 现代 Locator API 共享同一语义基础：

```
agent-browser snapshot 节点:
  { role: "button", name: "Login" }
  { role: "textbox", label: "Email" }

↓ 机械映射 ↓

Playwright locator:
  page.getByRole('button', { name: 'Login' })
  page.getByLabel('Email')
```

LLM 做这个翻译极其可靠——规则简单、无歧义。

## 变更范围

```
修改：
  plugins/zcode/skills/gen-test-scripts/SKILL.md
  plugins/zcode/skills/gen-test-scripts/templates/helpers.ts
  plugins/zcode/skills/gen-test-scripts/templates/package.json
  plugins/zcode/skills/run-e2e-tests/SKILL.md

新增：
  plugins/zcode/skills/gen-test-scripts/templates/playwright-ui.spec.ts

不变：
  templates/api.spec.ts, cli.spec.ts, tsconfig.json
  node:test 作为测试运行器
  API/CLI 测试路径完全不动
```

## 新工作流

```
旧流程：
  test-cases.md → [生成 ab() 调用] → ui.spec.ts → 执行

新流程：
  test-cases.md
       ↓
  Step 2: 探索阶段（agent-browser）
       ↓ 每个页面/路由做 snapshot
  page-map.json（语义节点清单）
       ↓
  Step 3: Locator 映射（LLM 翻译）
       ↓ role+name → getByRole / getByLabel / getByText
  locator-map（内存中，直接用于生成）
       ↓
  Step 4: 生成 Playwright spec
       ↓
  ui.spec.ts（无 @eN ref，无 ab() 调用）
```

## 探索阶段设计（Step 2）

### 探索目标

从 `test-cases.md` 提取所有涉及的页面路由，对每个路由：
1. 捕获基础状态快照
2. 识别并触发动态状态（Modal、Tab、Dropdown）
3. 捕获每个状态的快照
4. 重置状态，继续下一个

### 探索执行步骤

```
对每个 route:
  ab('open <baseUrl><route>')
  ab('wait --load networkidle')
  snap_base = abJson('snapshot -i')

  识别触发元素（role=button/tab/disclosure 且 name 非空）:
    for each trigger_element in snap_base:
      ab('click @eN')
      ab('wait --load networkidle')
      snap_state = abJson('snapshot -i')
      记录 { trigger: {role, name}, snapshot: snap_state }
      ab('press Escape') 或 ab('click @close_btn')  ← 重置

ab('close')
```

### page-map.json 格式

```json
{
  "exploredAt": "2026-04-24T10:00:00Z",
  "baseUrl": "http://localhost:3456",
  "pages": [
    {
      "route": "/",
      "title": "Dashboard",
      "baseElements": [
        { "role": "heading",  "name": "Dashboard", "level": 1 },
        { "role": "button",   "name": "Add Task" },
        { "role": "textbox",  "label": "Search",  "placeholder": "Search tasks..." },
        { "role": "link",     "name": "Settings" }
      ],
      "states": [
        {
          "name": "add-task-modal",
          "trigger": { "role": "button", "name": "Add Task" },
          "elements": [
            { "role": "dialog",  "name": "Add Task" },
            { "role": "textbox", "label": "Task Title" },
            { "role": "button",  "name": "Save" },
            { "role": "button",  "name": "Cancel" }
          ]
        }
      ]
    }
  ]
}
```

写入：`docs/features/<slug>/testing/page-map.json`

## Locator 映射规则

LLM 读取 page-map.json，按优先级生成 Playwright locator：

| 优先级 | 条件 | 生成代码 |
|--------|------|---------|
| 0（最稳定） | `role` ∈ {button,link,heading,...} + `name` 非空 | `page.getByRole('button', { name: 'Submit' })` |
| 1 | `role=textbox` + `label` 非空 | `page.getByLabel('Email address')` |
| 2 | `role=textbox` + `placeholder` 非空，label 为空 | `page.getByPlaceholder('Search...')` |
| 3 | 静态文本节点 | `page.getByText('No results found')` |
| 4 | `data-testid` 可见 | `page.getByTestId('user-avatar')` |
| 5（降级） | 以上均不满足 | `page.locator('.btn') // UNSTABLE: no semantic anchor` |

**降级率 > 30% 时**，SKILL 输出警告：`⚠ 页面语义覆盖率不足，建议为关键元素添加 aria-label 或使用语义化 HTML`。

## 新模板文件

### templates/helpers.ts（更新）

```typescript
import { execSync } from 'node:child_process';
import { writeFileSync, mkdirSync, existsSync } from 'node:fs';
import { join, dirname } from 'node:path';
import { fileURLToPath } from 'node:url';
import { chromium, type Browser, type Page } from 'playwright';

const __dirname = dirname(fileURLToPath(import.meta.url));
const DEFAULT_TIMEOUT = parseInt(process.env.PLAYWRIGHT_TIMEOUT ?? '30000');
const SCREENSHOTS_DIR = join(__dirname, '..', 'results', 'screenshots');

export const baseUrl = process.env.E2E_BASE_URL ?? 'http://localhost:3456';

// ── Browser lifecycle ──────────────────────────────────────────────
let _browser: Browser | null = null;
let _page: Page | null = null;

export async function setupBrowser(): Promise<Page> {
  _browser = await chromium.launch();
  _page = await _browser.newPage();
  _page.setDefaultTimeout(DEFAULT_TIMEOUT);
  return _page;
}

export async function teardownBrowser(): Promise<void> {
  await _browser?.close();
  _browser = null;
  _page = null;
}

export function getPage(): Page {
  if (!_page) throw new Error('Browser not initialized. Call setupBrowser() first.');
  return _page;
}

// ── Evidence ───────────────────────────────────────────────────────
export async function screenshot(page: Page, tcId: string): Promise<string> {
  if (!existsSync(SCREENSHOTS_DIR)) mkdirSync(SCREENSHOTS_DIR, { recursive: true });
  const path = join(SCREENSHOTS_DIR, `${tcId}.png`);
  await page.screenshot({ path, fullPage: true });
  return path;
}

// ── HTTP (unchanged) ───────────────────────────────────────────────
export interface CurlResponse { status: number; headers: Record<string, string>; body: string; }

export async function curl(method: string, url: string, opts?: {
  body?: string; headers?: Record<string, string>; timeout?: number;
}): Promise<CurlResponse> {
  const controller = new AbortController();
  const timer = setTimeout(() => controller.abort(), opts?.timeout ?? 10000);
  try {
    const res = await fetch(url, {
      method,
      headers: { 'Content-Type': 'application/json', ...opts?.headers },
      body: opts?.body,
      signal: controller.signal,
    });
    const headers: Record<string, string> = {};
    res.headers.forEach((v, k) => { headers[k] = v; });
    return { status: res.status, headers, body: await res.text() };
  } finally { clearTimeout(timer); }
}

// ── CLI (unchanged) ────────────────────────────────────────────────
export interface CliResult { stdout: string; stderr: string; exitCode: number; }

export function runCli(cmd: string, cwd?: string): CliResult {
  try {
    const stdout = execSync(cmd, { encoding: 'utf-8', timeout: DEFAULT_TIMEOUT, cwd: cwd ?? process.cwd() });
    return { stdout, stderr: '', exitCode: 0 };
  } catch (e: any) {
    return { stdout: e.stdout ?? '', stderr: e.stderr ?? '', exitCode: e.status ?? 1 };
  }
}
```

### templates/playwright-ui.spec.ts（新增）

```typescript
import { describe, test, before, after } from 'node:test';
import assert from 'node:assert/strict';
import type { Page } from 'playwright';
import { setupBrowser, teardownBrowser, screenshot, baseUrl } from './helpers.js';

describe('UI E2E Tests', () => {
  let page: Page;

  before(async () => {
    page = await setupBrowser();
    await page.goto(baseUrl);
    await page.waitForLoadState('networkidle');
  });

  after(async () => {
    await teardownBrowser();
  });

  // Traceability: TC-001 → Story 1 / AC-1
  test('TC-001: Page renders with expected heading', async () => {
    await page.goto(`${baseUrl}/index.html`);
    await page.waitForLoadState('networkidle');
    await page.getByRole('heading', { name: 'Dashboard' }).waitFor({ state: 'visible' });
    assert.ok(true, 'Heading visible');
    await screenshot(page, 'TC-001');
  });
});
```

### templates/package.json（更新）

```json
{
  "name": "e2e-tests",
  "private": true,
  "type": "module",
  "scripts": {
    "test:ui":  "tsx ui.spec.ts",
    "test:api": "tsx api.spec.ts",
    "test:cli": "tsx cli.spec.ts",
    "test:all": "tsx cli.spec.ts && tsx api.spec.ts && tsx ui.spec.ts"
  },
  "devDependencies": {
    "tsx": "^4.19.0",
    "playwright": "^1.44.0"
  }
}
```

## gen-test-scripts SKILL.md 工作流变更

在"Read Test Cases"之后插入两个新步骤：

### Step 2: Explore UI Pages（新增）

**仅当存在 UI 类型测试用例时执行。**

从 test-cases.md 提取所有涉及的页面路由，去重后逐一探索：

1. `ab('open <baseUrl><route>')` + `ab('wait --load networkidle')`
2. `abJson('snapshot -i')` — 捕获基础状态
3. 识别触发元素（role=button/tab/disclosure，name 非空）
4. 对每个触发元素：click → snapshot → reset（Escape 或 click Cancel）
5. `ab('close')`

将探索结果写入 `testing/page-map.json`。

**降级处理**：若某元素无 role 或 name，记录为 `{ "role": "generic", "unstable": true }`。
降级率 > 30% 时输出警告。

### Step 3: Map Locators（新增）

读取 page-map.json，按优先级规则为每个元素生成 Playwright locator，构建内存映射表供 Step 4 使用。

### Step 4: Generate Spec Files（UI 部分更新）

**UI tests** — 读取模板 `templates/playwright-ui.spec.ts`：
- 使用 Step 3 的 locator 映射表替换模板中的示例 locator
- 每个 `test()` 使用 `await page.xxx.waitFor()` 替代 `assert.ok(snapshotContains())`
- 保留 `screenshot(page, 'TC-NNN')` 调用
- 降级 locator 附注释 `// UNSTABLE: no semantic anchor`

**移除** agent-browser command reference 章节（UI 测试不再需要）。

## run-e2e-tests SKILL.md 变更

### Setup 新增

```bash
# Install Playwright browser (if not installed)
npx playwright install chromium
```

### Teardown 变更

移除 `agent-browser close --all`（Playwright 在 `after()` 钩子中自动关闭）。

## 实施顺序

```
1. 更新 templates/helpers.ts
2. 新增 templates/playwright-ui.spec.ts
3. 更新 templates/package.json
4. 更新 gen-test-scripts/SKILL.md
5. 更新 run-e2e-tests/SKILL.md
```

## 关键约束

- `node:test` 保持不变，Playwright 仅作浏览器驱动库（非测试框架）
- agent-browser 仅在探索阶段使用，不出现在最终 spec 文件中
- `@eN` ref 不得出现在任何生成的 spec 文件中
- API/CLI 测试路径零改动

## 风险

| 风险 | 影响 | 缓解 |
|------|------|------|
| 页面 a11y 质量差，降级率高 | 脚本稳定性回退 | 降级率警告 + 注释标记 |
| 动态状态探索不完整 | 遗漏 Modal/Tab 内元素 | 多轮 snapshot 策略 |
| 两工具同时依赖 | CI 安装成本增加 | agent-browser 仅探索阶段，可选 |
