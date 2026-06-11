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

**契约回归无持续机制**：跨 skill 的隐性契约（TC ID 格式、slug 传播、模板分派、CLI 输出格式）修完 N 个第 N+1 个会出现。当前管道只测"功能对不对"，不测"组件之间的约定还在不在"。（参见 `docs/proposals/forge-architecture-simplification/test-strategy-evaluation.md`）

**风险平权**：所有 TC 同一粒度，高严重性缺陷（终态绕过、并发双写）应该有更严格的测试（穷举矩阵、压力测试），但当前模型没有风险等级 → 测试强度的映射。

### Urgency

- forge-cli 快速迭代，跨命令回归风险在增长
- 当前 6 个硬编码 language profile 每新增一个 language 需要修改 forge-cli 代码（每个新 language 约 2 人天：新增 profile 文件 + 适配 generate/run 逻辑 + 测试验证），扩展成本线性增长
- 契约回归问题已在 forge-architecture-simplification 中暴露（11 个契约修复点，累计消耗约 3 人天排查和修复），无持续机制意味着同样的问题会反复出现
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

测试从用户视角出发：**用户做什么 → 期望什么结果**。Journey 的每个步骤都有一个对应的 Contract 验证用户期望是否被系统满足。每个 Journey 的用例数量不硬性限制，由工作流复杂度和上下文窗口自适应：当单个 Journey 的 Contract 数超过 15 或估计 token 数超过 50k 时，自动拆分为多次调用，每次处理一个 Outcome group（happy path 为一批，边缘场景为一批）。

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

**多 Outcome Contract**：每个步骤支持多个 Outcome（成功/错误变体），每个 Outcome 有独立的 Preconditions + Output/State/Side-effect 契约。Outcome 之间通过不同的 Preconditions 互斥——同一时刻只有一个 Outcome 的 Preconditions 成立，不存在组合爆炸：

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

gen-contracts 通过以下机制避免组合爆炸：(1) Outcome 按 Preconditions 互斥分组，同一组内的 Outcome 不会同时为真；(2) 每个 Outcome 独立生成 Contract，不需要交叉验证其他 Outcome；(3) 超过 5 个 Outcome 的步骤自动拆分为多批处理，每批处理一组语义相近的 Outcome。

**Side-effect 契约化**：声明式 Side-effect 契约（hook 名、参数、输出 schema）由 Contract 层定义，Journey 烟测试实际执行验证。无副作用的命令标注 `Side-effect: none`。

**语义描述符**：gen-contracts 阶段使用语义描述而非精确正则（如 "success confirmation containing feature-slug" 而非 `Feature\s+([\w-]+)\s+created`）。具体的正则模式在 gen-test-scripts 阶段结合代码侦察（Fact Table）生成。这避免 LLM 在无法看到代码库时猜测输出格式。转换管道示例：gen-contracts 输出 `Output: "success confirmation containing feature-slug"` → gen-test-scripts 查询 Fact Table 获取该命令的实际 stdout 样本（如 `Feature my-feature created successfully`）→ 生成正则 `Feature\s+([\w-]+)\s+created successfully` + 捕获组声明 `feature_slug: $1`。语义描述符让 gen-contracts 专注于业务意图，正则精确性由 gen-test-scripts 基于真实代码输出保证。

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

- **gen-test-scripts** 生成测试时注入 `@feature` 标签（标签以语言框架原生方式嵌入：Go 使用 `//go:build feature`、Python 使用 `@pytest.mark.feature`、JavaScript 使用 `describe("@feature", ...)`）
- **毕业** = `forge test promote <journey>` 标签更新（`@feature` → `@regression`），不移动文件；promote 命令扫描 `tests/<journey>/` 下所有 `@feature` 标签，替换为 `@regression`
- **CI 选择** = `forge test run --tags regression` 或 `--tags feature`；框架原生标签过滤机制执行（Go `-run`、pytest `-m`、Playwright `--grep`）
- **Quality Gate** = `forge test promote` 执行前自动运行该 Journey 的所有测试，全部通过才更新标签；存在失败测试时拒绝晋升并输出失败报告
- 消除文件迁移、import 重写、去重、分类的复杂度

### 契约断裂检测机制

