---
created: "2026-05-27"
author: "forge-brainstorm"
status: Draft
---

# Proposal: Unified Template Engine for Prompt and Task Templates

## Problem

Forge 的模板渲染在三个包中各自使用 `strings.ReplaceAll` 做占位符替换，无法表达条件逻辑。所有条件化行为都靠脆弱的后处理函数实现。此外，部分模板仍使用过时的 `scope` 概念（已迁移到 `surface-key`/`surface-type`），且同一目录下的模板 frontmatter 规范不统一。

### Evidence

**三套独立的渲染+后处理机制**：

| 渲染引擎 | 位置 | 后处理函数 | 模板数 | 过时概念 |
|---------|------|-----------|--------|---------|
| `renderTemplate()` | `pkg/prompt/prompt.go` | `cleanTemplateOutput()` | 21 | — |
| `ApplyVars()` | `pkg/template/template.go` | `injectSurfaceFrontmatter()`（替换 `surface-key: ""` 字面值） | 2 | `surface-key: ""` 硬编码 |
| `renderBody()` | `pkg/task/autogen.go` | `removeLineContaining()`, `removeSection()` | 12 | `{{SCOPE}}` 占位符 |

**过时概念残留**：

| 过时概念 | 当前替代 | 残留位置 |
|---------|---------|---------|
| `{{SCOPE}}` 占位符 | `surface-key`/`surface-type` | `pkg/task/data/` 中 4 个模板：doc-consolidate, test-gen-contracts, test-gen-journeys, test-run |
| `BodyContext.Scope` 字段 | `SurfaceKey`/`SurfaceType` | `pkg/task/autogen.go` |

**frontmatter 不一致**：`pkg/task/data/` 中 `record-*.md`（6 个）有 YAML frontmatter，其余 12 个模板没有。`pkg/prompt/data/` 中 4 个 doc 模板缺少 `SURFACE_KEY` header 声明。这种不一致导致模板变量与 Go struct 字段的对应关系无法自动验证——当 35 个模板分散在 3 个包中时，新增或修改模板时极易遗漏字段映射，且人工审查无法保证覆盖完整性。

**目录结构分散**：模板相关代码分散在三个包中（`pkg/task/`、`pkg/prompt/`、`pkg/template/`），其中 `pkg/template/` 仅有 2 个模板文件和一个简单函数，概念上属于 `pkg/task/`（生成任务文件），但物理上独立存在。`pkg/task/data/` 混装了两类用途不同的模板（autogen body + record），仅靠文件名前缀区分。

**Skill/Command 层概念残留**：`init-justfile/templates/mixed.just` 仍使用 `frontend`/`backend` scope 参数；`FrontmatterData` struct 的 `scope` 字段虽已标记 deprecated 但未清理。

**后处理函数脆弱假设**：
- `cleanTemplateOutput()` 的 `isLabelWithEmptyValue()` 对标签名有空格限制（`strings.Contains(before, " ")` → false）
- `removeSection()` 通过 `## Heading` 前缀匹配删除段落，模板内容变更可能破坏匹配
- `injectSurfaceFrontmatter()` 用 `strings.Replace` 将模板中硬编码的 `surface-key: ""` 和 `surface-type: ""` 字面值替换为推断值

### Urgency

`task-pipeline-precision` 提案已批准，需要在 prompt 模板中实现 complexity 条件分支。如果不引入模板引擎，将叠加第三层后处理 hack（标记注释）。同时 `pkg/task/autogen.go` 的 `renderBody()` 仍使用过时的 `{{SCOPE}}` 概念，与已迁移到 `surface-key`/`surface-type` 的系统其余部分不一致。现在迁移可以一次性统一三个渲染引擎、清理过时概念、建立可扩展的条件化基础设施。

## Proposed Solution

将三个包的渲染引擎统一替换为 Go 标准库 `text/template`，同时清理过时概念、统一 frontmatter 规范。

