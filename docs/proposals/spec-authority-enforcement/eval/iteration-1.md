# Proposal Evaluation — Iteration 1

## Reasoning Audit

**Pre-score anchors:**

1. **Problem-Solution 时序断裂**：提案要求在 task-executor.md 的 Step 5 之前插入 Reference Files 声明步骤（Success Criteria: "Step 5 前"），但 task-executor.md 的执行协议中，agent 在 Step 3 获取合成 prompt，Step 5 才将控制权交给模板执行。在 Step 5 之前，agent 尚未读取任务文件，任务文件中的 `## Reference Files` section 对 agent 不可见。这意味着提案要求 agent 在看不到 Reference Files 的情况下声明其为权威来源——这是一个逻辑矛盾。自由评审精确指出："Reference Files 的权威性声明必须在 agent 读到任务文件之后才能生效。"

2. **quick-tasks 与 breakdown-tasks 的 Reference Files 策略未分化**：提案对 quick-tasks（无 tech-design.md）和 breakdown-tasks（有 tech-design.md）提出了相同的 Reference Files 精确引用要求，但两者输入不同。quick-tasks 的输入只有 proposal.md，不存在 tech-design.md，无法提取 section 引用。

3. **部署一致性盲点**：提案声明"纯 Markdown 修改，不涉及 Go 代码变更"，但修改 `forge-cli/pkg/prompt/data/*.md` 后必须 `go build` 才能生效。如果 task-executor.md 已更新（即时生效）但 coding.* 模板未重编译，两层防护不一致。

4. **审计标准推迟到实施阶段**：Scope 列出"审计全部 19 个模板"作为交付项，但未定义判断哪些模板需要强化的标准。这等于将核心决策推迟。

5. **标记稀释风险已承认但未解决**：提案在 Key Risks 中承认 Agent 可能忽略 `EXTREMELY-IMPORTANT` 标记，但 Mitigation 仅是"放在关键位置"，没有解决多 `EXTREMELY-IMPORTANT` 块的稀释问题。

## Dimension Scores

### 1. Problem Definition: 82/110

**Problem stated clearly (35/40):**
核心问题明确——"Agent 执行任务时以现有代码为参照而非以权威规范文档为参照，导致大规模偏离"。引用具体事件（test-capability-v2，43 处偏差）。问题陈述清晰，但存在一个小模糊点：标题说"从模板和任务生成层面防止 Agent 偏离"，这预设了解决方案的层面，而非纯粹描述问题。扣 5 分。

**Evidence provided (32/40):**
引用了教训文档 `docs/lessons/gotcha-spec-authority-drift.md` 的 Level 0-4 五级溯源。Level 0 根因精确到"coding.* 模板的 Step 1 只说'read the task file'，未要求读 Reference Files"。Level 3 的"反应式局部修复而非主动式全局审计"也是可验证的过程描述。但证据存在一个关键遗漏：没有展示"如果当时 agent 读 Reference Files 就不会有偏差"的对照实验或反事实分析。43 处偏差的原始数据未在提案中呈现（路径列表等），需要跳转到教训文档才能看到。扣 8 分。

**Urgency justified (15/30):**
提案说"这是系统性缺陷"和"修复成本极低（纯文档修改），拖延无正当理由"。但这不是紧迫性论证——它只是说问题存在且修复成本低。"系统性缺陷"一词不构成紧迫性论证；什么事件在近期会触发这个问题？是否有即将到来的大型特性开发依赖此修复？没有具体的近期触发事件或成本量化（例如"下一次大规模重构将导致同等规模的偏差"）。扣 15 分。

---

### 2. Solution Clarity: 55/120

**Approach is concrete (25/40):**
两层防护的框架清晰：Agent 层（task-executor.md 增加步骤）+ 任务生成层（quick-tasks/breakdown-tasks 改进 Reference Files 质量）。但"增加两个强制步骤"的具体内容模糊——提案没有给出新增步骤的精确文本或伪代码。例如"声明 Reference Files 为权威来源，按需加载"是一条意图描述，不是可执行的步骤定义。一个读者无法复述出具体的步骤文本。扣 15 分。

**User-facing behavior described (10/45):**
这是本维度最严重的扣分点。提案几乎没有描述终端用户的体验变化。用户（使用 Forge 的开发者）在执行任务时会看到什么不同？agent 的输出格式会变吗？任务文件的 `## Reference Files` section 的外观会变吗？提案的四个场景（场景 1-4）描述的是 agent 的内部行为流程，不是用户感知到的变化。扣 35 分。

