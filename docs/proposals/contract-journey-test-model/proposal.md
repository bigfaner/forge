---
created: 2026-05-17
updated: 2026-05-17
author: "faner + Claude"
status: Draft
replaces: unified-workflow-test-model
reviewers:
  - testing-expert: critical issues identified, revisions applied
  - ai-expert: critical issues identified, revisions applied
---

# Proposal: Journey-Driven Test Model

## Problem

Forge 的测试管道以 "language × interfaces" 为组织轴，硬编码 6 种语言 profile 和 5 种接口类型。测试按技术接口分组（CLI 测试、API 测试、UI 测试），而不是按用户真实工作流分组。

### Evidence

**测试组织不反映用户工作流**：用户不会按"CLI 测试""API 测试""UI 测试"思考——他们按"我完成一个任务""我取消一个里程碑"等真实工作流思考。但当前 gen-test-cases 的输出按接口类型分文件（`cli-test-cases.md`、`api-test-cases.md`），强制把用户工作流拆散到不同技术分类中。

**语言枚举不完**：6 个 language profile（go, javascript, python, java, rust, mobile）是硬编码的，无法覆盖所有场景。语言角色也不同——Java 纯后端、JavaScript 纯前端——language 和 interface 不是正交的两个维度。

**测试层级混淆**：当前模型把 CLI、API、TUI、Web-UI、Mobile-UI 平铺为 5 种 interface types，但实际上只有 Web-UI 和 Mobile-UI 能做真正的端到端测试。CLI/API/TUI 的"集成测试"本质是组件间契约验证，与 E2E 是不同性质。

**契约回归无持续机制**：跨 skill 的隐性契约（TC ID 格式、slug 传播、模板分派、CLI 输出格式）修完 N 个第 N+1 个会出现。当前管道只测"功能对不对"，不测"组件之间的约定还在不在"。（参见 `docs/proposals/task-pipeline-hardening/test-strategy-evaluation.md` 第 4 节）

**风险平权**：所有 TC 同一粒度，高严重性缺陷（终态绕过、并发双写）应该有更严格的测试（穷举矩阵、压力测试），但当前模型没有风险等级 → 测试强度的映射。

### Urgency

- forge-cli 快速迭代，跨命令回归风险在增长
- 当前 6 个硬编码 language profile 每新增一个 language 需要修改 forge-cli 代码，扩展成本线性增长
- 契约回归问题已在 task-pipeline-hardening 中暴露（11 个契约修复点），无持续机制意味着同样的问题会反复出现
- 测试策略评估文档指出的敏捷测试象限 Q1（技术+支持）和 Q4（技术+评价）完全空白

## Proposed Solution

将测试管道从 "language × interfaces" 重新设计为 **Journey-Driven 模型**：以用户真实工作流为组织中心，Contract 为每步验证手段。

### Journey = 用户真实工作流

一个 Journey 描述用户完成某个目标的真实操作序列，**包含 happy path 和所有边缘/错误变体**：

```
Journey "开发者完成一个任务" [Risk: High]
  Happy Path:
    Step 1: forge feature my-feature         → 期望：功能创建成功
    Step 2: forge task claim                 → 期望：领取到一个任务
    Step 3: forge task submit --result success → 期望：任务状态变为已完成

  Edge Cases:
    Step 1b: forge task submit (未 claim)    → 期望：拒绝提交，提示无任务
    Step 2b: forge task claim (无可用任务)   → 期望：提示无可领取任务
    Step 3b: forge task submit (重复提交)    → 期望：拒绝重复提交
```

测试从用户视角出发：**用户做什么 → 期望什么结果**。Journey 的每个步骤都有一个对应的 Contract 验证用户期望是否被系统满足。每个 Journey 的用例数量不硬性限制，由工作流复杂度和上下文窗口自适应。

**Risk 分级**：每个 Journey 标注风险等级（High/Medium/Low），由 gen-journeys 根据 PRD 的严重性和失败影响推断。高风险 Journey（涉及状态突变、数据丢失风险）应有更密集的边缘案例覆盖，低风险 Journey（只读操作）可以适当精简。

**Journey 隔离**：每个 Journey 在独立的临时工作目录中执行，避免并行运行时互相干扰。Journey 前置条件在 Setup 阶段声明并准备。

### Contract = 验证手段

