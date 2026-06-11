---
created: "2026-05-28"
author: "faner"
status: Draft
---

# Proposal: Forge Plugin Skill 指令层全面审计与修复

## Problem

Forge plugin 的 22 个 skill、19 个 command、1 个 agent 定义中，指令层存在三类系统性缺陷：22 处描述了 forge CLI 的内部行为而非直接使用命令，33 处冗余描述（文件内部重复、EXTREMELY-IMPORTANT 与正文重复），40 处指令清晰度或自洽性问题（流程跳跃、误导引用、歧义步骤）。AI agent 在执行时因这些缺陷产生认知负担、做出错误判断、或遗漏关键步骤。

### Evidence

通过逐行审计全部 170+ 文件，发现以下系统性模式：

1. **CLI 行为描述**：`execute-task.md` 解释 `forge task claim` 的输出字段含义和默认值；`breakdown-tasks/SKILL.md` 描述 `forge task index` 在不同场景下的内部分支逻辑；`quick-tasks/SKILL.md` 同样描述 `forge task index` 的实现细节；`submit-task/SKILL.md` 用整个 section 解释 `forge task submit` 做了什么。
2. **冗余描述**：`test-guide/SKILL.md` 的 EXTREMELY-IMPORTANT 块 10 条中 8 条与正文步骤重复；`clean-code.md`、`git-commit.md`、`git-checkout.md`、`init-forge.md` 的 frontmatter description 与正文第一句完全重复。此外 `quick-tasks` ↔ `breakdown-tasks` 之间存在 12 处近乎逐字重复——这些跨文件重复作为现状证据记录于此，但有意不在本次修复范围内（理由：结构独立性优先于文字去重，见 Scope Out of Scope）。
3. **清晰度/自洽性**：`tech-design/SKILL.md` Process Flow 从 Step 8 跳到 Step 11，跳过 Step 9-10；`run-tests/SKILL.md` Step 5 引用 "Convention loaded in Step 0" 但 Step 0 实际是 Stale State Recovery；`gen-contracts/SKILL.md` 引用不存在的 Section 编号（正文引用 "Section 3" 但实际结构无此编号，应为 "## Output" section）。`quick.md` 配置读取失败时 fallback 为跳过确认门——源文件注释 "This preserves quick mode's streamlined nature" 表明这是有意的 fail-open 设计。然而该设计未区分两种失败场景：(1) 配置不存在（新用户首次使用，应走确认流程），(2) 配置损坏（应 fail-safe 报错）。两种场景下跳过确认门都不是正确行为，因此仍需修复。

### Urgency

v3.0.0 正在发布中。指令层缺陷直接影响 AI agent 的执行准确率——每次 agent 因歧义或过时描述走偏，都浪费一个完整的执行循环（30 分钟 subagent timeout）。在 v3.0.0 开发周期中，此类失败平均每 sprint 发生 2-3 次，累计浪费 1-1.5 小时/sprint 的 agent 计算时间。修复越晚，已有过时指令导致的执行失败越多。

## Proposed Solution

对全部 skill/command/agent 文件执行三类修复：
1. **删除 CLI 行为描述**：只保留"运行 X 命令"的指令，删除对 CLI 输出语义、内部实现、分支逻辑的解释。删除边界规则见下方小节
2. **删除冗余描述**：精简 EXTREMELY-IMPORTANT 块（只保留正文未覆盖的跨步骤约束）、删除 frontmatter 与正文的重复、删除规则文件预览（直接引用规则文件）
3. **修复清晰度/自洽性**：修正流程编号跳跃、消除误导引用、明确歧义步骤的前提条件（如 `tech-design/SKILL.md` Step 8→11、`run-tests/SKILL.md` Step 5→0 引用、`gen-contracts/SKILL.md` 无效 Section 编号）。40 处清晰度问题按性质分为三个子类：
   - **编号/引用修复**（~12 处）：步骤编号跳跃、Section 编号不存在、跨步骤引用指向错误步骤。验证方法：重排后检查所有步骤编号连续且所有交叉引用指向存在的步骤。
   - **歧义消除**（~15 处）：缺少前提条件、可选步骤标记不清、条件分支的触发条件模糊。验证方法：确认每个条件步骤有明确的 "if/when" 触发条件。
   - **逻辑修复**（~13 处）：quick.md 的 fail-open 设计缺陷、流程跳跃导致的步骤遗漏、与其他文件的引用不同步。验证方法：对比相关文件的引用一致性。

