# Freeform Expert Review: gen-journeys/gen-contracts 自动生成测试流水线任务

**Reviewer**: Test Pipeline Architect
**Date**: 2026-05-23
**Document**: `docs/proposals/auto-gen-journeys-contracts/proposal.md`

---

## Section 1: Background Assessment

这份提案解决的核心问题是：Forge 测试流水线中存在自动化断口——`eval-journey` 和 `eval-contract` 已作为 `autogen.go` 中的自动生成任务存在，但其上游的 `gen-journeys` 和 `gen-contracts` 仍需手动调用。提案计划将这两个 skill 升级为 autogen 任务，补全 Breakdown 和 Quick 两种模式的端到端流水线。

提案的架构思路是自然的：复用现有 `autogen.go` 的任务定义、依赖链解析（`findTaskIndex`/`findTaskIndexByPrefix`）和 `embed.FS` 模板机制，新增两个任务类型 `test.gen-journeys` 和 `test.gen-contracts`，插入到已有依赖链的正确位置。Breakdown 模式获得完整链路（含 eval 质量关卡），Quick 模式跳过 eval 关卡直接走 gen-journeys -> gen-contracts -> gen-scripts -> run -> verify。

提案的一个重要设计决策是用拆分的独立任务替换 Quick 模式的 `gen-and-run` 组合任务，使得两种模式共享 gen-contracts/gen-scripts/run 的任务定义，仅通过 eval 关卡的有无来区分。这是一个合理的简化——减少了任务类型的数量，同时让 Quick 模式也能受益于 Journey/Contract 的结构化流程。

提案假设 gen-journeys 的 SKILL.md 可以通过修改支持 proposal.md 作为替代输入源（Quick 模式没有 PRD），这是一个有实际复杂度的设计点，因为 SKILL.md 当前的 Step 1 明确依赖 `prd-user-stories.md` 作为 primary source。

---

## Section 2: Key Risk Identification

### 依赖链索引偏移风险——当前代码使用硬编码索引

风险：提案声明「resolveBreakdownDeps 使用 findTaskIndex 按ID查找，不依赖硬编码索引」，暗示新增任务不会引入索引偏移问题。然而当前 `resolveBreakdownDeps` 的 Breakdown E2E 分支实际上同时使用了硬编码索引和 `findTaskIndex`——E2E pipeline 的前几个任务（eval-journey=0, eval-contract=1, gen-scripts 从索引 2 开始）全部使用硬编码偏移量计算，只有后续的 validate/specs/clean-code 任务才用 `findTaskIndex`。

引用原文：「resolveBreakdownDeps 使用 findTaskIndex 按ID查找，不依赖硬编码索引」

实际代码（`autogen.go` L425-449）：
```go
evalJourneyIdx := 0
evalContractIdx := 1
genStart := 2
```

在 eval-journey 前面插入 gen-journeys（索引 0），会导致 eval-journey 变为索引 1，eval-contract 变为索引 2，gen-start 变为 3——所有硬编码索引全部失效。提案需要在实现时将所有硬编码索引迁移到 `findTaskIndex` 调用，但提案未明确指出这一迁移工作量，且在 Key Risks 中给出了错误的缓解描述。

### Quick 模式依赖链需要全面重写

风险：当前 Quick 模式的 `resolveQuickDeps` 假设 E2E 任务序列是 `gen-and-run[0..n-1] -> verify[n]`。替换为拆分的独立任务后，序列变为 `gen-journeys -> gen-contracts -> gen-scripts[0..n-1] -> run -> verify`，依赖拓扑从「全并行 -> 单 verify」变为「全串行的线性链 -> 并行 gen-scripts -> 串行 run -> verify」。`resolveQuickDeps` 的逻辑需要从头重写，且新的 Quick 依赖链本质上与 Breakdown 的 E2E 分支相同（只是少了 eval-journey 和 eval-contract）。

