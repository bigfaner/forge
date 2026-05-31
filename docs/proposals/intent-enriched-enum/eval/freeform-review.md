# Freeform Expert Review: Intent Enriched Enum

**Reviewer**: Intent-Driven Pipeline Branching Architect
**Date**: 2026-05-31
**Document**: `docs/proposals/intent-enriched-enum/proposal.md`
**Scope**: Enum-driven dispatch system, conditional pipeline branching, declarative configuration tables for multi-path execution

---

## Background Assessment

本提案试图解决 Forge skill 管线中一个真实且正在恶化的架构债务：proposal intent 的 3 值枚举（`new-feature`、`refactor`、`cleanup`）与 task type 的 8 值体系之间的映射断裂。提案通过三条路径同时推进——扩充枚举到 6 值、引入混合模式 pipeline 分支、简化 brainstorm 推断——形成一个连贯的改进方案。

在验证过程中，我逐一阅读了提案引用的 6 个 skill 文件（brainstorm/SKILL.md、write-prd/SKILL.md、tech-design/SKILL.md、breakdown-tasks/SKILL.md、quick-tasks/SKILL.md）以及辅助规则文件（write-prd/rules/self-check.md、tech-design/rules/design-quality-checks.md、brainstorm/templates/proposal.md）。以下评估基于这些文件的当前实际状态与提案声称之间的交叉验证。

提案对现状的诊断是准确的。例如 write-prd 的 Intent Detection 表确实只有 3 行（`new-feature`、`refactor`、`cleanup`），tech-design 的 Intent Detection 表同样只有 3 行。brainstorm 的 Step 4.5 对 `coding.fix` 确实使用了启发式判断（"does the fix introduce new user-observable behavior?"）。这些与提案 Evidence 部分的声明一致。

---

## Key Risks

`风险：` 提案声称"与 task type 形成干净的 1:1 映射"，但 6 值 intent 无法覆盖 8 值 task type。现有的 `coding.fix` 在 breakdown-tasks 和 quick-tasks 中被标注为 "Auto-generated for test failures via `forge task add`; do not assign manually"，`doc.consolidate` 和 `doc.drift` 也是手动创建的小众类型。提案的 6 值枚举（`new-feature`、`enhancement`、`refactor`、`cleanup`、`fix`、`doc`）实际上覆盖了 task type 的 6 个主要值，但"干净的 1:1 映射"这个措辞暗示完全对齐，而 `doc.consolidate` 和 `doc.drift` 仍然没有对应 intent。提案在 Risk 表中承认了这一点（"6 值已覆盖所有现有 task type（除 doc.consolidate/doc.drift 这两个小众类型）"），但这与 Proposed Solution 中"与 task type 形成干净的 1:1 映射"的表述自相矛盾。这不是一个阻塞性问题，但如果意图是"近似 1:1"，应该在 Proposed Solution 中明确说明覆盖范围而非使用"干净的 1:1"这种绝对化表述。

`风险：` 提案在 Scenario 3 中描述"改变外部 API 的 refactor"通过 Override Signals 覆盖默认行为，但 Override Signals 的具体规则在提案中只停留在概念层面（"PRD 内容中的明确信号可以覆盖默认值"），没有给出信号检测的具体条件。提案的 Success Criteria 要求"Override Signals 规则存在且可被 PRD 内容触发"，但没有定义触发条件的语法或语义。这意味着实现时 LLM 需要从 PRD prose 中自行判断是否存在覆盖信号——这恰好回到了提案在 Alternatives 中否决的"完全内容驱动 pipeline"方案的风险路径。提案拒绝了"完全内容驱动 pipeline"因为"依赖 LLM 判断力，不稳定"，但 Override Signals 的设计本质上也依赖 LLM 对 PRD 内容的解读，只是范围从全部决策缩小到例外场景。范围缩小是否足以补偿稳定性风险，需要在 Pipeline Configuration 表的设计中明确信号检测的结构化规则（如"PRD 中是否出现 CLI 命令变更 / API 变更 / 数据库变更"等布尔条件），否则这个"混合模式"的"混合"边界是模糊的。

