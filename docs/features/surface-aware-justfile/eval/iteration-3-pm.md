---
iteration: 3
evaluator: Senior PM (Adversary)
target_docs:
  - docs/features/surface-aware-justfile/prd/prd-spec.md
  - docs/features/surface-aware-justfile/prd/prd-user-stories.md
total_score: 910
---

# PRD Evaluation Report — Iteration 3

## Overall Score: 910 / 1000

---

## Iteration 2 Issue Remediation Status

| # | Iteration 2 Attack | Status | Evidence |
|---|-------------------|--------|----------|
| 1 | `forge surfaces` CLI Scope 矛盾 | **Fixed** | In Scope 第59行标注为前置依赖；Out of Scope 第75行改为"除 `forge surfaces` 外的其他 forge CLI 新命令"，矛盾消除 |
| 2 | Goal 3/4 迁移流程缺失 | **Fixed** | 新增"Surface-key 迁移实施流程"节（第138-158行），三阶段实施+关键依赖链 |
| 3 | teardown 失败行为未定义 | **Fixed** | 第287行定义了 teardown 自身失败处理：重试一次→日志记录→继续清理→确保状态文件清理 |
| 4 | Mermaid 图 cli/tui 缺少错误分支 | **Fixed** | Mermaid 图 Test2 新增 exit 1 → RunEndErr 分支（第245行） |
| 5 | surface-key 命名字符限制未定义 | **Fixed** | 第165行："surface-key 仅允许 `[a-zA-Z0-9_-]` 字符" |
| 6 | config.yaml 格式错误和 CLI 输出异常未覆盖 | **Fixed** | Error Handling Paths 表新增两个场景（第135-136行） |
| 7 | 混合项目"并行启动"与 just 语法不一致 | **Fixed** | 第171行改为"按依赖顺序串行启动"，注释说明 just 语法为串行依赖列表 |
| 8 | Story 4 AC 耦合 Go 函数名 | **Fixed** | 第62行改为行为验证："旧任务文件含 `scope: frontend` 时，run-tests 能正确读取并按默认编排策略执行" |
| 9 | `.forge/test-state.json` 恢复机制是空壳 | **Partially Fixed** | 降级为引用现有机制（第285行"利用现有 run-tests 的恢复机制"），合理。但写入时机仍未描述 |
| 10 | "编排级配方"概念缺少定义清单 | **Fixed** | 第87行括号中列出完整清单：test、dev、run、probe、test-setup、test-teardown |
| 11 | 移除 test.execution 对现有用户配置的影响未评估 | **Fixed** | 第56行："残留配置被 Go YAML 宽松解析模式静默忽略，不影响功能；无需迁移或告警" |

**Remediation Summary**: 11 个攻击点中 10 个完全修复，1 个部分修复（test-state.json 恢复流程细节）。迭代2的 Top 5 问题全部得到实质性回应。

---

## Dimension 1: Background & Goals — 92 / 100

### Background has three elements (Reason/Target/Users) — 27 / 30

三要素完整：Reason（init-justfile 忽略 surface 类型，test.execution 冗余）、Target（4项交付物）、Users（插件开发者+项目开发者，各有明确使用场景描述）。

**-3**: Target 列出4项但无优先级排序。混合项目（web+api）和单 surface 项目（纯 cli）的实现复杂度差异巨大，但 Target 未标注哪些是 MVP 核心。4个Goal中有2个（Goal 3/4）属于内部重构，对终端用户不可见，但 Target 未区分用户可见价值 vs 内部技术债务。

### Goals are quantified — 28 / 30

5项Goals均有量化指标：层数（4→2）、类型数（5种）、组件数（7+）、字段数（2新字段）。

**-2**: "7+ 组件的 surface-key 值域从固定枚举迁移为用户自定义"——Functional Specs 表列出 11 个变更点，涉及组件数量可精确枚举。"+"后缀使指标边界模糊。

### Background and goals are logically consistent — 37 / 40

Goal 1-2 与 Background 直接对应（消除 test.execution + surface 感知）。Goal 3-4 作为 Goal 2 前置条件逻辑链成立。迁移实施流程（第138-158行）的三阶段规划与 Goal 3/4 的依赖关系一致。