引用原文：「Quick 模式用拆分的独立任务替换 gen-and-run 组合任务」

这意味着 Quick 模式的 `resolveQuickDeps` 不能简单修改，而需要引入一种共享依赖链解析逻辑或参数化函数。提案未讨论这一架构变化。

### gen-contracts SKILL.md 的前置条件与自动任务模式冲突

问题：gen-contracts 的 SKILL.md 在 Prerequisites 中要求：「Eval report for all Journeys (`testing/<journey>/.eval-report.md`) | Run `/eval --type journey` first. **Blocker**: do not proceed if any Journey scored below target.」

这意味着 gen-contracts 的 SKILL.md 把 eval-journey 的输出视为阻塞前置条件。但在 Quick 模式中，提案明确跳过 eval 关卡（「Quick 模式（跳过 eval 质量关卡）」）。如果 Quick 模式的自动任务要调用 gen-contracts skill，gen-contracts SKILL.md 的硬性前置检查会失败或产生警告。提案提到要「修改 gen-journeys SKILL.md 支持 proposal.md 作为替代输入源」，但没有提到需要同步修改 gen-contracts SKILL.md 以支持跳过 eval 前置条件。

### gen-journeys 作为自动任务时的 SKILL.md 兼容性

问题：gen-journeys SKILL.md 的 Step 6（Review & Commit）声明了 HARD-RULE：「Do NOT commit documents automatically. Present all generated Journey files to the user for review and wait for explicit approval before committing.」

当 gen-journeys 作为自动任务由 `run-tasks` dispatch 执行时，这个 HARD-RULE 与自动化执行模式存在语义冲突。自动任务期望非交互式执行，但 SKILL.md 要求等待用户审批。提案的 Non-Functional Requirements 仅声明了 `MainSession=false`，但没有讨论 SKILL.md 的交互式 HARD-RULE 如何适配自动任务模式。这不是一个简单的「不改 SKILL.md 核心生成逻辑」（提案 Out of Scope 声明）能解决的问题。

### Quick 模式 gen-journeys 缺少 proposal.md 输入的具体设计

问题：提案声明「Quick 模式 proposal 输入：gen-journeys 需要支持 proposal.md 作为替代输入源（Quick 模式没有 PRD user stories）」，但 SKILL.md 的 Step 2（Identify User Workflows）核心逻辑完全围绕 Given/When/Then 格式的 user stories 构建。proposal.md 的「Key Scenarios」部分并非以 Given/When/Then 格式编写，需要完全不同的提取策略。

引用原文：「proposal.md 包含 scope、success criteria、key scenarios，足以提取精简版 Journeys。质量由下游 gen-contracts 和 gen-test-scripts 的代码侦察来补偿」

这依赖了下游补偿而非上游质量保证，是一种已知的降级策略，但提案未定义降级的具体边界——如果 proposal.md 的 Key Scenarios 信息量极低（例如只有一句话的描述），gen-journeys 的输出可能不构成有意义的 Journey 文档，后续的质量补偿也无从谈起。

### `T-test-gen-journeys` ID 命名与 `isTestTaskID` 检测的一致性

风险：当前 `isTestTaskID` 函数通过 `strings.HasPrefix(id, "T-test-")` 检测测试管道任务。如果新增任务的 ID 为 `T-test-gen-journeys` 和 `T-test-gen-contracts`，它们会正确匹配 `T-test-` 前缀。但提案 Success Criteria 中提到「`forge -h` 输出的类型列表包含 test.gen-journeys 和 test.test-gen-contracts」——注意第二个是 `test.test-gen-contracts`，多了一个 `test.`，这看起来是 typo。如果实际类型名是 `test.gen-contracts` 而非 `test.test-gen-contracts`，那么 Success Criteria 中的文本需要修正。如果是 `test.test-gen-contracts`，则命名规范不一致（gen-journeys 是 `test.gen-journeys`，gen-contracts 却是 `test.test-gen-contracts`）。

