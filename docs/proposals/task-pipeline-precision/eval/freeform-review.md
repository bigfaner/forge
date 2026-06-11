# Freeform Review: Task Pipeline Precision Tuning

**Reviewer**: Prompt Compliance Architect
**Date**: 2026-05-27

---

## Section 1: Background Assessment

This proposal addresses a concrete, measurable efficiency problem in the Forge task execution pipeline: task executors waste significant time on simple tasks because the pipeline lacks precision controls at three stages -- task generation, task metadata, and execution template routing.

The diagnosis is grounded. Four independent lesson documents trace back to a coherent root cause chain: quick-tasks SKILL.md merges tasks by time estimate rather than verifiability, prompt templates treat all tasks of the same `coding.*` type identically regardless of complexity, and Reference Files pointing to proposal sections expose executors to out-of-scope requirements. The proposal's observation that "87% time on thinking" and "25 min for a text substitution" are not isolated incidents but systemic symptoms is well-supported by the cross-referencing between lessons.

The core insight -- that task complexity, not task type, determines optimal exploration depth -- is a meaningful departure from the current type-based routing. The proposal correctly identifies that within a single type like `coding.enhancement`, complexity can range from a one-line path substitution to a cross-module behavioral refactor. The current template system is blind to this variance.

The three-pronged approach (split rules, complexity metadata, template branching) maps cleanly to the three stages where precision is lost. The proposal also correctly identifies the Reference Files scope-creep problem as a consequence of the prompt template's own instruction conflict: `<CRITICAL>`-level spec-code scan overrides `CODING_PRINCIPLES`-level surgical scope, so when Reference Files expose out-of-scope spec content, the executor's highest-priority instruction tells it to fix what it sees.

However, the proposal operates at a level of abstraction that leaves several critical implementation details underspecified, and some of its assumptions about how the existing codebase works do not survive contact with the actual source.

---

## Section 2: Key Risks

风险：complexity 字段注入到 prompt template 的机制未在 `renderTemplate()` 层面设计，提案声称"加一个 `{{COMPLEXITY}}` 占位符是最小扩展"但实际需要同步修改 `FrontmatterData` 结构体、`Task` 结构体、`index.json` schema 和 `prompt.go` 的 `renderTemplate()` 函数。

提案原文："`renderTemplate()` 已有 7 个占位符，加一个 `{{COMPLEXITY}}` 是最小扩展。"

查看 `prompt.go` 的 `renderTemplate()` 函数（第 113-161 行），当前的占位符替换链从 `Task` 结构体字段直接取值。但 `Task` 结构体（`types.go` 第 216-253 行）目前没有 `Complexity` 字段。要注入 `{{COMPLEXITY}}`，需要：(1) `FrontmatterData` 加 `Complexity` 字段（`frontmatter.go`），(2) `Task` 结构体加 `Complexity` 字段，(3) `index.json` 反序列化/序列化路径处理新字段，(4) `renderTemplate()` 加替换行。这不是"最小扩展" -- 它是一条从 frontmatter 解析到 JSON schema 到 prompt 合成的完整数据管道变更。提案对此严重低估了改动面。

风险：提案将 complexity 分支逻辑嵌入到 5 个 coding template 中，但未考虑 `strings.ReplaceAll` 模板系统的限制 -- 当 `{{COMPLEXITY}}` 为 `low` 时如何让模板跳过 Step 1.5 整个段落？

提案原文："low 跳过 Step 1.5 spec-code scan、简化探索"

当前的模板系统使用 `strings.ReplaceAll` 做纯文本替换（`prompt.go` 第 98-106 行的 WARNING 注释已明确说明这个限制）。要让 complexity=low 跳过 Step 1.5，不能简单用一个占位符替换 -- 你需要条件性地删除或保留整个 Step 1.5 段落（约 20 行）。当前 `renderTemplate()` 没有条件块机制。`cleanTemplateOutput()` 只处理空值标签行和空反引号行。要实现这个需求，要么引入一个新的清理规则（如按标记注释删除段落），要么使用预渲染的两个模板版本（low vs full），要么给 `renderTemplate()` 加条件段落注入逻辑。提案完全没有讨论这个架构约束。

