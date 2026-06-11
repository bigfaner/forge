---
reviewer: Prompt Compliance Architect
date: 2026-05-23
proposal: docs/proposals/spec-authority-enforcement/proposal.md
---

# Freeform Review: 规范权威性强制执行提案

## 背景评估

本提案旨在解决一个已由真实事件验证的系统性问题：LLM agent 执行任务时以现有代码为参照而非以规范文档为参照，导致 test-capability-v2 特性中出现 43 处偏差。提案的核心诊断来自教训文档的 Level 4 分析——"LLM agent 天然倾向局部一致性而非全局一致性"，这个诊断是准确的。

提案采取两层防护策略：在 agent 执行层（task-executor.md）增加"声明 Reference Files 为权威来源"和"AC 逐条验收"两个步骤，同时在任务生成层（quick-tasks / breakdown-tasks）改进 Reference Files 质量，要求精确到 section 引用。

提案明确定位为"非创新性改进"，是对已验证工程实践的制度化。这个自我定位是正确的——提案的本质是在 prompt 合成管线的正确位置插入行为锚点，而非发明新机制。

我审查了提案涉及的四个关键文件的实际当前状态：task-executor.md 的执行协议确实没有 Reference Files 权威性声明步骤；coding-enhancement.md、coding-feature.md、coding-refactor.md 的 Step 1 确实只说"read the task file"而未提及 Reference Files；doc.md 的 Step 1 虽然提到"Identify all reference files listed in the task and read them"，但缺少权威性声明和 AC 逐条验收。quick-tasks 的 task.md 模板中 Reference Files 确实硬编码为 `docs/proposals/<slug>/proposal.md — Source proposal` 单一条目。这些问题诊断是准确的。

---

## 关键风险

风险：提案在 task-executor.md 中增加 `<EXTREMELY-IMPORTANT>` Reference Files 声明步骤时，将与已有的 Hard Constraints `<EXTREMELY-IMPORTANT>` 块产生标记稀释效应。

当前 task-executor.md 中已有一个 `<EXTREMELY-IMPORTANT>` 块覆盖 7 条硬约束。提案在 Success Criteria 中要求"需强化的模板 Step 1 包含 `<EXTREMELY-IMPORTANT>` Reference Files 权威性声明"。这意味着 task-executor.md 将出现第二个 `<EXTREMELY-IMPORTANT>` 块（在执行协议的 Step 5 之前）。LLM agent 对同一级别的强调标记存在注意力稀释——当所有内容都是"极其重要"时，没有任何内容真正"极其重要"。提案的 Key Risks 表格承认了这一风险——"Agent 忽略 `<EXTREMELY-IMPORTANT>` 标记"——但将其 Mitigation 定位为"标记放在执行流程的关键位置"并未实质解决稀释问题。

风险：提案要求"每个任务引用 2-5 个 section，而非整个文件"，但未定义 section 引用的标准格式如何与现有 task 文件的 `## Reference Files` section 兼容。

提案 Success Criteria 最后一条要求格式为 `path/to/file.md#section — 简要说明该 section 定义了什么`。但 quick-tasks 的 `templates/task.md` 模板中 Reference Files 条目格式是 `- docs/proposals/<slug>/proposal.md — Source proposal`，使用的是破折号分隔而非 `#section` 锚点格式。这两个格式之间的迁移路径未被讨论。如果只改模板不改已生成的任务文件，旧格式和新格式将在同一系统中并存，agent 可能无法统一解析。

问题：提案的"Agent 层"改进实际上混淆了两个不同层次的职责——task-executor.md（agent 定义文件）和 coding.* 模板（forge-cli 嵌入的 prompt 模板）。

提案在 Solution 中说"在执行协议中增加两个强制步骤"，Success Criteria 中说"task-executor.md 的执行协议包含...步骤（Step 5 前）"。但 task-executor.md 的执行协议是 Step 1-11 的线性流程，Step 5 是"Follow every step in the synthesized strategy exactly"——也就是说，task-executor.md 的 Step 5 是将控制权转交给合成后的 coding-enhancement.md / coding-feature.md 等模板的 Workflow Steps。如果在 task-executor.md 的 Step 5 之前插入 Reference Files 声明步骤，此时 agent 尚未进入具体模板的执行流程，它看到的合成 prompt 中可能只有 `{{TASK_FILE}}` 路径但没有 `## Reference Files` 内容（因为 Reference Files section 在任务文件中，而任务文件要等到 Step 1 "Read Task Definition" 才被读取）。

