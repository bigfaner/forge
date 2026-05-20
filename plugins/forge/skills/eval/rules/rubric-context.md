# Rubric Context Frontmatter

Rubrics may declare a `context` frontmatter field to inject project reality files into the scorer prompt. Rubrics without `context` continue to work unchanged.

```yaml
context:
  conventions: [api, naming, ux]  # list of category strings (optional)
  business-rules: auto            # "auto" or list of filenames (optional)
```

| Sub-field | Type | Description |
|-----------|------|-------------|
| `conventions` | list of strings | Each string matches filenames in `docs/conventions/` by prefix. E.g., `api` matches `api*.md`. Non-matching strings are skipped silently. |
| `business-rules` | `"auto"` or list of strings | `auto` loads all `.md` files from `docs/business-rules/`. A list specifies exact filenames. Missing files are skipped silently. |

At least one sub-field must be present for context injection to activate.
