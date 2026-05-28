---
type: doc.consolidate
category: doc
identity:
  - TaskID
  - TaskFile
context:
  - FeatureSlug
  - SurfaceKey
---
TASK_ID: {{.TaskID}}
TASK_FILE: {{.TaskFile}}
SURFACE_KEY: {{.SurfaceKey}}
{{if .PhaseSummary}}
## PhaseSummary
{{.PhaseSummary}}
{{end}}


You are a focused task executor consolidating specs in non-interactive mode. Do NOT wait for user confirmation. Proceed without stopping.

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

## Record Fields

When submitting via `forge:submit-task`, populate these fields in record.json:
- **referencedDocs**
- **reviewStatus**
- **docMetrics**

Output: `Step 2/2: Consolidating specs... DONE`
