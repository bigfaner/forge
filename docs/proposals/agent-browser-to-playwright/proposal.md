---
title: agent-browser 探索 → Playwright 稳定脚本
status: proposed
date: 2026-04-24
---

# 提案：agent-browser 探索 → Playwright 稳定脚本

## 背景

当前 `gen-test-scripts` skill 使用 `agent-browser` CLI 生成 UI 测试脚本。核心问题是 `@eN` ref 不稳定——ref 是 snapshot 时的临时 ID，DOM 变化后即失效，导致测试脆弱。

### 问题量化

在实际使用中，`@eN` ref 的不稳定性表现如下：

- **重排即失效**：任何 DOM 插入/删除/重排操作都会改变元素的 `@eN` 编号。例如页面 `<nav>` 前插入一个 `<div>` 广告横幅，会导致其后所有 `@eN` 编号 +1，全部已生成的测试断言失效。
- **维护成本**：每次 UI 变更后平均需要手动修复 3-5 个测试文件中的 `@eN` ref，每次修复约耗时 15-30 分钟定位和替换。来源：2026-04-15 在 `feat-task-cli-scope-extension` 分支合并后，手动排查并修复了 `ui.spec.ts` 中 4 处 `@eN` ref 失效，耗时约 25 分钟（见 commit `857c9cc` 前后 diff）。
- **CI 误报**：2026-04-01 至 2026-04-20 期间的 CI 运行中，至少 3 次 UI 测试失败经排查确认为 `@eN` ref 失效所致，而非功能回归——分别是 04-03（`@e5` 偏移，因 `<header>` 新增面包屑组件）、04-12（`@e2` 失效，因列表渲染顺序变化）、04-18（`@e8` 失效，因 modal 内新增提示文案）。每次误报排查耗时 10-15 分钟。
- **不可复现性**：`@eN` ref 依赖运行时 DOM 快照顺序，不同环境（本地 vs CI）或不同数据集可能产生不同的编号，导致"本地通过 CI 失败"的不可预测行为。

**不解决此问题的现状量化**：当前项目包含 4 个 UI 页面，每次 DOM 变更平均影响 3-5 个测试文件、耗时 15-30 分钟修复。按项目路线图，Q3 前预计扩展至 12 个 UI 页面。若保持现状，单次 DOM 变更的测试维护成本将从 15-30 分钟增长至约 45-90 分钟（3x 页面数 × 相同的 `@eN` 修复模式），UI 测试的维护投入将超过功能开发本身。

本提案将 agent-browser 的角色限定为**探索工具**，利用其 Accessibility Tree 输出作为语义桥梁，生成基于 Playwright 稳定 Locator API 的测试脚本。

## 备选方案

### 方案 A：维持现状（Do Nothing）

继续使用 `@eN` ref 生成测试脚本，遇到失效后手动修复。

| 优点 | 缺点 |
|------|------|
| 零迁移成本，无需改动现有代码 | 维护成本随页面数量线性增长 |
| agent-browser 作为单一工具链，无额外依赖 | `@eN` ref 失效导致的 CI 误报持续发生（04 月已记录 3 次） |
| 短期内可继续产出测试 | 页面扩展至 12 个时，单次 DOM 变更维护成本将达 45-90 分钟 |

**结论**：排除。当前模式不可持续，问题会随规模扩大而恶化。

### 方案 B：Cypress 替代 agent-browser

使用 Cypress Test Runner 作为 UI 测试框架，通过 `cy.get()` + 自定义 selector 替代 `@eN` ref。

| 优点 | 缺点 |
|------|------|
| 成熟的生态和社区支持 | 引入完整测试框架，与现有 `node:test` 运行器冲突 |
| Time-travel 调试体验好 | Cypress 不支持多 tab、多窗口场景 |
| 内置截图和视频录制 | 需要 Node.js 以外的额外安装（Cypress binary ~200MB） |
| | 无法与 agent-browser 的 Accessibility Tree 输出对接 |
| | 需要重写全部测试基础设施，迁移范围远超本提案 |

**结论**：排除。Cypress 作为完整测试框架与现有 `node:test` 体系冲突，迁移成本过高，且无法利用 agent-browser 已有的 Accessibility Tree 能力。

### 方案 C：Playwright Codegen 录制

使用 `npx playwright codegen` 录制用户操作，自动生成测试脚本。

| 优点 | 缺点 |
|------|------|
| 零代码生成 locator | 录制产出的代码质量低，需要大量手动清理 |
| 官方工具，稳定可靠 | 每次页面变更需要重新录制，无法自动化集成到 SKILL 流程 |
| 生成的 locator 使用语义 API | 录制过程需要人工操作，无法批量化处理多个测试用例 |
| | 无法与 `test-cases.md` 驱动的工作流集成 |

**结论**：排除。Codegen 适合人工辅助编写，但不适合集成到自动化生成流程中。

### 方案 D（选定）：agent-browser 探索 + Playwright Locator API

保留 agent-browser 用于探索阶段的 Accessibility Tree 采集，通过 LLM 机械映射生成 Playwright 语义 locator。

| 优点 | 缺点 |
|------|------|
| 最小迁移范围：仅改动 UI 测试路径，API/CLI 不动 | 仍依赖 agent-browser（探索阶段） |
| 语义 locator（role + name）对 DOM 结构变化免疫 | 页面 a11y 质量差时降级率可能偏高 |
| 与现有 `node:test` 运行器完全兼容 | LLM 映射环节需要质量监控 |
| 利用已有工具能力，无需引入新框架 | |

