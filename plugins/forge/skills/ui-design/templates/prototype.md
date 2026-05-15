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

## Navigation Contract

Before generating any page, load the platform-specific navigation rules:

1. Read the `## Navigation Architecture` section from `prd-ui-functions.md`
2. Identify the target platform (see Platform Identification below)
3. Read the corresponding platform file and apply its navigation patterns

### Platform Identification

Determine platform by checking in order:

| Signal | Web | Mobile App | TUI |
|--------|-----|------------|-----|
| PRD Navigation Architecture `platform` field | `web` | `mobile` | `tui` |
| UI Function descriptions | "pages", "routes", "browser" | "screens", "tabs", "mini-program" | "panels", "terminal", "key bindings" |
| User explicit instruction | — | — | — |

If ambiguous, ask the user.

### Platform Reference Files

| Platform | File | Navigation Style |
|----------|------|-----------------|
| Web | `templates/platforms/web.md` | Top nav bar / sidebar, breadcrumbs |
| Mobile App | `templates/platforms/mobile.md` | Bottom tab bar, secondary pages with back button |
| TUI | `templates/platforms/tui.md` | Keyboard-driven panels, modes, keymaps |

### Platform-Agnostic Rules

These rules apply regardless of platform:

- Primary navigation HTML must be byte-identical across all pages that share that navigation type (e.g., all tab-bar pages share one tab bar, all top-nav pages share one top nav)
- All navigation labels and targets must exactly match the PRD's Navigation Architecture table
- No href may point to a non-existent file

## Code Layer Separation

| Layer | File | Contents | Prohibited |
|-------|------|----------|------------|
| Shared | app.js | Nav highlight, toast, modal, dropdown, accordion, slider, generic form validation | Any function that only one page uses |
| Page-specific | inline `<script>` in each HTML | Page data, page-specific event handlers, page-specific DOM manipulation | Re-implementing shared behaviors |

Rules:
- If a function is used by only ONE page, it goes in that page's inline script, NOT in app.js
- Shared functions must not accept page-specific parameters that only one page provides
- When in doubt, put it in the page's inline script

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
// Tab/segmented switching (generic, not page-specific)
// Toast notifications (auto-dismiss 3s)
// Form validation with inline errors
// Active page highlight in navigation
//
// PROHIBITED in app.js:
// - Page-specific data objects (e.g., dayData, exerciseData)
// - Page-specific event handlers (e.g., selectDay, toggleEditMode)
// - Any function referenced by onclick in only one HTML file
// - DOM queries for elements that only exist on one page
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
  <script>
    // Page-specific logic here (loads after app.js)
  </script>
</body>
</html>
```

## TUI Prototype Rules

TUI prototypes use HTML + CSS to simulate terminal appearance in the browser. They are a **human review tool** — the ASCII mockup with numeric dimensions in `ui-design.md` remains the precise specification for agents.

### Single-File Structure

Unlike web/mobile multi-file prototypes, TUI prototypes are a **single `index.html`** containing all panels rendered inside a terminal-window container.

```
ui/prototype/          (single TUI feature)
└── index.html         # All panels in one terminal-window div

ui/prototype/tui/      (multi-platform feature)
└── index.html
```

### Output Path

| Scenario | Path |
|----------|------|
| Single TUI feature | `docs/features/<slug>/ui/prototype/index.html` |
| Multi-platform feature (web + tui) | `docs/features/<slug>/ui/prototype/tui/index.html` |

### Terminal Window Container

All panels render inside a single `<div class="terminal-window">` that simulates the terminal viewport:

```html
<div class="terminal-window">
  <div class="terminal-header">
    <span class="terminal-title">{{Feature Name}}</span>
  </div>
  <div class="terminal-content">
    <!-- Header panel -->
    <div class="tui-panel tui-header">...</div>
    <!-- Content panel(s) rendered per ASCII mockup -->
    <div class="tui-panel tui-content">...</div>
    <!-- Status bar -->
    <div class="tui-panel tui-status-bar">...</div>
  </div>
  <div class="simulated-keys">
    <button class="key-btn" data-panel="header">[Tab]</button>
    <button class="key-btn" data-panel="content">[1]</button>
    <button class="key-btn" data-panel="sidebar">[2]</button>
    <button class="key-btn" data-panel="detail">[3]</button>
    <button class="key-btn" data-action="quit">[q]</button>
    <button class="key-btn" data-action="command">[:command]</button>
  </div>
</div>
```

### Terminal CSS

Use monospace font and dark background to approximate terminal appearance. Map the TUI theme's 256-color or 16-color values to CSS hex equivalents.

```css
.terminal-window {
  background: #1e1e1e;           /* Approximate terminal bg */
  color: #d4d4d4;                /* Approximate terminal fg */
  font-family: 'Cascadia Mono', 'Consolas', 'Menlo', 'Monaco', 'Courier New', monospace;
  font-size: 14px;
  line-height: 1.4;
  padding: 0;
  border-radius: 8px;
  overflow: hidden;
  max-width: 960px;              /* Approximate 120-col terminal */
  margin: 0 auto;
}