`问题：` 提案将 `enhancement` 作为新增 intent 引入，但在 write-prd 和 tech-design 的当前实现中，`enhancement` 的 pipeline 行为没有明确定义。提案 Key Scenarios 中说 enhancement "pipeline 默认跳过 user stories 但保留 test pipeline"，这实际上是一种全新的 pipeline 配置——既不同于 `new-feature`（全量 pipeline），也不同于 `refactor`/`cleanup`（spec-only）。提案没有讨论 enhancement 是否需要自己的 PRD 格式变体（类似 refactor/cleanup 的 spec-only PRD），也没有说明 tech-design 阶段 enhancement 走哪种决策类型表（new-feature 的全量表还是 refactor 的内部架构表）。这种中间状态如果不提前定义，实现时很可能出现 write-prd 和 tech-design 对 enhancement 行为不一致的情况。

`问题：` 提案称 `fix` "始终为 `fix`，移除启发式判断"，但 brainstorm 的 Step 4.5 当前启发式是针对 `coding.fix` task type 而非 `fix` intent——也就是说，brainstorm 中 `coding.fix` 没有独立的 intent 值，需要被映射到 `new-feature` 或 `refactor`。提案改为 `fix` 始终推断为 `fix` intent 是正确的简化，但 Success Criteria 说"fix 始终推断为 fix，不再使用启发式"，而 proposal template 的 intent 字段默认值是 `"new-feature"`。如果 template 不更新默认值（或者改为让 brainstorm 动态填写），新创建的 proposal 可能仍然以错误的默认值开始。

`风险：` 提案声称"变更限于 plugins/forge/ 目录下的 8 个文件"，但 write-prd 和 tech-design 的分支重写涉及的不只是 SKILL.md——还包括 rules/ 目录下的多个文件（self-check.md、design-quality-checks.md），以及 templates/ 目录。提案 Scope 中列出了 self-check.md 和 design-quality-checks.md，但 In Scope 只列了 8 个文件。经过验证，8 个文件确实包含了这两个 rules 文件，所以数量上是对的。但 write-prd 和 tech-design 的 rules/ 目录中可能还有其他文件引用了 intent 值（如 write-prd 的 rules/knowledge-extraction.md、rules/sc-consistency.md），这些文件如果也包含 intent-gated 逻辑但未被列入 In Scope，就会出现遗漏。

`问题：` breakdown-tasks 和 quick-tasks 的 Intent Propagation 当前表述是 "If `proposal.md` has `intent`, use as default type. Individual task `type` overrides. Missing intent -> per-task Type Assignment. 1:1 mapping." 提案要更新为"严格 1:1 映射（6 值）"，但当前 quick-tasks 的 Intent Propagation 只有一句话，没有映射表。从 `intent: refactor` 到 `type: coding.refactor` 的 1:1 映射虽然直觉上简单，但 `intent: enhancement` 应该映射到什么 task type？`coding.enhancement` 存在于 Type Assignment 表中，但 `intent: fix` 映射到什么？`coding.fix` 在 Type Assignment 中被标注为 "Auto-generated ... do not assign manually"。如果 intent `fix` 映射到 `coding.fix`，就打破了 "do not assign manually" 的规则。如果 intent `fix` 映射到其他类型（如 `coding.feature`），那 1:1 映射就不成立。这个矛盾需要在提案中解决。

`风险：` 提案在 Scope 中未列出 `brainstorm/templates/proposal.md`。验证发现该 template 的 intent 字段当前硬编码为 `"new-feature"`。提案 In Scope 中列了 "brainstorm/templates/proposal.md：更新 intent 有效值注释"，所以 template 确实在 8 个文件中。但如果 intent 扩展到 6 个值，template 的默认值应该怎么处理？保持 `"new-feature"` 作为默认是合理的，但 proposal template 中应该列出所有有效值作为注释（如 `<!-- valid values: new-feature, enhancement, refactor, cleanup, fix, doc -->`），否则用户或 LLM 在手动填写时没有参考。提案只说"更新 intent 有效值注释"，但没有具体说明注释格式。

---

## Improvement Suggestions

`建议：` 为 Pipeline Configuration 表定义一个具体的 schema。提案的核心创新是"intent 控制默认 pipeline 配置（一张表），PRD 内容中的明确信号可以覆盖默认值"。建议在提案中直接给出这张表的完整 6 行定义，包括每个 intent 对每个 pipeline 阶段的默认行为（如 brainstorm 的 "Default Intent" 表已有类似格式，但 write-prd 和 tech-design 缺少等效表格）。具体来说：

