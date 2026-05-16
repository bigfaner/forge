---
created: 2026-05-15
updated: 2026-05-16
author: "faner"
status: Draft
---

# Proposal: Unified Workflow Test Model

## Problem

Forge 的 e2e 测试管道以"单动作"为粒度生成 test case，每个 TC 验证一条命令/一个端点/一个页面交互。但真实用户场景是跨多个操作的**工作流**（如 `forge init → task claim → task submit → verify`），当前管道无法生成覆盖完整用户旅程的 e2e 测试。

### Evidence

- forge-cli 的 70+ 单元测试和 20+ e2e 测试全部是单命令粒度（`exec.Command("forge", <single-command>)`），无跨命令工作流测试
- gen-test-cases 从 PRD AC 提取验证点时，每个 AC 映射为单条 TC，即使 AC 涉及多步操作也被拆散
- 6 个 profile 的模板（go-test, web-playwright, pytest, java-junit, rust-test, maestro）全部围绕单动作设计，无步骤间状态传递机制
- 行业标准（Microsoft Engineering Playbook, IBM）定义 e2e testing 为验证"complete user workflows from start to end"

### Reference Projects

两个真实项目暴露了问题的不同侧面：

**agent-forensic**（Go TUI 应用，go-test profile）：Bubble Tea TUI 的工作流测试使用 `sendKey → viewContains` 模式验证键盘导航流程（如 session → call tree → diagnosis 面板），但当前 gen-test-scripts 无法生成此类测试。

**pm-work-tracker**（Go+React 混合 Web 应用，web-playwright profile）：Playwright 测试中存在大量"API 数据准备 → UI 交互验证"模式（如 TC-026 先通过 API 创建 milestone+MI，再通过 UI 取消 milestone 并验证 MI 解绑），但 gen-test-cases 无法识别此类跨步依赖。

### Urgency

- forge-cli 正在快速迭代，跨命令的回归风险在增长（如 task claim → submit 的状态链路）
- **手工状态传递成本已可量化**：当前 Integration TC（TC-067 至 TC-082）中 16 条测试每条需要约 15 分钟手工编写步骤间状态传递代码（extractField + 变量绑定 + 断言），合计约 4 小时/次。pm-work-tracker 的 TC-026 仅 API setup 就占测试代码的 60%（20 行 setup / 33 行总计）。随着 profile 从 6 个增长，每新增一个 profile 需重复相同的手工模式，成本线性增长。v3.0 迭代期间已出现一次跨命令回归（task claim 后 submit 失败）在单步测试中未被捕获，直到手工集成测试才发现（从引入到发现耗时 2 天，诊断+修复耗时 3 小时，导致 v3.0.1 patch release 延迟 1 天）
- 其他 profile 用户（web-playwright, pytest）同样面临单动作覆盖不足的问题

## Proposed Solution

将 e2e 测试管道从"单动作模型"升级为**双模型工作流系统**，按 TC 的 Type 字段自动选择对应模型：

1. **Call-chain 模型**（CLI + API）：步骤间通过 `{stepN.field}` 声明式引用实现状态传递
2. **Interactive 模型**（TUI + Web-UI + Mobile-UI）：步骤间通过 action→assertion 原语描述用户交互
3. 每个 TC 纯粹属于单一 Type，步骤不跨模型。跨 Type 工作流（如 CLI 创建资源 + Web-UI 验证）拆解为多个独立 TC，通过共享 Setup 声明复用测试数据准备
4. 独立的 `Setup` 声明（复用 Call-chain 格式），gen-test-scripts 生成 `beforeAll`/`beforeEach` 代码
5. gen-test-cases 根据 AC 语义自动判断步数，无需配置开关
6. 单动作 = 两种模型的 1-step 退化形式

### Model Architecture

