# Freeform Expert Review: Freeform Pre-Revision Proposal

**Expert Profile**: Eval Pipeline Information-Flow Architect
**Reviewed Document**: `docs/proposals/eval-freeform-pre-revision/proposal.md`
**Date**: 2026-05-24

---

## Section 1: Background Assessment

This proposal identifies a real information-fidelity defect in the Forge eval pipeline's proposal evaluation flow. The empirical evidence is strong: 47% of freeform findings across two evaluation runs suffered information degradation or total loss when routed through the Scorer's mapping layer. The three loss paths identified in the proposal -- semantic compression via rubric mapping, priority demotion via `[beyond-rubric]` tagging, and silent dropping -- constitute a genuine pipeline design flaw, not a tuning problem.

The proposed solution inserts a Phase 0.5 (Pre-Revision) between finding extraction and the Scorer loop. Freeform findings are reformatted as ATTACK_POINTS, fed directly to the existing Reviser, and the revised document is then passed to the Scorer for what the proposal terms "标注盲审" (annotated blind review). The total iteration budget remains constant, with pre-revision consuming iteration 0. The `freeform-injection.md` rule file is deprecated entirely, and the Scorer no longer receives freeform findings.

The design rests on several core assumptions: (1) the Reviser can produce quality edits from the flattened finding format without the expert's full narrative context; (2) the annotated blind review preserves enough information for the Scorer to evaluate revision quality without introducing confirmation bias; (3) losing one Scorer-Reviser cycle from the iteration budget is an acceptable trade-off; and (4) deprecating the injection pathway is low-cost to reverse. The proposal presents extensive evidence for assumption (1) through its Decision 4 analysis and the industry context section. The trade-off analysis across four design options (Decision 1) is thorough.

The proposal's strongest contribution is the Industry Context section, which maps three established peer-review practices (SIGPLAN meta-reviewer rules, Gerrit +2, MT-Bench multi-judge) to the Forge pipeline and explicitly identifies where the LLM-vs-human trust asymmetry breaks the analogy. This kind of honest analogy-failure analysis is uncommon in design proposals and gives the reader confidence in the author's intellectual rigor.

---

## Section 2: Key Risk Identification

The proposal's information-theoretic analysis of the current pipeline is convincing, but the proposed replacement introduces its own information-fidelity challenges that are not fully addressed.

风险：标注盲审的"折中"设计可能在实践中退化为全量溯源或完全盲审，缺乏对退化方向的预见和控制。提案 Decision 2 声称："Scorer 知道哪些区域被改过（`<!-- pre-revised: high -->`），但不知道为什么改（原始 findings 内容不暴露）"并补充"severity 标记帮助 Scorer 分配注意力权重：high 区域值得更仔细检查，low 区域快速扫过"。但在信息论层面，标注盲审引入了一个隐含的信息通道：severity 标记本身就是一种溯源。当 Scorer 看到 `<!-- pre-revised: high -->` 时，它知道该区域被专家认定为高严重性问题并做了修订。这不等于知道具体 finding 内容，但足以让 Scorer 在该区域投入更多审查注意力——这与提案批评的"确认偏误"在方向上是一致的，只是程度较轻。提案的 Scorer prompt 补充指令声明"severity 标记供注意力分配参考，不影响评分标准"，但 prompt 指令对 LLM 行为的约束力有限。如果实测中 Scorer 对 high-severity 标注区域的攻击率系统性高于同等质量但无标注的区域，这恰恰说明 severity 标记引入了提案本想消除的偏误——只是从"确认 findings 所指问题"变成了"确认标注区域有问题"。提案的 Key Risks 表格中"标注盲审假阳性"条目承认了这个风险但将其 Likelihood 评估为 Medium，Mitigation 仅依赖 prompt 指令约束，未提出可度量的检测机制。

问题：提案在 Decision 4 中宣称"复用现有 Reviser，最小 protocol 适配"，但实际描述的自适应工作已经偏离了"复用"的含义。提案声明"SKILL.md 编排层构造合成 eval report（`iteration: 0` + ATTACK_POINTS + 空 rubric）"并估计"SKILL.md 需新增约 20 行编排代码——这不是完全的'零新'，而是最小适配"。这个表述在两个层面上值得审视。第一，构造合成 eval report 本身是一种隐式 protocol 扩展：Reviser protocol Step 1 读取 eval report 并期望从中获取 attack points、rubric 分数和 dimension breakdown。合成 report 中 rubric 维度全部标记 N/A 的做法需要验证 Reviser protocol 确实能正确处理这种输入——提案声称"此行为已对照 Reviser protocol 验证"，但验证结果未在 proposal 中展示，仅是断言。第二，提案在 Implementation Estimate 中将 SKILL.md 新增代码修正为"~40 行代码"（原估算 20 行），翻倍的修正本身就说明对工作量的初始估计过于乐观。如果合成 report 构造的边缘情况比预期多（如 Reviser 对空 rubric 的 fallback 行为不一致），实际工作量可能继续增长。