Contract 规范存储为结构化文档（markdown with schema），包含所有六维度的断言。当以下事件发生时，自动触发契约重新验证：
- CLI 命令的 help text 或输出格式变更
- API endpoint 的 schema 变更
- TUI 组件的 Model 字段变更

重新验证机制：`forge test verify` 命令扫描 Contract 规范，对每个 Output 断言重新执行代码侦察，比较 Fact Table 中的实际值与 Contract 中声明的语义描述是否匹配。

**verify 报告输出示例**：

```
$ forge test verify
Scanning 23 Contracts against Fact Table...

BROKEN (2):
  tests/task-lifecycle/_contracts/step-2-task-claim.md
    Output dimension: expected "claimed task <task_id>" → actual stdout "Task <task_id> claimed"
    Semantic match: partial (word order changed)
  tests/task-lifecycle/_contracts/step-3-task-submit.md
    State dimension: expected "status → completed" → actual "status → done"

OK (21):
  ... (unchanged contracts omitted for brevity)

Summary: 2 broken, 21 OK, 0 false positives
```

**Contract 规范文件示例**：

```
tests/task-lifecycle/
  _contracts/
    step-1-feature-create.md        ← Contract 规范
    step-2-task-claim.md
    step-3-task-submit.md
```

```markdown
# Contract: task-lifecycle / Step 2: forge task claim

## Outcome "success"
- Preconditions: "feature exists with slug matching arg; at least one task available"
- Input: (positional: none; flags: none)
- Output: stdout contains "claimed task <task_id>", exit code 0
- State: tasks/<task_id>/status → "in_progress"; index.json updated
- Side-effect: none

## Outcome "no-tasks-available"
- Preconditions: "feature exists; no tasks available for claiming"
- Input: (positional: none; flags: none)
- Output: stderr contains "no tasks available", exit code 1
- State: unchanged

## Journey Invariants
- feature_slug consistent across all steps
- task_id stable once assigned
```

### TUI 异步 Cmd 处理

TUI 交互步骤的 `await` 语义：
- `await` 表示等待所有 pending Cmd 完成（最长 N 毫秒，默认由 config 声明）
- 超时时 fail-fast，报告超时的 Cmd 名称
- `tea.Batch(cmd1, cmd2)` 中多个并发 Cmd 全部完成后才进入下一步

### Innovation Highlights

**语义描述符延迟精确匹配**：gen-contracts 阶段使用业务语义描述（如 "success confirmation containing feature-slug"），精确的正则匹配延迟到 gen-test-scripts 阶段基于 Fact Table 的代码侦察生成。这解决了 LLM 无法直接观察代码输出格式的根本限制——将"理解意图"和"验证格式"分为两个认知步骤。

**Tag-Based Promotion 替代文件迁移**：测试直接生成到最终目录，通过标签（`@feature` → `@regression`）管理生命周期，消除了传统测试毕业流程中的文件移动、import 重写和去重问题。

**致谢与溯源**：Journey-first 组织借鉴了 Jeff Patton 的用户故事映射（Story Mapping）——用户故事按用户旅程组织而非按功能模块。Contract 验证借鉴了 Pact 的 consumer-driven 思想——消费者（Journey 步骤）定义期望契约。Forge 的创新在于将语义描述符、标签晋升和六维度契约结合为从 PRD 到可验证测试代码的端到端自动化管道。

## Requirements Analysis

### Key Scenarios

- **Journey "开发者完成一个任务"（CLI）[Risk: High]**：用户操作：`forge feature → task claim → task submit`。3 个 Contract 验证每步的输出和状态变更 + Invariants（task_id 全程一致），步骤间传递 feature-slug 和 task-id。Journey 烟测试端到端运行完整链路
- **Journey "分析师诊断会话"（TUI）[Risk: Medium]**：用户操作：`打开会话 → 浏览调用树 → 展开条目 → 打开诊断面板`。4 个 Contract 验证 Model 转换 + View 输出，`key "d" await 3000ms → view contains "诊断摘要"` 处理异步 Cmd
- **Journey "项目经理取消里程碑"（Web-UI + API）[Risk: High]**：用户操作：`导航到里程碑 → 点击卡片 → 更改状态为取消 → 验证关联事项已解绑`。Contract Setup 准备测试数据 + E2E 验证 UI 交互
- **Journey "用户注册到创建资源"（API）[Risk: Medium]**：用户操作：`注册 → 登录 → 创建资源`。3 个 API Contract 验证认证链路 + Invariants（user_id 跨步骤一致）
- **Contract regression**：CLI 输出格式变更时，`forge test verify` 自动检测引用该输出的下游 Contract 是否断裂
- **Edge case**：Contract 中间步骤验证失败，跳过后续步骤并报告失败位置和维度
- **Backward compat**：现有单步 TC 作为单步骤 Journey 退化形式
- **Tag promotion lifecycle**：`gen-test-scripts` 生成测试注入 `@feature` 标签 → CI 运行 `--tags feature` → `forge test promote` 通过后标签更新为 `@regression` → CI 运行 `--tags regression` 作为回归门禁
- **分批生成自动触发**：单个 Journey 的 Contract 数超过 15 或估计 token 数超过 50k 时，gen-contracts 自动拆分为多次调用，每次处理一个 Outcome group（happy path 为一批，边缘场景为一批），拆分后各批结果合并为完整 Contract 文档

