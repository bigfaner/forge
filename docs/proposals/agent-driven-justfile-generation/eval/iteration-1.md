# Proposal Evaluation — Iteration 1

**Proposal**: Agent-Driven Justfile Generation
**Date**: 2026-06-08
**Scorer**: CTO Adversarial Review
**Rubric**: proposal.md (1000-point scale)

---

## Phase 1: Reasoning Audit

### 1. Problem → Solution Trace

**Problem**: 语言模板是灵活性的主要瓶颈——非标准项目生成后必须手动编辑，模板默认值反而是障碍。

**Solution**: 移除所有语言模板，以 surface + Convention + LLM 知识为唯一驱动源。

**Verdict**: 解决方案直接解决了问题陈述。移除模板消除了模板默认值与非标准项目之间的冲突。链路完整。

### 2. Solution → Evidence Trace

**Claim**: agent 可直接从 surface rule + 项目检测生成命令，无需模板中间层。

**Evidence**: 提案引用了 surface rule 文件已定义 recipe 契约（5个）、Convention 机制已提供框架知识注入通道。SKILL.md 确认 surface rules 和 Convention 加载机制已存在。

**Gap**: 提案声称"Agent（Claude/GPT）已具备主流语言构建命令的知识"但未提供实证。这是可行性判断而非证据。不过对于 LLM 能力而言，这是一个合理的行业共识假设。

**Verdict**: 证据链基本完整，可行性论证依赖隐含的 LLM 能力假设。

### 3. Evidence → Success Criteria Trace

SC 要求通过 verification step（dry-run + actual execution）验证四种语言（Go/Node/Python/Rust）的 recipe。这确实测试了核心价值主张——无模板生成是否能产出正确命令。

**Gap**: NFR 要求"结构级一致性"（多次运行 recipe 名称、分组、边界标记相同），但 SC 没有对应的多次运行一致性测试条目。SC[7] 要求"与当前输出一致"但只是单次结构比对，非多次运行一致性验证。

**Verdict**: SC 覆盖了主要功能路径，但遗漏了 NFR 中声明的一致性可测试性。

### 4. Self-Contradiction Check + SC Consistency Deep-Dive

**Cross-cluster analysis**:

- **SC[1]** (替换 TODO stub) ↔ **InScope[5]** (简化 surface rule 文件): 一致，SC[1] 是 InScope[5] 的可测试版本。
- **SC[2]** (删除模板) ↔ **InScope[1]** (删除模板): 完全一致。
- **SC[3]** (server-lifecycle.md) ↔ **InScope[3]** (server-lifecycle.md): 一致。
- **SC[4]** (空 surfaces 提示) ↔ **InScope[4]** (SKILL.md 重写): 一致，SKILL.md 重写应包含此行为。
- **SC[5]** (四种语言 verification) ↔ **InScope[3]** (SKILL.md 重写): 一致，verification step 是 SKILL.md 流程的一部分。
- **SC[6]** (混合语言多 surface) ↔ **InScope[3]** + **InScope[5]**: 一致。
- **SC[7]** (结构一致性) ↔ **InScope[4]** (向后兼容 NFR): 一致。

**Potential tension**: NFR 声称"相同项目多次运行生成的 justfile 结构级属性一致"与 LLM 生成的不确定性本质之间存在张力。提案承认了这一点（Risk[1]: LLM 生成命令不一致），缓解措施是 surface rule 契约 + Convention 约束。但这不是完全消除不确定性的保证——命令体允许因 LLM 变化而不同，且"语义契约"层面的一致性缺乏可验证的定义。

**No direct contradiction found** within SC set itself.

---

## Phase 2: Rubric Scoring

### 1. Problem Definition (110 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Problem stated clearly | 38/40 | 问题陈述清晰：模板硬编码与非标准项目的冲突。"模板默认值反而是障碍"是一个具体的痛点。轻微扣分：未量化"大量覆盖"——是覆盖多少行/多少 recipe？ |
| Evidence provided | 35/40 | 提供了三个具体证据点（6 模板 ~800 行、mixed.just 230 行、三层生成流程冲突）。经实际验证，模板总行数 811、mixed.just 229 行，数据准确。扣分：缺乏用户反馈或实际案例说明"大量覆盖"的频率和影响程度。 |
| Urgency justified | 28/30 | "移除模板后，agent 直接针对实际项目结构生成正确命令，消除这层摩擦"——说明了不做会怎样。轻微扣分：未量化延迟成本（多少项目受影响？每月多少时间花在手动编辑上？）。 |

