---
id: "5"
title: "Migrate all template frontmatters + PhaseSummary"
priority: "P1"
estimated_time: "2h"
complexity: "medium"
dependencies: [4]
surface-key: "."
surface-type: "cli"
breaking: false
type: "coding.enhancement"
mainSession: false
---

# 5: Migrate all template frontmatters + PhaseSummary

## Description
Batch migrate all 41 template files from flat `variables` list to grouped frontmatter format. Simultaneously migrate PhaseSummary from frontmatter variable to independent body section in all 21 prompt templates. Make a small controlled change to prompt.go for phaseSummaryLine format rendering.

Template type distribution:
- 21 prompt templates: identity + context + conditional (conditional only for templates with `{{if .X}}` blocks) groups
- 14 task templates: identity + context groups
- 6 record templates: identity group

## Reference Files
- forge-cli/pkg/prompt/templates/*.md: 21 prompt templates — add groups + migrate PhaseSummary to body (source: proposal.md#Key-Scenarios-6,7)
- forge-cli/pkg/task/templates/*.md: 14 task templates — add identity + context groups (source: proposal.md#Key-Scenarios-8)
- forge-cli/pkg/task/records/*.md: 6 record templates — add identity group (source: proposal.md#Key-Scenarios-9)
- forge-cli/pkg/prompt/prompt.go: phaseSummaryLine format — remove PHASE_SUMMARY: prefix only (source: proposal.md#In-Scope-Frontmatter-重构)

## Acceptance Criteria
- [ ] SC-FM-1: All 41 templates migrated — prompt templates have identity+context (conditional where applicable), task templates have identity+context, record templates have identity (with required fields per type)
- [ ] SC-FM-4: PhaseSummary removed from frontmatter (grep frontmatter section returns 0 matches), added as `## PhaseSummary` body section in 21 prompt templates wrapped in `{{if .PhaseSummary}}...{{end}}`
- [ ] SC-FM-5: All grouped field names use PascalCase matching Go struct names (e.g., TaskID, SurfaceKey)
- [ ] SC-FM-6: Rendered frontmatter (second `---` block in task/record templates) unchanged — git diff confirms zero modification
- [ ] prompt.go phaseSummaryLine updated: PHASE_SUMMARY: prefix removed, body section format preserved

## Hard Rules
- **forge task index compatibility**: Rendered frontmatter (second `---` block) is parsed by `task/frontmatter.go` for `FrontmatterData` — must not change. The two frontmatter parsing paths are independent: metadata frontmatter by `prompt/metadata.go`, rendered frontmatter by `task/frontmatter.go`
- **Distribution path**: Prompt templates at `forge-cli/pkg/prompt/templates/`, task templates at `forge-cli/pkg/task/templates/`, record templates at `forge-cli/pkg/task/records/`

## Implementation Notes
- PhaseSummary body section position: immediately after TASK_ID/TASK_FILE/SURFACE_KEY lines, before role description/CODING_PRINCIPLES
- PhaseSummary condition block format:
  ```
  {{if .PhaseSummary}}
  ## PhaseSummary
  {{.PhaseSummary}}
  {{end}}
  ```
- Use the proposal's grouped field assignment tables (SC-FM-1 section) for each template type as the source of truth for which fields go in which group
- prompt.go change is minimal: only `phaseSummaryLine` formatting, not Synthesize() main logic
