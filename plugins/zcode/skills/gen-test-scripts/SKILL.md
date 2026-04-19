---
name: gen-test-scripts
description: Generate executable TypeScript e2e test scripts from test cases. Uses agent-browser for UI, fetch for API, child_process for CLI. Zero external deps (node:test + node:assert).
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

| 产物 | 缺失时提示 |
|------|-----------|
| `testing/test-cases.md` | 先执行 `/gen-test-cases` |

```bash
ls docs/features/<slug>/testing/test-cases.md
```

## When to Use

**Trigger:**
- User asks to "generate test scripts" or "create e2e scripts"
- User provides `/gen-test-scripts` command
- After `/gen-test-cases` has produced `testing/test-cases.md`

**Skip:**
- `testing/test-cases.md` doesn't exist (run `/gen-test-cases` first)

## Workflow

```
1. Read test cases → 2. Generate helpers.ts → 3. Generate spec files → 4. Generate config files
```

### Step 1: Read Test Cases

Read `docs/features/<slug>/testing/test-cases.md`.

Parse each test case:
- Extract TC ID, title, type, pre-conditions, steps, expected result, priority
- Group by type: UI / API / CLI
- Count each group

### Step 2: Generate helpers.ts

Read the template at `plugins/zcode/skills/gen-test-scripts/templates/helpers.ts`.

Customize based on the feature:
- **`baseUrl`**: If testing a prototype, use `http://localhost:3456` (served from `ui/prototype/`). If testing a real app, ask the user or detect from the project.
- **`screenshotsDir`**: Set to `../results/screenshots/`

Write to: `testing/scripts/helpers.ts`

### Step 3: Generate Spec Files

For each type group, generate a spec file from the corresponding template.

<EXTREMELY-IMPORTANT>
**UI 测试必须使用 `agent-browser` CLI**，禁止使用 Playwright、Puppeteer 或其他浏览器自动化框架。
agent-browser 是本项目的 UI 测试标准工具，通过 `ab()` helper 调用。

**API 测试使用 Node.js 内置 `fetch`**，禁止引入 axios、supertest 等外部 HTTP 库。

**CLI 测试使用 `child_process.execSync`**，通过 `runCli()` helper 调用。

**所有测试使用 `node:test` + `node:assert`**，禁止引入 jest、mocha、vitest 等测试框架。
</EXTREMELY-IMPORTANT>

**UI tests (`testing/scripts/ui.spec.ts`)**:
- Read template: `plugins/zcode/skills/gen-test-scripts/templates/ui.spec.ts`
- For each UI test case, generate a `test()` block:
  - Translate steps into `agent-browser` commands via `ab()`
  - Common patterns:
    - Navigate: `ab('open <url>')` + `ab('wait --load networkidle')`
    - Find element: `abJson('snapshot -i')` then parse snapshot for expected text
    - Interact: `ab('click @eN')`, `ab('fill @eN "value"')`
    - Verify: assert on snapshot content, element visibility, page URL
    - Screenshot: `screenshot('TC-NNN')` for visual evidence
  - Each test includes a comment: `// Traceability: TC-NNN → {Source}`

**API tests (`testing/scripts/api.spec.ts`)**:
- Read template: `plugins/zcode/skills/gen-test-scripts/templates/api.spec.ts`
- For each API test case, generate a `test()` block:
  - Use `curl()` helper for HTTP requests
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

### Step 4: Generate Config Files

Write `testing/scripts/package.json` from template:
- Minimal: only `tsx` as devDependency for running TypeScript
- Scripts: `"test:ui"`, `"test:api"`, `"test:cli"`, `"test:all"`

Write `testing/scripts/tsconfig.json` from template:
- Target ES2022, module NodeNext, strict mode

## Output Files

All files go to `docs/features/<slug>/testing/scripts/`:

```
scripts/
  helpers.ts       # Shared utilities (ab, curl, runCli, screenshot)
  ui.spec.ts       # UI tests via agent-browser
  api.spec.ts      # API tests via fetch
  cli.spec.ts      # CLI tests via child_process
  package.json     # Minimal config (tsx only)
  tsconfig.json    # ES2022 + NodeNext
```

## Agent-Browser Command Reference

When generating UI test code, use these `agent-browser` commands via the `ab()` helper:

| Action | Command | Example |
|--------|---------|---------|
| Navigate | `ab('open <url>')` | `ab('open http://localhost:3456/login.html')` |
| Wait | `ab('wait --load networkidle')` | After navigation |
| Snapshot | `ab('snapshot -i --json')` | Get interactive elements |
| Click | `ab('click @eN')` | Click element by ref |
| Fill input | `ab('fill @eN "value"')` | Type into input |
| Press key | `ab('press Enter')` | Keyboard action |
| Find text | `ab('find text "Login" --json')` | Verify text exists |
| Screenshot | `screenshot('TC-001')` | Capture visual evidence |
| Get text | `ab('get text @eN')` | Read element content |
| Is visible | `ab('is visible @eN --json')` | Check visibility |
| Close | `ab('close')` | Close browser |

## Related Skills

| Skill | Usage |
|-------|-------|
| `/gen-test-cases` | Generate test cases from PRD |
| `/run-e2e-tests` | Execute test scripts and report results |
