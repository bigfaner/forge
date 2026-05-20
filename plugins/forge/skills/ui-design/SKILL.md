---
name: ui-design
description: Use after PRD ui-functions are defined to create UI design specifications. Supports design style selection and HTML prototype generation.
effort: high
---

# UI Design

## Overview

Produce UI design specifications from PRD's UI function requirements, with HTML prototype generation.

**Core principle**: Define HOW the interface looks and behaves, separated from the PRD's WHAT (requirements layer).

<HARD-GATE>
Do NOT write any implementation code (except HTML prototypes). This skill produces design specification documents and optionally HTML/CSS/JS prototypes.
</HARD-GATE>

## Prerequisites

Check previous stage artifacts. Abort and prompt user if missing:

```bash
ls docs/features/<slug>/prd/prd-ui-functions.md
```

| Artifact                    | Missing prompt                                        |
| --------------------------- | ----------------------------------------------------- |
| `prd/prd-ui-functions.md` | Run `/write-prd` first and add UI function highlights |

## When to Use

**Trigger conditions:**

- PRD with `prd/prd-ui-functions.md` exists
- User asks to design the UI

**Skip when:**

- No UI functions defined (backend/API/CLI features)
- UI design already exists

## Process Flow

```
1. Read manifest → 2. Read UI functions → 2.5. Extract nav architecture → 3. Select design style → 4. Draft design → 5. Write design → 6. Update manifest → 7. Auto eval-ui → 8. Generate prototype → 9. Human review → 10. Reconcile PRD → 11. Next step
```

## Step 1: Read Manifest

Read `manifest.md` to locate `prd/prd-ui-functions.md`.

## Step 2: Read UI Functions

Read `prd/prd-ui-functions.md` to understand UI requirements.

**Placement awareness**: For each UI Function, note its Placement section:
- `new-page`: design the complete page including layout, navigation, and all contained components
- `existing-page:<route>`: design only the new component(s) that will be embedded. The Placement's Position field tells you WHERE in the existing page — this constrains the component's layout (e.g., width must match existing content area).

## Step 2.5: Extract Navigation Architecture

Read the `## Navigation Architecture` section from `prd-ui-functions.md`.

If it exists:
- Use it as the single source of truth for all inter-page navigation
- Note the `platform` field to determine navigation style

If it does NOT exist:
- Prompt the user: "PRD lacks a Navigation Architecture section. Consider running `/write-prd` to add it, or define it here:"
- Use AskUserQuestion to collect: platform (web/mobile), primary navigation entries, secondary pages, navigation rules
- Write the collected structure back into prd-ui-functions.md (append after UI Scope section)

Then read the platform-specific navigation rules from `templates/platforms/{web,mobile,tui}.md` based on the identified platform.

## Step 3: Select Design Style

Select design style based on the platform identified in Step 2.5. Follow the priority chain and platform-specific rules in `${CLAUDE_SKILL_DIR}/rules/style-selection.md`:

1. **Check for user-provided DESIGN.md** (project root or feature directory) -- use directly if found.
2. **If no user-provided style**: let the user choose from built-in styles (5 web/mobile styles or 2 TUI themes).
3. **If built-in styles insufficient**: clone additional styles from the awesome-design-md repo.
4. **Multi-platform**: select a style independently for each platform.

Inline the selected style content as the `Design System` section in `ui-design.md`.

## Step 4: Draft UI Design

For each UI function, define:

- Layout structure (component hierarchy)
- States (loading, empty, error, populated)
- Interactions (triggers, actions, feedback)
- Data binding (UI element → data field)

**For existing-page UI Functions**: the Component section in `ui-design.md` MUST include a Placement subsection that specifies where in the existing page this component will be placed. This comes from the PRD's Placement Position field but may be refined with layout constraints discovered during design (e.g., "full width of the content area, above the sub-items table, below the progress summary heading").

### TUI Panel Design Requirements

When platform=tui, each panel MUST include all mandatory structural requirements per `${CLAUDE_SKILL_DIR}/rules/tui-panel-requirements.md`:
- 5 mandatory structural items: ASCII mockup, dimensions, character palette, color mapping, edge cases
- Additional per-panel specs: states, key bindings, data binding

All design decisions must follow the design style selected in Step 3.

## Step 5: Write UI Design

Save using `templates/ui-design.md` (for web/mobile) or the TUI component template section within it (for TUI panels).

