---
created: "2026-05-26"
tags: [testing, architecture]
---

# fix-task Template Missing Type Default Causes CategoryForType Warning

## Problem
`forge task add --template fix-task` creates tasks with `type: ""`, causing `CategoryForType: unknown type "", defaulting to coding` log warnings at claim time. The markdown template `fix-task.md` has `type: "coding.fix"` hardcoded, but the CLI `add.go` doesn't read type from template defaults.

## Root Cause
1. `template.Defaults` struct has no `Type` field — only Priority, Breaking, EstimatedTime, IDPrefix
2. `add.go:164` sets `opts.Type = addType` directly from the `--type` flag (empty when not provided)
3. `add.go:168-186` applies template defaults but skips Type because the struct doesn't have it
4. The rendered `.md` file gets `type: "coding.fix"` from the template, but `index.json` gets `type: ""` from the CLI opts

## Solution
Either:
- Add `Type` field to `template.Defaults` and set `"fix-task": {Type: "coding.fix"}` in the defaults map
- Or in `add.go`, when `addType` is empty and template is specified, extract type from the rendered template frontmatter

## Reusable Pattern
When adding a new field to task frontmatter templates, ensure the CLI `add` command applies it as a default — otherwise the `.md` file and `index.json` will have different values for the same field.

## References
- `forge-cli/pkg/template/template.go:20-39` — Defaults struct missing Type field
- `forge-cli/internal/cmd/task/add.go:164-186` — defaults applied without Type
- `forge-cli/pkg/task/category.go:20-40` — CategoryForType logs warning on empty type
