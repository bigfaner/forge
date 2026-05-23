# Proposal Evaluation — Iteration 3

## Iteration 2 Issue Tracking

| # | Iteration 2 Attack | Status | How Addressed |
|---|---------------------|--------|---------------|
| 1 | 紧迫性论证薄弱，缺少延迟成本量化 | **Resolved** | Urgency section 新增量化推算："test-capability-v2 约 15 个任务产生 43 处偏差（~2.9 处/任务），下一个大型特性预期偏差 43 处，每处修复 15-30 分钟 review+rework，总计 10-20 小时 wasted effort。修复成本 1-2 小时，ROI 5-10 倍。"推算逻辑清晰，基于历史数据外推。 |
| 2 | `<IMPORTANT>` Reference Files 声明的精确文本未给出 | **Resolved** | Proposed Solution 中新增了完整的声明块模板文本（4 条指令 + 2 条 fallback 输出），包括加载确认格式和 fallback 行为。Self-Check AC 验收也有独立模板文本。读者可以完整复述声明措辞。 |
| 3 | Self-Check 步骤在 coding.* 模板中的精确位置未定义 | **Resolved** | 明确定义："AC 逐条验收不是新增步骤，而是插入到 coding.* 模板现有 Self-Check / Verify 步骤（如 coding-feature.md 的 Step 3 'Verify & Finalize'）的内部，作为该步骤的第一个子步骤"。给出了具体模板文件名和步骤名。 |
| 4 | 选择理由未回应"为什么不只用 Hard Constraints 就够了" | **Resolved** | Industry Benchmarking section 新增大段论述解释两层防护的必要性：Hard Constraints 是"全局锚定"（agent-level），模板层是"执行时锚定"（task-level），类比"具体条文 vs 宪法原则"。并分析了增量价值："模板层失效时 Hard Constraints 仍可触发；Hard Constraints 被忽略时模板层在执行流程关键位置再次提醒"。 |
| 5 | 新旧 Reference Files 格式共存问题未解决 | **Resolved** | NFR section 新增"新旧格式兼容性"条目："旧格式 `path/to/file.md`（无 section 锚点）继续有效。agent 遇到旧格式时，加载整个文件作为参考。新格式是推荐格式，本次修改仅影响新任务的生成模板，不要求迁移已有任务文件"。 |
| 6 | doc.* 任务的 Reference Files 策略未分化 | **Resolved** | 场景 2 扩展为三个子场景：(a) 从 tech-design.md 生成文档 → 引用 tech-design.md 相关 section；(b) 修改现有文档 → 引用文档本身 + 相关设计文档 section；(c) 仅涉及 proposal.md → 与场景 4 一致。分化策略清晰。 |
| 7 | 两层防护本质上是同一种机制 | **Resolved** | Innovation Highlights 新增"关于'两层防护'的诚实声明"段："本提案的'两层'本质上是同一种机制（prompt 指令）在不同粒度和时机上的重复强调，而非编译器类型检查 + 运行时断言那种独立失效模式的两层防护。"并明确将运行时 hook 验证列为 Out of Scope。这是罕见的自省式诚实。 |
| 8 | quick-tasks 修改的具体文件未指定 | **Resolved** | Scope In Scope 第 5 条明确指定："改进 plugins/forge/skills/quick-tasks/SKILL.md（不是 templates/——quick-tasks 的任务生成逻辑在 SKILL.md 的 workflow 步骤中）"。括号内的解释消除了歧义。 |
| 9 | SC 5 "非仅 proposal.md" 与场景 4 存在张力 | **Resolved** | SC 5 修改为："quick-tasks 生成的 coding 任务 Reference Files 包含 >=1 个精确 section 引用（格式为 `file.md#Section-Title`）；当输入仅有 proposal.md 时，引用 proposal.md 中的具体 section 即可满足此条件（'非仅 proposal.md'指不能只写裸文件路径 `proposal.md` 而无 section 锚点）"。括号内的澄清消除了歧义。 |
| 10 | 行为验证样本量不足（1 个任务） | **Resolved** | SC 9 修改为"执行 3-5 个 coding 任务后，每个任务的 agent 输出中均包含 Reference Files 加载确认和 AC 验收报告，且无规范偏离"。样本量从 1 提升至 3-5，且新增了"无规范偏离"的结果验证（不只是行为验证）。 |
| 11 | "agent 加载规范后就会以规范为参照"的假设未验证 | **Resolved** | Next Steps 新增预实验："先用 1 个 coding 任务手动在 prompt 中注入 Reference Files 声明（不修改模板），验证 agent 是否确实以规范为参照而非以代码为参照。若预实验显示 agent 仍以代码为准，则需要重新评估方案"。这是一个合理的假设验证策略。 |
| 12 | breakdown-tasks 的 section 提取逻辑未给出 | **Resolved** | Scope 第 6 条新增提取决策逻辑："对每个生成的任务：(1) 从任务 Affected Files 提取文件路径；(2) 在 tech-design.md 中搜索提及这些文件路径的 section；(3) 提取与任务描述关键词匹配的架构决策 section；(4) 合并去重保留 2-5 个最相关的 section。"并标记为"待实施时细化"——这是合理的处理方式，因为执行者是 LLM agent 而非确定性程序。 |