**核心变更**：
1. 三个包的渲染引擎统一迁移到 `text/template`：`pkg/prompt`（21 模板）、`pkg/template`（2 模板）、`pkg/task/autogen`（12 模板）
2. 占位符语法从 `{{X}}` 迁移到 `{{.X}}`（CamelCase dot-notation，与现有 record 模板一致）
3. 条件段落用 `{{if .Field}}...{{end}}` 声明
4. 移除三套后处理函数（`cleanTemplateOutput` 条件逻辑、`injectSurfaceFrontmatter`、`removeLineContaining`/`removeSection`）
5. `{{SCOPE}}` 替换为 `{{.SurfaceKey}}`/`{{.SurfaceType}}`，`BodyContext.Scope` 字段重命名
6. Surface 推断失败时任务创建硬性报错
7. 目录重组：`pkg/template/` 合并入 `pkg/task/`，`pkg/task/data/` 按类别拆分为 `templates/` 和 `records/`，`pkg/prompt/data/` 重命名为 `templates/`
8. 所有模板文件统一添加 metadata frontmatter（type、category、variables），record 模板的输出 frontmatter 移入 body
9. Skill/Command/Agent 层对齐最新概念，清除 scope 残留

### Innovation Highlights

这不是技术创新——Go `text/template` 是标准库，且已在 `pkg/task/data/`（typed-task-records）中使用。本提案的价值在于统一：将散布在两个包中的三种后处理机制替换为一个声明式模板系统。

**行业参照**：Helm charts、Go CLI 工具（Hugo、goreleaser）普遍使用 `text/template` 管理条件化模板。Forge 的场景更简单（无嵌套模板、无管道），但模式相同。

### User-Facing Behavior

用户无可见功能变化——模板输出的 markdown 内容与当前行为等价。唯一例外：quality gate 创建 fix/cleanup 任务时，若 surface 无法推断（无 surfaces 配置），任务创建将失败而非生成空 surface 字段的任务文件。错误信息引导用户运行 `forge surfaces detect` 配置 surfaces。

## Requirements Analysis

### Key Scenarios

1. **Prompt 模板条件渲染**：phase summary 有值时注入、无值时整个段落消失；coverage 仅 testable 类型渲染；surface key 为空时省略整行
2. **任务模板条件渲染**：surface 已推断时直接填充值并省略 Surface Inference 段落；surface 未推断时触发硬性报错
3. **Complexity 条件分支**：`{{if eq .Complexity "low"}}` 跳过 Step 1.5 spec-code scan；medium/high 渲染完整流程
4. **模板数据模型统一**：`pkg/prompt` 和 `pkg/template` 各定义一个模板数据 struct，暴露模板需要的全部字段
5. **目录按类别组织**：开发者在 `pkg/task/templates/` 找 autogen 和任务创建模板，在 `pkg/task/records/` 找 record 模板，在 `pkg/prompt/templates/` 找 prompt 模板——每个目录对应单一职责
6. **Metadata frontmatter 驱动校验**：新增模板时声明 variables 列表，启动时自动校验模板数据 struct 与声明一致。动机：35 个模板分散在 3 个包中，人工审查无法保证模板变量与 Go struct 字段的对应关系——`missingkey=error` 仅在运行时捕获执行阶段错误，metadata `variables` 声明使模板成为 self-documenting 并支持启动时声明-实现一致性校验，防止新增模板时遗漏字段

### Non-Functional Requirements

- **向后兼容**：现有 index.json 中无 complexity 字段的任务默认为 `medium`，模板渲染行为不变
- **字节等价迁移**：迁移后模板渲染路径的 prompt 输出与当前 `strings.ReplaceAll` + `cleanTemplateOutput()` 的输出在功能上等价（允许空白行差异）。不覆盖 Scope 中明确声明的行为变更路径（surface 硬性失败、`forge init` surface 配置集成）
- **零运行时性能退化**：`text/template.Parse()` 在 init 时执行一次，`Execute()` 在每次渲染时执行——性能与 `strings.ReplaceAll` 可比
- **可靠性权衡已确认**：`embed.FS` + `template.Parse()` 启动时解析是显意为之的权衡——启动时全量校验（fail-fast）优于运行时按需失败（fail-late）。当前所有 35 个模板经人工审查不含字面 `{{`，迁移后 `missingkey=error` + 零值 Execute 可在开发期捕获语法和字段名错误。若未来模板需包含字面 `{{`，使用 `` {{ "{{" }} `` 转义

### Constraints & Dependencies