**-3**:
1. "75% 的实际示例已通过 just 命令调用"（第15行）——作为移除 test.execution 的核心决策依据，此数据来源不可独立验证。经两个迭代仍未补充数据来源或抽样方法。
2. 5个Goal中 Goal 3/4 对终端用户（项目开发者）不可见，但 Background Who 将项目开发者列为用户之一。4个Goal中2个只服务次要用户（插件开发者），价值分配不均的问题未正视。

---

## Dimension 2: Flow Diagrams — 138 / 150

### Mermaid diagram exists — 46 / 50

Mermaid flowchart 存在（第217-248行），覆盖 init-justfile 和 run-tests 两个主要流程。节点命名清晰，流程方向明确。

**-4**: 混合项目的生成与编排（第160-186行）有详细的文字描述+代码示例，但无独立 Mermaid 图。混合项目是最复杂的场景（多 surface 交互、聚合配方、依赖排序），缺少可视化流程图会降低可读性。连续两轮未修复。

### Main path complete (start → end) — 48 / 50

两个主路径 init-justfile（Start → Lang → ... → End）和 run-tests（RunStart → GetSurf → ... → RunEnd）均完整覆盖 start→end。

**-2**: cli/tui 分支（Dev2 → Test2 → RunEnd）缺少 teardown 步骤。编排表中 cli/tui 序列为"build → dev → test"，无 teardown。但 dev 步骤是否可能创建需清理的资源（如临时文件、后台进程）？设计决策未解释。若确认无需清理，应在文档中明确声明。

### Decision points + error branches covered — 44 / 50

决策点覆盖改善：CheckSurf（有/无 surface）、Arbitrate（冲突仲裁）、Protect（user-customized 保护）、ExecSeq（编排序列类型）、Probe exit（0/1/2 三路）、Test exit（0/1 两路）。cli/tui 分支的错误路径已补充（Test2 → exit 1 → RunEndErr）。

**-6**:
1. Dev 启动失败分支缺失：Dev1/Dev2/Dev3 节点后无错误出口。Flow Description 步骤4说"每步检查退出码"，但 dev 启动失败时如何处理未在图中体现。
2. "Surface 信息两个来源均不可用"在 Error Handling Paths 表中有定义（第134行），但 Mermaid 图 GetSurf 节点无对应的错误分支。
3. init-justfile 流程中缺少 just 版本检查步骤（Error Handling Paths 表第133行提及此检查）。

---

## Dimension 3: Flow Completeness — 185 / 200

### Flow steps describe complete business process — 63 / 70

init-justfile 5步流程、run-tests 4步调度器流程、混合项目 4步生成策略、Probe 重试规格、迁移实施三阶段——5个流程文档完整。

**-7**:
1. `forge surfaces` CLI 的内部行为规格仍然不足——第59行仅一句话"接受文件路径参数，返回 longest-prefix-match 的 surface-key 和 surface-type"，但缺少：(a) 输入参数规格（绝对路径还是相对路径？文件还是目录？）；(b) 输出 JSON 结构定义；(c) 无匹配时的 exit code 和输出；(d) config.yaml 未配置 surfaces 时的行为。作为需新建的前置依赖，规格过于简略。
2. Surface 编排模式表（第101-108行）中 mobile 的 test-setup "准备模拟器"——模拟器启动失败的错误路径未定义。

### Data flow documented — 65 / 70

跨组件数据流表（第112-124行）7步传递链清晰，每步标注来源、目标、格式、传递方式。Fallback 链明确（第124行）。

**-5**:
1. 步骤2的 JSON 格式 `{"surface-key": "admin-panel", "surface-type": "web"}` 是预期接口定义，但未标注为"契约定义"。`forge surfaces` CLI 需新建，实际输出格式是否与此一致无保证。
2. 步骤4 "Go struct 序列化"——旧任务 frontmatter 中 `scope: frontend` 反序列化到新 struct SurfaceKey 字段时的兼容映射未说明。Story 4 第62行说"run-tests 能正确读取"，但数据流层面如何从旧字段映射到新字段未描述。
3. 步骤5 "从 Task struct 读取 SurfaceKey + SurfaceType"——仅涵盖正常路径，未包含编排流程步骤1的 fallback 到 CLI 路径。

