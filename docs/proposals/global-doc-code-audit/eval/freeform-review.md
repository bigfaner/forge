# Freeform Review: 全局文档-代码一致性审计与知识库清理

**Reviewer**: Documentation Accuracy & Distribution Integrity Specialist
**Date**: 2026-06-03
**Proposal**: docs/proposals/global-doc-code-audit/proposal.md

---

## Section 1: Background Assessment

本提案提出对 Forge 项目执行一次三层系统性文档-代码一致性审计：L1 覆盖用户文档（README.md、DESIGN.md、ARCHITECTURE.md 及 docs/user-guide/、docs/official-references/），L2 覆盖规范文档（docs/conventions/、docs/business-rules/、docs/reference/），L3 覆盖知识库（docs/lessons/ 的 133 条和 docs/decisions/ 的 10 条）。提案以 v3.0.0 发布前清理为时间驱动，承诺产出结构化问题报告和可执行 Task，自身不执行任何修复。

提案在问题定性上是准确的——已有的 5 个局部审计提案（plugin-consistency-audit、skill-ecosystem-audit、skill-instruction-audit、prompt-template-audit、test-pipeline-consistency-audit）确实仅覆盖了 plugin 层和 test pipeline，未涉及用户文档和知识库。提案从"AI 代理基于过时文档执行错误操作"这一实际危害出发论证审计必要性，这一动机链是合理的。

提案的分层审计方法论（声明提取 → 代码定位 → 逐条比对 → 缺失检测 → 结果记录）和层级间反馈机制设计体现了一定的结构化思维。P0-P3 的严重级别分类标准从"影响 AI 代理执行破坏性操作"到"措辞不一致"的降级逻辑也是合理的。

然而，提案在事实基础、范围完整性、与现有工具生态的衔接方面存在一系列需要审视的问题，其中部分问题如果按原样执行将直接导致审计结果不可靠。

---

## Section 2: Key Risks

### 事实声明与代码库实际状态不符

问题：提案 Evidence 部分声称 "docs/conventions/ 下 22 份规范文档（顶层15份 + testing/7份）"。实际验证结果：`docs/conventions/` 顶层有 15 个 .md 文件，`testing/` 子目录仅有 3 个 .md 文件（index.md、cli/index.md、cli/core.md），合计 18 份而非 22 份。差额为 4 份。对于一份以"一致性审计"为主题的提案，自身引用的基础数据就与实际不符，这严重削弱提案的可信度——如果连目录下的文件数量都数错了，如何保证审计流程的完备性？

问题：提案 Resource & Timeline 部分声称 "docs/reference/ 仅含 1 个文件 test-type-model.md，不单独建 Task，合并到其他 L2 Task 中一并审计"。但 `docs/reference/` 目录在文件系统中根本不存在。实际的 `test-type-model.md` 位于 `plugins/forge/skills/test-guide/references/test-type-model.md`——这是 plugin 内部的 skill 参考文件，属于分发包内容，而非 docs/ 下的项目文档。提案将一个不存在的目录纳入 L2 审计范围，这意味着要么审计执行时会发现目标不存在，要么提案的审计范围本身就包含错误。无论哪种情况，都指向提案撰写时未对文件系统做实际验证。

问题：提案 Scope 部分声称 "根目录下除 README.md 和 DESIGN.md 外无其他面向用户的文档文件（CONTRIBUTING.md、CHANGELOG.md 等不存在）"。实际根目录存在 CLAUDE.md 和 CLAUDE.template.md。虽然 CLAUDE.md 是 AI 代理指令文件而非传统用户文档，但它确实是新成员（包括 AI 代理作为虚拟成员）了解项目规范的第一个入口。提案在范围完整性说明中声称的"不会因目录层级混合而遗漏审计目标"，与实际上忽略了 CLAUDE.md 这一关键文件的事实相矛盾。

### 审计范围的关键遗漏

风险：提案将 docs/features/（182 个子目录）和 docs/proposals/（149 个子目录）完全排除在审计范围之外。但根据 `docs/conventions/forge-distribution.md` 的描述，`docs/conventions/` 中的规范是由 `/consolidate-specs` 从 feature 文档中提取的派生产物。这意味着 L2 审计 `docs/conventions/` 与代码的不一致时，根因可能在于 features 目录中的源文档过时。但 features 被排除在范围外，导致 L2 审计只能发现表面症状而无法追溯到根因。更严重的是，features 目录中的 manifest.md、prd/、design/、tasks/ 是 AI 代理执行任务时直接读取的输入——如果这些文档引用了已废弃的 CLI 命令或路径，其危害远大于 README 中的描述不一致。

