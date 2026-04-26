---
name: ui-design
description: Use after PRD ui-functions are defined to create UI design specifications. Supports design style selection and HTML prototype generation.
---

# UI Design

## Overview

Produce UI design specifications from PRD's UI function requirements, with optional HTML prototype generation.

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
1. Read manifest → 2. Read UI functions → 3. Select design style → 4. Draft design → 5. Write design → 6. Update manifest → 7. Eval prompt → 8. Prototype prompt → 9. Generate prototype (optional)
```

## Step 1: Read Manifest

Read `manifest.md` to locate `prd/prd-ui-functions.md`.

## Step 2: Read UI Functions

Read `prd/prd-ui-functions.md` to understand UI requirements.

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

All design decisions must follow the design style selected in Step 3.

## Step 5: Write UI Design

Save to `ui/ui-design.md` using `templates/ui-design.md`.

The Design System section references the style selected in Step 3 (inline core tokens, no external file references).

## Step 6: Update Manifest

Update `manifest.md`:

- Add UI Design row to Documents table
- Add traceability links from UI Functions to UI Design sections
- Advance status to `design` if `/tech-design` already completed

Use `templates/manifest-update-ui.md` for the update pattern.

## Step 7: Adversarial Eval Prompt

After updating manifest, use `AskUserQuestion` to ask:

> Run `/eval-ui` for adversarial evaluation of the UI design? (default 80/3)

- **Yes** → invoke `/eval-ui` via `Skill` tool
- **Custom** → invoke `/eval-ui --target X --iterations Y` via `Skill` tool
- **No** → proceed to prototype prompt

## Step 8: Prototype Prompt

Use `AskUserQuestion`:

> Generate HTML/CSS/JS interactive prototype from the UI design?

- **Yes** → proceed to Step 9
- **No** → done

## Step 9: Generate Prototype (Optional)

Generate HTML prototype from `ui-design.md` and the selected design style:

- Follow `templates/prototype.md` specifications
- Multi-file structure: shared CSS/JS + per-page HTML + index.html navigation
- Implement all component states (loading/empty/error/populated)
- Interactive behaviors functional (navigation, modals, tabs, etc.)
- Responsive layout

Save to: `docs/features/<slug>/ui/prototype/`

## Integration

Works well with:

- `/write-prd` — Produces `prd/prd-ui-functions.md` input
- `/tech-design` — Parallel skill; both must complete before breakdown
- `/eval-ui` — Adversarial evaluation after UI design is created
