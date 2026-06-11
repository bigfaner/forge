---
created: "2026-05-25"
author: fanhuifeng
status: Draft
---

# Proposal: SC Consistency Gate — Proposal Success Criteria 矛盾检测机制

## Problem

Proposal 的 Success Criteria (SC) 条目之间存在逻辑矛盾（满足 A 则 B 必然失败），但现有的 brainstorm 和 eval 流程均无法检测此类矛盾，导致矛盾提案通过 897/1000 的 adversarial eval 后在实现阶段才暴露。

### Evidence

- Feature `pipeline-integration-stitch` 的 proposal 通过 897/1000 adversarial eval，但 SC 要求 `grep -r "gen-and-run" forge-cli/` 返回零结果，同时 In Scope 要求保留迁移错误提示——两者互斥
- 全量审计发现 5 个活跃提案存在同类矛盾风险（`skill-ecosystem-audit`、`deduplicate-quality-gate`、`surface-aware-justfile` 等）
- Lesson 已记录：`docs/lessons/gotcha-proposal-success-criteria-contradiction.md`

### Urgency

当前 v3.0.0 分支有多个活跃提案待实现，若不修复检测机制，矛盾 SC 将持续流入实现阶段，导致返工。

## Proposed Solution

双层防御——brainstorm 预防 + eval 检测：

**Layer 1 (brainstorm 预防)**：在 Step 5 (Write Proposal) 中新增 SC Consistency Check 规则文件。Agent 写完 SC 和 In Scope 后，先按作用区域（文件/目录/模块）对条目聚类，再在组内执行可满足性检查，发现"满足 A 则 B 必然失败"的组合。

**Layer 2 (eval 检测)**：扩展 scorer-protocol Phase 1 "Self-contradiction check" 子步骤，要求 scorer 显式执行聚类 + 组内 SC↔SC 和 SC↔InScope 的可满足性检查，将矛盾标记为 attack point。

### Innovation Highlights

**聚类 + 组内检查**而非朴素逐对。SC 矛盾的必要条件是两个条目作用于同一代码区域（同文件、同目录、同模块）。先聚类再检查将 O(n²) 削减为 O(n)（每个条目只与同组条目比较），典型提案检查量从 125-246 对降至 40-60 对，削减 70-80%，且不遗漏任何真正矛盾。

相比"方向标签"方案（ADD/SUBTRACT/MODIFY），聚类方案不增加 SC 条目的标注负担，且能捕捉非方向型逻辑互斥（如"响应 < 100ms" vs "全量日志落盘"）。

## Requirements Analysis

### Key Scenarios

- **场景 1（预防）**：Agent 在 brainstorm Step 5 写完 SC 后，执行聚类 + 组内 check，发现 SC-3 与 InScope-2 作用于同一区域（`forge-cli/`）且逻辑互斥，提示用户选择其一或声明互斥区域
- **场景 2（检测）**：eval scorer 在 Phase 1 Reasoning Audit 中按区域聚类后检查 SC 可满足性，发现矛盾并生成 attack point 要求 reviser 修订
- **场景 3（正常通过）**：无矛盾的 SC 集合，聚类后各组内零发现，不产生任何输出或阻塞

### Non-Functional Requirements

- 检查协议不增加 proposal 模板结构（不改模板，只加规则）
- 聚类 + 组内检查，典型提案检查量 < 60 对（vs 朴素逐对 125-246 对），agent 执行时间 < 20 秒
- 不产生误报：明确无矛盾的 SC 集合不应被标记

### Constraints & Dependencies

- brainstorm skill 的 rules/ 目录结构已建立（已有 challenge-protocol.md）
- eval scorer-protocol.md Phase 1 已有 "self-contradiction check" 子步骤，需扩展而非替换
- 修改 plugins/forge/ 下的文件需遵循 `docs/conventions/forge-distribution.md`

## Alternatives & Industry Benchmarking

### Industry Solutions

