---
name: quick
description: Streamlined pipeline for features (1-10 tasks). Brainstorm -> tasks -> execute, no PRD or design.
allowed-tools: Bash Read Write Edit Grep Glob Agent Skill AskUserQuestion
---

# /quick

Streamlined pipeline for small features: brainstorm → tasks → execute.

## Architecture

```mermaid
flowchart TD
    INPUT(["User Input<br><i>any format: question, idea, feature request</i>"]) --> A
    A["1. /brainstorm<br>proposal.md"] --> B{"2. User confirms ✋"}
    B -->|"Yes"| C["3. /quick-tasks<br>tasks + index.json + manifest.md"]
    B -->|"Revise"| A
    B -->|"Abort"| STOP["Stop"]
    C --> D["4. /run-tasks<br>auto-execute"]
    C -->|">10 tasks"| FULL["STOP → recommend<br>full pipeline"]
```

## Core Rules

<EXTREMELY-IMPORTANT>
1. **Execute the pipeline in order.** Always start with Step 1 (brainstorm), regardless of how the user's input is phrased. Question-like inputs ("can we simplify X?", "should we refactor Y?") are NOT discussions — they are feature requests that brainstorm exists to shape into structured proposals. Do NOT substitute ad-hoc analysis (Explore agents, grep, file reads) for the pipeline.
2. Maximum 10 business tasks. If brainstorm produces a proposal that needs more, STOP and suggest the full pipeline.
3. ONE feature per invocation.
4. The /quick pipeline is for small, well-scoped features. If scope grows during brainstorm, recommend switching to full mode.
</EXTREMELY-IMPORTANT>

## Step 1: Brainstorm

Invoke the brainstorm skill:

```
Skill(skill="forge:brainstorm")
```

This produces `docs/proposals/<slug>/proposal.md` through interactive dialogue with the user. The brainstorm skill handles all user interaction and commits the proposal.

After brainstorm completes, extract the feature slug from the proposal directory path.

## Step 2: User Confirmation

Read the generated `docs/proposals/<slug>/proposal.md` and present a summary:

```
## Quick Mode: Proposal Summary

**Problem**: <one line from proposal>
**Solution**: <one line from proposal>
**Scope**:
- <In Scope bullets>
**Success Criteria**:
- <Success Criteria checkboxes>

Generate tasks from this proposal?
```

Use `AskUserQuestion` with three options:

| Option | Action |
|--------|--------|
| **Yes, generate tasks** | Update proposal status, then proceed to Step 3 |
| **Revise proposal** | Return to Step 1 (re-run brainstorm) |
| **Abort** | Stop cleanly |

<EXTREMELY-IMPORTANT>
This confirmation is MANDATORY. The proposal is the sole input for the entire quick mode pipeline — no PRD or design will be created to correct course. A wrong direction here means all downstream tasks are wasted.
</EXTREMELY-IMPORTANT>

### Status Transition: Draft → Approved

When the user selects **"Yes, generate tasks"**, update the proposal frontmatter status:

```
Edit(file_path="docs/proposals/<slug>/proposal.md",
     old_string="status: Draft",
     new_string="status: Approved")
```

This must be an atomic frontmatter edit targeting only the `status:` line. Do NOT rewrite the entire file.

## Step 3: Generate Tasks

Invoke the quick-tasks skill:

```
Skill(skill="forge:quick-tasks")
```

This produces:
- `docs/features/<slug>/tasks/*.md` — task files (1-10 business + optional T-quick-1~4, T-quick-specs-1)
- `docs/features/<slug>/tasks/index.json` — task index (compatible with `/run-tasks`)
- `docs/features/<slug>/manifest.md` — simplified manifest

If quick-tasks reports >10 tasks needed, STOP and recommend the full pipeline:

```
"This feature requires more than 10 tasks — too large for quick mode.
Recommend using the full pipeline: /write-prd → /tech-design → /breakdown-tasks"
```

## Step 3→4 Transition

<EXTREMELY-IMPORTANT>
After quick-tasks completes — including its Step 8 commit — you MUST **immediately proceed** to Step 4 (run-tasks) with **zero intermediate output**. Specifically:

1. **Do NOT** output any summary, recap, or status message between quick-tasks and run-tasks.
2. **Do NOT** pause for user confirmation. The user already confirmed in Step 2.
3. **Do NOT** ask the user anything. Invoke run-tasks directly.
4. **Do NOT** perform any intermediate actions (file reads, git status, exploratory analysis) between the two skills.

If quick-tasks Step 8 commit fails, stop and fix the issue before proceeding. Only proceed to run-tasks after the commit succeeds.
</EXTREMELY-IMPORTANT>

## Step 4: Execute Tasks

Invoke the run-tasks command to auto-execute all tasks:

```
Skill(skill="forge:run-tasks")
```

The existing run-tasks dispatcher will:
1. Read `index.json`
2. Claim tasks in dependency order
3. Dispatch to task-executor subagents
4. Run breaking gates (compile + fmt + lint + test)
5. Handle fix tasks on failure
6. Run all-completed hook as final safety net

## Error Handling

| Situation | Action |
|-----------|--------|
| Brainstorm fails | Stop, user can retry |
| User aborts at confirmation | Stop cleanly |
| quick-tasks exceeds 10 task limit | Stop, recommend full pipeline |
| `forge task validate-index` fails | Stop, fix index.json issues |
| run-tasks encounters failures | Handled by dispatcher (fix tasks, retries) |
