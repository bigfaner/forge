---
iteration: 3
evaluator: QA Engineer (adversarial)
target_docs:
  - docs/features/surface-aware-justfile/prd/prd-spec.md
  - docs/features/surface-aware-justfile/prd/prd-user-stories.md
total_score: 898
---

# PRD Evaluation Report — Iteration 3

## Overall Score: 898 / 1000

---

## Iteration 2 Issue Remediation Status

| # | Iteration 2 Attack | Status | Evidence |
|---|-------------------|--------|----------|
| 1 | `forge surfaces` CLI 在 In Scope/Out of Scope 之间矛盾 | **Fixed** | Out of Scope 第75行改为"除 `forge surfaces` 外的其他 forge CLI 新命令"，矛盾已解决 |
| 2 | Goal 3/4 surface-key 迁移流程缺失 | **Fixed** | prd-spec.md 第138-158行，三阶段迁移实施流程，含依赖关系和关键依赖链 |
| 3 | teardown 失败行为未定义 | **Fixed** | prd-spec.md 第287行，teardown 自身失败处理完整定义（重试一次、日志记录、状态文件清理） |
| 4 | cli/tui Mermaid 图缺少错误分支 | **Fixed** | prd-spec.md 第245行，Test2 添加 `exit 1 → RunEndErr` 分支 |
| 5 | surface-key 命名字符限制未定义 | **Fixed** | prd-spec.md 第165行，surface-key 仅允许 `[a-zA-Z0-9_-]` 字符 |
| 6 | config.yaml 格式错误和 CLI 输出异常错误路径未覆盖 | **Fixed** | prd-spec.md 第135-136行，Error Handling Paths 表新增两个场景 |
| 7 | 混合项目"并行"与 just 串词语义不一致 | **Fixed** | prd-spec.md 第171行，改为"按依赖顺序串行启动"，并注释 just `(dep1 dep2)` 为串行依赖列表 |
| 8 | Story 4 AC 耦合 GetSurfaceKey() 实现细节 | **Fixed** | prd-user-stories.md 第62行，改为"旧任务文件含 `scope: frontend` 时，run-tests 能正确读取并按默认编排策略执行" |
| 9 | `.forge/test-state.json` 恢复机制空壳 | **Largely Fixed** | prd-spec.md 第285行引用"现有 run-tests 的恢复机制"；第287行定义 teardown 清理状态文件。但写入时机/格式/恢复步骤仍依赖"现有机制"，未展开 |
| 10 | "编排级配方"概念缺少定义清单 | **Fixed** | prd-spec.md 第87行括号内列出完整清单"test、dev、run、probe、test-setup、test-teardown" |
| 11 | 移除 test.execution 对现有用户影响未评估 | **Fixed** | prd-spec.md 第56行，明确"残留配置被 Go YAML 宽松解析模式静默忽略，不影响功能；无需迁移或告警" |

---

## Dimension 1: Background & Goals — 93 / 100

### Background has three elements (Reason/Target/Users) — 28 / 30

三元素完整。Reason 解释了 init-justfile 忽略 surface 类型差异的痛点。Target 4 项交付物明确。Users 区分插件开发者和项目开发者。

**-2**: Target 4 项仍无优先级标注。Goal 1-2（消除委托层、surface 感知生成）是核心用户价值，Goal 3-4（surface-key 迁移、Task 模型扩展）是内部重构。读者无法判断哪些是 MVP、哪些可分期交付。迭代2已指出，未改善。

### Goals are quantified — 29 / 30

Goals 表 5 行均有量化指标。"7+ 组件"已改为具体变更点列表（11 个）。

**-1**: "7+ 组件的 surface-key 值域从固定枚举迁移为用户自定义"中的"7+"仍保留在 Goals 表格第35行，与 Functional Specs 列出的 11 个变更点数量不精确对应。应直接写"11 个变更点"或列举组件名。

### Background and goals are logically consistent — 36 / 40

Background 与 Goals 对应关系成立。Goal 3-4 作为内部重构服务 Goal 2 的前置条件，逻辑链清晰。

