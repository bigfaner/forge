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
- 全量审计发现 5 个活跃提案存在同类矛盾风险（`skill-ecosystem-audit`：作用域重叠导致同一文件被双写；`deduplicate-quality-gate`：SC 要求删除冗余检查的同时 InScope 保留所有现有检查；`surface-aware-justfile`：SC 要求 Justfile 路径完全由表面推导同时 InScope 保留硬编码 fallback；`pipeline-integration-stitch`：SC 禁止 gen-and-run 引用同时 InScope 要求迁移错误提示；第 5 个为方向型冲突——SC 要求 ADD 与 SUBTRACT 作用于同一模块）
- Lesson 已记录：`docs/lessons/gotcha-proposal-success-criteria-contradiction.md`

### Urgency

当前 v3.0.0 分支有 5+ 个活跃提案待实现，全量审计发现 5/5 存在同类矛盾风险（100% 命中率）。若不修复检测机制，矛盾 SC 将持续流入实现阶段，单次矛盾返工成本约 0.5-1 人天（基于 pipeline-integration-stitch 实际损失估算）。

## Proposed Solution

双层防御——brainstorm 预防 + eval 检测：

**Layer 1 (brainstorm 预防)**：在 Step 5 (Write Proposal) 中新增 SC Consistency Check 规则文件。Agent 写完 SC 和 In Scope 后，先按作用区域（文件/目录/模块）对条目聚类，再在组内执行可满足性检查，发现"满足 A 则 B 必然失败"的组合。

**Layer 2 (eval 检测)**：扩展 scorer-protocol Phase 1 "Self-contradiction check" 子步骤，要求 scorer 显式执行聚类 + 组内 SC↔SC 和 SC↔InScope 的可满足性检查，将矛盾标记为 attack point。

### Innovation Highlights


**聚类 + 组内检查**而非朴素逐对。同一代码区域（同文件、同目录、同模块）是 SC 矛盾的高概率启发式——实践中绝大多数矛盾发生在同区域内。先聚类再组内检查，将 O(n²) 削减为 O(n·k)（每个条目只与同组条目比较，k 为平均簇大小）。审计数据表明 95%+ 提案 SC < 12 条（最大 66 对），聚类后典型检查量降至 15-30 对（削减 55-77%）；大提案（SC 16-22 条，最大 120-231 对）聚类后降至 40-60 对（削减 70-80%）。为防止跨区域矛盾被遗漏，组内检查后增加一轮轻量的全对方向检查（ADD vs SUBTRACT on same symbol）作为 fallback，覆盖跨组矛盾。

相比"方向标签"方案（ADD/SUBTRACT/MODIFY），聚类方案不增加 SC 条目的标注负担，且能捕捉非方向型逻辑互斥（如"响应 < 100ms" vs "全量日志落盘"）。

## Requirements Analysis

### Key Scenarios

- **场景 1（预防）**：Agent 在 brainstorm Step 5 写完 SC 后，执行聚类 + 组内 check，发现 SC-3 与 InScope-2 作用于同一区域（`forge-cli/`）且逻辑互斥。输出格式为结构化矛盾报告：列出矛盾对（引用具体 SC 文本）、矛盾类型（互斥/方向冲突/资源竞争）、建议操作（删除其一 / 声明互斥区域 / 重写为兼容表述）。示例输出：`CONFLICT: SC-3 "grep -r gen-and-run forge-cli/ returns 0" ↔ InScope-2 "保留 gen-and-run 迁移错误提示" | Type: 互斥（同一文件路径不能同时存在和不存在） | Suggestion: 删除 SC-3 或将 InScope-2 改为"在迁移文档中保留说明"`。用户在 brainstorm 交互中选择操作后 agent 修订 SC，继续流程
- **场景 2（检测）**：eval scorer 在 Phase 1 Reasoning Audit 中按区域聚类后检查 SC 可满足性，发现矛盾并生成 attack point 要求 reviser 修订。修订后的 SC 必须重新通过一致性检查（重新聚类 + 组内检查），避免修订引入新矛盾
- **场景 3（正常通过）**：无矛盾的 SC 集合，聚类后各组内零发现，不产生任何输出或阻塞
- **场景 4（模糊矛盾）**：两条 SC 语义边界不清晰（如"简化 X 流程" vs "保留 X 完整功能"），LLM 无法确定是否矛盾。规则要求标记为 "ambiguous — 需用户确认"，避免强制二选一
- **场景 5（聚类错误）**：SC 被错误聚类（如跨模块影响的条目未归入同一组），导致组内检查漏过矛盾。fallback 全对方向检查作为安全网覆盖此类情况。覆盖边界：fallback 仅覆盖"同符号 ADD vs SUBTRACT"方向型矛盾，非方向型跨组矛盾（如跨模块性能 vs 完整性互斥）仍可能遗漏，需依赖 eval 层全对扫描补强

