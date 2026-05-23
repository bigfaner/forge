---
title: "Java JUnit 5 Testing Convention"
---

# Java JUnit 5 Testing Convention

Convention for generating Java test code using JUnit 5 (Jupiter) with standard assertion methods.

## framework

- **name**: JUnit 5 (Jupiter)
- **version**: JUnit 5.8+
- **language**: Java
- **runner_command**: `mvn test` / `gradle test`

## discovery

- **test_dir**: `src/test/java/`
- **file_pattern**: `*Test.java`, `*Tests.java`
- **exclude_pattern**: `target/`, `build/`, `node_modules/`

## structure

- **suite_pattern**: Test class — each `*Test.java` file contains a test class with `@Test` methods
- **case_pattern**: `@Test void methodName()` — annotated methods within test classes
- **hook_pattern**: `@BeforeEach` / `@AfterEach` / `@BeforeAll` / `@AfterAll`

### Test Class Naming

Pattern: `<Feature>Test` using PascalCase.

```java
class TaskLifecycleTest {
    @Test
    void claimTaskSuccessfully() {
        // test body
    }
}
```

### Nested Test Classes

Use `@Nested` for grouped, hierarchical tests:

```java
class TaskLifecycleTest {
    @Nested
    class ClaimingTasks {
        @Test
        void shouldClaimAvailableTask() {
            // ...
        }

        @Test
        void shouldFailWhenNoTasksAvailable() {
            // ...
        }
    }
}
```

### Parameterized Tests

Use `@ParameterizedTest` with source annotations:

```java
@ParameterizedTest
@CsvSource({
    "hello, HELLO",
    ", ''",
    "123, 123"
})
void shouldUpperCase(String input, String expected) {
    assertEquals(expected, input.toUpperCase());
}
```

### CLI Testing

Use `ProcessBuilder` to invoke CLI binaries:

```java
@Test
void shouldRunCliCommand() throws Exception {
    ProcessBuilder pb = new ProcessBuilder("forge", "subcommand", "--flag", "value");
    pb.environment().put("CLAUDE_PROJECT_DIR", projectDir.toString());
    Process process = pb.start();
    String output = new String(process.getInputStream().readAllBytes());
    int exitCode = process.waitFor();
    assertEquals(0, exitCode);
    assertTrue(output.contains("expected output"));
}
```

### API Testing

Use `java.net.http.HttpClient` for HTTP integration testing:

```java
@Test
void shouldReturnOkFromApiEndpoint() throws Exception {
    HttpClient client = HttpClient.newHttpClient();
    HttpRequest request = HttpRequest.newBuilder()
        .uri(URI.create("http://localhost:8080/api/resource"))
        .GET()
        .build();
    HttpResponse<String> response = client.send(request, HttpResponse.BodyHandlers.ofString());
    assertEquals(200, response.statusCode());
    assertTrue(response.body().contains("data"));
}
```

### Traceability

Each test method should include a traceability comment:

```java
@Test
void shouldLoginWithValidCredentials() {
    // Traceability: TC-001 -> PRD User Auth section
}
```

## assertions

- **style**: static methods from `org.junit.jupiter.api.Assertions`
- **library**: JUnit 5 built-in Assertions class
- **custom_matchers**: none

### Key Functions

- `assertEquals(expected, actual)` — equality check
- `assertNotEquals(unexpected, actual)` — inequality check
- `assertTrue(condition)` — boolean true
- `assertFalse(condition)` — boolean false
- `assertNull(object)` — null check
- `assertNotNull(object)` — not-null check
- `assertThrows(Exception.class, executable)` — expected exception
- `assertDoesNotThrow(executable)` — no exception expected
- `assertIterableEquals(expected, actual)` — collection equality
- `assertLinesMatch(expected, actual)` — line-by-line match
- `assertTimeout(duration, executable)` — timeout check
- `assertAll(executables...)` — grouped assertions (all run, report all failures)

### Grouped Assertions

Use `assertAll` to check multiple conditions without short-circuiting:

```java
assertAll("response validation",
    () -> assertEquals(200, response.statusCode()),
    () -> assertNotNull(response.body()),
    () -> assertTrue(response.body().contains("data"))
);
```

**Rule**: Use `org.junit.jupiter.api.Assertions` static methods. Do not mix with Hamcrest or AssertJ.

## Tags

- **Format**: `@Tag("name")` annotation
- **Built-in**: none — all tags are custom
- **Filtering**: `mvn test -Dgroups="e2e"` or `gradle test --tests "*Tag*"`

```java
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;

@Tag("e2e")
class TaskLifecycleTest {
    @Test
    @Tag("slow")
    void shouldProcessLargeDataset() {
        // ...
    }
}
```

## Result Format

- **Maven output**: `target/surefire-reports/` (XML)
- **Gradle output**: `build/test-results/test/` (XML)
- **Format type**: Surefire XML report

### Surefire XML Structure

```xml
<testsuite name="com.example.TaskLifecycleTest" tests="5" failures="0" errors="0" time="1.234">
  <testcase classname="com.example.TaskLifecycleTest" name="shouldClaimTask" time="0.123" />
</testsuite>
```

## Import Patterns

Standard imports for JUnit 5 e2e tests:

```java
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.AfterAll;
import org.junit.jupiter.api.Nested;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.io.TempDir;

import static org.junit.jupiter.api.Assertions.*;
```

- Parameterized tests add: `import org.junit.jupiter.params.ParameterizedTest` and source annotations
- HTTP tests add: `import java.net.http.*`
- File tests add: `import java.nio.file.*`

## Anti-patterns (Forbidden)

| Forbidden | Replacement |
|-----------|-------------|
| `Thread.sleep()` for synchronization | Retry loop with `Awaitility` library |
| JUnit 4 `@Test` (from `org.junit.Test`) | JUnit 5 `@Test` (from `org.junit.jupiter.api.Test`) |
| `Assert` class (JUnit 4) | `Assertions` class (JUnit 5) |
| Hamcrest `assertThat` | JUnit 5 `Assertions` static methods |
| Hardcoded ports | Dynamic port allocation with `@TempDir` or `ServerSocket(0)` |
| Real secrets/tokens in code | `System.getenv("E2E_API_TOKEN")` |
| `System.out.println` for debug | Remove or use proper logging |
| `@Disabled` without reason | Implement properly or don't generate |

## Helpers

### runCli helper

```java
import java.io.IOException;

record CliResult(String stdout, String stderr, int exitCode) {}

static CliResult runCli(String... args) throws IOException, InterruptedException {
    ProcessBuilder pb = new ProcessBuilder(args);
    Process process = pb.start();
    String stdout = new String(process.getInputStream().readAllBytes());
    String stderr = new String(process.getErrorStream().readAllBytes());
    int exitCode = process.waitFor();
    return new CliResult(stdout, stderr, exitCode);
}
```

### retry helper

```java
static void retry(Runnable action, int maxAttempts, long intervalMs) {
    Exception lastException = null;
    for (int i = 0; i < maxAttempts; i++) {
        try {
            action.run();
            return;
        } catch (Exception e) {
            lastException = e;
            if (i < maxAttempts - 1) {
                Thread.sleep(intervalMs);
            }
        }
    }
    throw new RuntimeException("Retry exhausted", lastException);
}
```
