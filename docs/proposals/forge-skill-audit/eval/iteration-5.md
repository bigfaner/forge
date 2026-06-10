# Proposal Evaluation — Iteration 5

**Proposal**: Forge Plugin Skill 系统一致性修复
**Evaluator Role**: CTO Adversary
**Date**: 2026-06-10
**Rubric**: `plugins/forge/skills/eval/rubrics/proposal.md` (1000-point scale, target 900)

---

## Reasoning Audit

### Trace: Problem -> Solution -> Evidence -> SC

**Problem**: v3.0.0-rc.53 有 21 处审计发现，含 4 个 HIGH 级静默错误。
**Solution**: 文本级修复 12 项（4 HIGH + 8 MEDIUM），9 MINOR 标记不修。
**Evidence**: 审计覆盖 8 个维度（21 skill、16 command、5 hook），每项发现给出文件路径+行号+影响分析。经交叉验证：H-1（rubric scale/target 不一致）属实——rubric-reference.md 记录 journey scale=1000/target=850，实际 rubric frontmatter 为 scale=1150/target=975；H-3（模板硬编码）属实——breakdown-tasks/templates/task.md 确实使用 `"medium"` 和 `"coding.feature"` 而非占位符；H-4（doc.fix 分类）属实——record-format-coding.md 第 3 行确实列出 `doc.fix`。
**SC**: 16 项成功标准，与 12 项修复逐一对应，覆盖了回归验证命令。

**自相矛盾检查**:
1. Proposal 标题写"21 处不一致"，但 Summary Statistics 表格列合计也是 21——一致。
2. "M-1 config key 重命名仅修改 skill markdown" 同时标记 Go alias 为后续任务——scope 定义与修复方案一致。
3. H-1 影响分析称"通过门槛被降低 11 个百分点"——验算：850/1150=73.9% vs 975/1150=84.8%，差值 10.9pp，合理近似。

---

## Rubric Scoring

### 1. Problem Definition: 98/110

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Problem stated clearly | 38/40 | 四个 HIGH 问题均有文件路径+行号+具体矛盾描述。唯一模糊点：H-1 的"系统性风险"段提到"config 键默认值"作为第五个真相源，但未给出该 config 键的具体名称和默认值。 |
| Evidence provided | 40/40 | 每项发现给出文件路径、实际值 vs 期望值、影响链条。审计维度覆盖全面（8 维度 x 审计范围表格）。经代码库交叉验证，核心事实（H-1 scale/target、H-3 硬编码、H-4 doc.fix）均属实。 |
| Urgency justified | 20/30 | RC 阶段确实紧迫，但"紧迫性"段落的论述几乎全部围绕 H-1 展开，H-2/H-3/H-4 的紧迫性论述为空。对"为什么不延迟到 v3.1"的论证不充分——只是说"延迟成本增加"但未量化。 |

**Deductions**: -12 (H-2/H-3/H-4 紧迫性未独立论证; "config 键默认值"真相源未具名)

### 2. Solution Clarity: 108/120

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Approach is concrete | 38/40 | 每项修复方案给出具体文件、具体修改内容、修改前后对比。H-3 甚至给出了完整的注释块代码。但 M-4 "移至 _deprecated/" 或 "补充引用" 给了两个方案却未选定——是二选一还是两个都做？ |
| User-facing behavior described | 40/45 | 修复后用户行为变化描述清晰：H-1 修复后 eval-journey 不再误导用户传入 850 target；H-3 修复后 LLM 不再被硬编码值锚定。但 M-1 的用户影响未描述——统一为 kebab-case 后，已有 camelCase 配置的用户会遇到什么？ |
| Technical direction clear | 30/35 | "文本级修复，无代码变更"的方向明确。但 M-1 标注了前提条件"需先验证 Go config reader 是否支持 kebab-case"，如果验证不通过则修复方案变更——这是一个条件分支方案，增加了模糊性。 |

