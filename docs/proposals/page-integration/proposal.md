---
created: "2026-05-04"
status: proposal
tags: [architecture, pipeline, ui-integration]
---

# Page Integration Pipeline — Full Chain Enhancement

## Problem

DecisionTimeline component was fully built (24 unit tests passing) but never integrated into MainItemDetailPage. Root cause: forge pipeline has no concept of "page composition." The `prd-ui-functions.md` template uses "UI Function" as its unit with no page placement. The `ui-design.md` template uses "Component" as its unit with no page grouping. The `breakdown-tasks` skill applies a one-to-one rule (one component = one task), so "build DecisionTimeline and embed it in the detail page" becomes one task that only covers the build.

The gap exists across the entire pipeline — PRD has no structured page/component relationship, UI design has no page composition, tech design has no integration spec, task breakdown has no build/integrate split, and test generation has no integration verification.

## Design Decisions

12 decisions confirmed through adversarial review:

| # | Decision | Choice |
|---|----------|--------|
| 1 | Fix scope | Full chain: PRD → UI Design → Tech Design → Breakdown → Tests |
| 2 | Where to declare page/component relationship | PRD explicit declaration with `placement` field |
| 3 | Organization in PRD | Per-UI-Function Placement field + Page Composition summary table |
| 4 | How PRD knows existing pages | Allow write-prd to read `sitemap.json` (business asset, not tech selection) |
| 5 | ui-design responsibility boundary | Don't read code, but must specify embedding position in existing page |
| 6 | tech-design Integration Spec granularity | Minimal: target file + insertion point + data source |
| 7 | Task split trigger | PRD placement directly drives split; tech-design Integration Spec as validation supplement |
| 8 | Integration verification safeguard | Double: gen-test-cases integration test case + gate task file-diff check |
| 9 | Backward compatibility | Not compatible — mandatory Placement field |
| 10 | sitemap trust level | Trust sitemap.json, optimize gen-sitemap for correctness |
| 11 | gen-sitemap fix scope | P0+P1+P2 full hardening |
| 12 | Common element dedup | Multi-page comparison during collection + post-collection cross-page dedup as safety net |

## Detailed Changes

### A. gen-sitemap (highest priority — foundation for downstream trust)

#### A1. Route discovery: dual-source (P0)

Current: link-crawling only (breadth-first, depth 3). Misses routes not reachable by clicking links.

Change: router file parsing as primary source, link-crawling as supplement.

**Step 2a (expanded)** — Read route definitions not just for layout nesting, but to build a complete route registry:

```
Step 2a: Build Complete Route Registry

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

3. Build route registry:
   {
     "routes": [
       { "template": "/items", "source": "router", "layout": "AppLayout" },
       { "template": "/items/:id", "source": "router", "layout": "AppLayout" },
       { "template": "/login", "source": "router", "layout": null }
     ]
   }
```

**Step 3 (revised)** — Link-crawling supplements router discovery:

```
Step 3: Discover Routes (revised)

1. Start with route registry from Step 2a as the BASE
2. Link-crawl to discover additional routes not in router files (e.g., dynamically generated routes)
3. Merge: router routes + crawled routes = complete route list
4. Report: "12 routes from router, 1 additional from crawling"
```

#### A2. Element roles: full ARIA coverage (P1)

Current: 16 roles `{button, link, heading, textbox, checkbox, radio, combobox, tab, dialog, alert, navigation, search, form, menuitem, switch}`

Change: expand to full ARIA role set relevant to web app testing:

```
Added roles:
table, row, cell, columnheader, rowheader,
grid, listbox, option, list, listitem,
tooltip, progressbar, meter, slider, spinbutton,
status, log, marquee, timer,
img, separator, group, region,
feed, article, figure, caption
```

Update the extraction filter in Step 4 "Base Element Extraction" to include these roles.

#### A3. Post-generation validation (P1)

Add Step 5.5 between Merge and Write:

```
Step 5.5: Validate

1. JSON schema validation:
   - All pages have route + title + elements (non-empty arrays)
   - All elements have id + role + name (non-empty)
   - All state triggers reference existing element IDs
   - Layout wraps reference existing page routes
   - No duplicate element IDs across the entire sitemap
   - No duplicate routes

2. Structural integrity:
   - Layout element IDs (L-NNN) are unique and sequential
   - Page element IDs (E-NNN) are unique and sequential
   - No gaps in ID numbering (orphan check)

3. Completeness check:
   - Every route from the route registry (Step 2a) has a corresponding page entry
   - Report MISSING routes: routes found in router but not explored

On validation failure: write sitemap anyway, but append a "## Validation Warnings" section to the change report.
```