- 模板文件通过 `//go:embed` 嵌入二进制，`text/template.Parse()` 须在启动时完成——这是 fail-fast 权衡：模板语法错误在启动时阻塞而非运行时按模板触发。当前 35 个模板经审查不含字面 `{{`，`missingkey=error` 配置确保字段名错误在开发期暴露
- `pkg/template` 的 `ApplyVars()` 同时被 `forge task add` 命令和 quality gate hook 调用
- `pkg/prompt` 的 `renderTemplate()` 被 `forge prompt get-by-task-id` 调用
- `text/template` 的 `{{}}` 定界符与当前 `{{PLACEHOLDER}}` 语法兼容但需加 `.` 前缀；若模板需包含字面 `{{`，使用 `` {{ "{{" }} `` 转义（当前 35 个模板无此情况）

- **实施顺序**：`task-pipeline-precision` 先实施（在模板中引入 `<!-- IF NOT_LOW -->` 标记注释），本提案随后实施（将标记注释替换为 `{{if ne .Complexity "low"}}...{{end}}` 条件块）。两个提案不可并行——若本提案先实施，`task-pipeline-precision` 的标记注释将无后处理函数支撑；若同时实施，merge 冲突风险高

## Alternatives & Industry Benchmarking

### Industry Solutions

Go 生态中 `text/template` 是模板渲染的标准选择。Helm、Hugo、goreleaser、buffalo 等项目均使用它。无竞争力的第三方替代品。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零改动 | 三层后处理 hack 持续膨胀，conditionality 需求无法优雅表达 | Rejected: 技术债务积累 |
| 标记注释 + 后处理 | task-pipeline-precision | 改动小 | 第三层 hack，与 cleanTemplateOutput() 叠加 | Rejected: 仅推迟问题 |
| 仅迁移 pkg/template | — | 2 个模板，迁移快 | prompt 层 21 个模板仍用 hack，两套渲染机制共存 | Rejected: 不统一 |
| **统一 text/template** | Go 标准库 | 声明式条件、移除所有后处理 hack、已在代码库中使用 | 35 个模板文件需更新占位符语法 | **Selected: 统一清理，一劳永逸** |

## Feasibility Assessment

### Technical Feasibility

`text/template` 已在 `pkg/task/data/`（typed-task-records，task #2）中使用相同模式：`embed.FS` → `template.Parse()` → `template.Execute()`。迁移路径已被验证。

占位符语法迁移：`{{TASK_ID}}` → `{{.TaskID}}`。35 个 `.md` 文件的机械性替换，可用 sed 脚本批量完成。

### Resource & Timeline

- 35 个模板文件的占位符语法迁移（21 prompt + 2 task creation + 12 autogen）：机械性，约 2 小时
- `prompt.go` 重构：移除 `cleanTemplateOutput()` 全部四种条件逻辑，改为 `text/template` 渲染：约 2 小时
- `template.go` 重构：移除 `ApplyVars()` 和 `injectSurfaceFrontmatter()`（替换 `surface-key: ""`/`surface-type: ""` 字面值），改为 `text/template`：约 1 小时
- `autogen.go` 重构：移除 `renderBody()` + `removeLineContaining()`/`removeSection()`，改为 `text/template`：约 2 小时
- Surface 硬性约束 + `{{SCOPE}}` 两种模式迁移 + `BodyContext.Scope` 字段删除：约 1 小时
- 目录重组：`pkg/template/` 合并入 `pkg/task/`，`pkg/task/data/` 拆分，`pkg/prompt/data/` 重命名，`//go:embed` + import 路径更新：约 2 小时
- Metadata frontmatter：41 个模板文件添加 metadata frontmatter，record 模板输出 frontmatter 移入 body，模板加载器解析 metadata：约 2 小时
- Skill/Command/Agent 概念清理：`mixed.just` scope 清除，`FrontmatterData.Scope` 标记 deprecated，skill/command 审查：约 1 小时
- 测试：golden-file 对比确保渲染等价 + 新增目录/frontmatter/变量校验测试：约 2 小时
- `forge init` surface 配置集成：约 1 小时

总计约 15 个 coding task，适合 quick mode。

### Dependency Readiness

