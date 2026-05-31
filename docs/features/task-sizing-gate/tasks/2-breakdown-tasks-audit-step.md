---
id: "2"
title: "Insert task sizing audit step in breakdown-tasks skill"
priority: "P1"
estimated_time: "1h"
dependencies: [1]
type: "doc"
mainSession: false
---

# 2: Insert task sizing audit step in breakdown-tasks skill

## Description
在 `breakdown-tasks` skill 的 SKILL.md 中，于写 task 文件后、`forge task index` 前插入一个独立的 task sizing audit step。该 step 要求 LLM 对每个已生成的 task 做聚焦自审：检查 multi-verb（title 中是否有连接独立动作的连词）、AC 跨不相关领域（一条 task 的 AC 是否覆盖多个不相关功能），发现问题自动拆分并输出报告。

同时更新 skill 中所有 `validate-index` 引用为 `validate`，并顺延后续 step 编号。

## Reference Files
- `docs/proposals/task-sizing-gate/proposal.md` — Proposed Solution, Scope > In Scope, Key Risks
- `plugins/forge/skills/breakdown-tasks/SKILL.md`: 在 Step 6（Generate index.json）前插入 audit step，更新 validate-index 引用，重新编号后续 steps (ref: Proposed Solution)
- `docs/conventions/forge-distribution.md`: 理解 skill 文件分发模型，确保修改符合路径解析机制 (ref: Constraints & Dependencies)

## Affected Files

### Create
| File | Description |
|------|-------------|
| (无新文件) | |

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/breakdown-tasks/SKILL.md` | 插入 audit step + 更新 validate-index → validate + renumber steps |

### Delete
| File | Reason |
|------|--------|
| (无删除) | |

## Acceptance Criteria
- [ ] `breakdown-tasks` 的 SKILL.md 在写 task 文件后、`forge task index` 前包含独立的 task sizing audit step
- [ ] audit step 指示 LLM 对每个 task 检查 multi-verb 和跨域 AC，发现问题自动拆分并输出报告
- [ ] 所有 `validate-index` 引用已更新为 `validate`
- [ ] 所有 step 编号连续无跳跃

## Implementation Notes
- audit step 应包含明确的检查清单：multi-verb 检测（title 含 and/or/and then 连接独立动词短语）、AC 跨域检测（AC 涵盖不相关功能域）、operational ceiling（>8 文件）。
- audit 发现问题时应自动拆分 task 文件并输出拆分报告（哪些 task 被拆分、原因、新 task 的 title 和 ID）。
- Key Risk: LLM 在 audit step 仍忽略 multi-verb — CLI 层的 AC ≤ 6 校验作为兜底。
