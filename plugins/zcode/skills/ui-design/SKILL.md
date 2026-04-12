---
name: ui-design
description: Use after PRD ui-functions are defined to create UI design specifications. Parallel to /design-tech, reads prd/prd-ui-functions.md.
---

# UI Design

## Overview

从 PRD 的 UI 功能需求产出 UI 设计规格文档。

**核心原则**：定义 HOW 界面呈现和交互，与 PRD 的 WHAT（需求层）分离。

<HARD-GATE>
Do NOT write any implementation code. This skill produces a design specification document only.
</HARD-GATE>

## Prerequisites

检查上一阶段产物，缺失则中止并提示用户：

```bash
ls docs/features/<slug>/prd/prd-ui-functions.md
```

| 产物 | 缺失时提示 |
|------|-----------|
| `prd/prd-ui-functions.md` | 先执行 `/write-prd` 并补充 UI 功能要点 |

## When to Use

**Trigger conditions:**
- PRD with `prd/prd-ui-functions.md` exists
- User asks to design the UI

**Skip when:**
- No UI functions defined (backend/API/CLI features)
- UI design already exists

## Process Flow

```
1. Read manifest → 2. Read UI functions → 3. Explore patterns → 4. Draft design → 5. Review → 6. Update manifest
```

## Step 1: Read Manifest

Read `manifest.md` to locate `prd/prd-ui-functions.md`.

## Step 2: Read UI Functions

Read `prd/prd-ui-functions.md` to understand UI requirements.

## Step 3: Explore Existing Patterns

- Check for existing design system or component library
- Review existing UI components in the project

## Step 4: Draft UI Design

For each UI function, define:
- Layout structure (component hierarchy)
- States (loading, empty, error, populated)
- Interactions (triggers, actions, feedback)
- Data binding (UI element → data field)

## Step 5: Write UI Design

Save to `ui/ui-design.md` using `templates/ui-design.md`.

## Step 6: Update Manifest

Update `manifest.md`:
- Add UI Design row to Documents table
- Add traceability links from UI Functions to UI Design sections
- Advance status to `design` if `/design-tech` already completed

## Integration

Works well with:
- `/write-prd` — Produces `prd/prd-ui-functions.md` input
- `/design-tech` — Parallel skill; both must complete before breakdown
- `/eval-design` — Evaluates UI design alongside tech design
