---
created: "2026-05-27"
author: "forge-brainstorm"
status: Draft
---

# Proposal: 精简任务 Prompt 模板

## Problem

Forge 的 15 个任务 prompt 模板包含大量非指令内容——注释、解释性描述、冗长的角色定义——这些内容不指导 agent 行为，只增加 token 消耗并稀释指令清晰度。同时 task-executor agent 的 Execution Protocol 存在步骤冗长、逻辑重叠的问题。

### Evidence

模板中冗余内容的量化分析：

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

### Urgency

- 每个 task 执行都在消耗这些冗余 token，日积月累规模可观
- 清晰的 prompt 减少 agent 误解和执行偏差
- Prompt 精简是持续优化的一部分，目前已有 prompt-template-audit 等基础，可以在此基础上推进

## Proposed Solution

**就地精简**：保持现有模板独立，在每个文件内部删除非指令内容，将模糊描述改为清晰指令。

不抽取公共模块，不改变现有分类体系。

### Innovation Highlights

本方案不是技术创新，而是对现有 prompt 的"清理"。核心原则是"prompt 是指令，不是文档"——删掉所有不能直接指导 agent 行动的文字。

**行业参照：** 本方案的设计哲学与以下行业实践一致：
- **LangChain Prompt Templates** 在模板中区分"指令（instructions）"与"上下文（context）"，推荐仅将直接影响模型行为的文本保留在系统 prompt 中，解释性描述移至外部文档。
- **Anthropic Prompt Engineering Guide** 强调"show, don't just tell"——通过示例约束行为而非通过自然语言角色描述；本方案中的 AC 验证块精简（保留 REQUIRED 指令、删除展开说明）遵循同一原则。
- **OpenAI GPTs Instructions** 模式的演变方向也是删除冗余的系统 prompt 装饰，改用精确的祈使句指令。

### User-Facing Behavior

本提案对用户（task 执行者）无可见功能变化——用户提交 task 后看到的是相同的执行流程、相同的输出结构、相同的结果质量。唯一的可观测差异是 token 消耗降低（详见 Token 估算），在计费侧表现为每次 task 执行的成本下降。agent 行为层面的"无行为变更"经 SC2 轨迹对比验证。

## Requirements Analysis

### Key Scenarios

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

5. **code-quality-simplify / validation-code / validation-ux** 模板（共 3 个，约 30-50 行/个）：
   - 角色描述精简（同 coding-* 模式）
   - 无 AC 验证块和 CODING_PRINCIPLES——冗余集中在角色描述和框架性说明行
   - code-quality-simplify：~35 行。角色描述 5 行（转为祈使句保留 2 行）、Record Fields 说明 3 行（删除 3 行）、框架描述行 4 行（合并至 2 行）→ 可精简 ~6 行
   - validation-code：~45 行。角色描述 4 行（转为祈使句保留 2 行）、AC 验证说明 4 行（合并至指令行保留 1 行）、Record Fields 说明 3 行（删除 3 行）、框架描述行 3 行（合并至 2 行）→ 可精简 ~8 行
   - validation-ux：~50 行。角色描述 4 行（转为祈使句保留 2 行）、UX 评估标准说明 5 行（核心标准保留 2 行、展开说明 3 行删除）、Record Fields 说明 3 行（删除 3 行）→ 可精简 ~8 行
   - 三者合计精简约 22 行

4. **task-executor agent**：
   - Execution Protocol 步骤合并（步骤 4/5/6 处理 prompt 获取逻辑可合并为 1 步）
   - Retry Strategy 与 Complex Error Pause Flow 去重合并
   - 输出格式合并为紧凑格式

   **步骤 4/5/6 合并前错误恢复分析：** Step 4 读取模板（失败不可恢复，终止），Step 5 替换变量（失败可降级继续），Step 6 组装（无独立失败场景）。Steps 4-6 构成严格顺序链，合并后错误恢复路径不变，合并安全。