#### A4. Dynamic state exploration: expanded triggers (P2)

Current: only `click` on `role=button/tab/disclosure`.

Change: add hover, scroll, and form submission triggers:

```
Expanded trigger types:

1. Click triggers (existing): role=button/tab/disclosure with non-empty name
2. Hover triggers: role=tooltip (hover on parent element), elements with title/aria-describedby
3. Scroll triggers: role=feed, role=list with overflow, containers with scroll indicators
4. Form submission: role=form — fill required fields, submit, capture response state

For each trigger type:
- Trigger the interaction
- Wait for networkidle
- Snapshot and diff against base
- Record new elements as state entry
- Reset to base state
```

#### A5. Stale route detection (P2)

Add to Step 5 Merge:

```
Step 5 (addition): Stale Route Cleanup

For each route in the existing sitemap:
1. Check if route exists in the new route registry (Step 2a)
2. If NOT found in registry AND NOT found by link crawling:
   - Mark as "potentially stale"
   - Attempt to open the route with agent-browser
   - If 404/redirect: remove from sitemap
   - If still loads: keep with a warning in the change report
3. Report: "2 stale routes removed: /old-page, /deprecated-feature"
```

#### A6. Common element dedup: dual mechanism (Decision 12)

**Mechanism A: Multi-page comparison during collection (Step 2b revised)**

```
Step 2b (revised): Identify Shared Elements via Multi-Page Comparison

1. After Step 2a identifies layout-wrapped routes, visit 3-5 wrapped routes:
   ab('open <baseUrl><route_1>')
   ab('wait --load networkidle')
   snapshot_1 = abJson('snapshot -i')

   ab('open <baseUrl><route_2>')
   ab('wait --load networkidle')
   snapshot_2 = abJson('snapshot -i')

   // ... up to 5 routes

2. Compute intersection: elements present in ALL visited snapshots with same role+name
   → These are shared elements (sidebar, nav, header, footer, etc.)

3. Classify shared elements:
   - Layout-level (sidebar, top nav, footer) → layout.elements with L-NNN IDs
   - Pattern-level (breadcrumb appearing on multiple but not all pages) → note as "partial-shared"
     (still assigned to individual pages, but marked for awareness)

4. Output: layout.elements + partial-shared element list
```

**Mechanism C: Post-collection dedup (Step 5 addition)**

```
Step 5 (addition): Post-Collection Dedup

After merging all elements, scan for duplicates:

1. For each element E-NNN across all pages:
   - If same role+name appears in ≥2 pages AND exists in layout.elements → already handled
   - If same role+name appears in ≥2 pages BUT NOT in layout.elements:
     → This is an uncaught shared element
     → Promote to layout.elements (if wrapped by layout) or mark as "shared-pattern"

2. Report: "3 elements promoted to shared: Breadcrumb (was on 5 pages), PageTitle (was on 3 pages)"

3. Remove promoted elements from individual page element lists
```

---

### B. write-prd

#### B1. SKILL.md changes

**Step 1 (addition):** Allow reading sitemap.json:

```markdown
## Step 1: Explore Project Context

- Check `docs/proposals/<slug>/proposal.md` if a proposal exists
- Check `docs/features/<slug>/tasks/index.json` for related tasks
- Review recent git commits for related work
- **Read `docs/sitemap/sitemap.json`** if it exists — this is a business asset catalog
  listing all existing pages and their elements. Use it to understand what pages
  already exist when the user's requirements involve modifying existing pages.

**Forbidden**: Do not read `ARCHITECTURE.md`, `DECISIONS.md` or other technical docs.
`sitemap.json` is explicitly allowed — it documents business-level page inventory,
not technical architecture decisions.
```

**Step 8 (revised):** Add Placement guidance:

```markdown
## Step 8: Write UI Functions (if applicable)

For features with UI surfaces, create `prd/prd-ui-functions.md` using `templates/prd-ui-functions.md`.

**Placement rules** (mandatory for every UI Function):

1. Read `docs/sitemap/sitemap.json` to understand existing page inventory
2. For each UI Function, declare its Placement:
   - `new-page` — this function creates a brand new page
   - `existing-page:<route>` — this function adds UI to an existing page (route from sitemap)
3. For `existing-page`, also specify the Position within the page (e.g., "above sub-items table")
4. After all UI Functions are defined, compile the Page Composition summary table

**Validation**: Every UI Function MUST have a Placement section. Missing Placement → error, do not proceed.
```

**Step 9.5 Self-Check (addition):**