**Technical direction clear (20/35):**
技术方向是"修改 Markdown 模板文件"，这本身足够清晰。但存在严重的时序问题：提案要求在 task-executor.md 的 Step 5 之前插入 Reference Files 声明步骤，但 task-executor.md 的执行协议中 agent 在 Step 3 获取合成 prompt，Step 5 之前 agent 尚未读取任务文件中的 `## Reference Files` section。这个时序矛盾意味着技术方向的可行性存疑。此外，`embed.FS` 的一致性问题意味着 coding.* 模板的修改需要 `go build`，但提案声称"纯 Markdown 修改不涉及 Go 代码"。扣 15 分。

---

### 3. Industry Benchmarking: 48/120

**Industry solutions referenced (15/40):**
提案提到"RAG 检索、prompt 模板强制引用、structured output 约束"，但这三个词组是概念列举而非引用。没有引用任何具体的产品、论文、框架或公开实践。例如 RAG 方面可以引用 LangChain 的 retrieval QA、LlamaIndex 的 query engine；structured output 可以引用 OpenAI 的 function calling 或 Anthropic 的 tool use。提案完全没有给出出处。扣 25 分。

**At least 3 meaningful alternatives (18/30):**
三个替代方案：Do nothing、CLI 层 auto-inline、Prompt 模板 + Agent 协议。"Do nothing" 是合格的基线。"CLI 层 auto-inline"被描述为"自研"——这意味着这不是一个已有的替代方案，而是提案者自己构思的。第三个是提案本身。真正缺少的是一个"仅修改 task-executor.md Hard Constraints 而不修改模板"的轻量替代方案——自由评审实际上提出了这个方案。扣 12 分。

**Honest trade-off comparison (5/25):**
比较表过于简略。CLI auto-inline 方案的 Cons 仅有"需改 Go 代码，增加 prompt 长度"，没有讨论它的可靠性优势（agent 无法跳过）。提案方案的 Cons 是"依赖 agent 遵守 `<EXTREMELY-IMPORTANT>` 标记"——这是核心弱点，但在比较中没有展开讨论其影响程度。"Do nothing"的 Pros 是"零成本"而非"现有流程已经部分工作"。扣 20 分。

**Chosen approach justified (10/25):**
选择的理由是"零代码改动，覆盖所有任务类型"和"最小有效改动"。但"覆盖所有任务类型"这个主张有问题——自由评审指出 coding.* 模板通过 `embed.FS` 嵌入二进制，修改后需要重新编译，且如果不重编译则两层防护不一致。"最小有效改动"缺少证据——自由评审建议只修改 task-executor.md Hard Constraints 加一条规则，这比提案的两层防护更小。扣 15 分。

---

### 4. Requirements Completeness: 52/110

**Scenario coverage (18/40):**
四个场景覆盖了 coding、doc、quick-tasks 生成、breakdown-tasks 生成。但缺少关键边缘场景：(1) 任务文件中 `## Reference Files` 为空或缺失时的行为；(2) Reference Files 引用的 section 标题不存在时的行为（自由评审指出）；(3) quick-tasks 没有 tech-design.md 时的行为（自由评审指出）；(4) 嵌入模板（coding.*）和即时代效模板（task-executor.md）更新不同步时的行为。扣 22 分。

**Non-functional requirements (14/40):**
仅提到两条 NFR："Reference Files 精确引用不应导致 prompt 过长"和"模板改动不改变现有任务的执行流程结构"。缺少：(1) 性能影响评估——每个任务增加 2-5 个文件读取操作的开销；(2) 向后兼容性——新格式 `path/to/file.md#section` 与旧格式 `path/to/file.md` 的共存问题；(3) 可维护性——section 引用在 design 文档更新后的维护成本。自由评审精确指出："每个任务的 Step 1 都多出 2-5 个文件读取操作，这可能影响执行效率。"扣 26 分。

**Constraints & dependencies (20/30):**
列出两条约束："修改 plugins/forge/ 下的文件前必须遵循 forge-distribution.md"和"纯 Markdown 文档修改，不涉及 Go 代码变更"。第二条约束在技术上不准确——修改 `forge-cli/pkg/prompt/data/*.md` 后需要 `go build` 才能使嵌入的模板生效，这意味着虽然不改 `.go` 文件，但需要执行 Go 编译流程。缺少的依赖：`embed.FS` 编译依赖、两层防护更新同步依赖。扣 10 分。

---

### 5. Solution Creativity: 40/100

**Novelty over industry baseline (10/40):**
提案本身声明"非创新性改进"。这不扣分——诚实是好事。但提案的创新点来自 Level 4 分析："LLM agent 天然倾向局部一致性而非全局一致性"。这个洞察本身有新意，但解决方案（在 prompt 中加强调标记和检查步骤）是 prompt engineering 的标准实践，没有超越行业基线。扣 30 分。

