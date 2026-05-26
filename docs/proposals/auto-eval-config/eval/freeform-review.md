# Freeform Expert Review

**Reviewer**: Developer Tooling & Configuration Architect
**Document**: `docs/proposals/auto-eval-config/proposal.md`
**Date**: 2026-05-26

---

## Background Assessment

This proposal addresses a real friction point in the Forge pipeline: four document evaluation skills (`eval-proposal`, `eval-prd`, `eval-ui`, `eval-design`) each require a manual `AskUserQuestion` prompt before running, creating unnecessary interaction cost. The proposal observes that the project already has a `ModeToggle` configuration pattern via `auto.runTasks`, `auto.test`, etc., and proposes extending this pattern with an `auto.eval` nested struct containing four independent `ModeToggle` fields — one per eval stage — each supporting `quick`/`full` sub-keys.

The core technical approach is to add an `EvalConfig` struct (or equivalent nested fields) to the existing `AutoConfig` Go struct, wire it into the `parseAutoRaw`/`applyDefaults`/`autoModeField` pipeline, and update four skill markdown files to read the config instead of (or before) calling `AskUserQuestion`.

The proposal rests on two key assumptions: (1) that the existing `auto.*` flat namespace can be cleanly extended with a nested `auto.eval.*` sub-namespace without breaking the current YAML parsing and dot-notation resolution logic, and (2) that users will benefit from per-stage granularity rather than a single `auto.eval` toggle. The proposal is internally consistent with the claimed problem statement, and the defaults are chosen to preserve backward compatibility for three of four skills, with `ui-design` being the deliberate exception.

## Key Risks

The proposal introduces a namespace depth that does not exist in the current `auto` block. All current `AutoConfig` fields — `test`, `consolidateSpecs`, `cleanCode`, `validation`, `runTasks`, `knowledgeSave` — are flat `ModeToggle` fields at `auto.{field}.quick/full`. The proposal's `auto.eval.proposal.quick` adds a third nesting level. This is a structural departure from the existing pattern.

风险：嵌套命名空间破坏现有 YAML 解析和 dot-notation 路由逻辑
> "在 `.forge/config.yaml` 的 `auto` 块中新增 `eval` 嵌套结构体，包含 4 个独立的 ModeToggle 字段" — 当前 `AutoConfig` 是纯 flat 结构，所有 `ModeToggle` 字段直接挂在 `AutoConfig` 上。`parseAutoRaw` 函数硬编码了 `modeFields := []string{"test", "consolidateSpecs", ...}` 遍历 `auto` 下的直接子节点，`autoModeField` 通过 `switch field` 匹配字段名，`getAutoKeyValue` 和 `setAutoConfigValue` 用 `strings.SplitN(rest, ".", 2)` 只拆一层。引入 `auto.eval.proposal` 意味着 `auto.eval` 是一个中间结构体而非 `ModeToggle`，现有的 `autoModeField("eval")` 返回 `nil`，`forge config get auto.eval.proposal` 会走入错误分支。这不是"确定性改动"，而是需要对配置系统的解析、默认值填充、get/set 三条路径做结构性重构。

风险：Go struct 嵌套与 YAML 序列化的零值陷阱
> "默认值：`proposal`: `quick: true, full: true`；`uiDesign`: `quick: false, full: false`" — `ModeToggle` 的零值是 `{false, false}`。当 `uiDesign` 的默认值恰好等于零值时，`AutoConfig.WithDefaults()` 和 `applyDefaults()` 无法区分"用户显式设置为 false"和"字段未配置"。当前代码注释明确警告了这一点："This cannot distinguish 'user explicitly set ModeToggle{false, false}' from 'field was never set' because both equal ModeToggle{}"。对于 `uiDesign` 默认 `false/false` 的场景，如果 `raw` 追踪逻辑不支持嵌套路径（当前 `parseAutoRaw` 只看一层），用户显式配置 `auto.eval.uiDesign: {quick: false, full: false}` 和完全不配置 `eval` 块会得到相同的结果，这在语义上是正确的（默认就是 false），但 `techDesign` 和 `prd` 同理，这意味着 `raw` 追踪必须被扩展到嵌套结构体内部，否则无法实现 proposal 默认为 true 而 prd 默认为 false 的差异化行为。

