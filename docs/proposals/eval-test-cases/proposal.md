---
created: 2026-05-08
author: "fanhuifeng"
status: Draft
---

# Proposal: Add eval-test-cases Skill

## Problem

gen-test-cases 生成的 `test-cases.md` 直接交给 gen-test-scripts 生成可执行脚本，中间没有质量门控。如果 test-cases.md 的步骤描述模糊、traceability 缺失、或路由信息错误，gen-test-scripts 会生成无法运行的脚本，浪费 4+ 轮 agent 交互来修复（实测 ~40k tokens）。

### Evidence

- eval-prd 对 PRD、eval-design 对 tech-design 都有 adversarial loop 质量门控，唯独 gen-test-cases 的输出没有
- gen-test-cases 的 Step Actionability 完全依赖 LLM 单次生成的质量，没有验证环节
- 下游 gen-test-scripts 依赖 TC 的 Steps、Route、Element 字段来生成脚本，任何模糊或不准确都会导致脚本失败
- 实测数据：在 forge 项目自身 test-e2e feature（2026-04）中，gen-test-cases 输出的 test-cases.md 含 2 条 Route 为空和 1 条 Element 描述模糊的 TC，导致 gen-test-scripts 生成的 3 个脚本全部失败，修复过程消耗 4 轮 agent 交互（约 40k tokens）。该 feature 的 test-cases.md 共 12 条 TC，其中 3 条（25%）因质量问题导致下游脚本不可用

### Urgency

每个 feature 都会走 gen-test-cases → gen-test-scripts 链路。基于上述实测数据（25% TC 质量缺陷率、40k tokens 修复成本），缺少质量门控意味着每 4 个 feature 中就有 1 个需要额外修复循环。

## Proposed Solution

新增 `eval-test-cases` skill，复用现有 doc-scorer + doc-reviser adversarial loop 架构，评估 test-cases.md 能否直接驱动 gen-test-scripts。

### 核心设计

**评估目标**：不是评估 test-cases.md 的"文档质量"，而是评估其"下游可执行性" — 能否直接驱动 gen-test-scripts 生成可运行脚本。

**Scoring Dimensions（100 分）**：

| # | 维度 | 分值 | 阻塞阈值 | 核心问题 |
|---|------|------|----------|----------|
| 1 | PRD Traceability | 25 | — | 每个 TC 是否可追溯到具体 PRD AC？ |
| 2 | Step Actionability | 25 | <20 阻塞 | Steps 和 Expected Result 是否足够具体？ |
| 3 | Route & Element Accuracy | 20 | — | Route 和 Element 是否可解析？ |
| 4 | Completeness | 20 | — | 是否覆盖所有接口类型、边界情况和集成场景？ |
| 5 | Structure & ID Integrity | 10 | — | TC 编号、traceability table、分类是否正确？ |

**架构决策**：

| 决策项 | 结论 | 理由 |
|--------|------|------|
| 评估架构 | 标准 adversarial loop（doc-scorer + doc-reviser） | 与 eval-prd/eval-design 保持一致 |
| 形态 | 独立 skill | 不增加 gen-test-cases 的复杂度 |
| 输入 | test-cases.md + PRD 源文件 | 需要交叉验证 traceability 和 completeness |
| Reviser 范围 | 仅修改 test-cases.md | 不越权修改已通过 eval-prd 的 PRD |
| 输出路径 | `docs/features/<slug>/eval/iteration-N.md` | 与其他 eval skill 一致 |
| 默认参数 | target 90, max iterations 6 | 6 轮配合 <20 阻塞阈值 |
| 失败处理 | 标记 blocked，输出报告让用户决策 | 与现有 eval skill 行为一致 |
| T-test-1b 执行位置 | main session（`main_session: true`） | subagent 不能 spawn subagent，doc-scorer/doc-reviser 需由 main session 直接调度 |

**Subagent 约束说明**：

Claude Code 中只有 main session 可以 spawn subagent。`doc-scorer` 和 `doc-reviser` 是 subagent，`task-executor` 也是 subagent。如果 `run-tasks` 将 T-test-1b dispatch 给 `task-executor`，`task-executor` 无法再 spawn `doc-scorer`/`doc-reviser`。

解决方案：T-test-1b 任务模板标记 `main_session: true`，`run-tasks` 的 dispatch 逻辑识别此标记后，不 dispatch 给 `task-executor`，而是在 main session 中直接执行 eval-test-cases skill。

### 任务集成

在 breakdown-tasks 中追加 T-test-1b（标记 `main_session: true`），形成新的依赖链：

```
T-test-1 (gen-test-cases) → T-test-1b (eval-test-cases, main_session) → T-test-2 (gen-test-scripts)
```

T-test-2 的依赖从 T-test-1 改为 T-test-1b。

`run-tasks` 的 dispatch 逻辑需增加判断：遇到 `main_session: true` 的 task 时，跳过 dispatch 给 task-executor，改为在 main session 中直接调用对应 skill。

## Alternatives Considered

