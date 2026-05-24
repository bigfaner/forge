# Freeform Narrative Review: Freeform Pre-Revision Proposal

**Expert Profile**: Eval Pipeline Information-Flow Architect
**Reviewed Document**: `docs/proposals/eval-freeform-pre-revision/proposal.md`
**Date**: 2026-05-24

---

## Section 1: Background Assessment

This proposal addresses a concrete information-fidelity problem in the Forge eval pipeline's proposal evaluation flow. The current architecture routes freeform expert findings (Phase 0) through the Scorer before they reach the Reviser. The Scorer acts as a mapping layer with three possible outcomes per finding: map to a rubric dimension, classify as `[beyond-rubric]`, or silently drop. The empirical trigger comes from the spec-authority-enforcement eval run, where two high-value domain findings ("标记稀释效应" and "Agent 层职责混淆") were both classified as `[beyond-rubric]` by the Scorer, raising the question of whether the Scorer's mapping decisions systematically discount precisely those findings that fall outside the rubric's coverage -- which is, by definition, where the freeform expert adds the most value.

The proposed solution inserts a Pre-Revision phase (labeled Phase 0.5) between the finding extraction step and the Scorer loop. Freeform findings are reformatted as ATTACK_POINTS and fed directly to the existing Reviser. The Scorer then performs a blind review of the pre-revised document (no freeform injection), and the standard Scorer-Reviser iteration loop proceeds from there. The total iteration budget remains constant; pre-revision consumes iteration 0, reducing the available Scorer cycles by one. The `freeform-injection.md` rule file is deprecated entirely.

The core assumptions underpinning this design are: (1) the Reviser can produce meaningful edits from the flattened finding format without the expert's full narrative context, (2) blind Scorer review after pre-revision provides an independent quality signal without losing critical context about what was changed and why, (3) losing one Scorer-Reviser cycle from the total budget is an acceptable trade for direct finding-to-revision fidelity, and (4) deprecating the injection pathway is reversible at negligible cost.

---

## Section 2: Key Risk Identification

问题：Scorer 盲审丧失了 Pre-Revision 变更的上下文，可能导致误判修订质量。提案 Decision 2 中声明："Pre-reviser 修好了问题，Scorer 自然给高分；如果修坏了，Scorer 自然扣分。带溯源可能引入确认偏误。" 这个论证在信息论上是不对称的。盲审确实消除了确认偏误，但同时也消除了 Scorer 理解"文档为什么在 pre-revision 阶段被修改"的能力。当 Pre-Reviser 对 proposal 做了较大结构变更后，Scorer 看到的是一个可能面目全非的文档，无法区分哪些内容是原始提案的一部分、哪些是 pre-revision 新增的。这种盲审不是"独立评估当前版本质量"，而是"对一份缺少编辑历史的文档做评分"。具体场景：如果 pre-reviser 为了回应某个 finding 删除了一个有争议的段落，Scorer 无法判断该段落的缺失是有意的修订还是原始文档的缺陷。Scorer 可能会基于该段落的缺失而生成新的攻击点，触发不必要的后续修订循环，浪费本已减少的 iteration 预算。

风险：ATTACK_POINTS 的扁平化格式 (`- **[severity]** summary | 原文引用: "quote"`) 会鼓励表面级修补而非深层修订。提案 Decision 4 声明这个格式"已提供足够的权重信号和溯源"，但对比 freeform-injection.md 中的完整注入块，现有格式包含了明确的指令：要求 Scorer 对每个 finding 做 rubric 维度映射、标注 beyond-rubric、以及标注与 rubric 的分歧。而 pre-revision 的 ATTACK_POINTS 只提供 severity + summary + quote，缺少 finding 的完整论证链（问题陈述 -> 证据 -> 影响分析 -> 建议改进方向）。Reviser protocol 的 attack-point 处理策略是机械的表格映射（Vague language -> Replace; Missing section -> Add; 等等）。对于复杂的领域级发现（如"标记稀释效应"），这种扁平格式可能导致 Reviser 选择一个最接近的 fix strategy 而非真正理解专家的深层意图。提案的 Key Risks 表格中也承认了这一点："Pre-reviser 机械回应 findings，不理解专家深层意图"，但给出的缓解措施"盲审 Scorer 独立验证修订质量"只是事后检测，而非预防性设计。

问题：Iteration 预算削减的具体影响未被量化。提案声明 `--iterations 3` 会变成 iteration 0 (pre-revision) + iteration 1-2 (Scorer loop)，即 Scorer-Reviser 循环从 3 次减为 2 次。但当前 SKILL.md 中 Iteration Initialization 的逻辑是 `ITERATION = 1`，`MAX_ITERATIONS = resolved value`，每次循环后 ITERATION 递增。提案将 INITIAL_SCORE 的记录时机从 iteration 1 移到 iteration 1（Scorer 循环的第一个），但 pre-revision 的 iteration 0 不产生分数。这意味着 rollback 比较基线从原始 proposal 变成了 pre-revised proposal。如果 pre-revision 本身就降低了文档质量（Key Risks 表格承认这个可能），但后续 Scorer 循环又进一步恶化，rollback 只能恢复到 pre-revised 版本而非原始版本。这是一个不一致的退化语义。