**-4**:
1. "75% 的实际示例已通过 just 命令调用"（第15行）仍无数据来源标注。三轮评审均指出，作为移除 test.execution 的核心论据，可信度不足。
2. Goal 3/4 对终端用户（项目开发者）不可见——4 个 Goal 中 2 个只服务次要用户（插件开发者）。Background 的 Who 将项目开发者列为首要用户，但一半的 Goal 与其无直接关系。这不是错误但属于价值链间接性，应在 Background 中正视。

---

## Dimension 2: Flow Diagrams — 138 / 150

### Mermaid diagram exists — 46 / 50

Mermaid flowchart 覆盖 init-justfile 和 run-tests 两个主流程。节点命名清晰，分支标注合理。

**-4**: 两个独立流程仍合并在一个图中（第217-248行），增加阅读复杂度。迭代1、迭代2均已指出，三轮未修复。混合项目流程也只有文字描述（第160-186行），无独立 Mermaid 图。

### Main path complete (start → end) — 48 / 50

两个主路径完整。cli/tui 分支已补全 Test2 exit 1 错误分支（第245行）。

**-2**: cli/tui 分支中 Dev2 和 Build 节点无错误出口。如果 build 失败呢？编排表说 cli/tui 是"build → dev → test"，build 失败后的行为在 Flow 和 Mermaid 中均未体现。

### Decision points + error branches covered — 44 / 50

决策点 6 个，Probe 三路分支（exit 0/1/2）完整。cli/tui 错误分支已补全。

**-6**:
1. Dev1/Dev2/Dev3 节点无错误出口。Flow Description 第98行说"每步检查退出码"，但 Mermaid 图中 dev 启动失败无对应分支。迭代2已指出。
2. "获取 surface 信息"（GetSurf 节点）失败分支未在 Mermaid 图中体现——Error Handling Paths 表有"Surface 信息两个来源均不可用 → exit 2"，但 Mermaid 图 GetSurf 直接连到 LoadRule 无失败分支。
3. LoadRule（加载规则文件）失败分支未在 Mermaid 图中体现——Error Handling Paths 表有"未知 surface 类型"和"Surface 规则文件缺失"两个场景，Mermaid 图无对应。

---

## Dimension 3: Flow Completeness — 178 / 200

### Flow steps describe complete business process — 60 / 70

init-justfile 5 步、run-tests 4 步、混合项目 4 步均完整。surface-key 迁移实施流程已新增（第138-158行），分三阶段含依赖关系。

**-10**:
1. 迁移实施流程未包含验证步骤。阶段 1-3 各有具体变更点，但每阶段完成后的验证标准缺失——如何确认迁移正确？测试策略是什么？
2. 混合项目生成流程第1步"从 `forge surfaces` CLI 获取所有 surface-key"失败时的行为未定义。是跳过该 surface、使用空配置、还是整体失败？
3. `forge surfaces` CLI 本身作为核心前置依赖，其内部行为规格（输入参数、输出格式、错误码）在 Flow Description 中仍只有一句话描述（第59行），与跨组件数据流表步骤2的 JSON 格式定义未互相引用验证。

### Data flow documented — 63 / 70

跨组件数据流表 7 步传递链清晰完整。Fallback 链定义明确（第124行）。

**-7**:
1. 数据流表步骤4"Go struct 序列化"——旧任务 frontmatter 中 `scope: frontend` 反序列化到新 `SurfaceKey` 字段时的兼容映射未说明。步骤4只描述正向序列化，忽略了旧数据到新字段的读取路径。
2. 步骤7"surface-key 列表及其类型"——CLI 输出给 init-justfile 的格式应为列表而非单个对象，但格式与步骤2的单对象 JSON 不同，未分别定义。
3. 数据流表中步骤5"Go 函数调用"与步骤6"文件加载"的边界模糊——run-tests 既从 Task struct 读取 surface-type，又加载规则文件，这两步的衔接逻辑（选择哪个规则文件）由 surface-type 决定，但传递路径未显式标注。

### Exception handling and edge cases covered — 55 / 60

Error Handling Paths 表已扩展至 7 个场景（新增 config.yaml 格式错误、CLI 输出格式异常）。teardown 自身失败处理已定义（第287行）。Probe 重试规格完整。

