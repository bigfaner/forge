---
feature: "test-capability-v2"
eval-type: prd
scorer: pm
mode: B (no UI)
iteration: final
score: 945
---

# PRD Final Evaluation Report — PM Perspective

**Date**: 2026-05-23
**Scorer**: Senior PM (adversarial)
**Mode**: B — Feature WITHOUT UI (no prd-ui-functions.md)
**Documents Evaluated**: `prd-spec.md`, `prd-user-stories.md`

## Score Summary

| # | Dimension | Score | Max | Verdict |
|---|-----------|-------|-----|---------|
| 1 | Background & Goals | 100 | 100 | PASS |
| 2 | Flow Diagrams | 150 | 150 | PASS |
| 3 | Flow Completeness | 195 | 200 | PASS |
| 4 | User Stories | 190 | 200 | PASS |
| 5 | Scenario Completeness | 140 | 150 | PASS |
| 6 | Edge Case Coverage | 90 | 100 | PASS |
| 7 | Scope Clarity | 80 | 100 | WARN |
| **Total** | | **945** | **1000** | **PASS** |

---

## Dimension 1: Background & Goals — 100/100

### Background 三要素 (30/30)
- **Reason**: 清晰列出三大结构性缺陷（双路径并存、测试深度不足、通用性有限），每个缺陷都有具体的现状描述，不是泛泛而谈。
- **Target**: 五个升级方向（管线统一、深度增强、通用扩展、评测补全、信息增强）目标明确。
- **Users**: 两类用户定义清晰——Forge 用户（项目开发者）和 Forge 维护者，各有使用场景描述。

### Goals 量化 (30/30)
全部 6 个目标均有可量化指标：
- "高风险 Journey 平均测试数 ≥ 13"、"高风险 ≥ 低风险 × 1.5"
- "Fact Table 覆盖率提升 ≥ 20 个百分点"
- "内置 ≥ 3 个新 Convention 文件"
- "eval-journey/eval-contract 评分 ≥ 850/1000"
- "新 Mobile 项目从零到可运行 Maestro 测试 ≤ 30 分钟"
- "deep link 测试覆盖 ≥ 2 个核心 Journey"

### 逻辑一致性 (40/40)
Reason → Target → Goals 三者形成完整因果链：三大问题直接映射到五个升级方向和六个量化目标。

---

## Dimension 2: Flow Diagrams — 150/150

### Mermaid 图存在 (50/50)
包含完整的 `flowchart TD` Mermaid 图，约 40 个节点。

### 主路径完整 (50/50)
从 `START`（用户运行测试生成）到 `END`（完成/完成含失败详情），覆盖全部四个阶段。

### 决策点与错误分支 (50/50)
菱形决策节点包括：PRD_CHECK、SCENE_DETECT、DETECT（Convention）、EVAL_J、EVAL_C、R2L_CHOICE、ENV_CHECK、FIX_DECIDE。错误/异常分支涵盖：PRD_FAIL、SCENE_FAIL、EVAL_J_SKIP、EVAL_C_SKIP、PAUSE_J（3 种恢复）、PAUSE_C（3 种恢复）、R2L_DEGRADE、ENV_FAIL。

---

## Dimension 3: Flow Completeness — 195/200

### 流程步骤描述完整业务过程 (65/70)
四阶段 14 步覆盖完整。包含前置条件、场景检测信号映射表（17 种组合）、eval 迭代逻辑、PAUSE 恢复路径、自动修复回退策略。**扣分原因**：Run-to-Learn 的触发条件（`run_to_learn: true` 或 CLI flag）仅在 Mermaid 图节点注释中出现，Flow Description 正文仅写"可选"，正文与图信息不对称。正文应明确写出触发条件。

### 数据流文档 (70/70)
Data Flow Table 包含 9 行记录，覆盖所有关键数据传递路径，包含 session.yaml、文件系统目录、Fact Table JSON 等多种传递方式。eval_result 的结构化格式定义完整。

### 异常处理与边界情况 (60/60)
六类异常处理机制：PRD 前置检查、eval-skipped 降级（4 步策略）、PAUSE 三种恢复、Schema 验证重试、gen-test-scripts 验证重试、Run-to-Learn 失败处理（4 种场景）。Pipeline Exit Codes 表定义了 6 种终止点的退出码。