Output file naming:
| Platform | Output File |
|----------|------------|
| Web only | `ui/ui-design.md` |
| TUI only | `ui/ui-design-tui.md` |
| Mobile only | `ui/ui-design.md` |
| Multi-platform (web + tui) | `ui/ui-design-web.md` + `ui/ui-design-tui.md` |

The Design System section references the style selected in Step 3 (inline core tokens, no external file references).

## Step 6: Update Manifest

Update `manifest.md`:

- Add UI Design row to Documents table
- Add traceability links from UI Functions to UI Design sections
- Advance status to `design`

For multi-platform features, add a row for each platform's design file:
- `ui/ui-design-web.md` — web design
- `ui/ui-design-tui.md` — TUI design

Use `templates/manifest-update-ui.md` for the update pattern.

## Step 7: Auto Eval UI Design

Automatically invoke `/eval-ui` to evaluate `ui-design.md`. Default: 950 points / 3 rounds.

Invoke via `Skill` tool: `eval-ui`

eval-ui runs its own score → gate → revise loop and produces a report at `docs/features/<slug>/ui/eval/report.md`.

After eval-ui completes, check the final score:

| Condition | Action |
|-----------|--------|
| Score >= 950 | Proceed to Step 8 (prototype generation) |
| Score < 950 | Report to user and ask: "UI design score is {{SCORE}}/1000, below the 950 threshold. Continue revising?" — if yes, re-invoke eval-ui with additional iterations; if no, proceed to prototype anyway |

The eval report is attached at `docs/features/<slug>/ui/eval/report.md` for reference.

## Step 8: Generate Prototype

Generate HTML prototype from `ui-design.md` and the selected design style.

### For Web / Mobile Platform

- Follow `templates/prototype.md` specifications
- Multi-file structure: shared CSS/JS + per-page HTML + index.html navigation
- Implement all component states (loading/empty/error/populated)
- Interactive behaviors functional (navigation, modals, tabs, etc.)
- Responsive layout

Save to: `docs/features/<slug>/ui/prototype/`

### For TUI Platform

- Follow `templates/prototype.md` specifications for file structure
- HTML simulates a terminal window: dark monospace background, box-drawing characters rendered via CSS
- All panels rendered in a single `index.html` within a black "terminal window" div
- Bottom area simulates key buttons (`[Tab]`, `[1]`, `[q]`, `[:command]`) to switch panels
- Each panel's ASCII mockup from ui-design-tui.md is rendered as pre-formatted text
- Implement panel states (loading/empty/error/populated) as toggleable views

Save to: `docs/features/<slug>/ui/prototype/`

### Multi-Platform

When multiple platforms are present:
- Each platform gets its own prototype directory:
  - `docs/features/<slug>/ui/prototype/web/` (or `prototype/` for single-platform)
  - `docs/features/<slug>/ui/prototype/tui/`

## Step 9: Human Review Gate

<HARD-GATE>
Present the prototype to the user for review. Do NOT proceed to tech-design until the user explicitly approves.
</HARD-GATE>

Present the prototype summary and use `AskUserQuestion`:

> Prototype generated at `docs/features/<slug>/ui/prototype/`. Open `index.html` in a browser to review. Approve the prototype?

- **Approved** → proceed to next step
- **Request changes** → revise prototype based on feedback, then re-present for approval
- **Skip prototype** → discard prototype files, proceed without prototype

## Step 10: Reconcile PRD (Optional)

After prototype is approved, compare the generated prototype against `prd-ui-functions.md`.

Check for discrepancies:
- Pages in prototype but NOT in PRD → suggest adding a new UF
- Interactions in prototype but different from PRD → suggest updating UF interaction flow
- Data fields in prototype but not in PRD → suggest updating UF data requirements
- Navigation structure in prototype but not in PRD → suggest adding Navigation Architecture

If discrepancies found:
- List them with specific file references
- Ask user: "Prototype has N differences from PRD. Write back to prd-ui-functions.md?"
- If yes → update prd-ui-functions.md to match the prototype (prototype is the source of truth)
- If no → note discrepancies in ui-design.md for future reference

## Step 11: Next Step

After prototype is approved or skipped, use `AskUserQuestion`:

> Proceed to `/tech-design` to create technical design?

- **Yes** → invoke `/tech-design` via `Skill` tool
- **No** → done

## Integration

Works well with:

- `/write-prd` — Produces `prd/prd-ui-functions.md` input
- `/tech-design` — Next skill after UI design; informed by UI decisions
- `/eval-ui` — Adversarial evaluation after UI design is created
