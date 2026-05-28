---
created: "2026-05-27"
updated: "2026-05-28"
author: "forge-brainstorm"
status: Approved
---

# Proposal: 精简任务 Prompt 模板 + Frontmatter 元数据结构重构

## Problem

Forge 的 15 个任务 prompt 模板包含大量非指令内容——注释、解释性描述、冗长的角色定义——这些内容不指导 agent 行为，只增加 token 消耗并稀释指令清晰度。同时 task-executor agent 的 Execution Protocol 存在步骤冗长、逻辑重叠的问题。

此外，在 unified-template-engine MR 中为所有 41 个模板添加了统一的 metadata frontmatter（type/category/variables），当前的 frontmatter 使用扁平的 variables list 列举所有模板变量，未区分元数据字段和内容字段，导致：
- frontmatter 对模板调用方的契约声明不清晰——无法判断哪些是关键元数据、哪些是普通内容
- phaseSummary 等大段内容字段被声明在 variables list 中，但其实际上是正文的一个条件性 section
- 三套模板（prompt/task/record）的 frontmatter 结构缺乏统一的语义分层

### Evidence

内容冗余的量化分析：

| 冗余类别 | 出现文件数 | 每处损失 | 总冗余 |
|---------|-----------|---------|-------|
| HTML 注释 | 1 | 4 行 | 4 行 |
| Step 2 解释性描述 | 5 (gen-* / test-run / verify) | 1-2 行 | ~7 行 |
| 冗长角色描述 | 10 (coding.*, gate, doc) | 1 行 | ~10 行 |
| CODING_PRINCIPLES 解释性冗余 | 5 (coding.*) | 每原则 2-5 行 | ~50 行 |
| AC 验证块冗余 | 9 (coding.*, gate, doc) | 每处 ~12 行可缩至 ~4 行 | ~70 行 |
| Record Fields 描述性文字 | 9 | ~3 行可缩至 ~1 行 | ~20 行 |
| task-executor Execution Protocol | 1 | 11 步可合并为 8 步 | ~30 行 |

总计：约 **190 行** 非指令冗余。注意并非每个任务都加载全部 200 行——单个 coding.* task 的模板包含约 80-100 行冗余（AC 验证块为主体），gate/doc/test 类更少（约 20-40 行）。200 行为全模板集上界，操作以单模板实际冗余为准。

**Token 估算**（以 Claude Sonnet tokenizer 为参考）：不同类型行的 token 密度差异显著——空行 ~1 token，纯文本指令 ~8-12 tokens，代码块约束 ~15-25 tokens，JSON 示例 ~20-40 tokens。coding.* 模板 ~80-100 行冗余，按加权平均 token 密度估算约 1200-1500 tokens/task。按 daily task 量估算：保守（10 个 coding.* task × 800 tokens）~ 激进（20 个 coding.* task × 1100 tokens），**每日 token 节省范围为 8K-22K tokens（输入侧）**，月度约 170K-450K tokens。数值为近似，精确测算见 SC8。

Frontmatter 冗余的定性分析：

| 模板类型 | 数量 | 当前 variables 数量 | 其中元数据字段数 | 其中内容字段数 | 精简后 variables 数量 |
|---------|------|-------------------|----------------|--------------|-------------------|
| prompt 模板 | 21 | ~8-12 个 | 4-7 个 | 3-4 个 | 2-4 个 |
| task 模板 | 14 | ~10-14 个 | 4-5 个 | 6-8 个 | 6-8 个 |
| record 模板 | 6 | ~16-22 个 | 3 个 | 13-19 个 | 13-19 个 |

frontmatter 重构不改变模板运行期行为，仅改善机器校验能力和人类可读性。具体收益量化：(a) 编译时字段校验——分组声明使 `validateMetadataVariables` 能在启动时捕获字段拼写错误（如 `TaskId` 误写为 `Taskid`），防止运行时静默变量替换失败导致的 task 执行异常；(b) 预计调试时间可减少（基于模板字段拼写错误在调试中的占比估算，待部署后验证）——当前模板变量拼写错误需在运行时通过输出缺失字段排查，分组校验在启动时即报错；(c) 消除一类运行时错误——变量名与 Go struct 字段对齐校验可防止 `reflect.FieldByName` 返回空值导致的 panic。

### Urgency

- 每个 task 执行都在消耗这些冗余 token，日积月累规模可观
- 清晰的 prompt 减少 agent 误解和执行偏差
- Prompt 精简是持续优化的一部分，目前已有 prompt-template-audit 等基础，可以在此基础上推进
- unified-template-engine 刚完成 frontmatter 添加，趁热打铁完成结构优化比后续再改更高效

## Proposed Solution

**双线并行**：

1. **内容精简**：保持现有模板独立，在每个文件内部删除非指令内容，将模糊描述改为清晰指令。不抽取公共模块，不改变现有分类体系。（与原提案一致）