问题：提案的 complexity 判定启发式过于依赖静态数量指标，忽略了任务的实际认知复杂度。

提案原文："AC 数量 + Hard Rules + Reference Files 数量" 作为判定依据，以及 "AC ≤ 3 且无 Hard Rules 且 Reference Files ≤ 1 的任务被标记为 `complexity: low`"。

三个 AC 可以是"将常量 A 重命名为 B"这样的机械操作（low），也可以是"设计一个能同时满足性能、安全性和向后兼容性的新接口"这样的深度设计工作（high）。Hard Rules 为空不代表任务简单 -- 可能是 quick-tasks 生成时没有识别出隐含的约束。Reference Files 数量也不可靠 -- 一个 Reference File 可能是一个 200 行的 spec 文档，覆盖了复杂的架构约束。用静态计数代替认知复杂度评估，本质上是将一个需要判断力的问题降级为机械规则，而 quick-tasks 的执行者本身就是 LLM -- 让 LLM 在任务生成阶段做这个判断比硬编码阈值更可靠。

风险：Reference Files 内联化后，proposal 变更时 task 文档不会同步更新，提案承认了这一点但未提供缓解措施。

提案原文（lesson gotcha-task-reference-files-scope-creep）："取舍：quick-tasks 生成阶段需要从 proposal 提取并内联（更多工作），且 proposal 变更时 task doc 不会自动同步。但对于 quick mode（≤15 任务），这个代价可接受。"

这个"代价可接受"的判断存在隐含前提：quick mode 的 proposal 在 task 生成后不会变。但实际开发中，task 执行过程中发现问题回退修改 proposal 是常见场景。如果 task 1 执行时发现 proposal 的某个假设不成立，修改了 proposal，但 task 3 的内联 Reference Files 仍然引用旧版 proposal 的描述 -- executor 会按照过时信息执行。当前指针式引用至少保证了 executor 读到的是 proposal 的最新版本。内联化消灭了 scope creep 但引入了 stale reference 风险，提案没有评估两者的权衡。

问题：提案的 "搜索策略引导" 概念过于模糊，在 5 个 template 中的具体位置和行为未定义。

提案原文："加'先收集后修改'搜索策略引导" 以及 "搜索策略引导出现在所有 5 个 coding template 的 implementation 步骤前"。

"先收集后修改"是一个方向性的指导原则，但 prompt template 中需要的是具体的、可执行的指令序列。当前 5 个 template 的 Step 2 各不相同（TDD cycle / Make Improvements / Locate → Fix / Impact Mapping → Refactor），"搜索策略引导"在每个 template 中应该具体指导什么行为？是"先用 grep -rl 列出所有引用点，再逐一修改"？还是一个搜索 checklist？放置在"implementation 步骤前"意味着在 Step 1.5 和 Step 2 之间，但 Step 1.5 本身就是搜索（spec-code scan），再插入搜索策略引导可能导致指令冗余。提案没有讨论这个层叠关系。

风险：quick-tasks 和 breakdown-tasks 两个 SKILL.md 的同步修改在提案中没有精确的对应关系。

提案原文："breakdown-tasks SKILL.md 同步修改拆分规则" 以及 "两个 SKILL.md 使用相同的判定规则描述，确保逻辑一致"。

查看当前 quick-tasks SKILL.md（271 行）和 breakdown-tasks SKILL.md（228 行），两者在拆分规则、Reference Files 生成、Type Assignment 等方面有大量重叠但不完全相同。breakdown-tasks 有 Phase & Gate Detection、PRD Coverage Verification 等环节，而 quick-tasks 没有。"同步修改"意味着需要同时修改两个文件中关于拆分规则的部分，但这两个文件的上下文结构不同（一个从 proposal 读，一个从 tech-design 读），简单地复制相同描述可能产生歧义。提案没有给出两个文件中各自需要修改的具体段落。

问题：提案移除 quick-tasks 15 coding task 上限的依据不充分。