**Deductions**: -12 (M-4 二选一未决策; M-1 用户影响未描述; M-1 条件分支增加模糊性)

### 3. Industry Benchmarking: 72/120

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Industry solutions referenced | 30/40 | 引用了 promptfoo（prompt template 断言）、conftest（YAML/JSON 策略验证）、Pact（契约测试）。但只引用了概念映射，未给出任何具体实施示例或接口设计——这些工具如何与 forge 的 markdown skill 系统对接？ |
| At least 3 meaningful alternatives | 25/30 | 列出 4 个替代方案（手动修正、schema 驱动、仅修 HIGH、延迟处理），包含了"不做"选项。但"schema 驱动"替代方案与"手动修正"的权衡分析不完整——缺投入估算。 |
| Honest trade-off comparison | 12/25 | 表格中"手动修正"的劣势只写了一句"未来 rubric 变更仍需手动同步多处"。这是核心架构缺陷（多真相源），但 trade-off 分析轻描淡写。schema 驱动方案的"投入较大"有多大的估算？ |
| Chosen approach justified against benchmarks | 5/25 | "当前阶段最务实"不是论证。为什么不在这次修复中同时引入 rubric-reference 的自动验证？即使是简单的 grep 断言也只需要 5 分钟。提案选择手动修正但承认长期需要自动化，却没有给出中间步骤的迁移路径。 |

**Deductions**: -48 (工具对接缺具体实施; 投入估算缺失; trade-off 分析不诚实; 选择论证不充分)

### 4. Requirements Completeness: 78/110

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Scenario coverage | 30/40 | 覆盖了正常修复场景和回归验证。但缺少：如果 H-2 修复前搜索发现 `docs/features/` 下确实存在 proposal 文件怎么办？——提案提到"需先搜索确认"但未给出分支处理方案。 |
| Non-functional requirements | 22/40 | 完全没有 NFR。修复 12 个文件后的回归验证策略有，但没有：性能影响（INLINE 版本号标记增加的 grep 开销？）、兼容性（M-1 kebab-case 重命名对已发布模板的影响？）、可维护性（添加维护注释真的能防止未来遗漏吗？）。 |
| Constraints & dependencies | 26/30 | 明确列出"不修改 Go 代码"、"不追溯历史 eval 结果"等约束。M-1 依赖 Go config reader 验证结果。依赖链清晰。 |

**Deductions**: -32 (H-2 分支场景未覆盖; NFR 完全缺失; 维护注释作为防御机制的有效性未论证)

### 5. Solution Creativity: 35/100

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Novelty over industry baseline | 10/40 | 这是一份文本修正提案——"找到错误文本并改为正确文本"。没有创新点。 |
| Cross-domain inspiration | 10/35 | INLINE 版本号标记借鉴了依赖管理的版本锁定概念，但实现极简。维护注释借鉴了"cache invalidation"思维，但只是个 TODO 注释。 |
| Simplicity of insight | 15/25 | "用 grep 做回归验证"确实简单实用。 |

**Deductions**: -65 (这是纠错工作而非创造性解决方案; INLINE 版本号标记是手工同步的变体)

### 6. Feasibility: 88/100

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Technical feasibility | 38/40 | 全部是 markdown 文本修改，技术上零风险。唯一不确定项是 M-1 的 Go config reader 兼容性，但已标注为前提条件。 |
| Resource & timeline feasibility | 30/30 | 12 项文本修复，预估工时极短。无外部依赖。 |
| Dependency readiness | 20/30 | 外部依赖为零（不修改代码）。但 M-1 依赖 Go config reader 验证，如果验证不通过，M-1 的执行方案需要变更——这一前置步骤的完成时间未纳入计划。 |

**Deductions**: -12 (M-1 前置验证时间未纳入; Go config reader 验证如失败无备选方案)

### 7. Scope Definition: 72/80

