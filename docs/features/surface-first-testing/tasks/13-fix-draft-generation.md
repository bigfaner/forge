---
id: "13"
title: "Rewrite test-guide draft-generation.md for surface-first + remove orphan template"
priority: "P0"
estimated_time: "1.5h"
dependencies: []
type: "doc"
mainSession: false
---

# 13: Rewrite test-guide draft-generation.md for surface-first + remove orphan template

## Description
`test-guide/rules/draft-generation.md` 整体描述的是旧 framework-first 模型（4-section schema: framework, discovery, structure, assertions），与 surface-first 模型的 7+1 section 结构严重矛盾。模板路径引用 `docs/conventions/testing/go.md` 等扁平文件，应改为 `testing/<surface>/core.md`。同时 `test-guide/templates/convention-template.md` 使用旧 4-section 格式，且未被 SKILL.md 任何步骤引用，是孤儿文件。

## Reference Files
- `docs/proposals/surface-first-testing/proposal.md` — Proposed Solution
- `plugins/forge/skills/test-guide/rules/draft-generation.md`: 4-section schema、旧路径、旧命名 (ref: Proposed Solution)
- `plugins/forge/skills/test-guide/templates/convention-template.md`: 旧 4-section 格式，未被引用
- `plugins/forge/skills/test-guide/rules/convention-structure.md`: surface-first 7+1 section 权威定义

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/test-guide/rules/draft-generation.md` | 重写为 surface-first 生成流程：7+1 section schema、`testing/<surface>/core.md` 路径、surface strategy template 引用 |

### Delete
| File | Reason |
|------|--------|
| `plugins/forge/skills/test-guide/templates/convention-template.md` | 旧 4-section 格式，已被 `templates/surfaces/*.md` 替代，未被引用 |

## Acceptance Criteria
- [ ] draft-generation.md 的 section schema 从 4-section 改为与 convention-structure.md 一致的 7+1 section
- [ ] 模板路径从 `docs/conventions/testing/<framework>.md` 改为 `docs/conventions/testing/<surface>/core.md`
- [ ] 文件命名从 `<scope>.md` 改为 surface-first 目录结构
- [ ] Built-in Template 查找表更新为 `templates/surfaces/<surface>.md`（而非旧 framework name）
- [ ] `templates/convention-template.md` 已删除
- [ ] draft-generation.md 与 convention-structure.md 和 SKILL.md 的生成流程一致

## Hard Rules
- 必须先加载 `docs/conventions/forge-distribution.md`
- section 定义必须与 convention-structure.md 严格一致
- 不引用其它 skill 的内部文件

## Implementation Notes
- convention-structure.md 是 surface-first Convention 结构的权威定义，draft-generation.md 的所有 section 和路径引用都应对齐
- 旧的 Inferred Defaults 表（runner/assertion/tag per framework）保留，但改为辅助参考而非驱动生成流程
