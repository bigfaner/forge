---
iteration: 2
evaluator: QA Engineer (adversarial)
target_docs:
  - docs/features/surface-aware-justfile/prd/prd-spec.md
  - docs/features/surface-aware-justfile/prd/prd-user-stories.md
total_score: 842
---

# PRD Evaluation Report — Iteration 2

## Overall Score: 842 / 1000

---

## Iteration 1 Issue Remediation Status

| # | Iteration 1 Attack | Status | Evidence |
|---|-------------------|--------|----------|
| 1 | 跨组件数据流文档缺失 | **Fixed** | prd-spec.md 第110-124行，7步数据流表 + Fallback 链 |
| 2 | Exit code 处理与 HARD-GATE 未定义 | **Fixed** | prd-spec.md 第176-187行，完整 exit code 表 + HARD-GATE 定义 |
| 3 | 错误路径覆盖不足 | **Largely Fixed** | prd-spec.md 第126-134行，Error Handling Paths 表（5场景）；缺少 config.yaml 格式错误、CLI 输出格式异常 |
| 4 | run-tests surface 信息 fallback 缺失 | **Fixed** | prd-spec.md 第93-94行，两来源均失败时的行为明确 |
| 5 | surface-key 自定义 vs 固定类型张力 | **Fixed** | prd-spec.md 第35行，"用户自定义的是 key 名称而非 type 枚举" |
| 6 | forge surfaces CLI 前置依赖未声明 | **Partially Fixed** | prd-spec.md 第59行声明前置依赖；但与 Out of Scope 第75行 "新增 forge CLI 命令" 矛盾 |
| 7 | 混合项目端到端流程缺失 | **Fixed** | prd-spec.md 第136-161行，完整混合项目生成与编排流程 |
| 8 | Probe 重试参数分散 | **Fixed** | prd-spec.md 第164-175行，统一 Probe 重试规格表 |

---

## Dimension 1: Background & Goals — 92 / 100

### Background has three elements (Reason/Target/Users) — 28 / 30

三元素完整：Reason（init-justfile 忽略 surface 类型差异）、Target（4项交付物）、Users（插件开发者 + 项目开发者）。

**-2**: Target 4项交付物无优先级标注。范围偏宽导致读者无法判断哪些是 MVP 核心、哪些可延后。

### Goals are quantified — 28 / 30

Goals 表格 5 行量化指标。surface-key/type 语义已在第35行澄清。

**-2**: "7+ 组件的 surface-key 值域" 仍不精确。Functional Specs 表列出 11 个变更点，涉及组件数量可枚举。

### Background and goals are logically consistent — 36 / 40

Goal 1-2 与 Background 直接对应。Goal 3-4 作为 Goal 2 前置条件逻辑链成立。

**-4**:
1. "75% 的实际示例已通过 just 命令调用"（第15行）无数据来源、无抽样方法。作为移除 test.execution 的核心依据，缺乏可信度。
2. Goal 3/4 对终端用户（项目开发者）不可见，属于内部重构，但 Background Who 部分将项目开发者列为首要用户。4个Goal中2个只服务次要用户，价值链间接性未正视。

---

## Dimension 2: Flow Diagrams — 135 / 150

### Mermaid diagram exists — 45 / 50

存在 Mermaid flowchart（第193-223行），覆盖 init-justfile 和 run-tests。

**-5**: 两个独立流程仍合并在一个图中，可读性降低。迭代1已指出，未修复。

### Main path complete (start → end) — 48 / 50

两个主路径从开始到结束完整。

**-2**: cli/tui 分支（Dev2 → Test2 → RunEnd）缺少 teardown。编排表说"无需 probe"合理，但 dev 步骤是否可能创建需清理的资源？设计决策未解释。迭代1已指出，未修复。

### Decision points + error branches covered — 42 / 50

决策点 6 个（CheckSurf、Arbitrate、Protect、ExecSeq、Probe exit、Test exit）。Probe 三路 exit 0/1/2 分支完整。

**-8**:
1. 缺少 dev 启动失败的分支。Dev1/Dev2/Dev3 节点后无错误出口。
2. CheckSurf 只有"有/无"两个分支，缺少"有但类型无效"的错误出口。Error Handling Paths 表有"未知 surface 类型"场景，但流程图未体现。
3. Error Handling Paths 表中的"Surface 规则文件缺失"、"just 版本 < 1.4.0"未在流程图中体现。