### Non-Functional Requirements

- **Backward compatibility**：现有 test-cases.md 和生成的测试脚本零改动继续工作（或提供迁移路径）
- **分批生成**：单次生成一个 Journey（含 happy path + 边缘场景），如果用例数超出上下文窗口则自动分批（同一 Journey 内）
- **项目自适应**：纯 CLI 项目只生成 Contract，无 E2E 层；Web 应用两层都有
- **契约断裂报告**：Contract 验证失败时，输出标识失败维度（Preconditions/Input/Output/State/Side-effect/Invariants）和具体不匹配内容
- **执行性能**：Contract 测试的执行时间不超过当前集成测试的 120%（不含 Setup 时间）
- **Journey 隔离**：每个 Journey 在独立临时目录执行，支持并行运行无干扰
- **认证/授权契约**：Web-UI Journey 涉及认证流程时，Contract 的 Preconditions 必须包含 auth 状态声明（authenticated / unauthenticated / expired token），gen-test-scripts 生成的 Setup 代码注入对应 auth 策略
- **测试数据安全**：gen-test-scripts 生成的测试代码中不硬编码真实 secret/token；Contract 的 Input 维度对敏感字段使用占位符（如 `token: <from-env>`），gen-test-scripts 从环境变量或 .env 文件读取

### Constraints & Dependencies

- 依赖 `.forge/config.yaml` 声明测试框架和执行命令
- 依赖 `docs/conventions/` 中用户提供的测试约定
- 内置模板（从现有 6 个 language profile 演化）作为默认值，项目可覆盖
- 不依赖外部服务或 API
- Contract 的 State 维度验证需要项目暴露状态查询接口（CLI flag、API endpoint、文件读取）；当项目未暴露状态查询接口时，gen-contracts 将 State 维度降级为 Output-only（仅验证 stdout/stderr 中包含状态信息），并在生成的 Contract 中标注 `state-verification: partial`。`partial` 的含义：gen-contracts 在 State 维度仅声明可以从 Output 推断的字段（如 stdout 中包含 "status: completed" → State 包含 status=completed），不声明无法从 Output 推断的字段。gen-contracts 通过代码侦察（Fact Table）自动判断哪些 State 字段有对应的 Output 表达，无法推断的字段标注为 `state-verification: deferred` 并记录到 Contract 的 limitations 部分

## Alternatives & Industry Benchmarking

### Industry Solutions

