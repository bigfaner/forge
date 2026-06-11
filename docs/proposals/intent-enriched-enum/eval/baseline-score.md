# Baseline Score: Intent Enriched Enum

**Evaluator**: CTO Adversary (Baseline)
**Date**: 2026-05-31
**Document**: `docs/proposals/intent-enriched-enum/proposal.md`

---

```
SCORE: 640/1000
DIMENSIONS:
  Problem Definition: 75/110
  Solution Clarity: 70/120
  Industry Benchmarking: 30/120
  Requirements Completeness: 55/110
  Solution Creativity: 25/100
  Feasibility: 75/100
  Scope Definition: 60/80
  Risk Assessment: 55/90
  Success Criteria: 60/80
  Logical Consistency: 95/90
ATTACKS:
```

---

## Detailed Analysis

### 1. Problem Definition: 75/110

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Problem stated clearly | 30/40 | 核心问题表述清晰：3 值 intent 与 8 值 task type 映射断裂。但"导致下游 pipeline 分支过粗"中的"过粗"是主观判断，缺少量化。什么是"合适的流程"？没有定义判定标准。 |
| Evidence provided | 25/40 | 列出 4 条证据，均可通过代码验证（已验证 brainstorm/SKILL.md Step 4.5 确实对 `coding.fix` 使用启发式；write-prd Intent Detection 表确实只有 3 行）。但缺少实际受影响案例——哪些具体 proposal/feature 因为 3 值枚举走了错误的 pipeline？缺少真实事件回溯。 |
| Urgency justified | 20/30 | "随着 Forge 处理的场景增多"是趋势陈述，不是紧迫性论证。没有回答"如果我们延迟 3 个月会怎样"——cost of delay 未量化。当前启发式虽然不完美，但它到底造成了多少实际返工？ |

**Attacks:**

1. **Problem Definition**: 缺少真实失败案例 — "refactor 和 cleanup 在 write-prd/tech-design 中被完全等同对待" — 需要补充：这个等同对待实际导致了哪次返工或遗漏？没有 concrete incident 的紧迫性论证是推测性的。

### 2. Solution Clarity: 70/120

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Approach is concrete | 25/40 | 三条路径（扩枚举、混合 pipeline、简化推断）方向明确，但核心机制"Pipeline Configuration 表 + Override Signals"没有给出表的具体 schema。读者无法复述 enhancement 在 write-prd 中走什么流程——因为提案自己也没定义。 |
| User-facing behavior described | 25/45 | brainstorm AskUserQuestion 选项从 3 个变为 6 个是明确的用户行为变化。但混合模式的"PRD 内容中的明确信号可以覆盖默认值"对执行 skill 的 LLM 来说是用户行为描述——这个"信号检测"的行为规范完全缺失。LLM 如何判断 PRD 中存在"CLI 命令重命名"信号？靠全文理解还是关键词匹配？ |
| Technical direction clear | 20/35 | "所有变更是 markdown 编辑"给出了技术方向。但 Pipeline Configuration 表的结构、Override Signals 的检测机制、enhancement 的完整 pipeline 行为都是空白。write-prd 当前有 16+ 处 intent 相关引用需要重写，提案只说"将二元分支替换为 Pipeline Configuration 表"——这低估了 write-prd SKILL.md 的改动复杂度。 |

**Attacks:**

2. **Solution Clarity**: Pipeline Configuration 表是核心交付物但未定义 — "intent 控制默认 pipeline 配置（一张表）" — 必须在提案中给出这张表的完整 6 行定义（每行对每个 pipeline 阶段的默认行为），否则实现时 write-prd 和 tech-design 的作者会给出不一致的默认值。

3. **Solution Clarity**: Override Signals 是概念而非规范 — "PRD 内容中的明确信号可以覆盖默认值" — 必须定义：(a) 信号是什么（关键词/标签/结构化字段）？(b) 检测机制是什么？(c) 覆盖动作是什么？否则这个"混合模式"只是把"完全内容驱动 pipeline"的风险从默认路径移到了异常路径。

4. **Solution Clarity**: `enhancement` 的 pipeline 行为未定义 — Key Scenarios 中说 enhancement "pipeline 默认跳过 user stories 但保留 test pipeline" — 但 write-prd 当前只有 `new-feature`（full PRD）和 `refactor`/`cleanup`（spec-only PRD）两种格式。enhancement 走哪种？需要第三种格式吗？提案没有回答。