无外部依赖。`text/template` 是 Go 标准库。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| `strings.ReplaceAll` + 后处理足以应对条件化需求 | 5 Whys | Overturned: 每新增一个条件维度就需扩展后处理逻辑。`task-pipeline-precision` 的 complexity 分支是第三个条件维度（前两个是 phase summary 和 coverage），标记注释方案将成为第三层 hack |
| surface 推断失败时空值可接受 | XY Detection | Overturned: 用户明确要求硬性失败。空值导致 agent 在执行时自行推断 surface，浪费执行时间且结果不确定 |
| 模板占位符语法 `{{X}}` 可以保持不变 | Assumption Flip | Overturned: `text/template` 使用 `{{.X}}` 访问 struct 字段，`{{X}}` 是模板变量引用（需通过 `Define()` 声明）。直接使用 `{{.X}}` 更自然，且与 `pkg/task/data/` 中已有的 `text/template` 用法一致 |

## Scope

### In Scope

**pkg/prompt（Prompt 渲染层）**：
- `pkg/prompt/prompt.go` 重构：`renderTemplate()` 从 `strings.ReplaceAll` 迁移到 `text/template.Execute()`
- `pkg/prompt/data/` 下 21 个模板文件占位符语法迁移（`{{X}}` → `{{.X}}`）

- 移除 `cleanTemplateOutput()` 中的全部四种条件删除逻辑（空标签行、空 backtick 条件句、`just` 命令尾部空白、`<!-- IF NOT_LOW -->...<!-- END_IF -->` 段落块）。函数保留但仅做空白行塌陷——`text/template` 的 `{{-`/`-}}` 可消除单个 `{{if}}` 块周围的空行，但无法处理连续空行（多个条件块在同一位置省略时产生的多行空白），保留 Go 级空白行塌陷作为最终后处理步骤是必要的
- 为每个模板添加 `{{if}}` 条件块，替换当前靠后处理删除的段落：
  - `{{if .PhaseSummary}}...{{end}}` 替换 `If {{PHASE_SUMMARY}} is non-empty` 模式——所有 18 个模板中 PHASE_SUMMARY 出现在两个独立位置（标签行 + 条件指令行，相隔 10-30 行），需分别用两个独立的 `{{if .PhaseSummary}}` 块包裹，不可合并为单个块
  - `{{if .CoverageStrategy}}...{{end}}` 替换 coverage IMPORTANT 块
  - `{{if .SurfaceKey}}...{{end}}` 替换空 surface 标签删除
  - `just compile{{if .SurfaceKey}} {{.SurfaceKey}}{{end}}` 模式处理 `just compile {{SURFACE_KEY}}` 的尾部空格——当 SurfaceKey 为空时无尾随空格，有值时正确插入空格+值。应用于 6 个包含 `just compile` 命令的相关模板，无需依赖后处理清理尾部空格

  - `{{if ne .Complexity "low"}}...{{end}}` 替换 `<!-- IF NOT_LOW -->...<!-- END_IF -->` 条件段落块（4 个 coding 模板）
- 4 个 doc 模板（doc-consolidate, doc-drift, doc, doc-review）补齐 `SURFACE_KEY` header 声明
- `task-pipeline-precision` 的 complexity 条件分支改用 `{{if}}` 实现，替代标记注释方案
- 模板数据 struct 定义：`promptTemplateData`
  ```go
  type promptTemplateData struct {
      TaskID           string // 当前任务 ID
      TaskFile         string // 任务文件路径（如 tasks/001-fix-login.md）
      TaskCategory     string // 任务分类（fix/cleanup/doc/test 等），空字符串时省略分类段落
      FeatureSlug      string // 功能标识（如 "auth-login"），用于模板中功能相关引用
      PhaseSummary     string // phase summary 文本，空字符串时省略整个段落
      CoverageStrategy string // coverage 策略：testable 类型为具体策略文本，non-testable 为空（省略段落），cleanup 类型为特殊指令（"No coverage..."）。三态值由 Go 代码解析完整文本，模板侧仅判断 `{{if .CoverageStrategy}}`——coverage 指令文本的 source-of-truth 在 Go 代码而非模板中，这是有意的设计权衡：将分支逻辑集中在调用端，避免模板中增加二级条件判断
      CoverageTarget   string // coverage 目标值（如 "80%"），用于模板中 coverage 相关指令
      TestTypeArg      string // 测试类型参数（如 "contracts"/"journeys"），用于 just 命令组装
      SurfaceKey       string // surface key，空字符串时省略 surface 标签行
      SurfaceType      string // surface type
      Complexity       string // 任务复杂度：low/medium/high
  }
  ```