问题：quick/full 区分在 eval 场景缺乏语义基础
> "每个 ModeToggle 支持 `quick`/`full` 子键，区分 quick 和 full 流水线的行为" — 现有 `ModeToggle` 的 quick/full 区分用于 `forge task index` 的两种流水线模式（`GetQuickTestTasks` vs `GetBreakdownTestTasks`）。但 eval 的触发点是 skill 内部（brainstorm、write-prd、tech-design、ui-design），这些 skill 本身并不区分 quick/full 流水线。`brainstorm` 可以在 `/quick` 流水线中被调用，也可以独立运行 `/brainstorm`。proposal 没有说明 skill 如何判断当前处于 quick 还是 full 模式。如果 skill 无法获取当前流水线上下文，quick/full 区分就是死配置。

风险：ui-design 行为变更的向后兼容性声明不完整
> "ui-design 从无条件自动运行改为读取配置，与其他 skill 行为一致" 以及 "默认 full:false 保持询问" — 这意味着所有现有用户升级后，`/ui-design` 将不再自动运行 eval-ui，而是弹出 `AskUserQuestion`。提案将此风险评级为 "M likelihood, M impact"，但 ui-design 是唯一一个当前无条件自动评估的 skill，这个行为可能恰恰是因为 ui 设计质量更需要自动化验证而有意为之。将 "无条件自动" 降级为 "默认询问" 是一个功能回退，不仅仅是不一致性问题。提案中的 Assumption Flip 表格仅仅说了 "应与其他 skill 一致"，但没有提供证据说明为什么 ui-design 原来的行为是错误的。

问题：skill 中的 config check 逻辑实现方式未说明
> "4 个 skill（brainstorm、write-prd、tech-design、ui-design）增加 config check 逻辑" — 当前 skill 是 markdown 文件（`SKILL.md`），它们通过自然语言指令指导 Claude Code agent 的行为。提案没有说明 config check 的实现机制：skill 中是否使用 Bash 命令调用 `forge config get`？还是依赖 agent 的内置配置感知？如果是前者，每个 skill 需要增加一个"先检查配置，再决定是否 AskUserQuestion"的条件分支，这对 skill markdown 的可维护性有直接影响。如果是后者，agent 如何获得 `auto.eval.*` 的值？提案的 Scope 部分提到"使用统一的配置检查模板，在 skill 中用 EXTREMELY-IMPORTANT 标注"，但这个模板不存在于当前代码库中，是一个隐含的新增工件。

风险：四个 skill 的 config check 一致性难以保证
> "使用统一的配置检查模板，在 skill 中用 EXTREMELY-IMPORTANT 标注" — 当前四个 skill 的 eval 触发方式已经不一致：brainstorm 在 Step 7 用 `AskUserQuestion`；write-prd 在 Step 11 用 `AskUserQuestion`；tech-design 在 Step 10 用 `AskUserQuestion`；ui-design 在 Step 7 直接调用 `Skill` tool 无条件运行。要在四种不同的控制流中插入统一的 config check 模板，需要修改每个 skill 的特定步骤，而 "统一模板" 在 skill markdown 语境下只是一个 copy-paste 的文本块，没有运行时保证。随着 skill 演进，这些分散的 config check 片段极易漂移。

问题：`forge config get auto.eval.proposal` 的预期输出格式与现有模式不一致
> 成功标准中写 "`forge config get auto.eval.proposal` 返回 `quick:true full:true`" — 当前 `getAutoKeyValue` 对 `auto.test` 等返回 `"quick:%v full:%v"` 格式的字符串。如果 `auto.eval` 是一个嵌套结构体，`forge config get auto.eval` 应该返回什么？如果返回整个 eval 块的 YAML 表示，现有的字符串格式化逻辑无法处理。如果 `auto.eval` 不是一个有效的查询目标（只有 `auto.eval.proposal` 等三段路径有效），那么 `forge config get auto.eval` 应该返回错误。提案没有定义这种三段式路径的行为边界。

