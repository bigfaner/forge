---
created: "2026-05-28"
author: "faner"
status: Draft
---

# Proposal: Forge Plugin Skill 指令层全面审计与修复

## Problem

Forge plugin 的 22 个 skill、19 个 command、1 个 agent 定义中，指令层存在三类系统性缺陷：约 22 处描述了 forge CLI 的内部行为而非直接使用命令，约 33 处冗余描述（文件内部重复、EXTREMELY-IMPORTANT 与正文重复），约 40 处指令清晰度或自洽性问题（流程跳跃、误导引用、歧义步骤）。AI agent 在执行时因这些缺陷产生认知负担、做出错误判断、或遗漏关键步骤。

### Evidence

通过逐行审计全部 170+ 文件，发现以下系统性模式：

1. **CLI 行为描述**：`execute-task.md` 解释 `forge task claim` 的输出字段含义和默认值；`breakdown-tasks/SKILL.md` 描述 `forge task index` 在不同场景下的内部分支逻辑；`quick-tasks/SKILL.md` 同样描述 `forge task index` 的实现细节；`submit-task/SKILL.md` 用整个 section 解释 `forge task submit` 做了什么。
2. **冗余描述**：`test-guide/SKILL.md` 的 EXTREMELY-IMPORTANT 块 10 条中 8 条与正文步骤重复；`quick-tasks` ↔ `breakdown-tasks` 之间 12 处近乎逐字重复；`clean-code.md`、`git-commit.md`、`git-checkout.md`、`init-forge.md` 的 frontmatter description 与正文第一句完全重复。
3. **清晰度/自洽性**：`tech-design/SKILL.md` Process Flow 从 Step 8 跳到 Step 11，跳过 Step 9-10；`run-tests/SKILL.md` Step 5 引用 "Convention loaded in Step 0" 但 Step 0 实际是 Stale State Recovery；`quick.md` 配置读取失败时 fallback 为跳过确认门（逻辑倒置）；`gen-contracts/SKILL.md` 引用不存在的 Section 编号。

### Urgency

v3.0.0 正在发布中。指令层缺陷直接影响 AI agent 的执行准确率——每次 agent 因歧义或过时描述走偏，都浪费一个完整的执行循环（30 分钟 subagent timeout）。修复越晚，已有过时指令导致的执行失败越多。

## Proposed Solution

对全部 skill/command/agent 文件执行三类修复：
1. **删除 CLI 行为描述**：只保留"运行 X 命令"的指令，删除对 CLI 输出语义、内部实现、分支逻辑的解释
2. **删除冗余描述**：精简 EXTREMELY-IMPORTANT 块（只保留正文未覆盖的跨步骤约束）、删除 frontmatter 与正文的重复、删除规则文件预览（直接引用规则文件）
3. **修复清晰度/自洽性**：修正流程编号跳跃、消除误导引用、明确歧义步骤的前提条件

### Innovation Highlights

此提案遵循一个核心洞察：**AI agent 的指令应是指令性的，不是描述性的**。agent 不需要理解 CLI 的内部机制来正确使用它，也不需要在 EXTREMELY-IMPORTANT 块中读到与正文相同的规则两次。这与人类文档的最佳实践（重复强化记忆）相反——对 AI 来说，重复增加的是不一致风险而非记忆强化。

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

## Alternatives & Industry Benchmarking

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零成本 | agent 执行准确率持续受损，维护成本随时间增长 | Rejected: 缺陷已量化，继续累积不可接受 |
| 按 issue 类型批量修复 | 本次审计 | 分组清晰，每批同类修复便于 review | 跨文件上下文切换多 | **Selected: 修复逻辑一致，便于验证** |
| 按文件逐个修复 | — | 每个文件一次完成所有修复 | 同类问题在不同文件中修复方式可能不一致 | Rejected: 容易导致修复风格不统一 |

## Feasibility Assessment

### Technical Feasibility

纯文本修改，无代码变更，无依赖风险。所有修改都是删除或简化现有文本，不引入新逻辑。

### Resource & Timeline

约 95 处修改分布在 ~40 个文件中。预计 8-12 个 coding task（按 skill/command 分组）。

## Scope

### In Scope

- 删除 22 处 CLI 行为描述，替换为指令性操作
- 精简 33 处冗余描述（EXTREMELY-IMPORTANT 块、frontmatter 重复、规则文件预览）
- 修复 40 处清晰度/自洽性问题（流程编号、误导引用、歧义步骤）
- 确保 quick-tasks 和 breakdown-tasks 各自内部自洽
- 确保 execute-task 和 run-tasks 各自内部自洽

### Out of Scope

- 跨文件去重（quick-tasks↔breakdown-tasks、execute-task↔run-tasks 保持独立）
- 功能性变更（不改变任何 skill 的输入/输出/副作用）
- 新 skill 或新 command 的创建
- hooks/ 目录的修改
- forge CLI 源码修改
- fix-bug 的 Knowledge Review section 抽取（独立提案处理）

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 删除 CLI 描述后 agent 丢失必要上下文 | L | M | 保留 exit code 契约（0=成功/1=失败）和输出字段名列表，只删除语义解释 |
| EXTREMELY-IMPORTANT 精简后遗漏关键约束 | M | H | 每个文件的精简结果逐一与正文步骤比对，确保无遗漏 |
| 流程编号重排引入新错误 | M | M | 重排后验证 Process Flow 与实际步骤编号一一对应 |

## Success Criteria

- [ ] 22 处 CLI 行为描述全部删除，替换为指令性操作（grep 无 "What .* Does" section）
- [ ] EXTREMELY-IMPORTANT 块中与正文重复的条目全部删除（每个 E-I 块条目在正文中无对应）
- [ ] tech-design Process Flow 包含完整步骤（Step 8→9→10→11 无跳跃）
- [ ] run-tests 无误导引用（"Convention loaded in Step 0" 不存在）
- [ ] quick.md 配置读取失败时 fallback 为显示确认门（非跳过）
- [ ] 所有 skill/command 文件的 frontmatter description 与正文第一句不重复
- [ ] quick-tasks 和 breakdown-tasks 各自包含完整的独立指令集，agent 只读其一即可执行

## Next Steps

- Proceed to `/quick-tasks` to generate task breakdown