**合并后认知分段设计：** 虽然步骤合并为一个，但须在步骤描述中保留子任务边界标识以避免 agent 将其视为一个模糊的"准备提示"操作。合并后的步骤描述采用结构化格式：(a) 读取模板 → (b) 替换变量 → (c) 组装为最终 prompt，各子任务之间设定流转条件——(a)若失败则终止，(b)若失败则降级使用未替换模板，(c)自动执行无失败场景。此格式保留了原 3 步的认知分段，仅在协议层级合并以减少步骤列表长度。 
   **Retry 与 Error Pause 正交性分析：** Retry 操作单次 LLM 调用的临时错误，Error Pause 操作整个 task 的持久错误，层级不同，正交可合并。

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
| 原则 2: 反例/边界说明 | 2-5 行 | 约束边界演示 | 每原则保留 1 个代表性示例 + 压缩边界说明为 1 行概括。同上——示例作为原则之间的结构边界，全部删除可能抹除原则间的分段线索 |
| 超原则通用说明（如作用域声明） | 1 行 | 元指令 | 保留 |

~50 行 → ~25 行（50%），每原则保留 1 行指令 + 1 行边界概括 + 1 个代表性示例。保留示例的动机：在密集指令排列中示例起到"视觉分隔"和"注意力重置"作用——原则之间的结构边界在无示例时容易被高密度的指令行模糊，导致原则混叠（agent 将原则 A 的边界条件错误地应用到原则 B）。保留 1 个示例比全部删除更安全。 

**边界概括的功能等价性说明：** 边界的原始内容由两部分组成——(a) 正反对比示例（0-2 行），(b) 自然语言边界描述（1-3 行）。精简策略将 (a) 浓缩为 1 个代表性示例（保留"示范学习路径"的锚点），将 (b) 压缩为 1 行概括。压缩后的边界概括与原始边界描述的功能等价性由以下措施保障：修改者在生成边界概括后，对照原始边界描述中的每一条排除条件，确认边界概括的排他性条款覆盖全部负例。此校验结果记录在功能快照清单的 CODING_PRINCIPLES 节点备注中。

### 指令分类标准 
在上述逐类型分析中已经隐式使用了分类框架，现将其显式声明为方法论基础：

**三类指令的操作性定义：**

| 类别 | 定义 | 示例 | 精简处理策略 | 方法论依据 |
|------|------|------|-------------|-----------|
| **A. 正面指令** | 告诉 agent 应该做什么的祈使句或模态动词句（must/should/need to） | "Keep the existing behavior unchanged" / "You must include tests" | 保留。仅删除修饰性前置语（"Note that..." → 保留核心动词） | 可直译的 agent 行为规则，删除即丢失功能 |
| **B. 负面约束** | 告诉 agent 不应该做什么的否定句或禁止性表述 | "Do NOT remove format markers" / "You must not skip tests" | 保留。仅删除双重否定和展开说明 | 同 A——agent 需要知道禁令边界 |
| **C. 行为示范** | 通过正例/反例展示期望行为而非直接指令 | CODING_PRINCIPLES 中的 "Good: `{...}` Bad: `{...}`" | 按原则保留 1 个代表性示例。见 CODING_PRINCIPLES 逐原则分析 | 作用于 LLM 的示范学习（few-shot）路径，与指令路径正交；全部删除则失去该调节手段 |

**"约束边界演示"的方法论声明：** CODING_PRINCIPLES 中的行为示例/边界说明（上述类别 C）在本分析框架中归类为"约束边界演示"——其核心功能不是传递直接指令（由原则首行指令行承担），而是通过正反对比展示核心指令的适用范围和例外条件。这是 LLM prompt 中除"指令路径"之外的"示范学习路径"的显式设计。在精简应保持"指令路径优先"——先确保所有原则的首行指令句保留（类别 A/B），再按需保留边界示例（类别 C）。**此区分在整个提案中作为统一方法论使用。**

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

**分析结论：** 当前模板的结构性特征主要为 task-executor agent 的内部遍历逻辑消费（通过标题和前缀语义识别指令类别），而非通过字符串匹配方式解析。prompt.go 仅通过 embed FS 按文件名加载完整内容，不做结构解析。因此精简后结构变形（如标题措辞微调或标记前缀格式保持）不会导致运行时组装断裂。**唯一需保持的隐式协议：** 章节标题和 `AC:` 前缀必须保持可识别（不需要完全一致，但需要语义可对应），因为 task-executor agent 通过文本内容中的这些标记来理解指令层次。**实施建议：** (1) 标题措辞精简后须保持语义一致性（如 `## Notes / Implementation Notes` 可统一为 `## Notes`但不可改为 `## Remarks`）；(2) `AC:REQUIRED` / `AC:STRONGLY` 前缀不可删除。

