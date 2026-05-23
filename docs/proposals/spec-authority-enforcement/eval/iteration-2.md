# Proposal Evaluation — Iteration 2

## Iteration 1 Issue Tracking

| # | Iteration 1 Attack | Status | How Addressed |
|---|---------------------|--------|---------------|
| 1 | Step 5 前时序矛盾：agent 未读取任务文件时无法声明 Reference Files | **Resolved** | 提案修订为"在 coding.* 模板的 Step 1（读取任务文件之后）插入声明"，而非"在 task-executor.md 的 Step 5 前插入"。成功标准第 1 条也改为"Step 1（读取任务文件之后）"。时序逻辑正确。 |
| 2 | quick-tasks 无 tech-design.md 时无法提取精确 section 引用 | **Resolved** | 新增场景 4 专门描述此情况："从 proposal.md 中提取关键技术约束和决策，将 proposal.md 中与当前实现直接相关的 section 写入 Reference Files"。成功标准第 5 条相应修改为"包含 ≥1 个精确 section 引用（非仅 proposal.md）"——但此条对 quick-tasks 场景 4 仍有歧义（见下文攻击）。 |
| 3 | embed.FS 编译依赖未声明，声称"纯 Markdown 修改不涉及 Go 代码" | **Resolved** | Constraints 新增"coring.* 模板通过 embed.FS 嵌入二进制，修改后必须执行 go build"。Feasibility Assessment 和 Success Criteria 也包含了部署一致性验证。 |
| 4 | 审计标准推迟到实施阶段 | **Resolved** | Scope In Scope 新增审计标准：(a) 模板用于 coding 或 doc 任务类型；(b) Step 1 包含"读取任务文件"步骤；(c) 涉及需要对照规范执行的实现/修改任务。 |
| 5 | EXTREMELY-IMPORTANT 标记稀释 | **Resolved** | 提案改为使用 `<IMPORTANT>` 标记（而非 `EXTREMELY-IMPORTANT`），并在多处说明"不使用 EXTREMELY-IMPORTANT 以避免与模板中已有的 EXTREMELY-IMPORTANT 块产生标记稀释"。 |
| 6 | Reference Files 引用过期后无降级行为 | **Resolved** | 新增 Edge Cases & Degradation section，定义了四种降级场景：Reference Files 为空、section 标题不存在、文件路径不存在、所有引用失效。 |
| 7 | breakdown-tasks 缺少 Reference Files 填充具体指引 | **Partially Resolved** | 提案要求"为非 UI 任务增加 Reference Files 填充指引，要求精确到 section"，但没有给出具体的操作步骤（如何从 tech-design.md 提取相关 section 的决策逻辑）。 |
| 8 | AC 验收对不同任务类型适配问题 | **Partially Resolved** | 新增场景 3（无 AC 任务），但只讨论了"AC 为空"的退化情况。对于 doc.* 任务（场景 2）的 AC 验收差异（"验收重点是文档结构合规而非路径命名"）仍是一句话概述，缺少具体定义。 |
| 9 | User-facing behavior 缺失 | **Resolved** | 新增 User-Visible Behavior Changes section，描述了三种可观察的行为变化：加载确认、AC 验收报告、降级提示。 |
| 10 | Industry benchmarking 缺少具体引用 | **Resolved** | 新增具体引用：LangChain RetrievalQA、LlamaIndex QueryEngine、MetaGPT SOP、CrewAI Knowledge、OpenAI Function Calling、Anthropic Tool Use、Claude Prompt Engineering 指南。 |
| 11 | 缺少回滚计划 | **Resolved** | 新增 Rollback Plan section，包含即时回滚、部分回滚、验证回滚三层策略。 |
| 12 | Reference Files 与 Hard Rules 优先级冲突 | **Resolved** | 新增 Priority Rules section，定义了 Hard Rules > Reference Files > 现有代码的三级优先级，并定义了冲突时的行为。 |
| 13 | 部署不一致风险未列入 Risk 表 | **Resolved** | Risk 表新增第 4 行："task-executor.md 与 coding.* 模板更新不同步导致两层防护不一致"，Likelihood M，Impact H。 |
| 14 | 成功标准衡量产出而非结果 | **Partially Resolved** | 新增第 9 条成功标准："行为验证：使用修改后的模板执行 1 个 coding 任务后，agent 输出中包含 Reference Files 加载确认和 AC 验收报告"。但仍缺少量化度量（偏差率降低）。 |
| 15 | 替代方案缺少"仅修改 Hard Constraints"的轻量方案 | **Resolved** | Comparison Table 新增"仅修改 task-executor.md Hard Constraints"行，标注为"Considered: 作为兜底规则采纳"。 |

