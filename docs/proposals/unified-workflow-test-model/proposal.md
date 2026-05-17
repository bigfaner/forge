---
created: 2026-05-15
updated: 2026-05-17
author: "faner"
status: Draft
---

# Proposal: Unified Workflow Test Model

## Problem

Forge 的 e2e 测试管道以"单动作"为粒度生成 test case，每个 TC 验证一条命令/一个端点/一个页面交互。但真实用户场景是跨多个操作的**工作流**（如 `forge feature wf-test-1 → task claim → task submit → verify`），当前管道无法生成覆盖完整用户旅程的 e2e 测试。

### Evidence

- forge-cli 的 e2e 测试（20+ 测试文件，126+ 测试用例）全部是单命令粒度（`exec.Command("forge", <single-command>)`），无跨命令工作流测试；6 条 Integration TC（TC-067 至 TC-071、TC-082）虽然描述了多步工作流，但生成的脚本仍缺少步骤间的声明式状态传递
- gen-test-cases 从 PRD AC 提取验证点时，多数 AC 映射为单条 TC；Integration TC 证明了 gen-test-cases 已能识别多步操作并生成多步 Steps 描述，但步骤间缺少 `{stepN.field}` 声明式引用机制
- 6 个 language 的模板（go, javascript, python, java, rust, mobile）以单动作为主要设计，无步骤间声明式状态传递机制（Go 的 generate.md 虽有 TUI golden file 多步模式，但非声明式状态传递）
- 行业标准（Microsoft Engineering Playbook, IBM）定义 e2e testing 为验证"complete user workflows from start to end"

### Reference Projects

两个真实项目暴露了问题的不同侧面：

**agent-forensic**（Go TUI 应用，go language）：Bubble Tea TUI 的工作流测试使用 `sendKey → viewContains` 模式验证键盘导航流程（如 session → call tree → diagnosis 面板），但当前 gen-test-scripts 无法生成此类测试。

**pm-work-tracker**（Go+React 混合 Web 应用，javascript language）：Playwright 测试中存在大量"API 数据准备 → UI 交互验证"模式（如 TC-026 先通过 API 创建 milestone+MI，再通过 UI 取消 milestone 并验证 MI 解绑），但 gen-test-cases 无法识别此类跨步依赖。

### Risk Quantification

- forge-cli 正在快速迭代，跨命令的回归风险在增长（如 task claim → submit 的状态链路）
- **手工状态传递成本已可量化**：当前 6 条 Integration TC（TC-067 至 TC-071、TC-082）每条需要约 15 分钟手工编写步骤间状态传递代码（变量绑定 + 断言），合计约 1.5 小时/次。这些 TC 已有多步操作描述，但生成的脚本中步骤间输出解析和变量传递仍需手工实现。当前 testkit 提供 `RunCLI`/`RunCLIWithResult`/`WithRetry`（CLI 执行）、`ReadProjectFile`/`FileContains`/`FileNotContains`（文件断言）等 helper，但缺少从命令输出中按正则/JSON 路径提取字段的 `extractField` helper，步骤间状态传递代码仍需手写。pm-work-tracker 的 TC-026 仅 API setup 就占测试代码的 60%（20 行 setup / 33 行总计）。随着 language 从 6 个增长，每新增一个 language 需重复相同的手工模式，成本线性增长
- **回归风险量化**：如果发生跨命令回归（如 task claim 后 submit 因状态链路断裂而失败），当前单步测试无法捕获。诊断成本预计 2-4 小时（需要手动复现多步流程），修复 1-3 小时，可能导致 patch release 延迟
- 其他 language 用户（javascript, python）同样面临单动作覆盖不足的问题

## Proposed Solution

将 e2e 测试管道从"单动作模型"升级为**双模型工作流系统**，由步骤内容自动推断执行模型：

1. **Call-chain 模型**（CLI + API）：步骤间通过 `{stepN.field}` 声明式引用实现状态传递
2. **Interactive 模型**（TUI + Web-UI + Mobile-UI）：步骤间通过 action→assertion 原语描述用户交互
3. 每个 TC 纯粹属于单一 Type，步骤不跨模型。跨 Type 工作流（如 CLI 创建资源 + Web-UI 验证）拆解为多个独立 TC，通过共享 Setup 声明复用测试数据准备
4. 独立的 `Setup` 声明（复用 Call-chain 格式），gen-test-scripts 生成 `beforeAll`/`beforeEach` 代码
5. gen-test-cases 根据 AC 语义自动判断步数和模型类型，无需配置开关或显式标志
6. 单动作 = 两种模型的 1-step 退化形式

### Type → Model Dispatch

TC 的 Type 字段（已存在于当前格式）驱动执行模型选择：

| Type | Model | 状态传递机制 |
|------|-------|-------------|
| CLI | Call-chain | `{stepN.field}` — stdout/stderr 解析 |
| API | Call-chain | `{stepN.field}` — JSON response body 解析 |
| TUI | Interactive | action→assertion — 进程内 Model 状态连续性 |
| Web-UI | Interactive | action→assertion — 浏览器内 DOM 状态连续性 |
| Mobile-UI | Interactive | action→assertion — 设备内 UI 状态连续性 |