.terminal-content {
  padding: 0;
  white-space: pre;              /* Preserve ASCII layout spacing */
}

.tui-panel {
  border: 1px solid #444;
}

.tui-panel.focused {
  border-color: #5f87af;         /* Border Focus approximation */
}

.tui-status-bar {
  background: #1e1e1e;
  color: #808080;
  padding: 0 4px;
}
```

### Panel Rendering

Each panel from `ui-design.md` maps to a `<div class="tui-panel">`. Panel layout inside the terminal-window must match the ASCII mockup structure:

1. Read the `### ASCII Layout Mockup` from each TUI panel in `ui-design.md`
2. Render the ASCII art as-is inside a `<pre>` block within the panel div
3. Apply the panel's `### Color Mapping` as inline styles or CSS classes
4. Use the `### Character Palette` characters exactly as specified

```html
<!-- Panel rendered from ASCII mockup -->
<div class="tui-panel tui-content focused">
  <pre class="tui-ascii">
┌─ Dashboard ────────────────────────────────┐
│  CPU  [████████░░] 78%    Mem [██████░░] 62%│
│  Disk [████░░░░░░] 41%    Net [██░░░░░░] 23%│
└────────────────────────────────────────────┘</pre>
</div>
```

### Simulated Key Buttons

Interactive buttons at the bottom simulate keyboard input. Each button switches the focused panel or triggers a mode change.

**Required buttons** (from proposal D6):

| Button | Behavior |
|--------|----------|
| `[Tab]` | Cycle focus to next panel |
| `[1]`-`[9]` | Jump to panel by index |
| `[q]` | Simulate quit / exit mode |
| `[:command]` | Show command input overlay |

```css
.simulated-keys {
  display: flex;
  gap: 8px;
  padding: 8px;
  background: #2d2d2d;
  justify-content: center;
}

.key-btn {
  background: #3c3c3c;
  color: #cccccc;
  border: 1px solid #555;
  border-radius: 4px;
  padding: 4px 10px;
  font-family: monospace;
  font-size: 13px;
  cursor: pointer;
}

.key-btn:hover {
  background: #505050;
  border-color: #5f87af;
}
```

```js
// Simulated key interaction (inline script)
document.querySelectorAll('.key-btn').forEach(btn => {
  btn.addEventListener('click', () => {
    const panel = btn.dataset.panel;
    const action = btn.dataset.action;
    // Remove focused class from all panels
    document.querySelectorAll('.tui-panel').forEach(p => p.classList.remove('focused'));
    // Apply focus to target panel
    if (panel) {
      const target = document.querySelector(`.tui-panel[data-name="${panel}"]`)
        || document.querySelectorAll('.tui-panel')[parseInt(panel) - 1];
      if (target) target.classList.add('focused');
    }
  });
});
```

### Panel Layout Matching

The HTML panel arrangement must match the ASCII mockup from `ui-design.md`:

- **Header**: top, full width
- **Content**: center, fills remaining space
- **Status Bar**: bottom, full width
- **Sidebar**: left side (if present in design)
- **Detail**: right of sidebar (if present in design)

Verify by comparing the rendered browser output against the ASCII mockup structure.

### TUI-Specific Post-Generation Checks

- [ ] All panels from `ui-design.md` are rendered in a single `index.html`
- [ ] Terminal-window div has dark background and monospace font
- [ ] Simulated key buttons are present: `[Tab]`, `[1]`-`[9]`, `[q]`, `[:command]`
- [ ] Panel focus toggles when clicking simulated key buttons
- [ ] ASCII art inside panels matches the mockup from `ui-design.md`
- [ ] Color mapping from the TUI theme is applied as CSS

## State Mocks

For each page, implement all relevant states as sections or toggleable views:
- **Loading**: skeleton pulse animation or spinner
- **Empty**: illustration placeholder + helpful message
- **Error**: red-tinted alert with retry action
- **Populated**: realistic data with proper formatting

## Post-Generation Verification

After generating ALL files, perform these checks:

### Navigation Consistency
- [ ] All pages sharing primary navigation have identical nav HTML (except active state)
- [ ] All Secondary Pages provide a way to return to their parent
- [ ] No navigation href points to a non-existent file

### Code Separation
- [ ] app.js does NOT contain any function that is only called from one page
- [ ] No function name collision between app.js and any inline script
- [ ] Inline scripts load AFTER app.js (or use unique function names)

### Cross-Page Links
- [ ] Every `<a href="...">` in every file points to an existing file
- [ ] Every `onclick="location.href='...'"` points to an existing file
- [ ] index.html links to all generated pages

## Output

| Scenario | Path |
|----------|------|
| Single web or mobile feature | `docs/features/<slug>/ui/prototype/` (multi-file) |
| Single TUI feature | `docs/features/<slug>/ui/prototype/index.html` (single file) |
| Multi-platform feature | `docs/features/<slug>/ui/prototype/web/`, `docs/features/<slug>/ui/prototype/tui/` |
