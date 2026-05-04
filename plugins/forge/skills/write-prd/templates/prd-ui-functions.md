---
feature: "{{FEATURE_NAME}}"
---

# {{FEATURE_NAME}} — UI Functions

> Requirements layer: defines WHAT the UI must do. Not HOW it looks (that's ui-design.md).

## UI Scope

<!-- Summary of all UI surfaces this feature requires -->

## UI Function 1: {{Function Name}}

### Placement

- **Mode**: new-page | existing-page
- **Target Page**: {{page route (for existing-page) or page name (for new-page)}}
- **Position**: {{for existing-page: where in the page. For new-page: describe page purpose}}

### Description
<!-- What this UI element does -->

### User Interaction Flow
<!-- Step-by-step interaction: user clicks X → system responds with Y -->

### Data Requirements
<!-- What data this UI element needs to display or collect -->

| Field | Type | Source | Notes |
|-------|------|--------|-------|
| <!-- --> | <!-- --> | <!-- --> | <!-- --> |

### States
<!-- States this UI element can be in (loading, empty, error, populated, etc.) -->

| State | Display | Trigger |
|-------|---------|---------|
| <!-- --> | <!-- --> | <!-- --> |

### Validation Rules
<!-- Input validation, conditional display, etc. -->

---

## UI Function 2: {{Function Name}}

<!-- Repeat pattern above for each UI surface -->

---

## Page Composition

| Page | Type | UI Functions | Position Notes |
|------|------|-------------|----------------|
| {{route or name}} | new | UF-1, UF-2 | New page for {{purpose}} |
| {{route}} | existing | UF-3 | {{UF-3}} embedded {{position}} |