> **Integration type migration**: 当前 forge-cli 的 6 条 Integration TC（TC-067 至 TC-071、TC-082）使用 `Type: Integration`。在新模型下，这些 TC 的实际操作全部是 CLI 命令链（forge feature → task claim → submit），应映射为 `Type: CLI`（Call-chain 模型）。迁移策略：新模型保持 `Integration` 作为向后兼容的 Type 值，在 gen-test-scripts 中将 `Integration` 自动映射到 `CLI` 的 Call-chain 处理逻辑。现有 TC 无需修改 Type 字段。

gen-test-scripts 根据 Type 字段值自动选择对应模型，无需额外的 `Workflow` 标志。当 TC 的 Steps 包含 2+ 步骤且步骤间存在状态传递（`{stepN.field}` 引用或连续 action→assertion 链），即被识别为工作流测试。

### `{stepN.field}` Reference Grammar

`{stepN.field}` 是 Call-chain 模型的核心状态传递机制，遵循以下语法规则：

```
reference := "{" step_ref "." field_name "}"
step_ref  := "step" positive_integer     // 引用前序步骤（N < 当前步骤编号）
field_name := [a-z][a-z0-9_]*            // 小写字母开头，允许小写字母/数字/下划线
```

**规则**：
- `N` 必须小于当前步骤编号（只能引用前序步骤的输出，不能前向引用）
- `field_name` 必须在引用步骤的 Fact Table `EXTRACT_*` key 中有对应定义
- 嵌套字段访问使用下划线分隔：`{step2.data_items_0_id}` 而非 `{step2.data.items[0].id}`
- 生成时校验：gen-test-scripts 必须验证所有 `{stepN.field}` 引用的 N 在范围内且 field_name 在 EXTRACT_* key 中存在，校验失败阻止脚本生成
- 运行时失败语义：`extractField` 匹配失败时 fail-fast（`t.Fatalf("Step N/M: field %q not found in output, raw: %s", field, rawOutput)`），跳过后续步骤

### Model Architecture

```
TC (Type-driven dispatch)
├── Setup (Call-chain format, optional)
│   ├── Named Setup block (reusable via Setup-Ref)
│   └── API calls + EXTRACT_* → beforeAll/beforeEach code
└── Steps (model determined by Type)
    ├── CLI / API → Call-chain model
    │   ├── discrete call → parse output → next call
    │   └── {stepN.field} declarative state passing
    └── TUI / Web-UI / Mobile-UI → Interactive model
        ├── action → assertion
        └── in-process / in-browser state continuity
```

**Cross-Type Setup Sharing**: 多个 TC 可通过 `Setup-Ref: <name>` 引用同一命名 Setup 块，复用测试数据准备。gen-test-scripts 为引用同一 Setup 的 TC 生成共享的 `beforeAll` 代码。Setup 块的认证自动注入现有 auth 分类体系。

### Workflow Recognition Heuristics (gen-test-cases)

gen-test-cases 从 PRD AC 中识别多步工作流的启发式规则：

1. **时序连接词模式**：AC 文本包含"先...然后..."、"之后"、"接着"、"after X, Y"等时序词，识别为多步工作流
2. **Given/When/Then 多阶段模式**：一个 AC 的 When 块包含多个按顺序执行的操作（如 "When user creates a feature AND claims a task AND submits the task"），拆解为多步 TC
3. **状态依赖模式**：AC 描述的操作需要前一步的输出作为输入（如 "submit the claimed task" 隐含依赖 claim 步骤的 task-id），识别为 Call-chain 工作流
4. **UI 交互序列模式**：AC 描述连续的页面交互（如 "navigate to milestones, click card, change status"），识别为 Interactive 工作流
5. **Fallback 策略**：当启发式无法确定时，保守降级为单步 TC，不强制拼合不确定的多步操作。生成时在 TC 的 Steps 中保留完整操作描述，由人工审阅时决定是否升级为工作流

### Innovation Highlights

**结构性统一**：所有 TC 使用相同的结构格式（Type + Setup + Steps），由 Type 字段驱动执行模型选择。单动作测试是两种模型的 1-step 退化形式，无需独立类别。

**双模型分治**：CLI/API 的状态传递本质是"跨进程输出解析"（stdout/JSON → 提取字段 → 下一进程参数），而 TUI/Web-UI/Mobile-UI 的状态传递本质是"进程内/浏览器内交互连续性"（按键/点击 → 渲染状态变化 → 下一次交互）。两种机制的差异不是实现细节，而是本质区别——强行统一会制造不自然的抽象。

**Setup 作为一等公民**：将测试数据准备从步骤中分离为独立声明，使 TC 聚焦于被测行为而非基础设施。Setup 复用 Call-chain 格式（API 调用链 + `{stepN.field}`），gen-test-scripts 生成 `beforeAll` 代码。

### Before/After Examples

#### Example 1: CLI Workflow (Call-chain model)

**Before**（当前单动作模型）：test-cases.md 只能描述独立步骤，步骤间状态传递依赖手工验证：

```markdown
## TC-067: Agent Flow — claim to submit lifecycle
- **Steps**:
  1. Run `forge feature wf-test-1`
  2. Run `forge task claim`
  3. Run `forge task submit T-wf-1 --result success --summary "done"`
```

**After**（Call-chain 模型）：test-cases.md 声明步骤间依赖：

```markdown
## TC-W01: Agent Flow — claim to submit lifecycle
- **Type**: CLI
- **Steps**:
  1. Run `forge feature {feature.slug}` → capture `context`
  2. Run `forge task claim` → capture `task_id` from stdout
  3. Run `forge task submit {step2.task_id} --result success --summary "done"`
     → verify `{step2.task_id}` status is "completed" in index.json
- **Expected**: All steps succeed; task transitions pending → in_progress → completed
```

