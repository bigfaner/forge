---
iteration: 1
evaluator: QA Engineer (adversarial)
target_docs:
  - docs/features/surface-aware-justfile/prd/prd-spec.md
  - docs/features/surface-aware-justfile/prd/prd-user-stories.md
total_score: 720
---

# PRD Evaluation Report — Iteration 1

## Overall Score: 720 / 1000

---

## Dimension 1: Background & Goals — 85 / 100

### Background has three elements (Reason/Target/Users) — 28 / 30

三元素（Why/What/Who）均存在且具体：

- **Reason**: init-justfile 仅根据语言生成配方，忽略了 surface 类型差异。test.execution 委托层形成冗余间接层。
- **Target**: 4 项交付物明确列出（surface 感知层、编排简化、surface-key 迁移、Task 模型扩展）。
- **Users**: Forge 插件开发者 + Forge 用户（项目开发者），两类用户角色清晰。

**-2**: Target 列出 4 项交付物但未标注优先级，范围偏宽导致读者难以判断哪些是 MVP 核心。

### Goals are quantified — 25 / 30

Goals 表格包含 5 行量化指标：

| Goal | Metric |
|------|--------|
| 消除 test.execution 冗余 | 编排链路从 4 层降至 2 层 |
| Surface 感知配方生成 | 覆盖 5 种 surface 类型 |
| Surface-key 值域统一 | 7+ 组件迁移 |
| 零回归保证 | diff 输出对比验证 |
| Task 数据模型扩展 | 新增双字段 |

**-5**: "7+" 不精确，应枚举具体组件数量。"diff 输出对比验证" 是验证方法而非量化目标。

### Background and goals are logically consistent — 32 / 40

Goal 1（编排简化）和 Goal 2（surface 感知）与 Background 中的问题直接对应。Goal 3（surface-key 迁移）和 Goal 4（Task 模型扩展）是实现 Goal 2 的前置条件，逻辑链成立。

**-8**:
- "75% 的实际示例已通过 just 命令调用" 作为移除 test.execution 的核心依据，但无数据来源、无抽样方法、无可复现的证据。
- Goal 3 和 Goal 4 属于架构重构性质，与 Background 中用户描述的痛点（"配置 surfaces 后自动获得 surface 感知的配方"）存在间接性——用户不关心 surface-key 是枚举还是自定义，只关心功能正确。

---

## Dimension 2: Flow Diagrams — 128 / 150

### Mermaid diagram exists — 45 / 50

存在一个 Mermaid flowchart，覆盖 init-justfile 和 run-tests 两个流程。

**-5**: 两个独立流程合并在一个图中。虽然用 Start/RunStart 区分入口，但图的可读性降低。建议分为两个独立子图。

### Main path complete (start → end) — 48 / 50

两个流程的主路径均从开始到结束完整覆盖：

- init-justfile: Start → 检测语言 → 加载模板 → 生成基础配方 → 检测 surface → 加载规则 → 生成 surface 配方 → 冲突仲裁 → user-customized 保护 → 组装 → 跨平台验证 → End
- run-tests: RunStart → 获取 surface → 加载规则 → 三种编排分支 → 各自到 End

**-2**: run-tests 的 cli/tui 分支（build → dev → test → End）缺少 test 后的 teardown 步骤，与编排表一致但与 web/api/mobile 不一致，未解释原因。

### Decision points + error branches covered — 35 / 50

决策点：CheckSurf、Arbitrate、Protect、ExecSeq、Probe exit、Test exit — 6 个，数量充足。

错误分支：Probe 失败 → teardown → 中止；Test 失败 → teardown。

**-15**:
1. 缺少 dev 启动失败（如端口占用、权限不足导致 dev 直接退出）的错误分支。图中 Probe 只处理了 "dev 已启动但服务未就绪" 的场景。
2. 缺少未知/不支持 surface 类型的错误处理分支。CheckSurf 只有 "有/无" 两个出口，没有 "有但无效" 的出口。
3. 缺少 surface 规则文件加载失败的分支。

---

## Dimension 3: Flow Completeness — 115 / 200

### Flow steps describe complete business process — 50 / 70

init-justfile 和 run-tests 两个核心流程的步骤描述完整，状态转换清晰。

**-20**: Goal 3（surface-key 统一迁移）和 Goal 4（Task 数据模型扩展）在 Flow Description 中完全没有对应的流程描述。这两个目标是 In Scope 的核心交付物，涉及 prompt.go 重写、Go struct 变更、模板迁移等多个步骤，读者无法从文档了解这些变更的实施流程。

