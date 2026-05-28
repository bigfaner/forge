---
type: doc.consolidate
category: doc
variables:
  - TaskID
  - TaskFile
  - FeatureSlug
  - PhaseSummary
  - SurfaceKey
  - SurfaceType
---
TASK_ID: {{.TaskID}}
TASK_FILE: {{.TaskFile}}
SURFACE_KEY: {{.SurfaceKey}}
{{if .PhaseSummary}}{{.PhaseSummary}}{{end}}

You are a focused task executor running a spec consolidation task in non-interactive (pipeline) mode. You are running under `/run-tasks` dispatcher — no user is present. The consolidate-specs skill will auto-integrate all CROSS items and commit with `[auto-specs]` tag. Do NOT wait for user confirmation. Proceed without stopping.

## Workflow (2 Steps)

### Step 1: Read Task Definition

Read the task file at `{{.TaskFile}}` to understand what specs to consolidate.

{{if .PhaseSummary}}If the Phase Summary file is non-empty, read that file for context from the previous phase.{{end}}

Output: `Step 1/2: Reading task definition... DONE`

### Step 2: Consolidate Specs

Invoke the skill:

```
Skill(skill="forge:consolidate-specs")
```

This extracts business rules and tech specs from feature docs into preview files, detects overlaps with existing knowledge, auto-integrates all CROSS items in non-interactive mode, and commits changes with `[auto-specs]` tag. All-LOCAL items auto-proceed without integration.

## Record Fields

When submitting via `forge:submit-task`, populate these record fields in record.json:
- **referencedDocs**: list of spec files consolidated
- **reviewStatus**: consolidation outcome (e.g. "completed")
- **docMetrics**: consolidation statistics (e.g. "3 rules extracted, 1 overlap found")

Output: `Step 2/2: Consolidating specs... DONE`
