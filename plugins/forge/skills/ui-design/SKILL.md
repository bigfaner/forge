---
name: ui-design
description: Use after PRD ui-functions are defined to create UI design specifications. Supports design style selection and HTML prototype generation.
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

Then read the platform-specific navigation rules from `templates/platforms/{web,mobile}.md` based on the identified platform.

## Step 3: Select Design Style

Select design style by priority:

### Priority 1: User-provided DESIGN.md

Check the following locations; use directly if found:

- Project root `DESIGN.md`
- Feature directory `docs/features/<slug>/ui/style.md`

If the user specifies a custom path, prioritize that.

### Priority 2: Built-in Styles

If no user-provided DESIGN.md, let the user choose from 5 built-in styles:

| Style           | Vibe                     | Best for                                     |
| --------------- | ------------------------ | -------------------------------------------- |
| **Vercel**      | Black & white minimal, developer-tool feel | Developer platforms, CLI tools, technical docs |
| **Shadcn**      | Zinc neutral, functional minimalism | SaaS, admin panels, tool applications        |
| **Tailwind UI** | Indigo primary, professional warmth | General SaaS, marketing pages, enterprise    |
| **Stripe**      | Purple gradients, light elegance | Fintech, brand sites, payment products       |
| **Apple**       | Generous whitespace, image-driven, premium | Consumer products, brand sites, marketing     |

Use `AskUserQuestion` tool for user selection with brief descriptions.

Built-in style files located at: `templates/styles/{vercel,shadcn,tailwind-ui,stripe,apple}.md`

### Priority 3: Clone More Styles from Repo

If the 5 built-in styles are insufficient, clone additional styles from the awesome-design-md repo:

```bash
# Clone to a temp directory (not into project)
git clone --depth 1 git@github.com:VoltAgent/awesome-design-md.git /tmp/awesome-design-md

# Then use npx to fetch a specific site's DESIGN.md:
npx getdesign@latest add <site-name>
```

> Note: The repo's DESIGN.md files are hosted externally at getdesign.md. The `npx getdesign@latest` CLI fetches them.

### Using the Selected Style

Inline the selected style content as the `Design System` section in `ui-design.md`. All subsequent component designs must follow the style's color, typography, and component specifications.

## Step 4: Draft UI Design

For each UI function, define:

- Layout structure (component hierarchy)
- States (loading, empty, error, populated)
- Interactions (triggers, actions, feedback)
- Data binding (UI element → data field)

**For existing-page UI Functions**: the Component section in `ui-design.md` MUST include a Placement subsection that specifies where in the existing page this component will be placed. This comes from the PRD's Placement Position field but may be refined with layout constraints discovered during design (e.g., "full width of the content area, above the sub-items table, below the progress summary heading").

All design decisions must follow the design style selected in Step 3.

## Step 5: Write UI Design

Save to `ui/ui-design.md` using `templates/ui-design.md`.

The Design System section references the style selected in Step 3 (inline core tokens, no external file references).

## Step 6: Update Manifest

Update `manifest.md`:

- Add UI Design row to Documents table
- Add traceability links from UI Functions to UI Design sections
- Advance status to `design`

Use `templates/manifest-update-ui.md` for the update pattern.

## Step 7: Auto Eval UI Design

Automatically invoke `/eval-ui` to evaluate `ui-design.md`. Default: 90 points / 3 rounds.

Invoke via `Skill` tool: `eval-ui`

eval-ui runs its own score → gate → revise loop and produces a report at `docs/features/<slug>/ui/eval/report.md`.

After eval-ui completes, proceed to prototype generation regardless of final score — the eval report is attached for reference during prototype generation.

## Step 8: Generate Prototype

Generate HTML prototype from `ui-design.md` and the selected design style:

- Follow `templates/prototype.md` specifications
- Multi-file structure: shared CSS/JS + per-page HTML + index.html navigation
- Implement all component states (loading/empty/error/populated)
- Interactive behaviors functional (navigation, modals, tabs, etc.)
- Responsive layout

Save to: `docs/features/<slug>/ui/prototype/`

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
