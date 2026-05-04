---
name: tech-design
description: Use after PRD is finalized to create technical design with architecture and implementation details.
---

# Tech Design

## Overview

Produce technical design from PRD, making technology decisions informed by the current project state.

**Core principle**: Resolve technical uncertainty during the design phase, avoiding rework during implementation.

<HARD-GATE>
Do NOT write any implementation code until tech-design.md is approved. The output of this skill is a design document, not code.
</HARD-GATE>

## Prerequisites

Check previous stage artifacts. Abort and prompt user if missing:

```bash
ls docs/features/<slug>/prd/prd-spec.md
```

| Artifact | Missing prompt |
|----------|----------------|
| `prd/prd-spec.md` | Run `/write-prd` first, then `/eval-prd` |

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
1. Read PRD → 2. Explore context → 3. Identify decisions → 4. Ask questions → 5. Draft design → 6. Review → 7. Archive decisions (optional) → 8. Finalize
```

## Step 1: Read Manifest → PRD

1. Read `manifest.md` to locate documents
2. Read `prd/prd-spec.md`:
   - Understand requirements
   - Note non-functional requirements — these are the **technical constraints** that drive your decisions
   - Identify acceptance criteria
3. Read `prd/prd-user-stories.md` — extract all Given/When/Then acceptance criteria into a checklist
   - Keep this AC list visible throughout the design process — every AC must map to a design element

> **Note**: The PRD intentionally excludes technology selection (brainstorm and write-prd phases forbid it). All technology decisions start from this phase. Use non-functional constraints from the PRD as input conditions for technology selection.

## Step 2: Explore Context

| Source                 | What to Look For                                  |
| ---------------------- | ------------------------------------------------- |
| `docs/ARCHITECTURE.md` | Layer constraints                                 |
| `docs/decisions/`      | Existing decisions (category-based directory)     |
| `docs/business-rules/` | Cross-feature business rules from prior features  |
| `docs/conventions/`    | Technical conventions from prior features         |
| Package manager files  | Current dependencies (package.json, go.mod, etc.) |
| Source directories     | Existing patterns (src/, internal/, lib/, etc.)   |

## Step 3: Identify Decisions

| Decision Type          | Example Questions        |
| ---------------------- | ------------------------ |
| Architecture           | Where does this fit?     |
| Interface              | What interfaces needed?  |
| Data Model             | What structures needed?  |
| Dependencies           | New dependencies?        |
| Error Handling         | How to handle errors?    |
| Testing                | Test strategy?           |
| Security               | Security considerations? |
| Local Dev & Deployment | Dev environment setup?   |

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
| Cross-layer consistency | If feature spans layers, does the Data Map cover every field that crosses boundaries? |

If any check fails, add the missing detail before presenting to the user.

### 5.3 Cross-Layer Data Map

If the feature touches more than one architectural layer (database, API, UI, CLI, etc.):
- Complete the "Cross-Layer Data Map" table in the template
- Every field that appears in multiple layers must have a row showing its type/shape at each layer
- This becomes the Ground Truth for type decisions during task execution

If the feature is single-layer (e.g., only affects CLI output formatting):
- Write "Single-layer feature. Cross-Layer Data Map not applicable." in the section

## Step 6: Get Approval

For each section, wait for user approval.

## Step 7: Archive Decisions (Optional)

Triggered automatically after user approves the tech-design in Step 6.

Follow the tech-design archiving flow defined in `plugins/forge/references/shared/decision-logging.md` (Section 2).

- If the approved document contains key decisions, display the candidate list and prompt the user to select which to archive.
- User may enter `none` to skip archiving entirely.
- If no key decisions exist in the document, skip this step silently.

## Step 8: Write Design Documents

Save to:
- `docs/features/<slug>/design/tech-design.md` — using `templates/tech-design.md`
- `docs/features/<slug>/design/api-handbook.md` — using `templates/api-handbook.md` (if feature has API surface)

## Step 9: Update Manifest

Update `manifest.md`:
- Add Tech Design and API Handbook rows to Documents table
- Add traceability links from PRD sections to design sections
- Advance status to `design` if `/ui-design` already completed or if UI is not applicable

## Step 10: Adversarial Eval Prompt

After committing, use `AskUserQuestion` to ask:

> Run `/eval-design` for adversarial evaluation? (default: 90 points / 3 rounds)

- **Yes** → invoke `/eval-design` via `Skill` tool
- **Custom** → invoke `/eval-design --target X --iterations Y` via `Skill` tool
- **No** → proceed to `/breakdown-tasks`

## Integration

Works well with skills:

- `/write-prd` - Creates PRD input and manifest
- `/ui-design` - Parallel skill for UI features
- `/eval-design` - Evaluate tech-design.md quality before handing off to breakdown-tasks
- `/breakdown-tasks` - Uses tech-design.md to create tasks