```
TC (Type-driven dispatch)
├── Setup (Call-chain format, optional)
│   └── API calls → beforeAll/beforeEach code
└── Steps (model determined by Type)
    ├── CLI / API → Call-chain model
    │   ├── discrete call → parse output → next call
    │   └── {stepN.field} declarative state passing
    └── TUI / Web-UI / Mobile-UI → Interactive model
        ├── action → assertion
        └── in-process / in-browser state continuity
```

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
- **Workflow**: true
- **Steps**:
  1. Run `forge feature {feature.slug}` → capture `context`
  2. Run `forge task claim` → capture `task_id` from stdout
  3. Run `forge task submit {step2.task_id} --result success --summary "done"`
     → verify `{step2.task_id}` status is "completed" in index.json
- **Expected**: All steps succeed; task transitions pending → in_progress → completed
```

gen-test-scripts 为 go-test profile 生成：

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
- **Workflow**: true
- **Init**: fixture "session_normal.jsonl", session selected
- **Steps**:
  1. key "2"              → view contains "●"
  2. key Enter            → view contains "▼"
  3. key "j"              → view contains "tool_use"
  4. key "d" await        → view contains "诊断摘要"
  5. key Escape           → view not-contains "诊断摘要"
```

gen-test-scripts 为 go-test profile 生成（Bubble Tea 专用 helper）：

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
- **Workflow**: true
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

gen-test-scripts 为 web-playwright profile 生成：

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
- **Workflow**: true
- **Steps**:
  1. POST /api/users {"name":"test"} → capture user_id, token
  2. POST /api/auth/login {"userId":"{step1.user_id}"} → capture authToken
  3. POST /api/resources {"authToken":"{step2.authToken}"} → verify status 201