Contract 是 Journey 步骤的验证机制，从系统视角定义"这一步的 I/O/状态/副作用/不变量是否符合预期"：

**Contract 六维度**（所有维度按 Outcome 级别声明，Invariants 额外在 Journey 级别声明）：

**必选维度**（每个 Outcome 声明）：
| 维度 | CLI | API | TUI |
|------|-----|-----|-----|
| Preconditions | 执行前必须满足的状态（如 "task status is in_progress"） | 请求前必须满足的条件（如 "user is authenticated"） | 交互前必须满足的 Model 状态 |
| Input | 命令参数 + 工作目录状态 | Request schema + auth | Model 状态 + 事件 |

**必选维度**（每个 Outcome 声明）：
| 维度 | CLI | API | TUI |
|------|-----|-----|-----|
| Output | stdout/stderr 语义描述 + exit code | Response schema + status code | Model 新状态 + View 输出 |
| State | 文件系统变更（JSON Schema） | 数据库/API 状态变更 | 内存 Model 字段 |

**可选维度**（gen-contracts 可省略，省略 = 无该项约束）：
| 维度 | CLI | API | TUI |
|------|-----|-----|-----|
| Invariants（步骤级） | 步骤内不变量（如 "index.json 保持合法 JSON"） | 步骤内不变量 | 步骤内不变量 |
| Side-effect | hook 触发、外部调用 | 网络副作用、消息发送 | 异步 Cmd 执行 |

**Journey 级别 Invariants**（必选）：跨步骤不变量在 Journey 级别声明，如 "task_id 全程不变"、"feature_slug 一致性"。

**多 Outcome Contract**：每个步骤支持多个 Outcome（成功/错误变体），每个 Outcome 有独立的 Preconditions + Output/State/Side-effect 契约：

```
Step: forge task submit
  Outcome "success":
    Preconditions: "task status is in_progress"
    Input: task_id, --result success
    Output: stdout "submitted successfully", exit code 0
    State: status → completed
    Side-effect: hook "post-submit" triggered with { task_id, result }
  Outcome "not-in-progress":
    Preconditions: "no task claimed"           ← 不同的前置条件
    Input: --result success
    Output: stderr "no task claimed", exit code 1
    State: unchanged
```

**Side-effect 契约化**：声明式 Side-effect 契约（hook 名、参数、输出 schema）由 Contract 层定义，Journey 烟测试实际执行验证。无副作用的命令标注 `Side-effect: none`。

**语义描述符**：gen-contracts 阶段使用语义描述而非精确正则（如 "success confirmation containing feature-slug" 而非 `Feature\s+([\w-]+)\s+created`）。具体的正则模式在 gen-test-scripts 阶段结合代码侦察（Fact Table）生成。这避免 LLM 在无法看到代码库时猜测输出格式。

**Contract 为主 + Journey 烟测试**：每个 Journey 的 Contract 覆盖 90% 的验证场景。此外，每个 Journey 保留一个端到端实际执行的烟测试，覆盖运行时交互正确性（环境配置、进程间通信、时序行为）。烟测试由 gen-test-scripts 自动生成，复用 Journey 的 happy path 步骤。

### Pipeline：4 步拆分

```
gen-journeys → gen-contracts → gen-test-scripts → run-tests
     (叙述性)      (技术性)        (代码生成)        (执行)
```

每步专注一个认知任务：

| 步骤 | 输入 | 输出 | 需要代码侦察 |
|------|------|------|-------------|
| gen-journeys | PRD 用户故事 | Journey 叙述文档（用户工作流 + 风险分级） | 否 |
| gen-contracts | Journey 文档 + 代码侦察 | Contract 规范（六维度，语义描述符） | 是 |
| gen-test-scripts | Contract 规范 + 模板 + 代码侦察 | 可执行测试代码 + Journey 烟测试 | 是 |
| run-tests | 测试代码 | 结果报告 | 否 |

分批生成：每次只处理一个 Journey。如果单个 Journey 的用例数超出上下文窗口，自动分批（同一 Journey 内）。

### 测试层级按项目实际接口决定

```
Unit (TDD)                        — 开发者手写，纯函数/状态机
    ↑
Contract (Journey 步骤验证)        — Forge 生成，CLI/API/TUI
    ↑
E2E (Web-UI / Mobile-UI)          — 仅前端项目，完整用户体验
```

