---
name: design-tech
description: Use after PRD is finalized to create technical design with architecture and implementation details.
---

# Design Tech

## Overview

从 PRD 产出技术设计文档，结合项目现状进行技术决策。

**核心原则**：在设计阶段解决技术不确定性，避免实现时的返工。

<HARD-GATE>
Do NOT write any implementation code until design.md is approved. The output of this skill is a design document, not code.
</HARD-GATE>

## Position in Workflow

```
/write-prd → /design-tech → /eval-design → /breakdown-tasks
     ↓              ↓              ↓               ↓
   prd.md      design.md    design-eval.md     tasks/*.md
```

## When to Use

**Trigger conditions:**

- PRD document exists at `docs/features/<slug>/prd.md`
- PRD is approved and ready for technical design

**Skip when:**

- No PRD exists (use `/write-prd` first)
- Design already exists for the feature

## Process Flow

```
1. Read PRD → 2. Explore context → 3. Identify decisions → 4. Ask questions → 5. Draft design → 6. Review → 7. Finalize
```

## Step 1: Read PRD

Read `docs/features/<slug>/prd.md`:

- Understand requirements
- Note non-functional requirements
- Identify acceptance criteria

## Step 2: Explore Context

| Source                 | What to Look For                                  |
| ---------------------- | ------------------------------------------------- |
| `docs/ARCHITECTURE.md` | Layer constraints                                 |
| `docs/DECISIONS.md`    | Existing decisions                                |
| Package manager files  | Current dependencies (package.json, go.mod, etc.) |
| Source directories     | Existing patterns (src/, internal/, lib/, etc.)   |

## Step 3: Identify Decisions

| Decision Type  | Example Questions        |
| -------------- | ------------------------ |
| Architecture   | Where does this fit?     |
| Interface      | What interfaces needed?  |
| Data Model     | What structures needed?  |
| Dependencies   | New dependencies?        |
| Error Handling | How to handle errors?    |
| Testing        | Test strategy?           |
| Security       | Security considerations? |

## Step 4: Ask Questions

Use `AskUserQuestion` for ALL uncertain areas.

## Step 5: Draft Design

Present incrementally, section by section:

| Section        | Content                 |
| -------------- | ----------------------- |
| Overview       | High-level approach     |
| Architecture   | Component diagram       |
| Interfaces     | Interface definitions   |
| Data Models    | Struct definitions      |
| Error Handling | Error strategy          |
| Testing        | Test strategy           |
| Security       | Security considerations |

## Step 6: Get Approval

For each section, wait for user approval.

## Step 7: Write design.md

Save to `docs/features/<slug>/design.md`

## Integration

Works well with skills:

- `/write-prd` - Creates PRD input
- `/eval-design` - Evaluate design.md quality before handing off to breakdown-tasks
- `/breakdown-tasks` - Uses design.md to create tasks