gen-test-scripts 为 go language 生成：

```go
func TestTC_W01_AgentClaimSubmitLifecycle(t *testing.T) {
    // Traceability: TC-W01 -> Spec Agent Task Execution Flow
    out1 := runCLI(t, "forge", "feature", "wf-test-1")

    out2 := runCLI(t, "forge", "task", "claim")
    taskID := extractField(t, out2, `Claimed task:\s+(\S+)`) // VERIFY: pattern from Fact Table

    out3 := runCLI(t, "forge", "task", "submit", taskID, "--result", "success", "--summary", "done")

    status := getTaskStatus(t, taskID)
    assert.Equal(t, "completed", status)
}
```

#### Example 2: TUI Workflow (Interactive model, from agent-forensic)

**Before**：手工编写 Bubble Tea TUI 测试，gen-test-scripts 无法生成：

```go
func TestSessionFlow_LoadToCallTreeToDetail(t *testing.T) {
    sessions := loadFixtureSessions(t, "session_normal.jsonl")
    m, cleanup := initAppWithSessions(t, sessions)
    defer cleanup()

    m, cmd := sendSpecialKey(m, tea.KeyEnter)
    m = dispatchCmd(m, cmd)
    view := m.View()
    viewContains(t, view, "●")
    // ... manual multi-step assertions
}
```

**After**（Interactive 模型）：test-cases.md 声明交互流程：

```markdown
## TC-T01: Navigate session → call tree → diagnosis
- **Type**: TUI
- **Init**: fixture "session_normal.jsonl", session selected
- **Steps**:
  1. key "2"              → view contains "●"
  2. key Enter            → view contains "▼"
  3. key "j"              → view contains "tool_use"
  4. key "d" await        → view contains "诊断摘要"
  5. key Escape           → view not-contains "诊断摘要"
```

gen-test-scripts 为 go language 生成（Bubble Tea 专用 helper）：

```go
func TestTC_T01_NavigateSessionCallTreeDiagnosis(t *testing.T) {
    // Traceability: TC-T01 -> Spec TUI Navigation Flow
    sessions := loadFixtureSessions(t, "session_normal.jsonl")
    m, cleanup := initAppWithSession(t, sessions)
    defer cleanup()

    // Step 1: Focus call tree
    m = sendKeys(m, "2")
    viewContains(t, m.View(), "●")

    // Step 2: Expand turn
    m, _ = sendSpecialKey(m, tea.KeyEnter)
    viewContains(t, m.View(), "▼")

    // Step 3: Move to tool entry
    m = sendKeys(m, "j")
    viewContains(t, m.View(), "tool_use")

    // Step 4: Open diagnosis (await for async cmd)
    m, cmd := sendKey(m, "d")
    m = dispatchCmd(m, cmd)
    viewContains(t, m.View(), "诊断摘要")

    // Step 5: Close diagnosis
    m, _ = sendSpecialKey(m, tea.KeyEscape)
    viewNotContains(t, m.View(), "诊断摘要")
}
```

#### Example 3: Web-UI Workflow (Interactive model + Setup, from pm-work-tracker)

**Before**：手工编写 Playwright 测试，混合 API setup 和 UI 交互：

```typescript
test('TC-026: Cancel milestone auto-unbinds associated MIs', async ({ page }) => {
    // Manual API setup: create map, milestone, MI, change status...
    const mapBizKey = await createMilestoneMap(authCurl, teamBizKey, 'Cancel MS Map');
    const msBizKey = await createMilestone(authCurl, teamBizKey, mapBizKey, 'Cancel MS', '2026-12-01');
    await changeMilestoneStatus(authCurl, teamBizKey, msBizKey, 'in_progress');
    const miRes = await authCurl('POST', `/v1/teams/${teamBizKey}/main-items`, {
        body: JSON.stringify({ title: 'Bound MI', milestoneKey: msBizKey, ... }),
    });
    // Manual UI interaction...
    await page.reload();
    await findCardByTitle(page, 'Cancel MS Map').click();
    // ... more manual steps
});
```

**After**（Interactive 模型 + Setup 声明）：test-cases.md 分离 setup 和测试步骤：

```markdown
## TC-W04: Cancel milestone auto-unbinds associated MIs
- **Type**: Web-UI
- **Setup**:
  1. POST /v1/teams/{team}/milestone-maps {"mapName":"Cancel MS Map"} → capture mapBizKey
  2. POST /v1/teams/{team}/milestone-maps/{step1.mapBizKey}/milestones
     {"milestoneName":"Cancel MS","expectedEndDate":"2026-12-01"} → capture msBizKey
  3. PUT /v1/teams/{team}/milestones/{step2.msBizKey}/status {"status":"in_progress"}
  4. POST /v1/teams/{team}/main-items {"title":"Bound MI","milestoneKey":"{step2.msBizKey}",...}
- **Steps**:
  1. navigate /milestones → page loaded
  2. click card "Cancel MS Map" → timeline view visible
  3. click milestone node → detail panel visible with "Bound MI"
  4. click badge-status → dropdown visible
  5. click option "已取消" → badge shows "已取消"
  6. verify detail panel contains no "Bound MI"
```

gen-test-scripts 为 javascript language 生成：

