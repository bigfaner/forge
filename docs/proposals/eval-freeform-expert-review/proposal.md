---
created: 2026-05-23
author: faner
status: Draft
---

# Proposal: eval-proposal 自由专家评审前置阶段

## Problem

eval-proposal 的评审体验过于机械化——固定 rubric 维度强制评分、三阶段流水线协议（推理审计→打分→盲点搜索），导致评审产出像一份检查清单而非专家洞察。rubric 覆盖了已知失败模式，但无法发现文档特有的、超出预设维度的问题。

### Evidence

- 现有 eval-proposal 使用 10 维度 / 1000 分 rubric，CTO 专家角色固定，评审流程完全模板化
- 用户在实际使用中感受到「机械感」：每一轮评审产出格式相同、关注点可预测
- rubric 维度的设计基于通用失败模式（如「隐藏成本」「回滚计划缺失」），无法针对特定提案的独特风险点进行深度挖掘

### Urgency

eval 是 Forge 质量保障体系的核心。如果评审产出缺乏真正的洞察力，提案质量的上限被 rubric 天花板锁死。越早引入自由评审，越早突破这个上限。

## Proposed Solution

为 eval-proposal 增加 `--freeform-expert` 参数，启用后在 rubric 循环之前插入 **Phase 0 自由专家评审**阶段。不传该参数时行为与现有 eval-proposal 完全一致：

1. **参数控制**：`forge eval --type proposal --freeform-expert` 启用自由专家阶段；不传参数时走标准 rubric 流程
2. **动态专家生成**：分析提案内容（domain、技术栈、复杂度、关键决策），推断最适合评审的专家档案（背景、专业领域、评审风格），用户确认后使用
3. **自由叙事评审**：该专家以纯叙事形式对提案进行深度评审——无 rubric、无评分、无预设维度，完全由专家自主决定关注什么
4. **发现提取与注入**：从自由评审叙事中提取结构化发现（key findings），注入后续 rubric scorer 的 prompt 作为额外上下文
5. **专家持久化与复用**：动态专家档案保存到 `docs/experts/` 全局目录，后续评审可复用已有专家

### Innovation Highlights

**动态专家生成**区别于行业常见的静态角色定义（如「你是一个 CTO」）。系统根据文档内容推断专家背景，使评审视角与文档特性匹配。例如：一个后端性能优化提案可能生成「分布式系统架构师，专注高并发场景」，而一个用户体验提案可能生成「产品心理学家，擅长行为设计分析」。

**自由叙事 → 结构化注入**的管道设计，兼顾了自由度和系统性：自由评审阶段不受 rubric 约束，但产出经过提取后进入 rubric 循环，确保后续评分能覆盖自由评审发现的盲点。

**专家库积累**：随使用积累的 `docs/experts/` 目录形成可复用的专家库，越用越丰富。

## Requirements Analysis

### Key Scenarios

- **标准评审（无参数）**：`forge eval --type proposal` → 走现有 rubric 流程，行为完全不变
- **自由专家评审**：`forge eval --type proposal --freeform-expert` → 先走 Phase 0 自由评审，再走 rubric 流程
- **新专家生成**：评审一篇关于「测试框架插件化」的提案 → 系统推断需要一位「测试工具链架构师」→ 用户确认 → 专家进行自由评审
- **已有专家复用**：评审另一篇类似领域的提案 → 系统在 `docs/experts/` 中找到匹配的已有专家 → 用户确认复用
- **用户修改专家**：系统推断的专家不合适 → 用户修改专家档案 → 保存修改后的版本 → 继续评审
- **自由评审发现注入**：自由评审发现了「提案的扩展性假设缺乏验证」→ 该发现被提取为 attack point → 注入 rubric scorer → 后续评分覆盖此维度

### Non-Functional Requirements

- **延迟容忍**：Phase 0 增加约一次 LLM 调用（专家生成 + 自由评审），eval 总时间增加约 30%——可接受
- **可审计性**：所有动态专家档案持久化到 `docs/experts/`，用户可事后审核
- **确定性**：同一提案 + 同一专家应产出方向一致的评审（非完全随机）

### Constraints & Dependencies

- 仅适用于 `eval --type proposal`
- 依赖现有 eval skill 的 scorer/reviser 子 agent 架构
- 专家档案格式需兼容现有 `experts/scorer/*.md` 的 prompt 格式
- 需遵守 `docs/conventions/forge-distribution.md` 的分发规范

## Alternatives & Industry Benchmarking

### Industry Solutions