引用原文：「`forge -h` 输出的类型列表包含 test.gen-journeys 和 test.test-gen-contracts」

### `autogenTypeToFile` 映射和 embed 模板文件命名

问题：提案声明「新增 embed 模板文件（test-gen-journeys.md, test-gen-contracts.md）」，这遵循了现有命名规范（`test-gen-scripts.md`, `test-gen-and-run.md`）。但当前 `autogenTypeToFile` 映射中 `TypeTestGenAndRun` 对应 `data/test-gen-and-run.md`，提案要替换 gen-and-run 的生成逻辑但保留类型定义——这意味着 `data/test-gen-and-run.md` 模板仍需存在于 embed.FS 中（否则编译失败），但不再有新任务使用它。这不是阻塞问题，但需要确认模板文件的处置策略（保留 vs 标记 deprecated vs 注释说明）。

引用原文：「Out of Scope: 删除 `TypeTestGenAndRun` 类型（保留以防历史引用，标记 deprecated）」

### Breakdown 模式 eval-journey 的前置依赖变化

问题：当前 `ResolveFirstTestDep` 的 breakdown 分支中，`firstTestIdx` 查找的是 `T-eval-journey` 作为第一个测试任务。提案将 gen-journeys 插入到 eval-journey 之前，使得第一个测试任务变为 `T-test-gen-journeys`。`ResolveFirstTestDep` 需要更新查找逻辑，但提案只在 Scope 中提到「更新 ResolveFirstTestDep 适配新的首个测试任务」，没有讨论 Quick 模式的 `firstTestIdx` 也需要从 `T-quick-gen-and-run` 更改为 `T-test-gen-journeys`（或新的 Quick 模式首任务 ID）。

引用原文：「更新 ResolveFirstTestDep 适配新的首个测试任务」

Quick 模式当前通过 `findTaskIndexByPrefix(tasks, "T-quick-gen-and-run")` 查找首任务。替换 gen-and-run 后，Quick 模式的首任务变为 gen-journeys（ID 可能是 `T-test-gen-journeys`，不含 `T-quick-` 前缀），这意味着 Quick 分支的查找逻辑也需要同步更新。

---

## Section 3: Improvement Suggestions

建议：将 `resolveBreakdownDeps` 中 E2E 管道的硬编码索引全部迁移为 `findTaskIndex` 调用，作为新增任务之前的前置重构。

这直接缓解第一个风险。具体做法是：将 `evalJourneyIdx := 0` 替换为 `findTaskIndex(tasks, "T-eval-journey")`，将 `evalContractIdx := 1` 替换为 `findTaskIndex(tasks, "T-eval-contract")`，将 `genStart` 改为基于 `findTaskIndexByPrefix(tasks, "T-test-gen-scripts")` 的计算。同样地，`runIdx` 和 `verifyIdx` 应使用 `findTaskIndex` 而非算术偏移。完成此重构后，新增 gen-journeys/gen-contracts 任务不会导致任何索引偏移问题，且未来的任务插入也不再需要手动调整索引。采纳后，提案的实现风险从「需要仔细调整多处硬编码索引」降低到「只需在正确位置 append 新任务」。

建议：抽取共享依赖链解析函数，避免 Quick 和 Breakdown 模式的重复逻辑。

这缓解第二个风险。当前 Breakdown 和 Quick 的依赖解析是两个完全独立的函数。替换 gen-and-run 后，Quick 模式的 E2E 依赖链变为 `gen-journeys -> gen-contracts -> gen-scripts[per-type] -> run -> verify`，这本质上是 Breakdown 链去掉 eval-journey 和 eval-contract 的子集。可以提取一个参数化的 `resolveE2eChain(tasks []AutoGenTaskDef, includeEval bool)` 函数，两种模式共享核心逻辑。采纳后，提案的代码变更更集中，且两种模式的依赖链一致性更容易保证。

