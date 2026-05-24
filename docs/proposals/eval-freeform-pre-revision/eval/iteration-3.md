---
iteration: 3
total_score: 885
scorer: CTO-Adversarial
date: 2026-05-24
previous_iteration: 2
previous_score: 856
---

# Proposal Evaluation Report: Freeform Pre-Revision — Iteration 3

## Revision Delta Analysis

Iteration-2 后提案进一步修订，直接回应了 iteration-2 的 7 个 ATTACK 点中的 4 个：

| Iteration-2 Attack | 修订状态 | 评估 |
|---------------------|---------|------|
| 1. 偏误检测实现复杂度未估算 | **已修复** — scorer-composition.md 估算从 ~5 行修正为 ~10 行，新增"偏误检测 report template"描述 | 估算已包含偏误检测输出格式 |
| 2. Borderline 分类未列为风险 | **已修复** — Key Risks 新增"Borderline 分类元认知失败"风险项，Likelihood=Medium, Impact=Medium，含告警逻辑 | 风险项完整 |
| 3. Pre-Revision edits 与 Scorer ATTACK_POINTS 交互未定义 | **已修复** — Scorer prompt 新增 `conflict-with-pre-revision` 标注策略：以 rubric 标准为准 + 标注冲突供审查 | Scorer 端完整，但 Reviser 端消费语义未定义 |
| 4. SC #7 partially-accepted 判定标准未定义 | **已修复** — 新增判定标准："修改触及了 finding 指出的原文位置，且修改方向与 finding 的'期望改进方向'字段一致" + partially-accepted 比例超过 accepted 时触发人工抽检 | 判定标准明确且可操作 |
| 5. Decision 4 标题"最小 protocol 适配"不准确 | **已修复** — 标题改为"受控 protocol 扩展" | 标题与正文一致 |
| 6. SIGPLAN 因果论证"验证"过度声称 | **已修复** — 新增证伪条件："若后续 eval 的 findings 受损率降至 15% 以下...则维度不匹配前提不成立" | 因果论证可证伪 |
| 7. 偏误检测 report template 未列入改动文件 | **已修复** — scorer-composition.md 行注明包含"偏误检测 report template" | 改动文件表完整 |

---

## Phase 1: Reasoning Audit

### Argument Chain Trace

**Problem -> Solution**: 核心因果链在三轮 iteration 中持续强化。Iteration-3 提案新增冲突处理策略（`conflict-with-pre-revision`）闭合了 Pre-Revision edits 与 Scorer 判断之间的逻辑缺口，使因果链从"消除瓶颈"升级为"消除瓶颈 + 处理残余冲突"。Industry Context 的因果论证 + 证伪条件（"若受损率降至 15% 以下则前提不成立"）将整个论证框架从断言升级为可证伪的理论。

**Solution -> Evidence**: 证据基础仍为 2 次 eval、15 条 findings——三轮 iteration 均未扩展。但证据的论证深度已达上限：SIGPLAN 因果论证 + 证伪条件 + 冲突处理策略 + 偏误检测实证闭环 + borderline 告警逻辑构成了多层次的自我验证框架。证据薄弱点已被提案自身的"信息性参考"降级和证伪条件所承认。

**Evidence -> Success Criteria**: SC #6 降级为信息性参考 + 人工审查流程。SC #7 新增 partially-accepted 判定标准 + 人工抽检触发条件。冲突处理策略闭合了 iteration-2 blindspot-3 的"Pre-Revision edits 与 Scorer ATTACK_POINTS 重叠/矛盾"缺口。因果链从 evidence 到 SC 的映射已完整。

