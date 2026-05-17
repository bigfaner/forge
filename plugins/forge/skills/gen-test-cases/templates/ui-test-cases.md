---
feature: "{{FEATURE_SLUG}}"
sources:
  - docs/features/{{FEATURE_SLUG}}/prd/prd-user-stories.md
  - docs/features/{{FEATURE_SLUG}}/prd/prd-spec.md
  - docs/features/{{FEATURE_SLUG}}/prd/prd-ui-functions.md
type: ui
generated: "{{DATE}}"
---

# UI Test Cases: {{FEATURE_SLUG}}

## Summary

| Type | Count |
|------|-------|
| UI   | {{UI_COUNT}}   |

---

## UI Test Cases

| TC ID | Source | Type | Target | Test ID | Pre-conditions | Route | Steps | Expected | Priority |
|-------|--------|------|--------|---------|----------------|-------|--------|----------|----------|
{{UI_TEST_CASES}}

---

## Traceability

| TC ID | Source | Type | Target | Priority |
|-------|--------|------|--------|----------|
{{TRACEABILITY_ROWS}}

---

## Route Validation

| Route | Status | TC IDs | Matched Route |
|-------|--------|--------|---------------|
{{ROUTE_VALIDATION_ROWS}}

_Omit this section if route files cannot be discovered in the project._