Record Fields 逐字段分析：

| 行类型 | 典型数量 | 处理策略 |
|--------|---------|---------|
| 字段名 + 值（如 `## Output\n{...}`） | 1 行 | 保留 |
| 字段用途说明（如 "This field describes..."） | 1-2 行 | 删除——字段名自解释 |
| 格式示例/占位符展开 | 1-2 行 | 删除——嵌入实际值即可 |

~3 行 → ~1 行（66%），字段名和值保留。

### Non-Functional Requirements

- 精简后模板的指令覆盖必须与精简前等价（不能遗漏 agents 需要知道的信息）
- 所有 task-executor 的行为不发生变化

### Constraints & Dependencies

- 模板文件位于 `forge-cli/pkg/prompt/data/*.md`，由 `prompt.go` 通过 embed FS 加载
- 修改模板不影响 Go 代码，只需修改 .md 文件
- task-executor agent 位于 `plugins/forge/agents/task-executor.md`

## Alternatives & Industry Benchmarking

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| 分层模板组合 | LangChain PromptTemplate / Vercel AI SDK | 语义分离（instruction/tool/context 分层），单层修改不影响其他，改一处影响所有 | 需重构模板分类体系并修改 `prompt.go` 加载逻辑，与"不改后端代码"约束冲突；对 15 个文件引入抽象层，改动面大于收益 | Rejected: 架构约束否决 |
| 引入 DSL 生成 | 模板引擎模式 | 声明式模板定义，通过编译生成最终 prompt，压缩逻辑集中在 DSL 层 | 需要增加 DSL 定义文件、解析器、编译管线，对 15 个小模板引入完整工具链成本过高——模板改动频次低（月级而非天级），DSL 抽象层在小规模场景下维护负担超过收益 | Rejected: 模板规模小、变更频次低，DSL 工具链成本不合理 |
| 什么都不做 | — | 零风险 | token 持续浪费、指令不够清晰 | Rejected: 成本太低 |
| 抽取公共模块 | DRY 模式 | 修改一处同步所有模板 | 需要改 `prompt.go` 逻辑，且被用户否决 | Rejected: 不满足就地要求 |
| **就地精简** | Forge 现有风格 | 零架构变更，每模板独立修改，风险隔离 | 每个文件都要改 | **Selected: 简单直接** |

## Feasibility Assessment

### Technical Feasibility

纯文本编辑，无技术风险。

### Resource & Timeline

文字精简本身为 1 次编码任务（约 0.5 天）。附加制品与验证工作：为每个模板创建 JSON 功能快照清单（逐节点标注 category/type/role，16 个文件约 2-3 小时）、修改者修改后逐项核对（约 1 小时）、reviewer 签署确认（约 1 小时）、SC2 trial run（16 个模板 × 2 runs = 32 次执行，约 0.5 天自动化 + 人工判定差异）。总计约 2 天，其中纯编码约 0.5 天、制品与验证约 1.5 天。

### Dependency Readiness

前置条件：本次 brainstorm 输出的 proposal 通过。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| "我需要保留角色描述让 agent 理解上下文" | Assumption Flip | 角色描述中的自然语言（"You are a focused..."）对 LLM 行为的影响可由祈使句替代——但这仍是假设而非定论。该领域存在争议（部分研究表明系统角色有效，亦有研究显示模型更遵循后续指令），需通过实施后的行为等价验证确认。**此假设与 NFR #2（"所有 task-executor 的行为不发生变化"）存在固有张力：前者承认不确定性，后者要求确定性。该张力通过 SC2 trial run 协议化解——若 SC2 检测到轨迹偏移（一致性 < 90%），回退角色描述修改部分，保留其他精简项。即角色描述精简为有回退机制的实验性变更，而非无条件承诺不变。** |
| "每个模板独立意味着不需要关注跨模板一致性" | XY Detection | 用户确认了「核心流程重复是允许的」，所以跨模板一致性不是问题，不需要抽取公共模块。 |

