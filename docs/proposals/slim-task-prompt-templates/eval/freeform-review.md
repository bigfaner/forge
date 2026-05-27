---
reviewer: "Prompt Template Engineering & Agent Protocol Design Specialist"
proposal: "docs/proposals/slim-task-prompt-templates/proposal.md"
date: "2026-05-27"
review_type: "freeform"
---

# Freeform Narrative Review: 精简任务 Prompt 模板

## Section 1: Background Assessment

本提案的核心叙事是清晰的：Forge 的 15 个任务 prompt 模板中存在大量非指令内容——HTML 注释、解释性描述、冗长的角色定义、AC 验证块的冗余文字、CODING_PRINCIPLES 中的举例和说明——这些内容不直接指导 agent 行为，却消耗 token 并稀释指令清晰度。提案给出的量化估计是约 200 行的冗余，覆盖了 coding.* 系列、gate、doc、test-* 系列以及 task-executor agent 的 Execution Protocol。

技术方案是"就地精简"：保持每个模板的独立文件结构不变，不抽取公共模块，不在 `prompt.go` 中做任何改动，仅逐个编辑 .md 文件删除非指令内容。核心原则被提炼为 "prompt 是指令，不是文档"——这个口号简洁有力，体现了正确的 prompt engineering 直觉。

提案在 Assumptions Challenged 部分做出了一项关键论断：角色描述中的自然语言 ("You are a focused...") 对 LLM 行为的影响可以通过祈使句替代，因为"Agent 的执行行为由后续的 Workflow 步骤定义，不是由角色描述定义的"。这个论断构成了整个精简操作的理论基础——如果它是对的，那么所有角色描述都可以安全地压缩；如果它是错的，那么大规模精简就可能在无意中改变 agent 的行为模式。

Success Criteria 的核心指标是行数减少（>=150 行）和行为无可见差异，但没有定义"行为无可见差异"的验证方法。

## Section 2: Key Risk Identification

问题：提案的核心论断"角色描述不定义 agent 行为"与 prompt engineering 的已知经验相悖，且提案未提供任何实验证据支撑。

> "Agent 的执行行为由后续的 Workflow 步骤定义，不是由角色描述定义的。"

这个论断在理论上是吸引人的，但实际 LLM 系统中，角色描述（system prompt 中的 persona 定义）对输出的影响是经过广泛验证的——角色描述不仅影响语气和风格，还影响 LLM 的推理深度、错误检测倾向性、以及遵循指令的严格程度。Workflow 步骤定义了"做什么"，角色描述定义了"以什么身份认知和判断"。移除角色描述中的"自然语言"部分，即使替换为祈使句，也会改变 LLM 在推理时的 latent activation pattern。提案没有引用任何实验或文献支持其论断，也没有计划在修改前后进行 A/B 测试来验证行为等价性。对于一个以"prompt template engineering"为核心领域的人看来，这是提案中最脆弱的一环。

风险：AC 验证块从 ~12 行精简到 ~3 行的压缩率（75% 缩减）缺乏风险分析，且可能破坏 AC 执行的关键路径。

> "AC 验证块冗余：9 (coding.*, gate, doc)，每处 ~12 行可缩至 ~3 行"

AC（Acceptance Criteria）验证是 coding.* 模板中最关键的合规门控环节。当前的 ~12 行结构很可能包含：AC 条目列表、验证步骤、失败处理策略、验证顺序约束、以及可能的引用回溯指令。将其压缩到 ~3 行意味着大量结构信息被删除或隐式化。问题是：当前 12 行中有多少行是真正的冗余，有多少行是确保 AC 验证被执行、被正确排序、以及在失败时有明确降级路径的必要指令？提案仅计算了行数但未分析每行的功能。如果某个 AC 块中实际上只有 3 行是核心指令，其他 9 行是解释性的，那么压缩是安全的——但提案没有展示这个分析过程。在缺乏逐块功能拆解的情况下，"~12 行可缩至 ~3 行"是一个断言而非结论。

风险：Execution Protocol 的步骤合并假设步骤分离是冗余而非功能设计，但未做协议级依赖分析。

> "Execution Protocol 步骤合并（步骤 4/5/6 处理 prompt 获取逻辑可合并为 1 步）"