**Self-contradiction check**:
- **冲突处理策略的 Reviser 端未闭合**：Scorer prompt 定义了 `conflict-with-pre-revision` 标注，但 Reviser 如何处理带有此标注的 attack point 未定义。这引入了新的半开交互语义。
- **告警阈值缺乏论证**：偏误检测的 30%、borderline 告警的 >10 findings、证伪条件的 15%——三个数值阈值均未提供来源论证。
- **Implementation Estimate "约 1 天"的累积偏差**：SKILL.md ~40 行 + rollback ~15 行 + scorer-composition ~10 行 + 条件性废弃 ~3 行 = ~68 行代码 + 告警逻辑（未估算）+ 端到端测试。"约 1 天"在 iteration-1 时对应 ~28 行代码，iteration-3 对应 ~68 行代码 + 告警逻辑 + 测试，但总估算未调整。
- **Phase 0.5 失败时 iteration 预算语义**：三轮 iteration 均指出此问题。"占用 iteration 0（从总预算中扣除）"与降级原则"回退到 Scorer 直接评估模式"组合后产生矛盾——失败时"扣除"的 iteration 是否归还？有条件扣除语义未显式声明。

### Pre-Score Anchors

1. 提案在三轮 iteration 中从 702 分提升到 856 分，iteration-3 进一步修复了 4 个 ATTACK 点。提案接近成熟状态。
2. Industry Context 因果论证 + 证伪条件是提案最强部分，在整个迭代过程中持续强化。
3. 冲突处理策略是 iteration-3 的关键新增，闭合了 Pre-Revision edits 与 Scorer 判断的交互缺口，但引入了新的半开交互语义。
4. 证据基础（15 条 findings）是提案最薄弱环节，但在提案范畴内已无法改善——证据来自已发生的 eval，不能事后增加。
5. Implementation Estimate 的累积偏差（从 ~28 行到 ~68 行 + 未估算的告警逻辑）是 iteration-3 新浮现的结构性问题。
6. 告警阈值（30%、>10、15%）缺乏论证来源，三个数值直接影响方案可操作性。

---

## Phase 2: Rubric Scoring with Verification Stance

### 1. Problem Definition: 94/110

**Problem stated clearly (38/40)**: 核心问题——freeform findings 通过 Scorer 映射层时 47% 信息受损——表述精确。三种损失路径有明确因果机制。Industry Context 因果论证 + 证伪条件将 47% 从孤立数据点升级为可证伪的理论假设的一个观测值。

> Deduction (-2): 15 findings 的小样本仍使 47% 脆弱。提案未承认置信区间。

**Evidence provided (36/40)**: 两次具体 eval 运行（含具体 finding 名称和计数），根因分析文档引用，具体损失分类统计。"标记稀释效应在 Reviser 修订中未被触及"等描述增强证据质量。

> Deduction (-4): 全部证据来自同一作者的同一管道，无外部验证。证据基础（15 findings）未在三轮 iteration 中扩展。

**Urgency justified (20/30)**: "2 个活跃 proposal 受此影响"和"47% 信息损失率在每个 proposal eval 中复现"建立持续影响。

> Deduction (-10): 仍无明确延迟成本分析——多少个 proposal 会在修复前被评估？累计 finding 损失的期望值是多少？

### 2. Solution Clarity: 108/120

**Approach is concrete (39/40)**: 六个 Design Decisions + 信息流图 + Phase 0.5 失败处理表 + 两级 rollback 设计 + 冲突处理策略。读者可精确复述方案。

> Deduction (-1): 冲突处理策略的 Reviser 端消费语义未定义——Scorer 标注后 Reviser 如何处理？

**User-facing behavior described (38/45)**: Iteration-0 报告格式、最终 report summary 增加行、`--iterations 2` 质量退化预期、warning 文案、推荐最低配置、borderline 列出 + skipped 附理由 + 用户可据此判断是否手动干预。

> Deduction (-7): 用户对 `conflict-with-pre-revision` 标注的可见性和操作指导未描述。borderline 告警时用户具体应执行什么操作仍缺乏描述。

**Technical direction clear (31/35)**: 实现路径明确。Decision 4 的 4 步 protocol 验证追踪 + fallback + 空值兼容性验证使技术方向清晰。