提案原文："移除 quick-tasks 15 coding task 上限" 以及 "In Scope" 中列出 "移除 15 coding task 上限"。

15 coding task 上限的存在是为了将复杂 feature 推向完整 pipeline（PRD + tech-design + breakdown-tasks）。提案没有解释为什么移除这个上限是安全的 -- 如果一个 feature 需要 >15 coding tasks，很可能是因为需求复杂度超出了 quick mode 的承载能力。提案将"移除上限"混入 precision tuning 的 scope，但它是一个独立的架构决策，与精度控制无关。移除上限可能导致更粗粒度的 feature 进入 quick pipeline，加剧而非缓解提案试图解决的问题。

风险：prompt template 的 `<CRITICAL>` 与 `CODING_PRINCIPLES` 的优先级矛盾未被提案直接解决，只是通过信息隔离间接规避。

提案原文（lesson）："CRITICAL 优先级高于一般原则，executor 遵循了冲突扫描结果而非 scope 边界。"

提案的解决方案是让 Reference Files 内联化，从而让 executor 看不到超出 scope 的 spec。但这没有解决根本矛盾 -- 如果内联的 Reference Files 仍然描述了一个与现有代码的差异，Step 1.5 的 `<CRITICAL>` 级别指令仍然会要求 executor 修复它。问题不在于 Reference Files 指向哪里，而在于 Step 1.5 的指令与 Surgical Changes 原则之间存在优先级冲突。内联化只是缩小了冲突的触发范围，没有消除冲突本身。如果未来有人在 task doc 的 Reference Files 中内联了跨任务范围的 spec 差异描述，同样的越界问题会再次出现。

问题：Success Criteria 中"complexity: low 的任务执行时跳过 Step 1.5 spec-code scan，探索阶段 < 30s"无法可靠验证。

提案原文："complexity: low 的任务执行时跳过 Step 1.5 spec-code scan，探索阶段 < 30s"

"探索阶段 < 30s" 依赖于 LLM 的执行行为和 API 响应时间，不是一个确定性可验证的属性。Step 1.5 被跳过是可以通过检查 `forge prompt get-by-task-id` 输出来验证的（如果模板分支正确），但 "探索阶段 < 30s" 需要时间测量，而 forge 当前的 record 结构（`RecordData`）没有记录 step-by-step 执行时间。这个 success criterion 形同虚设。

问题：提案假设 `coding.fix` 类型也需要复杂度分支，但 fix 类型的 template 已经有独立的 5-step 流程（含额外的 Locate 步骤），与 enhancement/feature 的 4-step 流程结构不同。

提案原文："5 个 coding prompt templates 加复杂度分支"

查看 `coding-fix.md`，它有 5 个步骤（Step 1 Read → Step 1.5 Scan → Step 2 Locate → Step 3 Fix → Step 4 Verify），且 Step 2 Locate 本身就包含定位逻辑。fix task 通常是由 dispatcher 自动生成的，来源是 failed quality gate 或 failed test，其复杂度由错误本身决定而非 proposal 分析。对 fix 类型应用与 enhancement 相同的 complexity 判定逻辑（AC 数量 + Hard Rules + Reference Files 数量）在语义上不匹配。提案没有讨论 fix 类型的特殊性。

---

## Section 3: Improvement Suggestions

建议：用 `renderTemplate()` 后处理的条件段落机制替代 `{{COMPLEXITY}}` 简单占位符，利用已有的 `cleanTemplateOutput()` 模式。

当前 `cleanTemplateOutput()` 已经实现了一套基于文本模式的段落清理逻辑（删除空值标签行、删除空反引号条件句）。扩展这个机制来支持条件段落：在 template 中用注释标记段落边界，如 `<!-- COMPLEXITY:low SKIP-START -->` 和 `<!-- COMPLEXITY:low SKIP-END -->`，`cleanTemplateOutput()` 根据 complexity 值删除对应段落。这样 template 本身保持完整的 4-step 结构，只是后处理时裁剪不需要的段落，避免引入新的模板引擎或拆分 template。

