---
created: 2026-05-31
author: faner
status: Draft
intent: refactor
---

# Proposal: Intent Enriched Enum

## Problem

Proposal intent 枚举只有 3 个值（`new-feature`、`refactor`、`cleanup`），与 task type 的 8 个值映射断裂，导致下游 pipeline 分支过粗——`refactor`、`cleanup`、`fix` 被等同对待，不同场景走了不合适的流程。

### Evidence

- `coding.fix`、`coding.enhancement`、`doc` 作为 task type 存在，但 proposal intent 没有对应值
- `fix` 在 brainstorm 中被启发式拆分为 `new-feature` 或 `refactor`，没有独立身份
- `refactor` 和 `cleanup` 在 write-prd/tech-design 中被完全等同对待（spec-only PRD），但两者的变更范围和风险等级可能很不同
- 一个改变外部 API 的 refactor 仍跳过 API handbook，因为 pipeline 分支只看 intent 不看内容

### Urgency

随着 Forge 处理的场景增多（bug fix、enhancement、纯文档变更），3 值枚举的覆盖缺口越来越明显。每次遇到非标准场景都需要 LLM 做启发式判断，增加了不一致的风险。

## Proposed Solution

1. **扩充 intent 为 6 值枚举**：`new-feature`、`enhancement`、`refactor`、`cleanup`、`fix`、`doc`，与 task type 形成干净的 1:1 映射
2. **混合模式 pipeline 分支**：intent 控制默认 pipeline 配置（一张表），PRD 内容中的明确信号可以覆盖默认值
3. **简化 brainstorm 推断**：`fix` 始终为 `fix`，移除启发式判断；每个 task type 直接映射到对应 intent

### Innovation Highlights

无特别创新。对标 task type 的现有分类体系，消除 intent 与 type 之间的映射鸿沟。混合模式 pipeline 是对当前二元分支的自然细化。

## Requirements Analysis

### Key Scenarios

1. **Bug fix 提案**：brainstorm 直接推断 `fix`，pipeline 默认 spec-only，跳过 user stories 和 API handbook
2. **Enhancement 提案**：brainstorm 推断 `enhancement`，pipeline 默认跳过 user stories 但保留 test pipeline（改善现有行为需要测试覆盖）
3. **改变外部 API 的 refactor**：brainstorm 推断 `refactor`（默认跳过 API handbook），但 PRD 内容包含"CLI 命令重命名"信号 → 覆盖开启 API handbook
4. **纯文档提案**：brainstorm 推断 `doc`，pipeline 全部跳过，直接进入 task 生成
5. **混合内容提案**：brainstorm 按核心目标推断主 intent，个别 task 通过 per-task type override 覆盖

### Non-Functional Requirements

- 向后兼容：现有的 `new-feature`、`refactor`、`cleanup` 值行为不变
- 一致性：intent-to-type 映射为严格 1:1，消除歧义

### Constraints & Dependencies

- Intent 分支逻辑全部在 skill markdown 中，无 Go 代码依赖
- 变更限于 plugins/forge/ 目录下的 8 个文件

## Alternatives & Industry Benchmarking

### Industry Solutions

分类系统的精确度随场景增长自然演进是常见模式——从粗粒度到细粒度。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零改动 | fix/enhancement/doc 无对应 intent，pipeline 过粗 | Rejected: 覆盖缺口随场景增长扩大 |
| 只扩枚举 | — | 最小改动 | pipeline 分支仍然过粗 | Rejected: 只解决一半问题 |
| **扩枚举 + 混合 pipeline** | CI lint gate 模式 | 完整解决两个动机 | 8 个文件变更 | **Selected: 双重改进** |
| 完全内容驱动 pipeline | — | 最精准 | 依赖 LLM 判断力，不稳定 | Rejected: 基线不可靠 |

## Feasibility Assessment

### Technical Feasibility

完全可行。所有变更是 markdown 编辑——更新推断表、分支表、覆盖规则。无 Go 代码变更。

### Resource & Timeline

中型变更：8 个 markdown 文件。write-prd 和 tech-design 变更量最大（需要重写 pipeline 分支逻辑），其余文件是小幅更新。预计 2-3 个任务可完成。

### Dependency Readiness

无外部依赖。所有 skill 文件已存在且结构清晰。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| "3 值够用" | 5 Whys | Overturned: task type 有 8 个值，3 值 intent 无法覆盖。缺口通过启发式弥补，增加了不一致风险 |
| "fix 需要启发式区分" | Occam's Razor | Overturned: fix 作为独立 intent 更简单。是否有新用户可见行为由 pipeline 覆盖规则处理，不需要在 intent 层面区分 |
| "pipeline 分支只能靠 intent" | Assumption Flip | Refined: intent 提供稳定基线，PRD 内容提供覆盖信号。两者结合比任何单一维度都可靠 |
| "refactor 和 cleanup 应该等同对待" | Stress Test | Overturned: refactor 可能涉及外部接口变更（需要 API handbook），cleanup 不会。等价对待导致 refactor 跳过必要的产物 |

## Scope

### In Scope

- **brainstorm/SKILL.md**：更新 Step 4.5 intent mapping 表（6 值），移除 fix 启发式，更新 AskUserQuestion 选项
- **brainstorm/templates/proposal.md**：更新 intent 有效值注释
- **write-prd/SKILL.md**：将二元分支（new-feature vs refactor/cleanup）替换为 Pipeline Configuration 表 + Override Signals
- **write-prd/rules/self-check.md**：更新 intent-gated 检查为 6 值
- **tech-design/SKILL.md**：将二元分支替换为 Pipeline Configuration 表 + Override Signals
- **tech-design/rules/design-quality-checks.md**：更新 intent-gated 检查为 6 值
- **breakdown-tasks/SKILL.md**：更新 Intent Propagation 为严格 1:1 映射（6 值）
- **quick-tasks/SKILL.md**：更新 Intent Propagation 为严格 1:1 映射（6 值）

### Out of Scope

- Go CLI 代码变更（CLI 不引用 intent 字段）
- 新 skill 或 command 创建
- 已有提案的迁移（旧 3 值行为不变）
- task-sizing-gate 提案（独立提案）

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 6 值枚举仍不够 | L | L | 可以后续继续扩充；6 值已覆盖所有现有 task type（除 doc.consolidate/doc.drift 这两个小众类型） |
| 混合模式的覆盖规则被 LLM 忽略 | M | M | 覆盖信号是明确的条件表，不是模糊指令；LLM 对结构化规则的遵守度高于 prose 描述 |
| write-prd/tech-design 分支重写引入不一致 | M | M | Pipeline Configuration 表统一两处逻辑，减少不一致可能性 |
| 旧提案与新规则不兼容 | L | L | 旧 3 值在表中仍有对应行，行为不变 |

## Success Criteria

- [ ] brainstorm 推断结果为 6 值之一，用户可在 AskUserQuestion 中选择全部 6 个值
- [ ] `fix` 始终推断为 `fix`，不再使用启发式
- [ ] write-prd 和 tech-design 使用统一的 Pipeline Configuration 表（6 行）
- [ ] Override Signals 规则存在且可被 PRD 内容触发
- [ ] breakdown-tasks 和 quick-tasks 的 Intent Propagation 为严格 1:1 映射
- [ ] 现有 `new-feature`、`refactor`、`cleanup` 值的 pipeline 行为不变
- [ ] 8 个文件全部更新，无遗漏

## Next Steps

- Proceed to `/quick-tasks` to generate implementation tasks
