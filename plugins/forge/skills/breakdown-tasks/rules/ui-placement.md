# UI Placement Rules

**Load condition**: load this file IF `ui/ui-design.md` exists OR `prd/prd-ui-functions.md` exists.

**Guard clause**: if the referenced artifact (`ui/ui-design.md` or `prd/prd-ui-functions.md`) has no parseable content, skip this rule and proceed.

## UI Element Mapping

When this rule is loaded (i.e., `ui/ui-design.md` OR `prd/prd-ui-functions.md` exists), add these rows to the Element Mapping table in Step 2:

| Design Element | Source | Task Type |
|---|---|---|
| UI Component (Layout + States + Interactions + Binding) | `ui/ui-design.md` | Implementation (UI) |
| Integration Spec (existing-page) | `design/tech-design.md` | Integration (UI) |
| Page composition (new-page) | `prd/prd-ui-functions.md` Page Composition | Page Assembly task |

## UI Prototype Reading

If `ui/ui-design.md` exists, also list `ui/prototype/` files and read `ui/prototype/index.html` for page inventory. Skip this step if no prototype directory exists.

## Placement Validation

This procedure activates when `prd/prd-ui-functions.md` exists.

1. Read the Page Composition table from `prd/prd-ui-functions.md`
2. Check if `docs/sitemap/sitemap.json` exists. If not, WARN: `"sitemap.json not found — cannot verify existing-page routes. Run /gen-sitemap for full validation."` and proceed without route verification (skip step 3).
3. For each `existing-page:<route>` entry, verify the route exists in `docs/sitemap/sitemap.json`
4. If route not found in sitemap, ERROR: abort with message `"Route <route> not found in sitemap.json. Run /gen-sitemap first or verify the route is correct."`
5. If no Placement sections found in any UI Function, ERROR: `"Missing Placement declarations. All UI Functions must have a Placement section. Edit prd/prd-ui-functions.md to add Placement sections, or re-run /write-prd."`

## UI Task Split Rules

These rules drive how UI tasks are generated based on PRD Placement values.

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

**Placement format note**: the PRD template stores Placement as two separate fields (`Mode: new-page | existing-page` and `Target Page: <page route or name>`). Downstream consumers use the combined canonical form: `<mode>:<target-page-value>`. For example, if Mode is `existing-page` and Target Page is `/dashboard`, the canonical placement is `existing-page:/dashboard`. If Mode is `new-page` and Target Page is `Analytics`, the canonical placement is `new-page:Analytics`.

## UI Dependency Layer

UI components form a natural dependency layer after data models and interfaces, but do NOT require backend implementation to be complete (can mock).

In Step 3 (Derive Phases & Dependencies):
- UI tasks depend on interfaces + models only (can mock backend; does NOT need the backend implementation phase)

## UI Reference File Requirements

### Build Component tasks

Reference Files must include:

1. Matching `ui/ui-design.md` Component section
2. Corresponding `ui/prototype/<page>.html` (or note "No HTML prototype available")
3. Data binding table for this component
4. Relevant `design/tech-design.md` interfaces

Example:

```
- ui/ui-design.md Component "Dashboard" — layout, states, interactions
- ui/prototype/dashboard.html — interactive prototype
- design/tech-design.md Interfaces — data contracts
```

### Integration tasks (existing-page)

Reference Files must include:

1. `design/tech-design.md` Integration Spec section
2. `ui/ui-design.md` Component Placement section
3. Target page file path (for file-diff verification)
4. Any relevant prototype file

### Page Assembly tasks (new-page)

Reference Files must include:

1. `prd/prd-ui-functions.md` Page Composition table
2. `ui/ui-design.md` Components for this page
3. Route configuration file (for route registration)
4. Navigation component file (for adding nav links)

## Maintenance Note

This rule file depends on the following sections in the skeleton SKILL.md:

- **Step 2: Map -> Tasks** — Element Mapping table (adds UI-specific rows)
- **Step 3: Derive Phases & Dependencies** — dependency layering (UI layer rules)
- **Step 4a: Create Task Files** — Reference Files population per task type

If any of these sections change in the skeleton, verify that the rules in this file remain consistent.