---

## Dimension 4: User Stories — 190/200

### 覆盖 (50/50)
Forge 用户（项目开发者）→ Stories 1, 2, 3, 5, 6；Forge 维护者 → Stories 4, 7。两类用户均有覆盖。

### 格式正确 (50/50)
全部 7 个故事遵循 As a / I want to / So that 格式，动词具体可操作。

### AC 格式 (50/50)
所有故事都有 Given/When/Then 格式的 AC。Story 3 覆盖 CLI 和 Mobile 两个场景。Story 4 覆盖 Journey 和 Contract 两个评测。Story 5 覆盖 3 种情况。

### AC 可验证性与边界覆盖 (40/50)
多数 AC 可量化验证（测试数、覆盖率、占比、评分）。

**扣分原因**：
- Story 5 "用户修改行数占草稿总行数比例 ≤ 20%" 标记为 `human-verified`——有数字但自动化回归验证不可行。
- Story 6 "边界/异常 Outcome 占比 ≥ 30%" 同样标记为 `human-verified`——依赖人工判断 Outcome 是否属于"边界/异常"范畴。
- 这两处 AC 虽有量化指标，但验证手段依赖人工审核，无法作为自动化验收门禁。

---

## Dimension 5: Scenario Completeness — 140/150

### 端到端场景覆盖 (58/60)
覆盖 10 种端到端场景：标准管线、PRD 缺失、场景未知/混合、Convention 不存在、eval 未达阈值、eval 迭代耗尽、LLM 解析失败、Run-to-Learn 失败、环境不就绪、测试执行失败。每种场景都有从触发到最终状态的完整描述。

### 隐含假设暴露 (32/40)
前置假设大部分已明确：PRD 质量前置检查、场景检测依赖项目文件信号、Fact Table 覆盖率策略、LLM 解析依赖、Convention 草稿审核。

**扣分原因**：
- 隐含假设："gen-journeys 从 PRD 用户故事提取 Journey"假设 PRD User Story 粒度足够细、质量足够高。前置检查只验证了 User Story 的结构存在性（有 As a/I want/So that + AC），但没有验证 User Story 的内容质量（粒度是否合适、是否覆盖完整业务流程）。如果 PRD 中一个 User Story 覆盖整个功能，Journey 提取可能产出粒度不当的结果。这是一个未暴露的环境依赖假设。

### 业务规则一致性 (50/50)
- BIZ-quality-gate-001：PRD 明确定义了 eval 门禁与 BIZ-quality-gate-001 的集成关系（串行执行、独立判定）。一致。
- BIZ-error-reporting-001：Pipeline Exit Codes 表完全遵循 0/1/2 语义。一致。
- BIZ-task-lifecycle-003：删除 test.graduate 和 test.gen-cases 类型明确列出。新增类型是 PRD 的新增内容，无冲突。

---

## Dimension 6: Edge Case Coverage — 90/100

### 错误路径文档化 (40/40)
10 种错误路径均有具体的检测方式和处理策略，包含退出码、重试次数、降级方案、用户操作指引。

### 边界条件覆盖 (25/35)
覆盖的边界：空/缺失输入（PRD、Convention、Fact Table）、迭代上限（eval 3 轮、重试 2 次）、场景检测边缘情况（未知/混合）、风险等级分类边界、Outcome 数量范围、置信度阈值。

**扣分原因**：
- 缺少"空 Journey"边界处理——PRD 有 User Story 但无法提取出有意义的 Journey 时的行为未定义。
- 缺少"Convention 文件存在但内容无效"（如缺少必需 section）的处理。
- 缺少 Fact Table 初始覆盖率 = 0%（全新项目无代码侦察结果）的特别处理。

### 失败恢复描述 (25/25)
每种失败都有明确恢复路径：PAUSE 三种选项、eval-skipped 手动清除、骨架测试降级、环境修复后重检、修复耗尽输出报告。

---

## Dimension 7: Scope Clarity — 80/100

### 范围内项目是具体可交付物 (35/35)
In Scope 列表中的每个条目都是可识别的功能模块或交付物，有具体的删除清单、文件列表和功能描述。

