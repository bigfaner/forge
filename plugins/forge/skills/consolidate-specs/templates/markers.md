# Integration Marker Templates

## Standard Integration Marker

Write `docs/features/{{SLUG}}/specs/.integrated`:

```yaml
feature: "{{SLUG}}"
integrated: "{{DATE}}"
biz_count: {{BIZ_COUNT}}
tech_count: {{TECH_COUNT}}
replaced:
  - decisions/error-handling.md row "Adopt AIError..." -> BIZ-auth-003
  - lessons/gotcha-error-handling.md -> TECH-error-005
```

The `replaced` field is omitted if no overlaps were resolved.

## Early-Exit Marker (all LOCAL)

```yaml
feature: "{{SLUG}}"
integrated: "{{DATE}}"
status: "skipped: all local"
biz_count: 0
tech_count: 0
```