风险：Iteration 预算削减对低配置场景的影响被低估。提案 Decision 5 承认"`--iterations 2` 是当前最低有用配置，pre-revision 将 Scorer 循环从 2 次减为 1 次（减少 50%）"并声明会输出 warning。但 proposal 的 Success Criteria 第 6 条要求"iteration-1 Scorer 盲审评分 >= 同一 proposal 无 pre-revision 的 Scorer 盲审评分"，这个成功标准的验证本身就需要额外的 eval 运行（一次带 pre-revision，一次不带）。如果验证结果显示 pre-reversion 版本得分低于 baseline，提案没有定义下一步行动——是回退方案，还是接受质量下降换取 finding 覆盖率？这是一个未闭合的决策分支。

问题：提案对 `freeform-injection.md` 废弃的风险评估存在内部矛盾。Decision 6 声称"Scorer 不再注入 freeform findings（盲审），因此 `rules/freeform-injection.md` 整个废弃"，Key Risks 表格中也将此列为 Low Likelihood / Medium Impact。但在 Key Risks 的详细描述中，同一行承认："恢复成本不仅是重建单个 rule 文件，还需同步恢复 scorer-composition.md 中的注入逻辑、SKILL.md 中 P0.5 的编排。依赖链横跨 3 个文件，恢复需全链路回归测试。" 这与 Low Likelihood 的评估存在张力：如果恢复确实涉及 3 个文件的协调修改和全链路回归测试，Impact 至少应为 High，因为回退成本直接决定了方案的 reversibility。提案的架构承诺部分试图缓解这个问题："恢复需同步修改 freeform-injection.md + scorer-composition.md + SKILL.md P0.5 逻辑共 3 个文件，但每处修改均为 prompt 组合编排，无复杂逻辑代码。" 但"无复杂逻辑代码"不等于"无复杂语义依赖"——freeform-injection.md 定义的 beyond-rubric 处理、contradiction annotation、partial extraction handling 等语义与 Scorer 的输出格式深度耦合，恢复时需要确保这些语义的完整还原。

风险：提案的 rollback 基线语义在 pre-revision 引入后发生了隐式变更，但提案对这一变更的处理不够彻底。当前 SKILL.md 的 rollback 逻辑以 INITIAL_SCORE 为对比基线，INITIAL_SCORE 在 iteration 1 记录。提案将 iteration 1 的输入从原始 proposal 变为 pre-revised proposal，因此 INITIAL_SCORE 反映的是 pre-revised 版本的质量。如果 pre-revision 降低了文档质量（Key Risks 表格承认此可能），Scorer 给出的 INITIAL_SCORE 可能低于原始文档的合理分数。后续 Scorer-Reviser 循环即使将分数提升到原始水平，也可能低于 target 而触发不必要的额外迭代。提案在 Key Risks 中列出"INITIAL_SCORE 基线漂移"并声明"Rollback 对比点改为 Phase 0 原始快照（非 INITIAL_SCORE）"，但 Scope 表格中 rollback 相关修改仅列为"~5 行"——这意味着变更仅是引用替换，而非重新设计 rollback 的语义层级。

问题：提案声称 Pre-Revision"占用 iteration 0"但 iteration 计数语义需要更精确的定义。当前 SKILL.md Iteration Initialization 设置 `ITERATION = 1`，循环在 `ITERATION <= MAX_ITERATIONS` 时继续。如果 pre-revision 被定义为 iteration 0，那么 Scorer 循环从 iteration 1 开始——但 iteration 0 是否计入 `ITERATION` 变量？如果计入，则 ITERATION 初始化为 0，pre-revision 执行后 ITERATION 递增为 1，Scorer 循环在 `ITERATION <= MAX_ITERATIONS` 时继续。这意味着 `--iterations 3` 实际给了 pre-revision 1 次 + Scorer 循环 3 次（iteration 1-3），与提案声明的"iteration 0 + iteration 1-2"不一致。提案需要明确 iteration 0 是否消耗 ITERATION 计数器，以及 gate 逻辑如何处理这个偏移。

风险：Decision 3 的三层分类策略（事实性修正/结构建议/主观偏好）依赖 Pre-Reviser 的判断能力，但未定义判断失败时的检测机制。提案声明"事实性修正（可定位原文缺陷）：直接编辑"和"主观偏好：标注 not actionable，不编辑"。分类边界本身是模糊的——一条 finding 是否"可验证"取决于 Reviser 对领域知识的掌握程度，而 Reviser 不接收专家 profile（Decision 4 明确声明"不注入专家 profile"）。例如，一条关于"分布式一致性场景下脑裂恢复"的 finding，对于理解分布式系统的 Reviser 是事实性修正，对于不熟悉该领域的 Reviser 可能被误分类为"结构建议"或"主观偏好"。提案的 Key Risks 中"Pre-reviser 机械回应 findings"条目将 Likelihood 评为 Medium，Mitigation 依赖"标注盲审 Scorer 对 `<!-- pre-revised -->` 区域检查修订质量"——但这只是事后检测，无法预防分类错误导致的 finding 被跳过。一旦 finding 被标记为 "not actionable"，它不会出现在 ATTACK_POINTS 中，Scorer 也不会知道它曾经存在。

