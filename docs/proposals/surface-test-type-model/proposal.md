---
created: "2026-05-26"
author: "fanhuifeng"
status: Approved
---

# Proposal: Surface-Specific Test Type Model（Surface 测试类型模型）

## Problem

Forge 把所有生成的测试统称为 "e2e 测试"，但不同 Surface 的测试在语义、执行模型和验证维度上存在本质差异。CLI 测试启动子进程验证输出，API 测试发送 HTTP 请求验证响应，Web 测试驱动浏览器验证交互——它们不是同一类测试，却共享同一个标签。

### Evidence

1. **目录结构**：测试代码统一放入 `tests/<journey>/`（部分项目用 `tests/e2e/`），不区分 surface 类型
2. **Justfile recipe**：使用 `just test` 或 `just test-e2e` 执行所有测试，无 surface 区分
3. **Task 类型名**：`test.gen-scripts`、`test.run` 等任务类型不携带 surface 信息
4. **文档混用**：ARCHITECTURE.md 中 "高级测试" 和 "e2e 测试" 交替使用，概念边界模糊
5. **质量门**：`just unit-test`（开发者自写）和 `just test`（Forge 生成）两层划分，但第二层缺乏测试类型语义
6. **gen-test-scripts 的 types/ 目录**：已有 CLI/API/TUI/Web/Mobile 五种生成策略，但生成的测试代码在命名和标签上不体现这一差异

**实际影响**：Forge 代码库中 `eval-design` skill 在评审 CLI 项目时，因测试标记为 "e2e" 而按完整用户旅程标准评判，导致评分偏低——这是分类标签与评审维度不匹配的直接证据（见 Urgency 段落详述）。此外，API 项目的 `just test-e2e` 实际只发送 HTTP 请求验证响应，不涉及客户端行为，但标签暗示了全链路覆盖。

### Urgency

当前已造成实际影响：Forge 代码库中 `eval-design` skill 在评审 CLI 项目时，因测试标记为 "e2e" 而按完整用户旅程标准评判，导致评分偏低——这是已确认的工具链行为偏差，而非假设场景。新 skill/规则文件编写时，作者需要自行判断 "此处的 e2e 实际指什么"，增加了概念对齐成本。

若推迟 3 个月，Forge 预计新增 2-3 种 surface 类型的项目支持，每种都会引入新的测试类型歧义。届时受影响文件将从当前的 ~15 个扩展到 ~25 个，迁移成本翻倍。术语混乱还会渗透到新增的 skill 规则文件中，增加后续清理的范围。

## Proposed Solution

取消 "e2e 测试" 作为统一标签，建立 **Surface → Test Type** 的映射模型。每种 surface 有自己的测试类型名称、语义定义和验证维度。

### Test Type Mapping

| Surface | Test Type（EN） | 测试类型（CN） | 验证维度 | 执行模型 |
|---------|-----------------|---------------|---------|---------|
| `cli` | CLI Functional Test | CLI 功能测试 | 进程退出码 + stdout 文本 + stderr 文本 | 子进程执行 |
| `tui` | Terminal Functional Test | 终端功能测试 | 终端输出文本 + stdin 交互响应序列 | 子进程 + stdin pipe |
| `api` | API Functional Test | API 功能测试 | HTTP 状态码 + 响应体 JSON + 响应 Header | HTTP 客户端 |
| `web` | Web E2E Test | Web 端到端测试 | DOM 元素可见性 + 用户操作响应 + 页面 URL 变更 + 元素属性值 | 浏览器自动化 |
| `mobile` | Mobile E2E Test | 移动端端到端测试 | UI 元素可见性 + 用户操作响应 + 屏幕 ID 变更 | Maestro YAML / 手动验证 |

### 分类标准声明

本模型的分类维度有两个：

1. **一级分类键：Surface**（cli / tui / api / web / mobile）— 决定测试的执行模型（子进程 / HTTP / 浏览器 / 设备自动化）
2. **二级属性：测试范围**（功能测试 / 端到端测试）— 由 Surface 的验证机制决定

