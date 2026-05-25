---
created: "2026-05-25"
author: "faner"
status: Draft
---

# Proposal: 领域级自由专家生成机制

## Problem

自由专家（freeform expert）的动态生成以单一 proposal 为锚点，生成的专家领域范围极窄（如 "Build Orchestration & Test Infrastructure Expert" 仅服务于 surface-aware-justfile proposal），导致跨 proposal 复用率近乎为零——11 个已生成的专家之间 Jaccard 相似度无法达到 0.3 的复用阈值。

### Evidence

- `docs/experts/` 下 11 个专家文件，每个的 `generated_for` 都指向唯一的 proposal
- `domain` 关键词高度特化：如 "pipeline-integration, type-system-categorization, dead-code-removal, go-backend, prompt-template-architecture"——5 个关键词中任何一个出现在其他 proposal 的概率极低
- 实际复用匹配从未成功过——每次评估都触发了全新专家生成

### Urgency

每次评估都经历完整的专家生成→确认循环（3 轮修改/拒绝上限），增加评审耗时。领域级专家一次生成即可服务同一领域的多个 proposal，显著降低评审启动成本。

## Proposed Solution

改造 `expert-inference.md` 的专家生成流程，引入预定义领域分类表实现分层识别：

1. **大领域匹配**：从 proposal 提取特征后，查分类表确定所属大领域（如"构建与测试基础设施"）
2. **领域内专家生成**：LLM 在该大领域范围内生成专家，关键词和背景覆盖整个领域而非单一 proposal

同时更新 `expert-template.md` 增加 `scope` 字段，更新 `freeform-expert-persistence.md` 的复用匹配逻辑以适配新格式。

### Innovation Highlights

**分层领域识别**：通过预定义分类表保证领域标签一致性。不同 proposal 对同一领域的识别结果相同（不会出现 "test infrastructure" vs "testing pipeline" 的不一致），从根本上解决复用匹配的可靠性问题。大多数专家系统要么用固定专家库（无灵活性），要么完全依赖 LLM 自由推断（无一致性）——分层方案在两者之间取得平衡。

## Requirements Analysis

### Key Scenarios

- **场景 1（新领域首次评估）**：评估一个属于"构建与测试基础设施"领域的 proposal → 分类表匹配成功 → 生成该领域专家 → 保存到 `docs/experts/`
- **场景 2（同领域再次评估）**：评估另一个同领域 proposal → 复用匹配命中已有专家 → 跳过生成 → 直接用于评审
- **场景 3（跨领域 proposal）**：proposal 涉及多个领域（如"Agent架构" + "配置Schema"）→ 匹配最相关的一个大领域 → 生成或复用对应专家
- **场景 4（分类表未覆盖）**：proposal 属于分类表外的领域 → LLM 自由推断领域 → 生成专家，新领域可后续加入分类表

### Non-Functional Requirements

- **可扩展性**：领域分类表应易于扩展，新增领域只需修改 prompt 文件
- **向后兼容**：现有 11 个 proposal-specific 专家文件保留不变，不被新系统干扰

### Constraints & Dependencies

- 改动仅限 `experts/freeform/` 目录下的 prompt 文件和 `rules/freeform-expert-persistence.md`
- 不影响 freeform-review-protocol、scorer-composition、reviser-composition
- 用户确认循环（Accept / Modify / Regenerate）保持不变

## Alternatives & Industry Benchmarking

### Industry Solutions

业内常见的专家/评审者选择机制：
1. **固定专家库**（如学术同行评审的 TPC 成员列表）——预先定义角色，不动态生成
2. **LLM 自由推断**（如 ChatGPT 的 Custom Instructions persona）——完全依赖模型判断
3. **混合方案**（如 Claude Code 的 multi-expert parallel scoring）——部分固定 + 部分动态

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零改动成本 | 每次评估重复造轮子，复用率 0% | Rejected: 核心问题不会自行消失 |
| 纯 Prompt 重写 | LLM self-guided | 改动最小 | 领域标签不一致，复用匹配不可靠 | Rejected: 没解决一致性问题 |
| 固定专家库 | 学术 TPC 模式 | 最强一致性 | 无法适应未知领域，维护成本高 | Rejected: 不够灵活 |
| **分类表引导** | 分层识别 | 一致性 + 灵活性兼顾 | 分类表需维护 | **Selected: 一致性与灵活性最优平衡** |

## Feasibility Assessment

### Technical Feasibility

完全可行。改动仅涉及 prompt 文件（Markdown），不涉及代码变更。

### Resource & Timeline

单次 prompt 改写 + 测试验证，工作量小。

### Dependency Readiness

无外部依赖。所有改动在 Forge plugin 内部完成。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| 分类表会限制专家的领域覆盖 | Assumption Flip | Refined: 分类表只定义大领域方向，LLM 在领域内仍有充分细化空间。未覆盖的领域有降级路径 |
| Jaccard 匹配足以区分新旧专家 | Occam's Razor | Confirmed: 新专家关键词更广，自然匹配度更高。旧专家会在自然竞争中逐渐被淘汰 |
| 每个 proposal 只需一个领域专家 | Stress Test | Refined: 跨领域 proposal 匹配最相关的一个大领域。若评审质量不足，用户可通过 Modify 循环调整焦点 |

## Scope

### In Scope
- `expert-inference.md` 嵌入领域分类表，改造为两步生成流程（大领域匹配 → 领域内专家生成）
- `expert-template.md` 增加 `scope` 字段（`domain-level` / `proposal-specific`）
- `freeform-expert-persistence.md` 更新复用匹配逻辑以适配 domain-level 专家

### Out of Scope
- 推广自由专家到 PRD / tech-design / ui-design 评估（todo #166）
- 修改 freeform-review-protocol、scorer-composition、reviser-composition
- 修改 freeform-pipeline 编排流程
- 迁移或废弃现有 11 个专家文件
- 复用匹配对旧专家的兼容

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 分类表初期覆盖不全，部分 proposal 无法匹配 | M | L | 提供 LLM 自由推断降级路径，未匹配时自动降级 |
| 领域级专家的评审深度不如 proposal-specific 专家 | M | M | LLM 在领域内细化专业方向时参考 proposal 内容，保持针对性 |
| 分类表随时间膨胀难以维护 | L | L | 分类表控制在大领域粒度（8-12 个），不细化到子领域 |

## Success Criteria

- [ ] 新生成的专家 `domain` 关键词覆盖范围 ≥ 2 个 proposal 的领域交集（验证：用已有 proposal 交叉比对）
- [ ] 同领域内的第二个 proposal 评估时，复用匹配成功（当前为 0 成功）
- [ ] 专家生成后经用户确认的轮次 ≤ 2（当前经常需要修改以扩大领域范围）
- [ ] 分类表覆盖 ≥ 80% 的已有 proposal（验证：将已有 proposal 逐一匹配分类表）

## Next Steps

- Proceed to `/quick-tasks` to generate implementation tasks
