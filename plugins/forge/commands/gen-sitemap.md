---
name: gen-sitemap
description: Auto-generate and maintain sitemap.json for a web app. Uses agent-browser to explore routes, capture accessibility tree, and discover dynamic states. Preserves element IDs across runs.
allowed-tools: Bash Read Write Grep Glob
argument-hint: "[base-url] [api-base-url]"
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
2. **If not found**: create `tests/e2e/config.yaml` from the following template, then **abort and prompt the user**:

```yaml
# E2E Test Environment Configuration
# Adjust values to match your local/dev environment
# IMPORTANT: add config.yaml to .gitignore to prevent credential leaks

# Application URLs
baseUrl: http://localhost:3456        # Frontend URL
apiBaseUrl: http://localhost:8080     # Backend API URL

# Default timeout (ms)
timeout: 30000

# Auth credentials (used by loginViaUI and getApiToken)
username: admin
password: password

# Login page locators (optional — overrides regex defaults in loginViaUI)
loginLocators:
  usernameField: "username|email"
  passwordField: "password"
  submitButton: "login|sign in|submit"
```

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

Full example:

```json
{
	"baseUrl": "http://localhost:3456",
	"updatedAt": "2026-04-25T10:00:00Z",
	"layout": {
		"name": "AppLayout",
		"wraps": ["/", "/settings", "/tasks/:id"],
		"elements": [
			{ "id": "L-001", "role": "navigation", "name": "Main nav" },
			{ "id": "L-002", "role": "link", "name": "Dashboard" },
			{ "id": "L-003", "role": "link", "name": "Settings" },
			{ "id": "L-004", "role": "link", "name": "Profile" },
			{ "id": "L-005", "role": "button", "name": "Logout" }
		]
	},
	"pages": [
		{
			"route": "/login",
			"title": "Login",
			"elements": [
				{ "id": "E-030", "role": "heading", "name": "Sign In", "level": 1 },
				{ "id": "E-031", "role": "textbox", "name": "Username", "label": "Username", "placeholder": "Enter username" },
				{ "id": "E-032", "role": "textbox", "name": "Password", "label": "Password", "placeholder": "Enter password" },
				{ "id": "E-033", "role": "button", "name": "Login" }
			]
		},
		{
			"route": "/",
			"title": "Dashboard",
			"elements": [
				{ "id": "E-001", "role": "heading", "name": "Dashboard", "level": 1 },
				{ "id": "E-002", "role": "button", "name": "Add Task" },
				{ "id": "E-003", "role": "textbox", "label": "Search", "placeholder": "Search tasks..." }
			],
			"states": [
				{
					"name": "add-task-modal",
					"trigger": "E-002",
					"elements": [
						{ "id": "E-004", "role": "dialog", "name": "Add Task" },
						{ "id": "E-005", "role": "textbox", "label": "Task Title" },
						{ "id": "E-006", "role": "button", "name": "Save" },
						{ "id": "E-007", "role": "button", "name": "Cancel" }
					]
				}
			]
		},
		{
			"route": "/settings",
			"title": "Settings",
			"elements": [
				{ "id": "E-010", "role": "heading", "name": "Settings", "level": 1 },
				{ "id": "E-011", "role": "tab", "name": "General" },
				{ "id": "E-012", "role": "tab", "name": "Notifications" },
				{ "id": "E-013", "role": "button", "name": "Save" }
			]
		},
		{
			"route": "/tasks/:id",
			"title": "Task Detail",
			"elements": [
				{ "id": "E-020", "role": "heading", "name": "Task Detail", "level": 1 },
				{ "id": "E-021", "role": "button", "name": "Edit" },
				{ "id": "E-022", "role": "button", "name": "Delete" },
				{ "id": "E-023", "role": "link", "name": "Back" }
			],
			"states": [
				{
					"name": "delete-confirm",
					"trigger": "E-022",
					"elements": [
						{ "id": "E-024", "role": "dialog", "name": "Confirm Delete" },
						{ "id": "E-025", "role": "button", "name": "Confirm" },
						{ "id": "E-026", "role": "button", "name": "Cancel" }
					]
				}
			]
		}
	]
}
```

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
1. Load existing sitemap → 2. Analyze layout & build route registry → 3. Discover routes → 4. Explore pages → 5. Merge & dedup & validate → 6. Write
```

### Step 1: Load Existing Sitemap

Read `docs/sitemap/sitemap.json` if it exists.

Build an ID index: use `route + role + name` as key, mapping to existing IDs.

If no existing sitemap, start numbering from `E-001`.

### Step 2: Analyze Layout & Build Route Registry

<EXTREMELY-IMPORTANT>
**This step MUST execute before page exploration.** The goal is to identify shared layout to avoid extracting layout elements on every page.

**Read code first, then explore.** Do not blindly explore first and then go back to deduplicate.
</EXTREMELY-IMPORTANT>

**2a. Build complete route registry**: Find the application's route configuration files (e.g. `App.tsx`, `router.ts`, `routes.tsx`, Next.js `pages/` or `app/` directory, Vue `router.ts`/`router.js`). Extract ALL registered routes:

```
1. Scan route configuration files:
   - React Router: App.tsx, router.tsx, routes.tsx
   - Next.js: pages/ or app/ directory structure
   - Vue Router: router.ts, router.js
   - Generic: grep for path patterns in route definitions