二级属性的判定标准基于**验证机制**而非覆盖的技术栈深度。**功能测试**通过协议级调用（子进程调用、HTTP 请求）验证输入-输出行为，测试工具在协议边界上观测被测系统的响应——CLI 测试观测子进程的退出码和 stdout/stderr，API 测试观测 HTTP 响应的状态码和 body。**端到端测试**通过设备级自动化（浏览器驱动、移动设备自动化）模拟真实用户操作序列，测试工具模拟用户的设备级交互——Web 测试通过 Playwright 模拟浏览器操作，Mobile 测试通过 Maestro 模拟移动设备操作。

关键区分：CLI 测试可以遍历完整技术栈（如 CLI 命令读写数据库后输出结果），API 测试也可以触发从 HTTP 请求到持久层再返回响应的完整调用链。它们在技术栈覆盖上可能是"端到端"的，但验证机制是在协议边界上的单次调用观测，而非通过设备级自动化模拟用户操作流程。本模型中的"功能测试"/"端到端测试"标签反映的是验证机制，不是技术栈覆盖深度。

验证维度在所有 Surface 上统一为"可观测输出属性"粒度：每个维度都是测试可以直接断言的外部可观测值。

### 语义定义

- **CLI 功能测试**：编译独立二进制，通过子进程调用，验证命令行参数解析、输出格式、退出码、错误处理。不测试内部函数，通过进程边界隔离。
- **终端功能测试**：编译独立二进制，通过 stdin pipe 模拟用户输入，验证终端渲染输出（ANSI 序列处理、布局、异步 Cmd 响应）。与 CLI 功能测试的区别在于需要模拟交互输入流。
- **API 功能测试**：启动 HTTP 服务器（或使用测试服务器），发送请求，验证响应符合 Contract 定义的六个维度。注意：此处的"Contract"指 Forge 在 gen-contracts 阶段生成的 API 行为规约，而非 Pact 意义上的消费者-提供者契约测试。
- **Web 端到端测试**：启动 dev server，通过浏览器自动化（Playwright）模拟用户操作，验证 UI 渲染、交互逻辑、跨页面导航和跨组件状态流转。覆盖从用户输入到持久层再回到 UI 的完整用户旅程。
- **移动端端到端测试**：通过 Maestro YAML 定义操作序列，驱动移动端 UI，验证渲染、交互和屏幕导航。与 Web 端到端测试逻辑同构，均通过设备级自动化覆盖完整用户旅程。Best-effort 模式，部分场景标记为 manual-only。

### User-Facing Experience

测试执行时，justfile recipe 输出标签从 `Running e2e tests...` 变为 `Running <surface-key> tests...`。测试报告中的 suite 名称从 `e2e/journey-name` 变为 `<surface-key>/journey-name`。CI dashboard 上每个测试类型的执行结果独立显示，不再混在单一的 "e2e" 分类下。

### Technical Direction

**Task type 解析**：当前 task-lifecycle parser 已支持带点的类型名（如 `eval.design`、`eval.contract`），`test.gen-scripts.cli` 遵循相同的 `{action}.{skill}.{surface}` 三段式解析规则，无需 parser 改动。

**Justfile recipe 命名**：`init-justfile` skill 当前使用 justfile 模板生成 recipe。实现方式：在每个 surface 的 justfile 模板中，recipe 使用 `<surface-key>-test` 命名格式（如 CLI surface 的 recipe 命名为 `cli-test`），不再使用 `test-e2e` 或 `test-<surface-type>-<scope>` 格式。

> 更新说明：本提案原定使用 `test-<surface-type>-<scope>` 命名格式（如 `test-cli-functional`）并保留 `test-e2e` alias 做 2 版本过渡期。后经 `test-pipeline-consistency-audit` 提案审核，v3.0.0 大版本允许直接破坏性变更，不再需要 alias 过渡期。recipe 命名统一为 `<surface-key>-test`（如 `cli-test`），向后兼容 alias 直接删除。

### Why This Needs a Proposal

表面上是"改名"，实际影响范围横跨 skill 文件、task 类型系统、justfile 模板、业务规则文档和概念参考文档，涉及多个 skill 的协作语义。不加提案直接改动会导致：(1) 各 skill 对新术语采用时间不一致，产生术语混用期；(2) task 类型名变更需与 task-lifecycle business rule 的保留类型列表协调，否则 parser 会拒绝新类型名。提案的价值在于统一变更节奏和协调依赖。