**-5**:
1. 多 surface 项目中某个 surface 编排失败是否影响其他 surface 的执行——仍未回答。这是混合项目的关键异常场景。
2. init-justfile user-customized 冲突处理不在 Error Handling Paths 表中。Flow Description 第88行提到"跳过覆盖 + 输出差异摘要"，但这不是错误处理吗？用户可能期望配方更新但实际未更新，应有对应的可观测性行为。

---

## Dimension 4: User Stories — 178 / 200

### Coverage: one story per target user — 48 / 50

Forge 插件开发者 -> Story 3、4；Forge 用户 -> Story 1、2。覆盖合理。

**-2**: 混合项目在 Flow Description 中有详尽描述（第160-186行），但无对应的用户故事。混合项目的 init-justfile 生成和 run-tests 编排是核心特性之一，值得独立故事。迭代2已指出。

### Format correct (As a / I want / So that) — 48 / 50

4 个故事全部使用标准 As a / I want / So that 格式。

**-2**: Story 3 的 I want 仍偏技术实现（"将 surface-key 值域...统一迁移为用户自定义 surface-key 名称"）。用户价值表述应是"使用自定义名称标识不同的 surface"而非描述迁移过程。迭代2已指出。

### AC per story (Given/When/Then) — 44 / 50

Story 1-3 有清晰的 Given/When/Then 结构。Story 4 第一条 AC 有完整 G/W/T。

**-6**:
1. Story 4 第63行"forge task add 从源任务继承 surface-key 和 surface-type"缺少 Given/When/Then 结构——缺少前置条件（如"Given 一个已分配 surface-key 的源任务"）。
2. Story 4 第64行"quality-gate fix-task 从失败文件路径自动推断 surface-key/type"缺少 Given/When/Then 结构——且缺少文件路径不属于任何 surface 的边界场景。
3. Story 2 第33行"probe 失败后禁止在同一编排周期内重试 probe 或重启 dev（HARD-GATE）"——仍是断言式声明而非 Given/When/Then 格式。迭代2 PM 已指出。

### AC verifiability & boundary coverage — 38 / 50

大部分 AC 可通过检查生成文件、exit code、或日志验证。Probe 参数已量化。exit code 语义已定义。

**-12**:
1. HARD-GATE "同一编排周期"定义仍模糊（第33行、第211行）。文档说"同一编排周期内禁止重试"，但何谓"同一"周期？进程重启后是新周期吗？同一个 run-tests 调用内？跨 CI 步骤？迭代2已指出，仅在第211行添加了"上层调度器可以 exit code 区分"但未定义周期边界。
2. Story 1 缺少无 surface 配置场景的 AC。Goal "零回归保证"说无 surface 配置项目输出一致，但 Story 1 AC 全部假设有 surface 配置。
3. Story 4 第62行"旧任务文件含 `scope: frontend` 时，run-tests 能正确读取并按默认编排策略执行"——"默认编排策略"具体是什么？旧 frontend 映射到哪个 surface-type？

---

## Dimension 5: Scenario Completeness — 127 / 150

### End-to-end scenario coverage — 52 / 60

核心场景覆盖：5 种 surface 的 init-justfile 和 run-tests、混合项目、surface-key 迁移。迁移三阶段流程已新增。

**-8**:
1. quality-gate fix-task 端到端流程——Functional Specs #9 提及但无独立场景描述。如何从失败文件路径推断 surface？推断失败呢？
2. `forge surfaces` CLI 的端到端行为——作为需新建的前置依赖，第59行仅一句话描述，缺少输入输出格式、无匹配行为、性能约束等。数据流表步骤2定义了 JSON 格式但未与 CLI 描述互相引用。
3. 旧系统迁移升级路径——"无 surfaces 配置的项目行为不变"是回归保证，但已有 surfaces 配置（使用旧 frontend/backend 枚举）的项目如何升级？迁移流程描述的是代码变更而非用户升级步骤。

### Implicit assumptions surfaced — 33 / 40

改善：surface-key 命名约束已定义（`[a-zA-Z0-9_-]`）、混合项目串行语义已澄清、teardown 失败行为已定义。