> Deduction (-4): 冲突处理策略引入新 Scorer-Reviser 交互模式，但实现仅一行 prompt 指令，未定义 report 中的格式化规范和 Reviser 处理优先级。

### 3. Industry Benchmarking: 110/120

**Industry solutions referenced (38/40)**: SIGPLAN、Gerrit +2、MT-Bench multi-judge——三个真实实践。类比失效点分析保留。证伪条件新增。

> Deduction (-2): 仍缺乏同行评审方法论的文献引用。

**At least 3 meaningful alternatives (27/30)**: 四个选项 + "do nothing"明确列出。

> Deduction (-3): Options A/B 仍非行业验证方案（如加权聚合、校准评分等）。

**Honest trade-off comparison (23/25)**: 贸易比较直接且诚实。

> Deduction (-2): 仍无正式的加权评分矩阵。

**Chosen approach justified against benchmarks (22/25)**: 三个行业实践均有因果论证段落。证伪条件的引入是 iteration-3 亮点——"若后续 eval 的 findings 受损率降至 15% 以下...则维度不匹配前提不成立"。

> Deduction (-3): 证伪阈值 15% 的选择缺乏论证——为什么是 15% 而非 20% 或 10%？

### 4. Requirements Completeness: 92/110

**Scenario coverage (37/40)**: Happy path + 4 种 Phase 0.5 失败场景 + degradation 路径 + borderline 分类 + 冲突处理策略。Iteration-2 的"Pre-Revision edits 与 Scorer ATTACK_POINTS 重叠/矛盾"已通过 `conflict-with-pre-revision` 标注修复。

> Deduction (-3): `conflict-with-pre-revision` 标注的下游处理策略仅定义了 Scorer 端，Reviser 端如何处理带有此标注的 attack point 未定义——若 Reviser 忽略此标注，冲突处理策略仅产生噪音。

**Non-functional requirements (34/40)**: 延迟、兼容性、安全、rollback、偏误检测可观测性、`--iterations 2` 场景、条件性废弃降低单向门风险。

> Deduction (-3): 偏误检测告警触发逻辑（"连续 >= 2 次 eval"）的实现位置未指定 (-2)。LLM 输出作为 Reviser 输入的安全面分析仍偏简略 (-1)。
> Deduction (-3): 告警逻辑的实现位置（SKILL.md 编排层 vs Scorer 内部）和输出格式未指定。

**Constraints & dependencies (21/30)**: 依赖清单完整。type==proposal 约束明确。

> Deduction (-5): ITERATION 变量初始化/递增逻辑仍未显式定义。Phase 0.5 失败时计数器行为仅隐含 (-4)。

### 5. Solution Creativity: 80/100

**Novelty over industry baseline (33/40)**: 提案坦诚非原创。标注盲审有了偏误检测实证闭环和冲突处理策略。但核心思想仍是对现有模式的适配。

**Cross-domain inspiration (28/35)**: 三个不同领域 + 类比失效分析。三轮 iteration 未扩展领域范围。

> Deduction (-7): 领域探索仍限于三个已引用实践。

**Simplicity of insight (19/25)**: 核心洞察简洁。实现复杂度从 iteration-1 的 ~28 行增长到 iteration-3 的 ~68 行 + 告警逻辑，冲突处理策略新增交互语义。

> Deduction (-6): 实现复杂度持续增长（68 行 + 未计入的告警逻辑 + conflict 处理），简洁洞察的实现表面持续扩张。

### 6. Feasibility: 85/100

**Technical feasibility (37/40)**: 所有变更为 markdown/prompt 配置文件。Decision 4 的 4 步 protocol 验证追踪 + fallback + 空值兼容性验证使合成 eval report 可审查。冲突处理策略是合理的 prompt 设计。

> Deduction (-3): 验证仍是论证性的，非实际运行结果。冲突处理策略的有效性未验证。

