# Freeform Review: 全局文档-代码一致性审计与知识库清理

**Reviewer**: Documentation-Implementation Drift Auditor
**Date**: 2026-06-02
**Proposal**: docs/proposals/global-doc-code-audit/proposal.md

---

## Section 1: Background Assessment

本提案旨在对 Forge 项目进行一次性的三层文档-代码一致性审计。提案者观察到用户文档、规范文档和知识库与实际代码实现之间存在未量化的不一致，且已存在 5 个局部审计提案但均未执行。提案采用分层策略：L1 审计用户文档层、L2 审计规范文档层、L3 审查知识库有效性，每层从"过时/错误"、"缺失"、"冗余"三个维度进行检查。

提案的核心理由是：v3.0.0 分支正处于开发阶段，延迟清理会导致文档-代码漂移持续恶化，而 AI 代理基于过时文档执行的风险日益增加。提案明确将审计范围限定为"只报告不修复"，产出的 Task 由人工确认后执行。

提案的整体思路清晰、层次分明，对问题的定性描述准确。然而，在事实基础、范围边界和工作量估算等方面存在若干需要审视的问题。

---

## Section 2: Key Risk Identification

### 事实基础不精确

问题：提案 Evidence 部分声称 "docs/conventions/ 下 16+ 份规范文档"，但实际审计发现 `docs/conventions/` 顶层目录包含 15 个 .md 文件和 1 个 `testing/` 子目录（含 5 个 .md 文件 + index.md），合计约 20 份。此外 `docs/business-rules/` 下还有 4 份。提案将这两个目录合并在 L2 层审计，但只引用了 conventions 的数量，遗漏了 business-rules 的 4 份。这种不精确在审计提案中格外刺眼——一个关于审计准确性的提案自身的数据就不准确。

问题：提案声称 "docs/decisions/（10条）"，实际 `docs/decisions/` 目录下确实有 10 个文件，数量吻合。但 "docs/lessons/ 下积累了 133 条经验教训"——实际 `ls | wc -l` 结果确为 133，吻合。然而提案在 Comparison Table 中写道 "执行现有5个提案的86个task"，这个数字无法从提案本身验证，读者需要逐一打开那 5 个提案才能确认。作为提案的事实基础，关键数据应当自包含或附注来源。

### 范围边界问题

问题：提案 Scope 部分将 "docs/features/（183个 feature 目录）的清理" 标记为 Out of Scope。实际文件系统中有 182 个 feature 子目录，而非 183。数字偏差不大，但这再次指向一个更根本的问题：提案完全排除了 docs/features/ 目录的审计。features 目录是 Forge 工作流的核心产出区域，包含 manifest.md、prd/、design/、tasks/ 等子结构。这些文档是 AI 代理执行任务时直接读取的输入，如果其中包含过时信息（比如引用了已废弃的 CLI 命令格式），其危害远大于 README 中的描述不一致。将如此关键的目录完全排除在审计之外，缺乏充分论证。

风险：提案将 "docs/proposals/（204个 proposal）的清理" 标记为 Out of Scope，实际有 181 个 proposal 子目录。数量偏差虽小，但 proposals 目录与 features 目录类似——已采纳的 proposal 会驱动后续工作流，过时的 proposal 中的技术假设可能误导后续开发。完全排除意味着一个重要的不一致来源被忽视。

### 工作量估算缺乏依据

问题：提案在 Resource & Timeline 部分估算 "L1 用户文档层：约 15-20 文件，预计 2-3 个 Task"。实际的 L1 文件清单包括：README.md（1）、ARCHITECTURE.md（1）、docs/user-guide/ 下 4 个文件、docs/official-references/ 下 5 个文件。合计 11 个文件，而非 15-20。类似地，L2 层：docs/business-rules/ 4 个文件 + docs/conventions/ 约 20 个文件 + docs/reference/ 1 个文件 = 约 25 个文件，而非提案声称的 "约 21 文件"。这些偏差虽不至于影响提案方向，但让人质疑估算方法——这些数字是基于什么得出的？是否有遗漏或重复计算？

### 成功标准不可验证

问题：提案的 Success Criteria 包含 "发现的不一致问题 100% 记录在报告中"。这是一个不可验证的声明——如何证明已经发现了 100% 的不一致？审计的完备性本身无法在审计框架内证明。更合理的措辞应该是 "对范围内所有目标文件逐一完成审计流程"（过程性标准），而非 "100% 发现"（结果性标准）。

风险：成功标准 "所有问题已转化为可执行 Task，每个 Task 可由 task-executor 独立执行" 暗示所有审计发现的问题都应转化为 task-executor 可执行的 Task。但 L3 知识库审查的核心结论是 "标记过时/重复/无参考价值"，这类判断性标注并不适合自动化 Task 执行。提案虽然声明 "知识库清理需人工确认"，但成功标准中仍然要求所有 Task 可由 task-executor 独立执行，两者存在矛盾。