## Requirements Analysis

### Key Scenarios

1. **概念查询**：用户查阅文档时，能快速找到自己项目 surface 对应的测试类型定义
2. **测试生成**：gen-test-scripts 生成的测试代码文件名/注释中体现测试类型（而非统一的 "e2e"）
3. **测试执行**：justfile recipe 按测试类型命名（如 `cli-test`、`api-test`），而非统一的 `test-e2e`
4. **任务追踪**：index.json 中 test 任务的类型名携带 surface 信息（如 `test.gen-scripts.cli`）和测试范围（functional / e2e）
5. **质量门**：质量门报告区分不同测试类型的执行结果

### Constraints & Dependencies

- 现有 skill 的 `types/` 和 `rules/surfaces/` 文件已按 surface 分化，是本提案的基础
- `forge surfaces` CLI 已提供 surface 检测能力
- Justfile recipe 命名需与 init-justfile skill 的 surface 规则同步更新
- 任务类型名的变更需与 task-lifecycle business rule 中的保留类型列表协调

### Non-Functional Requirements

1. **向后兼容**：[已覆盖] 旧 justfile recipe 名（`test`、`test-e2e`）在 v3.0.0 大版本中允许直接破坏性变更（`test-pipeline-consistency-audit` 提案覆盖本提案的 NFR1），不再需要 2 版本 alias 过渡期。旧 recipe 名（`test-e2e`）在本提案中直接被 `<surface-key>-test` 替代，无向后兼容层。
2. **迁移性能**：术语更新工作在单次 forge 执行周期内完成，不要求重新生成测试代码或重新运行测试
3. **可发现性**：用户通过 `just --list` 即可看到按测试类型命名的 recipe 列表，无需查阅外部文档即可理解每个 recipe 的用途
4. **测试执行性能零影响**：术语变更不引入新的测试执行步骤或额外的进程启动开销，每个 surface recipe 的执行路径与当前 `test-e2e` recipe 完全一致（注：v3.0.0 中 `test-e2e` 被 `<surface-key>-test` 直接替代）
5. **CI 集成无破坏性变更**：[已覆盖] CI pipeline 中引用旧 recipe 名（`just test-e2e`）在 v3.0.0 中由 `test-pipeline-consistency-audit` 提案直接替代为 `<surface-key>-test`，不再需要过渡期
6. **过渡期结束追踪**：[已覆盖] 由 `test-pipeline-consistency-audit` 提案的 v3.0.0 大版本一次性变更替代，不再需要分版本 alias 移除追踪

### Edge Cases

1. **多 surface 项目**：项目同时包含 CLI 和 API surface 时，justfile 中应包含 `cli-test` 和 `api-test` 两个独立 recipe，同时保留 `test` 作为运行所有测试的聚合 recipe
2. **无 surface 项目**：纯库项目（无 surface）不生成测试类型 recipe，`just test` 仍然映射到 `just unit-test`
3. **过渡期混用**：旧项目已生成的测试代码使用旧术语，新项目使用新术语。两套术语在 CI 中并行运行，不产生冲突——测试代码的注释/标签变更不影响执行结果

## Alternatives & Industry Benchmarking

### Industry Solutions

