# Surface: API (Application Programming Interface)

API surface 适用于提供 HTTP API 的后端服务（REST、GraphQL 等）。测试重点是 status code、response schema、认证/授权、幂等性。

## Detection Signals

| Signal | File Pattern | Dependency Pattern | Exclusion |
|--------|-------------|-------------------|-----------|
| Go HTTP API | `main.go` exists | `http.Handler`, `gin`, or `echo` in imports/dependencies | No frontend framework entry in `package.json` |
| Node.js HTTP API | `package.json` exists | `express`, `fastify`, or `koa` in `dependencies` | No frontend framework (`react`, `vue`, `svelte`) in dependencies |
| Python HTTP API | `pyproject.toml` or `setup.py` exists | `flask`, `fastapi`, `django`, or `starlette` in dependencies | No frontend entry (no browser DOM references) |
| Python API (test-driven) | `pyproject.toml` or `setup.py` exists | `pytest` or `unittest` in dependencies (no explicit HTTP framework, no CLI/TUI framework) | No frontend entry, no `click`/`typer`/`argparse`/`rich`/`textual`/`prompt_toolkit` |
| Java HTTP API | `pom.xml` or `build.gradle` exists | Spring Boot, JAX-RS, or Jersey in dependencies | No frontend entry (no server-side rendered HTML templates as primary output) |
| Java API (test-driven) | `pom.xml` or `build.gradle` exists | `JUnit` or `TestNG` in dependencies (no explicit HTTP framework) | No frontend entry |
| Rust HTTP API | `Cargo.toml` exists | `actix-web`, `axum`, `rocket`, or `warp` in `[dependencies]` | No frontend entry |

**Confidence Levels**:

- **High**: Language entry file + HTTP framework dependency + no frontend framework
- **Medium**: Language entry file + HTTP framework dependency + ambiguous frontend signal (e.g., template engine for error pages)
- **Low**: Only partial signals (e.g., `package.json` without recognized HTTP framework)

**Disambiguation Rules**:

1. If `package.json` contains both `express`/`fastify` and `react`/`vue`, check for a frontend build output. If the server serves a SPA, classify as WebUI (the API serves the frontend). If the server provides a pure API consumed by external clients, classify as API.
2. Server-side rendering (SSR) with template engines (EJS, Pug, Jinja2): if the primary output is HTML pages for browser consumption, classify as WebUI. If the primary output is JSON/XML API responses, classify as API.
3. GraphQL: If the project exposes a GraphQL endpoint, classify as API regardless of the underlying transport.
4. gRPC/protobuf: If the project only exposes gRPC services without HTTP endpoints, classify as API.

## General Testing Principles

1. **HTTP status code verification**: Every test must assert the response status code. Common assertions:
   - `200 OK` for successful GET/PUT
   - `201 Created` for successful POST
   - `204 No Content` for successful DELETE
   - `400 Bad Request` for validation errors
   - `401 Unauthorized` for missing/invalid authentication
   - `403 Forbidden` for valid auth but insufficient permissions
   - `404 Not Found` for non-existent resources
2. **Response schema validation**: Assert the structure of response bodies, not just individual field values. Validate required fields, data types, and nested object structures.
3. **Authentication/Authorization**: Test with valid credentials, invalid credentials, expired tokens, and insufficient permissions as separate scenarios.
4. **Idempotency**: For PUT and DELETE operations, verify that repeated identical requests produce the same result. For POST, verify that duplicate submissions are handled correctly (idempotency key, deduplication, or explicit duplicate error).
5. **Stateless test isolation**: Each test should set up its own data and clean up afterward. Avoid depending on test execution order. Use database transactions or test-specific namespaces for isolation.

## Test Strategy Guidance

**Test Level Emphasis**: Balanced 50/50 (Contract 50% / Journey smoke 50%)

API testing benefits from both Contract tests (individual endpoint behavior) and Journey smoke tests (multi-step API workflows like authenticate -> create resource -> update -> delete).

**Execution Model**: HTTP client

- Use the project's Convention-defined HTTP testing approach (e.g., `net/http/httptest` in Go, `supertest` in Node.js, `pytest` + `httpx` in Python)
- For integration tests: start the actual server on a random port
- For unit-level tests: use framework-provided test request/response mechanisms
- Each test sends HTTP requests and asserts on the response

**Environment Readiness Checks**:

| Check | How to Verify |
|-------|--------------|
| Server starts | Application server starts without error |
| Database connected | Database connection is established and migrations are current |
| Test fixtures loaded | Required seed data is available |
| Authentication configured | Test API keys or tokens are available for authentication |

**Why balanced 50/50**: API endpoints are highly testable at the Contract level due to well-defined request/response contracts. However, multi-step workflows (create -> list -> update -> delete) exercise cross-endpoint state management that Contract tests alone may miss. Both levels provide distinct value.

## Required Outcome Reference

**Mandatory derived Outcome** (must be considered for every API Journey involving authenticated endpoints):

- **unauthorized**: Request to an authenticated endpoint without valid credentials. Example: GET `/api/tasks` without Bearer token, or with expired token. Assert: `401 Unauthorized` status code, response body contains authentication error details, no sensitive data leaked in response.

**Additional common API boundary Outcomes**:

- **validation-error**: Request with invalid body/parameters. Assert: `400 Bad Request`, response body lists all validation failures.
- **not-found**: Request for non-existent resource. Assert: `404 Not Found`, response body contains resource identifier for debugging.
- **conflict**: Attempt to create duplicate resource. Assert: `409 Conflict`, response indicates which constraint was violated.
- **rate-limit**: Too many requests in a time window. Assert: `429 Too Many Requests`, response includes `Retry-After` header.
- **internal-error**: Unexpected server error. Assert: `500 Internal Server Error`, no stack trace leaked in production mode.
- **pagination**: Large result set is paginated. Assert: response includes pagination metadata (total, page, per_page), links work correctly.
