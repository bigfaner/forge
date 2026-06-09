---
id: "4"
title: "Delete language templates and project-detection rule"
priority: "P1"
estimated_time: "0.5h"
dependencies: [2]
type: "doc"
mainSession: false
# Note: surface-key and surface-type fields are intentionally absent from doc tasks.
# Doc tasks produce non-compilable output (markdown, specs, templates) and do not
# interact with the quality gate or test pipeline, so surface routing is unnecessary.
---

# 4: Delete language templates and project-detection rule

## Description

Delete the 6 language template files and the `rules/project-detection.md` rule file. After Task 2 rewrites SKILL.md to use agent-driven generation (no template references) and Task 3 simplifies surface rules (no TODO stub references to templates), these files become dead code. This task is sequenced after Task 2 to ensure no references remain before deletion.

Before deleting, verify with grep that no remaining files in the skill directory reference the deleted files.

## Reference Files
- `docs/proposals/agent-driven-justfile-generation/proposal.md` — Scope > In Scope, Feasibility Assessment > Resource & Timeline (ref: ## Scope > ### In Scope, ## Feasibility Assessment > ### Resource & Timeline)
- `plugins/forge/skills/init-justfile/SKILL.md` — verify no template references after Task 2 rewrite

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|

### Delete
| File | Reason |
|------|--------|
| `plugins/forge/skills/init-justfile/templates/go.just` | Go language template — replaced by agent-driven generation |
| `plugins/forge/skills/init-justfile/templates/node.just` | Node/TypeScript template — replaced by agent-driven generation |
| `plugins/forge/skills/init-justfile/templates/python.just` | Python template — replaced by agent-driven generation |
| `plugins/forge/skills/init-justfile/templates/rust.just` | Rust template — replaced by agent-driven generation |
| `plugins/forge/skills/init-justfile/templates/mixed.just` | Mixed frontend+backend template — replaced by agent-driven generation |
| `plugins/forge/skills/init-justfile/templates/generic.just` | Generic fallback template — replaced by agent-driven generation |
| `plugins/forge/skills/init-justfile/rules/project-detection.md` | Project type detection rule — no longer needed with agent-driven approach |

## Acceptance Criteria
- [ ] 6 个语言模板文件已删除：`go.just`、`node.just`、`python.just`、`rust.just`、`mixed.just`、`generic.just`
- [ ] `rules/project-detection.md` 已删除
- [ ] `templates/` 目录已清空或删除（确认无其他文件残留）
- [ ] grep 验证：`SKILL.md` 和所有 surface rule 文件中无对已删除文件的引用（搜索 `templates/`、`project-detection`、`generic.just` 等关键词）

## Implementation Notes

### 验证步骤
删除后运行以下验证：
```bash
# 确认文件已删除
ls plugins/forge/skills/init-justfile/templates/ 2>/dev/null
ls plugins/forge/skills/init-justfile/rules/project-detection.md 2>/dev/null

# 确认无残留引用
grep -r "templates/" plugins/forge/skills/init-justfile/ --include="*.md"
grep -r "project-detection" plugins/forge/skills/init-justfile/ --include="*.md"
grep -r "generic.just\|go.just\|node.just\|python.just\|rust.just\|mixed.just" plugins/forge/skills/init-justfile/ --include="*.md"
```

### 回退策略（来自 proposal）
- 模板删除前在 git 中打 tag（如 `pre-template-removal`）
- 若 agent 驱动生成在关键场景表现不足，可从 tag 恢复模板作为 fallback
