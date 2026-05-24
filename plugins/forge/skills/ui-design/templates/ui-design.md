---
created: "{{DATE}}"
source: prd/prd-ui-functions.md
status: Draft
---

# UI Design: {{FEATURE_NAME}}

## Design System

<!-- Reference to existing design system or component library -->

## Component: {{Component Name}}

### Placement

- **Mode**: new-page | existing-page
- **Target**: {{page route or name}}
- **Position**: {{where in the page — inherited from PRD, refined with design constraints}}

### Layout Structure
<!-- Component hierarchy, grid/flex layout description -->

### States

| State | Visual | Behavior |
|-------|--------|----------|
| Default | <!-- --> | <!-- --> |
| Loading | <!-- --> | <!-- --> |
| Empty | <!-- --> | <!-- --> |
| Error | <!-- --> | <!-- --> |

### Interactions

| Trigger | Action | Feedback |
|---------|--------|----------|
| <!-- --> | <!-- --> | <!-- --> |

### Data Binding

| UI Element | Data Field | Source |
|------------|-----------|--------|
| <!-- --> | <!-- --> | <!-- --> |

---

<!-- Repeat for each component -->

## TUI Component: {{Panel Name}}

> This section is for TUI panels only. Each panel MUST include all 5 mandatory structural requirements.

### Panel Placement (View + Position)

- **View**: {{view name — e.g., "All", "Detail"}}
- **Position**: {{position — e.g., "Center, main area"}}
- **Size Hint**: {{size — e.g., "Fills remaining space"}}

### ASCII Layout Mockup

<!-- Use box-drawing characters (Modern Dark theme) or ASCII characters (Minimal ASCII theme).
     Show the exact visual structure of this panel with sample data. -->

```
┌─ {{Panel Title}} ────────────────────────────┐
│                                                │
│  {{mockup content with real sample data}}      │
│                                                │
└────────────────────────────────────────────────┘
```

### Dimensions

<!-- Concrete numeric values. No "approximately" or "appropriate". -->

| Element | Value | Formula / Notes |
|---------|-------|-----------------|
| Panel width | {{N}} chars | viewport - 2 (borders) |
| Content area | {{N}} chars | panel width - 4 (border + padding) |
| Column widths | {{N}} chars each | explicit character counts |
| Bar/chart max width | {{N}} chars | explicit character count |
| String truncation maxLen | {{N}} chars | formula for path/text truncation |

### Character Palette

<!-- Every visual element MUST specify its Unicode character with code point. No "TBD". -->

| Element | Character | Unicode | Reason |
|---------|-----------|---------|--------|
| {{element}} | {{char}} | {{U+XXXX}} | {{why this character was chosen}} |

### Color Mapping

<!-- Foreground/background from the selected TUI theme palette. -->

| Element | Character | Foreground | Background |
|---------|-----------|------------|------------|
| {{element}} | {{char}} | {{color #}} | {{color # or "-"}} |

### Edge Cases

<!-- 5 mandatory scenarios. Feature with CJK text adds scenario 6. -->

| # | Scenario | Expected |
|---|----------|----------|
| 1 | Narrow terminal (80x24) | Layout does not overflow, charts/tables scale down |
| 2 | Wide terminal (140+ col) | Layout does not distort, charts max at specified width |
| 3 | Mixed numeric widths (1 vs 100) | Right-pad to widest value, column-aligned |
| 4 | Long paths/strings (>50 chars) | Truncate with ".../" preserving trailing segment |
| 5 | No data | Centered "No data" placeholder |

### States

| State | Visual | Behavior |
|-------|--------|----------|
| Default | <!-- --> | <!-- --> |
| Loading | <!-- --> | <!-- --> |
| Empty | <!-- --> | <!-- --> |
| Error | <!-- --> | <!-- --> |

### Key Bindings

| Key | Action | Context |
|-----|--------|---------|
| <!-- --> | <!-- --> | <!-- --> |

### Data Binding

| UI Element | Data Field | Source |
|------------|-----------|--------|
| <!-- --> | <!-- --> | <!-- --> |

---

<!-- Repeat TUI Component section for each TUI panel -->
