# Eval Report: Task Pipeline Precision Tuning

**Iteration**: 1
**Scorer**: CTO Adversarial
**Date**: 2026-05-27
**Document**: `docs/proposals/task-pipeline-precision/proposal.md`

---

## Phase 1: Reasoning Audit

### Argument Chain Trace

1. **Problem -> Solution**: 提案诊断了 3 类根因（任务粒度过粗、模板忽略复杂度差异、Reference Files 指向过宽），对应的 3 环节精度控制（拆分规则优化、complexity 元数据字段、模板分支+搜索策略引导）形成一一映射。论证链完整，solution 直接 address problem。

2. **Solution -> Evidence**: 4 个 lesson 表格提供了具体数据（18.3 min、25 min、11 条 AC 中 2 条未完成、越界修改 task 4），与 solution 各环节对应关系清晰。gotcha-task-executor-thinking-overhead -> 拆分规则+搜索策略，gotcha-prompt-template-complexity-agnostic -> complexity 分支，gotcha-quick-tasks-merge-threshold -> 拆分标准，gotcha-task-reference-files-scope-creep -> 内联+scope boundary。

3. **Evidence -> Success Criteria**: SC 覆盖了 complexity 判定（SC-1）、模板分支（SC-2、SC-6）、Reference Files 格式（SC-3）、AC 上限（SC-4）、搜索策略引导（SC-5）、向后兼容（SC-7）。Evidence 中的 4 个故障场景均有对应 SC 验证。SC 集合在内部逻辑上可满足、无矛盾。

4. **Self-contradiction check**: 经过 pre-revision 修改后，已移除"15 上限"scope creep，fix 模板已排除在 complexity routing 之外，搜索策略引导的具体位置和行为已定义。未发现 reintroduction。

### SC Consistency Deep-Dive

Cluster by affected area:

**Cluster A: Task Generation (quick-tasks/breakdown-tasks SKILL.md)**
- SC-1 (complexity: low 标记) <-> InScope: "复杂度自动判定逻辑" -> 可满足，静态指标+LLM 覆盖机制明确
- SC-4 (AC <= 6) <-> InScope: "AC 上限 6 条" -> 可满足，拆分规则直接约束
- SC-3 (Reference Files 内联) <-> InScope: "Reference Files 生成策略改为内联精确信息" -> 可满足
- SC-1 <-> SC-4: 无冲突，complexity 判定和 AC 上限是独立维度

**Cluster B: Prompt Templates (coding-*.md)**
- SC-2 (low 跳过 Step 1.5) <-> SC-6 (prompt 输出包含 complexity 分支) -> 可满足且一致
- SC-5 (搜索策略引导出现在 4 个非 fix template) <-> InScope: "4 个 coding prompt templates...加搜索策略引导" -> 可满足
- SC-5 <-> SC-2: 无冲突，搜索策略引导在 Step 1 之后 Step 2 之前，Step 1.5 跳过不影响搜索策略引导位置

**Cluster C: Data Pipeline (prompt.go, types)**
- SC-7 (无 complexity 字段默认 medium) <-> InScope: "prompt.go 传递 complexity 字段" -> 可满足
- SC-6 <-> SC-7: 无冲突，default medium 确保 backward compat

**Cross-cluster**: SC-2 要求 complexity: low 跳过 Step 1.5，依赖 Cluster A 的 SC-1 正确标记 low -> 可满足。SC-3 要求内联格式，依赖 Cluster A 的生成逻辑 -> 可满足。

结论：SC 集合内部无矛盾，与 In Scope 双向可满足。

---

## Phase 2: Rubric Scoring

### D1. Problem Definition (110 pts)

**Problem stated clearly (35/40)**: 核心问题"任务执行在简单任务上浪费大量时间"表述清晰。"87% 时间花在冗余探索"和"25 分钟做一个纯文本替换"是具体的量化描述。3 个根因（粒度过粗、模板忽略复杂度、Reference Files 越界）均无歧义。微扣 5 分：开篇"87% 时间花在冗余探索上，25 分钟做一个纯文本替换"中 87% 是一个 case 的数据点（gotcha-task-executor-thinking-overhead），不是全局统计，但写法暗示系统性。精确表述应为"单个案例中 87% 思考时间"。

**Evidence provided (38/40)**: 4 个 lesson 表格提供了具体的故障现象、根因映射，且注明了"4 个独立 but 相互关联"。数据具体（18.3 min、25 min、11 条 AC）。扣 2 分：lesson 表格中的数据依赖于 lesson 文档的存在，但提案未直接引用 lesson 文件路径，reader 无法自行追溯原始数据。