**Cross-domain inspiration (15/35):**
提案从认知科学/心理学中借鉴了"局部优化陷阱"的概念（Level 4 分析），这来自系统思维领域。但解决方案没有借鉴其他领域的验证机制——例如编译器的类型检查（自动验证）、数据库的约束（强制数据完整性）、CI/CD 的 gate check（自动化验收）。这些领域有成熟的"规范强制执行"方案，提案没有参考。扣 20 分。

**Simplicity of insight (15/25):**
核心洞察是简洁的："在 agent 的执行流程中嵌入规范权威性锚点"。但实施方案不够简洁——两层防护、19 个模板审计、两种任务生成策略的修改，对于一个"修复成本极低"的问题来说偏重。自由评审建议的方案（仅在 task-executor.md Hard Constraints 加一条规则 + 在模板 Step 1 的正确位置加 `<IMPORTANT>` 声明）更简洁。扣 10 分。

---

### 6. Feasibility: 55/100

**Technical feasibility (25/40):**
纯 Markdown 修改在技术上可行。但存在三个技术障碍：(1) 时序矛盾——在 task-executor.md Step 5 之前插入 Reference Files 声明，此时 agent 尚未读取任务文件；(2) `embed.FS` 一致性——coding.* 模板修改后需要 `go build`，但提案声称不需要；(3) quick-tasks 没有 tech-design.md 作为输入，无法提取精确 section 引用。这三个障碍都在自由评审中指出，提案没有回应。扣 15 分。

**Resource & timeline feasibility (20/30):**
"1-2 小时完成"的时间估计有问题。提案要求审计全部 19 个模板（没有标准），然后修改确认需要的模板，同时修改 task-executor.md、quick-tasks 和 breakdown-tasks。19 个模板的审计本身就需要 1-2 小时（每个模板 3-5 分钟的阅读和判断），加上修改和测试，1-2 小时估计偏紧。但提案承认"实际需修改数量待定"，这增加了不确定性。扣 10 分。

**Dependency readiness (10/30):**
提案说"无外部依赖"。但 `embed.FS` 的编译依赖被忽略了——修改 `forge-cli/pkg/prompt/data/*.md` 后必须 `go build` 才能使 coding.* 模板的修改生效。这不是外部依赖，但是一个被忽略的内部部署依赖。此外，自由评审指出的"两层防护不一致"问题意味着解决方案的生效依赖于用户同时更新所有层级——这个依赖条件未被声明。扣 20 分。

---

### 7. Scope Definition: 42/80

**In-scope items are concrete (18/30):**
五项 In Scope 条目中，前三项是具体可交付的（修改 task-executor.md、审计模板、修改模板）。但"审计全部 19 个模板"缺少审计标准——提案说"确定哪些需要 Reference Files 权威性声明"，但没有定义判断标准。自由评审指出："提案的 Scope 列出'审计全部 19 个 forge-cli/pkg/prompt/data/*.md 模板'但没有定义审计标准。等于将核心决策推迟到实施阶段。"扣 12 分。

**Out-of-scope explicitly listed (12/25):**
四项 Out of Scope 条目清晰：不改 Go 代码、不改 index.json schema、不添加 hooks、不改合成逻辑。但缺少一项关键的 Out of Scope：已生成的任务文件的 Reference Files 格式迁移。如果新格式是 `path/to/file.md#section`，旧格式是 `path/to/file.md`，已有任务文件是否需要迁移？扣 13 分。

**Scope is bounded (12/25):**
没有明确的完成时间或迭代边界。提案的 Next Steps 只说"Proceed to /quick-tasks"，没有分阶段交付计划。Success Criteria 中的"审计报告"是一个自然的检查点，但提案没有将其定义为 phase gate。扣 13 分。

---

### 8. Risk Assessment: 48/90

**Risks identified (20/30):**
三个风险列出：(1) Agent 忽略 `EXTREMELY-IMPORTANT` 标记，(2) Reference Files section 引用过时，(3) breakdown-tasks 生成任务时 Reference Files 填充不完整。这些都是有效的风险。但缺少以下风险：(1) 部署不一致（task-executor.md 已更新但 coding.* 模板未重编译），这是自由评审的高优先级发现；(2) 新旧 Reference Files 格式共存导致的解析混乱；(3) AC 验收步骤在非 coding 任务类型上的适配问题。扣 10 分。