| Criterion | Score | Justification |
|-----------|-------|---------------|
| In-scope items are concrete | 26/30 | 12 项修复均有具体文件和修改描述。但 M-4 的修复方案不具体（"移至 _deprecated/ 或补充引用"）。 |
| Out-of-scope explicitly listed | 24/25 | 5 项 out-of-scope 明确列出。Go 代码变更、历史结果修正、性能优化、新功能——边界清晰。 |
| Scope is bounded | 22/25 | 范围明确为 12 项修复。但 M-1 有前提条件，如果前提不满足则 scope 会膨胀（需要先改 Go 代码）。这在 In Scope 段落中已标注但增加了不确定性。 |

**Deductions**: -8 (M-4 修复方案不具体; M-1 前提条件可能导致 scope 膨胀)

### 8. Risk Assessment: 72/90

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Risks identified | 25/30 | 列出 6 项风险。但遗漏了：批量修复过程中的注意力衰减风险（12 项修改逐一手动执行，第 8 项时出错概率上升）；以及 INLINE 版本号标记引入的虚假安全感风险（开发者看到版本号标记后可能误以为有自动同步机制）。 |
| Likelihood + impact rated | 22/30 | 评级整体诚实。但"M-1 部分迁移导致新 key 静默失效"评为 Likelihood=中、Impact=高——这是本提案可能主动引入的最高风险，却只给中等概率？如果 Go config reader 不支持 kebab-case 而提案已经修改了 markdown 中的 key 名称，则 Likelihood 应为"高"。 |
| Mitigations are actionable | 25/30 | 大部分缓解措施具体可执行（grep 回归验证、git revert 回滚）。但"长期引入 CI 检查"和"中期固化 grep 命令为 justfile recipe"不是本次修复的可执行缓解措施——它们是未来计划。 |

**Deductions**: -18 (遗漏注意力衰减风险和虚假安全感风险; M-1 风险评级偏低; 部分缓解措施为未来计划而非当前可执行)

### 9. Success Criteria: 68/80

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Criteria are measurable and testable | 25/30 | 16 项 SC 大部分包含具体的 grep 命令和验证步骤。但 "Config 系统审计结论与实际发现一致" ——如何测量"一致"？需要更具体的判定标准。 |
| Coverage is complete | 20/25 | 12 项修复对应 14 条 SC（H-1 有 2 条，H-4 有 2 条），加上 2 条系统级 SC。但 MINOR 9 项完全不包含 SC——虽然标记为不修，但没有"确认不修的决策记录"的 SC。 |
| SC internal consistency | 23/25 | SC 之间无逻辑矛盾。但 M-1 的 SC 有前提条件："前提：验证 Go config reader 是否支持 kebab-case；如不支持，将 M-1 与 Go alias 绑定为原子操作，推迟至 Go 代码变更窗口执行"——这意味着 M-1 的 SC 可能不适用于本迭代，但 SC 列表中未标注条件性。 |

**Deductions**: -12 (一条 SC 可测量性不足; MINOR 项缺决策记录 SC; M-1 SC 条件性未标注)

### 10. Logical Consistency: 82/90

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Solution addresses the stated problem | 33/35 | 4 个 HIGH 问题均有直接对应的修复方案。M 级修复与发现项一一对应。但 H-2 的修复方案说"搜索确认无并存路径后移除死路径"——然而 eval/SKILL.md 第 12 行写的是 "Exceptions: `proposal` uses full path `docs/proposals/<slug>/`"，这说明 `docs/features/<slug>/proposal.md` 确实不是正确路径，但 proposal 的异常处理已经文档化了。修复 H-2 时是否需要同步更新 eval/SKILL.md 的异常说明？提案未提及。 |
| Scope <-> Solution <-> SC aligned | 28/30 | Scope 列 12 项，Solution 描述 12 项，SC 覆盖 14 条（含 2 条系统级）。基本对齐。但 SC 第 15 条 "Config 系统审计结论与实际发现一致" 在 Solution 中无对应修复动作——这是一个验证项而非修复项。 |
| Requirements <-> Solution coherent | 21/25 | 隐含需求（"修复不应引入新不一致"）在回归验证中覆盖。但提案未明确声明"修复的执行顺序有依赖关系"——H-1 应在 M-9 之前执行（因为 rubric scale 值可能影响 INLINE 版本号标记中的引用），但 Fix Order 仅按 H/M 优先级排列。 |