- **Pact**：Consumer-driven contract testing 框架。核心机制：消费者在测试中声明对提供者的期望（provider states + request/response），Pact 将这些期望序列化为 Pact 文件，提供者在独立验证流程中回放这些请求并校验响应（通过 Pact Broker 管理契约版本）。Forge 与 Pact 的根本差异：(1) Pact 面向微服务间 HTTP 通信，Forge 面向同一进程内 CLI/TUI/API 的 I/O 验证；(2) Pact 的契约是精确的 request/response schema，Forge 使用语义描述符延迟精确匹配；(3) Pact 需要 Broker 基础设施，Forge 的契约断裂通过 `forge test verify` 本地检测。但 Forge 采纳了 Pact 的核心洞察——contract test 与 integration test 互补而非替代（Contract 为主 + Journey 烟测试）
- **Playwright test.step**：多步骤测试，步骤间通过变量共享状态，UI 自动化领域的事实标准
- **Go testing**：社区实践推荐 table-driven subtests 模拟工作流。Bubble Tea 应用通过直接操作 Model 对象测试状态转换
- **Robot Framework**：keyword-driven 模式天然支持工作流，步骤间通过变量传递。Forge 不采用 Robot Framework 的理由：(1) Robot Framework 是独立测试框架，要求项目使用 Robot 的 DSL 编写测试，而 Forge 的用户项目已有各自的测试框架（Go testing、pytest、mocha），引入 Robot 意味着额外学习成本和双框架维护；(2) Robot 的 keyword 抽象层在 LLM 生成场景下增加了不必要的中间层——Forge 直接从 Contract 生成目标框架的原生代码更简单；(3) Robot 的执行报告格式是自有的，无法直接集成到 forge test run 的统一报告。Forge 借鉴了 Robot 的工作流组织思想（步骤间变量传递 → Contract 的 Outcome 间数据流），但将实现保留在项目原生框架内
- **Google Testing**：三层金字塔（Small/Unit 70%, Medium/Integration 20%, Large/E2E 10%）+ 手动探索性测试

### Comparison Table

| Approach | Source | 实现成本 | 维护开销 | 学习曲线 | 工作流对齐 | 契约验证 | Verdict |
|----------|--------|---------|---------|---------|-----------|---------|---------|
| **Do nothing** | — | 零改动 | 随 language 增长线性增长（2 人天/language） | 零学习 | 工作流被拆散到 5 个接口文件 | 无契约验证 | Rejected: 核心问题会随项目增长恶化 |
| Incremental Pyramid Overlay | — | 在现有系统上加 contract 层（约 3 周） | 双系统并行维护，同一工作流验证两次 | 只需学习 contract 概念 | 现有测试仍按接口分组 | 仅新增的 contract 层有验证 | Rejected: 增量叠加不解决根本问题，总成本更高 |
| Keep language×interfaces + add workflow | unified-workflow-test-model proposal | 最小改动（约 1 周） | 与现有系统等同 | 最低（仅新增工作流字段） | 工作流步骤仍被 language×interfaces 拆散 | 无契约验证 | Rejected: XY 问题 |
| **Journey-Driven model** | Pact + Robot Framework + Google Testing | 一次性重写（10 周），后续新增 language 零改动 | 单系统，按 Journey 组织 | 需理解 Journey/Contract 概念（约 0.5 天） | 一个用户工作流 = 一个文件 | 六维度 Contract + forge test verify | **Selected: 从根本上解决测试模型问题** |

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
  - **Phase 1（4 周 + 1 周风险缓冲 = 5 周）**：模型定义 + 目录规范 + 配置系统 + gen-journeys skill + gen-contracts skill（含六维度 + 语义描述符 + Risk 分级 + TUI await 形式化 + Journey 隔离）。Phase 1 含 3 个独立可交付物（模型定义 + gen-journeys + gen-contracts），预留 1 周风险缓冲用于 Dimension 验证和边界 case 处理
  - **Phase 2（3 周）**：gen-test-scripts 重写（配置驱动 + 内置模板迁移 + 语义→regex 转换 + Journey 烟测试生成）
  - **Phase 3（2 周）**：run-tests 重写 + `forge test verify` 契约断裂检测 + forge-cli 命令适配 + eval rubric 更新 + 端到端验证
- 总计 10 周日历时间

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
- TUI await 语义形式化：定义为 Contract 维度的 await 规范（超时阈值、fail-fast 行为、并发 Cmd 等待语义），gen-test-scripts 生成对应的测试代码
- 内置模板迁移：现有 6 个 language profile → 可覆盖的默认模板
- forge-cli `testing` 命令重命名和适配
- 废弃 graduate-tests，替换为 Tag-Based Promotion（标签管理测试生命周期）

### Out of Scope