此建议解决了 `renderTemplate()` 缺乏条件块机制的风险。改动集中在 `cleanTemplateOutput()` 一个函数，5 个 template 各加标记注释，不影响 `strings.ReplaceAll` 替换链。

建议：将 complexity 判定从硬编码阈值改为 quick-tasks/breakdown-tasks SKILL.md 中的 LLM 判断指引。

替代方案中"Prompt 内置启发式规则"被Rejected 的理由是"模板膨胀；每次执行都做判定"。但这个理由混淆了两个不同的位置：在 executor prompt 中做判定（每次执行时）和在 task generation prompt 中做判定（一次）。让 quick-tasks/breakdown-tasks 的 SKILL.md 包含一段判断指引（如"评估任务复杂度：如果任务只涉及机械文本替换且无架构决策影响，标记为 low；如果涉及跨模块接口变更或新的错误处理路径，标记为 high"），让 LLM 在生成任务时做认知判断，比硬编码 "AC ≤ 3 = low" 更准确。这比在 executor prompt 中做判定更好，因为任务生成只执行一次。

此建议回应了 complexity 判定启发式过于依赖静态指标的问题。改动仅影响两个 SKILL.md 的 instruction 文本。

建议：在 task doc 模板中同时保留 proposal section 指针（用于溯源）和内联精确定位信息（用于 scope 约束），并加 scope boundary 显式声明。

当前 proposal 的 Reference Files 格式是二选一（要么指向 section 要么内联）。建议改为两层：第一层是内联的精确定位信息（executor 实际使用的 spec 约束），第二层是指向 proposal section 的溯源链接（用 `Source: proposal.md#Section` 标注，明确标记为"仅供参考，不构成执行指令"）。同时在 `## Reference Files` 末尾加一段 scope boundary 声明："SCOPE BOUNDARY: 仅上述列出的文件和修改点属于本任务范围。发现的任何其他不一致应记录但不在本任务中修复。"这直接解决了 `<CRITICAL>` 与 Surgical Changes 的优先级矛盾 -- 不是通过信息隔离，而是通过显式的 scope 边界声明覆盖 spec-code scan 的越界倾向。

此建议同时回应了 Reference Files stale reference 风险和 prompt template 优先级矛盾两个问题。内联信息控制 scope，溯源链接保持可追溯性，scope boundary 声明直接解决指令优先级冲突。

建议：将 "移除 15 coding task 上限" 从本提案 scope 中移出，作为独立 proposal 处理。

上限移除与精度控制是两个独立的关注点。精度控制的核心是让已进入 pipeline 的任务执行更高效，而上限决定什么任务可以进入 pipeline。混在一起使得 success criteria 难以隔离验证 -- 如果移除上限后出现更多粗粒度任务，无法区分是 precision tuning 失效还是上限移除导致的。

此建议回应了移除上限依据不充分的风险。移出后本提案的 scope 更聚焦，10 个文件改动的估计也更可信。

建议：在 Success Criteria 中用 "Step 1.5 被跳过" 替代 "探索阶段 < 30s"，用 `forge prompt get-by-task-id` 输出验证替代时间测量。

"Step 1.5 被跳过" 是一个可以通过 `forge prompt get-by-task-id` 输出直接验证的属性 -- 如果 complexity=low 的任务的 prompt 输出不包含 Step 1.5 的 SPEC-CODE SCAN 段落，则验证通过。"探索阶段 < 30s" 应降级为 Non-Functional Requirement 中的期望目标而非 Success Criterion。

此建议回应了 "探索阶段 < 30s" 不可靠验证的问题。

建议：对 `coding.fix` 类型的 complexity 分支做差异化处理，或不对其应用 complexity 分支。

fix 类型由 dispatcher 自动生成，其 frontmatter 可能不包含 complexity 字段（取决于 fix-task 是如何创建的）。查看 `forge task add` 的参数（`--type coding.fix --title --source-task-id --block-source --var`），没有 complexity 相关参数。建议要么在 fix-task 生成时不设置 complexity（默认 medium，行为不变），要么在提案中明确 fix 类型不纳入 complexity routing 的 scope。

此建议回应了 fix 类型语义不匹配的问题。