```typescript
test.describe('TC-W04: Cancel milestone auto-unbinds associated MIs', () => {
  let mapBizKey: string;
  let msBizKey: string;

  test.beforeAll(async () => {
    const token = await getApiToken(apiBaseUrl);
    const authCurl = createAuthCurl(apiBaseUrl, token);

    // Setup Step 1: Create milestone map
    const res1 = await authCurl('POST', `/v1/teams/${teamBizKey}/milestone-maps`, {
      body: JSON.stringify({ mapName: 'Cancel MS Map' }),
    });
    mapBizKey = extractBizKey(parseApiBody(res1.body));

    // Setup Step 2: Create milestone
    const res2 = await authCurl('POST', `/v1/teams/${teamBizKey}/milestone-maps/${mapBizKey}/milestones`, {
      body: JSON.stringify({ milestoneName: 'Cancel MS', expectedEndDate: '2026-12-01' }),
    });
    msBizKey = extractBizKey(parseApiBody(res2.body));

    // Setup Step 3: Change status to in_progress
    await authCurl('PUT', `/v1/teams/${teamBizKey}/milestones/${msBizKey}/status`, {
      body: JSON.stringify({ status: 'in_progress' }),
    });

    // Setup Step 4: Create bound MI
    await authCurl('POST', `/v1/teams/${teamBizKey}/main-items`, {
      body: JSON.stringify({ title: 'Bound MI', milestoneKey: msBizKey, ... }),
    });
  });

  test('TC-W04: Cancel milestone auto-unbinds associated MIs', async ({ page }) => {
    // Step 1: Navigate to milestones page
    await page.goto(`${baseUrl}/milestones`);
    await page.waitForLoadState('networkidle');

    // Step 2: Click card to enter timeline view
    await findCardByTitle(page, 'Cancel MS Map').click();
    await expect(page.locator('[data-testid="timeline-view"]')).toBeVisible({ timeout: 5000 });

    // Step 3: Click milestone node to open detail panel
    await page.locator('[data-testid="milestone-node"]').first().click();
    await expect(page.locator('[data-testid="detail-panel"]')).toBeVisible({ timeout: 5000 });

    // Steps 4-5: Cancel milestone via status dropdown
    await page.locator('[data-testid="badge-status"]').click();
    const dropdown = page.locator('[data-testid="dropdown-status-options"]');
    await expect(dropdown).toBeVisible({ timeout: 3000 });
    await dropdown.getByText('已取消').click();

    // Step 6: Verify MI is unbound
    await expect(page.locator('[data-testid="detail-panel"]')).not.toContainText('Bound MI');
  });
});
```

#### Example 4: API Workflow (Call-chain model)

**Before**（当前单动作模型）：API 测试只能验证单个端点，跨端点依赖需要手工硬编码：

```markdown
## TC-045: API — user registration to resource creation
- **Steps**:
  1. POST /api/users {"name":"test"}
  2. POST /api/auth/login {"userId":"hardcoded-test-user"}
  3. POST /api/resources {"authToken":"hardcoded-token"}
```

**After**（Call-chain 模型）：步骤间通过 `{stepN.field}` 自动传递状态：

```markdown
## TC-W05: API — user registration to resource creation
- **Type**: API
- **Steps**:
  1. POST /api/users {"name":"test"} → capture user_id, token
  2. POST /api/auth/login {"userId":"{step1.user_id}"} → capture authToken
  3. POST /api/resources {"authToken":"{step2.authToken}"} → verify status 201
- **Expected**: All steps succeed; resource created for authenticated user
```

gen-test-scripts 为 javascript language（API 子模板）生成：

```typescript
test.describe('TC-W05: API — user registration to resource creation', () => {
  test('TC-W05: Full API workflow', async ({ request }) => {
    // Step 1: Create user
    const res1 = await request.post('/api/users', { data: { name: 'test' } });
    expect(res1.status()).toBe(201);
    const body1 = await res1.json();
    const userId = body1.user_id;
    const token = body1.token;

    // Step 2: Authenticate
    const res2 = await request.post('/api/auth/login', { data: { userId } });
    expect(res2.status()).toBe(200);
    const body2 = await res2.json();
    const authToken = body2.authToken;

    // Step 3: Create resource with auth
    const res3 = await request.post('/api/resources', {
      headers: { Authorization: `Bearer ${authToken}` },
      data: {},
    });
    expect(res3.status()).toBe(201);
  });
});
```

## Requirements Analysis

### Key Scenarios

- **Happy path (CLI)**: `forge feature → task claim → task submit`，步骤间传递 feature-slug 和 task-id（Call-chain 模型）
- **Happy path (TUI)**: `focus call tree → expand turn → navigate entry → open diagnosis`，步骤间传递键盘导航状态（Interactive 模型，agent-forensic）
- **Happy path (Web-UI)**: `navigate milestones → create map → create milestone → change status`，UI 交互 + API 数据准备（Interactive + Setup，pm-work-tracker）
- **Happy path (API)**: `create user → authenticate → create resource`，步骤间传递 user-id 和 auth-token（Call-chain 模型）。Before/after 示例见上方 Example 4
- **Edge case**: 工作流中间步骤失败，后续步骤应被跳过
- **Edge case**: 步骤输出格式变化导致字段解析失败（Call-chain 模型）
- **Backward compat**: 现有单步 TC 作为退化形式继续工作，无需迁移

### Non-Functional Requirements

