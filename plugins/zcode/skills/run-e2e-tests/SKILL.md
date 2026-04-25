---
name: run-e2e-tests
description: Execute e2e test scripts and generate a results report. Runs UI tests via Playwright, API tests via fetch, CLI tests via child_process. Produces evidence-backed pass/fail report.
---

# Run E2E Tests

执行 e2e 测试脚本并生成测试报告。

**核心原则**：只执行 `/gen-test-scripts` 生成的脚本，忠实收集结果，不做任何修改。

<HARD-GATE>
本 skill 只执行已有的测试脚本并报告结果。禁止：
- 修改测试脚本内容
- 跳过失败的测试（必须如实报告）
- 在执行过程中"修复"测试使之通过
</HARD-GATE>

## Prerequisites

检查上一阶段产物，缺失则中止并提示用户：

| 产物 | 缺失时提示 |
|------|-----------|
| `testing/scripts/` 目录 | 先执行 `/gen-test-scripts` |
| `testing/scripts/helpers.ts` | 先执行 `/gen-test-scripts` |
| 至少一个 `.spec.ts` 文件 | 先执行 `/gen-test-scripts` |
| `tests/e2e/config.yaml` | 先执行 `/gen-sitemap` 或手动创建 |

```bash
ls docs/features/<slug>/testing/scripts/
```

## When to Use

**Trigger:**
- User asks to "run e2e tests" or "run tests"
- User provides `/run-e2e-tests` command
- After `/gen-test-scripts` has generated test scripts

**Skip:**
- `testing/scripts/` doesn't exist (run `/gen-test-scripts` first)

## Workflow

```
1. Setup → 2. Run specs → 3. Collect results → 4. Generate report → 5. Teardown
```

### Step 1: Setup Environment

**Install dependencies** (if `node_modules` doesn't exist):

```bash
cd docs/features/<slug>/testing/scripts && npm install
```

**Install Playwright browser** (if not installed):

```bash
npx playwright install chromium
```

Verify installation:

```bash
npx playwright --version
```

**Start servers** (if needed):

Detect what needs to be started based on the project:

| Artifact | Server Command | When |
|----------|---------------|------|
| HTML prototype | `npx serve ../ui/prototype -p 3456 -s &` | `ui/prototype/` exists |
| Real application | User-specified or detected | App needs running server |

Start servers in background. Record PIDs for cleanup.

### Step 2: Run Test Specs

Run each spec file that exists:

```bash
cd docs/features/<slug>/testing/scripts
```

<EXTREMELY-IMPORTANT>
**执行顺序**: CLI → API → UI（由快到慢，环境依赖递增）

**Skip files that don't exist** — not all projects have all three types.

**每个 spec 的输出必须完整保存**（`tee` 到 results/ 目录），用于后续结果收集。
</EXTREMELY-IMPORTANT>

**CLI tests** (run first):

```bash
npx tsx cli.spec.ts 2>&1 | tee ../results/cli-output.txt
```

**API tests**:

```bash
npx tsx api.spec.ts 2>&1 | tee ../results/api-output.txt
```

**UI tests** (run last — requires Playwright browser):

```bash
npx tsx ui.spec.ts 2>&1 | tee ../results/ui-output.txt
```

### Step 3: Collect Results

Parse the Node.js test runner output (TAP-like format) from each spec:

- **Pass**: Lines containing `✓` or `ok N`
- **Fail**: Lines containing `✗` or `not ok N` followed by error details
- **Skip**: Lines containing `# skip`

For each test, capture:
- TC ID (from test name)
- Status: PASS / FAIL / SKIP
- Duration (from output)
- Error message (if failed)
- Screenshot path (if UI test)

### Step 4: Generate Report

Read the template at `plugins/zcode/skills/run-e2e-tests/templates/e2e-report.md`.

Fill in:
- Summary statistics (total/pass/fail/skip per type)
- Per-test-case results with evidence
- Failed test details with error messages
- Screenshot paths for UI tests

Write to: `docs/features/<slug>/testing/results/latest.md`

**For failed UI tests**: Use MCP image analysis tools to examine screenshots and add diagnostic notes:

```
Skill(skill="mcp__zai-mcp-server__analyze_image", image_source="<screenshot_path>", prompt="What is shown on this page? Why might the test have failed?")
```

### Step 5: Teardown

<HARD-RULE>
**Teardown 是强制步骤**，即使测试失败也必须执行：

1. `kill <server_pid>` — 终止所有在 Step 1 启动的服务器
2. 清理临时文件

Playwright 浏览器实例由 spec 文件中的 `after()` 钩子自动关闭，无需手动清理。
</HARD-RULE>

## Output

After completion, report to the user:

```
E2E Test Results: X/Y passed (Z failed)

Failed tests:
- TC-NNN: {failure reason}
- TC-NNN: {failure reason}

Report: docs/features/<slug>/testing/results/latest.md
```

If all tests pass:

```
E2E Test Results: X/X passed
Report: docs/features/<slug>/testing/results/latest.md
```

## Error Handling

| Situation | Action |
|-----------|--------|
| Playwright browser not installed | Run `npx playwright install chromium`, retry |
| Server won't start | Report error, skip tests that need it |
| Test timeout | Kill test process, mark as FAIL with timeout reason |
| `node_modules` missing | Run `npm install`, retry |
| Spec file doesn't compile | Report TypeScript error, skip that spec |

## Related Skills

| Skill | Usage |
|-------|-------|
| `/gen-test-cases` | Generate test cases from PRD |
| `/gen-test-scripts` | Generate executable test scripts |