2. **Frontmatter 重构**：将 metadata frontmatter 从扁平的 `variables` list 改为语义分组结构。新增 `identity`/`context`/`conditional` 三个分组组。`phaseSummary` 从 frontmatter 移除，改为正文独立 `## PhaseSummary` section。保留双 frontmatter 结构（task 和 record 模板的 rendered frontmatter 不变）。Go 侧解析和校验逻辑同步更新。保持三套模板（prompt/task/record）使用统一的语义分组规范。

**执行顺序约束**：双线不可真正同时修改相同的 21 个 prompt 模板——否则先完成的修改会使 SC-Pre 基线失效。执行顺序为：(a) 内容精简先执行，完成后以精简后的模板建立 SC-Pre 基线（token/行数快照）；(b) Frontmatter 重构随后执行，以内容精简后的状态作为其基线。两线各自独立验证，避免基线交叉污染。task 模板和 record 模板仅涉及 frontmatter 重构，可与 (b) 同步进行。

**分组判定规则**：
- **identity**：唯一标识模板实例的字段——缺少该字段导致模板无法定位具体实例（如 `TaskID`、`ID`、`Title`）
- **context**：运行时上下文字段——提供执行环境信息但不控制条件段落的显示/隐藏（如 `SurfaceKey`、`Complexity`、`FeatureSlug`）
- **conditional**：以 `{{if .X}}` 控制正文中整段内容显示/隐藏的字段（如 `CoverageStrategy`——模板正文中存在 `{{if .CoverageStrategy}}...{{end}}` 条件块）
- **variables**：其余所有在正文中以 `{{.X}}` 插值但不符合以上三类的模板变量

### Innovation Highlights

本方案不是技术创新，而是对现有 prompt 的"清理"。核心原则是"prompt 是指令，不是文档"——删掉所有不能直接指导 agent 行动的文字。

Frontmatter 重构的设计灵感来自 OpenAPI 规范的参数分类（path/query/header/body）——将模板变量按其语义角色（身份标识 / 运行上下文 / 条件分支控制 / 内容数据）分类声明，让 frontmatter 既是机器校验的契约也是人类阅读的目录。

**行业参照：** 本方案的设计哲学与以下行业实践一致：
- **LangChain Prompt Templates** 在模板中区分"指令（instructions）"与"上下文（context）"，推荐仅将直接影响模型行为的文本保留在系统 prompt 中，解释性描述移至外部文档。
- **Anthropic Prompt Engineering Guide** 强调"show, don't just tell"——通过示例约束行为而非通过自然语言角色描述；本方案中的 AC 验证块精简（保留 REQUIRED 指令、删除展开说明）遵循同一原则。
- **OpenAI GPTs Instructions** 模式的演变方向也是删除冗余的系统 prompt 装饰，改用精确的祈使句指令。
- **Kubernetes YAML 的 apiVersion/kind/metadata 分层**——强制性的元数据 vs 可选的 spec 内容分离，是 frontmatter 分组声明的行业参考模型。

### User-Facing Behavior

本提案对用户（task 执行者）无可见功能变化——用户提交 task 后看到的是相同的执行流程、相同的输出结构、相同的结果质量。唯一的可观测差异是 token 消耗降低（详见 Token 估算），在计费侧表现为每次 task 执行的成本下降。agent 行为层面的"无行为变更"经 SC2 轨迹对比验证。

Frontmatter 重构对用户完全透明——frontmatter 在渲染前被剥离，不影响最终的 agent prompt 内容。

## Requirements Analysis

### Key Scenarios

#### 内容精简场景

1. **coding-feature / coding-enhancement / coding-fix / coding-cleanup / coding-refactor** 五个核心模板：
   - 角色描述从自然语言改为祈使句
   - CODING_PRINCIPLES 去掉举例和解释，保留核心约束
   - AC 验证块从 ~12 行精简到 ~4 行
   - Step 2 的实现说明保留，只去修饰性语言

2. **gate / doc** 模板：
   - 角色描述精简
   - AC 验证块精简

3. **test-run / test-gen-scripts / test-gen-contracts / test-gen-journeys / test-verify-regression** 模板：
   - Step 2 中的 "This generates X from Y..." 解释性描述删除
   - 角色描述精简

4. **code-quality-simplify / validation-code / validation-ux** 模板（共 3 个，约 30-50 行/个）：
   - 角色描述精简（同 coding-* 模式）
   - 无 AC 验证块和 CODING_PRINCIPLES——冗余集中在角色描述和框架性说明行
   - 三者合计精简约 22 行

