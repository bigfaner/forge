---
feature: "{{FEATURE_SLUG}}"
sources:
  - docs/features/{{FEATURE_SLUG}}/prd/prd-user-stories.md
  - docs/features/{{FEATURE_SLUG}}/prd/prd-spec.md
  - docs/features/{{FEATURE_SLUG}}/prd/prd-ui-functions.md
generated: "{{DATE}}"
---

# Test Cases: {{FEATURE_SLUG}}

## Summary

| Type | Count |
|------|-------|
| UI   | {{UI_COUNT}}   |
| **Integration** | **{{INTEGRATION_COUNT}}** |
| API  | {{API_COUNT}}  |
| CLI  | {{CLI_COUNT}}  |
| **Total** | **{{TOTAL_COUNT}}** |

> **Note**: Integration test count is a subset of UI count. Integration tests verify that components are correctly wired into their parent pages, using the same Playwright framework as UI tests.

---

## UI Test Cases

{{UI_TEST_CASES}}

---

## API Test Cases

{{API_TEST_CASES}}

---

## CLI Test Cases

{{CLI_TEST_CASES}}

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
