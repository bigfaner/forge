---
name: gen-sitemap
description: Auto-generate and maintain sitemap.json for a web app. Uses agent-browser to explore routes, capture accessibility tree, and discover dynamic states. Preserves element IDs across runs.
argument-hints:
  - name: base-url
    description: Application base URL to explore (e.g. http://localhost:3456). Optional if config.yaml exists.
    required: false
  - name: api-base-url
    description: Backend API base URL (e.g. http://localhost:8080). Optional if config.yaml exists.
    required: false
---

# /gen-sitemap

Auto-generate and maintain `docs/sitemap/sitemap.json` for web applications.

**Core principle**: The sitemap is a complete structural map of the web application, serving as the sole input source for Playwright locator generation. Element IDs (E-NNN) are stable identifiers preserved across runs.

## Prerequisites

- Web application is running and accessible
- agent-browser is installed (optional tool for auto-exploring page structure and dynamic states)

**Verify agent-browser installation:**

```bash
npx agent-browser@latest --version
```

If it fails, install first:

```bash
npx agent-browser@latest install
```

<HARD-RULE>
Always use a version-qualified `npx agent-browser@latest` invocation. Do not use bare `npx agent-browser` (unpinned) — it may pull breaking changes between runs. For CI reproducibility, consider pinning to a specific version (e.g. `agent-browser@0.4.x`) in your environment.
</HARD-RULE>

After installation, re-run this command.

## Config Resolution

Before executing the workflow, resolve configuration sources:

1. Check if `tests/e2e/config.yaml` exists
2. **If not found**: copy from template `plugins/forge/references/shared/config.yaml` to `tests/e2e/config.yaml`, then **abort and prompt the user**:

```
Created tests/e2e/config.yaml (template). Please fill in values for your environment and re-run.
  baseUrl: actual application URL
  apiBaseUrl: actual backend API URL
  username/password: test account credentials (if app requires auth)
  loginLocators: uncomment and customize if login page locators differ from defaults
```

3. **If exists**: read `baseUrl`, `apiBaseUrl`, `username`, `password` fields

**base-url priority**: command-line argument > `baseUrl` from config.yaml > abort with error

**apiBaseUrl rule**: `apiBaseUrl` is the backend's raw base address (e.g. `http://localhost:8080`), **without any path prefix**. When constructing request URLs, path prefixes (e.g. `/v1`) must be read from backend route definitions (e.g. `backend/internal/handler/router.go` `r.Group(...)`), not assumed. Correct: `${apiBaseUrl}/v1/auth/login`, not `${apiBaseUrl}/api/v1/auth/login`.

**Authenticated page exploration**: If config.yaml has non-empty `username` and `password`, log in before exploring each page in Step 3:

```
ab('open <baseUrl>/login')
ab('wait --load networkidle')
// Fill login form with username/password from config.yaml
ab('click <login_button>')
ab('wait --load networkidle')
```

This ensures agent-browser can access authenticated pages without exploration interruptions.

## Schema

See `plugins/forge/references/shared/sitemap.json` for a full example.

**Key fields:**

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

**Layout elements vs page elements:**

- `layout.elements`: Elements shared across all wrapped pages (navbar, sidebar, footer), ID format `L-NNN`
- `pages[].elements`: Elements unique to the current page, ID format `E-NNN`
- During test script generation, layout elements are available in any wrapped page

## Workflow

```
1. Load existing sitemap → 2. Analyze layout → 3. Discover routes → 4. Explore pages → 5. Merge & write
```

### Step 1: Load Existing Sitemap

Read `docs/sitemap/sitemap.json` if it exists.

Build an ID index: use `route + role + name` as key, mapping to existing IDs.

If no existing sitemap, start numbering from `E-001`.

### Step 2: Analyze Layout

<EXTREMELY-IMPORTANT>
**This step MUST execute before page exploration.** The goal is to identify shared layout to avoid extracting layout elements on every page.

**Read code first, then explore.** Do not blindly explore first and then go back to deduplicate.
</EXTREMELY-IMPORTANT>

**2a. Read route definitions**: Find the application's route configuration files (e.g. `App.tsx`, `router.ts`, `routes.tsx`), identify layout nesting:

```
<Route element={<AppLayout />}>      ← shared layout
  <Route path="/" element={<Dashboard />} />
  <Route path="/settings" element={<Settings />} />
  ...
</Route>
<Route path="/login">                ← no layout (standalone page)
```

Record:
- `layout.name`: layout component name (e.g. `AppLayout`)
- `layout.wraps`: list of routes wrapped by the layout
- Routes not wrapped by any layout (e.g. `/login`) are standalone pages

**2b. Explore layout elements**: Take a snapshot of the first wrapped page, extracting **only layout-level** elements:

```
ab('open <baseUrl><first_wrapped_route>')
ab('wait --load networkidle')
snapshot = abJson('snapshot -i')
```

From the snapshot, identify layout regions (typically `role=navigation`, `role=banner`, sidebar containers, etc.) and extract their elements into `layout.elements`.

**If route code cannot be read** (pure HTML project or no source access): skip this step and treat all elements as page-level. After exploring the first two pages in Step 4, compare snapshots for elements with identical `role + name` as layout candidates, assign them to `layout.elements`, and filter them out during subsequent page exploration.

### Step 3: Discover Routes

Use agent-browser to navigate to the user-provided `base-url` and extract all links from the page:

```
ab('open <base-url>')
ab('wait --load networkidle')
links = abJson('snapshot -i')  // extract all role=link nodes' href
```

1. Filter to same-origin paths (exclude external links, `mailto:`, `javascript:`)
2. Deduplicate to get a route list
3. Recursively extract links from each new route (breadth-first, max depth 3)
4. Merge manually added routes from existing sitemap

**Dynamic route handling**: Routes with parameters (e.g. `/tasks/123`) are recorded as template form `/tasks/:id`. Parameterization rules:

| URL segment pattern | Replace with | Example |
|---------------------|-------------|---------|
| Pure numeric | `:id` | `/tasks/42` → `/tasks/:id` |
| UUID format | `:uuid` | `/orders/550e8400-...` → `/orders/:uuid` |
| 32-char hex | `:hash` | `/files/a1b2c3d4e5f6...` → `/files/:hash` |

During deduplication, routes with the same template keep only one entry. `layout.wraps` also uses template form.

### Step 4: Explore Pages

For each route, explore with agent-browser one by one:

```
ab('open <baseUrl><route>')
ab('wait --load networkidle')
snapshot = abJson('snapshot -i')
```

#### Layout Element Filtering

<HARD-RULE>
If Step 2 identified a shared layout, this step **MUST skip layout elements**.
Compare the snapshot against `layout.elements`, filtering out elements matching by `role + name`.
Only page-specific content area elements are included in `pages[].elements`.
</HARD-RULE>

#### Base Element Extraction

1. Get page title
2. Extract elements from snapshot with filters:
   - Exclude elements already in `layout.elements` (matched by `role + name`)
   - `role` ∈ {button, link, heading, textbox, checkbox, radio, combobox, tab, dialog, alert, navigation, search, form, menuitem, switch}
   - `name` is non-empty
3. For each element, record full attributes:
   - Common: `{ role, name }`
   - heading: additionally record `level`
   - textbox/combobox: additionally record `label` (associated label text) and `placeholder`

#### Dynamic State Exploration

For trigger elements with role=button/tab/disclosure and non-empty name:

```
ab('click @eN')
ab('wait --load networkidle')
state_snapshot = abJson('snapshot -i')
// Extract new elements (compare with base snapshot)
ab('press Escape')  // or ab('click @close_btn') to reset
```

1. Compare state snapshot with base snapshot, extract new elements
2. Record as `states` entry: `{ name, trigger: "<elementID>", elements: [...] }`
3. `trigger` references the trigger element's E-NNN ID (e.g. `"E-002"`)
4. In-state elements also receive E-NNN IDs

> **Note**: `@eN` is agent-browser CLI's element reference syntax, used only during sitemap generation. Generated test scripts (`*.spec.ts`) must NOT use `@eN`; they must use Playwright Locator API.

```
ab('close')
```

### Step 5: Merge & Write

For each element (including layout and in-states elements), match against existing sitemap using the `route + role + name` triplet:

- **Match found** → preserve existing ID
- **No match** → assign new ID (current max ID + 1)
- **Existing ID has no match** → element was removed, delete from sitemap

Write to `docs/sitemap/sitemap.json`.

**Report changes:**

```
Sitemap updated: docs/sitemap/sitemap.json
  Layout: AppLayout (5 shared elements, L-001..L-005)
  3 pages, 12 elements (page-specific) + 8 elements (states)
  +4 new elements (E-015..E-018)
  -2 removed elements (previously on /settings)
  17 unchanged
```

## Element ID Assignment Rules

- **Layout elements**: format `L-NNN`, globally unique, independent numbering space
- **Page elements**: format `E-NNN`, globally unique (base and states share ID space)
- First generation: start from `L-001` and `E-001` respectively
- Incremental update: new elements start from current max ID + 1
- IDs are never reused (deleted IDs are not recycled)