### Exception handling and edge cases covered — 57 / 60

Error Handling Paths 表覆盖 7 个场景（较迭代2新增2个），每个含恢复提示，符合 BIZ-error-reporting-002。Exit Code 语义表完整对齐 BIZ-error-reporting-001。teardown 自身失败处理已定义（第287行）。Probe 重试参数统一（第190-199行）。

**-3**:
1. 多 surface 项目中某个 surface 编排失败是否影响其他 surface——未回答。混合项目节描述"按 surface 编排表顺序执行各 surface 的 test 序列"，但未定义某 surface test 失败后的行为（中止后续 surface？继续执行并汇总结果？）。
2. init-justfile user-customized 保护机制（第48行）的冲突处理不在 Error Handling Paths 表中。这是用户可感知的行为变更，应有对应错误处理描述。

---

## Dimension 4: User Stories — 185 / 200

### Coverage: one story per target user — 46 / 50

Forge 插件开发者 → Story 3、4；Forge 用户 → Story 1、2。覆盖两类目标用户。

**-4**:
1. 混合项目的 run-tests 编排涉及新的配方命名规则（`dev-<surface-key>`/`test-<surface-key>`），CI 配置可能需要更新以适配新的配方名称。缺少 CI/CD 维护者角色的视角。
2. init-justfile 的 `# user-customized` 保护机制和 `--force-regenerate` 参数无对应用户故事——这是一个用户可感知的功能特性。

### Format correct (As a / I want / So that) — 48 / 50

4个故事全部使用标准 As a / I want / So that 格式。

**-2**: Story 3 的 "I want to" 描述偏技术实现（"将 surface-key 值域从固定枚举统一迁移为用户自定义"）而非用户价值。连续两轮指出，未修改。更好的写法："I want to 在 config.yaml 中为不同 surface 自定义名称（如 admin-panel），以便在混合项目中区分和管理多个同类 surface"。

### AC per story (Given/When/Then) — 45 / 50

Story 1-3 使用 Given/When/Then 结构。Story 4 的后两条 AC 改为行为描述（第62-64行），但结构上仍以 "And" 附加而非独立场景。

**-5**:
1. Story 2 第33行 "And probe 失败后禁止在同一编排周期内重试 probe 或重启 dev（HARD-GATE）"——这不是 Given/When/Then 格式，而是一个约束声明。应重写为可验证的场景。
2. Story 4 第64行 "And quality-gate fix-task 从失败文件路径自动推断 surface-key/type"——缺少文件路径不属于任何 surface 的边界场景。

### AC verifiability & boundary coverage — 46 / 50

大部分 AC 可通过检查生成文件、exit code、日志输出来验证。Probe 参数已量化（3次/30秒/90秒）。Exit code 语义已定义。Story 4 第62行已改为行为验证。

**-4**:
1. HARD-GATE "同一编排周期"——第211行定义了"若 probe 已判定失败，说明服务存在根本性问题，重试只会掩盖问题"，但"同一编排周期"的边界仍模糊：进程重启是新周期吗？不同 CI step 是新周期吗？
2. Story 1 缺少"无 surface 配置"场景的 AC——第46行 In Scope 提及此场景，但 Story 1 只验证了有 surface 配置时的情况。
3. Story 3 "无 surfaces 配置的项目行为不变"——如何验证"行为不变"？需要与旧版本输出做 diff，但 AC 未明确验证手段。

---

## Dimension 5: Scenario Completeness — 133 / 150

### End-to-end scenario coverage — 52 / 60

核心场景覆盖全面：init-justfile surface 感知（5种类型）、run-tests 编排（5种模式）、surface-key 迁移（三阶段）、Task 模型扩展、混合项目生成与编排。