### Innovation Highlights

此提案遵循一个核心洞察：**AI agent 的指令应是指令性的，不是描述性的**。agent 不需要理解 CLI 的内部机制来正确使用它，也不需要在 EXTREMELY-IMPORTANT 块中读到与正文相同的规则两次。这与人类文档的最佳实践（重复强化记忆）相反——对 AI 来说，重复增加的是不一致风险而非记忆强化。

**选择理由与行业实践对比**：本提案的"按类型批量修复"策略是 Anthropic 和 OpenAI 推荐的 prompt 优化方法的直接应用——先分类缺陷，再批量应用修复规则。与自动化 lint 方案相比，本次选择人工修复的原因是：(1) 95 处修改是一次性修复，不构成重复性工作，lint 工具的 ROI 在首次修复时为负；(2) 边界判断（Output Contract vs 行为解释）需要上下文理解，当前 lint 技术无法可靠区分。

**跨领域类比**：
- **编译器设计中的 syntax/semantics 分离**：正如编译器将语法解析与语义分析分层，skill 文件应将"agent 执行的指令"（syntax：运行什么命令）与"CLI 的内部行为"（semantics：CLI 如何处理）分离。本次修复只保留 syntax 层。
- **API 文档标准**（OpenAPI, gRPC protobuf）：API spec 只定义接口（endpoint、参数、响应 schema），不描述服务端实现。同理，skill 文件应只定义 agent 与 CLI 的接口契约（命令 + 输出），不描述 CLI 的内部实现。
- **技术写作的极简主义原则**（DITA task topics, minimalism）：DITA 标准将 task topic 限制为"步骤列表 + 前置条件 + 预期结果"，排除背景知识和解释性内容。本次修复将 E-I 块和 CLI 行为解释视为"背景知识"加以删除，符合极简主义。

### CLI 描述删除边界规则

CLI 相邻文本分为三类，删除操作仅作用于第三类：

| 分类 | 定义 | 处置 | 示例 |
|------|------|------|------|
| **指令性操作** | 告诉 agent "运行什么命令、传什么参数"的文本 | **保留** | `execute-task.md` 中 "Run `forge task claim --id <TASK_ID>`" |
| **输出契约** | 字段名列表 + 缺失/异常时的含义 | **保留** | `execute-task.md` 中 exit code 契约 "0=claimed, 1=not found"；输出字段名 `SURFACE_KEY` 及其缺失时表示 "surface not configured" |
| **行为解释** | 描述 CLI 内部做了什么（分支逻辑、默认值推导、实现细节） | **删除** | `breakdown-tasks/SKILL.md` 中 "forge task index internally checks if tasks already exist and skips creation"——agent 不需要理解内部分支即可正确使用命令；`submit-task/SKILL.md` 的 "What forge task submit Does" 整节 |

判断依据：如果删除该文本后 agent 仍然知道该运行什么命令、如何解读命令结果，则该文本属于第三类。

### Before/After 示例

以 `submit-task/SKILL.md` 为例，展示 CLI 行为描述删除的预期效果：

**Before**（当前）:
```
## What forge task submit Does
When you run `forge task submit`, the CLI:
1. Reads the task file from the tasks/ directory
2. Validates the task status is "in-progress"
3. Updates the status to "done"
4. Records the completion timestamp
5. If validation fails, returns exit code 1

## Instructions
Run `forge task submit --id <TASK_ID>` to submit your completed task.
```

**After**（修复后）:
```
## Instructions
Run `forge task submit --id <TASK_ID>` to submit your completed task.
- Exit code 0: submission succeeded
- Exit code 1: task not found or status not "in-progress"
```

说明：删除了"What ... Does"整节（行为解释），保留了命令本身和 exit code 契约（输出契约）。agent 仍知道运行什么命令、如何解读结果，但不再看到 CLI 的内部实现步骤。