在 agent 执行协议设计中，步骤是否应该分离不取决于它们是否"处理同一类事务"，而取决于它们是否有独立的错误处理路径、状态检查点和中止条件。步骤 4/5/6 被分开很有可能是因为每个步骤对应一个不同的系统交互点：步骤 4 可能负责检查 prompt 文件是否存在，步骤 5 负责加载并解析模板中的占位符，步骤 6 负责验证解析后的 prompt 是否完整。将它们合并为 1 步意味着失败时无法区分是哪个环节出了问题，日志粒度和可调试性都会下降。提案将步骤合并完全视为"行数优化"问题，未评估对可调试性和错误恢复的影响。

风险：Retry Strategy 与 Complex Error Pause Flow 的去重合并可能掩盖两种机制之间的关键区别。

> "Retry Strategy 与 Complex Error Pause Flow 去重合并"

Retry Strategy 和 Complex Error Pause Flow 听起来相似，但 serving 的目的可能是正交的：Retry Strategy 定义的是"什么条件下重试、重试多少次、退避策略是什么"——这是关于重试时机和频率的安排；Complex Error Pause Flow 定义的是"遇到不可恢复错误时如何暂停、如何保存状态、如何通知用户"——这是关于错误升级路径的安排。合并它们的前提是两者的语义和操作对象完全重叠，但提案没有给出这个证明。

问题：Success Criteria 中的行数减少指标（≥150 行）与被验证的 agent 行为等价性之间存在目标偏差。

> "15 个模板文件 + task-executor 共减少 **≥150 行**（去除注释、解释性描述、冗长定义）"

行数减少是一个可测量但不一定可取的指标。它衡量的是"删除动作的数量"而非"质量改进的程度"。如果某个修改删除了 20 行但改变了 agent 的行为，它在行数指标上成功了，但在实际目标上失败了。更根本的问题是：提案的 Success Criteria 中，"所有模板精简后，agent 执行相同 task 的行为无可见差异"被列为一个独立条件，却没有定义"可见差异"的检测方法——是人工观察 task 输出？是自动化 diff？是跑一遍全部 journey？没有定义检测方法的标准就等同于不可验证。这意味着 Success Criteria 实际上只有一个可操作的指标：行数。这是典型的目标替代（surrogate goal）问题。

问题：证据表中的总冗余量（~200 行）存在统计口径问题，可能给人过度的紧迫感。

> "总计：约 **200 行** 非指令冗余。每执行一个 task，agent 都要阅读这些无用 token，形成累积开销。"

"每执行一个 task"意味着每个 task 消耗所有 15 个模板 + task-executor 的全部 200 行。但实际执行中，一个任务只调用 1 个模板（例如 coding-feature）+ task-executor。每个 template 的实际冗余远低于 200 行。提案在 Evidence 表中已经按类别列出了冗余分布（CODING_PRINCIPLES ~50 行、AC 验证块 ~80 行等），但"总计 200 行"的表述方式容易让读者高估单个 task 的浪费。更精确的说法应该是"每个 coding.* task 包含约 50 行 CODING_PRINCIPLES 冗余 + AC 验证块 + 角色描述 = 约 80-100 行冗余"——这个数字依然值得优化，但比"200 行"更诚实。这个问题不改变提案的合理性，但它反映了提案在论证严谨性上的松懈。

风险：提案缺乏回归检测机制——编辑后如何确认没有引入行为退化。

> 整个提案中唯一的验证描述是 "每个模板修改后对比：所有功能点是否仍被覆盖；task-executor 的每个步骤的行为约束是否保持"

这个 Mitigation 策略描述的是一个手工对照过程（对比功能点、检查步骤约束），但没有定义具体的操作方式：谁来对比？用什么基准？功能点列表从哪里来？行为约束的完整规范在哪里？对于一个涉及 16 个文件、每个文件都有实质性内容删除的修改，依赖手工"对比"作为唯一的验证手段风险极高。在 prompt template engineering 中，每一次修改都应当有对应的 regression test——要么是一组已知的 task 输入-输出对，要么是行为断言式的 E2E 测试。提案没有提及任何形式的行为回归测试。