- **Backward compatibility**: 现有 test-cases.md 和生成的测试脚本零改动继续工作
- **Language parity**: 所有 6 个 language 统一支持双模型，无功能差距
- **Template simplicity**: 每个 language 的 generate.md 只增加工作流执行逻辑，不改变现有单步测试的生成路径
- **Error reporting for multi-step failures**: 工作流测试失败时，CI 输出必须标识失败步骤编号（如 "Step 2/3 failed"）并展示该步骤的输出内容。覆盖以下失败模式：
  - Call-chain: 字段解析失败（"Step N/M failed: field X not found in output"）、断言失败（"Step N/M assertion failed: expected Y, got Z"）
  - Interactive: action 超时（"Step N/M timeout: element X not found within 5s"）、assertion 失败（"Step N/M assertion failed: view missing expected text Y"）
  - Setup: beforeAll 失败（"Setup Step N failed: API returned status 500"）
- **Field reference coupling**: `{stepN.field}` 引用必须在生成时（gen-test-scripts）验证 field 名称在 Fact Table 中有对应的提取规则；生成时校验失败应报错阻止脚本生成
- **Framework-aware helper generation**: TUI 交互 helper 需要感知具体 TUI 框架（Bubble Tea、ratatui、textual），由 generate.md 定义框架映射规则
- **Generation and execution performance**: 多步 TC 的生成时间不超过同等单步 TC 生成时间的 2 倍；工作流测试的执行时间不超过对应单步测试执行时间之和的 120%（不含 Setup 的 beforeAll 时间），确保 e2e suite 总运行时间的增量可控

### Constraints & Dependencies

- 依赖现有 language 系统的 generate.md（位于 `forge-cli/pkg/profile/languages/<lang>/generate.md`）和 `forge testing` CLI 命令
- **TUI 能力前置条件**：当前 Go 的 `languageCapabilities` 仅声明 `{"api", "cli"}`（见 `forge-cli/pkg/profile/embed.go`），不包含 TUI。Go 的 `generate.md` 已有 TUI Testing 部分（golden file 比较模式），但 language 注册层面缺少 TUI 接口。Phase 1 需要将 TUI 添加到 Go 的 `languageCapabilities` 中作为前置步骤（改动量极小：`embed.go` 中 `{"api", "cli"}` → `{"api", "cli", "tui"}`）。java 和 rust 的 TUI 扩展留待 Phase 2（需先确认对应 TUI 框架：Java 为 Lanterna/Text-UI，Rust 为 ratatui）
- `{stepN.field}` 的 field 名称需要与 CLI/API 的输出格式对齐。**这是一个新增的运行时输出解析层**，需要为每种输出格式（CLI stdout text、API JSON body）定义提取规则。现有 Fact Table 是从源代码中提取静态值（端口、路由、data-testid），不直接适用于运行时输出解析。新机制将参考 Fact Table 的 Key/Value/Source 结构，但作用域不同（静态侦察 vs 运行时提取）。示例 `EXTRACT_*` keys：
  - `EXTRACT_TASK_ID` → pattern `Claimed task:\s+(\S+)`，source: `forge task claim` stdout
  - `EXTRACT_FEATURE_SLUG` → pattern `Feature\s+([\w-]+)\s+created`，source: `forge feature` stdout
  - `EXTRACT_TASK_STATUS` → JSON path `$.status`，source: `tasks/index.json` 文件内容
  - `EXTRACT_SUBMIT_RESULT` → pattern `Task\s+\S+\s+submitted\s+\((\w+)\)`，source: `forge task submit` stdout
- TUI 交互 helper 是框架特有的：Bubble Tea 用 `sendKey/dispatchCmd/viewContains`，ratatui 用 `TestBackend/Event::Key`，textual 用 `pilot.press/query_one`。generate.md 负责定义框架映射
- **Setup 步骤的认证集成**：Setup 步骤中的 API 调用复用 gen-test-scripts 现有的 auth 分类体系（login-test、auth-required、public、custom-auth）。Setup 生成的 `beforeAll` 代码自动注入对应 auth 策略（如 `storageState` 缓存凭据、`authCurl` helper），无需在 Setup 声明中重复认证配置
- **Setup-Ref 范围**：`Setup-Ref`（跨 TC 的命名 Setup 共享）纳入 **Phase 2 Out of Scope**——Phase 1 中每个 TC 独立声明自己的 Setup 块，不实现命名引用机制。Phase 2 根据实际需求决定是否引入 Setup-Ref
- 不改变 graduate-tests 和 run-e2e-tests 的工作流（工作流测试的毕业和执行流程与现有测试一致）

## Alternatives & Industry Benchmarking

### Industry Solutions

