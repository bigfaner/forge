---
reviewer: "Template Engine Migration & Rendering Pipeline Architect"
proposal: "docs/proposals/slim-task-prompt-templates/proposal.md"
date: "2026-05-28"
review_type: "freeform"
---

# Freeform Expert Review: 精简任务 Prompt 模板 + Frontmatter 元数据结构重构

## Background Assessment

本提案解决 Forge 插件中两个相互独立但同时推进的技术债务：第一，15 个 prompt 模板和 task-executor agent 中累积的约 190 行非指令性内容导致每次 task 执行浪费约 1200-1500 tokens（输入侧），日积月累形成可观的 token 消耗；第二，41 个模板的 metadata frontmatter 采用扁平 `variables` list，未区分元数据字段和内容字段的语义角色，导致 frontmatter 作为机器可校验契约的效力不足。

内容精简的核心方案是"就地精简"——逐文件删减非指令内容，不引入架构变更。Frontmatter 重构的核心方案是将扁平 `variables` list 改为 `identity`/`context`/`conditional`/`variables` 四分组结构，同时将 PhaseSummary 从 frontmatter 移至正文独立 section。

从渲染管线角度审视，当前流程是：`embed.FS` 加载模板 -> `placeholderReplacer` 将 `{{X}}` 转为 `{{.X}}` -> `stripMetadataFrontmatter` 剥离 metadata frontmatter -> `text/template` 执行渲染 -> `strings.Replace` 注入 TASK_CATEGORY -> `collapseBlankLines` 后处理。Frontmatter 重构只修改 `parseMetadataFrontmatter` 的解析逻辑和 `validateMetadataVariables` 的校验逻辑，理论上不触碰渲染管线本身。但"理论上"和"实际上"之间有几个关键缝隙，以下逐一分析。

提案的一个核心假设是"prompt 中被删除的内容不直接影响 agent 执行行为"，验证手段是功能快照清单的 100% 保留率和轨迹对比的 90% 一致性。另一个核心假设是"frontmatter 重构对渲染输出无影响"，依赖 `stripMetadataFrontmatter` 在渲染前剥离 frontmatter 的机制。

## Key Risks

提案在 frontmatter 重构的技术实现细节上存在若干风险，以下从渲染管线架构的角度逐一分析。

风险：`parseMetadataFrontmatter` 的行级 YAML 解析器扩展复杂度被低估——当前实现是一个无状态的手工行扫描器，新增分组支持需要引入有状态的分组上下文跟踪。

> "`parseMetadataFrontmatter` 是单文件内的行级 YAML 解析器，扩展分组支持复杂度低。" — Feasibility Assessment, Technical Feasibility

当前 metadata.go 第 54-77 行的实现是一个逐行 `strings.HasPrefix` 匹配器，唯一的"状态"是 `meta.Variables != nil` 用来判断是否处于变量列表模式。提案要求新增 `identity:`、`context:`、`conditional:` 三个分组，每个分组下有 `key: true` 形式的键值对。这意味着解析器必须跟踪"当前所在的分组"——遇到 `identity:` 行时进入 identity 分组模式，后续的 `taskID: true` 行应被解析为 identity 分组下的字段；遇到 `context:` 或空行时切换或退出分组。这需要引入 `currentSection string` 状态变量和对应的分支逻辑。虽然确实是"单文件范围内"的改动，但行级解析器在处理边缘情况时容易出错：空行是否重置分组状态？缩进不一致的行（如 `  identity:` 有前导空格）是否能正确匹配？注释行 `# identity:` 是否会误触发分组切换？提案没有讨论这些边缘情况。更稳健的替代方案是引入 `gopkg.in/yaml.v3` 标准库替代手工解析器，但提案没有评估这个选项。

问题：分组字段的 PascalCase vs camelCase 命名规则存在内在矛盾，可能导致反射校验失败。

> "YAML 使用 camelCase 但与 Go struct 的 PascalCase 字段名不同" — Key Risks 表
> "分组内的字段名使用 Go struct 的原始 PascalCase 名（用于反射校验）" — SC-FM-5