### Data flow documented — 30 / 70

这是一个跨多系统的特性（init-justfile skill、run-tests skill、Go 代码、CLI 命令、prompt 模板），但没有独立的数据流文档。

关键未文档化的数据流：
1. Surface 信息从 config.yaml → 任务 frontmatter → run-tests skill 的传递链
2. Surface-key 从 `forge surfaces` CLI → breakdown-tasks/quick-tasks → index.json 的流经路径
3. Task struct 的 Scope → SurfaceKey 字段迁移期间的数据兼容性保证

Functional Specs 的 Related Changes 表格部分涉及数据流，但无法替代显式的数据流文档。

**-40**: 对于涉及 7+ 组件、11 个变更点的跨系统特性，数据流文档的缺失是重大缺陷。

### Exception handling and edge cases covered — 35 / 60

Other Notes → Reliability 部分覆盖了部分异常：
- Dev server 崩溃 → probe 超时 → teardown
- 会话中断 → `.forge/test-state.json` 恢复
- Teardown 幂等（PID 不存在时跳过）

**-25**: 以下错误路径未覆盖：
1. config.yaml 中 surfaces 格式错误或包含未识别的 surface 类型
2. Surface 规则文件缺失或内容损坏
3. `forge surfaces` CLI 命令执行失败时的降级策略
4. 多 surface 项目中某个 surface 的编排失败是否影响其他 surface
5. init-justfile 在检测到 user-customized 配方时的具体提示内容和格式

---

## Dimension 4: User Stories — 166 / 200

### Coverage: one story per target user — 48 / 50

- Forge 插件开发者 → Story 3、Story 4
- Forge 用户（项目开发者）→ Story 1、Story 2

**-2**: 缺少一个 "系统管理员/DevOps" 角色的故事。Background 只定义了两类用户，但如果 surface-key 迁移影响了 CI/CD 配置，应有对应的用户故事。

### Format correct (As a / I want / So that) — 48 / 50

4 个故事全部使用标准格式。动作具体，不是模糊的 "manage" 或 "handle"。

**-2**: Story 3 的 "I want to 将 surface-key 值域...迁移为用户自定义 surface-key" 动作偏向技术实现而非用户价值。

### AC per story (Given/When/Then) — 40 / 50

Story 1-3 使用显式 Given/When/Then 格式。Story 4 的 AC 缺少显式格式标签。

**-10**: Story 4 的 AC 以无序列表形式呈现：
> "And 任务包含 `surface-key: "admin-panel"` 和 `surface-type: "web"`"
> "And 旧任务的 scope 字段通过 GetSurfaceKey() 兼容访问"

缺少显式的 "Given/When/Then" 关键词，不符合标准格式。且第 4 条 AC "And quality-gate fix-task 从失败文件路径自动推断 surface-key/type" 是一个独立的验证场景，应拆分为单独的 AC。

### AC verifiability & boundary coverage — 30 / 50

可验证性分析：
- Story 1: 可验证。但 "生成的配方包含 [linux]/[windows] 双平台变体" 缺少验证方法——是人工检查还是自动化断言？
- Story 2: 可验证。但 exit code 2/3 未在文档中定义含义。
- Story 3: "无 surfaces 配置的项目行为不变" 可验证但验证方法未具体化。
- Story 4: "旧任务的 scope 字段通过 GetSurfaceKey() 兼容访问" 测试的是实现细节（Go 函数名），而非行为。

**-20**:
1. Probe exit code 2/3 语义未定义，无法构造对应的测试用例。
2. Story 1 缺少无 surface 配置场景的 AC（仅在 Story 3 间接覆盖）。
3. Story 2 缺少 probe 超时数值边界（多少秒？多少次重试？）。
4. Story 4 缺少 surface-key/type 同时缺失场景的 AC。

---

## Dimension 5: Scenario Completeness — 90 / 150

### End-to-end scenario coverage — 40 / 60

核心场景覆盖：
- init-justfile surface 感知生成（Story 1 + Flow Description）
- run-tests 编排执行（Story 2 + Flow Description）
- Surface-key 迁移（Story 3）
- Task 模型扩展（Story 4）

**-20**: 以下端到端场景缺失：
1. 混合项目（多 surface，如 admin-panel:web + payment-service:api）的 init-justfile 生成——Scope 提到但无场景。
2. quality-gate fix-task 的端到端流程——仅在 Functional Specs 提及。
3. 从旧系统迁移（已有项目从固定枚举升级到自定义 surface-key）的升级路径。

