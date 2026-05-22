---
created: 2026-05-22
author: "faner"
status: Approved
---

# Proposal: 重构执行器的预期影响范围声明

## Problem

重构任务执行器在遇到定标测试（characterization tests）失败时，无法判断这是意图性的行为变化还是非预期的回归，导致执行器在两个选择之间反复横跳并最终卡住。

### Evidence

`docs/lessons/gotcha-characterization-test-vs-refactoring.md` 记录了 fix.7 卡死的完整因果链：task 2.8 消除了 SourceTaskID sentinel，改变了 `--block-source` 的行为。fix.7 被要求更新 4 个定标测试，但无法判断"source 1.1 should be blocked"断言的失败是预期变化还是 bug。

当前 `coding-refactor.md` 模板的 `BEHAVIOR_CHANGE_DETECTED` 机制是反应式的——检测到行为变化后跳过，不区分"预期变化"和"非预期回归"。

### Urgency

大规模重构（如 v3.0.0 的 architecture-simplification）中，行为改变型重构是常见需求。当前机制迫使重构任务"假装行为不变"，遇到行为变化就跳过，然后创建 fix 任务去处理被跳过的变化。fix 任务又因为缺乏上下文而卡住。这是一个结构性缺陷，每次大规模重构都会触发。

## Proposed Solution

在重构执行模板的 Impact Mapping 步骤后新增 **Impact Declaration** 子步骤：执行器在动手修改代码前，先分析并声明哪些测试的行为预期会发生变化（EVOLVE），哪些必须保持不变（PRESERVE）。

重构过程中，EVOLVE 测试的失败被视为预期变化，执行器直接更新测试断言。PRESERVE 测试的失败仍触发暂停。未出现在声明中的测试失败也触发暂停（安全网）。

### Innovation Highlights

这是对传统重构定义（"不改变外部行为"）的务实修正。行业标准的重构定义假设行为不变，但在实际的大型代码库演进中，"重构"经常包含有意图的行为简化。此方案在保留安全网（PRESERVE + 未声明 = 暂停）的同时，为 EVOLVE 类测试提供了显式通道。

灵感来自 Michael Feathers 的《Working Effectively with Legacy Code》中 characterization test 的概念：定标测试记录的是"当前行为"而非"正确行为"，当意图性重构改变行为时，测试应随之演进。

## Requirements Analysis

### Key Scenarios

- **Happy path**: 重构任务声明 3 个 EVOLVE 测试和 10 个 PRESERVE 测试。重构完成后，EVOLVE 测试被更新，PRESERVE 测试全部通过。
- **声明不完整**: 执行器声明了 3 个 EVOLVE 测试，但重构实际影响了第 4 个未声明的测试。该测试失败触发暂停，创建 fix 任务。安全网生效。
- **纯结构重构**: 重构只涉及重命名/移动，无行为变化。声明中所有测试为 PRESERVE，与当前行为一致。
- **大规模重构（>50 个测试受影响）**: 执行器声明大量 EVOLVE 测试。声明本身提供了影响范围的文档化记录，便于审查。

### Non-Functional Requirements

- **声明开销**: Impact Declaration 分析步骤应在重构总时间的 10% 以内。对于简单重构，声明应简洁。
- **可追溯性**: 声明必须作为执行记录的一部分输出，可事后审计。

### Constraints & Dependencies

- 仅修改 `forge-cli/pkg/prompt/data/coding-refactor.md` 模板文件
- 不修改 CLI 代码、任务生成器或其他模板
- 声明格式为纯文本输出，不引入新的结构化存储

## Alternatives & Industry Benchmarking

### Industry Solutions

Martin Fowler 的重构目录定义重构为"不改变可观察行为的代码变换"。Michael Feathers 的 characterization test 模式建议：当需要改变行为时，先写新的测试定义期望行为，再将 characterization test 降级为回归测试。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零成本 | 每次大规模重构都会卡死 | Rejected: 已有两次卡死记录 |
| 任务描述预填影响范围 | Task generator | 创建者意图明确 | 任务生成器无法准确预测影响范围；依赖人类自觉 | Rejected: 自动化程度低 |
| 拆分为分析+执行两任务 | 分阶段执行 | 天然检查点 | 增加调度复杂度，两个任务间状态传递困难 | Rejected: 过度工程化 |
| **执行器自分析+声明** | 本方案 | 执行者即分析者，信息最充分；声明即文档 | 执行器分析可能不完整 | **Selected: 安全网兜底不完整的情况** |

## Feasibility Assessment

### Technical Feasibility

完全可行。修改单个 markdown 模板文件，不涉及代码变更。

### Resource & Timeline

1 个 coding 任务（修改模板）+ 1 个 doc 任务（更新 lesson），预估 <1h。

### Dependency Readiness

无外部依赖。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| 执行器能准确预测重构影响范围 | Assumption Flip: 如果预测不完整会怎样？ | Confirmed: 不完整时声明范围外的测试失败仍触发暂停，安全网兜底。最差情况是多创建几个 fix 任务。 |
| 重构 = 不改变行为 | 5 Whys: 为什么行为变化会出现在"重构"任务中？ | Overturned: 实践中"重构"经常包含行为简化（如消除 sentinel 值）。模板需要承认这个现实。 |
| BEHAVIOR_CHANGE_DETECTED 够用 | Stress Test: 大规模重构中有 10+ 个行为变化点 | Refined: 反应式检测在单点变化时可以，多点变化时 fix 任务链会卡死。需要前置声明。 |

## Scope

### In Scope

- 修改 `coding-refactor.md`：在 Step 2 (Impact Mapping) 中新增 Impact Declaration 子步骤
- 修改 `coding-refactor.md` Step 3 和 Step 4：EVOLVE 测试失败时的处理逻辑
- Impact Declaration 的输出格式定义

### Out of Scope

- 调度器预编译检查（参考 `gotcha-pre-existing-syntax-errors-block-executor`）
- CLI 代码变更
- 其他执行模板（coding-feature.md, coding-fix.md）
- 任务生成器变更
- 结构化存储 Impact Declaration

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 执行器过度声明 EVOLVE（将所有测试标为 EVOLVE 以避免暂停） | M | M | 在模板中加约束：EVOLVE 必须附带 reason 和 expected_change，空理由的 EVOLVE 视为无效 |
| 分析步骤增加执行时间 | L | L | 分析是读代码+grep，不涉及编译，开销在秒级 |
| 声明格式不被执行器严格遵循 | M | M | 用结构化模板格式 + 示例，降低遵循门槛 |

## Success Criteria

- [ ] 重构执行器在 Step 2 后输出结构化的 IMPACT_DECLARATION（包含 PRESERVE/EVOLVE 分类）
- [ ] EVOLVE 分类的测试失败时，执行器自动更新测试断言而非输出 BEHAVIOR_CHANGE_DETECTED
- [ ] 未声明或 PRESERVE 分类的测试失败时，仍触发暂停并创建 fix 任务
- [ ] Impact Declaration 作为执行输出的一部分可追溯（在任务执行记录中可见）
