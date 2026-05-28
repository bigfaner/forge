---
type: eval.contract
category: eval
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


You are a focused task executor evaluating quality.

## Task Constraints

<TASK-CONSTRAINTS>
- MUST invoke `Skill(skill="forge:eval", args="--type contract --target 850")` to evaluate quality
- MUST NOT modify the files being evaluated
</TASK-CONSTRAINTS>

## Workflow (2 Steps)

### Step 1: Read Task Definition

Read the task file at `{{.TaskFile}}` to understand what to evaluate.

{{if .PhaseSummary}}If the Phase Summary file is non-empty, read that file for context from the previous phase.{{end}}

Output: `Step 1/2: Reading task definition... DONE`

### Step 2: Run Evaluation

Invoke the skill:

```
Skill(skill="forge:eval", args="--type contract --target 850")
```

## Record Fields

When submitting via `forge:submit-task`, populate these fields in record.json:
- **score**
- **findings**
- **severity**
- **passed**

Output: `Step 2/2: Running evaluation... DONE`
