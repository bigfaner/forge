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

Run `just build && just test`.

## Step 4: Record Task (MANDATORY)

Invoke the skill (it contains file location details):

```
Skill(skill="record-task")
```

The skill provides complete workflow including:
- File locations (process/record.json vs tmp/)
- JSON format and fields
- CLI command usage

## Step 5: Commit

```
Skill(skill="git-commit")
```

## Rules

<EXTREMELY-IMPORTANT>
- record-task is mandatory - No completion without it
- All verifications must pass
- Commit only after record
- ONE TASK PER INVOCATION — after Step 5, STOP immediately, no exceptions
- FORBIDDEN: run "task claim", read index.json, or start any subsequent task
</EXTREMELY-IMPORTANT>

<HARD-GATE>
Task is NOT complete until record-task CLI command succeeds. Commit is blocked until record exists.
</HARD-GATE>

## STOP

<HARD-RULE>
ONE TASK PER INVOCATION. This is absolute and non-negotiable.

After Step 5, you MUST stop immediately.

<PROHIBITIONS>
- Running `task claim` under any circumstances
- Reading the next task file
- Continuing with any additional work
</PROHIBITIONS>

Output your final summary and STOP.
</HARD-RULE>

## Related Commands

| Command | Usage |
|---------|-------|
| `/run-tasks` | Auto-execute all tasks |
| `/record-task` | Create record + update status |
