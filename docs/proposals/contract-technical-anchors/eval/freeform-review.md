# Freeform Review: Contract Technical Anchors

**Reviewer Persona:** Contract Pipeline & Test Specification Architect
**Document:** `docs/proposals/contract-technical-anchors/proposal.md`
**Date:** 2026-06-05

---

## Section 1: Background Assessment

这个 proposal 要解决的核心问题是 Forge 测试管道中一个真实且反复出现的痛点：gen-test-scripts 在生成测试代码时缺少确定性的技术锚点，被迫依赖 LLM 推断 HTTP method、CLI command name、Web page route 等关键细节。推断错误时，生成的测试与实际代码不匹配，而且三层测试体系（E2E + 单元 + 集成）都无法捕获这种语义-技术层面的错位。pm-work-tracker 中 POST vs PUT 的案例是有力的证据——api-handbook 里明确定义了 `PUT /teams/:teamId/sub-items/:subId/move`，但 Contract 未引用它，导致整个验证链条断裂。

proposal 的核心方案是建立一条从设计文档到 Contract 到测试代码的信息链：tech-design 自动生成各 surface 的 handbook（API 已有，CLI/Web/Mobile 新增），Contract frontmatter 增加锚点字段（endpoint、command、page、screen），gen-contracts 从 handbook 填充这些字段，gen-test-scripts 再与代码侦察结果做交叉验证。其中一个关键设计决策是"设计文档为 authority source"——当设计文档与代码实现不一致时，以设计文档为准修复 Contract，并将不一致标记为代码 bug。

这个 proposal 建立在几个核心假设之上：（1）设计文档本身是准确的且及时更新的；（2）新增的 cli-handbook / page-map / screen-map 格式可以与 api-handbook 的成熟模式对齐；（3）自动修复机制不会引入比它解决的更多的问题；（4）各 surface 的锚点字段定义足够覆盖实际场景。下面我将逐一审视这些假设中的薄弱环节。

---

## Section 2: Key Risk Identification

### 2.1 Authority Source 假设：设计文档永远正确

`问题：` proposal 在"创新亮点"中明确断言"交叉验证以设计文档为 authority source，设计-实现不一致时定位为代码 bug 而非测试问题"。这个断言在理想流程（先设计后实现）中是合理的，但现实中存在大量场景设计文档本身就有错误或过时。proposal 自身的 Risk Table 里也承认了这个可能：

> "自动修复覆盖了正确的 Contract（设计文档本身有误）| L | H | 修复前保存原始值到 Contract 的注释中，可回溯"

`风险：` "保存原始值到注释中可回溯"这个缓解措施严重不足。当一个错误的自动修复进入 Contract 后，它会影响后续所有 gen-test-scripts 的生成。除非有明确的回滚流程和检测机制，否则"可回溯"只意味着"理论上可以人工恢复"，而不是"系统自动防止错误传播"。更严重的是，`建议：` 后文中会详细讨论，这里要指出的是 proposal 低估了设计文档有误的概率。在快速迭代的项目中，设计文档经常滞后于代码变更，如果每次代码先行更新了接口但设计文档未同步，交叉验证就会以过时的设计文档为准去"修复"实际上已经正确的 Contract——这不是修复，是破坏。

### 2.2 自动修复的静默腐败风险

`问题：` proposal 中"gen-test-scripts 交叉验证并自动修复"这一机制缺乏足够的安全保障。具体来说：

> "将 Fact Table（代码侦察）与 Contract frontmatter 比对，不匹配时以设计文档为准自动修复 Contract"

`风险：` "自动修复"是整份 proposal 中最危险的操作。如果一个 Contract 的 endpoint 原本是正确的，但由于代码侦察的某个 edge case（比如路由注册使用了动态模式、装饰器注册、运行时注册等非静态可侦察方式）导致 Fact Table 得到了错误信息，交叉验证就会误判为不匹配，然后用一个可能同样有误的设计文档值去覆盖正确的值。这个 proposal 没有讨论代码侦察本身的精度和局限性。Fact Table 的代码侦察是静态分析，它无法覆盖所有路由注册模式——如果项目使用了插件系统、动态加载、反射机制等，侦察结果就是不完整的。当侦察结果不完整时，"不匹配"可能只是"没侦察到"，而非"真的不匹配"。

