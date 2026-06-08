---
id: "4"
title: "Add AC rendering section to auto-gen task templates"
priority: "P0"
estimated_time: "30min"
dependencies: []
type: "doc"
mainSession: false
---

# 4: Add AC rendering section to auto-gen task templates

## Description

3 个自动生成任务模板（`test-gen-scripts.md`、`doc-drift.md`、`doc-review.md`）在 frontmatter `context:` 中声明了 `AcceptanceCriteria`，但模板 body 中没有 `## Acceptance Criteria` section 和 `{{.AcceptanceCriteria}}` 渲染指令。导致 `forge task index` 生成的 .md 文件缺少 AC section，`forge task validate` 报 0 AC 错误。

`autogen.go` 的 `buildAutogenTemplateData` 在 AC 为空时已有 fallback（`"- [ ] All acceptance criteria met"`），只需在模板 body 中添加渲染指令即可生效。

注：`test-run.md` 模板由 Task 2 单独处理（同时补 AC 渲染 + 硬化 AC 内容）。

## Reference Files
- `docs/proposals/test-pipeline-interleaved/proposal.md` — Proposed Solution, Scope > In Scope, Success Criteria
- `docs/lessons/gotcha-autogen-template-missing-ac-section.md` — 根因分析
- `forge-cli/pkg/task/templates/test-gen-scripts.md`: add AC rendering section (ref: Scope > In Scope)
- `forge-cli/pkg/task/templates/doc-drift.md`: add AC rendering section (ref: Scope > In Scope)
- `forge-cli/pkg/task/templates/doc-review.md`: add AC rendering section (ref: Scope > In Scope)

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| `forge-cli/pkg/task/templates/test-gen-scripts.md` | body 末尾追加 `## Acceptance Criteria\n\n{{.AcceptanceCriteria}}` |
| `forge-cli/pkg/task/templates/doc-drift.md` | body 末尾追加 `## Acceptance Criteria\n\n{{.AcceptanceCriteria}}` |
| `forge-cli/pkg/task/templates/doc-review.md` | body 末尾追加 `## Acceptance Criteria\n\n{{.AcceptanceCriteria}}` |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria

- [ ] `test-gen-scripts.md` 模板 body 包含 `## Acceptance Criteria` section，渲染 `{{.AcceptanceCriteria}}`
- [ ] `doc-drift.md` 模板 body 包含 `## Acceptance Criteria` section，渲染 `{{.AcceptanceCriteria}}`
- [ ] `doc-review.md` 模板 body 包含 `## Acceptance Criteria` section，渲染 `{{.AcceptanceCriteria}}`
- [ ] `forge task validate` 对这 3 个自动生成任务不再报 "has 0 acceptance criteria"

## Implementation Notes

参考已有正确实现的模板：`validation-code.md`、`validation-ux.md`、`test-gen-contracts.md`、`test-gen-journeys.md`。这些模板都在 body 中包含 `{{.AcceptanceCriteria}}` 渲染。

格式：
```markdown
## Acceptance Criteria

{{.AcceptanceCriteria}}
```

`buildAutogenTemplateData` 在 `ctx.AcceptanceCriteria` 为空时自动 fallback 为 `"- [ ] All acceptance criteria met"`，所以无需在模板中硬编码默认值。
