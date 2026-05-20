# Sitemap Schema Reference

Read the full example at `templates/sitemap-example.json`.

## Key Fields

| Field | Description |
|-------|-------------|
| `baseUrl` | Application base URL (e.g. `http://localhost:3456`) |
| `updatedAt` | Last update time (RFC3339 format) |
| `layout.name` | Layout component name (e.g. `AppLayout`) |
| `layout.wraps` | List of routes sharing this layout |
| `layout.elements[]` | Layout-level shared elements (sidebar, top nav, etc.), ID format `L-NNN` |
| `pages[].elements[].role` | Accessibility role (button, heading, etc.) |
| `pages[].elements[].name` | Accessibility name |
| `pages[].elements[].level` | Heading level (heading role only) |
| `pages[].elements[].label` | Associated label text (form elements like textbox only) |
| `pages[].elements[].placeholder` | Placeholder text (textbox only) |
| `pages[].states[]` | Dynamic states (modal, tab panel, dropdown, etc.) |
| `pages[].states[].trigger` | Trigger element ID (e.g. `"E-002"`) |
| `pages[].states[].elements` | Elements within the state (also with E-NNN IDs) |

## Layout Elements vs Page Elements

- `layout.elements`: Elements shared across all wrapped pages (navbar, sidebar, footer), ID format `L-NNN`
- `pages[].elements`: Elements unique to the current page, ID format `E-NNN`
- During test script generation, layout elements are available in any wrapped page

## Element ID Assignment Rules

- **Layout elements**: format `L-NNN`, globally unique, independent numbering space
- **Page elements**: format `E-NNN`, globally unique (base and states share ID space)
- First generation: start from `L-001` and `E-001` respectively
- Incremental update: new elements start from current max ID + 1
- IDs are never reused (deleted IDs are not recycled)

## Dynamic Route Parameterization

| URL segment pattern | Replace with | Example |
|---------------------|-------------|---------|
| Pure numeric | `:id` | `/tasks/42` -> `/tasks/:id` |
| UUID format | `:uuid` | `/orders/550e8400-...` -> `/orders/:uuid` |
| 32-char hex | `:hash` | `/files/a1b2c3d4e5f6...` -> `/files/:hash` |

During deduplication, routes with the same template keep only one entry. `layout.wraps` also uses template form.

Also merge manually added routes from existing sitemap that are not found by either source.
