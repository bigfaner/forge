---
created: 2026-05-23
author: faner
status: Approved
---

# Proposal: 规范权威性强制执行——从模板和任务生成层面防止 Agent 偏离

## Problem

Agent 执行任务时以现有代码为参照而非以权威规范文档（tech-design.md 等）为参照，导致大规模偏离。在 test-capability-v2 特性中产生了 43 处与 tech-design.md 不一致的偏差。

### Evidence

- `docs/lessons/gotcha-spec-authority-drift.md` 记录了完整的根因分析和 5 级溯源
- Level 0 根因：coding.* 模板的 Step 1 只说"read the task file"，未要求读 Reference Files
- Level 3 过程原因：缺少"自顶向下验证"步骤——反应式局部修复而非主动式全局审计
- 探索发现：quick-tasks 模板硬编码 Reference Files 为单一 proposal.md；breakdown-tasks 对非 UI 任务没有 Reference Files 填充指引

### Urgency

这是系统性缺陷——每次涉及规范驱动修改的任务都可能重蹈覆辙。

**量化推算**：test-capability-v2 特性包含约 15 个 coding 任务，产生了 43 处偏差（平均每任务 ~2.9 处）。根因分析表明偏差直接归因于"模板 Step 1 未要求读 Reference Files"（Level 0）。假设下一个大型特性规模相当（15 个任务），不修复的预期偏差 = 15 × 2.9 ≈ 43 处，每处偏差的修复成本为平均 15-30 分钟的 review + rework（基于 test-capability-v2 事后修复的实际耗时），总计约 10-20 小时 wasted effort。修复成本：纯文档修改，1-2 小时。ROI：5-10 倍。

## Proposed Solution

两层防护（职责分离：模板层执行，Agent 层兜底）：

1. **模板层（coding.\* 模板 Step 1）**：在模板的 Step 1（读取任务文件之后）插入 `<IMPORTANT>` Reference Files 声明，要求 agent 将 Reference Files 列出的文档声明为权威来源并按需加载；在 Self-Check 步骤中插入 AC 逐条验收。声明块的精确模板文本如下：

```
<IMPORTANT>
## Spec Authority Enforcement

The task file's `## Reference Files` section lists authoritative specification sources.
You MUST:

1. Load each Reference File listed in `## Reference Files` immediately after reading the task file.
2. Treat these documents as the authoritative source of truth — when existing code conflicts with specifications in these documents, follow the specifications.
3. Priority when conflicts arise: task `## Hard Rules` > `## Reference Files` > existing code structure.
4. Output a confirmation after loading: "Loaded Reference Files: [list], treating them as authoritative sources."

If `## Reference Files` is empty or missing, output: "Reference Files empty — falling back to existing code structure and Hard Rules."
</IMPORTANT>
```
2. **Agent 层（task-executor.md Hard Constraints）**：在 Hard Constraints 中增加一条兜底规则——"若模板 Step 1 未包含 Reference Files 声明，agent 仍须主动读取任务文件中的 `## Reference Files` 并将其视为权威来源"

**Self-Check AC 验收的精确插入方式**：AC 逐条验收不是新增步骤，而是插入到 coding.* 模板现有 Self-Check / Verify 步骤（如 coding-feature.md 的 Step 3 "Verify & Finalize"）的内部，作为该步骤的第一个子步骤。具体来说，在现有 Self-Check 的开头插入以下指令：