### 3. Industry Benchmarking: 30/120

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Industry solutions referenced | 5/40 | "分类系统的精确度随场景增长自然演进是常见模式"是一句泛泛的陈述，没有引用任何具体产品、开源项目或已发表的模式。Comparison Table 中"CI lint gate 模式"作为 Source 被提及但未展开——哪个 CI 系统？什么 lint gate 模式？如何映射到 intent pipeline？ |
| At least 3 meaningful alternatives | 15/30 | 4 个选项存在。"只扩枚举"（只解决一半问题）和"完全内容驱动 pipeline"（不稳定）都是 straw man——它们被设计为明显劣于选中方案。真正有意义的替代方案如"保留 3 值但在 downstream 增加子类型标签"或"基于内容特征自动推断 pipeline 而非 intent"未被考虑。"Do nothing"虽然列出但只说"覆盖缺口随场景增长扩大"作为 reject 理由，缺乏具体的 cost of inaction 分析。 |
| Honest trade-off comparison | 5/25 | 选中方案的 Cons 只有"8 个文件变更"——这是工作量而非 trade-off。真正的 trade-off 未被讨论：6 值枚举增加 brainstorm 的认知负担（用户需要理解 6 个选项的区别）、混合模式引入 pipeline 行为的不确定性（同一个 intent 在不同 PRD 内容下可能走不同路径）。 |
| Chosen approach justified against benchmarks | 5/25 | "CI lint gate 模式"作为唯一参考，但映射关系未说明。CI lint gate 是"规则表 + 异常覆盖"，提案的"Pipeline Configuration 表 + Override Signals"确实是类似结构，但为什么这个模式适用于 LLM-driven pipeline 的具体论证缺失。 |

**Attacks:**

5. **Industry Benchmarking**: 零具体引用 — "分类系统的精确度随场景增长自然演进是常见模式——从粗粒度到细粒度" — 这不是 industry benchmarking，这是常识陈述。需要引用至少一个具体系统（如 Sentry 的 issue classification、GitHub 的 label taxonomy evolution、SemVer 的版本语义分类）并分析其演进路径与当前提案的对应关系。

6. **Industry Benchmarking**: Straw-man alternatives — "只扩枚举"的 Cons 是"pipeline 分支仍然过粗" — 但这恰好是很多系统选择的渐进改进路径（先扩分类再优化分支）。作为 straw man 被一棍子打死，没有分析其可行性。更严重的是，"完全内容驱动 pipeline"被否决的理由是"依赖 LLM 判断力，不稳定"，但选中的 Override Signals 同样依赖 LLM 判断——只是范围缩小了。这个 trade-off 分析不自洽。

### 4. Requirements Completeness: 55/110

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Scenario coverage | 25/40 | 5 个 Key Scenarios 覆盖了主要 happy path（fix、enhancement、refactor+override、doc、混合内容）。但缺少 edge cases：(a) 用户在 AskUserQuestion 中选择的 intent 与 brainstorm 推断不一致时怎么办？(b) proposal 没有对应 proposal.md 时（quick 模式直接生成 task）如何处理 intent？(c) 一个 proposal 的 In Scope 同时包含 doc 和 coding 类型的 task 时，主 intent 如何决定？ |
| Non-functional requirements | 15/40 | 向后兼容被提及（"现有的 new-feature、refactor、cleanup 值行为不变"），但只有这一条 NFR。缺少：LLM token 消耗影响（6 个选项 vs 3 个选项在 AskUserQuestion 中的差异）、pipeline 判断的确定性（Override Signals 的可重复性——同一 PRD 内容两次运行是否产生相同 pipeline 路径）、迁移成本（已有 features 目录下的旧 proposal 是否需要更新）。 |
| Constraints & dependencies | 15/30 | "Intent 分支逻辑全部在 skill markdown 中，无 Go 代码依赖"是正确的约束。但"变更限于 plugins/forge/ 目录下的 8 个文件"这个约束可能不准确——freeform review 已指出 write-prd 和 tech-design 的 rules/ 子目录中可能还有其他文件引用 intent。经代码验证，write-prd/rules/knowledge-extraction.md 和 write-prd/rules/ui-functions.md 不含 intent 引用，但 tech-design/rules/design-quality-checks.md 含 7 处 intent 引用（已被列入 Scope）。总体约束描述基本正确但缺少验证方法论。 |

