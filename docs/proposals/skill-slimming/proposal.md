---
created: 2026-05-20
author: "faner"
status: Draft
---

# Proposal: 逐 Skill 瘦身——拆分、精简、消歧

## Problem

Forge plugin 的 21 个 SKILL.md 文件总计 ~6700 行，其中多个大文件（consolidate-specs 607 行、tech-design 472 行、write-prd 407 行）混合了流程指令、业务规则和内联模板，导致 LLM 上下文浪费且维护困难。同时部分 skill 存在指令歧义（如 noTest vs doc* 概念混淆），增加 agent 执行偏差风险。

### Evidence

- Top 3 大文件平均 495 行，远超其余 skill 平均 245 行
- `consolidate-specs/SKILL.md` 607 行全部自包含，无任何拆分辅助文件
- 多个 skill 内嵌大量模板文本和解释性段落，可直接拆出
- `guide.md` 和多个 SKILL.md 对 `noTest`/`doc*` 的描述产生歧义

### Urgency

v3.0.0 重构窗口期。已有 5 个瘦身相关提案均未执行——方向分散、范围过大是主因。需要一个可立即落地的增量方案。

## Proposed Solution

按大小分层、逐组处理：大文件（400+ 行）独立拆分，中/小文件按领域分组合并处理。每个任务聚焦一组 skill，依次完成拆分结构、精简行数、消除歧义三项目标。

### Innovation Highlights

三层瘦身法：对每个 skill 按需施以拆分（大文件）、精简（冗余文本）、消歧（模糊指令）三种操作，而非一刀切。安全增量：每个任务独立 commit，可逐个验证回滚。

## Requirements Analysis

### Key Scenarios

- Agent 加载 skill 后获得精简、无歧义的指令
- 开发者维护 skill 时通过 SKILL.md 快速理解流程，通过 rules/templates 了解细节
- 大 skill（如 consolidate-specs）拆分后 SKILL.md 降至 300 行以内，关键指令不丢失

### Non-Functional Requirements

- 每个 SKILL.md 行数不超过 350 行（拆分后）
- 拆分产生的辅助文件放在 skill 目录内的 rules/ 或 templates/ 子目录
- 不改变 skill 的输入/输出契约

### Constraints & Dependencies

- 遵守 `docs/conventions/forge-distribution.md` 分发模型
- 遵守 `docs/conventions/skill-self-containment.md` 自洽原则——SKILL.md 必须包含完整流程步骤，辅助文件仅存放规则和模板细节
- 不涉及 Go 源码修改
- 不合并同类 skill（那是 skill-rationalization 的范畴）

## Alternatives & Industry Benchmarking

### Industry Solutions

大型 prompt 工程项目中，指令拆分是常见实践。OpenAI 的 GPT best practices 建议将 system prompt 控制在关键指令内，详细规则外置。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零成本 | 5 个提案已证明现状不可持续 | Rejected: 债务积累 |
| 模式审计 + 批量清理 | skill-slim-down 提案 | 全局视角 | 审计阶段重，已搁置 | Rejected: 太理论化 |
| **按大小分层逐组处理** | 增量重构最佳实践 | 安全可控，立即可执行 | 小组内 skill 可能需不同策略 | **Selected: 平衡效率与安全** |

## Feasibility Assessment

### Technical Feasibility

纯文本修改 + 文件拆分。git 提供完整回滚能力。已有多数 skill 使用 rules/templates 子目录的先例。

### Resource & Timeline

21 个 skill 分 9 组，预计 9 个任务。每个任务包含：分析 → 拆分/精简/消歧 → 验证引用完整性 → commit。

### Dependency Readiness

无外部依赖。所有文件已在本地。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| "逐个 skill 清理效率低" | XY Detection | 用户核心需求是安全增量，不是效率最优；逐组处理可在安全与效率间平衡 |
| "SKILL.md 必须 100% 自包含" | Assumption Flip | 自洽不等于单文件；SKILL.md 包含完整流程 + 辅助文件包含细节规则，仍是自洽的。现有 eval skill 已用此模式 |
| "需要先审计再行动" | Occam's Razor | 5 个提案已做过充分分析；直接动手 + 逐个验证更简单有效 |

## Scope

### In Scope

- 21 个 `skills/*/SKILL.md` 文件的拆分、精简、消歧
- 在各 skill 目录内新建 rules/ 或 templates/ 子目录（按需）
- 清理过时标签、路径引用和歧义描述

### Out of Scope

- Go 源码修改
- skill 输入/输出契约变更
- 合并同类 skill（skill-rationalization 范畴）
- commands/ 和 agents/ 目录的文件
- hooks/、references/、scripts/ 目录

### Task Grouping (9 tasks)

**Tier 1: 大文件独立任务（3 tasks）**
1. consolidate-specs (607 行) → 拆分为 SKILL.md + rules/
2. tech-design (472 行) → 拆分为 SKILL.md + rules/
3. write-prd (407 行) → 拆分为 SKILL.md + rules/

**Tier 2: 中文件按领域分组（3 tasks）**
4. eval (372) + gen-contracts (365) + test-guide (380) → 评测/质量域
5. gen-sitemap (395) + gen-journeys (211) + gen-test-cases (136) + gen-test-scripts (350) → 生成域
6. init-justfile (387) + ui-design (314) + extract-design-md (242) → 基础设施/设计域

**Tier 3: 小文件按领域分组（3 tasks）**
7. breakdown-tasks (144) + quick-tasks (208) + submit-task (156) → 任务管线域
8. brainstorm (139) + learn (259) + forensic (198) + improve-harness (163) → 元分析域
9. clean-code (190) + run-e2e-tests (299) → 工具域

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 拆分后 SKILL.md 丢失关键指令 | M | H | 拆分时 SKILL.md 保留流程骨架 + 关键约束，辅助文件仅放规则细节 |
| 辅助文件命名不统一 | M | L | 约定 rules/ 放规则、templates/ 放模板，不新建其他子目录类型 |
| 消歧时引入新歧义 | L | M | 每处消歧需在 commit message 中注明原文和修改理由 |

## Success Criteria

- [ ] 每个 SKILL.md 行数不超过 350 行
- [ ] 21 个 SKILL.md 总行数减少 25%+（当前 ~6700 行 → 目标 5000 行以下）
- [ ] 无内部引用断裂（所有 SKILL.md 中引用的文件路径均存在）
- [ ] 每个 commit 仅涉及 1 组 skill，可独立回滚
