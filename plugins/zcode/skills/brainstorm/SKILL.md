---
name: brainstorm
description: Use when a user has a vague idea or feature request and needs to explore it before formalizing into a PRD. Outputs a structured proposal document.
---

# Brainstorm

## Overview

从模糊想法到结构化提案，通过协作对话探索问题空间。

**核心原则**：在投入 PRD 之前，先确认问题值得解决、方案值得投入。

<HARD-GATE>
Do NOT write any code or take implementation action. This skill produces a proposal document only.
</HARD-GATE>

## When to Use

**Trigger conditions:**
- User describes an idea without clear specifications
- User says "I'm thinking about..." or "What if we..."
- Starting exploration before committing to a feature

**Skip when:**
- Requirements are already clear (use `/write-prd` directly)
- Bug fix or small tweak

## Process Flow

```
Understand idea → Explore context → Challenge assumptions → Define scope → Write proposal → Commit
```

## Step 1: Understand the Idea

Listen actively. Ask clarifying questions one at a time via `AskUserQuestion`:
- What problem does this solve?
- Who is affected?
- What does success look like?

## Step 2: Explore Context

| Source | What to Look For |
|--------|-----------------|
| Existing features | Is this already solved elsewhere? |
| Recent commits | Related recent changes |
| Project docs | Architecture constraints, existing decisions |

## Step 3: Challenge Assumptions

Play devil's advocate:
- Is this the right problem to solve?
- Are there simpler alternatives?
- What if we did nothing?

## Step 4: Define Scope

Propose in-scope and out-of-scope boundaries. Get user agreement.

## Step 5: Write Proposal

Save to `docs/proposals/<slug>/proposal.md` using `templates/proposal.md`.

## Step 6: Commit

```bash
git add docs/proposals/<slug>/
git commit -m "docs: add proposal for <feature-slug>"
```

## Integration

Works well with:
- `/write-prd` — Takes proposal as optional input to formalize into PRD
