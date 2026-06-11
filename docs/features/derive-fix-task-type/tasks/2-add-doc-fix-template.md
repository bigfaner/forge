---
id: "2"
title: "Add doc-fix task template"
priority: "P0"
estimated_time: "30min"
dependencies: [1]
type: "doc"
mainSession: false
---

# 2: Add doc-fix task template

## Description

Create the `doc-fix.md` task template for doc-category fix tasks. This template is used by `forge task add --type doc.fix` to generate properly scoped fix task files. Unlike `coding.fix` tasks, doc fix tasks must skip code quality gates (lint, compile, test) and focus on markdown/content fixes.

## Reference Files

- `forge-cli/pkg/task/templates/coding.fix.md`: Reference for template structure — frontmatter (type, category, identity, context sections) and body format with Go template placeholders (source: proposal.md#In-Scope)
- `forge-cli/pkg/task/tasktemplate.go`: Template loading uses `autogenTemplateFS.ReadFile("templates/" + name + ".md")` — file must be placed in this directory (source: proposal.md#Constraints-&-Dependencies)

## Affected Files

### Create
| File | Description |
|------|-------------|
| `forge-cli/pkg/task/templates/doc-fix.md` | Task template for doc-category fix tasks |

### Modify
| File | Changes |
|------|---------|

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria

- [ ] `doc-fix.md` template exists at `forge-cli/pkg/task/templates/doc-fix.md`
- [ ] Template contains fix instructions scoped to doc-type failures: no code quality gates, no test execution, only markdown/content fixes
- [ ] `GetTaskTemplate("doc.fix")` returns the template content without error

## Implementation Notes

### Template design

Adapt from `coding.fix.md` but remove:
- Surface inference section (doc tasks have no surface)
- Fix boundaries related to dev servers, npm, test execution
- Verification section with `go test` commands

Keep:
- Frontmatter with type/category/identity/context
- Go template placeholders for ID, Title, Description, SourceTaskID
- Root Cause and Reference Files sections
- Auto-restore note for source task unblocking

### Embed FS

The template must be placed in `forge-cli/pkg/task/templates/` which is embedded via `autogenTemplateFS`. No separate embed registration needed — Go's `embed` directive picks up all files in the directory.
