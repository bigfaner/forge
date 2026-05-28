---
type: test.gen-scripts
category: test
variables:
  - TaskID
  - TaskFile
  - FeatureSlug
  - PhaseSummary
  - TestTypeArg
  - SurfaceKey
  - SurfaceType
---
TASK_ID: {{.TaskID}}
TASK_FILE: {{.TaskFile}}
{{if .SurfaceKey}}SURFACE_KEY: {{.SurfaceKey}}{{end}}
{{if .PhaseSummary}}{{.PhaseSummary}}{{end}}

You are a focused task executor generating test scripts.

## Task Constraints

<TASK-CONSTRAINTS>
- MUST invoke `Skill(skill="forge:gen-test-scripts")` to generate scripts
- MUST NOT write test scripts manually — the skill generates them from test cases
</TASK-CONSTRAINTS>

## Workflow (2 Steps)

### Step 1: Read Task Definition

Read the task file at `{{.TaskFile}}` to understand what test scripts to generate.

{{if .PhaseSummary}}If the Phase Summary file is non-empty, read that file for context from the previous phase.{{end}}

Output: `Step 1/2: Reading task definition... DONE`

### Step 2: Generate Test Scripts

Invoke the skill:

```
Skill(skill="forge:gen-test-scripts"{{.TestTypeArg}})
```

## Record Fields

When submitting via `forge:submit-task`, populate these fields in record.json:
- **scriptsCreated**
- **casesGenerated**

Output: `Step 2/2: Generating test scripts... DONE`
