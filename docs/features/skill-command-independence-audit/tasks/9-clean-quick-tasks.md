---
id: "9"
title: "Delete Integration + reduce shared content in quick-tasks"
priority: "P1"
estimated_time: "1h"
dependencies: []
type: "doc"
complexity: "medium"
mainSession: false
---

# 9: Delete Integration + reduce shared content in quick-tasks

## Description
quick-tasks 仅需删除 ## Integration 段落（保留 ## Reference Files，是模板占位符 {{REFERENCE_FILES}} 的替换规则说明）。同时精简与 breakdown-tasks 共享的 ~150 行内容。

## Reference Files
- `docs/proposals/skill-command-independence-audit/proposal.md` — Scope > In Scope, Key Risks, Success Criteria
- plugins/forge/skills/quick-tasks/SKILL.md: 删除 Integration 段落并精简与 breakdown-tasks 共享的内容 (ref: Scope)
- plugins/forge/skills/breakdown-tasks/SKILL.md: 共享内容参考，确保精简后两 skill 各自自洽 (ref: Scope)

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| plugins/forge/skills/quick-tasks/SKILL.md | 删除 ## Integration 段落，精简与 breakdown-tasks 共享的内容 |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] ## Integration 段落已删除
- [ ] ## Reference Files 段落已保留（模板占位符替换规则说明）
- [ ] 与 breakdown-tasks 共享的 ~150 行内容已精简，保留 quick-tasks 专属逻辑

## Hard Rules
- 必须保留 ## Reference Files 段落（非 pipeline 上下游信息，而是模板占位符 {{REFERENCE_FILES}} 的替换规则说明）

## Implementation Notes
quick-tasks 和 breakdown-tasks 共享 task 生成逻辑但有不同的模板和流程。精简时保留各 skill 的独特部分。