**Attacks:**

7. **Requirements Completeness**: 缺少 override 行为的 edge case — Scenario 3 说"PRD 内容包含 CLI 命令重命名信号 → 覆盖开启 API handbook" — 但如果 PRD 中既有"CLI 命令重命名"也有"纯内部重构"的描述怎么办？Override 的优先级和冲突解决规则未定义。

8. **Requirements Completeness**: `enhancement` 的 PRD 格式未定义 — "brainstorm 推断 enhancement，pipeline 默认跳过 user stories 但保留 test pipeline" — write-prd 当前有 `new-feature` 格式（full PRD）和 `refactor`/`cleanup` 格式（spec-only PRD）。enhancement 走哪种？这是 requirements gap。

### 5. Solution Creativity: 25/100

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Novelty over industry baseline | 10/40 | 提案自己承认"无特别创新"。扩枚举 + 配置表 + 覆盖规则是标准的配置驱动设计，没有任何超越行业基准的创新点。 |
| Cross-domain inspiration | 5/35 | "CI lint gate 模式"是唯一跨域参考，但没有具体展开。没有从其他领域（如编译器 optimization pipeline 的 pass selection、规则引擎的 Rete 算法、路由表的 longest-prefix match）借鉴思路。 |
| Simplicity of insight | 10/25 | "消除 intent 与 type 之间的映射鸿沟"是一个合理的简化洞察，但不是"为什么我没想到"级别的洞察——它只是注意到两个已有的分类体系之间存在 gap。 |

**Attacks:**

9. **Solution Creativity**: 提案明确声明无创新 — "无特别创新。对标 task type 的现有分类体系，消除 intent 与 type 之间的映射鸿沟。" — 这不是扣分原因（不是所有提案都需要创新），但 25 分反映的是：如果这是一个纯工程改进提案，它应该在其他维度（如 Requirements Completeness、Solution Clarity）表现更好来补偿。当前提案在其他维度也有缺陷，说明"无创新"之外还有"不够完整"的问题。

### 6. Feasibility: 75/100

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Technical feasibility | 35/40 | "所有变更是 markdown 编辑"在技术上完全可行。write-prd/SKILL.md 和 tech-design/SKILL.md 确实只包含 markdown 中的条件逻辑。但 underestimated 的是 write-prd/SKILL.md 中 intent 相关逻辑的分散程度——Intent Detection 表、Process Flow、Checklist、Output Documents、Step 7、Step 7A、Step 8、Step 9 都有 intent 分支，改动量比"将二元分支替换为 Pipeline Configuration 表"这句话暗示的要大得多。 |
| Resource & timeline | 25/30 | "2-3 个任务可完成"的估算基本合理，但如果 enhancement 的 pipeline 行为需要定义新的 PRD 格式变体（而非简单选择 existing 格式），实际工作量可能翻倍。 |
| Dependency readiness | 15/30 | "无外部依赖"正确。但 internal dependency 未充分分析：breakdown-tasks 和 quick-tasks 的 Type Assignment 表中 `coding.fix` 被标注为 "do not assign manually"——如果 intent `fix` 要映射到 `coding.fix`，就必须修改 Type Assignment 的语义。这是一个内部一致性依赖，提案未处理。 |

**Attacks:**

10. **Feasibility**: `coding.fix` 的手动分配禁令 — breakdown-tasks/SKILL.md 和 quick-tasks/SKILL.md 的 Type Assignment 表明确说 `coding.fix` "Auto-generated for test failures via forge task add; do not assign manually" — 如果 intent `fix` 映射到 `coding.fix`，就违反了 "do not assign manually" 规则。如果映射到 `coding.feature`，1:1 映射就不成立。这个矛盾在提案中完全未提及。

### 7. Scope Definition: 60/80