这造成一个时序矛盾：提案要求 agent 在读取任务文件之前就声明 Reference Files 为权威来源，但 Reference Files 列表定义在任务文件内部。实际的信息流是 task-executor.md Step 3 获取合成 prompt → 合成 prompt 包含任务文件路径和模板 → agent 按模板 Step 1 读任务文件 → 任务文件中有 `## Reference Files`。Reference Files 的权威性声明必须在 agent 读到任务文件之后才能生效。

风险：提案在 coding.* 模板的 Step 1 插入 `<EXTREMELY-IMPORTANT>` Reference Files 声明，但 coding-enhancement.md 和 coding-feature.md 的 Step 1 已有明确的执行顺序——先读 conventions 目录（按 domains frontmatter 过滤），再读任务文件。

当前 coding-enhancement.md Step 1 的指令是"Check `docs/conventions/` and `docs/business-rules/` for project-specific knowledge"，然后"read the task file at `{{TASK_FILE}}`"。如果 `<EXTREMELY-IMPORTANT>` Reference Files 声明插入在 Step 1 开头，它将要求 agent 在读 conventions 之前就加载 Reference Files 中的文档，改变了现有的执行顺序。如果插入在"read the task file"之后，则 agent 在读到 Reference Files 之前已经形成了对任务的理解，锚定效应已生效。提案未明确 Reference Files 声明在 Step 1 内的精确插入位置。

问题：提案要求 quick-tasks 从 proposal 和相关文档中提取"精确 section 引用"，但 quick-tasks skill 的输入只有 proposal.md 一个文件，它可能没有读取 tech-design.md 的上下文。

提案场景 3 描述"从 proposal 提取关键技术约束，将相关 design 文档的精确 section 写入 Reference Files"。但 quick-tasks SKILL.md Step 1 明确说"Read `docs/proposals/<slug>/proposal.md` — the sole input document"。quick-tasks 不走 PRD/design 路径（那正是 quick pipeline 和 full pipeline 的区别）。当没有 tech-design.md 时，quick-tasks 无法提取"精确 section 引用"——因为根本不存在 design 文档。提案未区分有 design 文档和无 design 文档两种情况下的 Reference Files 生成策略。

问题：提案对 breakdown-tasks 的改进方向正确但缺少具体操作指引。

提案说"为非 UI 任务增加 Reference Files 填充指引，要求精确到 section"。但 breakdown-tasks 的 `templates/task.md` 模板中 Reference Files 使用 `{{REFERENCE_FILES}}` 占位符，SKILL.md 的 Step 4a "Business Tasks" 小节没有任何关于 Reference Files 填充的规则——它只提到了 User Stories、Hard Rules、breaking 标记、scope assignment 和 priority assignment。提案的改动需要在 Step 4a 中增加一个 Reference Files 填充子步骤，但提案未说明这个子步骤的操作逻辑：是从 tech-design.md 的哪些 section 提取？是按 design element 映射表（Step 2 的 Element Mapping）反向推导每个任务相关的 tech-design section？还是依赖 agent 自行判断？

风险：提案的 Key Risks 表格中"Reference Files section 引用过时"的 Mitigation 仅为"引用 section 标题而非行号，行号仅作为辅助定位"。

这个 Mitigation 隐含了一个假设：section 标题是稳定的。但在实际项目中，tech-design.md 在迭代过程中可能重组结构、重命名 section、或合并/split section。当 section 标题变更时，Reference Files 中的引用将指向不存在的 section——agent 收到指示去读一个不存在的 section 时，会怎么处理？它可能跳过这条 Reference File 条目（因为它无法定位），也可能误解为整个文件都是 Reference。提案没有定义 Reference Files 引用失效时的降级行为。

问题：提案将 doc.md 模板列为"优先"改进对象，但 doc.md 是唯一一个 Step 1 已经要求读 Reference Files 的模板。

