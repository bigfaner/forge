---
id: "2"
title: "Template engine infrastructure and coding template"
priority: "P0"
estimated_time: "1.5h"
dependencies: ["1"]
scope: "backend"
breaking: false
type: "coding.feature"
mainSession: false
---

# 2: Template engine infrastructure and coding template

## Description

Introduce `text/template` + `//go:embed` for record generation, replacing the string-concatenation in `fillRecordTemplate()`. Create the first template file (`record-coding.md`) that produces byte-identical output to the current format.

This is the infrastructure migration step. After this task, the rendering pipeline is template-based, and adding new category templates (tasks 3-5) is just adding new `.md` files.

## Reference Files
- `docs/proposals/typed-task-records/proposal.md` — Source proposal
- `forge-cli/internal/cmd/task/submit.go` — `fillRecordTemplate()` (lines 367-448)
- `forge-cli/pkg/prompt/prompt.go` — existing template engine pattern (reference for embed + render)
- `forge-cli/pkg/prompt/data/*.md` — existing embedded templates (reference)

## Acceptance Criteria
- [ ] `record-coding.md` template file created under `forge-cli/pkg/task/data/` (or similar)
- [ ] Template embedded via `//go:embed`
- [ ] `fillRecordTemplate()` refactored to: determine category → select template → render with `text/template`
- [ ] Coding task records are **byte-identical** to current output (backward compatible)
- [ ] Template data struct exposes all fields needed by the template (status, times, summary, files, decisions, tests, coverage, criteria, notes, reclassification)
- [ ] Helper functions (`formatList`, `formatCoverage`, `formatTestsExecuted`, `formatCriteria`, `formatDuration`) available in template context via `FuncMap`
- [ ] Unit test: golden-file comparison of coding template output vs old `fillRecordTemplate` output

## Hard Rules
- `fillRecordTemplate` signature unchanged: `(t *task.Task, rd *task.RecordData, startedTime string) string`
- Follow the `pkg/prompt/prompt.go` pattern for template loading (embed + parse + execute)
- Template files live under `forge-cli/pkg/task/data/` — same directory convention as `prompt/data/`
- Must handle the `TypeReclassification` conditional block in the template

## Implementation Notes
- The template needs access to task metadata (ID, Title, Type) plus record data plus computed fields (timeSpent, started/completed timestamps). Create a `recordTemplateData` struct that combines all these.
- Coding template covers the current format exactly. Use `{{if .TypeReclassification}}...{{end}}` for the reclass block.
- The template selection logic: `CategoryForType(t.Type)` → map to template file name → render.