### Iteration 2 Blindspot Tracking

| # | Blindspot | Status | How Addressed |
|---|-----------|--------|---------------|
| B1 | 新旧 Reference Files 格式共存（F2） | **Resolved** | NFR 新增兼容性声明。 |
| B2 | 行为验证样本量不足 | **Resolved** | SC 9 从 1 个任务提升至 3-5 个，且要求"无规范偏离"。 |
| B3 | 两层防护共享同一失效模式 | **Resolved** | Innovation Highlights 诚实声明承认了这一事实。 |
| B4 | AC 逐条验收可能引入新的 agent 行为问题 | **Partially Resolved** | Rollback Plan 新增了"过度遵守 Reference Files 导致忽视用户意图"和"AC 验收输出过于冗长"的回滚场景，但未讨论"agent 为了通过 AC 而迎合 AC"或"agent 忽视 AC 未覆盖的质量维度"的问题。不过，这属于边缘风险，不影响整体评分。 |
| B5 | breakdown-tasks section 提取逻辑 | **Resolved** | Scope 新增四步提取决策逻辑。 |

---

## Dimension Scores

### 1. Problem Definition: 102/110

**Problem stated clearly (38/40):**
核心问题清晰："Agent 执行任务时以现有代码为参照而非以权威规范文档为参照，导致大规模偏离"。引用了 43 处偏差的具体事件。根因分析到 Level 0（模板 Step 1 未要求读 Reference Files）和 Level 3（缺少自顶向下验证）。问题陈述不预设特定解决方案——虽然标题提到"从模板和任务生成层面防止"，但这已经是问题描述的一部分而非解决方案暗示。扣 2 分：标题中的"从模板和任务生成层面"仍然包含了一点方案暗示，如果纯描述问题可以改为"Agent 规范偏离的系统性根因与防护"之类。

**Evidence provided (38/40):**
强证据链：(1) 教训文档的 5 级溯源；(2) Level 0 根因精确到"模板 Step 1 未要求读 Reference Files"；(3) 探索发现的具体事实（quick-tasks 硬编码 proposal.md、breakdown-tasks 无填充指引）；(4) Urgency section 的量化推算（15 个任务 x 2.9 处/任务 ≈ 43 处，每处 15-30 分钟，总计 10-20 小时）。推算基于可验证的历史数据。扣 2 分：仍然缺少"43 处偏差中有多少直接归因于未读取规范文档"的精确归因——虽然 Level 0 根因已指向模板缺陷，但 43 处偏差可能包含其他类型的偏差（如 agent 理解错误而非参照源错误）。

**Urgency justified (26/30):**
iteration 3 的 Urgency section 有重大改进。量化推算清晰：预期下一次大型特性产生约 43 处偏差，修复成本 10-20 小时 vs 修复方案成本 1-2 小时，ROI 5-10 倍。"系统性缺陷——每次涉及规范驱动修改的任务都可能重蹈覆辙"说明了问题的持续性。扣 4 分：(1) 推算假设下一个大型特性规模相当（15 个任务），这个假设未经验证；(2) "10-20 小时 wasted effort"的计算依赖于"每处偏差修复 15-30 分钟"的经验值，但未说明这个经验值是否来自 test-capability-v2 的实际修复耗时记录；(3) 缺少"为什么不是下个迭代"的时间敏感性论证——ROI 高只能说明应该做，不能说明应该现在做。

---

### 2. Solution Clarity: 112/120

**Approach is concrete (39/40):**
三层结构完整且每层有精确描述：
1. 模板层：在 coding.* Step 1（读取任务文件之后）插入 `<IMPORTANT>` 声明——给出了完整的声明块模板文本（4 条 MUST 指令 + 2 条 fallback 输出格式）
2. Agent 层：task-executor.md Hard Constraints 增加兜底规则
3. 任务生成层：quick-tasks/breakdown-tasks 改进 Reference Files 生成质量