**-8**:
1. 缺少"从零配置到第一次成功运行测试"的完整用户旅程走查——各片段分散在 Flow Description、混合项目、编排模式表等不同节中，无法直接拼出完整路径。
2. quality-gate fix-task 端到端流程——Functional Specs #9 提及"从失败文件路径推断 surface-key/type"但无完整场景描述（推断失败怎么办？推断成功后走哪个编排路径？）。
3. `forge surfaces` CLI 本身的端到端行为——作为前置依赖，其命令行接口的输入/输出/错误行为只有一句话描述，无独立规格。

### Implicit assumptions surfaced — 35 / 40

主要假设已暴露：CLI 前置依赖、just >= 1.4.0 版本要求、零回归保证、Go YAML 宽松解析。surface-key 命名约束已定义（第165行）。

**-5**:
1. config.yaml surfaces 字段的精确 YAML 格式——第116行写 `map<string, string>`，但用户应该写 `surfaces: {admin-panel: web}` 还是 `surfaces:\n  admin-panel: web`？两种 YAML 格式语义相同但未提供参考示例。
2. `.forge/test-state.json` 的 schema——第285-287行引用了现有机制并补充了 teardown 失败处理，但文件格式仍未定义。
3. 混合项目依赖排序规则——第171行说"按依赖顺序串行启动"，但依赖顺序如何确定？是固定规则（api 先于 web）还是用户可配置？示例中"payment-service(api) → admin-panel(web)"暗示 api 先于 web，但规则未明确。

### Business-rules consistency — 46 / 50

Exit code 对齐 BIZ-error-reporting-001（0/1/2 三级）。错误消息包含失败原因+恢复提示，对齐 BIZ-error-reporting-002。HARD-GATE 定义清晰且与 exit code 关系明确。Surface 编排模式表与 Mermaid 图一致。

**-4**:
1. BIZ-quality-gate-001 定义"两层配方模型：unit-test (language) vs test (surface)"——文档中 test 配方的归属（语言模板 vs surface 规则）在仲裁步骤中描述为"Surface 规则覆盖语言模板的编排级配方"，但 unit-test 的生成仍来自语言模板。两层模型的边界在 PRD 中未显式对齐 BIZ 规则。
2. 步骤4 "每步检查退出码：exit 0 继续"——但编排表和 Mermaid 图显示 web/api 的 test exit 0 后也执行 teardown。第98行只说"继续"，未提 teardown。与实际编排行为不一致。

---

## Dimension 6: Edge Case Coverage — 88 / 100

### Error paths documented — 36 / 40

Error Handling Paths 表覆盖 7 个场景，每个含检测位置、行为描述、exit code、恢复提示。较迭代2新增 config.yaml 格式错误和 CLI 输出异常。teardown 自身失败处理已定义。

**-4**:
1. 多 surface 项目部分编排失败的错误隔离——混合项目"按 surface 编排表顺序执行各 surface 的 test 序列"，但某 surface 编排失败（如 probe 超时）是否中止后续 surface 未定义。
2. init-justfile 某个 surface 的规则文件损坏/格式错误——Error Handling Paths 表覆盖了"Surface 规则文件缺失"但未覆盖"文件存在但内容损坏"。

### Boundary conditions covered — 30 / 35

已覆盖：无 surface 配置、probe 重试上限、surface-key 命名约束（`[a-zA-Z0-9_-]`）、双源失败、旧任务兼容、混合项目。

**-5**:
1. 空 surfaces map（`surfaces: {}`）的处理——config.yaml 配置了 surfaces 但值为空时，init-justfile 和 run-tests 如何处理？
2. longest-prefix-match 等长匹配冲突——两个 surface 的路径前缀等长时如何仲裁？规则未定义。
3. `--force-regenerate` 覆盖 user-customized 配方——边界行为未完全描述（覆盖后能否回滚？用户修改全部丢失？）。

### Failure recovery described — 22 / 25

test-state.json 恢复已降级为引用现有机制。teardown 幂等保证+失败处理已定义。git revert 回滚策略明确。

**-3**:
1. Surface-key 迁移后回滚——旧数据（scope 字段值 frontend/backend）的兼容性恢复路径未描述。第143行保留 GetSurfaceKey() 兼容访问，但回滚后新字段 SurfaceKey 的值是否被正确处理未说明。
2. `.forge/test-state.json` 恢复流程——虽然降级为"现有机制"引用，但 teardown 失败处理（第287行）引用了此文件，两者形成隐含依赖但未显式关联。