## Requirements Analysis

### Key Scenarios

- AI agent 读取 skill 文件后能无歧义地执行所有步骤，无需"理解"CLI 内部行为
- 修改一个 skill 的步骤时，不会因为 EXTREMELY-IMPORTANT 块的遗漏同步而导致指令冲突
- 新 skill 编写时，有清晰的范式可遵循：指令性语言、引用规则文件、不描述 CLI 行为
- quick-tasks 和 breakdown-tasks 各自独立完整，agent 只需读取其中一个就能完整执行

### Constraints & Dependencies

- quick-tasks 和 breakdown-tasks 必须各自独立自洽，不抽取共享 rule 文件
- execute-task.md 和 run-tasks.md 同理保持独立
- 修改不能改变任何 skill 的外部行为（输入/输出/副作用）
- 遵循 forge-distribution.md 的路径规范（相对路径、不使用源码根路径）
- 依赖 CLI 行为在审计与修复期间保持稳定：若 CLI 的 exit code 或输出格式在此期间变更，"行为解释"分类可能需要重新评估。风险极低（v3.0.0 发布期间 CLI 接口冻结）

## Alternatives & Industry Benchmarking

### Industry Context

本提案的核心原则——"指令性而非描述性"——与多个行业实践一致：

- **Anthropic prompt engineering guidelines**（2024）：明确建议"Be clear and direct"、"Use examples rather than abstract definitions"，与本次提案删除行为解释、保留命令+输出的策略一致。
- **Cursor/Windsurf 的 `.cursorrules` 实践**：社区最佳实践强调指令文件应只包含"what to do"而非"how the system works"，避免在 agent 上下文中注入实现细节。
- **OpenAI 的 GPT best practices**（structured outputs）：推荐将复杂任务分解为声明式的步骤列表，而非描述性的解释段落。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零成本 | agent 执行准确率持续受损，维护成本随时间增长 | Rejected: 缺陷已量化，继续累积不可接受 |
| 按 issue 类型批量修复 | 本次审计 + 行业实践 | 分组清晰，每批同类修复便于 review，与 Anthropic "imperative instructions" 原则对齐 | 跨文件上下文切换多 | **Selected: 修复逻辑一致，便于验证** |
| 按文件逐个修复 | — | 每个文件一次完成所有修复 | 同类问题在不同文件中修复方式可能不一致 | Rejected: 容易导致修复风格不统一 |
| 自动化 lint 校验 | API 文档标准（OpenAPI spec 的 "description vs schema" 分离） | 用 CI 规则自动检测新增的 CLI 行为描述或 E-I 冗余，防止回归 | 需要定义规则语法（如正则匹配 "What .* Does" 标题、E-I 块与正文重复检测），开发成本高 | Deferred: 作为后续独立提案，本次修复的文件集可作为 lint 规则的训练数据 |

## Feasibility Assessment

### Technical Feasibility

纯文本修改，无代码变更，无依赖风险。所有修改都是删除或简化现有文本，不引入新逻辑。

### Resource & Timeline

约 95 处修改分布在 ~40 个文件中。预计 8-12 个 coding task（按 skill/command 分组）。每个 task 预计修改 3-5 个文件，review 负担约 10-15 分钟/task（对照边界规则表和子分类验证方法逐项确认）。总 review 时间约 1.5-3 小时。

## Scope

### In Scope

- 删除 22 处 CLI 行为描述，替换为指令性操作
- 精简 33 处冗余描述（EXTREMELY-IMPORTANT 块、frontmatter 重复、规则文件预览）
- 修复 40 处清晰度/自洽性问题（流程编号、误导引用、歧义步骤——具体实例包括 `tech-design/SKILL.md` Process Flow Step 8→11 跳跃、`run-tests/SKILL.md` Step 5 对 Step 0 的误导引用、`gen-contracts/SKILL.md` 对不存在 Section 编号的引用等）
- 确保 quick-tasks 和 breakdown-tasks 各自内部自洽
- 确保 execute-task 和 run-tasks 各自内部自洽