- **Go 社区 build tags**：Go 通过 `//go:build integration`、`//go:build e2e` 等 build tags 在源码中声明式标记测试类型，`go test -tags=e2e` 按标签选择执行。模式：开发者在每个测试文件的源码中添加编译约束标签，测试框架在编译阶段过滤。不可采用原因：Forge 的测试是代码生成而非手写，build tags 要求开发者在测试代码中做出标记选择，但 Forge 的生成流程在 `gen-test-scripts` 阶段就已确定测试类型——tags 的选择应在生成时决定而非由开发者事后标注。如果采用，需要在生成模板中硬编码 build tags，这与 Forge 已有的 surface 分化机制（`types/` 目录）功能重复。吸收：build tags 的"在代码中嵌入类型标记"思路被 Forge 的 task type 命名方案吸收，但以 `test.gen-scripts.cli` 三段式类型名替代编译标签。
- **Spring Boot @Tag 注解**：JUnit 5 的 `@Tag("integration")`、`@Tag("e2e")` 注解标记测试方法/类，Maven/Gradle Surefire 插件按标签过滤执行。模式：注解声明在测试代码上，构建工具在运行时按标签组合过滤。不可采用原因：Forge 不运行在 JVM 生态，注解语法不适用。更根本的差异是：Spring Boot 的标签是测试代码的元数据属性，而 Forge 的测试类型是生成时的结构性决策（由 surface 类型决定，不由开发者选择）。吸收："用标签而非目录结构区分测试类型"的思路已被 Forge 的 task type 命名方案吸收——task type 名即标签，无需注解或目录约定。
- **Playwright project 配置**：Playwright 通过 `playwright.config.ts` 中的 `projects` 数组定义不同测试套件，每个 project 有独立的测试目录、浏览器配置和测试参数。模式：同一配置文件中声明多个测试项目，每个项目绑定到特定目录和运行参数，`npx playwright test --project=chromium` 选择执行。不可采用原因：Playwright 的 project 维度是运行时配置（同一套测试可配不同浏览器/设备），而 Forge 需要的是语义维度（不同 surface 不同测试类型名称），两者正交。Forge 的 Web surface 测试已使用 Playwright，但 project 配置解决的是"如何运行"而非"如何命名"。吸收：Playwright 的"每个 project 独立配置"思路对应 Forge 中每个 surface 有独立的 `rules/surfaces/` 规则文件。
- **Postman/Newman collection run**：Postman 的测试集合称为 "collection"，执行命令为 `newman run collection.json`，测试报告中使用 collection 名和 folder 名组织结果，不使用 "e2e" 或 "integration" 等抽象标签。模式：按执行对象（collection）而非按测试类型命名，测试结果按请求分组而非按类型分组。不可采用原因：Postman 仅为 API 测试工具，不存在多 surface 分类问题。吸收：Newman 的命名原则——按执行方式命名而非按假设的覆盖范围命名——直接验证了本提案的核心主张：API 测试不应被称为 "e2e"。Forge 的 API 功能测试命名与 Newman 的命名逻辑一致。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零成本 | 术语混乱持续恶化，阻碍后续重构 | Rejected: 问题会随 Forge 支持的 surface 类型增多而加剧 |
| 统一改名 "高级测试" | ARCHITECTURE.md 现有术语 | 改动小 | 仍然是一个笼统概念，不解决类型错配问题 | Rejected: 只是换了一个模糊标签 |
| 引入标准测试分层（unit/integration/e2e） | 行业标准 | 概念通用 | 不适合 Forge 的场景——CLI 测试不是传统意义上的 integration test，API 测试也不是传统意义上的 contract test | Rejected: 行业术语有既定含义，强行复用会产生歧义（如"契约测试"已被 Pact 语义占据） |
| **Surface → Test Type 映射** | Forge 自身实践 | 精确、与已有分化一致、可扩展 | 需要更新多个文件和概念 | **Selected: 最小惊讶原则——名称匹配实际行为** |

### Trade-off Analysis

选择 Surface → Test Type 映射方案的真实代价：

1. **学习曲线**：现有用户需要从单一 "e2e" 概念切换到 5 种测试类型名称。缓解：分类规则是确定性的（按 surface key 查表即可），无主观判断。
2. **迁移成本**：所有引用 "e2e" 的文件需逐一更新，现有项目已生成的 justfile 和测试代码不会自动迁移。缓解：v3.0.0 大版本允许一次性破坏性变更（由 `test-pipeline-consistency-audit` 提案覆盖），不再需要 backward-compatible alias。
3. **"功能测试" 语义冲突**：行业中 "functional test" 泛指验证功能需求的测试（包含单元级），本提案中特指通过进程/HTTP 边界的黑盒测试。缓解：每次使用时都带 surface 前缀（如 "CLI 功能测试"），限定范围。
4. **二分法扩展性**：functional/e2e 的二分法可能在新增 surface 类型时遇到边界模糊的情况（如 desktop surface）。缓解：分类标准声明中明确定义了判定规则——是否通过设备级自动化覆盖完整用户旅程——可据此判定新 surface 的归类。

