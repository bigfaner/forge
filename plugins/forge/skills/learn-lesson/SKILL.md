---
name: learn-lesson
description: Use when you have solved an error or discovered a useful pattern. Extracts reusable knowledge from the current session.
---

# Learn Lesson

## Overview

Extract reusable knowledge from the current session and record it to `docs/lessons/`.

**Core principle**: Record "what to do next time you encounter a similar problem", not "what I did".

## When to Use

**Trigger conditions:**
- Solved a non-trivial error with reusable insights
- Discovered a pattern/technique worth remembering
- User explicitly requests `/learn-lesson`

**Skip when:**
- Trivial typo fixes
- One-off environment issues
- Information already documented elsewhere

## Workflow

```
1. Identify lesson → 2. Classify category → 3. Write doc → 4. User review → 5. Commit
```

## Step 1: Identify Lesson

- The problem encountered (symptoms)
- Root cause
- Solution
- **Reusable takeaway**

### Root Cause Investigation (Mandatory)

When an error occurs, **you must dig deep into the root cause** — never stop at surface symptoms. This is especially critical for errors related to Claude Code or agent behavior itself:

- **Don't accept the first explanation**: tool call failures, output truncation, context loss, etc. almost always have a deeper cause
- **Distinguish symptoms from root cause**: `tool call failed` is a symptom; *why* it failed is the root cause
- **Agent-related errors — key areas to investigate**:
  - Claude Code tool permissions / sandbox restrictions
  - Information loss due to context window compression
  - Hook interception or configuration conflicts
  - State desync between sub-agent and main agent
  - Mismatch between model output format and tool expectations
- **Causal chain**: symptom → direct cause → root cause → trigger condition — trace at least 3 levels deep

## Step 2: Classify Category

| Category | Prefix | Example |
|----------|--------|---------|
| Debugging | `debug-` | `debug-race-condition.md` |
| Architecture | `arch-` | `arch-dependency-direction.md` |
| Tooling | `tool-` | `tool-go-test-coverage.md` |
| Pattern | `pattern-` | `pattern-error-wrapping.md` |
| Gotcha | `gotcha-` | `gotcha-context-cancellation.md` |

## Step 3: Write Document

Output: `docs/lessons/<category-prefix><slug>.md`

Template:
```markdown
# <Title>

## Problem
<!-- Symptom description -->

## Root Cause
<!-- Why it happened -->

## Solution
<!-- How to fix it -->

## Key Takeaway
<!-- Reusable insight for next time -->
```

## Step 4: User Review

**Do not commit directly.** Show the generated lesson document and wait for user confirmation:

- Commit only after the user confirms the content is correct
- If the user requests changes, revise and show again

Only execute after explicit user approval:

```bash
git add docs/lessons/<filename>.md
git commit -m "docs(lessons): add <title>"
```

## Common Mistakes

| Mistake | Correction |
|---------|------------|
| Recording "what I did" | Focus on "what to do next time" |
| Too specific | Generalize to reusable pattern |
| Missing root cause | Always include why |
| Stopping at surface symptoms | Trace the causal chain at least 3 levels deep |
