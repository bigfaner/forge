# Freeform Expert Review

## Background Assessment

This proposal addresses a concrete structural gap in Forge's test generation pipeline: when all business tasks for a feature are surface-type `web` or `mobile`, the `gen-contracts` step produces no usable output for those journeys, and `gen-test-scripts` silently skips them because it finds no Contract files as input. The result is 71% test coverage loss for the milestone-map feature (5 out of 7 journeys) with zero pipeline errors or warnings.

The core insight is sound: the Contract abstraction was designed for protocol-level surfaces (API/CLI/TUI) where inputs and outputs are structured data, and Contract's six-dimension model (Preconditions, Input, Output, State, Side-effect, Invariants) maps naturally onto these surfaces. For interaction-level surfaces (Web/Mobile), the Journey document already contains all the information needed for test generation -- user actions, expected visual states, navigation flows -- and forcing these through a Contract intermediary provides no information gain.

The proposed solution has two layers: (1) at the pipeline level, skip `T-test-gen-contracts` and `T-eval-contract` when all business tasks have interaction-only surface types; (2) at the skill level, give `gen-test-scripts` a direct path from `journey.md` to test scripts for Web/Mobile journeys that lack Contract files, supplemented by a coverage self-check that hard-fails on gaps.

The proposal claims four scenarios (S1-S4) are exhaustive and proposes `CondHasProtocolSurfaceTask` as a new `GenerateCondition` function to detect when contracts should be skipped.

## Key Risks

The most fundamental issue in this proposal is the **asymmetry between the pipeline-level skip condition and the skill-level fallback path**. The pipeline decides to skip gen-contracts based on whether *any* business task has a protocol-level surface type. But the gen-test-scripts skill needs to make a per-journey decision: for a given journey, does a Contract file exist or not? These are two different granularities of decision-making, and the proposal does not fully explore the implications of this mismatch.

风险：Pipeline 层的 `CondHasProtocolSurfaceTask` 判断粒度与 Skill 层的直达路径判断粒度不一致，可能导致混合场景下的逻辑冲突。
> "检查 feature 所有业务任务的 `surface-type` 字段，若不存在 tui/cli/api 类型任务，跳过 `T-test-gen-contracts` 和 `T-eval-contract`" — 这个判断是在 feature 级别做的，但 S3 场景描述了混合 surface feature，其中部分 journey 是 web-only 的。在混合场景下，gen-contracts 会被保留（因为存在 api 任务），但 web-only journey 仍然没有 contract。Skill 层需要独立处理这种情况。这意味着两个层的决策逻辑存在冗余但不完全一致：pipeline 层决定"是否生成 contracts"，skill 层决定"单个 journey 是否需要 contract"。如果两个层的判断逻辑产生不一致（例如 pipeline 层因为某个 api 任务保留了 gen-contracts，但某个 web journey 确实没有对应 contract），skill 层的直达路径必须兜底。Proposal 虽然提到了这种场景，但没有明确声明这两层决策的优先级关系和一致性保证。

风险：覆盖率自检作为硬失败的兜底机制，其触发条件依赖于 `surface-type` 字段的正确填充，但该字段的错误可能导致"假成功"而非"假失败"。
> "覆盖率自检硬失败兜底：即使流水线判断错误，自检会发现缺失的测试脚本" — 这个声明只在测试脚本完全缺失时成立。但如果 `surface-type` 字段被错误地标记为 `api`（实际是 `web`），pipeline 会走 contract 路径，gen-contracts 会尝试为一个纯交互式 journey 生成协议级 contract。这不是"缺失测试脚本"，而是"生成了无意义的 contract，测试脚本基于错误的 contract 生成"。覆盖率自检计数会通过（`count(journeys) == count(test-script-sets)`），但测试质量可能严重下降。 Proposal 的风险表假设了二元失败模式（有/无），忽略了第三种状态：错误类型的测试。

问题：S3 和 S4 场景的边界条件定义不够精确，可能遗漏"部分任务缺失 surface-type 字段"的情况。
> "S3: 混合 surface feature（如 `backend=api, frontend=web`）：部分任务 web、部分 api → gen-contracts 正常生成" — 这里隐含假设所有任务都有 `surface-type` 字段。但在实际数据中，可能存在老任务或手动创建的任务缺少此字段。`CondHasProtocolSurfaceTask` 的行为在遇到 `surface-type` 为空的任务时应该如何处理？如果空字段被当作"没有协议级任务"，会导致 contract 被错误跳过；如果被当作"可能有协议级任务"，则保留了安全性但可能无法跳过本应跳过的场景。Proposal 应明确声明对缺失字段的 fallback 策略。

