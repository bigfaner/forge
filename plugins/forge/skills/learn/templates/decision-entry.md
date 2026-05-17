# Decision Entry Template

Used by `/learn` when writing to `docs/decisions/<type>.md`.

## Row Format

Append a single table row to the end of the target type file:

```
| YYYY-MM-DD | <feature-slug> | <Decision, one sentence> | <Rationale, one sentence> | <source> |
```

### Field Constraints

| Field | Format | Constraint |
|-------|--------|------------|
| Date | YYYY-MM-DD | ISO 8601 |
| Feature | slug or `-` | Feature slug, e.g. `feat-log-decisions`; use `-` if unknown |
| Decision | single sentence | Max 80 characters |
| Rationale | single sentence | Max 80 characters |
| Source | file path or `manual` | `<feature-slug>/<file>.md §<Section>` or `manual` |

## Manifest Update

After writing the decision row, update `docs/decisions/manifest.md`:

1. **Categories table**: Find the row matching the decision type. Increment `Decisions` count by 1. Set `Last Updated` to today.
2. **Recent Decisions table**: Insert a new row immediately below the table header (newest first). Keep max 10 rows; remove oldest if count exceeds 10.

Recent row format:

```
| YYYY-MM-DD | <feature-slug> | <Type Name> | <Decision, one sentence> | <source> |
```

## Type File Initial State

If the target type file does not exist, create it with:

```markdown
# <Type Name> Decisions

| Date | Feature | Decision | Rationale | Source |
|------|---------|----------|-----------|--------|
```

If the header row is missing (file corrupted or empty), prepend the standard header before appending the new row.

## Directory Bootstrap

If `docs/decisions/` does not exist, auto-create the directory plus all 8 type files and `manifest.md` from their initial templates.