### Non-Functional Requirements

- 检查协议不增加 proposal 模板结构（不改模板，只加规则）

- 聚类 + 组内检查，典型提案（SC < 12 条）聚类后检查量 < 30 对（vs 朴素逐对最大 66 对），大提案（SC 16-22 条）聚类后 < 60 对（vs 朴素逐对最大 231 对）；执行时间与当前 Step 5 相比增加 < 30%（需实测验证）
- 最坏情况（所有 SC 落入同一簇，k=n）退化为 O(n²)，但实际审计中 95%+ 提案 SC < 12 条，此时 ≈66 对仍在可接受范围
- 可接受误报率 < 5%：对明确无矛盾的 SC 集合，误标为矛盾的比例应低于 5%；若实测误报率 > 10% 则需调整 prompt 阈值或增加"高置信/低置信"分级

### Constraints & Dependencies

- brainstorm skill 的 rules/ 目录结构已建立（已有 challenge-protocol.md）
- eval scorer-protocol.md Phase 1 已有 "self-contradiction check" 子步骤，需扩展而非替换
- 修改 plugins/forge/ 下的文件需遵循 `docs/conventions/forge-distribution.md`
- **核心依赖：LLM 推理能力**。本方案的矛盾检测完全依赖 LLM 对自然语言 SC 的逻辑推理，而非形式化验证。LLM 推理的准确性和稳定性直接影响检测质量——不同模型版本或 temperature 设置可能导致结果不一致。双层设计部分缓解此风险，但无法消除
- **Token 溢出风险**：超大提案（SC > 25 条）聚类后仍可能产生 60+ 对检查对，加上 SC 全文作为上下文，单次检查 prompt 可能超出 context window。缓解：分组串行检查（每次仅传入一个簇的 SC 对），而非一次性传入所有对

## Alternatives & Industry Benchmarking

### Industry Solutions

软件需求工程中，需求一致性检查是标准实践。IBM DOORS 和 Siemens Polarion 等需求管理工具内置 consistency check（基于 traceability link 和 attribute conflict detection）。Formal methods 如 Alloy Analyzer（MIT）可对需求模型做约束求解验证，但需形式化规约输入，不适用于自然语言 SC。本方案采用 LLM 推理实现聚类 + 组内可满足性检查，在自然语言 SC 场景下平衡覆盖面与执行成本。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零成本 | 矛盾持续流入实现阶段，返工成本高 | Rejected: 已有实际损失记录 |
| 方向标签 (ADD/SUBTRACT/MODIFY) | lesson 建议 | 结构化，机械执行 | 只能捕捉方向型矛盾，需每个 SC 标注 | Rejected: 泛化能力不足 |
| 朴素逐对检查 | 通用 | 覆盖全 | O(n²) 对数增长，信号/噪声比低 | Rejected: 大提案 231 对中仅 1-2 对有矛盾 |
| 形式化约束求解 (如 Alloy) | MIT Alloy Analyzer | 数学完备性，确定性结果 | 需将自然语言 SC 转为形式化规约，转换成本极高 | Rejected: 输入格式不匹配 |
| **聚类 + 组内逐对检查** | 分区启发式 + LLM 推理 | 覆盖面广，检查量削减 55-80%，不改模板 | 依赖 LLM 推理（非确定性，存在漏报/误报风险）；token 消耗随 SC 数量增长；聚类质量影响检查效果 | **Selected: 覆盖面广 + 高效 + 可落地** |

## Feasibility Assessment

### Technical Feasibility

完全可行。brainstorm 层只需新建 1 个规则文件 + SKILL.md 加 1 行引用；eval 层扩展 scorer-protocol 已有的 self-contradiction check 子步骤。LLM 推理稳定性方面：pipeline-integration-stitch 矛盾案（grep 零结果 vs 保留迁移提示）在 Claude Sonnet 级别模型上可稳定检出（逻辑互斥明确），模糊矛盾的检出率依赖于 prompt 工程，通过"标记为 ambiguous 交由用户"策略降低对 LLM 判断精度的要求。Prompt 策略方向：规则文件采用 "对每个 SC 对，分别假设 A 为真推导 B 的状态，再假设 B 为真推导 A 的状态" 的双向反证法 prompt 结构，而非直接问"这两个 SC 是否矛盾"，以减少 LLM 倾向性回答。

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