5. **task-executor agent**：
   - Execution Protocol 步骤合并（11 步 → ≤8 步），具体合并方案：
     - 合并 4/5/6（prompt 获取/解析/注入）→ 1 步"获取并注入 prompt"
     - 合并 Retry Strategy 与 Complex Error Pause Flow → 1 步"异常处理"（去重：两者共享"暂停并通知用户"语义）
     - 输出格式从多行描述改为紧凑单行摘要
   - 合并后预估 8 步：初始化 → 任务验证 → 获取并注入 prompt → 执行任务 → 输出格式化 → 结果记录 → 异常处理 → 完成通知

AC 验证块逐行分析：

| 行类型 | 典型数量 | 处理策略 | 压缩后行数 |
|--------|---------|---------|-----------|
| AC:REQUIRED 指令 | 2-3 行 | 保留——义务级别最高，必须执行 | 2-3 行 |
| AC:STRONGLY 指令 | 1-2 行 | 保留——义务级别低于 REQUIRED，但仍是强制建议（建议而非命令） | 1-2 行 |
| 指令展开说明 | 3-5 行 | 合并至指令行 | 0 |
| 场景举例 | 0-2 行 | 删除 | 0 |
| 格式装饰 | 3-4 行 | 保留必要空行 | 1-2 行 |

~12 行 → ~5 行（58%），AC:REQUIRED 与 AC:STRONGLY 的区分被保留——两者在 agent 执行语义上对应不同的遵循强度，合并会导致 agent 对 AC 优先级层次的感知模糊。

CODING_PRINCIPLES 逐原则分析：

| 原则条目 | 行数 | 功能判定 | 处理策略 |
|---------|------|---------|---------|
| 原则 1: 纯指令行 | 1 行 | 核心约束 | 保留 |
| 原则 1: 行为示例/边界说明 | 2-5 行 | 约束边界演示——非核心指令，但作为"注意力分段锚点"在密集指令排列中提供视觉分隔和注意力重置作用 | 每原则保留 1 个代表性示例（视觉分隔功能）+ 压缩边界说明为 1 行概括 |
| 原则 2: 纯指令行 | 1 行 | 核心约束 | 保留 |
| 原则 2: 反例/边界说明 | 2-5 行 | 约束边界演示 | 每原则保留 1 个代表性示例 + 压缩边界说明为 1 行概括 |
| 超原则通用说明（如作用域声明） | 1 行 | 元指令 | 保留 |

~50 行 → ~25 行（50%），每原则保留 1 行指令 + 1 行边界概括 + 1 个代表性示例。

Record Fields 逐字段分析：

| 行类型 | 典型数量 | 处理策略 |
|--------|---------|---------|
| 字段名 + 值（如 `## Output\n{...}`） | 1 行 | 保留 |
| 字段用途说明（如 "This field describes..."） | 1-2 行 | 删除——字段名自解释 |
| 格式示例/占位符展开 | 1-2 行 | 删除——嵌入实际值即可 |

~3 行 → ~1 行（66%），字段名和值保留。

#### Frontmatter 重构场景

6. **prompt 模板 frontmatter**：所有 21 个 prompt 模板的 metadata frontmatter 改为分组声明格式：

```yaml
# 当前（扁平 variables list）
---
type: coding.feature
category: coding
variables:
  - TaskID
  - TaskFile
  - TaskCategory
  - FeatureSlug
  - PhaseSummary
  - CoverageStrategy
  - CoverageTarget
  - TestTypeArg
  - SurfaceKey
  - SurfaceType
  - Complexity
---

# 新版（语义分组声明，PhaseSummary 从 frontmatter 移除）
---
type: coding.feature
category: coding
identity:
  TaskID: true
context:
  TaskFile: true
  FeatureSlug: true
  SurfaceKey: true
  SurfaceType: true
  Complexity: true
conditional:
  CoverageStrategy: true
variables:
  - TaskCategory
  - CoverageTarget
  - TestTypeArg
---
```

7. **PhaseSummary 正文独立 section**：PhaseSummary 不再存在于 frontmatter，而是在模板正文中作为独立 Markdown 二级标题 section：

```markdown
TASK_ID: {{.TaskID}}
TASK_FILE: {{.TaskFile}}
{{if .SurfaceKey}}SURFACE_KEY: {{.SurfaceKey}}{{end}}
{{if .PhaseSummary}}
## PhaseSummary
{{.PhaseSummary}}
{{end}}
```

新 section 位置：紧接 TASK_ID/TASK_FILE 行之后、角色描述/CODING_PRINCIPLES 之前。

8. **task 模板 frontmatter**：

```yaml
# 当前
---
type: coding.cleanup
category: coding
variables:
  - ID
  - Title
  - Priority
  - EstimatedTime
  - Description
  - SourceTaskID
  - SurfaceKey
  - SurfaceType
  - SourceFiles
  - TestScript
  - TestResults
  - ScopeDescription
---
---  # rendered frontmatter 不变

# 新版
---
type: coding.cleanup
category: coding
identity:
  ID: true
  Title: true
context:
  SurfaceKey: true
  SurfaceType: true
variables:
  - Priority
  - EstimatedTime
  - Description
  - SourceTaskID
  - SourceFiles
  - TestScript
  - TestResults
  - ScopeDescription
---
---  # rendered frontmatter 不变
```