## Freeform Findings Tracking

| # | Finding | Status | Evidence in Revision |
|---|---------|--------|----------------------|
| F1 | [high] 增加 EXTREMELY-IMPORTANT 标记会产生标记稀释效应 | **Resolved** | 全文改用 `<IMPORTANT>`，且在 Scope 和 Success Criteria 中明确说明原因 |
| F2 | [high] section 引用格式未与现有 Reference Files 兼容 | **Unresolved** | 提案定义了新格式 `path/to/file.md#section — 简要说明`，但未讨论旧格式 `path/to/file.md` 的共存问题。已有任务文件使用旧格式时，agent 如何处理？ |
| F3 | [high] Reference Files 声明在 agent 读取任务文件之前无法生效 | **Resolved** | 修正为"在 coding.* 模板的 Step 1（读取任务文件之后）" |
| F4 | [high] Reference Files 引用过期后无降级行为定义 | **Resolved** | Edge Cases & Degradation section 定义了完整降级链 |
| F5 | [high] 修改 embed.FS 嵌入模板后不重编译会导致两层防护不一致 | **Resolved** | Constraints、Feasibility、Risk、Success Criteria 四处提及 |
| F6 | [medium] quick-tasks 无 tech-design.md 时无法提取精确 section 引用 | **Resolved** | 场景 4 专门处理 |
| F7 | [medium] breakdown-tasks 缺少 Reference Files 填充的具体操作指引 | **Partially Resolved** | Scope 列出"为非 UI 任务增加 Reference Files 填充指引"，但未给出提取 section 的决策逻辑 |
| F8 | [medium] AC 验收步骤未讨论对不同任务类型的适配问题 | **Partially Resolved** | 场景 3 覆盖"无 AC"情况，场景 2 提及 doc 任务差异但未展开 |

---

## Dimension Scores

### 1. Problem Definition: 92/110

**Problem stated clearly (37/40):**
核心问题清晰且不预设解决方案——"Agent 执行任务时以现有代码为参照而非以权威规范文档为参照，导致大规模偏离"。引用了 43 处偏差的具体事件。相比 iteration 1，标题中的"从模板和任务生成层面防止"仍然保留了方案层面的暗示，但整体问题陈述已足够准确。扣 3 分。

**Evidence provided (35/40):**
引用了教训文档的 5 级溯源，Level 0 根因精确到"模板 Step 1 未要求读 Reference Files"。Level 3 的"反应式局部修复"是可验证的过程描述。探索发现的具体事实（quick-tasks 硬编码 proposal.md、breakdown-tasks 无填充指引）是强证据。但仍然缺少"如果当时 agent 读 Reference Files 就不会有偏差"的对照分析——这需要展示 43 处偏差中有多少直接归因于未读取规范文档。扣 5 分。

**Urgency justified (20/30):**
提案说"系统性缺陷"和"修复成本极低，拖延无正当理由"。相比 iteration 1 无实质性改进——仍然缺少具体的近期触发事件或延迟成本的量化。"拖延无正当理由"是论证的薄弱形式：它没有解释为什么现在而不是下个迭代。"系统性缺陷"描述的是问题的性质，不是紧迫性。扣 10 分。

---

### 2. Solution Clarity: 93/120

