---
name: write-prd
description: Use when user provides requirements or feature requests that need to be formalized into a structured PRD document through collaborative dialogue.
---

# Write PRD

## Overview

从模糊需求产出清晰的 PRD（产品需求文档），通过协作对话逐步澄清需求。

**核心原则**：在编码前先澄清 "做什么" 和 "为什么做"，避免方向性错误。

<HARD-GATE>
Do NOT write any code, scaffold any project, or take any implementation action until the PRD is finalized and approved. Present the PRD and get user approval first.
</HARD-GATE>

## When to Use

**Trigger conditions:**
- User describes a feature/requirement without clear specifications
- User says "I want to..." or "We need..." without details
- Starting a new phase or major feature

**Skip when:**
- Clear task definitions already exist
- Simple bug fix or small tweak

## Process Flow

```
Explore context → Assess scope → Ask questions → Propose approaches → Present PRD → Write doc → Commit
```

## Checklist

1. **Explore project context** — check files, docs, recent commits
2. **Assess scope** — determine if request needs decomposition
3. **Ask clarifying questions** — one at a time via AskUserQuestion tool
4. **Propose 2-3 approaches** — with trade-offs and your recommendation
5. **Present PRD sections** — get approval after each section
6. **Write PRD document** — save to `docs/features/<feature-slug>/prd.md`
7. **Commit** — commit the PRD document

## Step 1: Explore Project Context

Before asking questions, understand the current state:
- Read `docs/ARCHITECTURE.md` for architecture constraints
- Read `docs/DECISIONS.md` for existing technical decisions
- Check `docs/features/<slug>/tasks/index.json` for related tasks
- Review recent git commits

## Step 2: Assess Scope

Evaluate if the request is appropriately scoped:
- If request describes multiple independent subsystems → **Decompose first**
- If single focused feature → **Proceed with questions**

## Step 3: Ask Clarifying Questions

**CRITICAL**: Use `AskUserQuestion` tool for ALL questions.

### Question Guidelines

- **One question at a time** — never batch questions
- **Prefer multiple choice** — easier to answer than open-ended
- **Focus on understanding**: purpose, constraints, success criteria
- **Go back when needed** — if something doesn't make sense, clarify

## Step 4: Propose Approaches

After understanding requirements, propose 2-3 implementation approaches:
1. **Present options conversationally** with your recommendation
2. **Lead with your recommended option** and explain why
3. **Include trade-offs** for each approach

## Step 5: Present PRD Sections

Present incrementally, getting approval after each section:

| Section | Content | When to Present |
|---------|---------|-----------------|
| Background | Problem statement, context | First |
| Goals | Primary goals, success metrics | After background approved |
| Scope | In/out of scope items | After goals approved |
| Requirements | Functional requirements | After scope approved |
| Acceptance Criteria | Testable conditions | Last |

## Step 6: Write PRD Document

**Directory structure:**
```
docs/features/<feature-slug>/
├── prd.md           # The PRD document
├── tasks/           # Task definitions (created by /breakdown-tasks)
└── records/         # Execution records (created by /record-task)
```

## Integration

Works well with:
- `/breakdown-tasks` - After PRD is finalized, break into tasks
- `docs/DECISIONS.md` - Record key decisions during PRD creation
