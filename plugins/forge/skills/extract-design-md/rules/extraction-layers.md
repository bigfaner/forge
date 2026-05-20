# Visual Style Extraction Strategy

Extract layer by layer -- **stop at the first successful layer, no need to proceed further**:

## Layer 1: Trace CSS bundle

Use `WebFetch` to get page HTML, extract CSS file URLs from `<link rel="stylesheet">` (React build output is typically `/static/css/main.xxxxxx.css`), then fetch that CSS file.

The CSS bundle contains complete style rules -- the most direct information source.

## Layer 2: Extract CSS custom properties (design tokens)

Search for `:root` blocks in the CSS bundle. Modern design systems almost always use CSS variables for tokens:

```css
:root {
  --color-primary: #635bff;
  --font-size-base: 16px;
  --radius-md: 8px;
}
```

These variables map directly to DESIGN.md fields -- the highest quality information source.

## Layer 3: Multi-page sampling

A single page may only have landing content, missing form, card, table component styles. Fetch additional paths (if publicly accessible):

- `/login` -- inputs, buttons
- `/dashboard` or `/app` -- cards, navigation, data display
- `/settings` -- forms, toggles, grouped sections

Compare multi-page extraction results to fill in missing component specs.

## Layer 4: agent-browser runtime extraction (local apps)

If the target is a locally running SPA (e.g. `http://localhost:3000`), use agent-browser to execute JS for runtime computed styles:

```
ab('open <url>')
ab('wait --load networkidle')
// Extract CSS custom properties
ab('eval document.documentElement getComputedStyle --color-* --font-* --radius-* variables')
// Extract key component computed styles
ab('eval find button[class*="primary"] get backgroundColor borderRadius padding')
ab('eval find [class*="card"] get background border boxShadow borderRadius')
```

## Layer 5: Visual inference (fallback)

When none of the above layers can obtain specific values, infer from page visual descriptions, screenshots, or HTML structure.

Mark uncertain values with `(estimated)` in the output.