2. Extract ALL registered routes:
   - Static routes: /items, /settings
   - Dynamic routes: /items/:id, /users/:userId
   - Nested routes: /items/:id/sub/:subId
   - Layout-wrapped vs standalone classification

3. Identify layout nesting:
   <Route element={<AppLayout />}>      ← shared layout
     <Route path="/" element={<Dashboard />} />
     <Route path="/settings" element={<Settings />} />
     ...
   </Route>
   <Route path="/login">                ← no layout (standalone page)

4. Record:
   - layout.name: layout component name (e.g. `AppLayout`)
   - layout.wraps: list of routes wrapped by the layout
   - Routes not wrapped by any layout (e.g. `/login`) are standalone pages

5. Build route registry with source tracking:
   routes = [
     { template: "/items", source: "router", layout: "AppLayout" },
     { template: "/items/:id", source: "router", layout: "AppLayout" },
     { template: "/login", source: "router", layout: null }
   ]
```

**2b. Identify shared elements via multi-page comparison**: Visit 3-5 wrapped routes to identify shared elements by intersection:

```
Visit up to 5 wrapped routes with agent-browser:
  ab('open <baseUrl><route_1>')
  ab('wait --load networkidle')
  snapshot_1 = abJson('snapshot -i')

  ab('open <baseUrl><route_2>')
  ab('wait --load networkidle')
  snapshot_2 = abJson('snapshot -i')

  // ... up to 5 routes total

Compute intersection: elements present in ALL visited snapshots with same role+name
→ These are shared elements (sidebar, nav, header, footer, etc.)

Classify shared elements:
- Layout-level (sidebar, top nav, footer) → layout.elements with L-NNN IDs
- Pattern-level (e.g. breadcrumb on multiple but not all pages) → noted as "partial-shared"
  (still assigned to individual pages, but identified for awareness)
```

**If route code cannot be read** (pure HTML project or no source access): skip this step and treat all elements as page-level. After exploring the first two pages in Step 4, compare snapshots for elements with identical `role + name` as layout candidates, assign them to `layout.elements`, and filter them out during subsequent page exploration.

### Step 3: Discover Routes (dual-source)

Combine router registry with link-crawling for complete coverage:

```
1. Start with route registry from Step 2a as the BASE
2. Navigate to baseUrl with agent-browser to crawl links:
   ab('open <base-url>')
   ab('wait --load networkidle')
   links = abJson('snapshot -i')  // extract all role=link nodes' href