**pkg/template（任务创建层）**：
- `pkg/template/template.go` 重构：`ApplyVars()` 替换为 `text/template.Execute()`
- `pkg/template/data/` 下 2 个模板文件（coding.fix, coding.cleanup）占位符语法迁移
- 移除 `pkg/task/add.go` 的 `injectSurfaceFrontmatter()`——该函数用 `strings.Replace` 替换模板中硬编码的 `surface-key: ""` 和 `surface-type: ""` 字面值，统一后由模板渲染直接填充，无需后处理替换
- 模板条件化：surface 有值时渲染字段 + 省略 Surface Inference 段落
- 模板数据 struct 定义：`taskTemplateData`
  ```go
  type taskTemplateData struct {
      TaskName     string // 任务名称（如 "coding.fix"）
      SurfaceKey   string // surface key，空字符串时渲染空行（CLI 命令路径保留软性行为）
      SurfaceType  string // surface type
      TaskGoal     string // 用户描述的任务目标
      ScopeDescription string // 用户提供的任务作用域描述（非 deprecated scope 概念；这里是 task-level 语境描述，非 surface-level 的 SurfaceKey）
  }
  ```

**pkg/task/autogen（自动生成任务层）**：
- `pkg/task/autogen.go` 重构：`renderBody()` 从 `strings.ReplaceAll` 迁移到 `text/template.Execute()`
- `pkg/task/data/` 下 12 个非 record 模板占位符语法迁移（`{{X}}` → `{{.X}}`）

- `{{SCOPE}}` 两种使用模式的迁移——统一替换为 `SurfaceKey`/`SurfaceType`，彻底消除 `scope` 概念：
  - **段落级**（test-gen-contracts, test-gen-journeys）：`{{SCOPE}}` 作为 `## Scope` 段落标题下的独立块 → 整段用 `{{if .SurfaceKey}}...{{end}}` 包裹，内容引用 `{{.SurfaceKey}}`
  - **行内值**（doc-consolidate, test-run 及其他使用点）：`{{SCOPE}}` 作为行内占位符 → 替换为 `{{.SurfaceKey}}`

- `BodyContext.Scope`（`[]string` 类型）字段直接删除——其语义已由 `SurfaceKey`/`SurfaceType` 覆盖，无需过渡字段
- 移除 `removeLineContaining()` 和 `removeSection()` 后处理函数，条件逻辑改用 `{{if}}`
- 模板数据 struct 定义：`autogenTemplateData`
  ```go
  type autogenTemplateData struct {
      TaskID             string // 任务 ID
      TaskType           string // 任务类型标识（如 "test-gen-contracts"）
      FeatureSlug        string // 功能标识（如 "auth-login"），用于模板中功能相关引用
      Mode               string // 生成模式（如 "create"/"update"），用于条件化模板行为
      SurfaceKey         string // surface key，用于模板中 {{.SurfaceKey}} 行内替换和 {{if .SurfaceKey}} 条件判断
      SurfaceType        string // surface type
      SurfaceTypes       string // 多 surface type 预格式化字符串（如 "frontend, backend"），供模板直接输出
      AcceptanceCriteria string // 验收标准预格式化文本，空字符串时省略相关段落
      DocTaskCriteria    string // 文档任务标准预格式化文本，空字符串时省略相关段落
  }
  ```


**renderTemplate() 中 TASK_CATEGORY 注入迁移**：
- `pkg/prompt/prompt.go` 第 136-137 行的 `TASK_CATEGORY` 字符串拼接逻辑迁移到 `promptTemplateData.TaskCategory` 字段
- submit-task 路由中的 `TASK_CATEGORY` 段落改用 `{{if .TaskCategory}}...{{end}}` 条件块渲染

**Surface 硬性约束**：

