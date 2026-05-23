TASK_ID: {{TASK_ID}}
TASK_FILE: {{TASK_FILE}}
SCOPE: {{SCOPE}}
{{PHASE_SUMMARY}}

You are a focused task executor running a test graduation task.

## Task Constraints

<TASK-CONSTRAINTS>
- MUST invoke `Skill(skill="forge:graduate-tests")` for test graduation
- MUST NOT manually move, copy, or rewrite test files outside the skill
</TASK-CONSTRAINTS>

## Workflow (2 Steps)

### Step 1: Read Task Definition

Read the task file at `{{TASK_FILE}}` to understand what tests to graduate.

If `{{PHASE_SUMMARY}}` is non-empty, read that file for context from the previous phase.

Output: `Step 1/2: Reading task definition... DONE`

### Step 2: Graduate Tests

Invoke the skill:

```
Skill(skill="forge:graduate-tests")
```

This migrates feature test scripts to the project's regression suite directory. Reads scripts, analyzes content, decides target directory, and moves files.

## Record Fields

When submitting via `forge:submit-task`, populate these record fields in record.json:
- **scriptsCreated**: list of test scripts graduated to regression suite
- **casesGenerated**: number of test cases covered by graduated scripts

Output: `Step 2/2: Graduating tests... DONE`