**Urgency justified (28/30)**: "每次 quick mode 执行都受这些问题影响" + "4 个 lesson 在同一次 feature 执行中发现，说明问题是系统性的"提供了为什么现在解决的理由。扣 2 分：缺少"不解决的持续成本"量化——如果每周执行 N 次 quick mode，每次浪费 M 分钟，累计成本是多少？

**D1 Total: 101/110**

---

### D2. Solution Clarity (120 pts)

**Approach is concrete (36/40)**: 3 个环节各给出具体方案——拆分规则改为"可独立验证"、complexity 字段 low/medium/high、模板分支+搜索策略引导。reader 可以复述将建什么。扣 4 分：搜索策略引导的具体指令内容在正文中给出了一段示例（"在修改任何文件前，先用 Grep/Glob 搜索所有需要修改的位置..."），但这段内容是"具体内容为..."的附带说明，而非正式的模板定义，位置和格式的精确规范仍有模糊空间。

**User-facing behavior described (42/45)**: 用户（开发者）体验变化清晰：quick-tasks 生成的任务会有 complexity 标签、低复杂度任务执行更快（跳过 spec-code scan、简化探索）、task 文件是自包含的。扣 3 分：缺少一个 end-to-end 场景走读——用户执行 `/quick` 后每个阶段（生成、显示、执行）的可观测变化是什么？

**Technical direction clear (33/35)**: Constraints 部分详细描述了 renderTemplate() 的 strings.ReplaceAll 机制、cleanTemplateOutput() 条件段落方案、embed.FS 加载方式。实现路径足够具体。扣 2 分：task.md frontmatter 的 complexity 字段 schema（值域、是否必须、默认值）应在 Requirements 或 Constraints 中明确声明，当前仅在 In Scope 中隐含。

**D2 Total: 111/120**

---

### D3. Industry Benchmarking (120 pts)

**Industry solutions referenced (25/40)**: 提案提到 Cursor, Copilot, Aider 对任务不区分复杂度，但仅用一句"统一 prompt 或完全依赖模型自行判断"概括。未引用任何具体产品文档、技术博客或开源仓库来支撑这一观察。扣 15 分：缺少可追溯的行业参考。"大多数 AI coding agent"是一个广泛断言，应附带至少 2-3 个具体来源。另外，类比领域（如 IDE 的 quick-fix vs refactor 区分、编译器的 -O1/-O2/-O3 优化级别）可作为 cross-domain 参考，但未提及。

**At least 3 meaningful alternatives (28/30)**: "Do nothing"、"Prompt 内置启发式规则"、"显式 complexity 字段 + 内联 Reference"三个替代方案各自代表不同策略层次（无为、运行时判定、生成时标注）。"Do nothing"不是 straw man，因有 lesson 证据支持 rejection。扣 2 分：缺少一个"混合方案"替代（如同时支持静态指标默认+人工 override），这是 industry 中常见的做法。

**Honest trade-off comparison (22/25)**: "显式 complexity 字段"的 Cons 是"需改 task 模板和生成逻辑"，诚实且合理。"Do nothing"的 Cons 是"executor 效率持续低下"。扣 3 分：comparison table 缺少量化对比——各方案预估的 token 节省/执行时间改善是多少？即使是粗略估计也比纯定性描述好。

**Chosen approach justified against benchmarks (20/25)**: "判定一次，执行时零成本"是明确的 selection rationale，且与 Rejected 方案"每次执行都做判定"形成对比。扣 5 分：选择理由仅基于效率，未讨论 accuracy trade-off——静态指标+LLM 覆盖方案的判定准确率是否真的比模型运行时自行判断更高？提案假设了"是"但未论证。

**D3 Total: 95/120**

---

### D4. Requirements Completeness (110 pts)

**Scenario coverage (35/40)**: 5 个关键场景覆盖了 low/medium/high 三种复杂度，加上多动词拆分和 Reference Files 生成。Happy path 完整。扣 5 分：缺少 error/edge case 场景——如果 complexity 判定逻辑遇到"3 AC + 1 Hard Rule + 0 Reference Files"这种边界情况怎么办？如果内联 Reference Files 的源 proposal section 被 rename/删除了呢？如果 cleanTemplateOutput() 的条件标记在某个 template 中格式错误呢？