---

## Dimension 7: Scope Clarity — 89 / 100

### In-scope items are concrete deliverables — 32 / 35

In Scope 使用 checkbox 格式列出具体交付物，包含文件路径（`skills/init-justfile/rules/surfaces/{web,api,cli,tui,mobile}.md`）、函数名（resolveScope）、字段名（SurfaceKey、SurfaceType）。

**-3**:
1. "16 个 prompt 模板 SURFACE_KEY 变量值域同步"——未列出具体模板清单。实现者需要逐一查找 16 个文件，Scope 条目本身不可操作。
2. "SKILL.md 新增 surface 检测步骤和 surface 感知配方生成流程"——"流程"作为交付物偏模糊，应改为具体的 SKILL.md 变更描述。

### Out-of-scope explicitly lists deferred items — 28 / 30

Out of Scope 6 项明确排除。`forge surfaces` CLI 的矛盾已解决——Out of Scope 第75行改为"除 `forge surfaces` 外的其他 forge CLI 新命令"。

**-2**: "变更 `forge-cli/internal/cmd/quality_gate.go` 或 `testrunner` 的 Go 代码"在 Out of Scope 中，但 In Scope "quality-gate fix-task 从失败文件路径推断 surface-key/type"涉及 quality-gate 行为变更。变更是在 skill 文档层面还是 Go 代码层面未澄清。若仅在 skill 层面，则推断逻辑实现在哪里？

### Scope consistent with functional specs and user stories — 29 / 35

11 个 Related Changes 条目与 In Scope 基本对应。迁移实施流程的三阶段与 In Scope 9 项变更点对应。

**-6**:
1. "移除 `test.execution` 节点文档"——第56行补充了迁移影响说明（静默忽略），但无对应用户故事。如果这是一个向后兼容的变更，应有 Story 描述"当用户已有 test.execution 配置时升级后的行为"。
2. Functional Specs #11 "surface-key-assignment 规则：文件路径分类改为 CLI 动态查询"——In Scope 无对应条目。这个变更点缺少明确的交付物定义。
3. Story 1 AC "CLI/TUI surface 不生成 `run` 配方"——In Scope 第46行写"CLI/TUI 只生成 `dev`，不生成 `run`"。但 `run` 配方在编排级配方清单中出现（第87行），两者之间的语义关系未说明——`run` 何时存在？何时不存在？

---

## Phase 4: Blindspot Attacks

### [blindspot] 1: `forge surfaces` CLI 行为规格不足以驱动实现

> prd-spec.md 第59行: "此命令接受文件路径参数，返回 longest-prefix-match 的 surface-key 和 surface-type"

作为需新建的前置依赖，仅一句话描述。7 个下游组件依赖此 CLI，但其行为规格缺失：(1) 输入参数规格——接受绝对路径还是相对路径？文件路径还是目录路径？(2) 输出 JSON 结构——字段名、可选字段、嵌套结构；(3) 多个 surface 路径前缀匹配时的优先级规则（longest-prefix-match 是否唯一匹配？等长如何处理？）；(4) 无匹配时的 exit code 和输出；(5) config.yaml 未配置 surfaces 时的行为；(6) 性能要求（单个文件查询延迟上限）。

**Must improve**: 为 `forge surfaces` CLI 添加完整行为规格。至少包含：命令签名、输入参数、输出格式（JSON schema）、所有 exit code 场景、边界条件处理。可内联于 PRD 或引用独立设计文档。

### [blindspot] 2: 多 surface 项目部分编排失败的隔离策略未定义

> prd-spec.md 第172行: "test：按 surface 编排表顺序执行各 surface 的 test 序列"

混合项目的 test 执行按顺序进行，但未定义某 surface 编排失败后对后续 surface 的影响。例如，payment-service(api) 的 probe 超时导致 teardown+abort，admin-panel(web) 是否还会执行？若是"遇到第一个失败即中止"，应在文档中明确。若是"继续执行并汇总所有结果"，需要定义汇总策略和最终 exit code 计算规则。