提案在 Key Risks 表和 SC-FM-5 中对字段命名有看似矛盾的说法。Key Risks 表提到"YAML 使用 camelCase"，但 SC-FM-5 要求"分组内的字段名使用 Go struct 的原始 PascalCase 名"。查看提案的示例（Scenario 6），frontmatter 中 identity 分组使用了 `taskID: true`（camelCase），而 `validateMetadataVariables` 通过 `reflect.Type.FieldByName`（metadata.go 第 119 行）校验字段名——这是一个大小写敏感的操作。Go struct `promptTemplateData` 中的字段名是 `TaskID`（PascalCase）。如果 frontmatter 写 `taskID` 但 Go struct 的字段名是 `TaskID`，`FieldByName("taskID")` 将返回 false，校验失败。提案在 Key Risks 表的 Mitigation 中提到"分组声明中的 key 使用 camelCase 仅当 key 是分组 label，不用于反射校验"——但这句话含义模糊。如果分组内的字段不参与反射校验，那新增分组的校验意义何在？如果参与校验，则必须使用 PascalCase。提案需要明确统一命名规范并消除这个矛盾。

风险：PhaseSummary 迁移引入的条件块格式变更与当前渲染逻辑生成的值格式不兼容。

> "在所有 21 个 prompt 模板正文中添加 `## PhaseSummary` 独立 section" — Scope, Frontmatter 重构
> "`phaseSummaryLine = 'PHASE_SUMMARY: ' + phaseSummaryPath`" — prompt.go 第 189 行

当前模板中 PhaseSummary 的渲染是 `{{if .PhaseSummary}}{{.PhaseSummary}}{{end}}`（见 coding-feature.md 第 20 行），渲染结果为单行文本如 `PHASE_SUMMARY: path/to/summary.md`。提案要求改为 `{{if .PhaseSummary}}\n## PhaseSummary\n{{.PhaseSummary}}\n{{end}}`（多行条件块）。但 `renderTemplate` 函数（prompt.go 第 187-191 行）生成的 PhaseSummary 值是 `"PHASE_SUMMARY: " + phaseSummaryPath` 形式的单行文本。如果新格式在其前面加了 `## PhaseSummary` 标题行，输出将变为：

```
## PhaseSummary
PHASE_SUMMARY: path/to/summary.md
```

这产生了语义重复：一个 Markdown 二级标题 `## PhaseSummary` 后面跟着的仍然是带 `PHASE_SUMMARY:` 前缀的旧格式文本行。标题和前缀在语义上是冗余的。更关键的是，提案声明"不修改 `prompt.go` 渲染逻辑（`Synthesize()` 等）"为 Out of Scope（Scope 章节），但 PhaseSummary 格式的一致性恰恰需要修改 `renderTemplate` 中 `phaseSummaryLine` 的生成逻辑。如果保持 `phaseSummaryLine` 不变，新的 section 格式在语义上就是冗余的；如果修改 `phaseSummaryLine`，就违反了 Out of Scope 约束。这是一个提案尚未意识到的自相矛盾。

风险：`validateMetadataVariables` 校验范围的扩展未匹配不同模板类型的 Go struct。

> "`validateMetadataVariables`（或新的 `validateGroupedMetadata`）校验所有分组（`Identity`/`Context`/`Conditional`/`Variables`）中的字段名集合均存在于对应的 Go struct 中" — SC-FM-3

当前 `validateMetadataVariables`（metadata.go 第 93-109 行）只校验 `Variables` 列表中的字段名是否存在于指定的 struct 中。调用点是 `ValidatePromptTemplates`（prompt.go 第 76-120 行），它只遍历 `task.ValidTypes` 并统一使用 `promptTemplateData` 进行校验。task 模板和 record 模板不在该函数的校验范围内——它们有不同的 Go struct（task 模板使用 `task.Task` 的字段，record 模板使用 `RecordTemplateData` 的字段）。提案要求校验扩展到所有分组，但没有说明：(1) task 模板和 record 模板的 frontmatter 应该使用哪个 struct 校验？(2) 是否需要新增 `ValidateTaskTemplates` 和 `ValidateRecordTemplates` 校验函数？(3) 如果不新增校验函数，那 task 和 record 模板中新增的分组字段如何被验证？缺乏这些定义意味着 frontmatter 重构可能引入未校验的字段名拼写错误，而启动时的 `ValidatePromptTemplates` 不会捕获这些错误。

问题：41 个模板的分组层级差异（prompt 三组、task 两组、record 一组）缺乏明确的判定规则。