---

## Dimension 3: Flow Completeness — 165 / 200

### Flow steps describe complete business process — 55 / 70

init-justfile 和 run-tests 流程步骤完整。混合项目生成流程已新增（第136-161行）。Probe 重试规格已统一（第164-175行）。

**-15**:
1. Goal 3（surface-key 统一迁移）和 Goal 4（Task 模型扩展）在 Flow Description 中仍无对应流程。9项迁移变更点的执行顺序、依赖关系、验证步骤未描述。迭代1已指出，未修复。
2. 混合项目生成流程缺少失败场景（某个 surface 的规则文件缺失时，是跳过该 surface 还是整体失败？）。

### Data flow documented — 62 / 70

跨组件数据流表（第110-124行）7步传递链清晰。Fallback 链明确。

**-8**:
1. 数据流表步骤2的 JSON 格式是预期接口定义，但未标注为"契约定义"。`forge surfaces` CLI 需新建，实际输出格式是否与此一致？
2. 旧任务 frontmatter 中的 `scope: frontend` 反序列化到新 struct SurfaceKey 字段时的兼容映射未说明。步骤4只说 "Go struct 序列化"，未提反向兼容。
3. 步骤5 "从 Task struct 读取" 未包含编排流程步骤1的 fallback 到 CLI 路径。

### Exception handling and edge cases covered — 48 / 60

Error Handling Paths 表覆盖5个场景，每个含恢复提示，符合 BIZ-error-reporting-002。

**-12**:
1. 多 surface 项目中某个 surface 编排失败是否影响其他 surface——未回答。
2. init-justfile user-customized 冲突处理不在 Error Handling Paths 表中。
3. config.yaml surfaces 格式错误（类型错误、结构错误）的检测和行为未记录。

---

## Dimension 4: User Stories — 173 / 200

### Coverage: one story per target user — 48 / 50

Forge 插件开发者 → Story 3、4；Forge 用户 → Story 1、2。

**-2**: 混合项目的 run-tests 编排涉及新的配方命名规则，CI 配置可能需更新，缺少 CI/CD 角色的故事。

### Format correct (As a / I want / So that) — 48 / 50

4个故事全部使用标准格式。

**-2**: Story 3 "I want to 将 surface-key 值域...迁移为用户自定义" 偏向技术实现而非用户价值。迭代1已指出，未修改。

### AC per story (Given/When/Then) — 42 / 50

Story 1-3 使用 Given/When/Then。Story 4 后两条 AC 缺少独立 Given/When/Then。

**-8**: Story 4 第63-64行 "And 旧任务的 scope 字段通过 GetSurfaceKey() 兼容访问"、"And quality-gate fix-task 从失败文件路径自动推断 surface-key/type" 缺少独立场景的 Given/When/Then 结构。后一条是一个独立验证场景，应拆分。

### AC verifiability & boundary coverage — 35 / 50

exit code 语义已定义（第176-186行），可构造测试用例。Probe 参数已量化。

**-15**:
1. HARD-GATE "同一编排周期"定义模糊（第33行）。如何区分"同一"和"新"周期？进程重启后是新周期吗？
2. Story 1 缺少无 surface 配置场景的 AC。
3. Story 4 "通过 GetSurfaceKey() 兼容访问"测试实现细节（Go 函数名），而非行为。
4. Story 4 "从失败文件路径推断 surface-key/type" 缺少文件路径不属于任何 surface 的边界场景。

---

## Dimension 5: Scenario Completeness — 116 / 150

### End-to-end scenario coverage — 48 / 60

核心场景覆盖：init-justfile surface 感知、run-tests 编排、surface-key 迁移、Task 模型扩展、混合项目生成。

**-12**:
1. quality-gate fix-task 端到端流程——Functional Specs #9 提及但无场景描述。
2. 旧系统迁移（frontend/backend → 自定义 surface-key）升级路径——迭代1已指出，未改善。
3. `forge surfaces` CLI 本身的端到端行为规格——作为需新建的前置依赖，其输入/输出/错误规格缺失。

### Implicit assumptions surfaced — 28 / 40

改善：`forge surfaces` CLI 已声明前置依赖，probe 参数已量化，exit code 已定义。

