---
created: "2026-05-27"
author: "forge-brainstorm"
status: Draft
---

# Proposal: 精简任务 Prompt 模板

## Problem

Forge 的 15 个任务 prompt 模板包含大量非指令内容——注释、解释性描述、冗长的角色定义——这些内容不指导 agent 行为，只增加 token 消耗并稀释指令清晰度。同时 task-executor agent 的 Execution Protocol 存在步骤冗长、逻辑重叠的问题。

### Evidence

模板中冗余内容的量化分析：

| 冗余类别 | 出现文件数 | 每处损失 | 总冗余 |
|---------|-----------|---------|-------|
| HTML 注释 | 1 | 4 行 | 4 行 |
| Step 2 解释性描述 | 5 (gen-* / test-run / verify) | 1-2 行 | ~7 行 |
| 冗长角色描述 | 10 (coding.*, gate, doc) | 1 行 | ~10 行 |
| CODING_PRINCIPLES 解释性冗余 | 5 (coding.*) | 每原则 2-5 行 | ~50 行 |
| AC 验证块冗余 | 9 (coding.*, gate, doc) | 每处 ~12 行可缩至 ~3 行 | ~80 行 |
| Record Fields 描述性文字 | 9 | ~3 行可缩至 ~1 行 | ~20 行 |
| task-executor Execution Protocol | 1 | 11 步可合并为 8 步 | ~30 行 |

总计：约 **200 行** 非指令冗余。每执行一个任务，agent 都要阅读这些无用 token，形成累积开销。

### Urgency

- 每个 task 执行都在消耗这些冗余 token，日积月累规模可观
- 清晰的 prompt 减少 agent 误解和执行偏差
- Prompt 精简是持续优化的一部分，目前已有 prompt-template-audit 等基础，可以在此基础上推进

## Proposed Solution

**就地精简**：保持现有模板独立，在每个文件内部删除非指令内容，将模糊描述改为清晰指令。

不抽取公共模块，不改变现有分类体系。

### Innovation Highlights

本方案不是技术创新，而是对现有 prompt 的"清理"。核心原则是"prompt 是指令，不是文档"——删掉所有不能直接指导 agent 行动的文字。

## Requirements Analysis

### Key Scenarios

1. **coding-feature / coding-enhancement / coding-fix / coding-cleanup / coding-refactor** 五个核心模板：
   - 角色描述从自然语言改为祈使句
   - CODING_PRINCIPLES 去掉举例和解释，保留核心约束
   - AC 验证块从 ~12 行精简到 ~3 行
   - Step 2 的实现说明保留，只去修饰性语言

2. **gate / doc** 模板：
   - 角色描述精简
   - AC 验证块精简

3. **test-run / test-gen-scripts / test-gen-contracts / test-gen-journeys / test-verify-regression** 模板：
   - Step 2 中的 "This generates X from Y..." 解释性描述删除
   - 角色描述精简

4. **task-executor agent**：
   - Execution Protocol 步骤合并（步骤 4/5/6 处理 prompt 获取逻辑可合并为 1 步）
   - Retry Strategy 与 Complex Error Pause Flow 去重合并
   - 输出格式合并为紧凑格式

### Non-Functional Requirements

- 精简后模板的指令覆盖必须与精简前等价（不能遗漏 agents 需要知道的信息）
- 所有 task-executor 的行为不发生变化

### Constraints & Dependencies

- 模板文件位于 `forge-cli/pkg/prompt/data/*.md`，由 `prompt.go` 通过 embed FS 加载
- 修改模板不影响 Go 代码，只需修改 .md 文件
- task-executor agent 位于 `plugins/forge/agents/task-executor.md`

## Alternatives & Industry Benchmarking

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| 什么都不做 | — | 零风险 | token 持续浪费、指令不够清晰 | Rejected: 成本太低 |
| 抽取公共模块 | DRY 模式 | 修改一处同步所有模板 | 需要改 `prompt.go` 逻辑，且被用户否决 | Rejected: 不满足就地要求 |
| 引入 DSL 生成 | 模板引擎 | 最灵活 | 过度工程化 | Rejected: 小题大做 |
| **就地精简** | Forge 现有风格 | 零架构变更，每模板独立修改，风险隔离 | 每个文件都要改 | **Selected: 简单直接** |

## Feasibility Assessment

### Technical Feasibility

纯文本编辑，无技术风险。

### Resource & Timeline

10-15 个文件的文字精简，1 次编码任务即可完成。

### Dependency Readiness

前置条件：本次 brainstorm 输出的 proposal 通过。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| "我需要保留角色描述让 agent 理解上下文" | Assumption Flip | 角色描述中的自然语言("You are a focused...")对 LLM 行为的影响可以通过祈使句替代。Agent 的执行行为由后续的 Workflow 步骤定义，不是由角色描述定义的。 |
| "每个模板独立意味着不需要关注跨模板一致性" | XY Detection | 用户确认了「核心流程重复是允许的」，所以跨模板一致性不是问题，不需要抽取公共模块。 |

## Scope

### In Scope
- 修改 `forge-cli/pkg/prompt/data/` 下全部 15 个模板文件：
  - coding-feature.md, coding-fix.md, coding-enhancement.md, coding-cleanup.md, coding-refactor.md
  - gate.md, doc.md
  - test-run.md, test-gen-scripts.md, test-gen-contracts.md, test-gen-journeys.md, test-verify-regression.md
  - code-quality-simplify.md, validation-code.md, validation-ux.md
- 修改 `plugins/forge/agents/task-executor.md`
- 删除 HTML 注释
- 删除 Step 2 解释性描述（"This generates X from Y"）
- 精简角色描述（自然语言 → 祈使句）
- 精简 AC 验证块（~12 行 → ~3 行）
- 精简化 Record Fields（去掉引导性描述，保留字段名和值）
- 精简 CODING_PRINCIPLES（去掉举例和解释）
- 精简 task-executor Execution Protocol（合并步骤）

### Out of Scope
- 不抽取公共模块文件
- 不修改 `prompt.go` 代码逻辑
- 不新增/删除模板文件
- 不改动模板占位符（`{{TASK_ID}}` 等）
- 不改动 Spec Authority Enforcement 逻辑（保留但可考虑精简其行数）
- 不改动 Hard Rules / CRITICAL 块的逻辑

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 精简过度导致 agent 遗漏关键行为 | Low | High | 每个模板修改后对比：所有功能点是否仍被覆盖；task-executor 的每个步骤的行为约束是否保持 |
| 多个模板同步修改，跨模板不一致 | Medium | Medium | 每个模板在修改时以 coding-feature 为基准对齐 |

## Success Criteria

- [ ] 15 个模板文件 + task-executor 共减少 **≥150 行**（去除注释、解释性描述、冗长定义）
- [ ] 所有模板精简后，agent 执行相同 task 的行为无可见差异
- [ ] 所有 Spec Authority Enforcement、Hard Rules 等关键约束块保留，逻辑不变
- [ ] task-executor 的 Execution Protocol 步骤数从 11 步减少到 ≤8 步

## Next Steps

- Proceed to `/write-prd` to formalize requirements