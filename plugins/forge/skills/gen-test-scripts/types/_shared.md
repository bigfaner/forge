---
type: _shared
conventions: []
---

# Cross-Type Universal Golden Rules

Framework-agnostic principles that apply to **all** interface types (CLI, TUI, UI, Mobile, API). Type files reference these principles in their Golden Rules section rather than duplicating them.

**Layer model**: `_shared.md` (abstract principles) → type file Golden Rules (type-specific constraints) → Convention (framework implementation).

## Principle: Isolation

Each test must execute in an isolated environment that does not depend on or interfere with other tests, the host system, or external state.

**Constraint**: No test may read or mutate shared global state. Each test creates its own working directory, environment variables, and resource scope. Tests must be independently runnable in any order without prerequisites from other tests.

**Rationale**: Shared state causes order-dependent failures that are impossible to reproduce in isolation. When Test A writes a file that Test B reads, removing or reordering Test A breaks Test B for reasons unrelated to the behavior under test.

**Antipattern guard**: Hardcoded configuration that couples tests to a specific environment (host paths, fixed ports, device dimensions, literal URLs). All configuration must come from environment variables, Fact Table entries, or test-scoped fixtures — never from literal values in test code.

## Principle: Determinism

Given the same test input and system under test, a test must produce the same result on every run, regardless of when or where it executes.

**Constraint**: Tests must not depend on non-reproducible values or external services. Determinism has three sub-dimensions:

**(a) No random dependency**: Timestamps, UUIDs, random numbers, and other non-deterministic values must be replaced with fixed values in test inputs and assertions. If the system under test generates non-deterministic output, assert on structure and type rather than exact values (e.g., assert a UUID-format regex, not a specific UUID; assert a timestamp is within a range, not an exact time).

**(b) No external service dependency**: Third-party APIs, email services, SMS gateways, and other external systems must not be called in tests. Use the Convention's mocking or interception mechanism to replace external calls with deterministic responses. Tests must not require network access to pass.

**(c) No order dependency**: Each test must be independently runnable. No test may assume another test has run before it. Test execution order must not affect outcomes. Setup and teardown are self-contained within each test.

**Rationale**: Non-deterministic tests erode trust. When a test fails intermittently, developers learn to ignore failures, masking real regressions. External service dependencies make tests fail when the service is down, not when the code is broken.

**Antipattern guard**: Tests that call real third-party services, assert on exact timestamps or generated IDs, or require a specific test execution sequence.

## Principle: Timeout Protection

All I/O operations in tests must have an explicit upper-bound timeout. No test may wait indefinitely for any condition.

**Constraint**: Every blocking operation — subprocess execution, HTTP requests, wait conditions, element visibility polls, stdin reads — must specify a timeout. The timeout value must be a named constant or configuration variable, never a bare literal.

**Timeout scope**:

- **Operation-level timeout**: Each individual I/O operation (subprocess spawn, HTTP call, element wait) has its own timeout upper bound.
- **Test-level timeout**: Each test function declares a maximum execution time. If the test exceeds this time, the test runner terminates it and reports failure.

**Default and override**: The Convention defines default timeout values appropriate for the framework and runtime. Type-specific constraints may override defaults (e.g., CLI subprocess timeout may differ from API HTTP timeout). Override values come from the Convention, not from test code.

**Rationale**: Tests without timeouts hang indefinitely in CI, consuming runner minutes and blocking the pipeline. The failure mode is a resource exhaustion timeout rather than an assertion failure — providing zero diagnostic value.

**Antipattern guard**: I/O operations without timeout parameters, or timeout values embedded as magic numbers in test code.

## Principle: Idempotency

Running a test multiple times against the same system under test must produce the same result. Repeated test execution must not accumulate side effects that break subsequent tests.

**Constraint**: For stateful interfaces (API, CLI), each test must create its own test data and clean it up, or use ephemeral data that does not persist. For stateful interfaces with persistent side effects (UI, TUI, Mobile), repeated interaction with the system under test must not leave the application in a state that breaks subsequent tests.

**Rationale**: Idempotent tests support retries in CI. When a test fails due to transient infrastructure issues, re-running it must have the same chance of success as the first run. State accumulation from prior runs makes retries meaningless.

**Antipattern guard**: Tests that create persistent resources without cleanup, or tests that depend on a clean initial state but do not enforce it.

## Principle: Resource Cleanup

Tests must not leave behind temporary files, background processes, database records, browser sessions, or any other side effects after execution.

**Constraint**: Every resource acquired during a test — temp directories, spawned subprocesses, database rows, network connections, browser pages, app instances — must be released or removed when the test completes, regardless of whether it passed or failed. Cleanup must be registered before the resource is used, not appended at the end of the test function.

**Rationale**: Leaked resources accumulate across test suites, causing disk exhaustion, port conflicts, zombie processes, and database pollution. These failures manifest hours or days later, making root cause analysis extremely difficult.

**Antipattern guard**: Tests that create files in the project directory instead of temp directories, spawn processes without termination guarantees, or open connections without closing them.

## Shared Antipattern Guards

The following antipatterns are universal across all types. Type files define only type-specific antipatterns in addition to these.

### 1. Sleep-Based Waits

**Pattern**: Using fixed-duration delays (`sleep`, `time.Sleep`, `setTimeout`, `wait`) to wait for asynchronous operations to complete.

**Why harmful**: Sleep duration is either too short (test flakes on slow CI) or too long (wastes time on fast machines). Masks real timing issues. Makes test suites arbitrarily slow.

**Instead**: Use event-driven waits or polling with timeout. Wait for a concrete observable condition (element visible, output produced, response received, exit code returned) within a timeout window. The observable condition must be defined by the test case's Expected field.

### 2. Hardcoded Configuration

**Pattern**: Embedding environment-specific values (URLs, ports, paths, credentials, device names) directly in test code.

**Why harmful**: Tests break when the environment changes. Cannot run against different environments (staging, CI, local). Couples tests to a specific deployment configuration.

**Instead**: All configuration comes from environment variables, Fact Table entries, or test-scoped configuration objects. The Convention defines how to access these values in the framework's syntax.

### 3. Vacuous Assertions

**Pattern**: Assertions that verify almost nothing — checking only status codes without body content, asserting `not null` without field verification, or asserting "response is successful" without any concrete value check.

**Why harmful**: Any response that satisfies the minimal condition passes the test, regardless of whether the actual behavior is correct. A 200 response with an empty or malformed body passes. Zero regression detection value.

**Instead**: Every test must assert at least one concrete, meaningful value from the test case's Expected field — a specific field value, a specific output string, a specific state change. Status codes or exit codes alone are necessary but not sufficient.

### 4. Source-Code-Level Testing

**Pattern**: Reading source code files and asserting on code text (variable names, function definitions, markdown content) rather than executing the system under test and verifying runtime behavior.

**Why harmful**: Tests the implementation structure, not the behavior. A refactoring that changes internal code without changing observable output breaks the test for no valid reason. Zero verification of actual runtime behavior.

**Instead**: Only test runtime behavior: invoke the interface (execute the binary, send the HTTP request, render the UI, launch the app) and assert on the observable output. Never read source files as test input.
