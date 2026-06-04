---
id: "5"
title: "Inline test-type model into init-justfile + reduce examples"
priority: "P0"
estimated_time: "1h"
dependencies: []
type: "doc"
complexity: "medium"
mainSession: false
---

# 5: Inline test-type model into init-justfile + reduce examples

## Description
init-justfile 引用 test-guide/references/test-type-model.md，需将 test-type 映射表内联。同时精简 justfile 示例。

## Reference Files
- `docs/proposals/skill-command-independence-audit/proposal.md` — Scope > In Scope, Key Risks
- plugins/forge/skills/init-justfile/SKILL.md: 需内联 test-type 映射表并精简 justfile 示例 (ref: Scope)
- plugins/forge/skills/test-guide/references/test-type-model.md: test-type 映射表 ~30 行需内联 (ref: Scope)

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| plugins/forge/skills/init-justfile/SKILL.md | 内联 test-type 映射表，精简 justfile 示例 |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] init-justfile/SKILL.md 包含内联的 test-type 映射表（~30 行），带 `<!-- INLINE:origin=... -->` 标记
- [ ] justfile 示例已精简，保留关键模式
- [ ] init-justfile 不再包含对 test-guide 内部文件的跨 skill 引用

## Hard Rules
- 内联段落必须添加 `<!-- INLINE:origin=test-guide/references/test-type-model.md -->` 标记以提供可追溯性

## Implementation Notes
INJECT: test-type 映射表 ~30 行；SKIP: 详细说明和示例。