**Approach is concrete (35/40):**
两层防护 + 任务生成层改进的三层结构清晰。每层的职责和修改位置明确：coding.* 模板 Step 1 插入 `<IMPORTANT>` 声明 + Self-Check 插入 AC 验收；task-executor.md Hard Constraints 加兜底规则；quick-tasks/breakdown-tasks 改进生成质量。相比 iteration 1 大幅改善。但仍未给出 `<IMPORTANT>` Reference Files 声明的精确文本——只有意图描述。一个读者可以说出"在 Step 1 插入声明"，但无法复述声明的具体措辞。扣 5 分。

**User-facing behavior described (35/45):**
新增的 User-Visible Behavior Changes section 是 iteration 1 最严重缺陷的直接修复。三种可观察行为（加载确认、AC 验收报告、降级提示）都有具体示例。但仍然缺少一个维度：当 agent 输出这些新增内容时，对整体输出长度和可读性的影响——用户可能因为额外的确认/报告而难以找到核心实现内容。此外，"降级提示"的格式仅为示例性描述（"Reference Files 为空，将以现有代码结构和 Hard Rules 为参照"），未定义为结构化输出。扣 10 分。

**Technical direction clear (23/35):**
技术方向基本清晰——修改 Markdown 模板 + go build。但有两个模糊点：
1. "在 Self-Check 步骤中插入 AC 逐条验收"——Self-Check 步骤在 coding.* 模板中的具体位置是哪里？coding-feature.md 的 workflow 只有 3 步（Read Task Definition → TDD Implementation → Verify & Finalize），Self-Check 是 Step 3 的一部分还是新增 Step 4？
2. quick-tasks/breakdown-tasks 的修改是改 SKILL.md 还是改模板？Scope 说"改进 plugins/forge/skills/quick-tasks/"但未指定具体文件。

扣 12 分。

---

### 3. Industry Benchmarking: 82/120

**Industry solutions referenced (30/40):**
iteration 2 引用了 LangChain RetrievalQA、LlamaIndex QueryEngine、MetaGPT SOP、CrewAI Knowledge、OpenAI Function Calling、Anthropic Tool Use、Claude Prompt Engineering 指南。每个方案都有出处（GitHub URL 或文档 URL）。相比 iteration 1 大幅改善。但引用方式是列举式的——没有深入分析每个方案的具体实现细节和与本提案的精确映射关系。例如"MetaGPT 在 SOP 流程中嵌入文档引用步骤"——MetaGPT 的哪个 SOP？引用步骤的具体形式是什么？扣 10 分。

**At least 3 meaningful alternatives (25/30):**
四个替代方案：Do nothing、CLI auto-inline、仅修改 Hard Constraints、本方案。"Do nothing" 是合格的基线。"仅修改 Hard Constraints" 是 iteration 2 新增的有效替代——它确实是一个更轻量的真实选择。"CLI auto-inline"标记为"自研"——这意味着它是提案者构思的，不是引用行业方案。但作为理论替代是有效的。扣 5 分——缺少"修改现有 Reference Files 解析逻辑而非插入新步骤"这一替代思路。

**Honest trade-off comparison (15/25):**
Comparison Table 的 Cons 列对核心弱点坦诚：本方案的 Cons 是"依赖 agent 遵守 `<IMPORTANT>` 标记；coding.* 模板需 go build"。"仅修改 Hard Constraints"的 Cons 是"仅一层防护，无模板层保障"。但比较仍然偏简略：
- "Do nothing"的 Pros 只说"零成本"——但现实是"现有流程已部分工作，43 处偏差中可能大部分可在 review 中捕获"，这更诚实。
- 本方案"覆盖主要任务类型"——哪些任务类型不被覆盖？未说明。
扣 10 分。

**Chosen approach justified (12/25):**
选择理由是"零代码改动，两层防护"。但 iteration 1 的攻击"仅修改 Hard Constraints 更小"已被部分采纳（作为兜底规则），选择理由未正面回应"为什么不只用 Hard Constraints 就够了"。提案选择的实际上是"两层防护"，应该论证"两层比一层好多少"——这需要引用一些关于 LLM 遵从 prompt 指令的可靠性数据。扣 13 分。

---