```
<IMPORTANT>
Before performing other verification checks, validate against each Acceptance Criteria item from the task file:
- For each AC item, output: "[AC-N] PASS/FAIL — [brief reason]"
- If any AC item is FAIL, address the failure before proceeding to other checks.
- If `## Acceptance Criteria` is empty or missing, output: "No AC defined — skipping per-item validation."
</IMPORTANT>
```
3. **任务生成层（quick-tasks / breakdown-tasks）**：改进 Reference Files 生成质量，要求精确到 section 或行号范围

### Innovation Highlights

非创新性改进，而是将已验证的工程实践制度化。关键洞察来自教训文档的 Level 4 分析：LLM agent 天然倾向局部一致性而非全局一致性，因此强制措施必须嵌入执行流程本身，而非依赖 agent 的自觉性。

**关于"两层防护"的诚实声明**：本提案的"两层"本质上是同一种机制（prompt 指令）在不同粒度和时机上的重复强调，而非编译器类型检查 + 运行时断言那种独立失效模式的两层防护。真正的第二层（如 hook 拦截 agent 输出并验证是否加载了 Reference Files）需要修改 Go 代码，属于 Out of Scope。因此本提案的防护强度应理解为"在执行流程关键节点重复锚定规范权威性"，而非"独立的两道防线"。

### User-Visible Behavior Changes

修改生效后，用户在使用 Forge 执行 coding 任务时将观察到以下行为变化：

1. **Reference Files 加载确认**：agent 在 Step 1 执行后会输出确认消息，例如 `"已加载 Reference Files: docs/designs/feature-x.md#Data-Flow, docs/designs/feature-x.md#API-Contract，将其视为权威来源"`
2. **AC 验收报告**：agent 在 Self-Check 步骤中输出逐条验收报告，格式为每个 AC 项的 pass/fail 状态及简要说明
3. **降级提示**：当任务文件中 `## Reference Files` 为空或引用的 section 标题不存在时，agent 将输出警告并说明降级行为（见 Requirements Analysis 中的边缘场景定义）

## Requirements Analysis

### Key Scenarios

- **场景 1：coding.\* 任务执行**：agent 收到合成 prompt → 进入 coding.\* 模板 Step 1 读取任务文件 → 声明 Reference Files 列出的文档为权威来源 → 按需加载 tech-design.md 中与当前实现相关的 section → 以规范为权威而非现有代码 → 实现完成后在 Self-Check 步骤中对照 AC 逐条验收
- **场景 2：doc.\* 任务执行**：同上，但验收重点是文档结构合规而非路径命名。doc.* 任务的 Reference Files 来源与 coding.* 不同：(a) 若任务涉及从 tech-design.md 生成文档，则引用 tech-design.md 的相关 section；(b) 若任务涉及修改现有文档（如更新 README），则引用该文档本身 + 相关的设计文档 section；(c) 若任务仅涉及 proposal.md，则与场景 4 策略一致
- **场景 3：无 AC 的任务**：当任务文件中 `## Acceptance Criteria` 为空或缺失时，跳过 AC 逐条验收步骤，agent 输出提示 `"无 AC 定义，跳过逐条验收"`，但仍执行 Reference Files 加载
- **场景 4：quick-tasks 生成任务（无 tech-design.md）**：quick-tasks 的输入只有 proposal.md，不存在 tech-design.md。策略：从 proposal.md 中提取关键技术约束和决策，将 proposal.md 中与当前实现直接相关的 section 写入 Reference Files（格式：`proposal.md#Section-Title — 该 section 定义了 X`）。若 proposal.md 中引用了外部设计文档且文件存在，则同时引用该文档的相关 section
- **场景 5：breakdown-tasks 生成任务（有 tech-design.md）**：从 tech-design.md 的架构决策中提取每个任务相关的 section，精确引用而非笼统指向整个文件

### Non-Functional Requirements

- Reference Files 精确引用不应导致 prompt 过长——每个任务引用 2-5 个 section，而非整个文件
- 模板改动不改变现有任务的执行流程结构，只是在已有步骤间插入新步骤
- **新旧格式兼容性**：旧格式 `path/to/file.md`（无 section 锚点）继续有效。agent 遇到旧格式时，加载整个文件作为参考（不视为精确引用）。新格式 `path/to/file.md#section — 说明` 是推荐格式，本次修改仅影响新任务的生成模板，不要求迁移已有任务文件

### Edge Cases & Degradation

- **`## Reference Files` 为空或缺失**：agent 输出警告 `"Reference Files 为空，将以现有代码结构和 Hard Rules 为参照"`，继续执行但不进行规范权威性声明
- **引用的 section 标题不存在**：agent 输出警告 `"Reference Files 引用的 section 'X' 在 Y 中不存在，降级为读取该文件的全局内容"`，退化为文件级引用
- **引用的文件路径不存在**：agent 输出警告 `"Reference Files 引用的文件 X 不存在，跳过该引用"`，继续处理剩余 Reference Files
- **降级时的兜底**：当所有 Reference Files 引用失效时，agent 回退到现有行为（以代码为参照），但在输出中明确标注"本次执行未加载权威规范来源"