3. Filter to same-origin paths (exclude external links, `mailto:`, `javascript:`)
4. Deduplicate to get crawled route list
5. Recursively extract links from each new route (breadth-first, max depth 3)
6. Merge: router routes + crawled routes = complete route list
7. Report source counts: "12 routes from router, 1 additional from crawling"
```

**Dynamic route handling**: Routes with parameters (e.g. `/tasks/123`) are recorded as template form `/tasks/:id`. Parameterization rules:

| URL segment pattern | Replace with | Example |
|---------------------|-------------|---------|
| Pure numeric | `:id` | `/tasks/42` → `/tasks/:id` |
| UUID format | `:uuid` | `/orders/550e8400-...` → `/orders/:uuid` |
| 32-char hex | `:hash` | `/files/a1b2c3d4e5f6...` → `/files/:hash` |

During deduplication, routes with the same template keep only one entry. `layout.wraps` also uses template form.

Also merge manually added routes from existing sitemap that are not found by either source.

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
   - `role` ∈ the full ARIA role set:
     {button, link, heading, textbox, checkbox, radio, combobox, tab, dialog, alert, navigation, search, form, menuitem, switch,
      table, row, cell, columnheader, rowheader,
      grid, listbox, option, list, listitem,
      tooltip, progressbar, meter, slider, spinbutton,
      status, log, marquee, timer,
      img, separator, group, region,
      feed, article, figure, caption}
   - `name` is non-empty
3. For each element, record full attributes:
   - Common: `{ role, name }`
   - heading: additionally record `level`
   - textbox/combobox: additionally record `label` (associated label text) and `placeholder`

#### Dynamic State Exploration

Explore dynamic states triggered by multiple interaction types:

**1. Click triggers** (existing): For elements with role=button/tab/disclosure and non-empty name:

```
ab('click @eN')
ab('wait --load networkidle')
state_snapshot = abJson('snapshot -i')
// Extract new elements (compare with base snapshot)
ab('press Escape')  // or ab('click @close_btn') to reset
```

**2. Hover triggers**: For elements with tooltip or aria-describedby attributes:

```
ab('hover @eN')
ab('wait --load networkidle')
state_snapshot = abJson('snapshot -i')
// Extract new elements (tooltip content)
ab('move 0 0')  // move mouse away to reset
```

**3. Scroll triggers**: For elements with role=feed, role=list with overflow, or scrollable containers:

```
ab('scroll @eN down')
ab('wait --load networkidle')
state_snapshot = abJson('snapshot -i')
// Extract new elements (lazy-loaded content)
ab('scroll @eN up')  // scroll back to reset
```

**4. Form submission triggers**: For elements with role=form:

```
// Fill required fields with test data
ab('fill @required_field "test value"')
ab('click @submit_button')
ab('wait --load networkidle')
state_snapshot = abJson('snapshot -i')
// Extract new elements (validation results, success/error messages)
ab('press Escape')  // or navigate back to reset
```

For each trigger type:
1. Compare state snapshot with base snapshot, extract new elements
2. Record as `states` entry: `{ name, trigger: "<elementID>", elements: [...] }`
3. `trigger` references the trigger element's E-NNN ID (e.g. `"E-002"`)
4. In-state elements also receive E-NNN IDs
5. Reset to base state before exploring next trigger

> **Note**: `@eN` is agent-browser CLI's element reference syntax, used only during sitemap generation. Generated test scripts (`*.spec.ts`) must NOT use `@eN`; they must use Playwright Locator API.

```
ab('close')
```

### Step 5: Merge, Dedup & Validate

#### 5a. Element Merge

For each element (including layout and in-states elements), match against existing sitemap using the `route + role + name` triplet:

- **Match found** → preserve existing ID
- **No match** → assign new ID (current max ID + 1)
- **Existing ID has no match** → element was removed, delete from sitemap

#### 5b. Post-Collection Dedup

After merging all elements, scan for uncaught shared elements:

```
1. For each element across all pages:
   - If same role+name appears in ≥2 pages AND exists in layout.elements → already handled
   - If same role+name appears in ≥2 pages BUT NOT in layout.elements:
     → Uncaught shared element
     → Promote to layout.elements (if wrapped by layout)
     → Remove from individual page element lists

2. Report promotions: "3 elements promoted to shared: Breadcrumb (was on 5 pages), PageTitle (was on 3 pages)"
```

#### 5c. Stale Route Detection

```
For each route in the existing sitemap:
1. Check if route exists in the new route registry (Step 2a) or was discovered by crawling (Step 3)
2. If NOT found in either source:
   - Mark as "potentially stale"
   - Attempt to open the route with agent-browser
   - If 404/redirect → remove from sitemap
   - If still loads → keep with a warning in the change report
3. Report: "2 stale routes removed: /old-page, /deprecated-feature"
```

#### 5d. Validation

```
1. JSON schema validation:
   - All pages have route + title + elements (non-empty arrays)
   - All elements have id + role + name (non-empty)
   - All state triggers reference existing element IDs
   - Layout wraps reference existing page routes
   - No duplicate element IDs across the entire sitemap
   - No duplicate routes

2. Structural integrity:
   - Layout element IDs (L-NNN) are unique
   - Page element IDs (E-NNN) are unique
   - No gaps in ID numbering (orphan check)

3. Completeness check:
   - Every route from the route registry (Step 2a) has a corresponding page entry
   - Report MISSING routes: routes found in router but not explored

On validation failure: write sitemap anyway, but include validation warnings in the change report.
```

### Step 6: Write

Write to `docs/sitemap/sitemap.json`.

**Report changes:**

```
Sitemap updated: docs/sitemap/sitemap.json
  Layout: AppLayout (5 shared elements, L-001..L-005)
  3 pages, 12 elements (page-specific) + 8 elements (states)
  Routes: 12 from router, 1 additional from crawling
  +4 new elements (E-015..E-018)
  -2 removed elements (previously on /settings)
  2 stale routes removed: /old-page, /deprecated
  3 elements promoted to shared: Breadcrumb, PageTitle, BackButton
  17 unchanged

Validation: PASS / Validation: 2 WARNINGS
  ⚠ Route /admin found in router but not explored (requires special auth)
  ⚠ Route /deprecated-feature still loads but not in router
```

## Element ID Assignment Rules

- **Layout elements**: format `L-NNN`, globally unique, independent numbering space
- **Page elements**: format `E-NNN`, globally unique (base and states share ID space)
- First generation: start from `L-001` and `E-001` respectively
- Incremental update: new elements start from current max ID + 1
- IDs are never reused (deleted IDs are not recycled)
