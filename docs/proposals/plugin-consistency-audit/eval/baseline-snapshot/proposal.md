---
created: "2026-05-29"
author: "faner"
status: Draft
intent: "refactor"
---

# Proposal: Forge Plugin 内部一致性审计

## Problem

v3.0.0 经历 test profile system 和 intent-driven pipeline branching 等大规模重构后，22 个 skill、18 个 command、1 个 agent 的 SKILL.md 与其各自的 templates/rules/data 文件之间可能存在指令矛盾、冗余信息或时序问题，且缺乏系统性验证。

### Evidence

- test profile 系统将 Playwright 硬编码替换为可插拔 profile，影响了 gen-test-scripts、run-tests、init-justfile 等多个 skill 的 rules 和 templates
- intent-driven branching 引入 `new-feature/refactor/cleanup` 三路分支，影响 breakdown-tasks、quick-tasks、run-tasks 等核心编排 skill
- 172+ 个 .md 文件通过手动维护交叉引用，重构过程中依赖局部修改而非全局验证

### Urgency

v3.0.0 发版在即，内部不一致会导致运行时行为异常（流程卡死、模板字段缺失、步骤时序错乱），发版后修复成本远高于现在。

## Proposed Solution

对 forge plugin 所有组件进行**内部逻辑自洽性审计**：逐一检查每个 skill 的 SKILL.md 与其 templates/rules/data 之间、每个 command 的内部流程、以及 agent 的指令之间是否存在矛盾、冗余或时序问题。输出结构化问题报告（含文件路径、问题描述、严重等级、修复建议），不做实际修复。

### Innovation Highlights

审计按"单一组件自洽"而非"跨组件协调"组织——这降低了审计复杂度，同时覆盖了最可能出问题的维度（组件内部重构后的残留不一致）。跨组件冗余是设计层面的合理重复，不在此次审计范围内。

## Requirements Analysis

### Key Scenarios

- SKILL.md 描述的步骤流程与 template 中假设的字段/结构不一致
- rules 文件中的约束条件与 SKILL.md 中的指令矛盾
- SKILL.md 引用的 template/rule 文件路径不存在或已过时
- 同一 skill 内重复描述同一行为（SKILL.md 和 rules 各说一遍）
- Command 内部流程步骤时序错误（如先读后写、先验证后检查）

### Non-Functional Requirements

- 审计覆盖率: 100% 的 skill（22个）、command（18个）、agent（1个）
- 问题分类: 矛盾(CONFLICT)、冗余(REDUNDANT)、时序(TIMING)、引用(REFERENCE)

### Constraints & Dependencies

- 审计基于当前 v3.0.0 分支代码，不依赖运行时测试
- 不修改任何代码，仅输出报告

## Alternatives & Industry Benchmarking

### Industry Solutions

大型 prompt-based system 通常通过 schema 验证和 lint 工具检查一致性，但 forge plugin 的 prompt 文件是自由格式的 markdown，无 schema 约束。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零成本 | 隐式不一致会在使用中暴露 | Rejected: 发版前风险不可接受 |
| 人工逐文件审读 | 传统方法 | 细致 | 172 文件量巨大，易遗漏 | Rejected: 效率太低 |
| 自动化 schema 验证 | 行业标准 | 可重复执行 | 需先定义 schema，成本高 | Rejected: ROI 不够 |
| **AI 辅助分层审计** | 本次方案 | 覆盖全面、可理解上下文语义 | 可能遗漏隐含假设 | **Selected: 最适合当前规模和紧迫性** |

## Feasibility Assessment

### Technical Feasibility

所有文件均为 markdown，可完整读取和分析。无技术障碍。

### Resource & Timeline

单次审计，预计产出 1 份结构化报告。

### Dependency Readiness

所有文件均在当前仓库中，无需外部依赖。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| 重构后的文件已经内部一致 | Assumption Flip: 假设不一致，逐一验证 | Refined: 重构是局部修改，残留不一致很可能存在 |
| 跨 skill 冗余是问题 | Occam's Razor: 跨 skill 重复可能是设计意图 | Confirmed: 用户确认忽略跨 skill 冗余，聚焦组件内部自洽 |

## Scope

### In Scope

- 22 个 skill 的 SKILL.md 与其各自的 templates/rules/data 之间的逻辑自洽性
- 18 个 command 的内部流程一致性
- 1 个 agent (task-executor) 的内部指令一致性
- hooks/guide.md 的内部一致性
- 问题分类: 矛盾(CONFLICT)、冗余(REDUNDANT)、时序(TIMING)、引用(REFERENCE)

### Out of Scope

- 跨 skill 之间的冗余内容（设计层面的合理重复）
- rules/rubrics/experts 的功能性质量审查
- Forge CLI Go 源码
- 用户项目目录结构
- 实际代码修复（仅产出报告）

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 语义层面的隐含矛盾难以通过文本对比发现 | M | M | 重点检查关键词不一致（如"必须"vs"可选"）、条件分支覆盖 |
| 报告问题过多导致修复优先级不清 | L | M | 每个问题标注严重等级（P0-P3） |
| 审计过程中遗漏某些文件 | L | L | 使用文件清单逐一勾选 |

## Success Criteria

- [ ] 22 个 skill 100% 覆盖审计，每个 skill 的 SKILL.md 与其 templates/rules/data 逐一对比
- [ ] 18 个 command 100% 覆盖审计
- [ ] 1 个 agent (task-executor) 完成审计
- [ ] 输出结构化问题报告，每个问题包含: 文件路径、问题描述、严重等级(P0-P3)、修复建议
- [ ] 问题按 CONFLICT/REDUNDANT/TIMING/REFERENCE 四类分类