### Priority Rules

当 Reference Files 内容与任务 Hard Rules 冲突时，优先级从高到低为：

1. **任务 `## Hard Rules`**：最高优先级——Hard Rules 是任务的硬约束（如"MUST NOT modify existing test files"），不可被 Reference Files 覆盖
2. **`## Reference Files`**：次高优先级——Reference Files 是规范权威来源，当 Hard Rules 未覆盖时以 Reference Files 为准
3. **现有代码结构**：最低优先级——仅当 Hard Rules 和 Reference Files 均未涉及时的参照

冲突时的行为：agent 输出冲突说明，以 Hard Rules 为准，并在 AC 验收报告中标注冲突项

### Constraints & Dependencies

- 修改 `plugins/forge/` 下的文件前必须遵循 `docs/conventions/forge-distribution.md`
- coding.\* 模板通过 `embed.FS` 嵌入二进制，修改 `forge-cli/pkg/prompt/data/*.md` 后必须执行 `go build` 才能生效；task-executor.md 为即时代效文件，不需要编译
- 两层防护的生效依赖于 coding.\* 模板和 task-executor.md 的同步更新——必须确保 `go build` 在部署前完成

## Alternatives & Industry Benchmarking

### Industry Solutions

业界解决 LLM agent 规范遵循问题的常见手段：

1. **RAG 检索**：LangChain RetrievalQA（github.com/langchain-ai/langchain）和 LlamaIndex QueryEngine（docs.llamaindex.ai）在运行时检索相关文档片段注入 prompt，确保 agent 有权访问规范内容。但 RAG 解决的是"信息可达性"而非"信息权威性"——agent 仍可能以现有代码为准。
2. **Prompt 模板强制引用**：MetaGPT（github.com/geekan/MetaGPT）在 SOP 流程中嵌入文档引用步骤；CrewAI（docs.crewai.com）的 Knowledge 系统要求每个 agent 在执行前加载指定知识源。这类方案与本提案最接近——通过执行流程强制 agent 在特定步骤加载规范。
3. **Structured Output 约束**：OpenAI Function Calling（platform.openai.com/docs/guides/function-calling）和 Anthropic Tool Use（docs.anthropic.com/en/docs/build-with-claude/tool-use）通过 JSON Schema 约束输出格式，可要求 agent 必须输出"已加载的 Reference Files 列表"和"AC 验收结果"作为结构化输出。但这改变 agent 的输出模式，侵入性较高。
4. **认知锚定理论**：Prompt 编写指南（docs.anthropic.com/en/docs/build-with-claude/prompt-engineering/overview）建议在 prompt 开头声明关键规则以提高遵从率，这与本提案在 Step 1 声明 Reference Files 权威性的策略一致。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零成本 | 问题必复发 | Rejected: 教训已证明代价高昂 |
| CLI 层 auto-inline | 自研 | 最可靠，agent 无法跳过 | 需改 Go 代码，增加 prompt 长度 | Rejected: 用户选择纯文档方案 |
| 仅修改 task-executor.md Hard Constraints | Claude Prompt Engineering 指南 | 最小改动，即时生效，无编译依赖 | 仅一层防护，无模板层保障 | Considered: 作为兜底规则采纳 |
| Prompt 模板 + Agent 协议（本方案） | MetaGPT SOP + Claude Prompt Engineering | 零代码改动，两层防护，覆盖主要任务类型 | 依赖 agent 遵守 `<IMPORTANT>` 标记；coding.\* 模板需 go build | **Selected: 最小有效改动** |

