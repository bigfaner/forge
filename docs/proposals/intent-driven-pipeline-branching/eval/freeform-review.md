# Freeform Expert Review: Intent-Driven Pipeline Branching

## Section 1: Background Assessment

This proposal addresses a concrete, recurring problem in the Forge autogen pipeline: when `coding.refactor` or `coding.cleanup` tasks flow through the full test pipeline (gen-journeys -> eval-journey -> gen-contracts -> eval-contract -> gen-scripts -> run-tests), the journey generation step is fundamentally broken because gen-journeys is prohibited from reading source code. The result is speculative, unverified claims that degrade through eval iterations rather than improving. The proposal cites the `unify-enum-constants` feature as evidence: three eval rounds produced scores of 466 -> 630 -> 585, trending downward.

The core technical approach is a two-dimensional branching matrix where the first dimension is the existing pipeline mode (Breakdown vs Quick, determined by document presence) and the second dimension is a new "feature intent" (`new-feature`, `refactor`, `cleanup`) determined during the brainstorm phase. The intent value is persisted in `proposal.md` frontmatter and read by `build.go` during `BuildIndex()`. When intent is `refactor` or `cleanup`, the test pipeline segment (journeys, contracts, scripts, evals) is skipped entirely; validation falls back to the existing quality-gate hook (compile + fmt + lint + test).

The proposal also describes adaptive PRD behavior: when intent is `refactor`, `write-prd` skips user stories and generates a "spec-only" format containing change scope, constraints, and verification criteria. Tech-design similarly shifts focus to internal architecture rather than API surface.

The approach rests on three assumptions: (1) the intent field can be reliably inferred by AI and confirmed by the user, (2) the quality-gate hook is sufficient to verify that refactors and cleanups do not introduce regressions, and (3) the downstream dependency chain can be cleanly rewired to skip the test pipeline without creating orphaned or dead-end tasks.

## Section 2: Key Risk Identification

风险：`IsTestableType()` 的修改范围被低估。提案在 "In Scope" 中写道 "forge-cli `build.go`：`IsTestableType()` 区分行为变更 vs 行为保持"，但当前 `IsTestableType()` 的实现是 `strings.HasPrefix(typ, "coding.") || typ == TypeCleanCode`。将 `coding.refactor` 和 `coding.cleanup` 从 testable 中排除意味着它们不再触发 `needsTestPipeline()` 返回 true，这直接影响 `GenerateTestTasks()` 的调用路径、stage-gate 生成（步骤 6.5 的 `if needsTest` 条件）、以及 `resolveTestDepsAndInjectReviewDoc()` 的执行。提案聚焦于 `build.go` 和 `autogen.go`，但 stage-gate 生成的跳过是一个连带效果，提案没有明确讨论 refactor/cleanup 特征是否仍然需要 phase gates。如果 refactor 任务有 5 个 coding.refactor 业务任务跨 3 个 phase，跳过 stage-gate 可能导致任务调度缺少串行化保障。

问题：refactor 的 Breakdown 管道描述存在依赖链断裂风险。提案描述 refactor（Breakdown）的管道为 "business tasks -> validate-code -> clean-code -> consolidate"，但当前 `resolveBreakdownDeps()` 中 `T-validate-code` 的依赖是 `lastRunID`（最后一个 run-test 任务的 ID）。当 refactor 跳过整个 test pipeline 时，`lastRunID` 将为空字符串，`T-validate-code` 不会获得任何依赖。提案提到 "下游任务（validate-code, clean-code, consolidate-specs）的依赖链需重新挂载"，但没有给出具体的接线逻辑——特别是，这些任务应该挂载到最后一个 business task、最后一个 gate，还是最后一个 summary？当前 `ResolveFirstTestDep()` 的逻辑对于 breakdown 模式是 `findHighestGateOrSummary()` 或 `findMaxBusinessTaskID()` 取较高者，但这个函数只在第一个测试任务上调用，不会影响 validate-code/clean-code 的依赖。

风险：Quick+refactor 管道中 `T-clean-code` 和 `T-quick-doc-drift` 的依赖来源不明确。提案描述 Quick+refactor 为 "proposal -> quick-tasks -> quality-gate -> done"，但在 autogen.go 的 `GetQuickTestTasks()` 中，当 `auto.Test.Quick` 为 true 时仍然会生成 gen-journeys 和 run-test 任务。提案说 refactor/cleanup 应该跳过这些，但如果 `auto.Test.Quick` 被禁用来跳过测试任务，那么 `resolveQuickDeps()` 中的 `lastRunID` 为空，`T-quick-doc-drift` 和 `T-validate-code` 的依赖设置被跳过。这恰好是 `ResolveDriftFallbackDep()` 存在的原因（提案确实引用了这个 fix），但 fallback 只处理 `T-quick-doc-drift` 和 `T-specs-consolidate`，不处理 `T-validate-code` 和 `T-clean-code`。如果 refactor Quick 模式仍然生成 validate 和 clean-code 任务（只是不生成 test 任务），它们的依赖链会是空的。