> "21 个 prompt 模板：添加 `identity`/`context`/`conditional` 分组；14 个 task 模板：添加 `identity`/`context` 分组；6 个 record 模板：添加 `identity` 分组" — Scope, Frontmatter 重构

提案为三类模板预设了不同的分组层级，但未解释判定依据。为什么 task 模板不需要 `conditional` 分组？查看 task 模板示例（coding.cleanup.md），其 frontmatter 中确实有 `SurfaceKey`、`SurfaceType` 等字段，这些在模板正文中以 `{{if .SurfaceKey}}` 条件块使用——按提案的分组语义，这些应该是 `conditional` 类型。但提案没有为 task 模板分配 `conditional` 分组。同样，record 模板中 `TypeReclassification` 字段在正文中以条件块使用（`{{if .TypeReclassification}}`），但没有被归入 `conditional` 分组。缺乏明确的分组判定规则意味着实施时遇到边界情况没有决策依据。

风险：向后兼容性的语义映射未定义——新格式的四分组字段并集与旧格式 variables 列表的等价关系不明。

> "`parseMetadataFrontmatter` 保持容错：优先解析分组字段，若无分组则 fallback 到旧 variables list 解析。测试覆盖新旧两种格式。" — Key Risks 表

SC-FM-2 要求"旧格式的模板文件能被正确解析"和"解析结果中的 `Variables` 字段与旧解析器一致"。但这只覆盖了"旧格式模板在新解析器下正常工作"的方向。反向的语义等价问题没有定义：新格式中 `identity.taskID: true` + `variables: [TaskFile, ...]` 的组合，在旧代码视角下等价于什么？`TemplateMetadata.Variables` 列表应该只包含 `variables` 分组下的字段，还是应该包含所有分组字段的并集？如果只包含 `variables` 分组的字段，那旧代码中依赖 `Variables` 列表包含完整字段集的逻辑（例如用于模板完整性检查）将失效。提案没有定义新格式的 `TemplateMetadata` 结构体应该如何表示分组信息——是保持扁平的 `Variables []string` 并新增 `Identity map[string]bool` 等字段，还是将所有分组的字段合并到统一的 `Variables` 列表中？这个数据结构决策直接影响向后兼容性的实现方式。

问题：SC-FM-1 的迁移检测逻辑存在逻辑矛盾——将"设计上不需要 identity 的模板"也判定为"未迁移"。

> "若重构后模板的 frontmatter 不含 `identity:` 键则算作未迁移（包括那些确实无需 identity 字段的模板——如空 frontmatter 的模板）" — SC-FM-1

SC-FM-1 使用 `identity:` 键的存在性作为"已迁移"的唯一判据，但同时承认"包括那些确实无需 identity 字段的模板"。这产生了逻辑矛盾：如果一个模板根据设计确实不需要 identity 分组，它会被错误地标记为"未迁移"。正确的做法是按模板类型分别定义迁移标准：prompt 模板必须有 `identity` + `context` + `conditional`，task 模板必须有 `identity` + `context`，record 模板必须有 `identity`。统一用 `identity:` 存在性做检测过于粗糙。

风险：renderTemplate 中 TASK_CATEGORY 的 `strings.Replace` 注入依赖精确字符串匹配，内容精简可能破坏这个隐式协议。

> "`result = strings.Replace(result, "TASK_FILE: "+td.TaskFile, "TASK_FILE: "+td.TaskFile+"\nTASK_CATEGORY: "+td.TaskCategory, 1)`" — prompt.go 第 251-252 行
> "不修改 `prompt.go` 渲染逻辑（`Synthesize()` 等）" — Scope, Out of Scope

当前 `renderTemplate` 通过 `strings.Replace` 在渲染后的输出中查找 `"TASK_FILE: " + taskFile` 字符串，并在其后注入 `TASK_CATEGORY:` 行。这个替换依赖精确匹配 TASK_FILE 行的格式。如果内容精简修改了 `TASK_FILE:` 行的格式（例如改变缩进、在行尾添加空格、将 TASK_FILE 行与其他行合并），`strings.Replace` 将静默失败——不会报错，只是不注入 TASK_CATEGORY 行，导致 task-executor 收到的 prompt 缺少 TASK_CATEGORY 信息。提案声明"不修改 `prompt.go` 渲染逻辑"为 Out of Scope，但没有将"保持 TASK_FILE 行格式不变"列为显式约束。