**Non-functional requirements (36/40)**: 向后兼容、零运行时开销、探索效率（< 30s 软目标）、模板一致性 4 个 NFR 均具体。扣 4 分：缺少一个 NFR——complexity 字段的可审计性/可调试性。当开发者认为 complexity 判定错误时，如何追溯判定依据？这是一个运维场景，对 acceptance 有关键影响。

**Constraints & dependencies (26/30)**: 列出了 embed.FS、strings.ReplaceAll、命名约定、cleanTemplateOutput() 条件段落方案等 4 个技术约束。扣 4 分：未说明 SKILL.md 修改对 breakdown-tasks 的 PRD Coverage Verification 和 Phase & Gate Detection 环节的影响——这两个环节在 breakdown-tasks 中存在但 quick-tasks 中不存在，"同步修改"是否需要适配这些差异？

**D4 Total: 97/110**

---

### D5. Solution Creativity (100 pts)

**Novelty over industry baseline (32/40)**: "Complexity-aware prompt routing"在 AI coding agent 领域确实不是标准做法，有创新性。区分点在于"生成时标注 vs 运行时判断"的设计选择。扣 8 分：创新程度适中而非突破性——complexity 分级是软件工程中的常见模式（编译器优化级别、测试 pyramid 分层、CI pipeline 的 fast/full 模式），将其应用于 AI prompt routing 是合理的跨领域迁移，但不是"why didn't I think of that"级别的洞察。

**Cross-domain inspiration (22/35)**: Scope boundary declaration 借鉴了最小权限原则（安全领域），Self-contained task documents 借鉴了微服务/模块化设计理念。但提案未显式引用这些领域知识。扣 13 分：未识别任何具体的跨领域灵感来源。编译器的优化级别选择、IDE 的 quick-fix vs full refactor 区分、数据库的 query plan optimization——这些都是高度相关的类比，提案未触及。

**Simplicity of insight (22/25)**: "任务粒度应由可独立验证性而非时间估算决定"和"复杂度决定探索深度而非任务类型"是两个简洁有力的洞察。Scope boundary declaration 作为主动防御机制也较优雅。扣 3 分：cleanTemplateOutput() 的标记注释方案（`<!-- IF NOT_LOW -->...<!-- END_IF -->`）引入了一个小型 DSL，虽然比引入模板引擎轻量，但仍然是额外的复杂度层。提案声称"更轻量"但未量化比 Go `text/template` 简单多少。

**D5 Total: 76/100**

---

### D6. Feasibility (100 pts)

**Technical feasibility (36/40)**: Constraints 部分准确识别了 4 处数据管道改动链（FrontmatterData -> Task struct -> index.json schema -> renderTemplate()）。cleanTemplateOutput() 条件段落方案在现有架构内可行。扣 4 分：10 个文件改动的估算似乎只覆盖了配置/模板层，未计入 prompt.go、types.go、frontmatter.go、cleanTemplateOutput() 等 Go 代码文件——这些才是改动链的核心。如果计入代码文件，实际改动面可能为 14-16 个文件。

**Resource & timeline feasibility (27/30)**: "6-10 个 coding task，适合 quick mode"给出了具体规模估计。扣 3 分：如果实际改动面为 14-16 文件（含 Go 代码），6-10 coding task 可能偏紧，尤其是 data pipeline 的 4 处同步修改需要高精度协调。未提供 fallback plan（如果某个环节改动比预期复杂怎么办）。

**Dependency readiness (28/30)**: "无外部依赖，所有修改均在 plugins/forge/ 和 forge-cli/ 内部"清晰且可验证。扣 2 分：未确认 quick-tasks 和 breakdown-tasks 的 SKILL.md 是否有共同依赖（如共享的 prompt fragment 或 config），如果有，同步修改时需注意修改顺序。

**D6 Total: 91/100**

---

### D7. Scope Definition (80 pts)

**In-scope items are concrete (27/30)**: 8 个 In Scope 条目中大多数是可交付的——"quick-tasks SKILL.md 拆分规则优化"可验证、"task frontmatter 加 complexity 字段"可验证。扣 3 分："prompt.go 传递 complexity 字段到模板渲染"是一个技术实现描述而非用户可交付物，更适合放在 Implementation Notes 中。但作为 quick mode 的 task 描述，勉强可接受。

**Out-of-scope explicitly listed (23/25)**: 6 个 Out of Scope 条目明确且合理——Template 体系重组、质量门基线测试、task-executor agent 定义、现有 proposal 合并、/quick 命令上限、移除 15 coding task 上限。扣 2 分："现有 proposal 合并或替代（slim-task-prompt-templates, prompt-template-audit）"提到了两个具体 proposal 名称，但未说明它们与当前 proposal 的关系——是 predecessor？competitor？superset？