Self-Check AC 验收也有独立的模板文本。AC 验收的插入方式精确到"插入到 coding-feature.md 的 Step 3 'Verify & Finalize' 的内部，作为第一个子步骤"。一个读者可以准确复述将要构建什么。扣 1 分：模板文本中 "Load each Reference File listed in `## Reference Files` immediately after reading the task file" 中的 "Load" 在 agent 的执行模型中是指调用 Read 工具读取文件还是指"声明为权威来源"？如果 Reference Files 列出的文件很长（如完整 tech-design.md），agent 是否需要全文读取？声明文本中 "load" 和 "treat as authoritative" 的区别未完全澄清。

**User-facing behavior described (40/45):**
三种可观察行为变化（加载确认、AC 验收报告、降级提示）都有具体示例输出。每种行为的触发条件和降级路径在 Edge Cases 中完整定义。扣 5 分：(1) 仍然缺少对输出长度影响的量化评估——新增的加载确认（1 行）+ AC 验收报告（每条 AC 1 行，假设 3-8 条 AC）+ 可能的降级提示（1-2 行）总共增加约 5-12 行输出，这对 agent 输出的可读性影响未讨论；(2) "降级提示"虽然给出了示例文本，但未定义为结构化输出格式——不同 agent 可能在降级时输出不同格式的警告，影响用户体验一致性。

**Technical direction clear (33/35):**
技术方向清晰：修改 Markdown 模板 + go build（coding.* 模板）+ 即时生效（task-executor.md）。具体的修改文件路径明确：`forge-cli/pkg/prompt/data/coding.*.md`、`plugins/forge/agents/task-executor.md`、`plugins/forge/skills/quick-tasks/SKILL.md`、`plugins/forge/skills/breakdown-tasks/SKILL.md`。AC 验收的插入位置精确定义。扣 2 分：(1) `<IMPORTANT>` vs `<EXTREMELY-IMPORTANT>` 的分层标记策略在提案中有说明（"避免标记稀释"），但未讨论模板中已有 `EXTREMELY-IMPORTANT` 块与新增 `IMPORTANT` 块的位置关系——如果两者紧邻，用户可能质疑分层标记的实际效果；(2) breakdown-tasks SKILL.md 的四步提取逻辑标记为"待实施时细化"，这意味着技术方向在此处仍然有设计空间。

---

### 3. Industry Benchmarking: 98/120

**Industry solutions referenced (35/40):**
引用了 7 个具体方案：LangChain RetrievalQA、LlamaIndex QueryEngine、MetaGPT SOP、CrewAI Knowledge、OpenAI Function Calling、Anthropic Tool Use、Claude Prompt Engineering 指南。每个方案有出处（GitHub URL 或文档 URL）。更重要的是，对每个方案给出了具体的分析维度：RAG 解决"信息可达性"而非"信息权威性"；MetaGPT/CrewAI 通过执行流程强制加载——"与本提案最接近"；Structured Output 可要求结构化输出但"侵入性较高"；Prompt Engineering 的认知锚定理论"与提案策略一致"。扣 5 分：引用方式虽然比 iteration 2 更深入，但对 MetaGPT 的引用仍然不够精确——"MetaGPT 在 SOP 流程中嵌入文档引用步骤"没有指定 MetaGPT 的哪个 SOP（ProductManager SOP? Architect SOP? Engineer SOP?），也没有说明嵌入的具体形式。

**At least 3 meaningful alternatives (28/30):**
四个替代方案：Do nothing、CLI auto-inline（自研）、仅修改 Hard Constraints、本方案（Prompt 模板 + Agent 协议）。"仅修改 Hard Constraints"是一个真正更轻量的替代——iteration 3 正面回应了它的存在并将其纳入为兜底层。Comparison Table 的结构清晰，每个方案有 Pros/Cons/Verdict。扣 2 分：缺少"修改现有 Reference Files 解析逻辑（如在 forge prompt get-by-task-id 合成时自动 inline Reference Files 内容）"这一替代——这介于"CLI auto-inline"和"纯 prompt 修改"之间，是一个值得考虑的中间路径。

**Honest trade-off comparison (20/25):**
Comparison Table 的 Cons 列对核心弱点坦诚。iteration 3 在 Innovation Highlights 中新增了"关于两层防护的诚实声明"，承认两层本质上是同一种机制的重复强调——这是非常罕见的自省。扣 5 分：(1) "Do nothing" 的 Pros 说"零成本"但未提及"现有流程已部分工作，review 可以捕获部分偏差"——这是更诚实的评估；(2) 本方案 Verdict "Selected: 最小有效改动"——但"仅修改 Hard Constraints"改动更小，Verdict 应该是"最小有效改动（考虑两层锚定的增量价值）"而非直接跳到"Selected"。