风险：内容精简的"功能快照清单"是验证体系的核心依赖但创建时机、粒度标准、分类字典均未定义。

> "功能约束保留率 **100%**——每个模板修改后，对照功能快照清单逐项比对，所有指令/约束/格式节点保留率为 100% 方可合并。" — SC1

SC1 要求 100% 的功能保留率，但"功能快照清单"本身的创建标准是模糊的。提案没有在 Scope 或 Next Steps 中定义谁在何时创建这个快照清单，清单中每个"节点"的粒度是什么（一条完整的原则是一个节点，还是原则中的每句话各是一个节点？），以及 `category`/`type` 的分类枚举是什么。如果不同创建者对粒度和分类有不同的理解，100% 的保留率就没有可比性。对于被标为 Low Likelihood / High Impact 的核心风险，其缓解措施的严谨性应该更高。

问题：task-executor Execution Protocol 步骤合并的具体操作定义不充分。

> "步骤 4/5/6 处理 prompt 获取逻辑可合并为 1 步" — 内容精简场景，Key Scenarios
> "Execution Protocol 步骤数从 11 步减少到 <=8 步" — SC7

提案只提到"步骤 4/5/6 可合并"，但没有给出合并后的步骤定义。从 11 步到 <=8 步需要减少 3 步——步骤 4/5/6 合并为 1 步省 2 步，那第 3 步的减少来自哪里？另外，Retry Strategy 与 Complex Error Pause Flow 的"去重合并"具体如何操作？提案中提到这两处可以合并但没有展开。实施者缺乏足够的操作指导。

## Improvement Suggestions

建议：引入 `gopkg.in/yaml.v3` 标准库替代手工行级解析器解析 metadata frontmatter。

Addresses: 行级 YAML 解析器扩展复杂度被低估。

> What changes: 将 `parseMetadataFrontmatter` 中的手工行解析替换为标准 YAML 库的 `yaml.Unmarshal`。先提取 `---` 之间的内容（保留当前的 `---` 分隔符检测逻辑），然后用 YAML 库解析到新的 `TemplateMetadata` 结构体。新结构体定义为：

```go
type TemplateMetadata struct {
    Type       string            `yaml:"type"`
    Category   string            `yaml:"category"`
    Identity   map[string]bool   `yaml:"identity,omitempty"`
    Context    map[string]bool   `yaml:"context,omitempty"`
    Conditional map[string]bool  `yaml:"conditional,omitempty"`
    Variables  []string          `yaml:"variables,omitempty"`
}
```

这样可以避免手工维护分组状态机的各种边缘情况，同时让解析逻辑更健壮、更易扩展。Go 生态中 `gopkg.in/yaml.v3` 已被广泛使用，直接引入即可。向后兼容性通过 `yaml:",omitempty"` 自然实现——旧格式只有 `variables` 字段，新格式有分组字段，两者都能正确解析。

建议：统一字段命名为 PascalCase，与 Go struct 保持一致，消除命名矛盾。

Addresses: PascalCase vs camelCase 命名矛盾。

> What changes: 在 frontmatter 的所有分组中，字段名一律使用 Go struct 的 PascalCase 原始名（如 `TaskID` 而非 `taskID`）。分组容器标签名（`identity`、`context`、`conditional`）使用 lowercase，但这些是分组的元名，不参与反射校验。提案应在显式约束中声明这一规则："frontmatter 分组内的字段名必须与 Go struct 的导出字段名完全一致（PascalCase），因为校验逻辑使用 `reflect.Type.FieldByName` 做大小写敏感匹配。"采纳后，SC-FM-5 的验证标准也相应明确化。

建议：将 PhaseSummary 迁移涉及的 Go 代码修改纳入 In Scope，修改 `phaseSummaryLine` 的生成逻辑以匹配新格式。

Addresses: PhaseSummary 条件块格式变更与旧格式不等价、Out of Scope 约束的自相矛盾。

> What changes: 将 `renderTemplate` 中 `phaseSummaryLine` 的生成从 `"PHASE_SUMMARY: " + phaseSummaryPath` 修改为纯路径 `phaseSummaryPath`，同时将正文 section 模板改为 `## PhaseSummary\nRead: {{.PhaseSummary}}`。这样 section 标题提供上下文，内容只包含路径，语义清晰无冗余。同时将"不修改 `prompt.go` 渲染逻辑"从 Out of Scope 中移除，替换为"仅修改 `phaseSummaryLine` 的格式前缀"这一受控变更。采纳后，PhaseSummary section 的格式自洽且无语义重复，同时提案的 Scope 不再包含自相矛盾。