### 4. Requirements Completeness: 78/110

**Scenario coverage (28/40):**
五个场景覆盖了 coding.* 任务、doc.* 任务、无 AC 任务、quick-tasks 生成（无 tech-design）、breakdown-tasks 生成。Edge Cases & Degradation 覆盖了四种退化情况。相比 iteration 1 大幅改善。但仍然缺少以下场景：
1. **doc.* 任务中 Reference Files 的差异**：场景 2 仅说"验收重点是文档结构合规而非路径命名"，但 doc.* 任务的 Reference Files 来源和格式是否与 coding.* 相同？doc.* 任务可能没有 tech-design.md。
2. **clean-code / coding-fix 等非特性任务的适配**：这些任务的 Reference Files 可能不来自 tech-design.md 而是来自 bug 报告或 PR 描述。提案未讨论。

扣 12 分。

**Non-functional requirements (25/40):**
两条 NFR："prompt 不过长"和"不改变执行流程结构"。相比 iteration 1 无新增。缺少：
1. **性能影响**：每个任务 Step 1 增加 2-5 个文件读取操作的开销未评估。
2. **向后兼容性**：F2 finding（新旧 Reference Files 格式共存）仍未解决。提案定义了新格式 `path/to/file.md#section — 简要说明`，但未讨论旧格式 `path/to/file.md` 是否继续有效、agent 如何区分新旧格式。
3. **输出体积增长**：新增的加载确认 + AC 验收报告 + 降级提示对 agent 输出的影响未量化。

扣 15 分。

**Constraints & dependencies (25/30):**
三条约束覆盖了 forge-distribution.md 遵循、embed.FS 编译依赖、两层防护同步更新。相比 iteration 1 大幅改善。但缺少一项：go build 后如何验证编译成功且模板内容正确？提案在 Success Criteria 第 8 条提及"通过 forge prompt get-by-task-id 验证"，但 Constraints section 本身没有列出这个验证步骤作为部署依赖。扣 5 分。

---

### 5. Solution Creativity: 42/100

**Novelty over industry baseline (10/40):**
提案本身声明"非创新性改进"。这仍是一个诚实的评估。核心方案（在 prompt 中加强调标记和检查步骤）是 prompt engineering 的标准实践。iteration 2 未在此维度引入任何新元素。扣 30 分。

**Cross-domain inspiration (17/35):**
Level 4 分析（"LLM agent 天然倾向局部一致性而非全局一致性"）来自系统思维。Priority Rules 的三级优先级（Hard Rules > Reference Files > 现有代码）类似于法律体系的层级原则。但仍然缺少编译器类型检查、数据库约束、CI/CD gate check 等成熟领域的借鉴——这些领域有"规范强制执行"的自动化方案，而本提案完全依赖 agent 的指令遵从。扣 18 分。

**Simplicity of insight (15/25):**
核心洞察仍然简洁——"在执行流程中嵌入规范权威性锚点"。但实施方案相比 iteration 1 变得更复杂：Edge Cases & Degradation（4 种场景）、Priority Rules（3 级优先级 + 冲突行为）、Rollback Plan（3 层）。这些内容本身合理，但它们增加了一个声称"修复成本极低"的提案的认知复杂度。扣 10 分。

---

### 6. Feasibility: 72/100

**Technical feasibility (32/40):**
时序问题已解决（在模板 Step 1 内部，读取任务文件之后）。embed.FS 编译依赖已声明。quick-tasks 无 tech-design.md 的场景已处理。三个 iteration 1 的技术障碍全部解决。但仍然存在一个小问题：coding-feature.md 的 workflow 只有 3 步，"在 Self-Check 步骤中插入 AC 逐条验收"意味着要修改 Step 3 的内部结构——这需要理解 Step 3 的现有内容并找到合适的插入点，提案没有展开这部分技术细节。扣 8 分。

**Resource & timeline feasibility (22/30):**
"1-2 小时完成"的估计。审计标准已定义，减少了不确定性。但 Edge Cases、Priority Rules、Rollback Plan 的复杂度增加了实现工作量。特别是 quick-tasks 和 breakdown-tasks 的 SKILL.md 修改需要仔细设计 section 提取逻辑，这部分工作可能超出"1-2 小时"。扣 8 分。

