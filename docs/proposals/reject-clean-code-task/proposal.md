---
created: 2026-05-15
author: "faner"
status: Draft
supersedes: clean-code-task
---

# Proposal: Reject Clean-Code Task Type — Use /simplify Instead

## Problem

已有的 `clean-code-task` 提案（`docs/proposals/clean-code-task/`）建议在 forge 任务流水线中新增 `clean-code` 任务类型，用于在 verify-regression 和 consolidate-specs 之间自动清理代码。经过分析，该提案所解决的问题可由 Claude Code 内置的 `/simplify` skill 完全覆盖，无需引入新的任务类型。

### Evidence

**功能对比：**

| 清理目标 | `/simplify` | clean-code 提案 |
|----------|-------------|-----------------|
| Dead imports | ✅ "reuse, quality" 涵盖 | ✅ |
| Commented-out code | ✅ "efficiency" 涵盖 | ✅ |
| 重复逻辑 | ✅ "reuse" 涵盖 | ✅ |
| 命名不一致 | ✅ "quality" 涵盖 | ✅ |
| 作用域限定 | 用户指定或当前 diff | record-driven 自动限定 |

clean-code 提案唯一的差异化是 **feature scope 自动限定**（通过 records 确定变更文件列表），但这只是 `git diff --name-only` 级别的工作，用户在调用 `/simplify` 时已经知道该审查哪些文件。

**提案宣称的 pipeline 时序保证（consolidate-specs 之前清理）被高估了：** unused imports、commented-out code 等问题不影响业务规则和技术规范的提取质量。consolidate-specs 关注的是接口契约、数据模型、流程逻辑等结构化信息，这些不受上述代码卫生问题影响。

### Urgency

应尽快关闭该提案，避免投入工程资源实现一个不必要的任务类型。已有提案已通过 3 轮评估（最终 927/1000），若不明确关闭，后续可能被误认为待实现。

## Proposed Solution

1. **标记已有提案为 Superseded** — 在 `docs/proposals/clean-code-task/proposal.md` 头部更新 status
2. **更新 todo.txt 第 43 条** — 标记为 done，注明使用 `/simplify` 替代
3. **添加使用引导** — 在 forge guide 或相关文档中提示用户在 feature 完成后可使用 `/simplify` 进行代码清理

## Requirements Analysis

### Key Scenarios

1. **Feature 完成后清理** — 用户在所有 business tasks 和 test tasks 完成后，手动执行 `/simplify`，审查本 feature 变更文件的代码质量
2. **consolidate-specs 前清理** — 用户希望 specs 基于干净代码提取，在执行 `/consolidate-specs` 前先执行 `/simplify`
3. **发现累积问题** — 多 task 并行实现后，用户发现代码质量问题，主动执行 `/simplify` 清理

### Constraints & Dependencies

- `/simplify` 是 Claude Code 内置 skill，无需额外开发
- 不需要修改 forge-cli 的任务类型系统

## Alternatives & Industry Benchmarking

### Comparison Table

| Approach | Pros | Cons | Verdict |
|----------|------|------|---------|
| 实现 clean-code 任务类型 | pipeline 自动化；record-driven scope | 与 `/simplify` 高度重叠；增加引擎复杂度（4 Go 文件 + 1 模板 + 测试）；pipeline 语义被非确定性任务污染 | Rejected: 功能重叠 |
| **使用 `/simplify` 替代** | 零开发成本；内置可用；功能覆盖完整 | 手动触发；无 pipeline 时序保证 | **Selected: ROI 最优** |
| 增强 `/simplify` 增加 feature scope | 保留自动作用域；比新建类型轻量 | `/simplify` 是 Claude Code 内置 skill，不可修改 | Rejected: 不可控 |
| 在 all-completed hook 中自动调用 `/simplify` | 保留 pipeline 集成 | hook 是非交互的；改变 hook 模型 | Rejected: 模型不匹配 |

## Feasibility Assessment

### Technical Feasibility

无需技术实现。仅需文档更新。

### Resource & Timeline

10 分钟完成。

## Scope

### In Scope

- 标记 `docs/proposals/clean-code-task/proposal.md` 为 Superseded
- 创建本 rejection 提案，记录决策理由
- 更新 todo.txt 第 43 条

### Out of Scope

- 修改 `/simplify` 的行为或增强其能力
- 在 forge pipeline 中集成代码清理步骤
- 修改 all-completed hook

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 用户忘记在 feature 结束时执行 `/simplify` | M | L | 在引导文档中明确建议；consolidate-specs 的提取质量不受代码卫生问题影响 |
| 未来发现 `/simplify` 无法覆盖特定清理需求 | L | L | 届时重新评估，可基于本提案的决策记录快速启动 |

## Success Criteria

- [x] `docs/proposals/clean-code-task/proposal.md` status 更新为 Superseded
- [x] 本提案完整记录了拒绝 clean-code 任务类型的理由
- [x] todo.txt 第 43 条更新为使用 `/simplify` 替代