**Chosen approach justified (15/25):**
iteration 3 新增大段论述解释两层防护的必要性（类比"具体条文 vs 宪法原则"），分析了触发时机差异（task-level vs agent-level）和统计上的失效概率降低。这是对 iteration 2 攻击的直接回应。扣 10 分：(1) "统计上降低了同时失效的概率"这个统计论断没有数据支撑——没有引用任何关于 LLM 对 prompt 不同位置指令遵从率的研究。如果有实验数据（如"agent 对 Step 内部 IMPORTANT 标记的遵从率为 X%，对 Hard Constraints 的遵从率为 Y%，两者独立时联合失效概率为 Z%"），论证会更有力；(2) 类比"具体条文 vs 宪法原则"是修辞而非论证——法律体系有司法机构强制执行，prompt 指令没有。

---

### 4. Requirements Completeness: 96/110

**Scenario coverage (35/40):**
五个场景 + 四种退化情况覆盖全面：
- 场景 1：coding.* 任务执行（完整流程）
- 场景 2：doc.* 任务执行（三个子场景：从 tech-design 生成文档 / 修改现有文档 / 仅涉及 proposal）
- 场景 3：无 AC 的任务
- 场景 4：quick-tasks 生成（无 tech-design.md）——从 proposal.md 提取 section
- 场景 5：breakdown-tasks 生成（有 tech-design.md）——精确 section 提取

Edge Cases 覆盖了四种退化：Reference Files 为空、section 不存在、文件不存在、全部失效。

扣 5 分：仍然缺少 **clean-code / coding-fix 等非特性任务** 的场景描述。提案在 In Scope 中说审计全部 19 个模板确定哪些需要强化，但 Requirements Analysis 只分析了 coding.* 和 doc.* 任务类型。如果 clean-code 模板被判定需要强化（满足审计标准 (a)(b)(c)），它的 Reference Files 策略是什么？bug 修复任务的 Reference Files 可能来自 issue 描述而非 tech-design.md。

**Non-functional requirements (33/40):**
iteration 3 新增了"新旧格式兼容性"——这是 iteration 2 blindspot F2 的直接修复。三条 NFR（prompt 不过长、不改变执行流程结构、新旧格式兼容）覆盖了主要关注点。扣 7 分：(1) 性能影响：每个任务 Step 1 增加 2-5 个文件读取操作的开销仍未评估——提案只在 NFR 中说"每个任务引用 2-5 个 section，而非整个文件"，但未讨论这些读取操作对 agent 执行时间的影响；(2) 输出体积增长：新增加载确认 + AC 验收报告 + 降级提示的总量仍末量化。假设一个任务有 5 条 AC，验收报告增加约 5 行；加上加载确认 1 行和可能的降级提示 1-2 行，总共约 7-8 行新增输出。对于通常 50-200 行的 agent 输出，这个增量约 4-15%，影响不大但值得声明。

**Constraints & dependencies (28/30):**
三条约束覆盖了 forge-distribution.md 遵循、embed.FS 编译依赖、两层防护同步更新。比 iteration 2 更精确地描述了编译依赖："coding.* 模板通过 embed.FS 嵌入二进制，修改 forge-cli/pkg/prompt/data/*.md 后必须执行 go build 才能生效；task-executor.md 为即时生效文件，不需要编译"。扣 2 分：go build 的验证步骤在 Constraints 中提及（"必须确保 go build 在部署前完成"），但具体的验证命令（forge prompt get-by-task-id）只在 Success Criteria 中出现，Constraints section 本身应该引用这个验证方法。

---

### 5. Solution Creativity: 48/100

**Novelty over industry baseline (12/40):**
提案本身声明"非创新性改进，而是将已验证的工程实践制度化"。核心方案（在 prompt 中加强调标记和检查步骤）是 prompt engineering 的标准实践。iteration 3 未在此维度引入任何新元素。"两层防护"的诚实声明实际上降低了 novelty 评分——因为它承认两层是同一机制的重复。但"制度化"（将教训文档的分析结果系统性地嵌入执行流程）本身有一定的工程价值。扣 28 分。

**Cross-domain inspiration (18/35):**
Level 4 分析（LLM agent 天然倾向局部一致性而非全局一致性）来自认知科学的系统思维。Priority Rules 的三级优先级（Hard Rules > Reference Files > 现有代码）类比法律体系层级原则（iteration 3 明确使用了"具体条文 vs 宪法原则"的类比）。AC 逐条验收借鉴了软件测试的 checklist 验证模式。扣 17 分：仍然缺少编译器类型检查、数据库约束、CI/CD gate check 等成熟领域"规范强制执行"方案的借鉴——提案完全依赖 agent 的指令遵从，没有从这些领域中获取任何灵感。iteration 3 的诚实声明实际上承认了这一局限。

