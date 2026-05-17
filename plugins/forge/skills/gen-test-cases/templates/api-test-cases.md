---
feature: "{{FEATURE_SLUG}}"
sources:
  - docs/features/{{FEATURE_SLUG}}/prd/prd-user-stories.md
  - docs/features/{{FEATURE_SLUG}}/prd/prd-spec.md
type: api
generated: "{{DATE}}"
---

# API Test Cases: {{FEATURE_SLUG}}

## Summary

| Type | Count |
|------|-------|
| API  | {{API_COUNT}}  |

---

## API Test Cases

| TC ID | Source | Type | Target | Test ID | Pre-conditions | Route | Steps | Expected | Priority |
|-------|--------|------|--------|---------|----------------|-------|--------|----------|----------|
{{API_TEST_CASES}}

---

## Traceability

| TC ID | Source | Type | Target | Priority |
|-------|--------|------|--------|----------|
{{TRACEABILITY_ROWS}}