**Dependency readiness (18/30):**
embed.FS 编译依赖已声明。但 Success Criteria 第 8 条"通过 forge prompt get-by-task-id 获取的合成 prompt 中包含 Reference Files 声明文本"依赖于 `forge prompt get-by-task-id` 命令的现有行为——提案假设该命令会原样输出模板内容，但没有验证这个假设。如果该命令对 prompt 有截断或格式化处理，验证可能不准确。扣 12 分。

---

### 7. Scope Definition: 62/80

**In-scope items are concrete (25/30):**
七项 In Scope 条目，每项指定了具体文件和修改内容。审计标准已定义（三条判定条件）。相比 iteration 1 大幅改善。但 "改进 plugins/forge/skills/quick-tasks/" 未指定具体修改的文件——是 SKILL.md 还是 templates/ 下的模板文件？扣 5 分。

**Out-of-scope explicitly listed (18/25):**
四项 Out of Scope 清晰。但 F2 finding 仍未解决：新格式 `path/to/file.md#section` 与旧格式 `path/to/file.md` 的迁移问题。如果旧任务文件使用旧格式，它们是否需要更新？这应该明确列入 In Scope 或 Out of Scope。扣 7 分。

**Scope is bounded (19/25):**
没有明确的完成时间或分阶段交付计划。Next Steps 仍然只说"Proceed to /quick-tasks"。Success Criteria 的"行为验证"（第 9 条）是一个自然的检查点，但提案没有将其定义为 phase gate。iteration 1 的 blindspot-1（缺少最小可验证增量的交付策略）仍未完全解决——虽然 Success Criteria 第 9 条要求执行 1 个任务验证，但这不是分阶段交付策略。扣 6 分。

---

### 8. Risk Assessment: 70/90

**Risks identified (25/30):**
四个风险：忽略标记、引用过时、填充不完整、更新不同步。iteration 1 的 blindspot-5（Reference Files 与 Hard Rules 优先级冲突）通过 Priority Rules section 解决，虽然未在 Risk 表中显式列为风险。缺少以下风险：
1. **旧格式 Reference Files 的解析混乱**：F2 finding 仍未解决。
2. **agent 输出体积显著增长**：新增加载确认 + AC 验收报告可能使输出过长，影响用户阅读效率。

扣 5 分。

**Likelihood + impact rated (23/30):**
风险 1（忽略标记）Likelihood M——提案已将标记从 EXTREMELY-IMPORTANT 降为 IMPORTANT，但未论证为什么 Likelihood 仍然是 M 而非 L。如果选择 IMPORTANT 就是为了降低被忽略的概率，Likelihood 应该相应调整。
风险 4（更新不同步）Likelihood M Impact H——这是新增的风险，评估合理。
扣 7 分。

**Mitigations are actionable (22/30):**
风险 1 的 Mitigation："标记放在 Step 1 和 Self-Check，使用 IMPORTANT 避免标记稀释，task-executor.md Hard Constraints 兜底"——三层防御清晰可操作。风险 2 的 Mitigation："引用 section 标题而非行号，定义降级行为"——Edge Cases 已定义。风险 3 的 Mitigation："SKILL.md 中加显式 checklist"——可操作。风险 4 的 Mitigation："Success Criteria 中加入部署一致性验证"——可操作。但风险 1 的核心缓解是"Hard Constraints 兜底"，而 Hard Constraints 本身也是 prompt 文本——如果 agent 忽略 IMPORTANT 标记，它为什么不会忽略 Hard Constraints 中的兜底规则？这个递归依赖未在 Mitigation 中讨论。扣 8 分。

---

### 9. Success Criteria: 67/80

