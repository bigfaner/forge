# Decision Entry Template

Single decision row format for appending to `docs/decisions/{{TYPE}}.md`.

## Table Row

| Date | Feature | Decision | Rationale | Source |
|------|---------|----------|-----------|--------|
| {{DATE}} | {{FEATURE_SLUG}} | {{DECISION}} | {{RATIONALE}} | {{SOURCE}} |

## Field Constraints

| Field | Format | Constraint |
|-------|--------|------------|
| Date | ISO 8601 | e.g. `2026-04-23` |
| Feature | slug | e.g. `feat-log-decisions` |
| Decision | string | One sentence, max 80 chars |
| Rationale | string | One sentence, max 80 chars |
| Source | string | `{{FEATURE_SLUG}}/{{FILE}}.md §{{SECTION}}` or `manual` |