- `quality_gate.go` 的 `addSingleFixTask()` 中 `inferSurface()` 失败时返回错误（硬性失败）——这是行为变更（当前行为为静默空值），而非纯重构。变更影响限于 `addSingleFixTask()` 调用路径，无其他调用路径受影响
- `forge task add` 命令路径保留软性行为——`injectSurfaceFrontmatter()` 移除后，该命令的 surface 注入逻辑由模板渲染统一处理，推断失败时使用空字符串而非报错，保持 CLI 命令的容错性
- `forge init` 中集成 surface 配置步骤（`forge surfaces detect`），确保新项目初始化即具备 surfaces 配置，从源头避免硬性失败场景

**Frontmatter 规范统一**：
- 所有模板文件添加 metadata frontmatter，格式统一为：
  ```yaml
  ---
  type: test.run        # 任务类型常量
  category: test        # 任务分类（coding/doc/test/eval/validation/gate/record）
  variables:            # 模板使用的变量列表（用于校验）
    - FeatureSlug
    - SurfaceKey
  ---
  ```
- Record 模板重构：输出 frontmatter（status、started 等）从 metadata 区移入 body 区，由 Go template 渲染。模板文件结构变为：metadata frontmatter → body（含 `---` 输出 frontmatter + 内容）
- metadata frontmatter 由模板加载器在 `template.Parse()` 前剥离，不参与渲染输出
- autogen body 模板的 metadata frontmatter 替代原有的 `<!-- body-only -->` 注释方案——frontmatter 的 `variables` 字段本身即可表达"此模板期望哪些变量"的语义
- 模板校验：`ValidateTemplates()` 使用 `variables` 字段与模板数据 struct 的反射字段做交叉校验，确保声明与实现一致

**目录重组**：
- `pkg/template/` 整包合并入 `pkg/task/`：其 2 个模板文件（coding-fix.md、coding-cleanup.md）迁入 `pkg/task/templates/`，代码逻辑迁入 `pkg/task/` 对应文件
- `pkg/task/data/` 拆分为：
  - `pkg/task/templates/` — autogen body 模板 + 迁入的任务创建模板（共 14 个）
  - `pkg/task/records/` — record 模板（6 个），文件名去掉 `record-` 前缀（`coding.md`、`doc.md` 等）
- `pkg/prompt/data/` 重命名为 `pkg/prompt/templates/`（21 个模板）
- 所有 `//go:embed` 路径和 `templatePath()` 函数相应更新
- 删除空包 `pkg/template/`

**Skill/Command/Agent 概念对齐**：
- `plugins/forge/skills/breakdown-tasks/rules/scope-to-surface-key.md` — 删除（迁移已完成，不再需要迁移文档）
- `plugins/forge/skills/init-justfile/templates/mixed.just` — 清除 `frontend`/`backend` scope 参数，改用 surface-aware recipe 命名（`<surface-key>-<verb>`）
- `plugins/forge/agents/task-executor.md` — 审查并确认无 scope 残留
- `plugins/forge/commands/*.md` — 审查并确认无 scope 残留
- `FrontmatterData` struct 的 `scope` 字段（`yaml:"scope"`）标记为 `deprecated` 并添加 `// Deprecated:` 注释，`CheckLegacyScope()` 保留用于迁移检测

### Out of Scope

- 模板内容精简（属于 `slim-task-prompt-templates` 提案）
- 任务文档模板（`plugins/forge/skills/*/templates/`）——由 LLM agent 渲染
- Record 模板引擎变更——已使用 `text/template`，仅做 metadata frontmatter 和目录重组
- 新增模板文件或合并/拆分现有模板（同名模板保持分离，任务描述和执行指令是不同关注点）
- `prompt-template-audit` 提案中的其他优化建议（如双重提交、Hard Rules 命名）
- `mixed.just` 之外的其他 init-justfile 模板重构
- `ARCHITECTURE.md` 中 `scope` 引用的文档更新（属于文档同步，非本提案交付物）

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|

