---
id: "2"
title: "Update write-prd SKILL.md for refactor intent branch"
priority: "P1"
estimated_time: "1h"
complexity: "low"
dependencies: [1]
type: "doc"
mainSession: false
---

# 2: Update write-prd SKILL.md for refactor intent branch

## Description

更新 write-prd skill 以支持 refactor intent 的 spec-only PRD 格式。当 intent 为 `refactor` 时，PRD 不生成 user stories（"As a user / I want / So that" 格式对纯重构语义为空），改为生成包含"变更范围 + 约束条件 + 验证标准"三个必需字段的 spec。

## Reference Files
- `plugins/forge/skills/write-prd/SKILL.md`: write-prd skill 定义，需添加 refactor 内部分支逻辑 (source: proposal.md#In-Scope, item 3)
- `docs/proposals/intent-driven-pipeline-branching/proposal.md#Proposed-Solution`: 定义了 spec-only PRD 的三个必需字段：变更范围（affected modules/files）、约束条件（behavioral invariants to preserve）、验证标准（regression acceptance criteria）

## Affected Files

### Create
| File | Description |
|------|-------------|
| (无新文件) | |

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/write-prd/SKILL.md` | 添加 refactor intent 内部分支：跳过 user stories 生成，改为生成 spec-only PRD |

### Delete
| File | Reason |
|------|--------|
| (无删除) | |

## Acceptance Criteria
- [ ] write-prd SKILL.md 包含 intent 检测逻辑：当 proposal.md frontmatter 的 `intent` 为 `refactor` 时，执行 spec-only PRD 分支
- [ ] spec-only PRD 格式包含三个必需字段：变更范围（affected modules/files）、约束条件（behavioral invariants to preserve）、验证标准（regression acceptance criteria）
- [ ] refactor 分支下不生成 `prd-user-stories.md` 文件

## Implementation Notes

- spec-only PRD 必须包含足够的信息供 tech-design 使用，无需依赖 user stories
- `intent: new-feature` 时 write-prd 行为完全不变