- **forge-cli**（纯 CLI）：Unit → Contract，无 E2E
- **agent-forensic**（TUI）：Unit → Contract，无 E2E
- **pm-work-tracker**（Web 应用）：Unit → Contract → E2E

### 配置驱动框架选择

项目通过 `.forge/config.yaml` 声明测试框架和执行命令，Forge 不硬编码语言。内置模板作为便利默认值。

### 测试生命周期：Tag-Based Promotion

测试直接生成到最终目录（`tests/<journey>/`），不再有 staging 区域。通过标签管理生命周期：

```
tests/
  task-lifecycle/                ← Journey 目录（领域，稳定）
    claim_submit_test.go         ← @feature（新生成，验证中）
    task_record_test.go          ← @regression（已验证，回归测试）
```

- **gen-test-scripts** 生成测试时注入 `@feature` 标签
- **毕业** = `forge test promote <journey>` 标签更新（`@feature` → `@regression`），不移动文件
- **CI 选择** = `forge test run --tags regression` 或 `--tags feature`
- **Quality Gate** = 测试必须至少通过一次才能晋升标签
- 消除文件迁移、import 重写、去重、分类的复杂度

### 契约断裂检测机制

Contract 规范存储为结构化文档（markdown with schema），包含所有六维度的断言。当以下事件发生时，自动触发契约重新验证：
- CLI 命令的 help text 或输出格式变更
- API endpoint 的 schema 变更
- TUI 组件的 Model 字段变更

重新验证机制：`forge test verify` 命令扫描 Contract 规范，对每个 Output 断言重新执行代码侦察，比较 Fact Table 中的实际值与 Contract 中声明的语义描述是否匹配。

### TUI 异步 Cmd 处理

TUI 交互步骤的 `await` 语义：
- `await` 表示等待所有 pending Cmd 完成（最长 N 毫秒，默认由 config 声明）
- 超时时 fail-fast，报告超时的 Cmd 名称
- `tea.Batch(cmd1, cmd2)` 中多个并发 Cmd 全部完成后才进入下一步

### Innovation Highlights

**Journey-first，不是 Contract-first**：测试的组织轴是用户工作流，不是技术接口。Journey 从 PRD 的用户故事中提取（用户做什么 → 期望什么），Contract 从 Journey 步骤推导（验证什么）。

**契约 + 烟测试互补**：Contract 验证接口兼容性和状态正确性，Journey 烟测试验证运行时交互正确性。两者互补而非互斥。

**配置驱动替代硬编码**：项目声明框架和执行方式，Forge 负责从 PRD 到测试代码的生成和调度。

**契约断裂持续检测**：Contract 既是测试规范也是可检测的契约，通过 `forge test verify` 实现持续的契约回归防护。

## Requirements Analysis

### Key Scenarios

- **Journey "开发者完成一个任务"（CLI）[Risk: High]**：用户操作：`forge feature → task claim → task submit`。3 个 Contract 验证每步的输出和状态变更 + Invariants（task_id 全程一致），步骤间传递 feature-slug 和 task-id。Journey 烟测试端到端运行完整链路
- **Journey "分析师诊断会话"（TUI）[Risk: Medium]**：用户操作：`打开会话 → 浏览调用树 → 展开条目 → 打开诊断面板`。4 个 Contract 验证 Model 转换 + View 输出，`key "d" await 3000ms → view contains "诊断摘要"` 处理异步 Cmd
- **Journey "项目经理取消里程碑"（Web-UI + API）[Risk: High]**：用户操作：`导航到里程碑 → 点击卡片 → 更改状态为取消 → 验证关联事项已解绑`。Contract Setup 准备测试数据 + E2E 验证 UI 交互
- **Journey "用户注册到创建资源"（API）[Risk: Medium]**：用户操作：`注册 → 登录 → 创建资源`。3 个 API Contract 验证认证链路 + Invariants（user_id 跨步骤一致）
- **Contract regression**：CLI 输出格式变更时，`forge test verify` 自动检测引用该输出的下游 Contract 是否断裂
- **Edge case**：Contract 中间步骤验证失败，跳过后续步骤并报告失败位置和维度
- **Backward compat**：现有单步 TC 作为单步骤 Journey 退化形式

### Non-Functional Requirements