建议：为 Quick 模式的 gen-contracts 添加 eval-skip 的显式支持路径。

这解决第三个问题。具体做法是在 gen-contracts 的 SKILL.md Prerequisites 表格中增加一个条件分支：当通过自动任务模式执行（且模式为 Quick）时，跳过 eval-report 前置检查，改为在任务模板（`data/test-gen-contracts.md`）中明确标注「Quick 模式跳过 eval 质量关卡，Journey 质量由上游 gen-journeys 保证，下游 gen-contracts 的代码侦察补偿」。这需要在 embed 模板中使用 `{{MODE}}` 占位符来区分两种模式的执行指令。采纳后，Quick 模式的自动任务不会因前置条件检查而失败。

建议：在提案 Scope 中明确 SKILL.md 的 HARD-RULE 适配策略。

这解决第四个问题。当前 Out of Scope 声明「不修改 gen-journeys/gen-contracts 的核心生成逻辑」，但 Step 6 的交互式审批 HARD-RULE 不属于「核心生成逻辑」，它属于执行控制流。建议将 SKILL.md 的交互式审批规则的适配纳入 In Scope，具体做法是为自动任务模板（`data/test-gen-journeys.md`）编写独立的执行指令，绕过 SKILL.md 的交互式等待逻辑。自动任务模板通过 embed.FS 渲染，其内容与 SKILL.md 的执行路径不同，因此不违反「不修改 SKILL.md」的约束。采纳后，自动任务模式的 gen-journeys 可以无交互地完成执行。

建议：为 proposal.md 输入源定义最低信息量阈值（fall back guard）。

这缓解第五个风险。具体做法是在 gen-journeys 的任务模板或 SKILL.md 的 proposal 输入分支中定义显式的最低要求：proposal.md 必须包含至少 N 个 Key Scenarios（建议 N=2）和至少 M 条 Success Criteria（建议 M=1），否则 gen-journeys 任务应 fail 并给出明确的诊断信息，而不是生成低质量的 Journey 文档后依赖下游补偿。采纳后，Quick 模式的 gen-journeys 在信息不足时会快速失败而非静默降级，减少了返工成本。

建议：在 Key Risks 表中补充 `resolveBreakdownDeps` 硬编码索引的实际状态描述。

这修正提案中关于依赖链索引风险的误导性描述。当前描述「resolveBreakdownDeps 使用 findTaskIndex 按ID查找，不依赖硬编码索引」与实际代码不符。应更新为：「resolveBreakdownDeps 的 E2E 管道分支使用硬编码索引（evalJourneyIdx=0, evalContractIdx=1, genStart=2），新增任务前需迁移为 findTaskIndex。后续的 validate/specs/clean-code 任务已使用 findTaskIndex。」采纳后，实现者不会因为提案的错误描述而忽略索引迁移工作。

建议：统一 Quick 模式新任务的 ID 前缀策略。

当前 Quick 模式的测试任务使用 `T-quick-` 前缀（`T-quick-gen-and-run-<type>`, `T-quick-verify-regression`），而提案计划将 Quick 模式的 gen-journeys/gen-contracts 作为 `T-test-gen-journeys`/`T-test-gen-contracts`——与 Breakdown 模式共享相同的 ID。这导致 `isTestTaskID` 中 `T-quick-` 前缀检测对新增的 Quick 模式首任务不再适用，`InferType` 函数也需要添加新的 case 分支。建议明确文档化：Quick 和 Breakdown 模式共享 gen-journeys/gen-contracts 的任务 ID（而非使用 `T-quick-gen-journeys` 独立 ID），并确认 `InferType` 中新增 `T-test-gen-journeys` 和 `T-test-gen-contracts` 的识别规则。采纳后，`isTestTaskID`、`InferType`、`SystemTypes` 的更新范围清晰可追溯。