## Scope

### In Scope
- 修改 `forge-cli/pkg/prompt/data/` 下全部 15 个模板文件：
  - coding-feature.md, coding-fix.md, coding-enhancement.md, coding-cleanup.md, coding-refactor.md
  - gate.md, doc.md
  - test-run.md, test-gen-scripts.md, test-gen-contracts.md, test-gen-journeys.md, test-verify-regression.md
  - code-quality-simplify.md, validation-code.md, validation-ux.md
- 修改 `plugins/forge/agents/task-executor.md`
- 删除 HTML 注释
- 删除 Step 2 解释性描述（"This generates X from Y"）
- 精简角色描述（自然语言 → 祈使句）
- 精简 AC 验证块（~12 行 → ~4 行）
- 精简化 Record Fields（去掉引导性描述，保留字段名和值）
- 精简 CODING_PRINCIPLES（去掉举例和解释）
- 精简 task-executor Execution Protocol（合并步骤）

### Out of Scope
- 不抽取公共模块文件
- 不修改 `prompt.go` 代码逻辑
- 不新增/删除模板文件
- 不改动模板占位符（`{{TASK_ID}}` 等）
- 不改动 Spec Authority Enforcement 逻辑结构（保留现有约束块内容，不增不减）
- 不改动 Hard Rules / CRITICAL 块的逻辑

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 精简过度导致 agent 遗漏关键行为 | Low | High | **制品：** 每模板建立"功能快照清单"——JSON 格式节点台账，每个节点包含：`{id, category, type, content_snippet, role}`。**节点粒度原则：** 以"最小不可拆分语义单位"为粒度——如果删除该内容后需要补充一条新指令来维持语义完整性，则它是一个不可拆分节点。例如 "You must include tests" 是 1 个节点（删除后需补充等价指令），而 "You must include tests. This ensures regression coverage." 是 1 个指令节点 + 1 个可删除的说明节点（删除说明后不需要补充新指令）。**分类枚举字典：** category 允许值——`instruction`（直接命令）、`constraint`（禁令边界）、`example`（行为示范）、`format`（结构装饰如空行/分隔线）、`separator`（原则间分段符）。type 允许值——`hard-rule`、`critical`、`ac-required`、`ac-explanation`、`role-desc`、`record-field`、`principle-core`、`principle-boundary`、`step-header`、`format-marker`。各 type 定义见附录（简化版：对应各分析表中的行类型命名）。role 标明 `keep`/`delete`/`merge`（合并至父节点）。**创建时机：** 修改前由修改者按模板逐行标注。**签署确认标准：** reviewer 逐节点核对：(a) 粒度合理——无应拆未拆的复合节点；(b) category/type 分类正确，与本文档"指令分类标准"一致；(c) role 标注与精简策略一致（如 `ac-explanation` 应标注 `delete` 而非 `keep`）。签署方式：reviewer 在 JSON 文件顶部添加 `{signed_by, signed_at, nodes_count}` 元数据字段。**存储：** 仓库 `scripts/function-snapshots/<template-name>.json`，版本化管理支持 PR diff。**流程：** (1) 修改前签署清单；(2) 修改后逐项比对，每项标记 pass/fail；(3) 全部 pass 方可合并。**判断标准：** 清单中任一项 fail → 回滚修改并重新调整。|
| 多个模板同步修改，跨模板不一致 | Medium | Medium | **基线文件：** 以 `coding-feature.md` 作为分类基准模板，其余 coding-* 模板修改后与其做 diff，确保同样角色结构、同样原则格式、同样 AC 精简模式。**操作者：** 修改者执行 diff，reviewer 确认 diff 差异为合理的特定语义差异（非格式漂移）。**pass/fail：** 非语义性结构差异 > 3 处 → fail。|
| 现有测试基础设施无法检测 prompt 层行为漂移——当前 forge 测试覆盖 Go 代码逻辑和 task 执行结果，但无机制检测因 prompt 修改导致的 agent 行为差异（如指令理解偏差、约束优先级变化） | Medium | High | **治理措施：** 上述 Risk 1 功能快照清单覆盖指令/约束/示例/格式全部节点，修改前后逐项比对。**补充手段：** (1) 修改后对每个模板对应的典型测试场景执行 2 次 trial run（配合 SC2 自动化轨迹 diff 脚本），对比输出一致性；(2) 非删除项语义一致性检查，100% 覆盖，否则回滚。**CI 化：** (a) 功能快照清单存储为版本化 JSON（`scripts/function-snapshots/` 目录），PR 自动 diff 检查节点不可被意外删除；(b) 轨迹对比脚本作为强制 PR check（阻塞合并——轨迹一致性 < 90% 时禁止合入）|
| prompt 变更为有状态修改，合入后发现影响需要回滚但无标准化流程 | Medium | Medium | **回滚流程：** (1) 每批模板修改独立提交，禁止单 commit 修改全部 16 个文件——分 3 批提交（coding-* 为 1 批、gate/doc 为 1 批、test-* 为 1 批），批间间隔至少 1 个 CI 周期以确保隔离。(2) 合入后观察期：合入后运行一轮完整 journey 测试（`just test-e2e`），若任一 journey 出现与 baseline 不同的行为且确认为 prompt 修改导致，立即 `git revert` 对应批次的 commit。(3) 无需 feature flag——模板为静态文件，revert 即恢复行为。(4) 应急备选：若 revert 冲突（如中间插入了其他 commit），从 baseline snapshot（`eval/baseline-snapshot/`）复制回原始模板重新提交。(5) **单模板级快速回滚：** 若问题仅限于单个模板（如某 coding-feature 的 AC 块精简导致验证遗漏），无需 revert 整批 commit——直接使用 `git checkout <commit-hash>^ -- forge-cli/pkg/prompt/data/coding-feature.md` 从该模板的上一个版本中单独恢复该文件，然后重新提交修复。此操作用于 CI 观察期内快速止血，不等同于完整的回归验证——单文件恢复后仍需运行一次 SC2 确认其他模板行为未受影响。|
| 精简导致信息密度提升，关键指令在密集文本中显著性降低 | Medium | Medium | **定性评估方法：** 在 SC2 轨迹对比的基础上增加"指令显著性"的定性分析——对比修改前后 prompt 文本中：(a) 核心指令行的"视觉隔离度"（前后空行/分隔线数量），(b) 关键约束在 prompt 中的"相对位置"（越靠后衰减越显著），(c) 每模板压缩前后的"指令行占比"（instruction_lines / total_lines → 若 > 70% 则标记为高注意力衰减风险，需回退部分删除以恢复间距）。**缓解措施：** 若某模板压缩后指令行占比 > 70%，在关键指令前增加空行分隔（恢复注意力锚点），而非回退全部删除。|
| 合并后 prompt 变更的长期累积行为效应——指令显著性衰减、约束优先级漂移——在单次合并前验证中无法捕获 | Low | Medium | **周期性轨迹重放检测：** 合并后每周自动对 SC2 选定的参考 task 集合执行一次 trial run（复用 SC2 脚本），对比执行轨迹与合并时基线的一致性。若轨迹一致性趋势持续下降（连续 2 周低于 85%），触发告警并安排人工审查。此检测不阻塞部署，仅为监测机制。**实施方式：** 将 SC2 轨迹对比脚本包装为 CI cron job（`scripts/compare-trajectory.sh --baseline <merge-commit> --current HEAD`），输出报告存储至 `eval/trajectory-monitor/` 目录。|

