---
feature: "{{FEATURE_SLUG}}"
sources:
  - docs/features/{{FEATURE_SLUG}}/prd/prd-user-stories.md
  - docs/features/{{FEATURE_SLUG}}/prd/prd-spec.md
type: cli
generated: "{{DATE}}"
---

# CLI Test Cases: {{FEATURE_SLUG}}

## Summary

| Type | Count |
|------|-------|
| CLI  | {{CLI_COUNT}}  |

---

## CLI Test Cases

| TC ID | Source | Type | Target | Test ID | Pre-conditions | Route | Steps | Expected | Priority |
|-------|--------|------|--------|---------|----------------|-------|--------|----------|----------|
{{CLI_TEST_CASES}}

---

## Traceability

| TC ID | Source | Type | Target | Priority |
|-------|--------|------|--------|----------|
{{TRACEABILITY_ROWS}}