9. **record 模板 frontmatter**：

```yaml
# 当前
---
type: record
category: record
variables:
  - Status
  - Started
  - Completed
  - ...
---
---  # rendered frontmatter

# 新版
---
type: record
category: record
identity:
  TaskID: true
  TaskTitle: true
  Status: true
variables:
  - Started
  - Completed
  - TimeSpent
  - Summary
  - FilesCreatedFormatted
  - ...
---
---  # rendered frontmatter 不变
```

### 指令分类标准

在逐类型分析中已经隐式使用了分类框架，现将其显式声明为方法论基础：

**三类指令的操作性定义：**

| 类别 | 定义 | 示例 | 精简处理策略 | 方法论依据 |
|------|------|------|-------------|-----------|
| **A. 正面指令** | 告诉 agent 应该做什么的祈使句或模态动词句（must/should/need to） | "Keep the existing behavior unchanged" / "You must include tests" | 保留。仅删除修饰性前置语（"Note that..." → 保留核心动词） | 可直译的 agent 行为规则，删除即丢失功能 |
| **B. 负面约束** | 告诉 agent 不应该做什么的否定句或禁止性表述 | "Do NOT remove format markers" / "You must not skip tests" | 保留。仅删除双重否定和展开说明 | 同 A——agent 需要知道禁令边界 |
| **C. 行为示范** | 通过正例/反例展示期望行为而非直接指令 | CODING_PRINCIPLES 中的 "Good: `{...}` Bad: `{...}`" | 按原则保留 1 个代表性示例。见 CODING_PRINCIPLES 逐原则分析 | 作用于 LLM 的示范学习（few-shot）路径，与指令路径正交；全部删除则失去该调节手段 |

### 隐式结构依赖审计

修改前置步骤：创建结构依赖矩阵，识别模板结构性特征是否被消费组件以字符串匹配方式依赖。

**结构依赖矩阵（非完整版，实施前须逐项核实）：**

| 结构性特征 | 示例 | task-executor agent | prompt.go 解析逻辑 | 测试脚本/CI 工具 |
|-----------|------|--------------------|------------------|----------------|
| 章节标题（`## X`） | `## Output`, `## Instructions` | 通过标题解析输出结构 | embed FS 按文件名索引，不解析标题 | 无依赖 |
| 标记前缀（`AC:`） | `AC:REQUIRED`, `AC:STRONGLY` | 通过标记前缀识别约束类型 | 无依赖 | 无依赖 |
| 格式约定（分隔线、缩进） | `---`, `- ` 列表 | 无依赖 | 无依赖 | 无依赖 |
| 占位符格式（`{{VAR}}`） | `{{TASK_ID}}`, `{{TASK_INSTRUCTION}}` | 通过占位符插入运行期变量 | prompt.go 执行字符串替换 | 无依赖 |
| CRITICAL 块标记 | `**CRITICAL**` | 通过标记识别高优先级指令 | 无依赖 | 无依赖 |
| frontmatter `---` 分隔 | `---` | 不消费 frontmatter | parseMetadataFrontmatter 剥离 | 无依赖 |

**分析结论：** 当前模板的结构性特征主要为 task-executor agent 的内部遍历逻辑消费（通过标题和前缀语义识别指令类别），而非通过字符串匹配方式解析。prompt.go 仅通过 embed FS 按文件名加载完整内容，不做结构解析。因此精简后结构变形不会导致运行时组装断裂。frontmatter 分组重构仅修改 metadata frontmatter 的 YAML 结构，`parseMetadataFrontmatter` 剥离 frontmatter 的行为不变，prompt 渲染不受影响。

### Non-Functional Requirements

- 精简后模板的指令覆盖必须与精简前等价（不能遗漏 agents 需要知道的信息）
- 所有 task-executor 的行为不发生变化
- Frontmatter 重构后 Go 侧校验逻辑必须正确验证所有分组中的字段
- 分组校验规则表（定义三类模板各自的校验 struct 和分组分配）：

| 模板类型 | 校验目标 Go Struct | identity 分组 | context 分组 | conditional 分组 | variables 分组 |
|---------|-------------------|-------------|------------|----------------|--------------|
| prompt 模板 | `promptTemplateData` | TaskID | TaskFile, FeatureSlug, SurfaceKey, SurfaceType, Complexity 等 | CoverageStrategy 等 | 模板正文占位符（非结构体字段的模板变量） |
| task 模板 | `TemplateData` | ID, Title | SurfaceKey, SurfaceType 等 | — | Priority, Description, SourceFiles 等 |
| record 模板 | `RecordTemplateData` | TaskID, TaskTitle, Status | — | — | Started, Completed, Summary 等 |
- `parseMetadataFrontmatter()` 解析 frontmatter 的行为必须保持向后兼容——不带 frontmatter 的模板文件应继续正常解析