风险：提案 Comparison Table 中声称 "执行现有5个提案的146个task"，但实际验证发现，这 5 个审计提案目录（plugin-consistency-audit、skill-ecosystem-audit、skill-instruction-audit、prompt-template-audit、test-pipeline-consistency-audit）中均不存在 tasks/ 子目录。每个提案目录仅包含 proposal.md 和（部分）eval/ 目录。146 这个数字在当前代码库中找不到对应的实体。提案以这个无法验证的数字作为排除"执行现有提案"这一替代方案的理由之一，其论证基础不可靠。

### 成功标准定义问题

问题：提案 Success Criteria 要求 "每层完成后随机抽取 10% 的审计结果进行人工复核，遗漏率不超过 20%"。这个指标在两个维度上有问题。第一，"遗漏率"的测量需要已知全集——只有知道所有不一致问题的完整清单，才能计算审计遗漏了百分之多少。但审计的目的恰恰是发现这个完整清单，所以遗漏率本身是一个无法直接测量的量。第二，20% 的阈值意味着允许每 5 个问题中遗漏 1 个，对于声称"为 v3.0.0 发布提供准确性保障"的审计来说，这个容忍度偏高。在缺乏独立基准（如完全独立的双重审计）的情况下，10% 抽样复核只能发现最明显的遗漏，对系统性盲点（如整类问题被忽视）无能为力。

问题：提案 Success Criteria 要求 "所有问题已转化为可执行 Task，修复类 Task 可由 task-executor 独立执行" 与 Constraints 部分的 "知识库清理需人工确认，不可自动删除" 存在张力。如果知识库审查 Task 需要"人工确认"环节，那它就不是完全"可由 task-executor 独立执行"的。提案在 Non-Functional Requirements 中补充了"需人工确认的 Task 在描述中明确标注确认步骤，task-executor 执行到确认点时暂停等待人工输入"，这在一定程度上缓解了矛盾，但"独立执行"的措辞仍会误导 task 生成时的粒度设计。

### 审计方法论与时间约束的矛盾

风险：提案承诺 "从审计启动到三份层级报告全部产出，不超过 2 个工作日（约 16 小时有效工作时间）"。同时，提案的审计执行流程要求对每份文档执行"声明提取 → 代码定位 → 逐条比对 → 缺失检测 → 结果记录"五个步骤，其中"逐条比对"包含"阅读代码逻辑，对比文档描述的步骤/顺序/参数是否与代码实际行为一致"。L2 层有约 27 个文件（business-rules 4 + conventions 18 + 假设 reference 有文件），L3 有 143 条知识库条目。在 16 小时内对 27 个规范文件逐一与代码交叉比对，同时对 143 条知识库条目逐条验证适用性，平均每个规范文件的审计时间不到 30 分钟，每条知识库条目的审查时间不到 4 分钟。这个时间预算对于"语义比对"（而非简单的路径存在性检查）来说极其紧张。如果审计代理为了赶时间而跳过深度代码阅读，审计质量将无法保证。

风险：提案的"层级间反馈机制"要求 L1/L2 发现的代码结构不一致须同步检查 L3 相关条目，L3 发现的过时引用也须通知 L2。这意味着三层审计之间存在交叉依赖。但提案同时要求 2 个工作日完成全部三层，且 L1/L2/L3 可"部分并行"。在实际执行中，如果 L1 发现了重大的架构描述不一致（如 hook 执行顺序），L3 审计人员需要先理解这个发现，再去检查引用该 hook 的 lessons 是否过时。这种跨层协调在紧凑时间线下的可操作性存疑。

### 分发约束未纳入审计框架

风险：提案在 L2 审计范围中包含了 `docs/conventions/` 和 `docs/business-rules/`，但未考虑这些文档在 Forge 分发模型中的角色。根据 `docs/conventions/forge-distribution.md`，`docs/` 目录的内容不分发到用户环境——分发的仅是 `plugins/forge/` 下的内容。AI 代理在用户项目中运行时，读取的是 plugin 分发包中的 skills/、hooks/、agents/，而非源码仓库中的 docs/ 目录。这意味着 docs/ 下的文档不一致主要影响的是源码仓库的维护者（包括在源码仓库中工作的 AI 代理），而非终端用户的 AI 代理。提案的 Problem 描述中 "AI 代理按旧术语生成测试代码" 的危害场景，需要区分是"在源码仓库中工作的代理"还是"在用户项目中工作的代理"——两者的文档消费路径完全不同。