## Success Criteria

**主要指标（保留率）和次要指标（token/行数）的双层结构：** 保留率为首要校验门禁，token 压缩为主要效率指标，行数压缩为次要参考——当保留率不达标时禁止合并，token 压缩不达标可接受，行数不达标不单独处理。

### 前置基线测量（修改前）

- **[SC-Pre] 修改前 token 和行数基线：** 修改开始前，对所有 In Scope 范围内的模板文件及 task-executor 执行 tokenize（使用 Claude Sonnet tokenizer），记录每文件的当前 token 数和行数，作为精简后节约量计算参考基准。**输出物：** 文件 `eval/baseline-token-counts.json`，格式为 `{file_name: {lines: N, tokens: N}}`。**签署确认：** 修改者生成基线文件，reviewer 确认 tokenize 命令正确（tokenizer 版本、模型参数与 SC8 一致）后签署。此基线数据在 SC6 和 SC8 中作为计算节约量的参考起点。

### 功能保留（首要门禁）

- [SC1] 功能约束保留率 **100%**——每个模板修改后，对照功能快照清单逐项比对，所有指令/约束/格式节点保留率为 100% 方可合并。节点分类包括：(1) Hard Rules、(2) CRITICAL 块、(3) Spec Authority Enforcement、(4) CODING_PRINCIPLES 各原则的指令行（不含边界说明——边界说明允许按 SC3 压缩）、(5) Record Fields 字段名与值结构、(6) AC:REQUIRED 指令。检测方法：修改者逐节点标注 pass/fail，reviewer 签署确认。