| Criterion | Score | Justification |
|-----------|-------|---------------|
| In-scope items are concrete | 25/30 | 8 个文件全部列出，每个文件的具体变更内容（"更新 Step 4.5 intent mapping 表"、"将二元分支替换为 Pipeline Configuration 表"）也给出了。这是提案中做得较好的部分。 |
| Out-of-scope explicitly listed | 20/25 | 4 个 out-of-scope 项清晰明确。但"已有提案的迁移（旧 3 值行为不变）"是一个隐含的向后兼容策略——如果旧 proposal 使用 `intent: refactor`，新的 write-prd 是否需要特殊处理？Out-of-scope 说"旧 3 值行为不变"，但这需要在 Pipeline Configuration 表中显式保证（3 个旧行的行必须保持原行为），提案没有强调这一点。 |
| Scope is bounded | 15/25 | "8 个文件"提供了边界。但 write-prd/SKILL.md 的改动量（替换二元分支为 6 行表 + Override Signals 规则）是否可以在一个任务中完成？如果 enhancement 的 PRD 格式需要新增 Step 7B（介于 full PRD 和 spec-only 之间），scope 可能膨胀。 |

**Attacks:**

11. **Scope Definition**: 文件列表可能不完整 — write-prd/SKILL.md 有 16+ 处 intent 引用分布在整个文件中（Intent Detection 表、Process Flow、Checklist、Output Documents、Step 7 Gate、Step 7A、Step 8、Step 9），但 In Scope 只说"将二元分支替换为 Pipeline Configuration 表 + Override Signals"。Step 7 的 Intent Gate、Step 8 的 Intent Gate、Step 9 的 Manifest 条件——这些都只是"替换 Pipeline Configuration 表"吗？还是需要逐一更新？scope 的边界定义不够精确。

### 8. Risk Assessment: 55/90

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Risks identified | 18/30 | 4 个风险列出，但缺少关键风险：(a) `enhancement` 的 undefined pipeline 行为导致 write-prd 和 tech-design 实现不一致；(b) Override Signals 的 LLM 可靠性问题（提案否决了"完全内容驱动 pipeline"因为不稳定，但 Override Signals 本质上是缩小版的内容驱动）；(c) `coding.fix` 的 "do not assign manually" 约束与 intent `fix` 的映射冲突。 |
| Likelihood + impact rated | 17/30 | "6 值枚举仍不够"的 L/L 评级是诚实的。"混合模式的覆盖规则被 LLM 忽略"的 M/M 也合理。但"write-prd/tech-design 分支重写引入不一致"被评级为 M/M——考虑到这是本提案改动量最大的部分（write-prd 有 16+ 处 intent 引用），Impact 应该是 H。 |
| Mitigations are actionable | 20/30 | "Pipeline Configuration 表统一两处逻辑"是可操作的。"覆盖信号是明确的条件表，不是模糊指令"也是可操作的——但前提是这张条件表被定义。提案中没有定义条件表，所以这个 mitigation 当前是 aspirational 而非 actionable。 |

**Attacks:**

12. **Risk Assessment**: 遗漏 intent-to-type 映射的核心矛盾 — `coding.fix` 在 Type Assignment 表中被标注为 "do not assign manually"，但提案新增 `fix` intent 并要求"严格 1:1 映射" — 这个矛盾是 proposal-level 的 design flaw，应在 Risk 表中显式列出。

13. **Risk Assessment**: 混合模式的稳定性风险被低估 — "覆盖信号是明确的条件表，不是模糊指令；LLM 对结构化规则的遵守度高于 prose 描述" — 但提案中这个"条件表"还不存在。如果条件表最终仍然是 prose 形式的"如果 PRD 提及 CLI 命令变更则启用 API handbook"，LLM 的遵守度并不比当前启发式更高。Mitigation 的有效性取决于条件表的具体形式——这又回到了 Solution Clarity 的缺陷。

### 9. Success Criteria: 60/80

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Criteria are measurable and testable | 22/30 | 大部分 SC 是可验证的（"brainstorm 推断结果为 6 值之一"、"fix 始终推断为 fix"、"write-prd 和 tech-design 使用统一的 Pipeline Configuration 表"）。但"Override Signals 规则存在且可被 PRD 内容触发"中"可被触发"的可测试性取决于触发条件的定义——当前未定义。 |
| Coverage is complete | 18/25 | 7 条 SC 覆盖了 brainstorm、write-prd、tech-design、breakdown-tasks、quick-tasks 和向后兼容。但缺少：(a) enhancement 的 pipeline 行为 SC；(b) Override Signals 的具体触发条件 SC；(c) write-prd 的 Step 7 Intent Gate 和 Step 8 Intent Gate 更新后的行为验证 SC。 |
| SC internal consistency | 20/25 | SC 集合内部无直接矛盾。但"严格 1:1 映射"（breakdown-tasks 和 quick-tasks）与 `coding.fix` 的 "do not assign manually" 存在隐含冲突——如果 `fix` intent 不能映射到 `coding.fix`（因为 do not assign manually），那"严格 1:1"就不成立。这不是 SC-to-SC 矛盾，而是 SC-to-existing-rule 矛盾。 |