**-12**:
1. config.yaml surfaces 字段的精确 YAML 格式——第116行写 `map<string,string>`，但用户应写 `surfaces: {admin-panel: web}` 还是其他格式？未提供示例。
2. `.forge/test-state.json` 的 schema——仅提文件名（第260行），未定义结构。
3. 混合项目 "按依赖序并行启动"（第147行）——依赖序如何确定？固定规则还是用户配置？
4. 第147行 "并行" 与 just 语法 `(dep1 dep2)` 的串行依赖语义不一致。

### Business-rules consistency — 40 / 50

Exit code 对齐 BIZ-error-reporting-001。HARD-GATE 与 exit code 关系已明确（第187行）。

**-10**:
1. BIZ-task-lifecycle-001 定义 7 状态状态机。"从源任务继承 surface-key"（第63行）发生在哪个状态转换点未说明。
2. 第98行 "每步检查退出码：exit 0 继续" —— 但编排表和流程图显示 test exit 0 后也执行 teardown。第98行只说"继续"，未提 teardown。与编排表和流程图不一致。
3. 规则文件的编排序列格式未定义——Functional Specs #2 说 "编排序列由规则文件定义"，但规则文件的内容格式/schema 未指定。

---

## Dimension 6: Edge Case Coverage — 70 / 100

### Error paths documented — 32 / 40

Error Handling Paths 表覆盖5个关键场景，每个含恢复提示。

**-8**:
1. config.yaml surfaces 格式错误的检测和行为未覆盖。
2. init-justfile user-customized 冲突不在 Error Handling Paths 表中。
3. 多 surface 项目部分编排失败的错误隔离未描述。
4. `forge surfaces` CLI 输出格式异常（非 JSON、字段缺失）的防御性处理未覆盖。

### Boundary conditions covered — 20 / 35

已覆盖：无 surface 配置、probe 重试上限（3次/90秒）。

**-15**:
1. 空 surfaces map（`surfaces: {}`）的处理。
2. Surface type 大小写敏感性（"Web" vs "web"）。迭代1已指出，未修复。
3. 两个 surface 共享同一文件路径前缀的 longest-prefix-match 冲突。
4. 混合项目并行 dev server 端口冲突。迭代1已指出，未修复。
5. init-justfile --force-regenerate 覆盖 user-customized 配方的边界行为。

### Failure recovery described — 18 / 25

已描述：test-state.json 恢复、teardown 幂等、git revert。

**-7**:
1. Surface-key 迁移后回滚——旧数据（scope 字段值）的兼容性恢复路径未描述。迭代1已指出，未修复。
2. 混合项目部分编排失败的恢复策略未描述。
3. `forge surfaces` CLI 执行失败的重试策略——Error Handling Paths 表标记 exit 1（retryable），但重试次数和间隔未定义。

---

## Dimension 7: Scope Clarity — 91 / 100

### In-scope items are concrete deliverables — 32 / 35

In Scope 列出具体文件路径、函数名、字段名。

**-3**: "死代码清理：extractTestTypeArg()、genScriptBases"——缺乏明确度（是删除还是替换调用方？删除后的回归风险？）。迭代1已指出，未改善。

### Out-of-scope explicitly lists deferred items — 27 / 30

Out of Scope 6 项。

**-3**: "变更 `forge-cli/internal/cmd/quality_gate.go` 或 `testrunner` 的 Go 代码"在 Out of Scope 中，但 In Scope "quality-gate fix-task 从失败文件路径推断 surface-key/type" 涉及 quality-gate 行为变更。变更是在 skill 文档层面还是 Go 代码层面未澄清。

### Scope consistent with functional specs and user stories — 32 / 35

11 个变更点与 In Scope 基本对应。

**-3**:
1. **严重矛盾**: In Scope 第59行声明 `forge surfaces` CLI "需新建"，Out of Scope 第75行写 "新增 forge CLI 命令"。`forge surfaces` 就是新增的 forge CLI 命令，二者直接矛盾。
2. Story 3 AC "无 surfaces 配置的项目行为不变" 的验证手段无对应 Scope 条目。

---

## Phase 4: Blindspot Attacks

### [blindspot] 1: `forge surfaces` CLI 在 In Scope 和 Out of Scope 之间存在直接矛盾

> In Scope 第59行: "**前置依赖：`forge surfaces` CLI 命令** — 当前状态：需新建。"
> Out of Scope 第75行: "新增 forge CLI 命令"