**Likelihood + impact rated (10/30):**
Likelihood 和 Impact 使用 H/M/L 标注。风险 1（忽略标记）的 Likelihood 是 M——但提案的核心假设是 agent 会遵守标记，如果 Likelihood 是 M（约 50%），解决方案的有效性就打了折扣。风险 2（引用过时）的 Mitigation 是"引用 section 标题而非行号"，但这假设 section 标题稳定——没有证据支持。风险 3（填充不完整）的 Likelihood 是 L——但 breakdown-tasks 的 SKILL.md 当前完全没有 Reference Files 填充指引，L 可能低估了。评分标准不一致。扣 20 分。

**Mitigations are actionable (18/30):**
风险 1 的 Mitigation 是"标记放在执行流程的关键位置（Step 1 和 submit 前），并声明'按需加载'降低遵从成本"——这没有解决核心问题（标记稀释）。风险 2 的 Mitigation 是"引用 section 标题而非行号"——这是部分可操作的，但没有降级行为（section 标题不存在时怎么办）。风险 3 的 Mitigation 是"在 SKILL.md 中加显式 checklist"——这是可操作的。扣 12 分。

---

### 9. Success Criteria: 50/80

**Criteria are measurable and testable (32/55):**
六条 Success Criteria 中，大部分有明确的可验证条件：(1) "task-executor.md 包含两个步骤"——可验证（检查文件内容）；(2) "全部 19 个模板完成审计，输出审计报告"——可验证（报告存在性）；(3) "模板 Step 1 包含 EXTREMELY-IMPORTANT Reference Files 权威性声明"——可验证。但以下问题影响可测试性：

- "精确 section 引用（非仅 proposal.md）"——"非仅 proposal.md"是排除性条件，不是正向定义。什么算"精确"？section 标题？行号范围？
- "Reference Files 填充规则"——什么标准判断规则是否完整？没有给出。
- "所有 Reference Files 条目格式统一"——格式定义为 `path/to/file.md#section — 简要说明`，但没有说明如何处理没有 section 概念的文件（如 proposal.md 本身）。

此外，成功标准缺少关键可测试项：(1) 没有"agent 在执行 coding 任务时实际读取了 Reference Files 中列出的文档"的验证机制；(2) 没有"AC 逐条验收步骤实际防止了偏差"的度量。成功标准衡量的是文档修改的完成度，而非问题解决的有效性。扣 23 分。

**Coverage is complete (18/25):**
成功标准覆盖了 task-executor.md 修改、模板审计、模板修改、quick-tasks 改进、breakdown-tasks 改进、格式统一。但缺少以下覆盖：(1) 部署一致性验证（coding.* 模板修改后 go build 的执行）；(2) 已有任务文件的 Reference Files 格式迁移。扣 7 分。

---

### 10. Logical Consistency: 45/90

**Solution addresses the stated problem (15/35):**
问题定义是"Agent 以现有代码为参照而非以规范文档为参照"。解决方案的两层防护在方向上正确——强制 agent 读 Reference Files 并按 AC 验收。但存在关键逻辑断裂：

1. 提案要求在 task-executor.md Step 5 之前插入 Reference Files 声明步骤，但此时 agent 尚未读取任务文件（控制权在 Step 5 才交给模板）。自由评审精确指出这个时序矛盾："Reference Files 的权威性声明必须在 agent 读到任务文件之后才能生效。"如果提案的实际意图是在 coding.* 模板的 Step 1 内部（读取任务文件之后）插入声明，那么成功标准中"Step 5 前"的定位是错误的。

2. 提案要求 quick-tasks 生成精确 section 引用，但 quick-tasks 的输入只有 proposal.md，可能没有 tech-design.md。自由评审指出："当没有 tech-design.md 时，quick-tasks 无法提取'精确 section 引用'——因为根本不存在 design 文档。"

3. 成功标准要求 AC 验收步骤在 Step 8 前（submit-task 前），但 task-executor.md 的 Step 8 是 submit-task 本身，Step 7 是 check blocked。AC 验收应该在哪里？如果是在 Step 7.5（check blocked 和 submit 之间），那当任务被 blocked 时 AC 验收是否跳过？

扣 20 分。

**Scope <-> Solution <-> Success Criteria aligned (15/30):**
Scope 列出五项 In Scope。Solution 描述两层防护。Success Criteria 有六条检查项。对齐问题：