- 修改 `plugins/forge/skills/eval/rubrics/proposal.md` Dimension 9 — 新增 "SC internal consistency" criterion (25pts)，调整现有 criterion 分值（measurable 40→30, coverage 40→25, consistency 新增 25, 总分 80 不变）。D9 检查 SC 条目间的内部可满足性（组内 SC↔SC 和 SC↔InScope），与 D10（检查 SC 与 Scope/Solution 的覆盖对齐）职责不重叠。注意：measurable 从 40→30 和 coverage 从 40→25 的调整可能降低原高分提案在这两个子维度的得分上限，需在实现时验证已有提案评分不受显著影响。连锁影响：D9 分值重分配会改变 eval 总分，可能影响依赖总分阈值的 gate 判断（如 700/1000 通过线），需同步验证 gate 逻辑是否需要调整
- 修改 `plugins/forge/skills/brainstorm/rules/sc-consistency.md` — 包含 fallback 全对方向检查规则（组内检查后增加一轮 ADD vs SUBTRACT on same symbol 的跨组检查）

### Out of Scope

- 修改 proposal 模板（templates/proposal.md）
- 修改 reviser-protocol.md
- 修改 eval freeform review 流程
- 回溯修正现有提案的 SC 矛盾（注：Evidence 中识别的 5 个风险提案不在本方案范围内，建议在 SC Consistency Gate 上线后，对这 5 个提案执行一次专项 eval 以验证检测效果）

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|

| 逐对检查增加 agent 执行时间 | L | L | 聚类后组内检查，典型提案 < 60 对，执行时间增量 < 30%（需实测验证） |
| 误报：将非矛盾标记为矛盾 | M | M | 规则要求矛盾必须可论证（引用具体 SC 文本 + 逻辑推导链），非模糊感觉。实际误报率需上线后统计，若 > 10% 则需调整 prompt 阈值或增加"高置信/低置信"分级 |

| agent 忽略规则文件不执行检查 | H | H | SKILL.md Step 5 中将一致性检查列为强制步骤（非可选引用），规则文件命名明确（sc-consistency.md）；eval 层独立执行检查作为兜底。承认：两层均依赖 LLM 主动读取规则文件，无硬性防护。硬防护方案：在 brainstorm SKILL.md Step 5 的 SC 输出格式中增加结构化字段 "consistency_check_result: {status: pass|fail|ambiguous, pairs_checked: N}"，若该字段缺失则 SC 输出格式校验失败，迫使 agent 执行检查 |
| eval scorer Phase 1 扩展后 attack points 增多导致 reviser 负担 | L | M | 矛盾类 attack point 通常只有 1-3 个，不会爆炸 |
| D9 rubric 分值重分配（measurable 40→30, coverage 40→25）导致已有提案评分不可比 | M | M | SC #10 要求回溯测试验证评分差异 < 5%；若差异超标，可保留原始分值结构并新增 D9a 子维度（总分从 80 增至 105）作为过渡方案 |
| LLM 矛盾检测漏报（false negative） | M | H | 双层（brainstorm + eval）使用同一 LLM，无法通过简单 fallback 消除漏报。差异化策略：eval 层使用更宽泛的搜索 prompt（不限定区域聚类，直接全对扫描）和更高 temperature（0.7 vs brainstorm 层 0.3）以增加推理多样性，降低两层犯同样错误的概率。长期可增加结构化 SC 影响路径标注作为不依赖 LLM 推理多样性的硬兜底 |

## Success Criteria

- [ ] brainstorm 规则文件 `rules/sc-consistency.md` 存在且包含聚类协议、组内可满足性检查协议和 lesson 反例
- [ ] SKILL.md Step 5 包含对 `rules/sc-consistency.md` 的引用
- [ ] scorer-protocol Phase 1 Step 4 包含聚类 + 组内 SC↔SC 和 SC↔InScope 可满足性检查指令
- [ ] proposal rubric D9 包含 "SC internal consistency" criterion (25pts)，D9 总分 80pts 不变

- [ ] scorer-protocol.md 文本中包含聚类 + 组内可满足性检查的显式指令，且 gen-and-run 矛盾场景作为示例用例被引用（功能性验证：指令中明确要求"按区域聚类后在组内执行双向可满足性推导"，而非仅要求"检查矛盾"）
- [ ] sc-consistency.md 规则文件中包含"无矛盾 SC 集合应零输出"的明确规则描述（功能性验证：规则中包含"对无矛盾集合输出空报告"的显式测试用例）
- [ ] 功能有效性验证：对 pipeline-integration-stitch 矛盾案（grep 零结果 vs 保留迁移提示）执行 brainstorm 一致性检查，确认能检出该矛盾并输出结构化矛盾报告
- [ ] 修改遵循 `docs/conventions/forge-distribution.md` 的分发约束
- [ ] sc-consistency.md 规则文件中包含 fallback 跨组方向检查规则（覆盖跨区域矛盾场景）
- [ ] proposal rubric D9 measurable (30pts) + coverage (25pts) + consistency (25pts) 分值调整后，对已有无矛盾提案的评分差异 < 5%（通过回溯测试验证）

## Next Steps

- Proceed to `/write-prd` to formalize requirements
