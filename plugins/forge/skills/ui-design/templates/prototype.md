# Prototype Generation Guide

Generate a multi-file HTML prototype from the UI design specification.

## Rules

1. **Multi-file structure** — shared CSS/JS + one HTML per page + index as navigation hub
2. **No frameworks** — vanilla HTML/CSS/JS only
3. **No build step** — open files directly in browser, all paths relative
4. **Responsive** — works on desktop (1200px+) and mobile (<768px)
5. **Interactive** — navigation between pages, modals, dropdowns, tabs must be functional
6. **Realistic content** — use realistic text, not lorem ipsum
7. **Self-documenting** — HTML structure mirrors the design spec section names

## File Structure

```
ui/prototype/
├── index.html          # Navigation hub — links to all pages with thumbnails
├── styles.css          # Shared: reset, CSS variables, layout, components, states, responsive
├── app.js              # Shared: nav toggle, modals, dropdowns, tabs, toasts
├── {page-name}.html    # One file per page from ui-design.md
└── ...
```

### index.html (Navigation Hub)

Sitemap-style page listing all prototype pages. Each entry shows:
- Page name (matches ui-design.md component name)
- Brief description of the page purpose
- Clickable link to the HTML file
- Optional: embedded screenshot/thumbnail placeholder

```html
<!-- Entry pattern -->
<a href="dashboard.html" class="page-link">
  <div class="page-card">
    <div class="page-preview"><!-- placeholder --></div>
    <h3>Dashboard</h3>
    <p>Main overview with stats, recent activity, and quick actions</p>
  </div>
</a>
```

### styles.css (Shared Styles)

```css
/* 1. Reset */
/* 2. CSS Variables from chosen design style */
/* 3. Typography */
/* 4. Layout (grid, container, spacing) */
/* 5. Navigation */
/* 6. Components (buttons, cards, inputs, badges, tables) */
/* 7. States (hover, focus, loading, empty, error) */
/* 8. Utility classes */
/* 9. Responsive overrides */
```

#### CSS Variables

Extract from the chosen design style:

```css
:root {
  --color-bg: ...;
  --color-surface: ...;
  --color-border: ...;
  --color-text-primary: ...;
  --color-text-secondary: ...;
  --color-accent: ...;
  --font-body: ...;
  --font-heading: ...;
  --font-mono: ...;
  --radius-sm: ...;
  --radius-md: ...;
  --radius-lg: ...;
  --shadow-sm: ...;
  --shadow-md: ...;
}
```

### app.js (Shared Interactions)

Shared behaviors loaded by every page:

```js
// Mobile nav toggle
// Modal open/close (backdrop click + Escape)
// Dropdown toggle (click-outside to close)
// Tab switching
// Toast notifications (auto-dismiss 3s)
// Form validation with inline errors
// Active page highlight in navigation
```

### Page HTML Pattern

Each `{page-name}.html` follows this structure:

```html
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>{{Page Name}} — {{Feature Name}}</title>
  <link rel="stylesheet" href="https://fonts.googleapis.com/css2?family=...">
  <link rel="stylesheet" href="styles.css">
</head>
<body>
  <nav><!-- Shared navigation with active page highlighted --></nav>
  <main>
    <!-- Page-specific content from ui-design.md -->
  </main>
  <script src="app.js"></script>
</body>
</html>
```

## State Mocks

For each page, implement all relevant states as sections or toggleable views:
- **Loading**: skeleton pulse animation or spinner
- **Empty**: illustration placeholder + helpful message
- **Error**: red-tinted alert with retry action
- **Populated**: realistic data with proper formatting

## Output

Save to: `docs/features/<slug>/ui/prototype/`