```markdown
| Placement completeness | Every UI Function has a Placement section with Mode and target |
| Placement consistency | existing-page routes exist in sitemap.json (if sitemap available) |
| Page Composition valid | Page Composition table lists all pages with correct UI Function references |
```

#### B2. Template: prd-ui-functions.md

**Current structure:**

```markdown
## UI Scope

## UI Function 1: {{Function Name}}
### Description
### User Interaction Flow
### Data Requirements
### States
### Validation Rules
```

**New structure:**

```markdown
## UI Scope

<!-- Summary of all UI surfaces this feature requires -->

## UI Function 1: {{Function Name}}

### Placement

- **Mode**: new-page | existing-page
- **Target Page**: {{page route (for existing-page) or page name (for new-page)}}
- **Position**: {{for existing-page: where in the page. For new-page: describe page purpose}}

### Description

### User Interaction Flow

### Data Requirements

| Field | Type | Source | Notes |
|-------|------|--------|-------|

### States

| State | Display | Trigger |
|-------|---------|---------|

### Validation Rules

---

<!-- Repeat for each UI Function -->

## Page Composition

| Page | Type | UI Functions | Position Notes |
|------|------|-------------|----------------|
| {{route or name}} | new | UF-1, UF-2 | New page for {{purpose}} |
| {{route}} | existing | UF-3 | {{UF-3}} embedded {{position}} |
```

---

### C. ui-design

#### C1. SKILL.md changes

**Step 2 (addition):** Consume PRD Placement:

```markdown
## Step 2: Read UI Functions

Read `prd/prd-ui-functions.md` to understand UI requirements.

**Placement awareness**: For each UI Function, note its Placement section:
- `new-page`: design the complete page including layout, navigation, and all contained components
- `existing-page:<route>`: design only the new component(s) that will be embedded.
  The Placement's Position field tells you WHERE in the existing page — this constrains
  the component's layout (e.g., width must match existing content area).
```

**Step 4 (addition):** Embedding position in output:

```markdown
## Step 4: Draft UI Design

For each UI function, define layout structure, states, interactions, and data binding.

**For existing-page UI Functions**: the Component section in ui-design.md MUST include
an "Embedding" subsection that specifies where in the existing page this component
will be placed. This comes from the PRD's Placement Position field but may be refined
with layout constraints discovered during design (e.g., "full width of the content area,
above the sub-items table, below the progress summary heading").
```

#### C2. Template: ui-design.md

**Add Embedding subsection to Component:**

```markdown
## Component: {{Component Name}}

### Placement

- **Mode**: new-page | existing-page
- **Target**: {{page route or name}}
- **Position**: {{where in the page — inherited from PRD, refined with design constraints}}

### Layout Structure

### States

| State | Visual | Behavior |
|-------|--------|----------|

### Interactions

| Trigger | Action | Feedback |
|---------|--------|----------|

### Data Binding

| UI Element | Data Field | Source |
|------------|-----------|--------|

---
```

---

### D. tech-design

#### D1. SKILL.md changes

**Add after Error Handling section:** Generate Integration Spec for existing-page components.

```markdown
## Integration Specs (when applicable)

For each UI Function with `placement: existing-page:<route>`, generate an Integration Spec
in the tech design document. Read the UI Design's Placement section for context.

The Integration Spec declares what file to modify and where:
- Do NOT specify implementation details (import statements, prop interfaces)
- Do specify: target file path, insertion point description, data source

This spec is consumed by breakdown-tasks to generate separate integration tasks.
```

#### D2. Template: tech-design.md

**Add new section between Cross-Layer Data Map and Testing Strategy:**

```markdown
## Integration Specs

<!-- Required when any UI Function has placement: existing-page.
     Skip with "No existing-page integrations — not applicable." otherwise. -->

### Integration: {{Component Name}} → {{Target Page}}

- **Target File**: {{file path of the existing page component}}
- **Insertion Point**: {{where to add the component in the page, e.g., "above SubItemsTable"}}
- **Data Source**: {{API call or data hook needed, e.g., "decision logs API by mainItemId"}}

<!-- Repeat for each existing-page integration -->
```

---

### E. breakdown-tasks

#### E1. SKILL.md changes

**Step 1 (addition):** Read sitemap for validation:

```markdown
<HAS_PLACEMENT>
If `prd/prd-ui-functions.md` exists, validate Placement:
1. Read the Page Composition table from `prd/prd-ui-functions.md`
2. For each `existing-page:<route>` entry, verify the route exists in `docs/sitemap/sitemap.json`
3. If route not found in sitemap → ERROR: abort with message "Route <route> not found in sitemap.json.
   Run /gen-sitemap first or verify the route is correct."
4. If no Placement sections found in any UI Function → ERROR: "Missing Placement declarations.
   All UI Functions must have a Placement section. Run /write-prd to update."
</HAS_PLACEMENT>
```