### 提案与现有工具的关系处理不清

问题：提案在 Comparison Table 中将 "增强现有工具 /consolidate-specs" 列为 Rejected，理由是 "开发成本高，当前需快速解决" 和 "本质是一次性工作"。然而 `consolidate-specs` skill 已经具备 "检测规范漂移" 的功能（其非交互模式可自动检测规范与代码的 drift）。提案将一次性的全量审计与工具化的持续 drift 检测对立起来，但两者并非互斥——可以先做一次全量审计建立基准，再通过 consolidate-specs 持续维护。提案未考虑这种组合方案。

风险：提案声明 "不修改任何代码或文档，只生成报告和 Task" 作为约束。这意味着审计产出的是一堆待执行 Task，而这些 Task 的执行又会产生新的文档变更。如果在 Task 执行期间代码继续变化（提案也承认这个风险），生成的 Task 可能需要重新审计。这是一个潜在的循环：审计 -> 生成 Task -> 执行 Task（改变文档/代码）-> 文档可能再次不一致。

### 分发约束与审计范围的张力

风险：提案 L2 层审计 `docs/conventions/` 下的文件。但 `forge-distribution.md` 明确指出 `docs/conventions/` 由 `/consolidate-specs` 从 feature 文档中提取，agent 在任务执行时读取这些规范。这意味着 conventions 目录本身不是原始来源，而是从 features 目录中提取的派生产物。如果 conventions 内容与代码不一致，根本原因可能在 features 目录中的源文档。但提案将 features 目录排除在审计范围之外，这导致 L2 审计可能只能发现表面症状而非根因。

### Next Steps 过于简略

问题：提案 Next Steps 仅写 "Proceed to `/quick-tasks` to generate tasks directly from this proposal"。对于如此大规模的审计工作，缺失了关键的执行计划信息：审计的执行顺序是什么？三个层级是否有依赖关系？审计代理的具体工作协议是什么（逐文件读取 vs. 交叉比对 vs. 搜索模式匹配）？这些信息对于 task 生成至关重要。

---

## Section 3: Improvement Suggestions

建议：修正提案中的数量事实。具体而言：将 "docs/conventions/ 下 16+ 份规范文档" 改为更精确的 "docs/conventions/ 下 15 份规范文件 + testing/ 子目录下 6 份，共约 21 份"，将 L1 文件数修正为 11 个，将 features 目录数量修正为 182 个，将 proposals 数量修正为 181 个。对于一个以准确性为生命的审计提案，自身数据必须精确。

建议：重新审视 docs/features/ 的排除决定。建议至少将 docs/features/ 下活跃 feature（manifest.md 状态为 in-progress 或 tasks 的）纳入 L1 审计范围，因为这些文档直接影响当前 AI 代理的行为。已完成或废弃的 feature 可以排除。这样既控制范围，又覆盖了最高风险的区域。

建议：将成功标准从结果性标准改为过程性标准。将 "发现的不一致问题 100% 记录在报告中" 改为 "对范围内每个目标文件完成以下审计步骤：(1) 提取文档中的所有事实性声明；(2) 逐一与代码/配置验证；(3) 记录所有不一致"。同样，将 "每个 Task 可由 task-executor 独立执行" 改为区分两类 Task：自动化可执行的修复类 Task，和需人工判断的知识库审查类 Task。

建议：在提案中补充审计方法论的简要说明。当前提案只说了 "AI 代理直接审计即可"，但未描述具体的审计流程。建议至少说明：(1) 每个文件提取事实性声明的方法；(2) 交叉比对的验证策略（如 CLI 命令文档 vs. Go 代码中的 cobra 命令定义、目录结构文档 vs. 实际目录）；(3) 严重级别 P0-P3 的分类标准。这些信息应在提案中明确，而非留给 task 生成阶段。

建议：考虑将 L3 知识库审查与 L1/L2 审计分离为独立阶段。知识库审查（133 条 lessons + 10 条 decisions 的有效性判断）与文档-代码一致性审计（事实比对）在方法论上有本质差异——前者是主观判断，后者是客观比对。合并在一起可能导致审计 Task 粒度不均，也增加了"审计期间代码变化导致结果过时"的风险。先完成 L1/L2 的客观审计，再基于审计结果中发现的已废弃代码路径，来指导 L3 的知识库审查，会更有效率。

建议：在提案中明确与 consolidate-specs 的后续衔接方案。建议增加一个 Next Step：审计完成后，基于审计结果建立的基准，通过 consolidate-specs 的 drift 检测功能建立持续监控机制，防止未来再次漂移。这将一次性审计的价值延伸为持续保障。
