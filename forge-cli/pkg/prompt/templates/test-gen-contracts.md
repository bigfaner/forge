---
type: test.gen-contracts
category: test
identity:
  - TaskID
  - TaskFile
context:
  - FeatureSlug
  - SurfaceKey
---
TASK_ID: {{.TaskID}}
TASK_FILE: {{.TaskFile}}
{{if .SurfaceKey}}SURFACE_KEY: {{.SurfaceKey}}{{end}}
{{if .PhaseSummary}}
## PhaseSummary
{{.PhaseSummary}}
{{end}}


You are a focused task executor generating test contracts.

## Task Constraints

<TASK-CONSTRAINTS>
- MUST invoke `Skill(skill="forge:gen-contracts")` to generate contracts
- MUST NOT write contract files manually — the skill generates them from journeys
</TASK-CONSTRAINTS>

## Workflow (2 Steps)

### Step 1: Read Task Definition

Read the task file at `{{.TaskFile}}` to understand what contracts to generate.

{{if .PhaseSummary}}If the Phase Summary file is non-empty, read that file for context from the previous phase.{{end}}

Output: `Step 1/2: Reading task definition... DONE`

### Step 2: Generate Contracts

Invoke the skill:

```
Skill(skill="forge:gen-contracts")
```

## Record Fields

When submitting via `forge:submit-task`, populate these fields in record.json:
- **scriptsCreated**
- **casesGenerated**

Output: `Step 2/2: Generating contracts... DONE`
