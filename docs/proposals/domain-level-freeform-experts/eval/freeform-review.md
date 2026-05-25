---
reviewer: "Expert System Design & Prompt Architecture Strategist"
proposal: "docs/proposals/domain-level-freeform-experts/proposal.md"
date: "2026-05-25"
review_type: "freeform"
---

# Freeform Narrative Review: 领域级自由专家生成机制

## Section 1: Background Assessment

本提案试图解决 Forge eval pipeline 中自由专家（freeform expert）复用率为零的核心问题。当前的 expert-inference 流程为每个 proposal 独立生成一个高度特化的专家，domain 关键词从 proposal 文本中直接提取，导致不同 proposal 生成的专家之间几乎没有关键词重叠。复用匹配依赖 Jaccard 相似度达到 0.3 的阈值，而实际上从未达到过——11 个已有专家的 Jaccard 相似度始终低于此值。

提案的核心方案是引入一个预定义的大领域分类表（8-12 个条目），将 expert-inference 的单步生成改为两步：先匹配大领域，再在大领域范围内生成专家。这样，同一个大领域内的不同 proposal 将共享领域级专家，避免重复生成。分类表未覆盖的领域有 LLM 自由推断降级路径。改动范围被明确限定在三个 prompt 文件：expert-inference.md、expert-template.md、freeform-expert-persistence.md。

我注意到提案选择了一条有意的设计路径：在"固定专家库"和"完全动态推断"之间取折中。这个定位是合理的——但折中方案的真正风险不在于两端，而在于中间地带的维护成本是否可控，以及分类表本身是否比它试图解决的问题更难管理。

## Section 2: Key Risk Identification

风险：分类表的"8-12 个大领域"假设缺乏实证支撑，且分类表会引入单点维护负担。

提案在 Key Risks 表中声称分类表应"控制在大领域粒度（8-12 个），不细化到子领域"，但未提供任何关于这个数字如何得出的依据。当前 11 个 proposal 横跨了多少个自然领域？提案没有做这个基础分析。如果实际自然领域数量是 15-20 个，强行压缩到 8-12 个会导致领域边界模糊；如果是 6-8 个，分类表的价值就值得质疑。更关键的是，分类表一旦嵌入 expert-inference.md（纯 prompt 文本），每次调整都需要修改 prompt 文件本身——这与提案声称的"可扩展性：新增领域只需修改 prompt 文件"看似简单，实际上意味着每次新增领域都要回归测试整个 inference prompt 的行为稳定性。分类表不是配置数据，它是 prompt 推理路径的一部分，改动它的副作用不可预测。

问题：提案假设"大领域匹配"的可靠性高于"关键词匹配"，但未分析为什么换一个匹配粒度就能解决一致性问题。

提案在 Innovation Highlights 中声称"通过预定义分类表保证领域标签一致性。不同 proposal 对同一领域的识别结果相同（不会出现 'test infrastructure' vs 'testing pipeline' 的不一致）"。这个论断跳过了一个关键环节：从 proposal 文本到分类表条目的映射本身也需要 LLM 推断。LLM 面对 proposal 提取特征后查分类表，这个过程与直接提取 domain 关键词一样，都是 LLM 对文本的语义判断。如果 LLM 对同一领域的不同 proposal 能可靠地映射到同一个分类表条目，那么理论上它也应该能可靠地生成一致的 domain 关键词——提案没有解释为什么前者比后者更可靠。分类表的价值应该在于缩小搜索空间（从开放词汇到固定列表），而非保证一致性。提案对"一致性"的承诺过度了。

风险：领域级专家的评审深度退化问题未得到充分分析，提案提供的缓解策略是循环论证。

提案在 Key Risks 表中识别了"领域级专家的评审深度不如 proposal-specific 专家"的风险，缓解措施是"LLM 在领域内细化专业方向时参考 proposal 内容，保持针对性"。但这恰好是当前系统已经在做的事情——如果 LLM 足以在领域内针对具体 proposal 细化，那么当前的 proposal-specific 专家本就不应存在问题。提案的真正风险是：一个"构建与测试基础设施"领域专家面对 surface-aware-justfile proposal 时，其评审焦点会被整个领域的广度稀释。提案的 Success Criteria 要求"新生成的专家 domain 关键词覆盖范围 >= 2 个 proposal 的领域交集"——但覆盖范围扩大和评审深度之间是零和关系，提案没有给出如何同时实现两者的机制。

问题：Scope 声明的"复用匹配对旧专家的兼容"被标记为 Out of Scope，但旧专家会持续干扰新的匹配流程。

提案在 Out of Scope 中明确排除了"复用匹配对旧专家的兼容"和"迁移或废弃现有 11 个专家文件"。然而，freeform-expert-persistence.md 的 Step 1 会加载 docs/experts/ 下的所有非 deprecated 专家。这意味着 11 个旧的 proposal-specific 专家会持续出现在候选列表中。它们不会被新系统生成，但它们的 domain 关键词会参与 Jaccard 计算，可能产生误匹配（例如 "Build Orchestration & Test Infrastructure" 专家可能匹配到任何涉及构建或测试的 proposal）。旧的 proposal-specific 专家本质上成为了新系统中的噪声数据，而提案没有计划处理这个噪声。

