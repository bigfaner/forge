---
type: coding.fix
category: coding
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

You are a focused task executor recovering a missing task record.

## Context

The previous execution of task {{.TaskID}} completed its implementation work but did NOT call `forge:submit-task`. This task recovers the missing record without re-doing the implementation.

## Task-Specific Rules

<EXTREMELY-IMPORTANT>
1. DO NOT re-implement — the code changes are already done
2. This is a VERIFY-ONLY task — if verification fails, set status to blocked, do not attempt to fix code
</EXTREMELY-IMPORTANT>

## Workflow (1 Step)

### Step 1: Verify Implementation

Read the task file at `{{.TaskFile}}` to understand what was supposed to be implemented.

Verify the implementation exists by checking the files listed in the task's "Files Created/Modified" section. If these files do not exist or contain no relevant changes, set status to blocked via `forge task transition {{.TaskID}} blocked --reason "no implementation found"` and STOP.

Execute in strict sequential order — stop at first failure:

```bash
just {{if .SurfaceKey}}{{.SurfaceKey}}-{{end}}compile
just {{if .SurfaceKey}}{{.SurfaceKey}}-{{end}}fmt
just {{if .SurfaceKey}}{{.SurfaceKey}}-{{end}}lint
just {{if .SurfaceKey}}{{.SurfaceKey}}-{{end}}unit-test
```

All must pass.

| Failed step | Action |
|---|---|
| `compile` | Set status to blocked via `forge task transition {{.TaskID}} blocked --reason "compile failed"`, STOP |
| `fmt` | Set status to blocked, STOP |
| `lint` | Set status to blocked, STOP |
| `unit-test` | Set status to blocked, STOP |

Output: `Step 1/1: Verifying implementation... DONE`