- **Backward compatibility**：现有 test-cases.md 和生成的测试脚本零改动继续工作（或提供迁移路径）
- **分批生成**：单次生成一个 Journey（含 happy path + 边缘场景），如果用例数超出上下文窗口则自动分批（同一 Journey 内）
- **项目自适应**：纯 CLI 项目只生成 Contract，无 E2E 层；Web 应用两层都有
- **契约断裂报告**：Contract 验证失败时，输出标识失败维度（Preconditions/Input/Output/State/Side-effect/Invariants）和具体不匹配内容
- **执行性能**：Contract 测试的执行时间不超过当前集成测试的 120%（不含 Setup 时间）
- **Journey 隔离**：每个 Journey 在独立临时目录执行，支持并行运行无干扰

### Constraints & Dependencies

- 依赖 `.forge/config.yaml` 声明测试框架和执行命令
- 依赖 `docs/conventions/` 中用户提供的测试约定
- 内置模板（从现有 6 个 language profile 演化）作为默认值，项目可覆盖
- 不依赖外部服务或 API
- Contract 的 State 维度验证需要项目暴露状态查询接口（CLI flag、API endpoint、文件读取）

## Alternatives & Industry Benchmarking

### Industry Solutions

- **Pact**：Consumer-driven contract testing。核心思想是消费者定义期望契约、提供者验证契约。但 Pact 明确指出 contract test 应与 integration test 并行运行，不是替代——我们采纳了这一点（Contract 为主 + Journey 烟测试）
- **Playwright test.step**：多步骤测试，步骤间通过变量共享状态，UI 自动化领域的事实标准
- **Go testing**：社区实践推荐 table-driven subtests 模拟工作流。Bubble Tea 应用通过直接操作 Model 对象测试状态转换
- **Robot Framework**：keyword-driven 模式天然支持工作流，步骤间通过变量传递
- **Google Testing**：三层金字塔（Small/Unit 70%, Medium/Integration 20%, Large/E2E 10%）+ 手动探索性测试

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| **Do nothing** | — | 零改动 | 测试不反映用户工作流；语言枚举问题不解决；契约回归无持续机制 | Rejected: 核心问题会随项目增长恶化 |
| Incremental Pyramid Overlay | — | 过渡平滑，风险可控 | 新旧模型并存，技术债累积 | Rejected: 最终还是要大改，总工作量更大 |
| Keep language×interfaces + add workflow | unified-workflow-test-model proposal | 最小改动，只加工作流能力 | 不解决组织维度和契约问题 | Rejected: XY 问题 |
| **Journey-Driven model** | Pact + Robot Framework + Google Testing | 用户工作流驱动；Contract+烟测试互补；配置驱动解耦语言；契约断裂持续检测 | 改动面大 | **Selected: 从根本上解决测试模型问题** |

## Feasibility Assessment

### Technical Feasibility

- **高**。4 步 pipeline 每步有明确的认知任务：
  - gen-journeys：纯叙述性提取，LLM 擅长
  - gen-contracts：结合代码侦察推导六维度，Fact Table 机制已成熟
  - gen-test-scripts：基于 Contract 规范 + 模板生成代码，当前已有基础
  - run-tests：项目声明命令，Forge 调度
- **配置驱动**已有基础：`.forge/config.yaml` 已有 `languages` 和 `interfaces` 字段
- **内置模板迁移**是机械性工作

### Resource & Timeline

- 1 名工程师，分 3 个阶段交付：
  - **Phase 1（4 周）**：模型定义 + 目录规范 + 配置系统 + gen-journeys skill + gen-contracts skill（含六维度 + 语义描述符 + Risk 分级 + TUI await 形式化 + Journey 隔离）
  - **Phase 2（3 周）**：gen-test-scripts 重写（配置驱动 + 内置模板迁移 + 语义→regex 转换 + Journey 烟测试生成）
  - **Phase 3（2 周）**：run-tests 重写 + `forge test verify` 契约断裂检测 + forge-cli 命令适配 + eval rubric 更新 + 端到端验证
- 总计 9 周日历时间

### Dependency Readiness

- `.forge/config.yaml` 配置系统已存在，只需扩展字段
- `docs/conventions/` 目录和加载机制已存在
- 现有 6 个 language profile 的 generate.md/run.md 可作为内置模板基础
- forge-cli 的 `testing` 子命令架构已稳定，改造是增量修改
- Fact Table 机制已稳定（20+ e2e 测试连续通过）

## Scope

### In Scope