软件需求工程中，需求一致性检查是标准实践（如 DOORS、Polarion 等需求管理工具内置 consistency check）。NLP 领域有基于约束求解的需求冲突检测研究。本方案借鉴了 SAT solver 的思路——将 SC 条目视为约束，检查可满足性。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零成本 | 矛盾持续流入实现阶段，返工成本高 | Rejected: 已有实际损失记录 |
| 方向标签 (ADD/SUBTRACT/MODIFY) | lesson 建议 | 结构化，机械执行 | 只能捕捉方向型矛盾，需每个 SC 标注 | Rejected: 泛化能力不足 |
| 朴素逐对检查 | 约束求解 | 覆盖全 | O(n²) 对数增长，信号/噪声比低 | Rejected: 大提案 246 对中仅 1-2 对有矛盾 |
| **聚类 + 组内逐对检查** | 约束求解 + 分区优化 | 覆盖全，检查量削减 70-80%，不改模板 | 无 | **Selected: 覆盖面广 + 高效** |

## Feasibility Assessment

### Technical Feasibility

完全可行。brainstorm 层只需新建 1 个规则文件 + SKILL.md 加 1 行引用；eval 层扩展 scorer-protocol 已有的 self-contradiction check 子步骤。

### Resource & Timeline

单次对话内可完成：规则文件起草 + scorer-protocol 修改。

### Dependency Readiness

无外部依赖。所有修改在 plugins/forge/ 内完成。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| "eval 的 Logical Consistency 维度能捕捉 SC 矛盾" | 事实核查 | Overturned: pipeline-integration-stitch 得分 897/1000，该维度未检出矛盾 |
| "方向标签足以捕捉所有 SC 矛盾" | Assumption Flip | Overturned: 非方向型矛盾（性能 vs 完整性）无法通过方向标签捕捉 |
| "SC 条目通常 < 15 条" | 事实核查 | Confirmed: 全量审计中 95%+ 的提案 SC 条目 < 12 条 |

## Scope

### In Scope

- 新建 `plugins/forge/skills/brainstorm/rules/sc-consistency.md` — SC 聚类 + 组内可满足性检查规则
- 修改 `plugins/forge/skills/brainstorm/SKILL.md` Step 5 — 添加规则引用
- 修改 `plugins/forge/skills/eval/experts/protocol/scorer-protocol.md` Phase 1 Step 4 — 扩展 self-contradiction check
- 修改 `plugins/forge/skills/eval/rubrics/proposal.md` Dimension 9 — 新增 "SC internal consistency" criterion (25pts)，调整现有 criterion 分值

### Out of Scope

- 修改 proposal 模板（templates/proposal.md）
- 修改 reviser-protocol.md
- 修改 eval freeform review 流程
- 回溯修正现有提案的 SC 矛盾

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 逐对检查增加 agent 执行时间 | L | L | 聚类后组内检查，典型提案 < 60 对，agent 执行 < 20 秒 |
| 误报：将非矛盾标记为矛盾 | M | M | 规则要求矛盾必须可论证（引用具体 SC 文本），非模糊感觉 |
| agent 忽略规则文件不执行检查 | M | H | SKILL.md 引用 + 规则文件命名清晰 + eval 层兜底 |
| eval scorer Phase 1 扩展后 attack points 增多导致 reviser 负担 | L | M | 矛盾类 attack point 通常只有 1-3 个，不会爆炸 |

## Success Criteria

- [ ] brainstorm 规则文件 `rules/sc-consistency.md` 存在且包含聚类协议、组内可满足性检查协议和 lesson 反例
- [ ] SKILL.md Step 5 包含对 `rules/sc-consistency.md` 的引用
- [ ] scorer-protocol Phase 1 Step 4 包含聚类 + 组内 SC↔SC 和 SC↔InScope 可满足性检查指令
- [ ] proposal rubric D9 包含 "SC internal consistency" criterion (25pts)，D9 总分 80pts 不变
- [ ] 对 lesson 中的 gen-and-run 矛盾场景，scorer 能在 D9 "SC internal consistency" 维度扣分并生成 attack point
- [ ] 对无矛盾的 SC 集合，规则不产生任何阻塞或警告，D9 "SC internal consistency" 满分
- [ ] 修改遵循 `docs/conventions/forge-distribution.md` 的分发约束

## Next Steps

- Proceed to `/write-prd` to formalize requirements