- **Playwright**: 原生支持多步骤测试（test.step），步骤间通过变量共享状态。UI 自动化领域的事实标准。
- **Go testing**: 社区实践推荐 table-driven subtests 模拟工作流，但无框架级支持。Bubble Tea 应用通过直接操作 Model 对象实现工作流测试。
- **Postman/Newman**: API 工作流测试通过 collection runner + 环境变量传递状态。
- **Robot Framework**: keyword-driven 模式天然支持工作流，步骤间通过变量传递。
- **Maestro**: YAML flow 天然支持多步交互（tapOn → assertVisible → inputText），步骤间通过 env 变量传递。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| **Adopt Robot Framework** as workflow orchestrator | Industry: Robot Framework keyword-driven testing | 成熟的 keyword-driven 模型天然支持工作流，变量传递是内置能力 | 引入 Python + Robot Framework 技术栈，与现有 6 个 language 的原生语言不一致；无法复用现有模板和 Fact Table 机制 | Rejected: 技术栈异构成本高于收益 |
| Single unified model with `{stepN.field}` for all types | Playwright state-passing pattern | 模型单一，实现简单 | 强行用输出解析模型描述进程内交互会产生不自然代码。例如 TUI 步骤 `key "j" → view contains "tool_use"` 在统一模型下需写成 `runAction("j") → parseViewOutput() → assertContains("tool_use")`，将 Bubble Tea 的 `Model.View()` 返回值当作 "stdout" 解析，丧失 Interactive 模型的语义清晰度。**实证**：agent-forensic 项目早期的 TUI 测试尝试了统一输出解析方式，结果 3 个测试中有 2 个因 `Model.View()` 返回值包含 ANSI 转义序列导致正则匹配失败，最终重写为 Interactive 模式的 `viewContains` 直接断言。每种类型的模板都需要大量条件分支来适配本质不同的状态传递机制 | Rejected: 单一模型在 TUI/Web-UI 上制造了不自然的抽象，模板复杂度反超双模型 |
| **Keep Forge single-action + use external tools for workflows** | Playwright test.step + Postman collection runner | 不修改 Forge 管道，直接使用成熟工具的工作流能力；零实现成本 | 工作流测试脱离 Forge 管道：无法从 PRD AC 自动生成工作流 TC，无法通过 eval-test-cases 评分验证质量，无法通过 graduate-tests 毕业，也无法利用 Fact Table 进行字段对齐校验。每个项目需独立维护外部工作流测试套件，与 Forge 的 language 系统和 CI 管道完全割裂 | Rejected: 放弃了 Forge 管道的核心价值（从 PRD 到测试的自动化），且维护两套测试体系成本更高 |
| **Dual model (Call-chain + Interactive)** | Playwright (call-chain) + Bubble Tea testing pattern (interactive) | 两种模型分别匹配两种本质不同的状态传递机制；Type 驱动自动分派；每种模型内部保持退化形式统一 | 需同时实现两套步骤执行逻辑，但每套逻辑更简单（各只处理一种模式） | **Selected: 双模型分治消除了单一抽象的强制统一，每种模型的表达力更强** |

## Feasibility Assessment

### Technical Feasibility

- **高**。两种模型都有成熟的实践验证：
  - Call-chain 模型：当前 gen-test-scripts 已有 helpers（runCLI, curl, withRetry），工作流执行只需顺序调用 + 变量传递
  - Interactive 模型：agent-forensic 的 `tests/e2e_go/` 已有完整的 TUI 工作流测试模式（sendKey + viewContains + dispatchCmd），pm-work-tracker 的 Playwright 测试已有完整的 Web-UI 工作流模式（page.click + expect + beforeAll API setup）
- **`{stepN.field}` 需要新建运行时输出解析层**：现有 Fact Table 从源代码中提取静态侦察值（API_PORT、CLI_BINARY、TESTID_* 等）。`{stepN.field}` 需要的是从 CLI stdout / API JSON response 中按提取规则获取运行时值——这是一个新机制。具体做法：在 Fact Table 中新增 `EXTRACT_*` 类 key（如 `EXTRACT_TASK_ID` → pattern `Claimed task:\s+(\S+)`，source 标注为 CLI help text 或已知输出格式），gen-test-scripts 生成 `extractField()` 调用时引用这些 key。**注意**：当前 testkit 有 `RunCLI`/`RunCLIWithResult`/`WithRetry`（CLI 执行）和 `ReadProjectFile`/`FileContains`/`FileNotContains`（文件断言）等 helper，但无输出字段解析 helper。`extractField`（stdout regex 提取）和 `extractJSONField`（JSON path 提取）需要作为新 helper 实现在 `testkit/extract.go`。这是 Phase 1 的核心交付物之一（T1-extract-helper）
- **Setup 复用 Call-chain 格式**：gen-test-scripts 已有 API 调用代码生成能力（javascript language 的 `authCurl` helper），Setup 只是将其生成到 `beforeAll` 而非测试函数内
- **generate.md 增量修改**：6 个 generate.md 已有成熟的结构，新增 workflow 分支是增量修改（在现有 test function 内添加步骤循环 + 变量绑定），不改变单步测试的生成路径

### Resource & Timeline