风险：提案声称 "审计产出的所有 Skill、Command、任务模板、提示词模板等统一采用英文撰写"。但提案本身是中文撰写的，且 docs/ 下的所有文档也是中文。如果审计产出的 Task 是英文的，但需要审计的源文档是中文的，task-executor 在执行时需要在两种语言间切换，增加了误读风险。这一约束的动机未在提案中说明。

---

## Section 3: Improvement Suggestions

建议：修正提案中所有可验证的事实声明。具体而言：(1) 将 "docs/conventions/ 下 22 份规范文档（顶层15份 + testing/7份）" 修正为 "docs/conventions/ 下 18 份规范文件（顶层15份 + testing/ 子目录3份）"；(2) 删除对 `docs/reference/test-type-model.md` 的引用，承认该目录不存在，并重新评估 L2 审计范围的完整性；(3) 将 L1 文件数从 "12 文件" 确认无误（3 根文档 + 4 user-guide + 5 official-references = 12），同时明确说明 CLAUDE.md 是否应纳入 L1 范围；(4) 修正 Comparison Table 中 "146个task" 的数字——当前 5 个审计提案目录中均无 tasks/ 子目录，这个数字缺乏来源。

建议：重新审视 docs/features/ 的排除决定。不需要审计全部 182 个 feature 目录，但建议至少将以下两类纳入审计范围：(1) 包含 manifest.md 且状态标记为活跃（如 in-progress）的 feature——这些是当前工作流正在使用的文档；(2) 被 docs/conventions/ 中规范引用的 feature 源文档——L2 审计发现不一致时需要追溯到这些源文件。这样可以控制范围扩大（预计增加 5-10 个文件），同时覆盖根因分析路径。

建议：重新定义成功标准。将不可验证的结果性标准改为可验证的过程性标准：(1) 将 "遗漏率不超过 20%" 改为 "对范围内每个目标文件逐一完成声明提取、代码定位、逐条比对三个步骤，过程记录存档"；(2) 增加一个更可靠的交叉验证机制——对 L1/L2 审计结果，由独立审计代理对同一文件执行相同流程，比较两份结果的重合度，以重合度代替"遗漏率"作为质量指标；(3) 将 "每个 Task 可由 task-executor 独立执行" 改为区分两类 Task 描述模板——修复类 Task 使用标准模板，审查类 Task 使用带人工确认节点的模板。

建议：补充审计方法论的详细说明。当前提案的"声明提取 → 代码定位 → 逐条比对"流程过于抽象。建议至少为三类比对提供具体的操作协议：(1) 路径引用验证——`find` + `grep` 确认目标存在且内容匹配；(2) 行为描述验证——定位代码中的函数/配置/hook，阅读实际逻辑，逐项对比文档声称的步骤、顺序、参数、错误处理；(3) 配置声明验证——对比文档中的版本号、依赖项、默认值与代码中实际定义的值。这些操作协议应在提案中明确，作为 task 生成的输入。

建议：将 L3 知识库审查与 L1/L2 文档审计在方法论上明确区分。L1/L2 是客观比对（文档声称 X，代码实际 Y），L3 是主观判断（这个经验教训在当前项目中是否仍有参考价值）。建议将 L3 的成功标准从"遗漏率"改为"一致性率"——由两名独立审查者对同一批条目做判定，以两人判定一致的比例作为质量指标。这也意味着 L3 的 Task 粒度和执行协议应与 L1/L2 不同，不应强行统一到同一套模板下。

建议：在提案中明确与 `consolidate-specs` skill 的关系。`consolidate-specs` 已具备 drift detection 功能（`plugins/forge/skills/consolidate-specs/rules/drift-detection.md`），可以自动检测规范与代码的漂移。提案应说明：(1) 本次审计与 consolidate-specs 的 drift detection 在覆盖范围上的差异（consolidate-specs 只覆盖从 feature 提取的规范，不覆盖用户文档和知识库）；(2) 审计完成后如何利用 consolidate-specs 建立持续监控，防止审计成果被后续开发再次漂移。这将一次性审计的价值从"当前快照"延伸到"持续保障"。

建议：考虑将提案中 "docs/reference/ 仅含 1 个文件 test-type-model.md" 这一错误发现为契机，对 L2 审计范围做一次完整的文件系统扫描验证。当前提案列出的 L2 范围（business-rules/4 + conventions/18 + reference/0）可能还有其他遗漏。例如 `docs/` 下还存在 `experts/`、`forensics/`、`harness-reports/`、`plan/`、`research/`、`self-evolution/`、`superpowers/`、`todos/` 等目录，提案未说明这些目录为何被排除。即使排除是合理的，也需要在 Scope 部分明确说明排除理由，否则审计范围的定义看起来像是基于不完整的目录扫描。