### Implicit assumptions surfaced — 18 / 40

**-22**: 以下隐含假设未明确浮现：
1. `forge surfaces` CLI 命令已存在且可用——文档多处引用但未声明为前置条件。
2. config.yaml 的 `surfaces` 字段格式和 schema——文档提到使用但未展示具体格式。
3. `.forge/test-state.json` 的 schema——仅在 Reliability 中提到文件名，未定义结构。
4. Probe 重试次数 30 和超时时间——Observability 中提到 "[retry 3/30]" 但未说明 30 是默认值还是可配置。
5. just >= 1.4.0 的版本检查失败时的行为——Performance Requirements 中提到但未说明失败处理。

### Business-rules consistency — 32 / 50

**-18**:
1. BIZ-error-reporting-001 定义 exit 1=retryable, exit 2=blocking。文档中 probe 使用 exit 1/2/3，但 exit 3 在业务规则中无定义。且 Story 2 AC 说 "probe 失败时执行 teardown 后中止"，这与 exit 1=retryable 的语义矛盾——如果是 retryable，为什么 HARD-GATE 禁止重试？
2. BIZ-error-reporting-002 要求每个错误消息包含 failure reason + recovery hint。文档的 Observability 部分只展示了状态输出格式，未验证错误消息是否符合此规则。

---

## Dimension 6: Edge Case Coverage — 52 / 100

### Error paths documented — 22 / 40

已记录：
- Probe 失败 → teardown → 中止
- Test 失败 → teardown
- Dev server 崩溃 → probe 超时 → teardown

**-18**: 缺失的错误路径：
1. Surface 检测失败（config.yaml surfaces 格式错误、包含未知类型）
2. Surface 规则文件缺失或内容损坏
3. `forge surfaces` CLI 不存在或执行失败
4. just 版本 < 1.4.0
5. 用户自定义配方的冲突处理（仅提到 "跳过覆盖 + 输出差异摘要"，未说明差异摘要格式）

### Boundary conditions covered — 15 / 35

已覆盖：
- 无 surface 配置的项目（零回归保证）
- Probe 重试上限（隐含 30 次）

**-20**: 缺失的边界条件：
1. 多 surface 项目的执行策略——顺序执行还是并行？某 surface 失败是否阻断后续？
2. 空 surfaces map（`surfaces: {}`）的处理
3. Surface type 大小写敏感性（"Web" vs "web"）
4. Surface 规则与语言模板的冲突优先级——Arbitrate 节点说 "Surface 规则覆盖"，但覆盖范围的定义不明确
5. init-justfile --force-regenerate 覆盖已有 user-customized 配方时的边界行为

### Failure recovery described — 15 / 25

已描述：
- `.forge/test-state.json` 恢复清理
- Teardown 幂等
- git revert 回滚

**-10**:
1. Surface-key 迁移后回滚，旧数据（scope 字段值）的兼容性恢复路径未描述。
2. 多 surface 项目中部分编排失败的恢复策略未描述。
3. init-justfile --force-regenerate 后用户自定义配方被覆盖的恢复手段未描述。

---

## Dimension 7: Scope Clarity — 84 / 100

### In-scope items are concrete deliverables — 30 / 35

In Scope 列出了具体交付物：文件路径、函数名、字段名、模板数量。

**-5**: "死代码清理：extractTestTypeArg()、genScriptBases" 缺乏明确度——是直接删除？还是需要替换调用方？删除后是否有回归风险？

### Out-of-scope explicitly lists deferred items — 24 / 30

Out of Scope 列出 6 项具体条目。

**-6**: "变更 `forge-cli/internal/cmd/quality_gate.go` 或 `testrunner` 的 Go 代码" 在 Out of Scope 中，但 In Scope 中 "quality-gate fix-task：从失败文件路径推断 surface-key/type" 涉及 quality-gate 的行为变更。如果变更仅在 skill 文档层面（不涉及 Go 代码），应明确说明；如果 fix-task 是 skill 层面的逻辑，应区分清楚。

### Scope consistent with functional specs and user stories — 30 / 35

Functional Specs 11 个变更点与 In Scope 条目基本一一对应。

**-5**: User Story 3 的 AC "无 surfaces 配置的项目行为不变" 需要验证手段，但 Scope 中没有对应的验证工具或测试基础设施变更。

---

## Phase 4: Blindspot Attacks

