---
title: "Testing Convention Index"
---

# Testing Convention Index

This index lists all built-in testing conventions. LLM reads this file first, then loads the relevant convention based on project context.

**Loading mechanism**: Two-level loading — read this index first, then load the specific `.md` file matching the project's test framework.

## Conventions

### Go Testing

- **File**: [go.md](go.md)
- **Framework**: Go `testing` package + testify/assert
- **Language**: Go
- **Runner**: `go test -v -json -tags=e2e`
- **When to use**: Go projects using the standard `testing` package with table-driven test patterns. Look for `*_test.go` files, `go.mod`, and `func TestXxx(t *testing.T)` signatures. Exclude if Ginkgo imports are present.

### Ginkgo

- **File**: [ginkgo.md](ginkgo.md)
- **Framework**: Ginkgo v2 + Gomega
- **Language**: Go
- **Runner**: `ginkgo -v --json-report=report.json -tags=e2e`
- **When to use**: Go projects using the Ginkgo BDD framework. Look for `Describe`, `Context`, `It` patterns, and imports of `github.com/onsi/ginkgo/v2`.

### Vitest

- **File**: [vitest.md](vitest.md)
- **Framework**: Vitest
- **Language**: TypeScript / JavaScript
- **Runner**: `vitest run --reporter=verbose`
- **When to use**: TypeScript/JavaScript projects using Vitest. Look for `vitest.config.ts`, `*.test.ts` / `*.spec.ts` files, and imports from `vitest`. Exclude if Jest or Mocha config is present.

### pytest

- **File**: [pytest.md](pytest.md)
- **Framework**: pytest
- **Language**: Python
- **Runner**: `pytest -v`
- **When to use**: Python projects using pytest. Look for `pytest.ini`, `pyproject.toml` with `[tool.pytest]`, `conftest.py`, and `test_*.py` / `*_test.py` files. Exclude if `unittest` or `nosetests` config is present.

### JUnit 5

- **File**: [junit.md](junit.md)
- **Framework**: JUnit 5 (Jupiter)
- **Language**: Java
- **Runner**: `mvn test` / `gradle test`
- **When to use**: Java projects using JUnit 5. Look for `@Test` from `org.junit.jupiter.api`, `*Test.java` files in `src/test/java/`, and Maven/Gradle build files with JUnit 5 dependencies.

### Rust / cargo test

- **File**: [rust.md](rust.md)
- **Framework**: cargo test (built-in Rust test framework)
- **Language**: Rust
- **Runner**: `cargo test`
- **When to use**: Rust projects using the built-in test framework. Look for `Cargo.toml`, `#[test]` attributes, `tests/` directory for integration tests, and `#[cfg(test)] mod` blocks in `src/`.