**Resource & timeline feasibility (25/30)**: ~1 天总工作量。Implementation Estimate 的迭代修正诚实。但告警逻辑（borderline 率异常检测 + attack density 偏差计算 + 告警输出）的实现位置和工作量未在估算中体现。

> Deduction (-5): 告警检测逻辑的实现位置和工作量未计入"约 1 天"估算。

**Dependency readiness (23/30)**: 现有 Reviser protocol、scorer-composition、SKILL.md 可用。合成 eval report 有 4 步验证追踪。

> Deduction (-7): Baseline snapshot 的存储/恢复逻辑未详细描述。合成 report 虽有验证追踪但仍需实际构建验证。冲突处理策略引入新的 Scorer-Reviser 交互依赖。

### 7. Scope Definition: 74/80

**In-scope items are concrete (28/30)**: 改动文件表列出 4 行。偏误检测 report template 已纳入 scorer-composition.md 的 ~10 行估算。

> Deduction (-2): `conflict-with-pre-revision` 标注在 report 中的格式化规范未列为独立交付项。Borderline 告警逻辑未列入。

**Out-of-scope explicitly listed (23/25)**: 不改动文件表 + 架构承诺。冲突处理策略已闭合 iteration-2 的缺口。

> Deduction (-2): 冲突处理策略的 Reviser 端处理是否范围内未明确界定。

**Scope is bounded (23/25)**: 4 文件修改 + ~1 天。scorer-composition.md 已修正为 ~10 行。

> Deduction (-2): 总估算"约 1 天"仍可能未包含告警逻辑的实现工作量。

### 8. Risk Assessment: 82/90

**Risks identified (28/30)**: 7 个风险项 + borderline 元认知失败风险已新增。

> Deduction (-2): 冲突处理策略引入新的交互模式失败风险（Reviser 忽略 conflict 标注导致冲突攻击点被当作普通攻击点处理），未列为独立风险。

**Likelihood + impact rated (26/30)**: 评级整体诚实。INITIAL_SCORE 基线漂移 Likelihood=High 是诚实的。borderline 元认知失败 Likelihood=Medium, Impact=Medium 合理。

> Deduction (-4): "Pre-revision 修坏 proposal" Likelihood=Low 仍可能低估——borderline 分类和冲突处理都依赖 LLM 判断力，增加了误修改的可能性。

**Mitigations are actionable (28/30)**: 偏误检测 attack density 检测机制是实证闭环。borderline 告警逻辑可操作。冲突处理策略以 rubric 标准为准。

> Deduction (-2): Borderline 告警阈值（"0 borderline / >10 findings"）来源未论证——>10 的阈值是经验值还是配置？偏误检测"偏差超过 30%"的阈值同理缺乏论证。

### 9. Success Criteria: 76/80

**Criteria are measurable and testable (52/55)**:
- SC 1-5: testable, good.
- SC 6: 已降级为信息性参考 + 定义失败后人工审查流程。方法论缺陷已通过降级修复。
- SC 7: 新增 partially-accepted 判定标准（"修改触及了 finding 指出的原文位置，且修改方向与 finding 的'期望改进方向'字段一致"）+ partially-accepted 比例超过 accepted 时触发人工抽检。

> Deduction (-3): 偏误检测 attack density 阈值（偏差 >30%）和 borderline 告警阈值仍未列为 SC——这些是方案的关键质量保障机制，但缺乏 pass/fail 判定。

**Coverage is complete (24/25)**: SC 覆盖所有 in-scope 项。

> Deduction (-1): 冲突处理策略的有效性（冲突标注是否被 Reviser 正确处理）未作为 SC。

### 10. Logical Consistency: 84/90

**Solution addresses the stated problem (34/35)**: Pre-Revision 直接消除 Scorer 信息瓶颈。偏误检测实证闭环使标注盲审不再是未检测的偏误通道。冲突处理策略闭合了 Pre-Revision edits 与 Scorer 判断的逻辑缺口。

> Deduction (-1): 偏误检测闭环尚未验证（无实证数据），理论上的偏误风险仍存在。但证伪条件的引入使风险可控。