**Deductions**: -8 (H-2 可能遗漏 eval/SKILL.md 同步更新; 修复顺序依赖关系未分析; SC 包含无对应修复的验证项)

---

## Score Summary

| Dimension | Score | Max |
|-----------|-------|-----|
| Problem Definition | 98 | 110 |
| Solution Clarity | 108 | 120 |
| Industry Benchmarking | 72 | 120 |
| Requirements Completeness | 78 | 110 |
| Solution Creativity | 35 | 100 |
| Feasibility | 88 | 100 |
| Scope Definition | 72 | 80 |
| Risk Assessment | 72 | 90 |
| Success Criteria | 68 | 80 |
| Logical Consistency | 82 | 90 |
| **Total** | **773** | **1000** |

---

## Attacks

1. **Industry Benchmarking**: 选择论证循环论证 — 引用"当前阶段最务实"但没有给出"务实"的量化标准。为什么不能在本次修复中同时引入一条 grep 断言到 justfile？这需要的投入约 10 分钟，但可以将 rubric-reference 同步问题从"依赖人工注释"转化为"依赖自动化检查"。提案在 Alternatives 表中承认长期需要自动化，却没有给出从手动到自动的迁移阶梯。

2. **Industry Benchmarking**: conftest/Pact 引用是装饰性的 — "可映射到 rubric-reference.md 与 rubric frontmatter 的一致性校验"只是一句概念映射，没有给出 conftest 如何解析 markdown frontmatter 的技术方案。Pact 的 provider-consumer 契约测试映射到"skill 间的输入输出契约验证"更是模糊——这些 skill 是 LLM prompt，不是 API endpoint。引用行业工具但未分析其适用性差异，是虚增可信度。

3. **Solution Creativity**: 维护注释不是解决方案 — H-1 的核心问题是多真相源（rubric frontmatter、rubric-reference.md、命令 argument-hint、命令 description、config 默认值）。修复方案是"手动更新五处 + 添加维护注释"。但维护注释是人类注意力的最低级保障——它依赖"下次修改 rubric 的人会读到这个注释"。这不是系统保证。方案应该是让 rubric-reference.md 从 rubric frontmatter 自动生成，或者至少让 eval-journey 命令从 rubric frontmatter 读取 scale/target 而非依赖文档缓存。

4. **Requirements Completeness**: NFR 完全缺失 — 这是 12 项文本修改，NFR 确实不太可能成为瓶颈。但 M-1 的 kebab-case 重命名有一个隐藏的 NFR：向后兼容性。如果用户项目中有 `auto.eval.uiDesign: true` 的 `.claude/settings.json`，重命名后该配置项失效。提案在 Risk 中提到了这个问题，但在 Requirements 中没有将其列为兼容性需求。

5. **Risk Assessment**: M-1 是陷阱项 — M-1 要求统一 config key 为 kebab-case，但 alias 兼容需要修改 Go 代码（out-of-scope）。这意味着 M-1 要么：a) 验证 Go config reader 支持 kebab-case 后执行（前提不确定），b) 与 Go 代码变更绑定（scope 膨胀），c) 跳过（修复不完整）。提案选择了条件执行但未给出优先级判断——如果前提不满足，M-1 应该从 In Scope 移至 Out of Scope。当前状态是 M-1 悬在 scope 边界上，增加了执行不确定性。

6. **Solution Clarity**: H-2 修复可能不完整 — H-2 只移除 tech-design/SKILL.md 中的死路径 `docs/features/<slug>/proposal.md`，但 eval/SKILL.md 第 12 行写的是 "Exceptions: `proposal` uses full path `docs/proposals/<slug>/`"。这个"异常说明"的存在本身就暗示 `docs/features/<slug>/proposal.md` 曾是默认路径。修复 H-2 时应同步审查 eval/SKILL.md 的路径解析逻辑是否需要更新。