### 行为等价性

- [SC2] 模板精简后，agent 执行相同 task 的行为无可见差异。**统计效力说明：** 2+2 次 trial run（修改前 2 次 + 修改后 2 次）可检测到中等及以上的行为漂移，但对于 LLM 输出中常见的细微差异（如措辞变化、推理路径非本质差异）可能存在漏报。若 SC2 首次通过后对结果信心不足，可将 run 次数增加至各 5 次（共 10 次 run），使用 Fisher 精确检验评估功能/非功能差异的比例分布差异。默认 2+2 为最小可行方案。**检测协议：** (1) 选取典型 task——覆盖规则：每个 template 选取 1 个典型 task，该 task 应预期覆盖该 template 功能快照清单中 80% 的 instruction/constraint 类别节点（如 coding-feature 的典型 task 应触发 AC 验证、compile/fmt/lint/test 全部 4 个步骤）。典型 task 由修改者从现有 task 库中选取，reviewer 基于 task 描述与模板功能快照清单前瞻性评估预期覆盖率。确认标准：reviewer 认为预期覆盖率 ≥ 80% 即可通过，无需执行前精确计算。(1a) **执行后覆盖率核定（验证而非筛选）：** task 执行完成后，从 agent 轨迹中提取实际触发的指令/约束类型，与功能快照清单做交集运算，计算实际覆盖率 = 触发的 instruction/constraint 节点数 / 清单中 instruction/constraint 节点总数。实际覆盖率用于验证前瞻性评估的准确性而非作为 task 选取条件。若实际覆盖率 < 80%，判定该次 run 的轨迹对比结果效力降低（可能遗漏部分节点的行为差异），但不自动使 run 无效。修改者可重新选取 task 重试，最多 3 次——若 3 次尝试后仍无 task 达到 80% 实际覆盖率，在报告中注明覆盖限制并继续使用最佳覆盖率 run 的数据进行轨迹对比。若未覆盖节点均为静态约束（如输出格式约束——agent 无需"触发"格式约束，它们始终隐式生效），可降级处理——reviewer 签署豁免。(2) 分别在修改前/后模板上执行该 task 各 2 次（共 4 次 run）；(3) 对比同一 task 的 agent 执行轨迹（关键步骤序列、工具调用参数、最终输出结构）；(4) 轨迹一致性 ≥ 90%（容差：步骤顺序因 LLM 生成随机性导致的非功能性差异）视为通过。

**功能性差异 vs 非功能性差异分类判定标准：** 以下为两类差异的定义示例及判定规则，用于指导轨迹 diff 评估中的人工判定环节：