**Simplicity of insight (18/25):**
核心洞察简洁——"在执行流程中嵌入规范权威性锚点"。iteration 3 的实施方案虽然比 iteration 2 更完整（增加了精确模板文本、提取逻辑、兼容性声明），但这些增加的内容是对实施细节的完善而非概念的复杂化。三层结构（模板层 + Agent 层 + 任务生成层）的职责分离清晰。扣 7 分：Edge Cases & Degradation（4 种场景）+ Priority Rules（3 级优先级 + 冲突行为）+ Rollback Plan（3 层）+ Assumptions Challenged 表格的认知复杂度较高，对于一个声称"修复成本极低"的提案来说略显沉重。

---

### 6. Feasibility: 85/100

**Technical feasibility (36/40):**
iteration 3 解决了 iteration 2 的所有技术障碍：
1. 时序问题——在模板 Step 1 内部，读取任务文件之后（Resolved in iteration 2）
2. AC 验收插入位置——精确定义为 Step 3 "Verify & Finalize" 的第一个子步骤
3. embed.FS 编译依赖——已声明且有验证方案
4. quick-tasks 无 tech-design.md——场景 4 专门处理
5. breakdown-tasks 提取逻辑——四步提取决策逻辑

扣 4 分：(1) `<IMPORTANT>` vs `<EXTREMELY-IMPORTANT>` 分层标记策略的实际效果未验证——提案声称使用 `<IMPORTANT>` 是为了避免标记稀释，但没有实验数据或文献引用证明 `<IMPORTANT>` 在 `<EXTREMELY-IMPORTANT>` 已存在时能保持独立影响力；(2) breakdown-tasks 的四步提取逻辑中，"在 tech-design.md 中搜索提及这些文件路径的 section" 依赖于 tech-design.md 中是否显式提及文件路径——如果设计文档没有引用具体文件路径，这一步会失效。

**Resource & timeline feasibility (25/30):**
"4-6 个文件修改 + go build，1-2 小时完成"。iteration 3 的 Scope 更明确：审计 19 个模板 + 修改需要强化的模板 + 修改 task-executor.md + 修改 quick-tasks SKILL.md + 修改 breakdown-tasks SKILL + go build。扣 5 分：审计 19 个模板本身可能需要 30-60 分钟（需要逐一读取并判断是否满足审计标准的三条条件）。如果审计后发现 8-10 个模板需要强化（而非"4-6 个文件修改"的估计），修改工作量可能翻倍。timeline 估计的置信区间较宽。

**Dependency readiness (24/30):**
无外部依赖。内部依赖：embed.FS 编译 + forge prompt get-by-task-id 验证命令。扣 6 分：`forge prompt get-by-task-id` 命令是否会原样输出模板中的 `<IMPORTANT>` 声明文本？如果该命令对 prompt 有截断、格式化或去重处理，验证可能不准确。提案假设命令原样输出但未验证这个假设。此外，Next Steps 中提到"预实验：先用 1 个 coding 任务手动在 prompt 中注入 Reference Files 声明"——如果预实验失败，整个方案的可行性需要重新评估，这意味着 feasibility 部分依赖于一个未执行的实验结果。

---

### 7. Scope Definition: 72/80

**In-scope items are concrete (28/30):**
七项 In Scope 条目，每项指定了具体文件路径和修改内容：
1. coding.* 模板 Step 1 + Self-Check 插入
2. task-executor.md Hard Constraints 兜底
3. 审计全部 19 个模板
4. 审计标准（三条判定条件）
5. quick-tasks SKILL.md（明确指定"不是 templates/"）
6. breakdown-tasks SKILL.md（含四步提取逻辑）
7. 格式统一声明

扣 2 分：审计标准 (b) "模板的 Step 1 包含'读取任务文件'步骤"——如何判断一个模板是否包含此步骤？是关键词搜索还是语义判断？审计标准 (c) "涉及需要对照规范执行的实现/修改任务"——"需要对照规范执行"的判定标准是什么？这些细节在实施时可能产生分歧。

**Out-of-scope explicitly listed (22/25):**
四项 Out of Scope：修改 Go 代码、修改 index.json schema、添加 hooks 或运行时验证、改变 forge prompt get-by-task-id 合成逻辑。清晰。iteration 3 还通过 Innovation Highlights 的诚实声明明确了"真正的第二层防护（hook 拦截）属于 Out of Scope"。扣 3 分：新旧格式迁移是否属于 Out of Scope？NFR 说"不要求迁移已有任务文件"，但未在 Out of Scope 中显式列出"迁移旧格式 Reference Files"。这是一个小遗漏。

