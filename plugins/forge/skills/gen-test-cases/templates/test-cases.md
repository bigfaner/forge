---
feature: "{{FEATURE_SLUG}}"
sources:
  - prd/prd-user-stories.md
  - prd/prd-spec.md
  - prd/prd-ui-functions.md
generated: "{{DATE}}"
---

# Test Cases: {{FEATURE_SLUG}}

## Summary

| Type | Count |
|------|-------|
| UI   | {{UI_COUNT}}   |
| API  | {{API_COUNT}}  |
| CLI  | {{CLI_COUNT}}  |
| **Total** | **{{TOTAL_COUNT}}** |

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
