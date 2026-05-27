---
created: "2026-05-27"
author: "forge-brainstorm"
status: Draft
---

# Proposal: Unified Template Engine for Prompt and Task Templates

## Problem

Forge 的模板渲染依赖 `strings.ReplaceAll` 做占位符替换，无法表达条件逻辑。所有条件化行为（空值省略、字段注入）都靠脆弱的后处理函数实现——`cleanTemplateOutput()` 删除空标签行和条件句，`injectSurfaceFrontmatter()` 替换硬编码的空值字段。随着条件化需求增长（complexity 分支、surface 差异化），这些后处理 hack 将持续膨胀。

### Evidence

| 后处理函数 | 位置 | 依赖的脆弱假设 |
|-----------|------|--------------|
| `cleanTemplateOutput()` | `pkg/prompt/prompt.go` | 按行匹配 `If \`\` is non-empty` 文本、`KEY:` 格式标签、`just ` 前缀命令 |
| `injectSurfaceFrontmatter()` | `pkg/task/add.go` | 模板必须硬编码 `surface-key: ""`，然后靠字符串替换覆盖 |
| 标记注释方案（计划中） | `task-pipeline-precision` | 用 `<!-- IF NOT_LOW -->...<!-- END_IF -->` HTML 注释包裹条件段落，由后处理删除 |

`cleanTemplateOutput()` 的 `isLabelWithEmptyValue()` 检测器甚至对标签名有空格限制（`strings.Contains(before, " ")` → 返回 false），这在未来添加新占位符时可能意外跳过清理。

### Urgency

`task-pipeline-precision` 提案已批准，需要在 prompt 模板中实现 complexity 条件分支。如果不引入模板引擎，将叠加第三层后处理 hack（标记注释）。现在迁移可以一次性清理技术债务，而非在三个后处理层之上继续堆叠。

## Proposed Solution

将 `pkg/prompt` 和 `pkg/template` 的渲染引擎统一替换为 Go 标准库 `text/template`。两个包已共享 `embed.FS` + `//go:embed` 模式，迁移路径清晰。

**核心变更**：
1. 占位符语法从 `{{X}}` 迁移到 `{{.X}}`（Go `text/template` 的 struct 字段访问语法）
2. 条件段落用 `{{if .Field}}...{{end}}` 声明
3. 移除 `cleanTemplateOutput()` 和 `injectSurfaceFrontmatter()` 的字符串匹配逻辑
4. Surface 推断失败时任务创建硬性报错

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

### Non-Functional Requirements

- **向后兼容**：现有 index.json 中无 complexity 字段的任务默认为 `medium`，模板渲染行为不变
- **字节等价迁移**：迁移后 prompt 输出与当前 `strings.ReplaceAll` + `cleanTemplateOutput()` 的输出在功能上等价（允许空白行差异）
- **零运行时性能退化**：`text/template.Parse()` 在 init 时执行一次，`Execute()` 在每次渲染时执行——性能与 `strings.ReplaceAll` 可比

### Constraints & Dependencies

- 模板文件通过 `//go:embed` 嵌入二进制，`text/template.Parse()` 须在启动时完成
- `pkg/template` 的 `ApplyVars()` 同时被 `forge task add` 命令和 quality gate hook 调用
- `pkg/prompt` 的 `renderTemplate()` 被 `forge prompt get-by-task-id` 调用
- `text/template` 的 `{{}}` 定界符与当前 `{{PLACEHOLDER}}` 语法兼容但需加 `.` 前缀

## Alternatives & Industry Benchmarking

### Industry Solutions

Go 生态中 `text/template` 是模板渲染的标准选择。Helm、Hugo、goreleaser、buffalo 等项目均使用它。无竞争力的第三方替代品。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零改动 | 三层后处理 hack 持续膨胀，conditionality 需求无法优雅表达 | Rejected: 技术债务积累 |
| 标记注释 + 后处理 | task-pipeline-precision | 改动小 | 第三层 hack，与 cleanTemplateOutput() 叠加 | Rejected: 仅推迟问题 |
| 仅迁移 pkg/template | — | 2 个模板，迁移快 | prompt 层 22 个模板仍用 hack，两套渲染机制共存 | Rejected: 不统一 |
| **统一 text/template** | Go 标准库 | 声明式条件、移除所有后处理 hack、已在代码库中使用 | 24 个模板文件需更新占位符语法 | **Selected: 统一清理，一劳永逸** |

## Feasibility Assessment

### Technical Feasibility

`text/template` 已在 `pkg/task/data/`（typed-task-records，task #2）中使用相同模式：`embed.FS` → `template.Parse()` → `template.Execute()`。迁移路径已被验证。

占位符语法迁移：`{{TASK_ID}}` → `{{.TaskID}}`。24 个 `.md` 文件的机械性替换，可用 sed 脚本批量完成。

### Resource & Timeline

