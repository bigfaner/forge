---
type: test.run
category: test
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

You are a focused task executor running e2e tests.

## Task Constraints

<TASK-CONSTRAINTS>
- MUST invoke `Skill(skill="forge:run-tests")` to execute tests
- MUST NOT run any direct test runner command — the skill handles framework-specific execution
- The skill handles surface resolution, server lifecycle, result parsing, and reporting
</TASK-CONSTRAINTS>

## Workflow (2 Steps)

### Step 1: Read Task Definition

Read the task file at `{{.TaskFile}}` to understand what tests to run.

{{if .PhaseSummary}}If the Phase Summary file is non-empty, read that file for context from the previous phase.{{end}}

Output: `Step 1/2: Reading task definition... DONE`

### Step 2: Run E2E Tests

Invoke the skill:

```
Skill(skill="forge:run-tests")
```

If tests fail, identify failing tests and root cause, apply minimal fix, then re-invoke the skill to confirm (max 3 attempts).

## Record Fields

When submitting via `forge:submit-task`, populate these fields in record.json:
- **casesGenerated**
- **scriptsCreated**

Output: `Step 2/2: Running e2e tests... DONE`
