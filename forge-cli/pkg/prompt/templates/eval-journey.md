---
type: eval.journey
category: eval
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
{{if .SurfaceKey}}SURFACE_KEY: {{.SurfaceKey}}{{end}}
{{if .PhaseSummary}}{{.PhaseSummary}}{{end}}

You are a focused task executor running a quality evaluation task.

## Task Constraints

<TASK-CONSTRAINTS>
- MUST invoke `Skill(skill="forge:eval", args="--type journey --target 850")` to evaluate quality
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
Skill(skill="forge:eval", args="--type journey --target 850")
```

This runs a quality evaluation using rubric scoring against the target threshold. The eval skill handles scoring, findings collection, and severity assessment internally.

## Record Fields

When submitting via `forge:submit-task`, populate these record fields in record.json:
- **score**: eval score (0-1000)
- **findings**: list of issues found during evaluation
- **severity**: overall severity level (critical/major/minor)
- **passed**: whether the evaluation passed the quality gate

Output: `Step 2/2: Running evaluation... DONE`