- 现有项目测试迁移（各项目按自己的节奏迁移）
- Unit 测试（开发者 TDD 手写，Forge 不生成）
- TUI 框架自动检测（框架映射由 config 或 convention 声明）
- 跨 Journey 测试依赖（每个 Journey 自包含）
- Static analysis 层（compile-time 防御由项目自己实现）
- Idempotency 契约（作为 Side-effect 的子维度。推迟理由：Idempotency 主要影响 CI 重试策略——决定失败的测试能否安全重试——而非测试生成的正确性。Phase 1 的 Contract 已覆盖 Side-effect 声明（hook 名、参数），Idempotency 标注（`idempotent: true/false`）可在 Phase 2 作为 Side-effect 维度的扩展字段加入，不影响 Phase 1 的 Contract 结构设计）
- Error propagation 契约（由 Journey 烟测试覆盖；Error propagation 由 Journey 烟测试的端到端执行覆盖，Phase 1 不引入新机制）

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Contract 六维度无法覆盖所有测试场景 | M | H | 六维度覆盖已知场景；gen-contracts 对无法归入已有维度的验证点自动归入 Invariants 维度并标注 `dimension: unclassified`；遇到新维度时扩展 Invariants schema 而非重新设计 |
| 内置模板迁移引入回归 | L | M | 逐个 language 迁移，每个 language 先通过现有 e2e 测试回归验证；现有 e2e 测试提供回归安全网 |
| 配置驱动系统过于复杂 | L | M | 内置模板作为默认值，简单项目零配置即可工作 |
| 4 步 pipeline 增加调用开销 | L | L | 每步专注单一任务，单步复杂度低于当前合并步骤，输出质量提升抵消调用成本；分批生成已在控制上下文 |
| Journey 烟测试不稳定（flaky） | M | M | 烟测试只跑 happy path；烟测试失败时自动重试 3 次（间隔 5s）；3 次均失败后标记为 @flaky 并生成修复建议，不阻塞 CI 但记录到报告中 |
| 语义描述符 → regex 转换不准确 | M | H | `forge test verify` 检测断裂；gen-test-scripts 使用 Fact Table 验证转换结果。verify 自身的准确性通过 bootstrap 策略保证：Phase 1 结束时，对 forge-cli 现有 126+ 测试用例的 stdout 生成 Fact Table 快照，verify 的初始验证集 = 这 126+ 已知正确的输出→Contract 映射，确保 verify 本身不从零开始 |
| Phase 1 范围过大（3 个独立可交付物 × 5 周含缓冲） | M | M | 3 个可交付物按依赖链排序：模型定义（1 周）→ gen-journeys（1.5 周）→ gen-contracts（2.5 周），模型定义先行解耦；如超期，gen-contracts 的 TUI await 形式化可降级为 CLI-only 首次交付，TUI 支持随 Phase 2 补齐 |
| 多 Outcome 步骤的 Outcome 数量失控 | L | M | gen-contracts 强制 Outcome 按 Preconditions 互斥（同一时刻只有一个 Outcome 的 Preconditions 为真），消除组合爆炸；超过 5 个 Outcome 的步骤触发 LLM 检查点，要求确认是否可合并语义相近的 Outcome |
| LLM 生成质量不稳定（gen-contracts / gen-journeys prompt 效果波动） | M | H | gen-contracts 的 prompt 使用结构化 few-shot 示例（从现有 126+ 测试用例提取的 golden examples）；每次 gen-contracts 运行后执行 Contract 格式校验（六维度完整性检查），格式不合格时自动重试；收集失败案例迭代改进 few-shot 示例库 |
| Fact Table 快照与代码库演进不同步 | M | M | verify 的 bootstrap 基于 Phase 1 结束时的 126+ 测试快照，随代码演进可能过时；缓解：`forge test verify` 每次运行时自动重新采集 Fact Table（基于当前代码库），不依赖历史快照；bootstrap 快照仅用于 verify 自身的首次正确性验证 |

## Success Criteria