风险：废弃 `freeform-injection.md` 是单向门，且提案对此风险的评估过于轻率。提案 Key Risks 表格声明："重新创建该 rule 文件的成本极低（纯 prompt 组合规则，无逻辑代码）。" 但 freeform-injection.md 不是孤立的文本文件，它定义了一套完整的信息传递语义：beyond-rubric 处理、contradiction annotation、partial extraction handling、degradation path。这些语义与 scorer-composition.md 中的注入步骤、SKILL.md 中 P0.5 的 `FREEFORM_INJECTION` 状态变量、以及 extraction-prompt.md 的验证逻辑构成了一个相互依赖的系统。废弃 freeform-injection.md 意味着这套语义的整体删除。如果未来发现 pre-revision 方案在某些场景下不如 injection 方案，恢复工作不仅是重建一个文件，而是重新验证所有依赖节点的集成。更重要的是，提案同时要修改 scorer-composition.md（移除 `<injected-freeform-findings>` 组合步骤）和 SKILL.md（P0.5 逻辑重写），这意味着回退涉及多文件协调。

问题：提案声称"仅影响 proposal 类型"但未分析这个约束的实现机制。`freeform-injection.md` 中的 When to Inject 条件明确包含 `The eval type is proposal`，scorer-composition.md 的 Freeform Findings Injection 段落也依赖 `FREEFORM_INJECTION` 状态变量。废弃 injection 机制后，其他 eval 类型（如 prd、design）若未来需要引入类似的多阶段评审，将没有现成的注入框架可复用。这不是一个阻塞性问题，但是一个架构前瞻性问题。

问题：Pre-revision 失败时的 degradation path 定义不够精确。提案 Success Criteria 第 5 条声明"Phase 0 失败时跳过 pre-revision，直接进 Scorer"，但这与当前 SKILL.md 中 Phase 0 Degradation Summary 的行为一致（Phase 0 失败 -> 标准rubric流程）。提案没有定义 Phase 0.5 本身的失败模式：如果 Pre-Reviser 返回 `REVISED: failed` 或产出空报告怎么办？Reviser protocol 有自己的质量检查（word count 不超过 30%、每个 attack point 都要被处理），但 pre-revision 作为"无 Scorer 监督的首次修订"，其失败时的回退路径未被设计。当前架构中，Reviser 的产出总是被 Scorer 验证；但在 pre-revision 中，Reviser 的产出直接成为 Scorer 的输入文档，中间没有质量关卡。

风险：INITIAL_SCORE 基线漂移导致 rollback 语义失效。当前 SKILL.md Step 5 的 rollback 逻辑是 `FINAL_SCORE vs INITIAL_SCORE`，其中 INITIAL_SCORE 在 iteration 1 记录。提案将 iteration 1 的输入从原始 proposal 变成了 pre-revised proposal。因此 INITIAL_SCORE 反映的是 pre-revised 版本的质量，而非原始提案。如果 pre-revision 修坏了文档（假设从"原始 60 分"变成"pre-revised 55 分"），Scorer 在 iteration 1 给出 55 分作为 INITIAL_SCORE，后续迭代即使进一步恶化到 50 分，rollback 也会恢复到 pre-revised 版本（55 分基线），而非用户真正想要的原始版本（60 分基线）。提案在 Key Risks 中提到"rollback 机制兜底（INITIAL_SCORE 在 iteration 1 记录，对比 final score）"，但未意识到 INITIAL_SCORE 的语义已经改变。

问题：Decision 3 "处理全部 findings（不限 severity）"与 Reviser 的保守修改指令存在张力。提案声明"Pre-reviser 的指令中要求只接受可以从文档中验证的事实性修正，对不确定的建议保持 conservative"，但同时又说"人为筛选 severity 会引入偏差。freeform 专家的核心价值是独立视角。"这里存在一个未解决的矛盾：如果所有 findings 都交给 pre-reviser 但 pre-reviser 被要求保守，那么 low-severity 的建议性 findings（如"建议补充竞品分析"）要么被保守地忽略（浪费了 finding 的传递成本），要么被勉强执行（违背保守原则）。提案没有说明 pre-reviser 如何在"处理全部 findings"和"保守修改"之间做出决策。