**Criteria are measurable and testable (42/55):**
九条成功标准大部分可验证：
1. 检查模板 Step 1 内容——可验证
2. 检查 task-executor.md Hard Constraints——可验证
3. 审计报告存在——可验证
4. 同第 1 条——可验证
5. "≥1 个精确 section 引用（非仅 proposal.md）"——对 quick-tasks 场景 4 有歧义：场景 4 说"从 proposal.md 中提取 section"，即 proposal.md 自身的 section 引用。但成功标准说"非仅 proposal.md"——这是否排除了 proposal.md 的 section 引用？如果是，则与场景 4 矛盾。
6. 检查 SKILL.md 内容——可验证
7. 格式统一——可验证
8. 部署一致性验证——可验证
9. 行为验证——可验证（执行 1 个任务后检查输出）

第 5 条的歧义是新增的逻辑问题。扣 8 分。

iteration 1 的"衡量产出而非结果"的攻击通过第 9 条部分解决——至少要求了行为验证。但仍然缺少量化度量：没有"偏差率从 X% 降至 Y%"或"执行 N 个任务后无偏差"的标准。执行 1 个任务不足以统计验证有效性。扣 5 分。

**Coverage is complete (25/25):**
九条成功标准覆盖了：模板层（1, 4）、Agent 层（2）、审计（3）、quick-tasks（5）、breakdown-tasks（6）、格式（7）、部署一致性（8）、行为验证（9）。覆盖了所有 In Scope 条目。完整。

---

### 10. Logical Consistency: 70/90

**Solution addresses the stated problem (28/35):**
问题定义是"Agent 以代码为参照而非以规范为参照"。解决方案的三层防护直接针对这个问题：强制加载 Reference Files（模板层 + Agent 层兜底）+ 提高生成质量（任务生成层）。逻辑链清晰。但有一个逻辑弱点：解决方案的核心假设是"agent 加载了 Reference Files 后就会以规范为参照"——但 agent 完全可以加载了规范文档然后仍然以代码为准。提案没有"规范 > 代码"的执行保证，只有 Priority Rules 的声明。Priority Rules 本身也是 prompt 文本，agent 可以忽略。扣 7 分。

**Scope <-> Solution <-> Success Criteria aligned (23/30):**
七项 In Scope 与九条 Success Criteria 的大致对齐关系：
- Scope 1 (模板) <-> SC 1, 4
- Scope 2 (Agent) <-> SC 2
- Scope 3 (审计) <-> SC 3
- Scope 5 (quick-tasks) <-> SC 5
- Scope 6 (breakdown-tasks) <-> SC 6
- Scope 7 (格式统一) <-> SC 7
- 部署 <-> SC 8
- 行为 <-> SC 9

但 SC 5 与场景 4 存在矛盾（见上）。Scope 说"改进 quick-tasks: Reference Files 从 proposal 和相关文档中提取精确 section 引用"，但 SC 5 说"非仅 proposal.md"——如果 quick-tasks 只有 proposal.md 输入，"非仅 proposal.md"意味着什么？这不是形式上的矛盾（因为场景 4 说"若 proposal.md 中引用了外部设计文档且文件存在，则同时引用"），但它要求 quick-tasks 必须找到至少一个非 proposal.md 的引用源——这个条件在纯代码修改的提案中可能无法保证。扣 7 分。

**Requirements <-> Solution coherent (19/25):**
五个场景与三层解决方案的映射基本清晰。但两个不一致：
1. 场景 2（doc.* 任务）的验收标准说"文档结构合规"，但解决方案未定义 doc.* 模板的 Reference Files 策略是否与 coding.* 相同。如果不同，需要分别描述。
2. 场景 5（breakdown-tasks）要求"从 tech-design.md 的架构决策中提取每个任务相关的 section"，但 Scope 中"为非 UI 任务增加 Reference Files 填充指引"未给出提取逻辑——如何判断一个 section 与特定任务相关？

扣 6 分。

---

## Cross-Dimension Coherence

1. **Problem Definition 说"修复成本极低"但 Edge Cases + Priority Rules + Rollback Plan 的复杂度不低**——Dimensions 1 和 2 之间存在张力。"极低成本"的描述需要考虑修订版增加的所有退化处理和优先级规则的实现成本。