- 1 名工程师，分 3 个阶段交付：
  - **Phase 1（3 周）**：核心双模型引擎 + go + javascript，以 **forge-cli 自身工作流**（feature → task claim → submit）作为主验证目标。
    - Task breakdown：
      1. **T1-TUI-capability**: go/java/rust 的 `languageCapabilities` 添加 TUI 接口（`embed.go` 三行改动 + 回归测试）
      2. **T1-extract-helper**: 新建 `extractField`（stdout regex）和 `extractJSONField`（JSON path）helper（`testkit/extract.go`），含单元测试
      3. **T1-EXTRACT-keys**: 定义 forge-cli 的 `EXTRACT_*` Fact Table keys（`EXTRACT_TASK_ID`、`EXTRACT_FEATURE_SLUG`、`EXTRACT_TASK_STATUS`、`EXTRACT_SUBMIT_RESULT`），更新 generate.md 使用新 key
      4. **T1-call-chain-engine**: gen-test-scripts SKILL.md + types/cli.md + types/api.md 新增 Call-chain 步骤处理（`{stepN.field}` 解析 + 变量绑定 + 顺序执行 + fail-fast）
      5. **T1-setup-generator**: gen-test-scripts 新增 Setup→beforeAll 代码生成（复用 Call-chain 格式，集成现有 auth 分类体系）
      6. **T1-go-workflow**: go generate.md + templates 更新，以 forge-cli TC-067（Agent Flow — claim to submit lifecycle）端到端验证
      7. **T1-js-workflow**: javascript generate.md + templates 更新，以 API workflow（user→auth→resource）端到端验证
    - 验收：forge-cli 的 task lifecycle 工作流 TC 从 test-cases.md 自动生成可编译运行的测试
  - **Phase 2（2 周）**：剩余 4 个 language 模板更新 + gen-test-cases 工作流识别 + 外部参考项目验证。交付物：python、java、rust、mobile 的 generate.md 更新、gen-test-cases 各 type 指令文件的工作流识别启发式支持。验收：agent-forensic 的 TUI 导航工作流 TC 和 pm-work-tracker 的 Web-UI milestone 工作流 TC 从 test-cases.md 自动生成可编译脚本
  - **Phase 3（0.5 周）**：eval-test-cases 评分维度更新 + 集成测试 + 向后兼容验证。交付物：eval rubric 工作流评分维度、现有 81 个 TC 零回归。验收：`just e2e-compile` 通过，现有单步 TC 输出与改动前 diff 为空
- 总计 12 个 task（Phase 1: 7, Phase 2: 3, Phase 3: 2），5.5 周日历时间

### Dependency Readiness

- language 系统和 `forge testing` CLI 命令已稳定运行（simplify-testing-model 重构后，6 个 language 的 generate.md 结构已统一）
- Fact Table 机制已稳定（字段提取在现有 20+ e2e 测试上连续通过），运行时提取规则将作为新 key 类型扩展 Fact Table
- 不依赖外部服务或 API

## Impact Analysis

### Tier 1: 必须修改

| 组件 | 文件 | 变更 |
|------|------|------|
| gen-test-cases | SKILL.md | 新增工作流识别启发式规则（Step 2-3 区域） |
| gen-test-cases | types/cli.md | CLI 工作流 TC 格式（Steps 多步 + `{stepN.field}`） |
| gen-test-cases | types/api.md | API 工作流 TC 格式（Steps 多步 + `{stepN.field}`） |
| gen-test-cases | types/tui.md | TUI 工作流 TC 格式（Steps 多步 + action→assertion） |
| gen-test-cases | types/ui.md | Web-UI 工作流 TC 格式（Steps 多步 + Setup 声明） |
| gen-test-cases | types/mobile.md | Mobile-UI 工作流 TC 格式（Steps 多步 + action→assertion） |
| gen-test-cases | templates/cli-test-cases.md | CLI 工作流模板结构 |
| gen-test-cases | templates/api-test-cases.md | API 工作流模板结构 |
| gen-test-cases | templates/tui-test-cases.md | TUI 工作流模板结构 |
| gen-test-cases | templates/ui-test-cases.md | Web-UI 工作流模板结构 |
| gen-test-cases | templates/mobile-test-cases.md | Mobile-UI 工作流模板结构 |
| gen-test-scripts | SKILL.md | Call-chain + Interactive 双模型处理、Setup→beforeAll 生成 |
| gen-test-scripts | types/cli.md | CLI Call-chain 步骤规则 |
| gen-test-scripts | types/api.md | API Call-chain 步骤规则 |
| gen-test-scripts | types/tui.md | TUI Interactive 步骤规则 |
| gen-test-scripts | types/ui.md | Web-UI Interactive 步骤规则 + Setup 生成 |
| gen-test-scripts | types/mobile.md | Mobile-UI Interactive 步骤规则 |
| go | generate.md | Call-chain(CLI) + Interactive(TUI) 步骤规则 |
| go | templates/test-file.go | 工作流测试函数骨架 |
| go | templates/helpers.go | 新增 TUI 交互 helper |
| javascript | generate.md | Call-chain(API) + Interactive(Web-UI) 步骤规则 |
| javascript | templates/helpers.ts | 新增工作流执行 helper |
| javascript | templates/api.spec.ts | Call-chain 工作流模板 |
| javascript | templates/playwright-ui.spec.ts | Interactive 工作流模板 |
| python | generate.md | Call-chain(API/CLI) 步骤规则 |
| python | templates/test_file.py | 工作流测试函数骨架 |
| python | templates/conftest.py | Setup fixture 生成 |
| java | generate.md | Call-chain + Interactive(TUI) 步骤规则 |
| java | templates/TestFile.java | 工作流测试函数骨架 |
| rust | generate.md | Call-chain(CLI) + Interactive(TUI) 步骤规则 |
| rust | templates/test_file.rs | 工作流测试函数骨架 |
| rust | templates/helpers.rs | TUI 交互 helper |
| mobile | generate.md | Interactive(Mobile-UI) 步骤规则 |

### Tier 2: 应该修改（6 文件）

| 文件 | 变更 |
|------|------|
| eval/rubrics/cli-test-cases.md | 新增工作流 TC 评分维度（步骤间状态传递有效性、Setup 正确性） |
| eval/rubrics/api-test-cases.md | 新增工作流 TC 评分维度 |
| eval/rubrics/tui-test-cases.md | 新增工作流 TC 评分维度 |
| eval/rubrics/ui-test-cases.md | 新增工作流 TC 评分维度（含 Setup 声明评分） |
| eval/rubrics/mobile-test-cases.md | 新增工作流 TC 评分维度 |
| eval/rubrics/test-cases.md | Legacy 单体 rubric 新增工作流评分维度（向后兼容） |