问题：Intent 字段在 `build.go` 中的注入时机与现有模式检测存在耦合。提案写道 "`BuildIndex()` 在 `detectMode()` 之后读取 `proposal.md` frontmatter 的 `intent` 字段"。当前 `detectMode()` 通过检查 `prd/prd-spec.md` 或 `docs/proposals/{slug}/proposal.md` 的存在来决定模式。但 proposal.md 的 intent 字段需要在 `needsTestPipeline()` 调用之前可用，这意味着需要在步骤 2（detectMode）和步骤 5.5.1（needsTestPipeline）之间添加 intent 读取逻辑。然而，proposal.md 的路径解析依赖 `setFeatureMetadata()`（步骤 3）的结果——特别是 `index.Proposal` 字段。当前 `detectMode()` 和 `setFeatureMetadata()` 各自独立构造路径（都硬编码了 `docs/proposals/{slug}/proposal.md`），如果 intent 读取在步骤 3 之后，需要从 `index.Proposal` 字段或直接从已知路径读取。这个时序问题虽然可解，但提案没有讨论。

风险：`coding.fix` 类型被忽略了。当前 `IsTestableType()` 对所有 `coding.*` 前缀返回 true，包括 `coding.fix`。Bug 修复同样可能不引入新的用户可观测行为（纯修复回归的 fix），但提案只区分了 `new-feature`、`refactor`、`cleanup` 三种 intent，没有为 `fix` 类型定义管道行为。提案的矩阵中，如果用户提交一个 intent 为 `refactor` 的 proposal 但 breakdown-tasks 生成了 `coding.fix` 类型的任务（因为重构过程中发现了 bug），`IsTestableType()` 的修改可能导致不一致：intent 说 refactor（不生成测试），但实际任务类型是 fix（通常应该有测试）。proposal 写道 "intent 在 proposal 阶段确定后，驱动整个 pipeline 拓扑选择"，但 pipeline 拓扑选择基于的是 proposal 的 intent 而非实际任务的 type，这种错配可能导致本应测试的 fix 任务被跳过。

风险：spec-only PRD 与下游 skill 的兼容性未经验证。提案在 "Key Risks" 表中承认 "refactor PRD spec-only 格式与下游 skill 不兼容" 的可能性为 M，影响为 M，但缓解措施只是 "write-prd 分支确保 spec 格式包含 tech-design 需要的字段"。当前 `tech-design` skill 会读取 `prd-spec.md` 中的 user stories 来推导 API surface 和 ER 图。如果 spec-only PRD 不包含 user stories，tech-design 的 prompt 模板可能产生空输出或错误。提案提到 tech-design 也会做 "内部分支"（refactor 侧重内部架构），但这是在 skill 层面的行为引导，依赖于 skill.md 的指令质量——没有结构性的保障来确保 tech-design 模板能正确处理缺少 user stories 的输入。

问题：Breakdown+cleanup 组合被标记为 "不适用"，但缺少防止用户触发此组合的机制。提案在矩阵中标注 cleanup 的 Breakdown 模式为 "(不适用，cleanup 不走 Breakdown 模式)"，但没有描述当用户对一个 cleanup proposal 执行 `write-prd` 时会发生什么。如果 brainstorm 的 AI 推断 intent 为 cleanup，但用户随后手动运行 `write-prd`（而非 `quick-tasks`），系统不会阻止这个操作。这会创建一个隐式的无效路径。提案说 "始终走 Quick 模式"，但没有说明是强制性的（代码层面阻止）还是引导性的（skill 指令层面建议）。

风险：`needsTestPipeline()` 的改动会影响混合类型特征。如果一个特征同时包含 `coding.feature` 和 `coding.refactor` 任务（例如重构某个模块并添加新功能），`needsTestPipeline()` 只要有任何一个 testable 类型就返回 true。但如果 intent 是 refactor，提案建议整个测试管道跳过——这意味着新增功能的任务也不会获得测试覆盖。提案在 "Out of Scope" 中提到 "混合 intent 支持（一个 proposal 包含多种 intent）——不支持，按主要意图归类"，但 "主要意图" 的判断标准没有定义。当一个 proposal 同时包含重构和新功能时，AI 推断的 "主 intent" 可能是错误的。