教训文档明确指出 doc.md 的 Step 1 "说'Identify all reference files listed in the task and read them'（较好）"。doc.md 已经做了正确的事——它要求读取所有 Reference Files。它的问题仅在于 Step 3 Self-Check 缺少 AC 逐条验收。如果提案对 doc.md 的主要改动是增加 `<EXTREMELY-IMPORTANT>` 声明，实际上是对一个已经合规的模板做冗余加固，而真正需要加固的 coding.* 模板（它们完全缺失 Reference Files 读取指令）反而需要更多注意力。优先级判断可能需要调整。

风险：提案声明"纯 Markdown 文档修改，不涉及 Go 代码变更"和"模板改动不改变现有任务的执行流程结构，只是在已有步骤间插入新步骤"，但对 coding.* 模板 Step 1 的修改可能隐式改变 `forge prompt get-by-task-id` 的合成 prompt 长度。

coding-enhancement.md 和 coding-feature.md 的 Step 1 目前包含 conventions 目录扫描逻辑。如果在该步骤中增加 `<EXTREMELY-IMPORTANT>` Reference Files 加载指令，不仅增加了 prompt 长度，还增加了 agent 在 Step 1 的 I/O 操作次数（需要额外读取 Reference Files 中的文档）。对于包含 15 个任务的 quick pipeline，每个任务的 Step 1 都多出 2-5 个文件读取操作，这可能影响执行效率。提案的 Non-Functional Requirements 说"Reference Files 精确引用不应导致 prompt 过长"，但未评估实际 I/O 开销。

问题：提案的 Scope 列出"审计全部 19 个 `forge-cli/pkg/prompt/data/*.md` 模板"但没有定义审计标准。

哪些模板需要 Reference Files 权威性声明？提案只讨论了 coding.* 和 doc.* 类型的模板。但 19 个模板中还包括 clean-code.md、coding-cleanup.md、coding-fix.md、test-gen-and-run.md、test-gen-scripts.md、test-run.md、test-verify-regression.md、validation-code.md、validation-ux.md、gate.md、fix-record-missed.md 等模板。这些模板是否需要 Reference Files 权威性声明？coding-fix.md 的 Step 1 已经有 conventions 目录扫描但没有 Reference Files 声明——fix task 的 Reference Files 是否同样重要？gate.md 是 stage-gate 验证任务——它的 Reference Files 需求是什么？提案将"输出审计报告（哪些模板需要强化 + 理由）"列为 Success Criteria，但在 Scope 中未提供审计的判断标准，等于将核心决策推迟到实施阶段。

风险：提案修改 `forge-cli/pkg/prompt/data/*.md` 模板文件，但这些文件通过 Go 的 `embed.FS` 嵌入二进制。虽然提案声明"纯 Markdown 修改不涉及 Go 代码"，但任何对 `forge-cli/pkg/prompt/data/*.md` 的修改都需要重新编译 `forge-cli` 二进制才能生效。

提案在 Constraints & Dependencies 中说"纯 Markdown 文档修改，不涉及 Go 代码变更"。这在技术上准确——不需要修改 `.go` 文件。但从部署角度看，用户必须重新 `go build` 才能让模板变更生效。如果用户修改了插件目录下的 `agents/task-executor.md`（即时生效）但没有重新编译 forge-cli，那么 task-executor.md 中新增的"声明 Reference Files 权威性"步骤会生效，但 coding.* 模板中新增的 `<EXTREMELY-IMPORTANT>` Reference Files 声明不会生效——因为 agent 看到的合成 prompt 仍使用旧的嵌入模板。这两层防护将处于不一致状态。

问题：提案的 Success Criteria 将 AC 验收步骤放在 Step 8 之前（submit-task 之前），但 task-executor.md 当前的 Step 7-8-9 是"check blocked → submit-task → git-commit"。

提案要求"task-executor.md 的执行协议包含'AC 逐条验收'步骤（Step 8 前）"。但 task-executor.md 的执行协议是 agent 级别的（11 步），而 AC 验收应该是模板级别的（在 coding-enhancement.md 的 Step 3 Self-Check 之后，submit-task 之前）。如果 AC 验收放在 task-executor.md 的 Step 8 之前，它将在所有任务类型上强制执行，包括 doc、test、validation 等类型——但这些类型的 AC 格式和验证逻辑与 coding 类型完全不同。提案未讨论 AC 验收步骤对不同任务类型的适配问题。