问题：提案缺少对 Pre-Revision 步骤在 SKILL.md 中的精确编排描述。Decision 5 说"Pre-revision 计入 MAX_ITERATIONS"，但当前 SKILL.md 的流程图中，从 Phase 0 到 Expert Dispatch Table 到 Iteration Initialization 是线性串行的。插入 Phase 0.5 意味着在 P0.5 (Inject Findings) 和 Expert Dispatch 之间需要新增：格式化 findings 为 ATTACK_POINTS -> 构造 EVAL_REPORT_PATH（指向什么？没有 Scorer 报告）-> 调用 Reviser。但 Reviser protocol 要求读取 evaluation report at `{{EVAL_REPORT_PATH}}`，而 pre-revision 阶段不存在 evaluation report。提案声称"复用现有 Reviser，零新 protocol"，但 Reviser protocol 的 Step 1 明确要求读取 eval report。需要一个合成的 eval report 或者修改 Reviser 的输入要求，这两者都不是"零新 protocol"。

---

## Section 3: Improvement Suggestions

建议：为 Pre-Revision 构造合成的 iteration-0 评估报告，解决 Reviser protocol 对 EVAL_REPORT_PATH 的硬性依赖。具体做法：在 Phase 0.5 中，将 freeform findings 格式化为标准的 eval report 格式（包含 SCORE: N/A、DIMENSIONS: from-freeform、ATTACKS: formatted findings），写入 `<DOC_DIR>/eval/iteration-0.md`。这样 Reviser 可以正常读取该文件而不需要修改 protocol。这直接解决了"零新 protocol"主张与 Reviser protocol 硬性要求之间的矛盾。采用此建议后，pre-revision 的改动范围限定在 SKILL.md 的流程编排逻辑和 findings 格式化步骤，Reviser protocol 和 reviser-composition.md 真正保持不变。

建议：在 Pre-Revision 阶段保留原始 proposal 的备份作为 rollback 基线，而非仅依赖 INITIAL_SCORE。当前 Step 1.5 已经创建了 `${DOC_DIR}.bak`，但该备份在 rollback 时用于恢复"pre-revised 版本"还是"原始版本"取决于 rollback 触发时机。具体方案：pre-revision 完成后、Scorer 循环启动前，保存一份 pre-revision checkpoint（如 `${DOC_DIR}.pre-revised`），并保留原始 `.bak` 不被覆盖。Rollback 逻辑区分两个层级：Scorer 循环内的 rollback 恢复到 pre-revised checkpoint；整体流程的 rollback 恢复到原始 `.bak`。这解决了 INITIAL_SCORE 基线漂移的风险，使得 rollback 语义与用户直觉一致（"恢复到我评估前的版本"）。

建议：为 ATTACK_POINTS 格式增加轻量级的论证上下文字段，缓解扁平化导致的信息损失。当前格式是 `- **[severity]** summary | 原文引用: "quote"`，建议扩展为 `- **[severity]** summary | 原文引用: "quote" | 期望改进方向: <one-line direction>`。这个 direction 字段不需要完整复述专家的论证链，只需要给出一个动词短语（如"补充脑裂恢复策略"或"量化回滚触发阈值"），帮助 Reviser 选择正确的 fix strategy 而非机械匹配。这个扩展保持格式的简洁性（不回退到注入完整 narrative），但提供了比当前扁平格式更多的语义信息。对应的风险是 Reviser 对复杂 findings 的表面级修补。

建议：将 freeform-injection.md 的废弃改为条件性禁用而非删除，降低单向门风险。具体做法：在 freeform-injection.md 文件头部添加 `deprecated: true` 标记和废弃说明，保留完整的注入语义定义以备回退。同时修改 scorer-composition.md 中的注入逻辑为条件判断：当 pre-revision 模式启用时跳过注入，否则走原路径。这样在验证 pre-revision 方案的效果之前，保留了低成本回退的能力。对应的风险是废弃 freeform-injection.md 作为单向门的不可逆性。

建议：定义 Phase 0.5 的显式 degradation path，弥补当前提案中 pre-revision 失败模式的空白。具体方案：如果 Pre-Reviser 返回失败或产出的变更报告为空，执行两个动作：(1) 丢弃 pre-revision 的所有文件变更（从 `.bak` 恢复），(2) 跳过 pre-revision，Scorer 循环从原始 proposal 开始，FULL_ITERATIONS 全部可用。同时在 eval report 中标注 "Pre-revision 未产出有效修订，已降级为标准流程"。这使得 pre-revision 的失败不消耗用户的迭代预算，保持与 Phase 0 degradation 语义的一致性。

建议：在 Scorer 盲审的 prompt 中增加一个轻量级的"变更区域标注"机制，缓解盲审的上下文丢失问题。不需要暴露 freeform findings 的具体内容，只需在 Scorer 看到的文档中用 HTML 注释标注 pre-revision 修改过的段落（如 `<!-- pre-revised -->`）。Scorer 在评估时可以意识到这些区域近期被修订过，但不知道修订的原因和来源。这保留了盲审的核心价值（不引入确认偏误），同时减少了 Scorer 对 pre-revision 变更区域的误判概率。对应的风险是盲审 Scorer 对 pre-revised 文档的误判。
