---
id: "3"
title: "Update tech-design SKILL.md for refactor intent branch"
priority: "P1"
estimated_time: "1h"
complexity: "low"
dependencies: [1]
type: "doc"
mainSession: false
---

# 3: Update tech-design SKILL.md for refactor intent branch

## Description

更新 tech-design skill 以支持 refactor intent 的内部架构侧重模式。当 intent 为 `refactor` 时，tech-design 跳过 API handbook 和 ER 图等面向外部接口的设计产出，侧重内部架构变更描述（模块重组、依赖关系调整、代码结构优化）。

## Reference Files
- `plugins/forge/skills/tech-design/SKILL.md`: tech-design skill 定义，需添加 refactor 内部分支逻辑 (source: proposal.md#In-Scope, item 4)
- `docs/proposals/intent-driven-pipeline-branching/proposal.md#Feasibility-Assessment`: 定义了 refactor 下 tech-design 的行为：侧重内部架构，跳过 API handbook / ER 图

## Affected Files

### Create
| File | Description |
|------|-------------|
| (无新文件) | |

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/tech-design/SKILL.md` | 添加 refactor intent 内部分支：侧重内部架构，跳过 API handbook 和 ER 图 |

### Delete
| File | Reason |
|------|--------|
| (无删除) | |

## Acceptance Criteria
- [ ] tech-design SKILL.md 包含 intent 检测逻辑：当 proposal.md frontmatter 的 `intent` 为 `refactor` 时，执行内部架构侧重分支
- [ ] refactor 分支下不生成 API handbook 文件和 ER 图文件
- [ ] refactor 分支下不生成 `prd-user-stories.md` 文件

## Implementation Notes

- `intent: new-feature` 时 tech-design 行为完全不变
- refactor 的 tech-design 仍需足够的架构信息供 breakdown-tasks 使用