`问题：` proposal 在 Non-Functional Requirements 中声称"交叉验证在 gen-test-scripts Step 1（代码侦察）中执行，无额外网络或 IO 开销"。这个说法忽略了交叉验证本身需要读取并解析 handbook 文件的 IO 开销，以及比对逻辑的计算开销。虽然可能不大，但不应该声称"无额外开销"。

### 2.3 Handbook 格式一致性与多 Surface 扩展

`问题：` proposal 对 cli-handbook、page-map、screen-map 的格式定义几乎一笔带过：

> "api-handbook 已稳定运行。cli-handbook / page-map / screen-map 是新增文档类型，无前置依赖。"

`风险：` API surface 的 endpoint 是一个高度结构化的标识符（HTTP method + URL path），天然适合做锚点。但 CLI command 的标识要复杂得多：子命令嵌套（`git remote add`）、命令别名（`git rm` vs `git remove`）、参数变体（短选项 `-f` vs 长选项 `--force`）等。page-map 和 screen-map 的挑战更大：Web 页面有路由、有组件名、有状态变体；Mobile screen 有导航栈位置、有平台差异（iOS vs Android）。proposal 没有讨论这些 surface 的锚点字段应该包含哪些子信息来确保唯一性和确定性。一个简单的 `page: "user-profile"` 可能根本不够用——是哪个路由下的 user-profile？带 query parameter 的变体呢？

`问题：` proposal 在 Scope 中承诺了"全 surface 覆盖"，但在 Key Risks 中又说"可分批实现"。这两者之间存在张力。如果分批实现，那么未实现 surface 的 Contract 缺少锚点时的行为是什么？proposal 只讨论了"缺少 handbook 时"的降级，但没有讨论"handbook 存在但锚点填充逻辑尚未实现"的中间状态。

### 2.4 信息链完整性中的断点

`问题：` proposal 描述的信息链是"tech-design -> gen-contracts -> gen-test-scripts"，但在这条链上有几个隐含的断点没有讨论：

1. **tech-design 更新但不重新生成 handbook**：当设计文档更新时，是否有机制触发 handbook 的重新生成？如果没有，handbook 就会过期，而 gen-contracts 从过期的 handbook 填充锚点就会引入错误。

2. **Contract 手动编辑后的锚点失效**：用户可能手动编辑 Contract，修改了操作描述但忘记更新 endpoint。proposal 没有讨论 Contract 被手动编辑后的锚点一致性检查。

3. **gen-contracts 到 gen-test-scripts 之间的时间差**：在这两个步骤之间，代码可能已经发生了变更（比如 refactor 改了路由），但 Contract 的锚点还是旧的。gen-test-scripts 的交叉验证是否能处理这种情况取决于代码侦察的覆盖度，而这又回到了 2.2 中讨论的侦察精度问题。

`风险：` proposal 在 Assumptions Challenged 中提到"Fact Table 代码侦察足以保证测试准确性"被证实不足，但它的解决方案（加入交叉验证）仍然依赖 Fact Table 的正确性。这本质上是用 Fact Table 的正确性来验证 handbook 的正确性，当两者都不可靠时，交叉验证就成了两个不可靠源之间的比对。

### 2.5 Scope 边界的模糊地带

`问题：` proposal 的 Out of Scope 中排除了"gen-journeys 改动（Journey 是语义层，不需要技术锚点）"。但 Journey 到 Contract 的映射过程中，Journey 的步骤描述是否包含了 surface 类型信息？如果 Journey 说"用户通过 API 创建团队"，gen-contracts 需要知道这是 API surface 才能去读 api-handbook。如果这个 surface 类型信息不在 Journey 中而在 Contract 中，那么 anchor 填充时就有一个先有鸡还是先有蛋的问题：需要知道 surface 类型才能读正确的 handbook，但 surface 类型本身可能需要从 handbook 中确认。

`问题：` proposal 排除了"代码不存在时的强制交叉验证"，但这恰恰是先设计后实现流程中最关键的时刻。在代码还没写之前，如果设计文档中的 endpoint 定义就有冲突（比如两个操作映射到同一个 endpoint），现在没有任何机制能在设计阶段就捕获这个问题。

### 2.6 向后兼容的边界条件

`问题：` proposal 承诺"缺少 handbook 时管道正常运行（向后兼容，零中断）"，但缺少 handbook 时的降级行为是"gen-test-scripts 发现缺少锚点时降级为 Fact Table 推断"。这意味着在没有 handbook 的项目中，行为和现在完全一样——包括现在已知的问题。proposal 没有讨论是否有机制提示用户缺少 handbook 以及建议生成 handbook。如果用户不知道 handbook 的存在和价值，向后兼容就变成了"安静地维持现状中的 bug"。