### Constraints & Dependencies

- 内容精简涉及的模板文件位于 `forge-cli/pkg/prompt/templates/*.md`，由 `prompt.go` 通过 embed FS 加载。修改模板不影响 Go 代码，只需修改 .md 文件
- Frontmatter 重构涉及的模板：
  - `forge-cli/pkg/prompt/templates/*.md`（21 个 prompt 模板）
  - `forge-cli/pkg/task/templates/*.md`（14 个 task 创建模板）
  - `forge-cli/pkg/task/records/*.md`（6 个 record 模板）
- Frontmatter 重构涉及的 Go 代码：
  - `forge-cli/pkg/prompt/metadata.go`（解析逻辑和数据结构）
  - `forge-cli/pkg/prompt/metadata_test.go`（单元测试）
- task-executor agent 位于 `plugins/forge/agents/task-executor.md`
- **行格式不变性约束**：`TASK_FILE`、`TASK_ID`、`SURFACE_KEY` 行在模板正文中的格式（`KEY: {{.Value}}`）和位置（正文首行区域）不可变更——`prompt.go` 使用 `strings.Replace` 对这些行做字符串替换，格式变更将导致静默失败（不报错但变量未被替换）
- **`forge task index` 兼容性约束**：`forge task index` 通过 `task/frontmatter.go` 的 `ParseFrontmatter()` 解析 task .md 文件（rendered frontmatter），读取 `FrontmatterData` 结构体字段（id, title, type, surface-key 等）生成 index.json。本提案的 frontmatter 重构仅影响模板文件的 **metadata frontmatter**（第一个 `---` 块，渲染前被 `parseMetadataFrontmatter` 剥离），不影响 task/record 模板的 **rendered frontmatter**（第二个 `---` 块，由 `forge task index` 消费）。两套 frontmatter 解析路径完全独立：metadata frontmatter 由 `prompt/metadata.go` 解析，rendered frontmatter 由 `task/frontmatter.go` 解析（使用 `gopkg.in/yaml.v3`）。本提案不修改 `FrontmatterData` 结构体、`ParseFrontmatter()` 函数、或 `build.go` 的 `BuildIndex()` 逻辑

## Alternatives & Industry Benchmarking

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| 分层模板组合 | LangChain PromptTemplate / Vercel AI SDK | 语义分离（instruction/tool/context 分层），单层修改不影响其他，改一处影响所有 | 需重构模板分类体系并修改 `prompt.go` 加载逻辑，与"不改后端代码"约束冲突；对 15 个文件引入抽象层，改动面大于收益 | Rejected: 架构约束否决 |
| 引入 DSL 生成 | 模板引擎模式 | 声明式模板定义，通过编译生成最终 prompt，压缩逻辑集中在 DSL 层 | 需要增加 DSL 定义文件、解析器、编译管线，对 15 个小模板引入完整工具链成本过高——模板改动频次低（月级而非天级），DSL 抽象层在小规模场景下维护负担超过收益 | Rejected: 模板规模小、变更频次低，DSL 工具链成本不合理 |
| 什么都不做 | — | 零风险 | token 持续浪费、指令不够清晰 | Rejected: 成本太低 |
| 抽取公共模块 | DRY 模式 | 修改一处同步所有模板 | 需要改 `prompt.go` 逻辑，且被用户否决 | Rejected: 不满足就地要求 |
| **就地精简（内容）** | Forge 现有风格 | 零架构变更，每模板独立修改，风险隔离 | 每个文件都要改 | **Selected: 简单直接** |
| **分组声明 frontmatter** | Kubernetes metadata 模型 | 清晰的契约声明，改善可读性和校验 | 需改 41 个模板 + Go 解析逻辑 | **Selected: 趁热打铁** |

## Feasibility Assessment

### Technical Feasibility

纯文本编辑 + Go 解析逻辑修改，技术风险可控但需注意以下：

`parseMetadataFrontmatter` 是单文件内的行级 YAML 解析器，当前通过逐行扫描 `---` 分隔符和 `key: value` 模式提取字段。扩展分组支持需要引入嵌套层级的状态跟踪（当前在 `identity:` 键内部需识别缩进层级），复杂度从"简单扩展"升级为"中等"——需要跟踪缩进深度以确定字段归属哪个分组。

**建议方案：** 引入 `gopkg.in/yaml.v3` 标准库替代手工行级解析，将 `TemplateMetadata` 直接映射为 YAML struct tag 反序列化。此方案：(1) 消除手工缩进状态跟踪的 bug 风险；(2) 天然支持向后兼容（YAML 反序列化时缺失字段为零值）；(3) 代码量预计减少。若引入第三方依赖不可接受，则需在现有行级解析器中增加缩进栈状态机，预估增加 ~50 行解析逻辑。

