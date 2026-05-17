---
type: api
conventions:
  - testing-api.md
---

# API Test Script Generation Instructions

Type-specific Steps for **API** (HTTP/network endpoint) test script generation. Loaded by the dispatcher when interface detection identifies API-type test cases.

## Classification Indicators

Classify test cases as **API** when they involve any of:

- Endpoints and URL paths
- Request/response structures and data contracts
- HTTP status codes (2xx, 4xx, 5xx)
- HTTP methods (GET, POST, PUT, PATCH, DELETE)
- Authentication headers and authorization
- Request body fields and response body schemas
- Query parameters, path parameters, headers

## Reconnaissance Strategy

API reconnaissance discovers the project's route definitions, handler signatures, request/response schemas, and authentication mechanisms from source code.

### Search Commands

Run these searches to discover API interface details. Adapt file extensions to the project's language.

| Target | Grep Command | What It Finds |
|--------|-------------|---------------|
| Go HTTP handlers | `grep -rn "HandleFunc\|http.Handle\|r.Get(\|r.Post(\|r.Put(\|r.Delete(\|r.Patch(" --include='*.go' .` | Route registration patterns (chi, gorilla/mux, stdlib) |
| Go middleware | `grep -rn "middleware\|Middleware" --include='*.go' .` | Auth middleware bindings, route groups |
| Go request/response | `grep -rn "json.Marshal\|json.Unmarshal\|Bind\|Render\|Respond" --include='*.go' .` | Request/response schema definitions |
| Express routes | `grep -rn "router.get(\|router.post(\|app.get(\|app.post(\|app.put(\|app.delete(" --include='*.ts' --include='*.js' .` | Express/Node route handler registration |
| FastAPI/Flask routes | `grep -rn "@app.get(\|@app.post(\|@router.get(\|@app.route" --include='*.py' .` | Python framework route decorators |
| Spring controllers | `grep -rn "@GetMapping\|@PostMapping\|@RequestMapping\|@RestController" --include='*.java' .` | Spring MVC endpoint annotations |
| Config files | `grep -rn "port\|PORT\|base_url\|BASE_URL\|host\|HOST" --include='*.yaml' --include='*.yml' --include='*.env' --include='*.toml' --include='*.json' .` | API port, base URL, host configuration |
| Auth endpoints | `grep -rn "login\|token\|auth\|session\|jwt\|bearer" --include='*.go' --include='*.ts' --include='*.js' --include='*.py' .` | Authentication endpoint paths, token field names, header formats |
| Response schemas | `grep -rn "Status\|StatusCode\|status_code\|response\|Response" --include='*.go' --include='*.ts' --include='*.js' --include='*.py' .` | Status code usage, response body shaping |

### Reconnaissance Procedure

1. **Detect framework**: Run the grep commands above. Identify which HTTP framework the project uses (chi, gorilla/mux, Express, FastAPI, Flask, Spring, etc.).
2. **Map route tree**: Extract all registered routes with their HTTP methods, path patterns, path parameters, and query parameters. Record each route's handler function name and source location.
3. **Identify middleware**: Discover auth middleware bindings and which route groups they protect. Record the authentication mechanism (Bearer token, API key, session cookie, OAuth).
4. **Locate config**: Find API port, base path prefix, and credential variable names from configuration files.
5. **Extract request/response schemas**: For each handler, identify the expected request body structure, required fields, validation rules, and response body shape including status codes.

### Required Reads

The following source categories must be read during reconnaissance. Adapt discovery to the project's structure.

| Source | What to Extract | Discovery Guidance |
|--------|----------------|---------------------|
| Router files | Route paths, path parameters, middleware bindings | Search for route registration patterns and configuration files. Look for URL path strings, HTTP method bindings, and path parameter definitions. |
| Config files | API port, base path prefix, auth credentials | Search for config/settings files (`.env`, `config.*`, `settings.*`). Look for port numbers, base URLs, and credential variable names. |
| API handlers | Request/response schemas, status codes, validation rules | Search for request handler functions and response definitions. Look for status code usage, input validation, and response body shaping. |
| Auth implementation | Login endpoint path, token field name, header format | Search for authentication/authorization modules. Look for login endpoints, token generation/parsing, and header middleware. |

## Fact Table Required Keys

After reconnaissance, the Fact Table must contain at least these API-specific entries for the completeness gate to pass:

| Key Pattern | Description | Example |
|-------------|-------------|---------|
| `API_PORT` | Port the API server listens on | `API_PORT` = `8080` |
| `AUTH_ENDPOINT` | Login/auth endpoint path with method | `AUTH_ENDPOINT` = `POST /v1/auth/login` |
| `ROUTE_*` | Route path entries for endpoints referenced in test cases | `ROUTE_USERS_LIST` = `GET /v1/users` |

**Minimum requirement**: Either `API_PORT` must be non-UNKNOWN, or at least one `ROUTE_*` entry must be non-UNKNOWN. If all API Fact Table keys are UNKNOWN, skip API test generation and emit a WARNING.

**Completeness gate rule** (from SKILL.md Step 1.5): If all required keys for API are UNKNOWN, do NOT generate API tests. Individual UNKNOWN keys are acceptable -- only skip when every API key is UNKNOWN.

## Verification Method

Before generating API test scripts, confirm the project actually exposes an HTTP API. A project that only has a CLI binary or frontend does not need API test scripts.

Run these checks -- first success is sufficient:

| Check | Command | Pass Condition |
|-------|---------|----------------|
| Go handler patterns | `grep -rn "HandleFunc\|http.Handle\|r.Get(\|r.Post(" --include='*.go' .` | At least one match found |
| Express/Fastify patterns | `grep -rn "router.get(\|app.get(\|app.post(" --include='*.ts' --include='*.js' .` | At least one match found |
| Python framework patterns | `grep -rn "@app.get(\|@router.post(\|@app.route" --include='*.py' .` | At least one match found |
| Spring annotations | `grep -rn "@GetMapping\|@PostMapping\|@RequestMapping" --include='*.java' .` | At least one match found |

**If all checks fail**: The project does not expose an HTTP API. Skip API test generation and emit a WARNING suggesting the user verify source structure.

## Generation Patterns

API test cases translate to executable scripts using HTTP client patterns. Follow the active profile's `generate.md` for framework-specific syntax (HTTP client imports, assertion library, test runner annotations). The type file describes *what* to generate; the profile determines the exact import syntax and assertion format.

### HTTP Request Construction

Each API test function constructs and sends an HTTP request:

1. **Build the request URL**: Combine the base URL (from Fact Table `API_PORT` or config) with the endpoint path from the test case. Use config/environment variables -- never hardcode URLs.
2. **Set HTTP method**: Use the method specified in the test case's Steps field (GET, POST, PUT, PATCH, DELETE).
3. **Set headers**: Include `Content-Type: application/json` for JSON bodies. Add auth headers if the test case requires authentication (from Auth Plan classification).
4. **Set request body**: Construct the body from the test case's Steps field. Use concrete field names and values -- do not invent request payloads.
5. **Send the request**: Use the HTTP client specified in the profile's `generate.md`.

### Response Assertion

API tests must include concrete assertions for each dimension specified in the test case's Expected field:

| Dimension | Assertion Pattern | Example |
|-----------|-------------------|---------|
| Status code | Assert exact HTTP status code | Assert response status is 200, 401, 404, etc. |
| Response body fields | Assert specific field values exist and match | `assert.Equal(body.name, "expected")` |
| Response body schema | Assert response structure matches expected shape | Verify required fields are present and types are correct |
| Response headers | Assert header values when specified | `assert.Contains(headers, "Content-Type", "application/json")` |

### Authentication Integration

Based on the Auth Plan from Step 1, configure authentication in generated API tests:

| Auth Category | Generation Strategy |
|---------------|---------------------|
| **login-test** | Generate independent login/logout logic. No shared auth. Must invalidate cached credentials after to prevent stale state. |
| **auth-required-test** | Use cached shared auth -- credentials acquired once, reused across all tests. Follow profile's `generate.md` for token caching mechanism. |
| **custom-auth-test** (API key, OAuth) | Detect auth mechanism from codebase during reconnaissance. Generate custom auth setup using the discovered mechanism (e.g., `X-API-Key` header, OAuth token exchange). |
| **public-test** | No auth headers needed. |

### Status Code Coverage

Each endpoint mentioned in test cases must have test coverage for:

- **Happy path**: Success response (2xx) with concrete response body assertions
- **Error cases**: Error responses (4xx, 5xx) specified in the test case, with specific error body assertions
- **Auth failures**: If the endpoint requires authentication, verify 401/403 responses when credentials are missing or invalid

### Error Contract Testing

API tests must verify error response contracts explicitly:

- Assert the error response body structure (error message field, error code field)
- Assert specific error messages when the test case specifies them
- Do not write vacuous assertions like "response is not successful" or "returns an error"
- Each error test case must assert the exact status code and at least one response body field

## API Antipattern Guards

Beyond the generic 6 antipattern guards in the main SKILL.md, API-specific generation must avoid these additional patterns:

### 1. Hardcoded URLs

**Pattern**: Embedding literal URLs like `http://localhost:8080/api/users` directly in test code.

**Why harmful**: Tests break when the port or host changes. Cannot run against different environments (staging, CI). Couples tests to a specific deployment configuration.

**Instead**: Use config/environment variables for all endpoint URLs. Read the base URL from the Fact Table (`API_PORT`, config-derived values) or test configuration. Construct full URLs at runtime by combining base URL with endpoint paths.

### 2. Missing Error Contract Tests

**Pattern**: Only testing happy-path (2xx) responses without verifying error response structure (4xx, 5xx).

**Why harmful**: Error responses are part of the API contract. If the error format changes silently (e.g., error message field renamed from `error` to `message`), clients break without any test catching it. API contracts include error responses.

**Instead**: For each endpoint, generate at least one error case test that asserts the error response body structure. If the test case file includes error scenarios, generate tests for them with concrete status code and body assertions.

### 3. Vacuous "Returns Success" Assertions

**Pattern**: Assertions like `assert.Equal(resp.status, 200)` with no response body verification, or `assert.NotNil(resp.body)` without checking any fields.

**Why harmful**: A 200 response with an empty or malformed body passes the test. The assertion verifies almost nothing about the actual API behavior -- any successful response passes regardless of content.

**Instead**: Every API test must assert at least one concrete response body field in addition to the status code. For list endpoints, assert the response is an array and contains at least the expected resource structure. For detail endpoints, assert specific field values from the test case's Expected field.

## Output

API test scripts are written to `tests/e2e/features/<feature>/` following the profile's template naming convention. Each test function includes a traceability comment linking back to the source test case ID.