**Scope is bounded (22/25):**
没有明确的完成时间或分阶段交付计划。Next Steps 仍然只说"预实验 → Proceed to /quick-tasks"。但 iteration 3 新增的预实验步骤实际上定义了一个 phase gate："若预实验显示 agent 仍以代码为准，则需要重新评估方案"。Success Criteria 第 9 条（3-5 个任务验证）是另一个自然的检查点。扣 3 分：虽然有隐式的 phase gate（预实验 → 实施 → 验证），但没有显式地将其定义为分阶段交付策略。建议将 Next Steps 扩展为"Phase 1: 预实验 → Phase 2: 模板修改 + 审计 → Phase 3: 验证"。

---

### 8. Risk Assessment: 80/90

**Risks identified (27/30):**
四个风险：忽略 `<IMPORTANT>` 标记（M/H）、Reference Files 引用过时（M/M）、breakdown-tasks 填充不完整（L/M）、更新不同步（M/H）。覆盖了主要技术风险。扣 3 分：仍然缺少以下风险：
1. **预实验失败风险**：Next Steps 中提到"若预实验显示 agent 仍以代码为准，则需要重新评估方案"——这意味着方案的核心假设可能不成立。这是一个 H/H 风险，但未出现在 Risk 表中。
2. **AC 验收导致 agent 行为扭曲**：agent 可能为了通过 AC 而只实现 AC 中明确列出的内容，忽视 AC 未覆盖的质量维度。这在 iteration 2 的 blindspot-4 中已指出但未纳入 Risk 表。

**Likelihood + impact rated (26/30):**
四个风险的 Likelihood/Impact 评估总体合理。扣 4 分：(1) 风险 1（忽略 `<IMPORTANT>` 标记）Likelihood M——提案选择了 `<IMPORTANT>` 而非 `<EXTREMELY-IMPORTANT>` 来避免标记稀释，但如果 `<IMPORTANT>` 的遵从率本身就低于 `<EXTREMELY-IMPORTANT>`，选择 IMPORTANT 可能反而增加了被忽略的概率。这个 trade-off 未在 Likelihood 评估中讨论；(2) 风险 3（填充不完整）Likelihood L——breakdown-tasks 的四步提取逻辑中步骤 (2) "在 tech-design.md 中搜索提及文件路径的 section" 的成功率取决于设计文档的编写风格，这个 Likelihood 评估可能过于乐观。

**Mitigations are actionable (27/30):**
四个风险的 Mitigation 均可操作。特别是：
- 风险 1：三层防御（Step 1 + Self-Check + Hard Constraints）+ 分层标记策略
- 风险 2：引用 section 标题而非行号 + 降级行为
- 风险 3：SKILL.md 中加 checklist
- 风险 4：Success Criteria 中加入部署一致性验证 + go build 后检查

扣 3 分：风险 1 的核心缓解是"Hard Constraints 兜底"，但 iteration 3 的诚实声明已承认"两层本质上是同一种机制"——Mitigation 应该更新以反映这一认知。更准确的 Mitigation 应该是："预实验验证 agent 对 IMPORTANT 标记的遵从率；如果遵从率不足，升级为 EXTREMELY-IMPORTANT 或考虑 hook 机制"。

---

### 9. Success Criteria: 75/80

**Criteria are measurable and testable (50/55):**
九条成功标准大部分可验证：
1. 检查模板 Step 1 内容——可验证（grep `<IMPORTANT>` 声明文本）
2. 检查 task-executor.md Hard Constraints——可验证
3. 审计报告存在——可验证
4. 需强化的模板包含声明——可验证
5. quick-tasks 生成任务的 Reference Files 格式——可验证（iteration 3 的括号澄清消除了歧义）
6. breakdown-tasks SKILL.md 内容——可验证
7. 格式统一——可验证
8. 部署一致性（go build + forge prompt get-by-task-id）——可验证
9. 行为验证（3-5 个任务，每个输出含确认 + 验收报告 + 无规范偏离）——可验证

扣 5 分：(1) SC 9 "无规范偏离（人工 review 对照 Reference Files 确认）"——"人工 review"的判定标准是什么？谁来做 review？review 者如何判断"无偏离"？如果 review 者与提案者是同一人，存在确认偏差。建议定义为"由至少一位非提案者进行独立 review，逐条对照 Reference Files 检查实现一致性"；(2) SC 9 的"3-5 个任务"给出了范围而非确定值——建议改为"至少 3 个，推荐 5 个"以消除歧义。

**Coverage is complete (25/25):**
九条成功标准覆盖了所有 In Scope 条目：模板层（SC 1, 4）、Agent 层（SC 2）、审计（SC 3）、quick-tasks（SC 5）、breakdown-tasks（SC 6）、格式（SC 7）、部署（SC 8）、行为验证（SC 9）。无遗漏。

---

### 10. Logical Consistency: 82/90