### Out of Scope

- 跨文件去重（quick-tasks↔breakdown-tasks、execute-task↔run-tasks 保持独立）
- 功能性变更（不改变任何 skill 的输入/输出/副作用）
- 新 skill 或新 command 的创建
- hooks/ 目录的修改
- forge CLI 源码修改
- fix-bug 的 Knowledge Review section 抽取（独立提案处理）
- **回归预防机制**（CI lint 规则自动检测新增的 CLI 行为描述或 E-I 冗余）。理由：lint 规则的开发是独立工程任务，需要设计规则语法、误报容忍度、CI 集成方案，超出本次纯文本修复的范围。本次修复的文件集可作为后续 lint 规则的 golden test 数据。

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 删除 CLI 描述后 agent 丢失必要上下文 | L | M | 保留 exit code 契约（0=成功/1=失败）和输出字段名列表，只删除语义解释 |
| EXTREMELY-IMPORTANT 精简后遗漏关键约束 | M | H | 每个文件的精简结果逐一与正文步骤比对。具体流程：(1) 提取 E-I 块每条约束的 key verb（MUST/NEVER/ALWAYS），(2) grep 正文是否包含该 verb 对应的约束，(3) 未匹配的约束标记为"保留"并注明理由 |
| 流程编号重排引入新错误 | M | M | 重排后验证 Process Flow 与实际步骤编号一一对应 |
| 审计分类错误：行为描述被误标为输出契约（或反之） | M | H | 执行删除前，每个文件由执行者按"边界规则表"逐条标注分类，在 task 的 PR 描述中列出标注结果。reviewer 对标注结果进行二次确认 |
| 40 文件批量修改的回归风险 | M | M | 修改不涉及代码/逻辑变更，回归风险仅限于指令文字。每个 task 完成后执行该 skill 的 dry-run（读文件验证无语法错误、无断裂引用），不运行完整测试套件 |
| 并行 task 间的一致性风险 | L | M | 选定"按 issue 类型批量修复"而非"按文件逐个修复"，确保同类修复由同一 task 处理，减少跨 task 的风格不一致 |

## Success Criteria

- [ ] 22 处 CLI 行为描述全部删除，替换为指令性操作。验证方法三层：(1) grep 无 "What .* Does" section 标题；(2) 每个被修改文件的 diff 中无输出契约字段名丢失（人工 spot-check：随机抽取 5 个被修改文件，按"删除边界规则"表中的第二类逐条比对，确认所有 Output Contract 字段名保留）；(3) 校准示例：`execute-task.md` 删除后保留 exit code 契约和 `SURFACE_KEY` 字段名，`submit-task/SKILL.md` 删除后仅保留 `forge task submit` 命令本身
- [ ] EXTREMELY-IMPORTANT 块精简完成，每个保留的 E-I 条目通过约束级别审计：满足以下之一 (a) 正文未包含该约束，或 (b) 正文包含该约束但 E-I 版本的强制等级更高（例如正文仅建议而 E-I 为 MUST/NEVER）
- [ ] tech-design Process Flow 包含完整步骤（Step 8→9→10→11 无跳跃）
- [ ] run-tests 无误导引用（"Convention loaded in Step 0" 不存在）
- [ ] quick.md 配置读取失败时 fallback 为显示确认门（非跳过）
- [ ] 所有 skill/command 文件的 frontmatter description 与正文第一句不重复
- [ ] quick-tasks 和 breakdown-tasks 各自包含完整的独立指令集，agent 只读其一即可执行。验证方法：(1) 提取 quick-tasks 文件中引用的所有外部文件路径，确认不引用 breakdown-tasks 目录下的任何文件（反之亦然）；(2) 每个 skill 的步骤链中无 "see breakdown-tasks/quick-tasks for..." 形式的交叉引用
- [ ] execute-task 和 run-tasks 各自内部自洽：提取各自引用的外部文件路径，确认 execute-task 不依赖 run-tasks 的步骤定义（反之亦然）；每个文件的步骤编号连续无跳跃

## Next Steps

- Proceed to `/quick-tasks` to generate task breakdown