**Subtotal: 101/110**

### 2. Solution Clarity (120 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Approach is concrete | 37/40 | 四步方案（检测→枚举→填充→复杂模式提取）清晰可执行。扣分：Step 3"填充内容"依赖 LLM 知识，未说明当 agent 知识与项目实际不匹配时的具体策略——仅靠 verification step 是事后纠正，不是生成策略。 |
| User-facing behavior described | 40/45 | 用户视角：运行 init-justfile → 得到正确 justfile。关键行为变化（--type 参数移除、空 surfaces 提示）都有说明。扣分：未明确描述用户从"当前流程"到"新流程"的迁移体验——已有的 justfile 怎么办？是否会自动重新生成？boundary marker 机制是否保留？ |
| Technical direction clear | 32/35 | 技术方向明确：surface rule 驱动 + Convention 注入 + LLM 知识填充。扣分：server-lifecycle.md 的"multi-service 场景指导"描述过于模糊——"端口感知启动"和"启动顺序依赖声明"是 rule 文件的指导还是可执行代码？ |

**Subtotal: 109/120**

### 3. Industry Benchmarking (120 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Industry solutions referenced | 28/40 | 引用了 cookiecutter、yeoman、hygen、Earthly、Taskfile。扣分：引用过于简略——没有说明这些工具如何解决类似问题，也没有链接到具体文档或模式。Earthly 和 Taskfile 的 DSL 抽象方式与本项目有何可比性？这些工具的"用户编写"模式与 agent-driven 模式的对比缺乏深度。 |
| At least 3 meaningful alternatives | 24/30 | 提供了三个：do nothing、parameterized templates、agent-driven。Do nothing 是有效基线。扣分：parameterized templates 被描述为"参数爆炸"但缺乏论证——hygen 的参数化模板实际复杂度如何？是否有行业案例证明参数化不可行？此 alternative 的 reject 理由"治标不治本"是定性判断而非实证。 |
| Honest trade-off comparison | 20/25 | 诚实承认 LLM 生成有轻微不确定性。扣分：对比表的 Cons 列过于简化——agent-driven 的 Cons 仅"LLM 生成有轻微不确定性"，但未列出生成延迟（LLM 调用 vs 模板读取）、token 成本、离线场景不可用等实际 trade-off。 |
| Chosen approach justified against benchmarks | 20/25 | "彻底解决灵活性问题"作为选择理由成立。扣分：未讨论为何不采用混合方案（如保留模板作为 fallback + agent 增强），直接跳到全 agent-driven 的激进方案。 |

**Subtotal: 92/120**

### 4. Requirements Completeness (110 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Scenario coverage | 35/40 | 覆盖了单 surface、多 surface、混合语言、非标准项目、空 surfaces。扣分：缺少以下场景：(a) 现有 justfile 已有 forge boundary marker 但部分 recipe 已被用户 `# user-customized` 标记——重新生成时如何处理？(b) 同一 surface 目录下检测到多种语言 marker（如既有 go.mod 又有 package.json）的处理策略。 |
| Non-functional requirements | 33/40 | 一致性（结构级属性相同）和向后兼容性是有效的 NFR。扣分：未提及性能 NFR——agent 驱动生成（需要 LLM 推理）vs 模板读取的延迟差异。未提及离线/无 LLM 场景下的退化策略。一致性要求中"语义契约"的验证手段未定义。 |
| Constraints & dependencies | 28/30 | 列出了四个依赖（forge surfaces、surface rules、Convention、just >= 1.50.0），均标注已就绪。扣分：未提及 LLM 能力作为约束——如果 agent 运行在弱模型上（如本地小模型），生成质量是否可接受？ |

**Subtotal: 96/110**

### 5. Solution Creativity (100 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Novelty over industry baseline | 32/40 | "LLM 驱动 + 结构约束"替代模板确实超越了传统脚手架工具的模式。Innovation Highlights 准确描述了差异化："Surface rule 文件定义'要生成什么'（契约），agent 决定'怎么生成'（内容）"。扣分：在 AI coding assistant 普及的 2026 年，用 LLM 替代模板算是一种"自然的演进"而非高度创新的飞跃。 |
| Cross-domain inspiration | 25/35 | 从 Forge 自身的 surface-first 设计理念中汲取灵感。扣分：未借鉴其他领域——例如编译器的 IR 中间表示（模板 → IR → 目标代码 vs 当前的模板直接到输出）、或者配置管理的 desired-state 模型（声明期望状态而非命令式步骤）。 |
| Simplicity of insight | 22/25 | "删掉模板，让 agent 直接生成"确实是"为什么不早这么做"级别的简洁洞察。扣分：这个简洁性的代价是引入了对 LLM 的强依赖，提案对此的分析不够充分。 |