风险：quality-gate hook 作为 refactor/cleanup 的唯一验证手段存在覆盖盲区。提案多次依赖 "已有 quality-gate hook（compile + fmt + lint + test）" 作为回归验证的替代方案。但 quality-gate hook 的 test 步骤运行的是已有测试套件——对于纯重构，已有测试通过仅证明行为未变，不证明重构的代码覆盖了所有变更点。如果重构涉及测试本身未覆盖的代码路径，quality-gate 不会检测到问题。对于清理死代码的场景更是如此：移除代码后编译通过、已有测试通过，但可能移除了被其他模块动态引用的代码（如反射调用、依赖注入等），这种错误在编译期不会暴露。

## Section 3: Improvement Suggestions

建议：在 `build.go` 的 `BuildIndex()` 中显式引入 `intent` 参数而非从文件系统读取。当前提案要求在 `detectMode()` 之后读取 `proposal.md` 的 frontmatter 来获取 intent，但这会在 `BuildIndex()` 内引入文件 I/O 来读取一个已经在上层（brainstorm/forge-cli command）可用的信息。更干净的做法是将 `intent` 作为 `BuildIndexOpts` 的一个字段，由调用方在构造 opts 时填充。这与 `FeatureSlug` 和 `ProjectRoot` 的传递模式一致，避免 `BuildIndex()` 内部直接解析 proposal.md frontmatter。调用方已经有了 proposal 路径信息（用于 `detectMode()`），可以在那里一并解析 intent。这种方案同时解决了上面提到的时序问题。

建议：为 refactor/cleanup 的下游任务依赖提供显式的接线函数。当前提案描述 "下游任务（validate-code, clean-code, consolidate-specs）的依赖链需重新挂载"，但没有给出具体实现。建议新增一个 `resolveRefactorDeps()` 函数（或在 `resolveBreakdownDeps()`/`resolveQuickDeps()` 中添加 intent 分支），当 intent 为 refactor/cleanup 时，将 validate-code/clean-code/consolidate-specs 的依赖直接接到最后一个业务任务（通过 `findMaxBusinessTaskID()` 或 `findHighestGateOrSummary()`），而非依赖 `lastRunID`。这确保了即使没有测试管道，下游任务仍然有正确的依赖链。应该同时扩展 `ResolveDriftFallbackDep()` 来覆盖 validate-code 和 clean-code 任务，而不仅仅是 drift 和 consolidate。

建议：将 intent 的作用域从 pipeline 拓扑扩展到 per-task 级别的 testable 判断。当前提案用 intent 控制 pipeline 级别的测试生成（全有或全无），但更精细的做法是让 `needsTestPipeline()` 同时检查 intent 和每个 task 的 type。具体来说：对于 `new-feature` intent，所有 `coding.*` 类型都是 testable（现状不变）；对于 `refactor` intent，只有 `coding.feature` 和 `coding.enhancement` 类型是 testable（`coding.refactor` 和 `coding.cleanup` 排除）；对于 `cleanup` intent，所有 coding 类型都排除。这样即使一个 refactor 特征中混入了少量 `coding.feature` 任务（如修复过程中发现的新功能需求），这些任务仍然能获得测试覆盖。这解决了混合类型特征的风险。

建议：在提案中明确定义 `coding.fix` 的 intent 映射和管道行为。fix 类型在当前系统中是 testable 的，且在 `InferType()` 中有专门的模式匹配（`strings.HasPrefix(id, "fix-")`）。建议增加 `fix` 作为第四种 intent，或者明确定义 fix 类型在各 intent 下的行为：`new-feature` 管道中 fix 任务走完整测试（不变），`refactor` 管道中 fix 任务由 quality-gate 覆盖。这消除了 intent 与 task type 之间的语义间隙。

建议：为 Breakdown+cleanup 组合添加代码级别的防护而非仅依赖 skill 引导。在 `forge task index` 命令中，当检测到 intent 为 `cleanup` 且 mode 为 `breakdown` 时，应该输出明确的 warning 或错误，引导用户使用 Quick 模式。这比依赖 brainstorm skill 的引导更可靠——用户可能跳过 brainstorm 直接手动创建 proposal.md。类似地，`write-prd` skill 在检测到 intent 为 cleanup 时应该提示用户切换到 quick-tasks。