| Approach | Pros | Cons | Verdict |
|----------|------|------|---------|
| **eval-test-cases 独立 skill** | 职责清晰，可独立重跑，架构一致 | 多一个 skill 需维护 | **Chosen** |
| 内置到 gen-test-cases | 用户无感，少一个调用步骤 | gen-test-cases 过于复杂，违反单一职责 | Rejected |
| 单次评分无迭代 | 实现简单 | 复杂问题一次修不好，与 adversarial loop 体系不一致 | Rejected |
| eval 同时修改 PRD | 根源修复 | 越权修改已评审的上游文档 | Rejected |
| 不做 eval，依赖 gen-test-scripts 自身验证 | 零额外成本 | 脚本生成后才发现问题，修复成本高 | Rejected |

## Scope

### In Scope

- 新增 `eval-test-cases` skill（SKILL.md + rubric.md + report.md 模板）— 复用现有 eval skill 结构，预计 ~1 session
- 新增 T-test-1b 任务模板（`breakdown-tasks/templates/eval-test-cases.md`，标记 `main_session: true`）— 参照现有模板，预计 ~15 min
- 修改 `run-tasks` dispatch 逻辑（识别 `main_session: true` 标记，在 main session 中直接执行）— 涉及 run-tasks 核心分发路径，需增加一个条件分支判断（~10 行逻辑 + 单元验证），是本 scope 中复杂度最高的改动，预计 ~1 session
- 修改 T-test-2 依赖（从 T-test-1 改为 T-test-1b）— 模板文件单行改动，预计 ~5 min
- 更新 SKILLS.md 注册表 — 单条目追加，预计 ~5 min
- 更新 CLAUDE.md 中的工作流图 — 单处修改，预计 ~10 min

**总估算**：~2 sessions（约 3-4 小时），其中 run-tasks dispatch 改动占约一半工作量

### Out of Scope

- 修改 gen-test-cases 的生成逻辑
- 修改 doc-scorer 或 doc-reviser agent 定义
- 修改 gen-test-scripts 的脚本生成逻辑
- 修改其他 eval skill 的任何内容
- 修改 task-executor agent 定义

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Reviser 修改 test-cases.md 时破坏 TC 编号或 traceability table | Medium | Medium | Reviser prompt 内置 invariant guard 指令：明确禁止修改 TC ID 编号格式和 traceability table 结构；每轮结束后 diff 检查 TC ID 集合是否与原始一致，不一致则拒绝本轮修订并回退 |
| 6 轮迭代仍无法达到 target 90 分 | Low | Low | 报告中列出未达标维度及具体扣分项，附带 reviser 建议修改方向（例如："Step Actionability 18/25，缺少 2 条 TC 的 expected result 描述"），降低用户手动修复的判断成本 |
| Step Actionability <20 阻塞阈值过于严格 | Low | Medium | 阈值可根据实战反馈调整；6 轮迭代给足修复空间 |
| eval 增加每个 feature 的 agent token 消耗 | High | Low | 估算单次 eval 消耗：doc-scorer ~5k tokens/轮 + doc-reviser ~8k tokens/轮 × 平均 3 轮 = ~40k tokens。修复成本同样 ~40k tokens（问题节已量化），但 eval 在问题上游拦截，避免 gen-test-scripts 的无效生成（~15k tokens/脚本 × 3 脚本 = ~45k）及后续修复（~40k）。综合比较：eval ~40k vs 无 eval ~85k（生成 + 修复），前置 eval 净节省 ~45k tokens |
| main-session task 阻塞 run-tasks dispatch 流水线 | Low | Medium | eval 有 max iterations 上限（默认 6），不会无限阻塞；且整个任务链只有一个 main-session task |

## Success Criteria

- [ ] `eval-test-cases` skill 可独立调用，输入 test-cases.md + PRD 源文件，输出 eval 报告到 `docs/features/<slug>/eval/iteration-N.md`
- [ ] Rubric 包含 5 个维度，总分 100，Step Actionability <20 阻塞下游
- [ ] adversarial loop 默认 target 90, max iterations 6，reviser 仅修改 test-cases.md
- [ ] T-test-1b 任务模板存在且依赖 T-test-1，T-test-2 依赖 T-test-1b
- [ ] run-tasks 遇到 `main_session: true` 标记的 task 时，在 main session 中直接调用 eval-test-cases skill（不 dispatch 给 task-executor），且 doc-scorer 和 doc-reviser subagent 正常 spawn 并完成评分与修订。**验收方法**：在一个测试 feature 上运行完整任务链（T-test-1 → T-test-1b → T-test-2），验证：(1) T-test-1b 执行时 task-executor 未被 spawn（检查 agent transcript 中无 task-executor 调用记录）；(2) eval 报告生成到 `docs/features/<slug>/eval/iteration-1.md`；(3) doc-scorer 和 doc-reviser 的 transcript 存在且包含评分/修订内容
- [ ] SKILLS.md 注册表包含 eval-test-cases 条目
- [ ] CLAUDE.md 工作流图反映 T-test-1 → T-test-1b → T-test-2 链路
- [ ] 不修改 doc-scorer.md 或 doc-reviser.md 的任何内容

## Next Steps

- Proceed to `/write-prd`