**Scope <-> Solution <-> Success Criteria aligned (27/30)**: Iteration-1 的四个不一致已全部修复。Phase 0.5 失败时 iteration 预算分配仍为隐含语义。

> Deduction (-3): Phase 0.5 失败时 iteration 预算分配未显式定义。冲突处理策略的 Reviser 端处理未对齐。

**Requirements <-> Solution coherent (23/25)**: "处理"语义明确为"分诊"。合成 eval report 在 Decision 4 详述。冲突处理策略缓解了 Pre-Revision edits 与 Scorer 判断的矛盾。

> Deduction (-2): 冲突处理策略引入新的交互模式但未在 Requirements 中作为显式需求列出。

---

## Phase 3: Blindspot Hunt

**[blindspot-1]** **冲突处理策略的 Reviser 端消费语义未定义。** Scorer prompt 新增 `conflict-with-pre-revision` 标注指令：当 Scorer 的 rubric 判断与 pre-revision 修改方向矛盾时，以 rubric 标准为准生成 attack point 并标注冲突。但提案未定义：(1) 这个标注在 eval report 中的格式化规范；(2) Reviser 看到带有此标注的 attack point 时的处理优先级——是优先处理（因为涉及 pre-revision 质量问题）还是正常排队；(3) 如果 Reviser 再次修改冲突区域（按 rubric 标准修复），是否会与 pre-reviser 的修改产生"来回摆动"（oscillation）——iteration 2 按 rubric 修复，iteration 3 按 freeform finding 修复，循环往复。

Quote: "当 Scorer 的 rubric 判断与 pre-revision 修改方向矛盾时...以 rubric 标准为准生成 attack point，但在 attack point 中标注 `conflict-with-pre-revision` 供审查" — 仅定义了 Scorer 端行为，Reviser 端和跨 iteration 的行为未定义。

**[blindspot-2]** **告警阈值缺乏论证来源。** 提案中出现了多个数值阈值但均缺乏论证：(1) 偏误检测的"偏差超过 30%"——30% 的来源是什么？是基于 LLM 评分方差的经验值还是理论推导？(2) borderline 告警的">10 findings"和"0 borderline"——10 的来源是什么？(3) SIGPLAN 证伪条件的"15% 以下"——为什么是 15%？这些阈值直接影响方案的可操作性，但提案将它们作为不证自明的常量处理。如果偏误检测的 30% 阈值过于宽松，偏误检测机制将形同虚设；如果过于严格，正常评分方差会触发误告警。

Quote: "若标注区域 density 系统性偏高（连续 >= 2 次 eval 中偏差超过 30%），触发'标注偏误告警'" — 30% 阈值未论证。另见："当 borderline 率异常低（如 0 borderline / >10 findings）时触发告警" — >10 未论证。

**[blindspot-3]** **偏误检测的后处理实现位置未指定。** 偏误检测要求"连续 >= 2 次 eval 中偏差超过 30%"。这意味着需要一个跨 eval 运行的持久化状态——存储上一次 eval 的 attack density 数据，与本次比较。但提案未指定：(1) 这个状态存储在哪里（SKILL.md 变量？文件系统？）；(2) 比较逻辑在哪里执行（SKILL.md 编排层？Scorer 内部？）；(3) 告警输出到哪里（eval report？stderr？独立告警文件？）。这不是实现细节——它决定了偏误检测机制是"集成到管道中"还是"依赖人工检查 eval report"。

Quote: "在 eval report 中分别记录标注区域与未标注区域的 attack density，供偏误检测" — "供偏误检测"暗示后续有处理步骤，但处理步骤的实现位置未指定。