风险：场景 3（跨领域 proposal）的"匹配最相关的一个大领域"策略可能导致系统性评审盲区。

提案在场景 3 中描述"proposal 涉及多个领域（如'Agent架构' + '配置Schema'）-> 匹配最相关的一个大领域"。这意味着跨领域 proposal 会被强制归入一个领域，由单一领域专家评审。例如一个同时涉及"Agent架构"和"配置Schema"的 proposal，如果匹配到"Agent架构"，那么配置 Schema 方面的评审深度将完全依赖专家的泛化能力。这个决策的后果是：越是跨领域的复杂 proposal（往往是最需要深度评审的），评审质量越可能不足。提案假设用户可以通过 Modify 循环调整焦点来弥补，但 Modify 循环的最大轮次是 3，且 Modify 的是专家的焦点而非领域归属。

问题：提案声称改动"不涉及代码变更"过于简化，忽略了 prompt 链联动的实际复杂度。

Feasibility Assessment 中称"改动仅涉及 prompt 文件（Markdown），不涉及代码变更"，Resource & Timeline 评估为"单次 prompt 改写 + 测试验证，工作量小"。然而，三个 prompt 文件的联动关系并非简单：expert-inference.md 的新两步流程需要与 freeform-expert-persistence.md 的 Jaccard 匹配逻辑协调（新的 scope 字段如何影响匹配？domain-level 专家的关键词更广，Jaccard 分母增大，匹配分数的分布会发生变化），expert-template.md 的 scope 字段需要被 persistence 规则识别和处理。此外，extraction-prompt.md（提案未提及但存在于 freeform 目录中）也可能需要调整，因为它从 freeform review 中提取 findings 时需要知道专家的 scope 级别。提案低估了 prompt 链的联动复杂度。

## Section 3: Improvement Suggestions

建议：在提案中增加一个"现有领域分布分析"章节，用 11 个已有 proposal 的实际数据验证 8-12 个大领域假设。

这个建议针对的是分类表规模假设缺乏实证依据的问题。具体做法是：将 11 个已有专家的 domain 字段提取出来，用人工或 LLM 辅助聚类，看自然形成的领域簇有多少个、每个簇包含多少个 proposal。如果自然簇是 6 个，说明 8-12 的范围合理；如果是 15 个，需要重新考虑分类表的粒度或是否应该用层次分类（大领域 + 子领域）而非扁平分类。这个分析不仅验证假设，还直接产出分类表的初始版本。采纳后，提案的 Feasibility Assessment 将从"假设合理"升级为"数据支撑"，Key Risks 中关于分类表膨胀的讨论也有了基线数据。

建议：将分类表从 prompt 内嵌改为外部数据文件（如 YAML/JSON），通过 prompt 引用机制加载，而非硬编码在 expert-inference.md 的正文中。

这针对的是分类表维护成本和维护方式的风险。当前提案将分类表设计为 prompt 文本的一部分，每次调整都需要修改 expert-inference.md 并回归测试整个 inference prompt。如果将分类表提取为独立文件（例如 `experts/freeform/domain-classification.yaml`），expert-inference.md 通过文件引用机制加载它，那么调整分类表条目不会触及 inference prompt 的推理逻辑，回归测试范围大幅缩小。采纳后，提案的 Non-Functional Requirements 中"可扩展性：新增领域只需修改 prompt 文件"变为"新增领域只需修改分类表文件"，风险隔离更彻底。

建议：明确旧专家的处理策略——要么标记为 proposal-specific scope 并在匹配时降权，要么在过渡期后批量 deprecated。

这针对的是旧专家成为匹配噪声的问题。具体方案：在 expert-template.md 的 scope 字段中区分 `domain-level` 和 `proposal-specific`，在 freeform-expert-persistence.md 的匹配逻辑中，对 proposal-specific 专家的匹配结果施加惩罚因子（例如 Jaccard 分数乘以 0.7），使其在候选排序中自然靠后但不被完全排除。这比直接标记 deprecated 更温和，也保留了 proposal-specific 专家在无 domain-level 专家可用时的降级价值。采纳后，Out of Scope 中可以去掉"复用匹配对旧专家的兼容"（因为它不再是"兼容"问题而是"共存"策略），同时解决了旧专家噪声问题。

建议：为跨领域 proposal 设计显式的多专家策略，而非依赖单专家 + Modify 循环。

这针对的是跨领域 proposal 评审深度不足的风险。当分类表匹配显示 proposal 涉及 2 个以上大领域时，可以考虑两种策略：（1）为每个匹配的大领域各生成/复用一个专家，以主专家 + 辅助专家的形式进行评审；（2）在专家生成阶段将跨领域 proposal 标记为 `scope: cross-domain`，由 LLM 生成一个跨领域专家，但在 prompt 中显式要求其覆盖所有匹配的大领域并标注各领域的评审深度权重。策略 1 的代价是多一次评审循环，策略 2 的代价是单专家的焦点分散。无论选哪种，都比当前的"匹配最相关的一个"更诚实。采纳后，场景 3 的描述从"匹配最相关的一个大领域"变为"触发跨领域专家策略"，评审盲区风险得到显式管理。