问题：gen-test-scripts 的前置条件检查当前是"必须有 Contract 文件"，直达路径需要修改这个硬性前置条件，但 proposal 没有详细描述这个检查点的具体修改方式。
> "当 journey 的 `surface_types` 仅含 web/mobile 且无对应 contract 文件时，直接从 journey.md + types/web.md 生成测试脚本，跳过 contract 前置检查" — 当前 gen-test-scripts SKILL.md 的 Prerequisites 部分明确要求"至少一个 Contract 文件"（第28行）。直达路径需要将这个前置条件改为条件性的：对于 protocol-level surface 仍然要求 contract，对于 interaction-level surface 允许跳过。这不仅仅是一个简单的 if-else -- 它影响了整个 Step 2（Read Contract Specifications）的流程，因为后续的 Step 2.5（Load Type Rules）需要从 Contract 的 `step-action` 字段提取 interface type。直达路径需要从 journey.md 而非 contract 文件提取这些信息，但 journey.md 的格式可能不包含等价的结构化字段。这是一个实现层面的重要gap，proposal 没有充分展开。

风险：gen-scripts 依赖链从 gen-contracts 改为 gen-journeys 时，如果 gen-contracts 仍被生成（混合场景），gen-scripts 的依赖解析需要支持动态路由——这在当前 `ResolveUpstream` 机制中可能无法直接表达。
> "依赖链调整：gen-contracts 跳过时，gen-scripts 依赖 gen-journeys / eval-journey" — 当前 PipelineRegistry 中 gen-scripts 的 `DependsOn` 使用 `ResolveUpstream`，它解析为 pipeline 中最近的上游任务。在当前实现中，这意味着 gen-contracts 或 eval-contract（如果存在）。跳过 gen-contracts 后，`ResolveUpstream` 应该自然回退到 eval-journey 或 gen-journeys——这在逻辑上是正确的。但在混合场景中（gen-contracts 被保留），gen-scripts 仍然依赖 gen-contracts，而 web-only journey 的测试生成又不需要 contract。这里没有实际的依赖问题（gen-contracts 处理的是 api journey 的 contract，不影响 web journey 的直达路径），但 proposal 应该明确验证 `ResolveUpstream` 的回退行为确实覆盖了纯交互场景。

风险：types/web.md 和 types/mobile.md 需要补充"直达生成规则"，但当前这两个文件的设计完全围绕 Contract 驱动生成。补充直达规则意味着在同一文件中维护两套生成策略，长期来看可能造成理解负担和分支爆炸。
> "gen-test-scripts 的 types/web.md 和 types/mobile.md 需要补充直达生成规则" — 审查了 types/web.md 和 types/mobile.md 的实际内容，它们目前定义的是 Golden Rules（框架无关约束）和 Reconnaissance Hints（发现辅助），这些规则在直达路径中仍然适用。但 Contract 驱动路径中的某些规则（如 Step 2.5.1 的"从 Contract step-action 字段提取 interface type"、Assertion Depth Enforcement 基于 Contract 的 fixture_spec）在直达路径中缺失输入源。直接在 types/web.md 中添加"当无 contract 时，从 journey.md 提取 X Y Z"的规则，会使该文件的职责从"Web 测试类型约束"扩展为"Web 测试生成路由器"，这是一个架构层面的变化，需要更审慎的设计。

问题：Proposal 声称 TUI 保持 Contract 路径，但未解释为什么 TUI 与 Web/Mobile 的区别足以支撑不同的处理策略。
> "TUI surface 处理方式变更（保持 Contract 路径）" — TUI 被归类为"协议级"，但 TUI 测试实际上也涉及用户交互（输入文字、读取终端输出）。Proposal 的核心论点是"交互级 surface 不需要 contract 中间层"，但 TUI 的交互模型介于 CLI（纯命令行）和 Web（视觉界面）之间。将 TUI 归类为"协议级"的边界条件可能不稳固：如果未来 TUI journey 包含复杂的终端交互场景（如 ncurses 界面），这些场景与 Web 交互的"无信息增益"论证同样适用。Proposal 应该明确 TUI 归类的具体判据，而不是依赖 surface type 作为二元开关。

风险：覆盖率自检的计数逻辑 `count(journeys) == count(test-script-sets)` 在多 surface 项目中需要按 surface 细分，否则会产生误判。
> "生成完成后执行覆盖率完整性自检，缺口硬失败" — 在多 surface 项目中（如 S3 场景），一个 journey 可能同时有 api 和 web 测试脚本。自检应该按 surface type 分别计数：api journey 的测试脚本应该来自 contract 路径，web journey 的测试脚本应该来自直达路径。如果自检只是简单比较 journey 总数和测试脚本总数，可能会掩盖单个 surface 的覆盖率缺口。例如，3 个 api journey 都生成了测试脚本，但 2 个 web journey 都没有生成 -- 总数可能恰好匹配（如果混叠了计数维度），导致假成功。