2. **Industry Benchmarking 引用了 Structured Output 方案但选择不采用**——Dimensions 3 和 5 的交互。提案正确地指出了 Structured Output 的侵入性，但未解释为什么不采用一个轻量版的结构化输出（如仅在 AC 验收步骤要求 JSON 格式）。

3. **Risk Assessment 风险 1 的 Mitigation 依赖 Hard Constraints 兜底，但 Hard Constraints 本身也是 prompt**——Dimensions 8 和 10 的递归依赖。如果 agent 忽略 IMPORTANT，它也可能忽略 Hard Constraints 中的兜底规则。两层防护共享同一个失效模式。

4. **Success Criteria SC 5 与场景 4 的张力**——Dimensions 9 和 4 的对齐问题。SC 5 要求"非仅 proposal.md"，但场景 4 的主要来源是 proposal.md。如果 proposal.md 没有引用外部文档，SC 5 可能无法通过。

## Blindspot Hunt

**[blindspot-1] 新旧 Reference Files 格式共存（F2 仍未解决）**

提案定义了新格式 `path/to/file.md#section — 简要说明该 section 定义了什么`。但已有任务文件使用旧格式 `path/to/file.md`。提案未定义：
- 旧格式是否继续有效？
- agent 如何区分新旧格式？
- 如果旧格式有效，"精确 section 引用"的改进目标对旧任务文件不适用，这意味着提案只改善新生成的任务，不修复已有任务。

引用："所有 Reference Files 条目格式统一：path/to/file.md#section — 简要说明该 section 定义了什么"——这条成功标准暗示所有条目都应该是新格式，但它只适用于修改后的模板生成的任务，不适用于已有任务。

**[blindspot-2] "执行 1 个 coding 任务"的行为验证样本量不足**

引用成功标准第 9 条："行为验证：使用修改后的模板执行 1 个 coding 任务后，agent 输出中包含 Reference Files 加载确认和 AC 验收报告"。

1 个任务不足以验证"防止 Agent 偏离"这个核心目标。43 处偏差是在一个大型特性中累积的——验证修复效果至少需要执行一个包含多个任务的特性开发流程。此外，"输出中包含加载确认"只验证了行为发生，不验证行为有效——agent 可以输出确认消息但仍然以代码为准。

**[blindspot-3] 两层防护共享同一失效模式——prompt 遵从率**

引用 Risk 表风险 1 的 Mitigation："task-executor.md Hard Constraints 兜底"。但 Hard Constraints 是在 `<EXTREMELY-IMPORTANT>` 标记内的 prompt 文本。如果 agent 忽略 `<IMPORTANT>` 标记（风险 1），它同样可能忽略 `<EXTREMELY-IMPORTANT>` 标记（Hard Constraints 所在位置）。两层防护在本质上是同一种机制（prompt 指令），共享同一个失效模式。这不是真正的"两层防护"——更像是同一层的重复强调。

**[blindspot-4] AC 逐条验收可能引入新的 agent 行为问题**

引用解决方案："在 Self-Check 步骤中插入 AC 逐条验收"。AC 逐条验收要求 agent 在完成实现后逐条检查 Acceptance Criteria。这可能导致：
1. Agent 过度关注 AC 的字面满足，忽视 AC 未覆盖的实现质量维度。
2. Agent 在 AC 验收中生成大量 pass/fail 输出，增加 prompt 长度，可能触发上下文窗口截断。
3. Agent 为了通过 AC 验收而"迎合"AC 条目（如只实现 AC 中明确列出的内容），忽视隐含的质量要求。

提案未讨论这些潜在副作用。Rollback Plan 仅考虑"过度遵守 Reference Files 导致忽视用户意图"和"AC 验收输出过于冗长"两种情况，但不包括"AC 验收导致 agent 行为扭曲"。

**[blindspot-5] breakdown-tasks 的 section 提取逻辑是未解决的核心难题**

引用场景 5："从 tech-design.md 的架构决策中提取每个任务相关的 section，精确引用而非笼统指向整个文件"。

