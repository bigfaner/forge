# JUnit 5 Generate Strategy

Profile-specific test generation rules for the `gen-test-scripts` skill.

## Test Runner & Imports

| Test type | Runner | Assertion | HTTP | Process |
|-----------|--------|-----------|------|---------|
| CLI | JUnit 5 (`@Test`) | `Assertions.*` | — | `ProcessBuilder` |
| API | JUnit 5 (`@Test`) | `Assertions.*` | `java.net.http.HttpClient` (Java 11+) | — |
| TUI | JUnit 5 (`@Test`) | `Assertions.*` | — | `ProcessBuilder` + output comparison |

All tests import from `org.junit.jupiter.api.*`. TestNG and JUnit 4 are **forbidden**.

## Test Class Naming

| Pattern | Use case |
|---------|----------|
| `*E2E.java` | End-to-end test classes (default) |
| `*Test.java` | General test classes |

## Annotations

| Annotation | Purpose |
|------------|---------|
| `@Test` | Marks a test method |
| `@DisplayName("TC-NNN: Description -> PRD Source")` | Human-readable name with traceability |
| `@BeforeEach` / `@AfterEach` | Per-test setup/teardown |
| `@BeforeAll` / `@AfterAll` | Per-class setup/teardown (must be `static`) |
| `@Tag("e2e")` | Tag all E2E test classes |

## Assertions

Use `static org.junit.jupiter.api.Assertions.*`:

| Assertion | Use case |
|-----------|----------|
| `assertEquals(expected, actual)` | Value equality |
| `assertTrue(condition)` | Boolean check |
| `assertNotNull(object)` | Null safety |
| `assertThrows(Exception.class, executable)` | Exception verification |
| `assertAll(...)` | Grouped assertions (all run, report all failures) |

## CLI Testing

```java
ProcessBuilder pb = new ProcessBuilder("command", "arg1", "arg2");
pb.redirectErrorStream(true);
Process process = pb.start();
String output = new String(process.getInputStream().readAllBytes());
int exitCode = process.waitFor();
assertEquals(0, exitCode);
assertTrue(output.contains("expected text"));
```

## API Testing

Use `java.net.http.HttpClient` (Java 11+ built-in, no external dependencies):

```java
HttpClient client = HttpClient.newHttpClient();
HttpRequest request = HttpRequest.newBuilder()
    .uri(URI.create(baseUrl + "/api/resource"))
    .header("Content-Type", "application/json")
    .build();
HttpResponse<String> response = client.send(request, HttpResponse.BodyHandlers.ofString());
assertEquals(200, response.statusCode());
```

## TUI Testing

Launch process via `ProcessBuilder`, capture stdout/stderr, assert output contains expected strings:

```java
ProcessBuilder pb = new ProcessBuilder("java", "-jar", "app.jar");
Process process = pb.start();
String output = new String(process.getInputStream().readAllBytes());
assertTrue(output.contains("Welcome"));
```

## Auth

HTTP headers via `HttpClient` with `Authenticator` or manual header injection:

```java
HttpRequest request = HttpRequest.newBuilder()
    .uri(URI.create(url))
    .header("Authorization", "Bearer " + token)
    .build();
```

## Anti-Patterns (Forbidden)

| Forbidden | Replacement |
|-----------|-------------|
| `Thread.sleep(N)` | Awaitility: `await().atMost(5, SECONDS).until(...)` |
| Hardcoded URLs | Read from config/ENV: `System.getenv("BASE_URL")` |
| `System.out.println` for debugging | Proper assertions only |
| JUnit 4 (`@org.junit.Test`) | JUnit 5 (`@org.junit.jupiter.api.Test`) |

## Test Organization

- One test class per functional domain
- Group related tests by feature area
- Each class annotated with `@Tag("e2e")`

## Traceability

Each `@Test` method must include a `@DisplayName` with traceability:

```java
@Test
@DisplayName("TC-001: Login with valid credentials -> PRD Section 2.1")
void testLoginWithValidCredentials() { ... }
```

## Compilation Check

After generating all test files:

```bash
just e2e-compile
```