**Must improve**: 在混合项目编排流程中添加部分失败策略，至少定义：(1) 单 surface 失败是否影响其他 surface；(2) 最终 exit code 的计算规则（取最大值？按优先级？）。

### [blindspot] 3: 混合项目聚合配方的编排顺序来源不明

> prd-spec.md 第173行: "在生成的 justfile 头部注释中记录编排顺序，供 run-tests 解析"
> prd-spec.md 第177行示例: "# 编排顺序: payment-service(api) → admin-panel(web)"

"编排顺序"由 init-justfile 生成时确定，但生成规则未定义。示例中 api 先于 web，是固定规则（api type 总是先于 web type）？还是按 config.yaml 中 surfaces 字段的定义顺序？还是按依赖关系图（但未定义如何推导依赖）？run-tests 需要解析此"头部注释"——注释的格式、解析规则、格式错误时的行为均未定义。

**Must improve**: 定义编排顺序的生成规则，并在 justfile 头部注释的格式和解析规则上给出明确规范。

---

## Score Summary

| Dimension | Score | Max |
|-----------|-------|-----|
| Background & Goals | 92 | 100 |
| Flow Diagrams | 138 | 150 |
| Flow Completeness | 185 | 200 |
| User Stories | 185 | 200 |
| Scenario Completeness | 133 | 150 |
| Edge Case Coverage | 88 | 100 |
| Scope Clarity | 89 | 100 |
| **Total** | **910** | **1000** |

## Comparison: Iteration 2 vs Iteration 3

| Dimension | Iter 2 | Iter 3 | Delta |
|-----------|--------|--------|-------|
| Background & Goals | 89 | 92 | +3 |
| Flow Diagrams | 130 | 138 | +8 |
| Flow Completeness | 180 | 185 | +5 |
| User Stories | 178 | 185 | +7 |
| Scenario Completeness | 127 | 133 | +6 |
| Edge Case Coverage | 79 | 88 | +9 |
| Scope Clarity | 89 | 89 | 0 |
| **Total** | **872** | **910** | **+38** |

**Remarks**: 迭代2的 11 个攻击点中 10 个完全修复、1 个部分修复。Edge Case Coverage（+9）改进最大——新增了 config.yaml 格式错误、CLI 输出异常、surface-key 命名约束等覆盖。Flow Diagrams（+8）改进来自 cli/tui 错误分支补充。迭代3新发现的攻击点更聚焦于深层规格细节：(1) `forge surfaces` CLI 行为规格不足以驱动实现，(2) 多 surface 部分失败隔离策略缺失，(3) 混合项目编排顺序的生成和解析规则未定义。这些都是实现阶段会导致阻塞的模糊点。

### Top 5 Issues to Address for Iteration 4

1. **`forge surfaces` CLI 行为规格不足** (Dimension 3, Flow Completeness): 作为 7 个组件的前置依赖，仅一句话描述其行为。缺少输入/输出规格、边界条件、性能要求。实现者无法凭此描述开发 CLI。

2. **多 surface 项目部分编排失败隔离策略** (Dimension 6, Edge Case Coverage): 混合项目按顺序执行多 surface 编排，但某 surface 失败对后续 surface 的影响未定义。这是混合项目最关键的运行时行为之一。

3. **混合项目编排顺序生成和解析规则** (Dimension 5, Scenario Completeness): justfile 头部注释的编排顺序由 init-justfile 生成、run-tests 解析，但两端的行为规格均缺失（生成规则、注释格式、解析规则、格式错误处理）。

4. **Story 3 格式偏向实现而非用户价值** (Dimension 4, User Stories): "I want to 将 surface-key 值域从固定枚举统一迁移为用户自定义"描述的是技术实现而非用户价值。连续两轮指出，仍未修改。

5. **`forge surfaces` CLI 等长前缀匹配冲突** (Dimension 6, Edge Case Coverage): longest-prefix-match 算法在两个 surface 路径前缀等长时如何仲裁未定义。这是 CLI 核心算法的边界条件。