问题：提案的 Success Criteria 第 6 条（质量可度量性）的 baseline 对比设计存在方法论缺陷。该标准声明"对同一 proposal 先运行一次 `--iterations 1 --no-freeform`（跳过 freeform review），记录 Scorer 盲审分数作为 baseline；再运行完整 pre-revision 流程，取 iteration-1 的 Scorer 盲审分数进行对比"。提案自己承认了"pre-revision 版本的文本长度和结构可能系统性差异"这一混淆因素，但缓解措施仅为"通过控制同一 proposal 内容来缓解"。实际上，两次运行之间不存在真正的控制变量：baseline 运行的 Scorer 看到的是原始文档，pre-revision 运行的 Scorer 看到的是修订后文档加上 `<!-- pre-revised -->` 标注。即使 Scorer 不受标注影响，文档长度差异本身就会影响 LLM 的注意力分配和评分行为。这不是一个可靠的 A/B 测试设计。

---

## Section 3: Improvement Suggestions

建议：在 Pre-Revision 阶段保存原始 proposal 的完整快照作为 rollback 的 ground truth 基线，而不仅仅是依赖 pre-revised 版本的 INITIAL_SCORE。具体做法：在 Phase 0.5 启动前，将原始 proposal 文件复制到 `<DOC_DIR>/eval/baseline-snapshot/`（或等价位置），确保无论 pre-revision 是否成功、后续 Scorer 循环是否恶化，都能恢复到用户提交评估前的原始状态。Rollback 逻辑区分两级：Scorer 循环内的 rollback 恢复到 pre-revised checkpoint；整体流程的 rollback 恢复到 baseline snapshot。这解决了 INITIAL_SCORE 基线漂移的风险，使得 rollback 语义与用户直觉一致。提案采纳此建议后，Scope 表格中的 rollback 修改从简单的"引用替换"变为两级 rollback 语义设计，但这正是该变更应有的完整度。

建议：为 freeform-injection.md 采用条件性废弃而非物理删除，降低单向门风险。具体做法：在文件头部添加 `status: deprecated` frontmatter 标记，保留完整的注入语义定义，并在 scorer-composition.md 中将注入逻辑改为条件分支（`if not pre_revision_mode: inject freeform findings`）。这保留了低成本回退路径，同时不增加运行时复杂度（条件分支在 SKILL.md 编排层而非每次 Scorer 调用时执行）。对应的风险是废弃 freeform-injection.md 作为单向门的不可逆性。提案采纳此建议后，Decision 6 的声明从"整个废弃"变为"条件性禁用"，Key Risks 表格中该条目的 Impact 可以合理降至 Low。

建议：扩展 ATTACK_POINTS 格式，增加一个轻量的"期望改进方向"字段。当前格式 `- **[severity]** summary | 原文引用: "quote"` 可扩展为 `- **[severity]** summary | 原文引用: "quote" | 期望改进方向: <动词短语>`。这个字段不回退到注入完整 narrative，仅为 Reviser 提供一个方向锚点，帮助它在 Reviser protocol 的 fix strategy 表格（Vague language -> Replace; Missing section -> Add 等）中选择更合适的策略。对应的风险是扁平化格式导致 Reviser 对复杂 findings 的表面级修补。提案采纳此建议后，Decision 4 的格式定义需要一行扩展，extraction-prompt.md 的输出规范需要增加一个可选字段。

建议：定义 Pre-Revision 失败时的完整 degradation path，确保 pre-revision 异常不消耗用户的 iteration 预算。当前提案的 Phase 0.5 失败处理表格列出了四种场景，但所有降级都"跳过 pre-revision，直接进 Scorer"。问题是：如果 pre-revision 已经修改了文档但产出格式异常，"丢弃 pre-revision 结果"意味着需要从备份恢复文档，然后 Scorer 循环使用完整的 MAX_ITERATIONS 预算。提案应明确声明：pre-revision 失败时，iteration 计数器不递增，Scorer 循环享有全部原始预算。这保持了与 Phase 0 degradation 语义的一致性。采纳此建议后，Phase 0.5 失败处理表格中的每种场景都需要增加一行"ITERATION 不递增，MAX_ITERATIONS 保持原值"。

建议：为标注盲审的偏误风险建立可度量的检测机制，而非仅依赖 prompt 指令约束。具体做法：在 eval report 中要求 Scorer 对每个 `<!-- pre-revised -->` 标注区域单独记录其 attack density（单位文本长度内的攻击点数量），并与未标注区域的 attack density 进行比较。如果标注区域的 density 系统性偏高（例如连续 N 次 eval 中偏高超过阈值），则触发"标注偏误告警"，提示需要调整标注策略或移除 severity 标记。这为 Decision 2 的折中设计提供了实证反馈闭环。采纳此建议后，Scorer 的 report template 需要增加一个 section 用于记录标注区域与未标注区域的 attack density 对比。