### 不受影响

| 组件 | 理由 |
|------|------|
| 所有 language 的 run.md | 运行时不区分工作流 vs 单步 |
| 所有 language 的 graduate.md | 毕业流程不变 |
| run-e2e-tests / graduate-tests SKILL.md | 工作流 TC 的执行和毕业与单步 TC 一致 |
| config.go | language 注册和检测机制不变 |
| embed.go | `languageCapabilities` 仅新增 TUI 到 go/java/rust（现有数据项修改，非结构变更） |

## Scope

### In Scope

- test-cases.md per-type 模板重设计：双模型结构（Type + Setup + Steps）
- gen-test-cases per-type 指令文件（5 个 types/*.md）：工作流识别启发式 + Setup 声明生成
- gen-test-scripts SKILL.md + per-type 指令文件（5 个 types/*.md）：Call-chain + Interactive 双模型步骤处理 + Setup→beforeAll 生成
- 全部 6 个 language 的 generate.md 策略更新
- 核心模板更新（go + javascript Phase 1，其余 Phase 2）
- eval-test-cases per-type rubric 评分维度更新（6 个 rubric 文件）

### Out of Scope

- 现有 test-cases.md 文件迁移（自然兼容，无需迁移）
- graduate-tests / run-e2e-tests 变更（工作流测试的毕业和执行流程不变）
- forge-cli 自身的单元/集成测试变更
- TUI 框架自动检测（generate.md 硬编码框架映射规则，不做框架推断）
- 跨 Type 工作流（如 CLI + Web-UI 组合步骤的单条 TC）。当前设计通过拆解为多条独立 TC 实现
- `Setup-Ref` 命名引用机制（Phase 1 中每个 TC 独立声明 Setup，不实现跨 TC 共享）
- java/rust 的 TUI 能力扩展（仅 go 在 Phase 1 添加 TUI；java/rust 的 TUI 框架映射留待 Phase 2 根据实际需求决定）

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 6 个 language 的双模型实现引入不一致 | M | H | 定义统一的 workflow 抽象接口规范，每个 language 只做语言+框架映射 |
| `{stepN.field}` 的 field 名称与实际输出不匹配 | M | M | 新增 `EXTRACT_*` Fact Table key 作为格式契约；匹配失败时输出 raw stdout 以便调试 |
| TUI helper 是框架特有的，Bubble Tea/ratatui/textual 需要不同实现 | M | M | 由 generate.md 定义框架→helper 映射，不建立上层抽象 |
| 工作流中间步骤失败导致误报 | L | H | Call-chain 用 require（fail-fast）；Interactive 按步骤隔离断言；失败输出标识步骤编号和失败类型 |
| gen-test-cases 的工作流识别质量不足 | M | M | Phase 1 实现 dry-run 模式：gen-test-cases 生成 TC 后先进行结构化 schema 校验（校验 Type 合法性、Steps 非空、`{stepN.field}` 引用的 step 序号在范围内），校验通过后才写入 test-cases.md；Phase 1 用 forge-cli 的 PRD 样本作为端到端验证 |
| 运行时输出解析层与静态 Fact Table 的概念混淆 | M | M | 在 Fact Table 文档中明确区分两类 key：静态侦察 key（`API_PORT`、`CLI_BINARY`）和运行时提取 key（`EXTRACT_TASK_ID`），前者 source 为源代码文件，后者 source 为 CLI help text 或已知输出格式 |

## Success Criteria

- [ ] go language 能从 CLI 工作流 TC 生成包含 2+ 步骤的 forge-cli 集成测试（feature → task claim → submit），且步骤间 task_id 正确传递
- [ ] go language 能从 TUI 工作流 TC 生成包含 2+ 步骤的 Bubble Tea 测试，且步骤间导航状态正确传递（以 agent-forensic 的 session→call tree→diagnosis 流程验证）
- [ ] javascript language 能从 Web-UI 工作流 TC + Setup 声明生成包含 beforeAll API 准备 + 2+ 步骤 UI 交互的 Playwright 测试（以 pm-work-tracker 的 milestone cancel→unbind MI 流程验证）
- [ ] 现有单步 TC 生成的测试代码与改动前完全一致（零回归）
- [ ] gen-test-cases 从 3 个不同 Type（CLI + TUI + Web-UI）的 PRD 样本生成结构合法的多步 TC（每条含 2+ 步骤），经 eval-test-cases 评分均 ≥ 75/100（当前单步 TC 基线平均分）
- [ ] test-cases.md 模板的 Steps 结构对全部 5 种 capability 类型（cli, api, tui, web-ui, mobile-ui）各生成一条可编译的工作流测试，证明每种类型都能正确表达至少 2 步工作流
- [ ] 所有 6 个 language 的 generate.md 支持双模型步骤生成，通过编译验证
- [ ] 工作流中间步骤解析失败时，生成的测试输出明确的错误标识（格式为 "Step N/M failed: {failure-type}: {details}"），且跳过后续步骤（Call-chain 用 fail-fast，Interactive 用步骤隔离）

## Next Steps

- Proceed to `/write-prd` to formalize requirements