| 占位符迁移遗漏导致模板渲染时出现 `{{.UnknownField}}` 错误 | M | H | 启动时 `ValidatePromptTemplates()` 和 `ValidateAutogenTemplates()` 使用 `template.Option("missingkey=error")` 配置 + 对零值 struct 执行 `Execute()` 到 `io.Discard`，编译期捕获字段拼写错误（`Parse()` 仅验证语法，无法检测字段名） |
| `text/template` 对 nil 指针解引用 panic | L | H | 模板数据 struct 使用值类型（string）而非指针；所有字段零值为空字符串（在模板中 `{{if .Field}}` 对空字符串为 false），无需 nil 检查 |
| Surface 硬性失败阻塞无 surfaces 配置的项目 | M | H | 错误信息包含 `forge surfaces detect` 命令指引；`forge init` 中集成 surface 配置步骤（见 Scope + SC），确保新项目不触发此场景。已有项目升级时需手动运行 `forge surfaces detect` |
| 迁移后 prompt 输出与当前行为有细微差异（空白行、格式） | M | L | Golden-file 测试对比迁移前后输出，允许空白行差异但要求内容等价 |

| `task-pipeline-precision` 的 complexity 分支实现需要同步调整 | M | M | `task-pipeline-precision` 先实施（引入 `<!-- IF NOT_LOW -->` 标记注释），本提案随后实施（将标记替换为 `{{if ne .Complexity "low"}}` 条件块）。顺序不可颠倒 |
| `BodyContext.Scope` 重命名影响 `BuildIndex()` 调用链 | M | M | `Scope` 仅在 `autogen.go` 内部消费，调用链为 `BuildIndex()` → `renderBody()`，影响面可控 |
| 35 个模板批量迁移可能引入格式错误（`text/template` 要求 `{{}}` 严格配对） | M | M | 迁移后每个模板执行 `template.Parse()` 验证语法，配合 golden-file 测试确保内容正确 |
| 迁移后发现细微渲染差异，需回退 | L | H | 回滚策略：每个包独立迁移并独立提交，保留旧函数（`ApplyVars`、`cleanTemplateOutput`、`renderBody`）在新函数通过全量 golden-file 测试前不删除。若需回滚，revert 对应包的提交即可，不影响其他包 |
| Surface 硬性失败行为变更导致生产环境阻塞 | L | H | 回滚策略：硬性失败逻辑独立封装为 `requireSurfaceInference()` 函数（`quality_gate.go`），可通过单行 revert 恢复为软性行为。硬性失败与模板引擎迁移解耦——即使回滚硬性失败，模板引擎迁移不受影响 |
| 目录重组导致 `//go:embed` 路径和 import 路径大面积变更 | M | M | 先完成引擎迁移和占位符替换（验证渲染等价），再执行目录重组（纯机械性文件移动+import 更新），两步分离降低风险 |
| Record 模板双 frontmatter 结构（metadata + 输出）增加理解成本 | L | M | 在 record 模板文件头部添加注释解释结构，并在 `docs/` 中记录约定 |
| `mixed.just` scope 参数清除影响现有项目 justfile 兼容性 | L | L | 仅影响 `forge init` 生成的模板，不修改用户现有 justfile；`init-justfile` 已有 surface-aware 路径，`mixed.just` 是旧路径的清理 |

## Success Criteria

- [ ] `pkg/prompt` 的 `renderTemplate()` 使用 `text/template.Execute()` 渲染，不再使用 `strings.ReplaceAll`
- [ ] `pkg/template` 的 `CreateTaskMarkdown()` 使用 `text/template.Execute()` 渲染，不再使用 `ApplyVars()` 和 `injectSurfaceFrontmatter()`
- [ ] `pkg/task/autogen.go` 的 `renderBody()` 使用 `text/template.Execute()` 渲染，不再使用 `strings.ReplaceAll`

- [ ] `cleanTemplateOutput()` 仅保留空白行塌陷逻辑，移除全部四种条件删除逻辑（空标签行、空 backtick 条件句、`just` 命令尾部空白、`<!-- IF NOT_LOW -->...<!-- END_IF -->` 段落块）。`just compile` 命令的尾部空格由模板级 `{{if .SurfaceKey}}` 条件处理
- [ ] `removeLineContaining()` 和 `removeSection()` 从 `autogen.go` 中移除，条件逻辑由模板 `{{if}}` 块处理
- [ ] 35 个模板文件中无 `{{PLACEHOLDER}}` 格式（全部为 `{{.Placeholder}}` 格式）
- [ ] `pkg/task/data/` 中无 `{{SCOPE}}` 残留
- [ ] `BodyContext` struct 中无 `Scope` 字段（已删除，语义由 `SurfaceKey`/`SurfaceType` 覆盖）
- [ ] `{{if .PhaseSummary}}` 条件块正确渲染：有值时同时注入标签行和条件指令行（两个独立 `{{if}}` 块），无值时两处均消失——覆盖所有 18 个使用 PHASE_SUMMARY 的模板