### Resource & Timeline

内容精简部分为 1 次编码任务（约 0.5 天）。附加制品与验证工作同原提案。

Frontmatter 重构部分：Go 端修改约 0.5 天（metadata.go + metadata_test.go），41 个模板 frontmatter 批量修改约 1 天（含 PhaseSummary section 迁移），验证约 0.5 天。总计约 2 天。

### Dependency Readiness

前置条件：本次 brainstorm 输出的 proposal 通过。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| "我需要保留角色描述让 agent 理解上下文" | Assumption Flip | 角色描述中的自然语言对 LLM 行为的影响可由祈使句替代——通过实施后的行为等价验证确认。 |
| "每个模板独立意味着不需要关注跨模板一致性" | XY Detection | 用户确认了核心流程重复是允许的，跨模板一致性不是问题，不需要抽取公共模块。 |
| "PhaseSummary 应该保留在 frontmatter 的 variables 中" | Need Gate | PhaseSummary 是大段内容文本，不属于元数据。移除到正文独立 section 更合理。 |
| "frontmatter 的分组应该避免改 Go 代码" | Assumption Flip | 改 Go 解析逻辑不可避免（支持分组解析和校验），但属于受控变更，metadata.go 是单文件解析器。 |
| "每原则保留 1 个代表性示例足以维持行为示范效果" | Parameter Tuning | few-shot 示例数量（当前设为 1 个/原则）是可调参数，非经验证的最优值。SC1 保留率验证需确认压缩后的原则仍保持功能等价——若 SC2 轨迹对比显示特定原则的行为漂移，则该原则的示例数量应上调至 2 个。初始值 1 个基于"注意力分段锚点"假设，实际效果由 SC2 实证检验。 |

## Scope

### In Scope

#### 内容精简
- 修改 `forge-cli/pkg/prompt/templates/` 下全部 15 个模板文件
- 修改 `plugins/forge/agents/task-executor.md`
- 删除 HTML 注释
- 删除 Step 2 解释性描述
- 精简角色描述（自然语言 → 祈使句）
- 精简 AC 验证块（~12 行 → ~4 行）
- 精简化 Record Fields（去掉引导性描述）
- 精简 CODING_PRINCIPLES（去掉举例和解释）
- 精简 task-executor Execution Protocol（合并步骤）

#### Frontmatter 重构
- 修改 `forge-cli/pkg/prompt/metadata.go`：
  - `TemplateMetadata` 结构体新增 `Identity`, `Context`, `Conditional` 字段，并新增 `AllFields()` 方法返回全部分组字段的并集（用于向后兼容的模板变量列表查询）
  - 新 `TemplateMetadata` 结构体定义：
    ```go
    type TemplateMetadata struct {
        Type       string
        Category   string
        Identity   map[string]bool  // 唯一标识模板实例的字段
        Context    map[string]bool  // 运行时上下文字段
        Conditional map[string]bool // 条件渲染控制字段
        Variables  []string         // 其余模板变量
    }
    func (m *TemplateMetadata) AllFields() []string {
        // 返回 Identity、Context、Conditional 的 key 并集 + Variables 列表
        // 用于向后兼容：旧代码通过 Variables 列表查询所有可用字段
    }
    ```
  - `parseMetadataFrontmatter()` 解析逻辑支持分组字段
  - `validateMetadataVariables()` 扩展为校验所有分组中的字段
  - 保持向后兼容——无 frontmatter 或旧格式 `variables` list 可继续解析
- 修改 `forge-cli/pkg/prompt/metadata_test.go`：
  - 新增分组解析测试用例
  - 新增分组校验测试用例
  - 新增向后兼容性测试（旧格式 variables list）
- 修改全部 41 个模板文件的 metadata frontmatter：
  - 21 个 prompt 模板：添加 `identity`/`context`/`conditional` 分组
  - 14 个 task 模板：添加 `identity`/`context` 分组
  - 6 个 record 模板：添加 `identity` 分组
- 将 PhaseSummary 从所有 21 个 prompt 模板的 frontmatter 变量中移除
- 在所有 21 个 prompt 模板正文中添加 `## PhaseSummary` 独立 section
- 修改 `forge-cli/pkg/prompt/prompt.go` 中 `phaseSummaryLine` 的格式渲染——仅去除 `PHASE_SUMMARY:` 前缀（受控修改，不涉及 Synthesize() 主逻辑）

### Out of Scope

