---
created: 2026-05-08
author: "fanhuifeng"
status: Draft
---

# Proposal: Add eval-test-cases Skill

## Problem

gen-test-cases 生成的 `test-cases.md` 直接交给 gen-test-scripts 生成可执行脚本，中间没有质量门控。如果 test-cases.md 的步骤描述模糊、traceability 缺失、或路由信息错误，gen-test-scripts 会生成无法运行的脚本，浪费 3-5 轮 agent 交互来修复（每轮 ~10k tokens）。

### Evidence

- eval-prd 对 PRD、eval-design 对 tech-design 都有 adversarial loop 质量门控，唯独 gen-test-cases 的输出没有
- gen-test-cases 的 Step Actionability 完全依赖 LLM 单次生成的质量，没有验证环节
- 下游 gen-test-scripts 依赖 TC 的 Steps、Route、Element 字段来生成脚本，任何模糊或不准确都会导致脚本失败

### Urgency

每个 feature 都会走 gen-test-cases → gen-test-scripts 链路，缺少质量门控意味着每个 feature 都承担脚本生成失败的风险。

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

### 任务集成

在 breakdown-tasks 中追加 T-test-1b，形成新的依赖链：

```
T-test-1 (gen-test-cases) → T-test-1b (eval-test-cases) → T-test-2 (gen-test-scripts)
```

T-test-2 的依赖从 T-test-1 改为 T-test-1b。

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

- 新增 `eval-test-cases` skill（SKILL.md + rubric.md + report.md 模板）
- 新增 T-test-1b 任务模板（`breakdown-tasks/templates/eval-test-cases.md`）
- 修改 T-test-2 依赖（从 T-test-1 改为 T-test-1b）
- 更新 SKILLS.md 注册表
- 更新 CLAUDE.md 中的工作流图

### Out of Scope

- 修改 gen-test-cases 的生成逻辑
- 修改 doc-scorer 或 doc-reviser agent 定义
- 修改 gen-test-scripts 的脚本生成逻辑
- 修改其他 eval skill 的任何内容

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Reviser 修改 test-cases.md 时破坏 TC 编号或 traceability table | Medium | Medium | Rubric 中 Structure & ID Integrity 维度（10分）会捕获此类问题，scorer 在下一轮扣分 |
| 6 轮迭代仍无法达到 target 90 分 | Low | Low | 任务标记为 blocked，用户可查看报告后决定是否接受或手动修复 |
| Step Actionability <20 阻塞阈值过于严格 | Low | Medium | 阈值可根据实战反馈调整；6 轮迭代给足修复空间 |
| eval 增加每个 feature 的 agent token 消耗 | High | Low | 相比 gen-test-scripts 失败后重试的消耗，前置 eval 的 ROI 为正 |

## Success Criteria

- [ ] `eval-test-cases` skill 可独立调用，输入 test-cases.md + PRD 源文件，输出 eval 报告
- [ ] Rubric 包含 5 个维度，总分 100，Step Actionability <20 阻塞下游
- [ ] adversarial loop 默认 target 90, max iterations 6，reviser 仅修改 test-cases.md
- [ ] T-test-1b 任务模板存在且依赖 T-test-1，T-test-2 依赖 T-test-1b
- [ ] SKILLS.md 注册表包含 eval-test-cases 条目
- [ ] CLAUDE.md 工作流图反映 T-test-1 → T-test-1b → T-test-2 链路
- [ ] 不修改 doc-scorer.md 或 doc-reviser.md 的任何内容

## Next Steps

- Proceed to `/write-prd`