## Feasibility Assessment

### Technical Feasibility

完全可行。Forge 已有 surface 分化的基础设施：
- `gen-test-scripts/types/` 下 5 种生成策略
- `run-tests/rules/surfaces/` 下 5 种编排规则
- `init-justfile/rules/surfaces/` 下 5 种 justfile 模板
- `gen-journeys/rules/surface-*.md` 下 5 种 surface 规则

本提案将这些已有分化在**概念层**和**命名层**统一表达。

**Task type 解析验证**：`test.gen-scripts.cli` 遵循 Forge 已有的 `{action}.{skill}.{surface}` 三段式命名。当前 task-lifecycle parser 已支持带点的类型名（`eval.design`、`eval.contract`、`gen.journeys` 等），新增 `test.gen-scripts.cli`/`test.run.cli` 等类型名无需 parser 改动。

### Blast Radius

受影响文件清单（基于代码库搜索）：
- `gen-test-scripts/`：5 个 `types/` 文件 + 1 个 SKILL.md = 6 文件
- `run-tests/`：5 个 `rules/surfaces/` 文件 + 1 个 SKILL.md = 6 文件
- `init-justfile/`：5 个 `rules/surfaces/` 文件 + 1 个 SKILL.md = 6 文件
- `gen-journeys/`：5 个 `rules/surface-*.md` 文件 = 5 文件
- 文档：ARCHITECTURE.md、guide.md = 2 文件
- Business rules：task-lifecycle.md = 1 文件
- 新增：测试类型概念参考文档 = 1 文件
- **总计约 27 个文件**，其中 21 个为术语/注释更新，6 个为模板/逻辑变更

### Resource & Timeline

纯文档 + 命名变更，不涉及核心逻辑重构。预计工作量：
- 概念参考文档：1 个 doc 任务
- 术语更新（skill 文件、文档）：5 个 doc 任务（按 skill 分组）
- 命名变更（task type、justfile recipe）：3 个 coding 任务（gen-test-scripts、run-tests、init-justfile 各 1 个）

总任务数 ≤ 9 个，使用 `/quick-tasks` 直接生成任务即可。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| "所有 Forge 生成的测试都是 e2e 测试" | Assumption Flip：端到端测试应覆盖从用户输入到持久层再回到用户可见输出的完整路径。逐一审视：CLI 测试只覆盖子进程调用边界，不涉及持久层；API 测试只覆盖 HTTP 请求-响应边界，客户端不在测试范围内；Web 测试通过浏览器自动化覆盖从用户操作到 UI 反馈的完整旅程；Mobile 测试通过 Maestro 覆盖从用户手势到 UI 反馈的完整旅程 | Overturned: Web 和 Mobile 覆盖完整用户旅程，CLI 和 API 只覆盖服务端进程边界 |
| "用户不需要知道测试类型差异" | Stress Test：当 CLI 项目的测试报告说 "e2e 覆盖率 100%" 时，用户会以为所有功能都端到端验证了，但实际上只验证了子进程层面 | Confirmed: 误导性命名影响用户判断 |
| "统一叫 e2e 可以简化概念" | Occam's Razor：简化了命名但增加了认知负担——用户需要自行区分 "这个 e2e 测试实际做了什么" | Refined: 统一名称 ≠ 简化概念，精确命名才是真正的简化 |

## Scope

### In Scope

- 定义 Surface → Test Type 的映射模型和语义
- 梳理当前所有使用 "e2e" 术语的文件和位置
- 编写测试类型概念参考文档
- **更新 guide.md（Terminology 部分），补充 Surface Type → Test Type 的简要说明**，使所有 agent 在任务执行时能正确使用测试类型术语
- 更新 ARCHITECTURE.md 中的测试相关章节
- 更新 skill SKILL.md 文件中的测试类型术语
- 更新 justfile recipe 命名（从 `test`/`test-e2e` 到 `<surface-key>-test`，如 `cli-test`）
- 更新 task type 命名（携带 surface 信息）
- 更新 business rules 文档中的测试相关术语
- 更新 gen-test-scripts 输出的测试代码中的注释/标签
- 更新 run-tests 的测试输出格式，使 suite 名称和标签使用 surface-specific 测试类型名称

