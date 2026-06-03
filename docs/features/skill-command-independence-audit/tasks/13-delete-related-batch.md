---
id: "13"
title: "Delete Related sections in consolidate-specs, run-tests, ui-design"
priority: "P1"
estimated_time: "1h"
dependencies: []
type: "doc"
complexity: "medium"
mainSession: false
---

# 13: Delete Related sections in consolidate-specs, run-tests, ui-design

## Description
consolidate-specs、run-tests、ui-design 的 Related Skills/Integration/References 章节内容均可从正文隐含推断，需删除。

## Reference Files
- `docs/proposals/skill-command-independence-audit/proposal.md` — Scope > In Scope, Key Risks
- plugins/forge/skills/consolidate-specs/SKILL.md: 删除 Related Skills/Integration/References 章节 (ref: Scope)
- plugins/forge/skills/run-tests/SKILL.md: 删除 Related Skills/Integration/References 章节 (ref: Scope)
- plugins/forge/skills/ui-design/SKILL.md: 删除 Related Skills/Integration/References 章节 (ref: Scope)

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| plugins/forge/skills/consolidate-specs/SKILL.md | 删除 Related Skills/Integration/References 章节 |
| plugins/forge/skills/run-tests/SKILL.md | 删除 Related Skills/Integration/References 章节 |
| plugins/forge/skills/ui-design/SKILL.md | 删除 Related Skills/Integration/References 章节 |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] 三个 skill 的 Related Skills、Integration、References 章节已删除
- [ ] 删除的内容均可从正文中隐含推断，无独有信息丢失
- [ ] 仅修改以下文件：consolidate-specs/SKILL.md、run-tests/SKILL.md、ui-design/SKILL.md

## Hard Rules
- 仅修改以下文件：plugins/forge/skills/consolidate-specs/SKILL.md、plugins/forge/skills/run-tests/SKILL.md、plugins/forge/skills/ui-design/SKILL.md

## Implementation Notes
这三个 skill 的 Related 章节都是纯 pipeline 上下游信息，无独有知识。删除前确认被删内容可从正文隐含推断。