问题："1-2 小时完成" 的时间估算未包含嵌套配置系统的工程复杂度
> "预计 1-2 小时完成所有改动（Go CLI + 4 个 skill + 测试）" — 这个估算假设改动是纯增量的。但 `parseAutoRaw` 需要支持嵌套解析、`applyDefaults` 需要处理嵌套 `raw` 追踪、`autoModeField` 需要区分子结构体和 ModeToggle、`getAutoKeyValue`/`setAutoConfigValue` 需要支持三段式路径。这些不是四个独立的简单改动，而是对配置系统核心路由逻辑的修改，改动之间有耦合。1-2 小时的估算还排除了 `config_test.go` 和 `config_schema_test.go` 的测试更新，这些测试需要覆盖嵌套路径的序列化/反序列化、默认值填充、以及向后兼容场景。

## Improvement Suggestions

建议：考虑使用扁平命名空间替代嵌套结构体
Addresses: 嵌套命名空间破坏现有解析逻辑的风险
> 具体改动：使用 `auto.evalProposal`、`auto.evalPrd`、`auto.evalUiDesign`、`auto.evalTechDesign` 作为四个独立的 flat ModeToggle 字段，直接挂在 `AutoConfig` 上。这完全复用现有的 `parseAutoRaw`/`autoModeField`/`getAutoKeyValue`/`setAutoConfigValue` 路径，无需引入嵌套解析。YAML 可读性略有下降（`auto.evalProposal` vs `auto.eval.proposal`），但消除了对配置核心路由逻辑的侵入式修改。提案的 Alternatives 表格中已经评估过这个方案并以 "命名不直观" 为由拒绝，但 "直觉" 是主观判断，而配置系统的一致性是客观约束。扁平方案将改动范围从"配置系统核心 + 4 个 skill"缩减到"4 个新字段注册 + 4 个 skill"，开发时间和回归风险大幅降低。

建议：补充 skill 中 config check 的具体实现规格
Addresses: skill 中 config check 逻辑实现方式未说明的问题
> 具体改动：在提案的 Scope 部分明确定义 config check 的实现机制。例如：在每个 skill 的 eval 触发步骤前增加一个 Bash 调用 `forge config get auto.eval{Proposal|Prd|UiDesign|TechDesign}.{quick|full}`，根据返回值决定跳过还是执行 AskUserQuestion。同时定义 agent 如何判断当前是 quick 还是 full 模式（例如通过环境变量、命令行参数、或 `forge config get` 查询）。没有这个规格，"4 个 skill 增加 config check 逻辑" 是不可执行的。

建议：为 ui-design 的默认行为变更提供迁移路径和用户通知机制
Addresses: ui-design 行为变更的向后兼容性风险
> 具体改动：有两种路径。路径 A：将 `uiDesign` 默认设为 `quick: true, full: true`（保持当前无条件自动行为），让用户显式关闭而非显式开启。这保留了向后兼容，代价是 ui-design 与其他三个 skill 的默认值不对称——但这种不对称恰恰反映了 ui-design 当前行为的特殊性。路径 B：如果坚持 `false/false` 默认值，应在 `ReadConfig` 或 `applyDefaults` 中检测首次使用新配置版本的升级场景，输出一次性提示（类似 `migrateOldE2eTestKey` 的模式），告知用户 ui-design 评估行为已变更以及如何恢复。

建议：在 proposal 中显式列出 `AutoConfig` Go struct 的变更伪代码
Addresses: 时间估算和工程复杂度问题
> 具体改动：在提案的 Technical Feasibility 或 Scope 部分增加一段伪代码，展示 `AutoConfig` struct 的新定义、`AutoConfigDefaults()` 的变更、`autoModeField` 的变更、以及 `parseAutoRaw` 的变更。这让评审者和实现者都能准确评估改动的影响面。当前的描述停留在概念层面（"新增 `Eval` 嵌套结构体"），但 Go 的 YAML 序列化行为、零值语义、和 dot-notation 路由逻辑都需要具体代码才能验证。

建议：增加 `forge config get auto.eval` （无子字段）的行为定义
Addresses: 三段式路径行为边界未定义的问题
> 具体改动：在提案中明确定义 `auto.eval` 本身是否是合法的配置路径。如果不合法，应明确说明 `forge config get auto.eval` 返回错误信息。如果合法，定义其输出格式（例如返回所有 eval 子字段的汇总）。同样定义 `forge config set auto.eval true` 的语义（是否等同于设置所有四个子字段为 true）。这些边界情况不影响核心功能，但影响 CLI 的用户体验一致性。
