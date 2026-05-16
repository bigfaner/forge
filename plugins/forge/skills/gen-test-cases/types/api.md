---
type: api
conventions:
  - testing-api.md
---

# API Test Case Generation Instructions

Type-specific Steps 3-4 for **API** (HTTP/network endpoint) test cases. Loaded by the dispatcher after Step 2.5 interface detection.

## Classification Indicators

Classify a PRD criterion as **API** when it involves any of:

- Endpoints and URL paths
- Request/response structures and data contracts
- HTTP status codes (2xx, 4xx, 5xx)
- HTTP methods (GET, POST, PUT, PATCH, DELETE)
- Authentication headers and authorization
- Request body fields and response body schemas
- Query parameters, path parameters, headers

## Target Derivation

- **Target format**: `api/<resource>`
- Derive `<resource>` from the endpoint path or resource name (e.g., `api/auth`, `api/users`, `api/orders`)

## Test ID Format

- **Test ID**: `<target>/<title-slug>`
- `title-slug` = lowercase title, spaces to hyphens, remove punctuation
- Example: `api/auth/post-login-returns-200-with-token`

## Priority Assignment

1. Criterion tied to a core/critical Given/When/Then in the PRD → **P0**
2. Criterion tied to a secondary story, or an explicit error/boundary case for a core story → **P1**
3. Nice-to-have verifications, minor edge cases → **P2**

If the PRD has no explicit priority marking, default P0 for the first story's ACs and P1 for all others.

## TC Format

```markdown
## TC-{NNN}: {Title}
- **Source**: {Story N / AC-N} or {Spec Section X.Y}
- **Type**: API
- **Target**: api/<resource>
- **Test ID**: api/<resource>/<title-slug>
- **Pre-conditions**: {What must be true before testing}
- **Steps**:
  1. {HTTP method and endpoint, e.g., POST /api/auth/login with body {"email": "...", "password": "..."}}
  2. {Additional request configuration if needed}
- **Expected**: {HTTP status code, response body structure, specific field values}
- **Priority**: P0 | P1 | P2
```

- API test cases do NOT include a `Route` field. They use `Target` and describe the HTTP method + endpoint in Steps.
- Steps must specify the exact HTTP method, endpoint path, and relevant request body/headers.
- Expected results must include concrete HTTP status codes and response body assertions.

## Contract Accuracy

API test cases must demonstrate contract accuracy:

- **Request contracts**: Specify required fields, data types, and valid values in Steps.
- **Response contracts**: Include status code, response body schema, and specific field assertions in Expected.
- **Error contracts**: Cover error responses (4xx, 5xx) with specific error body assertions.

Each endpoint mentioned in the PRD should have test cases covering at minimum:
- Happy path (success response)
- Error cases explicitly mentioned in the PRD

## Route Validation

Cross-reference each API test case's endpoint path against actual route definitions.

**Discovery patterns** (framework-specific):
- Go (chi/stdlib): `r.Get(`, `r.Post(`, `mux.HandleFunc`, `http.Handle`
- Express/Node: `router.get(`, `app.get(`, `app.post(`
- FastAPI/Flask: `@app.get(`, `@router.post(`, `@app.route`
- Spring: `@GetMapping`, `@PostMapping`, `@RequestMapping`

**Validation**: For each test case's endpoint:
- Match against discovered route definitions → annotate `Matched (source:line)`
- No match → annotate `Route not found -- verify path`

If no route files can be discovered, skip this step entirely. Do not fabricate validation results.

## Quality Rules

Apply the 6 Antipattern Prevention rules from the dispatcher's shared rules to every API test case. Key API-specific reminders:

- **Pre-conditions must be concrete and creatable**: Specify how to set up required data state (e.g., "POST /api/users with test user payload" not "user exists").
- **Expected results must be specific and verifiable**: State exact status code and response fields. Not "returns success" or "response is correct".
- **Steps describe runtime behavior**: Make actual HTTP requests to the running API server, not read source code or inspect handler definitions.

## Output

Write to `docs/features/<slug>/testing/api-test-cases.md`. Number test cases from TC-001 sequential. End the file with a traceability table:

```markdown
## Traceability

| TC ID | Source | Type | Target | Priority |
|-------|--------|------|--------|----------|
| TC-001 | Story 1 / AC-1 | API | api/auth | P0 |
```