**Subtotal: 79/100**

### 6. Feasibility (100 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Technical feasibility | 35/40 | 技术方案可行——SKILL.md 确认 surface rules 和 Convention 机制已存在，agent 已具备语言知识。扣分：server-lifecycle 的边界条件（PID 文件残留、进程被外部 kill、Windows `\r` 污染）被列为"已知限制"但缓解措施"提供可直接使用的可执行 bash 代码片段"并未实际消除风险——agent 可能仍然选择从头生成而非复用代码片段。 |
| Resource & timeline feasibility | 25/30 | 拆分为两个 phase 合理。13 个文件的改动范围明确。扣分：未估算每个 phase 的工时。Phase 1 的"核心架构迁移"涉及 SKILL.md 重写和模板删除——这是整个 skill 的核心逻辑重写，实际复杂度可能被低估。 |
| Dependency readiness | 28/30 | 所有依赖已就绪且已标注。轻微扣分：Convention 文件的"已有机制"和"已有知识"是不同层级——Convention 机制存在不代表每个语言/框架的 Convention 文件都已完善。 |

**Subtotal: 88/100**

### 7. Scope Definition (80 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| In-scope items are concrete | 26/30 | 六个 in-scope 项都是具体可交付物（删除文件、新增文件、重写文件）。扣分：InScope[5]"Post-generation 一致性验证"描述了功能但未说明这是一个独立步骤还是集成到 Step 4 verification 中的子步骤。 |
| Out-of-scope explicitly listed | 22/25 | 四个 out-of-scope 项清晰列出。扣分：未排除"init-justfile 的 --type 参数在非 surface 项目中的替代方案"——移除 --type 后，纯语言项目（无 surface）如何指定项目类型？依赖 agent 自动检测的准确率是多少？ |
| Scope is bounded | 22/25 | 两个 phase 的划分提供了时间边界。回退策略（git tag）提供了安全网。扣分：Phase 2 的"一致性验证"是开放性工作——验证到什么程度算完成？ |

**Subtotal: 70/80**

### 8. Risk Assessment (90 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Risks identified | 25/30 | 四个风险覆盖了主要场景：LLM 不一致、server lifecycle 复杂度、冷启动、罕见语言。扣分：缺少以下风险：(a) 迁移期间用户从旧模板生成的 justfile 与新 agent-driven 生成的不兼容；(b) agent 生成延迟（LLM 推理时间 vs 模板读取）对用户体验的影响。 |
| Likelihood + impact rated | 25/30 | 评级总体合理——LLM 不一致性为 M/L，server lifecycle 为 M/M。扣分：冷启动（无 Convention）评为 L/M 但实际很可能是常见场景（新项目初始化），likelihood 可能被低估。罕见语言评为 L/L 合理。 |
| Mitigations are actionable | 24/30 | 缓解措施具体：verification step、独立 rule 文件、error stub 回退。server-lifecycle 的缓解（"提供可直接使用的可执行 bash 代码片段，agent 优先复用而非从头生成"）是可操作的。扣分：(a) "agent 优先复用"是一个 instruction，不是 enforceable 约束——如何确保 agent 实际复用？(b) verification step 无法覆盖 server lifecycle 边界条件（提案自己也承认了），但未提供替代验证手段。 |

**Subtotal: 74/90**

### 9. Success Criteria (80 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Criteria are measurable and testable | 25/30 | SC[2-4] 完全可测试（文件存在/不存在）。SC[5-6] 通过 verification step 可测试。SC[7] 的"结构一致"需要精确比对工具但概念可测试。扣分：SC[1]（"替换为 Recipe Generation Requirements section"）——什么算合格的 section？需要包含哪些内容？SC[7] 提到 `[linux]/[windows] 双平台`——如何验证双平台？在 Windows 上测试还是只检查语法？ |
| Coverage is complete | 20/25 | SC 覆盖了大部分 in-scope 项。扣分：(a) InScope[5]（Post-generation 一致性验证）无对应 SC——如何验证"比对 surface rule 文件中的 Recipe Invocation Contract 与实际生成的 recipe 名称/参数"是否正确执行？(b) NFR 中的"多次运行结构级一致性"无对应 SC。 |
| SC internal consistency | 22/25 | SC 集合内部无矛盾。SC[5] 和 SC[6] 是渐进关系（单语言 → 混合语言）。扣分：SC[7] 要求"与当前输出一致"但 InScope[1] 删除了生成当前输出的模板——如果当前输出的某些细节（如注释风格、空白行格式）是由模板硬编码的，agent 驱动生成是否应保持这些细节？此处的"一致"定义不够精确。 |

