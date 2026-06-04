---
id: "6"
title: "Inline templates/rules into fix-bug command + reduce Knowledge Review"
priority: "P0"
estimated_time: "1h"
dependencies: []
type: "doc"
complexity: "medium"
mainSession: false
---

# 6: Inline templates/rules into fix-bug command + reduce Knowledge Review

## Description
fix-bug command 引用 learn/templates/ 和 consolidate-specs/rules/，需将模板决策点和 spec 提取规则内联。同时精简 Knowledge Review 段落。

## Reference Files
- `docs/proposals/skill-command-independence-audit/proposal.md` — Scope > In Scope, Key Risks
- plugins/forge/commands/fix-bug.md: 需内联模板和规则并精简 Knowledge Review (ref: Scope)
- plugins/forge/skills/learn/templates/: 模板决策点需内联 (ref: Scope)
- plugins/forge/skills/consolidate-specs/rules/: spec 提取规则需内联 (ref: Scope)

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| plugins/forge/commands/fix-bug.md | 内联模板决策点和 spec 提取规则，精简 Knowledge Review 段落 |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] fix-bug.md 包含内联的模板决策点和 spec 提取规则（~40 行），带 `<!-- INLINE:origin=... -->` 标记
- [ ] Knowledge Review 段落已精简
- [ ] fix-bug 不再包含对 learn 或 consolidate-specs 内部文件的跨引用

## Hard Rules
- 内联段落分别添加 `<!-- INLINE:origin=learn/templates/ -->` 和 `<!-- INLINE:origin=consolidate-specs/rules/ -->` 标记

## Implementation Notes
INJECT: 模板决策点 + spec 提取规则 ~40 行；SKIP: 其他 command 的上下文。