- 不抽取公共模块文件
- 不修改 `prompt.go` 的 `Synthesize()` 主逻辑
- 不改变 rendered frontmatter（task/record 模板的第二个 `---` 块不变）
- 不新增/删除模板文件
- 不改动模板占位符的语法和格式（`{{.X}}` 语法不变），但允许新增 PhaseSummary 条件块作为正文结构变更
- 不改动 Spec Authority Enforcement 逻辑结构
- 不改动 Hard Rules / CRITICAL 块的逻辑
- 不改变 task-executor 的执行语义（输出质量和行为结果不变），但允许步骤合并和输出格式精简
- 不改动 `promptTemplateData`、`TemplateData`、`RecordTemplateData` 等 Go struct

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 精简过度导致 agent 遗漏关键行为 | Low | High | 功能快照清单（见原提案 Risk 1） |
| 多个模板同步修改，跨模板不一致 | Medium | Medium | 以 coding-feature.md 为基准，diff 校验 |
| 现有测试基础设施无法检测 prompt 层行为漂移 | Medium | High | 功能快照清单 + SC2 trial run 轨迹对比 |
| prompt 变更为有状态修改，合入后发现影响需要回滚但无标准化流程 | Medium | Medium | 分批独立提交，CI 观察期，git revert |
| 精简导致信息密度提升，关键指令在密集文本中显著性降低 | Medium | Medium | 指令行占比评估 + 注意力锚点恢复 |
| 合并后 prompt 变更的长期累积行为效应 | Low | Medium | 周期性轨迹重放检测 |
| **frontmatter 分组解析向后兼容性**——现有 `variables` list 格式的模板文件（如尚未更新的文件）应能正常解析 | Low | High | `parseMetadataFrontmatter` 保持容错：优先解析分组字段，若无分组则 fallback 到旧 variables list 解析。测试覆盖新旧两种格式。 |
| **分组字段名与 Go struct 字段名对齐**——分组内的字段名直接使用 Go struct 的 PascalCase 名（如 `TaskID`、`SurfaceKey`），与 `reflect.FieldByName` 大小写敏感匹配保持一致 | Low | Medium | 分组内字段名必须使用 Go struct 原始 PascalCase 名，确保 `validateMetadataVariables` 通过反射校验时直接匹配。分组 label 本身（identity/context/conditional）使用小写。 |
| **PhaseSummary 迁移后 template 渲染不一致**——正文中新增 `{{if .PhaseSummary}}## PhaseSummary\n{{.PhaseSummary}}{{end}}` 后，当前 Go 渲染逻辑不改变，但需确认包裹条件块的格式与原有 `{{if .PhaseSummary}}{{.PhaseSummary}}{{end}}` 的行为等价 | Low | Low | 本质上是将单行条件替换为多行条件块，text/template 的 `{{if}}...{{end}}` 行为不变。仅需视觉确认渲染后格式正确。 |
| **双线并行的级联回退风险**——frontmatter 重构（stream b）发现模板正文需要重新基线时，可能使内容精简（stream a）已建立的 SC-Pre 基线失效 | Low | Medium | stream (b) 发现问题需要重新基线时，仅在 stream (a) 的 SC-Pre 尚未建立时回退；若 SC-Pre 已建立则在新基线上继续。两线通过 SC-Pre 建立状态作为回退门禁。 |

## Success Criteria

**主要指标（保留率）和次要指标（token/行数）的双层结构：** 保留率为首要校验门禁，token 压缩为主要效率指标，行数压缩为次要参考——当保留率不达标时禁止合并，token 压缩不达标可接受，行数不达标不单独处理。

### 前置基线测量（修改前）

- **[SC-Pre] 修改前 token 和行数基线：** 修改开始前，对所有 In Scope 范围内的模板文件及 task-executor 执行 tokenize，记录每文件的当前 token 数和行数。**输出物：** 文件 `eval/baseline-token-counts.json`。

### 功能快照清单创建标准

功能快照清单用于 SC1 的保留率验证，创建标准如下：

- **节点粒度**：以模板中的"语义节点"为最小记录单元。语义节点定义为一行或连续几行具有独立功能目的的文本块，例如：一条 AC 指令、一个 CODING_PRINCIPLES 原则条目、一个 Record Field 声明、一个章节标题。
- **分类字典**：每个节点标注类别标签，分类字典为：`{instruction: 正面指令, constraint: 负面约束, example: 行为示范, format: 格式约定, metadata: 元数据声明}`。
- **清单格式**：每模板生成一个 JSON 文件，结构为 `[{id, category, summary, sourceLine}]`，其中 `sourceLine` 记录原始行号用于定位。
- **[SC-FM-Pre] Frontmatter 基线确认：** 重构前对全部 41 个模板记录当前 frontmatter 结构（字段总数、variables list 数量），作为重构后验证的基线。**输出物：** 文件 `eval/frontmatter-baseline.json`。

### 功能保留（首要门禁）

- [SC1] 功能约束保留率 **100%**——每个模板修改后，对照功能快照清单逐项比对，所有指令/约束/格式节点保留率为 100% 方可合并。

