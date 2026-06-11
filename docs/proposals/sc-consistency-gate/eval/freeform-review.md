---
reviewer: requirements-consistency-forge-eval-specialist
date: "2026-05-25"
type: freeform
iteration: 0
---

# Freeform Review: SC Consistency Gate Proposal

## 1. Background Assessment

This proposal introduces a dual-layer defense mechanism to catch logical contradictions within a proposal's Success Criteria (SC) — a problem demonstrated by `pipeline-integration-stitch`, which passed adversarial eval at 897/1000 despite containing mutually exclusive SC items. The two layers are: (1) a brainstorm-skill rule that checks SC consistency at authoring time, and (2) an eval-scorer protocol extension that checks consistency during adversarial evaluation. The core innovation is scope-area clustering — grouping SC items by the file/directory/module they affect, then checking satisfiability only within clusters — which the proposal claims reduces checking from O(n^2) to O(n).

The proposal is well-structured, follows the standard proposal template, and addresses a real and documented failure mode. Its scope is appropriately narrow: one new rule file, one SKILL.md reference addition, one scorer-protocol extension, and one rubric modification. However, several critical assumptions warrant deeper scrutiny, particularly around the clustering heuristic's soundness guarantees, the dual-layer fallback completeness, and the feasibility of the claimed execution budget.

## 2. Key Risks

**风险：聚类启发式的完备性保证存在逻辑漏洞。** 提案声称"SC 矛盾的必要条件是两个条目作用于同一代码区域（同文件、同目录、同模块）"——这是一个未经证明的断言，被当作了公理使用。考虑 lesson 中引用的原始矛盾场景：`grep -r "gen-and-run" forge-cli/ 返回零结果` 和 `validate_index.go 对 test.gen-and-run 返回迁移错误`——这两个条目确实作用于同一区域。但该断言作为"必要条件"并不成立。两个 SC 可以作用于完全不同的文件却仍然逻辑互斥。例如："所有 API 响应时间 < 100ms"（作用于 `api/handlers/`）和"所有请求必须写入同步磁盘日志"（作用于 `lib/logger/`）——两者作用区域不同，但前者要求快速响应，后者引入同步 I/O 延迟，逻辑上可能矛盾。提案原文写的是"SC 矛盾的必要条件是两个条目作用于同一代码区域"，如果"必要条件"不成立，则聚类可能遗漏真正的矛盾对，产生假阴性。这不是理论上的极端 case——性能类 SC 和完整性/安全类 SC 分别作用于不同模块是常见的 proposal 结构。

**问题：O(n^2) 到 O(n) 的复杂度声明不准确。** 提案声称"先聚类再检查将 O(n^2) 削减为 O(n)（每个条目只与同组条目比较）"。但实际复杂度取决于平均簇大小 k，即 O(n·k)。如果所有 SC 条目都指向同一核心模块（如 `forge-cli/`），则 k = n，退化为 O(n^2)。提案确实提到"典型提案检查量从 125-246 对降至 40-60 对，削减 70-80%"，这暗示 k 值约为 5-7，但未给出 k 的上界约束或最坏情况分析。当提案专注于单一模块的深度重构时（如本提案本身就是对 eval pipeline 的集中改造），所有 7 条 SC 可能全部指向 `plugins/forge/skills/eval/`，此时聚类无法提供任何削减。

**风险：Layer 2 兜底机制覆盖不完整——eval scorer Phase 1 的 self-contradiction check 本身是 LLM 执行的非确定性过程。** 提案承认 Layer 1 的风险是"agent 忽略规则文件不执行检查"，并将此风险标记为 M likelihood / H impact，mitigation 是"SKILL.md 引用 + 规则文件命名清晰 + eval 层兜底"。但 Layer 2 同样依赖 LLM agent 执行聚类和可满足性推理——如果 LLM 能忽略 brainstorm 层的显式规则文件，它同样可能在 scorer protocol 中跳过或敷衍执行聚类检查。两层都依赖同一执行者（LLM）的同一能力（逻辑推理），这不构成真正的冗余，而是单点失效的伪装。真正的兜底需要一个不依赖 LLM 主动推理的机制（例如结构化的 SC 依赖声明模板或自动化冲突检测脚本）。

**问题：proposal rubric D9 的 25 分重新分配缺乏具体说明。** 提案在 In Scope 中写"修改 `plugins/forge/skills/eval/rubrics/proposal.md` Dimension 9 — 新增 'SC internal consistency' criterion (25pts)，调整现有 criterion 分值"，Success Criteria 中也写"proposal rubric D9 包含 'SC internal consistency' criterion (25pts)，D9 总分 80pts 不变"。但当前 D9 只有两个 criterion："Criteria are measurable and testable" (55pts) 和 "Coverage is complete" (25pts)，合计正好 80pts。如果新增 25pts 的 criterion 而总分不变，则现有 criterion 必须被压缩至少 25pts。然而提案没有说明 55pts 和 25pts 各压缩多少——是 measurability 从 55 降到 30，还是 coverage 从 25 降到 0？更严重的是，D10 "Logical Consistency" (90pts) 中已有 "Scope <-> Solution <-> Success Criteria aligned" (30pts) 这一 criterion，其描述正是检查 SC 之间的一致性。新增 D9 的 "SC internal consistency" 是否与 D10 的这个 criterion 功能重叠？提案未做任何区分说明。