## Improvement Suggestions

建议：明确 `CondHasProtocolSurfaceTask` 对 `surface-type` 缺失值和空值的处理策略，并编写针对边界条件的单元测试。
Addresses: "部分任务缺失 surface-type 字段"问题
> What changes: 在 Proposal 的 Constraints & Dependencies 部分，增加一条显式规则：当业务任务的 `surface-type` 字段为空或不存在时，`CondHasProtocolSurfaceTask` 应将其视为"存在协议级可能性"（保守策略），即不跳过 gen-contracts。这确保了错误/缺失数据时走安全路径，只有明确的 `surface-type: web` 或 `surface-type: mobile` 才触发跳过。同时在 PipelineRegistry 的 Go 实现中，`CondHasProtocolSurfaceTask` 函数应该有对应的 table-driven test 覆盖所有边界情况（空字段、未知类型、混合类型、全空）。

建议：将 Pipeline 层和 Skill 层的决策逻辑统一为"以 journey 的 contract 文件是否存在为单一真相源"，Pipeline 层的 skip 作为优化（减少无效任务生成），而非安全网。
Addresses: "Pipeline 层与 Skill 层判断粒度不一致"风险
> What changes: Proposal 应该明确声明架构决策的优先级：Skill 层（gen-test-scripts）的直达路径是完整且自足的，它不依赖 Pipeline 层是否跳过了 gen-contracts。Pipeline 层的 skip 纯粹是效率优化（避免运行无用的 gen-contracts 任务）。这意味着即使 Pipeline 层的判断逻辑有 bug（例如在边界条件下未正确跳过 gen-contracts），Skill 层仍然能正确处理 -- 它会找到 contract 目录为空，然后走直达路径。这种设计消除了两层逻辑的一致性耦合。在 Proposal 的 Innovation Highlights 中增加一段明确声明这一架构原则。

建议：在覆盖率自检中增加按 surface type 细分的计数逻辑，而不仅仅是 journey 级别的总数比较。
Addresses: "覆盖率自检计数逻辑在多 surface 项目中可能误判"风险
> What changes: 将自检逻辑从 `count(journeys) == count(test-script-sets)` 改为：对每个 surface type，检查该 surface 下所有 journey 是否都有对应的测试脚本集。具体而言，遍历所有 journey 的 `surface_types` 字段，为每个 surface type 建立一个 journey 列表，然后检查 `tests/<surfaceKey>/<journey>/` 或 `tests/<journey>/` 目录是否存在且包含至少一个测试文件。缺口报告应该按 surface type 列出缺失的 journey 名称。这比简单的总数比较更精确，能避免多 surface 场景下的假成功。

建议：为直达路径设计一个独立的规则文件（如 `rules/direct-generation.md`），而不是将直达生成规则嵌入 types/web.md 和 types/mobile.md。
Addresses: "types/web.md 补充直达规则导致职责扩展"风险
> What changes: 在 gen-test-scripts skill 中新增 `rules/direct-generation.md`，定义直达路径的完整生成流程：如何从 journey.md 提取 interface type 信息（替代 Contract 的 step-action）、如何从 journey.md 的步骤描述生成 fixture 数据（替代 Contract 的 fixture_spec）、如何应用 assertion depth 规则（不依赖 Contract 的 Outcome 定义）。types/web.md 和 types/mobile.md 只需引用这个规则文件（"当无 contract 文件时，按照 rules/direct-generation.md 执行"），保持自身的职责边界（Web/Mobile 特定约束）不被污染。这种分离使得直达路径的逻辑集中维护，避免在多个 type 文件中重复。

建议：在 Proposal 中补充 gen-test-scripts 前置条件修改的详细设计，特别是 Step 2（Read Contract Specifications）和 Step 2.5（Load Type Rules）在直达路径下的替代数据源。
Addresses: "gen-test-scripts 前置条件检查的具体修改方式"问题
> What changes: 在 Scope > In Scope 部分增加一条：gen-test-scripts 的 Prerequisites 检查需要从"必须有 Contract 文件"改为条件性检查：检测 journey 的 surface_types，如果是纯 web/mobile 则允许无 Contract 文件（前置条件降级为只需要 journey.md 存在）。同时，在 In Scope 中明确 Step 2 和 Step 2.5 在直达路径下的具体行为：Step 2 改为从 journey.md 解析步骤结构和 action 信息（替代 Contract 解析），Step 2.5 的 interface type 检测改为从 journey.md 的 surface_types 字段直接获取。这些修改点应该在 Proposal 中就有清晰的声明，而不是推迟到 PRD 或 tech-design 阶段。
