# Mobile App Platform Rules

Navigation patterns for mobile apps (native-like HTML, mini-programs, Capacitor/WKWebView wrappers).

## Primary Navigation

Bottom tab bar:

```html
<!-- Shared tab bar — identical across all tab pages, only active class moves -->
<nav class="tab-bar">
  <a href="home.html" class="tab-item active">
    <span class="tab-icon">{{icon}}</span>
    <span class="tab-label">Home</span>
  </a>
  <a href="stats.html" class="tab-item">
    <span class="tab-icon">{{icon}}</span>
    <span class="tab-label">Stats</span>
  </a>
</nav>
```

Rules:
- Follow the Platform-Agnostic Rules in prototype.md Navigation Contract
- Fixed to bottom of viewport, always visible on tab pages
- Tab pages should NOT display a top navigation bar (redundant with tab bar)

## Secondary Pages

Header with back button (push-style navigation):

```html
<header class="page-header">
  <a href="parent-page.html" class="back-btn" aria-label="Back">
    <span class="icon-chevron-left"></span>
  </a>
  <h1>{{Page Title}}</h1>
</header>
```

Rules:
- Every secondary page must have a back button in the top-left
- Back button href must match the entry page defined in Navigation Architecture
- No browser URL bar — all navigation must be explicit (tabs, buttons, links)

## Page Structure

```html
<body>
  <header class="page-header"><!-- back button + title --></header>
  <main class="page-content"><!-- scrollable content --></main>
  <!-- Tab pages only: -->
  <nav class="tab-bar"><!-- shared tab bar --></nav>
</body>
```