**Solution addresses the stated problem (32/35):**
问题定义是"Agent 以代码为参照而非以规范为参照"。三层解决方案直接针对：
1. 模板层——在执行流程关键节点强制加载 Reference Files
2. Agent 层——Hard Constraints 兜底
3. 任务生成层——提高 Reference Files 生成质量

逻辑链：问题根因是"模板 Step 1 未要求读 Reference Files"（Level 0）→ 模板层直接修复 Level 0 → Agent 层防止模板修复遗漏 → 任务生成层防止 Reference Files 本身质量不足。扣 3 分：核心假设"agent 加载规范后就会以规范为参照"仍未完全验证——虽然 Priority Rules 声明了规范优先级，但 agent 完全可以加载规范后仍以代码为准。Next Steps 的预实验是正确的验证策略，但在预实验完成前，"解决方案是否真的解决核心问题"仍有不确定性。

**Scope <-> Solution <-> Success Criteria aligned (26/30):**
七项 In Scope 与九条 Success Criteria 的对齐关系清晰（见 iteration 2 的映射，iteration 3 保持了同样的结构并修复了 SC 5 的歧义）。扣 4 分：(1) Scope 第 6 条"breakdown-tasks 提取逻辑"标记为"待实施时细化"——这意味着 Scope 的定义在实施时可能发生变化，与 Success Criteria 的确定性要求存在轻微张力；(2) Scope 的审计（第 3 条）可能发现更多需要修改的模板，导致实际 Scope 大于当前定义——提案没有为这种 Scope 扩展预留灵活性。

**Requirements <-> Solution coherent (24/25):**
五个场景与三层解决方案的映射基本完整：
- 场景 1 → 模板层（coding.* 模板 + Agent 层兜底）
- 场景 2 → 模板层（doc.* 模板适配）+ 三种子场景的 Reference Files 策略
- 场景 3 → 场景定义了降级行为
- 场景 4 → 任务生成层（quick-tasks SKILL.md）
- 场景 5 → 任务生成层（breakdown-tasks SKILL.md + 四步提取逻辑）

Edge Cases 与降级行为与 Priority Rules 一致。扣 1 分：场景 2 的三个子场景 (a)(b)(c) 在 Success Criteria 中没有对应的独立验证项——如果 doc.* 模板的 Reference Files 策略与 coding.* 不同，应该有针对 doc.* 任务的行为验证。

---

## Cross-Dimension Coherence

1. **Urgency 的量化推算与 Solution Creativity 的"非创新"声明一致**——提案正确定位为"低成本制度化改进"而非"创新突破"。ROI 计算基于实际历史数据而非假设收益。维度 1 和 5 之间无张力。

2. **Innovation Highlights 的诚实声明与 Risk Assessment 存在轻微不一致**——提案承认两层防护是同一机制的重复，但 Risk 表中仍然使用"两层防护"的表述（如风险 1 的 Mitigation "task-executor.md Hard Constraints 兜底"），暗示它们是独立防线。建议在 Risk 表中也承认"重复强调"而非"两层防护"的性质。

3. **Next Steps 的预实验与 Feasibility 的确定性存在张力**——预实验的核心目的是验证"agent 是否确实以规范为参照"，这直接决定了方案的可行性。如果预实验失败，整个提案需要重新评估。但 Feasibility Assessment 给出了"技术风险低"的结论——这个结论应该在预实验成功之后才能得出。

4. **Scope 审计标准与 Requirements 场景覆盖存在缺口**——审计标准 (a)(b)(c) 用于判断模板是否需要强化，但 Requirements Analysis 只详细分析了 coding.* 和 doc.* 任务类型。如果审计发现 clean-code 或 coding-fix 模板需要强化，Requirements 中没有对应的场景描述和 Reference Files 策略。

---

## Blindspot Hunt

**[blindspot-1] 预实验失败的回退方案缺失**

引用 Next Steps："若预实验显示 agent 仍以代码为准，则需要重新评估方案（如考虑 Structured Output 或 hook 机制）"。

预实验是方案可行性的关键验证点，但"重新评估方案"不是一个可操作的回退方案。具体来说：
- 如果 agent 加载了 Reference Files 但仍以代码为准，说明 prompt 指令层面的"权威性声明"不足以改变 agent 行为。此时整个提案的防护机制失效。
- "考虑 Structured Output 或 hook 机制"被列为 Out of Scope，但预实验失败意味着需要 In Scope 这些方案。
- 建议在提案中增加"预实验失败的决策树"：失败后是否立即升级为 hook 机制？还是先尝试更强的 prompt 策略（如 `<EXTREMELY-IMPORTANT>` + 结构化输出格式要求）？

**[blindspot-2] AC 验收模板文本缺少 doc.* 任务适配**