第59行说 `forge surfaces` CLI 需新建且为本特性前置依赖，第75行说 "新增 forge CLI 命令" 在 Scope 之外。`forge surfaces` 本身就是新增的 forge CLI 命令，二者不可调和。下游 agent 无法判断 `forge surfaces` 是否在本特性工作范围内。

**Must improve**: 解决矛盾。方案一：从 Out of Scope 移除 "新增 forge CLI 命令"，将 `forge surfaces` 纳入本特性 Scope。方案二：将 In Scope 标注改为 "前置条件：须在独立特性/迭代中先行完成"，并定义本特性无此 CLI 时的测试策略。

### [blindspot] 2: 混合项目并行 dev server 端口冲突无解决方案

> 第147行: "dev：按依赖序并行启动所有 dev server（如 api 先于 web）"

文档假设多个 dev server 可并行运行，但未解决端口冲突。且 "并行" 一词与 just 语法 `(dep1 dep2)` 的串行依赖语义不一致——just 的配方依赖列表是串行执行的。

**Must improve**: 澄清并行/串行语义（对齐 just 实际行为），说明端口分配策略或声明为用户责任。

### [blindspot] 3: Goal 3/4 迁移流程缺失，下游 agent 无法执行

> In Scope 第60-68行: 列出 9 项 surface-key 迁移的具体变更点
> Flow Description: 只描述了 init-justfile 和 run-tests 流程

9项迁移变更涉及 prompt.go 重写、Go struct 变更、模板迁移、CLI 命令更新等，但 Flow Description 部分无对应的实施流程。下游 agent 接到这些变更点后，无法推导执行顺序、依赖关系和验证步骤。迭代1已指出，未修复。

**Must improve**: 为 surface-key 迁移添加实施流程，至少包含变更顺序、依赖关系、每步验证方法。

### [blindspot] 4: Story 4 AC 仍耦合实现细节

> prd-user-stories.md 第63行: "And 旧任务的 scope 字段通过 GetSurfaceKey() 兼容访问"

GetSurfaceKey() 是 Go 函数名，属于实现细节。迭代1已指出，未修改。"兼容访问"语义模糊——是旧字段名作为别名保留？还是旧值自动映射到新字段？

**Must improve**: 重写为行为验证："Given 一个包含 scope: frontend 字段的旧任务文件，When run-tests 读取该任务，Then 该任务被正确识别并使用默认编排策略执行。"

### [blindspot] 5: `forge surfaces` CLI 行为规格缺失

> 第59行: "此命令接受文件路径参数，返回 longest-prefix-match 的 surface-key 和 surface-type"

作为需新建的前置依赖，仅一句话描述行为。缺少：(1) 输入参数规格（绝对/相对路径？文件/目录？）；(2) 输出格式规格（JSON 结构、字段名）；(3) 多个 surface 匹配时的优先级规则；(4) 无匹配时的 exit code 和输出；(5) config.yaml 未配置 surfaces 时的行为。

**Must improve**: 为 `forge surfaces` CLI 添加完整行为规格，或引用独立设计文档。

---

## Summary

| Dimension | Score | Max |
|-----------|-------|-----|
| Background & Goals | 92 | 100 |
| Flow Diagrams | 135 | 150 |
| Flow Completeness | 165 | 200 |
| User Stories | 173 | 200 |
| Scenario Completeness | 116 | 150 |
| Edge Case Coverage | 70 | 100 |
| Scope Clarity | 91 | 100 |
| **Total** | **842** | **1000** |

### Top 5 Issues to Address for Iteration 3

1. **`forge surfaces` CLI Scope 矛盾** (Dimension 7, Scope Clarity): In Scope "需新建" vs Out of Scope "新增 forge CLI 命令"直接矛盾，必须解决。

2. **Goal 3/4 迁移流程缺失** (Dimension 3, Flow Completeness): 9项 surface-key 迁移变更无实施流程描述，下游 agent 无法推导执行顺序。连续两轮未修复。

3. **错误路径未完全覆盖** (Dimension 6, Edge Case Coverage): config.yaml 格式错误、CLI 输出异常、多 surface 部分失败等关键场景仍缺失。

4. **混合项目端口冲突与并行语义** (Dimension 5, Scenario Completeness): "并行启动"与 just 串行语义不一致，端口冲突无解决方案。

5. **Story 4 AC 耦合实现细节** (Dimension 4, User Stories): GetSurfaceKey() 是 Go 函数名，AC 应验证行为而非内部方法。连续两轮未修复。
