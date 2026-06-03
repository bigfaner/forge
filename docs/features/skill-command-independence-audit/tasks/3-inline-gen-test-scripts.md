---
id: "3"
title: "Inline test isolation into gen-test-scripts + delete Related sections"
priority: "P0"
estimated_time: "1h"
dependencies: []
type: "doc"
complexity: "medium"
mainSession: false
---

# 3: Inline test isolation into gen-test-scripts + delete Related sections

## Description
gen-test-scripts 引用 run-tests/rules/test-isolation.md，需将隔离策略决策表内联。同时删除 Related Skills/Integration/References 章节。

## Reference Files
- `docs/proposals/skill-command-independence-audit/proposal.md` — Scope > In Scope, Key Risks
- plugins/forge/skills/gen-test-scripts/SKILL.md: 需内联隔离策略并删除 Related/Integration/References (ref: Scope)
- plugins/forge/skills/run-tests/rules/test-isolation.md: 隔离策略决策表 ~40 行需内联 (ref: Scope)

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| plugins/forge/skills/gen-test-scripts/SKILL.md | 内联隔离策略决策表，删除 Related Skills/Integration/References 章节 |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] gen-test-scripts/SKILL.md 包含内联的隔离策略决策表（~40 行），带 `<!-- INLINE:origin=... -->` 标记
- [ ] Related Skills、Integration、References 章节已删除
- [ ] gen-test-scripts 不再包含对 run-tests 内部文件的跨 skill 引用

## Hard Rules
- 内联段落必须添加 `<!-- INLINE:origin=run-tests/rules/test-isolation.md -->` 标记以提供可追溯性
- 保留所有 HARD-RULE / HARD-GATE / EXTREMELY-IMPORTANT / PROHIBITIONS 块计数不变

## Implementation Notes
INJECT: 隔离策略决策表 ~40 行；SKIP: run-tests 调度逻辑。