**Scope is bounded (22/25)**: "6-10 个 coding task"和"10 个文件改动"提供了量化边界。Out of Scope 列表防止了 scope creep。扣 3 分：缺少时间边界估计——"适合 quick mode"暗示可在一次 quick session 完成，但未明确估计小时数或 session 数。对于 CTO 审批，需要知道人力投入。

**D7 Total: 72/80**

---

### D8. Risk Assessment (90 pts)

**Risks identified (26/30)**: 5 个风险覆盖了判定准确性、文件过长、stale reference、模板维护成本、同步不一致。扣 4 分：缺少 2 个重要风险——(1) complexity 字段可能成为"永远 medium"的无效分类（如果启发式太保守，所有任务都标 medium，则整个方案沦为 no-op）；(2) 条件标记机制（`<!-- IF NOT_LOW -->`）本身可能引入 bug——如果标记格式错误导致 cleanTemplateOutput() 删除了不该删除的段落，后果比当前"多做一些探索"更严重。

**Likelihood + impact rated (26/30)**: Likelihood 和 Impact 评级基本诚实——"内联 Reference Files 使 task 文件过长"的 Likelihood: L、Impact: L 合理（因为有限制条目数 <=5 的缓解措施）。扣 4 分："prompt template 复杂度分支增加模板维护成本"标记为 Impact: L 偏乐观——4 个 template 各加条件块，每次修改 template 都需要考虑 3 个 complexity 分支的行为，维护成本实际为 M。

**Mitigations are actionable (24/30)**: "使用保守启发式"、"限制条目数 <=5"、"溯源标注 + 以代码为准"、"分支逻辑统一模板化"、"两个 SKILL.md 使用相同规则描述"。大部分可操作。扣 6 分：(1) "两个 SKILL.md 使用相同的判定规则描述，确保逻辑一致"不是一个 mitigation action，而是一个目标——如何确保？是共享一个 markdown fragment？还是 code review checklist？(2) "分支逻辑统一模板化（一段条件块，4 个 template 复制相同结构）"的"复制相同结构"本身引入维护风险（copy-paste drift），mitigation 应包含如何避免 drift（如共享 snippet 或 linter）。

**D8 Total: 76/90**

---

### D9. Success Criteria (80 pts)

**Criteria are measurable and testable (26/30)**: 7 个 SC 中大多数可通过 `forge prompt get-by-task-id` 输出、index.json 内容、template 文件内容客观验证。扣 4 分：SC-1 "AC <= 3 且无 Hard Rules 且 Reference Files <= 1 的任务被标记为 complexity: low"——如果 LLM 覆盖了静态指标的判定（In Scope 允许这样做），这个 SC 可能不满足（AC<=3 但 LLM 标为 medium），但方案仍然正确。SC 的验证条件与实际行为之间存在 gap。

**Coverage is complete (22/25)**: SC 覆盖了 task generation、template routing、Reference Files 格式、向后兼容。扣 3 分：缺少一个 SC 验证 scope boundary declaration——In Scope 中提到了"每个 task 文件头部嵌入显式 scope 边界声明"，但没有 SC 检查 task 文件确实包含此声明。这是一个 In Scope 条目无对应 SC 的 gap。

**SC internal consistency (23/25)**: 如 Phase 1 SC Consistency Deep-Dive 分析，7 个 SC 在内部逻辑上一致、可同时满足、无矛盾。扣 2 分：SC-6（"`forge prompt get-by-task-id` 输出包含 complexity 对应的流程分支内容"）与 SC-7（"现有无 complexity 字段的 index.json 任务执行时默认为 medium，行为不变"）之间的交互未显式验证——如果一个 old task（无 complexity 字段）通过 `get-by-task-id` 获取 prompt，输出是否包含 Step 1.5？SC-7 说"行为不变"暗示会包含 Step 1.5，SC-6 说"包含 complexity 对应的流程分支内容"——对于 medium 默认值，Step 1.5 应该包含。逻辑可推断但未显式声明。

**D9 Total: 71/80**

---

### D10. Logical Consistency (90 pts)

