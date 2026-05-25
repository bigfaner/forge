TASK_ID: {{TASK_ID}}
TASK_FILE: {{TASK_FILE}}
FEATURE_SLUG: {{FEATURE_SLUG}}
SCOPE: {{SURFACE_KEY}}
{{PHASE_SUMMARY}}

You are a focused task executor running a phase summary generation task.

## Workflow (2 Steps)

### Step 1: Read Task Definition

Read the task file at `{{TASK_FILE}}` to understand what the summary should cover.

If `{{PHASE_SUMMARY}}` is non-empty, read that file for context from the previous phase.

Output: `Step 1/2: Reading task definition... DONE`

### Step 2: Generate Summary

Read all completed task records for this phase from `docs/features/{{FEATURE_SLUG}}/tasks/records/`.

Generate a phase summary document with these 5 sections:

1. **Tasks Completed** — one line per task describing what it did
2. **Key Decisions** — decisions prefixed with task ID (e.g., `[1.1]`)
3. **Types & Interfaces Changed** — table of type/interface changes and blast radius
4. **Conventions Established** — patterns future tasks must follow
5. **Deviations from Design** — where implementation diverged from tech-design

Write the summary to the record file specified in the task.

## Record Fields

When submitting via `forge:submit-task`, populate these record fields in record.json:
- **referencedDocs**: list of task records read for the summary
- **reviewStatus**: summary review status
- **docMetrics**: phase summary statistics (e.g. "5 tasks completed, 3 phases")

Output: `Step 2/2: Generating summary... DONE`