| 分类 | 差异类型 | 示例 | 判定规则 |
|------|---------|------|---------|
| **非功能性差异** | 步骤顺序随机变化 | 修改前步骤序列为 [读取模板 → 生成代码 → 运行测试]，修改后为 [生成代码 → 读取模板 → 运行测试]——仅步骤 1/2 交换，无步骤缺失或新增 | 步骤名集合不变，仅执行顺序因 LLM 采样随机性导致的排列变化 |
| **非功能性差异** | 工具调用参数表达差异 | `cwd` 值相同但格式不同（`/path/to/project` vs `$PROJECT_ROOT` 解析后等价的路径） | 参数语义等价，仅字符串表现形式不同 |
| **非功能性差异** | 输出结构键顺序变化 | JSON 输出的 keys 顺序不同但 keys 集合相同、值语义等价 | 无序结构（如 JSON/dict）的键顺序变化不影响语义 |
| **功能性差异** | 关键步骤缺失或新增 | 修改前有 "运行 lint 检查" 步骤，修改后该步骤消失；或修改前无 "手动审核" 步骤，修改后新增 | 步骤名集合出现增减——合并范围外的步骤被跳过或引入额外操作 |
| **功能性差异** | 工具调用参数语义变化 | `path` 值从 `./src` 变为 `./lib`（不同目录）、或 `command` 值从 `go test ./...` 变为 `go build ./...`（不同操作） | 参数值语义不等价，会导致不同的执行行为和结果 |
| **功能性差异** | 约束引用缺失 | 修改前 agent 在步骤中引用 CRITICAL 约束（"根据 CRITICAL 规则，不能修改 X"），修改后未引用或引用方式语义等价性存疑 | 关键约束在 agent 推理中的显式引用消失——约束从 agent 的推理链中消失 |

**判定流程：** (1) 自动化脚本标记所有 diff 点；(2) 人工对照上表逐项分类——所有标记为非功能性差异的 diff 点须有对应的分类条目佐证，无法归入任一非功能性类别的 diff 点自动判定为功能性差异；(3) 功能性差异点占全体 diff 点的比例 > 10% → SC2 未通过。**自动化：** 轨迹对比通过脚本自动完成——提取每次 run 的步骤名序列、工具调用名、输出结构 key 集合，生成结构化 diff 报告，仅差异判定环节需人工介入。该脚本纳入仓库 `scripts/` 目录，作为 PR check 的强制验证（阻塞合并——轨迹一致性 < 90% 时禁止合入）。**注意力衰减定性评估（SC2 补充协议）：** 在每次 trial run 中额外记录 agent 对关键指令的"响应延迟"——工具调用中首次引用（mention）关键约束（如 CRITICAL 块内容、Hard Rules）的步数位置。若修改后的平均引用步数比修改前延迟 > 2 步，则标记为注意力衰减风险——参考 Risk 5 缓解措施评估是否需要恢复间距。**注意：** 此 SC 为熊市测试——通过不一定保证完全等价，不通过则明确失败。 
### 结构验证

- [SC3] CODING_PRINCIPLES 在 5 个 coding-* 模板中保留全部核心约束指令（每原则至少保留 1 行指令 + 1 行边界概括），通过 diff 确认对比精简前后原则覆盖无遗漏。
- [SC4] Record Fields 在所有出现模板中保留字段名和值结构，字段用途说明可删除。通过 diff 确认字段名节点未被误删。
- [SC5] Step 2 解释性描述（"This generates X from Y..."）在 5 个 test-* / verify 模板中全部删除，通过 grep 确认无残留。

### 效率指标

- [SC6] 15 个模板文件 + task-executor 共减少 **≥1800 tokens**（基于 SC-Pre 基线测算，精确数值以 SC8 实际 tokenize 为准），行数参考指标为 **≥150 行**（去除注释、解释性描述、冗长定义）。token 节省为主要效率指标，行数为次要参考——保留率门禁未通过时，此效率指标不构成合并理由。
- [SC7] task-executor 的 Execution Protocol 步骤数从 11 步减少到 ≤8 步。

### Token 验证

- [SC8] 精简完成后，对每个修改的模板文件执行实际 tokenize（使用 Claude Sonnet tokenizer，tokenizer 版本与 SC-Pre 一致），与 SC-Pre 基线对比输出报告：(a) 每模板的 tokens 对比（修改前 vs 修改后——修改前数据从 SC-Pre `eval/baseline-token-counts.json` 读取）；(b) 每模板总 token 节省；(c) 每模板行数对比（修改前 vs 修改后）；(d) 按日 task 量加权计算的每日 token 节省范围。实际 tokenize 值用于校准 Token 估算中的近似数值。 
## Next Steps

- Proceed to `/write-prd` to formalize requirements