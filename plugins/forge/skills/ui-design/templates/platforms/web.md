# Web Platform Rules

Navigation patterns for desktop and mobile-web browsers.

## Primary Navigation

Top navigation bar (header) or sidebar:

```html
<!-- Top nav bar pattern -->
<nav class="nav-bar">
  <a href="index.html" class="nav-brand">{{App Name}}</a>
  <div class="nav-links">
    <a href="page-a.html" class="active">Page A</a>
    <a href="page-b.html">Page B</a>
  </div>
</nav>
```

Rules:
- Follow the Platform-Agnostic Rules in prototype.md Navigation Contract
- Responsive: collapse to hamburger menu on mobile viewport (<768px)

## Secondary Pages

Breadcrumbs + standard `<a href>` links:

```html
<nav class="breadcrumb">
  <a href="dashboard.html">Dashboard</a> / <span>Detail</span>
</nav>
```

Rules:
- Secondary pages use breadcrumb trail, not explicit back buttons
- Browser back button is the primary "return" mechanism
- All `<a href>` targets must correspond to existing prototype files