- **Expected**: All steps succeed; resource created for authenticated user
```

gen-test-scripts 为 web-playwright profile（API 子模板）生成：

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
- **Profile parity**: 所有 6 个 profile 统一支持双模型，无功能差距
- **Template simplicity**: 每个 profile 的模板只增加工作流执行逻辑，不改变现有单步测试的生成路径
- **Error reporting for multi-step failures**: 工作流测试失败时，CI 输出必须标识失败步骤编号（如 "Step 2/3 failed"）并展示该步骤的输出内容
- **Field reference coupling**: `{stepN.field}` 引用必须在生成时（gen-test-scripts）验证 field 名称在对应 profile 的 Fact Table 中存在；生成时校验失败应报错阻止脚本生成
- **Framework-aware helper generation**: TUI 交互 helper 需要感知具体 TUI 框架（Bubble Tea、ratatui、textual），由 profile 的 generate.md 定义框架映射规则
- **Generation and execution performance**: 多步 TC 的生成时间不超过同等单步 TC 生成时间的 2 倍；工作流测试的执行时间不超过对应单步测试执行时间之和的 120%（不含 Setup 的 beforeAll 时间），确保 e2e suite 总运行时间的增量可控

### Constraints & Dependencies

- 依赖现有 profile 系统的 manifest.yaml 和 generate.md
- `{stepN.field}` 的 field 名称需要与 CLI/API 的输出格式对齐（依赖 Fact Table 的字段提取能力）
- TUI 交互 helper 是框架特有的：Bubble Tea 用 `sendKey/dispatchCmd/viewContains`，ratatui 用 `TestBackend/Event::Key`，textual 用 `pilot.press/query_one`。generate.md 负责定义框架映射。
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
| **Adopt Robot Framework** as workflow orchestrator | Industry: Robot Framework keyword-driven testing | 成熟的 keyword-driven 模型天然支持工作流，变量传递是内置能力 | 引入 Python + Robot Framework 技术栈，与现有 6 个 profile 的原生语言不一致；无法复用现有模板和 Fact Table 机制 | Rejected: 技术栈异构成本高于收益 |
| Single unified model with `{stepN.field}` for all types | Playwright state-passing pattern | 模型单一，实现简单 | 强行用输出解析模型描述进程内交互会产生不自然代码。例如 TUI 步骤 `key "j" → view contains "tool_use"` 在统一模型下需写成 `runAction("j") → parseViewOutput() → assertContains("tool_use")`，将 Bubble Tea 的 `Model.View()` 返回值当作 "stdout" 解析，丧失 Interactive 模型的语义清晰度。同理 Web-UI 的 `click → waitForSelector` 被扭曲为 `runCommand("click") → parseDOM()`。**实证**：agent-forensic 项目早期的 TUI 测试尝试了统一输出解析方式（将 `sendKey` 的返回值统一为 string 再断言），结果 3 个测试中有 2 个因 `Model.View()` 返回值包含 ANSI 转义序列导致正则匹配失败，最终重写为 Interactive 模式的 `viewContains` 直接断言。每种类型的模板都需要大量条件分支来适配本质不同的状态传递机制 | Rejected: 单一模型在 TUI/Web-UI 上制造了不自然的抽象，模板复杂度反超双模型 |
| **Keep Forge single-action + use external tools for workflows** | Playwright test.step + Postman collection runner | 不修改 Forge 管道，直接使用成熟工具的工作流能力；零实现成本 | 工作流测试脱离 Forge 管道：无法从 PRD AC 自动生成工作流 TC，无法通过 eval-test-cases 评分验证质量，无法通过 graduate-tests 毕业，也无法利用 Fact Table 进行字段对齐校验。每个项目需独立维护外部工作流测试套件，与 Forge 的 profile 系统和 CI 管道完全割裂 | Rejected: 放弃了 Forge 管道的核心价值（从 PRD 到测试的自动化），且维护两套测试体系成本更高 |
| **Dual model (Call-chain + Interactive)** | Playwright (call-chain) + Bubble Tea testing pattern (interactive) | 两种模型分别匹配两种本质不同的状态传递机制；Type 驱动自动分派；每种模型内部保持退化形式统一 | 需同时实现两套步骤执行逻辑，但每套逻辑更简单（各只处理一种模式） | **Selected: 双模型分治消除了单一抽象的强制统一，每种模型的表达力更强** |

## Feasibility Assessment

### Technical Feasibility

- **高**。两种模型都有成熟的实践验证：
  - Call-chain 模型：当前 gen-test-scripts 已有 helpers（runCLI, curl, withRetry），工作流执行只需顺序调用 + 变量传递
  - Interactive 模型：agent-forensic 的 `tests/e2e_go/` 已有完整的 TUI 工作流测试模式（sendKey + viewContains + dispatchCmd），pm-work-tracker 的 Playwright 测试已有完整的 Web-UI 工作流模式（page.click + expect + beforeAll API setup）
- **`{stepN.field}` 提取复用 Fact Table 的具体验证**：以 `forge task claim` 为例，当前 Fact Table 已定义字段提取规则 `Claimed task:\s+(\S+)` 匹配 stdout 输出。`{step2.task_id}` 的提取逻辑与此完全一致——两者都是从命令输出中按正则提取命名字段。区别仅为：Fact Table 在 TC 编写时静态定义 pattern，`{stepN.field}` 在 gen-test-scripts 时引用同一 pattern 并生成 `extractField()` 调用。go-test 和 web-playwright profile 已有 `extractField`/`extractBizKey` helper（见 Example 1 和 Example 3 的生成代码），覆盖了 stdout text 和 JSON body 两种输出格式
- **Setup 复用 Call-chain 格式**：gen-test-scripts 已有 API 调用代码生成能力（web-playwright 的 `authCurl` helper），Setup 只是将其生成到 `beforeAll` 而非测试函数内
- **Profile generate.md 增量修改**：6 个 generate.md 已有成熟的结构，新增 workflow 分支是增量修改（在现有 test function 内添加步骤循环 + 变量绑定），不改变单步测试的生成路径

### Resource & Timeline

- 1 名工程师，分 3 个阶段交付：
  - **Phase 1（3 周）**：核心双模型引擎 + go-test + web-playwright，用两个参考项目验证。交付物：test-cases.md 双模型结构定义、Call-chain `{stepN.field}` 解析器、Interactive action→assertion 步骤处理器、Setup→beforeAll 生成器、go-test 和 web-playwright 的 generate.md 更新。验收：agent-forensic 的 TUI 导航工作流 TC 和 pm-work-tracker 的 Web-UI milestone 工作流 TC 从 test-cases.md 自动生成可编译运行的测试
  - **Phase 2（1.5 周）**：剩余 4 个 profile 模板更新 + gen-test-cases 工作流识别。交付物：pytest、java-junit、rust-test、maestro 的 generate.md 更新、gen-test-cases SKILL.md 的双模型工作流识别支持。验收：所有 profile 对同一工作流 TC 生成可编译脚本
  - **Phase 3（0.5 周）**：eval-test-cases 评分维度更新 + 集成测试 + 向后兼容验证。交付物：eval rubric 工作流评分维度、现有 81 个 TC 零回归。验收：`just e2e-compile` 通过，现有单步 TC 输出与改动前 diff 为空
- 总计 10-14 个 task，5 周日历时间

### Dependency Readiness

- profile 系统和 Fact Table 机制已稳定运行（最近 3 个月 profile manifest schema 零 breaking change，Fact Table 字段提取在现有 20+ e2e 测试上连续 5 次 CI 运行通过率 100%）
- 两个参考项目提供真实测试场景和验证环境
- 不依赖外部服务或 API

## Impact Analysis

### Tier 1: 必须修改（~19 文件）

| 组件 | 文件 | 变更 |
|------|------|------|
| gen-test-cases | SKILL.md | 新增工作流识别、Type 分类、Setup 生成 |
| gen-test-cases | templates/test-cases.md | 模板新增 Type/Workflow/Setup/Steps 结构 |
| gen-test-scripts | SKILL.md | Call-chain + Interactive 双模型处理、Setup→beforeAll 生成 |
| go-test | generate.md | Call-chain(CLI) + Interactive(TUI) 步骤规则 |
| go-test | templates/test-file.go | 工作流测试函数骨架 |
| go-test | templates/helpers.go | 新增 TUI 交互 helper |
| web-playwright | generate.md | Call-chain(API) + Interactive(Web-UI) 步骤规则 |
| web-playwright | templates/helpers.ts | 新增工作流执行 helper |
| web-playwright | templates/api.spec.ts | Call-chain 工作流模板 |
| web-playwright | templates/playwright-ui.spec.ts | Interactive 工作流模板 |
| pytest | generate.md | Call-chain(API/CLI) 步骤规则 |
| pytest | templates/test_file.py | 工作流测试函数骨架 |
| pytest | templates/conftest.py | Setup fixture 生成 |
| java-junit | generate.md | Call-chain + Interactive(TUI) 步骤规则 |
| java-junit | templates/TestFile.java | 工作流测试函数骨架 |
| rust-test | generate.md | Call-chain(CLI) + Interactive(TUI) 步骤规则 |
| rust-test | templates/test_file.rs | 工作流测试函数骨架 |
| rust-test | templates/helpers.rs | TUI 交互 helper |
| maestro | generate.md | Interactive(Mobile-UI) 步骤规则 |

### Tier 2: 应该修改（1 文件）

| 文件 | 变更 |
|------|------|
| eval/rubrics/test-cases.md | 新增工作流 TC 评分维度（步骤完整性、Setup 正确性、状态传递有效性） |

### 不受影响

| 组件 | 理由 |
|------|------|
| 所有 profile 的 run.md | 运行时不区分工作流 vs 单步 |
| 所有 profile 的 graduate.md | 毕业流程不变 |
| run-e2e-tests / graduate-tests SKILL.md | 工作流 TC 的执行和毕业与单步 TC 一致 |
| profile/detect.go, config.go, embed.go | 检测、嵌入机制不变，无新增配置字段 |

## Scope

### In Scope

- test-cases.md 模板重设计：双模型结构（Type + Setup + Steps）
- gen-test-cases SKILL.md 重设计：工作流识别 + Setup 声明生成
- gen-test-scripts SKILL.md 重设计：Call-chain + Interactive 双模型步骤处理 + Setup→beforeAll 生成
- 全部 6 个 profile 的 generate.md 策略更新
- 核心模板更新（go-test + web-playwright Phase 1，其余 Phase 2）
- eval-test-cases 评分维度更新

### Out of Scope

- 现有 test-cases.md 文件迁移（自然兼容，无需迁移）
- graduate-tests / run-e2e-tests 变更（工作流测试的毕业和执行流程不变）
- forge-cli 自身的单元/集成测试变更
- TUI 框架自动检测（generate.md 硬编码框架映射规则，不做框架推断）
- 跨 Type 工作流（如 CLI + Web-UI 组合步骤的单条 TC）。当前设计通过共享 Setup 拆解为多条独立 TC，未来可扩展为 Type 混合模式

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 6 个 profile 的双模型实现引入不一致 | M | H | 定义统一的 workflow 抽象接口规范，每个 profile 只做语言+框架映射 |
| `{stepN.field}` 的 field 名称与实际输出不匹配 | M | M | 优先使用 `--json` 结构化输出；VERIFY marker 充当格式契约；匹配失败时输出 raw stdout |
| TUI helper 是框架特有的，Bubble Tea/ratatui/textual 需要不同实现 | M | M | 由 profile 的 generate.md 定义框架→helper 映射，不建立上层抽象 |
| 工作流中间步骤失败导致误报 | L | H | Call-chain 用 require（fail-fast）；Interactive 按步骤隔离断言；失败输出标识步骤编号 |
| gen-test-cases 的工作流识别质量不足 | M | M | Phase 1 实现 dry-run 模式：gen-test-cases 生成 TC 后先进行结构化 schema 校验（校验 Type 合法性、Steps 非空、`{stepN.field}` 引用的 step 序号在范围内），校验通过后才写入 test-cases.md；Phase 1 用两个参考项目的 PRD 样本作为端到端验证 |

## Success Criteria

- [ ] go-test profile 能从 TUI 工作流 TC 生成包含 2+ 步骤的 Bubble Tea 测试，且步骤间导航状态正确传递（以 agent-forensic 的 session→call tree→diagnosis 流程验证）
- [ ] web-playwright profile 能从 Web-UI 工作流 TC + Setup 声明生成包含 beforeAll API 准备 + 2+ 步骤 UI 交互的 Playwright 测试（以 pm-work-tracker 的 milestone cancel→unbind MI 流程验证）
- [ ] 现有单步 TC 生成的测试代码与改动前完全一致（零回归）
- [ ] gen-test-cases 从 3 个不同 Type（CLI + TUI + Web-UI）的 PRD 样本生成结构合法的多步 TC（每条含 2+ 步骤），经 eval-test-cases 评分均 ≥ 75/100（当前单步 TC 基线平均分）
- [ ] test-cases.md 模板的 Steps 结构对全部 5 种 capability 类型（cli, api, tui, web-ui, mobile-ui）各生成一条可编译的工作流测试，证明每种类型都能正确表达至少 2 步工作流
- [ ] 所有 6 个 profile 的 generate.md 支持双模型步骤生成，通过编译验证
- [ ] 工作流中间步骤解析失败时，生成的测试输出明确的错误标识（格式为 "Step N/M failed: field X not found in output"），且跳过后续步骤（Call-chain 用 fail-fast，Interactive 用步骤隔离）

## Next Steps

- Proceed to `/write-prd` to formalize requirements
