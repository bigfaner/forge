# Merge, Dedup & Validation Rules

## Element Merge (Step 5a)

For each element (including layout and in-states elements), match against existing sitemap using the `route + role + name` triplet:

- **Match found** -> preserve existing ID
- **No match** -> assign new ID (current max ID + 1)
- **Existing ID has no match** -> element was removed, delete from sitemap

## Post-Collection Dedup (Step 5b)

After merging all elements, scan for uncaught shared elements:

```
1. For each element across all pages:
   - If same role+name appears in >=2 pages AND exists in layout.elements -> already handled
   - If same role+name appears in >=2 pages BUT NOT in layout.elements:
     -> Uncaught shared element
     -> Promote to layout.elements (if wrapped by layout)
     -> Remove from individual page element lists

2. Report promotions: "3 elements promoted to shared: Breadcrumb (was on 5 pages), PageTitle (was on 3 pages)"
```

## Stale Route Detection (Step 5c)

```
For each route in the existing sitemap:
1. Check if route exists in the new route registry (Step 2a) or was discovered by crawling (Step 3)
2. If NOT found in either source:
   - Mark as "potentially stale"
   - Attempt to open the route with agent-browser
   - If 404/redirect -> remove from sitemap
   - If still loads -> keep with a warning in the change report
3. Report: "2 stale routes removed: /old-page, /deprecated-feature"
```

## Validation (Step 5d)

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