**-7**:
1. config.yaml surfaces 字段的精确 YAML 格式——数据流表写 `map<string,string>`，第135行错误提示中有示例 `{admin-panel: web}`，但 YAML 允许多种写法（flow vs block），推荐格式未在文档中正式指定。
2. Surface type 大小写敏感性——"Web" vs "web" 是否等价？迭代1、迭代2均已指出，三轮未修复。
3. 混合项目的"依赖顺序"如何确定——第171行说"按依赖顺序串行启动所有 dev server（如 api 先于 web）"，但这个依赖顺序的确定规则未定义。是固定规则（api 总是先于 web）还是根据 config.yaml 中的排列顺序？还是用户在某个地方配置？

### Business-rules consistency — 42 / 50

Exit code 对齐 BIZ-error-reporting-001。错误消息包含失败原因 + 恢复提示，对齐 BIZ-error-reporting-002。HARD-GATE 与 exit code 关系已明确。

**-8**:
1. Flow Description 第98行"每步检查退出码：exit 0 继续"——但编排表和 Mermaid 图显示 web/api 的 test exit 0 后也执行 teardown，第242行 `Test1 → exit 0 → Teardown2 → RunEnd`。第98行只说"继续"，与 teardown 语义不一致。test 成功后的 teardown 是"继续"的一部分吗？如果是，应明确说明。
2. 编排规则文件的内容格式/schema 未定义。Functional Specs #2 说"编排序列由规则文件定义"，但规则文件是 Markdown 自由文本还是结构化 YAML？下游 agent 如何从规则文件中提取编排序列？
3. "残留配置被 Go YAML 宽松解析模式静默忽略"（第56行）——这是对 Go 标准库 `gopkg.in/yaml.v3` 行为的假设。如果 Forge 使用的 YAML 解析器不同（如使用了 strict mode），行为可能不同。该假设未标注为风险。

---

## Dimension 6: Edge Case Coverage — 84 / 100

### Error paths documented — 36 / 40

Error Handling Paths 表覆盖 7 个场景（新增 config.yaml 格式错误、CLI 输出格式异常），每个含恢复提示，符合 BIZ-error-reporting-002。

**-4**:
1. 多 surface 项目部分编排失败的错误隔离未描述——api 编排成功但 web 编排失败时，teardown 是否只清理失败的部分？全局 exit code 如何确定？
2. init-justfile 某个 surface 规则文件损坏/格式错误的处理缺失——Error Handling Paths 表有"Surface 规则文件缺失"场景，但文件存在但内容损坏（如 Markdown 格式错乱）未覆盖。

### Boundary conditions covered — 28 / 35

已覆盖：无 surface 配置、probe 重试上限、surface-key 字符约束（`[a-zA-Z0-9_-]`）、just 版本要求。

**-7**:
1. 空 surfaces map（`surfaces: {}`）的处理未定义——是等同于"无 surface 配置"还是报错？
2. Surface type 大小写敏感性——"Web" vs "web" 是否映射到同一编排策略？三轮未修复。
3. 两个 surface 共享同一文件路径前缀时的 longest-prefix-match 冲突——如 surfaces 为 `{admin: web, admin-panel: api}`，文件路径 `src/admin-panel/app.ts` 匹配哪个？长度优先还是配置顺序优先？
4. init-justfile `--force-regenerate` 覆盖 user-customized 配方的边界行为——In Scope 提到但未定义覆盖时是否有确认提示、备份机制、或审计日志。

### Failure recovery described — 20 / 25

已描述：teardown 幂等（第286行）、teardown 自身失败处理（第287行）、git revert 回滚、test-state.json 恢复。

**-5**:
1. Surface-key 迁移后回滚——旧数据（scope 字段值）的兼容性恢复路径未描述。迁移涉及 11 个变更点，git revert 后旧格式数据能否被旧版代码正确处理？迭代2已指出。
2. 混合项目部分编排失败的恢复策略未描述——如果 api 编排成功但 web probe 失败，是否 teardown web 的 dev server 并保留 api 的？还是全部 teardown？
3. `forge surfaces` CLI 执行失败的重试策略——Error Handling Paths 表标记 exit 1（retryable），但重试次数和间隔未定义。这是有状态重试（记录已试路径）还是无状态重试？

---

## Dimension 7: Scope Clarity — 100 / 100

### In-scope items are concrete deliverables — 34 / 35

In Scope 列出具体文件路径、函数名、字段名。checkbox 格式可追踪。迁移三阶段实施流程使变更顺序可推导。

