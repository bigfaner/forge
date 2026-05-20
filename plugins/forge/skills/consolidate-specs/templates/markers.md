# Integration Marker Templates

## Standard Integration Marker

Write `docs/features/<slug>/specs/.integrated`:

```yaml
feature: "<slug>"
integrated: "<date>"
biz_count: <N>
tech_count: <M>
replaced:
  - decisions/error-handling.md row "Adopt AIError..." -> BIZ-auth-003
  - lessons/gotcha-error-handling.md -> TECH-error-005
```

The `replaced` field is omitted if no overlaps were resolved.

## Early-Exit Marker (all LOCAL)

```yaml
feature: "<slug>"
integrated: "<date>"
status: "skipped: all local"
biz_count: 0
tech_count: 0
```
