---
created: "2026-05-15"
tags: [architecture]
---

# Pipeline Skill Bypass: Skipping /quick Steps for Exploratory Input

## Problem

User invoked `/quick` with an exploratory argument ("can we simplify test tasks and gate tasks in quick-tasks and breakdown-tasks?"). Agent bypassed the entire /quick pipeline — no brainstorm, no proposal.md, no quick-tasks — and instead ran ad-hoc code analysis via Explore agents, grep, and file reads. Produced an architectural analysis document instead of following the pipeline.

## Root Cause

**Causal chain (4 levels)**:

1. **Symptom**: Agent launched 2 Explore agents + grep analysis, skipped `Skill(skill="forge:brainstorm")` entirely. No proposal.md created.

2. **Direct cause**: Agent judged the input as a "discussion question" rather than a feature request, and decided ad-hoc analysis was more appropriate than the brainstorm pipeline step.

3. **Root cause**: Insufficient trust in the pipeline skill's design. The brainstorm step exists precisely to convert vague/exploratory inputs into structured proposals. By skipping it, the agent circumvented the skill's purpose — and produced output that doesn't integrate with downstream steps (quick-tasks, run-tasks).

4. **Trigger condition**: When the user's input to a pipeline skill is phrased as a question ("is X possible?", "should we simplify Y?") rather than a declarative feature description, the agent defaults to analysis mode instead of pipeline mode.

## Solution

Follow the pipeline. When a user invokes `/quick` (or any pipeline skill), execute Step 1 as written. The brainstorm skill handles interactive dialogue to shape vague ideas into proposals — that's its core purpose.

For this specific case: the correct flow was brainstorm → proposal.md about "simplify test/gate task descriptions in skill files" → user confirms → quick-tasks generates implementation tasks → run-tasks executes them.

## Reusable Pattern

**Rule**: When a user invokes a pipeline skill (e.g., `/quick`, `/write-prd`, `/tech-design`), execute its steps in order. Do not substitute ad-hoc analysis for the skill's designed workflow.

**Why**: Pipeline skills produce structured outputs (proposal.md, prd-spec.md, tasks/*.md) that downstream steps depend on. Ad-hoc analysis produces knowledge that exists only in conversation context and can't be consumed by subsequent skills.

**How to apply**: Even when the input looks like a question rather than a feature request, start with the pipeline's first step. If the brainstorm reveals that the idea isn't suited for the pipeline (e.g., it's a pure discussion), the skill will naturally surface that and the user can redirect.

**Exception**: If the user explicitly says "just analyze this" or "don't run the pipeline", respect their instruction. But the default is to follow the pipeline.

## Example

```
# Wrong: skip pipeline, do ad-hoc analysis
User: /quick "can we simplify test tasks in skills?"
Agent: *launches Explore agents, greps files, produces analysis*

# Right: follow pipeline
User: /quick "can we simplify test tasks in skills?"
Agent: *invokes /brainstorm*
       → interactive dialogue shapes the idea
       → proposal.md written
       → user confirms
       → /quick-tasks generates tasks
       → /run-tasks executes
```

## Related Files

- `plugins/forge/commands/quick.md` — /quick pipeline definition
- `plugins/forge/skills/brainstorm/SKILL.md` — brainstorm skill (Step 1 of /quick)
- `plugins/forge/skills/quick-tasks/SKILL.md` — quick-tasks skill (Step 3 of /quick)