**Attacks:**

14. **Success Criteria**: "Override Signals 规则存在且可被 PRD 内容触发"不可测试 — "可被 PRD 内容触发"缺少触发条件的定义。什么内容的 PRD？包含什么关键词或结构？测试方法是"给一个包含 CLI 命令重命名的 PRD 看是否触发 API handbook"吗？如果是，需要在 SC 中明确测试用例。

15. **Success Criteria**: 缺少 enhancement 的 SC — In Scope 列出了 8 个文件需要更新，但 SC 中没有针对 enhancement 行为的验证。例如"enhancement intent 的 PRD 跳过 user stories 但保留 test pipeline"应该是 SC。

### 10. Logical Consistency: 95/90 (capped to dimension max)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Solution addresses the stated problem | 30/35 | 扩充枚举确实解决了"3 值 intent 无法覆盖 task type"的问题。混合模式确实解决了"pipeline 分支过粗"的问题。简化推断确实解决了"fix 启发式不稳定"的问题。三条路径都与问题陈述对应。但有一条松动：问题说"不同场景走了不合适的流程"，但"不合适"的判定标准从未定义。 |
| Scope <-> Solution <-> Success Criteria aligned | 30/30 | 8 个文件的 Scope、3 条路径的 Solution、7 条 SC 之间无矛盾。每个 In Scope 文件都有对应的 SC 验证。Out-of-scope 项（Go 代码、新 skill、迁移）在 SC 中无遗漏验证需求。 |
| Requirements <-> Solution coherent | 25/25 | 5 个 Key Scenarios 与 3 条 Solution 路径一一对应。NFR（向后兼容、1:1 映射）在 Solution 和 SC 中都有体现。约束（markdown-only、8 文件）与 Solution 一致。 |

**Note**: Raw score exceeds 90 (dimension max). Capped to 90.

**Attacks:**

16. **Logical Consistency**: "1:1 映射"声称与实际不一致 — "与 task type 形成干净的 1:1 映射" — 6 值 intent 无法覆盖 8 值 task type（`coding.fix` 的手动分配禁令、`doc.consolidate` 和 `doc.drift` 无对应 intent）。Risk 表承认了 `doc.consolidate`/`doc.drift` 的 gap，但 Proposed Solution 中的"干净的 1:1"措辞与 Risk 表的承认矛盾。应改为"覆盖所有主要 task type，近似 1:1 映射"。

17. **Logical Consistency**: 混合模式否决了"完全内容驱动"但自己也是内容驱动 — "完全内容驱动 pipeline...依赖 LLM 判断力，不稳定 → Rejected" vs "PRD 内容中的明确信号可以覆盖默认值" — 两者的区别仅在于范围（全部 vs 例外），但提案没有论证范围缩小是否足以补偿稳定性风险。如果 Override Signals 的条件表最终也是 LLM-interpreted prose，那混合模式就是带默认值的"完全内容驱动"——而这个选项在 Alternatives 中被否决了。

---

## Summary of Critical Gaps

1. **Pipeline Configuration 表未定义** — 这是提案的核心交付物，但 6 行表中 enhancement 和 fix 的行为完全空白
2. **Override Signals 是概念而非规范** — "PRD 内容中的明确信号"没有定义检测机制，依赖 LLM 判断与提案否决的"完全内容驱动"存在逻辑矛盾
3. **`enhancement` 的完整行为未定义** — PRD 格式、tech-design 决策类型、self-check 规则全部空白
4. **`coding.fix` 的 "do not assign manually" 约束与 1:1 映射矛盾** — 提案未识别这个冲突
5. **Industry benchmarking 几乎为零** — 只有一句泛泛的"常见模式"，没有引用任何具体系统或模式