这是提案最困难的实现细节——如何从一份完整的设计文档中，为每个子任务精确提取相关的 section？这需要理解：(a) 每个任务的范围和职责；(b) tech-design.md 中哪些 section 与该范围相关。对于一个 LLM agent（breakdown-tasks 的执行者），这本质上是一个信息检索 + 相关性判断的任务。提案没有给出任何指引（如"按任务 affected files 匹配 design 中涉及相同文件路径的 section"），只说了目标。

---

## Summary

SCORE: 628/1000
DIMENSIONS:
  Problem Definition: 92/110
  Solution Clarity: 93/120
  Industry Benchmarking: 82/120
  Requirements Completeness: 78/110
  Solution Creativity: 42/100
  Feasibility: 72/100
  Scope Definition: 62/80
  Risk Assessment: 70/90
  Success Criteria: 67/80
  Logical Consistency: 70/90
ATTACKS:
1. [Problem Definition]: 紧迫性论证仍然薄弱——引用"拖延无正当理由"是空洞论证，缺少近期触发事件或延迟成本的量化。需要具体的"如果不修复，下一次大型特性开发将导致 N 处偏差"的推算。
2. [Solution Clarity]: `<IMPORTANT>` Reference Files 声明的精确文本未给出——引用"在 Step 1（读取任务文件之后）插入 `<IMPORTANT>` Reference Files 权威性声明"是意图描述，不是可复述的具体措辞。需要给出声明的模板文本。
3. [Solution Clarity]: Self-Check 步骤在 coding.* 模板中的精确位置未定义——coding-feature.md 只有 3 步 workflow，AC 验收插入点是 Step 3 内部还是新增步骤？需要精确到行号或段落。
4. [Industry Benchmarking]: 选择理由未正面回应"为什么不只用 Hard Constraints 就够了"——引用 Verdict "Selected: 最小有效改动"，但"仅修改 Hard Constraints"更小。需要论证两层防护比一层的增量价值。
5. [Requirements Completeness]: 新旧 Reference Files 格式共存问题未解决——引用 Success Criteria "所有 Reference Files 条目格式统一"，但已有任务文件使用旧格式。需要声明旧格式是否继续有效。
6. [Requirements Completeness]: doc.* 任务的 Reference Files 策略未分化——引用场景 2 "验收重点是文档结构合规"，但未定义 doc.* 模板的 Reference Files 来源和格式。
7. [Solution Creativity]: 两层防护本质上是同一种机制（prompt 指令），共享同一失效模式——引用 Mitigation "task-executor.md Hard Constraints 兜底"，但 Hard Constraints 也是 prompt 文本。需要考虑真正的第二层防护（如运行时验证）或承认这是单层防护的重复强调。
8. [Scope Definition]: quick-tasks 修改的具体文件未指定——引用 Scope "改进 plugins/forge/skills/quick-tasks/"，是 SKILL.md 还是 templates/ 下的模板文件？需要明确。
9. [Success Criteria]: SC 5 "非仅 proposal.md" 与场景 4 "从 proposal.md 中提取 section" 存在张力——quick-tasks 只有 proposal.md 输入时，如果 proposal.md 未引用外部文档，SC 5 无法通过。需要调整 SC 5 或明确场景 4 的外部文档引用是必须的还是可选的。
10. [Success Criteria]: 行为验证样本量不足——引用 SC 9 "执行 1 个 coding 任务"，1 个任务无法统计验证"防止偏离"的有效性。需要至少 3-5 个任务的验证或量化偏差率度量。
11. [Logical Consistency]: "agent 加载规范后就会以规范为参照"的假设未验证——Priority Rules 是 prompt 文本，agent 可以忽略。需要至少一个小规模预实验验证 agent 对 Priority Rules 的遵从率。
12. [Logical Consistency]: breakdown-tasks 的 section 提取逻辑是核心难题但未解决——引用场景 5 "从 tech-design.md 的架构决策中提取每个任务相关的 section"，但未给出提取的决策逻辑。需要在 Scope 中明确这部分的设计细节或将其列为"待实施时细化"。