- 24 个模板文件的占位符语法迁移：机械性，约 1 小时
- `prompt.go` 重构：移除 `cleanTemplateOutput()` 中的条件删除逻辑，改为 `text/template` 渲染：约 2 小时
- `template.go` 重构：移除 `ApplyVars()` 和 `injectSurfaceFrontmatter()`，改为 `text/template`：约 1 小时
- Surface 硬性约束：修改 `quality_gate.go` 的 `addSingleFixTask()` 报错逻辑：约 0.5 小时
- 测试：golden-file 对比确保渲染等价：约 2 小时

总计约 6-8 个 coding task，适合 quick mode。

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

- `pkg/prompt/prompt.go` 重构：`renderTemplate()` 从 `strings.ReplaceAll` 迁移到 `text/template.Execute()`
- `pkg/prompt/data/` 下 22 个模板文件占位符语法迁移（`{{X}}` → `{{.X}}`）
- 移除 `cleanTemplateOutput()` 中的条件删除逻辑（空标签行、空 backtick 条件句、`just` 命令尾部空白）。函数保留但仅做空白行塌陷
- 为每个模板添加 `{{if}}` 条件块，替换当前靠后处理删除的段落：
  - `{{if .PhaseSummary}}...{{end}}` 替换 `If {{PHASE_SUMMARY}} is non-empty` 模式
  - `{{if .CoverageStrategy}}...{{end}}` 替换 coverage IMPORTANT 块
  - `{{if .SurfaceKey}}...{{end}}` 替换空 surface 标签删除
- `pkg/template/template.go` 重构：`ApplyVars()` 替换为 `text/template.Execute()`
- `pkg/template/data/` 下 2 个模板文件占位符语法迁移
- 移除 `pkg/task/add.go` 的 `injectSurfaceFrontmatter()`——surface 值直接由模板渲染
- 模板条件化：surface 有值时渲染字段 + 省略 Surface Inference 段落
- `quality_gate.go` 的 `addSingleFixTask()` 中 `inferSurface()` 失败时返回错误（硬性失败）
- `task-pipeline-precision` 的 complexity 条件分支改用 `{{if}}` 实现，替代标记注释方案
- 模板数据 struct 定义：`promptTemplateData`（pkg/prompt）和 `taskTemplateData`（pkg/template）

### Out of Scope

- 模板内容精简（属于 `slim-task-prompt-templates` 提案）
- 任务文档模板（`plugins/forge/skills/*/templates/`）——由 LLM agent 渲染
- `task-executor.md` agent 定义修改
- Record 模板（`pkg/task/data/`）——已使用 `text/template`
- 新增模板文件或合并/拆分现有模板
- `prompt-template-audit` 提案中的其他优化建议（如双重提交、Hard Rules 命名）

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 占位符迁移遗漏导致模板渲染时出现 `{{.UnknownField}}` 错误 | M | H | 启动时 `ValidatePromptTemplates()` 扩展为对所有模板执行 `template.Parse()` + 检查所有字段引用，编译期而非运行期捕获遗漏 |
| `text/template` 对 nil 指针解引用 panic | L | H | 模板数据 struct 使用值类型（string）而非指针；所有字段零值为空字符串（在模板中 `{{if .Field}}` 对空字符串为 false），无需 nil 检查 |
| Surface 硬性失败阻塞无 surfaces 配置的项目 | M | H | 错误信息包含 `forge surfaces detect` 命令指引；在 `forge init` 中增加 surface 配置步骤 |
| 迁移后 prompt 输出与当前行为有细微差异（空白行、格式） | M | L | Golden-file 测试对比迁移前后输出，允许空白行差异但要求内容等价 |
| `task-pipeline-precision` 的 complexity 分支实现需要同步调整 | M | M | 本提案与 `task-pipeline-precision` 共同实施，complexity 条件直接用 `{{if}}` 实现 |

## Success Criteria

- [ ] `pkg/prompt` 的 `renderTemplate()` 使用 `text/template.Execute()` 渲染，不再使用 `strings.ReplaceAll`
- [ ] `pkg/template` 的 `CreateTaskMarkdown()` 使用 `text/template.Execute()` 渲染，不再使用 `ApplyVars()` 和 `injectSurfaceFrontmatter()`
- [ ] `cleanTemplateOutput()` 仅保留空白行塌陷逻辑，移除所有条件删除逻辑（空标签行、空 backtick 条件句）
- [ ] 22 个 prompt 模板文件中无 `{{PLACEHOLDER}}` 格式（全部为 `{{.Placeholder}}` 格式）
- [ ] `{{if .PhaseSummary}}` 条件块正确渲染：有值时注入段落，无值时段落消失
- [ ] `{{if .CoverageStrategy}}` 条件块正确渲染：testable 类型渲染 coverage 指令，non-testable 类型无 coverage 段落
- [ ] Surface 推断失败时 `addSingleFixTask()` 返回错误而非创建空 surface 的任务文件
- [ ] `forge prompt get-by-task-id` 输出与迁移前功能等价（golden-file 对比，允许空白行差异）
- [ ] `ValidatePromptTemplates()` 在启动时对所有模板执行 `template.Parse()`，确保无语法错误

```
consistency_check_result:
  status: pass
  pairs_checked: 15
  conflicts_found: 0
```

## Next Steps

- Proceed to `/quick-tasks` to generate tasks from this proposal