- Journey-Driven 测试模型定义：Journey（用户真实工作流 + 风险分级）+ Contract（六维度验证）+ E2E（前端体验验证）
- 目录规范：`tests/<journey>/`（测试直接生成到最终目录，标签管理生命周期）
- 配置驱动框架系统：`.forge/config.yaml` 声明测试框架、执行命令、capabilities
- gen-journeys skill：从 PRD 用户故事提取 Journey（叙述性）+ 风险分级
- gen-contracts skill：从 Journey + 代码侦察推导 Contract（六维度 + 语义描述符 + 多 Outcome）+ Invariants
- gen-test-scripts skill：Contract → 可执行代码 + Journey 烟测试
- run-tests skill：项目声明执行命令，Forge 调度 + 结果报告
- Journey 隔离：独立临时工作目录
- `forge test verify`：契约断裂检测
- TUI await 语义形式化
- 内置模板迁移：现有 6 个 language profile → 可覆盖的默认模板
- forge-cli `testing` 命令重命名和适配
- 废弃 graduate-tests，替换为 Tag-Based Promotion（标签管理测试生命周期）

### Out of Scope

- 现有项目测试迁移（各项目按自己的节奏迁移）
- Unit 测试（开发者 TDD 手写，Forge 不生成）
- TUI 框架自动检测（框架映射由 config 或 convention 声明）
- 跨 Journey 测试依赖（每个 Journey 自包含）
- Static analysis 层（compile-time 防御由项目自己实现）
- Idempotency 契约（作为 Side-effect 的子维度，Phase 2 根据需求细化）
- Error propagation 契约（由 Journey 烟测试覆盖）

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Contract 六维度无法覆盖所有测试场景 | M | H | 六维度覆盖已知场景；遇到新维度时扩展而非重新设计 |
| 内置模板迁移引入回归 | M | M | 逐个 language 迁移，每个 language 先通过现有 e2e 测试回归验证 |
| 配置驱动系统过于复杂 | L | M | 内置模板作为默认值，简单项目零配置即可工作 |
| 4 步 pipeline 增加调用开销 | M | L | 每步专注单一任务，输出质量提升抵消调用成本；分批生成已在控制上下文 |
| Journey 烟测试不稳定（flaky） | M | M | 烟测试只跑 happy path；不稳定时标记为 known-flaky 不阻塞 CI |
| 语义描述符 → regex 转换不准确 | M | H | `forge test verify` 检测断裂；gen-test-scripts 使用 Fact Table 验证转换结果 |

## Success Criteria

- [ ] gen-journeys 从 PRD 用户故事提取 Journey（含 happy path + 边缘场景 + 风险分级），输出叙述性工作流文档
- [ ] gen-contracts 从 Journey + Fact Table 推导 Contract（六维度，语义描述符，多 Outcome），Invariants 覆盖跨步骤不变量
- [ ] gen-test-scripts 生成的测试带 `@feature` 标签，直接放入最终目录（无 staging）
- [ ] `forge test promote <journey>` 将通过的测试标签从 `@feature` 更新为 `@regression`（Tag-Based Promotion）
- [ ] gen-test-scripts 将语义描述符转换为精确正则（基于 Fact Table），生成可编译的测试代码 + Journey 烟测试
- [ ] forge-cli（纯 CLI 项目）的 Journey "开发者完成一个任务" 端到端验证：gen-journeys → gen-contracts → gen-test-scripts → run-tests 全链路通过
- [ ] 每个 Journey 在独立临时目录执行，支持并行运行无干扰
- [ ] `forge test verify` 检测到 CLI 输出格式变更后，报告受影响的 Contract
- [ ] gen-journeys 输出按 Journey 分文件（一个用户工作流一个文件），而非按接口类型分文件
- [ ] 现有单步 TC 作为单步骤 Journey 退化形式继续工作（零回归）
- [ ] Contract 验证失败时，输出标识失败维度和具体不匹配内容
- [ ] `.forge/config.yaml` 声明测试框架后，Forge 能正确加载并使用（零配置时使用内置模板默认值）
- [ ] TUI await 语义：异步 Cmd 等待超时时 fail-fast 并报告超时 Cmd 名称
- [ ] 现有 6 个 language profile 的 generate.md/run.md 作为默认模板工作，项目可通过 config 覆盖
- [ ] forge-cli `testing` 子命令重命名为 `test` 并适配完成

## Next Steps

- Proceed to `/write-prd` to formalize requirements