**问题：20 秒执行预算的可行性缺乏依据。** 提案声称"聚类 + 组内检查，典型提案检查量 < 60 对（vs 朴素逐对 125-246 对），agent 执行时间 < 20 秒"。但这个预算是针对 LLM agent 的推理时间，而非程序执行时间。LLM agent 需要：(1) 解析所有 SC 文本，(2) 识别每个 SC 的作用区域（文件/目录/模块），(3) 执行聚类分组，(4) 对每组内的 SC 对进行逻辑可满足性推理。仅第 (4) 步，对 40-60 对 SC 进行可满足性推理，每对需要 agent 理解两个自然语言描述的语义并判断联合可满足性——这不是简单的字符串匹配。在当前 LLM 推理延迟下（通常每 token 50-100ms，一个 SC 对的可满足性推理可能需要 200-500 tokens），60 对的纯推理时间就可能达到 30-60 秒。提案未提供任何测试数据或推理链来支撑 20 秒的估算。

**风险：提案自身存在潜在的 SC 内部矛盾——"自吃狗粮"检查。** 作为 SC 一致性检测机制的提案，其自身 SC 应当首先通过一致性检查。SC-4 要求"proposal rubric D9 包含 'SC internal consistency' criterion (25pts)，D9 总分 80pts 不变"，而 SC-5 要求"对 lesson 中的 gen-and-run 矛盾场景，scorer 能在 D9 'SC internal consistency' 维度扣分并生成 attack point"。但 SC-6 同时要求"对无矛盾的 SC 集合，规则不产生任何阻塞或警告，D9 'SC internal consistency' 满分"。SC-5 和 SC-6 在概念上验证的是同一机制的两种状态（有矛盾 vs 无矛盾），但 SC-5 测试的是 scorer 在 eval 流程中的行为（生成 attack point），而 SC-1 到 SC-4 测试的是文件是否存在和内容是否正确——这混淆了"交付物存在性验证"和"运行时行为验证"两类不同性质的 criterion。后者（SC-5、SC-6）在 proposal 阶段无法通过文件检查验证，需要实际运行 scorer，这超出了 proposal review 的验证能力。

## 3. Improvement Suggestions

**建议：将聚类启发式从"必要条件"降级为"高概率启发式"，并明确定义其覆盖边界。** 提案原文应修改"SC 矛盾的必要条件是两个条目作用于同一代码区域"这一断言，承认聚类方法可能遗漏跨区域的逻辑矛盾。具体做法：(1) 在规则文件中加入跨区域检查步骤——在组内检查完成后，对所有 SC 进行一次轻量的全对扫描，但仅检查"方向型矛盾"（ADD vs SUBTRACT on same symbol），这比完整的可满足性检查开销小得多；(2) 在 eval scorer 扩展中同样保留全对方向检查作为 fallback。这样聚类负责深度语义检查（组内可满足性），全对负责浅层结构检查（方向冲突），两者互补。

**建议：明确 D9 新 criterion 的分值分配方案，并与 D10 做去重区分。** 应在提案中具体说明：(1) 现有 55pts "measurable and testable" 压缩为多少分，25pts "coverage is complete" 压缩为多少分；(2) D9 新 "SC internal consistency" 与 D10 "Scope <-> Solution <-> Success Criteria aligned" 的职责边界——前者检查 SC 条目之间的内部可满足性（A 与 B 能否同时为真），后者检查 SC 与 Scope/Solution 的覆盖对齐关系（SC 是否覆盖了所有 Scope 条目），两者关注点不同但需要明文界定以避免 scorer 重复扣分。

**建议：为 Layer 2 兜底增加不依赖 LLM 主动推理的结构化机制。** 考虑在 proposal 模板（尽管提案声明不改模板，但可以添加可选字段）或 brainstorm Step 5 的规则中，要求 agent 为每个 SC 条目标注其"影响的代码路径"列表。这个结构化标注可以被后续检查流程（包括 eval scorer）直接消费，而不需要 LLM 重新推断作用区域。这同时解决了两个问题：(1) 将隐式推理转化为显式标注，降低 LLM 遗漏的风险；(2) 标注本身可以作为 agent 理解 SC 的辅助手段——如果 agent 无法为某个 SC 标注影响路径，说明该 SC 可能不够具体。

**建议：将 20 秒执行预算替换为"与当前 brainstorm Step 5 写入时间相比的增加量 < 30%"，并补充实测基准。** 绝对时间预算在 LLM 推理延迟波动较大的环境下缺乏可验证性。使用相对增长比例更实际，并且提案应该承认这是需要实测验证的假设而非已确认的事实。在规则文件的 SC 中也应将"agent 执行时间 < 20 秒"改为可实际验证的 criterion。

**建议：SC-5 和 SC-6 应拆分为实现阶段的集成测试 criterion，而非 proposal review 的交付物 criterion。** 当前 SC 混合了两类验证：文件存在性（SC-1 至 SC-4）和运行时行为（SC-5、SC-6）。建议将 SC-5 和 SC-6 移至 "Next Steps" 或标注为 "Post-implementation verification"，在 SC 中只保留文件级别和内容级别的可验证项。例如，SC-5 可以改为"scorer-protocol Phase 1 Step 4 的文本中包含对 lesson `gotcha-proposal-success-criteria-contradiction.md` 中 gen-and-run 场景的显式引用作为测试用例示例"——这是文件级别的验证，不需要实际运行 scorer。