**为什么两层而非仅 Hard Constraints 一层**：Hard Constraints（task-executor.md）在 agent 的合成 prompt 中位于顶层，距离具体执行步骤较远——它是一般性规则，不是执行流程中的即时锚点。模板层声明在 Step 1 内部，恰好是 agent 即将开始实现的时刻，具有更高的即时可见性。两者共享"prompt 指令"这一机制，但作用于不同的认知锚定时机：模板层是"执行时锚定"（task-level），Hard Constraints 是"全局锚定"（agent-level）。这类似于法律体系中"具体条文"（场景触发）与"宪法原则"（兜底适用）的关系——不是重复，而是不同粒度的约束。增量价值：模板层失效时（agent 跳过 Step 1 声明），Hard Constraints 仍可触发；Hard Constraints 被忽略时，模板层在执行流程关键位置再次提醒。虽然两者的失效模式相同（agent 不遵从 prompt），但触发时机不同，统计上降低了同时失效的概率。

## Feasibility Assessment

### Technical Feasibility

Markdown 编辑 + `go build`（coding.\* 模板通过 embed.FS 嵌入），技术风险低但需注意编译步骤。

### Resource & Timeline

预计 4-6 个文件修改（task-executor.md + 审计后确认的模板 + quick-tasks/breakdown-tasks skill）+ `go build`，1-2 小时完成。

### Dependency Readiness

无外部依赖。内部依赖：coding.\* 模板修改后需 `go build`（见 Constraints）。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| "模板层改 4 个文件就够了" | 事实核查（Explore agent 审计全部 19 个模板） | Refined: 需审计全部模板，实际需修改数量待定 |
| "quick-tasks/breakdown-tasks 已经在 Reference Files 中填入设计文档" | 事实核查（读取模板和 SKILL.md） | Overturned: quick-tasks 硬编码只有 proposal.md；breakdown-tasks 对非 UI 任务无填充指引 |
| "Reference Files 只需列文件路径" | XY Detection — 用户实际需要的是精确到 section 的引用 | Refined: 精确引用减少 agent 注意力分散，提高规范遵循率 |

## Scope

### In Scope

- 修改 `forge-cli/pkg/prompt/data/coding.*.md` 模板：在 Step 1（读取任务文件之后）插入 `<IMPORTANT>` Reference Files 权威性声明，在 Self-Check 步骤插入 AC 逐条验收
- 修改 `plugins/forge/agents/task-executor.md`：在 Hard Constraints 中增加兜底规则（若模板未声明 Reference Files，agent 仍须主动读取并视为权威来源）
- 审计全部 19 个 `forge-cli/pkg/prompt/data/*.md` 模板，确定哪些需要 Reference Files 权威性声明
- **审计标准**：模板需要强化的条件为满足以下任一：(a) 模板用于 coding 或 doc 任务类型；(b) 模板的 Step 1 包含"读取任务文件"步骤（因为 Reference Files 声明必须在读取任务文件之后）；(c) 模板涉及需要对照规范执行的实现/修改任务（排除纯信息查询模板）
- 对需要强化的模板，在 Step 1 加 `<IMPORTANT>` Reference Files 声明（使用 `<IMPORTANT>` 而非 `<EXTREMELY-IMPORTANT>` 以避免标记稀释——模板中已有的 `EXTREMELY-IMPORTANT` 块保持不变，Reference Files 声明使用 `<IMPORTANT>` + 行为指令的分层标记方式）
- 改进 `plugins/forge/skills/quick-tasks/SKILL.md`（不是 templates/——quick-tasks 的任务生成逻辑在 SKILL.md 的 workflow 步骤中）：在生成 `## Reference Files` 的步骤中，增加从 proposal.md 提取精确 section 引用的指引
- 改进 `plugins/forge/skills/breakdown-tasks/SKILL.md`：为非 UI 任务增加 Reference Files 填充指引，要求精确到 section。**提取决策逻辑**：对每个生成的任务，按以下步骤确定 Reference Files：(1) 从任务的 `## Affected Files` 中提取涉及的文件路径列表；(2) 在 tech-design.md 中搜索提及这些文件路径的 section；(3) 同时提取 tech-design.md 中与任务描述关键词匹配的架构决策 section；(4) 合并去重后保留 2-5 个最相关的 section。此逻辑在 SKILL.md 中作为指引描述（非算法实现），标记为"待实施时细化"——breakdown-tasks 的执行者是 LLM agent，提取质量取决于 agent 的理解能力，SKILL.md 中给出启发式策略而非确定性算法