**Solution addresses the stated problem (32/35)**: 3 个根因 -> 3 环节精度控制，映射完整。搜索策略引导直接 address "重复 grep 同一模式"问题。Scope boundary declaration 直接 address "越界修改"问题。扣 3 分：gotcha-prompt-template-complexity-agnostic 描述的根因是"coding-enhancement.md 对所有 enhancement 任务强制完整 Step 1 + Step 1.5"，但 solution 中 Step 1 的简化仅在 NFR 中以"简化探索"模糊提及，未在 In Scope 或 SC 中有对应条目。Step 1 简化具体指什么？简化到什么程度？这是 problem->solution 映射的一个 gap。

**Scope <-> Solution <-> Success Criteria aligned (25/30)**: In Scope 的 8 个条目中，7 个有对应 SC。扣 5 分：(1) In Scope "scope boundary declaration"（每个 task 文件头部嵌入显式 scope 边界声明）无对应 SC；(2) In Scope "搜索策略引导"有 SC-5 验证其出现在 template 中，但 SC-5 只检查"出现"不检查"有效"——如果搜索策略引导存在但 executor 仍然边搜边改，SC 不会检测到。

**Requirements <-> Solution coherent (22/25)**: Requirements 的 5 个场景均有对应 solution 元素。扣 3 分：NFR "模板一致性：5 个 coding template 的复杂度分支逻辑保持统一结构"——但 In Scope 明确说"4 个 coding prompt templates...coding.fix 模板不纳入"。NFR 说"5 个"而 In Scope 说"4 个"，数字不一致。

**D10 Total: 79/90**

---

## Phase 3: Blindspot Hunt

### [blindspot] Complexity 字段默认值的渐进式失效风险

提案假设"默认 medium"是安全的 backward compat 策略。但这个默认值隐含了一个假设：existing tasks 的行为是正确的，不需要优化。如果 existing tasks 中有本应是 low 的任务，它们将永远按 medium 执行，无法获得效率提升。这不是 bug，但意味着方案对存量任务的效率改善为零——仅对新任务有效。提案未讨论这个 rollout 策略。

### [blindspot] breakdown-tasks 的 Source Document 差异被忽略

quick-tasks 从 proposal 生成任务，breakdown-tasks 从 tech-design 生成任务。两者 source document 的结构、粒度、信息密度完全不同。In Scope 第 2 条要求"breakdown-tasks SKILL.md 同步修改拆分规则"，但 complexity 判定逻辑（AC 数量 + Hard Rules + Reference Files 数量）是基于 proposal 生成的经验得出的——从 tech-design 生成的任务是否有相同的复杂度分布？breakdown-tasks 生成的是一个完整 pipeline 中的任务，其 AC 粒度可能天然更细（因为 tech-design 已经做了设计分解），导致 complexity 判定阈值可能需要差异化。

### [blindspot] SC-1 判定条件与 In Scope 判定逻辑的语义 gap

SC-1 写的是"AC <= 3 且无 Hard Rules 且 Reference Files <= 1 的任务被标记为 complexity: low"，但 In Scope 第 4 条写的是"以静态指标作为默认启发式，同时提供 LLM 判断指引，允许 LLM 在静态指标与认知判断冲突时覆盖默认值"。这意味着 SC-1 的条件是 sufficient 但不是 necessary——一个 AC=5 但 LLM 判断为 low 的任务也应标为 low。SC-1 的措辞暗示"只有满足这三个条件才标 low"，与 In Scope 的"LLM 可覆盖"矛盾。

---

## Bias Detection Report

Annotated regions (pre-revised markers):
- Total annotated paragraphs: 7
- Attack points on annotated regions: 3 (D8 条件标记 bug 风险、D9 scope boundary SC gap、D7 prompt.go 描述)
- Attack density: 3/7 = 0.43

Unannotated regions:
- Total unannotated paragraphs: ~38
- Attack points on unannotated regions: 14
- Attack density: 14/38 = 0.37

Ratio (annotated/unannotated): 1.16

Interpretation: Annotated regions have slightly higher attack density (1.16x), close to parity. No significant bias detected — the pre-revision improved most marked regions, but a few residual gaps remain.

---

## Score Summary

| Dimension | Score | Max |
|-----------|-------|-----|
| D1. Problem Definition | 101 | 110 |
| D2. Solution Clarity | 111 | 120 |
| D3. Industry Benchmarking | 95 | 120 |
| D4. Requirements Completeness | 97 | 110 |
| D5. Solution Creativity | 76 | 100 |
| D6. Feasibility | 91 | 100 |
| D7. Scope Definition | 72 | 80 |
| D8. Risk Assessment | 76 | 90 |
| D9. Success Criteria | 71 | 80 |
| D10. Logical Consistency | 79 | 90 |
| **Total** | **869** | **1000** |