**Subtotal: 67/80**

### 10. Logical Consistency (90 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Solution addresses stated problem | 33/35 | 解决方案直接解决了模板作为灵活性瓶颈的问题。移除模板 → agent 直接针对项目结构生成 → 消除模板默认值与实际项目不匹配的摩擦。链路完整。扣分：问题中提到"三层生成流程中第一步（模板）和第三步（LLM 微调）经常冲突"——新方案消除了冲突，但新方案中 Convention 知识与 LLM 知识之间是否可能产生类似的冲突？ |
| Scope ↔ Solution ↔ SC aligned | 25/30 | 总体对齐良好。扣分：InScope[5]（Post-generation 一致性验证）在 SC 中无对应条目（如上所述）。Solution 描述的 server-lifecycle.md 包含 "multi-service 场景指导"但 SC[3] 只要求"PID 追踪、幂等启动、健康检查"——multi-service 指导部分无验证条目。 |
| Requirements ↔ Solution coherent | 22/25 | 需求到方案的映射清晰：5 个 key scenarios 对应 solution 的 4 个步骤。扣分：Scenario[4]（非标准项目）完全依赖 LLM 知识，但 Solution 中未说明 agent 如何区分"非标准但正确"和"错误检测"——例如项目同时存在 Makefile 和 package.json，agent 如何判断主构建系统？ |

**Subtotal: 80/90**

---

## Phase 3: Blindspot Hunt

### [blindspot] B1: LLM 调用成本和延迟未被评估

提案将"模板读取"替换为"agent 推理生成"，但从未分析这对用户体验的影响。模板读取是毫秒级操作，LLM 生成至少需要数秒。对于 `init-justfile` 这种用户主动触发的命令，延迟可能是可接受的，但应该在提案中明确承认并量化。

**Severity**: Medium — 影响用户决策但不是 blocker。

### [blindspot] B2: server-lifecycle.md 的"可直接使用的可执行 bash 代码片段"与 agent 自主性的矛盾

提案的核心价值是"agent 决定怎么生成"。但 server-lifecycle.md 提供了"可直接使用的代码片段，agent 优先复用"。这实际上是在用另一种形式的模板（bash 代码模板）替代被删除的语言模板。虽然范围更小（仅 server lifecycle），但逻辑上与"移除模板"的价值主张有张力。

提案应明确区分：server-lifecycle.md 是"参考实现"（agent 可以参考但不受约束）还是"强制模板"（agent 必须复用）。当前措辞暗示后者，这削弱了"零模板"的声明。

**Severity**: Medium — 设计哲学一致性问题。

### [blindspot] B3: 已有 justfile 的迁移路径不完整

Scope 列出了"Out of Scope: 其他 skill 更新（它们消费生成的 just recipe，不直接调用 init-justfile）"。但这意味着如果用户在迁移后重新运行 init-justfile，生成的 recipe 可能与之前模板生成的有细微差异（不同的命令体），而其他 skill 依赖的 recipe 名称和参数契约虽然不变，命令行为可能微妙变化。提案声称"向后兼容"但未讨论"重跑 init-justfile"对现有项目的影响。

**Severity**: Low — recipe 契约不变，行为变化在可接受范围。

### [blindspot] B4: 离线/无 LLM 场景的退化策略缺失

如果用户在无网络环境运行 init-justfile，或者 LLM 服务不可用，整个 skill 将无法工作。当前的模板方案不依赖 LLM 可用性（模板是本地文件）。这不是 risk 表中的任何一项。

**Severity**: Low — 对 Forge 的典型使用场景（开发者在线环境）不太可能发生，但作为工程考量应被提及。

---

## Bias Detection Report

- Annotated regions (`<!-- pre-revised -->`): 7 attack points / 8 annotated paragraphs = density 0.88
- Unannotated regions: 22 attack points / 34 unannotated paragraphs = density 0.65
- Ratio (annotated/unannotated): 1.35

**Interpretation**: Annotated regions show ~35% higher attack density. This is within acceptable range — pre-revised paragraphs addressed their targeted issues but introduced new surface area for attack (e.g., NFR consistency claims, server-lifecycle multi-service scope expansion). No systematic bias detected.

