---
name: gen-sitemap
description: Auto-generate and maintain sitemap.json for a web app. Uses agent-browser to explore routes, capture accessibility tree, and discover dynamic states. Preserves element IDs across runs.
allowed-tools: Bash Read Write Grep Glob
argument-hint: "[base-url] [api-base-url]"
---

# Gen Sitemap

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

If install also fails (network error, permission denied, incompatible Node version), proceed without agent-browser: skip Steps 3-4 link crawling and dynamic state exploration. Build sitemap from router registry only (Step 2a) and notify user: "agent-browser unavailable. Sitemap built from static route analysis only — dynamic states and SPA client-side routes may be incomplete."

<HARD-RULE>
Always use a version-qualified `npx agent-browser@latest` invocation. Do not use bare `npx agent-browser` (unpinned) — it may pull breaking changes between runs. For CI reproducibility, consider pinning to a specific version (e.g. `agent-browser@0.4.x`) in your environment.
</HARD-RULE>

After installation, re-run this command.

## Config Resolution

Before executing the workflow, resolve configuration sources:

1. Check if `tests/config.yaml` exists
2. **If not found**: read the template at `templates/test-config.yaml`, write it to `tests/config.yaml`, then **abort and prompt the user**:

```
Created tests/config.yaml (template). Please fill in values for your environment and re-run.
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

See `rules/schema.md` for the complete field reference, layout-vs-page element distinction, element ID assignment rules, and dynamic route parameterization rules.

## Process Flow

```
1. Load existing sitemap -> 2. Analyze layout & build route registry -> 3. Discover routes -> 4. Explore pages -> 5. Merge & dedup & validate -> 6. Write
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
   <Route element={<AppLayout />}>      <- shared layout
     <Route path="/" element={<Dashboard />} />
     <Route path="/settings" element={<Settings />} />
   </Route>
   <Route path="/login">                <- no layout (standalone page)

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
-> These are shared elements (sidebar, nav, header, footer, etc.)

Classify shared elements:
- Layout-level (sidebar, top nav, footer) -> layout.elements with L-NNN IDs
- Pattern-level (e.g. breadcrumb on multiple but not all pages) -> noted as "partial-shared"
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
3. Filter to same-origin paths:
   - Keep: paths under baseUrl's origin (compare `new URL(href).origin === new URL(baseUrl).origin`)
   - Keep: relative paths (no protocol, e.g., `/about`, `./page`)
   - Exclude: external links (different origin), `mailto:`, `tel:`, `javascript:`, `data:`, `#fragment-only`
4. Deduplicate to get crawled route list
5. Recursively extract links from each new route (breadth-first, max depth 3)
6. Merge: router routes + crawled routes = complete route list
7. Report source counts: "12 routes from router, 1 additional from crawling"
```

**Dynamic route handling**: Parameterization rules are defined in `rules/schema.md`.

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

#### Element Extraction & Dynamic State Exploration

See `rules/page-exploration.md` for element extraction rules, ARIA role filter set, and dynamic state exploration procedures (click, hover, scroll, form submission triggers).

```
ab('close')
```

### Step 5: Merge, Dedup & Validate

See `rules/merge-validation.md` for the complete rules on:
- 5a. Element merge (ID preservation/reassignment)
- 5b. Post-collection dedup (promote uncaught shared elements)
- 5c. Stale route detection
- 5d. Validation (JSON schema, structural integrity, completeness)

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
  Warning: Route /admin found in router but not explored (requires special auth)
  Warning: Route /deprecated-feature still loads but not in router
```
