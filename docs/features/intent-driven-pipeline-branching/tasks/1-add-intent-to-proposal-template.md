---
id: "1"
title: "Add intent field to proposal template and brainstorm skill"
priority: "P0"
estimated_time: "1-2h"
complexity: "high"
dependencies: []
type: "doc"
mainSession: false
---

# 1: Add intent field to proposal template and brainstorm skill

## Description

为 proposal 模板添加 `intent` frontmatter 字段，并为 brainstorm skill 添加 intent 推断步骤。这是 Intent-Driven Pipeline Branching 的基础——intent 字段是后续所有 pipeline 选择的数据源。

Proposal 模板当前只有 `created`/`author`/`status` 字段，需要添加 `intent` 字段（默认值 `new-feature`）。Brainstorm skill 需要在生成 proposal 时推断 intent，通过 AskUserQuestion 让用户确认后写入 frontmatter。

## Reference Files
- `plugins/forge/skills/brainstorm/templates/proposal.md`: 当前模板无 intent 字段，需在 frontmatter 中添加 `intent: new-feature` 作为默认值 (source: proposal.md#In-Scope, item 1)
- `plugins/forge/skills/brainstorm/SKILL.md`: brainstorm skill 定义，需添加 intent 推断步骤和 AskUserQuestion 确认 (source: proposal.md#In-Scope, item 2)
- `docs/proposals/intent-driven-pipeline-branching/proposal.md#In-Scope`: 定义了 intent 的三个有效值 (`new-feature`/`refactor`/`cleanup`) 及 task type → intent 映射规则

## Affected Files

### Create
| File | Description |
|------|-------------|
| (无新文件) | |

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/brainstorm/templates/proposal.md` | frontmatter 添加 `intent: "new-feature"` 字段 |
| `plugins/forge/skills/brainstorm/SKILL.md` | 添加 intent 推断步骤：AI 根据 proposal 内容推断 intent，用 AskUserQuestion 确认 |

### Delete
| File | Reason |
|------|--------|
| (无删除) | |

## Acceptance Criteria
- [ ] `plugins/forge/skills/brainstorm/templates/proposal.md` frontmatter 包含 `intent: "new-feature"` 字段，位于 `status` 之后
- [ ] brainstorm SKILL.md 包含 intent 推断步骤：AI 根据 proposal 内容和 task type → intent 映射规则推断 intent（`coding.feature`/`coding.enhancement` → `new-feature`，`coding.cleanup` → `cleanup`，`coding.refactor` → `refactor`，`coding.fix` 按"是否引入新的用户可观测行为"判断）
- [ ] brainstorm 使用 AskUserQuestion 展示推断的 intent，用户可覆盖确认后写入 proposal.md frontmatter
- [ ] 对 `coding.fix` 类型 proposal，intent 推断逻辑正确：引入新用户可观测行为 → `new-feature`，仅内部调整 → `refactor`
- [ ] 混合内容 proposal（既有新行为又有重构）按"是否引入新的用户可观测行为"判断主要 intent，用户可在确认阶段覆盖

## Implementation Notes

- proposal.md 模板的 `intent` 字段值必须用双引号包裹（`intent: "new-feature"`），与其他 frontmatter 字段风格一致
- intent 推断应基于 proposal 的 **Proposed Solution** 和 **Scope** 内容，而非仅看标题
- 映射规则来源：proposal.md 的 "Feature Intent" 定义段落