1. Scope 第 2 项"审计全部 19 个模板"和 Success Criteria 第 2 条"输出审计报告"对齐。但审计标准未定义（Scope 和 Success Criteria 都没有提供标准），所以这个对齐是形式上的而非实质上的。
2. Scope 第 1 项"修改 task-executor.md：增加步骤"和 Success Criteria 第 1 条"执行协议包含步骤（Step 5 前）"对齐。但自由评审指出 Step 5 前的定位有时序矛盾。
3. Solution 说"Agent 层增加两个步骤"和"任务生成层改进 Reference Files 质量"，但 Success Criteria 没有度量 Reference Files 质量改进的实际效果——只度量了文档是否修改。

扣 15 分。

**Requirements <-> Solution coherent (15/25):**
场景 1（coding.* 任务执行）和场景 3（quick-tasks 生成）的要求与解决方案部分对齐。但场景 3 要求"从 proposal 提取关键技术约束，将相关 design 文档的精确 section 写入 Reference Files"——这假设 quick-tasks 能获取 design 文档，但 quick-tasks 的输入只有 proposal.md。Requirements 中对 quick-tasks 的场景假设与实际约束不一致。扣 10 分。

---

## Cross-Dimension Coherence

1. **Problem Definition 说"修复成本极低"但 Scope 涉及审计 19 个模板 + 修改 4-6 个文件 + 无审计标准**——Dimensions 1 和 7 之间存在张力。"极低成本"与"审计全部 19 个模板"的数量级不匹配。

2. **Solution Clarity 说"在执行协议中增加两个强制步骤"但 Feasibility 发现时序矛盾**——Dimensions 2 和 6 之间存在矛盾。如果 Step 5 前的定位在技术上行不通，解决方案本身需要重新设计。

3. **Industry Benchmarking 的比较表声称"覆盖所有任务类型"但 Requirements Completeness 发现 AC 验收对不同任务类型不适配**——Dimensions 3 和 4 之间存在不一致。

4. **Risk Assessment 未包含部署不一致风险，但 Feasibility 发现 `embed.FS` 编译依赖**——Dimensions 6 和 8 的风险覆盖不完整。

5. **Success Criteria 衡量文档修改完成度而非问题解决有效性**——Dimensions 9 和 1 之间的目标偏差。问题定义是"Agent 偏离规范"，但成功标准没有度量偏离率的降低。

## Blindspot Hunt

**[blindspot-1] 缺少"最小可验证增量"的交付策略**
提案的 Next Steps 只说"Proceed to /quick-tasks to generate and execute tasks"。这意味着提案的所有改动将作为一批一次性交付。但提案涉及的文件横跨 agent 层、模板层、skill 层——修改任何一个都可能导致 agent 行为变化，一次性修改全部将无法定位问题。提案缺少分阶段交付策略（例如先只修改 coding-enhancement.md 和 coding-feature.md 两个最关键的模板，验证效果后再推广到其他模板）。

**[blindspot-2] 缺少回滚计划**
提案修改的是核心执行流程模板（task-executor.md）和合成 prompt 模板（coding.*）。如果修改后 agent 行为异常（例如过度遵守 Reference Files 导致忽视用户意图），如何回滚？提案没有讨论回滚策略。作为一个修改 agent 核心执行协议的提案，这是显著遗漏。

**[blindspot-3] 成功标准衡量的是产出而非结果**
所有六条成功标准衡量的都是"是否修改了文档"而非"修改后 agent 行为是否改善"。没有一条标准是"使用新模板执行 N 个任务后，偏差率从 X% 降至 Y%"。这意味着即使所有成功标准通过，问题可能仍未解决。提案将"修改完成"等同于"问题解决"，这是一个经典的"完成偏差"（completion bias）。

**[blindspot-4] 提案未讨论 `<EXTREMELY-IMPORTANT>` 标记在现有 Claude 模型上的实际遵从率**
提案的核心依赖是"Agent 会遵守 `<EXTREMELY-IMPORTANT>` 标记"。但提案没有提供任何证据表明这个标记对 Claude 模型的遵从率有可测量的影响。风险表将其 Likelihood 设为 M（中等），这意味着提案预期的成功率可能不到 50%——对于一个"系统性缺陷"的修复来说，这是不够的。提案应该引用 Claude 的 system prompt 遵从率数据或至少进行小规模验证。

**[blindspot-5] 提案未考虑 Reference Files 声明与 Hard Rules 的优先级冲突场景**
提案要求 Reference Files 为权威来源。但任务文件可能同时包含 `## Hard Rules`（硬规则）和 `## Reference Files`（规范文档）。如果 Hard Rules 说"MUST NOT modify existing test files"但 Reference Files 中的 tech-design.md 要求"重构测试文件路径"，agent 应该遵循哪个？提案没有定义 Reference Files 与 Hard Rules 的优先级关系。

## Summary
SCORE: 517/1000
