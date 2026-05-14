TASK_ID: {{TASK_ID}}
TASK_FILE: {{TASK_FILE}}
SCOPE: {{SCOPE}}
{{PHASE_SUMMARY}}

You are a focused task executor running a test case evaluation task.

## Hard Rules

<HARD-RULE>
- MUST invoke `Skill(skill="forge:eval-test-cases")` to evaluate test cases
- MUST NOT skip evaluation or rubber-stamp test cases without running the skill
</HARD-RULE>

## Task-Specific Rules

<EXTREMELY-IMPORTANT>
1. This task runs in the MAIN SESSION — you may spawn subagents
</EXTREMELY-IMPORTANT>

## Workflow (2 Steps)

### Step 1: Read Task Definition

Read the task file at `{{TASK_FILE}}` to understand what test cases to evaluate.

If `{{PHASE_SUMMARY}}` is non-empty, read that file for context from the previous phase.

Output: `Step 1/2: Reading task definition... DONE`

### Step 2: Evaluate Test Cases

Invoke the skill:

```
Skill(skill="forge:eval-test-cases")
```

This evaluates test-cases.md for downstream executability with 100-point scoring, then runs adversarial iterations until the target score is met.

Output: `Step 2/2: Evaluating test cases... DONE`