---

## 改进建议

建议：分层使用强调标记，避免 `<EXTREMELY-IMPORTANT>` 稀释。

当前 task-executor.md 的 Hard Constraints 使用 `<EXTREMELY-IMPORTANT>`。建议 Reference Files 权威性声明使用不同的强调机制——例如在模板的 Step 1 中使用 `<IMPORTANT>` 加上明确的行为指令（"You MUST read all files listed in ## Reference Files before proceeding"），而非另起一个 `<EXTREMELY-IMPORTANT>` 块。真正需要 `EXTREMELY-IMPORTANT` 的是 task-executor.md Hard Constraints 中新增一条："Reference Files listed in the task are the authoritative source of truth. Conflicts between Reference Files and existing code MUST be resolved in favor of Reference Files." 这样只增加一条规则，不增加一个强调块。

建议：将 Reference Files 权威性声明的位置明确为模板 Step 1 内部（"read the task file"之后），而非 task-executor.md 执行协议中。

正确的插入点是 coding-enhancement.md / coding-feature.md / coding-refactor.md 的 Step 1 中"read the task file"之后、`<IMPORTANT>` Hard Rules 块之前。此时 agent 已经读到任务文件中的 `## Reference Files` section，可以立即声明这些文件为权威来源并按需加载。task-executor.md 只需要在 Hard Constraints 中增加一条规则作为兜底，不需要在执行协议中增加独立步骤。

建议：为 quick-tasks 和 breakdown-tasks 定义不同的 Reference Files 策略。

quick-tasks（无 tech-design.md）的 Reference Files 应包含：proposal.md（必须）+ 现有的 conventions 文档路径（从 Step 0 的 language resolution 中获取）。如果有用户手动提供的 design 文档，也应包含。模板应改为 `{{REFERENCE_FILES}}` 占位符而非硬编码。

breakdown-tasks（有 tech-design.md）的 Reference Files 策略需要在 SKILL.md Step 4a 中新增填充规则：按 Step 2 Element Mapping 表的对应关系，为每个任务引用其 design element 在 tech-design.md 中所在 section 的标题。例如，如果任务是"实现 Convention 加载机制"，对应的 design element 在 tech-design.md 的 `## Convention Loading` section，则 Reference Files 应为 `design/tech-design.md#Convention Loading — 定义了加载机制的两级架构`。

建议：为 Reference Files 引用失效定义降级行为。

在模板中增加指令："If a referenced section title does not exist in the file, read the entire file and identify the closest matching section. Report the mismatch in your output." 这比静默跳过或猜测都更安全。

建议：在 task-executor.md 的 Hard Constraints 中增加 AC 验收规则，而非在执行协议中增加独立步骤。

将第 7 条改为："HARD RULES OVERRIDE ... Before submitting, verify every Acceptance Criteria checkbox is satisfied." 这样 AC 验收作为硬约束存在，由 agent 在 submit-task 之前自行执行，而不需要在 11 步协议中插入新步骤。具体如何验收（逐条检查、修复重验）由各模板的 Self-Check / Step 3 步骤定义。

建议：明确模板审计的分类标准。

审计标准应基于任务类型是否涉及规范遵循：(1) 所有 coding.* 模板——需要 Reference Files 权威性声明（因为编码任务最容易出现以现有代码为参照的问题）；(2) doc.* 模板——doc.md 已有 Reference Files 读取，其余 doc-* 模板按具体需要评估；(3) test.* 模板——test-gen-and-run.md 等可能需要 Reference Files 来确保生成的测试符合规范；(4) gate.md、fix-record-missed.md、validation-*.md——这些是流程控制型模板，Reference Files 的需求较低。

建议：在提案中增加部署一致性说明。

明确指出：修改 `forge-cli/pkg/prompt/data/*.md` 后需要重新 `go build` 以更新嵌入模板。两层防护的生效前提是 agent 层（task-executor.md）和模板层（coding.*）同时更新。如果只更新了 task-executor.md（即时生效）但未重新编译 forge-cli，Reference Files 权威性声明只在 task-executor.md 的 Hard Constraints 中存在，合成后的 coding.* 模板仍不包含读取指令——效果大打折扣。