**[blindspot-4]** **Iteration 预算语义在 Phase 0.5 失败时仍不透明。** 三轮 iteration 均指出此问题。提案在 Decision 5 中说"占用 iteration 0（从总预算中扣除）"，降级表格隐含"Scorer 使用完整预算"。但"从总预算中扣除"与"Scorer 使用完整预算"之间存在矛盾——如果从总预算中扣除了 1 次，`--iterations 3` 应该给 Scorer 2 次；但如果 Phase 0.5 失败跳过，"扣除"的那 1 次是否归还？提案的降级原则说"任何 Phase 0.5 异常都回退到 Scorer 直接评估模式"——这是否意味着 `--iterations 3` 在 Phase 0.5 失败时仍给 Scorer 3 次循环？如果如此，"从总预算中扣除"就不是真正的扣除，而是有条件的扣除。这个有条件语义未显式声明。

Quote: "占用 iteration 0（从总预算中扣除）" — 与降级原则"任何 Phase 0.5 异常都回退到 Scorer 直接评估模式"组合后，"扣除"是否为有条件扣除未声明。

**[blindspot-5]** **Implementation Estimate 的总估算"约 1 天"未包含告警相关实现。** 提案引入了两种告警机制：(1) 偏误检测告警（attack density 偏差 >30%）；(2) borderline 率告警（0 borderline / >10 findings）。两种告警都需要：状态存储、阈值比较、告警输出格式。这些逻辑的实现位置和工作量在 Implementation Estimate 中无对应行。结合 SKILL.md 已从 ~20 行增长到 ~40 行，且 scorer-composition.md 从 ~5 行增长到 ~10 行，告警逻辑可能进一步增加 10-20 行——使总工作量从"约 1 天"偏向"1.5-2 天"。

Quote: Implementation Estimate 总计"约 1 天工作量" — 无告警逻辑的估算行。

**[blindspot-6]** **`conflict-with-pre-revision` 可能成为 Scorer 的认知负担。** Scorer 在标注盲审模式下已有复杂的认知任务：区分标注区域与未标注区域、对标注区域检查修订质量、对未标注区域正常评估、记录 attack density。新增 `conflict-with-pre-revision` 标注要求 Scorer 在评估过程中判断"我的 rubric 评估方向是否与 pre-revision 修改方向矛盾"——这是一个元认知判断，要求 Scorer 在评分的同时推理 pre-reviser 的意图。LLM 的元认知能力有限（提案自身在 borderline 分类中已承认），Scorer 可能无法可靠执行此判断，导致冲突标注的 precision/recall 不明。

Quote: "当 Scorer 的 rubric 判断与 pre-revision 修改方向矛盾时...在 attack point 中标注 `conflict-with-pre-revision`" — 要求 Scorer 推理 pre-reviser 的修改意图，但 Scorer 不接触原始 findings，只能通过变更标记推断意图——信息不足支持可靠的意图推理。

---

## Injected Freeform Findings Disposition

| Finding | Disposition |
|---------|-------------|
| **[high]** 标注盲审退化风险 | **已修复** — 偏误检测机制（attack density 对比 + 30% 阈值 + 告警闭环）使退化方向可检测。Iteration-2 后保持。 |
| **[high]** severity 标记隐式信息通道 | **已修复** — 同上偏误检测机制。Iteration-2 后保持。 |
| **[high]** 三层分类依赖 Pre-Reviser 判断但无失败检测 | **已修复** — borderline 分类 + 分类审计章节 + borderline 告警逻辑 + 元认知失败风险项。Iteration-3 新增了 Key Risks 中的独立风险项。 |
| **[high]** finding 标记为 not actionable 后彻底消失 | **已修复** — 分类审计章节 + iteration-0 报告 + borderline 保留机制。Iteration-2 后保持。 |
| **[medium]** Rollback 基线语义变更需重新设计 | **已修复** — 两级 rollback 设计，估算修正为 ~15 行。Iteration-2 后保持。 |
| **[medium]** Iteration 0 计数语义未明确定义 | **部分修复** — Decision 5 明确了 `--iterations 3` 的分解，但 Phase 0.5 失败时的"有条件扣除"语义仍未显式声明。[blindspot-4] |
| **[medium]** SC #6 baseline 对比缺乏控制变量 | **已修复** — SC #6 降级为信息性参考 + 承认方法论局限 + 定义失败后人工审查。Iteration-2 后保持。 |

