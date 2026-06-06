---
created: "{{DATE}}"
related: design/tech-design.md
---

# CLI Handbook: {{FEATURE_NAME}}

## Command Overview

<!-- High-level CLI design summary -->

## Commands

### {{Command Name}}

**Usage**: `{{command}} {{subcommand}} [flags] [args]`
**Alias**: <!-- Comma-separated aliases, or "none" -->
**Auth**: <!-- Required role or "none" -->

#### Subcommands

| Subcommand | Description | Default Flags |
|------------|-------------|---------------|
| <!-- -->   | <!-- -->    | <!-- -->      |

#### Flags

| Flag | Short | Type | Required | Default | Description |
|------|-------|------|----------|---------|-------------|
| <!-- --> | <!-- --> | <!-- --> | <!-- --> | <!-- --> | <!-- --> |

#### Arguments

| Argument | Type | Required | Description |
|----------|------|----------|-------------|
| <!-- --> | <!-- --> | <!-- --> | <!-- --> |

#### Exit Codes

| Code | Description |
|------|-------------|
| 0 | Success |
| 1 | General error |
| <!-- --> | <!-- --> |

#### Examples

```bash
# {{Example description}}
{{command}} {{subcommand}} --flag value
```

---

<!-- Repeat for each command -->

## Shared Flags

| Flag | Short | Type | Description |
|------|-------|------|-------------|
| <!-- --> | <!-- --> | <!-- --> | <!-- --> |

## Environment Variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| <!-- --> | <!-- --> | <!-- --> | <!-- --> |
