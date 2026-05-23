# UI Functions Rules

Detailed rules for writing prd-ui-functions.md in Step 8.

## Placement Rules (mandatory for every UI Function)

1. Read `docs/sitemap/sitemap.json` to understand existing page inventory
2. For each UI Function, declare its Placement:
   - `new-page` — this function creates a brand new page
   - `existing-page:<route>` — this function adds UI to an existing page (route from sitemap)
3. For `existing-page`, also specify the Position within the page (e.g., "above sub-items table")
4. After all UI Functions are defined, compile the Page Composition summary table at the end of the document

**Validation**: Every UI Function MUST have a Placement section. Missing Placement → error, do not proceed.

## Navigation Architecture Platform Handling

The template contains two conditional navigation sections. Render exactly one based on platform:

- **platform=tui**: Render the **TUI Navigation** section (Keymap, Panel Layout, Modes, Navigation Rules). Omit the Pointer-Driven Navigation section entirely.
- **platform=web|mobile|mini-program|tablet**: Render the **Pointer-Driven Navigation** section (Primary Navigation, Secondary Pages, Navigation Rules). Omit the TUI Navigation section entirely.

TUI uses keyboard-driven navigation (keymap + panels + modes) instead of pointer-driven navigation (pages + routes + icons). Do not mix the two patterns.

## Downstream Impact

Placement determines how downstream skills structure their output:

- **`/breakdown-tasks`**: `existing-page` generates a Build task + an Integrate task (wire component into existing page); `new-page` generates a Build task + a Page Assembly task (create page file, register route, compose components)