- [ ] `{{if .CoverageStrategy}}` 条件块正确渲染：testable 类型渲染 coverage 指令，non-testable 类型无 coverage 段落，cleanup 类型渲染特殊 "No coverage..." 指令（三态：空/策略/特殊指令）

- [ ] Surface 推断失败时 `addSingleFixTask()` 返回错误而非创建空 surface 的任务文件（`quality_gate.go` 路径硬性失败）
- [ ] `forge task add` 命令路径保持软性行为——surface 推断失败时使用空字符串而非报错
- [ ] `forge init` 中集成 surface 配置步骤（`forge surfaces detect`），新项目初始化时具备 surfaces 配置
- [ ] `forge prompt get-by-task-id` 输出与迁移前功能等价（golden-file 对比，允许空白行差异）

- [ ] `ValidatePromptTemplates()` 和 `ValidateAutogenTemplates()` 使用 `missingkey=error` 选项 + 零值 struct `Execute()` 验证，确保无字段拼写错误
- [ ] 4 个 doc 类型 prompt 模板（doc-consolidate, doc-drift, doc, doc-review）补齐 `SURFACE_KEY` header 声明

- [ ] `renderTemplate()` 中的 `TASK_CATEGORY` 拼接逻辑迁移到 `promptTemplateData.TaskCategory` 字段，模板中用 `{{if .TaskCategory}}...{{end}}` 渲染

- [ ] `<!-- IF NOT_LOW -->...<!-- END_IF -->` 标记从 4 个 coding 模板中移除，替换为 `{{if ne .Complexity "low"}}...{{end}}`

- [ ] `BodyContext.Scope` 字段已删除，`{{SCOPE}}` 统一替换为 `{{.SurfaceKey}}`，模板中无 `{{range}}` 循环
- [ ] `{{if .SurfaceKey}}` 条件块在 prompt 模板中正确渲染：有值时显示 surface 标签行，无值时整行省略
- [ ] 任务创建模板（coding-fix.md、coding-cleanup.md）包含 `surface-key` 和 `surface-type` 字段行，surface 有值时渲染字段 + 省略 Surface Inference 段落
- [ ] 所有模板文件（35 + 6 record = 41 个）包含 metadata frontmatter（type、category、variables 字段）
- [ ] Record 模板文件结构正确：metadata frontmatter + body（含输出 frontmatter 渲染）
- [ ] 模板加载器在 `template.Parse()` 前正确剥离 metadata frontmatter，渲染输出不含 metadata
- [ ] autogen body 模板中无 `<!-- body-only -->` 注释（已被 metadata frontmatter 替代）
- [ ] `pkg/template/` 包已删除，其模板文件迁入 `pkg/task/templates/`，代码逻辑合并到 `pkg/task/` 对应文件
- [ ] `pkg/task/data/` 已拆分为 `pkg/task/templates/`（14 个模板）和 `pkg/task/records/`（6 个模板），record 文件无 `record-` 前缀
- [ ] `pkg/prompt/data/` 已重命名为 `pkg/prompt/templates/`
- [ ] `ValidateTemplates()` 使用 `variables` 字段与模板数据 struct 做正向交叉校验：每个模板 metadata 中声明的 variable 必须在对应 struct 中有匹配字段（missingkey=error 确保反向覆盖）
- [ ] `mixed.just` 模板中无 `scope` 参数和 `frontend`/`backend` 值
- [ ] `FrontmatterData` struct 的 `scope` 字段包含 `// Deprecated:` 注释
- [ ] `scope-to-surface-key.md` 已删除
- [ ] `task-executor.md` 和 `commands/*.md` 中无 deprecated scope 字段引用（prose 用法如 "well-scoped" 不算）
- [ ] 所有 `//go:embed` 路径和 import 路径已更新，`go build ./...` 通过

```
consistency_check_result:
  status: pass
  pairs_checked: 21
  conflicts_found: 0
```

## Next Steps

- Proceed to `/quick-tasks` to generate tasks from this proposal