**选择理由**：方案 D 在迁移成本、长期稳定性和与现有架构的兼容性之间取得最佳平衡。Playwright 仅作为浏览器驱动库（非测试框架），`node:test` 体系不受影响。语义 locator 的稳定性已在 Playwright 官方文档和社区实践中被广泛验证——`getByRole('button', { name: 'Login' })` 只要按钮的 role 和可访问名称不变，无论 DOM 结构如何重组都能稳定定位。

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

LLM 做这个翻译极其可靠——规则简单、无歧义。映射规则是确定性的一对一转换：`role` 直接映射到 `getByRole`，`label` 直接映射到 `getByLabel`，`placeholder` 直接映射到 `getByPlaceholder`。没有判断空间，没有歧义场景。成功标准第 2 条定义了可验证的阈值：降级 locator 占比 < 10%，即 ≥90% 的元素命中优先级 0-4 的语义映射。

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

## 开发者工作流

开发者通过 `/zcode:gen-test-scripts` 触发整个流程。以下是开发者视角的完整交互过程：

### Step 1: 读取测试用例（自动）

```
📖 Reading test cases from docs/features/task-cli/testing/test-cases.md
   Found 8 test cases: 3 UI, 3 API, 2 CLI
   UI routes detected: /, /settings, /tasks/:id
```

### Step 2: 探索 UI 页面（自动，开发者等待）

```
🔍 Exploring UI pages via agent-browser...
   [1/3] GET / ... captured base + 1 dynamic state (add-task-modal)
   [2/3] GET /settings ... captured base (no dynamic states)
   [3/3] GET /tasks/1 ... captured base + 2 dynamic states (edit-modal, delete-confirm)
✅ Exploration complete → docs/features/task-cli/testing/page-map.json
   14 semantic elements mapped, 2 degraded (14% degradation rate)
```

开发者可在此暂停，打开 `page-map.json` 审阅探索结果——确认动态状态是否遗漏、元素名称是否准确。

### Step 3: Locator 映射（自动）

```
🔄 Mapping elements to Playwright locators...
   Priority 0 (getByRole):   8 elements
   Priority 1 (getByLabel):  2 elements
   Priority 2 (getByPlaceholder): 1 element
   Priority 3 (getByText):   1 element
   Priority 5 (DEGRADED):    2 elements ← UNSTABLE, no semantic anchor
   ⚠ Degradation rate 14% (< 30% threshold, proceeding)
```

若降级率超过 30%，开发者看到阻断提示：
```
   🚫 Degradation rate 35% exceeds 30% threshold.
   建议：为标记 UNSTABLE 的元素添加 aria-label 或使用语义化 HTML。
   页面 page-map.json 已标记所有降级元素，修复后重新运行即可。
```

### Step 4: 生成测试文件（自动）

```
📝 Generating test scripts...
   → api.spec.ts (3 tests, unchanged path)
   → cli.spec.ts (2 tests, unchanged path)
   → ui.spec.ts  (3 tests, Playwright locators, zero @eN refs)

✅ All test scripts generated.
   Run with: npx tsx ui.spec.ts
```

### 开发者迭代

若生成的测试执行失败：
```
❌ TC-002 failed: locator.getByRole('button', { name: 'Submit' }) timed out
   → 检查 page-map.json 中该元素映射是否准确
   → 若页面结构已变更，重新运行 Step 2 探索
```

开发者修复方式：调整 `page-map.json` 中对应元素的映射，或为页面元素补充 `aria-label`，然后重新生成。

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

## 成功标准

| # | 标准 | 验证方式 |
|---|------|---------|
| 1 | 生成的 UI 测试文件中 `@eN` ref 出现次数 = 0 | `grep -c '@e' ui.spec.ts` 返回 0 |
| 2 | 降级 locator（优先级 5，使用 CSS 选择器）占比 < 10% | 统计 page-map.json 中 `unstable: true` 元素比例 |
| 3 | page-map.json 覆盖 test-cases.md 中 100% 的页面路由 | 逐一比对 test-cases.md 提及的 route 与 page-map.json 的 pages[].route |
| 4 | 生成的 UI 测试在目标应用上全部通过（exit code 0） | 执行 `npx tsx ui.spec.ts`，确认 0 failures |
| 5 | API/CLI 测试行为完全不变 | 对比迁移前后 `tsx api.spec.ts` 和 `tsx cli.spec.ts` 输出，结果一致 |
| 6 | 单次 UI 测试生成耗时 < 60 秒（不含 agent-browser 探索时间） | 计时 Step 3 + Step 4 执行时间 |

## 风险

| 风险 | 可能性 | 影响 | 缓解 |
|------|--------|------|------|
| 页面 a11y 质量差，降级率 > 30% | 中（依赖现有页面结构） | 高（脚本稳定性回退至接近现状） | 降级率警告机制 + 降级 locator 标记 `// UNSTABLE` 注释，便于后续定向修复 |
| 动态状态探索不完整（遗漏 Modal/Tab 内元素） | 中（复杂交互页面） | 中（部分元素未出现在 page-map.json 中） | 对每个触发元素做 click → snapshot → reset 循环；探索结果人工可审阅（page-map.json 是可读 JSON） |
| agent-browser + Playwright 双工具依赖增加 CI 安装成本 | 低（均为一次性安装） | 低（CI 时间增加 < 2 分钟） | agent-browser 仅探索阶段使用，可在 CI cache 中持久化 Playwright chromium binary |
| LLM 映射引入错误 locator（语义不匹配） | 低（规则为确定性映射） | 中（错误定位导致测试失败） | 映射规则表固定且无歧义；生成后可通过 `npx tsx ui.spec.ts` 即时验证 |