| Intent | User Stories | API Handbook | ER Diagram | Integration Specs | Test Pipeline |
|--------|-------------|-------------|------------|-------------------|---------------|
| `new-feature` | Generated | Generated | Conditional | Generated | Full |
| `enhancement` | ??? | ??? | ??? | ??? | ??? |
| `refactor` | Skipped | Skipped | Skipped | Skipped | Quality-gate |
| `cleanup` | Skipped | Skipped | Skipped | Skipped | Quality-gate |
| `fix` | ??? | ??? | ??? | ??? | ??? |
| `doc` | Skipped | Skipped | Skipped | Skipped | Skipped |

`enhancement` 和 `fix` 行的行为在提案中没有定义。如果不在提案阶段就明确，实现时 write-prd 和 tech-design 的作者可能给出不一致的默认值。

`建议：` 为 Override Signals 定义结构化的检测条件，而非依赖 LLM 对 PRD prose 的解读。参考提案 Scenario 3 的例子（"PRD 内容包含 CLI 命令重命名信号"），可以定义如下条件规则：

| Signal | Detection Condition | Default Override |
|--------|-------------------|-----------------|
| API change | PRD contains "CLI 命令重命名"、"API 变更"、"endpoint change" | Enable API Handbook |
| DB change | PRD contains "数据库变更"、"schema change"、"migration" | Enable ER Diagram |
| New user flow | PRD contains new user-observable workflow | Enable User Stories |

这种结构化条件的 LLM 遵守度远高于 prose 描述，与提案自身在 Key Risks 中的分析一致（"LLM 对结构化规则的遵守度高于 prose 描述"）。

`建议：` 解决 `intent: fix` 到 task type 的映射矛盾。当前 `coding.fix` 在 breakdown-tasks 和 quick-tasks 中被标注为 "do not assign manually"。有两种解决方案：(a) 保持 `coding.fix` 为 auto-generated only，`intent: fix` 映射到 `coding.feature`（破坏 1:1 映射但保持现有 type 语义）；(b) 允许手动分配 `coding.fix`，更新 breakdown-tasks 和 quick-tasks 的 Type Assignment 表（保持 1:1 映射但改变 type 语义）。建议选择 (b) 并在提案中说明理由——因为 `fix` 作为独立 intent 的核心价值就是区分于 `new-feature`，如果映射时又回到 `coding.feature`，fix 的独立身份就没有意义了。

`建议：` 将 `enhancement` 的完整 pipeline 行为在提案中显式定义。当前提案只说 enhancement "pipeline 默认跳过 user stories 但保留 test pipeline"。需要补充：(a) PRD 格式是否使用 spec-only（类似 refactor）还是简化版 full PRD？(b) tech-design 走 new-feature 决策类型表还是内部架构表？(c) enhancement 的 Success Criteria 格式是什么？建议 enhancement 走简化版 full PRD——保留 Background 和 Goals（因为改善现有功能仍需用户角色分析），跳过 User Stories（改善行为的用户故事通常是"作为用户我希望现有的更好"这种低价值描述），保留 Test Pipeline（改善行为需要测试覆盖来防止回归）。

`建议：` 补充对 write-prd 和 tech-design 中 rules/ 子目录的完整扫描。当前提案列了 self-check.md 和 design-quality-checks.md，但 write-prd 还有 rules/sc-consistency.md、rules/knowledge-extraction.md、rules/ui-functions.md 等文件。建议在提案中增加一个验证步骤：用 grep 扫描所有 rules/ 文件中包含 "intent" 关键字的位置，确保没有遗漏需要更新的 intent-gated 逻辑。

`建议：` 在 proposal template 的 intent 字段中添加有效值注释。当前 template 硬编码 `intent: "new-feature"`。建议改为：

```yaml
intent: "new-feature"  # valid: new-feature | enhancement | refactor | cleanup | fix | doc
```

这确保 LLM 和手动编辑者都能看到完整的有效值列表，减少填写错误。

`建议：` 提案的 Next Steps 写的是"Proceed to `/quick-tasks`"，但这是一个涉及 8 个文件的中型变更，且需要修改 write-prd 和 tech-design 的核心分支逻辑。建议走完整的 `/write-prd` -> `/tech-design` -> `/breakdown-tasks` 流程而非 quick-tasks，以确保设计阶段的充分性。quick-tasks 适合不需要 PRD 和 design 的简单功能，而本提案的变更触及管线的核心分支逻辑，值得更严谨的设计审查。