代码评审领域的趋势是从「检查清单」走向「上下文感知评审」。GitHub Copilot Code Review、Phabricator 的 Herald Rules 都在尝试根据代码内容调整评审策略。学术界（如 Microsoft Research 的 Code质量研究）也表明领域专家的直觉性判断常优于结构化检查表。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零成本，不增加复杂度 | 不解决「太机械」的痛点 | Rejected: 用户明确要求改变 |
| 混合首轮（自由+rubric 同步） | 内部方案 | 不增加迭代次数 | rubric 影响会渗透到「自由」评审，不够纯粹 | Rejected: 不够自由 |
| **Pre-scorer 前置阶段** | 本提案 | 纯增量修改，自由评审完全独立于 rubric | 多一轮 LLM 调用 | **Selected: 最纯粹的自由评审实现** |

## Feasibility Assessment

### Technical Feasibility

完全可行。核心改动点：
1. eval skill / eval-proposal command 增加 `--freeform-expert` 参数解析
2. eval skill 的 proposal 类型处理流程中增加 Phase 0（仅在参数启用时进入）
3. 新增一个「专家推断 + 自由评审」子 agent prompt
4. 新增发现提取逻辑
5. 修改 rubric scorer prompt 以接收注入的发现

所有改动均在现有 eval skill 架构内完成，无需新框架或外部依赖。

### Resource & Timeline

预计 4-6 个 coding task：
1. `--freeform-expert` 参数解析 + 条件分支（command + skill SKILL.md）
2. 专家推断逻辑 + 专家档案模板
3. 自由评审子 agent prompt + 协议
4. 发现提取 + 注入机制
5. eval skill 集成（Phase 0 编排）
6. 专家持久化与复用逻辑

加上 doc 类型 task（专家档案模板文档、协议文档），总量在 quick mode 范围内。

### Dependency Readiness

无外部依赖。所有改动在 `plugins/forge/skills/eval/` 目录内完成。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| 自由叙事评审比 rubric 打分更有洞察力 | Assumption Flip | Confirmed: rubric 覆盖已知模式但遗漏文档特有问题；自由评审捕获盲点但可能遗漏已知陷阱。两者互补而非替代 |
| 动态生成的专家足够可靠 | Stress Test | Refined: 需要用户确认环节作为安全网，避免不合适的专家产出低质量评审 |
| Phase 0 增加的延迟可接受 | Occam's Razor | Confirmed: 一次额外 LLM 调用（约 30% 时间增加）换来显著提升的评审质量，ROI 合理 |

## Scope

### In Scope

- `--freeform-expert` 参数解析与条件分支（启用 / 未启用两条路径）
- 动态专家档案推断机制（分析提案 → 推断专家 → 生成详细档案）
- 用户确认机制（接受 / 修改 / 重新生成）
- 专家档案持久化到 `docs/experts/` 全局目录
- 已有专家复用逻辑（匹配提案内容与已有专家档案）
- 自由叙事评审协议（纯叙事、无 rubric、无评分）
- 发现提取机制（从叙事中提取结构化 key findings）
- 注入机制（将发现作为额外上下文注入 rubric scorer）
- eval skill 的 proposal 类型集成（Phase 0 编排）

### Out of Scope

- 扩展到其他 eval 类型（prd、design、ui 等）——未来可复用相同架构
- 修改现有 proposal rubric 本身
- 多专家并行自由评审
- 对自由评审产出进行评分
- 任何 UI/交互式组件

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| AI 推断的专家角色不合适 | M | M | 用户确认环节可修改或重新生成 |
| 自由叙事产出难以提取结构化发现 | M | H | 设计明确的提取协议，要求评审产出中标注关键段落 |
| 动态专家引入不可预测性 | L | L | 这是特性的核心价值（发现未知盲点），不是缺陷 |
| 专家库膨胀导致匹配困难 | L | L | 命名约定 + domain 标签帮助检索 |

## Success Criteria

- [ ] eval-proposal 不传 `--freeform-expert` 时行为与现有版本完全一致（零回归）
- [ ] 传入 `--freeform-expert` 时进入 Phase 0，生成动态专家档案并经用户确认
- [ ] 自由评审产出为纯叙事格式，无 rubric 维度、无评分
- [ ] 自由评审的 key findings 成功提取并注入后续 rubric scorer
- [ ] 动态专家档案保存到 `docs/experts/`，可事后审核
- [ ] 后续评审可复用 `docs/experts/` 中的已有专家
- [ ] eval 总时间增加不超过 40%（基线：当前单轮 eval 时间）

## Next Steps

- Proceed to `/quick-tasks` to generate implementation tasks