**Step 2 (addition):** Task split rules driven by placement:

```markdown
<HAS_PLACEMENT>
<RULE> (replaces existing one-to-one rule)

UI Task Split Rules — driven by PRD Placement:

1. For each UI Function with `placement: new-page`:
   - Create one "Build Component" task per component (existing behavior)
   - Create one "Page Assembly" task: create page file, register route, compose all components
   - Build tasks depend on interfaces + models
   - Page Assembly depends on all Build tasks for its page

2. For each UI Function with `placement: existing-page:<route>`:
   - Create one "Build Component" task (component implementation + unit tests)
   - Create one "Integrate Component" task (wire component into existing page)
   - Build task depends on interfaces + models
   - Integrate task depends on Build task
   - Integrate task's acceptance criteria MUST reference:
     a. Target page file from tech-design Integration Spec
     b. Insertion point from tech-design Integration Spec
     c. Component visible at correct position (verifiable by e2e)

3. For mixed scenarios (some new-page, some existing-page):
   - Apply rules 1 and 2 independently per UI Function

4. NO fallback to one-to-one rule. Every UI component MUST have explicit Placement.
</RULE>
</HAS_PLACEMENT>
```

**Step 2 Element Mapping (addition):**

```markdown
| Integration Spec (existing-page) | tech-design.md | Integration (UI) |
| Page composition (new-page)      | prd-ui-functions.md Page Composition | Page Assembly task |
```

**Step 4a (addition):** Integration task template reference:

```markdown
For Integration tasks, Reference Files must include:
1. tech-design.md Integration Spec section
2. ui-design.md Component Placement section
3. The target page file path (for file-diff verification)
4. Any relevant prototype file

For Page Assembly tasks, Reference Files must include:
1. prd-ui-functions.md Page Composition table
2. ui-design.md Components for this page
3. Route configuration file (for route registration)
4. Navigation component file (for adding nav links)
```

#### E2. Template: task.md

No structural changes needed — integration tasks use the same template with:
- `title`: "Integrate {{Component}} into {{Page}}" or "Assemble {{Page}}"
- `description`: references the Integration Spec or Page Composition
- `acceptance criteria`: includes integration-specific checks

---

### F. gen-test-cases

#### F1. SKILL.md changes

**Step 3 (addition):** Auto-generate integration verification test cases:

```markdown
### Integration Test Case Generation

For each UI Function with `placement: existing-page:<route>`, generate a dedicated
integration verification test case:

## TC-{NNN}: Integration — {{Component}} visible on {{Page}}
- **Source**: PRD UI Function "{{Function Name}}" Placement + Integration Spec
- **Type**: UI
- **Target**: ui/<page-name>
- **Test ID**: ui/<page-name>/integration-<component-slug>
- **Pre-conditions**: Component build complete, integration task complete
- **Route**: <route>
- **Element**: {{sitemap element IDs for the insertion point area}}
- **Steps**:
  1. Navigate to <route>
  2. Verify {{Component}} is visible at {{Position}}
  3. Verify {{Component}} renders with expected data
- **Expected**: Component appears at the specified position and displays data correctly
- **Priority**: P0

This test case MUST exist for every existing-page integration. It serves as a safety net:
if the integration task is skipped, this test will fail.
```

#### F2. Template: test-cases.md

**Add Integration row to Summary table:**

```markdown
## Summary

| Type | Count |
|------|-------|
| UI | {{UI_COUNT}} |
| **Integration** | **{{INTEGRATION_COUNT}}** |
| API | {{API_COUNT}} |
| CLI | {{CLI_COUNT}} |
| **Total** | **{{TOTAL_COUNT}}** |
```

---

### G. gate-task

#### G1. Template: gate-task.md

**Add to Verification Checklist:**

```markdown
7. [ ] All Integration Specs from tech-design.md have corresponding code changes
     (for each Integration Spec: verify target file was modified since feature branch started)
8. [ ] All integration test cases pass (if gen-test-cases already ran)
```

---

### H. manifest traceability

#### H1. Template: manifest-update-tasks.md

**Expand traceability table:**