## Section 3: Improvement Suggestions

建议：在修改模板前，先为每个模板建立一份"行为等价性规范"，作为修改后的验证基准。

Addresses: 风险（AC 验证块压缩可能破坏 AC 执行）和风险（缺乏回归检测机制）。

对于每个要修改的模板（尤其是 coding.* 和 task-executor.md），在修改前先提取"该模板当前通过 prompt 实现的所有功能约束点"。例如，对于 coding-feature.md 的 AC 部分，逐行标注：这一行是"必须执行的指令"还是"对指令的解释"还是"示例"？建立一个功能点列表作为 golden reference。修改后，用这个列表逐项验证：每个功能约束在修改后的模板中是否仍被覆盖。这个列表本身就是后续任何修改的 regression test baseline。采纳后，Mitigation 表中的"每个模板修改后对比"从模糊的手工操作升级为有明确 checklists 的结构化验证，回归风险大幅降低。

建议：将 "≥150 行" 的成功指标替换或补充为 "每个模板功能约束完整保留率 100%"。

Addresses: 问题（Success Criteria 中的行数指标与行为等价性之间存在目标偏差）。

行数指标只测量删除数量，不测量保留质量。建议的替代方案：定义每个模板的"功能约束覆盖率"——即修改后保留的功能约束点数除以修改前的总约束点数，目标为 100%。行数减少可以保留为次级参考指标，但不作为决策 gate。这样设计者不会被激励去"为减少行数而减少行数"，而是被激励为"在保持功能覆盖的前提下去除冗余"。采纳后，Success Criteria 的第一项从 "减少 ≥150 行" 变为 "所有功能约束点 100% 保留 + 行数减少 ≥150 行"，形成约束-优化双层结构。

建议：在精简 CODING_PRINCIPLES 之前，识别其举例和解释是否实际上在扮演 few-shot demonstration 的角色。

Addresses: 风险（角色描述移除对行为的隐性影响）的延伸。

CODING_PRINCIPLES 中的"每原则 2-5 行"冗余被提案归为"解释性冗余"。但在 prompt engineering 实践中，principles + examples 是非常标准的 few-shot 结构——例子不是对原则的解释，而是对原则应用方式的示范。如果删除这些例子，可能出现"原则仍在但 agent 不再以预期方式应用它"的情况。建议的实际做法：将每个 CODING_PRINCIPLES 条目标注为 [rule-only] 或 [rule+example]，仅删除确认为纯解释性的内容。保留那些看似是"举例"但实际上在示范边界条件的条目。采纳后，CODING_PRINCIPLES 的精简不再是"硬砍 50 行"的批量操作，而是逐条分析后的精确修剪。

建议：在合并 Execution Protocol 步骤前，为每一步绘制错误恢复依赖图。

Addresses: 风险（Execution Protocol 步骤合并忽略功能分离的设计意图）。

对 task-executor.md 的步骤 4/5/6，逐步骤回答：如果这一步失败，当前步骤的中止条件是调用者（上一步/下一步）要知道的，还是可以在本步骤内部封装的？如果每一步都有独立的错误恢复路径，则不应合并；如果错误的处理方式相同（全部是 "retry from current step"），则合并是安全的。在修改之前绘制一个微型状态机，标注每个步骤的 entry condition、success exit、failure exit、和 side effect，然后基于这个状态机做合并决策。采纳后，协议修改不再依赖直觉判断，而是有明确的分析产出物作为决策依据。

建议：增加验证阶段的回滚标准和自动化检测手段。

Addresses: 风险（缺乏回归检测机制）。

提案的 Next Steps 是 "Proceed to /write-prd to formalize requirements"，没有提及验证阶段。建议在实现阶段明确指出：修改 16 个模板后，应当用当前的 forge 任务系统运行一组代表性的 journey 或 task（覆盖 coding-feature、coding-fix、gate、doc 各至少一个），通过对比修改前后的 agent 输出（或至少是 agent 的思考链和行为路径）来检测行为漂移。如果任何 journey 的输出出现与预期不符的变化，应终止修改并回滚到原始模板。采纳后，整个项目有了明确的 exit criteria 和 fail-fast 机制，不再依赖"但愿没有问题"的乐观假设。