**`conflict-with-pre-revision` tags**: None. All attacks on pre-revised regions align with rubric criteria rather than contradicting the pre-revision direction.

---

## Score Summary

| # | Dimension | Score | Max |
|---|-----------|-------|-----|
| 1 | Problem Definition | 101 | 110 |
| 2 | Solution Clarity | 109 | 120 |
| 3 | Industry Benchmarking | 92 | 120 |
| 4 | Requirements Completeness | 96 | 110 |
| 5 | Solution Creativity | 79 | 100 |
| 6 | Feasibility | 88 | 100 |
| 7 | Scope Definition | 70 | 80 |
| 8 | Risk Assessment | 74 | 90 |
| 9 | Success Criteria | 67 | 80 |
| 10 | Logical Consistency | 80 | 90 |
| | **Total** | **856** | **1000** |

**Gate**: FAIL (target: 900, achieved: 856, gap: 44 points)

---

## Top Attacks (Ordered by Impact)

1. **Industry Benchmarking**: 替代方案论证不充分 — "parameterized templates" 被 reject 的理由是"参数爆炸"和"治标不治本"，但这是定性判断而非实证。hygen 的参数化机制在实际项目中是否真的导致"参数爆炸"？需要提供案例或引用。 — (Industry Benchmarking -24 pts)

2. **Success Criteria**: SC 缺少 NFR 一致性的可测试条目 — NFR 声称"相同项目多次运行生成的 justfile 结构级属性一致"，但 7 条 SC 中没有一条要求验证多次运行一致性。SC[7] 是单次结构比对，非多次运行验证。需要新增 SC：对同一项目连续运行 3 次 init-justfile，比较输出结构。 — (Success Criteria -20 pts)

3. **Requirements Completeness**: 离线/LLM 不可用场景未覆盖 — NFR 未提及 LLM 不可用时的退化策略。当前模板方案不依赖 LLM，新方案将 LLM 从增强工具变为必要组件，这是一个隐含的架构约束变化。需在 Constraints 中明确声明"本方案要求 LLM 可用"。 — (Requirements Completeness -14 pts)

4. **Solution Clarity**: 迁移体验未描述 — 提案描述了新方案的行为，但未说明已有 justfile 的用户如何从模板生成迁移到 agent 生成。boundary marker 机制是否保留？重新运行 init-justfile 时是否自动替换？ — (Solution Clarity -13 pts)

5. **Industry Benchmarking**: Trade-off 分析不完整 — agent-driven 的 Cons 仅"LLM 生成有轻微不确定性"，未包括生成延迟（LLM 推理 vs 模板读取的毫秒级操作）、token 成本、离线不可用。需要补充完整的 trade-off 列表。 — (Industry Benchmarking -12 pts)

6. **Risk Assessment**: 缺少迁移兼容性风险 — 风险表未包含"已有 justfile 用户重新生成后命令体微妙变化"的风险。提案声称向后兼容但仅保证结构兼容，不保证行为兼容。 — (Risk Assessment -10 pts)

7. **Success Criteria**: SC[1] 的验证标准模糊 — "替换为 Recipe Generation Requirements section"但未定义合格 section 的最低内容要求。 — (Success Criteria -8 pts)

8. **[blindspot] B2**: server-lifecycle.md 本质上是 server lifecycle 领域的"模板" — 提案声称"零模板维护"但 server-lifecycle.md 提供的"可直接使用的代码片段"在 server lifecycle 领域扮演了模板角色。需要明确其定位是参考实现还是强制模板。 — (Solution Creativity -7 pts)

---

## Recommendations for Next Iteration

1. **补充一致性 SC**（+15-20 pts potential）：新增"对同一项目运行 3 次 init-justfile，比较 recipe 名称/分组/边界标记/退出码语义是否完全相同"的 SC。

2. **深化行业对比**（+15-20 pts potential）：为 parameterized templates 提供具体的 reject 证据（案例或数据），补充 agent-driven 的完整 trade-off 列表（延迟、成本、离线）。

3. **明确 Constraints 中的 LLM 依赖**（+8-10 pts potential）：将 LLM 可用性作为显式约束声明，讨论退化策略。

4. **补充迁移指南**（+8-10 pts potential）：描述已有用户从模板生成迁移到 agent 生成的具体步骤和预期变化。

5. **量化 SC[1] 和 SC[7]**（+5-8 pts potential）：定义合格 Recipe Generation Requirements section 的最低内容，定义"结构一致"的精确比对维度。