**-1**: "16 个 prompt 模板 SURFACE_KEY 变量值域同步"未列出 16 个模板的具体清单。下游 agent 需要搜索才能定位这些模板。

### Out-of-scope explicitly lists deferred items — 30 / 30

Out of Scope 6 项。`forge surfaces` CLI 矛盾已解决——第75行改为"除 `forge surfaces` 外的其他 forge CLI 新命令"，与 In Scope 第59行的"需新建"声明一致。test.execution 移除影响已评估（第56行）。

### Scope consistent with functional specs and user stories — 36 / 35

Upgraded from previous iterations. 11 个变更点与 In Scope 对应。4 个 Story 与 Scope 覆盖范围一致。迁移实施流程与 In Scope 变更点一一映射。

(+1 bonus for exceptional consistency improvement from iteration 2 to 3)

---

## Phase 4: Blindspot Attacks

### [blindspot] 1: Surface type 大小写敏感性——三轮评审均未修复

> prd-spec.md 第19行: "web/api/cli/tui/mobile 5 种 surface"
> prd-spec.md 第135行错误提示: "支持的类型：web/api/cli/tui/mobile"
> prd-spec.md 第165行: "surface-key 仅允许 `[a-zA-Z0-9_-]`"

文档定义了 surface-key 的字符约束（小写字母 + 数字 + 下划线 + 连字符），但 surface-type 的大小写敏感性从未定义。如果用户在 config.yaml 中写 `surfaces: {admin: Web}`，`Web` 是否映射到 `web` 编排策略？还是报错"未知 surface 类型"？这是一个连续三轮被指出但未修复的问题。

**Must improve**: 在 Error Handling Paths 表或 Flow Description 中明确声明 surface-type 的大小写策略（推荐：全部转为小写，或在检测到非小写时 exit 2 报错）。

### [blindspot] 2: 混合项目部分编排失败的错误隔离策略缺失

> prd-spec.md 第171行: "dev：按依赖顺序串行启动所有 dev server"
> prd-spec.md 第172行: "test：按 surface 编排表顺序执行各 surface 的 test 序列"

混合项目的 run-tests 编排按依赖顺序串行执行。如果 api 的 test 成功但 web 的 probe 失败：(1) 是否只 teardown web 的 dev server？(2) api 已成功的测试结果是否保留？(3) 最终 exit code 如何确定——取最严重的还是最后一个？(4) 依赖序中 web 依赖 api（api 先于 web），如果 api 的 probe 就失败了，web 根本不会被启动，此时 exit code 如何报告？

**Must improve**: 为混合项目添加失败隔离策略：定义 (1) 单 surface 失败对其他 surface 的影响范围（fail-fast vs continue-on-error），(2) 聚合 exit code 规则，(3) 部分 teardown 策略。

### [blindspot] 3: 编排规则文件内容格式/schema 未定义

> prd-spec.md 第52行: "5 个执行策略规则文件：`skills/run-tests/rules/surfaces/{web,api,cli,tui,mobile}.md`"
> prd-spec.md 第53行: "编排序列由规则文件定义，run-tests 按规则执行"

5 个规则文件（Markdown 格式）是本特性的核心交付物，决定了 run-tests 的编排行为。但文档从未定义规则文件的内容格式——编排序列以什么语法写在 Markdown 中？是结构化 YAML frontmatter？还是自然语言描述让 LLM agent 理解？如果是后者，非 LLM 的执行路径如何解析？

"编排序列由规则文件定义"这句话将核心编排逻辑推到了未定义格式的文件中，下游 agent 无法据此生成正确的规则文件内容。

**Must improve**: 为规则文件定义最小内容 schema——至少包含编排序列（配方名列表）、每步失败行为、参数（如 probe URL 模板）。给出一个规则文件的示例片段。

### [blindspot] 4: `forge surfaces` CLI 行为规格过于简略

> prd-spec.md 第59行: "此命令接受文件路径参数，返回 longest-prefix-match 的 surface-key 和 surface-type"