### [blindspot] 1: `forge surfaces` CLI 命令是关键前置依赖但未作为前置条件声明

> "prompt.go resolveScope() 完全重写为 surfaces map 集合查询"
> "surface-key-assignment 规则文件路径分类改为 CLI 动态查询"
> "forge surfaces <path> longest-prefix-match"

文档多处引用 `forge surfaces` CLI 命令，但从未说明此命令是否已存在。如果需要新建，这是一个显著的额外工作量，应纳入 Scope。如果已存在，应作为前置条件明确声明。当前状态导致下游 agent 无法确定依赖关系。

**Must improve**: 在 Scope 或 Background 中明确声明 `forge surfaces` CLI 的前置条件状态（已存在/需新建/需扩展）。

### [blindspot] 2: 混合项目的 init-justfile 生成策略缺少端到端描述

> "混合项目 dev 配方接受 surface-key 参数"

In Scope 提到混合项目 dev 配方接受 surface-key 参数，但 Flow Description 和 Flow Diagram 都只展示了单 surface 的生成路径。关键问题未解答：
- 一个 justfile 中如何区分不同 surface 的 dev/test 配方？
- 配方命名规则是什么？（`dev-admin-panel` vs `dev surface-key=admin-panel`？）
- 多 surface 的 probe 端口冲突如何处理？

**Must improve**: 添加混合项目的 init-justfile 生成流程描述，或在编排表中为混合项目增加一行。

### [blindspot] 3: cli/tui 缺少 teardown 的设计决策未解释

> 编排表: "cli: build → dev → test | 无服务启动，无需 probe"

web/api/mobile 都有 teardown 步骤，但 cli/tui 没有。图中 Test2（cli/tui）直接到 RunEnd。虽然 cli/tui 无需启动服务所以 "无需 probe" 是合理的，但 dev 步骤是否可能创建需要清理的资源（临时文件、构建缓存锁等）？这个设计决策未被解释。

**Must improve**: 在编排表或 Functional Specs 中明确说明 cli/tui 不需要 teardown 的理由，或将该决策记录为显式选择。

### [blindspot] 4: Story 2 AC 中的 HARD-GATE 语义未定义

> "probe 失败后禁止重试 probe 或重试 dev（HARD-GATE）"

HARD-GATE 不是通用术语，文档中未定义其含义。这是一个错误类型？一个状态标记？一个配置选项？且其语义（"禁止重试"）与 BIZ-error-reporting-001 中 exit 1=retryable 的定义矛盾——如果 probe 以 exit 1 退出（retryable），为什么 HARD-GATE 禁止重试？

**Must improve**: 定义 HARD-GATE 的精确语义，并与 exit code 体系对齐。

### [blindspot] 5: probe exit code 3 在业务规则中无定义

> "probe 失败时（exit 1/2/3）执行 teardown 后中止"

BIZ-error-reporting-001 定义了 exit 0/1/2 的语义，但 exit code 3 在任何业务规则中都未定义。probe 使用 exit 3 表示什么？超时？连接被拒绝？协议错误？这个值域扩展需要被显式定义和注册。

**Must improve**: 在文档中定义 probe 的 exit code 值域及每个 code 的精确语义，或复用已定义的 exit code 体系。

---

## Summary

| Dimension | Score | Max |
|-----------|-------|-----|
| Background & Goals | 85 | 100 |
| Flow Diagrams | 128 | 150 |
| Flow Completeness | 115 | 200 |
| User Stories | 166 | 200 |
| Scenario Completeness | 90 | 150 |
| Edge Case Coverage | 52 | 100 |
| Scope Clarity | 84 | 100 |
| **Total** | **720** | **1000** |

### Top 5 Issues to Address

1. **数据流文档缺失** (Dimension 3, -40 pts): 涉及 7+ 组件的跨系统特性无显式数据流文档，下游 agent 无法理解信息传递链。
2. **Goal 3/4 缺少流程描述** (Dimension 3, -20 pts): surface-key 迁移和 Task 模型扩展是 In Scope 核心交付物，但无实施流程。
3. **错误路径覆盖不足** (Dimension 6, -18 pts): 关键错误路径（surface 检测失败、规则文件缺失、CLI 失败）未记录。
4. **隐含假设未浮现** (Dimension 5, -22 pts): `forge surfaces` CLI、config.yaml surfaces schema、probe 超时参数等关键假设未声明。
5. **AC 边界条件缺失** (Dimension 4, -20 pts): probe exit code 语义未定义，超时数值边界缺失，Story 4 格式不标准。