### 范围外明确列出推迟项 (30/30)
Out of Scope 列出 11 项，每项命名明确，部分附带理由说明（如"仅做就绪检测"说明为什么排除自动准备）。

### 范围与功能规格和用户故事一致 (15/35)
大部分 In Scope 条目与 User Stories 一一对应。

**扣分原因（核心问题）**：
- **跨节不一致**：Story 7 描述了"通过添加场景类型配置文件接入新场景类型"的能力，但 In Scope 列表中没有对应条目。Story 7 的 AC 要求"配置文件放入 scenarios/ 目录"即可扩展——这是一个具体的可交付物（场景类型配置 schema + 目录结构 + eval rubric 动态适配），但 Scope 中遗漏了。
- **边界模糊**：Out of Scope 说"合约 6 维度模型（schema）修改"不在范围内，但 Story 7 要求 eval rubric 的"场景适配"维度根据新增场景类型的 `required_outcomes` 动态调整评分逻辑（"缺少任何一个必须 Outcome 扣 30 分/个"）。虽然文档注解说 `required_outcomes` 是"实例数据不属于 schema 变更"，但 eval 评分逻辑的动态适配能力暗示了 eval 技能的代码修改，Out of Scope 的边界表述不够精确。

---

## Attacks (Blindspot Hunt)

1. **[Scope Clarity] Story 7 缺少 In Scope 条目** — Story 7 描述了完整的"可扩展场景类型系统"（配置文件 schema + scenarios/ 目录 + eval rubric 动态适配），但 In Scope 列表中没有对应条目。引用：Story 7 AC "将配置文件放入 scenarios/ 目录，管线自动识别新场景类型"。需要在 In Scope 中新增"可扩展场景类型系统：场景类型配置 schema + scenarios/ 目录结构 + eval rubric 动态适配"。

2. **[Scope Clarity] Out of Scope 边界与 Story 7 矛盾** — Out of Scope 声明"合约 6 维度模型（schema）修改"不在范围内，但 Story 7 要求 eval rubric 根据新增场景的 `required_outcomes` 动态调整评分逻辑。引用：Out of Scope "合约 6 维度模型（schema）修改（注：required_outcomes 是按场景类型配置的实例数据，不属于 schema 变更）"vs Story 7 AC "eval 评分时按该场景类型的 required_outcomes 列表检查 Outcome 覆盖率，缺少任何一个必须 Outcome 扣 30 分/个"。需要明确 eval 评分逻辑的动态适配是否属于 In Scope，以及是否涉及 rubric schema 变更。

3. **[User Stories] human-verified AC 降低了可自动化验收性** — Story 5 和 Story 6 各有一个 AC 标记为 `human-verified`。引用：Story 5 "用户修改行数占草稿总行数的比例 ≤ 20%（此指标为人工审核参考，标记为 human-verified）"；Story 6 "重新生成的测试中边界/异常 Outcome 占比 ≥ 30%（标记为 human-verified）"。这些 AC 虽有量化数字，但无法作为自动化验收门禁。建议为这些 AC 增加可自动化的代理指标（如：diff 行数可通过 git diff --stat 自动计算）。

4. **[Scenario Completeness] User Story 粒度对 Journey 提取的隐含依赖** — gen-journeys "从 PRD 用户故事提取 Journey 叙事"的步骤假设 User Story 粒度足够细。引用：Flow Description 步骤 4 "gen-journeys 从 PRD 用户故事提取 Journey 叙事（含风险分级）"。PRD 质量前置检查只验证了结构（有 User Story + AC），未验证内容粒度。建议增加 Journey 提取质量的后置检查（如：提取的 Journey 数量与 User Story 数量的比例阈值）。

5. **[Edge Case Coverage] 缺少 Convention 文件内容无效的处理** — 文档覆盖了"Convention 不存在"的处理（test-guide 自动生成），但未覆盖"Convention 存在但缺少必需 section"的情况。引用：Flow Description 步骤 3 "系统检查 Convention 文件是否存在；若不存在，test-guide 自动检测框架并生成 Convention 草稿"。需要增加 Convention 内容有效性验证逻辑。