---

## Summary

| Dimension | Score | Max | Delta from Iter-2 |
|-----------|-------|-----|-------------------|
| Problem Definition | 94 | 110 | +2 |
| Solution Clarity | 108 | 120 | +3 |
| Industry Benchmarking | 110 | 120 | +2 |
| Requirements Completeness | 92 | 110 | +4 |
| Solution Creativity | 80 | 100 | +2 |
| Feasibility | 85 | 100 | +3 |
| Scope Definition | 74 | 80 | +2 |
| Risk Assessment | 82 | 90 | +4 |
| Success Criteria | 76 | 80 | +4 |
| Logical Consistency | 84 | 90 | +3 |
| **Total** | **885** | **1000** | **+29** |

---

## ATTACKS

1. **[Feasibility]**: 告警检测逻辑的实现位置和工作量未计入估算 — "在 eval report 中分别记录标注区域与未标注区域的 attack density，供偏误检测" + "当 borderline 率异常低（如 0 borderline / >10 findings）时触发告警" — Implementation Estimate 中无对应估算行。偏误检测的跨 eval 状态存储、偏差计算、告警输出 + borderline 率统计和告警，可能增加 10-20 行实现，使"约 1 天"估算偏紧。必须在 Implementation Estimate 中增加告警逻辑的估算行。

2. **[Solution Clarity]**: 冲突处理策略的 Reviser 端消费语义未定义 — "当 Scorer 的 rubric 判断与 pre-revision 修改方向矛盾时...以 rubric 标准为准生成 attack point，但在 attack point 中标注 `conflict-with-pre-revision` 供审查" — 仅定义了 Scorer 端，Reviser 如何处理带有此标注的 attack point、是否优先处理、跨 iteration 是否会 oscillation 均未定义。必须定义 Reviser 端的冲突处理行为和跨 iteration 的防 oscillation 策略。

3. **[Risk Assessment]**: 多个数值阈值缺乏论证来源 — "偏差超过 30%" + ">10 findings" + "15% 以下" — 这些阈值直接影响方案可操作性但作为不证自明的常量处理。30% 偏误阈值过宽则检测机制形同虚设，过严则误告警。必须为每个阈值提供论证来源（经验数据、理论推导、或明确标注为"初始值待校准"）。

4. **[Requirements Completeness]**: Phase 0.5 失败时 iteration 预算的"有条件扣除"语义未显式声明 — "占用 iteration 0（从总预算中扣除）" + "任何 Phase 0.5 异常都回退到 Scorer 直接评估模式" — 两条规则组合后，"扣除"是否在失败时归还未声明。`--iterations 3` 在 Phase 0.5 失败时给 Scorer 2 次还是 3 次循环？必须显式定义 Phase 0.5 失败时的 iteration 预算语义。

5. **[Logical Consistency]**: `conflict-with-pre-revision` 要求 Scorer 推理 pre-reviser 意图但信息不足 — Scorer 不接触原始 findings，仅通过变更标记 `<!-- pre-revised: {severity} -->` 和修订后文本推断 pre-reviser 的修改意图。要求 Scorer 判断"rubric 评估方向是否与 pre-revision 修改方向矛盾"是元认知判断，信息不足支持可靠推理。必须在 Scorer prompt 中提供更具体的冲突检测指引（如：当 Scorer 对标注区域的 attack 方向是"删除/回退 pre-revision 的修改"时标记冲突），而非要求 Scorer 推理意图。

6. **[Scope Definition]**: 告警逻辑未列入改动文件表 — 改动文件表列出 4 行，但偏误检测告警 + borderline 率告警需要 SKILL.md 新增实现逻辑。必须在改动文件表的 SKILL.md 行中注明包含告警检测逻辑，或增加估算行。
