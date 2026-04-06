---
name: execute-task
description: Execute single task with focused TDD workflow.
allowed_tools: ["Bash", "Read", "Write", "Edit", "Grep", "Glob", "Agent", "LSP"]
---

# /execute-task

Execute a single task with streamlined TDD workflow.

## Workflow (5 Steps)

```
Step 1: Read task definition
Step 2: TDD (RED → GREEN → REFACTOR)
Step 3: Full verification
Step 4: Record task (MANDATORY)
Step 5: Git commit
```

## Step 1: Claim & Read

```bash
task claim
```

Parse output for KEY, ID, FILE. Read task file.

## Step 2: TDD Implementation

```
RED      → Write failing test first
GREEN    → Implement minimal code to pass
REFACTOR → Clean up while keeping tests green
```

## Step 3: Full Verification

Run project-specific verification commands.

## Step 4: Record Task (MANDATORY)

```
Skill(skill="record-task")
```

## Step 5: Commit

```
Skill(skill="git-commit")
```

## Rules

- **record-task is mandatory** - No completion without it
- **All verifications must pass**
- **Commit only after record**

## Related Commands

| Command | Usage |
|---------|-------|
| `/run-tasks` | Auto-execute all tasks |
| `/claim-task` | Claim task only |
| `/record-task` | Create record + update status |