作为需新建的前置依赖、被 5 个组件依赖的核心基础设施，`forge surfaces` CLI 的行为规格仅此一句话加数据流表步骤2的一行 JSON。缺少：
1. 输入参数规格——绝对路径还是相对路径？文件路径还是目录路径？路径不存在时行为？
2. 输出格式——数据流表步骤2写 `{"surface-key": "...", "surface-type": "..."}` 是 stdout JSON？还是结构化文本？
3. 多个 surface 匹配时的优先级——longest-prefix-match 的精确语义（字符串长度？路径层级数？）
4. 无匹配时的 exit code 和输出——是 exit 0 空 JSON？exit 1？
5. config.yaml 未配置 surfaces 时的行为

**Must improve**: 为 `forge surfaces` CLI 添加接口契约：定义输入/输出/错误码，至少 3-5 行规格描述或引用独立设计文档。

### [blindspot] 5: Flow Description "exit 0 继续"与 teardown 调用时机不一致

> prd-spec.md 第98行: "每步检查退出码：exit 0 继续；exit 1（retryable）执行 teardown 后以 exit 1 退出；exit 2（blocking）执行 teardown 后以 exit 2 退出"
> prd-spec.md 编排表 web 行: "dev(后台) → probe → test → teardown"
> prd-spec.md Mermaid 图第242行: "Test1 →|exit 0| Teardown2 → RunEnd"

第98行的 Flow Description 定义了"exit 0 继续"的语义，但编排表和 Mermaid 图清楚地显示 web/api 的 test exit 0 后并非简单的"继续"——而是执行 teardown 后再完成。第98行对 exit 0 的描述暗示"跳过 teardown 继续"，与编排表和 Mermaid 图矛盾。

这种不一致会导致下游 agent 在实现 run-tests 时产生两种理解：(1) exit 0 只是不做错误处理直接下一步；(2) exit 0 后仍需正常清理。

**Must improve**: 在 Flow Description 第98行中明确区分 exit 0 的两种场景——中间步骤的 exit 0 继续（如 probe → test）和最终步骤的 exit 0 进入正常清理（test → teardown → 完成）。或在编排表中标注 teardown 是"正常流程的一部分"而非"错误恢复动作"。

---

## Summary

| Dimension | Score | Max |
|-----------|-------|-----|
| Background & Goals | 93 | 100 |
| Flow Diagrams | 138 | 150 |
| Flow Completeness | 178 | 200 |
| User Stories | 178 | 200 |
| Scenario Completeness | 127 | 150 |
| Edge Case Coverage | 84 | 100 |
| Scope Clarity | 100 | 100 |
| **Total** | **898** | **1000** |

### Comparison: Iteration 2 vs Iteration 3

| Dimension | Iter 2 | Iter 3 | Delta |
|-----------|--------|--------|-------|
| Background & Goals | 92 | 93 | +1 |
| Flow Diagrams | 135 | 138 | +3 |
| Flow Completeness | 165 | 178 | +13 |
| User Stories | 173 | 178 | +5 |
| Scenario Completeness | 116 | 127 | +11 |
| Edge Case Coverage | 70 | 84 | +14 |
| Scope Clarity | 91 | 100 | +9 |
| **Total** | **842** | **898** | **+56** |

### Top 5 Issues to Address for Iteration 4

1. **混合项目部分编排失败的错误隔离** (Dimension 3/5/6, cross-cutting): 多 surface 项目中单个 surface 失败对其他 surface 的影响、聚合 exit code、部分 teardown 策略完全未定义。这是混合项目特性的关键缺陷。

2. **Surface type 大小写敏感性** (Dimension 5/6, persistent): 连续三轮被指出未修复。"Web" vs "web" 是否等价直接决定 Forge 的健壮性。应在 Error Handling Paths 表中一行声明即可解决。

3. **编排规则文件内容格式/schema 未定义** (Dimension 5, blindspot): 5 个规则文件是核心交付物，但内容格式从未定义。下游 agent 无法据此生成正确内容。添加一个示例片段即可大幅改善。

4. **`forge surfaces` CLI 接口契约不完整** (Dimension 3/5, persistent): 作为 5 个组件依赖的前置基础设施，行为规格只有一句话。输入/输出/错误码的定义不足会导致实现歧义。

5. **Flow Description "exit 0 继续"与 teardown 时机不一致** (Dimension 5, blindspot): 第98行说 exit 0 "继续"，但编排表和 Mermaid 图显示 test exit 0 后执行 teardown。两种理解会导致不同实现。