---

## ATTACK_POINTS

1. **[D3: Industry Benchmarking]** 行业参考缺乏可追溯来源 — "大多数 AI coding agent（Cursor, Copilot, Aider）对任务不区分复杂度——统一 prompt 或完全依赖模型自行判断" — 需要引用具体产品文档、技术博客或开源仓库来支撑这一断言，至少 2-3 个可追溯来源

2. **[D3: Industry Benchmarking]** Trade-off 对比缺少量化 — "Pros: 精确、可审计、零运行时开销" — 需要粗略估计各方案的 token 节省/执行时间改善，纯定性对比不足以支撑 CTO 决策

3. **[D3: Industry Benchmarking]** 未论证判定准确率 — "判定一次，执行时零成本"作为 selection rationale 仅讨论效率 — 需要论证"生成时静态标注+LLM 覆盖"的判定准确率是否优于"运行时模型自行判断"

4. **[D5: Solution Creativity]** 未识别跨领域灵感来源 — 整个 Innovation Highlights 部分未引用任何跨领域类比 — 应识别并引用编译器优化级别、IDE quick-fix vs refactor、CI fast/full pipeline 等相关领域知识

5. **[D6: Feasibility]** 文件改动数量估算不准确 — "~10 个文件改动，均为配置/模板层" — 未计入 prompt.go、types.go、frontmatter.go、cleanTemplateOutput() 等 Go 代码文件，实际改动面可能为 14-16 文件，"均为配置/模板层"描述不准确

6. **[D8: Risk Assessment]** 缺少 complexity 字段"永远 medium"风险 — 当前风险表仅讨论"简单任务被标为 high"的反方向 — 如果启发式太保守，所有任务都标 medium，整个方案沦为 no-op，此风险应纳入评估

7. **[D8: Risk Assessment]** 缺少条件标记机制的格式错误风险 — "用标记注释包裹 Step 1.5 段落，由 cleanTemplateOutput() 根据 complexity 值删除标记块" — 标记格式错误可能导致 cleanTemplateOutput() 删除不该删除的段落，后果比当前"多做一些探索"更严重

8. **[D8: Risk Assessment]** Mitigation "两个 SKILL.md 使用相同的判定规则描述"不是可操作行动 — 应说明如何确保一致（共享 markdown fragment？code review checklist？linter？）

9. **[D9: Success Criteria]** SC-1 判定条件与 In Scope 判定逻辑矛盾 — SC-1 写"AC <= 3 且无 Hard Rules 且 Reference Files <= 1 的任务被标记为 complexity: low"，In Scope 允许 LLM 覆盖——SC-1 应改为"满足...条件的任务**默认**标为 low，LLM 可根据认知判断覆盖"

10. **[D9: Success Criteria]** In Scope "scope boundary declaration"无对应 SC — "每个 task 文件头部嵌入显式 scope 边界声明" — 需要一个 SC 验证 task 文件确实包含 `## Scope Boundary` 段落

11. **[D10: Logical Consistency]** NFR "5 个 coding template"与 In Scope "4 个...coding.fix 不纳入"数字不一致 — NFR 写"5 个 coding template 的复杂度分支逻辑保持统一结构"，In Scope 明确 fix 模板不参与 — 应统一为 4 个

12. **[D10: Logical Consistency]** Step 1 简化缺少对应 In Scope 条目和 SC — Problem 描述了"coding-enhancement.md 强制完整 Step 1 + Step 1.5"，solution 提到"low 简化 Step 1"，但 In Scope 和 SC 中无 Step 1 简化的具体描述和验证条件

13. **[blindspot]** breakdown-tasks 的 complexity 判定阈值可能需要差异化 — "以静态指标（AC 数量 + Hard Rules + Reference Files 数量）作为默认启发式"基于 proposal 生成经验 — 从 tech-design 生成的任务 AC 粒度天然更细，相同阈值可能不适配

14. **[blindspot]** 存量任务的 efficiency 改善为零 — "现有 index.json 中无 complexity 字段的任务默认为 medium" — 未讨论存量任务的 rollout 策略，默认 medium 意味着所有已有任务无法获得 low routing 的效率提升

15. **[D4: Requirements Completeness]** 缺少 complexity 判定边界情况场景 — "3 AC + 1 Hard Rule + 0 Reference Files"等边界条件未在 Key Scenarios 中覆盖 — 应增加 boundary case 的判定规则说明