7. **Risk Assessment**: INLINE 版本号标记的虚假安全感 — M-9 提议在 INLINE 引用处添加 `@ v3.0.0-rc.53` 版本号标记。这创造了一种"有版本管理"的幻觉，但实际上：a) 没有自动化机制检测版本不一致，b) 开发者更新源文件时大概率忘记更新 INLINE 引用处的版本号，c) grep 只能找到"有没有版本号标记"但不能验证"版本号是否正确"。这比没有标记更危险——因为它给人一种已受控的错误印象。

8. **Success Criteria**: MINOR 项的决策无追溯 — 9 项 MINOR 标记为不修，但没有 SC 记录"已审慎决定不修"这个决策。未来有人重新审计时无法判断"不修"是有意决定还是遗漏。应在 SC 中添加一条："MINOR 项的'不修'决策已在 proposal 中记录理由"。

9. **Logical Consistency**: 修复顺序缺乏依赖分析 — Fix Order 按 H->M 优先级排列，但 M-9（INLINE 版本号标记）应在所有内容修复完成后执行，因为修复可能改变 INLINE 引用的源文件内容。当前 Fix Order 将 M-9 排在 H-2 之后、M-1~M-7 之前，如果 M-2 修改了 ui-design 的 auto.eval 实现方式，gen-journeys 中内联的 journey-contract-model.md 可能受到影响——但 M-9 的版本号标记已经在 H-2 之后打上了，版本号可能已经过时。

10. **Problem Definition**: H-1 影响范围可能被低估 — 提案称"仅手动传参用户受影响"，但 rubric-reference.md 是 eval 系统的文档级单一参考，如果 LLM 在执行 eval 上下文中读取 rubric-reference.md（而非 rubric frontmatter）来理解评分体系，则所有 eval-journey 执行都可能受影响——不仅仅是手动传参用户。提案未分析 LLM 是否会读取 rubric-reference.md 作为上下文。

---

## Blindspot Hunt

**[blindspot-1]**: 提案没有考虑修复执行的原子性。12 项修复是 12 次独立的文件修改，如果执行到第 7 项时中断（LLM 会话超时、用户中断），系统将处于半修复状态。应定义一个最小可行修复子集（如仅 4 个 HIGH），使得任何中断点都不留下比修复前更差的状态。

**[blindspot-2]**: 提案没有评估 audit 本身的方法论可靠性。21 处发现是基于"逐文件完整读取"的手动审计，但没有交叉验证审计结论的完整性和准确性。例如：H-2 声称 `docs/features/<slug>/proposal.md` "从未被任何 skill 创建"，但 quick pipeline 是否有特殊的文件创建逻辑？提案提到"需先搜索确认"，这说明审计时并未确认这一点。

**[blindspot-3]**: 提案没有考虑 eval 命令的 description 字段被 LLM 用作 prompt 上下文的影响维度。eval-journey 的 description 写 "1000-point rubric"，这不仅是用户文档——它是 LLM 执行 eval 时的上下文锚点。如果 LLM 基于此理解评分体系是 1000 分制，它可能在生成评分时锚定在 1000 分而非 1150 分，影响评分分布。这个影响路径比"误导用户传入错误 target"更隐蔽。

**[blindspot-4]**: M-3 的修复（添加完整路径）与 H-2 的修复（移除死路径）都涉及 `breakdown-tasks/SKILL.md`，但提案未分析这两个修复是否在同一文件的相邻区域，是否可以合并为一次修改以减少 grep 回归验证的次数。

**[blindspot-5]**: 提案的"回滚计划"过于简单——"git revert 即可"。但如果修复分多次 commit（每项修复一次），git revert 需要精确知道回滚到哪个 commit。如果修复合并为一次 commit，则无法选择性地回滚单条修复。回滚粒度未定义。