### Out of Scope

- 测试管线流程的重构（gen-journeys → gen-contracts → gen-test-scripts → run-tests 的流程不变）
- 质量门判定逻辑的改动（两层门结构不变，通过/失败判定规则不变）
- eval 管线的改动（eval-journey、eval-contract 的评分维度不变）
- 测试目录结构的重新组织（保持 `tests/<journey>/` 或按 surface-key 分目录可作为后续优化）

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 术语变更导致已有文档/教程失效 | M | M | 在变更点添加术语映射表（旧术语 → 新术语），方便迁移 |
| justfile recipe 重命名影响已有项目的 CI 流程 | M | H | [已覆盖] 本提案原定提供 2 版本向后兼容 alias（旧名 → 新名），v3.0.0 中由 `test-pipeline-consistency-audit` 提案直接破坏性变更替代 |
| "功能测试"标签与行业术语冲突——行业中 "functional test" 泛指验证功能需求的测试 | M | M | 每次使用都带 surface 前缀限定范围（"CLI 功能测试" 而非 "功能测试"），在概念文档中明确本提案中的定义与行业定义的差异 |
| 新概念增加用户学习成本 | L | L | `just --list` 直接展示按测试类型命名的 recipe，用户无需查阅文档即可发现测试类型；概念文档以一页纸为限，映射表一目了然 |
| 分类二分法（功能/端到端）在新增 surface 类型时失效 | L | H | 分类标准声明以"验证机制"为判定依据（协议级调用 vs 设备级自动化），该判定规则可适用于任何新 surface；若出现混合验证机制的 surface（如 desktop SDK 同时支持 CLI 和 GUI），可在概念文档中新增第三分类而非强行归入现有二分法 |
| 术语回归——未来贡献者重新引入 "e2e" 作为通用标签 | M | M | SC3 在完成时验证术语一致性；长期防护通过在概念参考文档中明确标注"e2e 仅用于 Web/Mobile surface 的端到端测试"实现，review 流程中检查新增文件是否遵循测试类型命名规范 |
| 过渡期内部分 skill 文件未完成术语更新，导致新旧术语混用 | M | M | 所有 skill 文件的术语更新在单次 PR 中完成（约 27 个文件），不采用分批更新策略，消除混用窗口期 |

## Success Criteria

- [ ] 概念参考文档完成，包含 5 种 surface 的测试类型定义、语义、验证维度和执行模型
- [ ] guide.md Terminology 部分包含 Surface → Test Type 映射的简要说明，agent 可据此正确使用测试类型术语
- [ ] Forge 代码库中不再有将所有生成测试统称为 "e2e" 的地方（搜索英文 "e2e" 只出现在 Web/Mobile surface 的端到端测试上下文中；搜索中文 "端到端" 仅出现在 Web/Mobile surface 的测试类型名称和定义中）
- [ ] 所有涉及测试类型的 skill 文件使用 surface-specific 测试类型名称
- [ ] 所有包含测试规则的 skill rules 文件（gen-test-scripts、run-tests、init-justfile、gen-journeys 的 surface rules）引用概念文档中的测试类型定义
- [ ] task type 命名在 index.json 中携带 surface 信息（如 `test.gen-scripts.cli`），遵循 `{action}.{skill}.{surface}` 三段式格式
- [ ] business rules 文档（task-lifecycle.md）中的保留类型列表已包含新的 test 类型名，与 task type 命名变更同步
- [ ] 旧 justfile recipe 名（`test`、`test-e2e`）被 `<surface-key>-test` 直接替代（v3.0.0 破坏性变更，由 `test-pipeline-consistency-audit` 提案覆盖）
- [ ] 测试执行输出中的 suite 名称和标签使用 surface-specific 测试类型名称（如 `cli` 而非 `e2e`），质量门报告中不同测试类型的执行结果分类展示

## Next Steps

变更范围已量化为约 9 个任务（6 doc + 3 coding），使用 `/quick-tasks` 直接生成任务。