```markdown
| PRD Section | Design Section | UI Component | Placement | Tasks |
|-------------|----------------|--------------|-----------|-------|
| "{{PRD_SECTION}}" (prd-spec §N) | "{{DESIGN_SECTION}}" (tech-design §N) | "{{UI_COMPONENT}}" (ui-design §N, or "—") | new-page / existing-page:<route> / — | {{TASK_IDS}} |
```

---

## Implementation Order

Dependencies between changes determine execution sequence:

```
Phase 1: Foundation (gen-sitemap)
├── A1. Route discovery: dual-source
├── A2. Element roles: full ARIA coverage
├── A6. Common element dedup (A+C)
└── A3. Post-generation validation
    ↓
Phase 2: PRD Layer (write-prd)
├── B2. Template: prd-ui-functions.md (Placement + Page Composition)
├── B1. SKILL.md (read sitemap, Placement guidance, self-check)
└── Verify: generate a test PRD with both new-page and existing-page scenarios
    ↓
Phase 3: Design Layer (parallel)
├── C2. Template: ui-design.md (Placement subsection)
├── C1. SKILL.md (consume Placement, specify embedding position)
├── D2. Template: tech-design.md (Integration Spec section)
└── D1. SKILL.md (generate Integration Spec)
    ↓
Phase 4: Task Breakdown Layer
├── E1. SKILL.md (placement-driven split, sitemap validation, no-placement error)
└── Verify: run breakdown on a test feature, confirm build+integrate task split
    ↓
Phase 5: Verification Layer (parallel)
├── F1/F2. gen-test-cases (integration test cases)
├── G1. gate-task (Integration Spec file-diff check)
└── H1. manifest traceability (Placement column)
    ↓
Phase 6: gen-sitemap P2 hardening (can run in parallel with Phase 2+)
├── A4. Dynamic state exploration expansion
└── A5. Stale route detection
    ↓
Phase 7: Integration Test
└── End-to-end test: full pipeline with Decision Log feature
    (write-prd → ui-design → tech-design → breakdown-tasks → verify task split)
```

## Validation Criteria

The enhancement is complete when:

1. **gen-sitemap**: All routes from router files appear in sitemap.json; no duplicate shared elements; validation step passes
2. **write-prd**: Every generated `prd-ui-functions.md` has Placement on all UI Functions + Page Composition table
3. **ui-design**: Every component in `ui-design.md` has Placement subsection; existing-page components have explicit embedding position
4. **tech-design**: Every existing-page placement has a corresponding Integration Spec
5. **breakdown-tasks**: `existing-page` placements produce separate build + integrate tasks; `new-page` placements produce build + page assembly tasks; missing Placement causes error
6. **gen-test-cases**: Every existing-page integration has a P0 integration verification test case
7. **gate-task**: Checks Integration Spec coverage via file diff
8. **E2E validation**: Decision Log feature (from pm-work-tracker lesson) runs through the full pipeline and produces correct build + integrate tasks

## Appendix: Decision Log

Full traceability of the 12 design decisions and their rationale:

1. **Full chain** — Partial fixes leave gaps. The DecisionTimeline problem spans PRD → breakdown, so the fix must too.
2. **PRD explicit declaration** — "Component exists on page" is a business intent, not a design or implementation detail. The PRD author knows this intent earliest.
3. **Per-function + summary table** — Downstream skills consume differently: breakdown needs per-component placement, test generation needs page-level view. Both serve different readers.
4. **Read sitemap.json** — "App has a detail page for main items" is business knowledge, not tech selection. Comparable to knowing "we have a user management module."
5. **ui-design no code, must specify position** — UI design is about appearance and interaction, not code wiring. But position affects layout constraints (width, surrounding elements).
6. **Minimal Integration Spec** — Tech design declares "what needs to change" (file, position, data), not "how to code it." Executors retain implementation flexibility.
7. **PRD placement drives split** — PRD is the earliest and clearest intent signal. Tech-design Integration Spec is derived from PRD; if PRD says existing-page but tech-design misses it, that's a tech-design bug to catch.
8. **Double verification** — Test cases catch missing integration at runtime; gate tasks catch it structurally (file diff) before tests run.
9. **Not backward compatible** — The entire point is preventing a specific failure mode. Allowing old format to silently skip protection defeats the purpose.
10. **Trust sitemap** — Sitemap is the single source of truth for page inventory. Rather than adding alternative page registries, invest in making sitemap correct.
11. **Full gen-sitemap hardening** — P0 (route discovery) and P1 (elements, validation) directly affect sitemap reliability as PRD validation source. P2 (states, stale cleanup) improves overall quality.
12. **A+C dedup** — Multi-page comparison catches shared elements during collection (efficient); post-collection dedup catches any that leak through (safety net).
