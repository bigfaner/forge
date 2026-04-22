---
name: design-tech
description: Use after PRD is finalized to create technical design with architecture and implementation details.
---

# Design Tech

## Overview

从 PRD 产出技术设计文档，结合项目现状进行技术决策。

**核心原则**：在设计阶段解决技术不确定性，避免实现时的返工。

<HARD-GATE>
Do NOT write any implementation code until tech-design.md is approved. The output of this skill is a design document, not code.
</HARD-GATE>

## Prerequisites

检查上一阶段产物，缺失则中止并提示用户：

```bash
ls docs/features/<slug>/prd/prd-spec.md
```

| 产物 | 缺失时提示 |
|------|-----------|
| `prd/prd-spec.md` | 先执行 `/write-prd`，再执行 `/eval-prd` |

## When to Use

**Trigger conditions:**

- Manifest exists at `docs/features/<slug>/manifest.md` with status `prd`
- PRD Spec exists at `prd/prd-spec.md`
- PRD is approved and ready for technical design

**Skip when:**

- No PRD exists (use `/write-prd` first)
- Design already exists for the feature

## Process Flow

```
1. Read PRD → 2. Explore context → 3. Identify decisions → 4. Ask questions → 5. Draft design → 6. Review → 7. Finalize
```

## Step 1: Read Manifest → PRD

1. Read `manifest.md` to locate documents
2. Read `prd/prd-spec.md`:
   - Understand requirements
   - Note non-functional requirements — these are the **technical constraints** that drive your decisions
   - Identify acceptance criteria
3. Read `prd/prd-user-stories.md` — extract all Given/When/Then acceptance criteria into a checklist
   - Keep this AC list visible throughout the design process — every AC must map to a design element

> **Note**: PRD 故意不含技术选型（brainstorm 和 write-prd 阶段禁止引入）。所有技术决策从本阶段开始。用 PRD 中的非功能性约束作为技术选型的输入条件。

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

### 5.1 PRD Coverage Verification

After drafting each section, verify every PRD acceptance criterion is addressed:

1. For each AC from `prd-user-stories.md`, identify which interface, model, or component handles it
2. If an AC has no corresponding design element, add one
3. Document the mapping in the "PRD Coverage Map" section of the template

### 5.2 Breakdown-Readiness Check

Before seeking approval, verify the design can be directly decomposed into implementation tasks:

| Check | Requirement |
|-------|-------------|
| Components enumerable | Can you list and count all components/modules by name? |
| Interfaces → tasks | Does each interface map to at least one implementation task? |
| Models → tasks | Does each data model map to at least one schema/migration task? |
| PRD AC coverage | Are all acceptance criteria from user stories addressed? |

If any check fails, add the missing detail before presenting to the user.

## Step 6: Get Approval

For each section, wait for user approval.

## Step 7: Write Design Documents

Save to:
- `docs/features/<slug>/design/tech-design.md` — using `templates/tech-design.md`
- `docs/features/<slug>/design/api-handbook.md` — using `templates/api-handbook.md` (if feature has API surface)

## Step 8: Update Manifest

Update `manifest.md`:
- Add Tech Design and API Handbook rows to Documents table
- Add traceability links from PRD sections to design sections
- Advance status to `design` if `/ui-design` already completed or if UI is not applicable

## Step 9: Adversarial Eval Prompt

After committing, use `AskUserQuestion` to ask:

> 是否运行 `/eval-design` 对技术设计进行对抗性评估？（默认 80 分 / 3 轮）

- **Yes** → invoke `/eval-design` via `Skill` tool
- **Custom** → invoke `/eval-design --target X --iterations Y` via `Skill` tool
- **No** → proceed to `/breakdown-tasks`

## Integration

Works well with skills:

- `/write-prd` - Creates PRD input and manifest
- `/ui-design` - Parallel skill for UI features
- `/eval-design` - Evaluate tech-design.md quality before handing off to breakdown-tasks
- `/breakdown-tasks` - Uses tech-design.md to create tasks