- [ ] gen-journeys 从 PRD 用户故事提取 Journey，输出 Markdown 文档，每个 Journey 包含：名称、风险等级（High/Medium/Low）、至少 1 个 happy path 步骤、至少 1 个边缘场景步骤。文档格式符合 gen-contracts 可直接解析的结构（Journey 名称 + Step 序列 + 每步的用户操作和期望结果）
- [ ] gen-contracts 从 Journey + Fact Table 推导 Contract，每个 Outcome 包含 4 个必选维度（Preconditions/Input/Output/State），可省略 Side-effect 和步骤级 Invariants。所有生成的 Contract 通过六维度完整性校验（必选维度非空，语义描述符不含 regex 语法）
- [ ] gen-contracts 为每个 Journey 生成 Journey 级别 Invariants（跨步骤不变量声明，至少 1 条），在测试执行中验证（如 task_id 全程不变、feature_slug 一致性）
- [ ] gen-test-scripts 生成的测试带 `@feature` 标签，直接放入 `tests/<journey>/` 目录（无 staging 中间目录）
- [ ] `forge test promote <journey>` 将该 Journey 下所有 `@feature` 标签替换为 `@regression`，替换前后用 `git diff` 确认仅标签变更、无其他代码改动
- [ ] gen-test-scripts 将语义描述符转换为精确正则（基于 Fact Table），生成的测试代码通过编译（`go test -c` / `pytest --collect-only` / `tsc --noEmit`），且至少包含 1 个 Journey 烟测试
- [ ] forge-cli 的 Journey "开发者完成一个任务" 端到端验证：gen-journeys → gen-contracts → gen-test-scripts → run-tests 4 步依次执行成功（每步 exit code = 0），run-tests 报告所有测试通过（0 failed）
- [ ] 每个 Journey 在独立临时目录中执行（路径包含 journey 名称和随机后缀），执行完成后临时目录被清理；3 个 Journey 并行运行时结果与顺序运行完全一致（文件系统状态、退出码、输出内容 diff 为空）
- [ ] `forge test verify` 检测到 CLI 输出格式变更后，报告受影响的 Contract（列出 Contract 文件路径和断裂维度）；对 20+ 个未变更的 Contract 运行 verify，误报数 = 0
- [ ] gen-journeys 输出按 Journey 分文件（一个用户工作流一个文件），而非按接口类型分文件
- [ ] 现有单步 TC 作为单步骤 Journey 退化形式继续工作；所有现有 e2e 测试通过（`just e2e-compile` 编译通过 + 现有 126+ 测试用例输出与改动前 diff 为空）
- [ ] Contract 验证失败时，输出包含：失败维度名称、Contract 文件路径、期望值、实际值（如 "Output dimension FAILED: step-2-task-claim.md, expected 'exit code 0', got 'exit code 1'"）
- [ ] `.forge/config.yaml` 声明 `mocha` → 生成的测试文件包含 `describe/it` 结构；声明 `pytest` → 包含 `def test_` 函数；零配置时输出与 `template: default` 声明完全一致（diff 为空）
- [ ] TUI await 语义：异步 Cmd 等待超时时 fail-fast 并报告超时 Cmd 名称；`tea.Batch(cmd1, cmd2)` 等待全部完成后再进入下一步
- [ ] 现有 6 个 language profile → 内置模板：零配置时 gen-test-scripts 输出与现有 profile 输出 diff 为空；config 声明自定义模板路径时，gen-test-scripts 使用自定义模板
- [ ] forge-cli `testing` 子命令重命名为 `test`，所有子命令行为不变（`test detect`、`test interfaces`、`test run` 等功能等价，仅命令前缀变更）
- [ ] 多 Outcome Contract 正确性：生成的多 Outcome Contract 中，每个 Outcome 的 Preconditions 互斥（同一输入不会同时满足两个 Outcome 的 Preconditions），由 gen-contracts 的结构化 prompt 约束保证
- [ ] Contract 测试执行性能： forge-cli 项目的 Journey 测试执行总时间 ≤ 现有集成测试执行时间的 120%（排除 Setup 阶段，测量从第一个测试开始到最后一个测试结束的 wall time）
- [ ] gen-journeys 的 Risk 分级输出：每个 Journey 标注 High/Medium/Low 之一，高风险 Journey 的边缘场景数量 ≥ happy path 步骤数量
- [ ] 认证/授权契约：Web-UI Journey 的 Contract Preconditions 包含 auth 状态声明（authenticated/unauthenticated/expired token 之一），gen-test-scripts 生成的 Setup 代码包含对应的 auth 策略代码（如 cookie 注入或 token header 设置）
- [ ] 分批生成正确性：当单个 Journey 的 Contract 数 > 15 时，gen-contracts 自动拆分为多批调用，合并后的 Contract 文档与单批生成的 Contract 在六维度内容上完全一致（diff 为空）
- [ ] Journey 烟测试行为正确性：烟测试端到端运行 Journey 的 happy path，每步输出与 Contract 中 "success" Outcome 声明的 Output/State 完全匹配

## Next Steps

- Proceed to `/write-prd` to formalize requirements