### Out of Scope

- 修改 `forge-cli` Go 代码（auto-inline Reference Files）
- 修改任务文件格式（index.json schema）
- 添加 hooks 或运行时验证机制
- 改变 `forge prompt get-by-task-id` 的合成逻辑

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Agent 忽略 `<IMPORTANT>` 标记 | M | H | 标记放在执行流程的关键位置（Step 1 和 Self-Check），使用 `<IMPORTANT>` 而非 `<EXTREMELY-IMPORTANT>` 避免标记稀释；task-executor.md Hard Constraints 兜底 |
| Reference Files section 引用过时（design 文档更新后行号偏移） | M | M | 引用 section 标题而非行号，行号仅作为辅助定位；定义降级行为（见 Edge Cases） |
| breakdown-tasks 生成任务时 Reference Files 填充不完整 | L | M | 在 SKILL.md 中加显式 checklist：每个任务必须有 ≥1 个 design-level Reference File |
| task-executor.md 与 coding.\* 模板更新不同步导致两层防护不一致 | M | H | 在 Success Criteria 中加入部署一致性验证：`go build` 后检查 coding.\* 模板中的 Reference Files 声明是否存在；优先修改 coding.\* 模板，task-executor.md 仅作兜底 |

## Success Criteria

- [ ] coding.\* 模板 Step 1（读取任务文件之后）包含 `<IMPORTANT>` Reference Files 权威性声明；Self-Check 步骤包含 AC 逐条验收
- [ ] task-executor.md Hard Constraints 包含兜底规则：若模板未声明 Reference Files，agent 仍须主动读取并视为权威来源
- [ ] 全部 19 个模板完成审计，输出审计报告（哪些模板需要强化 + 理由）
- [ ] 需强化的模板 Step 1 包含 `<IMPORTANT>` Reference Files 权威性声明（不使用 `EXTREMELY-IMPORTANT` 以避免与模板中已有的 `EXTREMELY-IMPORTANT` 块产生标记稀释）
- [ ] quick-tasks 生成的 coding 任务 Reference Files 包含 ≥1 个精确 section 引用（格式为 `file.md#Section-Title`）；当输入仅有 proposal.md 时，引用 proposal.md 中的具体 section 即可满足此条件（"非仅 proposal.md"指不能只写裸文件路径 `proposal.md` 而无 section 锚点）
- [ ] breakdown-tasks SKILL.md 包含非 UI 任务的 Reference Files 填充规则
- [ ] 新生成任务的 Reference Files 条目格式统一：`path/to/file.md#section — 简要说明该 section 定义了什么`（旧任务文件的旧格式不受此条约束）
- [ ] 部署一致性验证：`go build` 后通过 `forge prompt get-by-task-id` 获取的合成 prompt 中包含 Reference Files 声明文本
- [ ] 行为验证：使用修改后的模板执行 3-5 个 coding 任务后，每个任务的 agent 输出中均包含 Reference Files 加载确认和 AC 验收报告（pass/fail 列表），且无规范偏离（人工 review 对照 Reference Files 确认）

## Rollback Plan

如果修改后 agent 行为异常（如过度遵守 Reference Files 导致忽视用户意图、加载确认或 AC 验收输出过于冗长），回滚策略如下：

1. **即时回滚**：task-executor.md 的 Hard Constraints 修改和 coding.\* 模板修改均通过 git revert 即时回滚，无需重新编译
2. **部分回滚**：若问题仅出在 AC 验收步骤（如 agent 在 Self-Check 中陷入循环），可仅移除 Self-Check 中的 AC 验收指令，保留 Reference Files 加载声明
3. **验证回滚**：回滚后通过 `forge prompt get-by-task-id` 验证合成 prompt 不再包含 Reference Files 声明文本，确认回滚生效

## Next Steps

- **预实验**：在全面修改前，先用 1 个 coding 任务手动在 prompt 中注入 Reference Files 声明（不修改模板），验证 agent 是否确实以规范为参照而非以代码为参照。若预实验显示 agent 仍以代码为准，则需要重新评估方案（如考虑 Structured Output 或 hook 机制）
- Proceed to `/quick-tasks` to generate and execute tasks