`风险：` 当项目只有部分 surface 有 handbook 时（比如有 api-handbook 但没有 cli-handbook），交叉验证只在部分 Contract 上生效。这种不一致可能导致用户产生虚假的安全感——以为所有 surface 都在被验证，但实际上只有 API 被覆盖了。

---

## Section 3: Improvement Suggestions

`建议：` **将自动修复改为"建议修复"加人工确认流程，至少在初始版本中如此。** 自动修复是 proposal 中风险最高的操作。改为输出一份"修复建议报告"，列出所有不匹配项及其建议值，由用户确认后批量应用。这不仅消除了静默腐败的风险，还能让用户在确认过程中发现设计文档本身的错误。proposal 可以设计一个渐进策略：初始版本需要确认，收集数据后分析自动修复的准确率，达到阈值后再开启全自动模式。这样 proposal 的 value 仍然成立（发现不匹配是最有价值的部分），但风险大幅降低。

`建议：` **为每个新 handbook 类型定义最小可行的锚点 schema，而非简单对齐 api-handbook。** API endpoint 的锚点格式很自然：`method + path`。但 CLI command 需要 `command + subcommand + alias`，Web page 需要 `route + component + state`，Mobile screen 需要 `screen-name + navigation-path + platform`。建议在 proposal 中增加一个 sub-section，为每种 surface 定义锚点字段的完整 schema，包括必填字段和可选字段。这不会大幅增加 proposal 的篇幅，但会显著降低后续实现的歧义。特别是 CLI surface，需要考虑命令层级嵌套的表示方式——是用空格分隔的字符串（`"team member add"`）还是结构化对象（`{command: "team", subcommand: "member", action: "add"}`）？

`建议：` **增加"handbook 新鲜度检查"机制。** 在 gen-contracts 读取 handbook 时，比对 handbook 的生成时间戳与 tech-design 文档的最后修改时间。如果 tech-design 比 handbook 更新，发出警告或自动触发 handbook 重新生成。这是信息链完整性的关键保障，防止过期的 handbook 成为错误锚点的来源。具体实现可以是在 handbook 文件的 frontmatter 中记录 `generated_from: <tech-design-path>` 和 `generated_at: <timestamp>`，gen-contracts 读取时做比对。

`建议：` **增加交叉验证结果的分类报告，而非简单的"修复"或"标记 bug"。** 当 Fact Table 和 handbook 不匹配时，可能有多种原因：（a）代码侦察不完整（动态路由）；（b）设计文档过期；（c）真正的代码 bug。proposal 当前只区分了（b）和（c），忽略（a）会直接导致误修复。建议将不匹配结果分为三类：高置信度不匹配（静态路由注册可确认）、低置信度不匹配（可能因侦察不完整）、无法验证（侦察未发现相关路由）。只有高置信度的不匹配才自动处理或建议修复。

`建议：` **在向后兼容降级路径中增加用户提示。** 当 handbook 不存在时，不仅在后台降级为 Fact Table 推断，还在输出中明确提示："Surface X 缺少 handbook，Contract 的技术锚点未被验证。建议运行 tech-design 生成 handbook 以启用锚点验证。" 这让向后兼容不再是一种静默降级，而是一个引导用户改善测试质量的契机。

`建议：` **明确分批实现的路线图和每批的 scope。** 如果全 surface 覆盖确实要分批（Risk Table 中提到），那么 proposal 应该明确第一批实现哪些 surface（建议先做 API，因为它最成熟、格式最稳定），后续 batch 的触发条件是什么（比如 API surface 的锚点验证准确率达到 X% 后扩展到 CLI）。这样既控制了 scope，又给了完整的愿景。

`建议：` **在 Out of Scope 中增加对"Contract 手动编辑后的一致性"的明确声明。** 如果这次不做，至少要明确记录这个 gap，以免后续假设 gen-test-scripts 的交叉验证能覆盖所有场景。实际上，一个轻量级的解决方案是在 Contract frontmatter 中记录 `last_anchor_sync: <timestamp>`，当 Contract 的 `updated_at` 大于这个值时，gen-test-scripts 输出一个"锚点可能过期"的警告。
