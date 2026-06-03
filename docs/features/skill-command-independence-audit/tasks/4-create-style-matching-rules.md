---
id: "4"
title: "Create style-matching rules in extract-design-md"
priority: "P0"
estimated_time: "30m"
dependencies: []
type: "doc"
complexity: "low"
mainSession: false
---

# 4: Create style-matching rules in extract-design-md

## Description
extract-design-md 对 ui-design/templates/styles/ 的引用是运行时数据读取而非知识引用（设计意图豁免）。创建 rules/style-matching.md 包含匹配特征摘要，风格文件保留在 ui-design 中。

## Reference Files
- `docs/proposals/skill-command-independence-audit/proposal.md` — Scope > In Scope
- plugins/forge/skills/extract-design-md/rules/match-strategy.md: 现有匹配策略，新文件需与之配合 (ref: Scope)

## Affected Files

### Create
| File | Description |
|------|-------------|
| plugins/forge/skills/extract-design-md/rules/style-matching.md | 各风格的匹配特征摘要 |

### Modify
| File | Changes |
|------|---------|
| plugins/forge/skills/extract-design-md/SKILL.md | 引用新创建的 rules/style-matching.md 替代对 ui-design/templates/styles/ 的直接引用说明 |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] `rules/style-matching.md` 已创建，包含各风格的匹配特征摘要
- [ ] SKILL.md 引用新 rules 文件替代对 ui-design/templates/styles/ 的路径说明

## Implementation Notes
此任务不消除跨 skill 引用（运行时数据读取是设计意图），而是通过创建匹配摘要使 extract-design-md 更自洽。风格文件保留在 ui-design 中不动。