建议：为三类模板分别定义校验 struct 和分组规则，扩展启动时校验覆盖范围。

Addresses: `validateMetadataVariables` 校验范围变更未充分设计、41 个模板分组归属规则不明确。

> What changes: 提案增加一个"分组规则表"：

| 模板类型 | 校验 struct | 分配分组 | 判定规则 |
|---------|-----------|---------|---------|
| prompt | `promptTemplateData` | identity + context + conditional + variables | identity: 唯一标识模板渲染的键字段（TaskID）；context: 描述运行时上下文的字段（SurfaceKey, FeatureSlug 等）；conditional: 控制正文条件块显示的字段（CoverageStrategy, PhaseSummary 等） |
| task | `task.Task` 对应的 struct | identity + context + variables | identity: 唯一标识任务实例的字段（ID, Title）；context: 任务执行上下文字段（SurfaceKey, SurfaceType） |
| record | `RecordTemplateData` | identity + variables | identity: 唯一标识记录实例的字段（TaskID, TaskTitle, Status） |

同时新增 `ValidateTaskTemplates` 和 `ValidateRecordTemplates` 校验函数，或在现有 `ValidatePromptTemplates` 中扩展校验循环覆盖三类模板。分组判定的核心规则应该是：identity 字段 = 唯一标识实例的字段（删除后无法定位具体实例），context 字段 = 影响行为但不控制条件渲染的字段，conditional 字段 = 在正文中以 `{{if .X}}` 控制段落显示的字段。采纳后，所有模板的分组有明确的判定依据，且所有分组的字段都能在启动时被校验。

建议：在 Constraints 章节中增加"TASK_FILE 行格式稳定性"约束。

Addresses: renderTemplate 中 TASK_CATEGORY 注入的字符串替换依赖。

> What changes: 在 Constraints & Dependencies 中增加一条约束："内容精简不得改变 `TASK_ID:`、`TASK_FILE:`、`SURFACE_KEY:` 行的格式（包括缩进、空格数量、行尾字符），因为 `renderTemplate` 的 `strings.Replace` 后处理依赖精确字符串匹配。"同时在 Success Criteria 中增加验证标准：精简后所有 prompt 模板中 `TASK_FILE:` 行的格式与精简前完全一致，通过 `git diff` 确认。更长期地，建议将 TASK_CATEGORY 注入从 `strings.Replace` 后处理迁移到模板正文中的 `{{.TaskCategory}}` 占位符，彻底消除这个隐式依赖。

建议：定义功能快照清单的创建标准和分类字典。

Addresses: 功能快照清单的粒度标准和分类字典未定义。

> What changes: 在提案 Scope 中增加前置工作项："在内容精简开始前，为每个模板创建功能快照清单。"同时定义：(1) 节点粒度规则——以"最小不可拆分语义单位"为粒度，判断标尺是"如果删除该内容后需要补充一条新指令来维持语义完整性，则它是一个节点"；(2) 分类枚举字典——明确定义所有允许的 `category` 值（instruction / constraint / example / format / separator）和 `type` 值（hard-rule / critical / ac-required / ac-explanation / role-desc / record-field / principle-core / principle-boundary / step-header / format-marker），每个值附带正反示例。采纳后，快照清单的创建有客观标准，不同创建者的分类一致性从"靠默契"提升为"靠规范"。

建议：为向后兼容性定义新格式与旧格式的语义映射规范。

Addresses: 向后兼容性的语义映射未定义。

> What changes: 在提案中增加"迁移语义映射"章节，明确定义：新格式的 `identity` + `context` + `conditional` + `variables` 四个分组的字段并集，应等于旧格式的 `variables` 列表。`TemplateMetadata` 结构体应提供 `AllFields()` 方法返回所有分组字段的并集，供需要完整字段列表的消费方使用。测试用例应覆盖三个场景：(1) 新格式模板——所有分组字段 + variables 字段都能通过校验；(2) 旧格式模板——Variables 列表照常校验；(3) 验证 `AllFields()` 的返回值与旧格式 Variables 列表的等价性。采纳后，向后兼容性有了完整的语义定义和验证标准。