### 行为等价性

- [SC2] 模板精简后，agent 执行相同 task 的行为无可见差异。轨迹一致性 ≥ 90%（容差：步骤顺序因 LLM 生成随机性导致的非功能性差异）视为通过。差异分类标准：

  | 差异类型 | 判定 | 示例 |
  |---------|------|------|
  | 步骤顺序互换 | 非功能性 | 前置检查在 lint 前执行 vs lint 后执行 |
  | 措辞/表述差异 | 非功能性 | "Running tests..." vs "Executing test suite..." |
  | 中间日志行数差异 | 非功能性 | 输出 3 行摘要 vs 5 行摘要，但信息等价 |
  | 缺失功能步骤 | **功能性** | 精简前执行 lint 检查，精简后跳过 lint |
  | 参数语义变更 | **功能性** | 原指令 "--strict mode"，精简后变为 "--loose mode" |
  | 约束遗漏 | **功能性** | 原指令要求 "must include tests"，精简后该约束消失 |

  任何功能性差异出现即视为 SC2 不通过，需回溯修正对应模板。

### 结构验证

- [SC3] CODING_PRINCIPLES 在 5 个 coding-* 模板中保留全部核心约束指令。
- [SC4] Record Fields 在所有出现模板中保留字段名和值结构。
- [SC5] Step 2 解释性描述全部删除，通过 grep 确认无残留。

### 效率指标

- [SC6] 15 个模板文件 + task-executor 共减少 **≥1800 tokens**，行数参考指标为 **≥150 行**。
- [SC7] task-executor 的 Execution Protocol 步骤数从 11 步减少到 ≤8 步。

### Token 验证

- [SC8] 精简完成后对每个修改的模板文件执行实际 tokenize，与 SC-Pre 基线对比输出报告。

### Frontmatter 重构验证

- **[SC-FM-1] 字段分组覆盖率 100%**：重构后全部 41 个模板的 metadata frontmatter 已从扁平 `variables` list 迁移为分组声明格式。按模板类型分别定义迁移检测标准：
  - **prompt 模板（21 个）**：frontmatter 必须包含 `identity:` 和 `context:` 分组，`PhaseSummary` 不出现在任何分组或 variables 中
  - **task 模板（14 个）**：frontmatter 必须包含 `identity:` 分组（至少含 `ID` 字段），可不包含 `conditional:` 分组
  - **record 模板（6 个）**：frontmatter 必须包含 `identity:` 分组（至少含 `TaskID` 和 `Status` 字段），可不包含 `context:` 和 `conditional:` 分组
  - **迁移工具**：自动化脚本按类型分别检测分组存在性和必需字段，而非统一检测 `identity:` 键
- **[SC-FM-2] `parseMetadataFrontmatter` 向前兼容**：旧格式（仅 `variables` list，无分组）的模板文件能被正确解析。通过单元测试验证——输入旧格式 frontmatter，解析结果中的 `Variables` 字段与旧解析器一致。
- **[SC-FM-3] 分组校验通过**：`validateMetadataVariables`（或新的 `validateGroupedMetadata`）校验所有分组（`Identity`/`Context`/`Conditional`/`Variables`）中的字段名集合均存在于对应的 Go struct 中。通过 `ValidatePromptTemplates()` 启动时验证和单元测试双重保障。
- **[SC-FM-4] PhaseSummary 迁移完整**：所有 21 个 prompt 模板中：(a) frontmatter 的 `variables` list 不再包含 `PhaseSummary` 或 `PhaseSummary`；(b) 正文在 TASK_ID/TASK_FILE/SURFACE_KEY 行之后包含 `## PhaseSummary` 独立 section，包裹在 `{{if .PhaseSummary}}...{{end}}` 条件块中。通过两条独立 grep 确认：**(a) frontmatter 内无 PhaseSummary**：提取每个模板的首对 `---` 之间的内容（即 frontmatter section），在提取内容中 `grep -c "PhaseSummary"` 输出 0；**(b) 正文包含 PhaseSummary section**：提取每个模板第二个 `---` 之后的内容（即 body），在提取内容中 `grep -c "## PhaseSummary"` 输出 > 0。
- **[SC-FM-5] Grouping field naming 规范统一**：所有模板的分组 label 使用小写（`identity`/`context`/`conditional`），分组内的字段名使用 Go struct 的原始 PascalCase 名（如 `TaskID`、`SurfaceKey`），与 `reflect.FieldByName` 大小写敏感匹配保持一致。通过自动化脚本检测分组内字段名的 PascalCase 合规性。
- **[SC-FM-6] Rendered frontmatter 不变**：task 模板和 record 模板的第二个 `---` 块（rendered YAML frontmatter）在重构前后完全一致。通过 git diff 确认——第二个 `---` 块及其内容无变更。

## Next Steps

- Proceed to `/write-prd` to formalize requirements