引用解决方案中的 Self-Check AC 验收模板文本："Before performing other verification checks, validate against each Acceptance Criteria item from the task file"。

场景 2 定义了 doc.* 任务的三个子场景，每个子场景有不同的 Reference Files 来源。但 AC 验收模板文本是通用的——没有区分 coding.* 和 doc.* 的验收重点。场景 2 说"验收重点是文档结构合规而非路径命名"，但 AC 验收模板中没有任何关于"文档结构合规"的特殊指令。如果 doc.* 任务的 AC 中没有显式包含"文档结构合规"的条目，AC 验收模板不会覆盖这个差异。

**[blindspot-3] `<IMPORTANT>` 标记在模板中的实际效果未经验证**

引用解决方案："使用 `<IMPORTANT>` 而非 `<EXTREMELY-IMPORTANT>` 以避免标记稀释"。

这是一个关键的设计决策——选择 `<IMPORTANT>` 而非 `<EXTREMELY-IMPORTANT>`。但提案没有提供任何证据证明：
1. `<IMPORTANT>` 在已有 `<EXTREMELY-IMPORTANT>` 块的模板中不会被"稀释"——为什么 `<IMPORTANT>` 不会被视为"不那么重要"？
2. `<IMPORTANT>` 的遵从率是否足够高——如果 agent 将 `<IMPORTANT>` 视为"可选"而非"必须"，防护效果会大打折扣。
3. 预实验（Next Steps）是否覆盖了这个标记选择的有效性验证——如果预实验只验证"Reference Files 声明是否改变 agent 行为"而不区分标记类型，无法回答这个问题。

---

## Summary

SCORE: 750/1000
DIMENSIONS:
  Problem Definition: 102/110
  Solution Clarity: 112/120
  Industry Benchmarking: 98/120
  Requirements Completeness: 96/110
  Solution Creativity: 48/100
  Feasibility: 85/100
  Scope Definition: 72/80
  Risk Assessment: 80/90
  Success Criteria: 75/80
  Logical Consistency: 82/90
ATTACKS:
1. [Industry Benchmarking]: 两层防护的"统计上降低失效概率"缺少数据支撑——引用"统计上降低了同时失效的概率"但未引用任何 LLM prompt 遵从率研究——需要提供实验数据或文献引用证明不同位置的 prompt 指令具有独立的失效概率分布。
2. [Solution Creativity]: 完全依赖 prompt 指令遵从，未借鉴编译器类型检查/CI gate check 等成熟领域的"规范强制执行"自动化方案——引用 Innovation Highlights "非创新性改进"——但"制度化"本身可以从其他领域引入更强的约束机制，而非仅停留在 prompt engineering 的范畴。
3. [Feasibility]: 预实验是方案可行性的关键 gate 但在 Feasibility Assessment 中未体现——引用 Next Steps "若预实验显示 agent 仍以代码为准，则需要重新评估方案"——如果预实验失败，Feasibility 的"技术风险低"结论不成立，建议将 Feasibility 评估分为"预实验前"和"预实验后"两个阶段。
4. [Risk Assessment]: 预实验失败风险未列入 Risk 表——引用 Next Steps "若预实验显示 agent 仍以代码为准"——这是一个可能导致整个方案推翻的高影响风险，Likelihood 未知但 Impact 极高，应纳入 Risk 表。
5. [Success Criteria]: SC 9 "人工 review 对照 Reference Files 确认"缺少独立性和判定标准——引用"无规范偏离（人工 review 对照 Reference Files 确认）"——谁来做 review？review 判定标准是什么？建议定义为"由至少一位非提案者进行独立 review"。
6. [Solution Clarity]: `<IMPORTANT>` vs `<EXTREMELY-IMPORTANT>` 的分层标记策略缺少有效性验证——引用"使用 `<IMPORTANT>` 而非 `<EXTREMELY-IMPORTANT>` 以避免标记稀释"——但 agent 可能将 `<IMPORTANT>` 视为比 `<EXTREMELY-IMPORTANT>` 更弱的信号，反而降低遵从率。需要在预实验中同时测试两种标记的效果差异。
7. [Requirements Completeness]: clean-code / coding-fix 等非特性任务的场景缺失——引用审计标准 (a)(b)(c) 暗示这些模板可能需要强化，但 Requirements Analysis 中无对应场景——需要补充非特性任务的 Reference Files 策略（如 bug 报告、issue 描述作为 Reference Files 来源）。
8. [Logical Consistency]: Risk 表仍使用"两层防护"表述但诚实声明已承认是"同一机制的重复强调"——引用风险 1 Mitigation "task-executor.md Hard Constraints 兜底"——诚实声明说"本质上是同一种机制"，Risk 表应改为"重复强调"而非暗示独立防线。
